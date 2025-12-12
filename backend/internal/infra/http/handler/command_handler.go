package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/command"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// CommandHandler agrupa os handlers de comandas
type CommandHandler struct {
	createUC             *command.CreateCommandUseCase
	getUC                *command.GetCommandUseCase
	listUC               *command.ListCommandsUseCase
	getByAppointmentUC   *command.GetCommandByAppointmentUseCase
	addItemUC            *command.AddCommandItemUseCase
	removeItemUC         *command.RemoveCommandItemUseCase
	addPaymentUC         *command.AddCommandPaymentUseCase
	removePaymentUC      *command.RemoveCommandPaymentUseCase
	closeUC              *command.CloseCommandUseCase
	finalizarIntegradaUC *command.FinalizarComandaIntegradaUseCase
	cancelUC             *command.CancelCommandUseCase // T-EST-003: Cancelamento com reversão de estoque
	logger               *zap.Logger
}

// NewCommandHandler cria um novo handler de comandas
func NewCommandHandler(
	createUC *command.CreateCommandUseCase,
	getUC *command.GetCommandUseCase,
	listUC *command.ListCommandsUseCase,
	getByAppointmentUC *command.GetCommandByAppointmentUseCase,
	addItemUC *command.AddCommandItemUseCase,
	removeItemUC *command.RemoveCommandItemUseCase,
	addPaymentUC *command.AddCommandPaymentUseCase,
	removePaymentUC *command.RemoveCommandPaymentUseCase,
	closeUC *command.CloseCommandUseCase,
	finalizarIntegradaUC *command.FinalizarComandaIntegradaUseCase,
	cancelUC *command.CancelCommandUseCase,
	logger *zap.Logger,
) *CommandHandler {
	return &CommandHandler{
		createUC:             createUC,
		getUC:                getUC,
		listUC:               listUC,
		getByAppointmentUC:   getByAppointmentUC,
		addItemUC:            addItemUC,
		removeItemUC:         removeItemUC,
		addPaymentUC:         addPaymentUC,
		removePaymentUC:      removePaymentUC,
		closeUC:              closeUC,
		finalizarIntegradaUC: finalizarIntegradaUC,
		cancelUC:             cancelUC,
		logger:               logger,
	}
}

