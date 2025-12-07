package repository

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
)

// CommissionItemRepository define as operações de repositório para itens de comissão
type CommissionItemRepository interface {
	// Create cria um novo item de comissão
	Create(ctx context.Context, item *entity.CommissionItem) (*entity.CommissionItem, error)

	// CreateBatch cria múltiplos itens de comissão de uma vez
	CreateBatch(ctx context.Context, items []*entity.CommissionItem) ([]*entity.CommissionItem, error)

	// GetByID busca um item de comissão por ID
	GetByID(ctx context.Context, tenantID, id string) (*entity.CommissionItem, error)

	// List lista itens de comissão com filtros
	List(ctx context.Context, tenantID string, professionalID *string, periodID *string, status *string, limit, offset int) ([]*entity.CommissionItem, error)

	// GetByProfessional busca itens de comissão por profissional
	GetByProfessional(ctx context.Context, tenantID, professionalID string) ([]*entity.CommissionItem, error)

	// GetByPeriod busca itens de comissão por período
	GetByPeriod(ctx context.Context, tenantID, periodID string) ([]*entity.CommissionItem, error)

	// GetByCommandItem busca item de comissão por item de comanda
	GetByCommandItem(ctx context.Context, tenantID, commandItemID string) (*entity.CommissionItem, error)

	// ListByCommand busca itens de comissão vinculados a uma comanda
	ListByCommand(ctx context.Context, tenantID string, commandID uuid.UUID) ([]*entity.CommissionItem, error)

	// GetPendingByProfessional busca itens pendentes de um profissional
	GetPendingByProfessional(ctx context.Context, tenantID, professionalID string) ([]*entity.CommissionItem, error)

	// GetByDateRange busca itens de comissão por intervalo de datas
	GetByDateRange(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*entity.CommissionItem, error)

	// GetTotalByPeriod retorna o total de comissão de um período
	GetTotalByPeriod(ctx context.Context, tenantID, periodID string) (float64, error)

	// GetTotalByProfessionalInRange retorna o total de comissão de um profissional em um intervalo
	GetTotalByProfessionalInRange(ctx context.Context, tenantID, professionalID string, startDate, endDate time.Time) (float64, error)

	// GetSummaryByProfessional retorna resumo de comissões por profissional
	GetSummaryByProfessional(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*entity.CommissionSummary, error)

	// GetSummaryByService retorna resumo de comissões por serviço
	GetSummaryByService(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*entity.CommissionByService, error)

	// Process processa um item (vincula a um período)
	Process(ctx context.Context, tenantID, id, periodID string) (*entity.CommissionItem, error)

	// AssignToPeriod vincula itens pendentes a um período
	AssignToPeriod(ctx context.Context, tenantID, professionalID, periodID string, startDate, endDate time.Time) (int64, error)

	// Update atualiza um item de comissão
	Update(ctx context.Context, item *entity.CommissionItem) (*entity.CommissionItem, error)

	// Delete remove um item de comissão (somente se PENDENTE)
	Delete(ctx context.Context, tenantID, id string) error

	// DeleteByCommandItem remove item de comissão por item de comanda
	DeleteByCommandItem(ctx context.Context, tenantID, commandItemID string) error
}
