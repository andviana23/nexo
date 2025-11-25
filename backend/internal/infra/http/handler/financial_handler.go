// Package handler contém os HTTP handlers da camada de infraestrutura.
// Implementa a interface web usando Echo framework seguindo Clean Architecture.
package handler

import (
	"net/http"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/financial"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// FinancialHandler agrupa os handlers financeiros
type FinancialHandler struct {
	// ContaPagar use cases
	createContaPagarUC *financial.CreateContaPagarUseCase
	getContaPagarUC    *financial.GetContaPagarUseCase
	listContasPagarUC  *financial.ListContasPagarUseCase
	updateContaPagarUC *financial.UpdateContaPagarUseCase
	deleteContaPagarUC *financial.DeleteContaPagarUseCase
	marcarPagamentoUC  *financial.MarcarPagamentoUseCase

	// ContaReceber use cases
	createContaReceberUC *financial.CreateContaReceberUseCase
	getContaReceberUC    *financial.GetContaReceberUseCase
	listContasReceberUC  *financial.ListContasReceberUseCase
	updateContaReceberUC *financial.UpdateContaReceberUseCase
	deleteContaReceberUC *financial.DeleteContaReceberUseCase
	marcarRecebimentoUC  *financial.MarcarRecebimentoUseCase

	// Compensação use cases
	createCompensacaoUC *financial.CreateCompensacaoUseCase
	getCompensacaoUC    *financial.GetCompensacaoUseCase
	listCompensacoesUC  *financial.ListCompensacoesUseCase
	deleteCompensacaoUC *financial.DeleteCompensacaoUseCase
	marcarCompensacaoUC *financial.MarcarCompensacaoUseCase

	// FluxoCaixa use cases
	generateFluxoDiarioUC *financial.GenerateFluxoDiarioUseCase
	getFluxoCaixaUC       *financial.GetFluxoCaixaUseCase
	listFluxoCaixaUC      *financial.ListFluxoCaixaUseCase

	// DRE use cases
	generateDREUC *financial.GenerateDREUseCase
	getDREUC      *financial.GetDREUseCase
	listDREUC     *financial.ListDREUseCase

	// Dashboard use case (TODO: implementar quando GetDashboardUseCase existir)
	// getDashboardUC *financial.GetDashboardUseCase

	logger *zap.Logger
}

// NewFinancialHandler cria um novo handler financeiro
func NewFinancialHandler(
	// ContaPagar
	createContaPagarUC *financial.CreateContaPagarUseCase,
	getContaPagarUC *financial.GetContaPagarUseCase,
	listContasPagarUC *financial.ListContasPagarUseCase,
	updateContaPagarUC *financial.UpdateContaPagarUseCase,
	deleteContaPagarUC *financial.DeleteContaPagarUseCase,
	marcarPagamentoUC *financial.MarcarPagamentoUseCase,
	// ContaReceber
	createContaReceberUC *financial.CreateContaReceberUseCase,
	getContaReceberUC *financial.GetContaReceberUseCase,
	listContasReceberUC *financial.ListContasReceberUseCase,
	updateContaReceberUC *financial.UpdateContaReceberUseCase,
	deleteContaReceberUC *financial.DeleteContaReceberUseCase,
	marcarRecebimentoUC *financial.MarcarRecebimentoUseCase,
	// Compensação
	createCompensacaoUC *financial.CreateCompensacaoUseCase,
	getCompensacaoUC *financial.GetCompensacaoUseCase,
	listCompensacoesUC *financial.ListCompensacoesUseCase,
	deleteCompensacaoUC *financial.DeleteCompensacaoUseCase,
	marcarCompensacaoUC *financial.MarcarCompensacaoUseCase,
	// FluxoCaixa
	generateFluxoDiarioUC *financial.GenerateFluxoDiarioUseCase,
	getFluxoCaixaUC *financial.GetFluxoCaixaUseCase,
	listFluxoCaixaUC *financial.ListFluxoCaixaUseCase,
	// DRE
	generateDREUC *financial.GenerateDREUseCase,
	getDREUC *financial.GetDREUseCase,
	listDREUC *financial.ListDREUseCase,
	// Dashboard (TODO: implementar quando GetDashboardUseCase existir)
	_ interface{}, // placeholder para getDashboardUC
	logger *zap.Logger,
) *FinancialHandler {
	return &FinancialHandler{
		// ContaPagar
		createContaPagarUC: createContaPagarUC,
		getContaPagarUC:    getContaPagarUC,
		listContasPagarUC:  listContasPagarUC,
		updateContaPagarUC: updateContaPagarUC,
		deleteContaPagarUC: deleteContaPagarUC,
		marcarPagamentoUC:  marcarPagamentoUC,
		// ContaReceber
		createContaReceberUC: createContaReceberUC,
		getContaReceberUC:    getContaReceberUC,
		listContasReceberUC:  listContasReceberUC,
		updateContaReceberUC: updateContaReceberUC,
		deleteContaReceberUC: deleteContaReceberUC,
		marcarRecebimentoUC:  marcarRecebimentoUC,
		// Compensação
		createCompensacaoUC: createCompensacaoUC,
		getCompensacaoUC:    getCompensacaoUC,
		listCompensacoesUC:  listCompensacoesUC,
		deleteCompensacaoUC: deleteCompensacaoUC,
		marcarCompensacaoUC: marcarCompensacaoUC,
		// FluxoCaixa
		generateFluxoDiarioUC: generateFluxoDiarioUC,
		getFluxoCaixaUC:       getFluxoCaixaUC,
		listFluxoCaixaUC:      listFluxoCaixaUC,
		// DRE
		generateDREUC: generateDREUC,
		getDREUC:      getDREUC,
		listDREUC:     listDREUC,
		// Dashboard (TODO)
		logger: logger,
	}
}

// CreateContaPagar godoc
// @Summary Criar conta a pagar
// @Description Cria uma nova conta a pagar (despesa)
// @Tags Financial
// @Accept json
// @Produce json
// @Param request body dto.CreateContaPagarRequest true "Dados da conta a pagar"
// @Success 201 {object} dto.ContaPagarResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/payables [post]
// @Security BearerAuth
func (h *FinancialHandler) CreateContaPagar(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id do contexto (definido por middleware JWT)
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		h.logger.Error("tenant_id não encontrado no contexto")
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Bind request
	var req dto.CreateContaPagarRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind da requisição", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	// Validar request
	if err := c.Validate(&req); err != nil {
		h.logger.Error("Erro de validação", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// Converter DTO para parâmetros do use case
	valor, tipo, dataVencimento, err := mapper.FromCreateContaPagarRequest(req)
	if err != nil {
		h.logger.Error("Erro ao converter request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "conversion_error",
			Message: err.Error(),
		})
	}

	// Executar use case
	conta, err := h.createContaPagarUC.Execute(ctx, financial.CreateContaPagarInput{
		TenantID:       tenantID,
		Descricao:      req.Descricao,
		CategoriaID:    req.CategoriaID,
		Fornecedor:     req.Fornecedor,
		Valor:          valor,
		Tipo:           tipo,
		DataVencimento: dataVencimento,
		Recorrente:     req.Recorrente,
		Periodicidade:  req.Periodicidade,
		PixCode:        req.PixCode,
		Observacoes:    req.Observacoes,
	})
	if err != nil {
		h.logger.Error("Erro ao criar conta a pagar", zap.Error(err), zap.String("tenant_id", tenantID))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao criar conta a pagar",
		})
	}

	// Converter entidade para response
	response := mapper.ToContaPagarResponse(conta)

	return c.JSON(http.StatusCreated, response)
}

