package port

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
)

// MovimentacaoEstoqueRepository define as operações de persistência para movimentações de estoque
type MovimentacaoEstoqueRepository interface {
	// CRUD básico
	Create(ctx context.Context, movimentacao *entity.MovimentacaoEstoque) error
	FindByID(ctx context.Context, tenantID, movimentacaoID uuid.UUID) (*entity.MovimentacaoEstoque, error)

	// Listagens
	ListByProduto(ctx context.Context, tenantID, produtoID uuid.UUID, limit int) ([]*entity.MovimentacaoEstoque, error)
	ListByTipo(ctx context.Context, tenantID uuid.UUID, tipo entity.TipoMovimentacao, limit int) ([]*entity.MovimentacaoEstoque, error)
	ListByPeriodo(ctx context.Context, tenantID uuid.UUID, inicio, fim time.Time) ([]*entity.MovimentacaoEstoque, error)

	// Relatórios
	GetTotalPorTipo(ctx context.Context, tenantID uuid.UUID, tipo entity.TipoMovimentacao, inicio, fim time.Time) (int64, error)
}
