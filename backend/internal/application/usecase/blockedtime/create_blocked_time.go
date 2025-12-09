package blockedtime

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
	"github.com/google/uuid"
)

// CreateBlockedTimeInput representa a entrada para criar um bloqueio
type CreateBlockedTimeInput struct {
	TenantID       string
	ProfessionalID string
	StartTime      time.Time
	EndTime        time.Time
	Reason         string
	UserID         *string
}

// CreateBlockedTimeOutput representa a saída da criação
type CreateBlockedTimeOutput struct {
	BlockedTime *entity.BlockedTime
}

// CreateBlockedTimeUseCase cria um novo bloqueio de horário
type CreateBlockedTimeUseCase struct {
	blockedTimeRepo repository.BlockedTimeRepository
}

// NewCreateBlockedTimeUseCase cria uma nova instância do use case
func NewCreateBlockedTimeUseCase(blockedTimeRepo repository.BlockedTimeRepository) *CreateBlockedTimeUseCase {
	return &CreateBlockedTimeUseCase{
		blockedTimeRepo: blockedTimeRepo,
	}
}

// Execute executa o use case
func (uc *CreateBlockedTimeUseCase) Execute(ctx context.Context, input CreateBlockedTimeInput) (*CreateBlockedTimeOutput, error) {
	// Converter tenant_id de string para uuid.UUID
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}

	// Cria a entidade validada
	blockedTime, err := entity.NewBlockedTime(
		tenantUUID,
		input.ProfessionalID,
		input.StartTime,
		input.EndTime,
		input.Reason,
	)
	if err != nil {
		return nil, err
	}

	// Define quem criou
	if input.UserID != nil {
		blockedTime.CreatedBy = input.UserID
	}

	// Verifica conflito com bloqueios existentes
	hasConflict, err := uc.blockedTimeRepo.CheckConflict(
		ctx,
		input.TenantID,
		input.ProfessionalID,
		input.StartTime,
		input.EndTime,
		nil, // Sem exclusão (é novo)
	)
	if err != nil {
		return nil, err
	}

	if hasConflict {
		return nil, entity.ErrTimeRangeOverlap
	}

	// Persiste
	created, err := uc.blockedTimeRepo.Create(ctx, blockedTime)
	if err != nil {
		return nil, err
	}

	return &CreateBlockedTimeOutput{BlockedTime: created}, nil
}
