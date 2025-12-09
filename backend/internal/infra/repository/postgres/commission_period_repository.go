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

// commissionPeriodRepository implementa o repository de períodos de comissão usando PostgreSQL
type commissionPeriodRepository struct {
	queries *db.Queries
}

// NewCommissionPeriodRepository cria uma nova instância do repository
func NewCommissionPeriodRepository(queries *db.Queries) repository.CommissionPeriodRepository {
	return &commissionPeriodRepository{
		queries: queries,
	}
}

// Create cria um novo período de comissão
func (r *commissionPeriodRepository) Create(ctx context.Context, period *entity.CommissionPeriod) (*entity.CommissionPeriod, error) {
	var unitID pgtype.UUID
	if period.UnitID != nil {
		uid, err := uuid.Parse(*period.UnitID)
		if err != nil {
			return nil, err
		}
		unitID = pgtype.UUID{Bytes: uid, Valid: true}
	}

	var professionalID pgtype.UUID
	if period.ProfessionalID != nil {
		pid, err := uuid.Parse(*period.ProfessionalID)
		if err != nil {
			return nil, err
		}
		professionalID = pgtype.UUID{Bytes: pid, Valid: true}
	}

	result, err := r.queries.CreateCommissionPeriod(ctx, db.CreateCommissionPeriodParams{
		TenantID:         entityUUIDToPgtype(period.TenantID),
		UnitID:           unitID,
		ReferenceMonth:   period.ReferenceMonth,
		ProfessionalID:   professionalID,
		TotalGross:       period.TotalGross,
		TotalCommission:  period.TotalCommission,
		TotalAdvances:    period.TotalAdvances,
		TotalAdjustments: period.TotalAdjustments,
		TotalNet:         period.TotalNet,
		ItemsCount:       int32(period.ItemsCount),
		Status:           period.Status,
		PeriodStart:      pgtype.Date{Time: period.PeriodStart, Valid: true},
		PeriodEnd:        pgtype.Date{Time: period.PeriodEnd, Valid: true},
		Notes:            period.Notes,
	})
	if err != nil {
		return nil, err
	}

	return commissionPeriodToDomain(result), nil
}

// GetByID busca um período de comissão por ID
func (r *commissionPeriodRepository) GetByID(ctx context.Context, tenantID, id string) (*entity.CommissionPeriod, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	pid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.GetCommissionPeriodByID(ctx, db.GetCommissionPeriodByIDParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		ID:       pgtype.UUID{Bytes: pid, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("período de comissão não encontrado")
		}
		return nil, err
	}

	return commissionPeriodToDomain(result), nil
}

// List lista períodos de comissão com filtros
func (r *commissionPeriodRepository) List(ctx context.Context, tenantID string, professionalID *string, status *string, limit, offset int) ([]*entity.CommissionPeriod, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	// Se status for fornecido, usa query por status
	if status != nil && *status != "" {
		results, err := r.queries.ListCommissionPeriodsByStatus(ctx, db.ListCommissionPeriodsByStatusParams{
			TenantID: pgtype.UUID{Bytes: tid, Valid: true},
			Status:   *status,
			Limit:    int32(limit),
			Offset:   int32(offset),
		})
		if err != nil {
			return nil, err
		}

		periods := make([]*entity.CommissionPeriod, 0, len(results))
		for _, result := range results {
			periods = append(periods, commissionPeriodRowToDomain(result))
		}
		return periods, nil
	}

	// Se professionalID for fornecido, usa query por profissional
	if professionalID != nil && *professionalID != "" {
		pid, err := uuid.Parse(*professionalID)
		if err != nil {
			return nil, err
		}

		results, err := r.queries.ListCommissionPeriodsByProfessional(ctx, db.ListCommissionPeriodsByProfessionalParams{
			TenantID:       pgtype.UUID{Bytes: tid, Valid: true},
			ProfessionalID: pgtype.UUID{Bytes: pid, Valid: true},
			Limit:          int32(limit),
			Offset:         int32(offset),
		})
		if err != nil {
			return nil, err
		}

		periods := make([]*entity.CommissionPeriod, 0, len(results))
		for _, result := range results {
			periods = append(periods, commissionPeriodByProfessionalRowToDomain(result))
		}
		return periods, nil
	}

	// Query padrão por tenant
	results, err := r.queries.ListCommissionPeriodsByTenant(ctx, db.ListCommissionPeriodsByTenantParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}

	periods := make([]*entity.CommissionPeriod, 0, len(results))
	for _, result := range results {
		periods = append(periods, commissionPeriodByTenantRowToDomain(result))
	}
	return periods, nil
}

