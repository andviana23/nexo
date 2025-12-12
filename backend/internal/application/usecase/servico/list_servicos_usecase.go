package servico

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// ListServicosUseCase lista serviços com filtros
type ListServicosUseCase struct {
	repo   port.ServicoRepository
	logger *zap.Logger
}

// NewListServicosUseCase cria uma nova instância do use case
func NewListServicosUseCase(repo port.ServicoRepository, logger *zap.Logger) *ListServicosUseCase {
	return &ListServicosUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute lista serviços com filtros
func (uc *ListServicosUseCase) Execute(
	ctx context.Context,
	tenantID string,
	req dto.ListServicosRequest,
) (*dto.ListServicosResponse, error) {
	filter := port.ServicoFilter{
		ApenasAtivos:   req.ApenasAtivos,
		CategoriaID:    req.CategoriaID,
		ProfissionalID: req.ProfissionalID,
		Search:         req.Search,
		OrderBy:        req.OrderBy,
		UnitID:         req.UnitID,
	}

	servicos, err := uc.repo.List(ctx, tenantID, filter.UnitID, filter)
	if err != nil {
		uc.logger.Error("Erro ao listar serviços",
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
		return nil, err
	}

	return &dto.ListServicosResponse{
		Servicos: mapper.ServicosToResponse(servicos),
		Total:    len(servicos),
	}, nil
}

// ListServicosByCategoriaUseCase lista serviços de uma categoria
type ListServicosByCategoriaUseCase struct {
	repo   port.ServicoRepository
	logger *zap.Logger
}

// NewListServicosByCategoriaUseCase cria uma nova instância
func NewListServicosByCategoriaUseCase(repo port.ServicoRepository, logger *zap.Logger) *ListServicosByCategoriaUseCase {
	return &ListServicosByCategoriaUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute lista serviços de uma categoria específica
func (uc *ListServicosByCategoriaUseCase) Execute(
	ctx context.Context,
	tenantID string,
	unitID string,
	categoriaID string,
) (*dto.ListServicosResponse, error) {
	servicos, err := uc.repo.ListByCategoria(ctx, tenantID, unitID, categoriaID)
	if err != nil {
		uc.logger.Error("Erro ao listar serviços por categoria",
			zap.String("tenant_id", tenantID),
			zap.String("unit_id", unitID),
			zap.String("categoria_id", categoriaID),
			zap.Error(err),
		)
		return nil, err
	}

	return &dto.ListServicosResponse{
		Servicos: mapper.ServicosToResponse(servicos),
		Total:    len(servicos),
	}, nil
}

// ListServicosByProfissionalUseCase lista serviços de um profissional
type ListServicosByProfissionalUseCase struct {
	repo   port.ServicoRepository
	logger *zap.Logger
}

// NewListServicosByProfissionalUseCase cria uma nova instância
func NewListServicosByProfissionalUseCase(repo port.ServicoRepository, logger *zap.Logger) *ListServicosByProfissionalUseCase {
	return &ListServicosByProfissionalUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute lista serviços que um profissional pode realizar
func (uc *ListServicosByProfissionalUseCase) Execute(
	ctx context.Context,
	tenantID string,
	unitID string,
	profissionalID string,
) (*dto.ListServicosResponse, error) {
	servicos, err := uc.repo.ListByProfissional(ctx, tenantID, unitID, profissionalID)
	if err != nil {
		uc.logger.Error("Erro ao listar serviços por profissional",
			zap.String("tenant_id", tenantID),
			zap.String("unit_id", unitID),
			zap.String("profissional_id", profissionalID),
			zap.Error(err),
		)
		return nil, err
	}

	return &dto.ListServicosResponse{
		Servicos: mapper.ServicosToResponse(servicos),
		Total:    len(servicos),
	}, nil
}
