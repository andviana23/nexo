package handler

import (
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/blockedtime"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// BlockedTimeHandler manipula requisições de bloqueios de horário
type BlockedTimeHandler struct {
	createUC *blockedtime.CreateBlockedTimeUseCase
	listUC   *blockedtime.ListBlockedTimesUseCase
	deleteUC *blockedtime.DeleteBlockedTimeUseCase
	logger   *zap.Logger
}

// NewBlockedTimeHandler cria uma nova instância do handler
func NewBlockedTimeHandler(
	createUC *blockedtime.CreateBlockedTimeUseCase,
	listUC *blockedtime.ListBlockedTimesUseCase,
	deleteUC *blockedtime.DeleteBlockedTimeUseCase,
	logger *zap.Logger,
) *BlockedTimeHandler {
	return &BlockedTimeHandler{
		createUC: createUC,
		listUC:   listUC,
		deleteUC: deleteUC,
		logger:   logger,
	}
}

// CreateBlockedTime godoc
// @Summary Criar bloqueio de horário
// @Description Cria um novo bloqueio de horário na agenda
// @Tags Bloqueios
// @Accept json
// @Produce json
// @Param request body dto.CreateBlockedTimeRequest true "Dados do bloqueio"
// @Success 201 {object} dto.BlockedTimeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Conflito com bloqueio existente"
// @Router /api/v1/blocked-times [post]
// @Security BearerAuth
func (h *BlockedTimeHandler) CreateBlockedTime(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrai tenant_id do contexto
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Extrai user_id (opcional)
	var userID *string
	if uid, ok := c.Get("user_id").(string); ok && uid != "" {
		userID = &uid
	}

	// Parse request
	var req dto.CreateBlockedTimeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	// Validação
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// Executa use case
	output, err := h.createUC.Execute(ctx, blockedtime.CreateBlockedTimeInput{
		TenantID:       tenantID,
		ProfessionalID: req.ProfessionalID,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Reason:         req.Reason,
		UserID:         userID,
	})

	if err != nil {
		h.logger.Error("Erro ao criar bloqueio de horário", zap.Error(err))

		if err == entity.ErrTimeRangeOverlap {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: "Conflito com bloqueio existente",
			})
		}

		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "error",
			Message: err.Error(),
		})
	}

	// Converte para response
	response := h.toResponse(output.BlockedTime)

	return c.JSON(http.StatusCreated, response)
}

// ListBlockedTimes godoc
// @Summary Listar bloqueios de horário
// @Description Lista bloqueios de horário com filtros opcionais
// @Tags Bloqueios
// @Accept json
// @Produce json
// @Param professional_id query string false "ID do profissional"
// @Param start_date query string false "Data inicial (RFC3339)"
// @Param end_date query string false "Data final (RFC3339)"
// @Success 200 {object} dto.ListBlockedTimesResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/blocked-times [get]
// @Security BearerAuth
func (h *BlockedTimeHandler) ListBlockedTimes(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrai tenant_id do contexto
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Parse query params
	var req dto.ListBlockedTimesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	// Executa use case
	output, err := h.listUC.Execute(ctx, blockedtime.ListBlockedTimesInput{
		TenantID:       tenantID,
		ProfessionalID: req.ProfessionalID,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
	})

	if err != nil {
		h.logger.Error("Erro ao listar bloqueios de horário", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "error",
			Message: "Erro ao listar bloqueios",
		})
	}

	// Converte para response
	response := dto.ListBlockedTimesResponse{
		Data:  make([]dto.BlockedTimeResponse, len(output.BlockedTimes)),
		Total: len(output.BlockedTimes),
	}

	for i, bt := range output.BlockedTimes {
		response.Data[i] = *h.toResponse(bt)
	}

	return c.JSON(http.StatusOK, response)
}

// DeleteBlockedTime godoc
// @Summary Deletar bloqueio de horário
// @Description Remove um bloqueio de horário
// @Tags Bloqueios
// @Accept json
// @Produce json
// @Param id path string true "ID do bloqueio"
// @Success 204
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/blocked-times/{id} [delete]
// @Security BearerAuth
func (h *BlockedTimeHandler) DeleteBlockedTime(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrai tenant_id do contexto
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Extrai ID do path
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do bloqueio é obrigatório",
		})
	}

	// Executa use case
	err := h.deleteUC.Execute(ctx, blockedtime.DeleteBlockedTimeInput{
		TenantID: tenantID,
		ID:       id,
	})

	if err != nil {
		h.logger.Error("Erro ao deletar bloqueio de horário", zap.Error(err))
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Bloqueio não encontrado",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// toResponse converte entidade para DTO de resposta
func (h *BlockedTimeHandler) toResponse(bt *entity.BlockedTime) *dto.BlockedTimeResponse {
	return &dto.BlockedTimeResponse{
		ID:             bt.ID,
		TenantID:       bt.TenantID.String(),
		ProfessionalID: bt.ProfessionalID,
		StartTime:      bt.StartTime,
		EndTime:        bt.EndTime,
		Reason:         bt.Reason,
		IsRecurring:    bt.IsRecurring,
		RecurrenceRule: bt.RecurrenceRule,
		CreatedAt:      bt.CreatedAt,
		UpdatedAt:      bt.UpdatedAt,
		CreatedBy:      bt.CreatedBy,
	}
}
