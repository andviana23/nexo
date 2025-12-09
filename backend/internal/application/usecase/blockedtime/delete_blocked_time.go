package blockedtime

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
)

// DeleteBlockedTimeInput representa a entrada para deletar um bloqueio
type DeleteBlockedTimeInput struct {
	TenantID string
	ID       string
}

// DeleteBlockedTimeUseCase deleta um bloqueio de horário
type DeleteBlockedTimeUseCase struct {
	blockedTimeRepo repository.BlockedTimeRepository
}

// NewDeleteBlockedTimeUseCase cria uma nova instância do use case
func NewDeleteBlockedTimeUseCase(blockedTimeRepo repository.BlockedTimeRepository) *DeleteBlockedTimeUseCase {
	return &DeleteBlockedTimeUseCase{
		blockedTimeRepo: blockedTimeRepo,
	}
}

// Execute executa o use case
func (uc *DeleteBlockedTimeUseCase) Execute(ctx context.Context, input DeleteBlockedTimeInput) error {
	// Verifica se existe
	_, err := uc.blockedTimeRepo.GetByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return err
	}

	// Deleta
	return uc.blockedTimeRepo.Delete(ctx, input.TenantID, input.ID)
}
