package dto

import "time"

// CreateBlockedTimeRequest representa a requisição para criar um bloqueio de horário
type CreateBlockedTimeRequest struct {
	ProfessionalID string    `json:"professional_id" validate:"required,uuid"`
	StartTime      time.Time `json:"start_time" validate:"required"`
	EndTime        time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	Reason         string    `json:"reason" validate:"required,min=3,max=255"`
	IsRecurring    bool      `json:"is_recurring"`
	RecurrenceRule *string   `json:"recurrence_rule,omitempty"`
}

// UpdateBlockedTimeRequest representa a requisição para atualizar um bloqueio
type UpdateBlockedTimeRequest struct {
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	Reason    string    `json:"reason" validate:"required,min=3,max=255"`
}

// ListBlockedTimesRequest representa os filtros para listar bloqueios
type ListBlockedTimesRequest struct {
	ProfessionalID *string    `query:"professional_id" validate:"omitempty,uuid"`
	StartDate      *time.Time `query:"start_date"`
	EndDate        *time.Time `query:"end_date"`
}

// BlockedTimeResponse representa um bloqueio de horário na resposta
type BlockedTimeResponse struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	ProfessionalID string    `json:"professional_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Reason         string    `json:"reason"`
	IsRecurring    bool      `json:"is_recurring"`
	RecurrenceRule *string   `json:"recurrence_rule,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedBy      *string   `json:"created_by,omitempty"`
}

// ListBlockedTimesResponse representa a resposta com lista de bloqueios
type ListBlockedTimesResponse struct {
	Data  []BlockedTimeResponse `json:"data"`
	Total int                   `json:"total"`
}
