package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MetaBarbeiro representa metas individuais de um barbeiro
type MetaBarbeiro struct {
	ID         string
	TenantID   uuid.UUID
	BarbeiroID string
	MesAno     valueobject.MesAno

	MetaServicosGerais valueobject.Money
	MetaServicosExtras valueobject.Money
	MetaProdutos       valueobject.Money

	// Realizado (calculado)
	RealizadoServicosGerais valueobject.Money
	RealizadoServicosExtras valueobject.Money
	RealizadoProdutos       valueobject.Money

	// Percentuais
	PercentualServicosGerais valueobject.Percentage
	PercentualServicosExtras valueobject.Percentage
	PercentualProdutos       valueobject.Percentage

	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewMetaBarbeiro cria uma nova meta para barbeiro
func NewMetaBarbeiro(
	tenantID uuid.UUID, barbeiroID string,
	mesAno valueobject.MesAno,
	metaServicosGerais, metaServicosExtras, metaProdutos valueobject.Money,
) (*MetaBarbeiro, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if barbeiroID == "" {
		return nil, domain.ErrInvalidID
	}
	if metaServicosGerais.IsNegative() || metaServicosExtras.IsNegative() || metaProdutos.IsNegative() {
		return nil, domain.ErrMetaNegativa
	}

	now := time.Now()
	return &MetaBarbeiro{
		ID:                       uuid.NewString(),
		TenantID:                 tenantID,
		BarbeiroID:               barbeiroID,
		MesAno:                   mesAno,
		MetaServicosGerais:       metaServicosGerais,
		MetaServicosExtras:       metaServicosExtras,
		MetaProdutos:             metaProdutos,
		RealizadoServicosGerais:  valueobject.Zero(),
		RealizadoServicosExtras:  valueobject.Zero(),
		RealizadoProdutos:        valueobject.Zero(),
		PercentualServicosGerais: valueobject.ZeroPercent(),
		PercentualServicosExtras: valueobject.ZeroPercent(),
		PercentualProdutos:       valueobject.ZeroPercent(),
		CriadoEm:                 now,
		AtualizadoEm:             now,
	}, nil
}

// CalcularProgresso calcula os percentuais realizados
func (m *MetaBarbeiro) CalcularProgresso(realizadoGerais, realizadoExtras, realizadoProdutos valueobject.Money) {
	m.RealizadoServicosGerais = realizadoGerais
	m.RealizadoServicosExtras = realizadoExtras
	m.RealizadoProdutos = realizadoProdutos

	// Calcula percentual de serviços gerais
	if m.MetaServicosGerais.IsPositive() {
		perc := realizadoGerais.Value().
			Div(m.MetaServicosGerais.Value()).
			Mul(decimal.NewFromInt(100))
		if perc.LessThanOrEqual(decimal.NewFromInt(200)) {
			m.PercentualServicosGerais = valueobject.NewPercentageUnsafe(perc)
		} else {
			m.PercentualServicosGerais, _ = valueobject.NewPercentage(decimal.NewFromInt(200))
		}
	}

	// Calcula percentual de serviços extras
	if m.MetaServicosExtras.IsPositive() {
		perc := realizadoExtras.Value().
			Div(m.MetaServicosExtras.Value()).
			Mul(decimal.NewFromInt(100))
		if perc.LessThanOrEqual(decimal.NewFromInt(200)) {
			m.PercentualServicosExtras = valueobject.NewPercentageUnsafe(perc)
		} else {
			m.PercentualServicosExtras, _ = valueobject.NewPercentage(decimal.NewFromInt(200))
		}
	}

	// Calcula percentual de produtos
	if m.MetaProdutos.IsPositive() {
		perc := realizadoProdutos.Value().
			Div(m.MetaProdutos.Value()).
			Mul(decimal.NewFromInt(100))
		if perc.LessThanOrEqual(decimal.NewFromInt(200)) {
			m.PercentualProdutos = valueobject.NewPercentageUnsafe(perc)
		} else {
			m.PercentualProdutos, _ = valueobject.NewPercentage(decimal.NewFromInt(200))
		}
	}

	m.AtualizadoEm = time.Now()
}

// Validate valida as regras de negócio
func (m *MetaBarbeiro) Validate() error {
	if m.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if m.BarbeiroID == "" {
		return domain.ErrInvalidID
	}
	if m.MesAno.String() == "" {
		return domain.ErrMesAnoRequired
	}
	if m.MetaServicosGerais.IsNegative() || m.MetaServicosExtras.IsNegative() || m.MetaProdutos.IsNegative() {
		return domain.ErrMetaNegativa
	}
	return nil
}

// MetaTotal retorna a meta total (soma de todas)
func (m *MetaBarbeiro) MetaTotal() valueobject.Money {
	return m.MetaServicosGerais.
		Add(m.MetaServicosExtras).
		Add(m.MetaProdutos)
}

// RealizadoTotal retorna o realizado total
func (m *MetaBarbeiro) RealizadoTotal() valueobject.Money {
	return m.RealizadoServicosGerais.
		Add(m.RealizadoServicosExtras).
		Add(m.RealizadoProdutos)
}

// PercentualGeral retorna o percentual geral de todas as metas
func (m *MetaBarbeiro) PercentualGeral() valueobject.Percentage {
	total := m.MetaTotal()
	realizado := m.RealizadoTotal()

	if total.IsPositive() {
		perc := realizado.Value().
			Div(total.Value()).
			Mul(decimal.NewFromInt(100))
		if perc.LessThanOrEqual(decimal.NewFromInt(200)) {
			return valueobject.NewPercentageUnsafe(perc)
		}
		p, _ := valueobject.NewPercentage(decimal.NewFromInt(200))
		return p
	}
	return valueobject.ZeroPercent()
}
