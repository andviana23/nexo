package handler

import (
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/servico"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ServicoHandler agrupa os handlers de serviços
type ServicoHandler struct {
	createUC       *servico.CreateServicoUseCase
	getUC          *servico.GetServicoUseCase
	listUC         *servico.ListServicosUseCase
	updateUC       *servico.UpdateServicoUseCase
	deleteUC       *servico.DeleteServicoUseCase
	toggleStatusUC *servico.ToggleServicoStatusUseCase
	getStatsUC     *servico.GetServicoStatsUseCase
	logger         *zap.Logger
}

// NewServicoHandler cria um novo handler de serviços
func NewServicoHandler(
	createUC *servico.CreateServicoUseCase,
	getUC *servico.GetServicoUseCase,
	listUC *servico.ListServicosUseCase,
	updateUC *servico.UpdateServicoUseCase,
	deleteUC *servico.DeleteServicoUseCase,
	toggleStatusUC *servico.ToggleServicoStatusUseCase,
	getStatsUC *servico.GetServicoStatsUseCase,
	logger *zap.Logger,
) *ServicoHandler {
	return &ServicoHandler{
		createUC:       createUC,
		getUC:          getUC,
		listUC:         listUC,
		updateUC:       updateUC,
		deleteUC:       deleteUC,
		toggleStatusUC: toggleStatusUC,
		getStatsUC:     getStatsUC,
		logger:         logger,
	}
}

// Create godoc
// @Summary Criar serviço
// @Description Cria um novo serviço
// @Tags Serviços
// @Accept json
// @Produce json
// @Param request body dto.CreateServicoRequest true "Dados do serviço"
// @Success 201 {object} dto.ServicoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Nome duplicado"
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/servicos [post]
// @Security BearerAuth
func (h *ServicoHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id do contexto
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
			Error:   "unit_required",
			Message: domain.ErrUnitIDRequired.Error(),
		})
	}

	// Bind e validação
	var req dto.CreateServicoRequest
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

	// Segurança: unit_id vem do contexto (JWT/Header), nunca do payload
	req.UnitID = &unitID

	// Executar use case
	result, err := h.createUC.Execute(ctx, tenantID, req)
	if err != nil {
		h.logger.Error("Erro ao criar serviço", zap.Error(err))

		switch err {
		case domain.ErrServicoNomeDuplicate:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "duplicate_name",
				Message: err.Error(),
			})
		case domain.ErrCategoriaNotFound:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "categoria_not_found",
				Message: err.Error(),
			})
		case domain.ErrServicoNomeRequired,
			domain.ErrServicoNomeTooShort,
			domain.ErrServicoNomeTooLong,
			domain.ErrServicoPrecoInvalido,
			domain.ErrServicoDuracaoInvalida,
			domain.ErrServicoComissaoInvalida,
			domain.ErrServicoCorInvalida:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Erro ao criar serviço",
			})
		}
	}

	return c.JSON(http.StatusCreated, result)
}