// CreateContaReceber godoc
// @Summary Criar conta a receber
// @Description Cria uma nova conta a receber (receita)
// @Tags Financial
// @Accept json
// @Produce json
// @Param request body dto.CreateContaReceberRequest true "Dados da conta a receber"
// @Success 201 {object} dto.ContaReceberResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/receivables [post]
// @Security BearerAuth
func (h *FinancialHandler) CreateContaReceber(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.CreateContaReceberRequest
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

	valor, dataVencimento, err := mapper.FromCreateContaReceberRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "conversion_error",
			Message: err.Error(),
		})
	}

	input := financial.CreateContaReceberInput{
		TenantID:        tenantID,
		Descricao:       req.DescricaoOrigem,
		Origem:          req.Origem,
		AssinaturaID:    req.AssinaturaID,
		Valor:           valor,
		DataVencimento:  dataVencimento,
		MetodoPagamento: "",
		Observacoes:     req.Observacoes,
		Subtipo:         valueobject.SubtipoReceita(req.Origem),
	}

	conta, err := h.createContaReceberUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao criar conta a receber", zap.Error(err), zap.String("tenant_id", tenantID))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao criar conta a receber",
		})
	}

	response := mapper.ToContaReceberResponse(conta)
	return c.JSON(http.StatusCreated, response)
}

