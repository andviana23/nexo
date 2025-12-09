package handler

import (
	"net/http"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/customer"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// CustomerHandler agrupa os handlers de clientes
type CustomerHandler struct {
	createUC         *customer.CreateCustomerUseCase
	updateUC         *customer.UpdateCustomerUseCase
	listUC           *customer.ListCustomersUseCase
	getUC            *customer.GetCustomerUseCase
	getWithHistoryUC *customer.GetCustomerWithHistoryUseCase
	inactivateUC     *customer.InactivateCustomerUseCase
	searchUC         *customer.SearchCustomersUseCase
	exportUC         *customer.ExportCustomerDataUseCase
	statsUC          *customer.GetCustomerStatsUseCase
	checkPhoneUC     *customer.CheckPhoneDuplicateUseCase
	checkCPFUC       *customer.CheckCPFDuplicateUseCase
	logger           *zap.Logger
}

// NewCustomerHandler cria um novo handler de clientes
func NewCustomerHandler(
	createUC *customer.CreateCustomerUseCase,
	updateUC *customer.UpdateCustomerUseCase,
	listUC *customer.ListCustomersUseCase,
	getUC *customer.GetCustomerUseCase,
	getWithHistoryUC *customer.GetCustomerWithHistoryUseCase,
	inactivateUC *customer.InactivateCustomerUseCase,
	searchUC *customer.SearchCustomersUseCase,
	exportUC *customer.ExportCustomerDataUseCase,
	statsUC *customer.GetCustomerStatsUseCase,
	checkPhoneUC *customer.CheckPhoneDuplicateUseCase,
	checkCPFUC *customer.CheckCPFDuplicateUseCase,
	logger *zap.Logger,
) *CustomerHandler {
	return &CustomerHandler{
		createUC:         createUC,
		updateUC:         updateUC,
		listUC:           listUC,
		getUC:            getUC,
		getWithHistoryUC: getWithHistoryUC,
		inactivateUC:     inactivateUC,
		searchUC:         searchUC,
		exportUC:         exportUC,
		statsUC:          statsUC,
		checkPhoneUC:     checkPhoneUC,
		checkCPFUC:       checkCPFUC,
		logger:           logger,
	}
}

// CreateCustomer godoc
// @Summary Criar cliente
// @Description Cria um novo cliente
// @Tags Clientes
// @Accept json
// @Produce json
// @Param request body dto.CreateCustomerRequest true "Dados do cliente"
// @Success 201 {object} dto.CustomerResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Telefone/CPF duplicado"
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers [post]
// @Security BearerAuth
func (h *CustomerHandler) CreateCustomer(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.CreateCustomerRequest
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

	result, err := h.createUC.Execute(ctx, tenantID, req)
	if err != nil {
		h.logger.Error("Erro ao criar cliente", zap.Error(err))

		switch err {
		case domain.ErrCustomerPhoneDuplicate:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "duplicate_phone",
				Message: err.Error(),
			})
		case domain.ErrCustomerCPFDuplicate:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "duplicate_cpf",
				Message: err.Error(),
			})
		case domain.ErrCustomerNameRequired,
			domain.ErrCustomerNameTooShort,
			domain.ErrCustomerPhoneRequired,
			domain.ErrCustomerPhoneInvalid,
			domain.ErrCustomerEmailInvalid,
			domain.ErrCustomerCPFInvalid,
			domain.ErrCustomerCEPInvalid,
			domain.ErrCustomerDateInvalid:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Erro ao criar cliente",
			})
		}
	}

	return c.JSON(http.StatusCreated, mapCustomerToResponse(result))
}

// UpdateCustomer godoc
// @Summary Atualizar cliente
// @Description Atualiza dados de um cliente existente
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path string true "ID do cliente"
// @Param request body dto.UpdateCustomerRequest true "Dados a atualizar"
// @Success 200 {object} dto.CustomerResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/{id} [put]
// @Security BearerAuth
func (h *CustomerHandler) UpdateCustomer(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	customerID := c.Param("id")
	if customerID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do cliente é obrigatório",
		})
	}

	var req dto.UpdateCustomerRequest
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

	result, err := h.updateUC.Execute(ctx, tenantID, customerID, req)
	if err != nil {
		h.logger.Error("Erro ao atualizar cliente", zap.Error(err))

		switch err {
		case domain.ErrCustomerNotFound:
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		case domain.ErrCustomerPhoneDuplicate:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "duplicate_phone",
				Message: err.Error(),
			})
		case domain.ErrCustomerCPFDuplicate:
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "duplicate_cpf",
				Message: err.Error(),
			})
		case domain.ErrCustomerNameRequired,
			domain.ErrCustomerNameTooShort,
			domain.ErrCustomerPhoneInvalid,
			domain.ErrCustomerEmailInvalid,
			domain.ErrCustomerCPFInvalid,
			domain.ErrCustomerCEPInvalid,
			domain.ErrCustomerDateInvalid:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Erro ao atualizar cliente",
			})
		}
	}

	return c.JSON(http.StatusOK, mapCustomerToResponse(result))
}