// CreateCommand godoc
// @Summary Criar comanda
// @Description Cria uma nova comanda para um agendamento
// @Tags Comandas
// @Accept json
// @Produce json
// @Param request body dto.CreateCommandRequest true "Dados da comanda"
// @Success 201 {object} dto.CommandResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/commands [post]
func (h *CommandHandler) CreateCommand(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id do contexto (JWT middleware)
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		h.logger.Error("failed to get tenant_id from context", zap.Error(err))
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	// Parse request
	var req dto.CreateCommandRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("failed to bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// Validar request
	if err := c.Validate(&req); err != nil {
		h.logger.Error("validation failed", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Executar use case
	response, err := h.createUC.Execute(ctx, tenantID, &req)
	if err != nil {
		h.logger.Error("failed to create command", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, response)
}

// GetCommand godoc
// @Summary Buscar comanda
// @Description Busca uma comanda por ID (com items e payments)
// @Tags Comandas
// @Produce json
// @Param id path string true "ID da comanda"
// @Success 200 {object} dto.CommandResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/commands/{id} [get]
func (h *CommandHandler) GetCommand(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	// Parse command ID
	commandID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid command_id"})
	}

	// Executar use case
	response, err := h.getUC.Execute(ctx, commandID, tenantID)
	if err != nil {
		if err.Error() == "command not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "command not found"})
		}
		h.logger.Error("failed to get command", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

// ListCommands godoc
// @Summary Listar comandas
// @Description Lista comandas com filtros e paginação
// @Tags Comandas
// @Produce json
// @Param status query string false "Filtrar por status (ABERTA, FECHADA, CANCELADA)"
// @Param customer_id query string false "Filtrar por cliente (UUID)"
// @Param date_from query string false "Data inicial (YYYY-MM-DD)"
// @Param date_to query string false "Data final (YYYY-MM-DD)"
// @Param page query int false "Página (default: 1)"
// @Param page_size query int false "Itens por página (default: 20, max: 100)"
// @Success 200 {object} dto.CommandListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/commands [get]
func (h *CommandHandler) ListCommands(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	// Parse query params
	input := command.ListCommandsInput{
		TenantID: tenantID,
		Page:     1,
		PageSize: 20,
	}

	// Status
	if status := c.QueryParam("status"); status != "" {
		input.Status = &status
	}

	// Customer ID
	if customerID := c.QueryParam("customer_id"); customerID != "" {
		input.CustomerID = &customerID
	}

	// Date range
	if dateFrom := c.QueryParam("date_from"); dateFrom != "" {
		input.DateFrom = &dateFrom
	}
	if dateTo := c.QueryParam("date_to"); dateTo != "" {
		input.DateTo = &dateTo
	}

	// Pagination
	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			input.Page = p
		}
	}
	if pageSize := c.QueryParam("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			input.PageSize = ps
		}
	}

	// Executar use case
	response, err := h.listUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("failed to list commands", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

// GetCommandByAppointment godoc
// @Summary Buscar comanda por agendamento
// @Description Busca uma comanda pelo ID do agendamento vinculado
// @Tags Comandas
// @Produce json
// @Param appointmentId path string true "ID do agendamento"
// @Success 200 {object} dto.CommandResponse
// @Failure 404 {object} map[string]string "Comanda não encontrada para este agendamento"
// @Failure 500 {object} map[string]string
// @Router /api/v1/commands/by-appointment/{appointmentId} [get]
func (h *CommandHandler) GetCommandByAppointment(c echo.Context) error {
	ctx := c.Request().Context()

	// Extrair tenant_id
	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	// Parse appointment ID
	appointmentID, err := uuid.Parse(c.Param("appointmentId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid appointment_id"})
	}

	// Executar use case
	response, err := h.getByAppointmentUC.Execute(ctx, appointmentID, tenantID)
	if err != nil {
		h.logger.Error("failed to get command by appointment", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if response == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "command not found for this appointment"})
	}

	return c.JSON(http.StatusOK, response)
}

// AddCommandItem godoc
// @Summary Adicionar item à comanda
// @Description Adiciona um item (serviço/produto) à comanda
// @Tags Comandas
// @Accept json
// @Produce json
// @Param id path string true "ID da comanda"
// @Param request body dto.AddCommandItemRequest true "Dados do item"
// @Success 200 {object} dto.CommandResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/commands/{id}/items [post]
func (h *CommandHandler) AddCommandItem(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	commandID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid command_id"})
	}

	var req dto.AddCommandItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// T-EST-001: Execute agora valida estoque para itens PRODUTO
	response, err := h.addItemUC.Execute(ctx, commandID, tenantID, userID, &req)
	if err != nil {
		h.logger.Error("failed to add item", zap.Error(err))
		// Retorna erro apropriado baseado no tipo
		if strings.Contains(err.Error(), "estoque") || strings.Contains(err.Error(), "inativo") {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

// RemoveCommandItem godoc
// @Summary Remover item da comanda
// @Description Remove um item da comanda
// @Tags Comandas
// @Produce json
// @Param id path string true "ID da comanda"
// @Param itemId path string true "ID do item"
// @Success 200 {object} dto.CommandResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/commands/{id}/items/{itemId} [delete]
func (h *CommandHandler) RemoveCommandItem(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	commandID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid command_id"})
	}

	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid item_id"})
	}

	response, err := h.removeItemUC.Execute(ctx, commandID, itemID, tenantID)
	if err != nil {
		h.logger.Error("failed to remove item", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

// AddCommandPayment godoc
// @Summary Adicionar pagamento à comanda
// @Description Registra um pagamento na comanda
// @Tags Comandas
// @Accept json
// @Produce json
// @Param id path string true "ID da comanda"
// @Param request body dto.AddCommandPaymentRequest true "Dados do pagamento"
// @Success 200 {object} dto.CommandResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/commands/{id}/payments [post]
func (h *CommandHandler) AddCommandPayment(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	commandID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid command_id"})
	}

	var req dto.AddCommandPaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// As taxas são buscadas automaticamente pelo UseCase a partir do MeioPagamento
	response, err := h.addPaymentUC.Execute(ctx, commandID, tenantID, userID, &req)
	if err != nil {
		h.logger.Error("failed to add payment", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

// RemoveCommandPayment godoc
// @Summary Remover pagamento da comanda
// @Description Remove um pagamento da comanda
// @Tags Comandas
// @Produce json
// @Param id path string true "ID da comanda"
// @Param paymentId path string true "ID do pagamento"
// @Success 200 {object} dto.CommandResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/commands/{id}/payments/{paymentId} [delete]
func (h *CommandHandler) RemoveCommandPayment(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	commandID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid command_id"})
	}

	paymentID, err := uuid.Parse(c.Param("paymentId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payment_id"})
	}

	response, err := h.removePaymentUC.Execute(ctx, commandID, paymentID, tenantID)
	if err != nil {
		h.logger.Error("failed to remove payment", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

// CloseCommand godoc
// @Summary Fechar comanda
// @Description Fecha a comanda e finaliza o atendimento
// @Tags Comandas
// @Accept json
// @Produce json
// @Param id path string true "ID da comanda"
// @Param request body dto.CloseCommandRequest true "Opções de fechamento"
// @Success 200 {object} dto.CommandResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/commands/{id}/close [post]
func (h *CommandHandler) CloseCommand(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	commandID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid command_id"})
	}

	var req dto.CloseCommandRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	response, err := h.closeUC.Execute(ctx, commandID, tenantID, userID, &req)
	if err != nil {
		h.logger.Error("failed to close command", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

// CloseCommandIntegrated godoc
// @Summary Fechar comanda com integração completa
// @Description Fecha a comanda com integração financeira, estoque e comissões
// @Tags Comandas
// @Accept json
// @Produce json
// @Param id path string true "ID da comanda"
// @Param request body dto.CloseCommandRequest true "Opções de fechamento"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/commands/{id}/close-integrated [post]
func (h *CommandHandler) CloseCommandIntegrated(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	commandID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid command_id"})
	}

	var req dto.CloseCommandRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// Executar UC de finalização integrada
	input := command.FinalizarComandaIntegradaInput{
		CommandID:          commandID,
		TenantID:           tenantID,
		UserID:             userID,
		DeixarTrocoGorjeta: req.DeixarTrocoGorjeta,
		DeixarSaldoDivida:  req.DeixarSaldoDivida,
		Observacoes:        req.Observacoes,
	}

	output, err := h.finalizarIntegradaUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("failed to close command with integration", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Retornar resposta completa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"command":               output.Command,
		"contas_receber":        output.ContasReceber,
		"operacoes_caixa":       output.OperacoesCaixa,
		"commission_items":      output.CommissionItems,
		"movimentacoes_estoque": output.MovimentacoesEstoque,
		"total_caixa":           output.TotalLancadoCaixa.String(),
		"total_contas_receber":  output.TotalContasReceber.String(),
		"total_comissoes":       output.TotalComissoes.String(),
	})
}

// CancelCommand godoc
// @Summary Cancelar comanda
// @Description Cancela uma comanda e reverte o estoque se necessário (T-EST-003)
// @Tags Comandas
// @Accept json
// @Produce json
// @Param id path string true "ID da comanda"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string "Comanda já cancelada"
// @Failure 500 {object} map[string]string
// @Router /api/v1/commands/{id}/cancel [post]
func (h *CommandHandler) CancelCommand(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, err := getTenantIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	commandID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid command_id"})
	}

	// Bind do motivo (opcional)
	var req struct {
		Motivo string `json:"motivo"`
	}
	c.Bind(&req)

	input := command.CancelCommandInput{
		CommandID: commandID,
		TenantID:  tenantID,
		UserID:    userID,
		Motivo:    req.Motivo,
	}

	output, err := h.cancelUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("failed to cancel command", zap.Error(err))
		// Verificar tipo de erro
		if strings.Contains(err.Error(), "já está cancelada") {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "não encontrada") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"command":                     output.Command,
		"estoque_revertido":           output.EstoqueRevertido,
		"quantidade_itens_revertidos": output.QuantidadeItensRevertidos,
		"movimentacoes_estoque":       output.MovimentacoesEstoque,
	})
}

// ============================================================================
// Helper Functions
// ============================================================================

func getTenantIDFromContext(c echo.Context) (uuid.UUID, error) {
	tenantIDStr, ok := c.Get("tenant_id").(string)
	if !ok {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "tenant_id not found in context")
	}
	return uuid.Parse(tenantIDStr)
}

func getUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "user_id not found in context")
	}
	return uuid.Parse(userIDStr)
}
