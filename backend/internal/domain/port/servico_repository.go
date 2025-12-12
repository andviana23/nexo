package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// ServicoFilter contém filtros para listagem de serviços
type ServicoFilter struct {
	ApenasAtivos   bool
	CategoriaID    string
	ProfissionalID string
	Search         string
	OrderBy        string // "nome", "preco", "duracao", "criado_em"
}

// ServicoStats contém estatísticas de serviços
type ServicoStats struct {
	TotalServicos    int64
	ServicosAtivos   int64
	ServicosInativos int64
	PrecoMedio       float64
	DuracaoMedia     float64
	ComissaoMedia    float64
}

// ServicoRepository define operações para Serviços
type ServicoRepository interface {
	// Create cria um novo serviço
	Create(ctx context.Context, servico *entity.Servico) error

	// FindByID busca um serviço por ID
	FindByID(ctx context.Context, tenantID, unitID, id string) (*entity.Servico, error)

	// List lista todos os serviços com filtros
	List(ctx context.Context, tenantID, unitID string, filter ServicoFilter) ([]*entity.Servico, error)

	// ListByCategoria lista serviços de uma categoria específica
	ListByCategoria(ctx context.Context, tenantID, unitID, categoriaID string) ([]*entity.Servico, error)

	// ListByProfissional lista serviços que um profissional pode realizar
	ListByProfissional(ctx context.Context, tenantID, unitID, profissionalID string) ([]*entity.Servico, error)

	// FindByIDs busca múltiplos serviços por IDs
	FindByIDs(ctx context.Context, tenantID, unitID string, ids []string) ([]*entity.Servico, error)

	// Update atualiza um serviço existente
	Update(ctx context.Context, servico *entity.Servico) error

	// Delete deleta um serviço
	Delete(ctx context.Context, tenantID, unitID, id string) error

	// DeleteByCategoria deleta todos os serviços de uma categoria
	DeleteByCategoria(ctx context.Context, tenantID, unitID, categoriaID string) error

	// CheckNomeExists verifica se já existe serviço com o mesmo nome
	CheckNomeExists(ctx context.Context, tenantID, unitID, nome, excludeID string) (bool, error)

	// ToggleStatus ativa/desativa um serviço
	ToggleStatus(ctx context.Context, tenantID, unitID, id string, ativo bool) error

	// UpdateCategoria atualiza a categoria de um serviço
	UpdateCategoria(ctx context.Context, tenantID, unitID, id, categoriaID string) error

	// UpdateProfissionais atualiza a lista de profissionais de um serviço
	UpdateProfissionais(ctx context.Context, tenantID, unitID, id string, profissionaisIDs []string) error

	// GetStats retorna estatísticas dos serviços
	GetStats(ctx context.Context, tenantID, unitID string) (*ServicoStats, error)

	// Count conta total de serviços do tenant
	Count(ctx context.Context, tenantID, unitID string) (int64, error)

	// CountAtivos conta total de serviços ativos do tenant
	CountAtivos(ctx context.Context, tenantID, unitID string) (int64, error)
}
