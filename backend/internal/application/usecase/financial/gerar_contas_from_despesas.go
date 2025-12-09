package financial

import (
	"context"

	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// GerarContasFromDespesasFixasInput define os dados de entrada
type GerarContasFromDespesasFixasInput struct {
	TenantID string // Se vazio, processa todos os tenants
	Ano      int
	Mes      int
}

// GerarContasFromDespesasFixasOutput define os dados de saída
type GerarContasFromDespesasFixasOutput struct {
	TotalDespesas   int
	ContasCriadas   int
	Erros           int
	DetalhesErros   []string
	TempoExecucaoMs int64
}

// GerarContasFromDespesasFixasUseCase gera contas a pagar a partir de despesas fixas
type GerarContasFromDespesasFixasUseCase struct {
	despesaFixaRepo port.DespesaFixaRepository
	contaPagarRepo  port.ContaPagarRepository
	logger          *zap.Logger
}

// NewGerarContasFromDespesasFixasUseCase cria nova instância do use case
func NewGerarContasFromDespesasFixasUseCase(
	despesaFixaRepo port.DespesaFixaRepository,
	contaPagarRepo port.ContaPagarRepository,
	logger *zap.Logger,
) *GerarContasFromDespesasFixasUseCase {
	return &GerarContasFromDespesasFixasUseCase{
		despesaFixaRepo: despesaFixaRepo,
		contaPagarRepo:  contaPagarRepo,
		logger:          logger,
	}
}

// Execute gera as contas a pagar para o mês especificado
func (uc *GerarContasFromDespesasFixasUseCase) Execute(ctx context.Context, input GerarContasFromDespesasFixasInput) (*GerarContasFromDespesasFixasOutput, error) {
	startTime := time.Now()

	output := &GerarContasFromDespesasFixasOutput{}

	// Validar mês/ano
	if input.Ano == 0 || input.Mes == 0 {
		// Se não especificado, usar mês atual
		now := time.Now()
		input.Ano = now.Year()
		input.Mes = int(now.Month())
	}

	// Validar range de mês
	if input.Mes < 1 || input.Mes > 12 {
		return nil, fmt.Errorf("mês inválido: %d", input.Mes)
	}

	uc.logger.Info("Iniciando geração de contas a pagar a partir de despesas fixas",
		zap.Int("ano", input.Ano),
		zap.Int("mes", input.Mes),
		zap.String("tenant_id", input.TenantID),
	)

	var despesas []*entity.DespesaFixa
	var err error

	if input.TenantID != "" {
		// Processar apenas um tenant específico
		despesas, err = uc.despesaFixaRepo.ListAtivas(ctx, input.TenantID)
	} else {
		// Processar todos os tenants (usado pelo cron job)
		despesasComTenant, err := uc.despesaFixaRepo.ListAtivasPorTenants(ctx)
		if err != nil {
			uc.logger.Error("Erro ao listar despesas fixas ativas por tenants",
				zap.Error(err),
			)
			return nil, err
		}
		for _, dt := range despesasComTenant {
			despesas = append(despesas, dt.DespesaFixa)
		}
	}

	if err != nil {
		uc.logger.Error("Erro ao listar despesas fixas ativas",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
		)
		return nil, err
	}

	output.TotalDespesas = len(despesas)

	// Gerar conta a pagar para cada despesa fixa
	for _, despesa := range despesas {
		conta, err := despesa.ToContaPagar(input.Ano, input.Mes)
		if err != nil {
			output.Erros++
			errMsg := fmt.Sprintf("Despesa %s (tenant: %s): %v", despesa.ID, despesa.TenantID, err)
			output.DetalhesErros = append(output.DetalhesErros, errMsg)
			uc.logger.Warn("Erro ao converter despesa fixa para conta a pagar",
				zap.Error(err),
				zap.String("despesa_id", despesa.ID),
				zap.String("tenant_id", despesa.TenantID.String()),
			)
			continue
		}

		// Persistir conta a pagar
		if err := uc.contaPagarRepo.Create(ctx, conta); err != nil {
			output.Erros++
			errMsg := fmt.Sprintf("Despesa %s (tenant: %s): erro ao persistir - %v", despesa.ID, despesa.TenantID, err)
			output.DetalhesErros = append(output.DetalhesErros, errMsg)
			uc.logger.Warn("Erro ao criar conta a pagar a partir de despesa fixa",
				zap.Error(err),
				zap.String("despesa_id", despesa.ID),
				zap.String("tenant_id", despesa.TenantID.String()),
			)
			continue
		}

		output.ContasCriadas++
		uc.logger.Debug("Conta a pagar criada a partir de despesa fixa",
			zap.String("conta_id", conta.ID),
			zap.String("despesa_id", despesa.ID),
			zap.String("tenant_id", despesa.TenantID.String()),
			zap.String("descricao", conta.Descricao),
			zap.Time("vencimento", conta.DataVencimento),
		)
	}

	output.TempoExecucaoMs = time.Since(startTime).Milliseconds()

	uc.logger.Info("Geração de contas a pagar concluída",
		zap.Int("total_despesas", output.TotalDespesas),
		zap.Int("contas_criadas", output.ContasCriadas),
		zap.Int("erros", output.Erros),
		zap.Int64("tempo_ms", output.TempoExecucaoMs),
	)

	return output, nil
}

// ExecuteForCurrentMonth é um atalho para executar no mês atual
func (uc *GerarContasFromDespesasFixasUseCase) ExecuteForCurrentMonth(ctx context.Context, tenantID string) (*GerarContasFromDespesasFixasOutput, error) {
	now := time.Now()
	return uc.Execute(ctx, GerarContasFromDespesasFixasInput{
		TenantID: tenantID,
		Ano:      now.Year(),
		Mes:      int(now.Month()),
	})
}

// ExecuteForNextMonth é um atalho para executar no próximo mês
func (uc *GerarContasFromDespesasFixasUseCase) ExecuteForNextMonth(ctx context.Context, tenantID string) (*GerarContasFromDespesasFixasOutput, error) {
	now := time.Now().AddDate(0, 1, 0)
	return uc.Execute(ctx, GerarContasFromDespesasFixasInput{
		TenantID: tenantID,
		Ano:      now.Year(),
		Mes:      int(now.Month()),
	})
}
