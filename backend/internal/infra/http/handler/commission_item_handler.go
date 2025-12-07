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
// Commission Item Handlers
// =============================================================================

// CreateCommissionItem godoc
// @Summary Criar item de comissão
// @Description Cria um novo item de comissão
// @Tags Comissões
// @Accept json
// @Produce json
// @Param request body dto.CreateCommissionItemRequest true "Dados do item"
// @Success 201 {object} dto.CommissionItemResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/commissions/items [post]
// @Security BearerAuth
func (h *CommissionHandler) CreateCommissionItem(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.CreateCommissionItemRequest
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

	// Parse data de referência
	referenceDate, err := time.Parse(time.RFC3339, req.ReferenceDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Data reference_date inválida",
		})
	}

	output, err := h.createItemUC.Execute(ctx, commission.CreateCommissionItemInput{
		TenantID:         tenantID,
		UnitID:           req.UnitID,
		ProfessionalID:   req.ProfessionalID,
		CommandID:        req.CommandID,
		CommandItemID:    req.CommandItemID,
		AppointmentID:    req.AppointmentID,
		ServiceID:        req.ServiceID,
		ServiceName:      req.ServiceName,
		GrossValue:       req.GrossValue,
		CommissionRate:   req.CommissionRate,
		CommissionType:   req.CommissionType,
		CommissionSource: req.CommissionSource,
		RuleID:           req.RuleID,
		ReferenceDate:    referenceDate,
		Description:      req.Description,
	})

	if err != nil {
		h.logger.Error("Erro ao criar item de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusCreated, h.itemToResponse(output.CommissionItem))
}

// CreateCommissionItemBatch godoc
// @Summary Criar múltiplos itens de comissão
// @Description Cria múltiplos itens de comissão de uma vez
// @Tags Comissões
// @Accept json
// @Produce json
// @Param request body dto.CreateCommissionItemBatchRequest true "Dados dos itens"
// @Success 201 {object} dto.ListCommissionItemsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/commissions/items/batch [post]
// @Security BearerAuth
func (h *CommissionHandler) CreateCommissionItemBatch(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.CreateCommissionItemBatchRequest
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

	// Converte requests para inputs
	inputs := make([]commission.CreateCommissionItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		referenceDate, err := time.Parse(time.RFC3339, item.ReferenceDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Data reference_date inválida em um dos itens",
			})
		}

		inputs = append(inputs, commission.CreateCommissionItemInput{
			TenantID:         tenantID,
			UnitID:           item.UnitID,
			ProfessionalID:   item.ProfessionalID,
			CommandID:        item.CommandID,
			CommandItemID:    item.CommandItemID,
			AppointmentID:    item.AppointmentID,
			ServiceID:        item.ServiceID,
			ServiceName:      item.ServiceName,
			GrossValue:       item.GrossValue,
			CommissionRate:   item.CommissionRate,
			CommissionType:   item.CommissionType,
			CommissionSource: item.CommissionSource,
			RuleID:           item.RuleID,
			ReferenceDate:    referenceDate,
			Description:      item.Description,
		})
	}

	output, err := h.createItemBatchUC.Execute(ctx, commission.CreateCommissionItemBatchInput{
		Items: inputs,
	})

	if err != nil {
		h.logger.Error("Erro ao criar itens de comissão em lote", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	items := make([]dto.CommissionItemResponse, 0, len(output.CommissionItems))
	for _, item := range output.CommissionItems {
		items = append(items, h.itemToResponse(item))
	}

	return c.JSON(http.StatusCreated, dto.ListCommissionItemsResponse{
		Data:  items,
		Total: len(items),
	})
}

// GetCommissionItem godoc
// @Summary Buscar item de comissão
// @Description Busca um item de comissão por ID
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do item"
// @Success 200 {object} dto.CommissionItemResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/commissions/items/{id} [get]
// @Security BearerAuth
func (h *CommissionHandler) GetCommissionItem(c echo.Context) error {
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
			Message: "ID do item é obrigatório",
		})
	}

	output, err := h.getItemUC.Execute(ctx, commission.GetCommissionItemInput{
		TenantID: tenantID,
		ID:       id,
	})

	if err != nil {
		h.logger.Error("Erro ao buscar item de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	if output.CommissionItem == nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Item de comissão não encontrado",
		})
	}

	return c.JSON(http.StatusOK, h.itemToResponse(output.CommissionItem))
}

