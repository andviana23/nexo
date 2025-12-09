// Package postgres implementa os repositórios usando PostgreSQL e sqlc.
package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/jackc/pgx/v5"
)

// DespesaFixaRepository implementa port.DespesaFixaRepository usando sqlc.
type DespesaFixaRepository struct {
	queries *db.Queries
}

// NewDespesaFixaRepository cria uma nova instância do repositório.
func NewDespesaFixaRepository(queries *db.Queries) *DespesaFixaRepository {
	return &DespesaFixaRepository{
		queries: queries,
	}
}

// Create persiste uma nova despesa fixa.
func (r *DespesaFixaRepository) Create(ctx context.Context, despesa *entity.DespesaFixa) error {
	tenantUUID := entityUUIDToPgtype(despesa.TenantID)
	unidadeUUID := uuidStringToPgtype(despesa.UnidadeID)
	categoriaUUID := uuidStringToPgtype(despesa.CategoriaID)

	params := db.CreateDespesaFixaParams{
		TenantID:      tenantUUID,
		UnidadeID:     unidadeUUID,
		Descricao:     despesa.Descricao,
		CategoriaID:   categoriaUUID,
		Fornecedor:    strPtrToNullable(despesa.Fornecedor),
		Valor:         moneyToRawDecimal(despesa.Valor),
		DiaVencimento: int32(despesa.DiaVencimento),
		Ativo:         &despesa.Ativo,
		Observacoes:   strPtrToNullable(despesa.Observacoes),
	}

	result, err := r.queries.CreateDespesaFixa(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar despesa fixa: %w", err)
	}

	despesa.ID = pgUUIDToString(result.ID)
	despesa.CriadoEm = timestamptzToTime(result.CriadoEm)
	despesa.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)

	return nil
}

// FindByID busca uma despesa fixa por ID.
func (r *DespesaFixaRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.DespesaFixa, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	result, err := r.queries.GetDespesaFixaByID(ctx, db.GetDespesaFixaByIDParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDespesaFixaNotFound
		}
		return nil, fmt.Errorf("erro ao buscar despesa fixa: %w", err)
	}

	return r.toDomain(&result), nil
}

// Update atualiza uma despesa fixa existente.
func (r *DespesaFixaRepository) Update(ctx context.Context, despesa *entity.DespesaFixa) error {
	tenantUUID := entityUUIDToPgtype(despesa.TenantID)
	idUUID := uuidStringToPgtype(despesa.ID)
	unidadeUUID := uuidStringToPgtype(despesa.UnidadeID)
	categoriaUUID := uuidStringToPgtype(despesa.CategoriaID)

	params := db.UpdateDespesaFixaParams{
		ID:            idUUID,
		TenantID:      tenantUUID,
		Descricao:     despesa.Descricao,
		CategoriaID:   categoriaUUID,
		Fornecedor:    strPtrToNullable(despesa.Fornecedor),
		Valor:         moneyToRawDecimal(despesa.Valor),
		DiaVencimento: int32(despesa.DiaVencimento),
		UnidadeID:     unidadeUUID,
		Observacoes:   strPtrToNullable(despesa.Observacoes),
	}

	result, err := r.queries.UpdateDespesaFixa(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrDespesaFixaNotFound
		}
		return fmt.Errorf("erro ao atualizar despesa fixa: %w", err)
	}

	despesa.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)
	return nil
}

// Delete remove uma despesa fixa.
func (r *DespesaFixaRepository) Delete(ctx context.Context, tenantID, id string) error {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	err := r.queries.DeleteDespesaFixa(ctx, db.DeleteDespesaFixaParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return fmt.Errorf("erro ao deletar despesa fixa: %w", err)
	}

	return nil
}

// Toggle alterna o status ativo/inativo.
func (r *DespesaFixaRepository) Toggle(ctx context.Context, tenantID, id string) (*entity.DespesaFixa, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	result, err := r.queries.ToggleDespesaFixa(ctx, db.ToggleDespesaFixaParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDespesaFixaNotFound
		}
		return nil, fmt.Errorf("erro ao alternar despesa fixa: %w", err)
	}

	return r.toDomain(&result), nil
}

