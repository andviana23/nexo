package financial

import (
"context"

"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type UpdateContaPagarInput struct {
	TenantID string
	ID       string
	Conta    *entity.ContaPagar
}

type UpdateContaPagarUseCase struct {
	repo   port.ContaPagarRepository
	logger *zap.Logger
}

func NewUpdateContaPagarUseCase(repo port.ContaPagarRepository, logger *zap.Logger) *UpdateContaPagarUseCase {
	return &UpdateContaPagarUseCase{repo: repo, logger: logger}
}

func (uc *UpdateContaPagarUseCase) Execute(ctx context.Context, input UpdateContaPagarInput) (*entity.ContaPagar, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if input.ID == "" {
		return nil, domain.ErrInvalidID
	}

	// Verificar se existe
	_, err := uc.repo.FindByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, fmt.Errorf("conta não encontrada: %w", err)
	}

	// Atualizar
	if err := uc.repo.Update(ctx, input.Conta); err != nil {
		return nil, fmt.Errorf("erro ao atualizar conta: %w", err)
	}

	uc.logger.Info("Conta a pagar atualizada", zap.String("tenant_id", input.TenantID), zap.String("id", input.ID))

	return input.Conta, nil
}