// ListCustomers godoc
// @Summary Listar clientes
// @Description Lista clientes com paginação e filtros
// @Tags Clientes
// @Accept json
// @Produce json
// @Param search query string false "Busca por nome, telefone, email ou CPF"
// @Param ativo query bool false "Filtrar por status ativo"
// @Param tags query []string false "Filtrar por tags"
// @Param order_by query string false "Ordenar por (nome, criado_em, atualizado_em)"
// @Param page query int false "Página" default(1)
// @Param page_size query int false "Tamanho da página" default(20)
// @Success 200 {object} dto.ListCustomersResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers [get]
// @Security BearerAuth
func (h *CustomerHandler) ListCustomers(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ListCustomersRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	// Defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	filter := port.CustomerFilter{
		Search:   req.Search,
		Ativo:    req.Ativo,
		Tags:     req.Tags,
		OrderBy:  req.OrderBy,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	customers, total, err := h.listUC.Execute(ctx, tenantID, filter)
	if err != nil {
		h.logger.Error("Erro ao listar clientes", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar clientes",
		})
	}

	data := make([]dto.CustomerResponse, 0, len(customers))
	for _, c := range customers {
		data = append(data, mapCustomerToResponse(c))
	}

	return c.JSON(http.StatusOK, dto.ListCustomersResponse{
		Data:     data,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	})
}

// GetCustomer godoc
// @Summary Buscar cliente
// @Description Busca um cliente pelo ID
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path string true "ID do cliente"
// @Success 200 {object} dto.CustomerResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/{id} [get]
// @Security BearerAuth
func (h *CustomerHandler) GetCustomer(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	customerID := c.Param("id")
	if customerID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do cliente é obrigatório",
		})
	}

	result, err := h.getUC.Execute(ctx, tenantID, customerID)
	if err != nil {
		h.logger.Error("Erro ao buscar cliente", zap.Error(err))

		if err == domain.ErrCustomerNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar cliente",
		})
	}

	return c.JSON(http.StatusOK, mapCustomerToResponse(result))
}

// GetCustomerWithHistory godoc
// @Summary Buscar cliente com histórico
// @Description Busca um cliente com histórico de atendimentos
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path string true "ID do cliente"
// @Success 200 {object} dto.CustomerWithHistoryResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/{id}/history [get]
// @Security BearerAuth
func (h *CustomerHandler) GetCustomerWithHistory(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	customerID := c.Param("id")
	if customerID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do cliente é obrigatório",
		})
	}

	result, err := h.getWithHistoryUC.Execute(ctx, tenantID, customerID)
	if err != nil {
		h.logger.Error("Erro ao buscar cliente com histórico", zap.Error(err))

		if err == domain.ErrCustomerNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar cliente",
		})
	}

	return c.JSON(http.StatusOK, mapCustomerWithHistoryToResponse(result))
}