// ListCommissionItems godoc
// @Summary Listar itens de comissão
// @Description Lista itens de comissão com filtros
// @Tags Comissões
// @Accept json
// @Produce json
// @Param professional_id query string false "ID do profissional"
// @Param period_id query string false "ID do período"
// @Param status query string false "Status do item"
// @Param limit query int false "Limite de resultados"
// @Param offset query int false "Offset para paginação"
// @Success 200 {object} dto.ListCommissionItemsResponse
// @Router /api/v1/commissions/items [get]
// @Security BearerAuth
func (h *CommissionHandler) ListCommissionItems(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ListCommissionItemsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	// T-SEC-002: RBAC - Barbeiro só vê suas próprias comissões
	professionalID := req.ProfessionalID
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		// Se barbeiro tentar filtrar por outro profissional, negar acesso
		if professionalID != nil && *professionalID != barberProfID {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "Acesso negado: você só pode ver suas próprias comissões",
			})
		}
		// Forçar filtro para o profissional atual
		professionalID = &barberProfID
	}

	output, err := h.listItemsUC.Execute(ctx, commission.ListCommissionItemsInput{
		TenantID:       tenantID,
		ProfessionalID: professionalID,
		PeriodID:       req.PeriodID,
		Status:         req.Status,
		Limit:          req.Limit,
		Offset:         req.Offset,
	})

	if err != nil {
		h.logger.Error("Erro ao listar itens de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	items := make([]dto.CommissionItemResponse, 0, len(output.CommissionItems))
	for _, item := range output.CommissionItems {
		items = append(items, h.itemToResponse(item))
	}

	return c.JSON(http.StatusOK, dto.ListCommissionItemsResponse{
		Data:  items,
		Total: len(items),
	})
}

// GetPendingCommissionItems godoc
// @Summary Listar itens pendentes de um profissional
// @Description Lista itens de comissão pendentes de processamento
// @Tags Comissões
// @Accept json
// @Produce json
// @Param professional_id path string true "ID do profissional"
// @Success 200 {object} dto.ListCommissionItemsResponse
// @Router /api/v1/commissions/items/pending/{professional_id} [get]
// @Security BearerAuth
func (h *CommissionHandler) GetPendingCommissionItems(c echo.Context) error {
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

	// T-SEC-002: RBAC - Barbeiro só vê suas próprias comissões pendentes
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		if professionalID != barberProfID {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "Acesso negado: você só pode ver suas próprias comissões",
			})
		}
	}

	output, err := h.getPendingItemsUC.Execute(ctx, commission.GetPendingCommissionItemsInput{
		TenantID:       tenantID,
		ProfessionalID: professionalID,
	})

	if err != nil {
		h.logger.Error("Erro ao listar itens pendentes", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	items := make([]dto.CommissionItemResponse, 0, len(output.CommissionItems))
	for _, item := range output.CommissionItems {
		items = append(items, h.itemToResponse(item))
	}

	return c.JSON(http.StatusOK, dto.ListCommissionItemsResponse{
		Data:  items,
		Total: len(items),
	})
}

// GetCommissionSummaryByProfessional godoc
// @Summary Resumo de comissões por profissional
// @Description Retorna resumo de comissões agrupado por profissional
// @Tags Comissões
// @Accept json
// @Produce json
// @Param start_date query string true "Data inicial (RFC3339)"
// @Param end_date query string true "Data final (RFC3339)"
// @Success 200 {object} dto.CommissionSummariesResponse
// @Router /api/v1/commissions/summary/by-professional [get]
// @Security BearerAuth
func (h *CommissionHandler) GetCommissionSummaryByProfessional(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	if startDateStr == "" || endDateStr == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Datas start_date e end_date são obrigatórias",
		})
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Data start_date inválida",
		})
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Data end_date inválida",
		})
	}

	output, err := h.getSummaryByProfessionalUC.Execute(ctx, commission.GetCommissionSummaryByProfessionalInput{
		TenantID:  tenantID,
		StartDate: startDate,
		EndDate:   endDate,
	})

	if err != nil {
		h.logger.Error("Erro ao buscar resumo por profissional", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	summaries := make([]dto.CommissionSummaryResponse, 0, len(output.Summaries))

	// T-SEC-002: RBAC - Barbeiro só vê seu próprio resumo
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		for _, s := range output.Summaries {
			if s.ProfessionalID == barberProfID {
				summaries = append(summaries, h.summaryToResponse(s))
				break
			}
		}
	} else {
		for _, s := range output.Summaries {
			summaries = append(summaries, h.summaryToResponse(s))
		}
	}

	return c.JSON(http.StatusOK, dto.CommissionSummariesResponse{
		ByProfessional: summaries,
		StartDate:      startDate,
		EndDate:        endDate,
	})
}

