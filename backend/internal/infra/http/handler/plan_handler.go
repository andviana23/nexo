package handler

import (
	"net/http"
	"strconv"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	planUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/plan"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// PlanHandler agrupa handlers de planos de assinatura
type PlanHandler struct {
	createUC     *planUC.CreatePlanUseCase
	getUC        *planUC.GetPlanUseCase
	listUC       *planUC.ListPlansUseCase
	updateUC     *planUC.UpdatePlanUseCase
	deactivateUC *planUC.DeactivatePlanUseCase
	logger       *zap.Logger
}

// NewPlanHandler cria instância
func NewPlanHandler(createUC *planUC.CreatePlanUseCase, getUC *planUC.GetPlanUseCase, listUC *planUC.ListPlansUseCase, updateUC *planUC.UpdatePlanUseCase, deactivateUC *planUC.DeactivatePlanUseCase, logger *zap.Logger) *PlanHandler {
	return &PlanHandler{
		createUC:     createUC,
		getUC:        getUC,
		listUC:       listUC,
		updateUC:     updateUC,
		deactivateUC: deactivateUC,
		logger:       logger,
	}
}

// Create cria um novo plano
func (h *PlanHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "Tenant ID não encontrado"})
	}

	var req dto.CreatePlanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "Dados inválidos"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: err.Error()})
	}

	res, err := h.createUC.Execute(ctx, tenantID, req)
	if err != nil {
		switch err {
		case domain.ErrPlanNameDuplicate:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "duplicate_name", Message: err.Error()})
		case domain.ErrPlanValueInvalid, domain.ErrPlanNameRequired, domain.ErrPlanNameTooShort, domain.ErrPlanNameTooLong:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: err.Error()})
		default:
			h.logger.Error("erro ao criar plano", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "Erro ao criar plano"})
		}
	}
	return c.JSON(http.StatusCreated, res)
}

// List retorna planos (com filtro ?ativo=true opcional)
func (h *PlanHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "Tenant ID não encontrado"})
	}

	onlyActive := false
	if v := c.QueryParam("ativo"); v != "" {
		parsed, _ := strconv.ParseBool(v)
		onlyActive = parsed
	}

	res, err := h.listUC.Execute(ctx, tenantID, onlyActive)
	if err != nil {
		h.logger.Error("erro ao listar planos", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "Erro ao listar planos"})
	}
	return c.JSON(http.StatusOK, res)
}

// Get retorna plano por ID
func (h *PlanHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "Tenant ID não encontrado"})
	}
	id := c.Param("id")
	res, err := h.getUC.Execute(ctx, tenantID, id)
	if err != nil {
		if err == domain.ErrPlanNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: err.Error()})
		}
		if err == domain.ErrInvalidID {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: err.Error()})
		}
		h.logger.Error("erro ao buscar plano", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "Erro ao buscar plano"})
	}
	return c.JSON(http.StatusOK, res)
}

// Update atualiza plano
func (h *PlanHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "Tenant ID não encontrado"})
	}
	id := c.Param("id")

	var req dto.UpdatePlanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "Dados inválidos"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: err.Error()})
	}

	res, err := h.updateUC.Execute(ctx, tenantID, id, req)
	if err != nil {
		switch err {
		case domain.ErrPlanNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: err.Error()})
		case domain.ErrPlanNameDuplicate:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "duplicate_name", Message: err.Error()})
		case domain.ErrPlanValueInvalid, domain.ErrPlanNameRequired, domain.ErrPlanNameTooShort, domain.ErrPlanNameTooLong:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: err.Error()})
		default:
			h.logger.Error("erro ao atualizar plano", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "Erro ao atualizar plano"})
		}
	}

	return c.JSON(http.StatusOK, res)
}

// Deactivate desativa plano
func (h *PlanHandler) Deactivate(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "Tenant ID não encontrado"})
	}
	id := c.Param("id")

	if err := h.deactivateUC.Execute(ctx, tenantID, id); err != nil {
		switch err {
		case domain.ErrInvalidID:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: err.Error()})
		case domain.ErrPlanNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: err.Error()})
		case domain.ErrPlanHasActiveSubscriptions:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "active_subscriptions", Message: err.Error()})
		default:
			h.logger.Error("erro ao desativar plano", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "Erro ao desativar plano"})
		}
	}

	return c.NoContent(http.StatusNoContent)
}
