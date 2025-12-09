package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// commissionRuleRepository implementa o repository de regras de comissão usando PostgreSQL
type commissionRuleRepository struct {
	queries *db.Queries
}

// NewCommissionRuleRepository cria uma nova instância do repository
func NewCommissionRuleRepository(queries *db.Queries) repository.CommissionRuleRepository {
	return &commissionRuleRepository{
		queries: queries,
	}
}

// Create cria uma nova regra de comissão
func (r *commissionRuleRepository) Create(ctx context.Context, rule *entity.CommissionRule) (*entity.CommissionRule, error) {
	var unitID pgtype.UUID
	if rule.UnitID != nil {
		uid, err := uuid.Parse(*rule.UnitID)
		if err != nil {
			return nil, err
		}
		unitID = pgtype.UUID{Bytes: uid, Valid: true}
	}

	var createdBy pgtype.UUID
	if rule.CreatedBy != nil {
		uid, err := uuid.Parse(*rule.CreatedBy)
		if err != nil {
			return nil, err
		}
		createdBy = pgtype.UUID{Bytes: uid, Valid: true}
	}

	// Converter Priority para *int32
	var priority *int32
	if rule.Priority != nil {
		p := int32(*rule.Priority)
		priority = &p
	}

	result, err := r.queries.CreateCommissionRule(ctx, db.CreateCommissionRuleParams{
		TenantID:        entityUUIDToPgtype(rule.TenantID),
		UnitID:          unitID,
		Name:            rule.Name,
		Description:     rule.Description,
		Type:            rule.Type,
		DefaultRate:     rule.DefaultRate,
		MinAmount:       decimalPtrToNumeric(rule.MinAmount),
		MaxAmount:       decimalPtrToNumeric(rule.MaxAmount),
		CalculationBase: rule.CalculationBase,
		EffectiveFrom:   pgtype.Date{Time: rule.EffectiveFrom, Valid: true},
		EffectiveTo:     timePtrToDate(rule.EffectiveTo),
		Priority:        priority,
		IsActive:        rule.IsActive,
		CreatedBy:       createdBy,
	})
	if err != nil {
		return nil, err
	}

	return commissionRuleToDomain(result), nil
}

// GetByID busca uma regra de comissão por ID
func (r *commissionRuleRepository) GetByID(ctx context.Context, tenantID, id string) (*entity.CommissionRule, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	rid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.GetCommissionRuleByID(ctx, db.GetCommissionRuleByIDParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		ID:       pgtype.UUID{Bytes: rid, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("regra de comissão não encontrada")
		}
		return nil, err
	}

	return commissionRuleToDomain(result), nil
}

// List lista todas as regras de comissão de um tenant
func (r *commissionRuleRepository) List(ctx context.Context, tenantID string) ([]*entity.CommissionRule, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.ListCommissionRulesByTenant(ctx, db.ListCommissionRulesByTenantParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		Limit:    1000,
		Offset:   0,
	})
	if err != nil {
		return nil, err
	}

	rules := make([]*entity.CommissionRule, 0, len(results))
	for _, result := range results {
		rules = append(rules, commissionRuleToDomain(result))
	}

	return rules, nil
}

// ListActive lista regras de comissão ativas de um tenant
func (r *commissionRuleRepository) ListActive(ctx context.Context, tenantID string) ([]*entity.CommissionRule, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.ListCommissionRulesActive(ctx, pgtype.UUID{Bytes: tid, Valid: true})
	if err != nil {
		return nil, err
	}

	rules := make([]*entity.CommissionRule, 0, len(results))
	for _, result := range results {
		rules = append(rules, commissionRuleToDomain(result))
	}

	return rules, nil
}

// GetEffective busca regra de comissão vigente em uma data
func (r *commissionRuleRepository) GetEffective(ctx context.Context, tenantID string, date time.Time) ([]*entity.CommissionRule, error) {
	// Usa ListActive que já filtra por data de vigência
	return r.ListActive(ctx, tenantID)
}

