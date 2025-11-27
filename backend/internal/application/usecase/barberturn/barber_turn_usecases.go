// Package barberturn contém os use cases da Lista da Vez do NEXO.
// Implementa operações de fila giratória para distribuição de clientes.
package barberturn

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// =============================================================================
// ListBarberTurnUseCase - Lista barbeiros na fila
// =============================================================================

// ListBarberTurnUseCase lista barbeiros na fila da vez
type ListBarberTurnUseCase struct {
	repo   port.BarberTurnRepository
	logger *zap.Logger
}

// NewListBarberTurnUseCase cria uma nova instância
func NewListBarberTurnUseCase(repo port.BarberTurnRepository, logger *zap.Logger) *ListBarberTurnUseCase {
	return &ListBarberTurnUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute lista todos os barbeiros na fila
func (uc *ListBarberTurnUseCase) Execute(ctx context.Context, tenantID string, isActive *bool) (*dto.ListBarbersTurnResponse, error) {
	uc.logger.Info("🔍 [DEBUG] ListBarberTurnUseCase.Execute iniciado",
		zap.String("tenantID", tenantID),
		zap.Any("isActive", isActive))

	// Buscar lista de barbeiros
	barbers, err := uc.repo.List(ctx, tenantID, isActive)
	if err != nil {
		uc.logger.Error("erro ao listar barbeiros na fila", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("🔍 [DEBUG] Barbeiros retornados do repositório",
		zap.Int("count", len(barbers)),
		zap.Any("barbers", barbers))

	// Buscar estatísticas
	stats, err := uc.repo.GetStats(ctx, tenantID)
	if err != nil {
		uc.logger.Error("erro ao buscar estatísticas", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("🔍 [DEBUG] Estatísticas retornadas", zap.Any("stats", stats))

	// Buscar próximo barbeiro
	var nextBarber *dto.NextBarberResponse
	next, err := uc.repo.GetNextBarber(ctx, tenantID)
	if err == nil && next != nil {
		nextBarber = &dto.NextBarberResponse{
			ProfessionalID:    next.ProfessionalID,
			ProfessionalName:  next.ProfessionalName,
			ProfessionalPhoto: next.ProfessionalPhoto,
			CurrentPoints:     next.CurrentPoints,
		}
	}

	// Mapear resposta
	response := &dto.ListBarbersTurnResponse{
		Barbers:    make([]dto.BarberTurnResponse, 0, len(barbers)),
		Total:      len(barbers),
		NextBarber: nextBarber,
		Stats: dto.BarberTurnStatsResponse{
			TotalAtivos:    stats.TotalAtivos,
			TotalPausados:  stats.TotalPausados,
			TotalGeral:     stats.TotalGeral,
			TotalPontosMes: stats.TotalPontosMes,
		},
	}

	for _, b := range barbers {
		response.Barbers = append(response.Barbers, mapBarberTurnToResponse(b))
	}

	uc.logger.Info("🔍 [DEBUG] Resposta final montada",
		zap.Int("total_barbers", len(response.Barbers)),
		zap.Any("response", response))

	return response, nil
}

// =============================================================================
// AddBarberToTurnListUseCase - Adiciona barbeiro à fila
// =============================================================================

// AddBarberToTurnListUseCase adiciona um barbeiro à lista da vez
type AddBarberToTurnListUseCase struct {
	repo   port.BarberTurnRepository
	logger *zap.Logger
}

// NewAddBarberToTurnListUseCase cria uma nova instância
func NewAddBarberToTurnListUseCase(repo port.BarberTurnRepository, logger *zap.Logger) *AddBarberToTurnListUseCase {
	return &AddBarberToTurnListUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute adiciona um barbeiro à lista da vez
func (uc *AddBarberToTurnListUseCase) Execute(ctx context.Context, tenantID string, req dto.AddBarberToTurnListRequest) (*dto.BarberTurnResponse, error) {
	// Verificar se profissional é do tipo BARBEIRO
	isBarber, err := uc.repo.CheckProfessionalIsBarber(ctx, tenantID, req.ProfessionalID)
	if err != nil {
		uc.logger.Error("erro ao verificar tipo do profissional", zap.Error(err))
		return nil, err
	}
	if !isBarber {
		return nil, domain.ErrBarberTurnProfessionalNotBarber
	}

	// Verificar se já está na lista
	inList, err := uc.repo.CheckProfessionalInList(ctx, tenantID, req.ProfessionalID)
	if err != nil {
		uc.logger.Error("erro ao verificar se profissional está na lista", zap.Error(err))
		return nil, err
	}
	if inList {
		return nil, domain.ErrBarberTurnAlreadyInList
	}

	// Criar entidade
	barberTurn, err := entity.NewBarberTurn(tenantID, req.ProfessionalID)
	if err != nil {
		return nil, err
	}

	// Persistir
	if err := uc.repo.Add(ctx, barberTurn); err != nil {
		uc.logger.Error("erro ao adicionar barbeiro à fila", zap.Error(err))
		return nil, err
	}

	// Buscar dados completos (com JOIN)
	result, err := uc.repo.FindByProfessionalID(ctx, tenantID, req.ProfessionalID)
	if err != nil {
		uc.logger.Error("erro ao buscar barbeiro adicionado", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("barbeiro adicionado à lista da vez",
		zap.String("professional_id", req.ProfessionalID),
		zap.String("tenant_id", tenantID),
	)

	resp := mapBarberTurnToResponse(result)
	return &resp, nil
}

// =============================================================================
// RecordTurnUseCase - Registra atendimento (incrementa pontos)
// =============================================================================

// RecordTurnUseCase registra um atendimento
type RecordTurnUseCase struct {
	repo   port.BarberTurnRepository
	logger *zap.Logger
}

// NewRecordTurnUseCase cria uma nova instância
func NewRecordTurnUseCase(repo port.BarberTurnRepository, logger *zap.Logger) *RecordTurnUseCase {
	return &RecordTurnUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute registra um atendimento para um barbeiro
func (uc *RecordTurnUseCase) Execute(ctx context.Context, tenantID string, req dto.RecordTurnRequest) (*dto.RecordTurnResponse, error) {
	// Buscar barbeiro atual
	barber, err := uc.repo.FindByProfessionalID(ctx, tenantID, req.ProfessionalID)
	if err != nil {
		uc.logger.Error("erro ao buscar barbeiro", zap.Error(err))
		return nil, domain.ErrBarberTurnNotFound
	}

	// Verificar se está ativo
	if !barber.IsActive {
		return nil, domain.ErrBarberTurnCannotRecord
	}

	// Guardar pontos anteriores
	previousPoints := barber.CurrentPoints

	// Registrar atendimento
	result, err := uc.repo.RecordTurn(ctx, tenantID, req.ProfessionalID)
	if err != nil {
		uc.logger.Error("erro ao registrar atendimento", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("atendimento registrado",
		zap.String("professional_id", req.ProfessionalID),
		zap.String("tenant_id", tenantID),
		zap.Int("previous_points", previousPoints),
		zap.Int("new_points", result.CurrentPoints),
	)

	return &dto.RecordTurnResponse{
		ProfessionalID:   result.ProfessionalID,
		ProfessionalName: result.ProfessionalName,
		PreviousPoints:   previousPoints,
		NewPoints:        result.CurrentPoints,
		LastTurnAt:       time.Now(),
		Message:          "Atendimento registrado com sucesso",
	}, nil
}

// =============================================================================
// ToggleBarberStatusUseCase - Pausar/ativar barbeiro
// =============================================================================

// ToggleBarberStatusUseCase alterna status ativo/inativo
type ToggleBarberStatusUseCase struct {
	repo   port.BarberTurnRepository
	logger *zap.Logger
}

// NewToggleBarberStatusUseCase cria uma nova instância
func NewToggleBarberStatusUseCase(repo port.BarberTurnRepository, logger *zap.Logger) *ToggleBarberStatusUseCase {
	return &ToggleBarberStatusUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute alterna o status de um barbeiro
func (uc *ToggleBarberStatusUseCase) Execute(ctx context.Context, tenantID, professionalID string) (*dto.ToggleStatusResponse, error) {
	// Verificar se existe
	barber, err := uc.repo.FindByProfessionalID(ctx, tenantID, professionalID)
	if err != nil {
		uc.logger.Error("erro ao buscar barbeiro", zap.Error(err))
		return nil, domain.ErrBarberTurnNotFound
	}

	// Alternar status
	result, err := uc.repo.ToggleStatus(ctx, tenantID, professionalID)
	if err != nil {
		uc.logger.Error("erro ao alternar status", zap.Error(err))
		return nil, err
	}

	var msg string
	if result.IsActive {
		msg = "Barbeiro ativado na fila"
	} else {
		msg = "Barbeiro pausado na fila"
	}

	uc.logger.Info("status do barbeiro alterado",
		zap.String("professional_id", professionalID),
		zap.String("tenant_id", tenantID),
		zap.Bool("is_active", result.IsActive),
	)

	return &dto.ToggleStatusResponse{
		ProfessionalID:   result.ProfessionalID,
		ProfessionalName: barber.ProfessionalName,
		IsActive:         result.IsActive,
		Message:          msg,
	}, nil
}

// =============================================================================
// RemoveBarberFromTurnListUseCase - Remove barbeiro da fila
// =============================================================================

// RemoveBarberFromTurnListUseCase remove um barbeiro da lista
type RemoveBarberFromTurnListUseCase struct {
	repo   port.BarberTurnRepository
	logger *zap.Logger
}

// NewRemoveBarberFromTurnListUseCase cria uma nova instância
func NewRemoveBarberFromTurnListUseCase(repo port.BarberTurnRepository, logger *zap.Logger) *RemoveBarberFromTurnListUseCase {
	return &RemoveBarberFromTurnListUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute remove um barbeiro da lista
func (uc *RemoveBarberFromTurnListUseCase) Execute(ctx context.Context, tenantID, professionalID string) (*dto.RemoveBarberResponse, error) {
	// Verificar se existe
	_, err := uc.repo.FindByProfessionalID(ctx, tenantID, professionalID)
	if err != nil {
		uc.logger.Error("erro ao buscar barbeiro", zap.Error(err))
		return nil, domain.ErrBarberTurnNotFound
	}

	// Remover
	if err := uc.repo.Remove(ctx, tenantID, professionalID); err != nil {
		uc.logger.Error("erro ao remover barbeiro", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("barbeiro removido da lista da vez",
		zap.String("professional_id", professionalID),
		zap.String("tenant_id", tenantID),
	)

	return &dto.RemoveBarberResponse{
		ProfessionalID: professionalID,
		Message:        "Barbeiro removido da lista com sucesso",
	}, nil
}

// =============================================================================
// ResetTurnListUseCase - Reset mensal
// =============================================================================

// ResetTurnListUseCase executa reset mensal
type ResetTurnListUseCase struct {
	repo   port.BarberTurnRepository
	logger *zap.Logger
}

// NewResetTurnListUseCase cria uma nova instância
func NewResetTurnListUseCase(repo port.BarberTurnRepository, logger *zap.Logger) *ResetTurnListUseCase {
	return &ResetTurnListUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute executa reset mensal da lista
func (uc *ResetTurnListUseCase) Execute(ctx context.Context, tenantID string, saveHistory bool) (*dto.ResetTurnListResponse, error) {
	// Buscar estatísticas atuais antes do reset
	stats, err := uc.repo.GetStats(ctx, tenantID)
	if err != nil {
		uc.logger.Error("erro ao buscar estatísticas", zap.Error(err))
		return nil, err
	}

	// Calcular mês/ano anterior
	now := time.Now()
	monthYear := fmt.Sprintf("%d-%02d", now.Year(), now.Month())

	var snapshot *dto.TurnResetSnapshot
	if saveHistory {
		// Salvar histórico
		if err := uc.repo.SaveHistoryBeforeReset(ctx, tenantID, monthYear); err != nil {
			uc.logger.Error("erro ao salvar histórico", zap.Error(err))
			return nil, err
		}

		snapshot = &dto.TurnResetSnapshot{
			MonthYear:             monthYear,
			TotalBarbers:          int(stats.TotalGeral),
			TotalPointsReset:      stats.TotalPontosMes,
			HistoryRecordsCreated: int(stats.TotalAtivos),
		}
	}

	// Executar reset
	if err := uc.repo.ResetAll(ctx, tenantID); err != nil {
		uc.logger.Error("erro ao executar reset", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("reset mensal executado",
		zap.String("tenant_id", tenantID),
		zap.String("month_year", monthYear),
		zap.Int64("total_points_reset", stats.TotalPontosMes),
	)

	return &dto.ResetTurnListResponse{
		Message:  "Reset mensal executado com sucesso",
		Snapshot: snapshot,
	}, nil
}

// =============================================================================
// GetTurnHistoryUseCase - Lista histórico mensal
// =============================================================================

// GetTurnHistoryUseCase busca histórico de atendimentos
type GetTurnHistoryUseCase struct {
	repo   port.BarberTurnRepository
	logger *zap.Logger
}

// NewGetTurnHistoryUseCase cria uma nova instância
func NewGetTurnHistoryUseCase(repo port.BarberTurnRepository, logger *zap.Logger) *GetTurnHistoryUseCase {
	return &GetTurnHistoryUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute busca histórico
func (uc *GetTurnHistoryUseCase) Execute(ctx context.Context, tenantID string, monthYear *string) (*dto.ListTurnHistoryResponse, error) {
	history, err := uc.repo.ListHistory(ctx, tenantID, monthYear)
	if err != nil {
		uc.logger.Error("erro ao buscar histórico", zap.Error(err))
		return nil, err
	}

	response := &dto.ListTurnHistoryResponse{
		History: make([]dto.TurnHistoryResponse, 0, len(history)),
		Total:   len(history),
	}

	for _, h := range history {
		response.History = append(response.History, dto.TurnHistoryResponse{
			ID:               h.ID,
			TenantID:         h.TenantID,
			ProfessionalID:   h.ProfessionalID,
			ProfessionalName: h.ProfessionalName,
			MonthYear:        h.MonthYear,
			TotalTurns:       h.TotalTurns,
			FinalPoints:      h.FinalPoints,
			CreatedAt:        h.CreatedAt,
		})
	}

	return response, nil
}

// =============================================================================
// GetHistorySummaryUseCase - Resumo dos últimos 12 meses
// =============================================================================

// GetHistorySummaryUseCase busca resumo do histórico
type GetHistorySummaryUseCase struct {
	repo   port.BarberTurnRepository
	logger *zap.Logger
}

// NewGetHistorySummaryUseCase cria uma nova instância
func NewGetHistorySummaryUseCase(repo port.BarberTurnRepository, logger *zap.Logger) *GetHistorySummaryUseCase {
	return &GetHistorySummaryUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute busca resumo do histórico
func (uc *GetHistorySummaryUseCase) Execute(ctx context.Context, tenantID string) (*dto.ListHistorySummaryResponse, error) {
	summary, err := uc.repo.GetHistorySummary(ctx, tenantID)
	if err != nil {
		uc.logger.Error("erro ao buscar resumo do histórico", zap.Error(err))
		return nil, err
	}

	response := &dto.ListHistorySummaryResponse{
		Summary: make([]dto.TurnHistorySummaryResponse, 0, len(summary)),
	}

	for _, s := range summary {
		response.Summary = append(response.Summary, dto.TurnHistorySummaryResponse{
			MonthYear:         s.MonthYear,
			TotalBarbeiros:    s.TotalBarbeiros,
			TotalAtendimentos: s.TotalAtendimentos,
			MediaAtendimentos: s.MediaAtendimentos,
		})
	}

	return response, nil
}

// =============================================================================
// GetAvailableBarbersUseCase - Lista barbeiros disponíveis para adicionar
// =============================================================================

// GetAvailableBarbersUseCase lista barbeiros não adicionados
type GetAvailableBarbersUseCase struct {
	repo   port.BarberTurnRepository
	logger *zap.Logger
}

// NewGetAvailableBarbersUseCase cria uma nova instância
func NewGetAvailableBarbersUseCase(repo port.BarberTurnRepository, logger *zap.Logger) *GetAvailableBarbersUseCase {
	return &GetAvailableBarbersUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute lista barbeiros disponíveis
func (uc *GetAvailableBarbersUseCase) Execute(ctx context.Context, tenantID string) (*dto.ListAvailableBarbersResponse, error) {
	barbers, err := uc.repo.GetAvailableBarbers(ctx, tenantID)
	if err != nil {
		uc.logger.Error("erro ao buscar barbeiros disponíveis", zap.Error(err))
		return nil, err
	}

	response := &dto.ListAvailableBarbersResponse{
		Barbers: make([]dto.AvailableBarberResponse, 0, len(barbers)),
		Total:   len(barbers),
	}

	for _, b := range barbers {
		response.Barbers = append(response.Barbers, dto.AvailableBarberResponse{
			ID:     b.ID,
			Nome:   b.Nome,
			Foto:   b.Foto,
			Status: b.Status,
		})
	}

	return response, nil
}

// =============================================================================
// Helpers
// =============================================================================

func mapBarberTurnToResponse(b *entity.BarberTurn) dto.BarberTurnResponse {
	return dto.BarberTurnResponse{
		ID:                b.ID,
		TenantID:          b.TenantID,
		ProfessionalID:    b.ProfessionalID,
		ProfessionalName:  b.ProfessionalName,
		ProfessionalType:  b.ProfessionalType,
		ProfessionalPhoto: b.ProfessionalPhoto,
		CurrentPoints:     b.CurrentPoints,
		LastTurnAt:        b.LastTurnAt,
		IsActive:          b.IsActive,
		Position:          b.Position,
		CreatedAt:         b.CreatedAt,
		UpdatedAt:         b.UpdatedAt,
	}
}
