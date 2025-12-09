package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

// Appointment representa um agendamento no sistema
type Appointment struct {
	ID             string
	TenantID       uuid.UUID
	ProfessionalID string
	CustomerID     string

	StartTime time.Time
	EndTime   time.Time

	CheckedInAt *time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time

	Status                valueobject.AppointmentStatus
	TotalPrice            valueobject.Money
	Notes                 string
	CanceledReason        string
	GoogleCalendarEventID string
	CommandID             string // Comanda vinculada ao agendamento

	// Relacionamentos (carregados via join)
	Services []AppointmentService

	// Dados denormalizados (read-only, vêm do JOIN)
	ProfessionalName string
	CustomerName     string
	CustomerPhone    string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// AppointmentService representa um serviço vinculado a um agendamento
type AppointmentService struct {
	AppointmentID     string
	ServiceID         string
	PriceAtBooking    valueobject.Money
	DurationAtBooking int // em minutos
	CreatedAt         time.Time

	// Dados do serviço (carregados via join)
	ServiceName string
}

// NewAppointment cria um novo agendamento validado
func NewAppointment(
	tenantID uuid.UUID,
	professionalID string,
	customerID string,
	startTime time.Time,
	services []AppointmentService,
) (*Appointment, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if professionalID == "" {
		return nil, domain.ErrAppointmentProfessionalRequired
	}
	if customerID == "" {
		return nil, domain.ErrAppointmentCustomerRequired
	}
	if startTime.IsZero() {
		return nil, domain.ErrAppointmentStartTimeRequired
	}
	if len(services) == 0 {
		return nil, domain.ErrAppointmentServicesRequired
	}

	// Calcula duração total e preço total
	var totalDuration int
	totalPrice := valueobject.Zero()
	for _, s := range services {
		totalDuration += s.DurationAtBooking
		totalPrice = totalPrice.Add(s.PriceAtBooking)
	}

	endTime := startTime.Add(time.Duration(totalDuration) * time.Minute)

	now := time.Now()
	return &Appointment{
		ID:             uuid.NewString(),
		TenantID:       tenantID,
		ProfessionalID: professionalID,
		CustomerID:     customerID,
		StartTime:      startTime,
		EndTime:        endTime,
		Status:         valueobject.AppointmentStatusCreated,
		TotalPrice:     totalPrice,
		Services:       services,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Confirm confirma o agendamento
func (a *Appointment) Confirm() error {
	if !a.Status.CanTransitionTo(valueobject.AppointmentStatusConfirmed) {
		return domain.ErrAppointmentInvalidStatusTransition
	}
	a.Status = valueobject.AppointmentStatusConfirmed
	a.UpdatedAt = time.Now()
	return nil
}

// CheckIn marca que o cliente chegou para o atendimento
func (a *Appointment) CheckIn() error {
	if !a.Status.CanTransitionTo(valueobject.AppointmentStatusCheckedIn) {
		return domain.ErrAppointmentInvalidStatusTransition
	}
	a.Status = valueobject.AppointmentStatusCheckedIn
	if a.CheckedInAt == nil {
		now := time.Now()
		a.CheckedInAt = &now
	}
	a.UpdatedAt = time.Now()
	return nil
}

// StartService inicia o atendimento
func (a *Appointment) StartService() error {
	if !a.Status.CanTransitionTo(valueobject.AppointmentStatusInService) {
		return domain.ErrAppointmentInvalidStatusTransition
	}
	a.Status = valueobject.AppointmentStatusInService
	if a.StartedAt == nil {
		now := time.Now()
		a.StartedAt = &now
	}
	a.UpdatedAt = time.Now()
	return nil
}

// FinishService marca o atendimento como finalizado, aguardando pagamento
func (a *Appointment) FinishService() error {
	if !a.Status.CanTransitionTo(valueobject.AppointmentStatusAwaitingPayment) {
		return domain.ErrAppointmentInvalidStatusTransition
	}
	a.Status = valueobject.AppointmentStatusAwaitingPayment
	if a.FinishedAt == nil {
		now := time.Now()
		a.FinishedAt = &now
	}
	a.UpdatedAt = time.Now()
	return nil
}

// Complete finaliza o atendimento (após pagamento confirmado)
func (a *Appointment) Complete() error {
	if !a.Status.CanTransitionTo(valueobject.AppointmentStatusDone) {
		return domain.ErrAppointmentInvalidStatusTransition
	}
	a.Status = valueobject.AppointmentStatusDone
	a.UpdatedAt = time.Now()
	return nil
}

// Cancel cancela o agendamento
func (a *Appointment) Cancel(reason string) error {
	if !a.Status.CanTransitionTo(valueobject.AppointmentStatusCanceled) {
		return domain.ErrAppointmentInvalidStatusTransition
	}
	a.Status = valueobject.AppointmentStatusCanceled
	a.CanceledReason = reason
	a.UpdatedAt = time.Now()
	return nil
}

// MarkNoShow marca como cliente não compareceu
func (a *Appointment) MarkNoShow() error {
	if !a.Status.CanTransitionTo(valueobject.AppointmentStatusNoShow) {
		return domain.ErrAppointmentInvalidStatusTransition
	}
	a.Status = valueobject.AppointmentStatusNoShow
	a.UpdatedAt = time.Now()
	return nil
}

// Reschedule reagenda o agendamento
func (a *Appointment) Reschedule(newStartTime time.Time) error {
	if a.Status == valueobject.AppointmentStatusDone ||
		a.Status == valueobject.AppointmentStatusCanceled ||
		a.Status == valueobject.AppointmentStatusNoShow {
		return domain.ErrAppointmentCannotReschedule
	}

	// Mantém a duração original
	duration := a.EndTime.Sub(a.StartTime)
	a.StartTime = newStartTime
	a.EndTime = newStartTime.Add(duration)
	a.UpdatedAt = time.Now()
	return nil
}

// SetNotes define observações
func (a *Appointment) SetNotes(notes string) {
	a.Notes = notes
	a.UpdatedAt = time.Now()
}

// SetGoogleCalendarEventID vincula ao evento do Google Calendar
func (a *Appointment) SetGoogleCalendarEventID(eventID string) {
	a.GoogleCalendarEventID = eventID
	a.UpdatedAt = time.Now()
}

// Duration retorna a duração total em minutos
func (a *Appointment) Duration() int {
	return int(a.EndTime.Sub(a.StartTime).Minutes())
}

// IsActive verifica se o agendamento está ativo (não finalizado/cancelado)
func (a *Appointment) IsActive() bool {
	return a.Status == valueobject.AppointmentStatusCreated ||
		a.Status == valueobject.AppointmentStatusConfirmed ||
		a.Status == valueobject.AppointmentStatusCheckedIn ||
		a.Status == valueobject.AppointmentStatusInService ||
		a.Status == valueobject.AppointmentStatusAwaitingPayment
}

// IsPast verifica se o agendamento é no passado
func (a *Appointment) IsPast() bool {
	return a.EndTime.Before(time.Now())
}

// IsFuture verifica se o agendamento é no futuro
func (a *Appointment) IsFuture() bool {
	return a.StartTime.After(time.Now())
}

// IsToday verifica se o agendamento é hoje
func (a *Appointment) IsToday() bool {
	now := time.Now()
	return a.StartTime.Year() == now.Year() &&
		a.StartTime.YearDay() == now.YearDay()
}

// ConflictsWith verifica se há conflito de horário com outro agendamento
func (a *Appointment) ConflictsWith(other *Appointment) bool {
	// Mesmo profissional
	if a.ProfessionalID != other.ProfessionalID {
		return false
	}
	// Ignora se for o mesmo agendamento
	if a.ID == other.ID {
		return false
	}
	// Verifica sobreposição de horários
	return a.StartTime.Before(other.EndTime) && a.EndTime.After(other.StartTime)
}

// Validate valida as regras de negócio
func (a *Appointment) Validate() error {
	if a.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if a.ProfessionalID == "" {
		return domain.ErrAppointmentProfessionalRequired
	}
	if a.CustomerID == "" {
		return domain.ErrAppointmentCustomerRequired
	}
	if a.StartTime.IsZero() {
		return domain.ErrAppointmentStartTimeRequired
	}
	if a.EndTime.IsZero() || !a.EndTime.After(a.StartTime) {
		return domain.ErrAppointmentInvalidTimeRange
	}
	if !a.Status.IsValid() {
		return domain.ErrAppointmentInvalidStatus
	}
	if a.TotalPrice.IsNegative() {
		return domain.ErrValorNegativo
	}
	return nil
}
