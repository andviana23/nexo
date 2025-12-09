package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// =============================================================================
// SUBSCRIPTION GUARD MIDDLEWARE - T-ASAAS-003
// Bloqueia tenants inadimplentes (assinatura vencida > N dias)
// =============================================================================

const (
	// DefaultGracePeriodDays é o período de tolerância padrão (5 dias)
	DefaultGracePeriodDays = 5
)

// SubscriptionChecker interface para verificar status de assinatura do tenant
type SubscriptionChecker interface {
	// GetTenantSubscriptionStatus retorna a assinatura ativa do tenant (ou a mais recente)
	GetTenantSubscriptionStatus(ctx context.Context, tenantID uuid.UUID) (*entity.Subscription, error)
	// HasActiveSubscription verifica se o tenant tem assinatura válida
	HasActiveSubscription(ctx context.Context, tenantID uuid.UUID, gracePeriodDays int) (bool, error)
}

// SubscriptionGuardConfig configura o middleware de proteção
type SubscriptionGuardConfig struct {
	Checker         SubscriptionChecker
	Logger          *zap.Logger
	GracePeriodDays int  // Dias de tolerância após vencimento (default: 5)
	BlockOnError    bool // Se true, bloqueia em caso de erro no check (default: false)
}

// SubscriptionGuard cria um middleware que bloqueia tenants inadimplentes
func SubscriptionGuard(config SubscriptionGuardConfig) echo.MiddlewareFunc {
	if config.GracePeriodDays <= 0 {
		config.GracePeriodDays = DefaultGracePeriodDays
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extrair tenant_id do contexto (injetado pelo JWTMiddleware)
			tenantIDStr, ok := c.Get("tenant_id").(string)
			if !ok || tenantIDStr == "" {
				if config.Logger != nil {
					config.Logger.Warn("TenantID não encontrado no token",
						zap.String("path", c.Request().URL.Path),
					)
				}
				// Se não tem tenant_id, deixa passar (pode ser rota pública ou sem tenant)
				return next(c)
			}

			tenantID, err := uuid.Parse(tenantIDStr)
			if err != nil {
				if config.Logger != nil {
					config.Logger.Error("TenantID inválido",
						zap.String("tenant_id", tenantIDStr),
						zap.Error(err),
					)
				}
				return echo.NewHTTPError(http.StatusBadRequest, "TenantID inválido")
			}

			// Verificar status da assinatura
			hasActive, err := config.Checker.HasActiveSubscription(c.Request().Context(), tenantID, config.GracePeriodDays)
			if err != nil {
				if config.Logger != nil {
					config.Logger.Error("Erro ao verificar assinatura do tenant",
						zap.String("tenant_id", tenantIDStr),
						zap.Error(err),
					)
				}
				// Em caso de erro, comportamento configurável
				if config.BlockOnError {
					return echo.NewHTTPError(http.StatusServiceUnavailable, "Erro ao verificar assinatura")
				}
				// Default: deixa passar em caso de erro
				return next(c)
			}

			if !hasActive {
				if config.Logger != nil {
					config.Logger.Warn("Tenant inadimplente bloqueado",
						zap.String("tenant_id", tenantIDStr),
						zap.String("path", c.Request().URL.Path),
						zap.Int("grace_period_days", config.GracePeriodDays),
					)
				}
				return echo.NewHTTPError(http.StatusPaymentRequired, map[string]interface{}{
					"code":    "SUBSCRIPTION_REQUIRED",
					"message": "Sua assinatura está vencida. Por favor, regularize para continuar usando o sistema.",
					"details": map[string]interface{}{
						"grace_period_days": config.GracePeriodDays,
						"action":            "Acesse a área de assinaturas para regularizar",
					},
				})
			}

			// Assinatura válida, continua
			if config.Logger != nil {
				config.Logger.Debug("Assinatura válida, acesso permitido",
					zap.String("tenant_id", tenantIDStr),
					zap.String("path", c.Request().URL.Path),
				)
			}

			return next(c)
		}
	}
}

// RequireActiveSubscription é um helper para criar middleware de verificação de assinatura
func RequireActiveSubscription(checker SubscriptionChecker, logger *zap.Logger) echo.MiddlewareFunc {
	return SubscriptionGuard(SubscriptionGuardConfig{
		Checker:         checker,
		Logger:          logger,
		GracePeriodDays: DefaultGracePeriodDays,
		BlockOnError:    false,
	})
}