// Update atualiza uma regra de comissão
func (r *commissionRuleRepository) Update(ctx context.Context, rule *entity.CommissionRule) (*entity.CommissionRule, error) {
	id, err := uuid.Parse(rule.ID)
	if err != nil {
		return nil, err
	}

	// Converter Priority para *int32
	var priority *int32
	if rule.Priority != nil {
		p := int32(*rule.Priority)
		priority = &p
	}

	result, err := r.queries.UpdateCommissionRule(ctx, db.UpdateCommissionRuleParams{
		ID:              pgtype.UUID{Bytes: id, Valid: true},
		TenantID:        entityUUIDToPgtype(rule.TenantID),
		Name:            rule.Name,
		Description:     rule.Description,
		Type:            rule.Type,
		DefaultRate:     rule.DefaultRate,
		MinAmount:       decimalPtrToNumeric(rule.MinAmount),
		MaxAmount:       decimalPtrToNumeric(rule.MaxAmount),
		CalculationBase: rule.CalculationBase,
		EffectiveFrom:   pgtype.Date{Time: rule.EffectiveFrom, Valid: true},
		EffectiveTo:     timePtrToDate(rule.EffectiveTo),
		Priority:        priority,
		IsActive:        rule.IsActive,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("regra de comissão não encontrada")
		}
		return nil, err
	}

	return commissionRuleToDomain(result), nil
}

// Delete remove uma regra de comissão
func (r *commissionRuleRepository) Delete(ctx context.Context, tenantID, id string) error {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return err
	}

	rid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteCommissionRule(ctx, db.DeleteCommissionRuleParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		ID:       pgtype.UUID{Bytes: rid, Valid: true},
	})
}

// Deactivate desativa uma regra de comissão
func (r *commissionRuleRepository) Deactivate(ctx context.Context, tenantID, id string) error {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return err
	}

	rid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	_, err = r.queries.DeactivateCommissionRule(ctx, db.DeactivateCommissionRuleParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		ID:       pgtype.UUID{Bytes: rid, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("regra de comissão não encontrada")
		}
		return err
	}

	return nil
}

// GetEffectiveByUnit busca regra vigente específica de uma unidade
// COM-001: Hierarquia nível 3 - Regra da unidade
func (r *commissionRuleRepository) GetEffectiveByUnit(ctx context.Context, tenantID, unitID string, date time.Time) (*entity.CommissionRule, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.Parse(unitID)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.GetCommissionRuleByUnit(ctx, db.GetCommissionRuleByUnitParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		UnitID:   pgtype.UUID{Bytes: uid, Valid: true},
		Column3:  pgtype.Date{Time: date, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Sem regra para esta unidade
		}
		return nil, err
	}

	return commissionRuleToDomain(result), nil
}

// GetEffectiveGlobal busca regra vigente global do tenant (sem unidade)
// COM-001: Hierarquia nível 4 - Regra global do tenant
func (r *commissionRuleRepository) GetEffectiveGlobal(ctx context.Context, tenantID string, date time.Time) (*entity.CommissionRule, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.GetGlobalCommissionRule(ctx, db.GetGlobalCommissionRuleParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		Column2:  pgtype.Date{Time: date, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Sem regra global
		}
		return nil, err
	}

	return commissionRuleToDomain(result), nil
}

// commissionRuleToDomain converte de sqlc model para entity domain
func commissionRuleToDomain(cr db.CommissionRule) *entity.CommissionRule {
	id := uuid.UUID(cr.ID.Bytes).String()
	tenantID := pgtypeToEntityUUID(cr.TenantID)

	var unitID *string
	if cr.UnitID.Valid {
		uid := uuid.UUID(cr.UnitID.Bytes).String()
		unitID = &uid
	}

	var createdBy *string
	if cr.CreatedBy.Valid {
		uid := uuid.UUID(cr.CreatedBy.Bytes).String()
		createdBy = &uid
	}

	var effectiveTo *time.Time
	if cr.EffectiveTo.Valid {
		effectiveTo = &cr.EffectiveTo.Time
	}

	var priority *int
	if cr.Priority != nil {
		p := int(*cr.Priority)
		priority = &p
	}

	// Converter pgtype.Numeric para decimal.Decimal
	var minAmount *decimal.Decimal
	if cr.MinAmount.Valid {
		d := numericToDecimal(cr.MinAmount)
		minAmount = &d
	}

	var maxAmount *decimal.Decimal
	if cr.MaxAmount.Valid {
		d := numericToDecimal(cr.MaxAmount)
		maxAmount = &d
	}

	return &entity.CommissionRule{
		ID:              id,
		TenantID:        tenantID,
		UnitID:          unitID,
		Name:            cr.Name,
		Description:     cr.Description,
		Type:            cr.Type,
		DefaultRate:     cr.DefaultRate,
		MinAmount:       minAmount,
		MaxAmount:       maxAmount,
		CalculationBase: cr.CalculationBase,
		EffectiveFrom:   cr.EffectiveFrom.Time,
		EffectiveTo:     effectiveTo,
		Priority:        priority,
		IsActive:        cr.IsActive,
		CreatedAt:       cr.CreatedAt.Time,
		UpdatedAt:       cr.UpdatedAt.Time,
		CreatedBy:       createdBy,
	}
}
