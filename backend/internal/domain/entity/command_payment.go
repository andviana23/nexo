package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// CommandPayment representa um pagamento da comanda
type CommandPayment struct {
	ID              uuid.UUID
	CommandID       uuid.UUID
	MeioPagamentoID uuid.UUID

	// Valores
	ValorRecebido  float64
	TaxaPercentual float64
	TaxaFixa       float64
	ValorLiquido   float64

	Observacoes *string
	CriadoEm    time.Time
	CriadoPor   *uuid.UUID
}

// NewCommandPayment cria um novo pagamento de comanda
func NewCommandPayment(
	commandID uuid.UUID,
	meioPagamentoID uuid.UUID,
	valorRecebido float64,
	taxaPercentual float64,
	taxaFixa float64,
	criadoPor *uuid.UUID,
) (*CommandPayment, error) {
	if commandID == uuid.Nil {
		return nil, errors.New("command_id é obrigatório")
	}
	if meioPagamentoID == uuid.Nil {
		return nil, errors.New("meio_pagamento_id é obrigatório")
	}
	if valorRecebido <= 0 {
		return nil, errors.New("valor recebido deve ser maior que zero")
	}
	if taxaPercentual < 0 || taxaPercentual > 100 {
		return nil, errors.New("taxa percentual deve estar entre 0 e 100")
	}
	if taxaFixa < 0 {
		return nil, errors.New("taxa fixa não pode ser negativa")
	}

	payment := &CommandPayment{
		ID:              uuid.New(),
		CommandID:       commandID,
		MeioPagamentoID: meioPagamentoID,
		ValorRecebido:   valorRecebido,
		TaxaPercentual:  taxaPercentual,
		TaxaFixa:        taxaFixa,
		CriadoEm:        time.Now(),
		CriadoPor:       criadoPor,
	}

	payment.CalculateValorLiquido()

	return payment, nil
}

// CalculateValorLiquido calcula o valor líquido após dedução de taxas
func (cp *CommandPayment) CalculateValorLiquido() {
	// Primeiro aplica taxa percentual
	valorTaxaPercentual := cp.ValorRecebido * (cp.TaxaPercentual / 100)

	// Depois deduz taxa fixa
	cp.ValorLiquido = cp.ValorRecebido - valorTaxaPercentual - cp.TaxaFixa

	// Garante que não fique negativo
	if cp.ValorLiquido < 0 {
		cp.ValorLiquido = 0
	}
}

// GetTotalTaxas retorna o total de taxas aplicadas
func (cp *CommandPayment) GetTotalTaxas() float64 {
	return cp.ValorRecebido - cp.ValorLiquido
}
