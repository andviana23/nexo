package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// MeioPagamentoRepository define operações para Meios de Pagamento
type MeioPagamentoRepository interface {
	// Create cria um novo meio de pagamento
	Create(ctx context.Context, meio *entity.MeioPagamento) error

	// FindByID busca um meio de pagamento por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.MeioPagamento, error)

	// List lista todos os meios de pagamento do tenant
	List(ctx context.Context, tenantID string) ([]*entity.MeioPagamento, error)

	// ListAtivos lista apenas os meios de pagamento ativos
	ListAtivos(ctx context.Context, tenantID string) ([]*entity.MeioPagamento, error)

	// ListByTipo lista meios de pagamento por tipo
	ListByTipo(ctx context.Context, tenantID string, tipo entity.TipoPagamento) ([]*entity.MeioPagamento, error)

	// Count retorna a quantidade total de meios de pagamento
	Count(ctx context.Context, tenantID string) (int64, error)

	// CountAtivos retorna a quantidade de meios de pagamento ativos
	CountAtivos(ctx context.Context, tenantID string) (int64, error)

	// Update atualiza um meio de pagamento existente
	Update(ctx context.Context, meio *entity.MeioPagamento) error

	// Toggle alterna o status ativo/inativo
	Toggle(ctx context.Context, tenantID, id string) (*entity.MeioPagamento, error)

	// Delete remove um meio de pagamento
	Delete(ctx context.Context, tenantID, id string) error

	// ExistsByNome verifica se existe meio de pagamento com o nome
	ExistsByNome(ctx context.Context, tenantID, nome string) (bool, error)
}
