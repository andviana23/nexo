package financial

import (
	"context"

	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"go.uber.org/zap"
)

// GenerateFluxoDiarioV2Input define os dados de entrada para gerar fluxo diário V2
type GenerateFluxoDiarioV2Input struct {
	TenantID string
	Data     time.Time
}

// GenerateFluxoDiarioV2UseCase implementa a geração de fluxo de caixa diário
// com suporte a received_at para assinaturas Asaas
// Alinhado com PLANO_AJUSTE_ASAAS.md - Sprint 3
type GenerateFluxoDiarioV2UseCase struct {
	fluxoRepo         port.FluxoCaixaDiarioRepository
	contasPagarRepo   port.ContaPagarRepository
	contasReceberRepo port.ContaReceberRepository
	compensacaoRepo   port.CompensacaoBancariaRepository
	logger            *zap.Logger
}

// NewGenerateFluxoDiarioV2UseCase cria nova instância do use case
func NewGenerateFluxoDiarioV2UseCase(
	fluxoRepo port.FluxoCaixaDiarioRepository,
	contasPagarRepo port.ContaPagarRepository,
	contasReceberRepo port.ContaReceberRepository,
	compensacaoRepo port.CompensacaoBancariaRepository,
	logger *zap.Logger,
) *GenerateFluxoDiarioV2UseCase {
	return &GenerateFluxoDiarioV2UseCase{
		fluxoRepo:         fluxoRepo,
		contasPagarRepo:   contasPagarRepo,
		contasReceberRepo: contasReceberRepo,
		compensacaoRepo:   compensacaoRepo,
		logger:            logger,
	}
}

// Execute gera ou atualiza o fluxo de caixa de um dia usando received_at para assinaturas
func (uc *GenerateFluxoDiarioV2UseCase) Execute(ctx context.Context, input GenerateFluxoDiarioV2Input) (*entity.FluxoCaixaDiario, error) {
	// Validações de entrada
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.Data.IsZero() {
		input.Data = time.Now()
	}

	// Normalizar data para início do dia
	data := time.Date(input.Data.Year(), input.Data.Month(), input.Data.Day(), 0, 0, 0, 0, input.Data.Location())

	// Buscar fluxo existente ou criar novo
	fluxo, err := uc.fluxoRepo.FindByData(ctx, input.TenantID, data)
	if err != nil {
		// Criar novo fluxo se não existir
		fluxo, err = entity.NewFluxoCaixaDiario(uuid.MustParse(input.TenantID), data)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar fluxo de caixa: %w", err)
		}
	}

	// Calcular saldo inicial (saldo final do dia anterior)
	dataAnterior := data.AddDate(0, 0, -1)
	fluxoAnterior, err := uc.fluxoRepo.FindByData(ctx, input.TenantID, dataAnterior)
	if err == nil && fluxoAnterior != nil {
		fluxo.SetSaldoInicial(fluxoAnterior.SaldoFinal)
	}

	// ===== ENTRADAS - REGIME DE CAIXA =====
	// Para fluxo de caixa, sempre usamos received_at (quando o dinheiro entrou)

	// 1. Entradas confirmadas de contas a receber (received_at = data)
	// Usar método do repositório SumByReceivedDate
	proximoDia := data.AddDate(0, 0, 1)

	// Somar por received_at (regime de caixa real)
	entradasAsaas, err := uc.contasReceberRepo.SumByReceivedDate(ctx, input.TenantID, data, proximoDia)
	if err != nil {
		uc.logger.Warn("erro ao buscar entradas por received_at, usando fallback",
			zap.Error(err),
		)
	}

	// 2. Entradas tradicionais (contas recebidas no dia)
	statusRecebido := valueobject.StatusContaRecebido
	entradasTradicionais, err := uc.contasReceberRepo.SumByPeriod(ctx, input.TenantID, data, data, &statusRecebido)
	if err != nil {
		uc.logger.Warn("erro ao calcular entradas tradicionais", zap.Error(err))
	}

	// 3. Entradas previstas (contas pendentes para o dia)
	statusPendente := valueobject.StatusContaPendente
	entradasPrevistas, err := uc.contasReceberRepo.SumByPeriod(ctx, input.TenantID, data, data, &statusPendente)
	if err != nil {
		uc.logger.Warn("erro ao calcular entradas previstas", zap.Error(err))
	}

	// Combinar entradas confirmadas (evitar duplicação)
	entradasConfirmadas := entradasTradicionais
	if entradasAsaas.IsPositive() {
		entradasConfirmadas = entradasAsaas
	}

	// 3.b Incluir compensações bancárias previstas/confirmadas/compensadas no dia
	if uc.compensacaoRepo != nil {
		comps, err := uc.compensacaoRepo.ListByDateRange(ctx, input.TenantID, data, data)
		if err != nil {
			uc.logger.Warn("erro ao listar compensações para fluxo diário V2", zap.Error(err))
		} else {
			for _, comp := range comps {
				switch comp.Status {
				case valueobject.StatusCompensacaoPrevisto, valueobject.StatusCompensacaoConfirmado:
					entradasPrevistas = entradasPrevistas.Add(comp.ValorLiquido)
				case valueobject.StatusCompensacaoCompensado:
					entradasConfirmadas = entradasConfirmadas.Add(comp.ValorLiquido)
				}
			}
		}
	}

	fluxo.EntradasPrevistas = entradasPrevistas
	fluxo.EntradasConfirmadas = entradasConfirmadas

	// ===== SAÍDAS =====
	// 4. Saídas pagas (contas pagas no dia)
	statusPago := valueobject.StatusContaPago
	saidasPagas, err := uc.contasPagarRepo.SumByPeriod(ctx, input.TenantID, data, data, &statusPago)
	if err != nil {
		uc.logger.Warn("erro ao calcular saídas pagas", zap.Error(err))
	}
	fluxo.SaidasPagas = saidasPagas

	// 5. Saídas previstas (contas pendentes para o dia)
	saidasPrevistas, err := uc.contasPagarRepo.SumByPeriod(ctx, input.TenantID, data, data, &statusPendente)
	if err != nil {
		uc.logger.Warn("erro ao calcular saídas previstas", zap.Error(err))
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

	uc.logger.Info("Fluxo de caixa diário V2 gerado",
		zap.String("tenant_id", input.TenantID),
		zap.String("data", data.Format("2006-01-02")),
		zap.String("entradas_confirmadas", fluxo.EntradasConfirmadas.String()),
		zap.String("entradas_previstas", fluxo.EntradasPrevistas.String()),
		zap.String("saidas_pagas", fluxo.SaidasPagas.String()),
		zap.String("saldo_final", fluxo.SaldoFinal.String()),
	)

	return fluxo, nil
}