// List lista despesas fixas com paginação e filtros.
func (r *DespesaFixaRepository) List(ctx context.Context, tenantID string, filters port.DespesaFixaListFilters) ([]*entity.DespesaFixa, int64, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	// Calcula offset
	offset := int32(0)
	limit := int32(20) // default
	if filters.PageSize > 0 {
		limit = int32(filters.PageSize)
	}
	if filters.Page > 1 {
		offset = int32((filters.Page - 1) * filters.PageSize)
	}

	// Busca lista
	results, err := r.queries.ListDespesasFixasByTenant(ctx, db.ListDespesasFixasByTenantParams{
		TenantID: tenantUUID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao listar despesas fixas: %w", err)
	}

	// Busca total
	total, err := r.queries.CountDespesasFixas(ctx, tenantUUID)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar despesas fixas: %w", err)
	}

	despesas := make([]*entity.DespesaFixa, 0, len(results))
	for _, row := range results {
		despesas = append(despesas, r.toDomain(&row))
	}

	return despesas, total, nil
}

// ListAtivas lista apenas despesas fixas ativas.
func (r *DespesaFixaRepository) ListAtivas(ctx context.Context, tenantID string) ([]*entity.DespesaFixa, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	results, err := r.queries.ListDespesasFixasAtivas(ctx, tenantUUID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar despesas fixas ativas: %w", err)
	}

	despesas := make([]*entity.DespesaFixa, 0, len(results))
	for _, row := range results {
		despesas = append(despesas, r.toDomain(&row))
	}

	return despesas, nil
}

// ListByUnidade lista despesas fixas de uma unidade específica.
func (r *DespesaFixaRepository) ListByUnidade(ctx context.Context, tenantID, unidadeID string) ([]*entity.DespesaFixa, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	unidadeUUID := uuidStringToPgtype(unidadeID)

	results, err := r.queries.ListDespesasFixasByUnidade(ctx, db.ListDespesasFixasByUnidadeParams{
		TenantID:  tenantUUID,
		UnidadeID: unidadeUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar despesas fixas por unidade: %w", err)
	}

	despesas := make([]*entity.DespesaFixa, 0, len(results))
	for _, row := range results {
		despesas = append(despesas, r.toDomain(&row))
	}

	return despesas, nil
}

// ListByCategoria lista despesas fixas de uma categoria específica.
func (r *DespesaFixaRepository) ListByCategoria(ctx context.Context, tenantID, categoriaID string) ([]*entity.DespesaFixa, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	categoriaUUID := uuidStringToPgtype(categoriaID)

	results, err := r.queries.ListDespesasFixasByCategoria(ctx, db.ListDespesasFixasByCategoriaParams{
		TenantID:    tenantUUID,
		CategoriaID: categoriaUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar despesas fixas por categoria: %w", err)
	}

	despesas := make([]*entity.DespesaFixa, 0, len(results))
	for _, row := range results {
		despesas = append(despesas, r.toDomain(&row))
	}

	return despesas, nil
}

// ListAtivasPorTenants lista todas as despesas ativas de todos os tenants ativos.
func (r *DespesaFixaRepository) ListAtivasPorTenants(ctx context.Context) ([]*port.DespesaFixaComTenant, error) {
	results, err := r.queries.ListDespesasFixasAtivasPorTenants(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar despesas fixas por tenants: %w", err)
	}

	despesas := make([]*port.DespesaFixaComTenant, 0, len(results))
	for _, row := range results {
		despesa := r.toDomainFromAtivasPorTenants(&row)
		despesas = append(despesas, &port.DespesaFixaComTenant{
			DespesaFixa: despesa,
			TenantNome:  row.TenantNome,
		})
	}

	return despesas, nil
}

// SumAtivas soma o valor total de despesas fixas ativas.
func (r *DespesaFixaRepository) SumAtivas(ctx context.Context, tenantID string) (valueobject.Money, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	result, err := r.queries.SumDespesasFixasAtivas(ctx, tenantUUID)
	if err != nil {
		return valueobject.Money{}, fmt.Errorf("erro ao somar despesas fixas ativas: %w", err)
	}

	return rawDecimalToMoney(result), nil
}

// SumByUnidade soma o valor total de despesas fixas ativas por unidade.
func (r *DespesaFixaRepository) SumByUnidade(ctx context.Context, tenantID, unidadeID string) (valueobject.Money, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	unidadeUUID := uuidStringToPgtype(unidadeID)

	result, err := r.queries.SumDespesasFixasByUnidade(ctx, db.SumDespesasFixasByUnidadeParams{
		TenantID:  tenantUUID,
		UnidadeID: unidadeUUID,
	})
	if err != nil {
		return valueobject.Money{}, fmt.Errorf("erro ao somar despesas fixas por unidade: %w", err)
	}

	return rawDecimalToMoney(result), nil
}

// Count conta total de despesas fixas do tenant.
func (r *DespesaFixaRepository) Count(ctx context.Context, tenantID string) (int64, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	count, err := r.queries.CountDespesasFixas(ctx, tenantUUID)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar despesas fixas: %w", err)
	}

	return count, nil
}

// CountAtivas conta despesas fixas ativas do tenant.
func (r *DespesaFixaRepository) CountAtivas(ctx context.Context, tenantID string) (int64, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	count, err := r.queries.CountDespesasFixasAtivas(ctx, tenantUUID)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar despesas fixas ativas: %w", err)
	}

	return count, nil
}

