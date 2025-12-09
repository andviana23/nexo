package financial

import (
	"context"

	"github.com/google/uuid"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"go.uber.org/zap"
)

// GenerateFluxoDiarioInput define os dados de entrada para gerar fluxo diário
type GenerateFluxoDiarioInput struct {
	TenantID string
	Data     time.Time
}

// GenerateFluxoDiarioUseCase implementa a geração de fluxo de caixa diário
// Este use case é executado por cron job diariamente
type GenerateFluxoDiarioUseCase struct {
	fluxoRepo         port.FluxoCaixaDiarioRepository
	contasPagarRepo   port.ContaPagarRepository
	contasReceberRepo port.ContaReceberRepository
	compensacaoRepo   port.CompensacaoBancariaRepository
	logger            *zap.Logger
}

// NewGenerateFluxoDiarioUseCase cria nova instância do use case
func NewGenerateFluxoDiarioUseCase(
	fluxoRepo port.FluxoCaixaDiarioRepository,
	contasPagarRepo port.ContaPagarRepository,
	contasReceberRepo port.ContaReceberRepository,
	compensacaoRepo port.CompensacaoBancariaRepository,
	logger *zap.Logger,
) *GenerateFluxoDiarioUseCase {
	return &GenerateFluxoDiarioUseCase{
		fluxoRepo:         fluxoRepo,
		contasPagarRepo:   contasPagarRepo,
		contasReceberRepo: contasReceberRepo,
		compensacaoRepo:   compensacaoRepo,
		logger:            logger,
	}
}

// Execute gera ou atualiza o fluxo de caixa de um dia
func (uc *GenerateFluxoDiarioUseCase) Execute(ctx context.Context, input GenerateFluxoDiarioInput) (*entity.FluxoCaixaDiario, error) {
	// Validações de entrada
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.Data.IsZero() {
		input.Data = time.Now()
	}

	// Buscar fluxo existente ou criar novo
	fluxo, err := uc.fluxoRepo.FindByData(ctx, input.TenantID, input.Data)
	if err != nil {
		// Criar novo fluxo se não existir
		fluxo, err = entity.NewFluxoCaixaDiario(uuid.MustParse(input.TenantID), input.Data)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar fluxo de caixa: %w", err)
		}
	}

	// Calcular saldo inicial (saldo final do dia anterior)
	dataAnterior := input.Data.AddDate(0, 0, -1)
	fluxoAnterior, err := uc.fluxoRepo.FindByData(ctx, input.TenantID, dataAnterior)
	if err == nil && fluxoAnterior != nil {
		fluxo.SetSaldoInicial(fluxoAnterior.SaldoFinal)
	}

	// Calcular entradas confirmadas (contas recebidas)
	statusPago := valueobject.StatusContaPago
	entradasConfirmadas, err := uc.contasReceberRepo.SumByPeriod(ctx, input.TenantID, input.Data, input.Data, &statusPago)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular entradas confirmadas: %w", err)
	}
	fluxo.EntradasConfirmadas = entradasConfirmadas

	// Calcular entradas previstas (contas pendentes para o dia)
	statusPendente := valueobject.StatusContaPendente
	entradasPrevistas, err := uc.contasReceberRepo.SumByPeriod(ctx, input.TenantID, input.Data, input.Data, &statusPendente)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular entradas previstas: %w", err)
	}

	// Incluir compensações bancárias previstas para o dia
	if uc.compensacaoRepo != nil {
		compensacoes, err := uc.compensacaoRepo.ListByDateRange(ctx, input.TenantID, input.Data, input.Data)
		if err == nil {
			for _, comp := range compensacoes {
				if comp.Status == valueobject.StatusCompensacaoPrevisto || comp.Status == valueobject.StatusCompensacaoConfirmado {
					// Adicionar valor líquido às entradas previstas (aguardando compensação)
					entradasPrevistas = entradasPrevistas.Add(comp.ValorLiquido)
				} else if comp.Status == valueobject.StatusCompensacaoCompensado {
					// Adicionar valor líquido às entradas confirmadas
					entradasConfirmadas = entradasConfirmadas.Add(comp.ValorLiquido)
				}
			}
		}
	}
	fluxo.EntradasPrevistas = entradasPrevistas
	fluxo.EntradasConfirmadas = entradasConfirmadas

	// Calcular saídas pagas (contas pagas)
	saidasPagas, err := uc.contasPagarRepo.SumByPeriod(ctx, input.TenantID, input.Data, input.Data, &statusPago)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular saídas pagas: %w", err)
	}
	fluxo.SaidasPagas = saidasPagas

	// Calcular saídas previstas (contas pendentes para o dia)
	saidasPrevistas, err := uc.contasPagarRepo.SumByPeriod(ctx, input.TenantID, input.Data, input.Data, &statusPendente)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular saídas previstas: %w", err)
	}
	fluxo.SaidasPrevistas = saidasPrevistas

	// Calcular saldo final
	fluxo.Calcular()

	// Persistir ou atualizar
	if fluxo.ProcessadoEm.IsZero() {
		if err := uc.fluxoRepo.Create(ctx, fluxo); err != nil {
			return nil, fmt.Errorf("erro ao salvar fluxo de caixa: %w", err)
		}
	} else {
		if err := uc.fluxoRepo.Update(ctx, fluxo); err != nil {
			return nil, fmt.Errorf("erro ao atualizar fluxo de caixa: %w", err)
		}
	}

	uc.logger.Info("Fluxo de caixa diário gerado",
		zap.String("tenant_id", input.TenantID),
		zap.String("data", input.Data.Format("2006-01-02")),
		zap.String("saldo_final", fluxo.SaldoFinal.String()),
	)

	return fluxo, nil
}
