package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/appointment"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AppointmentHandler agrupa os handlers de agendamentos
type AppointmentHandler struct {
	createUC            *appointment.CreateAppointmentUseCase
	listUC              *appointment.ListAppointmentsUseCase
	getUC               *appointment.GetAppointmentUseCase
	updateStatusUC      *appointment.UpdateAppointmentStatusUseCase
	rescheduleUC        *appointment.RescheduleAppointmentUseCase
	cancelUC            *appointment.CancelAppointmentUseCase
	finishWithCommandUC *appointment.FinishServiceWithCommandUseCase
	logger              *zap.Logger
}

// NewAppointmentHandler cria um novo handler de agendamentos
func NewAppointmentHandler(
	createUC *appointment.CreateAppointmentUseCase,
	listUC *appointment.ListAppointmentsUseCase,
	getUC *appointment.GetAppointmentUseCase,
	updateStatusUC *appointment.UpdateAppointmentStatusUseCase,
	rescheduleUC *appointment.RescheduleAppointmentUseCase,
	cancelUC *appointment.CancelAppointmentUseCase,
	finishWithCommandUC *appointment.FinishServiceWithCommandUseCase,
	logger *zap.Logger,
) *AppointmentHandler {
	return &AppointmentHandler{
		createUC:            createUC,
		listUC:              listUC,
		getUC:               getUC,
		updateStatusUC:      updateStatusUC,
		rescheduleUC:        rescheduleUC,
		cancelUC:            cancelUC,
		finishWithCommandUC: finishWithCommandUC,
		logger:              logger,
	}
}

