package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/commission"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// CommissionHandler agrupa os handlers de comissões
type CommissionHandler struct {
	// Commission Rule UseCases
	createRuleUC     *commission.CreateCommissionRuleUseCase
	getRuleUC        *commission.GetCommissionRuleUseCase
	listRulesUC      *commission.ListCommissionRulesUseCase
	getEffectiveUC   *commission.GetEffectiveCommissionRulesUseCase
	updateRuleUC     *commission.UpdateCommissionRuleUseCase
	deleteRuleUC     *commission.DeleteCommissionRuleUseCase
	deactivateRuleUC *commission.DeactivateCommissionRuleUseCase

	// Commission Period UseCases
	createPeriodUC     *commission.CreateCommissionPeriodUseCase
	getPeriodUC        *commission.GetCommissionPeriodUseCase
	getOpenPeriodUC    *commission.GetOpenCommissionPeriodUseCase
	getPeriodSummaryUC *commission.GetCommissionPeriodSummaryUseCase
	listPeriodsUC      *commission.ListCommissionPeriodsUseCase
	closePeriodUC      *commission.CloseCommissionPeriodUseCase
	markPeriodPaidUC   *commission.MarkPeriodAsPaidUseCase
	deletePeriodUC     *commission.DeleteCommissionPeriodUseCase

	// Advance UseCases
	createAdvanceUC       *commission.CreateAdvanceUseCase
	getAdvanceUC          *commission.GetAdvanceUseCase
	listAdvancesUC        *commission.ListAdvancesUseCase
	getPendingAdvancesUC  *commission.GetPendingAdvancesUseCase
	getApprovedAdvancesUC *commission.GetApprovedAdvancesUseCase
	approveAdvanceUC      *commission.ApproveAdvanceUseCase
	rejectAdvanceUC       *commission.RejectAdvanceUseCase
	markDeductedUC        *commission.MarkAdvanceDeductedUseCase
	cancelAdvanceUC       *commission.CancelAdvanceUseCase
	deleteAdvanceUC       *commission.DeleteAdvanceUseCase

	// Commission Item UseCases
	createItemUC               *commission.CreateCommissionItemUseCase
	createItemBatchUC          *commission.CreateCommissionItemBatchUseCase
	getItemUC                  *commission.GetCommissionItemUseCase
	listItemsUC                *commission.ListCommissionItemsUseCase
	getPendingItemsUC          *commission.GetPendingCommissionItemsUseCase
	getSummaryByProfessionalUC *commission.GetCommissionSummaryByProfessionalUseCase
	getSummaryByServiceUC      *commission.GetCommissionSummaryByServiceUseCase
	processItemUC              *commission.ProcessCommissionItemUseCase
	assignItemsUC              *commission.AssignItemsToPeriodUseCase
	deleteItemUC               *commission.DeleteCommissionItemUseCase

	logger *zap.Logger
}