// GetByProfessional busca períodos de comissão por profissional
func (r *commissionPeriodRepository) GetByProfessional(ctx context.Context, tenantID, professionalID string) ([]*entity.CommissionPeriod, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	pid, err := uuid.Parse(professionalID)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.ListCommissionPeriodsByProfessional(ctx, db.ListCommissionPeriodsByProfessionalParams{
		TenantID:       pgtype.UUID{Bytes: tid, Valid: true},
		ProfessionalID: pgtype.UUID{Bytes: pid, Valid: true},
		Limit:          1000,
		Offset:         0,
	})
	if err != nil {
		return nil, err
	}

	periods := make([]*entity.CommissionPeriod, 0, len(results))
	for _, result := range results {
		periods = append(periods, commissionPeriodByProfessionalRowToDomain(result))
	}
	return periods, nil
}

// GetOpenByProfessional busca o período aberto de um profissional
func (r *commissionPeriodRepository) GetOpenByProfessional(ctx context.Context, tenantID, professionalID string) (*entity.CommissionPeriod, error) {
	// Busca períodos abertos do tenant
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.ListOpenCommissionPeriods(ctx, pgtype.UUID{Bytes: tid, Valid: true})
	if err != nil {
		return nil, err
	}

	// Filtra pelo profissional
	for _, result := range results {
		if result.ProfessionalID.Valid {
			pid := uuid.UUID(result.ProfessionalID.Bytes).String()
			if pid == professionalID {
				return commissionPeriodOpenRowToDomain(result), nil
			}
		}
	}

	return nil, nil // Nenhum período aberto encontrado
}

// GetByDateRange busca períodos de comissão por intervalo de datas
func (r *commissionPeriodRepository) GetByDateRange(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*entity.CommissionPeriod, error) {
	// Usa List com filtro padrão
	return r.List(ctx, tenantID, nil, nil, 1000, 0)
}

// GetSummary retorna totais do período
func (r *commissionPeriodRepository) GetSummary(ctx context.Context, tenantID, periodID string) (*entity.CommissionPeriodSummary, error) {
	// Busca o período para obter os totais
	period, err := r.GetByID(ctx, tenantID, periodID)
	if err != nil {
		return nil, err
	}

	return &entity.CommissionPeriodSummary{
		TotalGross:       period.TotalGross,
		TotalCommission:  period.TotalCommission,
		TotalAdvances:    period.TotalAdvances,
		TotalAdjustments: period.TotalAdjustments,
		TotalNet:         period.TotalNet,
		ItemsCount:       period.ItemsCount,
	}, nil
}

// Update atualiza um período de comissão
func (r *commissionPeriodRepository) Update(ctx context.Context, period *entity.CommissionPeriod) (*entity.CommissionPeriod, error) {
	id, err := uuid.Parse(period.ID)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.UpdateCommissionPeriodTotals(ctx, db.UpdateCommissionPeriodTotalsParams{
		ID:               pgtype.UUID{Bytes: id, Valid: true},
		TenantID:         entityUUIDToPgtype(period.TenantID),
		TotalGross:       period.TotalGross,
		TotalCommission:  period.TotalCommission,
		TotalAdvances:    period.TotalAdvances,
		TotalAdjustments: period.TotalAdjustments,
		TotalNet:         period.TotalNet,
		ItemsCount:       int32(period.ItemsCount),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("período de comissão não encontrado")
		}
		return nil, err
	}

	return commissionPeriodToDomain(result), nil
}