// ExistsByDescricao verifica se já existe despesa fixa com mesma descrição.
func (r *DespesaFixaRepository) ExistsByDescricao(ctx context.Context, tenantID, descricao string, excludeID *string) (bool, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	var excludeUUID = uuidStringToPgtype("")
	if excludeID != nil {
		excludeUUID = uuidStringToPgtype(*excludeID)
	}

	exists, err := r.queries.ExistsDespesaFixaByDescricao(ctx, db.ExistsDespesaFixaByDescricaoParams{
		TenantID: tenantUUID,
		Lower:    descricao,
		ID:       excludeUUID,
	})
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência de despesa fixa: %w", err)
	}

	return exists, nil
}

// toDomain converte o resultado do banco para a entidade de domínio.
func (r *DespesaFixaRepository) toDomain(row *db.DespesasFixa) *entity.DespesaFixa {
	ativo := true
	if row.Ativo != nil {
		ativo = *row.Ativo
	}

	return &entity.DespesaFixa{
		ID:            pgUUIDToString(row.ID),
		TenantID: pgtypeToEntityUUID(row.TenantID),
		UnidadeID:     pgUUIDToString(row.UnidadeID),
		Descricao:     row.Descricao,
		CategoriaID:   pgUUIDToString(row.CategoriaID),
		Fornecedor:    strFromNullable(row.Fornecedor),
		Valor:         rawDecimalToMoney(row.Valor),
		DiaVencimento: int(row.DiaVencimento),
		Ativo:         ativo,
		Observacoes:   strFromNullable(row.Observacoes),
		CriadoEm:      timestamptzToTime(row.CriadoEm),
		AtualizadoEm:  timestamptzToTime(row.AtualizadoEm),
	}
}

// toDomainFromAtivasPorTenants converte o resultado da query de múltiplos tenants.
func (r *DespesaFixaRepository) toDomainFromAtivasPorTenants(row *db.ListDespesasFixasAtivasPorTenantsRow) *entity.DespesaFixa {
	ativo := true
	if row.Ativo != nil {
		ativo = *row.Ativo
	}

	return &entity.DespesaFixa{
		ID:            pgUUIDToString(row.ID),
		TenantID: pgtypeToEntityUUID(row.TenantID),
		UnidadeID:     pgUUIDToString(row.UnidadeID),
		Descricao:     row.Descricao,
		CategoriaID:   pgUUIDToString(row.CategoriaID),
		Fornecedor:    strFromNullable(row.Fornecedor),
		Valor:         rawDecimalToMoney(row.Valor),
		DiaVencimento: int(row.DiaVencimento),
		Ativo:         ativo,
		Observacoes:   strFromNullable(row.Observacoes),
		CriadoEm:      timestamptzToTime(row.CriadoEm),
		AtualizadoEm:  timestamptzToTime(row.AtualizadoEm),
	}
}

// strPtrToNullable converte string para *string (nil se vazio).
func strPtrToNullable(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// strFromNullable converte *string para string.
func strFromNullable(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
