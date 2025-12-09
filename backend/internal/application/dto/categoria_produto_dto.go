package dto

import "time"

// =============================================================================
// REQUEST DTOs
// =============================================================================

// CreateCategoriaProdutoRequest representa a request para criar categoria
type CreateCategoriaProdutoRequest struct {
	Nome        string `json:"nome" binding:"required,min=1,max=100"`
	Descricao   string `json:"descricao,omitempty"`
	Cor         string `json:"cor,omitempty" binding:"omitempty,len=7"`
	Icone       string `json:"icone,omitempty"`
	CentroCusto string `json:"centro_custo,omitempty" binding:"omitempty,oneof=CMV CUSTO_SERVICO DESPESA_OPERACIONAL"`
}

// UpdateCategoriaProdutoRequest representa a request para atualizar categoria
type UpdateCategoriaProdutoRequest struct {
	Nome        string `json:"nome" binding:"required,min=1,max=100"`
	Descricao   string `json:"descricao,omitempty"`
	Cor         string `json:"cor,omitempty" binding:"omitempty,len=7"`
	Icone       string `json:"icone,omitempty"`
	CentroCusto string `json:"centro_custo,omitempty" binding:"omitempty,oneof=CMV CUSTO_SERVICO DESPESA_OPERACIONAL"`
	Ativa       bool   `json:"ativa"`
}

// =============================================================================
// RESPONSE DTOs
// =============================================================================

// CategoriaProdutoResponse representa a resposta de categoria
type CategoriaProdutoResponse struct {
	ID           string    `json:"id"`
	Nome         string    `json:"nome"`
	Descricao    string    `json:"descricao,omitempty"`
	Cor          string    `json:"cor"`
	Icone        string    `json:"icone"`
	CentroCusto  string    `json:"centro_custo"`
	Ativa        bool      `json:"ativa"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// ListCategoriaProdutoResponse representa a resposta de listagem
type ListCategoriaProdutoResponse struct {
	Categorias []CategoriaProdutoResponse `json:"categorias"`
	Total      int                        `json:"total"`
}
