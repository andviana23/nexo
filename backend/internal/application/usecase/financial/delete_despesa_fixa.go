package financial

import (
	"context"


	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// DeleteDespesaFixaInput define os dados de entrada para deletar despesa fixa
type DeleteDespesaFixaInput struct {
	TenantID string
	ID       string
}

// DeleteDespesaFixaUseCase implementa a remoção de despesa fixa
type DeleteDespesaFixaUseCase struct {
	repo   port.DespesaFixaRepository
	logger *zap.Logger
}

// NewDeleteDespesaFixaUseCase cria nova instância do use case
func NewDeleteDespesaFixaUseCase(
	repo port.DespesaFixaRepository,
	logger *zap.Logger,
) *DeleteDespesaFixaUseCase {
	return &DeleteDespesaFixaUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute remove uma despesa fixa
func (uc *DeleteDespesaFixaUseCase) Execute(ctx context.Context, input DeleteDespesaFixaInput) error {
	if input.TenantID == "" {
		return domain.ErrTenantIDRequired
	}

	if input.ID == "" {
		return domain.ErrInvalidID
	}

	// Verificar se existe antes de deletar
	_, err := uc.repo.FindByID(ctx, input.TenantID, input.ID)
	if err != nil {
		uc.logger.Error("Erro ao buscar despesa fixa para deleção",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("id", input.ID),
		)
		return err
	}

	if err := uc.repo.Delete(ctx, input.TenantID, input.ID); err != nil {
		uc.logger.Error("Erro ao deletar despesa fixa",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("id", input.ID),
		)
		return err
	}

	uc.logger.Info("Despesa fixa deletada com sucesso",
		zap.String("id", input.ID),
		zap.String("tenant_id", input.TenantID),
	)

	return nil
}
