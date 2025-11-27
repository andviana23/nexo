package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/google/uuid"
)

// BarberTurn representa um profissional na lista da vez
type BarberTurn struct {
	ID             string
	TenantID       string
	ProfessionalID string

	// Estado atual na fila
	CurrentPoints int        // Pontuação acumulada no mês
	LastTurnAt    *time.Time // Último atendimento registrado
	IsActive      bool       // Se está participando ativamente da fila

	// Metadados
	CreatedAt time.Time
	UpdatedAt time.Time

	// Dados do profissional (preenchidos via JOIN)
	ProfessionalName   string
	ProfessionalType   string
	ProfessionalStatus string
	ProfessionalPhoto  *string

	// Posição calculada na fila (preenchida na query)
	Position int64
}

// NewBarberTurn cria um novo registro de barbeiro na lista da vez
func NewBarberTurn(tenantID, professionalID string) (*BarberTurn, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if professionalID == "" {
		return nil, domain.ErrBarberTurnProfessionalRequired
	}

	now := time.Now()
	return &BarberTurn{
		ID:             uuid.NewString(),
		TenantID:       tenantID,
		ProfessionalID: professionalID,
		CurrentPoints:  0,
		LastTurnAt:     nil,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// RecordTurn registra um atendimento (incrementa pontos)
func (b *BarberTurn) RecordTurn() error {
	if !b.IsActive {
		return domain.ErrBarberTurnCannotRecord
	}

	b.CurrentPoints++
	now := time.Now()
	b.LastTurnAt = &now
	b.UpdatedAt = now
	return nil
}

// ToggleActive alterna o status ativo/inativo
func (b *BarberTurn) ToggleActive() {
	b.IsActive = !b.IsActive
	b.UpdatedAt = time.Now()
}

// Pause pausa o barbeiro na fila
func (b *BarberTurn) Pause() {
	b.IsActive = false
	b.UpdatedAt = time.Now()
}

// Activate ativa o barbeiro na fila
func (b *BarberTurn) Activate() {
	b.IsActive = true
	b.UpdatedAt = time.Now()
}

// Reset zera os pontos e timestamps (usado no reset mensal)
func (b *BarberTurn) Reset() {
	b.CurrentPoints = 0
	b.LastTurnAt = nil
	b.UpdatedAt = time.Now()
}

// Validate valida as regras de negócio
func (b *BarberTurn) Validate() error {
	if b.TenantID == "" {
		return domain.ErrTenantIDRequired
	}
	if b.ProfessionalID == "" {
		return domain.ErrBarberTurnProfessionalRequired
	}
	if b.CurrentPoints < 0 {
		return domain.ErrBarberTurnInvalidPoints
	}
	return nil
}

// IsNext verifica se é o próximo da fila (posição 1)
func (b *BarberTurn) IsNext() bool {
	return b.Position == 1 && b.IsActive
}

// CanBeSelected verifica se pode ser selecionado para atendimento
func (b *BarberTurn) CanBeSelected() bool {
	return b.IsActive && b.ProfessionalStatus == "ATIVO"
}

// =============================================================================
// BarberTurnHistory representa o histórico mensal de um barbeiro
// =============================================================================

// BarberTurnHistory histórico mensal de atendimentos
type BarberTurnHistory struct {
	ID             string
	TenantID       string
	ProfessionalID string
	MonthYear      string // Formato: "YYYY-MM"
	TotalTurns     int
	FinalPoints    int
	CreatedAt      time.Time

	// Dados do profissional (preenchidos via JOIN)
	ProfessionalName string
}

// NewBarberTurnHistory cria um novo registro de histórico
func NewBarberTurnHistory(
	tenantID, professionalID, monthYear string,
	totalTurns, finalPoints int,
) (*BarberTurnHistory, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if professionalID == "" {
		return nil, domain.ErrBarberTurnProfessionalRequired
	}
	if len(monthYear) != 7 {
		return nil, domain.ErrBarberTurnMonthYearInvalid
	}

	return &BarberTurnHistory{
		ID:             uuid.NewString(),
		TenantID:       tenantID,
		ProfessionalID: professionalID,
		MonthYear:      monthYear,
		TotalTurns:     totalTurns,
		FinalPoints:    finalPoints,
		CreatedAt:      time.Now(),
	}, nil
}

// =============================================================================
// BarberTurnStats estatísticas da lista da vez
// =============================================================================

// BarberTurnStats estatísticas consolidadas
type BarberTurnStats struct {
	TotalAtivos       int64
	TotalPausados     int64
	TotalGeral        int64
	TotalPontosMes    int64
	AtendimentosHoje  int64
	UltimoAtendimento *time.Time
}

// =============================================================================
// TurnResetResult resultado do reset mensal
// =============================================================================

// TurnResetResult dados do resultado do reset
type TurnResetResult struct {
	MonthYear             string
	TotalBarbers          int
	TotalPointsReset      int64
	HistoryRecordsCreated int
}
