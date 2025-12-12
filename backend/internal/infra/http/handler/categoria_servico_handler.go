package handler

import (
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/categoria"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// CategoriaServicoHandler agrupa os handlers de categorias de serviço
type CategoriaServicoHandler struct {
	createUC *categoria.CreateCategoriaServicoUseCase
	getUC    *categoria.GetCategoriaServicoUseCase
	listUC   *categoria.ListCategoriasServicosUseCase
	updateUC *categoria.UpdateCategoriaServicoUseCase
	deleteUC *categoria.DeleteCategoriaServicoUseCase
	logger   *zap.Logger
}

// NewCategoriaServicoHandler cria um novo handler de categorias de serviço
func NewCategoriaServicoHandler(
	createUC *categoria.CreateCategoriaServicoUseCase,
	getUC *categoria.GetCategoriaServicoUseCase,
	listUC *categoria.ListCategoriasServicosUseCase,
	updateUC *categoria.UpdateCategoriaServicoUseCase,
	deleteUC *categoria.DeleteCategoriaServicoUseCase,
	logger *zap.Logger,
) *CategoriaServicoHandler {
	return &CategoriaServicoHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
		logger:   logger,
	}
}

// Create godoc
// @Summary Criar categoria de serviço
// @Description Cria uma nova categoria de serviço
// @Tags Categorias de Serviço
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoriaServicoRequest true "Dados da categoria"
// @Success 201 {object} dto.CategoriaServicoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Nome duplicado"
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categorias-servicos [post]
// @Security BearerAuth
func (h *CategoriaServicoHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id do contexto (middleware de autenticação)
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	unitID := middleware.GetUnitID(c)
	if unitID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: domain.ErrUnitIDRequired.Error(),
		})
	}

	// Bind e validação
	var req dto.CreateCategoriaServicoRequest
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

	// Executar use case
	result, err := h.createUC.Execute(ctx, tenantID, unitID, req)
	if err != nil {
		h.logger.Error("Erro ao criar categoria de serviço", zap.Error(err))

		switch err {
		case domain.ErrCategoriaNomeDuplicate:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "duplicate_name",
				Message: err.Error(),
			})
		case domain.ErrCategoriaNomeRequired,
			domain.ErrCategoriaNomeTooShort,
			domain.ErrCategoriaNomeTooLong,
			domain.ErrCategoriaCorInvalida,
			domain.ErrCategoriaDescricaoTooLong:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Erro ao criar categoria de serviço",
			})
		}
	}

	return c.JSON(http.StatusCreated, result)
}

// GetByID godoc
// @Summary Buscar categoria por ID
// @Description Busca uma categoria de serviço por ID
// @Tags Categorias de Serviço
// @Produce json
// @Param id path string true "ID da categoria"
// @Success 200 {object} dto.CategoriaServicoResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categorias-servicos/{id} [get]
// @Security BearerAuth
func (h *CategoriaServicoHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	unitID := middleware.GetUnitID(c)
	if unitID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: domain.ErrUnitIDRequired.Error(),
		})
	}

	categoriaID := c.Param("id")
	if categoriaID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID da categoria é obrigatório",
		})
	}

	result, err := h.getUC.Execute(ctx, tenantID, unitID, categoriaID)
	if err != nil {
		h.logger.Error("Erro ao buscar categoria de serviço", zap.Error(err))

		if err == domain.ErrCategoriaNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar categoria de serviço",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// List godoc
// @Summary Listar categorias de serviço
// @Description Lista todas as categorias de serviço com filtros opcionais
// @Tags Categorias de Serviço
// @Produce json
// @Param apenas_ativas query boolean false "Apenas categorias ativas"
// @Param order_by query string false "Ordenação (nome, criado_em)"
// @Success 200 {object} dto.ListCategoriasServicosResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categorias-servicos [get]
// @Security BearerAuth
func (h *CategoriaServicoHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	unitID := middleware.GetUnitID(c)
	if unitID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: domain.ErrUnitIDRequired.Error(),
		})
	}

	// Bind query parameters
	var req dto.ListCategoriasServicosRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind de query params", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	result, err := h.listUC.Execute(ctx, tenantID, unitID, req)
	if err != nil {
		h.logger.Error("Erro ao listar categorias de serviço", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar categorias de serviço",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// Update godoc
// @Summary Atualizar categoria de serviço
// @Description Atualiza dados de uma categoria de serviço existente
// @Tags Categorias de Serviço
// @Accept json
// @Produce json
// @Param id path string true "ID da categoria"
// @Param request body dto.UpdateCategoriaServicoRequest true "Dados a atualizar"
// @Success 200 {object} dto.CategoriaServicoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Nome duplicado"
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categorias-servicos/{id} [put]
// @Security BearerAuth
func (h *CategoriaServicoHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	unitID := middleware.GetUnitID(c)
	if unitID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: domain.ErrUnitIDRequired.Error(),
		})
	}

	categoriaID := c.Param("id")
	if categoriaID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID da categoria é obrigatório",
		})
	}

	var req dto.UpdateCategoriaServicoRequest
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

	result, err := h.updateUC.Execute(ctx, tenantID, unitID, categoriaID, req)
	if err != nil {
		h.logger.Error("Erro ao atualizar categoria de serviço", zap.Error(err))

		switch err {
		case domain.ErrCategoriaNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		case domain.ErrCategoriaNomeDuplicate:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "duplicate_name",
				Message: err.Error(),
			})
		case domain.ErrCategoriaNomeRequired,
			domain.ErrCategoriaNomeTooShort,
			domain.ErrCategoriaNomeTooLong,
			domain.ErrCategoriaCorInvalida,
			domain.ErrCategoriaDescricaoTooLong:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Erro ao atualizar categoria de serviço",
			})
		}
	}

	return c.JSON(http.StatusOK, result)
}

// Delete godoc
// @Summary Deletar categoria de serviço
// @Description Deleta uma categoria de serviço se não houver serviços vinculados
// @Tags Categorias de Serviço
// @Param id path string true "ID da categoria"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Categoria possui serviços vinculados"
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/categorias-servicos/{id} [delete]
// @Security BearerAuth
func (h *CategoriaServicoHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	unitID := middleware.GetUnitID(c)
	if unitID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: domain.ErrUnitIDRequired.Error(),
		})
	}

	categoriaID := c.Param("id")
	if categoriaID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID da categoria é obrigatório",
		})
	}

	err := h.deleteUC.Execute(ctx, tenantID, unitID, categoriaID)
	if err != nil {
		h.logger.Error("Erro ao deletar categoria de serviço", zap.Error(err))

		switch err {
		case domain.ErrCategoriaNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		case domain.ErrCategoriaHasServices:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "has_services",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Erro ao deletar categoria de serviço",
			})
		}
	}

	return c.NoContent(http.StatusNoContent)
}
