package dto

// === PRODUTO DTOs ===

// CreateProdutoRequest representa a requisição para criar um produto
type CreateProdutoRequest struct {
	SKU              string  `json:"sku" validate:"required"`
	Nome             string  `json:"nome" validate:"required"`
	Descricao        string  `json:"descricao"`
	Categoria        string  `json:"categoria" validate:"required,oneof=POMADA SHAMPOO CREME LAMINA TOALHA CONSUMIVEL REVENDA"`
	UnidadeMedida    string  `json:"unidade_medida" validate:"required,oneof=UN KG G ML L"`
	ValorUnitario    string  `json:"valor_unitario" validate:"required"` // String para receber do frontend
	QuantidadeMinima int     `json:"quantidade_minima"`
	FornecedorID     *string `json:"fornecedor_id,omitempty"`
}

// UpdateProdutoRequest representa a requisição para atualizar um produto
type UpdateProdutoRequest struct {
	Nome             string  `json:"nome"`
	Descricao        string  `json:"descricao"`
	Categoria        string  `json:"categoria" validate:"omitempty,oneof=POMADA SHAMPOO CREME LAMINA TOALHA CONSUMIVEL REVENDA"`
	UnidadeMedida    string  `json:"unidade_medida" validate:"omitempty,oneof=UN KG G ML L"`
	ValorUnitario    string  `json:"valor_unitario,omitempty"`
	QuantidadeMinima int     `json:"quantidade_minima"`
	FornecedorID     *string `json:"fornecedor_id,omitempty"`
}