// GetCommissionSummaryByService godoc
// @Summary Resumo de comissões por serviço
// @Description Retorna resumo de comissões agrupado por serviço
// @Tags Comissões
// @Accept json
// @Produce json
// @Param start_date query string true "Data inicial (RFC3339)"
// @Param end_date query string true "Data final (RFC3339)"
// @Success 200 {object} dto.CommissionSummariesResponse
// @Router /api/v1/commissions/summary/by-service [get]
// @Security BearerAuth
func (h *CommissionHandler) GetCommissionSummaryByService(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	if startDateStr == "" || endDateStr == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Datas start_date e end_date são obrigatórias",
		})
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Data start_date inválida",
		})
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Data end_date inválida",
		})
	}

	output, err := h.getSummaryByServiceUC.Execute(ctx, commission.GetCommissionSummaryByServiceInput{
		TenantID:  tenantID,
		StartDate: startDate,
		EndDate:   endDate,
	})

	if err != nil {
		h.logger.Error("Erro ao buscar resumo por serviço", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	summaries := make([]dto.CommissionByServiceResponse, 0, len(output.Summaries))
	for _, s := range output.Summaries {
		summaries = append(summaries, h.byServiceToResponse(s))
	}

	return c.JSON(http.StatusOK, dto.CommissionSummariesResponse{
		ByService: summaries,
		StartDate: startDate,
		EndDate:   endDate,
	})
}

// ProcessCommissionItem godoc
// @Summary Processar item de comissão
// @Description Processa um item de comissão vinculando a um período
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do item"
// @Param request body dto.ProcessCommissionItemRequest true "Dados do processamento"
// @Success 200 {object} dto.CommissionItemResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/commissions/items/{id}/process [post]
// @Security BearerAuth
func (h *CommissionHandler) ProcessCommissionItem(c echo.Context) error {
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
			Message: "ID do item é obrigatório",
		})
	}

	var req dto.ProcessCommissionItemRequest
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

	output, err := h.processItemUC.Execute(ctx, commission.ProcessCommissionItemInput{
		TenantID: tenantID,
		ItemID:   id,
		PeriodID: req.PeriodID,
	})

	if err != nil {
		h.logger.Error("Erro ao processar item de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusOK, h.itemToResponse(output.CommissionItem))
}

// AssignItemsToPeriod godoc
// @Summary Vincular itens a um período
// @Description Vincula todos os itens pendentes de um profissional a um período
// @Tags Comissões
// @Accept json
// @Produce json
// @Param request body dto.AssignItemsToPeriodRequest true "Dados da vinculação"
// @Success 200 {object} map[string]int64
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/commissions/items/assign [post]
// @Security BearerAuth
func (h *CommissionHandler) AssignItemsToPeriod(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.AssignItemsToPeriodRequest
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
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Data start_date inválida",
		})
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Data end_date inválida",
		})
	}

	output, err := h.assignItemsUC.Execute(ctx, commission.AssignItemsToPeriodInput{
		TenantID:       tenantID,
		ProfessionalID: req.ProfessionalID,
		PeriodID:       req.PeriodID,
		StartDate:      startDate,
		EndDate:        endDate,
	})

	if err != nil {
		h.logger.Error("Erro ao vincular itens ao período", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]int64{
		"items_assigned": output.ItemsAssigned,
	})
}

// DeleteCommissionItem godoc
// @Summary Excluir item de comissão
// @Description Exclui um item de comissão (apenas se estiver pendente)
// @Tags Comissões
// @Accept json
// @Produce json
// @Param id path string true "ID do item"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/commissions/items/{id} [delete]
// @Security BearerAuth
func (h *CommissionHandler) DeleteCommissionItem(c echo.Context) error {
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
			Message: "ID do item é obrigatório",
		})
	}

	_, err := h.deleteItemUC.Execute(ctx, commission.DeleteCommissionItemInput{
		TenantID: tenantID,
		ItemID:   id,
	})

	if err != nil {
		h.logger.Error("Erro ao excluir item de comissão", zap.Error(err))
		return h.handleDomainError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
