package dto

// CreateSubscriptionRequest representa o payload para nova assinatura
type CreateSubscriptionRequest struct {
	ClienteID       string  `json:"cliente_id" validate:"required,uuid"`
	PlanoID         string  `json:"plano_id" validate:"required,uuid"`
	FormaPagamento  string  `json:"forma_pagamento" validate:"required,oneof=CARTAO PIX DINHEIRO"`
	CodigoTransacao *string `json:"codigo_transacao,omitempty"`
}

// RenewSubscriptionRequest representa o payload de renovação manual
type RenewSubscriptionRequest struct {
	FormaPagamento  string  `json:"forma_pagamento" validate:"required,oneof=PIX DINHEIRO"`
	CodigoTransacao *string `json:"codigo_transacao,omitempty"`
}

// SubscriptionResponse representa a resposta de uma assinatura
type SubscriptionResponse struct {
	ID                 string  `json:"id"`
	ClienteID          string  `json:"cliente_id"`
	ClienteNome        string  `json:"cliente_nome"`
	ClienteTelefone    string  `json:"cliente_telefone"`
	PlanoID            string  `json:"plano_id"`
	PlanoNome          string  `json:"plano_nome"`
	FormaPagamento     string  `json:"forma_pagamento"`
	Status             string  `json:"status"`
	Valor              string  `json:"valor"`
	LinkPagamento      *string `json:"link_pagamento,omitempty"`
	DataAtivacao       *string `json:"data_ativacao,omitempty"`
	DataVencimento     *string `json:"data_vencimento,omitempty"`
	ServicosUtilizados int     `json:"servicos_utilizados"`
	CreatedAt          string  `json:"created_at"`
}

// SubscriptionMetricsResponse representa métricas agregadas
type SubscriptionMetricsResponse struct {
	TotalAssinantesAtivos   int     `json:"total_assinantes_ativos"`
	TotalInativas           int     `json:"total_inativas"`
	TotalInadimplentes      int     `json:"total_inadimplentes"`
	TotalPlanosAtivos       int     `json:"total_planos_ativos"`
	ReceitaMensal           float64 `json:"receita_mensal"`
	TaxaRenovacao           float64 `json:"taxa_renovacao"`
	RenovacoesProximos7Dias int     `json:"renovacoes_proximos_7_dias"`
}

// ReconcileAsaasResponse representa o resultado da reconciliação Asaas
// T-ASAAS-002: Reconciliação automática
type ReconcileAsaasResponse struct {
	TotalProcessed    int      `json:"total_processed"`
	TotalMatched      int      `json:"total_matched"`
	TotalMissingNexo  int      `json:"total_missing_nexo"`  // Existe no Asaas mas não no NEXO
	TotalMissingAsaas int      `json:"total_missing_asaas"` // Existe no NEXO mas não no Asaas
	TotalDivergent    int      `json:"total_divergent"`     // Valores diferentes
	TotalAutoFixed    int      `json:"total_auto_fixed"`    // Contas criadas automaticamente
	Errors            []string `json:"errors,omitempty"`
	ReconciliationID  string   `json:"reconciliation_id,omitempty"`
}
