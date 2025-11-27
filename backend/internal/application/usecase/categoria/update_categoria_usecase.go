package categoria

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// UpdateCategoriaServicoUseCase atualiza uma categoria de serviço
type UpdateCategoriaServicoUseCase struct {
	repo   port.CategoriaServicoRepository
	logger *zap.Logger
}

// NewUpdateCategoriaServicoUseCase cria uma nova instância do use case
func NewUpdateCategoriaServicoUseCase(repo port.CategoriaServicoRepository, logger *zap.Logger) *UpdateCategoriaServicoUseCase {
	return &UpdateCategoriaServicoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute atualiza uma categoria de serviço
func (uc *UpdateCategoriaServicoUseCase) Execute(
	ctx context.Context,
	tenantID, categoriaID string,
	req dto.UpdateCategoriaServicoRequest,
) (*dto.CategoriaServicoResponse, error) {
	// Buscar categoria existente
	categoria, err := uc.repo.FindByID(ctx, tenantID, categoriaID)
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

	// Verificar duplicidade de nome (exceto própria categoria)
	if req.Nome != categoria.Nome {
		exists, err := uc.repo.CheckNomeExists(ctx, tenantID, req.Nome, categoriaID)
		if err != nil {
			uc.logger.Error("erro ao verificar nome duplicado",
				zap.Error(err),
				zap.String("nome", req.Nome),
			)
			return nil, err
		}
		if exists {
			return nil, domain.ErrCategoriaNomeDuplicate
		}
	}

	// Atualizar entidade
	if err := categoria.Update(req.Nome, req.Descricao, req.Cor, req.Icone); err != nil {
		return nil, err
	}

	// Persistir alterações
	if err := uc.repo.Update(ctx, categoria); err != nil {
		uc.logger.Error("erro ao atualizar categoria de serviço",
			zap.Error(err),
			zap.String("tenant_id", tenantID),
			zap.String("categoria_id", categoriaID),
		)
		return nil, err
	}

	uc.logger.Info("categoria de serviço atualizada com sucesso",
		zap.String("categoria_id", categoriaID),
		zap.String("tenant_id", tenantID),
		zap.String("nome", categoria.Nome),
	)

	return mapper.CategoriaServicoToResponse(categoria), nil
}
