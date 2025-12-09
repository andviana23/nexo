package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
)

// PlanRepository define operações de persistência para planos de assinatura.
type PlanRepository interface {
	Create(ctx context.Context, plan *entity.Plan) error
	GetByID(ctx context.Context, id, tenantID uuid.UUID) (*entity.Plan, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.Plan, error)
	ListActiveByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.Plan, error)
	Update(ctx context.Context, plan *entity.Plan) error
	Deactivate(ctx context.Context, id, tenantID uuid.UUID) error
	CountActiveSubscriptions(ctx context.Context, planID, tenantID uuid.UUID) (int, error)
	CheckNameExists(ctx context.Context, tenantID uuid.UUID, nome string, excludeID *uuid.UUID) (bool, error)
}
