package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
)

// MovimentacaoEstoqueRepositoryPG implementa MovimentacaoEstoqueRepository usando PostgreSQL
type MovimentacaoEstoqueRepositoryPG struct {
	queries *db.Queries
}

// NewMovimentacaoEstoqueRepositoryPG cria nova instancia do repositorio
func NewMovimentacaoEstoqueRepositoryPG(queries *db.Queries) port.MovimentacaoEstoqueRepository {
	return &MovimentacaoEstoqueRepositoryPG{queries: queries}
}

// Create persiste uma nova movimentacao
func (r *MovimentacaoEstoqueRepositoryPG) Create(ctx context.Context, movimentacao *entity.MovimentacaoEstoque) error {
	params := db.CreateMovimentacaoEstoqueParams{
		TenantID:         uuidToPgUUID(movimentacao.TenantID),
		ProdutoID:        uuidToPgUUID(movimentacao.ProdutoID),
		TipoMovimentacao: string(movimentacao.Tipo),
		Quantidade:       movimentacao.Quantidade,    // Ja eh decimal.Decimal
		ValorUnitario:    movimentacao.ValorUnitario, // Ja eh decimal.Decimal
		ValorTotal:       movimentacao.ValorTotal,    // Ja eh decimal.Decimal
		FornecedorID:     uuidPtrToPgUUID(movimentacao.FornecedorID),
		UsuarioID:        uuidPtrToPgUUID(&movimentacao.UsuarioID),
		DataMovimentacao: timeToPgTimestamp(movimentacao.DataMovimentacao),
		Observacoes:      strPtrToPgText(movimentacao.Observacoes), // string -> *string
		Documento:        strPtrToPgText(movimentacao.Documento),   // string -> *string
	}

	created, err := r.queries.CreateMovimentacaoEstoque(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar movimentacao: %w", err)
	}

	movimentacao.ID = pgUUIDToUUID(created.ID)
	movimentacao.CreatedAt = created.CriadoEm.Time
	movimentacao.UpdatedAt = created.AtualizadoEm.Time

	return nil
}

// FindByID busca movimentacao por ID
func (r *MovimentacaoEstoqueRepositoryPG) FindByID(ctx context.Context, tenantID, movimentacaoID uuid.UUID) (*entity.MovimentacaoEstoque, error) {
	params := db.GetMovimentacaoByIDParams{
		ID:       uuidToPgUUID(movimentacaoID),
		TenantID: uuidToPgUUID(tenantID),
	}

	result, err := r.queries.GetMovimentacaoByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar movimentacao: %w", err)
	}

	return r.toDomain(&result), nil
}

// ListByProduto lista movimentacoes de um produto
func (r *MovimentacaoEstoqueRepositoryPG) ListByProduto(ctx context.Context, tenantID, produtoID uuid.UUID, limit int) ([]*entity.MovimentacaoEstoque, error) {
	params := db.ListMovimentacoesByProdutoParams{
		TenantID:  uuidToPgUUID(tenantID),
		ProdutoID: uuidToPgUUID(produtoID),
		Limit:     int32(limit),
		Offset:    0,
	}

	results, err := r.queries.ListMovimentacoesByProduto(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar movimentacoes por produto: %w", err)
	}

	return r.toDomainList(results), nil
}

// ListByTipo lista movimentacoes por tipo
func (r *MovimentacaoEstoqueRepositoryPG) ListByTipo(ctx context.Context, tenantID uuid.UUID, tipo entity.TipoMovimentacao, limit int) ([]*entity.MovimentacaoEstoque, error) {
	params := db.ListMovimentacoesByTipoParams{
		TenantID:         uuidToPgUUID(tenantID),
		TipoMovimentacao: string(tipo),
		Limit:            int32(limit),
		Offset:           0,
	}

	results, err := r.queries.ListMovimentacoesByTipo(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar movimentacoes por tipo: %w", err)
	}

	return r.toDomainList(results), nil
}

// ListByPeriodo lista movimentacoes em um periodo
func (r *MovimentacaoEstoqueRepositoryPG) ListByPeriodo(ctx context.Context, tenantID uuid.UUID, inicio, fim time.Time) ([]*entity.MovimentacaoEstoque, error) {
	params := db.ListMovimentacoesByPeriodoParams{
		TenantID:           uuidToPgUUID(tenantID),
		DataMovimentacao:   timeToPgTimestamp(inicio),
		DataMovimentacao_2: timeToPgTimestamp(fim),
	}

	results, err := r.queries.ListMovimentacoesByPeriodo(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar movimentacoes por periodo: %w", err)
	}

	return r.toDomainList(results), nil
}

// GetTotalPorTipo calcula total de quantidade por tipo em um periodo
func (r *MovimentacaoEstoqueRepositoryPG) GetTotalPorTipo(ctx context.Context, tenantID uuid.UUID, tipo entity.TipoMovimentacao, inicio, fim time.Time) (int64, error) {
	params := db.GetTotalPorTipoParams{
		TenantID:           uuidToPgUUID(tenantID),
		TipoMovimentacao:   string(tipo),
		DataMovimentacao:   timeToPgTimestamp(inicio),
		DataMovimentacao_2: timeToPgTimestamp(fim),
	}

	result, err := r.queries.GetTotalPorTipo(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("erro ao calcular total por tipo: %w", err)
	}

	return result.TotalQuantidade, nil
}

// toDomain converte modelo do sqlc para entidade de dominio
func (r *MovimentacaoEstoqueRepositoryPG) toDomain(m *db.MovimentacoesEstoque) *entity.MovimentacaoEstoque {
	return &entity.MovimentacaoEstoque{
		ID:               pgUUIDToUUID(m.ID),
		TenantID:         pgUUIDToUUID(m.TenantID),
		ProdutoID:        pgUUIDToUUID(m.ProdutoID),
		UsuarioID:        pgUUIDToUUID(m.UsuarioID),
		FornecedorID:     pgUUIDToUUIDPtr(m.FornecedorID),
		Tipo:             entity.TipoMovimentacao(m.TipoMovimentacao),
		Quantidade:       m.Quantidade,    // Ja eh decimal.Decimal
		ValorUnitario:    m.ValorUnitario, // Ja eh decimal.Decimal
		ValorTotal:       m.ValorTotal,    // Ja eh decimal.Decimal
		DataMovimentacao: m.DataMovimentacao.Time,
		Observacoes:      pgTextToStr(m.Observacoes), // *string -> string
		Documento:        pgTextToStr(m.Documento),   // *string -> string
		CreatedAt:        m.CriadoEm.Time,
		UpdatedAt:        m.AtualizadoEm.Time,
	}
}

// toDomainList converte lista de modelos para entidades
func (r *MovimentacaoEstoqueRepositoryPG) toDomainList(results []db.MovimentacoesEstoque) []*entity.MovimentacaoEstoque {
	movimentacoes := make([]*entity.MovimentacaoEstoque, len(results))
	for i, result := range results {
		movimentacoes[i] = r.toDomain(&result)
	}
	return movimentacoes
}
