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

// ContaReceberRepository implementa port.ContaReceberRepository usando sqlc.
type ContaReceberRepository struct {
	queries *db.Queries
}

// NewContaReceberRepository cria uma nova instância do repositório.
func NewContaReceberRepository(queries *db.Queries) *ContaReceberRepository {
	return &ContaReceberRepository{
		queries: queries,
	}
}

// Create persiste uma nova conta a receber.
func (r *ContaReceberRepository) Create(ctx context.Context, conta *entity.ContaReceber) error {
	tenantUUID := uuidStringToPgtype(conta.TenantID)

	origemStr := conta.Origem
	statusStr := string(conta.Status)

	var assinaturaUUID pgtype.UUID
	if conta.AssinaturaID != nil {
		assinaturaUUID = uuidStringToPgtype(*conta.AssinaturaID)
	}

	params := db.CreateContaReceberParams{
		TenantID:        tenantUUID,
		Origem:          &origemStr,
		AssinaturaID:    assinaturaUUID,
		ServicoID:       pgtype.UUID{Valid: false},
		Descricao:       conta.DescricaoOrigem,
		Valor:           moneyToRawDecimal(conta.Valor),
		ValorPago:       moneyToNumeric(conta.ValorPago),
		DataVencimento:  dateToDate(conta.DataVencimento),
		DataRecebimento: pgtype.Date{Valid: conta.DataRecebimento != nil},
		Status:          &statusStr,
		Observacoes:     &conta.Observacoes,
	}

	if conta.DataRecebimento != nil {
		params.DataRecebimento = dateToDate(*conta.DataRecebimento)
	}

	result, err := r.queries.CreateContaReceber(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar conta a receber: %w", err)
	}

	conta.ID = pgUUIDToString(result.ID)
	conta.CriadoEm = timestamptzToTime(result.CriadoEm)
	conta.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)

	return nil
}

// FindByID busca uma conta por ID.
func (r *ContaReceberRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.ContaReceber, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	result, err := r.queries.GetContaReceberByID(ctx, db.GetContaReceberByIDParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conta a receber: %w", err)
	}

	return r.toDomain(&result)
}

// Update atualiza uma conta existente.
func (r *ContaReceberRepository) Update(ctx context.Context, conta *entity.ContaReceber) error {
	tenantUUID := uuidStringToPgtype(conta.TenantID)
	idUUID := uuidStringToPgtype(conta.ID)

	statusStr := string(conta.Status)

	params := db.UpdateContaReceberParams{
		ID:              idUUID,
		TenantID:        tenantUUID,
		Descricao:       conta.DescricaoOrigem,
		Valor:           moneyToRawDecimal(conta.Valor),
		ValorPago:       moneyToNumeric(conta.ValorPago),
		DataVencimento:  dateToDate(conta.DataVencimento),
		DataRecebimento: pgtype.Date{Valid: conta.DataRecebimento != nil},
		Status:          &statusStr,
		Observacoes:     &conta.Observacoes,
	}

	if conta.DataRecebimento != nil {
		params.DataRecebimento = dateToDate(*conta.DataRecebimento)
	}

	result, err := r.queries.UpdateContaReceber(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar conta a receber: %w", err)
	}

	conta.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)
	return nil
}

// Delete remove uma conta.
func (r *ContaReceberRepository) Delete(ctx context.Context, tenantID, id string) error {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	err := r.queries.DeleteContaReceber(ctx, db.DeleteContaReceberParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return fmt.Errorf("erro ao deletar conta a receber: %w", err)
	}

	return nil
}

