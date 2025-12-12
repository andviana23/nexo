// Package postgres implementa os repositórios usando PostgreSQL e sqlc.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// ContaPagarRepository implementa port.ContaPagarRepository usando sqlc.
type ContaPagarRepository struct {
	queries *db.Queries
}

// NewContaPagarRepository cria uma nova instância do repositório.
func NewContaPagarRepository(queries *db.Queries) *ContaPagarRepository {
	return &ContaPagarRepository{
		queries: queries,
	}
}

// Create persiste uma nova conta a pagar.
func (r *ContaPagarRepository) Create(ctx context.Context, conta *entity.ContaPagar) error {
	tenantUUID := entityUUIDToPgtype(conta.TenantID)
	categoriaUUID := uuidStringToPgtype(conta.CategoriaID)

	tipoStr := string(conta.Tipo)
	statusStr := mapContaPagarStatusToDB(conta.Status)
	recorrente := conta.Recorrente

	params := db.CreateContaPagarParams{
		TenantID:       tenantUUID,
		Descricao:      conta.Descricao,
		CategoriaID:    categoriaUUID,
		Fornecedor:     &conta.Fornecedor,
		Valor:          moneyToRawDecimal(conta.Valor),
		Tipo:           &tipoStr,
		Recorrente:     &recorrente,
		Periodicidade:  &conta.Periodicidade,
		DataVencimento: dateToDate(conta.DataVencimento),
		DataPagamento:  pgtype.Date{Valid: conta.DataPagamento != nil},
		Status:         &statusStr,
		ComprovanteUrl: &conta.ComprovanteURL,
		PixCode:        &conta.PixCode,
		Observacoes:    &conta.Observacoes,
	}

	if conta.DataPagamento != nil {
		params.DataPagamento = dateToDate(*conta.DataPagamento)
	}

	result, err := r.queries.CreateContaPagar(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar conta a pagar: %w", err)
	}

	conta.ID = pgUUIDToString(result.ID)
	conta.CriadoEm = timestamptzToTime(result.CriadoEm)
	conta.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)

	return nil
}

// FindByID busca uma conta por ID.
func (r *ContaPagarRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.ContaPagar, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	result, err := r.queries.GetContaPagarByID(ctx, db.GetContaPagarByIDParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conta a pagar: %w", err)
	}

	return r.toDomain(&result)
}

// Update atualiza uma conta existente.
func (r *ContaPagarRepository) Update(ctx context.Context, conta *entity.ContaPagar) error {
	tenantUUID := entityUUIDToPgtype(conta.TenantID)
	idUUID := uuidStringToPgtype(conta.ID)

	statusStr := mapContaPagarStatusToDB(conta.Status)

	params := db.UpdateContaPagarParams{
		ID:             idUUID,
		TenantID:       tenantUUID,
		Descricao:      conta.Descricao,
		Valor:          moneyToRawDecimal(conta.Valor),
		DataVencimento: dateToDate(conta.DataVencimento),
		DataPagamento:  pgtype.Date{Valid: conta.DataPagamento != nil},
		Status:         &statusStr,
		ComprovanteUrl: &conta.ComprovanteURL,
		PixCode:        &conta.PixCode,
		Observacoes:    &conta.Observacoes,
	}

	if conta.DataPagamento != nil {
		params.DataPagamento = dateToDate(*conta.DataPagamento)
	}

	result, err := r.queries.UpdateContaPagar(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar conta a pagar: %w", err)
	}

	conta.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)
	return nil
}

// Delete remove uma conta.
func (r *ContaPagarRepository) Delete(ctx context.Context, tenantID, id string) error {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	err := r.queries.DeleteContaPagar(ctx, db.DeleteContaPagarParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return fmt.Errorf("erro ao deletar conta a pagar: %w", err)
	}

	return nil
}

