package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// CommandItemType representa os tipos de itens da comanda
type CommandItemType string

const (
	CommandItemTypeServico CommandItemType = "SERVICO"
	CommandItemTypeProduto CommandItemType = "PRODUTO"
	CommandItemTypePacote  CommandItemType = "PACOTE"
)

// CommandItem representa um item da comanda
type CommandItem struct {
	ID        uuid.UUID
	CommandID uuid.UUID

	Tipo      CommandItemType
	ItemID    uuid.UUID
	Descricao string

	// Preços
	PrecoUnitario      float64
	Quantidade         int
	DescontoValor      float64
	DescontoPercentual float64
	PrecoFinal         float64

	Observacoes *string
	CriadoEm    time.Time
}

// NewCommandItem cria um novo item de comanda
func NewCommandItem(
	commandID uuid.UUID,
	tipo CommandItemType,
	itemID uuid.UUID,
	descricao string,
	precoUnitario float64,
	quantidade int,
) (*CommandItem, error) {
	if commandID == uuid.Nil {
		return nil, errors.New("command_id é obrigatório")
	}
	if itemID == uuid.Nil {
		return nil, errors.New("item_id é obrigatório")
	}
	if descricao == "" {
		return nil, errors.New("descrição é obrigatória")
	}
	if precoUnitario < 0 {
		return nil, errors.New("preço unitário não pode ser negativo")
	}
	if quantidade <= 0 {
		return nil, errors.New("quantidade deve ser maior que zero")
	}

	item := &CommandItem{
		ID:                 uuid.New(),
		CommandID:          commandID,
		Tipo:               tipo,
		ItemID:             itemID,
		Descricao:          descricao,
		PrecoUnitario:      precoUnitario,
		Quantidade:         quantidade,
		DescontoValor:      0,
		DescontoPercentual: 0,
		CriadoEm:           time.Now(),
	}

	item.CalculatePrecoFinal()

	return item, nil
}

// ApplyDiscount aplica desconto ao item (valor ou percentual)
func (ci *CommandItem) ApplyDiscount(valor float64, percentual float64) error {
	if valor < 0 {
		return errors.New("desconto em valor não pode ser negativo")
	}
	if percentual < 0 || percentual > 100 {
		return errors.New("desconto percentual deve estar entre 0 e 100")
	}

	ci.DescontoValor = valor
	ci.DescontoPercentual = percentual
	ci.CalculatePrecoFinal()

	return nil
}

// CalculatePrecoFinal calcula o preço final do item
func (ci *CommandItem) CalculatePrecoFinal() {
	subtotal := ci.PrecoUnitario * float64(ci.Quantidade)

	// Aplicar desconto percentual primeiro
	if ci.DescontoPercentual > 0 {
		subtotal -= subtotal * (ci.DescontoPercentual / 100)
	}

	// Depois aplicar desconto em valor
	subtotal -= ci.DescontoValor

	if subtotal < 0 {
		subtotal = 0
	}

	ci.PrecoFinal = subtotal
}

// UpdateQuantity atualiza a quantidade do item
func (ci *CommandItem) UpdateQuantity(quantidade int) error {
	if quantidade <= 0 {
		return errors.New("quantidade deve ser maior que zero")
	}

	ci.Quantidade = quantidade
	ci.CalculatePrecoFinal()

	return nil
}

// UpdatePrice atualiza o preço unitário do item
func (ci *CommandItem) UpdatePrice(preco float64) error {
	if preco < 0 {
		return errors.New("preço não pode ser negativo")
	}

	ci.PrecoUnitario = preco
	ci.CalculatePrecoFinal()

	return nil
}