// enforceBarberScope garante que um barbeiro só atue nos próprios agendamentos.
func (h *AppointmentHandler) enforceBarberScope(ctx context.Context, c echo.Context, tenantID, appointmentID string) error {
	if !middleware.IsBarber(c) {
		return nil
	}

	barberProfID := middleware.GetProfessionalIDForBarber(c)
	if barberProfID == "" {
		return echo.NewHTTPError(http.StatusForbidden, dto.ErrorResponse{
			Error:   "forbidden",
			Message: "Acesso negado: profissional não associado",
		})
	}

	appt, err := h.getUC.Execute(ctx, appointment.GetAppointmentInput{
		TenantID:      tenantID,
		AppointmentID: appointmentID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAppointmentNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Agendamento não encontrado",
			})
		}
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "appointment_error",
			Message: err.Error(),
		})
	}

	if appt.ProfessionalID != barberProfID {
		return echo.NewHTTPError(http.StatusForbidden, dto.ErrorResponse{
			Error:   "forbidden",
			Message: "Acesso negado: você só pode agir nos seus agendamentos",
		})
	}

	return nil
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

	// RBAC: Barbeiro só pode criar para si mesmo
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		if barberProfID == "" {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "forbidden", Message: "Acesso negado: profissional não associado"})
		}
		if req.ProfessionalID != barberProfID {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "forbidden", Message: "Barbeiro só pode criar agendamento para si"})
		}
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
		switch {
		case errors.Is(err, domain.ErrAppointmentConflict),
			errors.Is(err, domain.ErrAppointmentBlockedTimeConflict),
			errors.Is(err, domain.ErrAppointmentMinimumInterval):
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "conflict", Message: err.Error()})
		case errors.Is(err, domain.ErrAppointmentProfessionalNotFound),
			errors.Is(err, domain.ErrAppointmentCustomerNotFound),
			errors.Is(err, domain.ErrAppointmentServiceNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: err.Error()})
		default:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "create_error", Message: err.Error()})
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

	// Parse dates - aceita YYYY-MM-DD ou ISO8601 (com timezone)
	var startDate, endDate time.Time
	if req.StartDate != "" {
		// Tentar parse ISO8601 primeiro
		parsed, err := time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			// Fallback para YYYY-MM-DD
			parsed, err = time.Parse("2006-01-02", req.StartDate)
			if err != nil {
				return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error:   "validation_error",
					Message: "Data inicial inválida (formato: YYYY-MM-DD ou ISO8601)",
				})
			}
		}
		// Normalizar para início do dia UTC
		startDate = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
	}
	if req.EndDate != "" {
		// Tentar parse ISO8601 primeiro
		parsed, err := time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			// Fallback para YYYY-MM-DD
			parsed, err = time.Parse("2006-01-02", req.EndDate)
			if err != nil {
				return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error:   "validation_error",
					Message: "Data final inválida (formato: YYYY-MM-DD ou ISO8601)",
				})
			}
		}
		// Normalizar para fim do dia UTC
		endDate = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 999999999, time.UTC)
	}

	// Parse status array
	var statuses []valueobject.AppointmentStatus
	if len(req.Status) > 0 {
		for _, s := range req.Status {
			parsed, valid := valueobject.ParseAppointmentStatus(s)
			if !valid {
				return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error:   "validation_error",
					Message: fmt.Sprintf("Status inválido: %s", s),
				})
			}
			statuses = append(statuses, parsed)
		}
	}

	// RBAC: Barbeiro só vê seus próprios agendamentos
	professionalID := req.ProfessionalID
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		// Se barbeiro tentar filtrar por outro profissional, negar acesso
		if professionalID != "" && professionalID != barberProfID {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "Acesso negado: você só pode ver seus próprios agendamentos",
			})
		}
		// Forçar filtro para o profissional atual
		professionalID = barberProfID
	}

	input := appointment.ListAppointmentsInput{
		TenantID:       tenantID,
		ProfessionalID: professionalID,
		CustomerID:     req.CustomerID,
		Statuses:       statuses,
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

	if err := h.enforceBarberScope(ctx, c, tenantID, appointmentID); err != nil {
		return err
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

	// RBAC: Barbeiro só pode ver seus próprios agendamentos
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		if result.ProfessionalID != barberProfID {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "Acesso negado: você só pode ver seus próprios agendamentos",
			})
		}
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

	if err := h.enforceBarberScope(ctx, c, tenantID, appointmentID); err != nil {
		return err
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
		switch {
		case errors.Is(err, domain.ErrAppointmentNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: "Agendamento não encontrado"})
		case errors.Is(err, domain.ErrAppointmentInvalidStatusTransition):
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "invalid_transition", Message: err.Error()})
		default:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "update_error", Message: err.Error()})
		}
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

	if err := h.enforceBarberScope(ctx, c, tenantID, appointmentID); err != nil {
		return err
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

	// RBAC: barbeiro não pode reagendar para outro profissional
	if middleware.IsBarber(c) {
		barberProfID := middleware.GetProfessionalIDForBarber(c)
		if req.ProfessionalID != "" && req.ProfessionalID != barberProfID {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "forbidden", Message: "Barbeiro não pode mover agendamento para outro profissional"})
		}
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

		switch {
		case errors.Is(err, domain.ErrAppointmentNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: "Agendamento não encontrado"})
		case errors.Is(err, domain.ErrAppointmentConflict),
			errors.Is(err, domain.ErrAppointmentBlockedTimeConflict),
			errors.Is(err, domain.ErrAppointmentMinimumInterval):
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "conflict", Message: err.Error()})
		case errors.Is(err, domain.ErrAppointmentInvalidStatusTransition),
			errors.Is(err, domain.ErrAppointmentCannotReschedule):
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "invalid_transition", Message: err.Error()})
		case errors.Is(err, domain.ErrAppointmentProfessionalNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: err.Error()})
		default:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "reschedule_error", Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, mapper.AppointmentToResponse(result))
}

