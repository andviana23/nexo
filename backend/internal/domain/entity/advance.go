package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Advance representa um adiantamento de comissão
type Advance struct {
	ID                string
	TenantID          uuid.UUID
	UnitID            *string
	ProfessionalID    string
	Amount            decimal.Decimal
	RequestDate       time.Time
	Reason            *string
	Status            string // PENDING, APPROVED, REJECTED, DEDUCTED, CANCELLED
	ApprovedAt        *time.Time
	ApprovedBy        *string
	RejectedAt        *time.Time
	RejectedBy        *string
	RejectionReason   *string
	DeductedAt        *time.Time
	DeductionPeriodID *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	CreatedBy         *string
}

// NewAdvance cria um novo adiantamento
func NewAdvance(
	tenantID uuid.UUID,
	professionalID string,
	amount decimal.Decimal,
) (*Advance, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if professionalID == "" {
		return nil, domain.ErrProfissionalObrigatorio
	}
	if amount.IsZero() || amount.IsNegative() {
		return nil, domain.ErrValorDeveSerPositivo
	}

	now := time.Now()
	return &Advance{
		ID:             uuid.NewString(),
		TenantID:       tenantID,
		ProfessionalID: professionalID,
		Amount:         amount,
		RequestDate:    now,
		Status:         "PENDING",
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Validate valida o adiantamento
func (a *Advance) Validate() error {
	if a.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if a.ProfessionalID == "" {
		return domain.ErrProfissionalObrigatorio
	}
	if a.Amount.IsZero() || a.Amount.IsNegative() {
		return domain.ErrValorDeveSerPositivo
	}
	return nil
}

// CanApprove verifica se o adiantamento pode ser aprovado
func (a *Advance) CanApprove() bool {
	return a.Status == "PENDING"
}

// CanReject verifica se o adiantamento pode ser rejeitado
func (a *Advance) CanReject() bool {
	return a.Status == "PENDING"
}

// CanDeduct verifica se o adiantamento pode ser deduzido
func (a *Advance) CanDeduct() bool {
	return a.Status == "APPROVED"
}

// CanCancel verifica se o adiantamento pode ser cancelado
func (a *Advance) CanCancel() bool {
	return a.Status == "PENDING" || a.Status == "APPROVED"
}

// Approve aprova o adiantamento
func (a *Advance) Approve(approvedBy string) error {
	if !a.CanApprove() {
		return domain.ErrAdiantamentoNaoPodeAprovar
	}

	now := time.Now()
	a.Status = "APPROVED"
	a.ApprovedAt = &now
	a.ApprovedBy = &approvedBy
	a.UpdatedAt = now
	return nil
}

// Reject rejeita o adiantamento
func (a *Advance) Reject(rejectedBy string, reason string) error {
	if !a.CanReject() {
		return domain.ErrAdiantamentoNaoPodeRejeitar
	}
	if reason == "" {
		return domain.ErrMotivoRejeicaoObrigatorio
	}

	now := time.Now()
	a.Status = "REJECTED"
	a.RejectedAt = &now
	a.RejectedBy = &rejectedBy
	a.RejectionReason = &reason
	a.UpdatedAt = now
	return nil
}

// MarkDeducted marca o adiantamento como deduzido
func (a *Advance) MarkDeducted(periodID string) error {
	if !a.CanDeduct() {
		return domain.ErrAdiantamentoNaoPodeDeduzir
	}

	now := time.Now()
	a.Status = "DEDUCTED"
	a.DeductedAt = &now
	a.DeductionPeriodID = &periodID
	a.UpdatedAt = now
	return nil
}

// Cancel cancela o adiantamento
func (a *Advance) Cancel() error {
	if !a.CanCancel() {
		return domain.ErrAdiantamentoNaoPodeCancelar
	}

	a.Status = "CANCELLED"
	a.UpdatedAt = time.Now()
	return nil
}
