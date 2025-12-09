package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CommissionPeriod representa um período de fechamento de comissões
type CommissionPeriod struct {
	ID               string
	TenantID         uuid.UUID
	UnitID           *string
	ReferenceMonth   string // formato: YYYY-MM
	ProfessionalID   *string
	TotalGross       decimal.Decimal
	TotalCommission  decimal.Decimal
	TotalAdvances    decimal.Decimal
	TotalAdjustments decimal.Decimal
	TotalNet         decimal.Decimal
	ItemsCount       int
	Status           string // ABERTO, PROCESSANDO, FECHADO, PAGO, CANCELADO
	PeriodStart      time.Time
	PeriodEnd        time.Time
	ClosedAt         *time.Time
	PaidAt           *time.Time
	ContaPagarID     *string
	ClosedBy         *string
	PaidBy           *string
	Notes            *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewCommissionPeriod cria um novo período de comissão
func NewCommissionPeriod(
	tenantID uuid.UUID,
	referenceMonth string,
	professionalID string,
	periodStart, periodEnd time.Time,
) (*CommissionPeriod, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if referenceMonth == "" {
		return nil, domain.ErrMesAnoObrigatorio
	}
	if professionalID == "" {
		return nil, domain.ErrProfissionalObrigatorio
	}
	if periodEnd.Before(periodStart) {
		return nil, domain.ErrDataFimAntesInicio
	}

	now := time.Now()
	return &CommissionPeriod{
		ID:               uuid.NewString(),
		TenantID:         tenantID,
		ReferenceMonth:   referenceMonth,
		ProfessionalID:   &professionalID,
		TotalGross:       decimal.Zero,
		TotalCommission:  decimal.Zero,
		TotalAdvances:    decimal.Zero,
		TotalAdjustments: decimal.Zero,
		TotalNet:         decimal.Zero,
		ItemsCount:       0,
		Status:           "ABERTO",
		PeriodStart:      periodStart,
		PeriodEnd:        periodEnd,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// Validate valida o período de comissão
func (p *CommissionPeriod) Validate() error {
	if p.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if p.ReferenceMonth == "" {
		return domain.ErrMesAnoObrigatorio
	}
	if p.PeriodEnd.Before(p.PeriodStart) {
		return domain.ErrDataFimAntesInicio
	}
	return nil
}

// CanClose verifica se o período pode ser fechado
func (p *CommissionPeriod) CanClose() bool {
	return p.Status == "ABERTO"
}

// CanPay verifica se o período pode ser pago
func (p *CommissionPeriod) CanPay() bool {
	return p.Status == "FECHADO"
}

// UpdateTotals atualiza os totais do período
func (p *CommissionPeriod) UpdateTotals(
	totalGross, totalCommission, totalAdvances decimal.Decimal,
	itemsCount int,
) {
	p.TotalGross = totalGross
	p.TotalCommission = totalCommission
	p.TotalAdvances = totalAdvances
	p.TotalNet = totalCommission.Sub(totalAdvances).Add(p.TotalAdjustments)
	p.ItemsCount = itemsCount
	p.UpdatedAt = time.Now()
}

// Close fecha o período
func (p *CommissionPeriod) Close(closedBy string) error {
	if !p.CanClose() {
		return domain.ErrPeriodoNaoPodeFechado
	}

	now := time.Now()
	p.Status = "FECHADO"
	p.ClosedAt = &now
	p.ClosedBy = &closedBy
	p.UpdatedAt = now
	return nil
}

// MarkAsPaid marca o período como pago
func (p *CommissionPeriod) MarkAsPaid(paidBy string) error {
	if !p.CanPay() {
		return domain.ErrPeriodoNaoPodePago
	}

	now := time.Now()
	p.Status = "PAGO"
	p.PaidAt = &now
	p.PaidBy = &paidBy
	p.UpdatedAt = now
	return nil
}

// Cancel cancela o período
func (p *CommissionPeriod) Cancel(notes string) error {
	if p.Status == "PAGO" {
		return domain.ErrPeriodoPagoNaoCancelavel
	}

	p.Status = "CANCELADO"
	p.Notes = &notes
	p.UpdatedAt = time.Now()
	return nil
}
