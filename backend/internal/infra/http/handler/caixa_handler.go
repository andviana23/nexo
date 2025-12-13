// Package handler contém os HTTP handlers da camada de infraestrutura.
package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/caixa"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	mw "github.com/andviana23/barber-analytics-backend/internal/infra/http/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// CaixaHandler agrupa os handlers do módulo Caixa Diário
type CaixaHandler struct {
	abrirCaixaUC     *caixa.AbrirCaixaUseCase
	sangriaUC        *caixa.SangriaUseCase
	reforcoUC        *caixa.ReforcoUseCase
	fecharCaixaUC    *caixa.FecharCaixaUseCase
	getCaixaAbertoUC *caixa.GetCaixaAbertoUseCase
	getCaixaByIDUC   *caixa.GetCaixaByIDUseCase
	listHistoricoUC  *caixa.ListHistoricoUseCase
	getTotaisUC      *caixa.GetTotaisCaixaUseCase
	logger           *zap.Logger
}

// NewCaixaHandler cria um novo handler de caixa
func NewCaixaHandler(
	abrirCaixaUC *caixa.AbrirCaixaUseCase,
	sangriaUC *caixa.SangriaUseCase,
	reforcoUC *caixa.ReforcoUseCase,
	fecharCaixaUC *caixa.FecharCaixaUseCase,
	getCaixaAbertoUC *caixa.GetCaixaAbertoUseCase,
	getCaixaByIDUC *caixa.GetCaixaByIDUseCase,
	listHistoricoUC *caixa.ListHistoricoUseCase,
	getTotaisUC *caixa.GetTotaisCaixaUseCase,
	logger *zap.Logger,
) *CaixaHandler {
	return &CaixaHandler{
		abrirCaixaUC:     abrirCaixaUC,
		sangriaUC:        sangriaUC,
		reforcoUC:        reforcoUC,
		fecharCaixaUC:    fecharCaixaUC,
		getCaixaAbertoUC: getCaixaAbertoUC,
		getCaixaByIDUC:   getCaixaByIDUC,
		listHistoricoUC:  listHistoricoUC,
		getTotaisUC:      getTotaisUC,
		logger:           logger,
	}
}

// RegisterRoutes registra as rotas do módulo Caixa
func (h *CaixaHandler) RegisterRoutes(g *echo.Group) {
	caixaGroup := g.Group("/caixa")

	// Status e consultas - Qualquer role autenticada pode visualizar
	caixaGroup.GET("/status", h.GetStatus, mw.RequireAnyRole(h.logger))
	caixaGroup.GET("/aberto", h.GetCaixaAberto, mw.RequireAnyRole(h.logger))
	caixaGroup.GET("/historico", h.ListHistorico, mw.RequireAnyRole(h.logger))
	caixaGroup.GET("/totais", h.GetTotais, mw.RequireAnyRole(h.logger))
	caixaGroup.GET("/:id", h.GetCaixaByID, mw.RequireAnyRole(h.logger))

	// Operações críticas - Apenas OWNER/MANAGER podem executar
	caixaGroup.POST("/abrir", h.AbrirCaixa, mw.RequireOwnerOrManager(h.logger))
	caixaGroup.POST("/sangria", h.Sangria, mw.RequireOwnerOrManager(h.logger))
	caixaGroup.POST("/reforco", h.Reforco, mw.RequireOwnerOrManager(h.logger))
	caixaGroup.POST("/fechar", h.FecharCaixa, mw.RequireOwnerOrManager(h.logger))
}