// ConfirmAppointment godoc
// @Summary Confirmar agendamento
// @Description Confirma um agendamento (CREATED -> CONFIRMED)
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do agendamento"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/appointments/{id}/confirm [post]
// @Security BearerAuth
func (h *AppointmentHandler) ConfirmAppointment(c echo.Context) error {
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

	if err := h.enforceBarberScope(ctx, c, tenantID, appointmentID); err != nil {
		return err
	}

	input := appointment.UpdateAppointmentStatusInput{
		TenantID:      tenantID,
		AppointmentID: appointmentID,
		NewStatus:     valueobject.AppointmentStatusConfirmed,
		Reason:        "Agendamento confirmado",
	}

	result, err := h.updateStatusUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao confirmar agendamento", zap.Error(err))
		switch {
		case errors.Is(err, domain.ErrAppointmentNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: "Agendamento não encontrado"})
		case errors.Is(err, domain.ErrAppointmentInvalidStatusTransition):
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "invalid_transition", Message: err.Error()})
		default:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "confirm_error", Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, mapper.AppointmentToResponse(result))
}

// CancelAppointment godoc
// @Summary Cancelar agendamento
// @Description Cancela um agendamento existente
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

	if err := h.enforceBarberScope(ctx, c, tenantID, appointmentID); err != nil {
		return err
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
		switch {
		case errors.Is(err, domain.ErrAppointmentNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: "Agendamento não encontrado"})
		case errors.Is(err, domain.ErrAppointmentInvalidStatusTransition):
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "invalid_transition", Message: err.Error()})
		default:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "cancel_error", Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, mapper.AppointmentToResponse(result))
}

// CheckInAppointment godoc
// @Summary Marcar cliente como chegou
// @Description Marca que o cliente chegou na barbearia (CONFIRMED/CREATED -> CHECKED_IN)
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do agendamento"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/appointments/{id}/check-in [post]
// @Security BearerAuth
func (h *AppointmentHandler) CheckInAppointment(c echo.Context) error {
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

	if err := h.enforceBarberScope(ctx, c, tenantID, appointmentID); err != nil {
		return err
	}

	input := appointment.UpdateAppointmentStatusInput{
		TenantID:      tenantID,
		AppointmentID: appointmentID,
		NewStatus:     valueobject.AppointmentStatusCheckedIn,
		Reason:        "Cliente chegou",
	}

	result, err := h.updateStatusUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao fazer check-in", zap.Error(err))
		switch {
		case errors.Is(err, domain.ErrAppointmentNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: "Agendamento não encontrado"})
		case errors.Is(err, domain.ErrAppointmentInvalidStatusTransition):
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "invalid_transition", Message: err.Error()})
		default:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "checkin_error", Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, mapper.AppointmentToResponse(result))
}

// StartServiceAppointment godoc
// @Summary Iniciar atendimento
// @Description Inicia o atendimento do cliente (CHECKED_IN/CONFIRMED -> IN_SERVICE)
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do agendamento"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/appointments/{id}/start [post]
// @Security BearerAuth
func (h *AppointmentHandler) StartServiceAppointment(c echo.Context) error {
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

	if err := h.enforceBarberScope(ctx, c, tenantID, appointmentID); err != nil {
		return err
	}

	input := appointment.UpdateAppointmentStatusInput{
		TenantID:      tenantID,
		AppointmentID: appointmentID,
		NewStatus:     valueobject.AppointmentStatusInService,
		Reason:        "Atendimento iniciado",
	}

	result, err := h.updateStatusUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao iniciar atendimento", zap.Error(err))
		switch {
		case errors.Is(err, domain.ErrAppointmentNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: "Agendamento não encontrado"})
		case errors.Is(err, domain.ErrAppointmentInvalidStatusTransition):
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "invalid_transition", Message: err.Error()})
		default:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "start_service_error", Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, mapper.AppointmentToResponse(result))
}

