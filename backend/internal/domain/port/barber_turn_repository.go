package port

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// BarberTurnRepository define operações para Lista da Vez
type BarberTurnRepository interface {
	// ==========================================================================
	// CREATE / ADD
	// ==========================================================================

	// Add adiciona um barbeiro à lista da vez
	Add(ctx context.Context, barberTurn *entity.BarberTurn) error

	// ==========================================================================
	// READ / LIST
	// ==========================================================================

	// FindByID busca um registro por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.BarberTurn, error)

	// FindByProfessionalID busca por professional_id
	FindByProfessionalID(ctx context.Context, tenantID, professionalID string) (*entity.BarberTurn, error)

	// List lista todos os barbeiros na fila com filtros
	List(ctx context.Context, tenantID string, isActive *bool) ([]*entity.BarberTurn, error)

	// ListActive lista apenas barbeiros ativos na fila
	ListActive(ctx context.Context, tenantID string) ([]*entity.BarberTurn, error)

	// GetNextBarber retorna o próximo barbeiro da fila
	GetNextBarber(ctx context.Context, tenantID string) (*entity.BarberTurn, error)

	// GetStats retorna estatísticas da lista da vez
	GetStats(ctx context.Context, tenantID string) (*entity.BarberTurnStats, error)

	// ==========================================================================
	// UPDATE
	// ==========================================================================

	// RecordTurn registra um atendimento (incrementa pontos)
	RecordTurn(ctx context.Context, tenantID, professionalID string) (*entity.BarberTurn, error)

	// ToggleStatus alterna status ativo/inativo
	ToggleStatus(ctx context.Context, tenantID, professionalID string) (*entity.BarberTurn, error)

	// SetActive ativa um barbeiro
	SetActive(ctx context.Context, tenantID, professionalID string) (*entity.BarberTurn, error)

	// SetInactive pausa um barbeiro
	SetInactive(ctx context.Context, tenantID, professionalID string) (*entity.BarberTurn, error)

	// ==========================================================================
	// DELETE
	// ==========================================================================

	// Remove remove um barbeiro da lista da vez
	Remove(ctx context.Context, tenantID, professionalID string) error

	// ==========================================================================
	// RESET MENSAL
	// ==========================================================================

	// ResetAll zera todos os pontos (reset mensal)
	ResetAll(ctx context.Context, tenantID string) error

	// SaveHistoryBeforeReset salva snapshot no histórico antes do reset
	SaveHistoryBeforeReset(ctx context.Context, tenantID, monthYear string) error

	// ==========================================================================
	// HISTÓRICO
	// ==========================================================================

	// ListHistory lista histórico mensal
	ListHistory(ctx context.Context, tenantID string, monthYear *string) ([]*entity.BarberTurnHistory, error)

	// GetHistoryByMonth busca histórico de um mês específico
	GetHistoryByMonth(ctx context.Context, tenantID, monthYear string) ([]*entity.BarberTurnHistory, error)

	// GetHistorySummary retorna resumo dos últimos 12 meses
	GetHistorySummary(ctx context.Context, tenantID string) ([]*HistorySummary, error)

	// ==========================================================================
	// VALIDAÇÕES
	// ==========================================================================

	// CheckProfessionalInList verifica se profissional já está na lista
	CheckProfessionalInList(ctx context.Context, tenantID, professionalID string) (bool, error)

	// CheckProfessionalIsBarber verifica se profissional é barbeiro ativo
	CheckProfessionalIsBarber(ctx context.Context, tenantID, professionalID string) (bool, error)

	// GetAvailableBarbers lista barbeiros disponíveis para adicionar
	GetAvailableBarbers(ctx context.Context, tenantID string) ([]*AvailableBarber, error)

	// ==========================================================================
	// RELATÓRIOS
	// ==========================================================================

	// GetTodayStats retorna estatísticas do dia
	GetTodayStats(ctx context.Context, tenantID string) (*TodayStats, error)
}

// =============================================================================
// Tipos auxiliares
// =============================================================================

// HistorySummary resumo mensal
type HistorySummary struct {
	MonthYear         string
	TotalBarbeiros    int64
	TotalAtendimentos int64
	MediaAtendimentos float64
}

// AvailableBarber barbeiro disponível para adicionar
type AvailableBarber struct {
	ID     string
	Nome   string
	Foto   *string
	Status string
}

// TodayStats estatísticas do dia
type TodayStats struct {
	AtendimentosHoje  int64
	TotalPontosMes    int64
	BarbeirosAtivos   int64
	UltimoAtendimento *time.Time
}
