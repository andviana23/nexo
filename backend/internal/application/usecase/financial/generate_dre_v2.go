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

// GenerateDREV2Input define os dados de entrada para gerar DRE V2
type GenerateDREV2Input struct {
	TenantID string
	MesAno   valueobject.MesAno
	// Regime define o regime de reconhecimento de receitas
	// "COMPETENCIA" = reconhece quando CONFIRMADO (padrão contábil)
	// "CAIXA" = reconhece quando RECEBIDO (regime de caixa)
	Regime string
}

// GenerateDREV2UseCase implementa a geração de DRE mensal com suporte a regime de competência vs caixa
// Alinhado com PLANO_AJUSTE_ASAAS.md - Sprint 3
type GenerateDREV2UseCase struct {
	dreRepo           port.DREMensalRepository
	contasPagarRepo   port.ContaPagarRepository
	contasReceberRepo port.ContaReceberRepository
	subscriptionRepo  port.SubscriptionPaymentRepository
	logger            *zap.Logger
}

// NewGenerateDREV2UseCase cria nova instância do use case
func NewGenerateDREV2UseCase(
	dreRepo port.DREMensalRepository,
	contasPagarRepo port.ContaPagarRepository,
	contasReceberRepo port.ContaReceberRepository,
	subscriptionRepo port.SubscriptionPaymentRepository,
	logger *zap.Logger,
) *GenerateDREV2UseCase {
	return &GenerateDREV2UseCase{
		dreRepo:           dreRepo,
		contasPagarRepo:   contasPagarRepo,
		contasReceberRepo: contasReceberRepo,
		subscriptionRepo:  subscriptionRepo,
		logger:            logger,
	}
}

// Execute gera ou atualiza o DRE de um mês usando regime de competência ou caixa
func (uc *GenerateDREV2UseCase) Execute(ctx context.Context, input GenerateDREV2Input) (*entity.DREMensal, error) {
	// Validações de entrada
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.MesAno.String() == "" {
		// Usar mês anterior se não informado
		input.MesAno = valueobject.NewMesAnoFromTime(time.Now().AddDate(0, -1, 0))
	}

	// Default: regime de competência (padrão contábil)
	if input.Regime == "" {
		input.Regime = "COMPETENCIA"
	}

	// Buscar DRE existente ou criar novo
	dre, err := uc.dreRepo.FindByMesAno(ctx, input.TenantID, input.MesAno)
	if err != nil {
		// Criar novo DRE se não existir
		dre, err = entity.NewDREMensal(uuid.MustParse(input.TenantID), input.MesAno)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar DRE: %w", err)
		}
	}

	// Calcular período do mês
	inicio := input.MesAno.PrimeiroDia()
	fim := input.MesAno.UltimoDia()
	competenciaMes := input.MesAno.String() // YYYY-MM

	// ===== RECEITAS =====
	var totalReceitas valueobject.Money
	var receitaAssinaturas valueobject.Money

	if input.Regime == "COMPETENCIA" {
		// Regime de competência: usa confirmed_at e competencia_mes
		// Considera receita quando CONFIRMADA, não quando recebida
		uc.logger.Debug("Calculando receitas por competência",
			zap.String("competencia", competenciaMes),
		)

		// Receitas de assinaturas por competência (status CONFIRMADO ou RECEBIDO)
		receitaAssinaturas, err = uc.contasReceberRepo.SumByCompetencia(ctx, input.TenantID, competenciaMes, nil)
		if err != nil {
			uc.logger.Warn("erro ao calcular receitas por competência, usando fallback",
				zap.Error(err),
			)
			// Fallback: usar método antigo
			statusRecebido := valueobject.StatusContaRecebido
			receitaAssinaturas, _ = uc.contasReceberRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusRecebido)
		}

		totalReceitas = receitaAssinaturas
	} else {
		// Regime de caixa: usa received_at (quando dinheiro entrou)
		uc.logger.Debug("Calculando receitas por regime de caixa",
			zap.Time("inicio", inicio),
			zap.Time("fim", fim),
		)

		statusRecebido := valueobject.StatusContaRecebido
		totalReceitas, err = uc.contasReceberRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusRecebido)
		if err != nil {
			return nil, fmt.Errorf("erro ao calcular receitas (caixa): %w", err)
		}
		receitaAssinaturas = totalReceitas
	}

	// Separar receitas por tipo (quando disponível)
	// TODO: Implementar filtro por origem quando disponível
	// Por enquanto, atribuindo tudo a "planos" (assinaturas são a receita principal)
	dre.SetReceitas(valueobject.Zero(), valueobject.Zero(), receitaAssinaturas)

	// ===== CUSTOS VARIÁVEIS =====
	// TODO: Buscar comissões e consumo de insumos do período
	dre.SetCustosVariaveis(valueobject.Zero(), valueobject.Zero())

	// ===== DESPESAS =====
	// Despesas sempre por regime de caixa (quando efetivamente pagas)
	statusPago := valueobject.StatusContaPago
	despesasTotal, err := uc.contasPagarRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusPago)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular despesas: %w", err)
	}

	// TODO: Separar despesas fixas e variáveis quando tipo disponível
	dre.SetDespesas(despesasTotal, valueobject.Zero())

	// Calcular resultado final
	dre.Calcular()

	// Marcar regime usado
	// TODO: Adicionar campo regime ao DREMensal se necessário

	// Persistir ou atualizar
	if dre.ProcessadoEm.IsZero() {
		if err := uc.dreRepo.Create(ctx, dre); err != nil {
			return nil, fmt.Errorf("erro ao salvar DRE: %w", err)
		}
	} else {
		if err := uc.dreRepo.Update(ctx, dre); err != nil {
			return nil, fmt.Errorf("erro ao atualizar DRE: %w", err)
		}
	}

	uc.logger.Info("DRE mensal V2 gerado",
		zap.String("tenant_id", input.TenantID),
		zap.String("mes_ano", input.MesAno.String()),
		zap.String("regime", input.Regime),
		zap.String("receita_total", dre.ReceitaTotal.String()),
		zap.String("lucro_liquido", dre.LucroLiquido.String()),
	)

	return dre, nil
}

// DefaultMesAnterior retorna o período YYYY-MM do mês anterior ao atual.
func (uc *GenerateDREV2UseCase) DefaultMesAnterior() valueobject.MesAno {
	return valueobject.NewMesAnoFromTime(time.Now().AddDate(0, -1, 0))
}