// ProdutoResponse representa a resposta de um produto
type ProdutoResponse struct {
	ID               string  `json:"id"`
	TenantID         string  `json:"tenant_id"`
	SKU              string  `json:"sku"`
	Nome             string  `json:"nome"`
	Descricao        string  `json:"descricao"`
	Categoria        string  `json:"categoria"`
	UnidadeMedida    string  `json:"unidade_medida"`
	ValorUnitario    string  `json:"valor_unitario"`
	QuantidadeAtual  int     `json:"quantidade_atual"`
	QuantidadeMinima int     `json:"quantidade_minima"`
	FornecedorID     *string `json:"fornecedor_id,omitempty"`
	EstaBaixo        bool    `json:"esta_baixo"`
	Ativa            bool    `json:"ativa"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// === MOVIMENTAÇÃO DTOs ===

// ItemEntradaRequest representa um item de entrada no estoque
type ItemEntradaRequest struct {
	ProdutoID     string `json:"produto_id" validate:"required,uuid"`
	Quantidade    int    `json:"quantidade" validate:"required,gt=0"`
	ValorUnitario string `json:"valor_unitario" validate:"required"`
}

// RegistrarEntradaRequest representa a requisição para registrar entrada de estoque
type RegistrarEntradaRequest struct {
	FornecedorID    string               `json:"fornecedor_id" validate:"required,uuid"`
	DataEntrada     string               `json:"data_entrada"` // YYYY-MM-DD
	Itens           []ItemEntradaRequest `json:"itens" validate:"required,min=1,dive"`
	Observacoes     string               `json:"observacoes"`
	GerarFinanceiro bool                 `json:"gerar_financeiro"` // Se true, cria conta a pagar
}

// RegistrarEntradaResponse representa a resposta de entrada de estoque
type RegistrarEntradaResponse struct {
	MovimentacoesIDs []string `json:"movimentacoes_ids"`
	ValorTotal       string   `json:"valor_total"`
	ItensProcessados int      `json:"itens_processados"`
}

// RegistrarSaidaRequest representa a requisição para registrar saída de estoque
type RegistrarSaidaRequest struct {
	ProdutoID   string `json:"produto_id" validate:"required,uuid"`
	Quantidade  string `json:"quantidade" validate:"required"` // Aceita decimal como string
	Motivo      string `json:"motivo" validate:"required,oneof=VENDA USO_INTERNO PERDA DEVOLUCAO"`
	Observacoes string `json:"observacoes"`
}

// AjustarEstoqueRequest representa a requisição para ajuste manual de estoque
type AjustarEstoqueRequest struct {
	ProdutoID      string `json:"produto_id" validate:"required,uuid"`
	NovaQuantidade string `json:"nova_quantidade" validate:"required"` // Aceita decimal como string
	Motivo         string `json:"motivo" validate:"required"`
}

// MovimentacaoResponse representa a resposta de uma movimentação
type MovimentacaoResponse struct {
	ID            string  `json:"id"`
	TenantID      string  `json:"tenant_id"`
	ProdutoID     string  `json:"produto_id"`
	ProdutoNome   string  `json:"produto_nome,omitempty"`
	UsuarioID     string  `json:"usuario_id"`
	FornecedorID  *string `json:"fornecedor_id,omitempty"`
	Tipo          string  `json:"tipo"`
	Quantidade    string  `json:"quantidade"` // Decimal como string
	ValorUnitario string  `json:"valor_unitario"`
	ValorTotal    string  `json:"valor_total"`
	Observacoes   string  `json:"observacoes"`
	Data          string  `json:"data"` // Data da movimentação
	CreatedAt     string  `json:"created_at"`
}

// === FORNECEDOR DTOs ===

// CreateFornecedorRequest representa a requisição para criar um fornecedor
type CreateFornecedorRequest struct {
	RazaoSocial  string `json:"razao_social" validate:"required"`
	NomeFantasia string `json:"nome_fantasia"`
	CNPJ         string `json:"cnpj"`
	Email        string `json:"email" validate:"omitempty,email"`
	Telefone     string `json:"telefone" validate:"required"`
	Endereco     string `json:"endereco"`
	Cidade       string `json:"cidade"`
	Estado       string `json:"estado" validate:"omitempty,len=2"`
	CEP          string `json:"cep"`
}

// UpdateFornecedorRequest representa a requisição para atualizar um fornecedor
type UpdateFornecedorRequest struct {
	RazaoSocial  string `json:"razao_social"`
	NomeFantasia string `json:"nome_fantasia"`
	CNPJ         string `json:"cnpj"`
	Email        string `json:"email" validate:"omitempty,email"`
	Telefone     string `json:"telefone"`
	Endereco     string `json:"endereco"`
	Cidade       string `json:"cidade"`
	Estado       string `json:"estado" validate:"omitempty,len=2"`
	CEP          string `json:"cep"`
}

// FornecedorResponse representa a resposta de um fornecedor
type FornecedorResponse struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenant_id"`
	RazaoSocial  string `json:"razao_social"`
	NomeFantasia string `json:"nome_fantasia"`
	CNPJ         string `json:"cnpj"`
	Email        string `json:"email"`
	Telefone     string `json:"telefone"`
	Endereco     string `json:"endereco"`
	Cidade       string `json:"cidade"`
	Estado       string `json:"estado"`
	CEP          string `json:"cep"`
	Ativo        bool   `json:"ativo"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// === RESPONSES COMUNS ===

// ListProdutosResponse representa a resposta da listagem de produtos
type ListProdutosResponse struct {
	Data  []ProdutoResponse `json:"data"`
	Total int               `json:"total"`
}

// ListFornecedoresResponse representa a resposta da listagem de fornecedores
type ListFornecedoresResponse struct {
	Data  []FornecedorResponse `json:"data"`
	Total int                  `json:"total"`
}

// ListMovimentacoesResponse representa a resposta da listagem de movimentações
type ListMovimentacoesResponse struct {
	Data  []MovimentacaoResponse `json:"data"`
	Total int                    `json:"total"`
}

// AlertaEstoqueBaixoResponse representa um alerta de estoque baixo
type AlertaEstoqueBaixoResponse struct {
	ProdutoID        string `json:"produto_id"`
	ProdutoNome      string `json:"produto_nome"`
	SKU              string `json:"sku"`
	QuantidadeAtual  int    `json:"quantidade_atual"`
	QuantidadeMinima int    `json:"quantidade_minima"`
	Diferenca        int    `json:"diferenca"`
}

// AlertaEstoqueBaixo representa um alerta de estoque baixo com detalhes
type AlertaEstoqueBaixo struct {
	ProdutoID        string `json:"produto_id"`
	SKU              string `json:"sku"`
	Nome             string `json:"nome"`
	Categoria        string `json:"categoria"`
	QuantidadeAtual  string `json:"quantidade_atual"`
	QuantidadeMinima string `json:"quantidade_minima"`
	UnidadeMedida    string `json:"unidade_medida"`
	Percentual       string `json:"percentual"` // Percentual do estoque (atual/mínimo * 100)
	Severidade       string `json:"severidade"` // crítico, alerta, baixo
}

// ListAlertasEstoqueBaixoResponse representa a resposta da listagem de alertas de estoque baixo
type ListAlertasEstoqueBaixoResponse struct {
	Total   int                  `json:"total"`
	Alertas []AlertaEstoqueBaixo `json:"alertas"`
}

// ListAlertasResponse representa a resposta da listagem de alertas (alias para compatibilidade)
type ListAlertasResponse struct {
	Alertas []AlertaEstoqueBaixoResponse `json:"alertas"`
	Total   int                          `json:"total"`
}