// NewCommissionHandler cria um novo handler de comissões
func NewCommissionHandler(
	// Commission Rule UseCases
	createRuleUC *commission.CreateCommissionRuleUseCase,
	getRuleUC *commission.GetCommissionRuleUseCase,
	listRulesUC *commission.ListCommissionRulesUseCase,
	getEffectiveUC *commission.GetEffectiveCommissionRulesUseCase,
	updateRuleUC *commission.UpdateCommissionRuleUseCase,
	deleteRuleUC *commission.DeleteCommissionRuleUseCase,
	deactivateRuleUC *commission.DeactivateCommissionRuleUseCase,
	// Commission Period UseCases
	createPeriodUC *commission.CreateCommissionPeriodUseCase,
	getPeriodUC *commission.GetCommissionPeriodUseCase,
	getOpenPeriodUC *commission.GetOpenCommissionPeriodUseCase,
	getPeriodSummaryUC *commission.GetCommissionPeriodSummaryUseCase,
	listPeriodsUC *commission.ListCommissionPeriodsUseCase,
	closePeriodUC *commission.CloseCommissionPeriodUseCase,
	markPeriodPaidUC *commission.MarkPeriodAsPaidUseCase,
	deletePeriodUC *commission.DeleteCommissionPeriodUseCase,
	// Advance UseCases
	createAdvanceUC *commission.CreateAdvanceUseCase,
	getAdvanceUC *commission.GetAdvanceUseCase,
	listAdvancesUC *commission.ListAdvancesUseCase,
	getPendingAdvancesUC *commission.GetPendingAdvancesUseCase,
	getApprovedAdvancesUC *commission.GetApprovedAdvancesUseCase,
	approveAdvanceUC *commission.ApproveAdvanceUseCase,
	rejectAdvanceUC *commission.RejectAdvanceUseCase,
	markDeductedUC *commission.MarkAdvanceDeductedUseCase,
	cancelAdvanceUC *commission.CancelAdvanceUseCase,
	deleteAdvanceUC *commission.DeleteAdvanceUseCase,
	// Commission Item UseCases
	createItemUC *commission.CreateCommissionItemUseCase,
	createItemBatchUC *commission.CreateCommissionItemBatchUseCase,
	getItemUC *commission.GetCommissionItemUseCase,
	listItemsUC *commission.ListCommissionItemsUseCase,
	getPendingItemsUC *commission.GetPendingCommissionItemsUseCase,
	getSummaryByProfessionalUC *commission.GetCommissionSummaryByProfessionalUseCase,
	getSummaryByServiceUC *commission.GetCommissionSummaryByServiceUseCase,
	processItemUC *commission.ProcessCommissionItemUseCase,
	assignItemsUC *commission.AssignItemsToPeriodUseCase,
	deleteItemUC *commission.DeleteCommissionItemUseCase,
	logger *zap.Logger,
) *CommissionHandler {
	return &CommissionHandler{
		createRuleUC:               createRuleUC,
		getRuleUC:                  getRuleUC,
		listRulesUC:                listRulesUC,
		getEffectiveUC:             getEffectiveUC,
		updateRuleUC:               updateRuleUC,
		deleteRuleUC:               deleteRuleUC,
		deactivateRuleUC:           deactivateRuleUC,
		createPeriodUC:             createPeriodUC,
		getPeriodUC:                getPeriodUC,
		getOpenPeriodUC:            getOpenPeriodUC,
		getPeriodSummaryUC:         getPeriodSummaryUC,
		listPeriodsUC:              listPeriodsUC,
		closePeriodUC:              closePeriodUC,
		markPeriodPaidUC:           markPeriodPaidUC,
		deletePeriodUC:             deletePeriodUC,
		createAdvanceUC:            createAdvanceUC,
		getAdvanceUC:               getAdvanceUC,
		listAdvancesUC:             listAdvancesUC,
		getPendingAdvancesUC:       getPendingAdvancesUC,
		getApprovedAdvancesUC:      getApprovedAdvancesUC,
		approveAdvanceUC:           approveAdvanceUC,
		rejectAdvanceUC:            rejectAdvanceUC,
		markDeductedUC:             markDeductedUC,
		cancelAdvanceUC:            cancelAdvanceUC,
		deleteAdvanceUC:            deleteAdvanceUC,
		createItemUC:               createItemUC,
		createItemBatchUC:          createItemBatchUC,
		getItemUC:                  getItemUC,
		listItemsUC:                listItemsUC,
		getPendingItemsUC:          getPendingItemsUC,
		getSummaryByProfessionalUC: getSummaryByProfessionalUC,
		getSummaryByServiceUC:      getSummaryByServiceUC,
		processItemUC:              processItemUC,
		assignItemsUC:              assignItemsUC,
		deleteItemUC:               deleteItemUC,
		logger:                     logger,
	}
}

// =============================================================================
// Commission Rule Handlers
// =============================================================================

