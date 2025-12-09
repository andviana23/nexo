package subscription

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// GetSubscriptionUseCase busca assinatura por ID
type GetSubscriptionUseCase struct {
	subRepo port.SubscriptionRepository
}

// NewGetSubscriptionUseCase cria instância
func NewGetSubscriptionUseCase(subRepo port.SubscriptionRepository) *GetSubscriptionUseCase {
	return &GetSubscriptionUseCase{subRepo: subRepo}
}

// Execute retorna assinatura
func (uc *GetSubscriptionUseCase) Execute(ctx context.Context, tenantID, subscriptionID string) (*dto.SubscriptionResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidTenantID
	}
	subUUID, err := uuid.Parse(subscriptionID)
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
	return mapper.SubscriptionToResponse(sub), nil
}
