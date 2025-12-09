// Package handler contém os HTTP handlers da camada de infraestrutura.
package handler

import (
	"net/http"
	"strconv"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/financial"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// DespesaFixaHandler agrupa os handlers de despesas fixas
type DespesaFixaHandler struct {
	createUC        *financial.CreateDespesaFixaUseCase
	getUC           *financial.GetDespesaFixaUseCase
	listUC          *financial.ListDespesasFixasUseCase
	updateUC        *financial.UpdateDespesaFixaUseCase
	toggleUC        *financial.ToggleDespesaFixaUseCase
	deleteUC        *financial.DeleteDespesaFixaUseCase
	gerarContasUC   *financial.GerarContasFromDespesasFixasUseCase
	despesaFixaRepo port.DespesaFixaRepository
	mapper          *mapper.DespesaFixaMapper
	logger          *zap.Logger
}

// NewDespesaFixaHandler cria um novo handler de despesas fixas
func NewDespesaFixaHandler(
	createUC *financial.CreateDespesaFixaUseCase,
	getUC *financial.GetDespesaFixaUseCase,
	listUC *financial.ListDespesasFixasUseCase,
	updateUC *financial.UpdateDespesaFixaUseCase,
	toggleUC *financial.ToggleDespesaFixaUseCase,
	deleteUC *financial.DeleteDespesaFixaUseCase,
	gerarContasUC *financial.GerarContasFromDespesasFixasUseCase,
	despesaFixaRepo port.DespesaFixaRepository,
	logger *zap.Logger,
) *DespesaFixaHandler {
	return &DespesaFixaHandler{
		createUC:        createUC,
		getUC:           getUC,
		listUC:          listUC,
		updateUC:        updateUC,
		toggleUC:        toggleUC,
		deleteUC:        deleteUC,
		gerarContasUC:   gerarContasUC,
		despesaFixaRepo: despesaFixaRepo,
		mapper:          mapper.NewDespesaFixaMapper(),
		logger:          logger,
	}
}

// Create godoc
// @Summary Criar despesa fixa
// @Description Cria uma nova despesa fixa recorrente
// @Tags Despesas Fixas
// @Accept json
// @Produce json
// @Param request body dto.CreateDespesaFixaRequest true "Dados da despesa fixa"
// @Success 201 {object} dto.DespesaFixaResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/fixed-expenses [post]
// @Security BearerAuth
func (h *DespesaFixaHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.CreateDespesaFixaRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind do request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Dados inválidos no corpo da requisição",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	input, err := h.mapper.ToCreateInput(req, tenantID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_value",
			Message: "Valor inválido: " + err.Error(),
		})
	}

	despesa, err := h.createUC.Execute(ctx, *input)
	if err != nil {
		h.logger.Error("Erro ao criar despesa fixa", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao criar despesa fixa",
		})
	}

	return c.JSON(http.StatusCreated, h.mapper.ToResponse(despesa))
}

// GetByID godoc
// @Summary Obter despesa fixa por ID
// @Description Retorna uma despesa fixa específica
// @Tags Despesas Fixas
// @Produce json
// @Param id path string true "ID da despesa fixa"
// @Success 200 {object} dto.DespesaFixaResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/fixed-expenses/{id} [get]
// @Security BearerAuth
func (h *DespesaFixaHandler) GetByID(c echo.Context) error {
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
			Error:   "invalid_id",
			Message: "ID da despesa fixa é obrigatório",
		})
	}

	despesa, err := h.getUC.Execute(ctx, financial.GetDespesaFixaInput{
		TenantID: tenantID,
		ID:       id,
	})
	if err != nil {
		if err == domain.ErrDespesaFixaNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Despesa fixa não encontrada",
			})
		}
		h.logger.Error("Erro ao buscar despesa fixa", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar despesa fixa",
		})
	}

	return c.JSON(http.StatusOK, h.mapper.ToResponse(despesa))
}

// List godoc
// @Summary Listar despesas fixas
// @Description Retorna lista paginada de despesas fixas
// @Tags Despesas Fixas
// @Produce json
// @Param page query int false "Página (default: 1)"
// @Param page_size query int false "Itens por página (default: 20)"
// @Param ativo query bool false "Filtrar por status ativo"
// @Param categoria_id query string false "Filtrar por categoria"
// @Param unidade_id query string false "Filtrar por unidade"
// @Success 200 {object} dto.DespesasFixasListResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/fixed-expenses [get]
// @Security BearerAuth
func (h *DespesaFixaHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Parse query params
	page := 1
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 20
	if ps := c.QueryParam("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	input := financial.ListDespesasFixasInput{
		TenantID: tenantID,
		Page:     page,
		PageSize: pageSize,
	}

	// Filtros opcionais
	if ativo := c.QueryParam("ativo"); ativo != "" {
		val := ativo == "true"
		input.Ativo = &val
	}
	if cat := c.QueryParam("categoria_id"); cat != "" {
		input.CategoriaID = &cat
	}
	if uni := c.QueryParam("unidade_id"); uni != "" {
		input.UnidadeID = &uni
	}

	output, err := h.listUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao listar despesas fixas", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar despesas fixas",
		})
	}

	return c.JSON(http.StatusOK, h.mapper.ToListResponse(output.Despesas, output.Total, output.Page, output.PageSize))
}

