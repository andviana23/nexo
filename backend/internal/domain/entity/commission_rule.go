package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CommissionRule representa uma regra de comissão
type CommissionRule struct {
	ID              string
	TenantID        uuid.UUID
	UnitID          *string
	Name            string
	Description     *string
	Type            string // PERCENTUAL, FIXO
	DefaultRate     decimal.Decimal
	MinAmount       *decimal.Decimal
	MaxAmount       *decimal.Decimal
	CalculationBase *string // BRUTO, LIQUIDO
	EffectiveFrom   time.Time
	EffectiveTo     *time.Time
	Priority        *int
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CreatedBy       *string
}

// NewCommissionRule cria uma nova regra de comissão
func NewCommissionRule(
	tenantID uuid.UUID,
	name string,
	ruleType string,
	defaultRate decimal.Decimal,
) (*CommissionRule, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if name == "" {
		return nil, domain.ErrNomeObrigatorio
	}
	if ruleType != "PERCENTUAL" && ruleType != "FIXO" {
		return nil, domain.ErrTipoComissaoInvalido
	}
	if ruleType == "PERCENTUAL" && defaultRate.GreaterThan(decimal.NewFromInt(100)) {
		return nil, domain.ErrPercentualInvalido
	}

	now := time.Now()
	calcBase := "BRUTO"
	return &CommissionRule{
		ID:              uuid.NewString(),
		TenantID:        tenantID,
		Name:            name,
		Type:            ruleType,
		DefaultRate:     defaultRate,
		CalculationBase: &calcBase,
		EffectiveFrom:   now,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Validate valida a regra de comissão
func (r *CommissionRule) Validate() error {
	if r.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if r.Name == "" {
		return domain.ErrNomeObrigatorio
	}
	if r.Type != "PERCENTUAL" && r.Type != "FIXO" {
		return domain.ErrTipoComissaoInvalido
	}
	if r.Type == "PERCENTUAL" && r.DefaultRate.GreaterThan(decimal.NewFromInt(100)) {
		return domain.ErrPercentualInvalido
	}
	if r.EffectiveTo != nil && r.EffectiveTo.Before(r.EffectiveFrom) {
		return domain.ErrDataFimAntesInicio
	}
	if r.MinAmount != nil && r.MaxAmount != nil && r.MaxAmount.LessThan(*r.MinAmount) {
		return domain.ErrMaxMenorQueMin
	}
	return nil
}

// IsEffective verifica se a regra está em vigência
func (r *CommissionRule) IsEffective() bool {
	if !r.IsActive {
		return false
	}
	now := time.Now()
	if now.Before(r.EffectiveFrom) {
		return false
	}
	if r.EffectiveTo != nil && now.After(*r.EffectiveTo) {
		return false
	}
	return true
}

// CalculateCommission calcula a comissão com base no valor bruto
func (r *CommissionRule) CalculateCommission(grossValue decimal.Decimal) decimal.Decimal {
	var commission decimal.Decimal

	if r.Type == "PERCENTUAL" {
		commission = grossValue.Mul(r.DefaultRate).Div(decimal.NewFromInt(100))
	} else {
		commission = r.DefaultRate
	}

	// Aplicar limites
	if r.MinAmount != nil && commission.LessThan(*r.MinAmount) {
		commission = *r.MinAmount
	}
	if r.MaxAmount != nil && commission.GreaterThan(*r.MaxAmount) {
		commission = *r.MaxAmount
	}

	return commission.Round(2)
}

// Deactivate desativa a regra
func (r *CommissionRule) Deactivate() {
	r.IsActive = false
	r.UpdatedAt = time.Now()
}

// Activate ativa a regra
func (r *CommissionRule) Activate() {
	r.IsActive = true
	r.UpdatedAt = time.Now()
}
