package handler

import (
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/commission"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// =============================================================================
// Advance Handlers
// =============================================================================

// CreateAdvance godoc
// @Summary Criar adiantamento
// @Description Cria um novo pedido de adiantamento para um profissional
// @Tags Comissões
// @Accept json
// @Produce json
// @Param request body dto.CreateAdvanceRequest true "Dados do adiantamento"
// @Success 201 {object} dto.AdvanceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/commissions/advances [post]
// @Security BearerAuth
func (h *CommissionHandler) CreateAdvance(c echo.Context) error {
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

	var req dto.CreateAdvanceRequest
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

	// T-SEC-002: RBAC - Barbeiro só pode criar adiantamento para si mesmo
	professionalID := req.ProfessionalID
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		if professionalID != barberProfID {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "Acesso negado: você só pode solicitar adiantamentos para você mesmo",
			})
		}
	}

	output, err := h.createAdvanceUC.Execute(ctx, commission.CreateAdvanceInput{
		TenantID:       tenantID,
		UnitID:         req.UnitID,
		ProfessionalID: professionalID,
		Amount:         req.Amount,
		Reason:         req.Reason,
		CreatedBy:      userID,
	})

	if err != nil {
		h.logger.Error("Erro ao criar adiantamento", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusCreated, h.advanceToResponse(output.Advance))
}

// GetAdvance godoc
// @Summary Buscar adiantamento
// @Description Busca um adiantamento por ID
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do adiantamento"
// @Success 200 {object} dto.AdvanceResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/commissions/advances/{id} [get]
// @Security BearerAuth
func (h *CommissionHandler) GetAdvance(c echo.Context) error {
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
			Message: "ID do adiantamento é obrigatório",
		})
	}

	output, err := h.getAdvanceUC.Execute(ctx, commission.GetAdvanceInput{
		TenantID: tenantID,
		ID:       id,
	})

	if err != nil {
		h.logger.Error("Erro ao buscar adiantamento", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	if output.Advance == nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Adiantamento não encontrado",
		})
	}

	return c.JSON(http.StatusOK, h.advanceToResponse(output.Advance))
}

// ListAdvances godoc
// @Summary Listar adiantamentos
// @Description Lista adiantamentos com filtros
// @Tags Comissões
// @Accept json
// @Produce json
// @Param professional_id query string false "ID do profissional"
// @Param status query string false "Status do adiantamento"
// @Param limit query int false "Limite de resultados"
// @Param offset query int false "Offset para paginação"
// @Success 200 {object} dto.ListAdvancesResponse
// @Router /api/v1/commissions/advances [get]
// @Security BearerAuth
func (h *CommissionHandler) ListAdvances(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ListAdvancesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	// T-SEC-002: RBAC - Barbeiro só vê seus próprios adiantamentos
	professionalID := req.ProfessionalID
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		if professionalID != nil && *professionalID != barberProfID {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "Acesso negado: você só pode ver seus próprios adiantamentos",
			})
		}
		professionalID = &barberProfID
	}

	output, err := h.listAdvancesUC.Execute(ctx, commission.ListAdvancesInput{
		TenantID:       tenantID,
		ProfessionalID: professionalID,
		Status:         req.Status,
		Limit:          req.Limit,
		Offset:         req.Offset,
	})

	if err != nil {
		h.logger.Error("Erro ao listar adiantamentos", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	advances := make([]dto.AdvanceResponse, 0, len(output.Advances))
	for _, advance := range output.Advances {
		advances = append(advances, h.advanceToResponse(advance))
	}

	return c.JSON(http.StatusOK, dto.ListAdvancesResponse{
		Data:  advances,
		Total: len(advances),
	})
}

// GetPendingAdvances godoc
// @Summary Listar adiantamentos pendentes
// @Description Lista adiantamentos pendentes de aprovação de um profissional
// @Tags Comissões
// @Accept json
// @Produce json
// @Param professional_id path string true "ID do profissional"
// @Success 200 {object} dto.AdvancesTotalsResponse
// @Router /api/v1/commissions/advances/pending/{professional_id} [get]
// @Security BearerAuth
func (h *CommissionHandler) GetPendingAdvances(c echo.Context) error {
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

	// T-SEC-002: RBAC - Barbeiro só vê seus próprios adiantamentos pendentes
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		if professionalID != barberProfID {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "Acesso negado: você só pode ver seus próprios adiantamentos",
			})
		}
	}

	output, err := h.getPendingAdvancesUC.Execute(ctx, commission.GetPendingAdvancesInput{
		TenantID:       tenantID,
		ProfessionalID: professionalID,
	})

	if err != nil {
		h.logger.Error("Erro ao listar adiantamentos pendentes", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	advances := make([]dto.AdvanceResponse, 0, len(output.Advances))
	for _, advance := range output.Advances {
		advances = append(advances, h.advanceToResponse(advance))
	}

	return c.JSON(http.StatusOK, dto.AdvancesTotalsResponse{
		Advances:     advances,
		TotalPending: formatMoney(output.TotalPending),
	})
}

// GetApprovedAdvances godoc
// @Summary Listar adiantamentos aprovados
// @Description Lista adiantamentos aprovados de um profissional (ainda não descontados)
// @Tags Comissões
// @Accept json
// @Produce json
// @Param professional_id path string true "ID do profissional"
// @Success 200 {object} dto.AdvancesTotalsResponse
// @Router /api/v1/commissions/advances/approved/{professional_id} [get]
// @Security BearerAuth
func (h *CommissionHandler) GetApprovedAdvances(c echo.Context) error {
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

	// T-SEC-002: RBAC - Barbeiro só vê seus próprios adiantamentos aprovados
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		if professionalID != barberProfID {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "Acesso negado: você só pode ver seus próprios adiantamentos",
			})
		}
	}

	output, err := h.getApprovedAdvancesUC.Execute(ctx, commission.GetApprovedAdvancesInput{
		TenantID:       tenantID,
		ProfessionalID: professionalID,
	})

	if err != nil {
		h.logger.Error("Erro ao listar adiantamentos aprovados", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	advances := make([]dto.AdvanceResponse, 0, len(output.Advances))
	for _, advance := range output.Advances {
		advances = append(advances, h.advanceToResponse(advance))
	}

	return c.JSON(http.StatusOK, dto.AdvancesTotalsResponse{
		Advances:      advances,
		TotalApproved: formatMoney(output.TotalApproved),
	})
}

// ApproveAdvance godoc
// @Summary Aprovar adiantamento
// @Description Aprova um pedido de adiantamento
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do adiantamento"
// @Success 200 {object} dto.AdvanceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/commissions/advances/{id}/approve [post]
// @Security BearerAuth
func (h *CommissionHandler) ApproveAdvance(c echo.Context) error {
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
			Message: "ID do adiantamento é obrigatório",
		})
	}

	output, err := h.approveAdvanceUC.Execute(ctx, commission.ApproveAdvanceInput{
		TenantID:   tenantID,
		AdvanceID:  id,
		ApprovedBy: userID,
	})

	if err != nil {
		h.logger.Error("Erro ao aprovar adiantamento", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusOK, h.advanceToResponse(output.Advance))
}

