package dto

import "time"

// ============================================================================
// Request DTOs - Despesas Fixas
// ============================================================================

// CreateDespesaFixaRequest representa o payload de criação de despesa fixa
type CreateDespesaFixaRequest struct {
	Descricao     string `json:"descricao" validate:"required,min=3,max=255"`
	CategoriaID   string `json:"categoria_id,omitempty"`
	Fornecedor    string `json:"fornecedor,omitempty"`
	Valor         string `json:"valor" validate:"required"` // Dinheiro sempre string
	DiaVencimento int    `json:"dia_vencimento" validate:"required,min=1,max=31"`
	UnidadeID     string `json:"unidade_id,omitempty"`
	Observacoes   string `json:"observacoes,omitempty"`
}

// UpdateDespesaFixaRequest representa o payload de atualização de despesa fixa
type UpdateDespesaFixaRequest struct {
	Descricao     string `json:"descricao" validate:"required,min=3,max=255"`
	CategoriaID   string `json:"categoria_id,omitempty"`
	Fornecedor    string `json:"fornecedor,omitempty"`
	Valor         string `json:"valor" validate:"required"` // Dinheiro sempre string
	DiaVencimento int    `json:"dia_vencimento" validate:"required,min=1,max=31"`
	UnidadeID     string `json:"unidade_id,omitempty"`
	Observacoes   string `json:"observacoes,omitempty"`
}

// GerarContasRequest representa o payload para geração de contas a partir de despesas fixas
type GerarContasRequest struct {
	Ano int `json:"ano,omitempty"` // Se vazio, usa ano atual
	Mes int `json:"mes,omitempty"` // Se vazio, usa mês atual
}

// ============================================================================
// Response DTOs - Despesas Fixas
// ============================================================================

// DespesaFixaResponse representa a resposta de uma despesa fixa
type DespesaFixaResponse struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id,omitempty"` // Opcional, pode ser omitido
	UnidadeID     string    `json:"unidade_id,omitempty"`
	Descricao     string    `json:"descricao"`
	CategoriaID   string    `json:"categoria_id,omitempty"`
	Fornecedor    string    `json:"fornecedor,omitempty"`
	Valor         string    `json:"valor"` // Dinheiro sempre string
	DiaVencimento int       `json:"dia_vencimento"`
	Ativo         bool      `json:"ativo"`
	Observacoes   string    `json:"observacoes,omitempty"`
	CriadoEm      time.Time `json:"criado_em"`
	AtualizadoEm  time.Time `json:"atualizado_em"`
}

// DespesasFixasListResponse representa a resposta paginada de listagem
type DespesasFixasListResponse struct {
	Data       []DespesaFixaResponse `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

// DespesasFixasSummaryResponse representa um resumo das despesas fixas
type DespesasFixasSummaryResponse struct {
	Total       int64  `json:"total"`
	TotalAtivas int64  `json:"total_ativas"`
	ValorTotal  string `json:"valor_total"` // Dinheiro sempre string
}

// GerarContasResponse representa a resposta da geração de contas
type GerarContasResponse struct {
	TotalDespesas   int      `json:"total_despesas"`
	ContasCriadas   int      `json:"contas_criadas"`
	Erros           int      `json:"erros"`
	DetalhesErros   []string `json:"detalhes_erros,omitempty"`
	TempoExecucaoMs int64    `json:"tempo_execucao_ms"`
}
