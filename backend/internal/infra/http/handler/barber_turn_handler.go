package handler

import (
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/barberturn"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// BarberTurnHandler agrupa os handlers da Lista da Vez
type BarberTurnHandler struct {
	listUC             *barberturn.ListBarberTurnUseCase
	addUC              *barberturn.AddBarberToTurnListUseCase
	recordTurnUC       *barberturn.RecordTurnUseCase
	toggleStatusUC     *barberturn.ToggleBarberStatusUseCase
	removeUC           *barberturn.RemoveBarberFromTurnListUseCase
	resetUC            *barberturn.ResetTurnListUseCase
	historyUC          *barberturn.GetTurnHistoryUseCase
	historySummaryUC   *barberturn.GetHistorySummaryUseCase
	availableBarbersUC *barberturn.GetAvailableBarbersUseCase
	logger             *zap.Logger
}

// NewBarberTurnHandler cria um novo handler da Lista da Vez
func NewBarberTurnHandler(
	listUC *barberturn.ListBarberTurnUseCase,
	addUC *barberturn.AddBarberToTurnListUseCase,
	recordTurnUC *barberturn.RecordTurnUseCase,
	toggleStatusUC *barberturn.ToggleBarberStatusUseCase,
	removeUC *barberturn.RemoveBarberFromTurnListUseCase,
	resetUC *barberturn.ResetTurnListUseCase,
	historyUC *barberturn.GetTurnHistoryUseCase,
	historySummaryUC *barberturn.GetHistorySummaryUseCase,
	availableBarbersUC *barberturn.GetAvailableBarbersUseCase,
	logger *zap.Logger,
) *BarberTurnHandler {
	return &BarberTurnHandler{
		listUC:             listUC,
		addUC:              addUC,
		recordTurnUC:       recordTurnUC,
		toggleStatusUC:     toggleStatusUC,
		removeUC:           removeUC,
		resetUC:            resetUC,
		historyUC:          historyUC,
		historySummaryUC:   historySummaryUC,
		availableBarbersUC: availableBarbersUC,
		logger:             logger,
	}
}

// ListBarbersTurn godoc
// @Summary Listar barbeiros na fila
// @Description Lista todos os barbeiros na Lista da Vez ordenados por pontuação
// @Tags Lista da Vez
// @Produce json
// @Param is_active query bool false "Filtrar por status ativo"
// @Success 200 {object} dto.ListBarbersTurnResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/barber-turn/list [get]
// @Security BearerAuth
func (h *BarberTurnHandler) ListBarbersTurn(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ListBarbersTurnRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind", zap.Error(err))
	}

	result, err := h.listUC.Execute(ctx, tenantID, req.IsActive)
	if err != nil {
		h.logger.Error("Erro ao listar barbeiros", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar barbeiros na fila",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// AddBarberToTurnList godoc
// @Summary Adicionar barbeiro à fila
// @Description Adiciona um profissional do tipo BARBEIRO à Lista da Vez
// @Tags Lista da Vez
// @Accept json
// @Produce json
// @Param request body dto.AddBarberToTurnListRequest true "ID do profissional"
// @Success 201 {object} dto.BarberTurnResponse
// @Failure 400 {object} dto.ErrorResponse "Profissional não é barbeiro ou já está na lista"
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/barber-turn/add [post]
// @Security BearerAuth
func (h *BarberTurnHandler) AddBarberToTurnList(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.AddBarberToTurnListRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	result, err := h.addUC.Execute(ctx, tenantID, req)
	if err != nil {
		h.logger.Error("Erro ao adicionar barbeiro", zap.Error(err))

		switch err {
		case domain.ErrBarberTurnProfessionalNotBarber:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "not_barber",
				Message: err.Error(),
			})
		case domain.ErrBarberTurnAlreadyInList:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "already_in_list",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Erro ao adicionar barbeiro à fila",
			})
		}
	}

	return c.JSON(http.StatusCreated, result)
}

// RecordTurn godoc
// @Summary Registrar atendimento
// @Description Incrementa pontos do barbeiro (+1) e reordena a fila
// @Tags Lista da Vez
// @Accept json
// @Produce json
// @Param request body dto.RecordTurnRequest true "ID do profissional"
// @Success 200 {object} dto.RecordTurnResponse
// @Failure 400 {object} dto.ErrorResponse "Barbeiro pausado"
// @Failure 404 {object} dto.ErrorResponse "Barbeiro não encontrado"
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/barber-turn/record [post]
// @Security BearerAuth
func (h *BarberTurnHandler) RecordTurn(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.RecordTurnRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	result, err := h.recordTurnUC.Execute(ctx, tenantID, req)
	if err != nil {
		h.logger.Error("Erro ao registrar atendimento", zap.Error(err))

		switch err {
		case domain.ErrBarberTurnNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		case domain.ErrBarberTurnCannotRecord:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "barber_paused",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Erro ao registrar atendimento",
			})
		}
	}

	return c.JSON(http.StatusOK, result)
}

