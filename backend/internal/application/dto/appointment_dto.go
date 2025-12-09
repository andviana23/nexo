package dto

import "time"

// =============================================================================
// DTOs para Agendamentos (Appointments)
// =============================================================================

// CreateAppointmentRequest requisição para criar agendamento
type CreateAppointmentRequest struct {
	ProfessionalID string    `json:"professional_id" validate:"required,uuid"`
	CustomerID     string    `json:"customer_id" validate:"required,uuid"`
	StartTime      time.Time `json:"start_time" validate:"required"`
	ServiceIDs     []string  `json:"service_ids" validate:"required,min=1,dive,uuid"`
	Notes          string    `json:"notes,omitempty"`
}

// UpdateAppointmentStatusRequest requisição para atualizar status
type UpdateAppointmentStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=CREATED CONFIRMED CHECKED_IN IN_SERVICE AWAITING_PAYMENT DONE NO_SHOW CANCELED"`
	Reason string `json:"reason,omitempty"`
}

// RescheduleAppointmentRequest requisição para reagendar
type RescheduleAppointmentRequest struct {
	NewStartTime   time.Time `json:"new_start_time" validate:"required"`
	ProfessionalID string    `json:"professional_id,omitempty" validate:"omitempty,uuid"`
}

// CancelAppointmentRequest requisição para cancelar
type CancelAppointmentRequest struct {
	Reason string `json:"reason,omitempty"`
}

// ListAppointmentsRequest query params para listagem
// Aceita datas em formato YYYY-MM-DD ou ISO8601 (com timezone)
// Status pode ser string única ou array (query param repetido: ?status=CREATED&status=CONFIRMED)
type ListAppointmentsRequest struct {
	ProfessionalID string   `query:"professional_id" validate:"omitempty,uuid"`
	CustomerID     string   `query:"customer_id" validate:"omitempty,uuid"`
	Status         []string `query:"status" validate:"omitempty,dive,oneof=CREATED CONFIRMED CHECKED_IN IN_SERVICE AWAITING_PAYMENT DONE NO_SHOW CANCELED"`
	StartDate      string   `query:"start_date" validate:"omitempty"`
	EndDate        string   `query:"end_date" validate:"omitempty"`
	Page           int      `query:"page" validate:"omitempty,min=1"`
	PageSize       int      `query:"page_size" validate:"omitempty,min=1,max=100"`
}

// =============================================================================
// Response DTOs
// =============================================================================

// AppointmentResponse resposta de agendamento
type AppointmentResponse struct {
	ID                    string                       `json:"id"`
	TenantID              string                       `json:"tenant_id"`
	ProfessionalID        string                       `json:"professional_id"`
	ProfessionalName      string                       `json:"professional_name,omitempty"`
	CustomerID            string                       `json:"customer_id"`
	CustomerName          string                       `json:"customer_name,omitempty"`
	CustomerPhone         string                       `json:"customer_phone,omitempty"`
	StartTime             time.Time                    `json:"start_time"`
	EndTime               time.Time                    `json:"end_time"`
	Duration              int                          `json:"duration"`
	CheckedInAt           *time.Time                   `json:"checked_in_at,omitempty"`
	StartedAt             *time.Time                   `json:"started_at,omitempty"`
	FinishedAt            *time.Time                   `json:"finished_at,omitempty"`
	Status                string                       `json:"status"`
	StatusDisplay         string                       `json:"status_display"`
	StatusColor           string                       `json:"status_color"`
	TotalPrice            string                       `json:"total_price"`
	Notes                 string                       `json:"notes,omitempty"`
	CanceledReason        string                       `json:"canceled_reason,omitempty"`
	GoogleCalendarEventID string                       `json:"google_calendar_event_id,omitempty"`
	CommandID             string                       `json:"command_id,omitempty"` // Comanda vinculada (quando status = AWAITING_PAYMENT)
	Services              []AppointmentServiceResponse `json:"services,omitempty"`
	CreatedAt             time.Time                    `json:"created_at"`
	UpdatedAt             time.Time                    `json:"updated_at"`
}

// AppointmentServiceResponse resposta de serviço do agendamento
type AppointmentServiceResponse struct {
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	Price       string `json:"price"`
	Duration    int    `json:"duration"`
}

// ListAppointmentsResponse resposta de listagem paginada
type ListAppointmentsResponse struct {
	Data     []AppointmentResponse `json:"data"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
	Total    int64                 `json:"total"`
}

// AppointmentDailyStatsResponse estatísticas diárias
type AppointmentDailyStatsResponse struct {
	Date              string `json:"date"`
	TotalAppointments int64  `json:"total_appointments"`
	CompletedCount    int64  `json:"completed_count"`
	CanceledCount     int64  `json:"canceled_count"`
	NoShowCount       int64  `json:"no_show_count"`
	TotalRevenue      string `json:"total_revenue"`
}

// =============================================================================
// Calendar-specific DTOs (para integração com FullCalendar)
// =============================================================================

// CalendarEventResponse evento formatado para FullCalendar
type CalendarEventResponse struct {
	ID              string                     `json:"id"`
	Title           string                     `json:"title"`
	Start           string                     `json:"start"`
	End             string                     `json:"end"`
	ResourceID      string                     `json:"resourceId"`
	BackgroundColor string                     `json:"backgroundColor"`
	BorderColor     string                     `json:"borderColor"`
	ExtendedProps   CalendarEventExtendedProps `json:"extendedProps"`
}

// CalendarEventExtendedProps propriedades extras do evento
type CalendarEventExtendedProps struct {
	Status           string   `json:"status"`
	CustomerName     string   `json:"customerName"`
	CustomerPhone    string   `json:"customerPhone"`
	ProfessionalName string   `json:"professionalName"`
	Services         []string `json:"services"`
	TotalPrice       string   `json:"totalPrice"`
	Notes            string   `json:"notes,omitempty"`
}

// CalendarResourceResponse recurso (profissional) para FullCalendar
type CalendarResourceResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Color string `json:"eventColor,omitempty"`
}

// CalendarEventsResponse lista de eventos para calendário
type CalendarEventsResponse struct {
	Events    []CalendarEventResponse    `json:"events"`
	Resources []CalendarResourceResponse `json:"resources"`
}