// MarcarPagamento godoc
// @Summary Marcar pagamento
// @Description Marca uma conta a pagar como paga
// @Tags Financial
// @Accept json
// @Produce json
// @Param id path string true "ID da conta a pagar"
// @Param request body dto.MarcarPagamentoRequest true "Dados do pagamento"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/payables/{id}/pay [post]
// @Security BearerAuth
func (h *FinancialHandler) MarcarPagamento(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	contaID := c.Param("id")
	if contaID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID da conta obrigatório",
		})
	}

	var req dto.MarcarPagamentoRequest
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

	dataPagamento, err := time.Parse("2006-01-02", req.DataPagamento)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "data_pagamento inválida",
		})
	}

	if _, err := h.marcarPagamentoUC.Execute(ctx, financial.MarcarPagamentoInput{
		TenantID:       tenantID,
		ContaID:        contaID,
		DataPagamento:  dataPagamento,
		ComprovanteURL: req.ComprovanteURL,
	}); err != nil {
		h.logger.Error("Erro ao marcar pagamento", zap.Error(err), zap.String("conta_id", contaID))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao marcar pagamento",
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Pagamento marcado com sucesso",
	})
}

// MarcarRecebimento godoc
// @Summary Marcar recebimento
// @Description Marca uma conta a receber como recebida (total ou parcial)
// @Tags Financial
// @Accept json
// @Produce json
// @Param id path string true "ID da conta a receber"
// @Param request body dto.MarcarRecebimentoRequest true "Dados do recebimento"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/receivables/{id}/receive [post]
// @Security BearerAuth
func (h *FinancialHandler) MarcarRecebimento(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	contaID := c.Param("id")
	if contaID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID da conta obrigatório",
		})
	}

	var req dto.MarcarRecebimentoRequest
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

	if _, err := decimal.NewFromString(req.ValorPago); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "valor_pago inválido",
		})
	}

	dataRecebimento, err := time.Parse("2006-01-02", req.DataRecebimento)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "data_recebimento inválida",
		})
	}

	if _, err := h.marcarRecebimentoUC.Execute(ctx, financial.MarcarRecebimentoInput{
		TenantID:        tenantID,
		ContaID:         contaID,
		DataRecebimento: dataRecebimento,
	}); err != nil {
		h.logger.Error("Erro ao marcar recebimento", zap.Error(err), zap.String("conta_id", contaID))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao marcar recebimento",
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Recebimento marcado com sucesso",
	})
}

// GetContaPagar godoc
// @Summary Buscar conta a pagar
// @Tags Financial
// @Produce json
// @Param id path string true "ID da conta"
// @Success 200 {object} dto.ContaPagarResponse
// @Router /api/v1/financial/payables/{id} [get]
// GetContaPagar godoc
// @Summary Buscar conta a pagar por ID
// @Tags Financial
// @Produce json
// @Param id path string true "ID da conta"
// @Success 200 {object} dto.ContaPagarResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/payables/{id} [get]
// @Security BearerAuth
func (h *FinancialHandler) GetContaPagar(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id do contexto
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		h.logger.Error("tenant_id não encontrado no contexto")
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Extrair ID da URL
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID é obrigatório",
		})
	}

	// Executar use case
	conta, err := h.getContaPagarUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar conta a pagar", zap.Error(err), zap.String("id", id))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar conta a pagar",
		})
	}

	// Converter para response
	response := mapper.ToContaPagarResponse(conta)
	return c.JSON(http.StatusOK, response)
}

