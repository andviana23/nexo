package appointment

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"go.uber.org/zap"
)

// CreateAppointmentInput dados de entrada para criar agendamento
type CreateAppointmentInput struct {
	TenantID       string
	ProfessionalID string
	CustomerID     string
	StartTime      time.Time
	ServiceIDs     []string
	Notes          string
}

// CreateAppointmentUseCase implementa a criação de agendamentos
type CreateAppointmentUseCase struct {
	appointmentRepo    port.AppointmentRepository
	serviceReader      port.ServiceReader
	professionalReader port.ProfessionalReader
	customerReader     port.CustomerReader
	logger             *zap.Logger
}

// NewCreateAppointmentUseCase cria nova instância do use case
func NewCreateAppointmentUseCase(
	appointmentRepo port.AppointmentRepository,
	serviceReader port.ServiceReader,
	professionalReader port.ProfessionalReader,
	customerReader port.CustomerReader,
	logger *zap.Logger,
) *CreateAppointmentUseCase {
	return &CreateAppointmentUseCase{
		appointmentRepo:    appointmentRepo,
		serviceReader:      serviceReader,
		professionalReader: professionalReader,
		customerReader:     customerReader,
		logger:             logger,
	}
}

// Execute cria um novo agendamento
func (uc *CreateAppointmentUseCase) Execute(ctx context.Context, input CreateAppointmentInput) (*entity.Appointment, error) {
	// 1. Validações básicas
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if input.ProfessionalID == "" {
		return nil, domain.ErrAppointmentProfessionalRequired
	}
	if input.CustomerID == "" {
		return nil, domain.ErrAppointmentCustomerRequired
	}
	if input.StartTime.IsZero() {
		return nil, domain.ErrAppointmentStartTimeRequired
	}
	if len(input.ServiceIDs) == 0 {
		return nil, domain.ErrAppointmentServicesRequired
	}

	// 2. Verificar se profissional existe e está ativo
	profExists, err := uc.professionalReader.Exists(ctx, input.TenantID, input.ProfessionalID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar profissional: %w", err)
	}
	if !profExists {
		return nil, domain.ErrAppointmentProfessionalNotFound
	}

	// 3. Verificar se cliente existe e está ativo
	customerExists, err := uc.customerReader.Exists(ctx, input.TenantID, input.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar cliente: %w", err)
	}
	if !customerExists {
		return nil, domain.ErrAppointmentCustomerNotFound
	}

	// 4. Buscar dados dos serviços
	services, err := uc.serviceReader.FindByIDs(ctx, input.TenantID, input.ServiceIDs)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar serviços: %w", err)
	}
	if len(services) != len(input.ServiceIDs) {
		return nil, domain.ErrAppointmentServiceNotFound
	}

	// 5. Montar lista de serviços do agendamento
	appointmentServices := make([]entity.AppointmentService, 0, len(services))
	for _, svc := range services {
		if !svc.Active {
			return nil, fmt.Errorf("serviço %s está inativo", svc.Name)
		}
		appointmentServices = append(appointmentServices, entity.AppointmentService{
			ServiceID:         svc.ID,
			ServiceName:       svc.Name,
			PriceAtBooking:    svc.Price,
			DurationAtBooking: svc.Duration,
		})
	}

	// 6. Criar entidade de agendamento (calcula end_time e total_price automaticamente)
	appointment, err := entity.NewAppointment(
		input.TenantID,
		input.ProfessionalID,
		input.CustomerID,
		input.StartTime,
		appointmentServices,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar agendamento: %w", err)
	}

	// 7. Verificar conflito de horário
	hasConflict, err := uc.appointmentRepo.CheckConflict(
		ctx,
		input.TenantID,
		input.ProfessionalID,
		appointment.StartTime,
		appointment.EndTime,
		"", // Não excluir nenhum agendamento
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar conflito: %w", err)
	}
	if hasConflict {
		return nil, domain.ErrAppointmentConflict
	}

	// 8. Definir observações
	if input.Notes != "" {
		appointment.SetNotes(input.Notes)
	}

	// 9. Persistir
	if err := uc.appointmentRepo.Create(ctx, appointment); err != nil {
		return nil, fmt.Errorf("erro ao salvar agendamento: %w", err)
	}

	uc.logger.Info("Agendamento criado",
		zap.String("tenant_id", input.TenantID),
		zap.String("appointment_id", appointment.ID),
		zap.String("professional_id", input.ProfessionalID),
		zap.String("customer_id", input.CustomerID),
		zap.Time("start_time", appointment.StartTime),
		zap.String("total_price", appointment.TotalPrice.String()),
	)

	return appointment, nil
}

// ListAppointmentsInput dados de entrada para listar agendamentos
type ListAppointmentsInput struct {
	TenantID       string
	ProfessionalID string
	CustomerID     string
	Status         valueobject.AppointmentStatus
	StartDate      time.Time
	EndDate        time.Time
	Page           int
	PageSize       int
}

// ListAppointmentsOutput saída da listagem
type ListAppointmentsOutput struct {
	Appointments []*entity.Appointment
	Total        int64
	Page         int
	PageSize     int
}

// ListAppointmentsUseCase implementa a listagem de agendamentos
type ListAppointmentsUseCase struct {
	repo   port.AppointmentRepository
	logger *zap.Logger
}

// NewListAppointmentsUseCase cria nova instância do use case
func NewListAppointmentsUseCase(
	repo port.AppointmentRepository,
	logger *zap.Logger,
) *ListAppointmentsUseCase {
	return &ListAppointmentsUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute lista agendamentos com filtros
func (uc *ListAppointmentsUseCase) Execute(ctx context.Context, input ListAppointmentsInput) (*ListAppointmentsOutput, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	// Defaults
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	filter := port.AppointmentFilter{
		ProfessionalID: input.ProfessionalID,
		CustomerID:     input.CustomerID,
		Status:         input.Status,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		Page:           input.Page,
		PageSize:       input.PageSize,
	}

	appointments, total, err := uc.repo.List(ctx, input.TenantID, filter)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar agendamentos: %w", err)
	}

	return &ListAppointmentsOutput{
		Appointments: appointments,
		Total:        total,
		Page:         input.Page,
		PageSize:     input.PageSize,
	}, nil
}
