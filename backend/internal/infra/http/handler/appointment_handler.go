package handler

import (
	"net/http"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/appointment"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AppointmentHandler agrupa os handlers de agendamentos
type AppointmentHandler struct {
	createUC       *appointment.CreateAppointmentUseCase
	listUC         *appointment.ListAppointmentsUseCase
	getUC          *appointment.GetAppointmentUseCase
	updateStatusUC *appointment.UpdateAppointmentStatusUseCase
	rescheduleUC   *appointment.RescheduleAppointmentUseCase
	cancelUC       *appointment.CancelAppointmentUseCase
	logger         *zap.Logger
}

// NewAppointmentHandler cria um novo handler de agendamentos
func NewAppointmentHandler(
	createUC *appointment.CreateAppointmentUseCase,
	listUC *appointment.ListAppointmentsUseCase,
	getUC *appointment.GetAppointmentUseCase,
	updateStatusUC *appointment.UpdateAppointmentStatusUseCase,
	rescheduleUC *appointment.RescheduleAppointmentUseCase,
	cancelUC *appointment.CancelAppointmentUseCase,
	logger *zap.Logger,
) *AppointmentHandler {
	return &AppointmentHandler{
		createUC:       createUC,
		listUC:         listUC,
		getUC:          getUC,
		updateStatusUC: updateStatusUC,
		rescheduleUC:   rescheduleUC,
		cancelUC:       cancelUC,
		logger:         logger,
	}
}

// CreateAppointment godoc
// @Summary Criar agendamento
// @Description Cria um novo agendamento
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param request body dto.CreateAppointmentRequest true "Dados do agendamento"
// @Success 201 {object} dto.AppointmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Conflito de horário"
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/appointments [post]
// @Security BearerAuth
func (h *AppointmentHandler) CreateAppointment(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.CreateAppointmentRequest
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

	input := appointment.CreateAppointmentInput{
		TenantID:       tenantID,
		ProfessionalID: req.ProfessionalID,
		CustomerID:     req.CustomerID,
		StartTime:      req.StartTime,
		ServiceIDs:     req.ServiceIDs,
		Notes:          req.Notes,
	}

	result, err := h.createUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao criar agendamento", zap.Error(err))

		// Verificar tipo de erro para status code apropriado
		switch err.Error() {
		case "conflito de horário com outro agendamento":
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "create_error",
				Message: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusCreated, mapper.AppointmentToResponse(result))
}

// ListAppointments godoc
// @Summary Listar agendamentos
// @Description Lista agendamentos com filtros e paginação
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param professional_id query string false "ID do profissional"
// @Param customer_id query string false "ID do cliente"
// @Param status query string false "Status (CREATED, CONFIRMED, IN_SERVICE, DONE, NO_SHOW, CANCELED)"
// @Param start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param end_date query string false "Data final (YYYY-MM-DD)"
// @Param page query int false "Página" default(1)
// @Param page_size query int false "Tamanho da página" default(20)
// @Success 200 {object} dto.ListAppointmentsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/appointments [get]
// @Security BearerAuth
func (h *AppointmentHandler) ListAppointments(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ListAppointmentsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Parâmetros inválidos",
		})
	}

	// Parse dates
	var startDate, endDate time.Time
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Data inicial inválida (formato: YYYY-MM-DD)",
			})
		}
		startDate = parsed
	}
	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Data final inválida (formato: YYYY-MM-DD)",
			})
		}
		endDate = parsed.Add(24 * time.Hour) // Incluir o dia todo
	}

	// Parse status
	var status valueobject.AppointmentStatus
	if req.Status != "" {
		parsed, valid := valueobject.ParseAppointmentStatus(req.Status)
		if !valid {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Status inválido",
			})
		}
		status = parsed
	}

	input := appointment.ListAppointmentsInput{
		TenantID:       tenantID,
		ProfessionalID: req.ProfessionalID,
		CustomerID:     req.CustomerID,
		Status:         status,
		StartDate:      startDate,
		EndDate:        endDate,
		Page:           req.Page,
		PageSize:       req.PageSize,
	}

	result, err := h.listUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao listar agendamentos", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "list_error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.ListAppointmentsResponse{
		Data:     mapper.AppointmentsToResponse(result.Appointments),
		Page:     result.Page,
		PageSize: result.PageSize,
		Total:    result.Total,
	})
}

