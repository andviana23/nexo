package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

// ContaPagar representa uma conta a pagar (despesa)
type ContaPagar struct {
	ID       string
	TenantID uuid.UUID

	Descricao   string
	CategoriaID string
	Fornecedor  string
	Valor       valueobject.Money

	Tipo          valueobject.TipoCusto // FIXO ou VARIAVEL
	Recorrente    bool
	Periodicidade string // MENSAL, TRIMESTRAL, ANUAL

	DataVencimento time.Time
	DataPagamento  *time.Time
	Status         valueobject.StatusConta // PENDENTE, PAGO, CANCELADO, ATRASADO

	ComprovanteURL string
	PixCode        string
	Observacoes    string

	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewContaPagar cria uma nova conta a pagar
func NewContaPagar(
	tenantID uuid.UUID, descricao, categoriaID, fornecedor string,
	valor valueobject.Money,
	tipo valueobject.TipoCusto,
	dataVencimento time.Time,
	recorrente bool,
	periodicidade string,
) (*ContaPagar, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if descricao == "" {
		return nil, domain.ErrInvalidID
	}
	if valor.IsNegative() {
		return nil, domain.ErrValorNegativo
	}
	if !tipo.IsValid() {
		return nil, domain.ErrStatusInvalido
	}
	if dataVencimento.IsZero() {
		return nil, domain.ErrDataVencimentoInvalida
	}

	now := time.Now()
	return &ContaPagar{
		ID:             uuid.NewString(),
		TenantID:       tenantID,
		Descricao:      descricao,
		CategoriaID:    categoriaID,
		Fornecedor:     fornecedor,
		Valor:          valor,
		Tipo:           tipo,
		Recorrente:     recorrente,
		Periodicidade:  periodicidade,
		DataVencimento: dataVencimento,
		Status:         valueobject.StatusContaPendente,
		CriadoEm:       now,
		AtualizadoEm:   now,
	}, nil
}

// MarcarComoPago marca a conta como paga
func (c *ContaPagar) MarcarComoPago(dataPagamento time.Time, comprovante string) error {
	if c.Status == valueobject.StatusContaPago {
		return domain.ErrContaJaPaga
	}
	if c.Status == valueobject.StatusContaCancelado {
		return domain.ErrContaCancelada
	}

	c.DataPagamento = &dataPagamento
	c.ComprovanteURL = comprovante
	c.Status = valueobject.StatusContaPago
	c.AtualizadoEm = time.Now()
	return nil
}

// Cancelar cancela a conta
func (c *ContaPagar) Cancelar() error {
	if c.Status == valueobject.StatusContaPago {
		return domain.ErrContaJaPaga
	}
	c.Status = valueobject.StatusContaCancelado
	c.AtualizadoEm = time.Now()
	return nil
}

// VerificarAtraso verifica se está atrasada e atualiza o status
func (c *ContaPagar) VerificarAtraso() {
	if c.Status == valueobject.StatusContaPendente && time.Now().After(c.DataVencimento) {
		c.Status = valueobject.StatusContaAtrasado
		c.AtualizadoEm = time.Now()
	}
}

// DiasAteVencimento calcula dias até vencimento (negativo se atrasado)
func (c *ContaPagar) DiasAteVencimento() int {
	diff := c.DataVencimento.Sub(time.Now())
	return int(diff.Hours() / 24)
}

// EstaAtrasada verifica se está atrasada
func (c *ContaPagar) EstaAtrasada() bool {
	return c.Status == valueobject.StatusContaAtrasado ||
		(c.Status == valueobject.StatusContaPendente && time.Now().After(c.DataVencimento))
}

// VenceEmBreve verifica se vence em até N dias
func (c *ContaPagar) VenceEmBreve(dias int) bool {
	if c.Status != valueobject.StatusContaPendente {
		return false
	}
	diasAte := c.DiasAteVencimento()
	return diasAte >= 0 && diasAte <= dias
}

// Validate valida as regras de negócio
func (c *ContaPagar) Validate() error {
	if c.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if c.Descricao == "" {
		return domain.ErrInvalidID
	}
	if c.Valor.IsNegative() {
		return domain.ErrValorNegativo
	}
	if !c.Tipo.IsValid() {
		return domain.ErrStatusInvalido
	}
	if !c.Status.IsValid() {
		return domain.ErrStatusInvalido
	}
	if c.DataVencimento.IsZero() {
		return domain.ErrDataVencimentoInvalida
	}
	if c.Recorrente && c.Periodicidade == "" {
		return domain.ErrMetaInvalida
	}
	return nil
}

// SetPixCode define o código PIX para pagamento
func (c *ContaPagar) SetPixCode(pixCode string) {
	c.PixCode = pixCode
	c.AtualizadoEm = time.Now()
}

// AddObservacao adiciona uma observação
func (c *ContaPagar) AddObservacao(obs string) {
	if c.Observacoes == "" {
		c.Observacoes = obs
	} else {
		c.Observacoes = c.Observacoes + "; " + obs
	}
	c.AtualizadoEm = time.Now()
}
