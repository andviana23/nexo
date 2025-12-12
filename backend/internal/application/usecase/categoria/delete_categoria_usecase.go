package categoria

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// DeleteCategoriaServicoUseCase deleta uma categoria de serviço
type DeleteCategoriaServicoUseCase struct {
	repo   port.CategoriaServicoRepository
	logger *zap.Logger
}

// NewDeleteCategoriaServicoUseCase cria uma nova instância do use case
func NewDeleteCategoriaServicoUseCase(repo port.CategoriaServicoRepository, logger *zap.Logger) *DeleteCategoriaServicoUseCase {
	return &DeleteCategoriaServicoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute deleta uma categoria de serviço
func (uc *DeleteCategoriaServicoUseCase) Execute(
	ctx context.Context,
	tenantID, unitID, categoriaID string,
) error {
	// Verificar se categoria existe
	categoria, err := uc.repo.FindByID(ctx, tenantID, unitID, categoriaID)
	if err != nil {
		uc.logger.Error("erro ao buscar categoria de serviço",
			zap.Error(err),
			zap.String("tenant_id", tenantID),
			zap.String("categoria_id", categoriaID),
		)
		return err
	}

	if categoria == nil {
		return domain.ErrCategoriaNotFound
	}

	// Verificar se há serviços vinculados
	count, err := uc.repo.CountServicos(ctx, tenantID, unitID, categoriaID)
	if err != nil {
		uc.logger.Error("erro ao contar serviços da categoria",
			zap.Error(err),
			zap.String("tenant_id", tenantID),
			zap.String("categoria_id", categoriaID),
		)
		return err
	}

	if count > 0 {
		uc.logger.Warn("tentativa de deletar categoria com serviços vinculados",
			zap.String("tenant_id", tenantID),
			zap.String("categoria_id", categoriaID),
			zap.Int64("total_servicos", count),
		)
		return domain.ErrCategoriaHasServices
	}

	// Deletar categoria
	if err := uc.repo.Delete(ctx, tenantID, unitID, categoriaID); err != nil {
		uc.logger.Error("erro ao deletar categoria de serviço",
			zap.Error(err),
			zap.String("tenant_id", tenantID),
			zap.String("categoria_id", categoriaID),
		)
		return err
	}

	uc.logger.Info("categoria de serviço deletada com sucesso",
		zap.String("categoria_id", categoriaID),
		zap.String("tenant_id", tenantID),
		zap.String("nome", categoria.Nome),
	)

	return nil
}
