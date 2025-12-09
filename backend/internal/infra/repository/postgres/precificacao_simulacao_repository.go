// Package postgres implementa os repositórios usando PostgreSQL e sqlc.
package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
)

// PrecificacaoSimulacaoRepository implementa port.PrecificacaoSimulacaoRepository usando sqlc.
type PrecificacaoSimulacaoRepository struct {
	queries *db.Queries
}

// NewPrecificacaoSimulacaoRepository cria uma nova instância do repositório.
func NewPrecificacaoSimulacaoRepository(queries *db.Queries) *PrecificacaoSimulacaoRepository {
	return &PrecificacaoSimulacaoRepository{
		queries: queries,
	}
}

// Create persiste uma nova simulação de precificação.
func (r *PrecificacaoSimulacaoRepository) Create(ctx context.Context, simulacao *entity.PrecificacaoSimulacao) error {
	tenantUUID := entityUUIDToPgtype(simulacao.TenantID)
	itemUUID := uuidStringToPgtype(simulacao.ItemID)

	// CriadoPor pode vir de contexto - usar UUID zero por padrão
	criadoPorUUID := uuidStringToPgtype("00000000-0000-0000-0000-000000000000")

	params := db.CreatePrecificacaoSimulacaoParams{
		TenantID:            tenantUUID,
		ItemID:              itemUUID,
		TipoItem:            &simulacao.TipoItem,
		CustoMateriais:      moneyToNumeric(simulacao.CustoMateriais),
		CustoMaoDeObra:      moneyToNumeric(simulacao.CustoMaoDeObra),
		CustoTotal:          moneyToNumeric(simulacao.CustoTotal),
		MargemDesejada:      percentageToRawDecimal(simulacao.MargemDesejada),
		ComissaoPercentual:  percentageToRawDecimal(simulacao.ComissaoPercentual),
		ImpostoPercentual:   percentageToRawDecimal(simulacao.ImpostoPercentual),
		PrecoSugerido:       moneyToRawDecimal(simulacao.PrecoSugerido),
		PrecoAtual:          moneyToNumeric(simulacao.PrecoAtual),
		DiferencaPercentual: percentageToNumeric(simulacao.DiferencaPercentual),
		LucroEstimado:       moneyToNumeric(simulacao.LucroEstimado),
		MargemFinal:         percentageToNumeric(simulacao.MargemFinal),
		ParametrosJson:      []byte(simulacao.ParametrosJSON),
		CriadoPor:           criadoPorUUID,
	}

	result, err := r.queries.CreatePrecificacaoSimulacao(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar simulação de precificação: %w", err)
	}

	simulacao.ID = pgUUIDToString(result.ID)
	simulacao.CriadoEm = timestamptzToTime(result.CriadoEm)

	return nil
}

// FindByID busca uma simulação por ID.
func (r *PrecificacaoSimulacaoRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.PrecificacaoSimulacao, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	params := db.GetPrecificacaoSimulacaoByIDParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	}

	model, err := r.queries.GetPrecificacaoSimulacaoByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar simulação: %w", err)
	}

	return r.toDomain(&model)
}

// ListByTenant lista todas as simulações de um tenant.
func (r *PrecificacaoSimulacaoRepository) ListByTenant(ctx context.Context, tenantID string, limit, offset int32) ([]*entity.PrecificacaoSimulacao, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	params := db.ListSimulacoesByTenantParams{
		TenantID: tenantUUID,
		Limit:    limit,
		Offset:   offset,
	}

	models, err := r.queries.ListSimulacoesByTenant(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar simulações: %w", err)
	}

	return r.toDomainList(models)
}

// ListByItemInternal lista simulações de um item específico (método interno).
func (r *PrecificacaoSimulacaoRepository) ListByItemInternal(ctx context.Context, tenantID, itemID string, limit, offset int32) ([]*entity.PrecificacaoSimulacao, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	itemUUID := uuidStringToPgtype(itemID)

	params := db.ListSimulacoesByItemParams{
		TenantID: tenantUUID,
		ItemID:   itemUUID,
		Limit:    limit,
		Offset:   offset,
	}

	models, err := r.queries.ListSimulacoesByItem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar simulações do item: %w", err)
	}

	return r.toDomainList(models)
}

// ListByItem lista simulações de um item específico com filtros (interface port).
func (r *PrecificacaoSimulacaoRepository) ListByItem(ctx context.Context, tenantID, itemID, tipoItem string, filters port.SimulacaoListFilters) ([]*entity.PrecificacaoSimulacao, error) {
	limit := int32(50)
	offset := int32(0)
	if filters.PageSize > 0 {
		limit = int32(filters.PageSize)
	}
	if filters.Page > 0 {
		offset = int32((filters.Page - 1) * filters.PageSize)
	}

	// tipoItem é informativo aqui - a query já filtra por itemID
	return r.ListByItemInternal(ctx, tenantID, itemID, limit, offset)
}

// ListByTipoItem lista simulações por tipo (SERVICO ou PRODUTO).
func (r *PrecificacaoSimulacaoRepository) ListByTipoItem(ctx context.Context, tenantID, tipoItem string, limit, offset int32) ([]*entity.PrecificacaoSimulacao, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	params := db.ListSimulacoesByTipoItemParams{
		TenantID: tenantUUID,
		TipoItem: &tipoItem,
		Limit:    limit,
		Offset:   offset,
	}

	models, err := r.queries.ListSimulacoesByTipoItem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar simulações por tipo: %w", err)
	}

	return r.toDomainList(models)
}