// FinishServiceAppointment godoc
// @Summary Finalizar atendimento
// @Description Finaliza o atendimento do cliente (IN_SERVICE -> AWAITING_PAYMENT) e cria comanda automaticamente
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do agendamento"
// @Success 200 {object} dto.FinishServiceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/appointments/{id}/finish [post]
// @Security BearerAuth
func (h *AppointmentHandler) FinishServiceAppointment(c echo.Context) error {
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

	if err := h.enforceBarberScope(ctx, c, tenantID, appointmentID); err != nil {
		return err
	}

	// Usar o novo use case que cria comanda automaticamente
	if h.finishWithCommandUC != nil {
		input := appointment.FinishServiceWithCommandInput{
			TenantID:      tenantID,
			AppointmentID: appointmentID,
		}

		result, err := h.finishWithCommandUC.Execute(ctx, input)
		if err != nil {
			h.logger.Error("Erro ao finalizar atendimento com comanda", zap.Error(err))
			switch {
			case errors.Is(err, domain.ErrAppointmentNotFound):
				return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: "Agendamento não encontrado"})
			case errors.Is(err, domain.ErrAppointmentInvalidStatusTransition):
				return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "invalid_transition", Message: err.Error()})
			default:
				return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "finish_service_error", Message: err.Error()})
			}
		}

		// Retornar resposta com informação da comanda
		response := mapper.AppointmentToResponse(result.Appointment)
		if result.Command != nil {
			response.CommandID = result.Command.ID.String()
		}

		return c.JSON(http.StatusOK, response)
	}

	// Fallback: usar o use case antigo se finishWithCommandUC não estiver configurado
	input := appointment.UpdateAppointmentStatusInput{
		TenantID:      tenantID,
		AppointmentID: appointmentID,
		NewStatus:     valueobject.AppointmentStatusAwaitingPayment,
		Reason:        "Atendimento finalizado",
	}

	result, err := h.updateStatusUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao finalizar atendimento", zap.Error(err))
		switch {
		case errors.Is(err, domain.ErrAppointmentNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: "Agendamento não encontrado"})
		case errors.Is(err, domain.ErrAppointmentInvalidStatusTransition):
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "invalid_transition", Message: err.Error()})
		default:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "finish_service_error", Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, mapper.AppointmentToResponse(result))
}

// CompleteAppointment godoc
// @Summary Concluir agendamento
// @Description Marca o agendamento como concluído (AWAITING_PAYMENT/IN_SERVICE -> DONE)
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do agendamento"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/appointments/{id}/complete [post]
// @Security BearerAuth
func (h *AppointmentHandler) CompleteAppointment(c echo.Context) error {
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

	if err := h.enforceBarberScope(ctx, c, tenantID, appointmentID); err != nil {
		return err
	}

	input := appointment.UpdateAppointmentStatusInput{
		TenantID:      tenantID,
		AppointmentID: appointmentID,
		NewStatus:     valueobject.AppointmentStatusDone,
		Reason:        "Pagamento recebido",
	}

	result, err := h.updateStatusUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao concluir agendamento", zap.Error(err))
		switch {
		case errors.Is(err, domain.ErrAppointmentNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: "Agendamento não encontrado"})
		case errors.Is(err, domain.ErrAppointmentInvalidStatusTransition):
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "invalid_transition", Message: err.Error()})
		default:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "complete_error", Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, mapper.AppointmentToResponse(result))
}

// NoShowAppointment godoc
// @Summary Marcar cliente como faltou
// @Description Marca que o cliente não compareceu (-> NO_SHOW)
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do agendamento"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/appointments/{id}/no-show [post]
// @Security BearerAuth
func (h *AppointmentHandler) NoShowAppointment(c echo.Context) error {
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

	if err := h.enforceBarberScope(ctx, c, tenantID, appointmentID); err != nil {
		return err
	}

	input := appointment.UpdateAppointmentStatusInput{
		TenantID:      tenantID,
		AppointmentID: appointmentID,
		NewStatus:     valueobject.AppointmentStatusNoShow,
		Reason:        "Cliente não compareceu",
	}

	result, err := h.updateStatusUC.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Erro ao marcar no-show", zap.Error(err))
		switch {
		case errors.Is(err, domain.ErrAppointmentNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: "Agendamento não encontrado"})
		case errors.Is(err, domain.ErrAppointmentInvalidStatusTransition):
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "invalid_transition", Message: err.Error()})
		default:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "no_show_error", Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, mapper.AppointmentToResponse(result))
}