// RejectAdvance godoc
// @Summary Rejeitar adiantamento
// @Description Rejeita um pedido de adiantamento
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do adiantamento"
// @Param request body dto.RejectAdvanceRequest true "Motivo da rejeição"
// @Success 200 {object} dto.AdvanceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/commissions/advances/{id}/reject [post]
// @Security BearerAuth
func (h *CommissionHandler) RejectAdvance(c echo.Context) error {
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
			Message: "ID do adiantamento é obrigatório",
		})
	}

	var req dto.RejectAdvanceRequest
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

	output, err := h.rejectAdvanceUC.Execute(ctx, commission.RejectAdvanceInput{
		TenantID:        tenantID,
		AdvanceID:       id,
		RejectedBy:      userID,
		RejectionReason: req.RejectionReason,
	})

	if err != nil {
		h.logger.Error("Erro ao rejeitar adiantamento", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusOK, h.advanceToResponse(output.Advance))
}

// MarkAdvanceDeducted godoc
// @Summary Marcar adiantamento como deduzido
// @Description Marca um adiantamento como deduzido em um período
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do adiantamento"
// @Param request body dto.MarkAdvanceDeductedRequest true "Dados da dedução"
// @Success 200 {object} dto.AdvanceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/commissions/advances/{id}/deduct [post]
// @Security BearerAuth
func (h *CommissionHandler) MarkAdvanceDeducted(c echo.Context) error {
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
			Message: "ID do adiantamento é obrigatório",
		})
	}

	var req dto.MarkAdvanceDeductedRequest
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

	output, err := h.markDeductedUC.Execute(ctx, commission.MarkAdvanceDeductedInput{
		TenantID:  tenantID,
		AdvanceID: id,
		PeriodID:  req.PeriodID,
	})

	if err != nil {
		h.logger.Error("Erro ao marcar adiantamento como deduzido", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusOK, h.advanceToResponse(output.Advance))
}

// CancelAdvance godoc
// @Summary Cancelar adiantamento
// @Description Cancela um pedido de adiantamento
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do adiantamento"
// @Success 200 {object} dto.AdvanceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/commissions/advances/{id}/cancel [post]
// @Security BearerAuth
func (h *CommissionHandler) CancelAdvance(c echo.Context) error {
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
			Message: "ID do adiantamento é obrigatório",
		})
	}

	output, err := h.cancelAdvanceUC.Execute(ctx, commission.CancelAdvanceInput{
		TenantID:  tenantID,
		AdvanceID: id,
	})

	if err != nil {
		h.logger.Error("Erro ao cancelar adiantamento", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusOK, h.advanceToResponse(output.Advance))
}

// DeleteAdvance godoc
// @Summary Excluir adiantamento
// @Description Exclui um adiantamento (apenas se estiver pendente)
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do adiantamento"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/commissions/advances/{id} [delete]
// @Security BearerAuth
func (h *CommissionHandler) DeleteAdvance(c echo.Context) error {
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
			Message: "ID do adiantamento é obrigatório",
		})
	}

	_, err := h.deleteAdvanceUC.Execute(ctx, commission.DeleteAdvanceInput{
		TenantID:  tenantID,
		AdvanceID: id,
	})

	if err != nil {
		h.logger.Error("Erro ao excluir adiantamento", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
