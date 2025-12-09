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
)

// advanceRepository implementa o repository de adiantamentos usando PostgreSQL
type advanceRepository struct {
	queries *db.Queries
}

// NewAdvanceRepository cria uma nova instância do repository
func NewAdvanceRepository(queries *db.Queries) repository.AdvanceRepository {
	return &advanceRepository{
		queries: queries,
	}
}

// Create cria um novo adiantamento
func (r *advanceRepository) Create(ctx context.Context, advance *entity.Advance) (*entity.Advance, error) {
	professionalID, err := uuid.Parse(advance.ProfessionalID)
	if err != nil {
		return nil, err
	}

	var unitID pgtype.UUID
	if advance.UnitID != nil {
		uid, err := uuid.Parse(*advance.UnitID)
		if err != nil {
			return nil, err
		}
		unitID = pgtype.UUID{Bytes: uid, Valid: true}
	}

	var createdBy pgtype.UUID
	if advance.CreatedBy != nil {
		uid, err := uuid.Parse(*advance.CreatedBy)
		if err != nil {
			return nil, err
		}
		createdBy = pgtype.UUID{Bytes: uid, Valid: true}
	}

	result, err := r.queries.CreateAdvance(ctx, db.CreateAdvanceParams{
		TenantID:       entityUUIDToPgtype(advance.TenantID),
		UnitID:         unitID,
		ProfessionalID: pgtype.UUID{Bytes: professionalID, Valid: true},
		Amount:         advance.Amount,
		RequestDate:    pgtype.Date{Time: advance.RequestDate, Valid: true},
		Reason:         advance.Reason,
		Status:         advance.Status,
		CreatedBy:      createdBy,
	})
	if err != nil {
		return nil, err
	}

	return advanceToDomain(result), nil
}

// GetByID busca um adiantamento por ID
func (r *advanceRepository) GetByID(ctx context.Context, tenantID, id string) (*entity.Advance, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	aid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.GetAdvanceByID(ctx, db.GetAdvanceByIDParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		ID:       pgtype.UUID{Bytes: aid, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("adiantamento não encontrado")
		}
		return nil, err
	}

	return advanceWithNamesToDomain(result), nil
}

// List lista adiantamentos com filtros
func (r *advanceRepository) List(ctx context.Context, tenantID string, professionalID *string, status *string, limit, offset int) ([]*entity.Advance, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	// Se status for fornecido, usa query por status
	if status != nil && *status != "" {
		results, err := r.queries.ListAdvancesByStatus(ctx, db.ListAdvancesByStatusParams{
			TenantID: pgtype.UUID{Bytes: tid, Valid: true},
			Status:   *status,
			Limit:    int32(limit),
			Offset:   int32(offset),
		})
		if err != nil {
			return nil, err
		}

		advances := make([]*entity.Advance, 0, len(results))
		for _, result := range results {
			advances = append(advances, advanceByStatusRowToDomain(result))
		}
		return advances, nil
	}

	// Se professionalID for fornecido, usa query por profissional
	if professionalID != nil && *professionalID != "" {
		pid, err := uuid.Parse(*professionalID)
		if err != nil {
			return nil, err
		}

		results, err := r.queries.ListAdvancesByProfessional(ctx, db.ListAdvancesByProfessionalParams{
			TenantID:       pgtype.UUID{Bytes: tid, Valid: true},
			ProfessionalID: pgtype.UUID{Bytes: pid, Valid: true},
			Limit:          int32(limit),
			Offset:         int32(offset),
		})
		if err != nil {
			return nil, err
		}

		advances := make([]*entity.Advance, 0, len(results))
		for _, result := range results {
			advances = append(advances, advanceByProfessionalRowToDomain(result))
		}
		return advances, nil
	}

	// Query padrão por tenant
	results, err := r.queries.ListAdvancesByTenant(ctx, db.ListAdvancesByTenantParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}

	advances := make([]*entity.Advance, 0, len(results))
	for _, result := range results {
		advances = append(advances, advanceByTenantRowToDomain(result))
	}
	return advances, nil
}

