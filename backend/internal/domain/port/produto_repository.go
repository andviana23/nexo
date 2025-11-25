package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProdutoRepository define as operações de persistência para produtos
type ProdutoRepository interface {
	// CRUD básico
	Create(ctx context.Context, produto *entity.Produto) error
	FindByID(ctx context.Context, tenantID, produtoID uuid.UUID) (*entity.Produto, error)
	FindBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*entity.Produto, error)
	Update(ctx context.Context, produto *entity.Produto) error
	Delete(ctx context.Context, tenantID, produtoID uuid.UUID) error

	// Listagens
	ListAll(ctx context.Context, tenantID uuid.UUID) ([]*entity.Produto, error)
	ListByCategoria(ctx context.Context, tenantID uuid.UUID, categoria entity.CategoriaProduto) ([]*entity.Produto, error)
	ListAbaixoDoMinimo(ctx context.Context, tenantID uuid.UUID) ([]*entity.Produto, error)

	// Operações de estoque
	AtualizarQuantidade(ctx context.Context, tenantID, produtoID uuid.UUID, novaQuantidade decimal.Decimal) error
}
