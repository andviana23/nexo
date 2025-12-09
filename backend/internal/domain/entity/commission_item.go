package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CommissionItem representa um item individual de comissão
type CommissionItem struct {
	ID               string
	TenantID         uuid.UUID
	UnitID           *string
	ProfessionalID   string
	CommandID        *string
	CommandItemID    *string
	AppointmentID    *string
	ServiceID        *string
	ServiceName      *string
	GrossValue       decimal.Decimal
	CommissionRate   decimal.Decimal
	CommissionType   string // PERCENTUAL, FIXO
	CommissionValue  decimal.Decimal
	CommissionSource string // SERVICO, PROFISSIONAL, REGRA, MANUAL
	RuleID           *string
	ReferenceDate    time.Time
	Description      *string
	Status           string // PENDENTE, PROCESSADO, PAGO, CANCELADO, ESTORNADO
	PeriodID         *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ProcessedAt      *time.Time
}

// NewCommissionItem cria um novo item de comissão
func NewCommissionItem(
	tenantID uuid.UUID,
	professionalID string,
	grossValue decimal.Decimal,
	commissionRate decimal.Decimal,
	commissionType string,
	commissionSource string,
	referenceDate time.Time,
) (*CommissionItem, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if professionalID == "" {
		return nil, domain.ErrProfissionalObrigatorio
	}
	if grossValue.IsZero() || grossValue.IsNegative() {
		return nil, domain.ErrValorDeveSerPositivo
	}
	if commissionType != "PERCENTUAL" && commissionType != "FIXO" {
		return nil, domain.ErrTipoComissaoInvalido
	}

	// Calcular valor da comissão
	var commissionValue decimal.Decimal
	if commissionType == "PERCENTUAL" {
		commissionValue = grossValue.Mul(commissionRate).Div(decimal.NewFromInt(100))
	} else {
		commissionValue = commissionRate
	}

	now := time.Now()
	return &CommissionItem{
		ID:               uuid.NewString(),
		TenantID:         tenantID,
		ProfessionalID:   professionalID,
		GrossValue:       grossValue,
		CommissionRate:   commissionRate,
		CommissionType:   commissionType,
		CommissionValue:  commissionValue.Round(2),
		CommissionSource: commissionSource,
		ReferenceDate:    referenceDate,
		Status:           "PENDENTE",
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// Validate valida o item de comissão
func (i *CommissionItem) Validate() error {
	if i.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if i.ProfessionalID == "" {
		return domain.ErrProfissionalObrigatorio
	}
	if i.GrossValue.IsZero() || i.GrossValue.IsNegative() {
		return domain.ErrValorDeveSerPositivo
	}
	if i.CommissionType != "PERCENTUAL" && i.CommissionType != "FIXO" {
		return domain.ErrTipoComissaoInvalido
	}
	return nil
}

// CanProcess verifica se o item pode ser processado
func (i *CommissionItem) CanProcess() bool {
	return i.Status == "PENDENTE"
}

// CanPay verifica se o item pode ser pago
func (i *CommissionItem) CanPay() bool {
	return i.Status == "PROCESSADO"
}

// CanCancel verifica se o item pode ser cancelado
func (i *CommissionItem) CanCancel() bool {
	return i.Status == "PENDENTE" || i.Status == "PROCESSADO"
}

// CanReverse verifica se o item pode ser estornado
func (i *CommissionItem) CanReverse() bool {
	return i.Status == "PAGO"
}

// Process processa o item
func (i *CommissionItem) Process(periodID string) error {
	if !i.CanProcess() {
		return domain.ErrItemNaoPodeProcessar
	}

	now := time.Now()
	i.Status = "PROCESSADO"
	i.PeriodID = &periodID
	i.ProcessedAt = &now
	i.UpdatedAt = now
	return nil
}

// MarkPaid marca o item como pago
func (i *CommissionItem) MarkPaid() error {
	if !i.CanPay() {
		return domain.ErrItemNaoPodePagar
	}

	i.Status = "PAGO"
	i.UpdatedAt = time.Now()
	return nil
}

// Cancel cancela o item
func (i *CommissionItem) Cancel() error {
	if !i.CanCancel() {
		return domain.ErrItemNaoPodeCancelar
	}

	i.Status = "CANCELADO"
	i.UpdatedAt = time.Now()
	return nil
}

// Reverse estorna o item
func (i *CommissionItem) Reverse() error {
	if !i.CanReverse() {
		return domain.ErrItemNaoPodeEstornar
	}

	i.Status = "ESTORNADO"
	i.UpdatedAt = time.Now()
	return nil
}

// Recalculate recalcula o valor da comissão
func (i *CommissionItem) Recalculate(newRate decimal.Decimal, ruleID *string) {
	if i.CommissionType == "PERCENTUAL" {
		i.CommissionValue = i.GrossValue.Mul(newRate).Div(decimal.NewFromInt(100)).Round(2)
	} else {
		i.CommissionValue = newRate.Round(2)
	}
	i.CommissionRate = newRate
	i.RuleID = ruleID
	i.UpdatedAt = time.Now()
}
