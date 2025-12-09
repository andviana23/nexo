// Package postgres contém implementações de repositórios usando PostgreSQL.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
)

// BarberTurnRepository implementa BarberTurnRepository usando PostgreSQL
type BarberTurnRepository struct {
	queries *db.Queries
}

// NewBarberTurnRepository cria uma nova instância
func NewBarberTurnRepository(queries *db.Queries) *BarberTurnRepository {
	return &BarberTurnRepository{queries: queries}
}

// =============================================================================
// CREATE / ADD
// =============================================================================

// Add adiciona um barbeiro à lista da vez
func (r *BarberTurnRepository) Add(ctx context.Context, barberTurn *entity.BarberTurn) error {
	params := db.AddBarberToTurnListParams{
		ID:             stringToUUID(barberTurn.ID),
		TenantID:       uuidToPgUUID(barberTurn.TenantID),
		ProfessionalID: stringToUUID(barberTurn.ProfessionalID),
	}

	_, err := r.queries.AddBarberToTurnList(ctx, params)
	return err
}

// =============================================================================
// READ / LIST
// =============================================================================

// FindByID busca um registro por ID
func (r *BarberTurnRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.BarberTurn, error) {
	params := db.GetBarberTurnByIDParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
	}

	row, err := r.queries.GetBarberTurnByID(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapBarberTurnByIDRowToEntity(row), nil
}

// FindByProfessionalID busca por professional_id
func (r *BarberTurnRepository) FindByProfessionalID(ctx context.Context, tenantID, professionalID string) (*entity.BarberTurn, error) {
	params := db.GetBarberTurnByProfessionalIDParams{
		ProfessionalID: stringToUUID(professionalID),
		TenantID:       stringToUUID(tenantID),
	}

	row, err := r.queries.GetBarberTurnByProfessionalID(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapBarberTurnByProfessionalIDRowToEntity(row), nil
}

// List lista todos os barbeiros na fila com filtros
func (r *BarberTurnRepository) List(ctx context.Context, tenantID string, isActive *bool) ([]*entity.BarberTurn, error) {
	params := db.ListBarbersTurnListParams{
		TenantID: stringToUUID(tenantID),
		IsActive: isActive,
	}

	fmt.Printf("🔍 [REPO DEBUG] List() chamado - tenantID: %s, isActive: %v\n", tenantID, isActive)

	rows, err := r.queries.ListBarbersTurnList(ctx, params)
	if err != nil {
		fmt.Printf("❌ [REPO DEBUG] Erro ao executar query: %v\n", err)
		return nil, err
	}

	fmt.Printf("🔍 [REPO DEBUG] Query retornou %d rows\n", len(rows))

	result := make([]*entity.BarberTurn, 0, len(rows))
	for i, row := range rows {
		entity := mapBarberTurnListRowToEntity(row)
		fmt.Printf("🔍 [REPO DEBUG] Row %d mapeada - ID: %s, Nome: %s\n", i, entity.ID, entity.ProfessionalName)
		result = append(result, entity)
	}

	fmt.Printf("🔍 [REPO DEBUG] Retornando %d barbeiros\n", len(result))
	return result, nil
}

// ListActive lista apenas barbeiros ativos na fila
func (r *BarberTurnRepository) ListActive(ctx context.Context, tenantID string) ([]*entity.BarberTurn, error) {
	rows, err := r.queries.ListActiveBarbersTurnList(ctx, stringToUUID(tenantID))
	if err != nil {
		return nil, err
	}

	result := make([]*entity.BarberTurn, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapBarberTurnActiveRowToEntity(row))
	}

	return result, nil
}

// GetNextBarber retorna o próximo barbeiro da fila
func (r *BarberTurnRepository) GetNextBarber(ctx context.Context, tenantID string) (*entity.BarberTurn, error) {
	row, err := r.queries.GetNextBarber(ctx, stringToUUID(tenantID))
	if err != nil {
		return nil, err
	}

	return mapNextBarberRowToEntity(row), nil
}

// GetStats retorna estatísticas da lista da vez
func (r *BarberTurnRepository) GetStats(ctx context.Context, tenantID string) (*entity.BarberTurnStats, error) {
	row, err := r.queries.CountBarbersTurnList(ctx, stringToUUID(tenantID))
	if err != nil {
		return nil, err
	}

	var totalPontos int64
	if v, ok := row.TotalPontos.(int64); ok {
		totalPontos = v
	}

	return &entity.BarberTurnStats{
		TotalAtivos:    row.TotalAtivos,
		TotalPausados:  row.TotalPausados,
		TotalGeral:     row.TotalGeral,
		TotalPontosMes: totalPontos,
	}, nil
}

// =============================================================================
// UPDATE
// =============================================================================

