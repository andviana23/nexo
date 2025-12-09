package subscription

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// ListSubscriptionsUseCase lista assinaturas
type ListSubscriptionsUseCase struct {
	subRepo port.SubscriptionRepository
}

// NewListSubscriptionsUseCase cria instância
func NewListSubscriptionsUseCase(subRepo port.SubscriptionRepository) *ListSubscriptionsUseCase {
	return &ListSubscriptionsUseCase{subRepo: subRepo}
}

// Execute lista assinaturas, opcionalmente filtrando por status
func (uc *ListSubscriptionsUseCase) Execute(ctx context.Context, tenantID string, status *string) ([]*dto.SubscriptionResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidTenantID
	}

	var subs []*entity.Subscription
	if status != nil && *status != "" {
		st := entity.SubscriptionStatus(*status)
		subs, err = uc.subRepo.ListByStatus(ctx, tenantUUID, st)
	} else {
		subs, err = uc.subRepo.ListByTenant(ctx, tenantUUID)
	}
	if err != nil {
		return nil, err
	}

	return mapper.SubscriptionsToResponse(subs), nil
}
