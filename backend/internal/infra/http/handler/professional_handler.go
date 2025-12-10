// Package handler contém o handler HTTP para profissionais
package handler

import (
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/infra/repository/postgres"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ProfessionalHandler manipula requisições HTTP de profissionais
type ProfessionalHandler struct {
	repo   *postgres.ProfessionalRepository
	logger *zap.Logger
}

// NewProfessionalHandler cria uma nova instância do handler
func NewProfessionalHandler(repo *postgres.ProfessionalRepository, logger *zap.Logger) *ProfessionalHandler {
	return &ProfessionalHandler{
		repo:   repo,
		logger: logger,
	}
}

// ListProfessionals godoc
// @Summary Listar profissionais
// @Description Lista profissionais com filtros e paginação
// @Tags Profissionais
// @Accept json
// @Produce json
// @Param search query string false "Busca por nome, email ou CPF"
// @Param status query string false "Filtrar por status"
// @Param tipo query string false "Filtrar por tipo"
// @Param order_by query string false "Ordenar por"
// @Param page query int false "Página"
// @Param page_size query int false "Tamanho da página"
// @Success 200 {object} dto.ListProfessionalsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/professionals [get]
// @Security BearerAuth
func (h *ProfessionalHandler) ListProfessionals(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.ListProfessionalsRequest
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

	professionals, total, err := h.repo.List(ctx, tenantID, req)
	if err != nil {
		h.logger.Error("Erro ao listar profissionais", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao listar profissionais",
		})
	}

	return c.JSON(http.StatusOK, dto.ListProfessionalsResponse{
		Data:     professionals,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	})
}

// GetProfessional godoc
// @Summary Buscar profissional
// @Description Busca um profissional pelo ID
// @Tags Profissionais
// @Accept json
// @Produce json
// @Param id path string true "ID do profissional"
// @Success 200 {object} dto.ProfessionalResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/professionals/{id} [get]
// @Security BearerAuth
func (h *ProfessionalHandler) GetProfessional(c echo.Context) error {
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

	professional, err := h.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar profissional", zap.Error(err))
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Profissional não encontrado",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": professional,
	})
}

// CreateProfessional godoc
// @Summary Criar profissional
// @Description Cria um novo profissional
// @Tags Profissionais
// @Accept json
// @Produce json
// @Param request body dto.CreateProfessionalRequest true "Dados do profissional"
// @Success 201 {object} dto.ProfessionalResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/professionals [post]
// @Security BearerAuth
func (h *ProfessionalHandler) CreateProfessional(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	var req dto.CreateProfessionalRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	if err := c.Validate(&req); err != nil {
		h.logger.Error("Erro de validação", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// Verificar se email já existe
	emailExists, err := h.repo.CheckEmailExists(ctx, tenantID, req.Email, nil)
	if err != nil {
		h.logger.Error("Erro ao verificar email", zap.Error(err))
	}
	if emailExists {
		return c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error:   "conflict",
			Message: "Este email já está cadastrado",
		})
	}

	// Verificar se CPF já existe
	if req.CPF != "" {
		cpfExists, err := h.repo.CheckCpfExists(ctx, tenantID, req.CPF, nil)
		if err != nil {
			h.logger.Error("Erro ao verificar CPF", zap.Error(err))
		}
		if cpfExists {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: "Este CPF já está cadastrado",
			})
		}
	}

	professional, err := h.repo.Create(ctx, tenantID, req)
	if err != nil {
		h.logger.Error("Erro ao criar profissional", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao criar profissional",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": professional,
	})
}

// UpdateProfessional godoc
// @Summary Atualizar profissional
// @Description Atualiza um profissional existente
// @Tags Profissionais
// @Accept json
// @Produce json
// @Param id path string true "ID do profissional"
// @Param request body dto.UpdateProfessionalRequest true "Dados do profissional"
// @Success 200 {object} dto.ProfessionalResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/professionals/{id} [put]
// @Security BearerAuth
func (h *ProfessionalHandler) UpdateProfessional(c echo.Context) error {
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

	var req dto.UpdateProfessionalRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	if err := c.Validate(&req); err != nil {
		h.logger.Error("Erro de validação", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// Verificar se email já existe (excluindo o próprio)
	emailExists, err := h.repo.CheckEmailExists(ctx, tenantID, req.Email, &id)
	if err != nil {
		h.logger.Error("Erro ao verificar email", zap.Error(err))
	}
	if emailExists {
		return c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error:   "conflict",
			Message: "Este email já está cadastrado",
		})
	}

	// Verificar se CPF já existe (se foi alterado)
	if req.CPF != "" {
		cpfExists, err := h.repo.CheckCpfExists(ctx, tenantID, req.CPF, &id)
		if err != nil {
			h.logger.Error("Erro ao verificar CPF", zap.Error(err))
		}
		if cpfExists {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: "Este CPF já está cadastrado",
			})
		}
	}

	professional, err := h.repo.Update(ctx, tenantID, id, req)
	if err != nil {
		h.logger.Error("Erro ao atualizar profissional", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao atualizar profissional",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": professional,
	})
}

