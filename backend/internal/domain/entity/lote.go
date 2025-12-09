package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Erros de Lote
var (
	ErrLoteQuantidadeInvalida = errors.New("quantidade do lote deve ser maior que zero")
	ErrLoteVencido            = errors.New("lote vencido")
	ErrLoteEsgotado           = errors.New("lote esgotado")
)

// Lote representa um lote de produto com validade
type Lote struct {
	ID                uuid.UUID
	ProdutoID         uuid.UUID
	CodigoLote        string
	DataValidade      time.Time
	QuantidadeInicial int
	QuantidadeAtual   int
	Ativo             bool
	CriadoEm          time.Time
}

// NewLote cria um novo lote
func NewLote(
	produtoID uuid.UUID,
	codigoLote string,
	dataValidade time.Time,
	quantidade int,
) (*Lote, error) {
	if quantidade <= 0 {
		return nil, ErrLoteQuantidadeInvalida
	}

	return &Lote{
		ID:                uuid.New(),
		ProdutoID:         produtoID,
		CodigoLote:        codigoLote,
		DataValidade:      dataValidade,
		QuantidadeInicial: quantidade,
		QuantidadeAtual:   quantidade,
		Ativo:             true,
		CriadoEm:          time.Now(),
	}, nil
}

// EstaVencido verifica se o lote está vencido
func (l *Lote) EstaVencido() bool {
	return time.Now().After(l.DataValidade)
}

// Consumir remove quantidade do lote
func (l *Lote) Consumir(quantidade int) error {
	if quantidade <= 0 {
		return ErrLoteQuantidadeInvalida
	}
	if l.QuantidadeAtual < quantidade {
		return ErrLoteEsgotado
	}
	if l.EstaVencido() {
		return ErrLoteVencido
	}

	l.QuantidadeAtual -= quantidade
	if l.QuantidadeAtual == 0 {
		l.Ativo = false
	}
	return nil
}
