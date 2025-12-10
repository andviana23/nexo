package dto

import "time"

// =============================================================================
// REQUEST DTOs
// =============================================================================

// CreateMeioPagamentoRequest representa a requisição para criar um meio de pagamento
type CreateMeioPagamentoRequest struct {
	Nome          string `json:"nome,omitempty" validate:"max=100"`
	Tipo          string `json:"tipo" validate:"required,oneof=DINHEIRO PIX CREDITO DEBITO TRANSFERENCIA BOLETO OUTRO"`
	Bandeira      string `json:"bandeira,omitempty" validate:"max=50"`
	Taxa          string `json:"taxa,omitempty"`      // Percentual (0-100)
	TaxaFixa      string `json:"taxa_fixa,omitempty"` // Valor fixo R$
	DMais         int    `json:"d_mais,omitempty"`    // Dias para compensação
	Icone         string `json:"icone,omitempty"`     // Material Icons name
	Cor           string `json:"cor,omitempty"`       // Hexadecimal #RRGGBB
	OrdemExibicao int    `json:"ordem_exibicao,omitempty"`
	Observacoes   string `json:"observacoes,omitempty" validate:"max=500"`
	Ativo         *bool  `json:"ativo,omitempty"`
}

// UpdateMeioPagamentoRequest representa a requisição para atualizar um meio de pagamento
type UpdateMeioPagamentoRequest struct {
	Nome          string `json:"nome,omitempty" validate:"max=100"`
	Tipo          string `json:"tipo,omitempty" validate:"omitempty,oneof=DINHEIRO PIX CREDITO DEBITO TRANSFERENCIA BOLETO OUTRO"`
	Bandeira      string `json:"bandeira,omitempty" validate:"max=50"`
	Taxa          string `json:"taxa,omitempty"`
	TaxaFixa      string `json:"taxa_fixa,omitempty"`
	DMais         *int   `json:"d_mais,omitempty"`
	Icone         string `json:"icone,omitempty"`
	Cor           string `json:"cor,omitempty"`
	OrdemExibicao *int   `json:"ordem_exibicao,omitempty"`
	Observacoes   string `json:"observacoes,omitempty" validate:"max=500"`
	Ativo         *bool  `json:"ativo,omitempty"`
}

// =============================================================================
// RESPONSE DTOs
// =============================================================================

// MeioPagamentoResponse representa a resposta com dados de um meio de pagamento
type MeioPagamentoResponse struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	Nome          string    `json:"nome"`
	Tipo          string    `json:"tipo"`
	TipoLabel     string    `json:"tipo_label"` // Nome amigável do tipo
	Bandeira      string    `json:"bandeira,omitempty"`
	Taxa          string    `json:"taxa"`      // Percentual como string
	TaxaFixa      string    `json:"taxa_fixa"` // Valor fixo como string
	DMais         int       `json:"d_mais"`
	DMaisLabel    string    `json:"d_mais_label"` // Ex: "D+30"
	Icone         string    `json:"icone,omitempty"`
	Cor           string    `json:"cor,omitempty"`
	OrdemExibicao int       `json:"ordem_exibicao"`
	Observacoes   string    `json:"observacoes,omitempty"`
	Ativo         bool      `json:"ativo"`
	CriadoEm      time.Time `json:"criado_em"`
	AtualizadoEm  time.Time `json:"atualizado_em"`
}

// ListMeiosPagamentoResponse representa a lista paginada de meios de pagamento
type ListMeiosPagamentoResponse struct {
	Data       []MeioPagamentoResponse `json:"data"`
	Total      int64                   `json:"total"`
	TotalAtivo int64                   `json:"total_ativo"`
}

// =============================================================================
// FILTER
// =============================================================================

// MeioPagamentoFilter filtros para listagem
type MeioPagamentoFilter struct {
	Tipo         string `query:"tipo"`          // Filtrar por tipo
	ApenasAtivos bool   `query:"apenas_ativos"` // Apenas ativos
	Search       string `query:"search"`        // Busca por nome
}
