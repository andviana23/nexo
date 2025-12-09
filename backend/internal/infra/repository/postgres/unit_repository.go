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

// UnitRepositoryPG implementa port.UnitRepository usando PostgreSQL
type UnitRepositoryPG struct {
	queries *db.Queries
}

// NewUnitRepository cria nova instância do repositório
func NewUnitRepository(queries *db.Queries) port.UnitRepository {
	return &UnitRepositoryPG{queries: queries}
}

// Create cria uma nova unidade
func (r *UnitRepositoryPG) Create(ctx context.Context, unit *entity.Unit) error {
	params := db.CreateUnitParams{
		TenantID:       uuidToPgUUID(unit.TenantID),
		Nome:           unit.Nome,
		Apelido:        unit.Apelido,
		Descricao:      unit.Descricao,
		EnderecoResumo: unit.EnderecoResumo,
		Cidade:         unit.Cidade,
		Estado:         unit.Estado,
		Timezone:       unit.Timezone,
		Ativa:          unit.Ativa,
		IsMatriz:       unit.IsMatriz,
	}

	result, err := r.queries.CreateUnit(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar unidade: %w", err)
	}

	unit.ID = pgUUIDToUUID(result.ID)
	unit.CriadoEm = result.CriadoEm.Time
	unit.AtualizadoEm = result.AtualizadoEm.Time

	return nil
}

// FindByID busca unidade por ID
func (r *UnitRepositoryPG) FindByID(ctx context.Context, tenantID, unitID uuid.UUID) (*entity.Unit, error) {
	result, err := r.queries.GetUnitByID(ctx, db.GetUnitByIDParams{
		ID:       uuidToPgUUID(unitID),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar unidade: %w", err)
	}

	return r.toDomain(result), nil
}

// FindByName busca unidade por nome
func (r *UnitRepositoryPG) FindByName(ctx context.Context, tenantID uuid.UUID, nome string) (*entity.Unit, error) {
	result, err := r.queries.GetUnitByName(ctx, db.GetUnitByNameParams{
		TenantID: uuidToPgUUID(tenantID),
		Nome:     nome,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar unidade por nome: %w", err)
	}

	return r.toDomain(result), nil
}

// FindMatriz busca a unidade matriz do tenant
func (r *UnitRepositoryPG) FindMatriz(ctx context.Context, tenantID uuid.UUID) (*entity.Unit, error) {
	result, err := r.queries.GetMatrizUnit(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar unidade matriz: %w", err)
	}

	return r.toDomain(result), nil
}

// List lista todas as unidades do tenant
func (r *UnitRepositoryPG) List(ctx context.Context, tenantID uuid.UUID) ([]*entity.Unit, error) {
	results, err := r.queries.ListUnitsByTenant(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar unidades: %w", err)
	}

	units := make([]*entity.Unit, len(results))
	for i, result := range results {
		units[i] = r.toDomain(result)
	}

	return units, nil
}

// ListActive lista apenas unidades ativas
func (r *UnitRepositoryPG) ListActive(ctx context.Context, tenantID uuid.UUID) ([]*entity.Unit, error) {
	results, err := r.queries.ListActiveUnitsByTenant(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar unidades ativas: %w", err)
	}

	units := make([]*entity.Unit, len(results))
	for i, result := range results {
		units[i] = r.toDomain(result)
	}

	return units, nil
}

// Update atualiza uma unidade
func (r *UnitRepositoryPG) Update(ctx context.Context, unit *entity.Unit) error {
	params := db.UpdateUnitParams{
		ID:             uuidToPgUUID(unit.ID),
		TenantID:       uuidToPgUUID(unit.TenantID),
		Column3:        unit.Nome,
		Apelido:        unit.Apelido,
		Descricao:      unit.Descricao,
		EnderecoResumo: unit.EnderecoResumo,
		Cidade:         unit.Cidade,
		Estado:         unit.Estado,
		Column9:        unit.Timezone,
	}

	result, err := r.queries.UpdateUnit(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ErrUnitNaoEncontrada
		}
		return fmt.Errorf("erro ao atualizar unidade: %w", err)
	}

	unit.AtualizadoEm = result.AtualizadoEm.Time
	return nil
}

// ToggleStatus alterna o status ativo/inativo
func (r *UnitRepositoryPG) ToggleStatus(ctx context.Context, tenantID, unitID uuid.UUID) (*entity.Unit, error) {
	result, err := r.queries.ToggleUnitStatus(ctx, db.ToggleUnitStatusParams{
		ID:       uuidToPgUUID(unitID),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrUnitNaoEncontrada
		}
		return nil, fmt.Errorf("erro ao alternar status: %w", err)
	}

	return r.toDomain(result), nil
}

// SetMatriz define uma unidade como matriz
func (r *UnitRepositoryPG) SetMatriz(ctx context.Context, tenantID, unitID uuid.UUID) error {
	err := r.queries.SetMatrizUnit(ctx, db.SetMatrizUnitParams{
		ID:       uuidToPgUUID(unitID),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		return fmt.Errorf("erro ao definir unidade matriz: %w", err)
	}
	return nil
}

// Delete exclui uma unidade (não pode ser matriz)
func (r *UnitRepositoryPG) Delete(ctx context.Context, tenantID, unitID uuid.UUID) error {
	err := r.queries.DeleteUnit(ctx, db.DeleteUnitParams{
		ID:       uuidToPgUUID(unitID),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		return fmt.Errorf("erro ao excluir unidade: %w", err)
	}
	return nil
}

// Count conta unidades do tenant
func (r *UnitRepositoryPG) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	count, err := r.queries.CountUnitsByTenant(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return 0, fmt.Errorf("erro ao contar unidades: %w", err)
	}
	return count, nil
}

// toDomain converte modelo sqlc para entidade de domínio
func (r *UnitRepositoryPG) toDomain(u db.Unit) *entity.Unit {
	return &entity.Unit{
		ID:             pgUUIDToUUID(u.ID),
		TenantID:       pgUUIDToUUID(u.TenantID),
		Nome:           u.Nome,
		Apelido:        u.Apelido,
		Descricao:      u.Descricao,
		EnderecoResumo: u.EnderecoResumo,
		Cidade:         u.Cidade,
		Estado:         u.Estado,
		Timezone:       u.Timezone,
		Ativa:          u.Ativa,
		IsMatriz:       u.IsMatriz,
		CriadoEm:       u.CriadoEm.Time,
		AtualizadoEm:   u.AtualizadoEm.Time,
	}
}
