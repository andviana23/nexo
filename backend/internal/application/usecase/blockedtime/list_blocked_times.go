package blockedtime

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
)

// ListBlockedTimesInput representa a entrada para listar bloqueios
type ListBlockedTimesInput struct {
	TenantID       string
	ProfessionalID *string
	StartDate      *time.Time
	EndDate        *time.Time
}

// ListBlockedTimesOutput representa a saída da listagem
type ListBlockedTimesOutput struct {
	BlockedTimes []*entity.BlockedTime
}

// ListBlockedTimesUseCase lista bloqueios de horário
type ListBlockedTimesUseCase struct {
	blockedTimeRepo repository.BlockedTimeRepository
}

// NewListBlockedTimesUseCase cria uma nova instância do use case
func NewListBlockedTimesUseCase(blockedTimeRepo repository.BlockedTimeRepository) *ListBlockedTimesUseCase {
	return &ListBlockedTimesUseCase{
		blockedTimeRepo: blockedTimeRepo,
	}
}

// Execute executa o use case
func (uc *ListBlockedTimesUseCase) Execute(ctx context.Context, input ListBlockedTimesInput) (*ListBlockedTimesOutput, error) {
	blockedTimes, err := uc.blockedTimeRepo.List(
		ctx,
		input.TenantID,
		input.ProfessionalID,
		input.StartDate,
		input.EndDate,
	)
	if err != nil {
		return nil, err
	}

	return &ListBlockedTimesOutput{BlockedTimes: blockedTimes}, nil
}
