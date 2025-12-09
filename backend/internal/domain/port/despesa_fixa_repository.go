package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
)

// DespesaFixaRepository define operações para Despesas Fixas
type DespesaFixaRepository interface {
	// Create cria uma nova despesa fixa
	Create(ctx context.Context, despesa *entity.DespesaFixa) error

	// FindByID busca uma despesa fixa por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.DespesaFixa, error)

	// Update atualiza uma despesa fixa existente
	Update(ctx context.Context, despesa *entity.DespesaFixa) error

	// Delete remove uma despesa fixa
	Delete(ctx context.Context, tenantID, id string) error

	// Toggle alterna o status ativo/inativo
	Toggle(ctx context.Context, tenantID, id string) (*entity.DespesaFixa, error)

	// List lista despesas fixas com paginação
	List(ctx context.Context, tenantID string, filters DespesaFixaListFilters) ([]*entity.DespesaFixa, int64, error)

	// ListAtivas lista apenas despesas fixas ativas (para o cron job)
	ListAtivas(ctx context.Context, tenantID string) ([]*entity.DespesaFixa, error)

	// ListByUnidade lista despesas fixas de uma unidade específica
	ListByUnidade(ctx context.Context, tenantID, unidadeID string) ([]*entity.DespesaFixa, error)

	// ListByCategoria lista despesas fixas de uma categoria específica
	ListByCategoria(ctx context.Context, tenantID, categoriaID string) ([]*entity.DespesaFixa, error)

	// ListAtivasPorTenants lista todas as despesas ativas de todos os tenants ativos
	// Usado pelo cron job para geração em massa
	ListAtivasPorTenants(ctx context.Context) ([]*DespesaFixaComTenant, error)

	// SumAtivas soma o valor total de despesas fixas ativas
	SumAtivas(ctx context.Context, tenantID string) (valueobject.Money, error)

	// SumByUnidade soma o valor total de despesas fixas ativas por unidade
	SumByUnidade(ctx context.Context, tenantID, unidadeID string) (valueobject.Money, error)

	// Count conta total de despesas fixas do tenant
	Count(ctx context.Context, tenantID string) (int64, error)

	// CountAtivas conta despesas fixas ativas do tenant
	CountAtivas(ctx context.Context, tenantID string) (int64, error)

	// ExistsByDescricao verifica se já existe despesa fixa com mesma descrição
	ExistsByDescricao(ctx context.Context, tenantID, descricao string, excludeID *string) (bool, error)
}

// DespesaFixaListFilters filtros para listagem de despesas fixas
type DespesaFixaListFilters struct {
	Ativo       *bool
	CategoriaID *string
	UnidadeID   *string
	Fornecedor  *string
	Page        int
	PageSize    int
	OrderBy     string // descricao, valor, dia_vencimento, criado_em
}

// DespesaFixaComTenant representa uma despesa fixa com informações do tenant
// Usado para o cron job que processa múltiplos tenants
type DespesaFixaComTenant struct {
	DespesaFixa *entity.DespesaFixa
	TenantNome  string
}
