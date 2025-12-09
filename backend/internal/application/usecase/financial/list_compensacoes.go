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

type ListCompensacoesInput struct {
	TenantID   string
	DataInicio time.Time
	DataFim    time.Time
}

type ListCompensacoesUseCase struct {
	repo   port.CompensacaoBancariaRepository
	logger *zap.Logger
}

func NewListCompensacoesUseCase(repo port.CompensacaoBancariaRepository, logger *zap.Logger) *ListCompensacoesUseCase {
	return &ListCompensacoesUseCase{repo: repo, logger: logger}
}

func (uc *ListCompensacoesUseCase) Execute(ctx context.Context, input ListCompensacoesInput) ([]*entity.CompensacaoBancaria, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	comps, err := uc.repo.ListByDateRange(ctx, input.TenantID, input.DataInicio, input.DataFim)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar compensações: %w", err)
	}

	return comps, nil
}
