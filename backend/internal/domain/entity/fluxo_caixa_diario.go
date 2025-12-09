package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

// FluxoCaixaDiario representa o fluxo de caixa compensado de um dia
type FluxoCaixaDiario struct {
	ID       string
	TenantID uuid.UUID
	Data     time.Time

	SaldoInicial        valueobject.Money
	EntradasConfirmadas valueobject.Money
	EntradasPrevistas   valueobject.Money
	SaidasPagas         valueobject.Money
	SaidasPrevistas     valueobject.Money
	SaldoFinal          valueobject.Money

	ProcessadoEm time.Time
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewFluxoCaixaDiario cria um novo fluxo de caixa diário
func NewFluxoCaixaDiario(tenantID uuid.UUID, data time.Time) (*FluxoCaixaDiario, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}

	now := time.Now()
	return &FluxoCaixaDiario{
		ID:                  uuid.NewString(),
		TenantID:            tenantID,
		Data:                data,
		SaldoInicial:        valueobject.Zero(),
		EntradasConfirmadas: valueobject.Zero(),
		EntradasPrevistas:   valueobject.Zero(),
		SaidasPagas:         valueobject.Zero(),
		SaidasPrevistas:     valueobject.Zero(),
		SaldoFinal:          valueobject.Zero(),
		CriadoEm:            now,
		AtualizadoEm:        now,
	}, nil
}

// Calcular recalcula o saldo final do fluxo de caixa
func (f *FluxoCaixaDiario) Calcular() {
	f.SaldoFinal = f.SaldoInicial.
		Add(f.EntradasConfirmadas).
		Add(f.EntradasPrevistas).
		Sub(f.SaidasPagas).
		Sub(f.SaidasPrevistas)

	f.ProcessadoEm = time.Now()
	f.AtualizadoEm = time.Now()
}

// Validate valida as regras de negócio
func (f *FluxoCaixaDiario) Validate() error {
	if f.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if f.Data.IsZero() {
		return domain.ErrDataVencimentoInvalida
	}
	// Entradas e saídas não podem ser negativas
	if f.EntradasConfirmadas.IsNegative() || f.EntradasPrevistas.IsNegative() {
		return domain.ErrValorNegativo
	}
	if f.SaidasPagas.IsNegative() || f.SaidasPrevistas.IsNegative() {
		return domain.ErrValorNegativo
	}
	return nil
}

// SetSaldoInicial define o saldo inicial
func (f *FluxoCaixaDiario) SetSaldoInicial(saldo valueobject.Money) {
	f.SaldoInicial = saldo
	f.AtualizadoEm = time.Now()
}

// AddEntradaConfirmada adiciona uma entrada confirmada
func (f *FluxoCaixaDiario) AddEntradaConfirmada(valor valueobject.Money) {
	f.EntradasConfirmadas = f.EntradasConfirmadas.Add(valor)
	f.AtualizadoEm = time.Now()
}

// AddEntradaPrevista adiciona uma entrada prevista
func (f *FluxoCaixaDiario) AddEntradaPrevista(valor valueobject.Money) {
	f.EntradasPrevistas = f.EntradasPrevistas.Add(valor)
	f.AtualizadoEm = time.Now()
}

// AddSaidaPaga adiciona uma saída paga
func (f *FluxoCaixaDiario) AddSaidaPaga(valor valueobject.Money) {
	f.SaidasPagas = f.SaidasPagas.Add(valor)
	f.AtualizadoEm = time.Now()
}

// AddSaidaPrevista adiciona uma saída prevista
func (f *FluxoCaixaDiario) AddSaidaPrevista(valor valueobject.Money) {
	f.SaidasPrevistas = f.SaidasPrevistas.Add(valor)
	f.AtualizadoEm = time.Now()
}
