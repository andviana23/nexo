package financial

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// ToggleDespesaFixaInput define os dados de entrada para alternar status
type ToggleDespesaFixaInput struct {
	TenantID string
	ID       string
}

// ToggleDespesaFixaUseCase implementa a alternância de status ativo/inativo
type ToggleDespesaFixaUseCase struct {
	repo   port.DespesaFixaRepository
	logger *zap.Logger
}

// NewToggleDespesaFixaUseCase cria nova instância do use case
func NewToggleDespesaFixaUseCase(
	repo port.DespesaFixaRepository,
	logger *zap.Logger,
) *ToggleDespesaFixaUseCase {
	return &ToggleDespesaFixaUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute alterna o status ativo/inativo da despesa fixa
func (uc *ToggleDespesaFixaUseCase) Execute(ctx context.Context, input ToggleDespesaFixaInput) (*entity.DespesaFixa, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.ID == "" {
		return nil, domain.ErrInvalidID
	}

	despesa, err := uc.repo.Toggle(ctx, input.TenantID, input.ID)
	if err != nil {
		uc.logger.Error("Erro ao alternar status de despesa fixa",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("id", input.ID),
		)
		return nil, err
	}

	uc.logger.Info("Status de despesa fixa alternado com sucesso",
		zap.String("id", despesa.ID),
		zap.String("tenant_id", despesa.TenantID.String()),
		zap.Bool("ativo", despesa.Ativo),
	)

	return despesa, nil
}
