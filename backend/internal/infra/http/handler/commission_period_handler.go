package handler

import (
	"net/http"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/commission"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// =============================================================================
// Commission Period Handlers
// =============================================================================

// CreateCommissionPeriod godoc
// @Summary Criar período de comissão
// @Description Cria um novo período de comissão para um profissional
// @Tags Comissões
// @Accept json
// @Produce json
// @Param request body dto.CreateCommissionPeriodRequest true "Dados do período"
// @Success 201 {object} dto.CommissionPeriodResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/commissions/periods [post]
// @Security BearerAuth
func (h *CommissionHandler) CreateCommissionPeriod(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.CreateCommissionPeriodRequest
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

	// Parse datas
	periodStart, err := time.Parse(time.RFC3339, req.PeriodStart)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Data period_start inválida",
		})
	}
	periodEnd, err := time.Parse(time.RFC3339, req.PeriodEnd)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Data period_end inválida",
		})
	}

	output, err := h.createPeriodUC.Execute(ctx, commission.CreateCommissionPeriodInput{
		TenantID:       tenantID,
		UnitID:         req.UnitID,
		ReferenceMonth: req.ReferenceMonth,
		ProfessionalID: req.ProfessionalID,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		Notes:          req.Notes,
	})

	if err != nil {
		h.logger.Error("Erro ao criar período de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusCreated, h.periodToResponse(output.CommissionPeriod))
}

// GetCommissionPeriod godoc
// @Summary Buscar período de comissão
// @Description Busca um período de comissão por ID
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do período"
// @Success 200 {object} dto.CommissionPeriodResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/commissions/periods/{id} [get]
// @Security BearerAuth
func (h *CommissionHandler) GetCommissionPeriod(c echo.Context) error {
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
			Message: "ID do período é obrigatório",
		})
	}

	output, err := h.getPeriodUC.Execute(ctx, commission.GetCommissionPeriodInput{
		TenantID: tenantID,
		ID:       id,
	})

	if err != nil {
		h.logger.Error("Erro ao buscar período de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	if output.CommissionPeriod == nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Período de comissão não encontrado",
		})
	}

	return c.JSON(http.StatusOK, h.periodToResponse(output.CommissionPeriod))
}

// GetCommissionPeriodSummary godoc
// @Summary Resumo de um período de comissão
// @Description Retorna o resumo com totais de um período de comissão
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do período"
// @Success 200 {object} dto.CommissionPeriodSummaryResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/commissions/periods/{id}/summary [get]
// @Security BearerAuth
func (h *CommissionHandler) GetCommissionPeriodSummary(c echo.Context) error {
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
			Message: "ID do período é obrigatório",
		})
	}

	output, err := h.getPeriodSummaryUC.Execute(ctx, commission.GetCommissionPeriodSummaryInput{
		TenantID: tenantID,
		PeriodID: id,
	})

	if err != nil {
		h.logger.Error("Erro ao buscar resumo do período", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	if output.Summary == nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Período de comissão não encontrado",
		})
	}

	return c.JSON(http.StatusOK, dto.CommissionPeriodSummaryResponse{
		TotalGross:      output.Summary.TotalGross.String(),
		TotalCommission: output.Summary.TotalCommission.String(),
		TotalAdvances:   output.Summary.TotalAdvances.String(),
		TotalNet:        output.Summary.TotalNet.String(),
		ItemsCount:      output.Summary.ItemsCount,
	})
}

// GetOpenCommissionPeriod godoc
// @Summary Buscar período aberto de um profissional
// @Description Busca o período de comissão aberto de um profissional
// @Tags Comissões
// @Accept json
// @Produce json
// @Param professional_id path string true "ID do profissional"
// @Success 200 {object} dto.CommissionPeriodResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/commissions/periods/open/{professional_id} [get]
// @Security BearerAuth
func (h *CommissionHandler) GetOpenCommissionPeriod(c echo.Context) error {
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

	// T-SEC-002: RBAC - Barbeiro só vê seu próprio período aberto
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		if professionalID != barberProfID {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "Acesso negado: você só pode ver seus próprios períodos de comissão",
			})
		}
	}

	output, err := h.getOpenPeriodUC.Execute(ctx, commission.GetOpenCommissionPeriodInput{
		TenantID:       tenantID,
		ProfessionalID: professionalID,
	})

	if err != nil {
		h.logger.Error("Erro ao buscar período aberto", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	if output.CommissionPeriod == nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Nenhum período aberto encontrado para este profissional",
		})
	}

	return c.JSON(http.StatusOK, h.periodToResponse(output.CommissionPeriod))
}

