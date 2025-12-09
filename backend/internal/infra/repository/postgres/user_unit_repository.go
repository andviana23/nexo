package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// UserUnitRepositoryPG implementa port.UserUnitRepository usando PostgreSQL
type UserUnitRepositoryPG struct {
	queries *db.Queries
}

// NewUserUnitRepository cria nova instância do repositório
func NewUserUnitRepository(queries *db.Queries) port.UserUnitRepository {
	return &UserUnitRepositoryPG{queries: queries}
}

// Create cria um novo vínculo
func (r *UserUnitRepositoryPG) Create(ctx context.Context, userUnit *entity.UserUnit) error {
	params := db.CreateUserUnitParams{
		UserID:       uuidToPgUUID(userUnit.UserID),
		UnitID:       uuidToPgUUID(userUnit.UnitID),
		IsDefault:    userUnit.IsDefault,
		RoleOverride: userUnit.RoleOverride,
	}

	result, err := r.queries.CreateUserUnit(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar vínculo usuário-unidade: %w", err)
	}

	userUnit.ID = pgUUIDToUUID(result.ID)
	userUnit.CriadoEm = result.CriadoEm.Time
	userUnit.AtualizadoEm = result.AtualizadoEm.Time

	return nil
}

// FindByUserAndUnit busca vínculo específico
func (r *UserUnitRepositoryPG) FindByUserAndUnit(ctx context.Context, userID, unitID uuid.UUID) (*entity.UserUnit, error) {
	result, err := r.queries.GetUserUnit(ctx, db.GetUserUnitParams{
		UserID: uuidToPgUUID(userID),
		UnitID: uuidToPgUUID(unitID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar vínculo: %w", err)
	}

	return r.toDomain(result), nil
}

// FindUserDefaultUnit busca a unidade padrão do usuário
func (r *UserUnitRepositoryPG) FindUserDefaultUnit(ctx context.Context, userID uuid.UUID) (*entity.UserUnitWithDetails, error) {
	result, err := r.queries.GetUserDefaultUnit(ctx, uuidToPgUUID(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar unidade padrão: %w", err)
	}

	return r.toDetailsDomain(result), nil
}

// ListByUser lista unidades do usuário
func (r *UserUnitRepositoryPG) ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.UserUnitWithDetails, error) {
	results, err := r.queries.ListUserUnits(ctx, uuidToPgUUID(userID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar unidades do usuário: %w", err)
	}

	units := make([]*entity.UserUnitWithDetails, len(results))
	for i, result := range results {
		units[i] = r.listToDetailsDomain(result)
	}

	return units, nil
}

// ListByUnit lista usuários da unidade
func (r *UserUnitRepositoryPG) ListByUnit(ctx context.Context, unitID uuid.UUID) ([]*entity.UserUnit, error) {
	results, err := r.queries.ListUnitUsers(ctx, uuidToPgUUID(unitID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar usuários da unidade: %w", err)
	}

	userUnits := make([]*entity.UserUnit, len(results))
	for i, result := range results {
		userUnits[i] = &entity.UserUnit{
			ID:           pgUUIDToUUID(result.ID),
			UserID:       pgUUIDToUUID(result.UserID),
			UnitID:       pgUUIDToUUID(result.UnitID),
			IsDefault:    result.IsDefault,
			RoleOverride: result.RoleOverride,
			CriadoEm:     result.CriadoEm.Time,
			AtualizadoEm: result.AtualizadoEm.Time,
		}
	}

	return userUnits, nil
}

// CheckAccess verifica se usuário tem acesso à unidade
func (r *UserUnitRepositoryPG) CheckAccess(ctx context.Context, userID, unitID uuid.UUID) (bool, error) {
	result, err := r.queries.CheckUserUnitAccess(ctx, db.CheckUserUnitAccessParams{
		UserID: uuidToPgUUID(userID),
		UnitID: uuidToPgUUID(unitID),
	})
	if err != nil {
		return false, fmt.Errorf("erro ao verificar acesso: %w", err)
	}

	return result, nil
}

// SetDefault define unidade padrão do usuário
func (r *UserUnitRepositoryPG) SetDefault(ctx context.Context, userID, unitID uuid.UUID) error {
	err := r.queries.SetUserDefaultUnit(ctx, db.SetUserDefaultUnitParams{
		UserID: uuidToPgUUID(userID),
		UnitID: uuidToPgUUID(unitID),
	})
	if err != nil {
		return fmt.Errorf("erro ao definir unidade padrão: %w", err)
	}
	return nil
}

// UpdateRole atualiza o papel do usuário na unidade
func (r *UserUnitRepositoryPG) UpdateRole(ctx context.Context, userID, unitID uuid.UUID, role *string) error {
	_, err := r.queries.UpdateUserUnitRole(ctx, db.UpdateUserUnitRoleParams{
		UserID:       uuidToPgUUID(userID),
		UnitID:       uuidToPgUUID(unitID),
		RoleOverride: role,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ErrUserUnitNaoEncontrado
		}
		return fmt.Errorf("erro ao atualizar papel: %w", err)
	}
	return nil
}

// Delete remove vínculo
func (r *UserUnitRepositoryPG) Delete(ctx context.Context, userID, unitID uuid.UUID) error {
	err := r.queries.DeleteUserUnit(ctx, db.DeleteUserUnitParams{
		UserID: uuidToPgUUID(userID),
		UnitID: uuidToPgUUID(unitID),
	})
	if err != nil {
		return fmt.Errorf("erro ao remover vínculo: %w", err)
	}
	return nil
}

// DeleteAllByUser remove todos os vínculos de um usuário
func (r *UserUnitRepositoryPG) DeleteAllByUser(ctx context.Context, userID uuid.UUID) error {
	err := r.queries.DeleteAllUserUnits(ctx, uuidToPgUUID(userID))
	if err != nil {
		return fmt.Errorf("erro ao remover vínculos do usuário: %w", err)
	}
	return nil
}

// DeleteAllByUnit remove todos os vínculos de uma unidade
func (r *UserUnitRepositoryPG) DeleteAllByUnit(ctx context.Context, unitID uuid.UUID) error {
	err := r.queries.DeleteAllUnitUsers(ctx, uuidToPgUUID(unitID))
	if err != nil {
		return fmt.Errorf("erro ao remover vínculos da unidade: %w", err)
	}
	return nil
}

// CountByUser conta unidades do usuário
func (r *UserUnitRepositoryPG) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := r.queries.CountUserUnits(ctx, uuidToPgUUID(userID))
	if err != nil {
		return 0, fmt.Errorf("erro ao contar unidades: %w", err)
	}
	return count, nil
}

// CountByUnit conta usuários da unidade
func (r *UserUnitRepositoryPG) CountByUnit(ctx context.Context, unitID uuid.UUID) (int64, error) {
	count, err := r.queries.CountUnitUsers(ctx, uuidToPgUUID(unitID))
	if err != nil {
		return 0, fmt.Errorf("erro ao contar usuários: %w", err)
	}
	return count, nil
}

// toDomain converte modelo sqlc para entidade de domínio
func (r *UserUnitRepositoryPG) toDomain(uu db.UserUnit) *entity.UserUnit {
	return &entity.UserUnit{
		ID:           pgUUIDToUUID(uu.ID),
		UserID:       pgUUIDToUUID(uu.UserID),
		UnitID:       pgUUIDToUUID(uu.UnitID),
		IsDefault:    uu.IsDefault,
		RoleOverride: uu.RoleOverride,
		CriadoEm:     uu.CriadoEm.Time,
		AtualizadoEm: uu.AtualizadoEm.Time,
	}
}

// toDetailsDomain converte resultado de GetUserDefaultUnit
func (r *UserUnitRepositoryPG) toDetailsDomain(row db.GetUserDefaultUnitRow) *entity.UserUnitWithDetails {
	return &entity.UserUnitWithDetails{
		UserUnit: entity.UserUnit{
			ID:           pgUUIDToUUID(row.ID),
			UserID:       pgUUIDToUUID(row.UserID),
			UnitID:       pgUUIDToUUID(row.UnitID),
			IsDefault:    row.IsDefault,
			RoleOverride: row.RoleOverride,
			CriadoEm:     row.CriadoEm.Time,
			AtualizadoEm: row.AtualizadoEm.Time,
		},
		UnitNome:    row.UnitNome,
		UnitApelido: row.UnitApelido,
		TenantID:    pgUUIDToUUID(row.TenantID),
	}
}

// listToDetailsDomain converte resultado de ListUserUnits
func (r *UserUnitRepositoryPG) listToDetailsDomain(row db.ListUserUnitsRow) *entity.UserUnitWithDetails {
	return &entity.UserUnitWithDetails{
		UserUnit: entity.UserUnit{
			ID:           pgUUIDToUUID(row.ID),
			UserID:       pgUUIDToUUID(row.UserID),
			UnitID:       pgUUIDToUUID(row.UnitID),
			IsDefault:    row.IsDefault,
			RoleOverride: row.RoleOverride,
			CriadoEm:     row.CriadoEm.Time,
			AtualizadoEm: row.AtualizadoEm.Time,
		},
		UnitNome:    row.UnitNome,
		UnitApelido: row.UnitApelido,
		UnitMatriz:  row.IsMatriz,
		UnitAtiva:   row.UnitAtiva,
		TenantID:    pgUUIDToUUID(row.TenantID),
	}
}
