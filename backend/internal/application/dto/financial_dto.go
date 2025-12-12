// Package dto contém os Data Transfer Objects para a camada de aplicação
package dto

// CreateContaPagarRequest representa a requisição para criar conta a pagar
type CreateContaPagarRequest struct {
	Descricao      string `json:"descricao" validate:"required,min=3"`
	CategoriaID    string `json:"categoria_id" validate:"required,uuid"`
	Fornecedor     string `json:"fornecedor" validate:"required"`
	Valor          string `json:"valor" validate:"required"`
	Tipo           string `json:"tipo" validate:"required,oneof=FIXO VARIAVEL"`
	DataVencimento string `json:"data_vencimento" validate:"required"`
	Recorrente     bool   `json:"recorrente"`
	Periodicidade  string `json:"periodicidade,omitempty"`
	PixCode        string `json:"pix_code,omitempty"`
	Observacoes    string `json:"observacoes,omitempty"`
}

// UpdateContaPagarRequest representa a requisição para atualizar conta a pagar
type UpdateContaPagarRequest struct {
	Descricao      *string `json:"descricao,omitempty" validate:"omitempty,min=3"`
	CategoriaID    *string `json:"categoria_id,omitempty" validate:"omitempty,uuid"`
	Fornecedor     *string `json:"fornecedor,omitempty"`
	Valor          *string `json:"valor,omitempty"`
	Tipo           *string `json:"tipo,omitempty" validate:"omitempty,oneof=FIXO VARIAVEL"`
	DataVencimento *string `json:"data_vencimento,omitempty"`
	Recorrente     *bool   `json:"recorrente,omitempty"`
	Periodicidade  *string `json:"periodicidade,omitempty"`
	PixCode        *string `json:"pix_code,omitempty"`
	Observacoes    *string `json:"observacoes,omitempty"`
}

// ListContasPagarRequest representa filtros para listagem de contas a pagar
type ListContasPagarRequest struct {
	// Status canônicos: PENDENTE, PAGO, ATRASADO, CANCELADO
	// Legado (compat): ABERTO, VENCIDO
	Status      *string `query:"status" validate:"omitempty,oneof=PENDENTE ABERTO PAGO ATRASADO VENCIDO CANCELADO"`
	CategoriaID *string `query:"categoria_id" validate:"omitempty,uuid"`
	DataInicio  *string `query:"data_inicio"`
	DataFim     *string `query:"data_fim"`
	Tipo        *string `query:"tipo" validate:"omitempty,oneof=FIXO VARIAVEL"`
	Page        int     `query:"page" validate:"min=1"`
	PageSize    int     `query:"page_size" validate:"min=1,max=100"`
}

// ContaPagarResponse representa a resposta de conta a pagar
type ContaPagarResponse struct {
	ID             string  `json:"id"`
	Descricao      string  `json:"descricao"`
	CategoriaID    string  `json:"categoria_id"`
	Fornecedor     string  `json:"fornecedor"`
	Valor          string  `json:"valor"`
	Tipo           string  `json:"tipo"`
	Recorrente     bool    `json:"recorrente"`
	Periodicidade  string  `json:"periodicidade,omitempty"`
	DataVencimento string  `json:"data_vencimento"`
	DataPagamento  *string `json:"data_pagamento,omitempty"`
	Status         string  `json:"status"`
	ComprovanteURL string  `json:"comprovante_url,omitempty"`
	PixCode        string  `json:"pix_code,omitempty"`
	Observacoes    string  `json:"observacoes,omitempty"`
	CriadoEm       string  `json:"criado_em"`
	AtualizadoEm   string  `json:"atualizado_em"`
}

// CreateContaReceberRequest representa a requisição para criar conta a receber
type CreateContaReceberRequest struct {
	Origem          string  `json:"origem" validate:"required"`
	AssinaturaID    *string `json:"assinatura_id,omitempty" validate:"omitempty,uuid"`
	DescricaoOrigem string  `json:"descricao_origem" validate:"required"`
	Valor           string  `json:"valor" validate:"required"`
	DataVencimento  string  `json:"data_vencimento" validate:"required"`
	Observacoes     string  `json:"observacoes,omitempty"`
}

// UpdateContaReceberRequest representa a requisição para atualizar conta a receber
type UpdateContaReceberRequest struct {
	Origem          *string `json:"origem,omitempty"`
	AssinaturaID    *string `json:"assinatura_id,omitempty" validate:"omitempty,uuid"`
	DescricaoOrigem *string `json:"descricao_origem,omitempty"`
	Valor           *string `json:"valor,omitempty"`
	DataVencimento  *string `json:"data_vencimento,omitempty"`
	Observacoes     *string `json:"observacoes,omitempty"`
}

// ListContasReceberRequest representa filtros para listagem de contas a receber
type ListContasReceberRequest struct {
	// Status canônicos: PENDENTE, CONFIRMADO, RECEBIDO, ATRASADO, ESTORNADO, CANCELADO
	// Legado (compat): PAGO
	Status     *string `query:"status" validate:"omitempty,oneof=PENDENTE CONFIRMADO RECEBIDO ATRASADO ESTORNADO CANCELADO PAGO"`
	Origem     *string `query:"origem"`
	DataInicio *string `query:"data_inicio"`
	DataFim    *string `query:"data_fim"`
	Page       int     `query:"page" validate:"min=1"`
	PageSize   int     `query:"page_size" validate:"min=1,max=100"`
}

