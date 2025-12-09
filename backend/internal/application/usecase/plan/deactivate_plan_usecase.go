package plan

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// DeactivatePlanUseCase desativa um plano respeitando PL-003
type DeactivatePlanUseCase struct {
	repo port.PlanRepository
}

// NewDeactivatePlanUseCase cria instância
func NewDeactivatePlanUseCase(repo port.PlanRepository) *DeactivatePlanUseCase {
	return &DeactivatePlanUseCase{repo: repo}
}

// Execute desativa o plano
func (uc *DeactivatePlanUseCase) Execute(ctx context.Context, tenantID, planID string) error {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return domain.ErrInvalidTenantID
	}
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		return domain.ErrInvalidID
	}

	plan, err := uc.repo.GetByID(ctx, planUUID, tenantUUID)
	if err != nil {
		return err
	}
	if plan == nil {
		return domain.ErrPlanNotFound
	}

	count, err := uc.repo.CountActiveSubscriptions(ctx, planUUID, tenantUUID)
	if err != nil {
		return err
	}
	if !plan.CanBeDeleted(count) {
		return domain.ErrPlanHasActiveSubscriptions
	}

	return uc.repo.Deactivate(ctx, planUUID, tenantUUID)
}