// AbrirCaixa abre um novo caixa diário
// @Summary Abrir caixa
// @Description Abre um novo caixa diário com saldo inicial
// @Tags Caixa
// @Accept json
// @Produce json
// @Param request body dto.AbrirCaixaRequest true "Dados de abertura"
// @Success 201 {object} dto.CaixaDiarioResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string "Já existe caixa aberto"
// @Router /api/v1/caixa/abrir [post]
func (h *CaixaHandler) AbrirCaixa(c echo.Context) error {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "usuário não identificado"})
	}

	var req dto.AbrirCaixaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dados inválidos"})
	}

	saldoInicial, err := decimal.NewFromString(req.SaldoInicial)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "saldo_inicial inválido"})
	}

	result, err := h.abrirCaixaUC.Execute(c.Request().Context(), caixa.AbrirCaixaInput{
		TenantID:     tenantID,
		UsuarioID:    userID,
		SaldoInicial: saldoInicial,
	})
	if err != nil {
		h.logger.Error("Erro ao abrir caixa", zap.Error(err))
		return handleCaixaError(c, err)
	}

	return c.JSON(http.StatusCreated, mapper.ToCaixaDiarioResponse(result))
}

// Sangria registra uma sangria no caixa
// @Summary Registrar sangria
// @Description Registra uma retirada de dinheiro do caixa
// @Tags Caixa
// @Accept json
// @Produce json
// @Param request body dto.SangriaRequest true "Dados da sangria"
// @Success 201 {object} dto.OperacaoCaixaResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string "Caixa não aberto"
// @Router /api/v1/caixa/sangria [post]
func (h *CaixaHandler) Sangria(c echo.Context) error {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "usuário não identificado"})
	}

	var req dto.SangriaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dados inválidos"})
	}

	valor, err := decimal.NewFromString(req.Valor)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "valor inválido"})
	}

	result, err := h.sangriaUC.Execute(c.Request().Context(), caixa.SangriaInput{
		TenantID:  tenantID,
		UsuarioID: userID,
		Valor:     valor,
		Destino:   req.Destino,
		Descricao: req.Descricao,
	})
	if err != nil {
		h.logger.Error("Erro ao registrar sangria", zap.Error(err))
		return handleCaixaError(c, err)
	}

	return c.JSON(http.StatusCreated, mapper.ToOperacaoCaixaResponse(result))
}

// Reforco registra um reforço no caixa
// @Summary Registrar reforço
// @Description Registra uma adição de dinheiro ao caixa
// @Tags Caixa
// @Accept json
// @Produce json
// @Param request body dto.ReforcoRequest true "Dados do reforço"
// @Success 201 {object} dto.OperacaoCaixaResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string "Caixa não aberto"
// @Router /api/v1/caixa/reforco [post]
func (h *CaixaHandler) Reforco(c echo.Context) error {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "usuário não identificado"})
	}

	var req dto.ReforcoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dados inválidos"})
	}

	valor, err := decimal.NewFromString(req.Valor)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "valor inválido"})
	}

	result, err := h.reforcoUC.Execute(c.Request().Context(), caixa.ReforcoInput{
		TenantID:  tenantID,
		UsuarioID: userID,
		Valor:     valor,
		Origem:    req.Origem,
		Descricao: req.Descricao,
	})
	if err != nil {
		h.logger.Error("Erro ao registrar reforço", zap.Error(err))
		return handleCaixaError(c, err)
	}

	return c.JSON(http.StatusCreated, mapper.ToOperacaoCaixaResponse(result))
}

// FecharCaixa fecha o caixa diário atual
// @Summary Fechar caixa
// @Description Fecha o caixa com conferência de saldo real
// @Tags Caixa
// @Accept json
// @Produce json
// @Param request body dto.FecharCaixaRequest true "Dados de fechamento"
// @Success 200 {object} dto.CaixaDiarioResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string "Caixa não aberto"
// @Router /api/v1/caixa/fechar [post]
func (h *CaixaHandler) FecharCaixa(c echo.Context) error {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "usuário não identificado"})
	}

	var req dto.FecharCaixaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dados inválidos"})
	}

	saldoReal, err := decimal.NewFromString(req.SaldoReal)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "saldo_real inválido"})
	}

	result, err := h.fecharCaixaUC.Execute(c.Request().Context(), caixa.FecharCaixaInput{
		TenantID:      tenantID,
		UsuarioID:     userID,
		SaldoReal:     saldoReal,
		Justificativa: req.Justificativa,
	})
	if err != nil {
		h.logger.Error("Erro ao fechar caixa", zap.Error(err))
		return handleCaixaError(c, err)
	}

	return c.JSON(http.StatusOK, mapper.ToCaixaDiarioResponse(result))
}

