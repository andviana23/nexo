package repository

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// AdvanceRepository define as operações de repositório para adiantamentos
type AdvanceRepository interface {
	// Create cria um novo adiantamento
	Create(ctx context.Context, advance *entity.Advance) (*entity.Advance, error)

	// GetByID busca um adiantamento por ID
	GetByID(ctx context.Context, tenantID, id string) (*entity.Advance, error)

	// List lista adiantamentos com filtros
	List(ctx context.Context, tenantID string, professionalID *string, status *string, limit, offset int) ([]*entity.Advance, error)

	// GetByProfessional busca adiantamentos por profissional
	GetByProfessional(ctx context.Context, tenantID, professionalID string) ([]*entity.Advance, error)

	// GetPendingByProfessional busca adiantamentos pendentes de um profissional
	GetPendingByProfessional(ctx context.Context, tenantID, professionalID string) ([]*entity.Advance, error)

	// GetApprovedByProfessional busca adiantamentos aprovados (não deduzidos) de um profissional
	GetApprovedByProfessional(ctx context.Context, tenantID, professionalID string) ([]*entity.Advance, error)

	// GetByDateRange busca adiantamentos por intervalo de datas
	GetByDateRange(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*entity.Advance, error)

	// GetTotalPendingByProfessional retorna o total pendente de um profissional
	GetTotalPendingByProfessional(ctx context.Context, tenantID, professionalID string) (float64, error)

	// GetTotalApprovedByProfessional retorna o total aprovado (não deduzido) de um profissional
	GetTotalApprovedByProfessional(ctx context.Context, tenantID, professionalID string) (float64, error)

	// Approve aprova um adiantamento
	Approve(ctx context.Context, tenantID, id, approvedBy string) (*entity.Advance, error)

	// Reject rejeita um adiantamento com motivo
	Reject(ctx context.Context, tenantID, id, rejectedBy, reason string) (*entity.Advance, error)

	// MarkDeducted marca um adiantamento como deduzido
	MarkDeducted(ctx context.Context, tenantID, id, periodID string) (*entity.Advance, error)

	// Cancel cancela um adiantamento
	Cancel(ctx context.Context, tenantID, id string) (*entity.Advance, error)

	// Update atualiza um adiantamento
	Update(ctx context.Context, advance *entity.Advance) (*entity.Advance, error)

	// Delete remove um adiantamento (somente se PENDING)
	Delete(ctx context.Context, tenantID, id string) error
}
