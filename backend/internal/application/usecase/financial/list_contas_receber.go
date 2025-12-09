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

type ListContasReceberInput struct {
	TenantID    string
	DataInicio  time.Time
	DataFim     time.Time
	ApenasAtivo bool
}

type ListContasReceberUseCase struct {
	repo   port.ContaReceberRepository
	logger *zap.Logger
}

func NewListContasReceberUseCase(repo port.ContaReceberRepository, logger *zap.Logger) *ListContasReceberUseCase {
	return &ListContasReceberUseCase{repo: repo, logger: logger}
}

func (uc *ListContasReceberUseCase) Execute(ctx context.Context, input ListContasReceberInput) ([]*entity.ContaReceber, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	contas, err := uc.repo.ListByDateRange(ctx, input.TenantID, input.DataInicio, input.DataFim)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas a receber: %w", err)
	}

	return contas, nil
}
