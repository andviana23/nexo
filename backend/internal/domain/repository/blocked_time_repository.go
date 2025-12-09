package repository

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// BlockedTimeRepository define as operações de repositório para bloqueios de horário
type BlockedTimeRepository interface {
	// Create cria um novo bloqueio
	Create(ctx context.Context, blockedTime *entity.BlockedTime) (*entity.BlockedTime, error)

	// GetByID busca um bloqueio por ID
	GetByID(ctx context.Context, tenantID, id string) (*entity.BlockedTime, error)

	// List lista bloqueios com filtros opcionais
	List(ctx context.Context, tenantID string, professionalID *string, startDate, endDate *time.Time) ([]*entity.BlockedTime, error)

	// GetInRange busca bloqueios em um intervalo de tempo específico
	GetInRange(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time) ([]*entity.BlockedTime, error)

	// CheckConflict verifica se há conflito de horário
	CheckConflict(ctx context.Context, tenantID, professionalID string, startTime, endTime time.Time, excludeID *string) (bool, error)

	// Update atualiza um bloqueio existente
	Update(ctx context.Context, blockedTime *entity.BlockedTime) (*entity.BlockedTime, error)

	// Delete remove um bloqueio
	Delete(ctx context.Context, tenantID, id string) error
}