// CreateCommissionRule godoc
// @Summary Criar regra de comissão
// @Description Cria uma nova regra de comissão
// @Tags Comissões
// @Accept json
// @Produce json
// @Param request body dto.CreateCommissionRuleRequest true "Dados da regra"
// @Success 201 {object} dto.CommissionRuleResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/commissions/rules [post]
// @Security BearerAuth
func (h *CommissionHandler) CreateCommissionRule(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var userID *string
	if uid, ok := c.Get("user_id").(string); ok && uid != "" {
		userID = &uid
	}

	var req dto.CreateCommissionRuleRequest
	if err := c.Bind(&req); err != nil {
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

	// Parse datas opcionais
	var effectiveFrom, effectiveTo *time.Time
	if req.EffectiveFrom != nil {
		t, err := time.Parse(time.RFC3339, *req.EffectiveFrom)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Data effective_from inválida",
			})
		}
		effectiveFrom = &t
	}
	if req.EffectiveTo != nil {
		t, err := time.Parse(time.RFC3339, *req.EffectiveTo)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Data effective_to inválida",
			})
		}
		effectiveTo = &t
	}

	output, err := h.createRuleUC.Execute(ctx, commission.CreateCommissionRuleInput{
		TenantID:        tenantID,
		UnitID:          req.UnitID,
		Name:            req.Name,
		Description:     req.Description,
		Type:            req.Type,
		DefaultRate:     req.DefaultRate,
		MinAmount:       req.MinAmount,
		MaxAmount:       req.MaxAmount,
		CalculationBase: req.CalculationBase,
		EffectiveFrom:   effectiveFrom,
		EffectiveTo:     effectiveTo,
		Priority:        req.Priority,
		CreatedBy:       userID,
	})

	if err != nil {
		h.logger.Error("Erro ao criar regra de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusCreated, h.ruleToResponse(output.CommissionRule))
}

// GetCommissionRule godoc
// @Summary Buscar regra de comissão
// @Description Busca uma regra de comissão por ID
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID da regra"
// @Success 200 {object} dto.CommissionRuleResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/commissions/rules/{id} [get]
// @Security BearerAuth
func (h *CommissionHandler) GetCommissionRule(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID da regra é obrigatório",
		})
	}

	output, err := h.getRuleUC.Execute(ctx, commission.GetCommissionRuleInput{
		TenantID: tenantID,
		ID:       id,
	})

	if err != nil {
		h.logger.Error("Erro ao buscar regra de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	if output.CommissionRule == nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Regra de comissão não encontrada",
		})
	}

	return c.JSON(http.StatusOK, h.ruleToResponse(output.CommissionRule))
}

