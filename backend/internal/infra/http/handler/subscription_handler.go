package handler

import (
	"net/http"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	subUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/subscription"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// SubscriptionHandler agrupa handlers de assinaturas
type SubscriptionHandler struct {
	createUC    *subUC.CreateSubscriptionUseCase
	listUC      *subUC.ListSubscriptionsUseCase
	getUC       *subUC.GetSubscriptionUseCase
	cancelUC    *subUC.CancelSubscriptionUseCase
	renewUC     *subUC.RenewSubscriptionUseCase
	metricsUC   *subUC.GetSubscriptionMetricsUseCase
	reconcileUC *subUC.ReconcileAsaasUseCase // T-ASAAS-002
	logger      *zap.Logger
}

// NewSubscriptionHandler cria instância
func NewSubscriptionHandler(
	createUC *subUC.CreateSubscriptionUseCase,
	listUC *subUC.ListSubscriptionsUseCase,
	getUC *subUC.GetSubscriptionUseCase,
	cancelUC *subUC.CancelSubscriptionUseCase,
	renewUC *subUC.RenewSubscriptionUseCase,
	metricsUC *subUC.GetSubscriptionMetricsUseCase,
	reconcileUC *subUC.ReconcileAsaasUseCase, // T-ASAAS-002
	logger *zap.Logger,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		createUC:    createUC,
		listUC:      listUC,
		getUC:       getUC,
		cancelUC:    cancelUC,
		renewUC:     renewUC,
		metricsUC:   metricsUC,
		reconcileUC: reconcileUC,
		logger:      logger,
	}
}

// Create cria nova assinatura
func (h *SubscriptionHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "Tenant ID não encontrado"})
	}

	var req dto.CreateSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "Dados inválidos"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: err.Error()})
	}

	res, err := h.createUC.Execute(ctx, tenantID, req)
	if err != nil {
		switch err {
		case domain.ErrCustomerNotFound, domain.ErrPlanNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: err.Error()})
		case domain.ErrPlanInactive, domain.ErrSubscriptionPaymentMethodInvalid, domain.ErrSubscriptionDuplicateActive:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: err.Error()})
		default:
			h.logger.Error("erro ao criar assinatura", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "Erro ao criar assinatura"})
		}
	}

	return c.JSON(http.StatusCreated, res)
}

// List retorna assinaturas (opcional status=?status)
func (h *SubscriptionHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "Tenant ID não encontrado"})
	}

	var statusPtr *string
	if v := c.QueryParam("status"); v != "" {
		statusPtr = &v
	}

	res, err := h.listUC.Execute(ctx, tenantID, statusPtr)
	if err != nil {
		h.logger.Error("erro ao listar assinaturas", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "Erro ao listar assinaturas"})
	}

	return c.JSON(http.StatusOK, res)
}

// Get retorna assinatura por ID
func (h *SubscriptionHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "Tenant ID não encontrado"})
	}
	id := c.Param("id")

	res, err := h.getUC.Execute(ctx, tenantID, id)
	if err != nil {
		if err == domain.ErrSubscriptionNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: err.Error()})
		}
		if err == domain.ErrInvalidID {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: err.Error()})
		}
		h.logger.Error("erro ao buscar assinatura", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "Erro ao buscar assinatura"})
	}

	return c.JSON(http.StatusOK, res)
}

// Renew renova assinatura manual
func (h *SubscriptionHandler) Renew(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "Tenant ID não encontrado"})
	}
	id := c.Param("id")

	var req dto.RenewSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Message: "Dados inválidos"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: err.Error()})
	}

	res, err := h.renewUC.Execute(ctx, tenantID, id, req)
	if err != nil {
		switch err {
		case domain.ErrSubscriptionNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: err.Error()})
		case domain.ErrSubscriptionPaymentMethodInvalid, domain.ErrSubscriptionCannotReactivate:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: err.Error()})
		default:
			h.logger.Error("erro ao renovar assinatura", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "Erro ao renovar assinatura"})
		}
	}

	return c.JSON(http.StatusOK, res)
}

// Cancel cancela assinatura
func (h *SubscriptionHandler) Cancel(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "Tenant ID não encontrado"})
	}
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "User ID não encontrado"})
	}
	id := c.Param("id")

	res, err := h.cancelUC.Execute(ctx, tenantID, id, userID)
	if err != nil {
		switch err {
		case domain.ErrSubscriptionNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: err.Error()})
		case domain.ErrSubscriptionAlreadyCanceled:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "already_canceled", Message: err.Error()})
		default:
			h.logger.Error("erro ao cancelar assinatura", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "Erro ao cancelar assinatura"})
		}
	}

	return c.JSON(http.StatusOK, res)
}

// Metrics retorna métricas agregadas
func (h *SubscriptionHandler) Metrics(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "Tenant ID não encontrado"})
	}

	res, err := h.metricsUC.Execute(ctx, tenantID)
	if err != nil {
		h.logger.Error("erro ao obter métricas de assinatura", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error", Message: "Erro ao obter métricas"})
	}
	return c.JSON(http.StatusOK, res)
}

// Reconcile executa reconciliação entre Asaas e NEXO
// @Summary Reconciliar pagamentos Asaas
// @Description Executa reconciliação automática entre pagamentos do Asaas e contas a receber do NEXO
// @Tags Assinaturas
// @Accept json
// @Produce json
// @Param start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param end_date query string false "Data final (YYYY-MM-DD)"
// @Param full_sync query bool false "Força sincronização completa"
// @Success 200 {object} dto.ReconcileAsaasResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/subscriptions/reconcile [post]
// @Security BearerAuth
func (h *SubscriptionHandler) Reconcile(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "Tenant ID não encontrado"})
	}

	// Parse query params
	input := subUC.ReconcileAsaasInput{
		TenantID: tenantID,
		FullSync: c.QueryParam("full_sync") == "true",
		AutoFix:  c.QueryParam("auto_fix") == "true",
	}

	// Parse dates
	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Data inicial inválida (formato: YYYY-MM-DD)",
			})
		}
		input.DataInicio = &startDate
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Data final inválida (formato: YYYY-MM-DD)",
			})
		}
		input.DataFim = &endDate
	}

	// Execute reconciliation
	output, err := h.reconcileUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("erro ao executar reconciliação Asaas", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao executar reconciliação",
		})
	}

	return c.JSON(http.StatusOK, dto.ReconcileAsaasResponse{
		TotalProcessed:    output.TotalProcessed,
		TotalMatched:      output.TotalMatched,
		TotalMissingNexo:  output.TotalMissingNexo,
		TotalMissingAsaas: output.TotalMissingAsaas,
		TotalDivergent:    output.TotalDivergent,
		TotalAutoFixed:    output.TotalAutoFixed,
		Errors:            output.Errors,
		ReconciliationID:  output.ReconciliationID,
	})
}