// GetByProfessional busca adiantamentos por profissional
func (r *advanceRepository) GetByProfessional(ctx context.Context, tenantID, professionalID string) ([]*entity.Advance, error) {
	return r.List(ctx, tenantID, &professionalID, nil, 1000, 0)
}

// GetPendingByProfessional busca adiantamentos pendentes de um profissional
func (r *advanceRepository) GetPendingByProfessional(ctx context.Context, tenantID, professionalID string) ([]*entity.Advance, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	// Usa query de pendentes do tenant e filtra por profissional
	results, err := r.queries.ListPendingAdvances(ctx, pgtype.UUID{Bytes: tid, Valid: true})
	if err != nil {
		return nil, err
	}

	advances := make([]*entity.Advance, 0)
	for _, result := range results {
		pid := uuid.UUID(result.ProfessionalID.Bytes).String()
		if pid == professionalID {
			advances = append(advances, advancePendingRowToDomain(result))
		}
	}
	return advances, nil
}

// GetApprovedByProfessional busca adiantamentos aprovados (não deduzidos) de um profissional
func (r *advanceRepository) GetApprovedByProfessional(ctx context.Context, tenantID, professionalID string) ([]*entity.Advance, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	pid, err := uuid.Parse(professionalID)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.ListApprovedAdvancesForProfessional(ctx, db.ListApprovedAdvancesForProfessionalParams{
		TenantID:       pgtype.UUID{Bytes: tid, Valid: true},
		ProfessionalID: pgtype.UUID{Bytes: pid, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	advances := make([]*entity.Advance, 0, len(results))
	for _, result := range results {
		advances = append(advances, advanceToDomain(result))
	}
	return advances, nil
}

// GetByDateRange busca adiantamentos por intervalo de datas
func (r *advanceRepository) GetByDateRange(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*entity.Advance, error) {
	// Usa List com filtro padrão
	return r.List(ctx, tenantID, nil, nil, 1000, 0)
}

// GetTotalPendingByProfessional retorna o total pendente de um profissional
func (r *advanceRepository) GetTotalPendingByProfessional(ctx context.Context, tenantID, professionalID string) (float64, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return 0, err
	}

	pid, err := uuid.Parse(professionalID)
	if err != nil {
		return 0, err
	}

	result, err := r.queries.SumPendingAdvancesByProfessional(ctx, db.SumPendingAdvancesByProfessionalParams{
		TenantID:       pgtype.UUID{Bytes: tid, Valid: true},
		ProfessionalID: pgtype.UUID{Bytes: pid, Valid: true},
	})
	if err != nil {
		return 0, err
	}

	total, _ := result.Float64()
	return total, nil
}

// GetTotalApprovedByProfessional retorna o total aprovado (não deduzido) de um profissional
func (r *advanceRepository) GetTotalApprovedByProfessional(ctx context.Context, tenantID, professionalID string) (float64, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return 0, err
	}

	pid, err := uuid.Parse(professionalID)
	if err != nil {
		return 0, err
	}

	result, err := r.queries.SumApprovedAdvancesByProfessional(ctx, db.SumApprovedAdvancesByProfessionalParams{
		TenantID:       pgtype.UUID{Bytes: tid, Valid: true},
		ProfessionalID: pgtype.UUID{Bytes: pid, Valid: true},
	})
	if err != nil {
		return 0, err
	}

	total, _ := result.Float64()
	return total, nil
}

// Approve aprova um adiantamento
func (r *advanceRepository) Approve(ctx context.Context, tenantID, id, approvedBy string) (*entity.Advance, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	aid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	approvedByUUID, err := uuid.Parse(approvedBy)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.ApproveAdvance(ctx, db.ApproveAdvanceParams{
		ID:         pgtype.UUID{Bytes: aid, Valid: true},
		TenantID:   pgtype.UUID{Bytes: tid, Valid: true},
		ApprovedBy: pgtype.UUID{Bytes: approvedByUUID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("adiantamento não encontrado ou não pode ser aprovado")
		}
		return nil, err
	}

	return advanceToDomain(result), nil
}

// Reject rejeita um adiantamento com motivo
func (r *advanceRepository) Reject(ctx context.Context, tenantID, id, rejectedBy, reason string) (*entity.Advance, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	aid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	rejectedByUUID, err := uuid.Parse(rejectedBy)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.RejectAdvance(ctx, db.RejectAdvanceParams{
		ID:              pgtype.UUID{Bytes: aid, Valid: true},
		TenantID:        pgtype.UUID{Bytes: tid, Valid: true},
		RejectedBy:      pgtype.UUID{Bytes: rejectedByUUID, Valid: true},
		RejectionReason: &reason,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("adiantamento não encontrado ou não pode ser rejeitado")
		}
		return nil, err
	}

	return advanceToDomain(result), nil
}

// MarkDeducted marca um adiantamento como deduzido
func (r *advanceRepository) MarkDeducted(ctx context.Context, tenantID, id, periodID string) (*entity.Advance, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	aid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	periodUUID, err := uuid.Parse(periodID)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.DeductAdvance(ctx, db.DeductAdvanceParams{
		ID:                pgtype.UUID{Bytes: aid, Valid: true},
		TenantID:          pgtype.UUID{Bytes: tid, Valid: true},
		DeductionPeriodID: pgtype.UUID{Bytes: periodUUID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("adiantamento não encontrado ou não pode ser deduzido")
		}
		return nil, err
	}

	return advanceToDomain(result), nil
}

// Cancel cancela um adiantamento
func (r *advanceRepository) Cancel(ctx context.Context, tenantID, id string) (*entity.Advance, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	aid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.CancelAdvance(ctx, db.CancelAdvanceParams{
		ID:       pgtype.UUID{Bytes: aid, Valid: true},
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("adiantamento não encontrado ou não pode ser cancelado")
		}
		return nil, err
	}

	return advanceToDomain(result), nil
}

// Update atualiza um adiantamento
func (r *advanceRepository) Update(ctx context.Context, advance *entity.Advance) (*entity.Advance, error) {
	// Por enquanto retorna o adiantamento sem atualizar
	// Implementação completa requer query de UPDATE específica
	return advance, nil
}

// Delete remove um adiantamento (somente se PENDING)
func (r *advanceRepository) Delete(ctx context.Context, tenantID, id string) error {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return err
	}

	aid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteAdvance(ctx, db.DeleteAdvanceParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		ID:       pgtype.UUID{Bytes: aid, Valid: true},
	})
}

// === Funções auxiliares de conversão ===

// advanceToDomain converte de db.Advance para entity
func advanceToDomain(a db.Advance) *entity.Advance {
	id := uuid.UUID(a.ID.Bytes).String()
	tenantID := pgtypeToEntityUUID(a.TenantID)
	professionalID := uuid.UUID(a.ProfessionalID.Bytes).String()

	var unitID *string
	if a.UnitID.Valid {
		uid := uuid.UUID(a.UnitID.Bytes).String()
		unitID = &uid
	}

	var approvedAt *time.Time
	if a.ApprovedAt.Valid {
		approvedAt = &a.ApprovedAt.Time
	}

	var approvedBy *string
	if a.ApprovedBy.Valid {
		uid := uuid.UUID(a.ApprovedBy.Bytes).String()
		approvedBy = &uid
	}

	var rejectedAt *time.Time
	if a.RejectedAt.Valid {
		rejectedAt = &a.RejectedAt.Time
	}

	var rejectedBy *string
	if a.RejectedBy.Valid {
		uid := uuid.UUID(a.RejectedBy.Bytes).String()
		rejectedBy = &uid
	}

	var deductedAt *time.Time
	if a.DeductedAt.Valid {
		deductedAt = &a.DeductedAt.Time
	}

	var deductionPeriodID *string
	if a.DeductionPeriodID.Valid {
		pid := uuid.UUID(a.DeductionPeriodID.Bytes).String()
		deductionPeriodID = &pid
	}

	var createdBy *string
	if a.CreatedBy.Valid {
		uid := uuid.UUID(a.CreatedBy.Bytes).String()
		createdBy = &uid
	}

	return &entity.Advance{
		ID:                id,
		TenantID:          tenantID,
		UnitID:            unitID,
		ProfessionalID:    professionalID,
		Amount:            a.Amount,
		RequestDate:       a.RequestDate.Time,
		Reason:            a.Reason,
		Status:            a.Status,
		ApprovedAt:        approvedAt,
		ApprovedBy:        approvedBy,
		RejectedAt:        rejectedAt,
		RejectedBy:        rejectedBy,
		RejectionReason:   a.RejectionReason,
		DeductedAt:        deductedAt,
		DeductionPeriodID: deductionPeriodID,
		CreatedAt:         a.CreatedAt.Time,
		UpdatedAt:         a.UpdatedAt.Time,
		CreatedBy:         createdBy,
	}
}

// advanceWithNamesToDomain converte de GetAdvanceByIDRow para entity
func advanceWithNamesToDomain(row db.GetAdvanceByIDRow) *entity.Advance {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := pgtypeToEntityUUID(row.TenantID)
	professionalID := uuid.UUID(row.ProfessionalID.Bytes).String()

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var approvedAt *time.Time
	if row.ApprovedAt.Valid {
		approvedAt = &row.ApprovedAt.Time
	}

	var approvedBy *string
	if row.ApprovedBy.Valid {
		uid := uuid.UUID(row.ApprovedBy.Bytes).String()
		approvedBy = &uid
	}

	var rejectedAt *time.Time
	if row.RejectedAt.Valid {
		rejectedAt = &row.RejectedAt.Time
	}

	var rejectedBy *string
	if row.RejectedBy.Valid {
		uid := uuid.UUID(row.RejectedBy.Bytes).String()
		rejectedBy = &uid
	}

	var deductedAt *time.Time
	if row.DeductedAt.Valid {
		deductedAt = &row.DeductedAt.Time
	}

	var deductionPeriodID *string
	if row.DeductionPeriodID.Valid {
		pid := uuid.UUID(row.DeductionPeriodID.Bytes).String()
		deductionPeriodID = &pid
	}

	var createdBy *string
	if row.CreatedBy.Valid {
		uid := uuid.UUID(row.CreatedBy.Bytes).String()
		createdBy = &uid
	}

	return &entity.Advance{
		ID:                id,
		TenantID:          tenantID,
		UnitID:            unitID,
		ProfessionalID:    professionalID,
		Amount:            row.Amount,
		RequestDate:       row.RequestDate.Time,
		Reason:            row.Reason,
		Status:            row.Status,
		ApprovedAt:        approvedAt,
		ApprovedBy:        approvedBy,
		RejectedAt:        rejectedAt,
		RejectedBy:        rejectedBy,
		RejectionReason:   row.RejectionReason,
		DeductedAt:        deductedAt,
		DeductionPeriodID: deductionPeriodID,
		CreatedAt:         row.CreatedAt.Time,
		UpdatedAt:         row.UpdatedAt.Time,
		CreatedBy:         createdBy,
	}
}

// advanceByStatusRowToDomain converte ListAdvancesByStatusRow para entity
func advanceByStatusRowToDomain(row db.ListAdvancesByStatusRow) *entity.Advance {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := pgtypeToEntityUUID(row.TenantID)
	professionalID := uuid.UUID(row.ProfessionalID.Bytes).String()

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var approvedAt *time.Time
	if row.ApprovedAt.Valid {
		approvedAt = &row.ApprovedAt.Time
	}

	var approvedBy *string
	if row.ApprovedBy.Valid {
		uid := uuid.UUID(row.ApprovedBy.Bytes).String()
		approvedBy = &uid
	}

	var rejectedAt *time.Time
	if row.RejectedAt.Valid {
		rejectedAt = &row.RejectedAt.Time
	}

	var rejectedBy *string
	if row.RejectedBy.Valid {
		uid := uuid.UUID(row.RejectedBy.Bytes).String()
		rejectedBy = &uid
	}

	var deductedAt *time.Time
	if row.DeductedAt.Valid {
		deductedAt = &row.DeductedAt.Time
	}

	var deductionPeriodID *string
	if row.DeductionPeriodID.Valid {
		pid := uuid.UUID(row.DeductionPeriodID.Bytes).String()
		deductionPeriodID = &pid
	}

	var createdBy *string
	if row.CreatedBy.Valid {
		uid := uuid.UUID(row.CreatedBy.Bytes).String()
		createdBy = &uid
	}

	return &entity.Advance{
		ID:                id,
		TenantID:          tenantID,
		UnitID:            unitID,
		ProfessionalID:    professionalID,
		Amount:            row.Amount,
		RequestDate:       row.RequestDate.Time,
		Reason:            row.Reason,
		Status:            row.Status,
		ApprovedAt:        approvedAt,
		ApprovedBy:        approvedBy,
		RejectedAt:        rejectedAt,
		RejectedBy:        rejectedBy,
		RejectionReason:   row.RejectionReason,
		DeductedAt:        deductedAt,
		DeductionPeriodID: deductionPeriodID,
		CreatedAt:         row.CreatedAt.Time,
		UpdatedAt:         row.UpdatedAt.Time,
		CreatedBy:         createdBy,
	}
}

// advanceByProfessionalRowToDomain converte ListAdvancesByProfessionalRow para entity
func advanceByProfessionalRowToDomain(row db.ListAdvancesByProfessionalRow) *entity.Advance {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := pgtypeToEntityUUID(row.TenantID)
	professionalID := uuid.UUID(row.ProfessionalID.Bytes).String()

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var approvedAt *time.Time
	if row.ApprovedAt.Valid {
		approvedAt = &row.ApprovedAt.Time
	}

	var approvedBy *string
	if row.ApprovedBy.Valid {
		uid := uuid.UUID(row.ApprovedBy.Bytes).String()
		approvedBy = &uid
	}

	var rejectedAt *time.Time
	if row.RejectedAt.Valid {
		rejectedAt = &row.RejectedAt.Time
	}

	var rejectedBy *string
	if row.RejectedBy.Valid {
		uid := uuid.UUID(row.RejectedBy.Bytes).String()
		rejectedBy = &uid
	}

	var deductedAt *time.Time
	if row.DeductedAt.Valid {
		deductedAt = &row.DeductedAt.Time
	}

	var deductionPeriodID *string
	if row.DeductionPeriodID.Valid {
		pid := uuid.UUID(row.DeductionPeriodID.Bytes).String()
		deductionPeriodID = &pid
	}

	var createdBy *string
	if row.CreatedBy.Valid {
		uid := uuid.UUID(row.CreatedBy.Bytes).String()
		createdBy = &uid
	}

	return &entity.Advance{
		ID:                id,
		TenantID:          tenantID,
		UnitID:            unitID,
		ProfessionalID:    professionalID,
		Amount:            row.Amount,
		RequestDate:       row.RequestDate.Time,
		Reason:            row.Reason,
		Status:            row.Status,
		ApprovedAt:        approvedAt,
		ApprovedBy:        approvedBy,
		RejectedAt:        rejectedAt,
		RejectedBy:        rejectedBy,
		RejectionReason:   row.RejectionReason,
		DeductedAt:        deductedAt,
		DeductionPeriodID: deductionPeriodID,
		CreatedAt:         row.CreatedAt.Time,
		UpdatedAt:         row.UpdatedAt.Time,
		CreatedBy:         createdBy,
	}
}

// advanceByTenantRowToDomain converte ListAdvancesByTenantRow para entity
func advanceByTenantRowToDomain(row db.ListAdvancesByTenantRow) *entity.Advance {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := pgtypeToEntityUUID(row.TenantID)
	professionalID := uuid.UUID(row.ProfessionalID.Bytes).String()

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var approvedAt *time.Time
	if row.ApprovedAt.Valid {
		approvedAt = &row.ApprovedAt.Time
	}

	var approvedBy *string
	if row.ApprovedBy.Valid {
		uid := uuid.UUID(row.ApprovedBy.Bytes).String()
		approvedBy = &uid
	}

	var rejectedAt *time.Time
	if row.RejectedAt.Valid {
		rejectedAt = &row.RejectedAt.Time
	}

	var rejectedBy *string
	if row.RejectedBy.Valid {
		uid := uuid.UUID(row.RejectedBy.Bytes).String()
		rejectedBy = &uid
	}

	var deductedAt *time.Time
	if row.DeductedAt.Valid {
		deductedAt = &row.DeductedAt.Time
	}

	var deductionPeriodID *string
	if row.DeductionPeriodID.Valid {
		pid := uuid.UUID(row.DeductionPeriodID.Bytes).String()
		deductionPeriodID = &pid
	}

	var createdBy *string
	if row.CreatedBy.Valid {
		uid := uuid.UUID(row.CreatedBy.Bytes).String()
		createdBy = &uid
	}

	return &entity.Advance{
		ID:                id,
		TenantID:          tenantID,
		UnitID:            unitID,
		ProfessionalID:    professionalID,
		Amount:            row.Amount,
		RequestDate:       row.RequestDate.Time,
		Reason:            row.Reason,
		Status:            row.Status,
		ApprovedAt:        approvedAt,
		ApprovedBy:        approvedBy,
		RejectedAt:        rejectedAt,
		RejectedBy:        rejectedBy,
		RejectionReason:   row.RejectionReason,
		DeductedAt:        deductedAt,
		DeductionPeriodID: deductionPeriodID,
		CreatedAt:         row.CreatedAt.Time,
		UpdatedAt:         row.UpdatedAt.Time,
		CreatedBy:         createdBy,
	}
}

// advancePendingRowToDomain converte ListPendingAdvancesRow para entity
func advancePendingRowToDomain(row db.ListPendingAdvancesRow) *entity.Advance {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := pgtypeToEntityUUID(row.TenantID)
	professionalID := uuid.UUID(row.ProfessionalID.Bytes).String()

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var approvedAt *time.Time
	if row.ApprovedAt.Valid {
		approvedAt = &row.ApprovedAt.Time
	}

	var approvedBy *string
	if row.ApprovedBy.Valid {
		uid := uuid.UUID(row.ApprovedBy.Bytes).String()
		approvedBy = &uid
	}

	var rejectedAt *time.Time
	if row.RejectedAt.Valid {
		rejectedAt = &row.RejectedAt.Time
	}

	var rejectedBy *string
	if row.RejectedBy.Valid {
		uid := uuid.UUID(row.RejectedBy.Bytes).String()
		rejectedBy = &uid
	}

	var deductedAt *time.Time
	if row.DeductedAt.Valid {
		deductedAt = &row.DeductedAt.Time
	}

	var deductionPeriodID *string
	if row.DeductionPeriodID.Valid {
		pid := uuid.UUID(row.DeductionPeriodID.Bytes).String()
		deductionPeriodID = &pid
	}

	var createdBy *string
	if row.CreatedBy.Valid {
		uid := uuid.UUID(row.CreatedBy.Bytes).String()
		createdBy = &uid
	}

	return &entity.Advance{
		ID:                id,
		TenantID:          tenantID,
		UnitID:            unitID,
		ProfessionalID:    professionalID,
		Amount:            row.Amount,
		RequestDate:       row.RequestDate.Time,
		Reason:            row.Reason,
		Status:            row.Status,
		ApprovedAt:        approvedAt,
		ApprovedBy:        approvedBy,
		RejectedAt:        rejectedAt,
		RejectedBy:        rejectedBy,
		RejectionReason:   row.RejectionReason,
		DeductedAt:        deductedAt,
		DeductionPeriodID: deductionPeriodID,
		CreatedAt:         row.CreatedAt.Time,
		UpdatedAt:         row.UpdatedAt.Time,
		CreatedBy:         createdBy,
	}
}
