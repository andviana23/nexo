package categoria

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// ListCategoriasServicosUseCase lista categorias de serviço
type ListCategoriasServicosUseCase struct {
	repo   port.CategoriaServicoRepository
	logger *zap.Logger
}

// NewListCategoriasServicosUseCase cria uma nova instância do use case
func NewListCategoriasServicosUseCase(repo port.CategoriaServicoRepository, logger *zap.Logger) *ListCategoriasServicosUseCase {
	return &ListCategoriasServicosUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute lista categorias de serviço com filtros
func (uc *ListCategoriasServicosUseCase) Execute(
	ctx context.Context,
	tenantID string,
	req dto.ListCategoriasServicosRequest,
) (*dto.ListCategoriasServicosResponse, error) {
	// Construir filtros
	filter := port.CategoriaServicoFilter{
		ApenasAtivas: req.ApenasAtivas,
		OrderBy:      req.OrderBy,
	}

	// Definir ordenação padrão
	if filter.OrderBy == "" {
		filter.OrderBy = "nome"
	}

	// Buscar categorias
	categorias, err := uc.repo.List(ctx, tenantID, filter)
	if err != nil {
		uc.logger.Error("erro ao listar categorias de serviço",
			zap.Error(err),
			zap.String("tenant_id", tenantID),
		)
		return nil, err
	}

	// Converter para DTO
	response := &dto.ListCategoriasServicosResponse{
		Categorias: mapper.CategoriasServicosToResponse(categorias),
		Total:      len(categorias),
	}

	uc.logger.Info("categorias listadas com sucesso",
		zap.String("tenant_id", tenantID),
		zap.Int("total", response.Total),
	)

	return response, nil
}
