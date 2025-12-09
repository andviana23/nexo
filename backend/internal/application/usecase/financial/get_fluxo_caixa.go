package financial

import (
"context"

"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type GetFluxoCaixaUseCase struct {
	repo   port.FluxoCaixaDiarioRepository
	logger *zap.Logger
}

func NewGetFluxoCaixaUseCase(repo port.FluxoCaixaDiarioRepository, logger *zap.Logger) *GetFluxoCaixaUseCase {
	return &GetFluxoCaixaUseCase{repo: repo, logger: logger}
}

func (uc *GetFluxoCaixaUseCase) Execute(ctx context.Context, tenantID, id string) (*entity.FluxoCaixaDiario, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	fluxo, err := uc.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar fluxo de caixa: %w", err)
	}

	return fluxo, nil
}
