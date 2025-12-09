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

// blockedTimeRepository implementa o repository de blocked times usando PostgreSQL
type blockedTimeRepository struct {
	queries *db.Queries
}

// NewBlockedTimeRepository cria uma nova instância do repository
func NewBlockedTimeRepository(queries *db.Queries) repository.BlockedTimeRepository {
	return &blockedTimeRepository{
		queries: queries,
	}
}

// Create cria um novo bloqueio
func (r *blockedTimeRepository) Create(ctx context.Context, blockedTime *entity.BlockedTime) (*entity.BlockedTime, error) {
	id, err := uuid.Parse(blockedTime.ID)
	if err != nil {
		return nil, err
	}

	professionalID, err := uuid.Parse(blockedTime.ProfessionalID)
	if err != nil {
		return nil, err
	}

	var createdBy pgtype.UUID
	if blockedTime.CreatedBy != nil {
		userID, err := uuid.Parse(*blockedTime.CreatedBy)
		if err != nil {
			return nil, err
		}
		createdBy = pgtype.UUID{Bytes: userID, Valid: true}
	}

	result, err := r.queries.CreateBlockedTime(ctx, db.CreateBlockedTimeParams{
		ID:             pgtype.UUID{Bytes: id, Valid: true},
		TenantID:       entityUUIDToPgtype(blockedTime.TenantID),
		ProfessionalID: pgtype.UUID{Bytes: professionalID, Valid: true},
		StartTime:      pgtype.Timestamptz{Time: blockedTime.StartTime, Valid: true},
		EndTime:        pgtype.Timestamptz{Time: blockedTime.EndTime, Valid: true},
		Reason:         blockedTime.Reason,
		IsRecurring:    blockedTime.IsRecurring,
		RecurrenceRule: blockedTime.RecurrenceRule,
		CreatedBy:      createdBy,
	})
	if err != nil {
		return nil, err
	}

	return toDomain(result), nil
}

// GetByID busca um bloqueio por ID
func (r *blockedTimeRepository) GetByID(ctx context.Context, tenantID, id string) (*entity.BlockedTime, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	bid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.GetBlockedTimeByID(ctx, db.GetBlockedTimeByIDParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		ID:       pgtype.UUID{Bytes: bid, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("bloqueio não encontrado")
		}
		return nil, err
	}

	return toDomain(result), nil
}

// List lista bloqueios com filtros opcionais
func (r *blockedTimeRepository) List(ctx context.Context, tenantID string, professionalID *string, startDate, endDate *time.Time) ([]*entity.BlockedTime, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	params := db.ListBlockedTimesParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
	}

	if professionalID != nil {
		pid, err := uuid.Parse(*professionalID)
		if err != nil {
			return nil, err
		}
		params.ProfessionalID = pgtype.UUID{Bytes: pid, Valid: true}
	}

	if startDate != nil {
		params.StartDate = pgtype.Timestamptz{Time: *startDate, Valid: true}
	}

	if endDate != nil {
		params.EndDate = pgtype.Timestamptz{Time: *endDate, Valid: true}
	}

	results, err := r.queries.ListBlockedTimes(ctx, params)
	if err != nil {
		return nil, err
	}

	blockedTimes := make([]*entity.BlockedTime, 0, len(results))
	for _, result := range results {
		blockedTimes = append(blockedTimes, toDomain(result))
	}

	return blockedTimes, nil
}

// GetInRange busca bloqueios em um intervalo
func (r *blockedTimeRepository) GetInRange(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time) ([]*entity.BlockedTime, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}

	pid, err := uuid.Parse(professionalID)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.GetBlockedTimesInRange(ctx, db.GetBlockedTimesInRangeParams{
		TenantID:       pgtype.UUID{Bytes: tid, Valid: true},
		ProfessionalID: pgtype.UUID{Bytes: pid, Valid: true},
		StartTime:      pgtype.Timestamptz{Time: startTime, Valid: true},
		EndTime:        pgtype.Timestamptz{Time: endTime, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	blockedTimes := make([]*entity.BlockedTime, 0, len(results))
	for _, result := range results {
		blockedTimes = append(blockedTimes, toDomain(result))
	}

	return blockedTimes, nil
}

// CheckConflict verifica conflito
func (r *blockedTimeRepository) CheckConflict(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time, excludeID *string) (bool, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return false, err
	}

	pid, err := uuid.Parse(professionalID)
	if err != nil {
		return false, err
	}

	params := db.CheckBlockedTimeConflictParams{
		TenantID:       pgtype.UUID{Bytes: tid, Valid: true},
		ProfessionalID: pgtype.UUID{Bytes: pid, Valid: true},
		StartTime:      pgtype.Timestamptz{Time: startTime, Valid: true},
		EndTime:        pgtype.Timestamptz{Time: endTime, Valid: true},
	}

	if excludeID != nil {
		eid, err := uuid.Parse(*excludeID)
		if err != nil {
			return false, err
		}
		params.ExcludeID = pgtype.UUID{Bytes: eid, Valid: true}
	}

	return r.queries.CheckBlockedTimeConflict(ctx, params)
}

// Update atualiza um bloqueio
func (r *blockedTimeRepository) Update(ctx context.Context, blockedTime *entity.BlockedTime) (*entity.BlockedTime, error) {
	// Por enquanto, vamos apenas retornar o bloqueio
	// Implementação completa futura
	return blockedTime, nil
}

// Delete remove um bloqueio
func (r *blockedTimeRepository) Delete(ctx context.Context, tenantID, id string) error {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return err
	}

	bid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteBlockedTime(ctx, db.DeleteBlockedTimeParams{
		TenantID: pgtype.UUID{Bytes: tid, Valid: true},
		ID:       pgtype.UUID{Bytes: bid, Valid: true},
	})
}

// toDomain converte de sqlc model para entity domain
func toDomain(bt db.BlockedTime) *entity.BlockedTime {
	var createdBy *string
	if bt.CreatedBy.Valid {
		uid, _ := uuid.FromBytes(bt.CreatedBy.Bytes[:])
		uidStr := uid.String()
		createdBy = &uidStr
	}

	var recurrenceRule *string
	if bt.RecurrenceRule != nil {
		recurrenceRule = bt.RecurrenceRule
	}

	id, _ := uuid.FromBytes(bt.ID.Bytes[:])
	tenantID := pgtypeToEntityUUID(bt.TenantID)
	professionalID, _ := uuid.FromBytes(bt.ProfessionalID.Bytes[:])

	return &entity.BlockedTime{
		ID:             id.String(),
		TenantID:       tenantID,
		ProfessionalID: professionalID.String(),
		StartTime:      bt.StartTime.Time,
		EndTime:        bt.EndTime.Time,
		Reason:         bt.Reason,
		IsRecurring:    bt.IsRecurring,
		RecurrenceRule: recurrenceRule,
		CreatedAt:      bt.CreatedAt.Time,
		UpdatedAt:      bt.UpdatedAt.Time,
		CreatedBy:      createdBy,
	}
}