// ListContasPagar godoc
// @Summary Listar contas a pagar
// @Tags Financial
// @Produce json
// @Param status query string false "Filtrar por status"
// @Param data_inicio query string false "Data início (YYYY-MM-DD)"
// @Param data_fim query string false "Data fim (YYYY-MM-DD)"
// @Success 200 {array} dto.ContaPagarResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/payables [get]
// @Security BearerAuth
func (h *FinancialHandler) ListContasPagar(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Bind query params
	var req dto.ListContasPagarRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	// Defaults para paginação
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// Parse datas
	var dataInicio, dataFim time.Time
	if req.DataInicio != nil {
		parsed, err := time.Parse("2006-01-02", *req.DataInicio)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "bad_request",
				Message: "Data início inválida (use YYYY-MM-DD)",
			})
		}
		dataInicio = parsed
	}
	if req.DataFim != nil {
		parsed, err := time.Parse("2006-01-02", *req.DataFim)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "bad_request",
				Message: "Data fim inválida (use YYYY-MM-DD)",
			})
		}
		dataFim = parsed
	}

	// Executar use case
	contas, err := h.listContasPagarUC.Execute(ctx, financial.ListContasPagarInput{
		TenantID:   tenantID,
		DataInicio: dataInicio,
		DataFim:    dataFim,
	})
	if err != nil {
		h.logger.Error("Erro ao listar contas a pagar", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar contas",
		})
	}

	// Converter para response
	responses := make([]dto.ContaPagarResponse, len(contas))
	for i, conta := range contas {
		responses[i] = mapper.ToContaPagarResponse(conta)
	}

	return c.JSON(http.StatusOK, responses)
}

// UpdateContaPagar godoc
// @Summary Atualizar conta a pagar
// @Tags Financial
// @Accept json
// @Produce json
// @Param id path string true "ID da conta"
// @Param request body dto.UpdateContaPagarRequest true "Dados para atualização"
// @Success 200 {object} dto.ContaPagarResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/payables/{id} [put]
// @Security BearerAuth
func (h *FinancialHandler) UpdateContaPagar(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Extrair ID
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID é obrigatório",
		})
	}

	// Bind request
	var req dto.UpdateContaPagarRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	// Buscar conta atual
	contaAtual, err := h.getContaPagarUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Conta não encontrada", zap.Error(err))
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Conta não encontrada",
		})
	}

	// Aplicar atualizações parciais
	if req.Descricao != nil {
		contaAtual.Descricao = *req.Descricao
	}
	if req.Fornecedor != nil {
		contaAtual.Fornecedor = *req.Fornecedor
	}
	if req.Observacoes != nil {
		contaAtual.Observacoes = *req.Observacoes
	}

	// Executar use case
	contaAtualizada, err := h.updateContaPagarUC.Execute(ctx, financial.UpdateContaPagarInput{
		TenantID: tenantID,
		ID:       id,
		Conta:    contaAtual,
	})
	if err != nil {
		h.logger.Error("Erro ao atualizar conta", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao atualizar conta",
		})
	}

	response := mapper.ToContaPagarResponse(contaAtualizada)
	return c.JSON(http.StatusOK, response)
}

// DeleteContaPagar godoc
// @Summary Deletar conta a pagar
// @Tags Financial
// @Param id path string true "ID da conta"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/payables/{id} [delete]
// @Security BearerAuth
func (h *FinancialHandler) DeleteContaPagar(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	// Extrair ID
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID é obrigatório",
		})
	}

	// Executar use case
	if err := h.deleteContaPagarUC.Execute(ctx, tenantID, id); err != nil {
		h.logger.Error("Erro ao deletar conta", zap.Error(err), zap.String("id", id))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao deletar conta",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetContaReceber godoc
// @Summary Buscar conta a receber por ID
// @Tags Financial
// @Produce json
// @Param id path string true "ID da conta"
// @Success 200 {object} dto.ContaReceberResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/receivables/{id} [get]
// @Security BearerAuth
func (h *FinancialHandler) GetContaReceber(c echo.Context) error {
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

	conta, err := h.getContaReceberUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar conta a receber", zap.Error(err), zap.String("id", id))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar conta",
		})
	}

	response := mapper.ToContaReceberResponse(conta)
	return c.JSON(http.StatusOK, response)
}

