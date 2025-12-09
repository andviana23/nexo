package financial

import (
"context"

"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type GetDREUseCase struct {
	repo   port.DREMensalRepository
	logger *zap.Logger
}

func NewGetDREUseCase(repo port.DREMensalRepository, logger *zap.Logger) *GetDREUseCase {
	return &GetDREUseCase{repo: repo, logger: logger}
}

func (uc *GetDREUseCase) Execute(ctx context.Context, tenantID, id string) (*entity.DREMensal, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	dre, err := uc.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar DRE: %w", err)
	}

	return dre, nil
}