// InactivateCustomer godoc
// @Summary Inativar cliente
// @Description Inativa um cliente (soft delete)
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path string true "ID do cliente"
// @Success 204 "Sem conteúdo"
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/{id} [delete]
// @Security BearerAuth
func (h *CustomerHandler) InactivateCustomer(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	customerID := c.Param("id")
	if customerID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do cliente é obrigatório",
		})
	}

	err := h.inactivateUC.Execute(ctx, tenantID, customerID)
	if err != nil {
		h.logger.Error("Erro ao inativar cliente", zap.Error(err))

		if err == domain.ErrCustomerNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao inativar cliente",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// SearchCustomers godoc
// @Summary Buscar clientes
// @Description Busca rápida de clientes por nome, telefone ou email
// @Tags Clientes
// @Accept json
// @Produce json
// @Param q query string true "Termo de busca"
// @Success 200 {array} dto.CustomerSummaryResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/search [get]
// @Security BearerAuth
func (h *CustomerHandler) SearchCustomers(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	query := c.QueryParam("q")
	if len(query) < 2 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Termo de busca deve ter pelo menos 2 caracteres",
		})
	}

	results, err := h.searchUC.Execute(ctx, tenantID, query)
	if err != nil {
		h.logger.Error("Erro ao buscar clientes", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao buscar clientes",
		})
	}

	data := make([]dto.CustomerSummaryResponse, 0, len(results))
	for _, r := range results {
		data = append(data, dto.CustomerSummaryResponse{
			ID:       r.ID,
			Nome:     r.Nome,
			Telefone: r.Telefone,
			Email:    r.Email,
			Tags:     r.Tags,
		})
	}

	return c.JSON(http.StatusOK, data)
}

// ExportCustomerData godoc
// @Summary Exportar dados do cliente (LGPD)
// @Description Exporta todos os dados de um cliente para LGPD
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path string true "ID do cliente"
// @Success 200 {object} dto.CustomerExportResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/{id}/export [get]
// @Security BearerAuth
func (h *CustomerHandler) ExportCustomerData(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	customerID := c.Param("id")
	if customerID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do cliente é obrigatório",
		})
	}

	result, err := h.exportUC.Execute(ctx, tenantID, customerID)
	if err != nil {
		h.logger.Error("Erro ao exportar dados do cliente", zap.Error(err))

		if err == domain.ErrCustomerNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao exportar dados",
		})
	}

	return c.JSON(http.StatusOK, mapCustomerExportToResponse(result))
}

// GetCustomerStats godoc
// @Summary Estatísticas de clientes
// @Description Retorna estatísticas dos clientes do tenant
// @Tags Clientes
// @Accept json
// @Produce json
// @Success 200 {object} dto.CustomerStatsResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/stats [get]
// @Security BearerAuth
func (h *CustomerHandler) GetCustomerStats(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	result, err := h.statsUC.Execute(ctx, tenantID)
	if err != nil {
		h.logger.Error("Erro ao obter estatísticas", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao obter estatísticas",
		})
	}

	return c.JSON(http.StatusOK, dto.CustomerStatsResponse{
		TotalAtivos:        result.TotalAtivos,
		TotalInativos:      result.TotalInativos,
		NovosUltimos30Dias: result.NovosUltimos30Dias,
		TotalGeral:         result.TotalGeral,
	})
}

// CheckPhoneExists godoc
// @Summary Verificar telefone
// @Description Verifica se telefone já está cadastrado
// @Tags Clientes
// @Accept json
// @Produce json
// @Param telefone query string true "Telefone a verificar"
// @Param exclude_id query string false "ID a excluir da verificação"
// @Success 200 {object} dto.CheckExistsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/check-phone [get]
// @Security BearerAuth
func (h *CustomerHandler) CheckPhoneExists(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	phone := c.QueryParam("telefone")
	if phone == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Telefone é obrigatório",
		})
	}

	var excludeID *string
	if eid := c.QueryParam("exclude_id"); eid != "" {
		excludeID = &eid
	}

	exists, err := h.checkPhoneUC.Execute(ctx, tenantID, phone, excludeID)
	if err != nil {
		h.logger.Error("Erro ao verificar telefone", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao verificar telefone",
		})
	}

	return c.JSON(http.StatusOK, dto.CheckExistsResponse{Exists: exists})
}

// CheckCPFExists godoc
// @Summary Verificar CPF
// @Description Verifica se CPF já está cadastrado
// @Tags Clientes
// @Accept json
// @Produce json
// @Param cpf query string true "CPF a verificar"
// @Param exclude_id query string false "ID a excluir da verificação"
// @Success 200 {object} dto.CheckExistsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/check-cpf [get]
// @Security BearerAuth
func (h *CustomerHandler) CheckCPFExists(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	cpf := c.QueryParam("cpf")
	if cpf == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "CPF é obrigatório",
		})
	}

	var excludeID *string
	if eid := c.QueryParam("exclude_id"); eid != "" {
		excludeID = &eid
	}

	exists, err := h.checkCPFUC.Execute(ctx, tenantID, cpf, excludeID)
	if err != nil {
		h.logger.Error("Erro ao verificar CPF", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao verificar CPF",
		})
	}

	return c.JSON(http.StatusOK, dto.CheckExistsResponse{Exists: exists})
}