// UpdateProfessionalStatus godoc
// @Summary Atualizar status do profissional
// @Description Atualiza o status de um profissional
// @Tags Profissionais
// @Accept json
// @Produce json
// @Param id path string true "ID do profissional"
// @Param request body dto.UpdateProfessionalStatusRequest true "Novo status"
// @Success 200 {object} dto.ProfessionalResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/professionals/{id}/status [put]
// @Security BearerAuth
func (h *ProfessionalHandler) UpdateProfessionalStatus(c echo.Context) error {
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

	var req dto.UpdateProfessionalStatusRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("Erro ao fazer bind", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Dados inválidos",
		})
	}

	if err := c.Validate(&req); err != nil {
		h.logger.Error("Erro de validação", zap.Error(err))
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	professional, err := h.repo.UpdateStatus(ctx, tenantID, id, req)
	if err != nil {
		h.logger.Error("Erro ao atualizar status", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao atualizar status",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": professional,
	})
}

// DeleteProfessional godoc
// @Summary Deletar profissional
// @Description Remove um profissional (soft delete alterando status)
// @Tags Profissionais
// @Accept json
// @Produce json
// @Param id path string true "ID do profissional"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/professionals/{id} [delete]
// @Security BearerAuth
func (h *ProfessionalHandler) DeleteProfessional(c echo.Context) error {
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

	// Deletar permanentemente o profissional
	err := h.repo.Delete(ctx, tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao deletar profissional permanentemente", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao deletar profissional",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// CheckEmailExists godoc
// @Summary Verificar email
// @Description Verifica se email já existe no tenant
// @Tags Profissionais
// @Accept json
// @Produce json
// @Param email query string true "Email para verificar"
// @Param exclude_id query string false "ID para excluir da verificação"
// @Success 200 {object} dto.CheckExistsProfessionalResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/professionals/check-email [get]
// @Security BearerAuth
func (h *ProfessionalHandler) CheckEmailExists(c echo.Context) error {
	ctx := c.Request().Context()

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "Tenant ID não encontrado",
		})
	}

	email := c.QueryParam("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Email é obrigatório",
		})
	}

	excludeID := c.QueryParam("exclude_id")
	var excludeIDPtr *string
	if excludeID != "" {
		excludeIDPtr = &excludeID
	}

	exists, err := h.repo.CheckEmailExists(ctx, tenantID, email, excludeIDPtr)
	if err != nil {
		h.logger.Error("Erro ao verificar email", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao verificar email",
		})
	}

	return c.JSON(http.StatusOK, dto.CheckExistsProfessionalResponse{
		Exists: exists,
	})
}

// CheckCpfExists godoc
// @Summary Verificar CPF
// @Description Verifica se CPF já existe no tenant
// @Tags Profissionais
// @Accept json
// @Produce json
// @Param cpf query string true "CPF para verificar"
// @Param exclude_id query string false "ID para excluir da verificação"
// @Success 200 {object} dto.CheckExistsProfessionalResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/professionals/check-cpf [get]
// @Security BearerAuth
func (h *ProfessionalHandler) CheckCpfExists(c echo.Context) error {
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

	excludeID := c.QueryParam("exclude_id")
	var excludeIDPtr *string
	if excludeID != "" {
		excludeIDPtr = &excludeID
	}

	exists, err := h.repo.CheckCpfExists(ctx, tenantID, cpf, excludeIDPtr)
	if err != nil {
		h.logger.Error("Erro ao verificar CPF", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Erro ao verificar CPF",
		})
	}

	return c.JSON(http.StatusOK, dto.CheckExistsProfessionalResponse{
		Exists: exists,
	})
}
