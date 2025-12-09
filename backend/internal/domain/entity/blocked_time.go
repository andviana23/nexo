package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidBlockedTimeID   = errors.New("ID do bloqueio inválido")
	ErrInvalidTimeRange       = errors.New("horário de fim deve ser posterior ao horário de início")
	ErrBlockedTimeReasonEmpty = errors.New("motivo do bloqueio é obrigatório")
	ErrTimeRangeOverlap       = errors.New("conflito com bloqueio existente")
)

// BlockedTime representa um bloqueio de horário na agenda
type BlockedTime struct {
	ID             string
	TenantID       uuid.UUID
	ProfessionalID string
	StartTime      time.Time
	EndTime        time.Time
	Reason         string
	IsRecurring    bool
	RecurrenceRule *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      *string
}

// NewBlockedTime cria um novo bloqueio validado
func NewBlockedTime(
	tenantID uuid.UUID,
	professionalID string,
	startTime time.Time,
	endTime time.Time,
	reason string,
) (*BlockedTime, error) {
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant_id inválido")
	}

	if _, err := uuid.Parse(professionalID); err != nil {
		return nil, errors.New("professional_id inválido")
	}

	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return nil, ErrInvalidTimeRange
	}

	if reason == "" || len(reason) < 3 {
		return nil, ErrBlockedTimeReasonEmpty
	}

	now := time.Now()
	return &BlockedTime{
		ID:             uuid.New().String(),
		TenantID:       tenantID,
		ProfessionalID: professionalID,
		StartTime:      startTime,
		EndTime:        endTime,
		Reason:         reason,
		IsRecurring:    false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Update atualiza os dados do bloqueio
func (bt *BlockedTime) Update(startTime, endTime time.Time, reason string) error {
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return ErrInvalidTimeRange
	}

	if reason == "" || len(reason) < 3 {
		return ErrBlockedTimeReasonEmpty
	}

	bt.StartTime = startTime
	bt.EndTime = endTime
	bt.Reason = reason
	bt.UpdatedAt = time.Now()

	return nil
}

// OverlapsWith verifica se este bloqueio sobrepõe outro período
func (bt *BlockedTime) OverlapsWith(otherStart, otherEnd time.Time) bool {
	return bt.StartTime.Before(otherEnd) && bt.EndTime.After(otherStart)
}