// ListCommissionRules godoc
// @Summary Listar regras de comissão
// @Description Lista todas as regras de comissão do tenant
// @Tags Comissões
// @Accept json
// @Produce json
// @Param active_only query bool false "Apenas regras ativas"
// @Success 200 {object} dto.ListCommissionRulesResponse
// @Router /api/v1/commissions/rules [get]
// @Security BearerAuth
func (h *CommissionHandler) ListCommissionRules(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ListCommissionRulesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	output, err := h.listRulesUC.Execute(ctx, commission.ListCommissionRulesInput{
		TenantID:   tenantID,
		ActiveOnly: req.ActiveOnly,
	})

	if err != nil {
		h.logger.Error("Erro ao listar regras de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	rules := make([]dto.CommissionRuleResponse, 0, len(output.CommissionRules))
	for _, rule := range output.CommissionRules {
		rules = append(rules, h.ruleToResponse(rule))
	}

	return c.JSON(http.StatusOK, dto.ListCommissionRulesResponse{
		Data:  rules,
		Total: len(rules),
	})
}

// GetEffectiveCommissionRules godoc
// @Summary Listar regras efetivas
// @Description Lista regras de comissão efetivas para uma data específica
// @Tags Comissões
// @Accept json
// @Produce json
// @Param reference_date query string false "Data de referência (RFC3339)"
// @Success 200 {object} dto.ListCommissionRulesResponse
// @Router /api/v1/commissions/rules/effective [get]
// @Security BearerAuth
func (h *CommissionHandler) GetEffectiveCommissionRules(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Parse data de referência (opcional, usa data atual se não fornecida)
	referenceDate := time.Now()
	if dateStr := c.QueryParam("reference_date"); dateStr != "" {
		parsed, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Data reference_date inválida",
			})
		}
		referenceDate = parsed
	}

	output, err := h.getEffectiveUC.Execute(ctx, commission.GetEffectiveCommissionRulesInput{
		TenantID: tenantID,
		Date:     referenceDate,
	})

	if err != nil {
		h.logger.Error("Erro ao listar regras efetivas", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	rules := make([]dto.CommissionRuleResponse, 0, len(output.CommissionRules))
	for _, rule := range output.CommissionRules {
		rules = append(rules, h.ruleToResponse(rule))
	}

	return c.JSON(http.StatusOK, dto.ListCommissionRulesResponse{
		Data:  rules,
		Total: len(rules),
	})
}

// UpdateCommissionRule godoc
// @Summary Atualizar regra de comissão
// @Description Atualiza uma regra de comissão existente
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID da regra"
// @Param request body dto.UpdateCommissionRuleRequest true "Dados para atualização"
// @Success 200 {object} dto.CommissionRuleResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/commissions/rules/{id} [put]
// @Security BearerAuth
func (h *CommissionHandler) UpdateCommissionRule(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID da regra é obrigatório",
		})
	}

	var req dto.UpdateCommissionRuleRequest
	if err := c.Bind(&req); err != nil {
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

	// Parse datas opcionais
	var effectiveFrom, effectiveTo *time.Time
	if req.EffectiveFrom != nil {
		t, err := time.Parse(time.RFC3339, *req.EffectiveFrom)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Data effective_from inválida",
			})
		}
		effectiveFrom = &t
	}
	if req.EffectiveTo != nil {
		t, err := time.Parse(time.RFC3339, *req.EffectiveTo)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Data effective_to inválida",
			})
		}
		effectiveTo = &t
	}

	output, err := h.updateRuleUC.Execute(ctx, commission.UpdateCommissionRuleInput{
		TenantID:        tenantID,
		ID:              id,
		Name:            req.Name,
		Description:     req.Description,
		Type:            req.Type,
		DefaultRate:     req.DefaultRate,
		MinAmount:       req.MinAmount,
		MaxAmount:       req.MaxAmount,
		CalculationBase: req.CalculationBase,
		EffectiveFrom:   effectiveFrom,
		EffectiveTo:     effectiveTo,
		Priority:        req.Priority,
		IsActive:        req.IsActive,
	})

	if err != nil {
		h.logger.Error("Erro ao atualizar regra de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusOK, h.ruleToResponse(output.CommissionRule))
}

// DeleteCommissionRule godoc
// @Summary Excluir regra de comissão
// @Description Exclui uma regra de comissão
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID da regra"
// @Success 204 "No Content"
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/commissions/rules/{id} [delete]
// @Security BearerAuth
func (h *CommissionHandler) DeleteCommissionRule(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID da regra é obrigatório",
		})
	}

	_, err := h.deleteRuleUC.Execute(ctx, commission.DeleteCommissionRuleInput{
		TenantID: tenantID,
		ID:       id,
	})

	if err != nil {
		h.logger.Error("Erro ao excluir regra de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeactivateCommissionRule godoc
// @Summary Desativar regra de comissão
// @Description Desativa uma regra de comissão sem excluí-la
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID da regra"
// @Success 204 "No Content"
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/commissions/rules/{id}/deactivate [post]
// @Security BearerAuth
func (h *CommissionHandler) DeactivateCommissionRule(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID da regra é obrigatório",
		})
	}

	_, err := h.deactivateRuleUC.Execute(ctx, commission.DeactivateCommissionRuleInput{
		TenantID: tenantID,
		ID:       id,
	})

	if err != nil {
		h.logger.Error("Erro ao desativar regra de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// =============================================================================
// Helper Functions
// =============================================================================

func (h *CommissionHandler) handleDomainError(c echo.Context, err error) error {
	switch err {
	case domain.ErrCommissionRuleNotFound:
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Regra de comissão não encontrada",
		})
	case domain.ErrCommissionPeriodNotFound:
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Período de comissão não encontrado",
		})
	case domain.ErrAdvanceNotFound:
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Adiantamento não encontrado",
		})
	case domain.ErrCommissionItemNotFound:
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Item de comissão não encontrado",
		})
	case domain.ErrTipoComissaoInvalido:
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Tipo de comissão inválido",
		})
	case domain.ErrPercentualInvalido:
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Percentual inválido",
		})
	case domain.ErrPeriodoNaoPodeFechado:
		return c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error:   "conflict",
			Message: "Período não pode ser fechado neste status",
		})
	case domain.ErrPeriodoNaoPodePago:
		return c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error:   "conflict",
			Message: "Período não pode ser pago neste status",
		})
	case domain.ErrAdiantamentoNaoPodeAprovar:
		return c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error:   "conflict",
			Message: "Adiantamento não pode ser aprovado neste status",
		})
	case domain.ErrAdiantamentoNaoPodeRejeitar:
		return c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error:   "conflict",
			Message: "Adiantamento não pode ser rejeitado neste status",
		})
	case domain.ErrItemNaoPodeProcessar:
		return c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error:   "conflict",
			Message: "Item não pode ser processado neste status",
		})
	default:
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "error",
			Message: err.Error(),
		})
	}
}

