package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
)

// CategoriaProdutoRepository define as operações de persistência para categorias de produtos
type CategoriaProdutoRepository interface {
	// Create cria uma nova categoria de produto
	Create(ctx context.Context, categoria *entity.CategoriaProdutoEntity) error

	// FindByID busca uma categoria por ID
	FindByID(ctx context.Context, tenantID, unitID, id uuid.UUID) (*entity.CategoriaProdutoEntity, error)

	// FindByNome busca uma categoria por nome (único por tenant)
	FindByNome(ctx context.Context, tenantID, unitID uuid.UUID, nome string) (*entity.CategoriaProdutoEntity, error)

	// ListAll lista todas as categorias do tenant
	ListAll(ctx context.Context, tenantID, unitID uuid.UUID) ([]*entity.CategoriaProdutoEntity, error)

	// ListAtivas lista apenas categorias ativas do tenant
	ListAtivas(ctx context.Context, tenantID, unitID uuid.UUID) ([]*entity.CategoriaProdutoEntity, error)

	// Update atualiza uma categoria existente
	Update(ctx context.Context, categoria *entity.CategoriaProdutoEntity) error

	// Delete remove uma categoria (soft delete via ativa=false ou hard delete)
	Delete(ctx context.Context, tenantID, unitID, id uuid.UUID) error

	// ExistsWithNome verifica se já existe categoria com o nome (para validação de duplicidade)
	ExistsWithNome(ctx context.Context, tenantID, unitID uuid.UUID, nome string, excludeID *uuid.UUID) (bool, error)

	// CountProdutosVinculados conta quantos produtos estão vinculados à categoria
	CountProdutosVinculados(ctx context.Context, tenantID, unitID, categoriaID uuid.UUID) (int64, error)
}
