package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// CategoriaServicoFilter contém filtros para listagem de categorias
type CategoriaServicoFilter struct {
	ApenasAtivas bool
	OrderBy      string // "nome", "criado_em"
}

// CategoriaServicoRepository define operações para Categorias de Serviço
type CategoriaServicoRepository interface {
	// Create cria uma nova categoria
	Create(ctx context.Context, categoria *entity.CategoriaServico) error

	// FindByID busca uma categoria por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.CategoriaServico, error)

	// List lista todas as categorias com filtros
	List(ctx context.Context, tenantID string, filter CategoriaServicoFilter) ([]*entity.CategoriaServico, error)

	// Update atualiza uma categoria existente
	Update(ctx context.Context, categoria *entity.CategoriaServico) error

	// Delete deleta uma categoria
	Delete(ctx context.Context, tenantID, id string) error

	// CheckNomeExists verifica se já existe categoria com o mesmo nome
	CheckNomeExists(ctx context.Context, tenantID, nome, excludeID string) (bool, error)

	// CountServicos conta quantos serviços estão vinculados à categoria
	CountServicos(ctx context.Context, tenantID, categoriaID string) (int64, error)

	// ToggleStatus ativa/desativa uma categoria
	ToggleStatus(ctx context.Context, tenantID, id string, ativa bool) error
}
