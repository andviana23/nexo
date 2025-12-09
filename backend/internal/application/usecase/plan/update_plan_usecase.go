package plan

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UpdatePlanUseCase atualiza dados de um plano
type UpdatePlanUseCase struct {
	repo   port.PlanRepository
	logger *zap.Logger
}

// NewUpdatePlanUseCase constrói use case
func NewUpdatePlanUseCase(repo port.PlanRepository, logger *zap.Logger) *UpdatePlanUseCase {
	return &UpdatePlanUseCase{repo: repo, logger: logger}
}

// Execute atualiza o plano
func (uc *UpdatePlanUseCase) Execute(ctx context.Context, tenantID, planID string, req dto.UpdatePlanRequest) (*dto.PlanResponse, error) {
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

	exists, err := uc.repo.CheckNameExists(ctx, tenantUUID, req.Nome, &planUUID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrPlanNameDuplicate
	}

	if err := mapper.UpdatePlanRequestToEntity(plan, &req); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, plan); err != nil {
		uc.logger.Error("erro ao atualizar plano", zap.Error(err))
		return nil, err
	}

	return mapper.PlanToResponse(plan), nil
}