// RecordTurn registra um atendimento (incrementa pontos)
func (r *BarberTurnRepository) RecordTurn(ctx context.Context, tenantID, professionalID string) (*entity.BarberTurn, error) {
	params := db.RecordTurnParams{
		ProfessionalID: stringToUUID(professionalID),
		TenantID:       stringToUUID(tenantID),
	}

	row, err := r.queries.RecordTurn(ctx, params)
	if err != nil {
		return nil, err
	}

	// Usar dados do row ou buscar dados completos com JOIN
	_ = row // evitar erro de não uso

	// Busca dados completos com JOIN
	return r.FindByProfessionalID(ctx, tenantID, professionalID)
}

// ToggleStatus alterna status ativo/inativo
func (r *BarberTurnRepository) ToggleStatus(ctx context.Context, tenantID, professionalID string) (*entity.BarberTurn, error) {
	params := db.ToggleBarberTurnStatusParams{
		ProfessionalID: stringToUUID(professionalID),
		TenantID:       stringToUUID(tenantID),
	}

	_, err := r.queries.ToggleBarberTurnStatus(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.FindByProfessionalID(ctx, tenantID, professionalID)
}

// SetActive ativa um barbeiro
func (r *BarberTurnRepository) SetActive(ctx context.Context, tenantID, professionalID string) (*entity.BarberTurn, error) {
	params := db.SetBarberTurnActiveParams{
		ProfessionalID: stringToUUID(professionalID),
		TenantID:       stringToUUID(tenantID),
	}

	_, err := r.queries.SetBarberTurnActive(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.FindByProfessionalID(ctx, tenantID, professionalID)
}

// SetInactive pausa um barbeiro
func (r *BarberTurnRepository) SetInactive(ctx context.Context, tenantID, professionalID string) (*entity.BarberTurn, error) {
	params := db.SetBarberTurnInactiveParams{
		ProfessionalID: stringToUUID(professionalID),
		TenantID:       stringToUUID(tenantID),
	}

	_, err := r.queries.SetBarberTurnInactive(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.FindByProfessionalID(ctx, tenantID, professionalID)
}

// =============================================================================
// DELETE
// =============================================================================

// Remove remove um barbeiro da lista da vez
func (r *BarberTurnRepository) Remove(ctx context.Context, tenantID, professionalID string) error {
	params := db.RemoveBarberFromTurnListParams{
		ProfessionalID: stringToUUID(professionalID),
		TenantID:       stringToUUID(tenantID),
	}

	return r.queries.RemoveBarberFromTurnList(ctx, params)
}

// =============================================================================
// RESET MENSAL
// =============================================================================

// ResetAll zera todos os pontos (reset mensal)
func (r *BarberTurnRepository) ResetAll(ctx context.Context, tenantID string) error {
	return r.queries.ResetAllTurnPoints(ctx, stringToUUID(tenantID))
}

// SaveHistoryBeforeReset salva snapshot no histórico antes do reset
func (r *BarberTurnRepository) SaveHistoryBeforeReset(ctx context.Context, tenantID, monthYear string) error {
	params := db.SaveTurnHistoryBeforeResetParams{
		TenantID: stringToUUID(tenantID),
		Column2:  monthYear,
	}

	return r.queries.SaveTurnHistoryBeforeReset(ctx, params)
}

// =============================================================================
// HISTÓRICO
// =============================================================================

// ListHistory lista histórico mensal
func (r *BarberTurnRepository) ListHistory(ctx context.Context, tenantID string, monthYear *string) ([]*entity.BarberTurnHistory, error) {
	var my string
	if monthYear != nil {
		my = *monthYear
	}

	params := db.ListTurnHistoryParams{
		TenantID: stringToUUID(tenantID),
		Column2:  my,
	}

	rows, err := r.queries.ListTurnHistory(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.BarberTurnHistory, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapTurnHistoryRowToEntity(row))
	}

	return result, nil
}

// GetHistoryByMonth busca histórico de um mês específico
func (r *BarberTurnRepository) GetHistoryByMonth(ctx context.Context, tenantID, monthYear string) ([]*entity.BarberTurnHistory, error) {
	params := db.GetTurnHistoryByMonthParams{
		TenantID:  stringToUUID(tenantID),
		MonthYear: monthYear,
	}

	rows, err := r.queries.GetTurnHistoryByMonth(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.BarberTurnHistory, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapTurnHistoryByMonthRowToEntity(row))
	}

	return result, nil
}

// GetHistorySummary retorna resumo dos últimos 12 meses
func (r *BarberTurnRepository) GetHistorySummary(ctx context.Context, tenantID string) ([]*port.HistorySummary, error) {
	rows, err := r.queries.GetTurnHistorySummary(ctx, stringToUUID(tenantID))
	if err != nil {
		return nil, err
	}

	result := make([]*port.HistorySummary, 0, len(rows))
	for _, row := range rows {
		result = append(result, &port.HistorySummary{
			MonthYear:         row.MonthYear,
			TotalBarbeiros:    row.TotalBarbeiros,
			TotalAtendimentos: row.TotalAtendimentos,
			MediaAtendimentos: row.MediaAtendimentos,
		})
	}

	return result, nil
}

// =============================================================================
// VALIDAÇÕES
// =============================================================================

// CheckProfessionalInList verifica se profissional já está na lista
func (r *BarberTurnRepository) CheckProfessionalInList(ctx context.Context, tenantID, professionalID string) (bool, error) {
	params := db.CheckProfessionalInTurnListParams{
		TenantID:       stringToUUID(tenantID),
		ProfessionalID: stringToUUID(professionalID),
	}

	return r.queries.CheckProfessionalInTurnList(ctx, params)
}

// CheckProfessionalIsBarber verifica se profissional é barbeiro ativo
func (r *BarberTurnRepository) CheckProfessionalIsBarber(ctx context.Context, tenantID, professionalID string) (bool, error) {
	params := db.CheckProfessionalIsBarberParams{
		ID:       stringToUUID(professionalID),
		TenantID: stringToUUID(tenantID),
	}

	return r.queries.CheckProfessionalIsBarber(ctx, params)
}

// GetAvailableBarbers lista barbeiros disponíveis para adicionar
func (r *BarberTurnRepository) GetAvailableBarbers(ctx context.Context, tenantID string) ([]*port.AvailableBarber, error) {
	rows, err := r.queries.GetAvailableBarbersForTurnList(ctx, stringToUUID(tenantID))
	if err != nil {
		return nil, err
	}

	result := make([]*port.AvailableBarber, 0, len(rows))
	for _, row := range rows {
		var status string
		if row.Status != nil {
			status = *row.Status
		}

		result = append(result, &port.AvailableBarber{
			ID:     uuidToString(row.ID),
			Nome:   row.Nome,
			Foto:   row.Foto,
			Status: status,
		})
	}

	return result, nil
}

// =============================================================================
// RELATÓRIOS
// =============================================================================

// GetTodayStats retorna estatísticas do dia
func (r *BarberTurnRepository) GetTodayStats(ctx context.Context, tenantID string) (*port.TodayStats, error) {
	row, err := r.queries.GetTodayStats(ctx, stringToUUID(tenantID))
	if err != nil {
		return nil, err
	}

	result := &port.TodayStats{
		AtendimentosHoje: row.AtendimentosHoje,
		TotalPontosMes:   row.TotalPontosMes,
		BarbeirosAtivos:  row.BarbeirosAtivos,
	}

	// Parse último atendimento
	if row.UltimoAtendimento != nil {
		if t, ok := row.UltimoAtendimento.(time.Time); ok {
			result.UltimoAtendimento = &t
		}
	}

	return result, nil
}

// =============================================================================
// Mappers
// =============================================================================

func mapBarberTurnByIDRowToEntity(row db.GetBarberTurnByIDRow) *entity.BarberTurn {
	var status string
	if row.ProfessionalStatus != nil {
		status = *row.ProfessionalStatus
	}

	return &entity.BarberTurn{
		ID:                 uuidToString(row.ID),
		TenantID:           pgtypeToEntityUUID(row.TenantID),
		ProfessionalID:     uuidToString(row.ProfessionalID),
		CurrentPoints:      int(row.CurrentPoints),
		LastTurnAt:         pgTimestamptzToTimePtr(row.LastTurnAt),
		IsActive:           row.IsActive,
		CreatedAt:          pgTimestamptzToTime(row.CreatedAt),
		UpdatedAt:          pgTimestamptzToTime(row.UpdatedAt),
		ProfessionalName:   row.ProfessionalName,
		ProfessionalType:   row.ProfessionalType,
		ProfessionalStatus: status,
	}
}

func mapBarberTurnByProfessionalIDRowToEntity(row db.GetBarberTurnByProfessionalIDRow) *entity.BarberTurn {
	var status string
	if row.ProfessionalStatus != nil {
		status = *row.ProfessionalStatus
	}

	return &entity.BarberTurn{
		ID:                 uuidToString(row.ID),
		TenantID:           pgtypeToEntityUUID(row.TenantID),
		ProfessionalID:     uuidToString(row.ProfessionalID),
		CurrentPoints:      int(row.CurrentPoints),
		LastTurnAt:         pgTimestamptzToTimePtr(row.LastTurnAt),
		IsActive:           row.IsActive,
		CreatedAt:          pgTimestamptzToTime(row.CreatedAt),
		UpdatedAt:          pgTimestamptzToTime(row.UpdatedAt),
		ProfessionalName:   row.ProfessionalName,
		ProfessionalType:   row.ProfessionalType,
		ProfessionalStatus: status,
	}
}

func mapBarberTurnListRowToEntity(row db.ListBarbersTurnListRow) *entity.BarberTurn {
	var status string
	if row.ProfessionalStatus != nil {
		status = *row.ProfessionalStatus
	}

	return &entity.BarberTurn{
		ID:                 uuidToString(row.ID),
		TenantID:           pgtypeToEntityUUID(row.TenantID),
		ProfessionalID:     uuidToString(row.ProfessionalID),
		CurrentPoints:      int(row.CurrentPoints),
		LastTurnAt:         pgTimestamptzToTimePtr(row.LastTurnAt),
		IsActive:           row.IsActive,
		CreatedAt:          pgTimestamptzToTime(row.CreatedAt),
		UpdatedAt:          pgTimestamptzToTime(row.UpdatedAt),
		ProfessionalName:   row.ProfessionalName,
		ProfessionalType:   row.ProfessionalType,
		ProfessionalStatus: status,
		ProfessionalPhoto:  row.ProfessionalPhoto,
		Position:           row.Position,
	}
}

func mapBarberTurnActiveRowToEntity(row db.ListActiveBarbersTurnListRow) *entity.BarberTurn {
	var status string
	if row.ProfessionalStatus != nil {
		status = *row.ProfessionalStatus
	}

	return &entity.BarberTurn{
		ID:                 uuidToString(row.ID),
		TenantID:           pgtypeToEntityUUID(row.TenantID),
		ProfessionalID:     uuidToString(row.ProfessionalID),
		CurrentPoints:      int(row.CurrentPoints),
		LastTurnAt:         pgTimestamptzToTimePtr(row.LastTurnAt),
		IsActive:           row.IsActive,
		CreatedAt:          pgTimestamptzToTime(row.CreatedAt),
		UpdatedAt:          pgTimestamptzToTime(row.UpdatedAt),
		ProfessionalName:   row.ProfessionalName,
		ProfessionalType:   row.ProfessionalType,
		ProfessionalStatus: status,
		ProfessionalPhoto:  row.ProfessionalPhoto,
		Position:           row.Position,
	}
}

func mapNextBarberRowToEntity(row db.GetNextBarberRow) *entity.BarberTurn {
	var status string
	if row.ProfessionalStatus != nil {
		status = *row.ProfessionalStatus
	}

	return &entity.BarberTurn{
		ID:                 uuidToString(row.ID),
		TenantID:           pgtypeToEntityUUID(row.TenantID),
		ProfessionalID:     uuidToString(row.ProfessionalID),
		CurrentPoints:      int(row.CurrentPoints),
		LastTurnAt:         pgTimestamptzToTimePtr(row.LastTurnAt),
		IsActive:           row.IsActive,
		CreatedAt:          pgTimestamptzToTime(row.CreatedAt),
		UpdatedAt:          pgTimestamptzToTime(row.UpdatedAt),
		ProfessionalName:   row.ProfessionalName,
		ProfessionalType:   row.ProfessionalType,
		ProfessionalStatus: status,
		ProfessionalPhoto:  row.ProfessionalPhoto,
		Position:           1, // Sempre 1 pois é o próximo
	}
}

func mapTurnHistoryRowToEntity(row db.ListTurnHistoryRow) *entity.BarberTurnHistory {
	return &entity.BarberTurnHistory{
		ID:               uuidToString(row.ID),
		TenantID:         pgtypeToEntityUUID(row.TenantID),
		ProfessionalID:   uuidToString(row.ProfessionalID),
		MonthYear:        row.MonthYear,
		TotalTurns:       int(row.TotalTurns),
		FinalPoints:      int(row.FinalPoints),
		CreatedAt:        pgTimestamptzToTime(row.CreatedAt),
		ProfessionalName: row.ProfessionalName,
	}
}

func mapTurnHistoryByMonthRowToEntity(row db.GetTurnHistoryByMonthRow) *entity.BarberTurnHistory {
	return &entity.BarberTurnHistory{
		ID:               uuidToString(row.ID),
		TenantID:         pgtypeToEntityUUID(row.TenantID),
		ProfessionalID:   uuidToString(row.ProfessionalID),
		MonthYear:        row.MonthYear,
		TotalTurns:       int(row.TotalTurns),
		FinalPoints:      int(row.FinalPoints),
		CreatedAt:        pgTimestamptzToTime(row.CreatedAt),
		ProfessionalName: row.ProfessionalName,
	}
}

// Verifica se BarberTurnRepository implementa a interface
var _ port.BarberTurnRepository = (*BarberTurnRepository)(nil)
