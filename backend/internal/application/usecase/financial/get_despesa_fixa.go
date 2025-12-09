package financial

import (
	"context"


	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// GetDespesaFixaInput define os dados de entrada para buscar despesa fixa
type GetDespesaFixaInput struct {
	TenantID string
	ID       string
}

// GetDespesaFixaUseCase implementa a busca de despesa fixa por ID
type GetDespesaFixaUseCase struct {
	repo   port.DespesaFixaRepository
	logger *zap.Logger
}

// NewGetDespesaFixaUseCase cria nova instância do use case
func NewGetDespesaFixaUseCase(
	repo port.DespesaFixaRepository,
	logger *zap.Logger,
) *GetDespesaFixaUseCase {
	return &GetDespesaFixaUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute busca uma despesa fixa por ID
func (uc *GetDespesaFixaUseCase) Execute(ctx context.Context, input GetDespesaFixaInput) (*entity.DespesaFixa, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.ID == "" {
		return nil, domain.ErrInvalidID
	}

	despesa, err := uc.repo.FindByID(ctx, input.TenantID, input.ID)
	if err != nil {
		uc.logger.Error("Erro ao buscar despesa fixa",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("id", input.ID),
		)
		return nil, err
	}

	return despesa, nil
}