// GetByID godoc
// @Summary Buscar serviço por ID
// @Description Busca um serviço por ID
// @Tags Serviços
// @Produce json
// @Param id path string true "ID do serviço"
// @Success 200 {object} dto.ServicoResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/servicos/{id} [get]
// @Security BearerAuth
func (h *ServicoHandler) GetByID(c echo.Context) error {
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
			Error:   "unit_required",
			Message: domain.ErrUnitIDRequired.Error(),
		})
	}

	servicoID := c.Param("id")
	if servicoID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do serviço é obrigatório",
		})
	}

	result, err := h.getUC.Execute(ctx, tenantID, unitID, servicoID)
	if err != nil {
		h.logger.Error("Erro ao buscar serviço", zap.Error(err))

		if err == domain.ErrServicoNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar serviço",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// List godoc
// @Summary Listar serviços
// @Description Lista serviços com filtros opcionais
// @Tags Serviços
// @Produce json
// @Param apenas_ativos query bool false "Apenas serviços ativos"
// @Param categoria_id query string false "Filtrar por categoria"
// @Param profissional_id query string false "Filtrar por profissional"
// @Param search query string false "Buscar por nome/descrição/tags"
// @Param order_by query string false "Ordenar por (nome, preco, duracao, criado_em)"
// @Success 200 {object} dto.ListServicosResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/servicos [get]
// @Security BearerAuth
func (h *ServicoHandler) List(c echo.Context) error {
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
			Error:   "unit_required",
			Message: domain.ErrUnitIDRequired.Error(),
		})
	}

	// Query params
	var req dto.ListServicosRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind dos query params", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// Segurança: forçar filtro por unit_id do contexto
	req.UnitID = unitID

	result, err := h.listUC.Execute(ctx, tenantID, req)
	if err != nil {
		h.logger.Error("Erro ao listar serviços", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar serviços",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// Update godoc
// @Summary Atualizar serviço
// @Description Atualiza um serviço existente
// @Tags Serviços
// @Accept json
// @Produce json
// @Param id path string true "ID do serviço"
// @Param request body dto.UpdateServicoRequest true "Dados do serviço"
// @Success 200 {object} dto.ServicoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Nome duplicado"
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/servicos/{id} [put]
// @Security BearerAuth
func (h *ServicoHandler) Update(c echo.Context) error {
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
			Error:   "unit_required",
			Message: domain.ErrUnitIDRequired.Error(),
		})
	}

	servicoID := c.Param("id")
	if servicoID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do serviço é obrigatório",
		})
	}

	// Bind e validação
	var req dto.UpdateServicoRequest
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

	result, err := h.updateUC.Execute(ctx, tenantID, unitID, servicoID, req)
	if err != nil {
		h.logger.Error("Erro ao atualizar serviço", zap.Error(err))

		switch err {
		case domain.ErrServicoNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		case domain.ErrServicoNomeDuplicate:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "duplicate_name",
				Message: err.Error(),
			})
		case domain.ErrCategoriaNotFound:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "categoria_not_found",
				Message: err.Error(),
			})
		case domain.ErrServicoNomeRequired,
			domain.ErrServicoNomeTooShort,
			domain.ErrServicoNomeTooLong,
			domain.ErrServicoPrecoInvalido,
			domain.ErrServicoDuracaoInvalida,
			domain.ErrServicoComissaoInvalida,
			domain.ErrServicoCorInvalida:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Erro ao atualizar serviço",
			})
		}
	}

	// Buscar serviço completo com categoria para resposta consistente
	servico, err := h.getUC.Execute(ctx, tenantID, unitID, servicoID)
	if err != nil {
		// Se update funcionou mas busca falhou, retorna o result do update
		return c.JSON(http.StatusOK, result)
	}

	return c.JSON(http.StatusOK, servico)
}

// Delete godoc
// @Summary Deletar serviço
// @Description Deleta um serviço
// @Tags Serviços
// @Produce json
// @Param id path string true "ID do serviço"
// @Success 204 "No Content"
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/servicos/{id} [delete]
// @Security BearerAuth
func (h *ServicoHandler) Delete(c echo.Context) error {
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
			Error:   "unit_required",
			Message: domain.ErrUnitIDRequired.Error(),
		})
	}

	servicoID := c.Param("id")
	if servicoID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do serviço é obrigatório",
		})
	}

	err := h.deleteUC.Execute(ctx, tenantID, unitID, servicoID)
	if err != nil {
		h.logger.Error("Erro ao deletar serviço", zap.Error(err))

		if err == domain.ErrServicoNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao deletar serviço",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ToggleStatus godoc
// @Summary Ativar/Desativar serviço
// @Description Altera o status ativo/inativo de um serviço
// @Tags Serviços
// @Accept json
// @Produce json
// @Param id path string true "ID do serviço"
// @Success 200 {object} dto.ServicoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/servicos/{id}/toggle-status [patch]
// @Security BearerAuth
func (h *ServicoHandler) ToggleStatus(c echo.Context) error {
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
			Error:   "unit_required",
			Message: domain.ErrUnitIDRequired.Error(),
		})
	}

	servicoID := c.Param("id")
	if servicoID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do serviço é obrigatório",
		})
	}

	result, err := h.toggleStatusUC.Execute(ctx, tenantID, unitID, servicoID)
	if err != nil {
		h.logger.Error("Erro ao alterar status do serviço", zap.Error(err))

		if err == domain.ErrServicoNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao alterar status do serviço",
		})
	}

	// Buscar serviço completo com categoria para resposta consistente
	servico, err := h.getUC.Execute(ctx, tenantID, unitID, servicoID)
	if err != nil {
		// Se toggle funcionou mas busca falhou, retorna o result do toggle
		return c.JSON(http.StatusOK, result)
	}

	return c.JSON(http.StatusOK, servico)
}

// GetStats godoc
// @Summary Estatísticas de serviços
// @Description Retorna estatísticas agregadas dos serviços
// @Tags Serviços
// @Produce json
// @Success 200 {object} dto.ServicoStatsResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/servicos/stats [get]
// @Security BearerAuth
func (h *ServicoHandler) GetStats(c echo.Context) error {
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
			Error:   "unit_required",
			Message: domain.ErrUnitIDRequired.Error(),
		})
	}

	result, err := h.getStatsUC.Execute(ctx, tenantID, unitID)
	if err != nil {
		h.logger.Error("Erro ao buscar estatísticas de serviços", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar estatísticas",
		})
	}

	return c.JSON(http.StatusOK, result)
}
