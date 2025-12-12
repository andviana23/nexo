package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

// ContaReceber representa uma conta a receber (receita futura)
type ContaReceber struct {
	ID       string
	TenantID uuid.UUID

	Origem          string  // ASSINATURA, SERVICO, OUTRO
	AssinaturaID    *string // Apenas se origem = ASSINATURA (tabela antiga assinaturas)
	DescricaoOrigem string

	Valor       valueobject.Money
	ValorPago   valueobject.Money
	ValorAberto valueobject.Money

	DataVencimento  time.Time
	DataRecebimento *time.Time
	Status          valueobject.StatusConta // PENDENTE, CONFIRMADO, RECEBIDO, ESTORNADO, CANCELADO, ATRASADO

	Observacoes string

	CriadoEm     time.Time
	AtualizadoEm time.Time

	// Vínculos explícitos com comanda/pagamento (Fase 1)
	CommandID        *string // FK para commands.id
	CommandPaymentID *string // FK para command_payments.id

	// Campos de integração Asaas (Migration 041)
	SubscriptionID *string    // FK para subscriptions (nova tabela)
	AsaasPaymentID *string    // ID da cobrança no Asaas (idempotência)
	CompetenciaMes *string    // Mês de competência YYYY-MM para DRE
	ConfirmedAt    *time.Time // Timestamp de confirmação (CONFIRMED)
	ReceivedAt     *time.Time // Timestamp de recebimento (RECEIVED)
}

// NewContaReceber cria uma nova conta a receber
func NewContaReceber(
	tenantID uuid.UUID, origem string,
	assinaturaID *string,
	descricaoOrigem string,
	valor valueobject.Money,
	dataVencimento time.Time,
) (*ContaReceber, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if origem == "" {
		return nil, domain.ErrInvalidID
	}
	if valor.IsNegative() {
		return nil, domain.ErrValorNegativo
	}
	if dataVencimento.IsZero() {
		return nil, domain.ErrDataVencimentoInvalida
	}

	now := time.Now()
	return &ContaReceber{
		ID:              uuid.NewString(),
		TenantID:        tenantID,
		Origem:          origem,
		AssinaturaID:    assinaturaID,
		DescricaoOrigem: descricaoOrigem,
		Valor:           valor,
		ValorPago:       valueobject.Zero(),
		ValorAberto:     valor,
		DataVencimento:  dataVencimento,
		Status:          valueobject.StatusContaPendente,
		CriadoEm:        now,
		AtualizadoEm:    now,
	}, nil
}

// MarcarComoRecebido marca a conta como totalmente recebida
func (c *ContaReceber) MarcarComoRecebido(dataRecebimento time.Time) error {
	if c.Status == valueobject.StatusContaRecebido {
		return domain.ErrContaJaPaga
	}
	if c.Status == valueobject.StatusContaCancelado {
		return domain.ErrContaCancelada
	}

	c.DataRecebimento = &dataRecebimento
	c.ValorPago = c.Valor
	c.ValorAberto = valueobject.Zero()
	c.Status = valueobject.StatusContaRecebido
	c.AtualizadoEm = time.Now()
	return nil
}

// RegistrarPagamentoParcial registra um pagamento parcial
func (c *ContaReceber) RegistrarPagamentoParcial(valorPago valueobject.Money) error {
	if c.Status == valueobject.StatusContaRecebido {
		return domain.ErrContaJaPaga
	}
	if c.Status == valueobject.StatusContaCancelado {
		return domain.ErrContaCancelada
	}
	if valorPago.IsNegative() || valorPago.IsZero() {
		return domain.ErrValorInvalido
	}

	c.ValorPago = c.ValorPago.Add(valorPago)
	c.ValorAberto = c.Valor.Sub(c.ValorPago)

	// Se pagou tudo, marca como pago
	if c.ValorAberto.IsZero() || c.ValorAberto.IsNegative() {
		c.ValorAberto = valueobject.Zero()
		c.Status = valueobject.StatusContaRecebido
		now := time.Now()
		c.DataRecebimento = &now
	}

	c.AtualizadoEm = time.Now()
	return nil
}

// Cancelar cancela a conta
func (c *ContaReceber) Cancelar() error {
	if c.Status == valueobject.StatusContaRecebido {
		return domain.ErrContaJaPaga
	}
	c.Status = valueobject.StatusContaCancelado
	c.AtualizadoEm = time.Now()
	return nil
}

// VerificarAtraso verifica se está atrasada
func (c *ContaReceber) VerificarAtraso() {
	if c.Status == valueobject.StatusContaPendente && time.Now().After(c.DataVencimento) {
		c.Status = valueobject.StatusContaAtrasado
		c.AtualizadoEm = time.Now()
	}
}

// DiasAteVencimento calcula dias até vencimento
func (c *ContaReceber) DiasAteVencimento() int {
	diff := c.DataVencimento.Sub(time.Now())
	return int(diff.Hours() / 24)
}

// EstaAtrasada verifica se está atrasada
func (c *ContaReceber) EstaAtrasada() bool {
	return c.Status == valueobject.StatusContaAtrasado ||
		(c.Status == valueobject.StatusContaPendente && time.Now().After(c.DataVencimento))
}

// VenceEmBreve verifica se vence em breve
func (c *ContaReceber) VenceEmBreve(dias int) bool {
	if c.Status != valueobject.StatusContaPendente {
		return false
	}
	diasAte := c.DiasAteVencimento()
	return diasAte >= 0 && diasAte <= dias
}

// PercentualRecebido retorna o percentual recebido
func (c *ContaReceber) PercentualRecebido() valueobject.Percentage {
	if c.Valor.IsZero() || c.Valor.IsNegative() {
		return valueobject.ZeroPercent()
	}
	perc := c.ValorPago.Value().
		Div(c.Valor.Value()).
		Mul(valueobject.HundredPercent().Value())
	return valueobject.NewPercentageUnsafe(perc)
}

// Validate valida as regras de negócio
func (c *ContaReceber) Validate() error {
	if c.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if c.Origem == "" {
		return domain.ErrInvalidID
	}
	if c.Valor.IsNegative() {
		return domain.ErrValorNegativo
	}
	if c.ValorPago.IsNegative() {
		return domain.ErrValorNegativo
	}
	if !c.Status.IsValid() {
		return domain.ErrStatusInvalido
	}
	if c.DataVencimento.IsZero() {
		return domain.ErrDataVencimentoInvalida
	}
	// Validar: se origem = ASSINATURA, assinaturaID é obrigatório
	if c.Origem == "ASSINATURA" {
		assinaturaAntigaOk := c.AssinaturaID != nil && *c.AssinaturaID != ""
		subscriptionOk := c.SubscriptionID != nil && *c.SubscriptionID != ""
		if !assinaturaAntigaOk && !subscriptionOk {
			return domain.ErrInvalidID
		}
	}
	return nil
}

// AddObservacao adiciona uma observação
func (c *ContaReceber) AddObservacao(obs string) {
	if c.Observacoes == "" {
		c.Observacoes = obs
	} else {
		c.Observacoes = c.Observacoes + "; " + obs
	}
	c.AtualizadoEm = time.Now()
}
