package repository

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// CommissionPeriodRepository define as operações de repositório para períodos de comissão
type CommissionPeriodRepository interface {
	// Create cria um novo período de comissão
	Create(ctx context.Context, period *entity.CommissionPeriod) (*entity.CommissionPeriod, error)

	// GetByID busca um período de comissão por ID
	GetByID(ctx context.Context, tenantID, id string) (*entity.CommissionPeriod, error)

	// List lista períodos de comissão com filtros
	List(ctx context.Context, tenantID string, professionalID *string, status *string, limit, offset int) ([]*entity.CommissionPeriod, error)

	// GetByProfessional busca períodos de comissão por profissional
	GetByProfessional(ctx context.Context, tenantID, professionalID string) ([]*entity.CommissionPeriod, error)

	// GetOpenByProfessional busca o período aberto de um profissional
	GetOpenByProfessional(ctx context.Context, tenantID, professionalID string) (*entity.CommissionPeriod, error)

	// GetByDateRange busca períodos de comissão por intervalo de datas
	GetByDateRange(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*entity.CommissionPeriod, error)

	// GetSummary retorna totais do período
	GetSummary(ctx context.Context, tenantID, periodID string) (*entity.CommissionPeriodSummary, error)

	// Update atualiza um período de comissão
	Update(ctx context.Context, period *entity.CommissionPeriod) (*entity.CommissionPeriod, error)

	// Close fecha um período de comissão
	Close(ctx context.Context, tenantID, id, closedBy string) (*entity.CommissionPeriod, error)

	// MarkAsPaid marca um período como pago
	MarkAsPaid(ctx context.Context, tenantID, id, paidBy string) (*entity.CommissionPeriod, error)

	// Delete remove um período de comissão (somente se ABERTO)
	Delete(ctx context.Context, tenantID, id string) error
}
