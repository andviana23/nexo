package handler

import (
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/categoriaproduto"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// CategoriaProdutoHandler gerencia os endpoints de categorias de produtos
type CategoriaProdutoHandler struct {
	createUC *categoriaproduto.CreateCategoriaProdutoUseCase
	listUC   *categoriaproduto.ListCategoriasProdutosUseCase
	getUC    *categoriaproduto.GetCategoriaProdutoUseCase
	updateUC *categoriaproduto.UpdateCategoriaProdutoUseCase
	deleteUC *categoriaproduto.DeleteCategoriaProdutoUseCase
	toggleUC *categoriaproduto.ToggleCategoriaProdutoUseCase
	logger   *zap.Logger
}

// NewCategoriaProdutoHandler cria uma nova instância do handler
func NewCategoriaProdutoHandler(
	createUC *categoriaproduto.CreateCategoriaProdutoUseCase,
	listUC *categoriaproduto.ListCategoriasProdutosUseCase,
	getUC *categoriaproduto.GetCategoriaProdutoUseCase,
	updateUC *categoriaproduto.UpdateCategoriaProdutoUseCase,
	deleteUC *categoriaproduto.DeleteCategoriaProdutoUseCase,
	toggleUC *categoriaproduto.ToggleCategoriaProdutoUseCase,
	logger *zap.Logger,
) *CategoriaProdutoHandler {
	return &CategoriaProdutoHandler{
		createUC: createUC,
		listUC:   listUC,
		getUC:    getUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
		toggleUC: toggleUC,
		logger:   logger,
	}
}

// Create cria uma nova categoria de produto
// @Summary Criar categoria de produto
// @Description Cria uma nova categoria de produto customizada
// @Tags Categorias de Produtos
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoriaProdutoRequest true "Dados da categoria"
// @Success 201 {object} dto.CategoriaProdutoResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/categorias-produtos [post]
// @Security BearerAuth
func (h *CategoriaProdutoHandler) Create(c echo.Context) error {
	tenantIDStr, ok := c.Get("tenant_id").(string)
	if !ok || tenantIDStr == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant_id inválido"})
	}

	unitIDStr := middleware.GetUnitID(c)
	if unitIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": domain.ErrUnitIDRequired.Error()})
	}
	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unit_id inválido"})
	}

	var req dto.CreateCategoriaProdutoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	input := categoriaproduto.CreateCategoriaProdutoInput{
		TenantID:    tenantID,
		UnitID:      unitID,
		Nome:        req.Nome,
		Descricao:   req.Descricao,
		Cor:         req.Cor,
		Icone:       req.Icone,
		CentroCusto: req.CentroCusto,
	}

	// Default centro_custo
	if input.CentroCusto == "" {
		input.CentroCusto = "CMV"
	}

	categoria, err := h.createUC.Execute(c.Request().Context(), input)
	if err != nil {
		h.logger.Error("Erro ao criar categoria de produto", zap.Error(err))
		if err.Error() == "já existe uma categoria com este nome" {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, mapper.CategoriaProdutoToResponse(categoria))
}

// List lista todas as categorias de produtos
// @Summary Listar categorias de produtos
// @Description Lista todas as categorias de produtos do tenant
// @Tags Categorias de Produtos
// @Produce json
// @Param ativas query bool false "Filtrar apenas ativas"
// @Success 200 {object} dto.ListCategoriaProdutoResponse
// @Router /api/v1/categorias-produtos [get]
// @Security BearerAuth
func (h *CategoriaProdutoHandler) List(c echo.Context) error {
	tenantIDStr, ok := c.Get("tenant_id").(string)
	if !ok || tenantIDStr == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant_id inválido"})
	}

	unitIDStr := middleware.GetUnitID(c)
	if unitIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": domain.ErrUnitIDRequired.Error()})
	}
	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unit_id inválido"})
	}

	apenasAtivas := c.QueryParam("ativas") == "true"

	categorias, err := h.listUC.Execute(c.Request().Context(), tenantID, unitID, apenasAtivas)
	if err != nil {
		h.logger.Error("Erro ao listar categorias de produtos", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, mapper.CategoriasProdutosToListResponse(categorias))
}

