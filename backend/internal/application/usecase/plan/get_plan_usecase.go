package plan

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// GetPlanUseCase busca um plano por ID
type GetPlanUseCase struct {
	repo port.PlanRepository
}

// NewGetPlanUseCase cria instância
func NewGetPlanUseCase(repo port.PlanRepository) *GetPlanUseCase {
	return &GetPlanUseCase{repo: repo}
}

// Execute executa a busca
func (uc *GetPlanUseCase) Execute(ctx context.Context, tenantID, planID string) (*dto.PlanResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidTenantID
	}
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	plan, err := uc.repo.GetByID(ctx, planUUID, tenantUUID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, domain.ErrPlanNotFound
	}
	return mapper.PlanToResponse(plan), nil
}
