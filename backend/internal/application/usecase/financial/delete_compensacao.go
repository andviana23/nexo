package financial

import (
"context"

"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type DeleteCompensacaoUseCase struct {
	repo   port.CompensacaoBancariaRepository
	logger *zap.Logger
}

func NewDeleteCompensacaoUseCase(repo port.CompensacaoBancariaRepository, logger *zap.Logger) *DeleteCompensacaoUseCase {
	return &DeleteCompensacaoUseCase{repo: repo, logger: logger}
}

func (uc *DeleteCompensacaoUseCase) Execute(ctx context.Context, tenantID, id string) error {
	if tenantID == "" {
		return domain.ErrTenantIDRequired
	}
	if id == "" {
		return domain.ErrInvalidID
	}

	if _, err := uc.repo.FindByID(ctx, tenantID, id); err != nil {
		return fmt.Errorf("compensação não encontrada: %w", err)
	}

	if err := uc.repo.Delete(ctx, tenantID, id); err != nil {
		return fmt.Errorf("erro ao deletar compensação: %w", err)
	}

	uc.logger.Info("Compensação deletada", zap.String("tenant_id", tenantID), zap.String("id", id))

	return nil
}