// ListContasReceber godoc
// @Summary Listar contas a receber
// @Tags Financial
// @Produce json
// @Param status query string false "Filtrar por status"
// @Param data_inicio query string false "Data início (YYYY-MM-DD)"
// @Param data_fim query string false "Data fim (YYYY-MM-DD)"
// @Success 200 {array} dto.ContaReceberResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/receivables [get]
// @Security BearerAuth
func (h *FinancialHandler) ListContasReceber(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ListContasReceberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	var dataInicio, dataFim time.Time
	if req.DataInicio != nil {
		parsed, err := time.Parse("2006-01-02", *req.DataInicio)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "bad_request",
				Message: "Data início inválida",
			})
		}
		dataInicio = parsed
	}
	if req.DataFim != nil {
		parsed, err := time.Parse("2006-01-02", *req.DataFim)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "bad_request",
				Message: "Data fim inválida",
			})
		}
		dataFim = parsed
	}

	contas, err := h.listContasReceberUC.Execute(ctx, financial.ListContasReceberInput{
		TenantID:   tenantID,
		DataInicio: dataInicio,
		DataFim:    dataFim,
	})
	if err != nil {
		h.logger.Error("Erro ao listar contas a receber", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar contas",
		})
	}

	responses := make([]dto.ContaReceberResponse, len(contas))
	for i, conta := range contas {
		responses[i] = mapper.ToContaReceberResponse(conta)
	}

	return c.JSON(http.StatusOK, responses)
}

// UpdateContaReceber godoc
// @Summary Atualizar conta a receber
// @Tags Financial
// @Accept json
// @Produce json
// @Param id path string true "ID da conta"
// @Param request body dto.UpdateContaReceberRequest true "Dados para atualização"
// @Success 200 {object} dto.ContaReceberResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/receivables/{id} [put]
// @Security BearerAuth
func (h *FinancialHandler) UpdateContaReceber(c echo.Context) error {
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

	var req dto.UpdateContaReceberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	contaAtual, err := h.getContaReceberUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Conta não encontrada", zap.Error(err))
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Conta não encontrada",
		})
	}

	if req.Origem != nil {
		contaAtual.Origem = *req.Origem
	}
	if req.DescricaoOrigem != nil {
		contaAtual.DescricaoOrigem = *req.DescricaoOrigem
	}
	if req.Observacoes != nil {
		contaAtual.Observacoes = *req.Observacoes
	}

	contaAtualizada, err := h.updateContaReceberUC.Execute(ctx, financial.UpdateContaReceberInput{
		TenantID: tenantID,
		ID:       id,
		Conta:    contaAtual,
	})
	if err != nil {
		h.logger.Error("Erro ao atualizar conta", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao atualizar conta",
		})
	}

	response := mapper.ToContaReceberResponse(contaAtualizada)
	return c.JSON(http.StatusOK, response)
}

// DeleteContaReceber godoc
// @Summary Deletar conta a receber
// @Tags Financial
// @Param id path string true "ID da conta"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/receivables/{id} [delete]
// @Security BearerAuth
func (h *FinancialHandler) DeleteContaReceber(c echo.Context) error {
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

	if err := h.deleteContaReceberUC.Execute(ctx, tenantID, id); err != nil {
		h.logger.Error("Erro ao deletar conta", zap.Error(err), zap.String("id", id))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao deletar conta",
		})
	}

	return c.NoContent(http.StatusNoContent)
} // GetCompensacao godoc
// GetCompensacao busca compensação bancária por ID
// @Summary Buscar compensação bancária
// @Tags Financial
// @Produce json
// @Param id path string true "ID da compensação"
// @Success 200 {object} dto.CompensacaoBancariaResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/compensations/{id} [get]
// @Security BearerAuth
func (h *FinancialHandler) GetCompensacao(c echo.Context) error {
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

	comp, err := h.getCompensacaoUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar compensação", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar compensação",
		})
	}

	response := mapper.ToCompensacaoBancariaResponse(comp)
	return c.JSON(http.StatusOK, response)
}

