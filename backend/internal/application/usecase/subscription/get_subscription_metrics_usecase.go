package subscription

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// GetSubscriptionMetricsUseCase retorna métricas agregadas
type GetSubscriptionMetricsUseCase struct {
	subRepo port.SubscriptionRepository
}

// NewGetSubscriptionMetricsUseCase cria instância
func NewGetSubscriptionMetricsUseCase(subRepo port.SubscriptionRepository) *GetSubscriptionMetricsUseCase {
	return &GetSubscriptionMetricsUseCase{subRepo: subRepo}
}

// Execute obtém métricas
func (uc *GetSubscriptionMetricsUseCase) Execute(ctx context.Context, tenantID string) (*dto.SubscriptionMetricsResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidTenantID
	}
	metrics, err := uc.subRepo.GetMetrics(ctx, tenantUUID)
	if err != nil {
		return nil, err
	}
	return mapper.SubscriptionMetricsToResponse(metrics), nil
}
