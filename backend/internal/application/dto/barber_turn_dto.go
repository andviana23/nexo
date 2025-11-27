package dto

import "time"

// =============================================================================
// DTOs para Lista da Vez (Barber Turn)
// Conforme FLUXO_LISTA_DA_VEZ.md
// =============================================================================

// =============================================================================
// Request DTOs
// =============================================================================

// AddBarberToTurnListRequest requisição para adicionar barbeiro à lista
type AddBarberToTurnListRequest struct {
	ProfessionalID string `json:"professional_id" validate:"required,uuid"`
}

// RecordTurnRequest requisição para registrar atendimento (incrementar pontos)
type RecordTurnRequest struct {
	ProfessionalID string `json:"professional_id" validate:"required,uuid"`
}

// ToggleBarberStatusRequest requisição para pausar/ativar barbeiro
type ToggleBarberStatusRequest struct {
	ProfessionalID string `json:"professional_id" validate:"required,uuid"`
}

// ListBarbersTurnRequest query params para listagem
type ListBarbersTurnRequest struct {
	IsActive *bool `query:"is_active"`
}

// GetTurnHistoryRequest query params para histórico
type GetTurnHistoryRequest struct {
	MonthYear string `query:"month_year" validate:"omitempty,len=7"` // YYYY-MM
}

// ResetTurnListRequest requisição para reset manual (admin)
type ResetTurnListRequest struct {
	SaveHistory bool `json:"save_history"` // Se true, salva snapshot antes de resetar
}

// =============================================================================
// Response DTOs
// =============================================================================

// BarberTurnResponse resposta de barbeiro na lista da vez
type BarberTurnResponse struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id"`
	ProfessionalID    string     `json:"professional_id"`
	ProfessionalName  string     `json:"professional_name"`
	ProfessionalType  string     `json:"professional_type"`
	ProfessionalPhoto *string    `json:"professional_photo,omitempty"`
	CurrentPoints     int        `json:"current_points"`
	LastTurnAt        *time.Time `json:"last_turn_at,omitempty"`
	IsActive          bool       `json:"is_active"`
	Position          int64      `json:"position"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// NextBarberResponse resposta do próximo barbeiro da vez
type NextBarberResponse struct {
	ProfessionalID    string  `json:"professional_id"`
	ProfessionalName  string  `json:"professional_name"`
	ProfessionalPhoto *string `json:"professional_photo,omitempty"`
	CurrentPoints     int     `json:"current_points"`
}

// ListBarbersTurnResponse resposta da lista de barbeiros na fila
type ListBarbersTurnResponse struct {
	Barbers    []BarberTurnResponse    `json:"barbers"`
	Total      int                     `json:"total"`
	NextBarber *NextBarberResponse     `json:"next_barber,omitempty"`
	Stats      BarberTurnStatsResponse `json:"stats"`
}

// BarberTurnStatsResponse estatísticas da lista da vez
type BarberTurnStatsResponse struct {
	TotalAtivos       int64      `json:"total_ativos"`
	TotalPausados     int64      `json:"total_pausados"`
	TotalGeral        int64      `json:"total_geral"`
	TotalPontosMes    int64      `json:"total_pontos_mes"`
	AtendimentosHoje  int64      `json:"atendimentos_hoje,omitempty"`
	UltimoAtendimento *time.Time `json:"ultimo_atendimento,omitempty"`
}

// RecordTurnResponse resposta após registrar atendimento
type RecordTurnResponse struct {
	ProfessionalID   string    `json:"professional_id"`
	ProfessionalName string    `json:"professional_name"`
	PreviousPoints   int       `json:"previous_points"`
	NewPoints        int       `json:"new_points"`
	LastTurnAt       time.Time `json:"last_turn_at"`
	Message          string    `json:"message"`
}

// ToggleStatusResponse resposta após pausar/ativar barbeiro
type ToggleStatusResponse struct {
	ProfessionalID   string `json:"professional_id"`
	ProfessionalName string `json:"professional_name"`
	IsActive         bool   `json:"is_active"`
	Message          string `json:"message"`
}

// RemoveBarberResponse resposta após remover barbeiro da lista
type RemoveBarberResponse struct {
	ProfessionalID string `json:"professional_id"`
	Message        string `json:"message"`
}

// ResetTurnListResponse resposta após reset mensal
type ResetTurnListResponse struct {
	Message  string             `json:"message"`
	Snapshot *TurnResetSnapshot `json:"snapshot,omitempty"`
}

// TurnResetSnapshot snapshot do reset
type TurnResetSnapshot struct {
	MonthYear             string `json:"month_year"`
	TotalBarbers          int    `json:"total_barbers"`
	TotalPointsReset      int64  `json:"total_points_reset"`
	HistoryRecordsCreated int    `json:"history_records_created"`
}

// =============================================================================
// Histórico DTOs
// =============================================================================

// TurnHistoryResponse histórico de atendimentos mensais
type TurnHistoryResponse struct {
	ID               string    `json:"id"`
	TenantID         string    `json:"tenant_id"`
	ProfessionalID   string    `json:"professional_id"`
	ProfessionalName string    `json:"professional_name"`
	MonthYear        string    `json:"month_year"`
	TotalTurns       int       `json:"total_turns"`
	FinalPoints      int       `json:"final_points"`
	CreatedAt        time.Time `json:"created_at"`
}

// ListTurnHistoryResponse resposta da listagem de histórico
type ListTurnHistoryResponse struct {
	History []TurnHistoryResponse `json:"history"`
	Total   int                   `json:"total"`
}

// TurnHistorySummaryResponse resumo do histórico mensal
type TurnHistorySummaryResponse struct {
	MonthYear         string  `json:"month_year"`
	TotalBarbeiros    int64   `json:"total_barbeiros"`
	TotalAtendimentos int64   `json:"total_atendimentos"`
	MediaAtendimentos float64 `json:"media_atendimentos"`
}

// ListHistorySummaryResponse lista de resumos mensais
type ListHistorySummaryResponse struct {
	Summary []TurnHistorySummaryResponse `json:"summary"`
}

// =============================================================================
// Relatórios DTOs
// =============================================================================

// DailyReportRequest query params para relatório diário
type DailyReportRequest struct {
	Date string `query:"date" validate:"omitempty"` // YYYY-MM-DD, default: hoje
}

// DailyReportResponse relatório diário
type DailyReportResponse struct {
	Date              string              `json:"date"`
	TotalAtendimentos int                 `json:"total_atendimentos"`
	Barbers           []DailyBarberReport `json:"barbers"`
}

// DailyBarberReport dados do barbeiro no relatório diário
type DailyBarberReport struct {
	ProfessionalID   string     `json:"professional_id"`
	ProfessionalName string     `json:"professional_name"`
	CurrentPoints    int        `json:"current_points"`
	AttendedToday    bool       `json:"attended_today"`
	LastTurnAt       *time.Time `json:"last_turn_at,omitempty"`
}

// =============================================================================
// Barbeiros Disponíveis para adicionar à Lista
// =============================================================================

// AvailableBarberResponse barbeiro disponível para adicionar
type AvailableBarberResponse struct {
	ID     string  `json:"id"`
	Nome   string  `json:"nome"`
	Foto   *string `json:"foto,omitempty"`
	Status string  `json:"status"`
}

// ListAvailableBarbersResponse lista de barbeiros disponíveis
type ListAvailableBarbersResponse struct {
	Barbers []AvailableBarberResponse `json:"barbers"`
	Total   int                       `json:"total"`
}