// ListCompensacoes godoc
// ListCompensacoes lista compensações bancárias
// @Summary Listar compensações bancárias
// @Tags Financial
// @Produce json
// @Param data_inicio query string false "Data início"
// @Param data_fim query string false "Data fim"
// @Success 200 {array} dto.CompensacaoBancariaResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/compensations [get]
// @Security BearerAuth
func (h *FinancialHandler) ListCompensacoes(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ListCompensacoesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	var dataInicio, dataFim time.Time
	if req.DataInicio != nil {
		parsed, _ := time.Parse("2006-01-02", *req.DataInicio)
		dataInicio = parsed
	}
	if req.DataFim != nil {
		parsed, _ := time.Parse("2006-01-02", *req.DataFim)
		dataFim = parsed
	}

	comps, err := h.listCompensacoesUC.Execute(ctx, financial.ListCompensacoesInput{
		TenantID:   tenantID,
		DataInicio: dataInicio,
		DataFim:    dataFim,
	})
	if err != nil {
		h.logger.Error("Erro ao listar compensações", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar compensações",
		})
	}

	responses := make([]dto.CompensacaoBancariaResponse, len(comps))
	for i, comp := range comps {
		responses[i] = mapper.ToCompensacaoBancariaResponse(comp)
	}

	return c.JSON(http.StatusOK, responses)
}

// DeleteCompensacao godoc
// DeleteCompensacao deleta compensação bancária
// @Summary Deletar compensação bancária
// @Tags Financial
// @Param id path string true "ID da compensação"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/compensations/{id} [delete]
// @Security BearerAuth
func (h *FinancialHandler) DeleteCompensacao(c echo.Context) error {
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

	if err := h.deleteCompensacaoUC.Execute(ctx, tenantID, id); err != nil {
		h.logger.Error("Erro ao deletar compensação", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao deletar compensação",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetFluxoCaixa busca fluxo de caixa diário por ID
// @Summary Buscar fluxo de caixa diário
// @Tags Financial
// @Produce json
// @Param id path string true "ID do fluxo"
// @Success 200 {object} dto.FluxoCaixaDiarioResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/cashflow/{id} [get]
// @Security BearerAuth
func (h *FinancialHandler) GetFluxoCaixa(c echo.Context) error {
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

	fluxo, err := h.getFluxoCaixaUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar fluxo de caixa", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar fluxo de caixa",
		})
	}

	response := mapper.ToFluxoCaixaDiarioResponse(fluxo)
	return c.JSON(http.StatusOK, response)
}

// ListFluxoCaixa lista fluxos de caixa diários
// @Summary Listar fluxo de caixa
// @Tags Financial
// @Produce json
// @Param data_inicio query string false "Data início"
// @Param data_fim query string false "Data fim"
// @Success 200 {array} dto.FluxoCaixaResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/cashflow [get]
// @Security BearerAuth
func (h *FinancialHandler) ListFluxoCaixa(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ListFluxoCaixaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	var dataInicio, dataFim time.Time
	if req.DataInicio != nil {
		parsed, _ := time.Parse("2006-01-02", *req.DataInicio)
		dataInicio = parsed
	}
	if req.DataFim != nil {
		parsed, _ := time.Parse("2006-01-02", *req.DataFim)
		dataFim = parsed
	}

	fluxos, err := h.listFluxoCaixaUC.Execute(ctx, financial.ListFluxoCaixaInput{
		TenantID:   tenantID,
		DataInicio: dataInicio,
		DataFim:    dataFim,
	})
	if err != nil {
		h.logger.Error("Erro ao listar fluxo de caixa", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar fluxo de caixa",
		})
	}

	responses := make([]dto.FluxoCaixaDiarioResponse, len(fluxos))
	for i, fluxo := range fluxos {
		responses[i] = mapper.ToFluxoCaixaDiarioResponse(fluxo)
	}

	return c.JSON(http.StatusOK, responses)
}

