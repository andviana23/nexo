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

type ListContasPagarInput struct {
	TenantID    string
	DataInicio  time.Time
	DataFim     time.Time
	ApenasAtivo bool
}

type ListContasPagarUseCase struct {
	repo   port.ContaPagarRepository
	logger *zap.Logger
}

func NewListContasPagarUseCase(repo port.ContaPagarRepository, logger *zap.Logger) *ListContasPagarUseCase {
	return &ListContasPagarUseCase{repo: repo, logger: logger}
}

func (uc *ListContasPagarUseCase) Execute(ctx context.Context, input ListContasPagarInput) ([]*entity.ContaPagar, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	contas, err := uc.repo.ListByDateRange(ctx, input.TenantID, input.DataInicio, input.DataFim)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas a pagar: %w", err)
	}

	return contas, nil
}
