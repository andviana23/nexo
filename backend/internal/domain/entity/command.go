package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// CommandStatus representa os possíveis status de uma comanda
type CommandStatus string

const (
	CommandStatusOpen     CommandStatus = "OPEN"
	CommandStatusClosed   CommandStatus = "CLOSED"
	CommandStatusCanceled CommandStatus = "CANCELED"
)

// Command representa uma comanda de atendimento
type Command struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	AppointmentID *uuid.UUID
	CustomerID    uuid.UUID
	Numero        *string
	Status        CommandStatus

	// Valores financeiros
	Subtotal      float64
	Desconto      float64
	Total         float64
	TotalRecebido float64
	Troco         float64
	SaldoDevedor  float64

	// Opções de fechamento
	Observacoes        *string
	DeixarTrocoGorjeta bool
	DeixarSaldoDivida  bool

	// Auditoria
	CriadoEm     time.Time
	AtualizadoEm time.Time
	FechadoEm    *time.Time
	FechadoPor   *uuid.UUID

	// Relacionamentos (carregados quando necessário)
	Items    []CommandItem
	Payments []CommandPayment
}

// NewCommand cria uma nova comanda
func NewCommand(tenantID, customerID uuid.UUID, appointmentID *uuid.UUID) (*Command, error) {
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant_id é obrigatório")
	}
	if customerID == uuid.Nil {
		return nil, errors.New("customer_id é obrigatório")
	}

	now := time.Now()
	return &Command{
		ID:                 uuid.New(),
		TenantID:           tenantID,
		CustomerID:         customerID,
		AppointmentID:      appointmentID,
		Status:             CommandStatusOpen,
		Subtotal:           0,
		Desconto:           0,
		Total:              0,
		TotalRecebido:      0,
		Troco:              0,
		SaldoDevedor:       0,
		DeixarTrocoGorjeta: false,
		DeixarSaldoDivida:  false,
		CriadoEm:           now,
		AtualizadoEm:       now,
		Items:              []CommandItem{},
		Payments:           []CommandPayment{},
	}, nil
}

// CanClose verifica se a comanda pode ser fechada
func (c *Command) CanClose() error {
	if c.Status != CommandStatusOpen {
		return errors.New("comanda já está fechada ou cancelada")
	}

	if len(c.Items) == 0 {
		return errors.New("comanda sem itens não pode ser fechada")
	}

	// Se não permitir saldo devedor, deve estar totalmente pago
	if !c.DeixarSaldoDivida && c.TotalRecebido < c.Total {
		falta := c.Total - c.TotalRecebido
		return errors.New("falta receber R$ " + formatMoney(falta))
	}

	return nil
}

// Close fecha a comanda
func (c *Command) Close(userID uuid.UUID) error {
	if err := c.CanClose(); err != nil {
		return err
	}

	now := time.Now()
	c.Status = CommandStatusClosed
	c.FechadoEm = &now
	c.FechadoPor = &userID
	c.AtualizadoEm = now

	// Calcular troco ou saldo devedor
	c.CalculateBalance()

	return nil
}

// Cancel cancela a comanda
func (c *Command) Cancel() error {
	if c.Status != CommandStatusOpen {
		return errors.New("apenas comandas abertas podem ser canceladas")
	}

	c.Status = CommandStatusCanceled
	c.AtualizadoEm = time.Now()

	return nil
}

// AddItem adiciona um item à comanda
func (c *Command) AddItem(item CommandItem) error {
	if c.Status != CommandStatusOpen {
		return errors.New("não é possível adicionar itens a uma comanda fechada")
	}

	c.Items = append(c.Items, item)
	c.RecalculateTotals()

	return nil
}

// RemoveItem remove um item da comanda
func (c *Command) RemoveItem(itemID uuid.UUID) error {
	if c.Status != CommandStatusOpen {
		return errors.New("não é possível remover itens de uma comanda fechada")
	}

	for i, item := range c.Items {
		if item.ID == itemID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.RecalculateTotals()
			return nil
		}
	}

	return errors.New("item não encontrado")
}

// AddPayment adiciona um pagamento à comanda
func (c *Command) AddPayment(payment CommandPayment) error {
	if c.Status != CommandStatusOpen {
		return errors.New("não é possível adicionar pagamentos a uma comanda fechada")
	}

	c.Payments = append(c.Payments, payment)
	c.TotalRecebido += payment.ValorRecebido
	c.CalculateBalance()

	return nil
}

// RemovePayment remove um pagamento da comanda
func (c *Command) RemovePayment(paymentID uuid.UUID) error {
	if c.Status != CommandStatusOpen {
		return errors.New("não é possível remover pagamentos de uma comanda fechada")
	}

	for i, payment := range c.Payments {
		if payment.ID == paymentID {
			c.TotalRecebido -= payment.ValorRecebido
			c.Payments = append(c.Payments[:i], c.Payments[i+1:]...)
			c.CalculateBalance()
			return nil
		}
	}

	return errors.New("pagamento não encontrado")
}

// RecalculateTotals recalcula os totais da comanda
func (c *Command) RecalculateTotals() {
	c.Subtotal = 0
	for _, item := range c.Items {
		c.Subtotal += item.PrecoFinal
	}

	c.Total = c.Subtotal - c.Desconto
	if c.Total < 0 {
		c.Total = 0
	}

	c.CalculateBalance()
	c.AtualizadoEm = time.Now()
}

// CalculateBalance calcula troco ou saldo devedor
func (c *Command) CalculateBalance() {
	diferenca := c.TotalRecebido - c.Total

	if diferenca > 0 {
		// Cliente pagou a mais -> tem troco
		if c.DeixarTrocoGorjeta {
			c.Troco = 0
			c.SaldoDevedor = 0
		} else {
			c.Troco = diferenca
			c.SaldoDevedor = 0
		}
	} else if diferenca < 0 {
		// Cliente pagou a menos -> tem saldo devedor
		c.Troco = 0
		c.SaldoDevedor = -diferenca
	} else {
		// Pagamento exato
		c.Troco = 0
		c.SaldoDevedor = 0
	}
}

// Helpers
func formatMoney(value float64) string {
	return string(rune(int(value*100) / 100))
}
