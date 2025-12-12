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
	tenantUUID := entityUUIDToPgtype(conta.TenantID)

	origemStr := conta.Origem
	statusStr := mapContaReceberStatusToDB(conta.Status)

	var assinaturaUUID pgtype.UUID
	if conta.AssinaturaID != nil {
		assinaturaUUID = uuidStringToPgtype(*conta.AssinaturaID)
	}

	var commandUUID pgtype.UUID
	if conta.CommandID != nil {
		commandUUID = uuidStringToPgtype(*conta.CommandID)
	}

	var commandPaymentUUID pgtype.UUID
	if conta.CommandPaymentID != nil {
		commandPaymentUUID = uuidStringToPgtype(*conta.CommandPaymentID)
	}

	params := db.CreateContaReceberParams{
		TenantID:         tenantUUID,
		Origem:           &origemStr,
		AssinaturaID:     assinaturaUUID,
		ServicoID:        pgtype.UUID{Valid: false},
		CommandID:        commandUUID,
		CommandPaymentID: commandPaymentUUID,
		Descricao:        conta.DescricaoOrigem,
		Valor:            moneyToRawDecimal(conta.Valor),
		ValorPago:        moneyToNumeric(conta.ValorPago),
		DataVencimento:   dateToDate(conta.DataVencimento),
		DataRecebimento:  pgtype.Date{Valid: conta.DataRecebimento != nil},
		Status:           &statusStr,
		Observacoes:      &conta.Observacoes,
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
	tenantUUID := entityUUIDToPgtype(conta.TenantID)
	idUUID := uuidStringToPgtype(conta.ID)

	statusStr := mapContaReceberStatusToDB(conta.Status)

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

	var statusStr *string
	if filters.Status != nil {
		s := mapContaReceberStatusToDB(*filters.Status)
		statusStr = &s
	}

	var origemStr *string
	if filters.Origem != nil && *filters.Origem != "" {
		o := *filters.Origem
		origemStr = &o
	}

	assinaturaUUID := pgtype.UUID{Valid: false}
	if filters.AssinaturaID != nil && *filters.AssinaturaID != "" {
		assinaturaUUID = uuidStringToPgtype(*filters.AssinaturaID)
	}

	dataInicio := pgtype.Date{Valid: false}
	if filters.DataInicio != nil && !filters.DataInicio.IsZero() {
		dataInicio = dateToDate(*filters.DataInicio)
	}

	dataFim := pgtype.Date{Valid: false}
	if filters.DataFim != nil && !filters.DataFim.IsZero() {
		dataFim = dateToDate(*filters.DataFim)
	}

	results, err := r.queries.ListContasReceberFiltered(ctx, db.ListContasReceberFilteredParams{
		TenantID:     tenantUUID,
		Limit:        limit,
		Offset:       offset,
		Status:       statusStr,
		Origem:       origemStr,
		AssinaturaID: assinaturaUUID,
		DataInicio:   dataInicio,
		DataFim:      dataFim,
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
	tenantUUID := uuidStringToPgtype(tenantID)

	if status != nil && (*status == valueobject.StatusContaRecebido || *status == valueobject.StatusContaPago) {
		// Usar query específica para contas recebidas
		params := db.SumContasRecebidasByPeriodParams{
			TenantID:          tenantUUID,
			DataRecebimento:   dateToDate(inicio),
			DataRecebimento_2: dateToDate(fim),
		}
		result, err := r.queries.SumContasRecebidasByPeriod(ctx, params)
		if err != nil {
			return valueobject.Zero(), fmt.Errorf("erro ao somar contas recebidas por período: %w", err)
		}
		return interfaceToMoney(result), nil
	}

	// Query geral por vencimento
	params := db.SumContasReceberByPeriodParams{
		TenantID:         tenantUUID,
		DataVencimento:   dateToDate(inicio),
		DataVencimento_2: dateToDate(fim),
	}
	result, err := r.queries.SumContasReceberByPeriod(ctx, params)
	if err != nil {
		return valueobject.Zero(), fmt.Errorf("erro ao somar contas a receber por período: %w", err)
	}
	return interfaceToMoney(result), nil
}

// SumByOrigem soma valores por origem
func (r *ContaReceberRepository) SumByOrigem(ctx context.Context, tenantID, origem string, inicio, fim time.Time) (valueobject.Money, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	origemStr := origem
	params := db.SumContasReceberByOrigemParams{
		TenantID:         tenantUUID,
		Origem:           &origemStr,
		DataVencimento:   dateToDate(inicio),
		DataVencimento_2: dateToDate(fim),
	}
	result, err := r.queries.SumContasReceberByOrigem(ctx, params)
	if err != nil {
		return valueobject.Zero(), fmt.Errorf("erro ao somar contas por origem: %w", err)
	}
	return interfaceToMoney(result), nil
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
		status = mapContaReceberStatusFromDB(*model.Status)
	}

	var observacoes string
	if model.Observacoes != nil {
		observacoes = *model.Observacoes
	}

	var commandID *string
	if model.CommandID.Valid {
		cid := pgUUIDToString(model.CommandID)
		commandID = &cid
	}

	var commandPaymentID *string
	if model.CommandPaymentID.Valid {
		cpid := pgUUIDToString(model.CommandPaymentID)
		commandPaymentID = &cpid
	}

	conta := &entity.ContaReceber{
		ID:               pgUUIDToString(model.ID),
		TenantID:         pgtypeToEntityUUID(model.TenantID),
		Origem:           origem,
		AssinaturaID:     assinaturaID,
		DescricaoOrigem:  model.Descricao,
		Valor:            rawDecimalToMoney(model.Valor),
		ValorPago:        numericToMoney(model.ValorPago),
		ValorAberto:      valueobject.Zero(), // Calculado no domínio
		DataVencimento:   dateToTime(model.DataVencimento),
		DataRecebimento:  dataRecebimento,
		Status:           status,
		Observacoes:      observacoes,
		CommandID:        commandID,
		CommandPaymentID: commandPaymentID,
		CriadoEm:         timestamptzToTime(model.CriadoEm),
		AtualizadoEm:     timestamptzToTime(model.AtualizadoEm),
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

// ListByCommandID lista contas vinculadas a uma comanda (para idempotência/estorno).
func (r *ContaReceberRepository) ListByCommandID(ctx context.Context, tenantID, commandID string) ([]*entity.ContaReceber, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	commandUUID := uuidStringToPgtype(commandID)

	results, err := r.queries.ListContasReceberByCommandID(ctx, db.ListContasReceberByCommandIDParams{
		TenantID:  tenantUUID,
		CommandID: commandUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas por comanda: %w", err)
	}

	return r.toDomainList(results)
}

// ============================================================================
// Métodos de integração Asaas (Migration 041)
// ============================================================================

// UpsertByAsaasPaymentID cria ou atualiza conta via webhook (idempotente)
func (r *ContaReceberRepository) UpsertByAsaasPaymentID(ctx context.Context, conta *entity.ContaReceber) error {
	tenantUUID := entityUUIDToPgtype(conta.TenantID)
	origemStr := conta.Origem
	statusStr := string(conta.Status)

	var assinaturaUUID pgtype.UUID
	if conta.AssinaturaID != nil {
		assinaturaUUID = uuidStringToPgtype(*conta.AssinaturaID)
	}

	var subscriptionUUID pgtype.UUID
	if conta.SubscriptionID != nil {
		subscriptionUUID = uuidStringToPgtype(*conta.SubscriptionID)
	}

	params := db.UpsertContaReceberByAsaasPaymentIDParams{
		TenantID:        tenantUUID,
		Origem:          &origemStr,
		AssinaturaID:    assinaturaUUID,
		SubscriptionID:  subscriptionUUID,
		AsaasPaymentID:  conta.AsaasPaymentID,
		ServicoID:       pgtype.UUID{Valid: false},
		Descricao:       conta.DescricaoOrigem,
		Valor:           moneyToRawDecimal(conta.Valor),
		ValorPago:       moneyToNumeric(conta.ValorPago),
		DataVencimento:  dateToDate(conta.DataVencimento),
		DataRecebimento: pgtype.Date{Valid: conta.DataRecebimento != nil},
		CompetenciaMes:  conta.CompetenciaMes,
		ConfirmedAt:     timestamptzFromTimePtr(conta.ConfirmedAt),
		ReceivedAt:      timestamptzFromTimePtr(conta.ReceivedAt),
		Status:          &statusStr,
		Observacoes:     &conta.Observacoes,
	}

	if conta.DataRecebimento != nil {
		params.DataRecebimento = dateToDate(*conta.DataRecebimento)
	}

	result, err := r.queries.UpsertContaReceberByAsaasPaymentID(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao upsert conta a receber: %w", err)
	}

	conta.ID = pgUUIDToString(result.ID)
	conta.CriadoEm = timestamptzToTime(result.CriadoEm)
	conta.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)

	return nil
}

// GetByAsaasPaymentID busca conta pelo payment ID do Asaas
func (r *ContaReceberRepository) GetByAsaasPaymentID(ctx context.Context, tenantID, asaasPaymentID string) (*entity.ContaReceber, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	result, err := r.queries.GetContaReceberByAsaasPaymentID(ctx, db.GetContaReceberByAsaasPaymentIDParams{
		TenantID:       tenantUUID,
		AsaasPaymentID: &asaasPaymentID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar conta por asaas_payment_id: %w", err)
	}

	return r.toDomain(&result)
}

// ListBySubscriptionID lista contas de uma subscription
func (r *ContaReceberRepository) ListBySubscriptionID(ctx context.Context, tenantID, subscriptionID string) ([]*entity.ContaReceber, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	subUUID := uuidStringToPgtype(subscriptionID)

	results, err := r.queries.GetContaReceberBySubscriptionID(ctx, db.GetContaReceberBySubscriptionIDParams{
		TenantID:       tenantUUID,
		SubscriptionID: subUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas por subscription: %w", err)
	}

	return r.toDomainList(results)
}

// SumByCompetencia soma valores por competência (para DRE)
func (r *ContaReceberRepository) SumByCompetencia(ctx context.Context, tenantID, competenciaMes string, status *valueobject.StatusConta) (valueobject.Money, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	if status != nil {
		statusStr := mapContaReceberStatusToDB(*status)
		result, err := r.queries.SumContasReceberByCompetenciaAndStatus(ctx, db.SumContasReceberByCompetenciaAndStatusParams{
			TenantID:       tenantUUID,
			CompetenciaMes: &competenciaMes,
			Status:         &statusStr,
		})
		if err != nil {
			return valueobject.Zero(), fmt.Errorf("erro ao somar por competência e status: %w", err)
		}
		return rawDecimalToMoney(result), nil
	}

	result, err := r.queries.SumContasReceberByCompetencia(ctx, db.SumContasReceberByCompetenciaParams{
		TenantID:       tenantUUID,
		CompetenciaMes: &competenciaMes,
	})
	if err != nil {
		return valueobject.Zero(), fmt.Errorf("erro ao somar por competência: %w", err)
	}
	return rawDecimalToMoney(result.TotalBruto), nil
}

// mapContaReceberStatusToDB converte status do domínio para o status esperado no banco.
// Mantém compatibilidade com código legado que ainda usa "PAGO" para contas a receber.
func mapContaReceberStatusToDB(status valueobject.StatusConta) string {
	if status == valueobject.StatusContaPago {
		return string(valueobject.StatusContaRecebido)
	}
	return string(status)
}

// mapContaReceberStatusFromDB converte status do banco para status canônico do domínio.
func mapContaReceberStatusFromDB(status string) valueobject.StatusConta {
	if status == string(valueobject.StatusContaPago) {
		return valueobject.StatusContaRecebido
	}
	return valueobject.StatusConta(status)
}

// MarcarRecebidaViaAsaas quita conta quando webhook RECEIVED chegar
func (r *ContaReceberRepository) MarcarRecebidaViaAsaas(ctx context.Context, tenantID, asaasPaymentID string, dataRecebimento time.Time, valorPago valueobject.Money) error {
	tenantUUID := uuidStringToPgtype(tenantID)

	_, err := r.queries.MarcarContaReceberRecebidaViaAsaas(ctx, db.MarcarContaReceberRecebidaViaAsaasParams{
		TenantID:        tenantUUID,
		AsaasPaymentID:  &asaasPaymentID,
		DataRecebimento: dateToDate(dataRecebimento),
		ValorPago:       moneyToNumeric(valorPago),
	})
	if err != nil {
		return fmt.Errorf("erro ao marcar conta como recebida: %w", err)
	}

	return nil
}

// EstornarViaAsaas estorna conta quando webhook REFUNDED chegar
func (r *ContaReceberRepository) EstornarViaAsaas(ctx context.Context, tenantID, asaasPaymentID, observacao string) error {
	tenantUUID := uuidStringToPgtype(tenantID)

	_, err := r.queries.EstornarContaReceberViaAsaas(ctx, db.EstornarContaReceberViaAsaasParams{
		TenantID:       tenantUUID,
		AsaasPaymentID: &asaasPaymentID,
		Observacoes:    &observacao,
	})
	if err != nil {
		return fmt.Errorf("erro ao estornar conta: %w", err)
	}

	return nil
}

// SumByReceivedDate soma valores recebidos em um período (para fluxo de caixa)
func (r *ContaReceberRepository) SumByReceivedDate(ctx context.Context, tenantID string, inicio, fim time.Time) (valueobject.Money, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	result, err := r.queries.SumContasReceberByReceivedDate(ctx, db.SumContasReceberByReceivedDateParams{
		TenantID:     tenantUUID,
		ReceivedAt:   timestamptzFromTimePtr(&inicio),
		ReceivedAt_2: timestamptzFromTimePtr(&fim),
	})
	if err != nil {
		return valueobject.Zero(), fmt.Errorf("erro ao somar por data de recebimento: %w", err)
	}
	return rawDecimalToMoney(result), nil
}

// SumByConfirmedDate soma valores confirmados em um período (para DRE regime competência)
func (r *ContaReceberRepository) SumByConfirmedDate(ctx context.Context, tenantID string, inicio, fim time.Time) (valueobject.Money, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	result, err := r.queries.SumContasReceberByConfirmedDate(ctx, db.SumContasReceberByConfirmedDateParams{
		TenantID:      tenantUUID,
		ConfirmedAt:   timestamptzFromTimePtr(&inicio),
		ConfirmedAt_2: timestamptzFromTimePtr(&fim),
	})
	if err != nil {
		return valueobject.Zero(), fmt.Errorf("erro ao somar por data de confirmação: %w", err)
	}
	return rawDecimalToMoney(result), nil
}

// timestamptzFromTimePtr converte *time.Time para pgtype.Timestamptz
func timestamptzFromTimePtr(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}
