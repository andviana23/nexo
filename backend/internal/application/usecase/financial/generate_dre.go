package financial

import (
	"context"

	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// GenerateDREInput define os dados de entrada para gerar DRE
type GenerateDREInput struct {
	TenantID string
	MesAno   valueobject.MesAno
}

// GenerateDREUseCase implementa a geração de DRE mensal
// Este use case é executado por cron job mensalmente
type GenerateDREUseCase struct {
	dreRepo            port.DREMensalRepository
	contasPagarRepo    port.ContaPagarRepository
	contasReceberRepo  port.ContaReceberRepository
	commissionItemRepo repository.CommissionItemRepository
	logger             *zap.Logger
}

// NewGenerateDREUseCase cria nova instância do use case
func NewGenerateDREUseCase(
	dreRepo port.DREMensalRepository,
	contasPagarRepo port.ContaPagarRepository,
	contasReceberRepo port.ContaReceberRepository,
	commissionItemRepo repository.CommissionItemRepository,
	logger *zap.Logger,
) *GenerateDREUseCase {
	return &GenerateDREUseCase{
		dreRepo:            dreRepo,
		contasPagarRepo:    contasPagarRepo,
		contasReceberRepo:  contasReceberRepo,
		commissionItemRepo: commissionItemRepo,
		logger:             logger,
	}
}

// Execute gera ou atualiza o DRE de um mês
func (uc *GenerateDREUseCase) Execute(ctx context.Context, input GenerateDREInput) (*entity.DREMensal, error) {
	// Validações de entrada
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.MesAno.String() == "" {
		// Usar mês anterior se não informado
		input.MesAno = valueobject.NewMesAnoFromTime(time.Now().AddDate(0, -1, 0))
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
	statusPago := valueobject.StatusContaPago

	// ===== RECEITAS POR ORIGEM =====
	// Receitas de serviços
	receitaServicos, err := uc.contasReceberRepo.SumByOrigem(ctx, input.TenantID, "SERVICO", inicio, fim)
	if err != nil {
		uc.logger.Warn("erro ao calcular receitas de serviços", zap.Error(err))
		receitaServicos = valueobject.Zero()
	}

	// Receitas de produtos
	receitaProdutos, err := uc.contasReceberRepo.SumByOrigem(ctx, input.TenantID, "PRODUTO", inicio, fim)
	if err != nil {
		uc.logger.Warn("erro ao calcular receitas de produtos", zap.Error(err))
		receitaProdutos = valueobject.Zero()
	}

	// Receitas de assinaturas
	receitaAssinaturas, err := uc.contasReceberRepo.SumByOrigem(ctx, input.TenantID, "ASSINATURA", inicio, fim)
	if err != nil {
		uc.logger.Warn("erro ao calcular receitas de assinaturas", zap.Error(err))
		receitaAssinaturas = valueobject.Zero()
	}

	// Se todas as receitas por origem forem zero, usar receita total como fallback
	if receitaServicos.IsZero() && receitaProdutos.IsZero() && receitaAssinaturas.IsZero() {
		totalReceitas, err := uc.contasReceberRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusPago)
		if err != nil {
			return nil, fmt.Errorf("erro ao calcular receitas: %w", err)
		}
		receitaServicos = totalReceitas
	}

	dre.SetReceitas(receitaServicos, receitaProdutos, receitaAssinaturas)

	// ===== CUSTOS VARIÁVEIS =====
	// Buscar total de comissões do período usando o módulo de comissões
	var custoComissoes valueobject.Money
	if uc.commissionItemRepo != nil {
		totalComissoes, err := uc.commissionItemRepo.SumByDateRange(ctx, input.TenantID, inicio, fim)
		if err != nil {
			uc.logger.Warn("erro ao calcular comissões do período", zap.Error(err))
			custoComissoes = valueobject.Zero()
		} else {
			custoComissoes = valueobject.NewMoneyFromDecimal(decimal.NewFromFloat(totalComissoes))
		}
	} else {
		uc.logger.Warn("commissionItemRepo não configurado, usando comissões zeradas")
		custoComissoes = valueobject.Zero()
	}

	// Custo de insumos (por enquanto zerado - TODO: integrar com módulo de estoque)
	dre.SetCustosVariaveis(custoComissoes, valueobject.Zero())

	// ===== DESPESAS =====
	// Buscar despesas pagas no período
	despesasTotais, err := uc.contasPagarRepo.SumByPeriod(ctx, input.TenantID, inicio, fim, &statusPago)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular despesas: %w", err)
	}

	// Por enquanto, separando 70% fixas e 30% variáveis como estimativa
	// TODO: Implementar separação real por tipo (FIXA vs VARIAVEL) quando houver query específica
	despesasFixasVal := despesasTotais.Value().Mul(decimalFromFloat(0.7))
	despesasVariaveisVal := despesasTotais.Value().Mul(decimalFromFloat(0.3))

	despesasFixas := valueobject.NewMoneyFromDecimal(despesasFixasVal)
	despesasVariaveis := valueobject.NewMoneyFromDecimal(despesasVariaveisVal)

	dre.SetDespesas(despesasFixas, despesasVariaveis)

	// Calcular resultado final
	dre.Calcular()

	// Persistir usando UPSERT (CreateOrUpdate) para garantir última versão
	// Se dre.ID for vazio, é criação. Se não, é update.
	// No caso de DRE, a chave única é (TenantID, MesAno).

	// Verificar se já existe DRE para este mês/ano novamente (double check lock)
	existingDRE, err := uc.dreRepo.FindByMesAno(ctx, input.TenantID, input.MesAno)
	if err == nil && existingDRE != nil {
		dre.ID = existingDRE.ID // Garantir que usamos o ID existente para o Update
		if err := uc.dreRepo.Update(ctx, dre); err != nil {
			return nil, fmt.Errorf("erro ao atualizar DRE: %w", err)
		}
	} else {
		if err := uc.dreRepo.Create(ctx, dre); err != nil {
			// Se falhar no create por chave duplicada (race condition), tentamos update
			if existingDRE, err := uc.dreRepo.FindByMesAno(ctx, input.TenantID, input.MesAno); err == nil && existingDRE != nil {
				dre.ID = existingDRE.ID
				if err := uc.dreRepo.Update(ctx, dre); err != nil {
					return nil, fmt.Errorf("erro ao atualizar DRE após retry: %w", err)
				}
			} else {
				return nil, fmt.Errorf("erro ao salvar DRE: %w", err)
			}
		}
	}

	uc.logger.Info("DRE mensal gerado",
		zap.String("tenant_id", input.TenantID),
		zap.String("mes_ano", input.MesAno.String()),
		zap.String("receita_total", dre.ReceitaTotal.String()),
		zap.String("lucro_liquido", dre.LucroLiquido.String()),
	)

	return dre, nil
}

// decimalFromFloat converte float64 para decimal.Decimal
func decimalFromFloat(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

// DefaultMesAnterior retorna o período YYYY-MM do mês anterior ao atual.
func (uc *GenerateDREUseCase) DefaultMesAnterior() valueobject.MesAno {
	return valueobject.NewMesAnoFromTime(time.Now().AddDate(0, -1, 0))
}
