package financial

import (
"context"

"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type GetContaPagarUseCase struct {
	repo   port.ContaPagarRepository
	logger *zap.Logger
}

func NewGetContaPagarUseCase(repo port.ContaPagarRepository, logger *zap.Logger) *GetContaPagarUseCase {
	return &GetContaPagarUseCase{repo: repo, logger: logger}
}

func (uc *GetContaPagarUseCase) Execute(ctx context.Context, tenantID, id string) (*entity.ContaPagar, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	conta, err := uc.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conta a pagar: %w", err)
	}

	return conta, nil
}
