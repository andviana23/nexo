package scheduler

import (
	"context"
	"os"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/financial"
	subscriptionUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/subscription"
	"go.uber.org/zap"
)

// FinancialJobDeps agrega use cases necessários para os cron jobs financeiros.
type FinancialJobDeps struct {
	GenerateDRE              *financial.GenerateDREUseCase
	GenerateDREV2            *financial.GenerateDREV2UseCase
	GenerateFluxoDiario      *financial.GenerateFluxoDiarioUseCase
	GenerateFluxoDiarioV2    *financial.GenerateFluxoDiarioV2UseCase
	MarcarCompensacoes       *financial.MarcarCompensacaoUseCase
	GerarContasDespesasFixas *financial.GerarContasFromDespesasFixasUseCase
}

// SubscriptionJobDeps agrega use cases do módulo de assinaturas
type SubscriptionJobDeps struct {
	ProcessOverdue *subscriptionUC.ProcessOverdueSubscriptionsUseCase
}

// RegisterFinancialJobs registra os cron jobs financeiros com base nas envs.
func RegisterFinancialJobs(s *Scheduler, logger *zap.Logger, deps FinancialJobDeps, tenants []string) error {
	// DRE mensal (mês anterior, todo dia 1 às 03:00)
	if deps.GenerateDREV2 != nil {
		if err := s.AddJob(JobConfig{
			Name:        "GenerateDREMonthly",
			Schedule:    getEnvSchedule("CRON_DRE_MONTHLY_SCHEDULE", "0 0 3 1 * *"),
			Enabled:     getEnvBool("CRON_DRE_MONTHLY_ENABLED", true),
			FeatureFlag: "FF_CRON_DRE_MONTHLY",
			Tenants:     tenants,
			TenantRunner: func(ctx context.Context, tenantID string) error {
				_, err := deps.GenerateDREV2.Execute(ctx, financial.GenerateDREV2Input{
					TenantID: tenantID,
					MesAno:   deps.GenerateDREV2.DefaultMesAnterior(),
					Regime:   "COMPETENCIA",
				})
				return err
			},
		}); err != nil {
			return err
		}
	} else if deps.GenerateDRE != nil {
		if err := s.AddJob(JobConfig{
			Name:        "GenerateDREMonthly",
			Schedule:    getEnvSchedule("CRON_DRE_MONTHLY_SCHEDULE", "0 0 3 1 * *"),
			Enabled:     getEnvBool("CRON_DRE_MONTHLY_ENABLED", true),
			FeatureFlag: "FF_CRON_DRE_MONTHLY",
			Tenants:     tenants,
			TenantRunner: func(ctx context.Context, tenantID string) error {
				_, err := deps.GenerateDRE.Execute(ctx, financial.GenerateDREInput{
					TenantID: tenantID,
					MesAno:   deps.GenerateDRE.DefaultMesAnterior(),
				})
				return err
			},
		}); err != nil {
			return err
		}
	} else {
		logger.Warn("GenerateDREUseCase não configurado, job não registrado")
	}

	// Fluxo diário (00:05 todos os dias)
	if deps.GenerateFluxoDiarioV2 != nil {
		if err := s.AddJob(JobConfig{
			Name:        "GenerateFluxoDiario",
			Schedule:    getEnvSchedule("CRON_FLUXO_DIARIO_SCHEDULE", "0 5 0 * * *"),
			Enabled:     getEnvBool("CRON_FLUXO_DIARIO_ENABLED", true),
			FeatureFlag: "FF_CRON_FLUXO_DIARIO",
			Tenants:     tenants,
			TenantRunner: func(ctx context.Context, tenantID string) error {
				_, err := deps.GenerateFluxoDiarioV2.Execute(ctx, financial.GenerateFluxoDiarioV2Input{
					TenantID: tenantID,
					Data:     time.Now(),
				})
				return err
			},
		}); err != nil {
			return err
		}
	} else if deps.GenerateFluxoDiario != nil {
		if err := s.AddJob(JobConfig{
			Name:        "GenerateFluxoDiario",
			Schedule:    getEnvSchedule("CRON_FLUXO_DIARIO_SCHEDULE", "0 5 0 * * *"),
			Enabled:     getEnvBool("CRON_FLUXO_DIARIO_ENABLED", true),
			FeatureFlag: "FF_CRON_FLUXO_DIARIO",
			Tenants:     tenants,
			TenantRunner: func(ctx context.Context, tenantID string) error {
				_, err := deps.GenerateFluxoDiario.Execute(ctx, financial.GenerateFluxoDiarioInput{
					TenantID: tenantID,
					Data:     time.Now(),
				})
				return err
			},
		}); err != nil {
			return err
		}
	} else {
		logger.Warn("GenerateFluxoDiarioUseCase não configurado, job não registrado")
	}

	// Marcar compensações (diário)
	if deps.MarcarCompensacoes != nil {
		if err := s.AddJob(JobConfig{
			Name:        "MarcarCompensacoes",
			Schedule:    getEnvSchedule("CRON_COMPENSACOES_SCHEDULE", "0 10 0 * * *"),
			Enabled:     getEnvBool("CRON_COMPENSACOES_ENABLED", true),
			FeatureFlag: "FF_CRON_COMPENSACOES",
			Tenants:     tenants,
			TenantRunner: func(ctx context.Context, tenantID string) error {
				_, err := deps.MarcarCompensacoes.ExecuteBatch(ctx, tenantID)
				return err
			},
		}); err != nil {
			return err
		}
	} else {
		logger.Warn("MarcarCompensacoesUseCase não configurado, job não registrado")
	}

	// Gerar Contas de Despesas Fixas (todo dia 1 às 01:00)
	// Cria automaticamente ContasPagar com base nas despesas fixas ativas
	if deps.GerarContasDespesasFixas != nil {
		if err := s.AddJob(JobConfig{
			Name:        "GerarContasDespesasFixas",
			Schedule:    getEnvSchedule("CRON_DESPESAS_FIXAS_SCHEDULE", "0 0 1 1 * *"),
			Enabled:     getEnvBool("CRON_DESPESAS_FIXAS_ENABLED", true),
			FeatureFlag: "FF_CRON_DESPESAS_FIXAS",
			Tenants:     tenants,
			TenantRunner: func(ctx context.Context, tenantID string) error {
				_, err := deps.GerarContasDespesasFixas.ExecuteForCurrentMonth(ctx, tenantID)
				return err
			},
		}); err != nil {
			return err
		}
	} else {
		logger.Warn("GerarContasDespesasFixasUseCase não configurado, job não registrado")
	}

	// Jobs ainda não implementados em domínio: registrar placeholders para não quebrar build.
	placeholderJobs := []JobConfig{
		{
			Name:        "NotifyPayables",
			Schedule:    getEnvSchedule("CRON_NOTIFY_PAYABLES_SCHEDULE", "0 30 8 * * *"),
			Enabled:     getEnvBool("CRON_NOTIFY_PAYABLES_ENABLED", false),
			FeatureFlag: "FF_CRON_NOTIFY_PAYABLES",
			Job: func(ctx context.Context) error {
				logger.Info("NotifyPayables não implementado - noop", zap.String("job", "NotifyPayables"))
				return nil
			},
		},
		{
			Name:        "CheckEstoqueMinimo",
			Schedule:    getEnvSchedule("CRON_CHECK_ESTOQUE_SCHEDULE", "0 15 1 * * *"),
			Enabled:     getEnvBool("CRON_CHECK_ESTOQUE_ENABLED", false),
			FeatureFlag: "FF_CRON_CHECK_ESTOQUE",
			Job: func(ctx context.Context) error {
				logger.Info("CheckEstoqueMinimo não implementado - noop", zap.String("job", "CheckEstoqueMinimo"))
				return nil
			},
		},
		{
			Name:        "CalculateComissoes",
			Schedule:    getEnvSchedule("CRON_COMISSOES_SCHEDULE", "0 0 2 1 * *"),
			Enabled:     getEnvBool("CRON_COMISSOES_ENABLED", false),
			FeatureFlag: "FF_CRON_COMISSOES",
			Job: func(ctx context.Context) error {
				logger.Info("CalculateComissoes não implementado - noop", zap.String("job", "CalculateComissoes"))
				return nil
			},
		},
	}

	for _, job := range placeholderJobs {
		if err := s.AddJob(job); err != nil {
			return err
		}
	}

	return nil
}

