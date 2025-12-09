package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
)

// UserUnitRepository define as operações de persistência para vínculos usuário-unidade
type UserUnitRepository interface {
	// Create cria um novo vínculo
	Create(ctx context.Context, userUnit *entity.UserUnit) error

	// FindByUserAndUnit busca vínculo específico
	FindByUserAndUnit(ctx context.Context, userID, unitID uuid.UUID) (*entity.UserUnit, error)

	// FindUserDefaultUnit busca a unidade padrão do usuário
	FindUserDefaultUnit(ctx context.Context, userID uuid.UUID) (*entity.UserUnitWithDetails, error)

	// ListByUser lista unidades do usuário
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.UserUnitWithDetails, error)

	// ListByUnit lista usuários da unidade
	ListByUnit(ctx context.Context, unitID uuid.UUID) ([]*entity.UserUnit, error)

	// CheckAccess verifica se usuário tem acesso à unidade
	CheckAccess(ctx context.Context, userID, unitID uuid.UUID) (bool, error)

	// SetDefault define unidade padrão do usuário
	SetDefault(ctx context.Context, userID, unitID uuid.UUID) error

	// UpdateRole atualiza o papel do usuário na unidade
	UpdateRole(ctx context.Context, userID, unitID uuid.UUID, role *string) error

	// Delete remove vínculo
	Delete(ctx context.Context, userID, unitID uuid.UUID) error

	// DeleteAllByUser remove todos os vínculos de um usuário
	DeleteAllByUser(ctx context.Context, userID uuid.UUID) error

	// DeleteAllByUnit remove todos os vínculos de uma unidade
	DeleteAllByUnit(ctx context.Context, unitID uuid.UUID) error

	// CountByUser conta unidades do usuário
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)

	// CountByUnit conta usuários da unidade
	CountByUnit(ctx context.Context, unitID uuid.UUID) (int64, error)
}
