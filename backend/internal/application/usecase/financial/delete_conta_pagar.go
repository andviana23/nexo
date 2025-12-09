package financial

import (
"context"

"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type DeleteContaPagarUseCase struct {
	repo   port.ContaPagarRepository
	logger *zap.Logger
}

func NewDeleteContaPagarUseCase(repo port.ContaPagarRepository, logger *zap.Logger) *DeleteContaPagarUseCase {
	return &DeleteContaPagarUseCase{repo: repo, logger: logger}
}

func (uc *DeleteContaPagarUseCase) Execute(ctx context.Context, tenantID, id string) error {
	if tenantID == "" {
		return domain.ErrTenantIDRequired
	}
	if id == "" {
		return domain.ErrInvalidID
	}

	if _, err := uc.repo.FindByID(ctx, tenantID, id); err != nil {
		return fmt.Errorf("conta não encontrada: %w", err)
	}

	if err := uc.repo.Delete(ctx, tenantID, id); err != nil {
		return fmt.Errorf("erro ao deletar conta: %w", err)
	}

	uc.logger.Info("Conta a pagar deletada", zap.String("tenant_id", tenantID), zap.String("id", id))

	return nil
}
