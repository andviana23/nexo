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
			Price:       svc.PriceAtBooking.String(),
			Duration:    svc.DurationAtBooking,
		})
	}

	return dto.AppointmentResponse{
		ID:                    a.ID,
		TenantID:              a.TenantID,
		ProfessionalID:        a.ProfessionalID,
		CustomerID:            a.CustomerID,
		StartTime:             a.StartTime,
		EndTime:               a.EndTime,
		Duration:              a.Duration(),
		Status:                a.Status.String(),
		StatusDisplay:         a.Status.DisplayName(),
		StatusColor:           a.Status.Color(),
		TotalPrice:            a.TotalPrice.String(),
		Notes:                 a.Notes,
		CanceledReason:        a.CanceledReason,
		GoogleCalendarEventID: a.GoogleCalendarEventID,
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
			CustomerName:     "", // Precisa vir do join
			CustomerPhone:    "", // Precisa vir do join
			ProfessionalName: "", // Precisa vir do join
			Services:         serviceNames,
			TotalPrice:       a.TotalPrice.String(),
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
