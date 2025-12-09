package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
)

// UnitRepository define as operações de persistência para unidades
type UnitRepository interface {
	// Create cria uma nova unidade
	Create(ctx context.Context, unit *entity.Unit) error

	// FindByID busca unidade por ID
	FindByID(ctx context.Context, tenantID, unitID uuid.UUID) (*entity.Unit, error)

	// FindByName busca unidade por nome
	FindByName(ctx context.Context, tenantID uuid.UUID, nome string) (*entity.Unit, error)

	// FindMatriz busca a unidade matriz do tenant
	FindMatriz(ctx context.Context, tenantID uuid.UUID) (*entity.Unit, error)

	// List lista todas as unidades do tenant
	List(ctx context.Context, tenantID uuid.UUID) ([]*entity.Unit, error)

	// ListActive lista apenas unidades ativas
	ListActive(ctx context.Context, tenantID uuid.UUID) ([]*entity.Unit, error)

	// Update atualiza uma unidade
	Update(ctx context.Context, unit *entity.Unit) error

	// ToggleStatus alterna o status ativo/inativo
	ToggleStatus(ctx context.Context, tenantID, unitID uuid.UUID) (*entity.Unit, error)

	// SetMatriz define uma unidade como matriz
	SetMatriz(ctx context.Context, tenantID, unitID uuid.UUID) error

	// Delete exclui uma unidade (não pode ser matriz)
	Delete(ctx context.Context, tenantID, unitID uuid.UUID) error

	// Count conta unidades do tenant
	Count(ctx context.Context, tenantID uuid.UUID) (int64, error)
}