// GetByID busca uma categoria por ID
// @Summary Buscar categoria por ID
// @Description Retorna os detalhes de uma categoria específica
// @Tags Categorias de Produtos
// @Produce json
// @Param id path string true "ID da categoria"
// @Success 200 {object} dto.CategoriaProdutoResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/categorias-produtos/{id} [get]
// @Security BearerAuth
func (h *CategoriaProdutoHandler) GetByID(c echo.Context) error {
	tenantIDStr, ok := c.Get("tenant_id").(string)
	if !ok || tenantIDStr == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant_id inválido"})
	}

	unitIDStr := middleware.GetUnitID(c)
	if unitIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": domain.ErrUnitIDRequired.Error()})
	}
	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unit_id inválido"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	categoria, err := h.getUC.Execute(c.Request().Context(), tenantID, unitID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar categoria de produto", zap.Error(err))
		if err.Error() == "categoria não encontrada" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, mapper.CategoriaProdutoToResponse(categoria))
}

// Update atualiza uma categoria existente
// @Summary Atualizar categoria
// @Description Atualiza os dados de uma categoria existente
// @Tags Categorias de Produtos
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Param request body dto.UpdateCategoriaProdutoRequest true "Dados da categoria"
// @Success 200 {object} dto.CategoriaProdutoResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/categorias-produtos/{id} [put]
// @Security BearerAuth
func (h *CategoriaProdutoHandler) Update(c echo.Context) error {
	tenantIDStr, ok := c.Get("tenant_id").(string)
	if !ok || tenantIDStr == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant_id inválido"})
	}

	unitIDStr := middleware.GetUnitID(c)
	if unitIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": domain.ErrUnitIDRequired.Error()})
	}
	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unit_id inválido"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	var req dto.UpdateCategoriaProdutoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	input := categoriaproduto.UpdateCategoriaProdutoInput{
		TenantID:    tenantID,
		UnitID:      unitID,
		ID:          id,
		Nome:        req.Nome,
		Descricao:   req.Descricao,
		Cor:         req.Cor,
		Icone:       req.Icone,
		CentroCusto: req.CentroCusto,
		Ativa:       req.Ativa,
	}

	// Default centro_custo
	if input.CentroCusto == "" {
		input.CentroCusto = "CMV"
	}

	categoria, err := h.updateUC.Execute(c.Request().Context(), input)
	if err != nil {
		h.logger.Error("Erro ao atualizar categoria de produto", zap.Error(err))
		switch err.Error() {
		case "categoria não encontrada":
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		case "já existe uma categoria com este nome":
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, mapper.CategoriaProdutoToResponse(categoria))
}

// Delete remove uma categoria
// @Summary Excluir categoria
// @Description Remove uma categoria de produto (apenas se não houver produtos vinculados)
// @Tags Categorias de Produtos
// @Param id path string true "ID da categoria"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/categorias-produtos/{id} [delete]
// @Security BearerAuth
func (h *CategoriaProdutoHandler) Delete(c echo.Context) error {
	tenantIDStr, ok := c.Get("tenant_id").(string)
	if !ok || tenantIDStr == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant_id inválido"})
	}

	unitIDStr := middleware.GetUnitID(c)
	if unitIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": domain.ErrUnitIDRequired.Error()})
	}
	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unit_id inválido"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	err = h.deleteUC.Execute(c.Request().Context(), tenantID, unitID, id)
	if err != nil {
		h.logger.Error("Erro ao excluir categoria de produto", zap.Error(err))
		switch err.Error() {
		case "categoria não encontrada":
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		case "não é possível excluir categoria com produtos vinculados":
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// Toggle ativa/desativa uma categoria
// @Summary Ativar/Desativar categoria
// @Description Alterna o status ativo/inativo de uma categoria
// @Tags Categorias de Produtos
// @Produce json
// @Param id path string true "ID da categoria"
// @Success 200 {object} dto.CategoriaProdutoResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/categorias-produtos/{id}/toggle [patch]
// @Security BearerAuth
func (h *CategoriaProdutoHandler) Toggle(c echo.Context) error {
	tenantIDStr, ok := c.Get("tenant_id").(string)
	if !ok || tenantIDStr == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant_id inválido"})
	}

	unitIDStr := middleware.GetUnitID(c)
	if unitIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": domain.ErrUnitIDRequired.Error()})
	}
	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unit_id inválido"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	categoria, err := h.toggleUC.Execute(c.Request().Context(), tenantID, unitID, id)
	if err != nil {
		h.logger.Error("Erro ao alternar status da categoria de produto", zap.Error(err))
		if err.Error() == "categoria não encontrada" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, mapper.CategoriaProdutoToResponse(categoria))
}
