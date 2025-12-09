package financial

import (
	"context"


	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// ListDespesasFixasInput define os dados de entrada para listar despesas fixas
type ListDespesasFixasInput struct {
	TenantID    string
	CategoriaID *string
	UnidadeID   *string
	Ativo       *bool
	Page        int
	PageSize    int
}

// ListDespesasFixasOutput define os dados de saída
type ListDespesasFixasOutput struct {
	Despesas []*entity.DespesaFixa
	Total    int64
	Page     int
	PageSize int
}

// ListDespesasFixasUseCase implementa a listagem de despesas fixas
type ListDespesasFixasUseCase struct {
	repo   port.DespesaFixaRepository
	logger *zap.Logger
}

// NewListDespesasFixasUseCase cria nova instância do use case
func NewListDespesasFixasUseCase(
	repo port.DespesaFixaRepository,
	logger *zap.Logger,
) *ListDespesasFixasUseCase {
	return &ListDespesasFixasUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute lista despesas fixas com paginação
func (uc *ListDespesasFixasUseCase) Execute(ctx context.Context, input ListDespesasFixasInput) (*ListDespesasFixasOutput, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	// Defaults
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}
	if input.PageSize > 100 {
		input.PageSize = 100
	}

	filters := port.DespesaFixaListFilters{
		CategoriaID: input.CategoriaID,
		UnidadeID:   input.UnidadeID,
		Ativo:       input.Ativo,
		Page:        input.Page,
		PageSize:    input.PageSize,
	}

	despesas, total, err := uc.repo.List(ctx, input.TenantID, filters)
	if err != nil {
		uc.logger.Error("Erro ao listar despesas fixas",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
		)
		return nil, err
	}

	return &ListDespesasFixasOutput{
		Despesas: despesas,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}