// GetStatus retorna o status do caixa (aberto/fechado)
// @Summary Status do caixa
// @Description Retorna se há caixa aberto e informações básicas
// @Tags Caixa
// @Produce json
// @Success 200 {object} dto.CaixaStatusResponse
// @Router /api/v1/caixa/status [get]
func (h *CaixaHandler) GetStatus(c echo.Context) error {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}

	caixaAberto, err := h.getCaixaAbertoUC.Execute(c.Request().Context(), tenantID)
	if err != nil {
		h.logger.Error("Erro ao buscar status do caixa", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro interno"})
	}

	return c.JSON(http.StatusOK, mapper.ToCaixaStatusResponse(caixaAberto, nil))
}

// GetCaixaAberto retorna o caixa aberto com operações
// @Summary Caixa aberto
// @Description Retorna o caixa aberto com todas as operações
// @Tags Caixa
// @Produce json
// @Success 200 {object} dto.CaixaDiarioResponse
// @Failure 404 {object} map[string]string "Nenhum caixa aberto"
// @Router /api/v1/caixa/aberto [get]
func (h *CaixaHandler) GetCaixaAberto(c echo.Context) error {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}

	result, err := h.getCaixaAbertoUC.Execute(c.Request().Context(), tenantID)
	if err != nil {
		h.logger.Error("Erro ao buscar caixa aberto", zap.Error(err))
		return handleCaixaError(c, err)
	}

	if result == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "nenhum caixa aberto"})
	}

	return c.JSON(http.StatusOK, mapper.ToCaixaDiarioResponse(result))
}

// GetCaixaByID retorna um caixa específico
// @Summary Buscar caixa por ID
// @Description Retorna um caixa específico com operações
// @Tags Caixa
// @Produce json
// @Param id path string true "ID do caixa"
// @Success 200 {object} dto.CaixaDiarioResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/caixa/{id} [get]
func (h *CaixaHandler) GetCaixaByID(c echo.Context) error {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}

	caixaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id inválido"})
	}

	result, err := h.getCaixaByIDUC.Execute(c.Request().Context(), caixaID, tenantID)
	if err != nil {
		h.logger.Error("Erro ao buscar caixa", zap.Error(err))
		return handleCaixaError(c, err)
	}

	return c.JSON(http.StatusOK, mapper.ToCaixaDiarioResponse(result))
}

// ListHistorico retorna o histórico de caixas fechados
// @Summary Histórico de caixas
// @Description Lista caixas fechados com paginação
// @Tags Caixa
// @Produce json
// @Param data_inicio query string false "Data início (YYYY-MM-DD)"
// @Param data_fim query string false "Data fim (YYYY-MM-DD)"
// @Param usuario_id query string false "ID do usuário"
// @Param page query int false "Página" default(1)
// @Param page_size query int false "Itens por página" default(20)
// @Success 200 {object} dto.ListCaixaHistoricoResponse
// @Router /api/v1/caixa/historico [get]
func (h *CaixaHandler) ListHistorico(c echo.Context) error {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}

	var req dto.ListCaixaHistoricoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "parâmetros inválidos"})
	}

	// Valores default para paginação
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Parse datas (YYYY-MM-DD)
	var dataInicioPtr, dataFimPtr *time.Time
	if req.DataInicio != nil {
		parsed, err := time.Parse("2006-01-02", *req.DataInicio)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "data_inicio inválida (use YYYY-MM-DD)"})
		}
		dataInicioPtr = &parsed
	}
	if req.DataFim != nil {
		parsed, err := time.Parse("2006-01-02", *req.DataFim)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "data_fim inválida (use YYYY-MM-DD)"})
		}
		dataFimPtr = &parsed
	}

	input := caixa.ListHistoricoInput{
		TenantID:   tenantID,
		DataInicio: dataInicioPtr,
		DataFim:    dataFimPtr,
		UsuarioID:  nil,
		Page:       page,
		PageSize:   pageSize,
	}
	if req.UsuarioID != nil {
		u, err := uuid.Parse(*req.UsuarioID)
		if err == nil {
			input.UsuarioID = &u
		}
	}

	result, err := h.listHistoricoUC.Execute(c.Request().Context(), input)
	if err != nil {
		h.logger.Error("Erro ao listar histórico", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro interno"})
	}

	return c.JSON(http.StatusOK, mapper.ToListCaixaHistoricoResponse(result.Caixas, result.Total, input.Page, input.PageSize))
}

