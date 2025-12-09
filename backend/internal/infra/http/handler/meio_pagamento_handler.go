package handler

import (
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/meiopagamento"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// MeioPagamentoHandler agrupa os handlers de meios de pagamento
type MeioPagamentoHandler struct {
	createUC *meiopagamento.CreateMeioPagamentoUseCase
	getUC    *meiopagamento.GetMeioPagamentoUseCase
	listUC   *meiopagamento.ListMeiosPagamentoUseCase
	updateUC *meiopagamento.UpdateMeioPagamentoUseCase
	toggleUC *meiopagamento.ToggleMeioPagamentoUseCase
	deleteUC *meiopagamento.DeleteMeioPagamentoUseCase
	logger   *zap.Logger
}

// NewMeioPagamentoHandler cria um novo handler de meios de pagamento
func NewMeioPagamentoHandler(
	createUC *meiopagamento.CreateMeioPagamentoUseCase,
	getUC *meiopagamento.GetMeioPagamentoUseCase,
	listUC *meiopagamento.ListMeiosPagamentoUseCase,
	updateUC *meiopagamento.UpdateMeioPagamentoUseCase,
	toggleUC *meiopagamento.ToggleMeioPagamentoUseCase,
	deleteUC *meiopagamento.DeleteMeioPagamentoUseCase,
	logger *zap.Logger,
) *MeioPagamentoHandler {
	return &MeioPagamentoHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		toggleUC: toggleUC,
		deleteUC: deleteUC,
		logger:   logger,
	}
}

// Create godoc
// @Summary Criar meio de pagamento
// @Description Cria um novo meio de pagamento
// @Tags Meios de Pagamento
// @Accept json
// @Produce json
// @Param request body dto.CreateMeioPagamentoRequest true "Dados do meio de pagamento"
// @Success 201 {object} dto.MeioPagamentoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/payment-methods [post]
// @Security BearerAuth
func (h *MeioPagamentoHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.CreateMeioPagamentoRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	// DEBUG: Log do payload recebido
	h.logger.Info("Payload recebido",
		zap.String("nome", req.Nome),
		zap.String("tipo", req.Tipo),
		zap.String("bandeira", req.Bandeira),
		zap.String("taxa", req.Taxa),
		zap.String("taxa_fixa", req.TaxaFixa),
		zap.Int("d_mais", req.DMais),
	)

	if err := c.Validate(&req); err != nil {
		h.logger.Error("Erro de validação", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	result, err := h.createUC.Execute(ctx, tenantID, req)
	if err != nil {
		h.logger.Error("Erro ao criar meio de pagamento",
			zap.Error(err),
			zap.String("tipo_erro", err.Error()),
		)

		switch err {
		case entity.ErrMeioPagamentoNomeVazio,
			entity.ErrMeioPagamentoNomeMuitoLongo,
			entity.ErrMeioPagamentoTipoInvalido,
			entity.ErrMeioPagamentoTaxaInvalida,
			entity.ErrMeioPagamentoTaxaFixaInvalida,
			entity.ErrMeioPagamentoDMaisInvalido:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Erro ao criar meio de pagamento",
			})
		}
	}

	return c.JSON(http.StatusCreated, result)
}

// Get godoc
// @Summary Buscar meio de pagamento por ID
// @Description Retorna um meio de pagamento pelo ID
// @Tags Meios de Pagamento
// @Produce json
// @Param id path string true "ID do meio de pagamento"
// @Success 200 {object} dto.MeioPagamentoResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/payment-methods/{id} [get]
// @Security BearerAuth
func (h *MeioPagamentoHandler) Get(c echo.Context) error {
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
			Message: "ID é obrigatório",
		})
	}

	result, err := h.getUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar meio de pagamento", zap.Error(err))
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Meio de pagamento não encontrado",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// List godoc
// @Summary Listar meios de pagamento
// @Description Lista todos os meios de pagamento do tenant
// @Tags Meios de Pagamento
// @Produce json
// @Param tipo query string false "Filtrar por tipo (DINHEIRO, PIX, CREDITO, DEBITO, TRANSFERENCIA)"
// @Param apenas_ativos query bool false "Apenas ativos"
// @Param search query string false "Buscar por nome"
// @Success 200 {object} dto.ListMeiosPagamentoResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/payment-methods [get]
// @Security BearerAuth
func (h *MeioPagamentoHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var filter dto.MeioPagamentoFilter
	if err := c.Bind(&filter); err != nil {
		h.logger.Warn("Erro ao fazer bind dos filtros", zap.Error(err))
	}

	result, err := h.listUC.Execute(ctx, tenantID, filter)
	if err != nil {
		h.logger.Error("Erro ao listar meios de pagamento", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar meios de pagamento",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// Update godoc
// @Summary Atualizar meio de pagamento
// @Description Atualiza um meio de pagamento existente
// @Tags Meios de Pagamento
// @Accept json
// @Produce json
// @Param id path string true "ID do meio de pagamento"
// @Param request body dto.UpdateMeioPagamentoRequest true "Dados a atualizar"
// @Success 200 {object} dto.MeioPagamentoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/payment-methods/{id} [patch]
// @Security BearerAuth
func (h *MeioPagamentoHandler) Update(c echo.Context) error {
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
			Message: "ID é obrigatório",
		})
	}

	var req dto.UpdateMeioPagamentoRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	result, err := h.updateUC.Execute(ctx, tenantID, id, req)
	if err != nil {
		h.logger.Error("Erro ao atualizar meio de pagamento", zap.Error(err))

		switch err {
		case entity.ErrMeioPagamentoTipoInvalido,
			entity.ErrMeioPagamentoTaxaInvalida,
			entity.ErrMeioPagamentoTaxaFixaInvalida,
			entity.ErrMeioPagamentoDMaisInvalido:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Meio de pagamento não encontrado",
			})
		}
	}

	return c.JSON(http.StatusOK, result)
}

// Toggle godoc
// @Summary Alternar status do meio de pagamento
// @Description Alterna entre ativo e inativo
// @Tags Meios de Pagamento
// @Produce json
// @Param id path string true "ID do meio de pagamento"
// @Success 200 {object} dto.MeioPagamentoResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/payment-methods/{id}/toggle [patch]
// @Security BearerAuth
func (h *MeioPagamentoHandler) Toggle(c echo.Context) error {
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
			Message: "ID é obrigatório",
		})
	}

	result, err := h.toggleUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao alternar status", zap.Error(err))
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Meio de pagamento não encontrado",
		})
	}

	return c.JSON(http.StatusOK, result)
}

// Delete godoc
// @Summary Excluir meio de pagamento
// @Description Remove um meio de pagamento
// @Tags Meios de Pagamento
// @Param id path string true "ID do meio de pagamento"
// @Success 204 "No Content"
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/payment-methods/{id} [delete]
// @Security BearerAuth
func (h *MeioPagamentoHandler) Delete(c echo.Context) error {
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
			Message: "ID é obrigatório",
		})
	}

	if err := h.deleteUC.Execute(ctx, tenantID, id); err != nil {
		h.logger.Error("Erro ao excluir meio de pagamento", zap.Error(err))
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Meio de pagamento não encontrado",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
