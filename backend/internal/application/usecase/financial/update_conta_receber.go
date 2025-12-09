package financial

import (
"context"

"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type UpdateContaReceberInput struct {
	TenantID string
	ID       string
	Conta    *entity.ContaReceber
}

type UpdateContaReceberUseCase struct {
	repo   port.ContaReceberRepository
	logger *zap.Logger
}

func NewUpdateContaReceberUseCase(repo port.ContaReceberRepository, logger *zap.Logger) *UpdateContaReceberUseCase {
	return &UpdateContaReceberUseCase{repo: repo, logger: logger}
}

func (uc *UpdateContaReceberUseCase) Execute(ctx context.Context, input UpdateContaReceberInput) (*entity.ContaReceber, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if input.ID == "" {
		return nil, domain.ErrInvalidID
	}

	_, err := uc.repo.FindByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, fmt.Errorf("conta não encontrada: %w", err)
	}

	if err := uc.repo.Update(ctx, input.Conta); err != nil {
		return nil, fmt.Errorf("erro ao atualizar conta: %w", err)
	}

	uc.logger.Info("Conta a receber atualizada", zap.String("tenant_id", input.TenantID), zap.String("id", input.ID))

	return input.Conta, nil
}
