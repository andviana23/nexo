package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PlanRepositoryPG implementa port.PlanRepository usando sqlc/Postgres.
type PlanRepositoryPG struct {
	queries *db.Queries
}

// NewPlanRepository cria nova instância do repositório de planos.
func NewPlanRepository(queries *db.Queries) port.PlanRepository {
	return &PlanRepositoryPG{queries: queries}
}

// Create persiste um novo plano.
func (r *PlanRepositoryPG) Create(ctx context.Context, plan *entity.Plan) error {
	params := db.CreatePlanParams{
		TenantID:        uuidToPgUUID(plan.TenantID),
		Nome:            plan.Nome,
		Descricao:       plan.Descricao,
		Valor:           plan.Valor,
		Periodicidade:   plan.Periodicidade,
		QtdServicos:     intToInt32Ptr(plan.QtdServicos),
		LimiteUsoMensal: intToInt32Ptr(plan.LimiteUsoMensal),
		Ativo:           plan.Ativo,
	}

	result, err := r.queries.CreatePlan(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar plano: %w", err)
	}

	r.applyRowToEntity(&result, plan)
	return nil
}

// GetByID retorna um plano por ID+tenant.
func (r *PlanRepositoryPG) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*entity.Plan, error) {
	row, err := r.queries.GetPlanByID(ctx, db.GetPlanByIDParams{
		ID:       uuidToPgUUID(id),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrPlanNotFound
		}
		return nil, fmt.Errorf("erro ao buscar plano: %w", err)
	}
	plan := r.rowToEntity(&row)
	return plan, nil
}

// ListByTenant lista todos os planos do tenant.
func (r *PlanRepositoryPG) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.Plan, error) {
	rows, err := r.queries.ListPlansByTenant(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar planos: %w", err)
	}
	return r.rowsToEntities(rows), nil
}

// ListActiveByTenant lista planos ativos (PL-002).
func (r *PlanRepositoryPG) ListActiveByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.Plan, error) {
	rows, err := r.queries.ListActivePlansByTenant(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar planos ativos: %w", err)
	}
	return r.rowsToEntities(rows), nil
}

// Update atualiza um plano existente.
func (r *PlanRepositoryPG) Update(ctx context.Context, plan *entity.Plan) error {
	params := db.UpdatePlanParams{
		ID:              uuidToPgUUID(plan.ID),
		TenantID:        uuidToPgUUID(plan.TenantID),
		Nome:            plan.Nome,
		Descricao:       plan.Descricao,
		Valor:           plan.Valor,
		QtdServicos:     intToInt32Ptr(plan.QtdServicos),
		LimiteUsoMensal: intToInt32Ptr(plan.LimiteUsoMensal),
		Ativo:           plan.Ativo,
	}

	result, err := r.queries.UpdatePlan(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar plano: %w", err)
	}

	plan.UpdatedAt = result.UpdatedAt.Time
	return nil
}

// Deactivate desativa um plano (soft delete).
func (r *PlanRepositoryPG) Deactivate(ctx context.Context, id, tenantID uuid.UUID) error {
	if err := r.queries.DeactivatePlan(ctx, db.DeactivatePlanParams{
		ID:       uuidToPgUUID(id),
		TenantID: uuidToPgUUID(tenantID),
	}); err != nil {
		return fmt.Errorf("erro ao desativar plano: %w", err)
	}
	return nil
}

// CountActiveSubscriptions retorna total de assinaturas ativas ligadas ao plano.
func (r *PlanRepositoryPG) CountActiveSubscriptions(ctx context.Context, planID, tenantID uuid.UUID) (int, error) {
	count, err := r.queries.CountActiveSubscriptionsByPlan(ctx, db.CountActiveSubscriptionsByPlanParams{
		PlanoID:  uuidToPgUUID(planID),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		return 0, fmt.Errorf("erro ao contar assinaturas ativas do plano: %w", err)
	}
	return int(count), nil
}

// CheckNameExists verifica duplicidade do nome (PL-005)
func (r *PlanRepositoryPG) CheckNameExists(ctx context.Context, tenantID uuid.UUID, nome string, excludeID *uuid.UUID) (bool, error) {
	exclude := uuid.Nil
	if excludeID != nil {
		exclude = *excludeID
	}

	exists, err := r.queries.CheckPlanNameExists(ctx, db.CheckPlanNameExistsParams{
		TenantID: uuidToPgUUID(tenantID),
		Nome:     nome,
		ID:       uuidToPgUUID(exclude),
	})
	if err != nil {
		return false, fmt.Errorf("erro ao verificar nome do plano: %w", err)
	}
	return exists, nil
}

// Helpers
func (r *PlanRepositoryPG) rowToEntity(row *db.Plan) *entity.Plan {
	if row == nil {
		return nil
	}
	plan := &entity.Plan{}
	r.applyRowToEntity(row, plan)
	return plan
}

func (r *PlanRepositoryPG) rowsToEntities(rows []db.Plan) []*entity.Plan {
	result := make([]*entity.Plan, 0, len(rows))
	for i := range rows {
		result = append(result, r.rowToEntity(&rows[i]))
	}
	return result
}

func (r *PlanRepositoryPG) applyRowToEntity(row *db.Plan, plan *entity.Plan) {
	plan.ID = pgUUIDToUUID(row.ID)
	plan.TenantID = pgUUIDToUUID(row.TenantID)
	plan.Nome = row.Nome
	plan.Descricao = row.Descricao
	plan.Valor = row.Valor
	plan.Periodicidade = row.Periodicidade
	plan.QtdServicos = int32PtrToIntPtr(row.QtdServicos)
	plan.LimiteUsoMensal = int32PtrToIntPtr(row.LimiteUsoMensal)
	plan.Ativo = row.Ativo
	plan.CreatedAt = row.CreatedAt.Time
	plan.UpdatedAt = row.UpdatedAt.Time
}

// intToInt32Ptr converte *int para *int32 (sqlc)
func intToInt32Ptr(v *int) *int32 {
	if v == nil {
		return nil
	}
	vv := int32(*v)
	return &vv
}

// int32PtrToIntPtr converte *int32 para *int
func int32PtrToIntPtr(v *int32) *int {
	if v == nil {
		return nil
	}
	vv := int(*v)
	return &vv
}