// GetTotais retorna os totais do caixa aberto
// @Summary Totais do caixa
// @Description Retorna os totais por tipo de operação
// @Tags Caixa
// @Produce json
// @Success 200 {object} dto.TotaisCaixaResponse
// @Failure 404 {object} map[string]string "Caixa não aberto"
// @Router /api/v1/caixa/totais [get]
func (h *CaixaHandler) GetTotais(c echo.Context) error {
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
	}

	result, err := h.getTotaisUC.Execute(c.Request().Context(), tenantID)
	if err != nil {
		h.logger.Error("Erro ao buscar totais", zap.Error(err))
		return handleCaixaError(c, err)
	}

	return c.JSON(http.StatusOK, dto.TotaisCaixaResponse{
		TotalVendas:   result.TotalVendas.String(),
		TotalSangrias: result.TotalSangrias.String(),
		TotalReforcos: result.TotalReforcos.String(),
		TotalDespesas: result.TotalDespesas.String(),
		SaldoAtual:    result.SaldoAtual.String(),
	})
}

// ============================================================
// HELPERS
// ============================================================

// handleCaixaError trata erros específicos do módulo caixa
func handleCaixaError(c echo.Context, err error) error {
	// Preferir errors.Is para não depender de mensagens string (mais robusto)
	switch {
	case errors.Is(err, domain.ErrCaixaJaAberto):
		return c.JSON(http.StatusConflict, map[string]string{"error": "já existe um caixa aberto"})
	case errors.Is(err, domain.ErrCaixaNaoAberto):
		return c.JSON(http.StatusNotFound, map[string]string{"error": "nenhum caixa aberto"})
	case errors.Is(err, domain.ErrCaixaJaFechado):
		return c.JSON(http.StatusConflict, map[string]string{"error": "caixa já está fechado"})
	case errors.Is(err, domain.ErrCaixaJustificativaObrigatoria):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "justificativa obrigatória para divergência maior que R$ 5,00"})
	case errors.Is(err, domain.ErrCaixaNotFound):
		return c.JSON(http.StatusNotFound, map[string]string{"error": "caixa não encontrado"})
	case errors.Is(err, domain.ErrSangriaDestinoObrigatorio):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "destino é obrigatório para sangria"})
	case errors.Is(err, domain.ErrReforcoOrigemObrigatoria):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "origem é obrigatória para reforço"})
	}

	switch err.Error() {
	case "caixa já aberto para este tenant":
		return c.JSON(http.StatusConflict, map[string]string{"error": "já existe um caixa aberto"})
	case "nenhum caixa aberto":
		return c.JSON(http.StatusNotFound, map[string]string{"error": "nenhum caixa aberto"})
	case "caixa já fechado":
		return c.JSON(http.StatusConflict, map[string]string{"error": "caixa já está fechado"})
	case "justificativa obrigatória para divergência maior que R$ 5,00":
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "justificativa obrigatória para divergência maior que R$ 5,00"})
	case "destino é obrigatório para sangria":
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "destino é obrigatório para sangria"})
	case "origem é obrigatória para reforço":
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "origem é obrigatória para reforço"})
	case "caixa não encontrado":
		return c.JSON(http.StatusNotFound, map[string]string{"error": "caixa não encontrado"})
	default:
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro interno"})
	}
}
