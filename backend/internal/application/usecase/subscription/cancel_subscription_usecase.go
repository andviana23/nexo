package subscription

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CancelSubscriptionUseCase cancela assinatura
// Reference: FLUXO_ASSINATURA.md — Seção 6.5 (Cancelar Assinatura)
type CancelSubscriptionUseCase struct {
	subRepo      port.SubscriptionRepository
	asaasGateway port.AsaasGateway // Pode ser nil
	logger       *zap.Logger
}

// NewCancelSubscriptionUseCase cria instância
func NewCancelSubscriptionUseCase(
	subRepo port.SubscriptionRepository,
	asaasGateway port.AsaasGateway,
	logger *zap.Logger,
) *CancelSubscriptionUseCase {
	return &CancelSubscriptionUseCase{
		subRepo:      subRepo,
		asaasGateway: asaasGateway,
		logger:       logger,
	}
}

// Execute cancela assinatura
func (uc *CancelSubscriptionUseCase) Execute(ctx context.Context, tenantID, subscriptionID, userID string) (*dto.SubscriptionResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidTenantID
	}
	subUUID, err := uuid.Parse(subscriptionID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	sub, err := uc.subRepo.GetByID(ctx, subUUID, tenantUUID)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, domain.ErrSubscriptionNotFound
	}
	if sub.Status == entity.StatusCancelado {
		return nil, domain.ErrSubscriptionAlreadyCanceled
	}

	// Integração Asaas: cancelar no gateway se for pagamento via cartão (RN-CANC-002)
	// Reference: REGRA AS-006
	if sub.FormaPagamento == entity.PaymentMethodCartao && sub.AsaasSubscriptionID != nil && uc.asaasGateway != nil {
		if err := uc.asaasGateway.CancelSubscription(ctx, *sub.AsaasSubscriptionID); err != nil {
			// Log error but continue with local cancellation (soft fail)
			uc.logger.Warn("failed to cancel subscription in Asaas, continuing with local cancellation",
				zap.String("subscription_id", sub.ID.String()),
				zap.String("asaas_subscription_id", *sub.AsaasSubscriptionID),
				zap.Error(err),
			)
		} else {
			uc.logger.Info("subscription canceled in Asaas",
				zap.String("subscription_id", sub.ID.String()),
				zap.String("asaas_subscription_id", *sub.AsaasSubscriptionID),
			)
		}
	}

	// Cancelar localmente
	if err := uc.subRepo.Cancel(ctx, subUUID, tenantUUID, userUUID); err != nil {
		return nil, err
	}

	// Atualizar flag do cliente se necessário
	count, err := uc.subRepo.CountActiveSubscriptionsByCliente(ctx, sub.ClienteID, tenantUUID)
	if err == nil && count == 0 {
		_ = uc.subRepo.SetClienteAsSubscriber(ctx, sub.ClienteID, tenantUUID, false)
	}

	now := time.Now()
	sub.Status = entity.StatusCancelado
	sub.DataCancelamento = &now
	sub.CanceladoPor = &userUUID

	return mapper.SubscriptionToResponse(sub), nil
}
