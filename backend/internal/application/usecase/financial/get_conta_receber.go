package financial

import (
"context"

"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type GetContaReceberUseCase struct {
	repo   port.ContaReceberRepository
	logger *zap.Logger
}

func NewGetContaReceberUseCase(repo port.ContaReceberRepository, logger *zap.Logger) *GetContaReceberUseCase {
	return &GetContaReceberUseCase{repo: repo, logger: logger}
}

func (uc *GetContaReceberUseCase) Execute(ctx context.Context, tenantID, id string) (*entity.ContaReceber, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	conta, err := uc.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conta a receber: %w", err)
	}

	return conta, nil
}