// ListCommissionPeriods godoc
// @Summary Listar períodos de comissão
// @Description Lista períodos de comissão com filtros
// @Tags Comissões
// @Accept json
// @Produce json
// @Param professional_id query string false "ID do profissional"
// @Param status query string false "Status do período"
// @Param limit query int false "Limite de resultados"
// @Param offset query int false "Offset para paginação"
// @Success 200 {object} dto.ListCommissionPeriodsResponse
// @Router /api/v1/commissions/periods [get]
// @Security BearerAuth
func (h *CommissionHandler) ListCommissionPeriods(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ListCommissionPeriodsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	// T-SEC-002: RBAC - Barbeiro só vê seus próprios períodos
	professionalID := req.ProfessionalID
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		if professionalID != nil && *professionalID != barberProfID {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "Acesso negado: você só pode ver seus próprios períodos de comissão",
			})
		}
		professionalID = &barberProfID
	}

	output, err := h.listPeriodsUC.Execute(ctx, commission.ListCommissionPeriodsInput{
		TenantID:       tenantID,
		ProfessionalID: professionalID,
		Status:         req.Status,
		Limit:          req.Limit,
		Offset:         req.Offset,
	})

	if err != nil {
		h.logger.Error("Erro ao listar períodos de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	periods := make([]dto.CommissionPeriodResponse, 0, len(output.CommissionPeriods))
	for _, period := range output.CommissionPeriods {
		periods = append(periods, h.periodToResponse(period))
	}

	return c.JSON(http.StatusOK, dto.ListCommissionPeriodsResponse{
		Data:  periods,
		Total: len(periods),
	})
}

// CloseCommissionPeriod godoc
// @Summary Fechar período de comissão
// @Description Fecha um período de comissão
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do período"
// @Success 200 {object} dto.CommissionPeriodResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/commissions/periods/{id}/close [post]
// @Security BearerAuth
func (h *CommissionHandler) CloseCommissionPeriod(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	userID := ""
	if uid, ok := c.Get("user_id").(string); ok {
		userID = uid
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do período é obrigatório",
		})
	}

	output, err := h.closePeriodUC.Execute(ctx, commission.CloseCommissionPeriodInput{
		TenantID: tenantID,
		PeriodID: id,
		ClosedBy: userID,
	})

	if err != nil {
		h.logger.Error("Erro ao fechar período de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusOK, h.periodToResponse(output.CommissionPeriod))
}

// MarkPeriodAsPaid godoc
// @Summary Marcar período como pago
// @Description Marca um período de comissão como pago
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do período"
// @Success 200 {object} dto.CommissionPeriodResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/commissions/periods/{id}/pay [post]
// @Security BearerAuth
func (h *CommissionHandler) MarkPeriodAsPaid(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	userID := ""
	if uid, ok := c.Get("user_id").(string); ok {
		userID = uid
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do período é obrigatório",
		})
	}

	output, err := h.markPeriodPaidUC.Execute(ctx, commission.MarkPeriodAsPaidInput{
		TenantID: tenantID,
		PeriodID: id,
		PaidBy:   userID,
	})

	if err != nil {
		h.logger.Error("Erro ao marcar período como pago", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusOK, h.periodToResponse(output.CommissionPeriod))
}

// DeleteCommissionPeriod godoc
// @Summary Excluir período de comissão
// @Description Exclui um período de comissão (apenas se estiver aberto)
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do período"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/commissions/periods/{id} [delete]
// @Security BearerAuth
func (h *CommissionHandler) DeleteCommissionPeriod(c echo.Context) error {
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
			Message: "ID do período é obrigatório",
		})
	}

	_, err := h.deletePeriodUC.Execute(ctx, commission.DeleteCommissionPeriodInput{
		TenantID: tenantID,
		PeriodID: id,
	})

	if err != nil {
		h.logger.Error("Erro ao excluir período de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