// GetAppointment godoc
// @Summary Buscar agendamento
// @Description Busca um agendamento por ID
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do agendamento"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/appointments/{id} [get]
// @Security BearerAuth
func (h *AppointmentHandler) GetAppointment(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	appointmentID := c.Param("id")
	if appointmentID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do agendamento é obrigatório",
		})
	}

	input := appointment.GetAppointmentInput{
		TenantID:      tenantID,
		AppointmentID: appointmentID,
	}

	result, err := h.getUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao buscar agendamento", zap.Error(err))
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Agendamento não encontrado",
		})
	}

	return c.JSON(http.StatusOK, mapper.AppointmentToResponse(result))
}

// UpdateAppointmentStatus godoc
// @Summary Atualizar status do agendamento
// @Description Atualiza o status de um agendamento
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do agendamento"
// @Param request body dto.UpdateAppointmentStatusRequest true "Novo status"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/appointments/{id}/status [patch]
// @Security BearerAuth
func (h *AppointmentHandler) UpdateAppointmentStatus(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	appointmentID := c.Param("id")
	if appointmentID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do agendamento é obrigatório",
		})
	}

	var req dto.UpdateAppointmentStatusRequest
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

	status, valid := valueobject.ParseAppointmentStatus(req.Status)
	if !valid {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Status inválido",
		})
	}

	input := appointment.UpdateAppointmentStatusInput{
		TenantID:      tenantID,
		AppointmentID: appointmentID,
		NewStatus:     status,
		Reason:        req.Reason,
	}

	result, err := h.updateStatusUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao atualizar status", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "update_error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, mapper.AppointmentToResponse(result))
}

// RescheduleAppointment godoc
// @Summary Reagendar agendamento
// @Description Reagenda um agendamento para novo horário
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do agendamento"
// @Param request body dto.RescheduleAppointmentRequest true "Novo horário"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Conflito de horário"
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/appointments/{id}/reschedule [patch]
// @Security BearerAuth
func (h *AppointmentHandler) RescheduleAppointment(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	appointmentID := c.Param("id")
	if appointmentID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do agendamento é obrigatório",
		})
	}

	var req dto.RescheduleAppointmentRequest
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

	input := appointment.RescheduleAppointmentInput{
		TenantID:       tenantID,
		AppointmentID:  appointmentID,
		NewStartTime:   req.NewStartTime,
		ProfessionalID: req.ProfessionalID,
	}

	result, err := h.rescheduleUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao reagendar", zap.Error(err))

		if err.Error() == "conflito de horário com outro agendamento" {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "reschedule_error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, mapper.AppointmentToResponse(result))
}

// CancelAppointment godoc
// @Summary Cancelar agendamento
// @Description Cancela um agendamento
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do agendamento"
// @Param request body dto.CancelAppointmentRequest true "Motivo do cancelamento"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/appointments/{id}/cancel [post]
// @Security BearerAuth
func (h *AppointmentHandler) CancelAppointment(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	appointmentID := c.Param("id")
	if appointmentID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "ID do agendamento é obrigatório",
		})
	}

	var req dto.CancelAppointmentRequest
	if err := c.Bind(&req); err != nil {
		// Permitir body vazio
		req = dto.CancelAppointmentRequest{}
	}

	input := appointment.CancelAppointmentInput{
		TenantID:      tenantID,
		AppointmentID: appointmentID,
		Reason:        req.Reason,
	}

	result, err := h.cancelUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao cancelar agendamento", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "cancel_error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, mapper.AppointmentToResponse(result))
}
