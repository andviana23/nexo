package handler

import (
	"net/http"
	"strings"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/stock"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// StockHandler gerencia endpoints de estoque
type StockHandler struct {
	produtoRepo        port.ProdutoRepository
	criarProdutoUC     *stock.CriarProdutoUseCase
	registrarEntradaUC *stock.RegistrarEntradaUseCase
	registrarSaidaUC   *stock.RegistrarSaidaUseCase
	ajustarEstoqueUC   *stock.AjustarEstoqueUseCase
	listarAlertasUC    *stock.ListarAlertasEstoqueBaixoUseCase
}

// NewStockHandler cria nova instancia do handler
func NewStockHandler(
	produtoRepo port.ProdutoRepository,
	criarProdutoUC *stock.CriarProdutoUseCase,
	registrarEntradaUC *stock.RegistrarEntradaUseCase,
	registrarSaidaUC *stock.RegistrarSaidaUseCase,
	ajustarEstoqueUC *stock.AjustarEstoqueUseCase,
	listarAlertasUC *stock.ListarAlertasEstoqueBaixoUseCase,
) *StockHandler {
	return &StockHandler{
		produtoRepo:        produtoRepo,
		criarProdutoUC:     criarProdutoUC,
		registrarEntradaUC: registrarEntradaUC,
		registrarSaidaUC:   registrarSaidaUC,
		ajustarEstoqueUC:   ajustarEstoqueUC,
		listarAlertasUC:    listarAlertasUC,
	}
}

// CreateProduto godoc
// @Summary Criar novo produto
// @Description Cria um novo produto no estoque
// @Tags Estoque
// @Accept json
// @Produce json
// @Param request body dto.CreateProdutoRequest true "Dados do produto"
// @Success 201 {object} dto.ProdutoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "SKU já existe"
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/stock/products [post]
func (h *StockHandler) CreateProduto(c echo.Context) error {
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

	// Bind request
	var req dto.CreateProdutoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	// Validar request
	if req.Nome == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Nome é obrigatório",
		})
	}
	if req.CategoriaProdutoID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Categoria do produto é obrigatória",
		})
	}

	// Converter DTO para Input do use case via mapper
	input, err := mapper.FromCreateProdutoRequest(tenantID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "conversion_error",
			Message: err.Error(),
		})
	}

	// Executar use case
	output, err := h.criarProdutoUC.Execute(ctx, input)
	if err != nil {
		// Verificar se é erro de conflito (produto duplicado)
		if strings.Contains(err.Error(), "já está em uso") {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	// Converter output para response
	resp := mapper.CriarProdutoOutputToResponse(output)

	return c.JSON(http.StatusCreated, resp)
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

// ListProdutos godoc
// @Summary Listar produtos do estoque
// @Description Lista todos os produtos do estoque do tenant
// @Tags Estoque
// @Produce json
// @Success 200 {object} dto.ListProdutosResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/stock/items [get]
func (h *StockHandler) ListProdutos(c echo.Context) error {
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
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_tenant_id",
			Message: "Tenant ID inválido",
		})
	}

	// Buscar produtos
	produtos, err := h.produtoRepo.ListAll(ctx, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}

	// Converter para response e calcular contagens
	responses := make([]dto.ProdutoResponse, len(produtos))
	lowStockCount := 0
	outOfStockCount := 0

	for i, p := range produtos {
		responses[i] = mapper.ToProdutoResponse(p)

		// Calcular contagens de estoque
		qtdAtual := p.QuantidadeAtual.InexactFloat64()
		qtdMinima := p.QuantidadeMinima.InexactFloat64()

		if qtdAtual == 0 {
			outOfStockCount++
		} else if qtdAtual <= qtdMinima {
			lowStockCount++
		}
	}

	return c.JSON(http.StatusOK, dto.ListProdutosResponse{
		Data:            responses,
		Total:           len(responses),
		LowStockCount:   lowStockCount,
		OutOfStockCount: outOfStockCount,
	})
}

// GetProduto godoc
// @Summary Buscar produto por ID
// @Description Busca um produto específico do estoque
// @Tags Estoque
// @Produce json
// @Param id path string true "ID do produto"
// @Success 200 {object} dto.ProdutoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/stock/items/{id} [get]
func (h *StockHandler) GetProduto(c echo.Context) error {
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
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_tenant_id",
			Message: "Tenant ID inválido",
		})
	}

	// Extrair ID do produto
	produtoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "ID do produto inválido",
		})
	}

	// Buscar produto
	produto, err := h.produtoRepo.FindByID(ctx, tenantID, produtoID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}
	if produto == nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Produto não encontrado",
		})
	}

	return c.JSON(http.StatusOK, mapper.ToProdutoResponse(produto))
}
