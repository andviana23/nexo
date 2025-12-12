package financial

import (
	"context"

	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"go.uber.org/zap"
)

// MarcarCompensacaoInput define os dados de entrada
type MarcarCompensacaoInput struct {
	TenantID        string
	CompensacaoID   string
	DataConfirmacao time.Time
}

// MarcarCompensacaoUseCase marca compensações como confirmadas/compensadas
// Este use case é executado manualmente ou por cron job diário
type MarcarCompensacaoUseCase struct {
	repo             port.CompensacaoBancariaRepository
	contaReceberRepo port.ContaReceberRepository
	logger           *zap.Logger
}

// NewMarcarCompensacaoUseCase cria nova instância do use case
func NewMarcarCompensacaoUseCase(
	repo port.CompensacaoBancariaRepository,
	contaReceberRepo port.ContaReceberRepository,
	logger *zap.Logger,
) *MarcarCompensacaoUseCase {
	return &MarcarCompensacaoUseCase{
		repo:             repo,
		contaReceberRepo: contaReceberRepo,
		logger:           logger,
	}
}

// Execute marca uma compensação como compensada
func (uc *MarcarCompensacaoUseCase) Execute(ctx context.Context, input MarcarCompensacaoInput) (*entity.CompensacaoBancaria, error) {
	// Validações de entrada
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.CompensacaoID == "" {
		return nil, fmt.Errorf("ID da compensação é obrigatório")
	}

	// Buscar compensação existente
	comp, err := uc.repo.FindByID(ctx, input.TenantID, input.CompensacaoID)
	if err != nil {
		uc.logger.Error("Erro ao buscar compensação",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("comp_id", input.CompensacaoID),
		)
		return nil, fmt.Errorf("erro ao buscar compensação: %w", err)
	}

	// Marcar como compensado (método do domínio)
	if err := comp.MarcarComoCompensado(); err != nil {
		uc.logger.Error("Erro ao marcar compensação como compensada",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("comp_id", input.CompensacaoID),
		)
		return nil, fmt.Errorf("erro ao marcar compensação: %w", err)
	}

	// Atualizar no repositório
	if err := uc.repo.Update(ctx, comp); err != nil {
		uc.logger.Error("Erro ao atualizar compensação",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("comp_id", input.CompensacaoID),
		)
		return nil, fmt.Errorf("erro ao atualizar compensação: %w", err)
	}

	uc.logger.Info("Compensação marcada como compensada",
		zap.String("tenant_id", input.TenantID),
		zap.String("comp_id", comp.ID),
		zap.String("data_compensado", comp.DataCompensado.Format("2006-01-02")),
	)

	// Atualizar conta a receber vinculada (quando receita_id aponta para contas_a_receber)
	uc.atualizarContaReceberCompensada(ctx, input.TenantID, comp)

	return comp, nil
}

// ExecuteBatch processa automaticamente compensações pendentes que já venceram
func (uc *MarcarCompensacaoUseCase) ExecuteBatch(ctx context.Context, tenantID string) (int, error) {
	if tenantID == "" {
		return 0, domain.ErrTenantIDRequired
	}

	// Buscar compensações pendentes
	compensacoes, err := uc.repo.ListPendentesCompensacao(ctx, tenantID)
	if err != nil {
		uc.logger.Error("Erro ao buscar compensações pendentes",
			zap.Error(err),
			zap.String("tenant_id", tenantID),
		)
		return 0, fmt.Errorf("erro ao buscar compensações pendentes: %w", err)
	}

	count := 0

	for _, comp := range compensacoes {
		if err := comp.MarcarComoCompensado(); err != nil {
			uc.logger.Warn("Erro ao marcar compensação individual",
				zap.Error(err),
				zap.String("comp_id", comp.ID),
			)
			continue
		}

		if err := uc.repo.Update(ctx, comp); err != nil {
			uc.logger.Warn("Erro ao atualizar compensação individual",
				zap.Error(err),
				zap.String("comp_id", comp.ID),
			)
			continue
		}

		uc.atualizarContaReceberCompensada(ctx, tenantID, comp)
		count++
	}

	uc.logger.Info("Compensações processadas em lote",
		zap.String("tenant_id", tenantID),
		zap.Int("total", len(compensacoes)),
		zap.Int("marcadas", count),
	)

	return count, nil
}

// atualizarContaReceberCompensada marca a conta vinculada como recebida quando a compensação liquida.
func (uc *MarcarCompensacaoUseCase) atualizarContaReceberCompensada(ctx context.Context, tenantID string, comp *entity.CompensacaoBancaria) {
	if uc.contaReceberRepo == nil || comp == nil || comp.ReceitaID == "" {
		return
	}

	conta, err := uc.contaReceberRepo.FindByID(ctx, tenantID, comp.ReceitaID)
	if err != nil || conta == nil {
		uc.logger.Debug("conta a receber vinculada não encontrada para compensação",
			zap.String("tenant_id", tenantID),
			zap.String("receita_id", comp.ReceitaID),
			zap.Error(err),
		)
		return
	}

	if conta.Status == valueobject.StatusContaRecebido ||
		conta.Status == valueobject.StatusContaEstornado ||
		conta.Status == valueobject.StatusContaCancelado {
		return
	}

	dataRecebimento := time.Now()
	if comp.DataCompensado != nil {
		dataRecebimento = *comp.DataCompensado
	}

	if err := conta.MarcarComoRecebido(dataRecebimento); err != nil {
		uc.logger.Warn("erro ao marcar conta a receber como recebida via compensação",
			zap.String("tenant_id", tenantID),
			zap.String("conta_receber_id", conta.ID),
			zap.Error(err),
		)
		return
	}

	conta.ReceivedAt = &dataRecebimento
	if err := uc.contaReceberRepo.Update(ctx, conta); err != nil {
		uc.logger.Warn("erro ao atualizar conta a receber após compensação",
			zap.String("tenant_id", tenantID),
			zap.String("conta_receber_id", conta.ID),
			zap.Error(err),
		)
	}
}
