package port

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// CustomerFilter contém filtros para listagem de clientes
type CustomerFilter struct {
	Search   string
	Ativo    *bool
	Tags     []string
	OrderBy  string // "nome", "criado_em", "atualizado_em"
	Page     int
	PageSize int
}

// CustomerRepository define operações para Clientes
type CustomerRepository interface {
	// Create cria um novo cliente
	Create(ctx context.Context, customer *entity.Customer) error

	// FindByID busca um cliente por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.Customer, error)

	// FindByPhone busca cliente por telefone
	FindByPhone(ctx context.Context, tenantID, phone string) (*entity.Customer, error)

	// FindByCPF busca cliente por CPF
	FindByCPF(ctx context.Context, tenantID, cpf string) (*entity.Customer, error)

	// Update atualiza um cliente existente
	Update(ctx context.Context, customer *entity.Customer) error

	// UpdateTags atualiza apenas as tags do cliente
	UpdateTags(ctx context.Context, tenantID, id string, tags []string) error

	// List lista clientes com filtros e paginação
	List(ctx context.Context, tenantID string, filter CustomerFilter) ([]*entity.Customer, int64, error)

	// ListActive lista todos os clientes ativos (para selects)
	ListActive(ctx context.Context, tenantID string) ([]*CustomerSummary, error)

	// Search busca rápida de clientes
	Search(ctx context.Context, tenantID, query string) ([]*CustomerSummary, error)

	// Inactivate inativa um cliente (soft delete)
	Inactivate(ctx context.Context, tenantID, id string) error

	// Reactivate reativa um cliente
	Reactivate(ctx context.Context, tenantID, id string) error

	// CheckPhoneExists verifica se telefone já existe
	CheckPhoneExists(ctx context.Context, tenantID, phone string, excludeID *string) (bool, error)

	// CheckCPFExists verifica se CPF já existe
	CheckCPFExists(ctx context.Context, tenantID, cpf string, excludeID *string) (bool, error)

	// CheckEmailExists verifica se email já existe
	CheckEmailExists(ctx context.Context, tenantID, email string, excludeID *string) (bool, error)

	// GetStats retorna estatísticas de clientes
	GetStats(ctx context.Context, tenantID string) (*CustomerStats, error)

	// GetWithHistory busca cliente com histórico de atendimentos
	GetWithHistory(ctx context.Context, tenantID, id string) (*CustomerWithHistory, error)

	// GetDataForExport busca dados completos para exportação LGPD
	GetDataForExport(ctx context.Context, tenantID, id string) (*CustomerExport, error)
}

// CustomerSummary dados resumidos de um cliente (para listagens/selects)
type CustomerSummary struct {
	ID       string
	Nome     string
	Telefone string
	Email    *string
	Tags     []string
}

// CustomerStats estatísticas de clientes
type CustomerStats struct {
	TotalAtivos        int64
	TotalInativos      int64
	NovosUltimos30Dias int64
	TotalGeral         int64
}

// CustomerWithHistory cliente com métricas de histórico
type CustomerWithHistory struct {
	Customer            *entity.Customer
	TotalAtendimentos   int64
	TotalGasto          string // Formatado em R$
	TicketMedio         string
	UltimoAtendimento   *time.Time
	FrequenciaMediaDias *int
}

// CustomerExport dados para exportação LGPD
type CustomerExport struct {
	Customer              *entity.Customer
	HistoricoAtendimentos []CustomerAppointmentHistory
	TotalGasto            string
	TicketMedio           string
	TotalVisitas          int64
}

// CustomerAppointmentHistory histórico de atendimento para exportação
type CustomerAppointmentHistory struct {
	Data         time.Time
	Status       string
	Profissional string
	ValorTotal   string
}
