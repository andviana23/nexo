package financial

import (
	"context"

	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"go.uber.org/zap"
)

type ListContasReceberInput struct {
	TenantID   string
	Status     *valueobject.StatusConta
	Origem     *string
	DataInicio *time.Time
	DataFim    *time.Time
	Page       int
	PageSize   int
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

	contas, err := uc.repo.List(ctx, input.TenantID, port.ContaReceberListFilters{
		Status:     input.Status,
		Origem:     input.Origem,
		DataInicio: input.DataInicio,
		DataFim:    input.DataFim,
		Page:       input.Page,
		PageSize:   input.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas a receber: %w", err)
	}

	return contas, nil
}