// Close fecha um período de comissão
func (r *commissionPeriodRepository) Close(ctx context.Context, tenantID, id, closedBy string) (*entity.CommissionPeriod, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	pid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	closedByUUID, err := uuid.Parse(closedBy)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.CloseCommissionPeriod(ctx, db.CloseCommissionPeriodParams{
		ID:       pgtype.UUID{Bytes: pid, Valid: true},
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		ClosedBy: pgtype.UUID{Bytes: closedByUUID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("período de comissão não encontrado")
		}
		return nil, err
	}

	return commissionPeriodToDomain(result), nil
}

// MarkAsPaid marca um período como pago
func (r *commissionPeriodRepository) MarkAsPaid(ctx context.Context, tenantID, id, paidBy string) (*entity.CommissionPeriod, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	pid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	paidByUUID, err := uuid.Parse(paidBy)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.MarkCommissionPeriodAsPaid(ctx, db.MarkCommissionPeriodAsPaidParams{
		ID:       pgtype.UUID{Bytes: pid, Valid: true},
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		PaidBy:   pgtype.UUID{Bytes: paidByUUID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("período de comissão não encontrado")
		}
		return nil, err
	}

	return commissionPeriodToDomain(result), nil
}

// Delete remove um período de comissão (somente se ABERTO)
func (r *commissionPeriodRepository) Delete(ctx context.Context, tenantID, id string) error {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return err
	}

	pid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteCommissionPeriod(ctx, db.DeleteCommissionPeriodParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		ID:       pgtype.UUID{Bytes: pid, Valid: true},
	})
}

// === Funções auxiliares de conversão ===

// commissionPeriodToDomain converte de db.CommissionPeriod para entity
func commissionPeriodToDomain(cp db.CommissionPeriod) *entity.CommissionPeriod {
	id := uuid.UUID(cp.ID.Bytes).String()
	tenantID := pgtypeToEntityUUID(cp.TenantID)

	var unitID *string
	if cp.UnitID.Valid {
		uid := uuid.UUID(cp.UnitID.Bytes).String()
		unitID = &uid
	}

	var professionalID *string
	if cp.ProfessionalID.Valid {
		pid := uuid.UUID(cp.ProfessionalID.Bytes).String()
		professionalID = &pid
	}

	var closedAt *time.Time
	if cp.ClosedAt.Valid {
		closedAt = &cp.ClosedAt.Time
	}

	var paidAt *time.Time
	if cp.PaidAt.Valid {
		paidAt = &cp.PaidAt.Time
	}

	var contaPagarID *string
	if cp.ContaPagarID.Valid {
		cid := uuid.UUID(cp.ContaPagarID.Bytes).String()
		contaPagarID = &cid
	}

	var closedBy *string
	if cp.ClosedBy.Valid {
		uid := uuid.UUID(cp.ClosedBy.Bytes).String()
		closedBy = &uid
	}

	var paidBy *string
	if cp.PaidBy.Valid {
		uid := uuid.UUID(cp.PaidBy.Bytes).String()
		paidBy = &uid
	}

	return &entity.CommissionPeriod{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ReferenceMonth:   cp.ReferenceMonth,
		ProfessionalID:   professionalID,
		TotalGross:       cp.TotalGross,
		TotalCommission:  cp.TotalCommission,
		TotalAdvances:    cp.TotalAdvances,
		TotalAdjustments: cp.TotalAdjustments,
		TotalNet:         cp.TotalNet,
		ItemsCount:       int(cp.ItemsCount),
		Status:           cp.Status,
		PeriodStart:      cp.PeriodStart.Time,
		PeriodEnd:        cp.PeriodEnd.Time,
		ClosedAt:         closedAt,
		PaidAt:           paidAt,
		ContaPagarID:     contaPagarID,
		ClosedBy:         closedBy,
		PaidBy:           paidBy,
		Notes:            cp.Notes,
		CreatedAt:        cp.CreatedAt.Time,
		UpdatedAt:        cp.UpdatedAt.Time,
	}
}