// =============================================================================
// Mappers
// =============================================================================

func mapCustomerToResponse(c *entity.Customer) dto.CustomerResponse {
	if c == nil {
		return dto.CustomerResponse{}
	}

	var dataNasc *string
	if c.DataNascimento != nil {
		d := c.DataNascimento.Format("2006-01-02")
		dataNasc = &d
	}

	tags := c.Tags
	if tags == nil {
		tags = []string{}
	}

	return dto.CustomerResponse{
		ID:                  c.ID,
		TenantID:            c.TenantID.String(),
		Nome:                c.Nome,
		Telefone:            c.Telefone,
		Email:               c.Email,
		CPF:                 c.CPF,
		DataNascimento:      dataNasc,
		Genero:              c.Genero,
		EnderecoLogradouro:  c.EnderecoLogradouro,
		EnderecoNumero:      c.EnderecoNumero,
		EnderecoComplemento: c.EnderecoComplemento,
		EnderecoBairro:      c.EnderecoBairro,
		EnderecoCidade:      c.EnderecoCidade,
		EnderecoEstado:      c.EnderecoEstado,
		EnderecoCEP:         c.EnderecoCEP,
		Observacoes:         c.Observacoes,
		Tags:                tags,
		Ativo:               c.Ativo,
		CreatedAt:           c.CreatedAt,
		UpdatedAt:           c.UpdatedAt,
	}
}

func mapCustomerWithHistoryToResponse(cwh *port.CustomerWithHistory) dto.CustomerWithHistoryResponse {
	if cwh == nil {
		return dto.CustomerWithHistoryResponse{}
	}

	return dto.CustomerWithHistoryResponse{
		CustomerResponse:    mapCustomerToResponse(cwh.Customer),
		TotalAtendimentos:   cwh.TotalAtendimentos,
		TotalGasto:          cwh.TotalGasto,
		TicketMedio:         cwh.TicketMedio,
		UltimoAtendimento:   cwh.UltimoAtendimento,
		FrequenciaMediaDias: cwh.FrequenciaMediaDias,
	}
}

func mapCustomerExportToResponse(ce *port.CustomerExport) dto.CustomerExportResponse {
	if ce == nil {
		return dto.CustomerExportResponse{}
	}

	var endereco *dto.CustomerAddressExport
	if ce.Customer.EnderecoLogradouro != nil {
		endereco = &dto.CustomerAddressExport{
			Logradouro:  ce.Customer.EnderecoLogradouro,
			Numero:      ce.Customer.EnderecoNumero,
			Complemento: ce.Customer.EnderecoComplemento,
			Bairro:      ce.Customer.EnderecoBairro,
			Cidade:      ce.Customer.EnderecoCidade,
			Estado:      ce.Customer.EnderecoEstado,
			CEP:         ce.Customer.EnderecoCEP,
		}
	}

	var dataNasc *string
	if ce.Customer.DataNascimento != nil {
		d := ce.Customer.DataNascimento.Format("2006-01-02")
		dataNasc = &d
	}

	response := dto.CustomerExportResponse{
		DataExportacao: time.Now(),
	}

	response.DadosPessoais.Nome = ce.Customer.Nome
	response.DadosPessoais.Email = ce.Customer.Email
	response.DadosPessoais.Telefone = ce.Customer.Telefone
	response.DadosPessoais.CPF = ce.Customer.CPF
	response.DadosPessoais.DataNascimento = dataNasc
	response.DadosPessoais.Genero = ce.Customer.Genero
	response.DadosPessoais.Endereco = endereco

	response.Metricas.TotalGasto = ce.TotalGasto
	response.Metricas.TicketMedio = ce.TicketMedio
	response.Metricas.TotalVisitas = ce.TotalVisitas

	historico := make([]dto.CustomerAppointmentExport, 0, len(ce.HistoricoAtendimentos))
	for _, h := range ce.HistoricoAtendimentos {
		historico = append(historico, dto.CustomerAppointmentExport{
			Data:         h.Data,
			Status:       h.Status,
			Profissional: h.Profissional,
			ValorTotal:   h.ValorTotal,
		})
	}
	response.HistoricoAtendimentos = historico

	return response
}
