// Package dto contém os Data Transfer Objects para a camada de aplicação
package dto

// ============================================================
// CAIXA DIÁRIO - DTOs
// Módulo de controle operacional da gaveta de dinheiro
// ============================================================

// ---------- ABERTURA ----------

// AbrirCaixaRequest representa a requisição para abrir o caixa
type AbrirCaixaRequest struct {
	SaldoInicial string `json:"saldo_inicial" validate:"required"`
}

// ---------- SANGRIA ----------

// SangriaRequest representa a requisição para registrar sangria
type SangriaRequest struct {
	Valor     string `json:"valor" validate:"required"`
	Destino   string `json:"destino" validate:"required,oneof=DEPOSITO PAGAMENTO COFRE OUTROS"`
	Descricao string `json:"descricao" validate:"required,min=5"`
}

// ---------- REFORÇO ----------

// ReforcoRequest representa a requisição para registrar reforço
type ReforcoRequest struct {
	Valor     string `json:"valor" validate:"required"`
	Origem    string `json:"origem" validate:"required,oneof=TROCO CAPITAL_GIRO TRANSFERENCIA OUTROS"`
	Descricao string `json:"descricao" validate:"required,min=5"`
}

// ---------- FECHAMENTO ----------

// FecharCaixaRequest representa a requisição para fechar o caixa
type FecharCaixaRequest struct {
	SaldoReal     string  `json:"saldo_real" validate:"required"`
	Justificativa *string `json:"justificativa,omitempty"`
}

// ---------- LISTAGEM ----------

// ListCaixaHistoricoRequest representa filtros para listagem do histórico
type ListCaixaHistoricoRequest struct {
	DataInicio *string `query:"data_inicio"`
	DataFim    *string `query:"data_fim"`
	UsuarioID  *string `query:"usuario_id" validate:"omitempty,uuid"`
	Page       int     `query:"page" validate:"min=1"`
	PageSize   int     `query:"page_size" validate:"min=1,max=100"`
}

// ---------- RESPONSES ----------

// CaixaDiarioResponse representa a resposta de um caixa diário
type CaixaDiarioResponse struct {
	ID                       string                  `json:"id"`
	UsuarioAberturaID        string                  `json:"usuario_abertura_id"`
	UsuarioAberturaNome      string                  `json:"usuario_abertura_nome"`
	UsuarioFechamentoID      *string                 `json:"usuario_fechamento_id,omitempty"`
	UsuarioFechamentoNome    string                  `json:"usuario_fechamento_nome,omitempty"`
	DataAbertura             string                  `json:"data_abertura"`
	DataFechamento           *string                 `json:"data_fechamento,omitempty"`
	SaldoInicial             string                  `json:"saldo_inicial"`
	TotalEntradas            string                  `json:"total_entradas"`
	TotalSaidas              string                  `json:"total_saidas"`
	TotalSangrias            string                  `json:"total_sangrias"`
	TotalReforcos            string                  `json:"total_reforcos"`
	SaldoEsperado            string                  `json:"saldo_esperado"`
	SaldoReal                *string                 `json:"saldo_real,omitempty"`
	Divergencia              *string                 `json:"divergencia,omitempty"`
	Status                   string                  `json:"status"`
	JustificativaDivergencia *string                 `json:"justificativa_divergencia,omitempty"`
	CreatedAt                string                  `json:"created_at"`
	UpdatedAt                string                  `json:"updated_at"`
	Operacoes                []OperacaoCaixaResponse `json:"operacoes,omitempty"`
}

// CaixaDiarioResumoResponse representa um resumo do caixa (para listagens)
type CaixaDiarioResumoResponse struct {
	ID                    string  `json:"id"`
	UsuarioAberturaNome   string  `json:"usuario_abertura_nome"`
	UsuarioFechamentoNome string  `json:"usuario_fechamento_nome,omitempty"`
	DataAbertura          string  `json:"data_abertura"`
	DataFechamento        *string `json:"data_fechamento,omitempty"`
	SaldoInicial          string  `json:"saldo_inicial"`
	SaldoEsperado         string  `json:"saldo_esperado"`
	SaldoReal             *string `json:"saldo_real,omitempty"`
	Divergencia           *string `json:"divergencia,omitempty"`
	Status                string  `json:"status"`
	TemDivergencia        bool    `json:"tem_divergencia"`
}

// OperacaoCaixaResponse representa uma operação do caixa
type OperacaoCaixaResponse struct {
	ID          string  `json:"id"`
	Tipo        string  `json:"tipo"`
	Valor       string  `json:"valor"`
	Descricao   string  `json:"descricao"`
	Destino     *string `json:"destino,omitempty"`
	Origem      *string `json:"origem,omitempty"`
	UsuarioID   string  `json:"usuario_id"`
	UsuarioNome string  `json:"usuario_nome"`
	CreatedAt   string  `json:"created_at"`
}

// CaixaStatusResponse representa o status atual do caixa
type CaixaStatusResponse struct {
	Aberto           bool                 `json:"aberto"`
	CaixaAtual       *CaixaDiarioResponse `json:"caixa_atual,omitempty"`
	UltimoFechamento *string              `json:"ultimo_fechamento,omitempty"`
}

// ListCaixaHistoricoResponse representa a resposta paginada do histórico
type ListCaixaHistoricoResponse struct {
	Items      []CaixaDiarioResumoResponse `json:"items"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	PageSize   int                         `json:"page_size"`
	TotalPages int                         `json:"total_pages"`
}

// TotaisCaixaResponse representa os totais do caixa por tipo de operação
type TotaisCaixaResponse struct {
	TotalVendas   string `json:"total_vendas"`
	TotalSangrias string `json:"total_sangrias"`
	TotalReforcos string `json:"total_reforcos"`
	TotalDespesas string `json:"total_despesas"`
	SaldoAtual    string `json:"saldo_atual"`
}