// ContaReceberResponse representa a resposta de conta a receber
type ContaReceberResponse struct {
	ID              string  `json:"id"`
	Origem          string  `json:"origem"`
	AssinaturaID    *string `json:"assinatura_id,omitempty"`
	DescricaoOrigem string  `json:"descricao_origem"`
	Valor           string  `json:"valor"`
	ValorPago       string  `json:"valor_pago"`
	ValorAberto     string  `json:"valor_aberto"`
	DataVencimento  string  `json:"data_vencimento"`
	DataRecebimento *string `json:"data_recebimento,omitempty"`
	Status          string  `json:"status"`
	Observacoes     string  `json:"observacoes,omitempty"`
	CriadoEm        string  `json:"criado_em"`
	AtualizadoEm    string  `json:"atualizado_em"`
}

// MarcarPagamentoRequest representa a requisição para marcar pagamento
type MarcarPagamentoRequest struct {
	DataPagamento  string `json:"data_pagamento" validate:"required"`
	ComprovanteURL string `json:"comprovante_url,omitempty"`
}

// MarcarRecebimentoRequest representa a requisição para marcar recebimento
type MarcarRecebimentoRequest struct {
	ValorPago       string `json:"valor_pago" validate:"required"`
	DataRecebimento string `json:"data_recebimento" validate:"required"`
}

// FluxoCaixaDiarioResponse representa a resposta de fluxo de caixa diário
type FluxoCaixaDiarioResponse struct {
	ID                  string `json:"id"`
	Data                string `json:"data"`
	SaldoInicial        string `json:"saldo_inicial"`
	EntradasConfirmadas string `json:"entradas_confirmadas"`
	EntradasPrevistas   string `json:"entradas_previstas"`
	SaidasPagas         string `json:"saidas_pagas"`
	SaidasPrevistas     string `json:"saidas_previstas"`
	SaldoFinal          string `json:"saldo_final"`
	ProcessadoEm        string `json:"processado_em"`
}

// CompensacaoBancariaResponse representa a resposta de compensação bancária
type CompensacaoBancariaResponse struct {
	ID              string  `json:"id"`
	ReceitaID       string  `json:"receita_id"`
	DataTransacao   string  `json:"data_transacao"`
	DataCompensacao string  `json:"data_compensacao"`
	DataCompensado  *string `json:"data_compensado,omitempty"`
	ValorBruto      string  `json:"valor_bruto"`
	TaxaPercentual  string  `json:"taxa_percentual"`
	TaxaFixa        string  `json:"taxa_fixa"`
	ValorLiquido    string  `json:"valor_liquido"`
	MeioPagamentoID string  `json:"meio_pagamento_id"`
	DMais           int     `json:"d_mais"`
	Status          string  `json:"status"`
}

// ListCompensacoesRequest representa filtros para listagem de compensações
type ListCompensacoesRequest struct {
	Status     *string `query:"status" validate:"omitempty,oneof=PREVISTO CONFIRMADO COMPENSADO CANCELADO"`
	DataInicio *string `query:"data_inicio"`
	DataFim    *string `query:"data_fim"`
	Page       int     `query:"page" validate:"min=1"`
	PageSize   int     `query:"page_size" validate:"min=1,max=100"`
}

// ListFluxoCaixaRequest representa filtros para listagem de fluxo de caixa
type ListFluxoCaixaRequest struct {
	DataInicio *string `query:"data_inicio"`
	DataFim    *string `query:"data_fim"`
	Page       int     `query:"page" validate:"min=1"`
	PageSize   int     `query:"page_size" validate:"min=1,max=100"`
}

// ListDRERequest representa filtros para listagem de DRE
type ListDRERequest struct {
	MesAnoInicio *string `query:"mes_ano_inicio"`
	MesAnoFim    *string `query:"mes_ano_fim"`
	Page         int     `query:"page" validate:"min=1"`
	PageSize     int     `query:"page_size" validate:"min=1,max=100"`
}

// DREMensalResponse representa a resposta de DRE mensal
type DREMensalResponse struct {
	ID                   string `json:"id"`
	MesAno               string `json:"mes_ano"`
	ReceitaServicos      string `json:"receita_servicos"`
	ReceitaProdutos      string `json:"receita_produtos"`
	ReceitaPlanos        string `json:"receita_planos"`
	ReceitaTotal         string `json:"receita_total"`
	CustoComissoes       string `json:"custo_comissoes"`
	CustoInsumos         string `json:"custo_insumos"`
	CustoVariavelTotal   string `json:"custo_variavel_total"`
	DespesaFixa          string `json:"despesa_fixa"`
	DespesaVariavel      string `json:"despesa_variavel"`
	DespesaTotal         string `json:"despesa_total"`
	ResultadoBruto       string `json:"resultado_bruto"`
	ResultadoOperacional string `json:"resultado_operacional"`
	MargemBruta          string `json:"margem_bruta"`
	MargemOperacional    string `json:"margem_operacional"`
	LucroLiquido         string `json:"lucro_liquido"`
	ProcessadoEm         string `json:"processado_em"`
}

// ErrorResponse representa uma resposta de erro padronizada
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// SuccessResponse representa uma resposta de sucesso genérica
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse representa uma resposta paginada
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}
