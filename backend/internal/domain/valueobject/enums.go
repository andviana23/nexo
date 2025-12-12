package valueobject

// StatusCompensacao representa o status de uma compensação bancária
type StatusCompensacao string

const (
	StatusCompensacaoPrevisto   StatusCompensacao = "PREVISTO"
	StatusCompensacaoConfirmado StatusCompensacao = "CONFIRMADO"
	StatusCompensacaoCompensado StatusCompensacao = "COMPENSADO"
	StatusCompensacaoCancelado  StatusCompensacao = "CANCELADO"
)

// IsValid verifica se o status é válido
func (s StatusCompensacao) IsValid() bool {
	switch s {
	case StatusCompensacaoPrevisto, StatusCompensacaoConfirmado, StatusCompensacaoCompensado, StatusCompensacaoCancelado:
		return true
	}
	return false
}

// String retorna a string do status
func (s StatusCompensacao) String() string {
	return string(s)
}

// StatusConta representa o status de uma conta a pagar/receber
type StatusConta string

const (
	// Status canônicos (Domínio)
	StatusContaPendente   StatusConta = "PENDENTE"   // Conta criada, ainda não confirmada/quitada
	StatusContaConfirmado StatusConta = "CONFIRMADO" // Pagamento confirmado (ex.: cartão/Asaas), aguardando compensação
	StatusContaRecebido   StatusConta = "RECEBIDO"   // Receita recebida/compensada
	StatusContaPago       StatusConta = "PAGO"       // Despesa paga
	StatusContaEstornado  StatusConta = "ESTORNADO"  // Receita estornada/refund
	StatusContaCancelado  StatusConta = "CANCELADO"  // Conta cancelada/sem efeito
	StatusContaAtrasado   StatusConta = "ATRASADO"   // Vencida sem quitação
)

// IsValid verifica se o status é válido
func (s StatusConta) IsValid() bool {
	switch s {
	case StatusContaPendente,
		StatusContaConfirmado,
		StatusContaRecebido,
		StatusContaPago,
		StatusContaEstornado,
		StatusContaCancelado,
		StatusContaAtrasado:
		return true
	}
	return false
}

// String retorna a string do status
func (s StatusConta) String() string {
	return string(s)
}

// TipoCusto representa o tipo de custo (fixo ou variável)
type TipoCusto string

const (
	TipoCustoFixo     TipoCusto = "FIXO"
	TipoCustoVariavel TipoCusto = "VARIAVEL"
)

// IsValid verifica se o tipo é válido
func (t TipoCusto) IsValid() bool {
	return t == TipoCustoFixo || t == TipoCustoVariavel
}

// String retorna a string do tipo
func (t TipoCusto) String() string {
	return string(t)
}

// SubtipoReceita representa o subtipo de receita
type SubtipoReceita string

const (
	SubtipoReceitaServico SubtipoReceita = "SERVICO"
	SubtipoReceitaProduto SubtipoReceita = "PRODUTO"
	SubtipoReceitaPlano   SubtipoReceita = "PLANO"
)

// IsValid verifica se o subtipo é válido
func (s SubtipoReceita) IsValid() bool {
	switch s {
	case SubtipoReceitaServico, SubtipoReceitaProduto, SubtipoReceitaPlano:
		return true
	}
	return false
}

// String retorna a string do subtipo
func (s SubtipoReceita) String() string {
	return string(s)
}

// OrigemMeta representa a origem de uma meta (manual ou automática)
type OrigemMeta string

const (
	OrigemMetaManual     OrigemMeta = "MANUAL"
	OrigemMetaAutomatica OrigemMeta = "AUTOMATICA"
)

// IsValid verifica se a origem é válida
func (o OrigemMeta) IsValid() bool {
	return o == OrigemMetaManual || o == OrigemMetaAutomatica
}

// String retorna a string da origem
func (o OrigemMeta) String() string {
	return string(o)
}

// TipoMetaTicket representa o tipo de meta de ticket médio
type TipoMetaTicket string

const (
	TipoMetaTicketGeral    TipoMetaTicket = "GERAL"
	TipoMetaTicketBarbeiro TipoMetaTicket = "BARBEIRO"
)

// IsValid verifica se o tipo é válido
func (t TipoMetaTicket) IsValid() bool {
	return t == TipoMetaTicketGeral || t == TipoMetaTicketBarbeiro
}

// String retorna a string do tipo
func (t TipoMetaTicket) String() string {
	return string(t)
}
