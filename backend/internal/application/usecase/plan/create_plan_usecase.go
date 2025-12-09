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

// CreatePlanUseCase cria um novo plano de assinatura
type CreatePlanUseCase struct {
	repo   port.PlanRepository
	logger *zap.Logger
}

// NewCreatePlanUseCase constrói o use case
func NewCreatePlanUseCase(repo port.PlanRepository, logger *zap.Logger) *CreatePlanUseCase {
	return &CreatePlanUseCase{repo: repo, logger: logger}
}

// Execute executa a criação do plano
func (uc *CreatePlanUseCase) Execute(ctx context.Context, tenantID string, req dto.CreatePlanRequest) (*dto.PlanResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidTenantID
	}

	exists, err := uc.repo.CheckNameExists(ctx, tenantUUID, req.Nome, nil)
	if err != nil {
		uc.logger.Error("erro ao verificar nome do plano", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, domain.ErrPlanNameDuplicate
	}

	plan, err := mapper.CreatePlanRequestToEntity(&req, tenantUUID)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, plan); err != nil {
		uc.logger.Error("erro ao criar plano", zap.Error(err))
		return nil, err
	}

	return mapper.PlanToResponse(plan), nil
}
