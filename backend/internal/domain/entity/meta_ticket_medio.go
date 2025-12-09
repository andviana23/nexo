package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MetaTicketMedio representa a meta de ticket médio
type MetaTicketMedio struct {
	ID         string
	TenantID   uuid.UUID
	MesAno     valueobject.MesAno
	Tipo       valueobject.TipoMetaTicket // GERAL ou BARBEIRO
	BarbeiroID *string                    // Apenas se tipo = BARBEIRO

	MetaValor valueobject.Money

	// Realizado (calculado)
	TicketMedioRealizado valueobject.Money
	Percentual           valueobject.Percentage

	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewMetaTicketMedio cria uma nova meta de ticket médio
func NewMetaTicketMedio(
	tenantID uuid.UUID,
	mesAno valueobject.MesAno,
	tipo valueobject.TipoMetaTicket,
	metaValor valueobject.Money,
	barbeiroID *string,
) (*MetaTicketMedio, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if !tipo.IsValid() {
		return nil, domain.ErrMetaInvalida
	}
	if metaValor.IsNegative() || metaValor.IsZero() {
		return nil, domain.ErrMetaNegativa
	}

	// Validação: se tipo = BARBEIRO, barbeiroID é obrigatório
	if tipo == valueobject.TipoMetaTicketBarbeiro && (barbeiroID == nil || *barbeiroID == "") {
		return nil, domain.ErrInvalidID
	}

	// Validação: se tipo = GERAL, barbeiroID deve ser nil
	if tipo == valueobject.TipoMetaTicketGeral && barbeiroID != nil {
		return nil, domain.ErrMetaInvalida
	}

	now := time.Now()
	return &MetaTicketMedio{
		ID:                   uuid.NewString(),
		TenantID:             tenantID,
		MesAno:               mesAno,
		Tipo:                 tipo,
		BarbeiroID:           barbeiroID,
		MetaValor:            metaValor,
		TicketMedioRealizado: valueobject.Zero(),
		Percentual:           valueobject.ZeroPercent(),
		CriadoEm:             now,
		AtualizadoEm:         now,
	}, nil
}

// CalcularProgresso calcula o percentual de atingimento
func (m *MetaTicketMedio) CalcularProgresso(ticketRealizado valueobject.Money) {
	m.TicketMedioRealizado = ticketRealizado

	if m.MetaValor.IsPositive() {
		perc := ticketRealizado.Value().
			Div(m.MetaValor.Value()).
			Mul(decimal.NewFromInt(100))

		// Limita a 200%
		if perc.GreaterThan(decimal.NewFromInt(200)) {
			perc = decimal.NewFromInt(200)
		}
		m.Percentual = valueobject.NewPercentageUnsafe(perc)
	}

	m.AtualizadoEm = time.Now()
}

// Validate valida as regras de negócio
func (m *MetaTicketMedio) Validate() error {
	if m.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if m.MesAno.String() == "" {
		return domain.ErrMesAnoRequired
	}
	if !m.Tipo.IsValid() {
		return domain.ErrMetaInvalida
	}
	if m.MetaValor.IsNegative() || m.MetaValor.IsZero() {
		return domain.ErrMetaNegativa
	}

	// Validação de consistência tipo vs barbeiro
	if m.Tipo == valueobject.TipoMetaTicketBarbeiro && (m.BarbeiroID == nil || *m.BarbeiroID == "") {
		return domain.ErrInvalidID
	}
	if m.Tipo == valueobject.TipoMetaTicketGeral && m.BarbeiroID != nil {
		return domain.ErrMetaInvalida
	}

	return nil
}

// Atingiu verifica se atingiu a meta (>= 100%)
func (m *MetaTicketMedio) Atingiu() bool {
	cem, _ := valueobject.NewPercentage(decimal.NewFromInt(100))
	return m.Percentual.GreaterThan(cem) || m.Percentual.Equals(cem)
}

// ForaDoAlvo verifica se está abaixo de 80%
func (m *MetaTicketMedio) ForaDoAlvo() bool {
	oitenta, _ := valueobject.NewPercentage(decimal.NewFromInt(80))
	return m.Percentual.LessThan(oitenta)
}
