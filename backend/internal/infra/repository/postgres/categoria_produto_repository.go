package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Garantir que implementa a interface
var _ port.CategoriaProdutoRepository = (*CategoriaProdutoRepository)(nil)

// CategoriaProdutoRepository implementa port.CategoriaProdutoRepository usando PostgreSQL
type CategoriaProdutoRepository struct {
	queries *db.Queries
}

// NewCategoriaProdutoRepository cria uma nova instância
func NewCategoriaProdutoRepository(queries *db.Queries) *CategoriaProdutoRepository {
	return &CategoriaProdutoRepository{queries: queries}
}

// =============================================================================
// CREATE
// =============================================================================

// Create persiste uma nova categoria de produto
func (r *CategoriaProdutoRepository) Create(ctx context.Context, categoria *entity.CategoriaProdutoEntity) error {
	centroCusto := string(categoria.CentroCusto)
	params := db.CreateCategoriaProdutoParams{
		ID:          uuidToPgUUID(categoria.ID),
		TenantID:    uuidToPgUUID(categoria.TenantID),
		Nome:        categoria.Nome,
		Descricao:   strPtr(categoria.Descricao),
		Cor:         strPtr(categoria.Cor),
		Icone:       strPtr(categoria.Icone),
		CentroCusto: &centroCusto,
		Ativa:       categoria.Ativa,
	}

	_, err := r.queries.CreateCategoriaProduto(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar categoria de produto: %w", err)
	}
	return nil
}

// =============================================================================
// READ
// =============================================================================

// FindByID busca categoria por ID
func (r *CategoriaProdutoRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.CategoriaProdutoEntity, error) {
	params := db.GetCategoriaProdutoByIDParams{
		ID:       uuidToPgUUID(id),
		TenantID: uuidToPgUUID(tenantID),
	}

	row, err := r.queries.GetCategoriaProdutoByID(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar categoria: %w", err)
	}

	return r.toDomain(row), nil
}

// FindByNome busca categoria por nome (único por tenant)
func (r *CategoriaProdutoRepository) FindByNome(ctx context.Context, tenantID uuid.UUID, nome string) (*entity.CategoriaProdutoEntity, error) {
	params := db.GetCategoriaProdutoByNomeParams{
		TenantID: uuidToPgUUID(tenantID),
		Lower:    nome,
	}

	row, err := r.queries.GetCategoriaProdutoByNome(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar categoria por nome: %w", err)
	}

	return r.toDomain(row), nil
}

// ListAll lista todas as categorias do tenant
func (r *CategoriaProdutoRepository) ListAll(ctx context.Context, tenantID uuid.UUID) ([]*entity.CategoriaProdutoEntity, error) {
	rows, err := r.queries.ListCategoriasProdutos(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar categorias: %w", err)
	}

	categorias := make([]*entity.CategoriaProdutoEntity, len(rows))
	for i, row := range rows {
		categorias[i] = r.toDomain(row)
	}

	return categorias, nil
}

// ListAtivas lista apenas categorias ativas do tenant
func (r *CategoriaProdutoRepository) ListAtivas(ctx context.Context, tenantID uuid.UUID) ([]*entity.CategoriaProdutoEntity, error) {
	rows, err := r.queries.ListCategoriasProdutosAtivas(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar categorias ativas: %w", err)
	}

	categorias := make([]*entity.CategoriaProdutoEntity, len(rows))
	for i, row := range rows {
		categorias[i] = r.toDomain(row)
	}

	return categorias, nil
}

// =============================================================================
// UPDATE
// =============================================================================

// Update atualiza uma categoria existente
func (r *CategoriaProdutoRepository) Update(ctx context.Context, categoria *entity.CategoriaProdutoEntity) error {
	centroCusto := string(categoria.CentroCusto)
	params := db.UpdateCategoriaProdutoParams{
		ID:          uuidToPgUUID(categoria.ID),
		TenantID:    uuidToPgUUID(categoria.TenantID),
		Nome:        categoria.Nome,
		Descricao:   strPtr(categoria.Descricao),
		Cor:         strPtr(categoria.Cor),
		Icone:       strPtr(categoria.Icone),
		CentroCusto: &centroCusto,
		Ativa:       categoria.Ativa,
	}

	_, err := r.queries.UpdateCategoriaProduto(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar categoria: %w", err)
	}

	return nil
}

// =============================================================================
// DELETE
// =============================================================================

// Delete remove uma categoria
func (r *CategoriaProdutoRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	params := db.DeleteCategoriaProdutoParams{
		ID:       uuidToPgUUID(id),
		TenantID: uuidToPgUUID(tenantID),
	}

	err := r.queries.DeleteCategoriaProduto(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao deletar categoria: %w", err)
	}

	return nil
}

// =============================================================================
// QUERIES AUXILIARES
// =============================================================================

// ExistsWithNome verifica se já existe categoria com o nome (para validação de duplicidade)
func (r *CategoriaProdutoRepository) ExistsWithNome(ctx context.Context, tenantID uuid.UUID, nome string, excludeID *uuid.UUID) (bool, error) {
	var idUUID pgtype.UUID

	if excludeID != nil {
		idUUID = uuidToPgUUID(*excludeID)
	} else {
		idUUID = pgtype.UUID{Valid: false}
	}

	params := db.CheckCategoriaProdutoNomeExistsParams{
		TenantID: uuidToPgUUID(tenantID),
		Lower:    nome,
		ID:       idUUID,
	}

	exists, err := r.queries.CheckCategoriaProdutoNomeExists(ctx, params)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar nome: %w", err)
	}

	return exists, nil
}

// CountProdutosVinculados conta quantos produtos estão vinculados à categoria
func (r *CategoriaProdutoRepository) CountProdutosVinculados(ctx context.Context, tenantID, categoriaID uuid.UUID) (int64, error) {
	params := db.CountProdutosByCategoriaParams{
		TenantID:           uuidToPgUUID(tenantID),
		CategoriaProdutoID: uuidToPgUUID(categoriaID),
	}

	count, err := r.queries.CountProdutosByCategoria(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar produtos: %w", err)
	}

	return count, nil
}

// =============================================================================
// MAPPERS
// =============================================================================

// toDomain converte db.CategoriasProduto para entity.CategoriaProdutoEntity
func (r *CategoriaProdutoRepository) toDomain(row db.CategoriasProduto) *entity.CategoriaProdutoEntity {
	return &entity.CategoriaProdutoEntity{
		ID:           pgUUIDToUUID(row.ID),
		TenantID:     pgUUIDToUUID(row.TenantID),
		Nome:         row.Nome,
		Descricao:    derefString(row.Descricao),
		Cor:          derefStringDefault(row.Cor, "#6B7280"),
		Icone:        derefStringDefault(row.Icone, "package"),
		CentroCusto:  entity.CentroCusto(derefStringDefault(row.CentroCusto, "CMV")),
		Ativa:        row.Ativa,
		CriadoEm:     row.CriadoEm.Time,
		AtualizadoEm: row.AtualizadoEm.Time,
	}
}

// strPtr converte string para *string, retorna nil se vazio
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// derefString dereference *string, retorna "" se nil
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// derefStringDefault dereference *string com valor default se nil
func derefStringDefault(s *string, defaultVal string) string {
	if s == nil {
		return defaultVal
	}
	return *s
}
