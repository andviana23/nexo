package dto

// =============================================================================
// DTOs para Serviços
// Conforme Sprint 1.4.2 - Módulo de Serviços
// =============================================================================

// =============================================================================
// Request DTOs
// =============================================================================

// CreateServicoRequest requisição para criar serviço
type CreateServicoRequest struct {
	CategoriaID      *string  `json:"categoria_id,omitempty" validate:"omitempty,uuid"`
	Nome             string   `json:"nome" validate:"required,min=2,max=255"`
	Descricao        *string  `json:"descricao,omitempty" validate:"omitempty,max=1000"`
	Preco            string   `json:"preco" validate:"required"` // Dinheiro como string
	Duracao          int      `json:"duracao" validate:"required,min=5"`
	Comissao         *string  `json:"comissao,omitempty" validate:"omitempty"` // Percentual como string
	Cor              *string  `json:"cor,omitempty" validate:"omitempty,hexcolor"`
	Imagem           *string  `json:"imagem,omitempty" validate:"omitempty,url|base64"`
	ProfissionaisIDs []string `json:"profissionais_ids,omitempty" validate:"omitempty,dive,uuid"`
	Observacoes      *string  `json:"observacoes,omitempty" validate:"omitempty,max=2000"`
	Tags             []string `json:"tags,omitempty" validate:"omitempty,dive,min=1,max=50"`
}

// UpdateServicoRequest requisição para atualizar serviço
type UpdateServicoRequest struct {
	CategoriaID      *string  `json:"categoria_id,omitempty" validate:"omitempty,uuid"`
	Nome             string   `json:"nome" validate:"required,min=2,max=255"`
	Descricao        *string  `json:"descricao,omitempty" validate:"omitempty,max=1000"`
	Preco            string   `json:"preco" validate:"required"`
	Duracao          int      `json:"duracao" validate:"required,min=5"`
	Comissao         *string  `json:"comissao,omitempty" validate:"omitempty"`
	Cor              *string  `json:"cor,omitempty" validate:"omitempty,hexcolor"`
	Imagem           *string  `json:"imagem,omitempty" validate:"omitempty,url|base64"`
	ProfissionaisIDs []string `json:"profissionais_ids,omitempty" validate:"omitempty,dive,uuid"`
	Observacoes      *string  `json:"observacoes,omitempty" validate:"omitempty,max=2000"`
	Tags             []string `json:"tags,omitempty" validate:"omitempty,dive,min=1,max=50"`
}

// ListServicosRequest query params para listagem
type ListServicosRequest struct {
	ApenasAtivos   bool   `query:"apenas_ativos"`
	CategoriaID    string `query:"categoria_id" validate:"omitempty,uuid"`
	ProfissionalID string `query:"profissional_id" validate:"omitempty,uuid"`
	Search         string `query:"search" validate:"omitempty,max=100"`
	OrderBy        string `query:"order_by" validate:"omitempty,oneof=nome preco duracao criado_em"`
}

// ToggleServicoStatusRequest requisição para ativar/desativar serviço
// UpdateServicoCategoriaRequest requisição para atualizar categoria de um serviço
type UpdateServicoCategoriaRequest struct {
	CategoriaID string `json:"categoria_id" validate:"required,uuid"`
}

// UpdateServicoProfissionaisRequest requisição para atualizar profissionais de um serviço
type UpdateServicoProfissionaisRequest struct {
	ProfissionaisIDs []string `json:"profissionais_ids" validate:"required,dive,uuid"`
}

// =============================================================================
// Response DTOs
// =============================================================================

// ServicoResponse resposta de serviço
type ServicoResponse struct {
	ID               string   `json:"id"`
	TenantID         string   `json:"tenant_id"`
	CategoriaID      *string  `json:"categoria_id,omitempty"`
	CategoriaNome    *string  `json:"categoria_nome,omitempty"`
	CategoriaCor     *string  `json:"categoria_cor,omitempty"`
	Nome             string   `json:"nome"`
	Descricao        *string  `json:"descricao,omitempty"`
	Preco            string   `json:"preco"` // Retornado como string para precisão
	PrecoCentavos    int64    `json:"preco_centavos"`
	Duracao          int      `json:"duracao"` // em minutos
	DuracaoFormatada string   `json:"duracao_formatada"`
	Comissao         string   `json:"comissao"` // Percentual como string
	Cor              *string  `json:"cor,omitempty"`
	Imagem           *string  `json:"imagem,omitempty"`
	ProfissionaisIDs []string `json:"profissionais_ids,omitempty"`
	Observacoes      *string  `json:"observacoes,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	Ativo            bool     `json:"ativo"`
	CriadoEm         string   `json:"criado_em"`
	AtualizadoEm     string   `json:"atualizado_em"`
}

// ServicoSimplificadoResponse resposta simplificada para listagens
type ServicoSimplificadoResponse struct {
	ID            string  `json:"id"`
	Nome          string  `json:"nome"`
	Preco         string  `json:"preco"`
	Duracao       int     `json:"duracao"`
	CategoriaNome *string `json:"categoria_nome,omitempty"`
	Cor           *string `json:"cor,omitempty"`
	Ativo         bool    `json:"ativo"`
}

// ListServicosResponse resposta de listagem
type ListServicosResponse struct {
	Servicos []*ServicoResponse `json:"servicos"`
	Total    int                `json:"total"`
}

// ServicoStatsResponse estatísticas de serviços
type ServicoStatsResponse struct {
	TotalServicos    int64   `json:"total_servicos"`
	ServicosAtivos   int64   `json:"servicos_ativos"`
	ServicosInativos int64   `json:"servicos_inativos"`
	PrecoMedio       string  `json:"preco_medio"`
	DuracaoMedia     float64 `json:"duracao_media"`
	ComissaoMedia    string  `json:"comissao_media"`
}
