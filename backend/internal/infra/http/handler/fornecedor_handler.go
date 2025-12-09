package handler

import (
	"net/http"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// =============================================================================
// DTOs
// =============================================================================

// CreateFornecedorRequest request para criar fornecedor
type CreateFornecedorRequest struct {
	RazaoSocial  string `json:"razao_social" validate:"required"`
	NomeFantasia string `json:"nome_fantasia,omitempty"`
	CNPJ         string `json:"cnpj,omitempty"`
	Email        string `json:"email,omitempty" validate:"omitempty,email"`
	Telefone     string `json:"telefone" validate:"required"`
	Celular      string `json:"celular,omitempty"`
	// Endereço
	EnderecoLogradouro  string `json:"endereco_logradouro,omitempty"`
	EnderecoNumero      string `json:"endereco_numero,omitempty"`
	EnderecoComplemento string `json:"endereco_complemento,omitempty"`
	EnderecoBairro      string `json:"endereco_bairro,omitempty"`
	EnderecoCidade      string `json:"endereco_cidade,omitempty"`
	EnderecoEstado      string `json:"endereco_estado,omitempty"`
	EnderecoCEP         string `json:"endereco_cep,omitempty"`
	// Banco
	Banco       string `json:"banco,omitempty"`
	Agencia     string `json:"agencia,omitempty"`
	Conta       string `json:"conta,omitempty"`
	Observacoes string `json:"observacoes,omitempty"`
}

// UpdateFornecedorRequest request para atualizar fornecedor
type UpdateFornecedorRequest struct {
	RazaoSocial  string `json:"razao_social,omitempty"`
	NomeFantasia string `json:"nome_fantasia,omitempty"`
	CNPJ         string `json:"cnpj,omitempty"`
	Email        string `json:"email,omitempty" validate:"omitempty,email"`
	Telefone     string `json:"telefone,omitempty"`
	Celular      string `json:"celular,omitempty"`
	// Endereço
	EnderecoLogradouro  string `json:"endereco_logradouro,omitempty"`
	EnderecoNumero      string `json:"endereco_numero,omitempty"`
	EnderecoComplemento string `json:"endereco_complemento,omitempty"`
	EnderecoBairro      string `json:"endereco_bairro,omitempty"`
	EnderecoCidade      string `json:"endereco_cidade,omitempty"`
	EnderecoEstado      string `json:"endereco_estado,omitempty"`
	EnderecoCEP         string `json:"endereco_cep,omitempty"`
	// Banco
	Banco       string `json:"banco,omitempty"`
	Agencia     string `json:"agencia,omitempty"`
	Conta       string `json:"conta,omitempty"`
	Observacoes string `json:"observacoes,omitempty"`
	Ativo       *bool  `json:"ativo,omitempty"`
}

// FornecedorResponse response de fornecedor
type FornecedorResponse struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenant_id"`
	RazaoSocial  string `json:"razao_social"`
	NomeFantasia string `json:"nome_fantasia,omitempty"`
	Nome         string `json:"nome"` // Alias para compatibilidade com frontend
	CNPJ         string `json:"cnpj,omitempty"`
	Email        string `json:"email,omitempty"`
	Telefone     string `json:"telefone"`
	Celular      string `json:"celular,omitempty"`
	// Endereço
	EnderecoLogradouro  string `json:"endereco_logradouro,omitempty"`
	EnderecoNumero      string `json:"endereco_numero,omitempty"`
	EnderecoComplemento string `json:"endereco_complemento,omitempty"`
	EnderecoBairro      string `json:"endereco_bairro,omitempty"`
	EnderecoCidade      string `json:"endereco_cidade,omitempty"`
	EnderecoEstado      string `json:"endereco_estado,omitempty"`
	EnderecoCEP         string `json:"endereco_cep,omitempty"`
	Endereco            string `json:"endereco,omitempty"` // Alias formatado
	// Banco
	Banco       string `json:"banco,omitempty"`
	Agencia     string `json:"agencia,omitempty"`
	Conta       string `json:"conta,omitempty"`
	Observacoes string `json:"observacoes,omitempty"`
	Ativo       bool   `json:"ativo"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListFornecedoresResponse response de lista de fornecedores
type ListFornecedoresResponse struct {
	Fornecedores []FornecedorResponse `json:"fornecedores"`
	Total        int                  `json:"total"`
}

// =============================================================================
// HANDLER
// =============================================================================

// FornecedorHandler gerencia endpoints de fornecedores
type FornecedorHandler struct {
	repo   port.FornecedorRepository
	logger *zap.Logger
}

// NewFornecedorHandler cria nova instância do handler
func NewFornecedorHandler(repo port.FornecedorRepository, logger *zap.Logger) *FornecedorHandler {
	return &FornecedorHandler{
		repo:   repo,
		logger: logger,
	}
}

// RegisterRoutes registra as rotas do handler
func (h *FornecedorHandler) RegisterRoutes(g *echo.Group) {
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.PATCH("/:id/ativar", h.Ativar)
	g.PATCH("/:id/desativar", h.Desativar)
}

// =============================================================================
// ENDPOINTS
// =============================================================================

// Create cria um novo fornecedor
// @Summary Criar fornecedor
// @Description Cria um novo fornecedor
// @Tags Fornecedores
// @Accept json
// @Produce json
// @Param request body CreateFornecedorRequest true "Dados do fornecedor"
// @Success 201 {object} FornecedorResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/fornecedores [post]
// @Security BearerAuth
func (h *FornecedorHandler) Create(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	var req CreateFornecedorRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dados inválidos"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Criar entidade
	fornecedor, err := entity.NewFornecedor(tenantID, req.RazaoSocial, req.Telefone)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Preencher campos opcionais
	fornecedor.NomeFantasia = req.NomeFantasia
	fornecedor.CNPJ = req.CNPJ
	fornecedor.Email = req.Email
	fornecedor.Celular = req.Celular
	fornecedor.EnderecoLogradouro = req.EnderecoLogradouro
	fornecedor.EnderecoNumero = req.EnderecoNumero
	fornecedor.EnderecoComplemento = req.EnderecoComplemento
	fornecedor.EnderecoBairro = req.EnderecoBairro
	fornecedor.EnderecoCidade = req.EnderecoCidade
	fornecedor.EnderecoEstado = req.EnderecoEstado
	fornecedor.EnderecoCEP = req.EnderecoCEP
	fornecedor.Banco = req.Banco
	fornecedor.Agencia = req.Agencia
	fornecedor.Conta = req.Conta
	fornecedor.Observacoes = req.Observacoes

	if err := h.repo.Create(c.Request().Context(), fornecedor); err != nil {
		h.logger.Error("Erro ao criar fornecedor", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro ao criar fornecedor"})
	}

	return c.JSON(http.StatusCreated, toFornecedorResponse(fornecedor))
}

// List lista todos os fornecedores
// @Summary Listar fornecedores
// @Description Lista todos os fornecedores do tenant
// @Tags Fornecedores
// @Produce json
// @Param ativos query bool false "Filtrar apenas ativos"
// @Success 200 {object} ListFornecedoresResponse
// @Router /api/v1/fornecedores [get]
// @Security BearerAuth
func (h *FornecedorHandler) List(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	apenasAtivos := c.QueryParam("ativos") == "true"

	var fornecedores []*entity.Fornecedor
	if apenasAtivos {
		fornecedores, err = h.repo.ListAtivos(c.Request().Context(), tenantID)
	} else {
		fornecedores, err = h.repo.ListAll(c.Request().Context(), tenantID)
	}

	if err != nil {
		h.logger.Error("Erro ao listar fornecedores", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro ao listar fornecedores"})
	}

	responses := make([]FornecedorResponse, len(fornecedores))
	for i, f := range fornecedores {
		responses[i] = toFornecedorResponse(f)
	}

	return c.JSON(http.StatusOK, ListFornecedoresResponse{
		Fornecedores: responses,
		Total:        len(responses),
	})
}

// GetByID busca fornecedor por ID
// @Summary Buscar fornecedor
// @Description Busca um fornecedor por ID
// @Tags Fornecedores
// @Produce json
// @Param id path string true "ID do fornecedor"
// @Success 200 {object} FornecedorResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/fornecedores/{id} [get]
// @Security BearerAuth
func (h *FornecedorHandler) GetByID(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	fornecedor, err := h.repo.FindByID(c.Request().Context(), tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar fornecedor", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro ao buscar fornecedor"})
	}
	if fornecedor == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "fornecedor não encontrado"})
	}

	return c.JSON(http.StatusOK, toFornecedorResponse(fornecedor))
}

// Update atualiza um fornecedor
// @Summary Atualizar fornecedor
// @Description Atualiza um fornecedor existente
// @Tags Fornecedores
// @Accept json
// @Produce json
// @Param id path string true "ID do fornecedor"
// @Param request body UpdateFornecedorRequest true "Dados atualizados"
// @Success 200 {object} FornecedorResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/fornecedores/{id} [put]
// @Security BearerAuth
func (h *FornecedorHandler) Update(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	var req UpdateFornecedorRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dados inválidos"})
	}

	// Buscar fornecedor existente
	fornecedor, err := h.repo.FindByID(c.Request().Context(), tenantID, id)
	if err != nil {
		h.logger.Error("Erro ao buscar fornecedor", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro ao buscar fornecedor"})
	}
	if fornecedor == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "fornecedor não encontrado"})
	}

	// Atualizar campos
	if req.RazaoSocial != "" {
		fornecedor.RazaoSocial = req.RazaoSocial
	}
	if req.NomeFantasia != "" {
		fornecedor.NomeFantasia = req.NomeFantasia
	}
	if req.CNPJ != "" {
		fornecedor.CNPJ = req.CNPJ
	}
	if req.Email != "" {
		fornecedor.Email = req.Email
	}
	if req.Telefone != "" {
		fornecedor.Telefone = req.Telefone
	}
	if req.Celular != "" {
		fornecedor.Celular = req.Celular
	}
	if req.EnderecoLogradouro != "" {
		fornecedor.EnderecoLogradouro = req.EnderecoLogradouro
	}
	if req.EnderecoNumero != "" {
		fornecedor.EnderecoNumero = req.EnderecoNumero
	}
	if req.EnderecoComplemento != "" {
		fornecedor.EnderecoComplemento = req.EnderecoComplemento
	}
	if req.EnderecoBairro != "" {
		fornecedor.EnderecoBairro = req.EnderecoBairro
	}
	if req.EnderecoCidade != "" {
		fornecedor.EnderecoCidade = req.EnderecoCidade
	}
	if req.EnderecoEstado != "" {
		fornecedor.EnderecoEstado = req.EnderecoEstado
	}
	if req.EnderecoCEP != "" {
		fornecedor.EnderecoCEP = req.EnderecoCEP
	}
	if req.Banco != "" {
		fornecedor.Banco = req.Banco
	}
	if req.Agencia != "" {
		fornecedor.Agencia = req.Agencia
	}
	if req.Conta != "" {
		fornecedor.Conta = req.Conta
	}
	if req.Observacoes != "" {
		fornecedor.Observacoes = req.Observacoes
	}
	if req.Ativo != nil {
		fornecedor.Ativo = *req.Ativo
	}

	fornecedor.UpdatedAt = time.Now()

	if err := h.repo.Update(c.Request().Context(), fornecedor); err != nil {
		h.logger.Error("Erro ao atualizar fornecedor", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro ao atualizar fornecedor"})
	}

	return c.JSON(http.StatusOK, toFornecedorResponse(fornecedor))
}

// Delete remove um fornecedor
// @Summary Excluir fornecedor
// @Description Remove um fornecedor (soft delete)
// @Tags Fornecedores
// @Param id path string true "ID do fornecedor"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /api/v1/fornecedores/{id} [delete]
// @Security BearerAuth
func (h *FornecedorHandler) Delete(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	if err := h.repo.Delete(c.Request().Context(), tenantID, id); err != nil {
		h.logger.Error("Erro ao excluir fornecedor", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro ao excluir fornecedor"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Ativar reativa um fornecedor
func (h *FornecedorHandler) Ativar(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	fornecedor, err := h.repo.FindByID(c.Request().Context(), tenantID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro ao buscar fornecedor"})
	}
	if fornecedor == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "fornecedor não encontrado"})
	}

	fornecedor.Reativar()
	if err := h.repo.Update(c.Request().Context(), fornecedor); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro ao ativar fornecedor"})
	}

	return c.JSON(http.StatusOK, toFornecedorResponse(fornecedor))
}

// Desativar desativa um fornecedor
func (h *FornecedorHandler) Desativar(c echo.Context) error {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	fornecedor, err := h.repo.FindByID(c.Request().Context(), tenantID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro ao buscar fornecedor"})
	}
	if fornecedor == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "fornecedor não encontrado"})
	}

	fornecedor.Desativar()
	if err := h.repo.Update(c.Request().Context(), fornecedor); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro ao desativar fornecedor"})
	}

	return c.JSON(http.StatusOK, toFornecedorResponse(fornecedor))
}

// =============================================================================
// HELPERS
// =============================================================================

func (h *FornecedorHandler) getTenantID(c echo.Context) (uuid.UUID, error) {
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

func toFornecedorResponse(f *entity.Fornecedor) FornecedorResponse {
	// Formatar endereço completo
	endereco := ""
	if f.EnderecoLogradouro != "" {
		endereco = f.EnderecoLogradouro
		if f.EnderecoNumero != "" {
			endereco += ", " + f.EnderecoNumero
		}
		if f.EnderecoComplemento != "" {
			endereco += " - " + f.EnderecoComplemento
		}
		if f.EnderecoBairro != "" {
			endereco += ", " + f.EnderecoBairro
		}
		if f.EnderecoCidade != "" && f.EnderecoEstado != "" {
			endereco += ", " + f.EnderecoCidade + "/" + f.EnderecoEstado
		}
		if f.EnderecoCEP != "" {
			endereco += " - CEP: " + f.EnderecoCEP
		}
	}

	// Nome para compatibilidade com frontend (usa razao_social ou nome_fantasia)
	nome := f.RazaoSocial
	if f.NomeFantasia != "" {
		nome = f.NomeFantasia
	}

	return FornecedorResponse{
		ID:                  f.ID.String(),
		TenantID:            f.TenantID.String(),
		RazaoSocial:         f.RazaoSocial,
		NomeFantasia:        f.NomeFantasia,
		Nome:                nome,
		CNPJ:                f.CNPJ,
		Email:               f.Email,
		Telefone:            f.Telefone,
		Celular:             f.Celular,
		EnderecoLogradouro:  f.EnderecoLogradouro,
		EnderecoNumero:      f.EnderecoNumero,
		EnderecoComplemento: f.EnderecoComplemento,
		EnderecoBairro:      f.EnderecoBairro,
		EnderecoCidade:      f.EnderecoCidade,
		EnderecoEstado:      f.EnderecoEstado,
		EnderecoCEP:         f.EnderecoCEP,
		Endereco:            endereco,
		Banco:               f.Banco,
		Agencia:             f.Agencia,
		Conta:               f.Conta,
		Observacoes:         f.Observacoes,
		Ativo:               f.Ativo,
		CreatedAt:           f.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           f.UpdatedAt.Format(time.RFC3339),
	}
}
