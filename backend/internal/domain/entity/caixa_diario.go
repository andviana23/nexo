package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// StatusCaixa representa os possíveis status de um caixa diário
type StatusCaixa string

const (
	StatusCaixaAberto  StatusCaixa = "ABERTO"
	StatusCaixaFechado StatusCaixa = "FECHADO"
)

// ValidarStatusCaixa verifica se o status é válido
func ValidarStatusCaixa(s string) bool {
	switch StatusCaixa(s) {
	case StatusCaixaAberto, StatusCaixaFechado:
		return true
	}
	return false
}

// CaixaDiario representa o controle operacional da gaveta de dinheiro
type CaixaDiario struct {
	ID                       uuid.UUID
	TenantID                 uuid.UUID
	UsuarioAberturaID        uuid.UUID
	UsuarioFechamentoID      *uuid.UUID
	DataAbertura             time.Time
	DataFechamento           *time.Time
	SaldoInicial             decimal.Decimal
	TotalEntradas            decimal.Decimal
	TotalSaidas              decimal.Decimal
	TotalSangrias            decimal.Decimal
	TotalReforcos            decimal.Decimal
	SaldoEsperado            decimal.Decimal // Calculado: inicial + entradas - sangrias + reforços
	SaldoReal                *decimal.Decimal
	Divergencia              *decimal.Decimal
	Status                   StatusCaixa
	JustificativaDivergencia *string
	CreatedAt                time.Time
	UpdatedAt                time.Time

	// Relacionamentos (carregados quando necessário)
	Operacoes             []OperacaoCaixa
	UsuarioAberturaNome   string
	UsuarioFechamentoNome string
}

// NewCaixaDiario cria um novo caixa diário para abertura
func NewCaixaDiario(tenantID, usuarioID uuid.UUID, saldoInicial decimal.Decimal) (*CaixaDiario, error) {
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant_id é obrigatório")
	}
	if usuarioID == uuid.Nil {
		return nil, errors.New("usuario_abertura_id é obrigatório")
	}
	if saldoInicial.IsNegative() {
		return nil, errors.New("saldo inicial não pode ser negativo")
	}

	now := time.Now()
	return &CaixaDiario{
		ID:                uuid.New(),
		TenantID:          tenantID,
		UsuarioAberturaID: usuarioID,
		DataAbertura:      now,
		SaldoInicial:      saldoInicial,
		TotalEntradas:     decimal.Zero,
		TotalSaidas:       decimal.Zero,
		TotalSangrias:     decimal.Zero,
		TotalReforcos:     decimal.Zero,
		SaldoEsperado:     saldoInicial, // Inicialmente igual ao saldo inicial
		Status:            StatusCaixaAberto,
		CreatedAt:         now,
		UpdatedAt:         now,
		Operacoes:         []OperacaoCaixa{},
	}, nil
}

// CalcularSaldoEsperado recalcula o saldo esperado
func (c *CaixaDiario) CalcularSaldoEsperado() decimal.Decimal {
	// Saldo Esperado = Inicial + Entradas - Sangrias + Reforços
	c.SaldoEsperado = c.SaldoInicial.
		Add(c.TotalEntradas).
		Sub(c.TotalSangrias).
		Add(c.TotalReforcos)
	return c.SaldoEsperado
}

// CanFechar verifica se o caixa pode ser fechado
func (c *CaixaDiario) CanFechar() error {
	if c.Status != StatusCaixaAberto {
		return errors.New("caixa já está fechado")
	}
	return nil
}

// Fechar fecha o caixa com o saldo real informado
func (c *CaixaDiario) Fechar(usuarioID uuid.UUID, saldoReal decimal.Decimal, justificativa *string) error {
	if err := c.CanFechar(); err != nil {
		return err
	}

	// Calcular saldo esperado atualizado
	c.CalcularSaldoEsperado()

	// Calcular divergência
	divergencia := saldoReal.Sub(c.SaldoEsperado)

	// Validar justificativa obrigatória se divergência > R$ 5,00
	tolerancia := decimal.NewFromFloat(5.00)
	if divergencia.Abs().GreaterThan(tolerancia) {
		if justificativa == nil || *justificativa == "" {
			return errors.New("justificativa obrigatória para divergência maior que R$ 5,00")
		}
	}

	now := time.Now()
	c.Status = StatusCaixaFechado
	c.UsuarioFechamentoID = &usuarioID
	c.DataFechamento = &now
	c.SaldoReal = &saldoReal
	c.Divergencia = &divergencia
	c.JustificativaDivergencia = justificativa
	c.UpdatedAt = now

	return nil
}

// RegistrarSangria atualiza os totais após uma sangria
func (c *CaixaDiario) RegistrarSangria(valor decimal.Decimal) error {
	if c.Status != StatusCaixaAberto {
		return errors.New("não é possível registrar sangria em caixa fechado")
	}
	if valor.IsNegative() || valor.IsZero() {
		return errors.New("valor da sangria deve ser positivo")
	}

	c.TotalSangrias = c.TotalSangrias.Add(valor)
	c.CalcularSaldoEsperado()
	c.UpdatedAt = time.Now()

	return nil
}

// RegistrarReforco atualiza os totais após um reforço
func (c *CaixaDiario) RegistrarReforco(valor decimal.Decimal) error {
	if c.Status != StatusCaixaAberto {
		return errors.New("não é possível registrar reforço em caixa fechado")
	}
	if valor.IsNegative() || valor.IsZero() {
		return errors.New("valor do reforço deve ser positivo")
	}

	c.TotalReforcos = c.TotalReforcos.Add(valor)
	c.CalcularSaldoEsperado()
	c.UpdatedAt = time.Now()

	return nil
}

// RegistrarEntrada atualiza os totais após uma venda em dinheiro
func (c *CaixaDiario) RegistrarEntrada(valor decimal.Decimal) error {
	if c.Status != StatusCaixaAberto {
		return errors.New("não é possível registrar entrada em caixa fechado")
	}
	if valor.IsNegative() || valor.IsZero() {
		return errors.New("valor da entrada deve ser positivo")
	}

	c.TotalEntradas = c.TotalEntradas.Add(valor)
	c.CalcularSaldoEsperado()
	c.UpdatedAt = time.Now()

	return nil
}

// TemDivergencia verifica se houve divergência no fechamento
func (c *CaixaDiario) TemDivergencia() bool {
	if c.Divergencia == nil {
		return false
	}
	return !c.Divergencia.IsZero()
}

// DivergenciaPositiva verifica se sobrou dinheiro (saldo real > esperado)
func (c *CaixaDiario) DivergenciaPositiva() bool {
	if c.Divergencia == nil {
		return false
	}
	return c.Divergencia.IsPositive()
}

// DivergenciaNegativa verifica se faltou dinheiro (saldo real < esperado)
func (c *CaixaDiario) DivergenciaNegativa() bool {
	if c.Divergencia == nil {
		return false
	}
	return c.Divergencia.IsNegative()
}