// ToggleBarberStatus godoc
// @Summary Pausar/Ativar barbeiro
// @Description Alterna o status ativo/inativo de um barbeiro na fila
// @Tags Lista da Vez
// @Produce json
// @Param professional_id path string true "ID do profissional"
// @Success 200 {object} dto.ToggleStatusResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/barber-turn/{professional_id}/toggle-status [put]
// @Security BearerAuth
func (h *BarberTurnHandler) ToggleBarberStatus(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	professionalID := c.Param("professional_id")
	if professionalID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do profissional é obrigatório",
		})
	}

	result, err := h.toggleStatusUC.Execute(ctx, tenantID, professionalID)
	if err != nil {
		h.logger.Error("Erro ao alternar status", zap.Error(err))

		switch err {
		case domain.ErrBarberTurnNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Erro ao alternar status",
			})
		}
	}

	return c.JSON(http.StatusOK, result)
}

// RemoveBarberFromTurnList godoc
// @Summary Remover barbeiro da fila
// @Description Remove um barbeiro da Lista da Vez (não preserva pontos)
// @Tags Lista da Vez
// @Produce json
// @Param professional_id path string true "ID do profissional"
// @Success 204 "No Content"
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/barber-turn/{professional_id} [delete]
// @Security BearerAuth
func (h *BarberTurnHandler) RemoveBarberFromTurnList(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	professionalID := c.Param("professional_id")
	if professionalID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do profissional é obrigatório",
		})
	}

	_, err := h.removeUC.Execute(ctx, tenantID, professionalID)
	if err != nil {
		h.logger.Error("Erro ao remover barbeiro", zap.Error(err))

		switch err {
		case domain.ErrBarberTurnNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Erro ao remover barbeiro",
			})
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// ResetTurnList godoc
// @Summary Reset mensal
// @Description Zera todos os pontos e opcionalmente salva histórico
// @Tags Lista da Vez
// @Accept json
// @Produce json
// @Param request body dto.ResetTurnListRequest true "Configurações do reset"
// @Success 200 {object} dto.ResetTurnListResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/barber-turn/reset [post]
// @Security BearerAuth
func (h *BarberTurnHandler) ResetTurnList(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ResetTurnListRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind", zap.Error(err))
		req.SaveHistory = true // Default: salvar histórico
	}

	result, err := h.resetUC.Execute(ctx, tenantID, req.SaveHistory)
	if err != nil {
		h.logger.Error("Erro ao executar reset", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao executar reset mensal",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// GetTurnHistory godoc
// @Summary Listar histórico
// @Description Lista histórico mensal de atendimentos
// @Tags Lista da Vez
// @Produce json
// @Param month_year query string false "Mês/ano (YYYY-MM)"
// @Success 200 {object} dto.ListTurnHistoryResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/barber-turn/history [get]
// @Security BearerAuth
func (h *BarberTurnHandler) GetTurnHistory(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.GetTurnHistoryRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind", zap.Error(err))
	}

	var monthYear *string
	if req.MonthYear != "" {
		monthYear = &req.MonthYear
	}

	result, err := h.historyUC.Execute(ctx, tenantID, monthYear)
	if err != nil {
		h.logger.Error("Erro ao buscar histórico", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar histórico",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// GetHistorySummary godoc
// @Summary Resumo do histórico
// @Description Retorna resumo dos últimos 12 meses
// @Tags Lista da Vez
// @Produce json
// @Success 200 {object} dto.ListHistorySummaryResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/barber-turn/history/summary [get]
// @Security BearerAuth
func (h *BarberTurnHandler) GetHistorySummary(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	result, err := h.historySummaryUC.Execute(ctx, tenantID)
	if err != nil {
		h.logger.Error("Erro ao buscar resumo do histórico", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar resumo do histórico",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// GetAvailableBarbers godoc
// @Summary Listar barbeiros disponíveis
// @Description Lista barbeiros ativos que ainda não estão na Lista da Vez
// @Tags Lista da Vez
// @Produce json
// @Success 200 {object} dto.ListAvailableBarbersResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/barber-turn/available [get]
// @Security BearerAuth
func (h *BarberTurnHandler) GetAvailableBarbers(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	result, err := h.availableBarbersUC.Execute(ctx, tenantID)
	if err != nil {
		h.logger.Error("Erro ao buscar barbeiros disponíveis", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar barbeiros disponíveis",
		})
	}

	return c.JSON(http.StatusOK, result)
}
