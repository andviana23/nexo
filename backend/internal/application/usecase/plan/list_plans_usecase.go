package plan

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// ListPlansUseCase lista planos de um tenant
type ListPlansUseCase struct {
	repo port.PlanRepository
}

// NewListPlansUseCase constrói use case
func NewListPlansUseCase(repo port.PlanRepository) *ListPlansUseCase {
	return &ListPlansUseCase{repo: repo}
}

// Execute lista planos (todos ou apenas ativos conforme flag)
func (uc *ListPlansUseCase) Execute(ctx context.Context, tenantID string, onlyActive bool) ([]*dto.PlanResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidTenantID
	}

	var entities []*entity.Plan
	if onlyActive {
		entities, err = uc.repo.ListActiveByTenant(ctx, tenantUUID)
	} else {
		entities, err = uc.repo.ListByTenant(ctx, tenantUUID)
	}
	if err != nil {
		return nil, err
	}

	return mapper.PlansToResponse(entities), nil
}