func (h *CommissionHandler) ruleToResponse(rule *entity.CommissionRule) dto.CommissionRuleResponse {
	var minAmount, maxAmount, effectiveTo *string
	if rule.MinAmount != nil {
		v := rule.MinAmount.String()
		minAmount = &v
	}
	if rule.MaxAmount != nil {
		v := rule.MaxAmount.String()
		maxAmount = &v
	}
	if rule.EffectiveTo != nil {
		v := rule.EffectiveTo.Format(time.RFC3339)
		effectiveTo = &v
	}

	return dto.CommissionRuleResponse{
		ID:              rule.ID,
		TenantID:        rule.TenantID.String(),
		UnitID:          rule.UnitID,
		Name:            rule.Name,
		Description:     rule.Description,
		Type:            rule.Type,
		DefaultRate:     rule.DefaultRate.String(),
		MinAmount:       minAmount,
		MaxAmount:       maxAmount,
		CalculationBase: rule.CalculationBase,
		EffectiveFrom:   rule.EffectiveFrom.Format(time.RFC3339),
		EffectiveTo:     effectiveTo,
		Priority:        rule.Priority,
		IsActive:        rule.IsActive,
		CreatedAt:       rule.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       rule.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *CommissionHandler) periodToResponse(period *entity.CommissionPeriod) dto.CommissionPeriodResponse {
	var closedAt, paidAt *string
	if period.ClosedAt != nil {
		v := period.ClosedAt.Format(time.RFC3339)
		closedAt = &v
	}
	if period.PaidAt != nil {
		v := period.PaidAt.Format(time.RFC3339)
		paidAt = &v
	}

	return dto.CommissionPeriodResponse{
		ID:               period.ID,
		TenantID:         period.TenantID.String(),
		UnitID:           period.UnitID,
		ReferenceMonth:   period.ReferenceMonth,
		ProfessionalID:   period.ProfessionalID,
		TotalGross:       period.TotalGross.String(),
		TotalCommission:  period.TotalCommission.String(),
		TotalAdvances:    period.TotalAdvances.String(),
		TotalAdjustments: period.TotalAdjustments.String(),
		TotalNet:         period.TotalNet.String(),
		ItemsCount:       period.ItemsCount,
		Status:           period.Status,
		PeriodStart:      period.PeriodStart.Format(time.RFC3339),
		PeriodEnd:        period.PeriodEnd.Format(time.RFC3339),
		ClosedAt:         closedAt,
		PaidAt:           paidAt,
		ClosedBy:         period.ClosedBy,
		PaidBy:           period.PaidBy,
		Notes:            period.Notes,
		CreatedAt:        period.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        period.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *CommissionHandler) advanceToResponse(advance *entity.Advance) dto.AdvanceResponse {
	var approvedAt, rejectedAt, deductedAt *string
	if advance.ApprovedAt != nil {
		v := advance.ApprovedAt.Format(time.RFC3339)
		approvedAt = &v
	}
	if advance.RejectedAt != nil {
		v := advance.RejectedAt.Format(time.RFC3339)
		rejectedAt = &v
	}
	if advance.DeductedAt != nil {
		v := advance.DeductedAt.Format(time.RFC3339)
		deductedAt = &v
	}

	return dto.AdvanceResponse{
		ID:                advance.ID,
		TenantID:          advance.TenantID.String(),
		UnitID:            advance.UnitID,
		ProfessionalID:    advance.ProfessionalID,
		Amount:            advance.Amount.String(),
		RequestDate:       advance.RequestDate.Format(time.RFC3339),
		Reason:            advance.Reason,
		Status:            advance.Status,
		ApprovedAt:        approvedAt,
		ApprovedBy:        advance.ApprovedBy,
		RejectedAt:        rejectedAt,
		RejectedBy:        advance.RejectedBy,
		RejectionReason:   advance.RejectionReason,
		DeductedAt:        deductedAt,
		DeductionPeriodID: advance.DeductionPeriodID,
		CreatedAt:         advance.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         advance.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *CommissionHandler) itemToResponse(item *entity.CommissionItem) dto.CommissionItemResponse {
	var processedAt *string
	if item.ProcessedAt != nil {
		v := item.ProcessedAt.Format(time.RFC3339)
		processedAt = &v
	}

	return dto.CommissionItemResponse{
		ID:               item.ID,
		TenantID:         item.TenantID.String(),
		UnitID:           item.UnitID,
		ProfessionalID:   item.ProfessionalID,
		CommandID:        item.CommandID,
		CommandItemID:    item.CommandItemID,
		AppointmentID:    item.AppointmentID,
		ServiceID:        item.ServiceID,
		ServiceName:      item.ServiceName,
		GrossValue:       item.GrossValue.String(),
		CommissionRate:   item.CommissionRate.String(),
		CommissionType:   item.CommissionType,
		CommissionValue:  item.CommissionValue.String(),
		CommissionSource: item.CommissionSource,
		RuleID:           item.RuleID,
		ReferenceDate:    item.ReferenceDate.Format(time.RFC3339),
		Description:      item.Description,
		Status:           item.Status,
		PeriodID:         item.PeriodID,
		ProcessedAt:      processedAt,
		CreatedAt:        item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        item.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *CommissionHandler) summaryToResponse(summary *entity.CommissionSummary) dto.CommissionSummaryResponse {
	return dto.CommissionSummaryResponse{
		ProfessionalID:   summary.ProfessionalID,
		ProfessionalName: summary.ProfessionalName,
		TotalGross:       summary.TotalGross.String(),
		TotalCommission:  summary.TotalCommission.String(),
		ItemsCount:       summary.ItemsCount,
	}
}

func (h *CommissionHandler) byServiceToResponse(summary *entity.CommissionByService) dto.CommissionByServiceResponse {
	return dto.CommissionByServiceResponse{
		ServiceID:       summary.ServiceID,
		ServiceName:     summary.ServiceName,
		TotalGross:      summary.TotalGross.String(),
		TotalCommission: summary.TotalCommission.String(),
		ItemsCount:      summary.ItemsCount,
	}
}

// formatMoney formata valor decimal como string de dinheiro
func formatMoney(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

// RegisterRoutes registra todas as rotas de comissões
func (h *CommissionHandler) RegisterRoutes(g *echo.Group) {
	// Commission Rules
	rules := g.Group("/rules")
	rules.POST("", h.CreateCommissionRule)
	rules.GET("", h.ListCommissionRules)
	rules.GET("/:id", h.GetCommissionRule)
	rules.GET("/effective", h.GetEffectiveCommissionRules)
	rules.PUT("/:id", h.UpdateCommissionRule)
	rules.DELETE("/:id", h.DeleteCommissionRule)
	rules.POST("/:id/deactivate", h.DeactivateCommissionRule)

	// Commission Periods
	periods := g.Group("/periods")
	periods.POST("", h.CreateCommissionPeriod)
	periods.GET("", h.ListCommissionPeriods)
	periods.GET("/:id", h.GetCommissionPeriod)
	periods.GET("/:id/summary", h.GetCommissionPeriodSummary)
	periods.GET("/open/:professional_id", h.GetOpenCommissionPeriod)
	periods.POST("/:id/close", h.CloseCommissionPeriod)
	periods.POST("/:id/pay", h.MarkPeriodAsPaid)
	periods.DELETE("/:id", h.DeleteCommissionPeriod)

	// Advances
	advances := g.Group("/advances")
	advances.POST("", h.CreateAdvance)
	advances.GET("", h.ListAdvances)
	advances.GET("/:id", h.GetAdvance)
	advances.GET("/pending/:professional_id", h.GetPendingAdvances)
	advances.GET("/approved/:professional_id", h.GetApprovedAdvances)
	advances.POST("/:id/approve", h.ApproveAdvance)
	advances.POST("/:id/reject", h.RejectAdvance)
	advances.POST("/:id/deduct", h.MarkAdvanceDeducted)
	advances.POST("/:id/cancel", h.CancelAdvance)
	advances.DELETE("/:id", h.DeleteAdvance)

	// Commission Items
	items := g.Group("/items")
	items.POST("", h.CreateCommissionItem)
	items.POST("/batch", h.CreateCommissionItemBatch)
	items.GET("", h.ListCommissionItems)
	items.GET("/:id", h.GetCommissionItem)
	items.GET("/pending/:professional_id", h.GetPendingCommissionItems)
	items.POST("/:id/process", h.ProcessCommissionItem)
	items.POST("/assign", h.AssignItemsToPeriod)
	items.DELETE("/:id", h.DeleteCommissionItem)

	// Summaries
	summary := g.Group("/summary")
	summary.GET("/by-professional", h.GetCommissionSummaryByProfessional)
	summary.GET("/by-service", h.GetCommissionSummaryByService)
}
