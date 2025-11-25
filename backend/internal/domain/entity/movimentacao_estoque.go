package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TipoMovimentacao representa os tipos de movimentação de estoque
type TipoMovimentacao string

const (
	MovimentacaoEntrada        TipoMovimentacao = "ENTRADA"
	MovimentacaoSaida          TipoMovimentacao = "SAIDA"
	MovimentacaoConsumoInterno TipoMovimentacao = "CONSUMO_INTERNO"
	MovimentacaoAjuste         TipoMovimentacao = "AJUSTE"
	MovimentacaoDevolucao      TipoMovimentacao = "DEVOLUCAO"
	MovimentacaoPerda          TipoMovimentacao = "PERDA"
)

// Erros do domínio
var (
	ErrMovimentacaoQuantidadeInvalida    = errors.New("quantidade deve ser maior que zero")
	ErrMovimentacaoObservacaoObrigatoria = errors.New("observação é obrigatória para este tipo de movimentação")
)

// MovimentacaoEstoque representa uma movimentação de estoque
type MovimentacaoEstoque struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	ProdutoID        uuid.UUID
	UsuarioID        uuid.UUID
	FornecedorID     *uuid.UUID // Apenas para ENTRADA
	Tipo             TipoMovimentacao
	Quantidade       decimal.Decimal // Numeric(15,3) no banco
	ValorUnitario    decimal.Decimal // Numeric(15,2) no banco
	ValorTotal       decimal.Decimal // Numeric(15,2) no banco
	DataMovimentacao time.Time       // Data da movimentação (default NOW)
	Observacoes      string
	Documento        string // Número de nota fiscal, etc
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewMovimentacaoEstoque cria uma nova movimentação com validações
func NewMovimentacaoEstoque(
	tenantID, produtoID, usuarioID uuid.UUID,
	tipo TipoMovimentacao,
	quantidade decimal.Decimal,
	valorUnitario decimal.Decimal,
	observacoes string,
) (*MovimentacaoEstoque, error) {
	// Validações
	if quantidade.LessThanOrEqual(decimal.Zero) {
		return nil, ErrMovimentacaoQuantidadeInvalida
	}

	// RN-EST-003: AJUSTE e PERDA exigem observações
	if (tipo == MovimentacaoAjuste || tipo == MovimentacaoPerda) && observacoes == "" {
		return nil, ErrMovimentacaoObservacaoObrigatoria
	}

	valorTotal := quantidade.Mul(valorUnitario)

	now := time.Now()
	return &MovimentacaoEstoque{
		ID:               uuid.New(),
		TenantID:         tenantID,
		ProdutoID:        produtoID,
		UsuarioID:        usuarioID,
		Tipo:             tipo,
		Quantidade:       quantidade,
		ValorUnitario:    valorUnitario,
		ValorTotal:       valorTotal,
		DataMovimentacao: now,
		Observacoes:      observacoes,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// DefinirFornecedor associa um fornecedor à movimentação (para ENTRADA)
func (m *MovimentacaoEstoque) DefinirFornecedor(fornecedorID uuid.UUID) {
	m.FornecedorID = &fornecedorID
}

// IsEntrada verifica se é uma movimentação de entrada
func (m *MovimentacaoEstoque) IsEntrada() bool {
	return m.Tipo == MovimentacaoEntrada
}

// IsSaida verifica se é uma movimentação de saída
func (m *MovimentacaoEstoque) IsSaida() bool {
	return m.Tipo == MovimentacaoSaida ||
		m.Tipo == MovimentacaoConsumoInterno ||
		m.Tipo == MovimentacaoPerda
}
