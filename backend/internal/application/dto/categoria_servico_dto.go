package dto

import "time"

// =============================================================================
// DTOs para Categorias de Serviço
// Conforme Sprint 1.4.1 - Módulo de Serviços
// =============================================================================

// =============================================================================
// Request DTOs
// =============================================================================

// CreateCategoriaServicoRequest requisição para criar categoria
type CreateCategoriaServicoRequest struct {
	Nome      string  `json:"nome" validate:"required,min=2,max=100"`
	Descricao *string `json:"descricao,omitempty" validate:"omitempty,max=500"`
	Cor       *string `json:"cor,omitempty" validate:"omitempty,hexcolor"`
	Icone     *string `json:"icone,omitempty" validate:"omitempty,max=50"`
}

// UpdateCategoriaServicoRequest requisição para atualizar categoria
type UpdateCategoriaServicoRequest struct {
	Nome      string  `json:"nome" validate:"required,min=2,max=100"`
	Descricao *string `json:"descricao,omitempty" validate:"omitempty,max=500"`
	Cor       *string `json:"cor,omitempty" validate:"omitempty,hexcolor"`
	Icone     *string `json:"icone,omitempty" validate:"omitempty,max=50"`
}

// ListCategoriasServicosRequest query params para listagem
type ListCategoriasServicosRequest struct {
	ApenasAtivas bool   `query:"apenas_ativas"`
	OrderBy      string `query:"order_by" validate:"omitempty,oneof=nome criado_em"`
}

// ToggleCategoriaServicoStatusRequest requisição para ativar/desativar categoria
type ToggleCategoriaServicoStatusRequest struct {
	Ativa bool `json:"ativa"`
}

// =============================================================================
// Response DTOs
// =============================================================================

// CategoriaServicoResponse resposta de categoria
type CategoriaServicoResponse struct {
	ID           string  `json:"id"`
	TenantID     string  `json:"tenant_id"`
	Nome         string  `json:"nome"`
	Descricao    *string `json:"descricao,omitempty"`
	Cor          *string `json:"cor,omitempty"`
	Icone        *string `json:"icone,omitempty"`
	Ativa        bool    `json:"ativa"`
	CriadoEm     string  `json:"criado_em"`
	AtualizadoEm string  `json:"atualizado_em"`
}

// CategoriaServicoWithCountResponse resposta com contagem de serviços
type CategoriaServicoWithCountResponse struct {
	CategoriaServicoResponse
	TotalServicos int64 `json:"total_servicos"`
}

// ListCategoriasServicosResponse resposta de listagem
type ListCategoriasServicosResponse struct {
	Categorias []*CategoriaServicoResponse `json:"categorias"`
	Total      int                         `json:"total"`
}

// =============================================================================
// Mappers
// =============================================================================

// ToCategoriaServicoResponse converte entity para response DTO
func ToCategoriaServicoResponse(categoria interface{}) *CategoriaServicoResponse {
	// Este método será implementado no mapper separado
	// Aqui deixamos apenas a assinatura para referência
	return nil
}

// FormatTimestamp formata time.Time para string ISO8601
func FormatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}