// Update godoc
// @Summary Atualizar despesa fixa
// @Description Atualiza uma despesa fixa existente
// @Tags Despesas Fixas
// @Accept json
// @Produce json
// @Param id path string true "ID da despesa fixa"
// @Param request body dto.UpdateDespesaFixaRequest true "Dados atualizados"
// @Success 200 {object} dto.DespesaFixaResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/fixed-expenses/{id} [put]
// @Security BearerAuth
func (h *DespesaFixaHandler) Update(c echo.Context) error {
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
			Error:   "invalid_id",
			Message: "ID da despesa fixa é obrigatório",
		})
	}

	var req dto.UpdateDespesaFixaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Dados inválidos no corpo da requisição",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	input, err := h.mapper.ToUpdateInput(req, tenantID, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_value",
			Message: "Valor inválido: " + err.Error(),
		})
	}

	despesa, err := h.updateUC.Execute(ctx, *input)
	if err != nil {
		if err == domain.ErrDespesaFixaNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Despesa fixa não encontrada",
			})
		}
		h.logger.Error("Erro ao atualizar despesa fixa", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao atualizar despesa fixa",
		})
	}

	return c.JSON(http.StatusOK, h.mapper.ToResponse(despesa))
}

// Toggle godoc
// @Summary Alternar status ativo/inativo
// @Description Ativa ou desativa uma despesa fixa
// @Tags Despesas Fixas
// @Produce json
// @Param id path string true "ID da despesa fixa"
// @Success 200 {object} dto.DespesaFixaResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/fixed-expenses/{id}/toggle [patch]
// @Security BearerAuth
func (h *DespesaFixaHandler) Toggle(c echo.Context) error {
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
			Error:   "invalid_id",
			Message: "ID da despesa fixa é obrigatório",
		})
	}

	despesa, err := h.toggleUC.Execute(ctx, financial.ToggleDespesaFixaInput{
		TenantID: tenantID,
		ID:       id,
	})
	if err != nil {
		if err == domain.ErrDespesaFixaNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Despesa fixa não encontrada",
			})
		}
		h.logger.Error("Erro ao alternar despesa fixa", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao alternar despesa fixa",
		})
	}

	return c.JSON(http.StatusOK, h.mapper.ToResponse(despesa))
}

// Delete godoc
// @Summary Excluir despesa fixa
// @Description Remove uma despesa fixa
// @Tags Despesas Fixas
// @Produce json
// @Param id path string true "ID da despesa fixa"
// @Success 204 "No Content"
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/fixed-expenses/{id} [delete]
// @Security BearerAuth
func (h *DespesaFixaHandler) Delete(c echo.Context) error {
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
			Error:   "invalid_id",
			Message: "ID da despesa fixa é obrigatório",
		})
	}

	err := h.deleteUC.Execute(ctx, financial.DeleteDespesaFixaInput{
		TenantID: tenantID,
		ID:       id,
	})
	if err != nil {
		if err == domain.ErrDespesaFixaNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Despesa fixa não encontrada",
			})
		}
		h.logger.Error("Erro ao excluir despesa fixa", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao excluir despesa fixa",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetSummary godoc
// @Summary Obter resumo de despesas fixas
// @Description Retorna totais e soma de valores das despesas fixas
// @Tags Despesas Fixas
// @Produce json
// @Success 200 {object} dto.DespesasFixasSummaryResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/fixed-expenses/summary [get]
// @Security BearerAuth
func (h *DespesaFixaHandler) GetSummary(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	total, err := h.despesaFixaRepo.Count(ctx, tenantID)
	if err != nil {
		h.logger.Error("Erro ao contar despesas fixas", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao obter resumo",
		})
	}

	totalAtivas, err := h.despesaFixaRepo.CountAtivas(ctx, tenantID)
	if err != nil {
		h.logger.Error("Erro ao contar despesas ativas", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao obter resumo",
		})
	}

	valorTotal, err := h.despesaFixaRepo.SumAtivas(ctx, tenantID)
	if err != nil {
		h.logger.Error("Erro ao somar despesas ativas", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao obter resumo",
		})
	}

	return c.JSON(http.StatusOK, h.mapper.ToSummaryResponse(total, totalAtivas, valorTotal))
}

// GenerateContas godoc
// @Summary Gerar contas a pagar a partir de despesas fixas
// @Description Cria contas a pagar para o mês especificado a partir das despesas fixas ativas
// @Tags Despesas Fixas
// @Accept json
// @Produce json
// @Param request body dto.GerarContasRequest true "Mês/Ano para geração"
// @Success 200 {object} dto.GerarContasResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/fixed-expenses/generate [post]
// @Security BearerAuth
func (h *DespesaFixaHandler) GenerateContas(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.GerarContasRequest
	if err := c.Bind(&req); err != nil {
		// Se não houver body, usar mês atual
		req = dto.GerarContasRequest{}
	}

	output, err := h.gerarContasUC.Execute(ctx, financial.GerarContasFromDespesasFixasInput{
		TenantID: tenantID,
		Ano:      req.Ano,
		Mes:      req.Mes,
	})
	if err != nil {
		h.logger.Error("Erro ao gerar contas", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao gerar contas a pagar",
		})
	}

	return c.JSON(http.StatusOK, h.mapper.ToGerarContasResponse(output))
}

// RegisterRoutes registra as rotas de despesas fixas
func (h *DespesaFixaHandler) RegisterRoutes(g *echo.Group) {
	// CRUD básico
	g.POST("/fixed-expenses", h.Create)
	g.GET("/fixed-expenses", h.List)
	g.GET("/fixed-expenses/summary", h.GetSummary)
	g.GET("/fixed-expenses/:id", h.GetByID)
	g.PUT("/fixed-expenses/:id", h.Update)
	g.PATCH("/fixed-expenses/:id/toggle", h.Toggle)
	g.DELETE("/fixed-expenses/:id", h.Delete)

	// Geração automática
	g.POST("/fixed-expenses/generate", h.GenerateContas)
}
