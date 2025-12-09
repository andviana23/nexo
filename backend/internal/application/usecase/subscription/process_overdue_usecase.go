package subscription

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// ProcessOverdueSubscriptionsUseCase marca assinaturas vencidas como inadimplentes
type ProcessOverdueSubscriptionsUseCase struct {
	subRepo port.SubscriptionRepository
}

// NewProcessOverdueSubscriptionsUseCase cria instância
func NewProcessOverdueSubscriptionsUseCase(subRepo port.SubscriptionRepository) *ProcessOverdueSubscriptionsUseCase {
	return &ProcessOverdueSubscriptionsUseCase{subRepo: subRepo}
}

// Execute processa assinaturas vencidas para um tenant, retorna quantidade atualizada
func (uc *ProcessOverdueSubscriptionsUseCase) Execute(ctx context.Context, tenantID string) (int, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return 0, domain.ErrInvalidTenantID
	}

	subs, err := uc.subRepo.ListByStatus(ctx, tenantUUID, entity.StatusAtivo)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	updated := 0
	for _, sub := range subs {
		if sub.ShouldBecomeInadimplente(now) {
			if err := uc.subRepo.UpdateStatus(ctx, sub.ID, tenantUUID, entity.StatusInadimplente); err == nil {
				updated++
			}
		}
	}
	return updated, nil
}