// commissionPeriodRowToDomain converte ListCommissionPeriodsByStatusRow para entity
func commissionPeriodRowToDomain(row db.ListCommissionPeriodsByStatusRow) *entity.CommissionPeriod {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := pgtypeToEntityUUID(row.TenantID)

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var professionalID *string
	if row.ProfessionalID.Valid {
		pid := uuid.UUID(row.ProfessionalID.Bytes).String()
		professionalID = &pid
	}

	var closedAt *time.Time
	if row.ClosedAt.Valid {
		closedAt = &row.ClosedAt.Time
	}

	var paidAt *time.Time
	if row.PaidAt.Valid {
		paidAt = &row.PaidAt.Time
	}

	var contaPagarID *string
	if row.ContaPagarID.Valid {
		cid := uuid.UUID(row.ContaPagarID.Bytes).String()
		contaPagarID = &cid
	}

	var closedBy *string
	if row.ClosedBy.Valid {
		uid := uuid.UUID(row.ClosedBy.Bytes).String()
		closedBy = &uid
	}

	var paidBy *string
	if row.PaidBy.Valid {
		uid := uuid.UUID(row.PaidBy.Bytes).String()
		paidBy = &uid
	}

	return &entity.CommissionPeriod{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ReferenceMonth:   row.ReferenceMonth,
		ProfessionalID:   professionalID,
		TotalGross:       row.TotalGross,
		TotalCommission:  row.TotalCommission,
		TotalAdvances:    row.TotalAdvances,
		TotalAdjustments: row.TotalAdjustments,
		TotalNet:         row.TotalNet,
		ItemsCount:       int(row.ItemsCount),
		Status:           row.Status,
		PeriodStart:      row.PeriodStart.Time,
		PeriodEnd:        row.PeriodEnd.Time,
		ClosedAt:         closedAt,
		PaidAt:           paidAt,
		ContaPagarID:     contaPagarID,
		ClosedBy:         closedBy,
		PaidBy:           paidBy,
		Notes:            row.Notes,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}
}

// commissionPeriodByProfessionalRowToDomain converte ListCommissionPeriodsByProfessionalRow para entity
func commissionPeriodByProfessionalRowToDomain(row db.ListCommissionPeriodsByProfessionalRow) *entity.CommissionPeriod {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := pgtypeToEntityUUID(row.TenantID)

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var professionalID *string
	if row.ProfessionalID.Valid {
		pid := uuid.UUID(row.ProfessionalID.Bytes).String()
		professionalID = &pid
	}

	var closedAt *time.Time
	if row.ClosedAt.Valid {
		closedAt = &row.ClosedAt.Time
	}

	var paidAt *time.Time
	if row.PaidAt.Valid {
		paidAt = &row.PaidAt.Time
	}

	var contaPagarID *string
	if row.ContaPagarID.Valid {
		cid := uuid.UUID(row.ContaPagarID.Bytes).String()
		contaPagarID = &cid
	}

	var closedBy *string
	if row.ClosedBy.Valid {
		uid := uuid.UUID(row.ClosedBy.Bytes).String()
		closedBy = &uid
	}

	var paidBy *string
	if row.PaidBy.Valid {
		uid := uuid.UUID(row.PaidBy.Bytes).String()
		paidBy = &uid
	}

	return &entity.CommissionPeriod{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ReferenceMonth:   row.ReferenceMonth,
		ProfessionalID:   professionalID,
		TotalGross:       row.TotalGross,
		TotalCommission:  row.TotalCommission,
		TotalAdvances:    row.TotalAdvances,
		TotalAdjustments: row.TotalAdjustments,
		TotalNet:         row.TotalNet,
		ItemsCount:       int(row.ItemsCount),
		Status:           row.Status,
		PeriodStart:      row.PeriodStart.Time,
		PeriodEnd:        row.PeriodEnd.Time,
		ClosedAt:         closedAt,
		PaidAt:           paidAt,
		ContaPagarID:     contaPagarID,
		ClosedBy:         closedBy,
		PaidBy:           paidBy,
		Notes:            row.Notes,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}
}

