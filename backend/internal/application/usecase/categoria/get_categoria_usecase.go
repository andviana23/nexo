package categoria

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// GetCategoriaServicoUseCase busca uma categoria por ID
type GetCategoriaServicoUseCase struct {
	repo   port.CategoriaServicoRepository
	logger *zap.Logger
}

// NewGetCategoriaServicoUseCase cria uma nova instância do use case
func NewGetCategoriaServicoUseCase(repo port.CategoriaServicoRepository, logger *zap.Logger) *GetCategoriaServicoUseCase {
	return &GetCategoriaServicoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute busca uma categoria por ID
func (uc *GetCategoriaServicoUseCase) Execute(
	ctx context.Context,
	tenantID, unitID, categoriaID string,
) (*dto.CategoriaServicoResponse, error) {
	// Buscar categoria
	categoria, err := uc.repo.FindByID(ctx, tenantID, unitID, categoriaID)
	if err != nil {
		uc.logger.Error("erro ao buscar categoria de serviço",
			zap.Error(err),
			zap.String("tenant_id", tenantID),
			zap.String("categoria_id", categoriaID),
		)
		return nil, err
	}

	if categoria == nil {
		return nil, domain.ErrCategoriaNotFound
	}

	return mapper.CategoriaServicoToResponse(categoria), nil
}
