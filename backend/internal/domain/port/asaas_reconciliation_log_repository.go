package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// AsaasReconciliationLogRepository define operações de persistência para logs de conciliação Asaas
type AsaasReconciliationLogRepository interface {
	// Create insere um novo log de conciliação
	Create(ctx context.Context, log *entity.AsaasReconciliationLog) error

	// Update atualiza um log de conciliação existente
	Update(ctx context.Context, log *entity.AsaasReconciliationLog) error

	// GetByID busca log por ID
	GetByID(ctx context.Context, tenantID, id string) (*entity.AsaasReconciliationLog, error)

	// ListByTenant lista logs de conciliação de um tenant
	ListByTenant(ctx context.Context, tenantID string, limit int) ([]*entity.AsaasReconciliationLog, error)

	// GetLatest busca o último log de conciliação de um tenant
	GetLatest(ctx context.Context, tenantID string) (*entity.AsaasReconciliationLog, error)
}