// RegisterSubscriptionJobs registra cron jobs do módulo de assinaturas
func RegisterSubscriptionJobs(s *Scheduler, logger *zap.Logger, deps SubscriptionJobDeps, tenants []string) error {
	// Verificar vencimentos (RN-VENC-003/RN-VENC-004)
	if deps.ProcessOverdue != nil {
		if err := s.AddJob(JobConfig{
			Name:        "ProcessOverdueSubscriptions",
			Schedule:    getEnvSchedule("CRON_SUBSCRIPTIONS_OVERDUE_SCHEDULE", "0 5 0 * * *"),
			Enabled:     getEnvBool("CRON_SUBSCRIPTIONS_OVERDUE_ENABLED", true),
			FeatureFlag: "FF_CRON_SUBSCRIPTIONS_OVERDUE",
			Tenants:     tenants,
			TenantRunner: func(ctx context.Context, tenantID string) error {
				_, err := deps.ProcessOverdue.Execute(ctx, tenantID)
				return err
			},
		}); err != nil {
			return err
		}
	} else {
		logger.Warn("ProcessOverdueSubscriptionsUseCase não configurado, job não registrado")
	}
	return nil
}

func getEnvSchedule(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func getEnvBool(key string, def bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	switch val {
	case "0", "false", "FALSE", "False":
		return false
	default:
		return true
	}
}
