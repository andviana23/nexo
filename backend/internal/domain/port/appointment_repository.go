package port

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
)

// AppointmentFilter contém filtros para listagem de agendamentos
type AppointmentFilter struct {
	ProfessionalID string
	CustomerID     string
	Statuses       []valueobject.AppointmentStatus // Array de status para filtrar (OR)
	StartDate      time.Time
	EndDate        time.Time
	Page           int
	PageSize       int
}

// AppointmentRepository define operações para Agendamentos
type AppointmentRepository interface {
	// Create cria um novo agendamento com seus serviços (transação)
	Create(ctx context.Context, appointment *entity.Appointment) error

	// FindByID busca um agendamento por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.Appointment, error)

	// Update atualiza um agendamento existente
	Update(ctx context.Context, appointment *entity.Appointment) error

	// Delete remove um agendamento (soft delete via status CANCELED)
	Delete(ctx context.Context, tenantID, id string) error

	// List lista agendamentos com filtros
	List(ctx context.Context, tenantID string, filter AppointmentFilter) ([]*entity.Appointment, int64, error)

	// ListByProfessionalAndDateRange lista agendamentos de um profissional em um período
	ListByProfessionalAndDateRange(
		ctx context.Context,
		tenantID string,
		professionalID string,
		startDate, endDate time.Time,
	) ([]*entity.Appointment, error)

	// ListByCustomer lista agendamentos de um cliente
	ListByCustomer(ctx context.Context, tenantID, customerID string) ([]*entity.Appointment, error)

	// CheckConflict verifica se há conflito de horário
	// Retorna true se houver conflito
	CheckConflict(
		ctx context.Context,
		tenantID string,
		professionalID string,
		startTime, endTime time.Time,
		excludeAppointmentID string, // Para ignorar o próprio agendamento ao reagendar
	) (bool, error)

	// CheckBlockedTimeConflict verifica se há conflito com horários bloqueados
	// Retorna true se houver conflito com blocked_times
	CheckBlockedTimeConflict(
		ctx context.Context,
		tenantID string,
		professionalID string,
		startTime, endTime time.Time,
	) (bool, error)

	// CheckMinimumIntervalConflict verifica se há conflito de intervalo mínimo
	// Retorna true se o agendamento estiver muito próximo de outro (menos de intervalMinutes)
	CheckMinimumIntervalConflict(
		ctx context.Context,
		tenantID string,
		professionalID string,
		startTime, endTime time.Time,
		excludeAppointmentID string,
		intervalMinutes int,
	) (bool, error)

	// CountByStatus conta agendamentos por status (para dashboard)
	CountByStatus(ctx context.Context, tenantID string, status valueobject.AppointmentStatus) (int64, error)

	// GetDailyStats retorna estatísticas diárias
	GetDailyStats(ctx context.Context, tenantID string, date time.Time) (*AppointmentDailyStats, error)
}

// AppointmentDailyStats estatísticas diárias de agendamentos
type AppointmentDailyStats struct {
	TotalAppointments int64
	CompletedCount    int64
	CanceledCount     int64
	NoShowCount       int64
	TotalRevenue      valueobject.Money
}

// ProfessionalRepository define operações para consulta de profissionais (leitura)
// Usado pelo módulo de agendamentos para validar profissionais
type ProfessionalReader interface {
	// Exists verifica se um profissional existe e está ativo
	Exists(ctx context.Context, tenantID, professionalID string) (bool, error)

	// FindByID busca dados básicos do profissional
	FindByID(ctx context.Context, tenantID, professionalID string) (*ProfessionalInfo, error)

	// ListActive lista profissionais ativos
	ListActive(ctx context.Context, tenantID string) ([]*ProfessionalInfo, error)
}

// ProfessionalInfo dados básicos de um profissional
type ProfessionalInfo struct {
	ID           string
	Name         string
	Status       string
	Color        string  // Cor para exibição no calendário
	Comissao     *string // Taxa de comissão (decimal como string)
	TipoComissao *string // PERCENTUAL ou FIXO
}

// CustomerReader define operações para consulta de clientes (leitura)
type CustomerReader interface {
	// Exists verifica se um cliente existe e está ativo
	Exists(ctx context.Context, tenantID, customerID string) (bool, error)

	// FindByID busca dados básicos do cliente
	FindByID(ctx context.Context, tenantID, customerID string) (*CustomerInfo, error)
}

// CustomerInfo dados básicos de um cliente
type CustomerInfo struct {
	ID    string
	Name  string
	Phone string
	Email string
}

// ServiceReader define operações para consulta de serviços (leitura)
type ServiceReader interface {
	// Exists verifica se um serviço existe e está ativo
	Exists(ctx context.Context, tenantID, serviceID string) (bool, error)

	// FindByID busca dados do serviço
	FindByID(ctx context.Context, tenantID, serviceID string) (*ServiceInfo, error)

	// FindByIDs busca múltiplos serviços
	FindByIDs(ctx context.Context, tenantID string, serviceIDs []string) ([]*ServiceInfo, error)
}

// ServiceInfo dados de um serviço
type ServiceInfo struct {
	ID       string
	Name     string
	Price    valueobject.Money
	Duration int // minutos
	Active   bool
	Comissao *string // Taxa de comissão específica do serviço (decimal como string)
}