// RequireActiveSubscriptionStrict é igual ao anterior, mas bloqueia em caso de erro
func RequireActiveSubscriptionStrict(checker SubscriptionChecker, logger *zap.Logger) echo.MiddlewareFunc {
	return SubscriptionGuard(SubscriptionGuardConfig{
		Checker:         checker,
		Logger:          logger,
		GracePeriodDays: DefaultGracePeriodDays,
		BlockOnError:    true,
	})
}

// =============================================================================
// Default Subscription Checker Implementation
// =============================================================================

// DefaultSubscriptionChecker implementa SubscriptionChecker usando o repositório
type DefaultSubscriptionChecker struct {
	repo   TenantSubscriptionRepo
	logger *zap.Logger
}

// TenantSubscriptionRepo define os métodos necessários do repositório de assinaturas
type TenantSubscriptionRepo interface {
	ListByStatus(ctx context.Context, tenantID uuid.UUID, status entity.SubscriptionStatus) ([]*entity.Subscription, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.Subscription, error)
}

// NewDefaultSubscriptionChecker cria um checker com o repositório de assinaturas
func NewDefaultSubscriptionChecker(repo TenantSubscriptionRepo, logger *zap.Logger) *DefaultSubscriptionChecker {
	return &DefaultSubscriptionChecker{
		repo:   repo,
		logger: logger,
	}
}

// GetTenantSubscriptionStatus retorna a assinatura mais relevante do tenant
func (c *DefaultSubscriptionChecker) GetTenantSubscriptionStatus(ctx context.Context, tenantID uuid.UUID) (*entity.Subscription, error) {
	// Primeiro tenta buscar assinatura ativa
	subs, err := c.repo.ListByStatus(ctx, tenantID, entity.StatusAtivo)
	if err != nil {
		return nil, err
	}

	if len(subs) > 0 {
		// Retorna a primeira ativa (ordenada por mais recente)
		return subs[0], nil
	}

	// Se não tem ativa, busca todas e retorna a mais recente
	allSubs, err := c.repo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	if len(allSubs) == 0 {
		return nil, nil // Tenant sem assinatura
	}

	// Retorna a mais recente
	return allSubs[0], nil
}

// HasActiveSubscription verifica se o tenant tem assinatura válida considerando período de tolerância
func (c *DefaultSubscriptionChecker) HasActiveSubscription(ctx context.Context, tenantID uuid.UUID, gracePeriodDays int) (bool, error) {
	sub, err := c.GetTenantSubscriptionStatus(ctx, tenantID)
	if err != nil {
		return false, err
	}

	// Sem assinatura = sem acesso
	if sub == nil {
		if c.logger != nil {
			c.logger.Debug("Tenant sem assinatura",
				zap.String("tenant_id", tenantID.String()),
			)
		}
		return false, nil
	}

	// Assinatura ativa sem vencimento = acesso permitido
	if sub.Status == entity.StatusAtivo && sub.DataVencimento == nil {
		return true, nil
	}

	// Assinatura ativa com vencimento dentro do período de tolerância = acesso permitido
	if sub.Status == entity.StatusAtivo && sub.DataVencimento != nil {
		deadline := sub.DataVencimento.AddDate(0, 0, gracePeriodDays)
		if time.Now().Before(deadline) {
			return true, nil
		}

		// Vencido há mais de N dias
		if c.logger != nil {
			c.logger.Debug("Assinatura vencida além do período de tolerância",
				zap.String("tenant_id", tenantID.String()),
				zap.Time("vencimento", *sub.DataVencimento),
				zap.Int("grace_period_days", gracePeriodDays),
			)
		}
		return false, nil
	}

	// Status INADIMPLENTE = sem acesso
	if sub.Status == entity.StatusInadimplente {
		if c.logger != nil {
			c.logger.Debug("Assinatura inadimplente",
				zap.String("tenant_id", tenantID.String()),
				zap.String("status", string(sub.Status)),
			)
		}
		return false, nil
	}

	// Status AGUARDANDO_PAGAMENTO com período de carência
	if sub.Status == entity.StatusAguardandoPagamento {
		// Permite acesso durante 7 dias após criação para nova assinatura
		carencia := sub.CreatedAt.AddDate(0, 0, 7)
		if time.Now().Before(carencia) {
			return true, nil
		}

		if c.logger != nil {
			c.logger.Debug("Assinatura aguardando pagamento há muito tempo",
				zap.String("tenant_id", tenantID.String()),
				zap.Time("created_at", sub.CreatedAt),
			)
		}
		return false, nil
	}

	// Outros status (INATIVO, CANCELADO) = sem acesso
	if c.logger != nil {
		c.logger.Debug("Assinatura em status não permitido",
			zap.String("tenant_id", tenantID.String()),
			zap.String("status", string(sub.Status)),
		)
	}
	return false, nil
}
