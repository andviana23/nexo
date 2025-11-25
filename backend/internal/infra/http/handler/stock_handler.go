package handler

import (
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/stock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// StockHandler gerencia endpoints de estoque
type StockHandler struct {
	registrarEntradaUC *stock.RegistrarEntradaUseCase
	registrarSaidaUC   *stock.RegistrarSaidaUseCase
	ajustarEstoqueUC   *stock.AjustarEstoqueUseCase
	listarAlertasUC    *stock.ListarAlertasEstoqueBaixoUseCase
}

// NewStockHandler cria nova instancia do handler
func NewStockHandler(
	registrarEntradaUC *stock.RegistrarEntradaUseCase,
	registrarSaidaUC *stock.RegistrarSaidaUseCase,
	ajustarEstoqueUC *stock.AjustarEstoqueUseCase,
	listarAlertasUC *stock.ListarAlertasEstoqueBaixoUseCase,
) *StockHandler {
	return &StockHandler{
		registrarEntradaUC: registrarEntradaUC,
		registrarSaidaUC:   registrarSaidaUC,
		ajustarEstoqueUC:   ajustarEstoqueUC,
		listarAlertasUC:    listarAlertasUC,
	}
}

// RegistrarEntrada godoc
// @Summary Registrar entrada de estoque
// @Description Registra entrada de produtos no estoque (compra/recebimento)
// @Tags Estoque
// @Accept json
// @Produce json
// @Param request body dto.RegistrarEntradaRequest true "Dados da entrada"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/stock/entries [post]
func (h *StockHandler) RegistrarEntrada(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id do contexto (middleware JWT)
	tenantIDStr, ok := c.Get("tenant_id").(string)
	if !ok || tenantIDStr == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_tenant",
			Message: "Tenant ID inválido",
		})
	}

	// Extrair user_id do contexto
	userIDStr, ok := c.Get("user_id").(string)
	if !ok || userIDStr == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User ID não encontrado",
		})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "User ID inválido",
		})
	}

	// Bind request
	var req dto.RegistrarEntradaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	// Converter DTO para Input do use case via mapper
	input, err := mapper.FromRegistrarEntradaRequest(tenantID, userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "conversion_error",
			Message: err.Error(),
		})
	}

	// Executar use case
	output, err := h.registrarEntradaUC.Execute(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	// Converter output para response
	resp := mapper.ToRegistrarEntradaResponse(output)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Entrada registrada com sucesso",
		"data":    resp,
	})
}

// RegistrarSaida godoc
// @Summary Registrar saida de estoque
// @Description Registra saida de produto (venda, consumo interno, perda)
// @Tags Estoque
// @Accept json
// @Produce json
// @Param request body dto.RegistrarSaidaRequest true "Dados da saida"
// @Success 201 {object} dto.MovimentacaoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/stock/exit [post]
func (h *StockHandler) RegistrarSaida(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id do contexto
	tenantIDStr, ok := c.Get("tenant_id").(string)
	if !ok || tenantIDStr == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_tenant",
			Message: "Tenant ID inválido",
		})
	}

	// Extrair user_id do contexto
	userIDStr, ok := c.Get("user_id").(string)
	if !ok || userIDStr == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User ID não encontrado",
		})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "User ID inválido",
		})
	}

	// Bind request
	var req dto.RegistrarSaidaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	// Validar request
	_, err = mapper.FromRegistrarSaidaRequest(tenantID, userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "conversion_error",
			Message: err.Error(),
		})
	}

	// Executar use case
	output, err := h.registrarSaidaUC.Execute(ctx, tenantID, userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	// Converter output para response
	resp := mapper.ToMovimentacaoResponse(output)

	return c.JSON(http.StatusCreated, resp)
}

// AjustarEstoque godoc
// @Summary Ajustar estoque manualmente
// @Description Realiza ajuste manual de quantidade com auditoria
// @Tags Estoque
// @Accept json
// @Produce json
// @Param request body dto.AjustarEstoqueRequest true "Dados do ajuste"
// @Success 200 {object} dto.MovimentacaoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/stock/adjust [post]
func (h *StockHandler) AjustarEstoque(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id do contexto
	tenantIDStr, ok := c.Get("tenant_id").(string)
	if !ok || tenantIDStr == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_tenant",
			Message: "Tenant ID inválido",
		})
	}

	// Extrair user_id do contexto
	userIDStr, ok := c.Get("user_id").(string)
	if !ok || userIDStr == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User ID não encontrado",
		})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "User ID inválido",
		})
	}

	// Bind request
	var req dto.AjustarEstoqueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	// Validar request
	_, err = mapper.FromAjustarEstoqueRequest(tenantID, userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "conversion_error",
			Message: err.Error(),
		})
	}

	// Executar use case
	output, err := h.ajustarEstoqueUC.Execute(ctx, tenantID, userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	// Converter output para response
	resp := mapper.ToMovimentacaoResponse(output)

	return c.JSON(http.StatusOK, resp)
}

// ListarAlertas godoc
// @Summary Listar alertas de estoque baixo
// @Description Lista produtos com estoque abaixo do minimo
// @Tags Estoque
// @Accept json
// @Produce json
// @Success 200 {object} dto.ListAlertasEstoqueBaixoResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/stock/alerts [get]
func (h *StockHandler) ListarAlertas(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id do contexto
	tenantIDStr, ok := c.Get("tenant_id").(string)
	if !ok || tenantIDStr == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Executar use case
	alertas, err := h.listarAlertasUC.Execute(ctx, tenantIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	// Converter para response via mapper
	resp := mapper.ToListAlertasEstoqueBaixoResponse(alertas)

	return c.JSON(http.StatusOK, resp)
}
