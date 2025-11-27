package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// CategoriaServicoRepository implementa port.CategoriaServicoRepository usando PostgreSQL
type CategoriaServicoRepository struct {
	queries *db.Queries
}

// NewCategoriaServicoRepository cria uma nova instância
func NewCategoriaServicoRepository(queries *db.Queries) *CategoriaServicoRepository {
	return &CategoriaServicoRepository{queries: queries}
}

// =============================================================================
// CREATE
// =============================================================================

// Create persiste uma nova categoria de serviço
func (r *CategoriaServicoRepository) Create(ctx context.Context, categoria *entity.CategoriaServico) error {
	params := db.CreateCategoriaServicoParams{
		ID:        stringToUUID(categoria.ID.String()),
		TenantID:  stringToUUID(categoria.TenantID.String()),
		Nome:      categoria.Nome,
		Descricao: categoria.Descricao,
		Cor:       categoria.Cor,
		Icone:     categoria.Icone,
		Ativa:     &categoria.Ativa,
	}

	_, err := r.queries.CreateCategoriaServico(ctx, params)
	return err
}

// =============================================================================
// READ
// =============================================================================

// FindByID busca categoria por ID
func (r *CategoriaServicoRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.CategoriaServico, error) {
	params := db.GetCategoriaServicoByIDParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
	}

	row, err := r.queries.GetCategoriaServicoByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("categoria não encontrada: %w", err)
	}

	return mapDBToCategoria(row), nil
}

// List lista todas as categorias com filtros
func (r *CategoriaServicoRepository) List(ctx context.Context, tenantID string, filter port.CategoriaServicoFilter) ([]*entity.CategoriaServico, error) {
	var rows []db.CategoriasServico
	var err error

	tenantUUID := stringToUUID(tenantID)

	if filter.ApenasAtivas {
		rows, err = r.queries.ListCategoriasServicosAtivas(ctx, tenantUUID)
	} else {
		rows, err = r.queries.ListCategoriasServicos(ctx, tenantUUID)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao listar categorias: %w", err)
	}

	categorias := make([]*entity.CategoriaServico, 0, len(rows))
	for _, row := range rows {
		categorias = append(categorias, mapDBToCategoria(row))
	}

	return categorias, nil
}

// =============================================================================
// UPDATE
// =============================================================================

// Update atualiza uma categoria existente
func (r *CategoriaServicoRepository) Update(ctx context.Context, categoria *entity.CategoriaServico) error {
	params := db.UpdateCategoriaServicoParams{
		ID:        stringToUUID(categoria.ID.String()),
		TenantID:  stringToUUID(categoria.TenantID.String()),
		Nome:      categoria.Nome,
		Descricao: categoria.Descricao,
		Cor:       categoria.Cor,
		Icone:     categoria.Icone,
	}

	_, err := r.queries.UpdateCategoriaServico(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar categoria: %w", err)
	}

	return nil
}

// ToggleStatus ativa/desativa uma categoria
func (r *CategoriaServicoRepository) ToggleStatus(ctx context.Context, tenantID, id string, ativa bool) error {
	params := db.ToggleCategoriaServicoStatusParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
		Ativa:    &ativa,
	}

	_, err := r.queries.ToggleCategoriaServicoStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao alterar status da categoria: %w", err)
	}

	return nil
}

// =============================================================================
// DELETE
// =============================================================================

// Delete deleta uma categoria
func (r *CategoriaServicoRepository) Delete(ctx context.Context, tenantID, id string) error {
	params := db.DeleteCategoriaServicoParams{
		ID:       stringToUUID(id),
		TenantID: stringToUUID(tenantID),
	}

	err := r.queries.DeleteCategoriaServico(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao deletar categoria: %w", err)
	}

	return nil
}

// =============================================================================
// QUERIES AUXILIARES
// =============================================================================

// CheckNomeExists verifica se já existe categoria com o mesmo nome
func (r *CategoriaServicoRepository) CheckNomeExists(ctx context.Context, tenantID, nome, excludeID string) (bool, error) {
	var idUUID pgtype.UUID

	if excludeID != "" {
		idUUID = stringToUUID(excludeID)
	} else {
		// UUID nulo para novo registro
		idUUID = pgtype.UUID{Valid: false}
	}

	params := db.CheckCategoriaServicoNomeExistsParams{
		TenantID: stringToUUID(tenantID),
		Lower:    nome, // A query já faz LOWER()
		ID:       idUUID,
	}

	exists, err := r.queries.CheckCategoriaServicoNomeExists(ctx, params)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar nome: %w", err)
	}

	return exists, nil
}

// CountServicos conta quantos serviços estão vinculados à categoria
func (r *CategoriaServicoRepository) CountServicos(ctx context.Context, tenantID, categoriaID string) (int64, error) {
	params := db.CountServicosInCategoriaParams{
		CategoriaID: stringToUUID(categoriaID),
		TenantID:    stringToUUID(tenantID),
	}

	count, err := r.queries.CountServicosInCategoria(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar serviços: %w", err)
	}

	return count, nil
}

// =============================================================================
// MAPPERS
// =============================================================================

// mapDBToCategoria converte db.CategoriasServico para entity.CategoriaServico
func mapDBToCategoria(row db.CategoriasServico) *entity.CategoriaServico {
	categoria := &entity.CategoriaServico{
		ID:           row.ID.Bytes,
		TenantID:     row.TenantID.Bytes,
		Nome:         row.Nome,
		Descricao:    row.Descricao,
		Cor:          row.Cor,
		Icone:        row.Icone,
		Ativa:        row.Ativa != nil && *row.Ativa,
		CriadoEm:     row.CriadoEm.Time,
		AtualizadoEm: row.AtualizadoEm.Time,
	}

	return categoria
}