// GetUltimaByItem busca a última simulação de um item.
func (r *PrecificacaoSimulacaoRepository) GetUltimaByItem(ctx context.Context, tenantID, itemID string) (*entity.PrecificacaoSimulacao, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	itemUUID := uuidStringToPgtype(itemID)

	params := db.GetUltimaSimulacaoByItemParams{
		TenantID: tenantUUID,
		ItemID:   itemUUID,
	}

	model, err := r.queries.GetUltimaSimulacaoByItem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar última simulação: %w", err)
	}

	return r.toDomain(&model)
}

// Delete remove uma simulação.
func (r *PrecificacaoSimulacaoRepository) Delete(ctx context.Context, tenantID, id string) error {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	params := db.DeletePrecificacaoSimulacaoParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	}

	err := r.queries.DeletePrecificacaoSimulacao(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao deletar simulação: %w", err)
	}

	return nil
}

// List lista todas as simulações com filtros (adapta ListByTenant para interface port).
func (r *PrecificacaoSimulacaoRepository) List(ctx context.Context, tenantID string, filters port.SimulacaoListFilters) ([]*entity.PrecificacaoSimulacao, error) {
	limit := int32(50)
	offset := int32(0)
	if filters.PageSize > 0 {
		limit = int32(filters.PageSize)
	}
	if filters.Page > 0 {
		offset = int32((filters.Page - 1) * filters.PageSize)
	}

	if filters.TipoItem != nil && *filters.TipoItem != "" {
		return r.ListByTipoItem(ctx, tenantID, *filters.TipoItem, limit, offset)
	}

	return r.ListByTenant(ctx, tenantID, limit, offset)
}

// GetLatestByItem alias para GetUltimaByItem (compatibilidade com interface port).
func (r *PrecificacaoSimulacaoRepository) GetLatestByItem(ctx context.Context, tenantID, itemID, tipoItem string) (*entity.PrecificacaoSimulacao, error) {
	return r.GetUltimaByItem(ctx, tenantID, itemID)
}

// Update atualiza uma simulação (stub - simulações geralmente são imutáveis).
func (r *PrecificacaoSimulacaoRepository) Update(ctx context.Context, simulacao *entity.PrecificacaoSimulacao) error {
	// Simulações são geralmente imutáveis - não atualizamos, apenas criamos novas
	return fmt.Errorf("atualização de simulações não implementada - crie uma nova simulação")
}

// toDomain converte modelo sqlc para entidade de domínio.
func (r *PrecificacaoSimulacaoRepository) toDomain(model *db.PrecificacaoSimulaco) (*entity.PrecificacaoSimulacao, error) {
	var tipoItem string
	if model.TipoItem != nil {
		tipoItem = *model.TipoItem
	}

	margemDesejada, err := rawDecimalToPercentage(model.MargemDesejada)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter margem_desejada: %w", err)
	}

	comissaoPercentual, err := rawDecimalToPercentage(model.ComissaoPercentual)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter comissao_percentual: %w", err)
	}

	impostoPercentual, err := rawDecimalToPercentage(model.ImpostoPercentual)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter imposto_percentual: %w", err)
	}

	diferencaPercentual, err := numericToPercentage(model.DiferencaPercentual)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter diferenca_percentual: %w", err)
	}

	margemFinal, err := numericToPercentage(model.MargemFinal)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter margem_final: %w", err)
	}

	simulacao := &entity.PrecificacaoSimulacao{
		ID:                  pgUUIDToString(model.ID),
		TenantID: pgtypeToEntityUUID(model.TenantID),
		ItemID:              pgUUIDToString(model.ItemID),
		TipoItem:            tipoItem,
		CustoMateriais:      numericToMoney(model.CustoMateriais),
		CustoMaoDeObra:      numericToMoney(model.CustoMaoDeObra),
		CustoTotal:          numericToMoney(model.CustoTotal),
		MargemDesejada:      margemDesejada,
		ComissaoPercentual:  comissaoPercentual,
		ImpostoPercentual:   impostoPercentual,
		PrecoSugerido:       rawDecimalToMoney(model.PrecoSugerido),
		PrecoAtual:          numericToMoney(model.PrecoAtual),
		DiferencaPercentual: diferencaPercentual,
		LucroEstimado:       numericToMoney(model.LucroEstimado),
		MargemFinal:         margemFinal,
		ParametrosJSON:      string(model.ParametrosJson),
		CriadoEm:            timestamptzToTime(model.CriadoEm),
	}

	return simulacao, nil
}

// toDomainList converte lista de modelos para lista de entidades.
func (r *PrecificacaoSimulacaoRepository) toDomainList(models []db.PrecificacaoSimulaco) ([]*entity.PrecificacaoSimulacao, error) {
	result := make([]*entity.PrecificacaoSimulacao, 0, len(models))

	for i := range models {
		simulacao, err := r.toDomain(&models[i])
		if err != nil {
			return nil, err
		}
		result = append(result, simulacao)
	}

	return result, nil
}
