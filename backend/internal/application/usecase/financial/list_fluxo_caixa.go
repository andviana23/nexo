package financial

import (
	"context"

	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

type ListFluxoCaixaInput struct {
	TenantID   string
	DataInicio time.Time
	DataFim    time.Time
}

type ListFluxoCaixaUseCase struct {
	repo   port.FluxoCaixaDiarioRepository
	logger *zap.Logger
}

func NewListFluxoCaixaUseCase(repo port.FluxoCaixaDiarioRepository, logger *zap.Logger) *ListFluxoCaixaUseCase {
	return &ListFluxoCaixaUseCase{repo: repo, logger: logger}
}

func (uc *ListFluxoCaixaUseCase) Execute(ctx context.Context, input ListFluxoCaixaInput) ([]*entity.FluxoCaixaDiario, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	fluxos, err := uc.repo.ListByDateRange(ctx, input.TenantID, input.DataInicio, input.DataFim)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar fluxo de caixa: %w", err)
	}

	return fluxos, nil
}