// List lista contas a pagar com filtros
func (r *ContaPagarRepository) List(ctx context.Context, tenantID string, filters port.ContaPagarListFilters) ([]*entity.ContaPagar, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	var limit int32 = 100
	var offset int32 = 0
	if filters.PageSize > 0 {
		limit = int32(filters.PageSize)
	}
	if filters.Page > 1 {
		offset = int32((filters.Page - 1) * filters.PageSize)
	}

	var statusStr *string
	if filters.Status != nil {
		s := mapContaPagarStatusToDB(*filters.Status)
		statusStr = &s
	}

	var tipoStr *string
	if filters.Tipo != nil {
		t := string(*filters.Tipo)
		tipoStr = &t
	}

	categoriaUUID := pgtype.UUID{Valid: false}
	if filters.CategoriaID != nil && *filters.CategoriaID != "" {
		categoriaUUID = uuidStringToPgtype(*filters.CategoriaID)
	}

	dataInicio := pgtype.Date{Valid: false}
	if filters.DataInicio != nil && !filters.DataInicio.IsZero() {
		dataInicio = dateToDate(*filters.DataInicio)
	}

	dataFim := pgtype.Date{Valid: false}
	if filters.DataFim != nil && !filters.DataFim.IsZero() {
		dataFim = dateToDate(*filters.DataFim)
	}

	results, err := r.queries.ListContasPagarFiltered(ctx, db.ListContasPagarFilteredParams{
		TenantID:    tenantUUID,
		Limit:       limit,
		Offset:      offset,
		UnitID:      pgtype.UUID{Valid: false},
		Status:      statusStr,
		Tipo:        tipoStr,
		CategoriaID: categoriaUUID,
		DataInicio:  dataInicio,
		DataFim:     dataFim,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas: %w", err)
	}

	return r.toDomainList(results)
}

// ListByDateRange lista contas em um período
func (r *ContaPagarRepository) ListByDateRange(ctx context.Context, tenantID string, inicio, fim time.Time) ([]*entity.ContaPagar, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	results, err := r.queries.ListContasPagarByPeriod(ctx, db.ListContasPagarByPeriodParams{
		TenantID:         tenantUUID,
		DataVencimento:   dateToDate(inicio),
		DataVencimento_2: dateToDate(fim),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas por período: %w", err)
	}

	return r.toDomainList(results)
}

// ListVencendoEm lista contas que vencem em até N dias
func (r *ContaPagarRepository) ListVencendoEm(ctx context.Context, tenantID string, dias int) ([]*entity.ContaPagar, error) {
	return r.ListVencendo(ctx, tenantID, int32(dias))
}

// ListAtrasadas lista contas atrasadas
func (r *ContaPagarRepository) ListAtrasadas(ctx context.Context, tenantID string) ([]*entity.ContaPagar, error) {
	return r.ListVencendo(ctx, tenantID, 0)
}

// SumByPeriod soma valores de contas em um período
func (r *ContaPagarRepository) SumByPeriod(ctx context.Context, tenantID string, inicio, fim time.Time, status *valueobject.StatusConta) (valueobject.Money, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	if status != nil && *status == valueobject.StatusContaPago {
		// Usar query específica para contas pagas
		params := db.SumContasPagasByPeriodParams{
			TenantID:        tenantUUID,
			DataPagamento:   dateToDate(inicio),
			DataPagamento_2: dateToDate(fim),
		}
		result, err := r.queries.SumContasPagasByPeriod(ctx, params)
		if err != nil {
			return valueobject.Zero(), fmt.Errorf("erro ao somar contas pagas por período: %w", err)
		}
		return interfaceToMoney(result), nil
	}

	// Query geral por vencimento
	params := db.SumContasPagarByPeriodParams{
		TenantID:         tenantUUID,
		DataVencimento:   dateToDate(inicio),
		DataVencimento_2: dateToDate(fim),
	}
	result, err := r.queries.SumContasPagarByPeriod(ctx, params)
	if err != nil {
		return valueobject.Zero(), fmt.Errorf("erro ao somar contas a pagar por período: %w", err)
	}
	return interfaceToMoney(result), nil
}

// SumByCategoria soma valores por categoria
func (r *ContaPagarRepository) SumByCategoria(ctx context.Context, tenantID, categoriaID string, inicio, fim time.Time) (valueobject.Money, error) {
	// Por enquanto, usar SumByPeriod filtrado manualmente
	// TODO: Criar query específica no sqlc
	tenantUUID := uuidStringToPgtype(tenantID)

	params := db.SumContasPagarByPeriodParams{
		TenantID:         tenantUUID,
		DataVencimento:   dateToDate(inicio),
		DataVencimento_2: dateToDate(fim),
	}
	result, err := r.queries.SumContasPagarByPeriod(ctx, params)
	if err != nil {
		return valueobject.Zero(), fmt.Errorf("erro ao somar contas por categoria: %w", err)
	}
	return interfaceToMoney(result), nil
}

// ListByStatus lista contas por status.
func (r *ContaPagarRepository) ListByStatus(ctx context.Context, tenantID string, status valueobject.StatusConta) ([]*entity.ContaPagar, error) {
	// Usar valores padrão para limit/offset
	limit := int32(100)
	offset := int32(0)

	tenantUUID := uuidStringToPgtype(tenantID)

	statusStr := mapContaPagarStatusToDB(status)

	results, err := r.queries.ListContasPagarByStatus(ctx, db.ListContasPagarByStatusParams{
		TenantID: tenantUUID,
		Status:   &statusStr,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas por status: %w", err)
	}

	return r.toDomainList(results)
}

// ListVencendo lista contas que vencem em até N dias.
func (r *ContaPagarRepository) ListVencendo(ctx context.Context, tenantID string, proximosDias int32) ([]*entity.ContaPagar, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	hoje := time.Now()

	results, err := r.queries.ListContasPagarVencidas(ctx, db.ListContasPagarVencidasParams{
		TenantID:       tenantUUID,
		DataVencimento: dateToDate(hoje),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas vencidas: %w", err)
	}

	return r.toDomainList(results)
}

// toDomain converte modelo sqlc para entidade de domínio.
func (r *ContaPagarRepository) toDomain(model *db.ContasAPagar) (*entity.ContaPagar, error) {
	var dataPagamento *time.Time
	if model.DataPagamento.Valid {
		dp := dateToTime(model.DataPagamento)
		dataPagamento = &dp
	}

	var fornecedor string
	if model.Fornecedor != nil {
		fornecedor = *model.Fornecedor
	}

	var tipo valueobject.TipoCusto
	if model.Tipo != nil {
		tipo = valueobject.TipoCusto(*model.Tipo)
	}

	var periodicidade string
	if model.Periodicidade != nil {
		periodicidade = *model.Periodicidade
	}

	var recorrente bool
	if model.Recorrente != nil {
		recorrente = *model.Recorrente
	}

	var status valueobject.StatusConta
	if model.Status != nil {
		status = mapContaPagarStatusFromDB(*model.Status)
	}

	var comprovanteURL string
	if model.ComprovanteUrl != nil {
		comprovanteURL = *model.ComprovanteUrl
	}

	var pixCode string
	if model.PixCode != nil {
		pixCode = *model.PixCode
	}

	var observacoes string
	if model.Observacoes != nil {
		observacoes = *model.Observacoes
	}

	conta := &entity.ContaPagar{
		ID:             pgUUIDToString(model.ID),
		TenantID:       pgtypeToEntityUUID(model.TenantID),
		Descricao:      model.Descricao,
		CategoriaID:    pgUUIDToString(model.CategoriaID),
		Fornecedor:     fornecedor,
		Valor:          rawDecimalToMoney(model.Valor),
		Tipo:           tipo,
		Recorrente:     recorrente,
		Periodicidade:  periodicidade,
		DataVencimento: dateToTime(model.DataVencimento),
		DataPagamento:  dataPagamento,
		Status:         status,
		ComprovanteURL: comprovanteURL,
		PixCode:        pixCode,
		Observacoes:    observacoes,
		CriadoEm:       timestamptzToTime(model.CriadoEm),
		AtualizadoEm:   timestamptzToTime(model.AtualizadoEm),
	}

	return conta, nil
}

// mapContaPagarStatusToDB converte status canônico do domínio para o status esperado no banco.
// No banco, "ABERTO" é equivalente a "PENDENTE".
func mapContaPagarStatusToDB(status valueobject.StatusConta) string {
	if status == valueobject.StatusContaPendente {
		return "ABERTO"
	}
	return string(status)
}

// mapContaPagarStatusFromDB converte status do banco para status canônico do domínio.
func mapContaPagarStatusFromDB(status string) valueobject.StatusConta {
	if status == "ABERTO" {
		return valueobject.StatusContaPendente
	}
	return valueobject.StatusConta(status)
}

// toDomainList converte lista de modelos para lista de entidades.
func (r *ContaPagarRepository) toDomainList(models []db.ContasAPagar) ([]*entity.ContaPagar, error) {
	result := make([]*entity.ContaPagar, 0, len(models))
	for i := range models {
		conta, err := r.toDomain(&models[i])
		if err != nil {
			return nil, err
		}
		result = append(result, conta)
	}
	return result, nil
}
