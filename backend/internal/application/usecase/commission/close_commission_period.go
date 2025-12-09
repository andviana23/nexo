package commission

import (
	"context"

	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// CloseCommissionPeriodInput representa a entrada para fechar um período de comissão
type CloseCommissionPeriodInput struct {
	TenantID string
	PeriodID string
	ClosedBy string
}

// CloseCommissionPeriodOutput representa a saída do fechamento
type CloseCommissionPeriodOutput struct {
	CommissionPeriod    *entity.CommissionPeriod
	ContaPagar          *entity.ContaPagar
	AdvancesDeducted    int    // COM-004: Quantidade de adiantamentos deduzidos
	TotalAdvancesAmount string // COM-004: Total deduzido em adiantamentos
}

// CloseCommissionPeriodUseCase fecha um período de comissão
// T-COM-002: Gera ContaPagar automaticamente ao fechar período
// COM-004: Deduz adiantamentos aprovados automaticamente
type CloseCommissionPeriodUseCase struct {
	commissionPeriodRepo repository.CommissionPeriodRepository
	commissionItemRepo   repository.CommissionItemRepository
	advanceRepo          repository.AdvanceRepository // COM-004: Repositório de adiantamentos
	contaPagarRepo       port.ContaPagarRepository
	professionalReader   port.ProfessionalReader
	logger               *zap.Logger
}

// NewCloseCommissionPeriodUseCase cria uma nova instância do use case
func NewCloseCommissionPeriodUseCase(
	commissionPeriodRepo repository.CommissionPeriodRepository,
	commissionItemRepo repository.CommissionItemRepository,
	advanceRepo repository.AdvanceRepository, // COM-004: Novo parâmetro
	contaPagarRepo port.ContaPagarRepository,
	professionalReader port.ProfessionalReader,
	logger *zap.Logger,
) *CloseCommissionPeriodUseCase {
	return &CloseCommissionPeriodUseCase{
		commissionPeriodRepo: commissionPeriodRepo,
		commissionItemRepo:   commissionItemRepo,
		advanceRepo:          advanceRepo,
		contaPagarRepo:       contaPagarRepo,
		professionalReader:   professionalReader,
		logger:               logger,
	}
}

// Execute executa o use case
// 1. Valida se período pode ser fechado
// 2. COM-004: Busca e deduz adiantamentos aprovados do profissional
// 3. Soma todas comissões do período
// 4. Cria ContaPagar para o profissional
// 5. Fecha o período vinculando à ContaPagar
func (uc *CloseCommissionPeriodUseCase) Execute(ctx context.Context, input CloseCommissionPeriodInput) (*CloseCommissionPeriodOutput, error) {
	// Verifica se existe
	period, err := uc.commissionPeriodRepo.GetByID(ctx, input.TenantID, input.PeriodID)
	if err != nil {
		return nil, err
	}

	if period == nil {
		return nil, domain.ErrCommissionPeriodNotFound
	}

	// Verifica se pode fechar
	if !period.CanClose() {
		return nil, domain.ErrPeriodoNaoPodeFechado
	}

	output := &CloseCommissionPeriodOutput{}

	// COM-004: Buscar e deduzir adiantamentos aprovados do profissional
	totalAdvancesDeducted := decimal.Zero
	advancesDeductedCount := 0

	if period.ProfessionalID != nil {
		advances, err := uc.advanceRepo.GetApprovedByProfessional(ctx, input.TenantID, *period.ProfessionalID)
		if err != nil {
			uc.logger.Warn("erro ao buscar adiantamentos do profissional",
				zap.String("professional_id", *period.ProfessionalID),
				zap.Error(err))
			// Continua mesmo com erro - não bloqueia fechamento
		} else if len(advances) > 0 {
			uc.logger.Info("adiantamentos encontrados para dedução",
				zap.String("professional_id", *period.ProfessionalID),
				zap.Int("quantidade", len(advances)))

			// Marcar cada adiantamento como deduzido
			for _, advance := range advances {
				_, err := uc.advanceRepo.MarkDeducted(ctx, input.TenantID, advance.ID, input.PeriodID)
				if err != nil {
					uc.logger.Warn("erro ao marcar adiantamento como deduzido",
						zap.String("advance_id", advance.ID),
						zap.Error(err))
					continue
				}
				totalAdvancesDeducted = totalAdvancesDeducted.Add(advance.Amount)
				advancesDeductedCount++
			}

			uc.logger.Info("adiantamentos deduzidos com sucesso",
				zap.Int("quantidade", advancesDeductedCount),
				zap.String("total", totalAdvancesDeducted.String()))
		}
	}

	// Atualizar o TotalAdvances do período com os adiantamentos reais deduzidos
	if advancesDeductedCount > 0 {
		period.TotalAdvances = totalAdvancesDeducted
	}

	output.AdvancesDeducted = advancesDeductedCount
	output.TotalAdvancesAmount = totalAdvancesDeducted.String()

	// Buscar sumário do período para obter totais de comissões
	summary, err := uc.commissionPeriodRepo.GetSummary(ctx, input.TenantID, input.PeriodID)
	if err != nil {
		uc.logger.Warn("erro ao buscar sumário do período", zap.Error(err))
		// Continua mesmo sem sumário
	}

	// Calcular valor líquido (comissões - adiantamentos + ajustes)
	// COM-004: Usar o totalAdvancesDeducted calculado ao invés do valor do período
	totalCommission := period.TotalCommission
	if summary != nil {
		totalCommission = summary.TotalCommission
	}
	totalNet := totalCommission.Sub(totalAdvancesDeducted).Add(period.TotalAdjustments)

	// T-COM-002: Criar ContaPagar se houver valor a pagar e profissional definido
	if !totalNet.IsZero() && totalNet.IsPositive() && period.ProfessionalID != nil {
		// Buscar nome do profissional
		professionalName := "Profissional"
		professional, err := uc.professionalReader.FindByID(ctx, input.TenantID, *period.ProfessionalID)
		if err == nil && professional != nil {
			professionalName = professional.Name
		}

		// Criar ContaPagar
		valorMoney := valueobject.NewMoneyFromDecimal(totalNet)

		// Vencimento: 5 dias úteis após fechamento
		dataVencimento := time.Now().AddDate(0, 0, 7)

		descricao := fmt.Sprintf("Comissão %s - %s", period.ReferenceMonth, professionalName)

		tenantUUID, err := uuid.Parse(input.TenantID)
		if err != nil {
			uc.logger.Error("erro ao converter tenant_id para uuid", zap.Error(err))
			return nil, err
		}

		contaPagar, err := entity.NewContaPagar(
			tenantUUID,
			descricao,
			"COMISSAO",       // CategoriaID padrão para comissões
			professionalName, // Fornecedor = profissional
			valorMoney,
			valueobject.TipoCustoVariavel, // Comissão é custo variável
			dataVencimento,
			false, // Não recorrente
			"",    // Sem periodicidade
		)
		if err != nil {
			uc.logger.Error("erro ao criar conta a pagar para comissão", zap.Error(err))
			// Não bloqueia fechamento do período
		} else {
			contaPagar.Observacoes = fmt.Sprintf("Período de comissão: %s a %s",
				period.PeriodStart.Format("02/01/2006"),
				period.PeriodEnd.Format("02/01/2006"))

			if err := uc.contaPagarRepo.Create(ctx, contaPagar); err != nil {
				uc.logger.Error("erro ao persistir conta a pagar", zap.Error(err))
			} else {
				output.ContaPagar = contaPagar
				// Vincular ContaPagar ao período
				period.ContaPagarID = &contaPagar.ID

				uc.logger.Info("conta a pagar criada para comissão",
					zap.String("conta_pagar_id", contaPagar.ID),
					zap.String("period_id", period.ID),
					zap.String("valor", totalNet.String()))
			}
		}
	}

	// Fecha o período
	closed, err := uc.commissionPeriodRepo.Close(ctx, input.TenantID, input.PeriodID, input.ClosedBy)
	if err != nil {
		return nil, err
	}

	// COM-005: Atualizar status dos commission_items para PROCESSADO
	if period.ProfessionalID != nil {
		itemsProcessed, err := uc.commissionItemRepo.AssignToPeriod(
			ctx,
			input.TenantID,
			*period.ProfessionalID,
			input.PeriodID,
			period.PeriodStart,
			period.PeriodEnd,
		)
		if err != nil {
			uc.logger.Warn("erro ao processar itens de comissão do período",
				zap.String("period_id", input.PeriodID),
				zap.Error(err))
			// Não bloqueia - o período já foi fechado
		} else {
			uc.logger.Info("itens de comissão processados",
				zap.String("period_id", input.PeriodID),
				zap.Int64("itens_processados", itemsProcessed))
		}
	}

	output.CommissionPeriod = closed

	uc.logger.Info("período de comissão fechado",
		zap.String("period_id", input.PeriodID),
		zap.String("total_comissao", totalNet.String()),
		zap.Int("adiantamentos_deduzidos", advancesDeductedCount),
		zap.String("total_adiantamentos", totalAdvancesDeducted.String()),
		zap.Bool("conta_pagar_criada", output.ContaPagar != nil))

	return output, nil
}

// MarkPeriodAsPaidInput representa a entrada para marcar um período como pago
type MarkPeriodAsPaidInput struct {
	TenantID string
	PeriodID string
	PaidBy   string
}

// MarkPeriodAsPaidOutput representa a saída da marcação como pago
type MarkPeriodAsPaidOutput struct {
	CommissionPeriod *entity.CommissionPeriod
}

// MarkPeriodAsPaidUseCase marca um período de comissão como pago
type MarkPeriodAsPaidUseCase struct {
	commissionPeriodRepo repository.CommissionPeriodRepository
}

// NewMarkPeriodAsPaidUseCase cria uma nova instância do use case
func NewMarkPeriodAsPaidUseCase(commissionPeriodRepo repository.CommissionPeriodRepository) *MarkPeriodAsPaidUseCase {
	return &MarkPeriodAsPaidUseCase{
		commissionPeriodRepo: commissionPeriodRepo,
	}
}

// Execute executa o use case
func (uc *MarkPeriodAsPaidUseCase) Execute(ctx context.Context, input MarkPeriodAsPaidInput) (*MarkPeriodAsPaidOutput, error) {
	// Verifica se existe
	period, err := uc.commissionPeriodRepo.GetByID(ctx, input.TenantID, input.PeriodID)
	if err != nil {
		return nil, err
	}

	if period == nil {
		return nil, domain.ErrCommissionPeriodNotFound
	}

	// Verifica se pode marcar como pago
	if !period.CanPay() {
		return nil, domain.ErrPeriodoNaoPodePago
	}

	// Marca como pago
	paid, err := uc.commissionPeriodRepo.MarkAsPaid(ctx, input.TenantID, input.PeriodID, input.PaidBy)
	if err != nil {
		return nil, err
	}

	return &MarkPeriodAsPaidOutput{CommissionPeriod: paid}, nil
}

// DeleteCommissionPeriodInput representa a entrada para deletar um período de comissão
type DeleteCommissionPeriodInput struct {
	TenantID string
	PeriodID string
}

// DeleteCommissionPeriodOutput representa a saída da exclusão
type DeleteCommissionPeriodOutput struct {
	Success bool
}

// DeleteCommissionPeriodUseCase deleta um período de comissão
type DeleteCommissionPeriodUseCase struct {
	commissionPeriodRepo repository.CommissionPeriodRepository
}

// NewDeleteCommissionPeriodUseCase cria uma nova instância do use case
func NewDeleteCommissionPeriodUseCase(commissionPeriodRepo repository.CommissionPeriodRepository) *DeleteCommissionPeriodUseCase {
	return &DeleteCommissionPeriodUseCase{
		commissionPeriodRepo: commissionPeriodRepo,
	}
}

// Execute executa o use case
func (uc *DeleteCommissionPeriodUseCase) Execute(ctx context.Context, input DeleteCommissionPeriodInput) (*DeleteCommissionPeriodOutput, error) {
	// Verifica se existe
	period, err := uc.commissionPeriodRepo.GetByID(ctx, input.TenantID, input.PeriodID)
	if err != nil {
		return nil, err
	}

	if period == nil {
		return nil, domain.ErrCommissionPeriodNotFound
	}

	// Só pode deletar se estiver aberto
	if period.Status != "ABERTO" {
		return nil, domain.ErrPeriodoNaoPodeFechado
	}

	// Deleta
	err = uc.commissionPeriodRepo.Delete(ctx, input.TenantID, input.PeriodID)
	if err != nil {
		return nil, err
	}

	return &DeleteCommissionPeriodOutput{Success: true}, nil
}
