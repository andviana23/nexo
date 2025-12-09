package handler

import (
	"net/http"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/unit"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// UnitHandler gerencia os endpoints de unidades (filiais)
type UnitHandler struct {
	createUC        *unit.CreateUnitUseCase
	listUC          *unit.ListUnitsUseCase
	getUC           *unit.GetUnitUseCase
	updateUC        *unit.UpdateUnitUseCase
	deleteUC        *unit.DeleteUnitUseCase
	toggleUC        *unit.ToggleUnitUseCase
	linkUserUC      *unit.LinkUserToUnitUseCase
	unlinkUserUC    *unit.UnlinkUserFromUnitUseCase
	listUserUnitsUC *unit.ListUserUnitsUseCase
	setDefaultUC    *unit.SetDefaultUnitUseCase
	getDefaultUC    *unit.GetDefaultUnitUseCase
	checkAccessUC   *unit.CheckUserAccessToUnitUseCase
	listUnitUsersUC *unit.ListUnitUsersUseCase
	logger          *zap.Logger
}

// NewUnitHandler cria uma nova instância do handler
func NewUnitHandler(
	createUC *unit.CreateUnitUseCase,
	listUC *unit.ListUnitsUseCase,
	getUC *unit.GetUnitUseCase,
	updateUC *unit.UpdateUnitUseCase,
	deleteUC *unit.DeleteUnitUseCase,
	toggleUC *unit.ToggleUnitUseCase,
	linkUserUC *unit.LinkUserToUnitUseCase,
	unlinkUserUC *unit.UnlinkUserFromUnitUseCase,
	listUserUnitsUC *unit.ListUserUnitsUseCase,
	setDefaultUC *unit.SetDefaultUnitUseCase,
	getDefaultUC *unit.GetDefaultUnitUseCase,
	checkAccessUC *unit.CheckUserAccessToUnitUseCase,
	listUnitUsersUC *unit.ListUnitUsersUseCase,
	logger *zap.Logger,
) *UnitHandler {
	return &UnitHandler{
		createUC:        createUC,
		listUC:          listUC,
		getUC:           getUC,
		updateUC:        updateUC,
		deleteUC:        deleteUC,
		toggleUC:        toggleUC,
		linkUserUC:      linkUserUC,
		unlinkUserUC:    unlinkUserUC,
		listUserUnitsUC: listUserUnitsUC,
		setDefaultUC:    setDefaultUC,
		getDefaultUC:    getDefaultUC,
		checkAccessUC:   checkAccessUC,
		listUnitUsersUC: listUnitUsersUC,
		logger:          logger,
	}
}

// ============================================================================
// UNIT CRUD ENDPOINTS
// ============================================================================

// Create cria uma nova unidade
// @Summary Criar unidade
// @Description Cria uma nova unidade/filial para o tenant
// @Tags Unidades
// @Accept json
// @Produce json
// @Param request body dto.CreateUnitRequest true "Dados da unidade"
// @Success 201 {object} dto.UnitResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/units [post]
// @Security BearerAuth
func (h *UnitHandler) Create(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	var req dto.CreateUnitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Default timezone
	timezone := req.Timezone
	if timezone == "" {
		timezone = "America/Sao_Paulo"
	}

	input := unit.CreateUnitInput{
		TenantID:       tenantID,
		Nome:           req.Nome,
		Apelido:        strPtr(req.Apelido),
		Descricao:      strPtr(req.Descricao),
		EnderecoResumo: strPtr(req.EnderecoResumo),
		Cidade:         strPtr(req.Cidade),
		Estado:         strPtr(req.Estado),
		Timezone:       timezone,
		IsMatriz:       req.IsMatriz,
	}

	unitEntity, err := h.createUC.Execute(c.Request().Context(), input)
	if err != nil {
		h.logger.Error("Erro ao criar unidade", zap.Error(err))
		if err.Error() == "já existe uma unidade com este nome" {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, toUnitResponse(unitEntity))
}

// List lista todas as unidades do tenant
// @Summary Listar unidades
// @Description Lista todas as unidades/filiais do tenant
// @Tags Unidades
// @Produce json
// @Param ativas query bool false "Filtrar apenas ativas"
// @Success 200 {object} dto.ListUnitsResponse
// @Router /api/v1/units [get]
// @Security BearerAuth
func (h *UnitHandler) List(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	apenasAtivas := c.QueryParam("ativas") == "true"

	units, err := h.listUC.Execute(c.Request().Context(), tenantID, apenasAtivas)
	if err != nil {
		h.logger.Error("Erro ao listar unidades", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]dto.UnitResponse, len(units))
	for i := range units {
		responses[i] = toUnitResponse(&units[i])
	}

	return c.JSON(http.StatusOK, dto.ListUnitsResponse{
		Units: responses,
		Total: len(responses),
	})
}

// GetByID busca uma unidade por ID
// @Summary Buscar unidade por ID
// @Description Retorna os detalhes de uma unidade específica
// @Tags Unidades
// @Produce json
// @Param id path string true "ID da unidade"
// @Success 200 {object} dto.UnitResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/units/{id} [get]
// @Security BearerAuth
func (h *UnitHandler) GetByID(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	unitEntity, err := h.getUC.Execute(c.Request().Context(), id, tenantID)
	if err != nil {
		h.logger.Error("Erro ao buscar unidade", zap.Error(err))
		if err.Error() == "unidade não encontrada" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toUnitResponse(unitEntity))
}

// Update atualiza uma unidade
// @Summary Atualizar unidade
// @Description Atualiza os dados de uma unidade existente
// @Tags Unidades
// @Accept json
// @Produce json
// @Param id path string true "ID da unidade"
// @Param request body dto.UpdateUnitRequest true "Dados atualizados"
// @Success 200 {object} dto.UnitResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/units/{id} [put]
// @Security BearerAuth
func (h *UnitHandler) Update(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	var req dto.UpdateUnitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	input := unit.UpdateUnitInput{
		ID:             id,
		TenantID:       tenantID,
		Nome:           strPtr(req.Nome),
		Apelido:        strPtr(req.Apelido),
		Descricao:      strPtr(req.Descricao),
		EnderecoResumo: strPtr(req.EnderecoResumo),
		Cidade:         strPtr(req.Cidade),
		Estado:         strPtr(req.Estado),
		Timezone:       strPtr(req.Timezone),
	}

	unitEntity, err := h.updateUC.Execute(c.Request().Context(), input)
	if err != nil {
		h.logger.Error("Erro ao atualizar unidade", zap.Error(err))
		switch err.Error() {
		case "unidade não encontrada":
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		case "já existe uma unidade com este nome":
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toUnitResponse(unitEntity))
}

// Delete remove uma unidade
// @Summary Remover unidade
// @Description Remove uma unidade do tenant (soft delete)
// @Tags Unidades
// @Param id path string true "ID da unidade"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/units/{id} [delete]
// @Security BearerAuth
func (h *UnitHandler) Delete(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	if err := h.deleteUC.Execute(c.Request().Context(), id, tenantID); err != nil {
		h.logger.Error("Erro ao deletar unidade", zap.Error(err))
		if err.Error() == "unidade não encontrada" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		if err.Error() == "não é possível excluir a última unidade ativa" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Toggle ativa/desativa uma unidade
// @Summary Alternar status da unidade
// @Description Ativa ou desativa uma unidade
// @Tags Unidades
// @Param id path string true "ID da unidade"
// @Success 200 {object} dto.UnitResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/units/{id}/toggle [patch]
// @Security BearerAuth
func (h *UnitHandler) Toggle(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	unitEntity, err := h.toggleUC.Execute(c.Request().Context(), id, tenantID)
	if err != nil {
		h.logger.Error("Erro ao alternar status da unidade", zap.Error(err))
		switch err.Error() {
		case "unidade não encontrada":
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		case "não é possível desativar a última unidade ativa":
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toUnitResponse(unitEntity))
}

// ============================================================================
// USER UNIT ENDPOINTS (para /me/units)
// ============================================================================

// ListMyUnits lista as unidades do usuário logado
// @Summary Listar minhas unidades
// @Description Lista todas as unidades às quais o usuário tem acesso
// @Tags Unidades
// @Produce json
// @Success 200 {object} dto.ListUserUnitsResponse
// @Router /api/v1/me/units [get]
// @Security BearerAuth
func (h *UnitHandler) ListMyUnits(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return err
	}

	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	userUnits, err := h.listUserUnitsUC.Execute(c.Request().Context(), userID)
	if err != nil {
		h.logger.Error("Erro ao listar unidades do usuário", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]dto.UserUnitResponse, len(userUnits))
	for i, uu := range userUnits {
		// Buscar detalhes da unidade
		unitEntity, err := h.getUC.Execute(c.Request().Context(), uu.UnitID, tenantID)
		if err != nil {
			continue
		}
		responses[i] = toUserUnitResponse(&uu, unitEntity)
	}

	return c.JSON(http.StatusOK, dto.ListUserUnitsResponse{
		Units: responses,
		Total: len(responses),
	})
}

// SwitchUnit troca a unidade ativa do usuário
// @Summary Trocar de unidade
// @Description Define uma nova unidade ativa para o usuário
// @Tags Unidades
// @Accept json
// @Produce json
// @Param request body dto.SwitchUnitRequest true "ID da unidade"
// @Success 200 {object} dto.SwitchUnitResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/v1/me/switch-unit [post]
// @Security BearerAuth
func (h *UnitHandler) SwitchUnit(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return err
	}

	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	var req dto.SwitchUnitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	unitID, err := uuid.Parse(req.UnitID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unit_id inválido"})
	}

	// Verificar se usuário tem acesso à unidade
	hasAccess, err := h.checkAccessUC.Execute(c.Request().Context(), userID, unitID)
	if err != nil {
		h.logger.Error("Erro ao verificar acesso à unidade", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !hasAccess {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "usuário não tem acesso a esta unidade"})
	}

	// Definir como unidade padrão
	if err := h.setDefaultUC.Execute(c.Request().Context(), userID, unitID); err != nil {
		h.logger.Error("Erro ao definir unidade padrão", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Buscar detalhes da unidade e user_unit
	unitEntity, err := h.getUC.Execute(c.Request().Context(), unitID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	userUnits, err := h.listUserUnitsUC.Execute(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var userUnit *dto.UserUnitResponse
	for _, uu := range userUnits {
		if uu.UnitID == unitID {
			resp := toUserUnitResponse(&uu, unitEntity)
			userUnit = &resp
			break
		}
	}

	if userUnit == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "vínculo não encontrado"})
	}

	// TODO: Gerar novo token JWT com unit_id
	// Por ora, retornamos sem o novo token (será implementado junto com auth)
	return c.JSON(http.StatusOK, dto.SwitchUnitResponse{
		Unit:        *userUnit,
		AccessToken: "", // Será preenchido quando integrar com auth
	})
}

// SetDefaultUnit define a unidade padrão do usuário
// @Summary Definir unidade padrão
// @Description Define qual unidade será a padrão para o usuário
// @Tags Unidades
// @Accept json
// @Produce json
// @Param request body dto.SetDefaultUnitRequest true "ID da unidade"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/v1/me/default-unit [put]
// @Security BearerAuth
func (h *UnitHandler) SetDefaultUnit(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return err
	}

	var req dto.SetDefaultUnitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	unitID, err := uuid.Parse(req.UnitID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unit_id inválido"})
	}

	// Verificar acesso
	hasAccess, err := h.checkAccessUC.Execute(c.Request().Context(), userID, unitID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !hasAccess {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "usuário não tem acesso a esta unidade"})
	}

	if err := h.setDefaultUC.Execute(c.Request().Context(), userID, unitID); err != nil {
		h.logger.Error("Erro ao definir unidade padrão", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// ============================================================================
// ADMIN - Gerenciamento de usuários em unidades
// ============================================================================

// AddUserToUnit adiciona um usuário a uma unidade
// @Summary Adicionar usuário à unidade
// @Description Vincula um usuário a uma unidade (admin only)
// @Tags Unidades
// @Accept json
// @Produce json
// @Param id path string true "ID da unidade"
// @Param request body dto.AddUserToUnitRequest true "Dados do vínculo"
// @Success 201 {object} dto.UserUnitResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/units/{id}/users [post]
// @Security BearerAuth
func (h *UnitHandler) AddUserToUnit(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	unitIDStr := c.Param("id")
	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID da unidade inválido"})
	}

	var req dto.AddUserToUnitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id inválido"})
	}

	input := unit.LinkUserToUnitInput{
		UserID:       userID,
		UnitID:       unitID,
		IsDefault:    req.IsDefault,
		RoleOverride: req.RoleOverride,
	}

	userUnit, err := h.linkUserUC.Execute(c.Request().Context(), input)
	if err != nil {
		h.logger.Error("Erro ao adicionar usuário à unidade", zap.Error(err))
		if err.Error() == "usuário já está vinculado a esta unidade" {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Buscar detalhes da unidade
	unitEntity, err := h.getUC.Execute(c.Request().Context(), unitID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	detailedUserUnit := &entity.UserUnitWithDetails{UserUnit: *userUnit}

	return c.JSON(http.StatusCreated, toUserUnitResponse(detailedUserUnit, unitEntity))
}

// RemoveUserFromUnit remove um usuário de uma unidade
// @Summary Remover usuário da unidade
// @Description Remove o vínculo de um usuário com uma unidade
// @Tags Unidades
// @Param id path string true "ID da unidade"
// @Param userId path string true "ID do usuário"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/units/{id}/users/{userId} [delete]
// @Security BearerAuth
func (h *UnitHandler) RemoveUserFromUnit(c echo.Context) error {
	unitIDStr := c.Param("id")
	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID da unidade inválido"})
	}

	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID do usuário inválido"})
	}

	if err := h.unlinkUserUC.Execute(c.Request().Context(), userID, unitID); err != nil {
		h.logger.Error("Erro ao remover usuário da unidade", zap.Error(err))
		if err.Error() == "vínculo não encontrado" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// ListUnitUsers lista os usuários de uma unidade
// @Summary Listar usuários da unidade
// @Description Lista todos os usuários vinculados a uma unidade
// @Tags Unidades
// @Produce json
// @Param id path string true "ID da unidade"
// @Success 200 {object} dto.ListUserUnitsResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/units/{id}/users [get]
// @Security BearerAuth
func (h *UnitHandler) ListUnitUsers(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	unitIDStr := c.Param("id")
	unitID, err := uuid.Parse(unitIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID da unidade inválido"})
	}

	userUnits, err := h.listUnitUsersUC.Execute(c.Request().Context(), unitID)
	if err != nil {
		h.logger.Error("Erro ao listar usuários da unidade", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Buscar detalhes da unidade
	unitEntity, err := h.getUC.Execute(c.Request().Context(), unitID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]dto.UserUnitResponse, len(userUnits))
	for i, uu := range userUnits {
		responses[i] = toUserUnitResponse(&uu, unitEntity)
	}

	return c.JSON(http.StatusOK, dto.ListUserUnitsResponse{
		Units: responses,
		Total: len(responses),
	})
}

// ============================================================================
// HELPERS
// ============================================================================

func (h *UnitHandler) getTenantID(c echo.Context) (uuid.UUID, error) {
	tenantIDStr, ok := c.Get("tenant_id").(string)
	if !ok || tenantIDStr == "" {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant não identificado"})
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "tenant não identificado")
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "tenant_id inválido"})
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "tenant_id inválido")
	}
	return tenantID, nil
}

func (h *UnitHandler) getUserID(c echo.Context) (uuid.UUID, error) {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "usuário não identificado"})
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "usuário não identificado")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "user_id inválido"})
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "user_id inválido")
	}
	return userID, nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func toUnitResponse(u *unit.UnitOutput) dto.UnitResponse {
	resp := dto.UnitResponse{
		ID:           u.ID.String(),
		TenantID:     u.TenantID.String(),
		Nome:         u.Nome,
		Timezone:     u.Timezone,
		Ativa:        u.Ativa,
		IsMatriz:     u.IsMatriz,
		CriadoEm:     u.CriadoEm.Format("2006-01-02T15:04:05Z07:00"),
		AtualizadoEm: u.AtualizadoEm.Format("2006-01-02T15:04:05Z07:00"),
	}

	if u.Apelido != nil {
		resp.Apelido = u.Apelido
	}
	if u.Descricao != nil {
		resp.Descricao = u.Descricao
	}
	if u.EnderecoResumo != nil {
		resp.EnderecoResumo = u.EnderecoResumo
	}
	if u.Cidade != nil {
		resp.Cidade = u.Cidade
	}
	if u.Estado != nil {
		resp.Estado = u.Estado
	}

	return resp
}

func toUserUnitResponse(uu *entity.UserUnitWithDetails, u *unit.UnitOutput) dto.UserUnitResponse {
	resp := dto.UserUnitResponse{
		ID:         uu.ID.String(),
		UserID:     uu.UserID.String(),
		UnitID:     uu.UnitID.String(),
		UnitNome:   u.Nome,
		UnitMatriz: u.IsMatriz,
		UnitAtiva:  u.Ativa,
		IsDefault:  uu.IsDefault,
		TenantID:   u.TenantID.String(),
	}

	if u.Apelido != nil {
		resp.UnitApelido = u.Apelido
	}
	if uu.RoleOverride != nil {
		resp.RoleOverride = uu.RoleOverride
	}

	return resp
}