// GetDRE busca DRE mensal por ID
// @Summary Buscar DRE mensal
// @Tags Financial
// @Produce json
// @Param month path string true "Mês (YYYY-MM)"
// @Success 200 {object} dto.DREResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/dre/{month} [get]
// @Security BearerAuth
func (h *FinancialHandler) GetDRE(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	id := c.Param("month")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Mês/Ano é obrigatório",
		})
	}

	dre, err := h.getDREUC.Execute(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar DRE", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar DRE",
		})
	}

	response := mapper.ToDREMensalResponse(dre)
	return c.JSON(http.StatusOK, response)
}

// ListDRE lista DREs mensais
// @Summary Listar DREs mensais
// @Tags Financial
// @Produce json
// @Param mes_ano_inicio query string false "Mês/Ano início (YYYY-MM)"
// @Param mes_ano_fim query string false "Mês/Ano fim (YYYY-MM)"
// @Success 200 {array} dto.DREMensalResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/dre [get]
// @Security BearerAuth
func (h *FinancialHandler) ListDRE(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ListDRERequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	// Parse MesAno (implementação simplificada - ajustar conforme valueobject.MesAno)
	var inicio, fim valueobject.MesAno
	if req.MesAnoInicio != nil {
		inicio, _ = valueobject.NewMesAno(*req.MesAnoInicio)
	}
	if req.MesAnoFim != nil {
		fim, _ = valueobject.NewMesAno(*req.MesAnoFim)
	}

	dres, err := h.listDREUC.Execute(ctx, financial.ListDREInput{
		TenantID: tenantID,
		Inicio:   inicio,
		Fim:      fim,
	})
	if err != nil {
		h.logger.Error("Erro ao listar DREs", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar DREs",
		})
	}

	responses := make([]dto.DREMensalResponse, len(dres))
	for i, dre := range dres {
		responses[i] = mapper.ToDREMensalResponse(dre)
	}

	return c.JSON(http.StatusOK, responses)
}

// GetDashboard godoc
// @Summary Obter dashboard financeiro
// @Description Retorna dados agregados de payables, receivables, cashflow e DRE
// @Tags Financial
// @Produce json
// @Param start_date query string true "Data de início (YYYY-MM-DD)"
// @Param end_date query string true "Data de fim (YYYY-MM-DD)"
// @Param month query string true "Mês no formato YYYY-MM"
// @Success 200 {object} dto.FinancialDashboardResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/financial/dashboard [get]
// @Security BearerAuth
func (h *FinancialHandler) GetDashboard(c echo.Context) error {
	// TODO: Implementar quando GetDashboardUseCase existir
	return c.JSON(http.StatusNotImplemented, dto.ErrorResponse{
		Error:   "not_implemented",
		Message: "Dashboard financeiro ainda não implementado",
	})
}

// RegisterRoutes registra todas as rotas financeiras
func (h *FinancialHandler) RegisterRoutes(g *echo.Group) {
	// Contas a pagar
	g.POST("/payables", h.CreateContaPagar)
	g.GET("/payables/:id", h.GetContaPagar)
	g.GET("/payables", h.ListContasPagar)
	g.PUT("/payables/:id", h.UpdateContaPagar)
	g.DELETE("/payables/:id", h.DeleteContaPagar)
	g.POST("/payables/:id/pay", h.MarcarPagamento)

	// Contas a receber
	g.POST("/receivables", h.CreateContaReceber)
	g.GET("/receivables/:id", h.GetContaReceber)
	g.GET("/receivables", h.ListContasReceber)
	g.PUT("/receivables/:id", h.UpdateContaReceber)
	g.DELETE("/receivables/:id", h.DeleteContaReceber)
	g.POST("/receivables/:id/receive", h.MarcarRecebimento)

	// Compensações bancárias
	g.GET("/compensations/:id", h.GetCompensacao)
	g.GET("/compensations", h.ListCompensacoes)
	g.DELETE("/compensations/:id", h.DeleteCompensacao)

	// Relatórios
	g.GET("/cashflow/:date", h.GetFluxoCaixa)
	g.GET("/cashflow", h.ListFluxoCaixa)
	g.GET("/dre/:month", h.GetDRE)
	g.GET("/dre", h.ListDRE)

	// Dashboard
	g.GET("/dashboard", h.GetDashboard)
}
