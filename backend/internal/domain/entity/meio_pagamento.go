package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Erros do domínio de Meio de Pagamento
var (
	ErrMeioPagamentoNomeVazio        = errors.New("nome do meio de pagamento não pode ser vazio")
	ErrMeioPagamentoNomeMuitoLongo   = errors.New("nome deve ter no máximo 100 caracteres")
	ErrMeioPagamentoTipoInvalido     = errors.New("tipo de pagamento inválido")
	ErrMeioPagamentoTaxaInvalida     = errors.New("taxa deve estar entre 0 e 100")
	ErrMeioPagamentoTaxaFixaInvalida = errors.New("taxa fixa deve ser maior ou igual a zero")
	ErrMeioPagamentoDMaisInvalido    = errors.New("D+ deve ser maior ou igual a zero")
)

// TipoPagamento representa os tipos de pagamento aceitos
type TipoPagamento string

const (
	TipoPagamentoDinheiro      TipoPagamento = "DINHEIRO"
	TipoPagamentoPIX           TipoPagamento = "PIX"
	TipoPagamentoCredito       TipoPagamento = "CREDITO"
	TipoPagamentoDebito        TipoPagamento = "DEBITO"
	TipoPagamentoTransferencia TipoPagamento = "TRANSFERENCIA"
	TipoPagamentoBoleto        TipoPagamento = "BOLETO"
	TipoPagamentoOutro         TipoPagamento = "OUTRO"
)

// IsValid verifica se o tipo de pagamento é válido
func (t TipoPagamento) IsValid() bool {
	switch t {
	case TipoPagamentoDinheiro, TipoPagamentoPIX, TipoPagamentoCredito,
		TipoPagamentoDebito, TipoPagamentoTransferencia, TipoPagamentoBoleto,
		TipoPagamentoOutro:
		return true
	}
	return false
}

// DisplayName retorna o nome amigável do tipo
func (t TipoPagamento) DisplayName() string {
	switch t {
	case TipoPagamentoDinheiro:
		return "Dinheiro"
	case TipoPagamentoPIX:
		return "PIX"
	case TipoPagamentoCredito:
		return "Crédito"
	case TipoPagamentoDebito:
		return "Débito"
	case TipoPagamentoTransferencia:
		return "Transferência"
	case TipoPagamentoBoleto:
		return "Boleto"
	case TipoPagamentoOutro:
		return "Outro"
	default:
		return string(t)
	}
}

// RequiresBandeira verifica se o tipo requer bandeira
func (t TipoPagamento) RequiresBandeira() bool {
	return t == TipoPagamentoCredito || t == TipoPagamentoDebito
}

// DefaultDPlus retorna o D+ padrão para o tipo
func (t TipoPagamento) DefaultDPlus() int {
	switch t {
	case TipoPagamentoDinheiro, TipoPagamentoPIX:
		return 0
	case TipoPagamentoDebito:
		return 1
	case TipoPagamentoCredito:
		return 30
	case TipoPagamentoTransferencia:
		return 1
	default:
		return 0
	}
}

// MeioPagamento representa um meio de pagamento cadastrado
type MeioPagamento struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	Nome          string
	Tipo          TipoPagamento
	Bandeira      string          // Visa, Master, Elo, etc (opcional)
	Taxa          decimal.Decimal // Taxa percentual (0-100)
	TaxaFixa      decimal.Decimal // Taxa fixa em R$
	DMais         int             // Dias para compensação (D+0, D+1, D+30, etc)
	Icone         string          // Nome do ícone Material Icons
	Cor           string          // Cor hexadecimal
	OrdemExibicao int             // Ordem na UI
	Observacoes   string
	Ativo         bool
	CriadoEm      time.Time
	AtualizadoEm  time.Time
}

// NewMeioPagamento cria um novo meio de pagamento com validações
func NewMeioPagamento(tenantID uuid.UUID, nome string, tipo TipoPagamento) (*MeioPagamento, error) {
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return nil, ErrMeioPagamentoNomeVazio
	}
	if len(nome) > 100 {
		return nil, ErrMeioPagamentoNomeMuitoLongo
	}

	if !tipo.IsValid() {
		return nil, ErrMeioPagamentoTipoInvalido
	}

	return &MeioPagamento{
		ID:            uuid.New(),
		TenantID:      tenantID,
		Nome:          nome,
		Tipo:          tipo,
		Taxa:          decimal.Zero,
		TaxaFixa:      decimal.Zero,
		DMais:         tipo.DefaultDPlus(),
		OrdemExibicao: 0,
		Ativo:         true,
		CriadoEm:      time.Now(),
		AtualizadoEm:  time.Now(),
	}, nil
}

// SetTaxa define a taxa percentual com validação
func (m *MeioPagamento) SetTaxa(taxa decimal.Decimal) error {
	if taxa.LessThan(decimal.Zero) || taxa.GreaterThan(decimal.NewFromInt(100)) {
		return ErrMeioPagamentoTaxaInvalida
	}
	m.Taxa = taxa
	m.AtualizadoEm = time.Now()
	return nil
}

// SetTaxaFixa define a taxa fixa com validação
func (m *MeioPagamento) SetTaxaFixa(taxaFixa decimal.Decimal) error {
	if taxaFixa.LessThan(decimal.Zero) {
		return ErrMeioPagamentoTaxaFixaInvalida
	}
	m.TaxaFixa = taxaFixa
	m.AtualizadoEm = time.Now()
	return nil
}

// SetDMais define os dias para compensação
func (m *MeioPagamento) SetDMais(dMais int) error {
	if dMais < 0 {
		return ErrMeioPagamentoDMaisInvalido
	}
	m.DMais = dMais
	m.AtualizadoEm = time.Now()
	return nil
}

// SetBandeira define a bandeira do cartão
func (m *MeioPagamento) SetBandeira(bandeira string) {
	m.Bandeira = strings.TrimSpace(bandeira)
	m.AtualizadoEm = time.Now()
}

// Activate ativa o meio de pagamento
func (m *MeioPagamento) Activate() {
	m.Ativo = true
	m.AtualizadoEm = time.Now()
}

// Deactivate desativa o meio de pagamento
func (m *MeioPagamento) Deactivate() {
	m.Ativo = false
	m.AtualizadoEm = time.Now()
}

// CalculateNetValue calcula o valor líquido após taxas
func (m *MeioPagamento) CalculateNetValue(grossValue decimal.Decimal) decimal.Decimal {
	// valor_liquido = valor_bruto - (valor_bruto × taxa / 100) - taxa_fixa
	percentualDeduction := grossValue.Mul(m.Taxa).Div(decimal.NewFromInt(100))
	return grossValue.Sub(percentualDeduction).Sub(m.TaxaFixa)
}

// CalculateSettlementDate calcula a data de compensação baseada em D+
// Considera apenas dias úteis (pula sábados e domingos)
func (m *MeioPagamento) CalculateSettlementDate(transactionDate time.Time) time.Time {
	if m.DMais == 0 {
		return transactionDate
	}

	result := transactionDate
	daysAdded := 0

	for daysAdded < m.DMais {
		result = result.AddDate(0, 0, 1)
		weekday := result.Weekday()

		// Pula sábado e domingo
		if weekday != time.Saturday && weekday != time.Sunday {
			daysAdded++
		}
	}

	return result
}

// FormatDMais retorna o D+ formatado (ex: "D+30")
func (m *MeioPagamento) FormatDMais() string {
	return "D+" + string(rune('0'+m.DMais/10)) + string(rune('0'+m.DMais%10))
}
