package financial

import (
"context"

"fmt"

"github.com/andviana23/barber-analytics-backend/internal/domain"
"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
"github.com/andviana23/barber-analytics-backend/internal/domain/port"
"go.uber.org/zap"
)

type GetCompensacaoUseCase struct {
	repo   port.CompensacaoBancariaRepository
	logger *zap.Logger
}

func NewGetCompensacaoUseCase(repo port.CompensacaoBancariaRepository, logger *zap.Logger) *GetCompensacaoUseCase {
	return &GetCompensacaoUseCase{repo: repo, logger: logger}
}

func (uc *GetCompensacaoUseCase) Execute(ctx context.Context, tenantID, id string) (*entity.CompensacaoBancaria, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	comp, err := uc.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar compensação: %w", err)
	}

	return comp, nil
}