// commissionPeriodByTenantRowToDomain converte ListCommissionPeriodsByTenantRow para entity
func commissionPeriodByTenantRowToDomain(row db.ListCommissionPeriodsByTenantRow) *entity.CommissionPeriod {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := pgtypeToEntityUUID(row.TenantID)

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var professionalID *string
	if row.ProfessionalID.Valid {
		pid := uuid.UUID(row.ProfessionalID.Bytes).String()
		professionalID = &pid
	}

	var closedAt *time.Time
	if row.ClosedAt.Valid {
		closedAt = &row.ClosedAt.Time
	}

	var paidAt *time.Time
	if row.PaidAt.Valid {
		paidAt = &row.PaidAt.Time
	}

	var contaPagarID *string
	if row.ContaPagarID.Valid {
		cid := uuid.UUID(row.ContaPagarID.Bytes).String()
		contaPagarID = &cid
	}

	var closedBy *string
	if row.ClosedBy.Valid {
		uid := uuid.UUID(row.ClosedBy.Bytes).String()
		closedBy = &uid
	}

	var paidBy *string
	if row.PaidBy.Valid {
		uid := uuid.UUID(row.PaidBy.Bytes).String()
		paidBy = &uid
	}

	return &entity.CommissionPeriod{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ReferenceMonth:   row.ReferenceMonth,
		ProfessionalID:   professionalID,
		TotalGross:       row.TotalGross,
		TotalCommission:  row.TotalCommission,
		TotalAdvances:    row.TotalAdvances,
		TotalAdjustments: row.TotalAdjustments,
		TotalNet:         row.TotalNet,
		ItemsCount:       int(row.ItemsCount),
		Status:           row.Status,
		PeriodStart:      row.PeriodStart.Time,
		PeriodEnd:        row.PeriodEnd.Time,
		ClosedAt:         closedAt,
		PaidAt:           paidAt,
		ContaPagarID:     contaPagarID,
		ClosedBy:         closedBy,
		PaidBy:           paidBy,
		Notes:            row.Notes,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}
}

// commissionPeriodOpenRowToDomain converte ListOpenCommissionPeriodsRow para entity
func commissionPeriodOpenRowToDomain(row db.ListOpenCommissionPeriodsRow) *entity.CommissionPeriod {
	id := uuid.UUID(row.ID.Bytes).String()
	tenantID := pgtypeToEntityUUID(row.TenantID)

	var unitID *string
	if row.UnitID.Valid {
		uid := uuid.UUID(row.UnitID.Bytes).String()
		unitID = &uid
	}

	var professionalID *string
	if row.ProfessionalID.Valid {
		pid := uuid.UUID(row.ProfessionalID.Bytes).String()
		professionalID = &pid
	}

	var closedAt *time.Time
	if row.ClosedAt.Valid {
		closedAt = &row.ClosedAt.Time
	}

	var paidAt *time.Time
	if row.PaidAt.Valid {
		paidAt = &row.PaidAt.Time
	}

	var contaPagarID *string
	if row.ContaPagarID.Valid {
		cid := uuid.UUID(row.ContaPagarID.Bytes).String()
		contaPagarID = &cid
	}

	var closedBy *string
	if row.ClosedBy.Valid {
		uid := uuid.UUID(row.ClosedBy.Bytes).String()
		closedBy = &uid
	}

	var paidBy *string
	if row.PaidBy.Valid {
		uid := uuid.UUID(row.PaidBy.Bytes).String()
		paidBy = &uid
	}

	return &entity.CommissionPeriod{
		ID:               id,
		TenantID:         tenantID,
		UnitID:           unitID,
		ReferenceMonth:   row.ReferenceMonth,
		ProfessionalID:   professionalID,
		TotalGross:       row.TotalGross,
		TotalCommission:  row.TotalCommission,
		TotalAdvances:    row.TotalAdvances,
		TotalAdjustments: row.TotalAdjustments,
		TotalNet:         row.TotalNet,
		ItemsCount:       int(row.ItemsCount),
		Status:           row.Status,
		PeriodStart:      row.PeriodStart.Time,
		PeriodEnd:        row.PeriodEnd.Time,
		ClosedAt:         closedAt,
		PaidAt:           paidAt,
		ContaPagarID:     contaPagarID,
		ClosedBy:         closedBy,
		PaidBy:           paidBy,
		Notes:            row.Notes,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}
}
