package financial

import (
	"context"

	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"go.uber.org/zap"
)

type ListDREInput struct {
	TenantID string
	Inicio   valueobject.MesAno
	Fim      valueobject.MesAno
}

type ListDREUseCase struct {
	repo   port.DREMensalRepository
	logger *zap.Logger
}

func NewListDREUseCase(repo port.DREMensalRepository, logger *zap.Logger) *ListDREUseCase {
	return &ListDREUseCase{repo: repo, logger: logger}
}

func (uc *ListDREUseCase) Execute(ctx context.Context, input ListDREInput) ([]*entity.DREMensal, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	dres, err := uc.repo.ListByPeriod(ctx, input.TenantID, input.Inicio, input.Fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar DREs: %w", err)
	}

	return dres, nil
}