// ListByStatus lista contas por status.
func (r *ContaReceberRepository) ListByStatus(ctx context.Context, tenantID string, status valueobject.StatusConta) ([]*entity.ContaReceber, error) {
	// Usar valores padrão para limit/offset
	limit := int32(100)
	offset := int32(0)

	tenantUUID := uuidStringToPgtype(tenantID)

	statusStr := string(status)

	results, err := r.queries.ListContasReceberByStatus(ctx, db.ListContasReceberByStatusParams{
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

// ListAtrasadas lista contas vencidas (atrasadas).
func (r *ContaReceberRepository) ListAtrasadas(ctx context.Context, tenantID string) ([]*entity.ContaReceber, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	hoje := time.Now()

	results, err := r.queries.ListContasReceberVencidas(ctx, db.ListContasReceberVencidasParams{
		TenantID:       tenantUUID,
		DataVencimento: dateToDate(hoje),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas vencidas: %w", err)
	}

	return r.toDomainList(results)
}

// List lista contas a receber com filtros
func (r *ContaReceberRepository) List(ctx context.Context, tenantID string, filters port.ContaReceberListFilters) ([]*entity.ContaReceber, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	var limit int32 = 100
	var offset int32 = 0
	if filters.PageSize > 0 {
		limit = int32(filters.PageSize)
	}
	if filters.Page > 1 {
		offset = int32((filters.Page - 1) * filters.PageSize)
	}

	results, err := r.queries.ListContasReceberByTenant(ctx, db.ListContasReceberByTenantParams{
		TenantID: tenantUUID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas: %w", err)
	}

	return r.toDomainList(results)
}

// ListByDateRange lista contas em um período
func (r *ContaReceberRepository) ListByDateRange(ctx context.Context, tenantID string, inicio, fim time.Time) ([]*entity.ContaReceber, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	results, err := r.queries.ListContasReceberByPeriod(ctx, db.ListContasReceberByPeriodParams{
		TenantID:         tenantUUID,
		DataVencimento:   dateToDate(inicio),
		DataVencimento_2: dateToDate(fim),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas por período: %w", err)
	}

	return r.toDomainList(results)
}

// ListByAssinatura lista contas de uma assinatura
func (r *ContaReceberRepository) ListByAssinatura(ctx context.Context, tenantID, assinaturaID string) ([]*entity.ContaReceber, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	assUUID := uuidStringToPgtype(assinaturaID)

	results, err := r.queries.ListContasReceberByAssinatura(ctx, db.ListContasReceberByAssinaturaParams{
		TenantID:     tenantUUID,
		AssinaturaID: assUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas por assinatura: %w", err)
	}

	return r.toDomainList(results)
}

// ListVencendoEm lista contas que vencem em até N dias
func (r *ContaReceberRepository) ListVencendoEm(ctx context.Context, tenantID string, dias int) ([]*entity.ContaReceber, error) {
	return r.ListAtrasadas(ctx, tenantID) // Implementação simplificada
}

// SumByPeriod soma valores de contas em um período
func (r *ContaReceberRepository) SumByPeriod(ctx context.Context, tenantID string, inicio, fim time.Time, status *valueobject.StatusConta) (valueobject.Money, error) {
	return valueobject.Money{}, nil // TODO: Implementar
}

// SumByOrigem soma valores por origem
func (r *ContaReceberRepository) SumByOrigem(ctx context.Context, tenantID, origem string, inicio, fim time.Time) (valueobject.Money, error) {
	return valueobject.Money{}, nil // TODO: Implementar
}

// toDomain converte modelo sqlc para entidade de domínio.
func (r *ContaReceberRepository) toDomain(model *db.ContasAReceber) (*entity.ContaReceber, error) {
	var dataRecebimento *time.Time
	if model.DataRecebimento.Valid {
		dr := dateToTime(model.DataRecebimento)
		dataRecebimento = &dr
	}

	var assinaturaID *string
	if model.AssinaturaID.Valid {
		aid := pgUUIDToString(model.AssinaturaID)
		assinaturaID = &aid
	}

	var origem string
	if model.Origem != nil {
		origem = *model.Origem
	}

	var status valueobject.StatusConta
	if model.Status != nil {
		status = valueobject.StatusConta(*model.Status)
	}

	var observacoes string
	if model.Observacoes != nil {
		observacoes = *model.Observacoes
	}

	conta := &entity.ContaReceber{
		ID:              pgUUIDToString(model.ID),
		TenantID:        pgUUIDToString(model.TenantID),
		Origem:          origem,
		AssinaturaID:    assinaturaID,
		DescricaoOrigem: model.Descricao,
		Valor:           rawDecimalToMoney(model.Valor),
		ValorPago:       numericToMoney(model.ValorPago),
		ValorAberto:     valueobject.Zero(), // Calculado no domínio
		DataVencimento:  dateToTime(model.DataVencimento),
		DataRecebimento: dataRecebimento,
		Status:          status,
		Observacoes:     observacoes,
		CriadoEm:        timestamptzToTime(model.CriadoEm),
		AtualizadoEm:    timestamptzToTime(model.AtualizadoEm),
	}

	// Calcular valor em aberto
	conta.ValorAberto = conta.Valor.Sub(conta.ValorPago)

	return conta, nil
}

// toDomainList converte lista de modelos para lista de entidades.
func (r *ContaReceberRepository) toDomainList(models []db.ContasAReceber) ([]*entity.ContaReceber, error) {
	result := make([]*entity.ContaReceber, 0, len(models))
	for i := range models {
		conta, err := r.toDomain(&models[i])
		if err != nil {
			return nil, err
		}
		result = append(result, conta)
	}
	return result, nil
}
