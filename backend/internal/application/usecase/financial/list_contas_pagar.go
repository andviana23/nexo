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

type ListContasPagarInput struct {
	TenantID    string
	Status      *valueobject.StatusConta
	CategoriaID *string
	Tipo        *valueobject.TipoCusto
	DataInicio  *time.Time
	DataFim     *time.Time
	Page        int
	PageSize    int
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

	contas, err := uc.repo.List(ctx, input.TenantID, port.ContaPagarListFilters{
		Status:      input.Status,
		Tipo:        input.Tipo,
		CategoriaID: input.CategoriaID,
		DataInicio:  input.DataInicio,
		DataFim:     input.DataFim,
		Page:        input.Page,
		PageSize:    input.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas a pagar: %w", err)
	}

	return contas, nil
}
