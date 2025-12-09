package mapper

import (
	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// AppointmentToResponse converte entidade para DTO de resposta
func AppointmentToResponse(a *entity.Appointment) dto.AppointmentResponse {
	services := make([]dto.AppointmentServiceResponse, 0, len(a.Services))
	for _, svc := range a.Services {
		services = append(services, dto.AppointmentServiceResponse{
			ServiceID:   svc.ServiceID,
			ServiceName: svc.ServiceName,
			Price:       svc.PriceAtBooking.Raw(), // Retorna "50.00" em vez de "R$ 50,00" para evitar NaN no frontend
			Duration:    svc.DurationAtBooking,
		})
	}

	return dto.AppointmentResponse{
		ID:                    a.ID,
		TenantID:              a.TenantID.String(),
		ProfessionalID:        a.ProfessionalID,
		ProfessionalName:      a.ProfessionalName,
		CustomerID:            a.CustomerID,
		CustomerName:          a.CustomerName,
		CustomerPhone:         a.CustomerPhone,
		StartTime:             a.StartTime,
		EndTime:               a.EndTime,
		Duration:              a.Duration(),
		CheckedInAt:           a.CheckedInAt,
		StartedAt:             a.StartedAt,
		FinishedAt:            a.FinishedAt,
		Status:                a.Status.String(),
		StatusDisplay:         a.Status.DisplayName(),
		StatusColor:           a.Status.Color(),
		TotalPrice:            a.TotalPrice.Raw(), // Retorna "50.00" em vez de "R$ 50,00" para evitar NaN no frontend
		Notes:                 a.Notes,
		CanceledReason:        a.CanceledReason,
		GoogleCalendarEventID: a.GoogleCalendarEventID,
		CommandID:             a.CommandID, // Campo command_id adicionado
		Services:              services,
		CreatedAt:             a.CreatedAt,
		UpdatedAt:             a.UpdatedAt,
	}
}

// AppointmentsToResponse converte lista de entidades para lista de DTOs
func AppointmentsToResponse(appointments []*entity.Appointment) []dto.AppointmentResponse {
	result := make([]dto.AppointmentResponse, 0, len(appointments))
	for _, a := range appointments {
		result = append(result, AppointmentToResponse(a))
	}
	return result
}

// AppointmentToCalendarEvent converte agendamento para evento de calendário
func AppointmentToCalendarEvent(a *entity.Appointment) dto.CalendarEventResponse {
	// Extrair nomes dos serviços
	serviceNames := make([]string, 0, len(a.Services))
	for _, svc := range a.Services {
		serviceNames = append(serviceNames, svc.ServiceName)
	}

	return dto.CalendarEventResponse{
		ID:              a.ID,
		Title:           buildEventTitle(a),
		Start:           a.StartTime.Format("2006-01-02T15:04:05"),
		End:             a.EndTime.Format("2006-01-02T15:04:05"),
		ResourceID:      a.ProfessionalID,
		BackgroundColor: a.Status.Color(),
		BorderColor:     a.Status.Color(),
		ExtendedProps: dto.CalendarEventExtendedProps{
			Status:           a.Status.String(),
			CustomerName:     a.CustomerName,
			CustomerPhone:    a.CustomerPhone,
			ProfessionalName: a.ProfessionalName,
			Services:         serviceNames,
			TotalPrice:       a.TotalPrice.Raw(), // Retorna "50.00" em vez de "R$ 50,00"
			Notes:            a.Notes,
		},
	}
}

// AppointmentsToCalendarEvents converte lista para eventos de calendário
func AppointmentsToCalendarEvents(appointments []*entity.Appointment) []dto.CalendarEventResponse {
	events := make([]dto.CalendarEventResponse, 0, len(appointments))
	for _, a := range appointments {
		events = append(events, AppointmentToCalendarEvent(a))
	}
	return events
}

// buildEventTitle constrói o título do evento
func buildEventTitle(a *entity.Appointment) string {
	if len(a.Services) == 0 {
		return "Agendamento"
	}
	if len(a.Services) == 1 {
		return a.Services[0].ServiceName
	}
	return a.Services[0].ServiceName + " +"
}
