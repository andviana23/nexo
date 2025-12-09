package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TipoOperacaoCaixa representa os tipos de operação no caixa
type TipoOperacaoCaixa string

const (
	TipoOperacaoVenda   TipoOperacaoCaixa = "VENDA"
	TipoOperacaoSangria TipoOperacaoCaixa = "SANGRIA"
	TipoOperacaoReforco TipoOperacaoCaixa = "REFORCO"
	TipoOperacaoDespesa TipoOperacaoCaixa = "DESPESA"
)

// ValidarTipoOperacaoCaixa verifica se o tipo é válido
func ValidarTipoOperacaoCaixa(t string) bool {
	switch TipoOperacaoCaixa(t) {
	case TipoOperacaoVenda, TipoOperacaoSangria, TipoOperacaoReforco, TipoOperacaoDespesa:
		return true
	}
	return false
}

// DestinoSangria representa os destinos válidos para sangria
type DestinoSangria string

const (
	DestinoDeposito  DestinoSangria = "DEPOSITO"
	DestinoPagamento DestinoSangria = "PAGAMENTO"
	DestinoCofre     DestinoSangria = "COFRE"
	DestinoOutros    DestinoSangria = "OUTROS"
)

// ValidarDestinoSangria verifica se o destino é válido
func ValidarDestinoSangria(d string) bool {
	switch DestinoSangria(d) {
	case DestinoDeposito, DestinoPagamento, DestinoCofre, DestinoOutros:
		return true
	}
	return false
}

// OrigemReforco representa as origens válidas para reforço
type OrigemReforco string

const (
	OrigemTroco         OrigemReforco = "TROCO"
	OrigemCapitalGiro   OrigemReforco = "CAPITAL_GIRO"
	OrigemTransferencia OrigemReforco = "TRANSFERENCIA"
	OrigemOutros        OrigemReforco = "OUTROS"
)

// ValidarOrigemReforco verifica se a origem é válida
func ValidarOrigemReforco(o string) bool {
	switch OrigemReforco(o) {
	case OrigemTroco, OrigemCapitalGiro, OrigemTransferencia, OrigemOutros:
		return true
	}
	return false
}

// OperacaoCaixa representa uma operação registrada no caixa
type OperacaoCaixa struct {
	ID        uuid.UUID
	CaixaID   uuid.UUID
	TenantID  uuid.UUID
	Tipo      TipoOperacaoCaixa
	Valor     decimal.Decimal
	Descricao string
	Destino   *string // Para sangrias
	Origem    *string // Para reforços
	UsuarioID uuid.UUID
	CreatedAt time.Time

	// Relacionamento (carregado quando necessário)
	UsuarioNome string
}

// NewOperacaoSangria cria uma nova operação de sangria
func NewOperacaoSangria(caixaID, tenantID, usuarioID uuid.UUID, valor decimal.Decimal, destino, descricao string) (*OperacaoCaixa, error) {
	if caixaID == uuid.Nil {
		return nil, errors.New("caixa_id é obrigatório")
	}
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant_id é obrigatório")
	}
	if usuarioID == uuid.Nil {
		return nil, errors.New("usuario_id é obrigatório")
	}
	if valor.IsNegative() || valor.IsZero() {
		return nil, errors.New("valor deve ser positivo")
	}
	if !ValidarDestinoSangria(destino) {
		return nil, errors.New("destino inválido para sangria")
	}
	if len(descricao) < 5 {
		return nil, errors.New("descrição deve ter pelo menos 5 caracteres")
	}

	return &OperacaoCaixa{
		ID:        uuid.New(),
		CaixaID:   caixaID,
		TenantID:  tenantID,
		Tipo:      TipoOperacaoSangria,
		Valor:     valor,
		Descricao: descricao,
		Destino:   &destino,
		UsuarioID: usuarioID,
		CreatedAt: time.Now(),
	}, nil
}

// NewOperacaoReforco cria uma nova operação de reforço
func NewOperacaoReforco(caixaID, tenantID, usuarioID uuid.UUID, valor decimal.Decimal, origem, descricao string) (*OperacaoCaixa, error) {
	if caixaID == uuid.Nil {
		return nil, errors.New("caixa_id é obrigatório")
	}
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant_id é obrigatório")
	}
	if usuarioID == uuid.Nil {
		return nil, errors.New("usuario_id é obrigatório")
	}
	if valor.IsNegative() || valor.IsZero() {
		return nil, errors.New("valor deve ser positivo")
	}
	if !ValidarOrigemReforco(origem) {
		return nil, errors.New("origem inválida para reforço")
	}
	if len(descricao) < 5 {
		return nil, errors.New("descrição deve ter pelo menos 5 caracteres")
	}

	return &OperacaoCaixa{
		ID:        uuid.New(),
		CaixaID:   caixaID,
		TenantID:  tenantID,
		Tipo:      TipoOperacaoReforco,
		Valor:     valor,
		Descricao: descricao,
		Origem:    &origem,
		UsuarioID: usuarioID,
		CreatedAt: time.Now(),
	}, nil
}

// NewOperacaoVenda cria uma nova operação de venda em dinheiro
func NewOperacaoVenda(caixaID, tenantID, usuarioID uuid.UUID, valor decimal.Decimal, descricao string) (*OperacaoCaixa, error) {
	if caixaID == uuid.Nil {
		return nil, errors.New("caixa_id é obrigatório")
	}
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant_id é obrigatório")
	}
	if usuarioID == uuid.Nil {
		return nil, errors.New("usuario_id é obrigatório")
	}
	if valor.IsNegative() || valor.IsZero() {
		return nil, errors.New("valor deve ser positivo")
	}
	if descricao == "" {
		return nil, errors.New("descrição é obrigatória")
	}

	return &OperacaoCaixa{
		ID:        uuid.New(),
		CaixaID:   caixaID,
		TenantID:  tenantID,
		Tipo:      TipoOperacaoVenda,
		Valor:     valor,
		Descricao: descricao,
		UsuarioID: usuarioID,
		CreatedAt: time.Now(),
	}, nil
}

// IsEntrada verifica se a operação é uma entrada de dinheiro
func (o *OperacaoCaixa) IsEntrada() bool {
	return o.Tipo == TipoOperacaoVenda || o.Tipo == TipoOperacaoReforco
}

// IsSaida verifica se a operação é uma saída de dinheiro
func (o *OperacaoCaixa) IsSaida() bool {
	return o.Tipo == TipoOperacaoSangria || o.Tipo == TipoOperacaoDespesa
}
