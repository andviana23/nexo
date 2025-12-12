package appointment

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateAppointmentInput dados de entrada para criar agendamento
type CreateAppointmentInput struct {
	TenantID       string
	UnitID         string
	ProfessionalID string
	CustomerID     string
	StartTime      time.Time
	ServiceIDs     []string
	Notes          string
}

// CreateAppointmentOutput dados de saída da criação de agendamento
type CreateAppointmentOutput struct {
	Appointment *entity.Appointment
	Command     *entity.Command // G-001: Comanda criada automaticamente
}

// CreateAppointmentUseCase implementa a criação de agendamentos
type CreateAppointmentUseCase struct {
	appointmentRepo    port.AppointmentRepository
	commandRepo        port.CommandRepository // G-001: Para criar comanda automaticamente
	serviceReader      port.ServiceReader
	professionalReader port.ProfessionalReader
	customerReader     port.CustomerReader
	logger             *zap.Logger
}

// NewCreateAppointmentUseCase cria nova instância do use case
func NewCreateAppointmentUseCase(
	appointmentRepo port.AppointmentRepository,
	commandRepo port.CommandRepository, // G-001: Adicionado CommandRepository
	serviceReader port.ServiceReader,
	professionalReader port.ProfessionalReader,
	customerReader port.CustomerReader,
	logger *zap.Logger,
) *CreateAppointmentUseCase {
	return &CreateAppointmentUseCase{
		appointmentRepo:    appointmentRepo,
		commandRepo:        commandRepo,
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
	if input.UnitID == "" {
		return nil, domain.ErrUnitIDRequired
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
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}
	unitUUID, err := uuid.Parse(input.UnitID)
	if err != nil {
		return nil, fmt.Errorf("unit_id inválido: %w", err)
	}

	appointment, err := entity.NewAppointment(
		tenantUUID,
		unitUUID,
		input.ProfessionalID,
		input.CustomerID,
		input.StartTime,
		appointmentServices,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar agendamento: %w", err)
	}

	// 7. Verificar conflito de horário com outros agendamentos
	hasConflict, err := uc.appointmentRepo.CheckConflict(
		ctx,
		input.TenantID,
		input.UnitID, // Added UnitID
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

	// 8. Verificar conflito com horários bloqueados (blocked_times)
	hasBlockedConflict, err := uc.appointmentRepo.CheckBlockedTimeConflict(
		ctx,
		input.TenantID,
		input.UnitID, // Added UnitID
		input.ProfessionalID,
		appointment.StartTime,
		appointment.EndTime,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar bloqueio: %w", err)
	}
	if hasBlockedConflict {
		return nil, domain.ErrAppointmentBlockedTimeConflict
	}

	// 9. Verificar intervalo mínimo entre agendamentos (RN-AGE-003: 10 minutos)
	const minimumIntervalMinutes = 10
	hasIntervalConflict, err := uc.appointmentRepo.CheckMinimumIntervalConflict(
		ctx,
		input.TenantID,
		input.UnitID, // Added UnitID
		input.ProfessionalID,
		appointment.StartTime,
		appointment.EndTime,
		"",
		minimumIntervalMinutes,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar intervalo mínimo: %w", err)
	}
	if hasIntervalConflict {
		return nil, domain.ErrAppointmentMinimumInterval
	}

	// 10. Definir observações
	if input.Notes != "" {
		appointment.SetNotes(input.Notes)
	}

	// 11. Persistir agendamento
	if err := uc.appointmentRepo.Create(ctx, appointment); err != nil {
		return nil, fmt.Errorf("erro ao salvar agendamento: %w", err)
	}

	// 12. G-001: Criar comanda automaticamente vinculada ao agendamento
	appointmentUUID, _ := uuid.Parse(appointment.ID)
	customerUUID, _ := uuid.Parse(input.CustomerID)

	command, err := entity.NewCommand(tenantUUID, customerUUID, &appointmentUUID)
	if err != nil {
		uc.logger.Warn("Falha ao criar entidade de comanda",
			zap.String("appointment_id", appointment.ID),
			zap.Error(err))
		// Não bloqueia criação do agendamento - comanda pode ser criada depois
	} else {
		// Adicionar os serviços como itens da comanda
		for _, svc := range services {
			item := entity.CommandItem{
				ID:            uuid.New(),
				CommandID:     command.ID,
				Tipo:          entity.CommandItemTypeServico,
				ItemID:        uuid.MustParse(svc.ID),
				Descricao:     svc.Name,
				PrecoUnitario: svc.Price.Value().InexactFloat64(),
				Quantidade:    1,
				PrecoFinal:    svc.Price.Value().InexactFloat64(),
				CriadoEm:      time.Now(),
			}
			if err := command.AddItem(item); err != nil {
				uc.logger.Warn("Falha ao adicionar item à comanda",
					zap.String("service_id", svc.ID),
					zap.Error(err))
			}
		}

		// Persistir comanda (o número será gerado automaticamente pelo repositório)
		if err := uc.commandRepo.Create(ctx, command); err != nil {
			uc.logger.Warn("Falha ao salvar comanda",
				zap.String("appointment_id", appointment.ID),
				zap.Error(err))
			// Não bloqueia - comanda pode ser criada manualmente depois
			command = nil
		} else {
			uc.logger.Info("Comanda criada automaticamente",
				zap.String("appointment_id", appointment.ID),
				zap.String("command_id", command.ID.String()),
				zap.String("numero", *command.Numero),
				zap.Int("itens", len(command.Items)),
			)
		}
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
	UnitID         string
	ProfessionalID string
	CustomerID     string
	Statuses       []valueobject.AppointmentStatus // Array de status
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
		UnitID:         input.UnitID,
		ProfessionalID: input.ProfessionalID,
		CustomerID:     input.CustomerID,
		Statuses:       input.Statuses,
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
