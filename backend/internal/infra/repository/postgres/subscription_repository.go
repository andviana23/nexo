package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SubscriptionRepositoryPG implementa port.SubscriptionRepository
type SubscriptionRepositoryPG struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

// NewSubscriptionRepository cria uma instância do repositório de assinaturas
func NewSubscriptionRepository(queries *db.Queries, pool *pgxpool.Pool) port.SubscriptionRepository {
	return &SubscriptionRepositoryPG{queries: queries, pool: pool}
}

// NewSubscriptionPaymentRepository cria uma instância do repositório de pagamentos de assinatura
// Usa o mesmo repositório interno, mas retorna a interface específica para pagamentos
func NewSubscriptionPaymentRepository(queries *db.Queries, pool *pgxpool.Pool) port.SubscriptionPaymentRepository {
	return &SubscriptionPaymentRepositoryPG{queries: queries}
}

// SubscriptionPaymentRepositoryPG implementa port.SubscriptionPaymentRepository
type SubscriptionPaymentRepositoryPG struct {
	queries *db.Queries
}

// Create registra novo pagamento de assinatura
func (r *SubscriptionPaymentRepositoryPG) Create(ctx context.Context, payment *entity.SubscriptionPayment) error {
	params := db.CreateSubscriptionPaymentParams{
		TenantID:        uuidToPgUUID(payment.TenantID),
		SubscriptionID:  uuidToPgUUID(payment.SubscriptionID),
		AsaasPaymentID:  payment.AsaasPaymentID,
		Valor:           payment.Valor,
		FormaPagamento:  string(payment.FormaPagamento),
		Status:          string(payment.Status),
		DataPagamento:   timestampToTimestamptzPtr(payment.DataPagamento),
		CodigoTransacao: payment.CodigoTransacao,
		Observacao:      payment.Observacao,
	}
	row, err := r.queries.CreateSubscriptionPayment(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar pagamento de assinatura: %w", err)
	}
	payment.ID = pgUUIDToUUID(row.ID)
	payment.CreatedAt = row.CreatedAt.Time
	return nil
}

// ListBySubscription lista pagamentos de uma assinatura
func (r *SubscriptionPaymentRepositoryPG) ListBySubscription(ctx context.Context, subscriptionID, tenantID uuid.UUID) ([]*entity.SubscriptionPayment, error) {
	rows, err := r.queries.ListPaymentsBySubscription(ctx, db.ListPaymentsBySubscriptionParams{
		SubscriptionID: uuidToPgUUID(subscriptionID),
		TenantID:       uuidToPgUUID(tenantID),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pagamentos: %w", err)
	}
	result := make([]*entity.SubscriptionPayment, 0, len(rows))
	for i := range rows {
		result = append(result, mapPaymentRowStatic(&rows[i]))
	}
	return result, nil
}

// UpdateStatus atualiza status do pagamento
func (r *SubscriptionPaymentRepositoryPG) UpdateStatus(ctx context.Context, paymentID, tenantID uuid.UUID, status entity.PaymentStatus, dataPagamento *time.Time) error {
	if status == entity.PaymentStatusConfirmado && dataPagamento != nil {
		return r.queries.ConfirmPayment(ctx, db.ConfirmPaymentParams{
			ID:            uuidToPgUUID(paymentID),
			TenantID:      uuidToPgUUID(tenantID),
			DataPagamento: timestampToTimestamptz(*dataPagamento),
		})
	}
	return r.queries.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
		ID:       uuidToPgUUID(paymentID),
		TenantID: uuidToPgUUID(tenantID),
		Status:   string(status),
	})
}

// GetByAsaasID busca pagamento pelo ID do Asaas
func (r *SubscriptionPaymentRepositoryPG) GetByAsaasID(ctx context.Context, asaasPaymentID string) (*entity.SubscriptionPayment, error) {
	row, err := r.queries.GetPaymentByAsaasID(ctx, &asaasPaymentID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar pagamento por asaas_id: %w", err)
	}
	return &entity.SubscriptionPayment{
		ID:              pgUUIDToUUID(row.ID),
		TenantID:        pgUUIDToUUID(row.TenantID),
		SubscriptionID:  pgUUIDToUUID(row.SubscriptionID),
		AsaasPaymentID:  row.AsaasPaymentID,
		Valor:           row.Valor,
		FormaPagamento:  entity.PaymentMethod(row.FormaPagamento),
		Status:          entity.PaymentStatus(row.Status),
		DataPagamento:   timestamptzToTimePtr(row.DataPagamento),
		CodigoTransacao: row.CodigoTransacao,
		Observacao:      row.Observacao,
		CreatedAt:       row.CreatedAt.Time,
	}, nil
}

// mapPaymentRowStatic mapeia row de pagamento (função estática)
func mapPaymentRowStatic(row *db.SubscriptionPayment) *entity.SubscriptionPayment {
	return &entity.SubscriptionPayment{
		ID:                  pgUUIDToUUID(row.ID),
		TenantID:            pgUUIDToUUID(row.TenantID),
		SubscriptionID:      pgUUIDToUUID(row.SubscriptionID),
		AsaasPaymentID:      row.AsaasPaymentID,
		Valor:               row.Valor,
		FormaPagamento:      entity.PaymentMethod(row.FormaPagamento),
		Status:              entity.PaymentStatus(row.Status),
		DataPagamento:       timestamptzToTimePtr(row.DataPagamento),
		CodigoTransacao:     row.CodigoTransacao,
		Observacao:          row.Observacao,
		CreatedAt:           row.CreatedAt.Time,
		StatusAsaas:         row.StatusAsaas,
		DueDate:             dateToTimePtr(row.DueDate),
		ConfirmedDate:       timestamptzToTimePtr(row.ConfirmedDate),
		ClientPaymentDate:   dateToTimePtr(row.ClientPaymentDate),
		CreditDate:          dateToTimePtr(row.CreditDate),
		EstimatedCreditDate: dateToTimePtr(row.EstimatedCreditDate),
		BillingType:         row.BillingType,
		NetValue:            numericToDecimal(row.NetValue),
		InvoiceURL:          row.InvoiceUrl,
		BankSlipURL:         row.BankSlipUrl,
		PixQRCode:           row.PixQrcode,
	}
}

// UpsertByAsaasID cria ou atualiza pagamento via webhook (idempotente)
func (r *SubscriptionPaymentRepositoryPG) UpsertByAsaasID(ctx context.Context, payment *entity.SubscriptionPayment) error {
	params := db.UpsertPaymentByAsaasIDParams{
		TenantID:            uuidToPgUUID(payment.TenantID),
		SubscriptionID:      uuidToPgUUID(payment.SubscriptionID),
		AsaasPaymentID:      payment.AsaasPaymentID,
		Valor:               payment.Valor,
		FormaPagamento:      string(payment.FormaPagamento),
		Status:              string(payment.Status),
		StatusAsaas:         payment.StatusAsaas,
		DueDate:             timePtrToDate(payment.DueDate),
		ConfirmedDate:       timestampToTimestamptzPtr(payment.ConfirmedDate),
		ClientPaymentDate:   timePtrToDate(payment.ClientPaymentDate),
		CreditDate:          timePtrToDate(payment.CreditDate),
		EstimatedCreditDate: timePtrToDate(payment.EstimatedCreditDate),
		BillingType:         payment.BillingType,
		NetValue:            decimalToNumeric(payment.NetValue),
		InvoiceUrl:          payment.InvoiceURL,
		BankSlipUrl:         payment.BankSlipURL,
		PixQrcode:           payment.PixQRCode,
		DataPagamento:       timestampToTimestamptzPtr(payment.DataPagamento),
		CodigoTransacao:     payment.CodigoTransacao,
		Observacao:          payment.Observacao,
	}

	row, err := r.queries.UpsertPaymentByAsaasID(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao upsert pagamento: %w", err)
	}
	payment.ID = pgUUIDToUUID(row.ID)
	payment.CreatedAt = row.CreatedAt.Time
	return nil
}

// UpdatePaymentConfirmed atualiza pagamento para CONFIRMADO
func (r *SubscriptionPaymentRepositoryPG) UpdatePaymentConfirmed(ctx context.Context, paymentID, tenantID uuid.UUID, statusAsaas string, confirmedDate time.Time) error {
	return r.queries.UpdatePaymentConfirmed(ctx, db.UpdatePaymentConfirmedParams{
		ID:            uuidToPgUUID(paymentID),
		TenantID:      uuidToPgUUID(tenantID),
		StatusAsaas:   &statusAsaas,
		ConfirmedDate: timestampToTimestamptz(confirmedDate),
	})
}

// UpdatePaymentReceived atualiza pagamento para RECEBIDO
func (r *SubscriptionPaymentRepositoryPG) UpdatePaymentReceived(ctx context.Context, paymentID, tenantID uuid.UUID, statusAsaas string, clientPaymentDate, creditDate time.Time) error {
	return r.queries.UpdatePaymentReceived(ctx, db.UpdatePaymentReceivedParams{
		ID:                uuidToPgUUID(paymentID),
		TenantID:          uuidToPgUUID(tenantID),
		StatusAsaas:       &statusAsaas,
		ClientPaymentDate: timeToDate(clientPaymentDate),
		CreditDate:        timeToDate(creditDate),
	})
}

// UpdatePaymentOverdue atualiza pagamento para VENCIDO
func (r *SubscriptionPaymentRepositoryPG) UpdatePaymentOverdue(ctx context.Context, paymentID, tenantID uuid.UUID, statusAsaas string) error {
	return r.queries.UpdatePaymentOverdue(ctx, db.UpdatePaymentOverdueParams{
		ID:          uuidToPgUUID(paymentID),
		TenantID:    uuidToPgUUID(tenantID),
		StatusAsaas: &statusAsaas,
	})
}

// UpdatePaymentRefunded atualiza pagamento para ESTORNADO
func (r *SubscriptionPaymentRepositoryPG) UpdatePaymentRefunded(ctx context.Context, paymentID, tenantID uuid.UUID, statusAsaas string) error {
	return r.queries.UpdatePaymentRefunded(ctx, db.UpdatePaymentRefundedParams{
		ID:          uuidToPgUUID(paymentID),
		TenantID:    uuidToPgUUID(tenantID),
		StatusAsaas: &statusAsaas,
	})
}

// ListNeedingReconciliation lista pagamentos que precisam de conciliação
func (r *SubscriptionPaymentRepositoryPG) ListNeedingReconciliation(ctx context.Context, tenantID string) ([]*entity.SubscriptionPayment, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}

	rows, err := r.queries.ListPaymentsNeedingReconciliation(ctx, uuidToPgUUID(tenantUUID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pagamentos para conciliação: %w", err)
	}

	result := make([]*entity.SubscriptionPayment, 0, len(rows))
	for i := range rows {
		result = append(result, mapPaymentRowStatic(&rows[i]))
	}
	return result, nil
}

// Create cria uma nova assinatura
func (r *SubscriptionRepositoryPG) Create(ctx context.Context, sub *entity.Subscription) error {
	params := db.CreateSubscriptionParams{
		TenantID:            uuidToPgUUID(sub.TenantID),
		ClienteID:           uuidToPgUUID(sub.ClienteID),
		PlanoID:             uuidToPgUUID(sub.PlanoID),
		AsaasCustomerID:     sub.AsaasCustomerID,
		AsaasSubscriptionID: sub.AsaasSubscriptionID,
		FormaPagamento:      string(sub.FormaPagamento),
		Status:              string(sub.Status),
		Valor:               sub.Valor,
		LinkPagamento:       sub.LinkPagamento,
		CodigoTransacao:     sub.CodigoTransacao,
		DataAtivacao:        timestampToTimestamptzPtr(sub.DataAtivacao),
		DataVencimento:      timestampToTimestamptzPtr(sub.DataVencimento),
	}

	row, err := r.queries.CreateSubscription(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar assinatura: %w", err)
	}

	r.applyRowToEntity(row, sub)
	return nil
}

// GetByID busca assinatura por ID
func (r *SubscriptionRepositoryPG) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*entity.Subscription, error) {
	row, err := r.queries.GetSubscriptionByID(ctx, db.GetSubscriptionByIDParams{
		ID:       uuidToPgUUID(id),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("erro ao buscar assinatura: %w", err)
	}
	return r.mapGetByIDRow(&row), nil
}

// Update atualiza uma assinatura completa
func (r *SubscriptionRepositoryPG) Update(ctx context.Context, sub *entity.Subscription) error {
	return r.queries.UpdateSubscription(ctx, db.UpdateSubscriptionParams{
		ID:                  uuidToPgUUID(sub.ID),
		TenantID:            uuidToPgUUID(sub.TenantID),
		PlanoID:             uuidToPgUUID(sub.PlanoID),
		FormaPagamento:      string(sub.FormaPagamento),
		Status:              string(sub.Status),
		Valor:               sub.Valor,
		DataAtivacao:        timePtrToPgTimestamptz(sub.DataAtivacao),
		DataVencimento:      timePtrToPgTimestamptz(sub.DataVencimento),
		AsaasCustomerID:     sub.AsaasCustomerID,
		AsaasSubscriptionID: sub.AsaasSubscriptionID,
		LinkPagamento:       sub.LinkPagamento,
		CodigoTransacao:     sub.CodigoTransacao,
		ServicosUtilizados:  int32(sub.ServicosUtilizados),
	})
}

// ListByTenant lista assinaturas do tenant
func (r *SubscriptionRepositoryPG) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.Subscription, error) {
	rows, err := r.queries.ListSubscriptionsByTenant(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar assinaturas: %w", err)
	}
	return r.mapListRows(rows), nil
}

// ListByStatus lista assinaturas por status
func (r *SubscriptionRepositoryPG) ListByStatus(ctx context.Context, tenantID uuid.UUID, status entity.SubscriptionStatus) ([]*entity.Subscription, error) {
	rows, err := r.queries.ListSubscriptionsByStatus(ctx, db.ListSubscriptionsByStatusParams{
		TenantID: uuidToPgUUID(tenantID),
		Status:   string(status),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar assinaturas por status: %w", err)
	}
	return r.mapStatusRows(rows), nil
}

// UpdateStatus atualiza apenas o status
func (r *SubscriptionRepositoryPG) UpdateStatus(ctx context.Context, id, tenantID uuid.UUID, status entity.SubscriptionStatus) error {
	if err := r.queries.UpdateSubscriptionStatus(ctx, db.UpdateSubscriptionStatusParams{
		ID:       uuidToPgUUID(id),
		TenantID: uuidToPgUUID(tenantID),
		Status:   string(status),
	}); err != nil {
		return fmt.Errorf("erro ao atualizar status da assinatura: %w", err)
	}
	return nil
}

// Activate marca assinatura como ativa e define datas
func (r *SubscriptionRepositoryPG) Activate(ctx context.Context, id, tenantID uuid.UUID, dataAtivacao, dataVencimento time.Time) error {
	return r.queries.ActivateSubscription(ctx, db.ActivateSubscriptionParams{
		ID:             uuidToPgUUID(id),
		TenantID:       uuidToPgUUID(tenantID),
		DataAtivacao:   timestampToTimestamptz(dataAtivacao),
		DataVencimento: timestampToTimestamptz(dataVencimento),
	})
}

// Cancel cancela assinatura
func (r *SubscriptionRepositoryPG) Cancel(ctx context.Context, id, tenantID uuid.UUID, canceladoPor uuid.UUID) error {
	return r.queries.CancelSubscription(ctx, db.CancelSubscriptionParams{
		ID:           uuidToPgUUID(id),
		TenantID:     uuidToPgUUID(tenantID),
		CanceladoPor: uuidToPgUUID(canceladoPor),
	})
}

// UpdateAsaasFields atualiza campos Asaas após webhook
func (r *SubscriptionRepositoryPG) UpdateAsaasFields(ctx context.Context, id, tenantID uuid.UUID, nextDueDate *time.Time, asaasStatus *string, lastConfirmedAt *time.Time) error {
	return r.queries.UpdateSubscriptionAsaasFields(ctx, db.UpdateSubscriptionAsaasFieldsParams{
		ID:              uuidToPgUUID(id),
		TenantID:        uuidToPgUUID(tenantID),
		NextDueDate:     timePtrToDate(nextDueDate),
		AsaasStatus:     asaasStatus,
		LastConfirmedAt: timestampToTimestamptzPtr(lastConfirmedAt),
	})
}

// UpdateStatusWithAsaas atualiza status interno e status Asaas juntos
func (r *SubscriptionRepositoryPG) UpdateStatusWithAsaas(ctx context.Context, id, tenantID uuid.UUID, status entity.SubscriptionStatus, asaasStatus *string) error {
	return r.queries.UpdateSubscriptionStatusWithAsaas(ctx, db.UpdateSubscriptionStatusWithAsaasParams{
		ID:          uuidToPgUUID(id),
		TenantID:    uuidToPgUUID(tenantID),
		Status:      string(status),
		AsaasStatus: asaasStatus,
	})
}

// CheckActiveExists verifica se cliente já tem assinatura ativa do plano
func (r *SubscriptionRepositoryPG) CheckActiveExists(ctx context.Context, clienteID, planoID uuid.UUID) (bool, error) {
	exists, err := r.queries.CheckActiveSubscriptionExists(ctx, db.CheckActiveSubscriptionExistsParams{
		ClienteID: uuidToPgUUID(clienteID),
		PlanoID:   uuidToPgUUID(planoID),
	})
	if err != nil {
		return false, fmt.Errorf("erro ao verificar assinatura ativa existente: %w", err)
	}
	return exists, nil
}

// UpdateAsaasIDs atualiza IDs externos
func (r *SubscriptionRepositoryPG) UpdateAsaasIDs(ctx context.Context, id, tenantID uuid.UUID, customerID, subscriptionID, linkPagamento *string) error {
	if err := r.queries.UpdateSubscriptionAsaasIDs(ctx, db.UpdateSubscriptionAsaasIDsParams{
		ID:                  uuidToPgUUID(id),
		TenantID:            uuidToPgUUID(tenantID),
		AsaasCustomerID:     customerID,
		AsaasSubscriptionID: subscriptionID,
		LinkPagamento:       linkPagamento,
	}); err != nil {
		return fmt.Errorf("erro ao atualizar IDs do asaas: %w", err)
	}
	return nil
}

// GetByAsaasSubscriptionID busca por ID externo
func (r *SubscriptionRepositoryPG) GetByAsaasSubscriptionID(ctx context.Context, asaasSubscriptionID *string) (*entity.Subscription, error) {
	row, err := r.queries.GetSubscriptionByAsaasID(ctx, asaasSubscriptionID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("erro ao buscar assinatura por asaas_subscription_id: %w", err)
	}
	return r.mapSubscriptionRow(row), nil
}

// IncrementServicosUtilizados incrementa contador
func (r *SubscriptionRepositoryPG) IncrementServicosUtilizados(ctx context.Context, id, tenantID uuid.UUID) error {
	return r.queries.IncrementServicosUtilizados(ctx, db.IncrementServicosUtilizadosParams{
		ID:       uuidToPgUUID(id),
		TenantID: uuidToPgUUID(tenantID),
	})
}

// ResetServicosUtilizados zera contador
func (r *SubscriptionRepositoryPG) ResetServicosUtilizados(ctx context.Context, id, tenantID uuid.UUID) error {
	return r.queries.ResetServicosUtilizados(ctx, db.ResetServicosUtilizadosParams{
		ID:       uuidToPgUUID(id),
		TenantID: uuidToPgUUID(tenantID),
	})
}

// ListOverdue retorna assinaturas vencidas (todas tenants)
func (r *SubscriptionRepositoryPG) ListOverdue(ctx context.Context) ([]*entity.Subscription, error) {
	rows, err := r.queries.ListOverdueSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar assinaturas vencidas: %w", err)
	}
	subs := make([]*entity.Subscription, 0, len(rows))
	for i := range rows {
		row := rows[i]
		subs = append(subs, r.mapOverdueRow(&row))
	}
	return subs, nil
}

// ListExpiringSoon busca assinaturas que vencem em N dias
func (r *SubscriptionRepositoryPG) ListExpiringSoon(ctx context.Context, tenantID uuid.UUID, days int) ([]*entity.Subscription, error) {
	daysStr := strconv.Itoa(days)
	rows, err := r.queries.ListExpiringSoon(ctx, db.ListExpiringSoonParams{
		TenantID: uuidToPgUUID(tenantID),
		Column2:  &daysStr,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar assinaturas a vencer: %w", err)
	}

	subs := make([]*entity.Subscription, 0, len(rows))
	for i := range rows {
		rows[i].TenantID = uuidToPgUUID(tenantID) // já vem no row, mas garantir
		subs = append(subs, r.mapListExpiringRow(&rows[i]))
	}
	return subs, nil
}

// GetMetrics retorna métricas agregadas
func (r *SubscriptionRepositoryPG) GetMetrics(ctx context.Context, tenantID uuid.UUID) (*entity.SubscriptionMetrics, error) {
	row, err := r.queries.GetSubscriptionMetrics(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao obter métricas de assinatura: %w", err)
	}
	return &entity.SubscriptionMetrics{
		TotalAtivas:             int(row.TotalAtivas),
		TotalInativas:           int(row.TotalInativas),
		TotalInadimplentes:      int(row.TotalInadimplentes),
		TotalPlanosAtivos:       int(row.TotalPlanosAtivos),
		ReceitaMensal:           row.ReceitaMensal,
		TaxaRenovacao:           row.TaxaRenovacao,
		RenovacoesProximos7Dias: int(row.RenovacoesProximos7Dias),
	}, nil
}

// CreatePayment registra histórico de pagamento
func (r *SubscriptionRepositoryPG) CreatePayment(ctx context.Context, payment *entity.SubscriptionPayment) error {
	params := db.CreateSubscriptionPaymentParams{
		TenantID:        uuidToPgUUID(payment.TenantID),
		SubscriptionID:  uuidToPgUUID(payment.SubscriptionID),
		AsaasPaymentID:  payment.AsaasPaymentID,
		Valor:           payment.Valor,
		FormaPagamento:  string(payment.FormaPagamento),
		Status:          string(payment.Status),
		DataPagamento:   timestampToTimestamptzPtr(payment.DataPagamento),
		CodigoTransacao: payment.CodigoTransacao,
		Observacao:      payment.Observacao,
	}
	row, err := r.queries.CreateSubscriptionPayment(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar pagamento de assinatura: %w", err)
	}
	payment.ID = pgUUIDToUUID(row.ID)
	payment.CreatedAt = row.CreatedAt.Time
	return nil
}

// ListPaymentsBySubscription lista histórico
func (r *SubscriptionRepositoryPG) ListPaymentsBySubscription(ctx context.Context, subscriptionID, tenantID uuid.UUID) ([]*entity.SubscriptionPayment, error) {
	rows, err := r.queries.ListPaymentsBySubscription(ctx, db.ListPaymentsBySubscriptionParams{
		SubscriptionID: uuidToPgUUID(subscriptionID),
		TenantID:       uuidToPgUUID(tenantID),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pagamentos: %w", err)
	}
	result := make([]*entity.SubscriptionPayment, 0, len(rows))
	for i := range rows {
		result = append(result, r.mapPaymentRow(&rows[i]))
	}
	return result, nil
}

// UpdatePaymentStatus atualiza status (e opcionalmente data)
func (r *SubscriptionRepositoryPG) UpdatePaymentStatus(ctx context.Context, paymentID, tenantID uuid.UUID, status entity.PaymentStatus, dataPagamento *time.Time) error {
	if status == entity.PaymentStatusConfirmado && dataPagamento != nil {
		return r.queries.ConfirmPayment(ctx, db.ConfirmPaymentParams{
			ID:            uuidToPgUUID(paymentID),
			TenantID:      uuidToPgUUID(tenantID),
			DataPagamento: timestampToTimestamptz(*dataPagamento),
		})
	}

	return r.queries.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
		ID:       uuidToPgUUID(paymentID),
		TenantID: uuidToPgUUID(tenantID),
		Status:   string(status),
	})
}

// GetPaymentByAsaasID busca pagamento por ID externo
func (r *SubscriptionRepositoryPG) GetPaymentByAsaasID(ctx context.Context, asaasPaymentID *string) (*entity.SubscriptionPayment, error) {
	row, err := r.queries.GetPaymentByAsaasID(ctx, asaasPaymentID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar pagamento por asaas_payment_id: %w", err)
	}
	return r.mapPaymentRow(&row), nil
}

// UpdateClienteAsaasID atualiza ID externo no cliente
func (r *SubscriptionRepositoryPG) UpdateClienteAsaasID(ctx context.Context, clienteID, tenantID uuid.UUID, asaasCustomerID *string) error {
	return r.queries.UpdateClienteAsaasID(ctx, db.UpdateClienteAsaasIDParams{
		ID:              uuidToPgUUID(clienteID),
		TenantID:        uuidToPgUUID(tenantID),
		AsaasCustomerID: asaasCustomerID,
	})
}

// SetClienteAsSubscriber marca/desmarca flag de assinante
func (r *SubscriptionRepositoryPG) SetClienteAsSubscriber(ctx context.Context, clienteID, tenantID uuid.UUID, isSubscriber bool) error {
	return r.queries.SetClienteAsSubscriber(ctx, db.SetClienteAsSubscriberParams{
		ID:           uuidToPgUUID(clienteID),
		TenantID:     uuidToPgUUID(tenantID),
		IsSubscriber: isSubscriber,
	})
}

// CountActiveSubscriptionsByCliente conta assinaturas ativas do cliente
func (r *SubscriptionRepositoryPG) CountActiveSubscriptionsByCliente(ctx context.Context, clienteID, tenantID uuid.UUID) (int, error) {
	count, err := r.queries.CountActiveSubscriptionsByCliente(ctx, db.CountActiveSubscriptionsByClienteParams{
		ClienteID: uuidToPgUUID(clienteID),
		TenantID:  uuidToPgUUID(tenantID),
	})
	if err != nil {
		return 0, fmt.Errorf("erro ao contar assinaturas ativas do cliente: %w", err)
	}
	return int(count), nil
}

// ============================
// Mapeamentos auxiliares
// ============================

func (r *SubscriptionRepositoryPG) applyRowToEntity(row db.Subscription, sub *entity.Subscription) {
	sub.ID = pgUUIDToUUID(row.ID)
	sub.TenantID = pgUUIDToUUID(row.TenantID)
	sub.ClienteID = pgUUIDToUUID(row.ClienteID)
	sub.PlanoID = pgUUIDToUUID(row.PlanoID)
	sub.AsaasCustomerID = row.AsaasCustomerID
	sub.AsaasSubscriptionID = row.AsaasSubscriptionID
	sub.FormaPagamento = entity.PaymentMethod(row.FormaPagamento)
	sub.Status = entity.SubscriptionStatus(row.Status)
	sub.Valor = row.Valor
	sub.LinkPagamento = row.LinkPagamento
	sub.CodigoTransacao = row.CodigoTransacao
	sub.DataAtivacao = timestamptzToTimePtr(row.DataAtivacao)
	sub.DataVencimento = timestamptzToTimePtr(row.DataVencimento)
	sub.DataCancelamento = timestamptzToTimePtr(row.DataCancelamento)
	sub.CanceladoPor = pgUUIDToUUIDPtr(row.CanceladoPor)
	sub.ServicosUtilizados = int(row.ServicosUtilizados)
	sub.CreatedAt = row.CreatedAt.Time
	sub.UpdatedAt = row.UpdatedAt.Time
}

func (r *SubscriptionRepositoryPG) mapListRows(rows []db.ListSubscriptionsByTenantRow) []*entity.Subscription {
	subs := make([]*entity.Subscription, 0, len(rows))
	for i := range rows {
		row := rows[i]
		subs = append(subs, &entity.Subscription{
			ID:                  pgUUIDToUUID(row.ID),
			TenantID:            pgUUIDToUUID(row.TenantID),
			ClienteID:           pgUUIDToUUID(row.ClienteID),
			PlanoID:             pgUUIDToUUID(row.PlanoID),
			AsaasCustomerID:     row.AsaasCustomerID,
			AsaasSubscriptionID: row.AsaasSubscriptionID,
			FormaPagamento:      entity.PaymentMethod(row.FormaPagamento),
			Status:              entity.SubscriptionStatus(row.Status),
			Valor:               row.Valor,
			LinkPagamento:       row.LinkPagamento,
			CodigoTransacao:     row.CodigoTransacao,
			DataAtivacao:        timestamptzToTimePtr(row.DataAtivacao),
			DataVencimento:      timestamptzToTimePtr(row.DataVencimento),
			DataCancelamento:    timestamptzToTimePtr(row.DataCancelamento),
			CanceladoPor:        pgUUIDToUUIDPtr(row.CanceladoPor),
			ServicosUtilizados:  int(row.ServicosUtilizados),
			CreatedAt:           row.CreatedAt.Time,
			UpdatedAt:           row.UpdatedAt.Time,
			PlanoNome:           row.PlanoNome,
			ClienteNome:         row.ClienteNome,
			ClienteTelefone:     row.ClienteTelefone,
		})
	}
	return subs
}

func (r *SubscriptionRepositoryPG) mapStatusRows(rows []db.ListSubscriptionsByStatusRow) []*entity.Subscription {
	subs := make([]*entity.Subscription, 0, len(rows))
	for i := range rows {
		row := rows[i]
		subs = append(subs, &entity.Subscription{
			ID:                  pgUUIDToUUID(row.ID),
			TenantID:            pgUUIDToUUID(row.TenantID),
			ClienteID:           pgUUIDToUUID(row.ClienteID),
			PlanoID:             pgUUIDToUUID(row.PlanoID),
			AsaasCustomerID:     row.AsaasCustomerID,
			AsaasSubscriptionID: row.AsaasSubscriptionID,
			FormaPagamento:      entity.PaymentMethod(row.FormaPagamento),
			Status:              entity.SubscriptionStatus(row.Status),
			Valor:               row.Valor,
			LinkPagamento:       row.LinkPagamento,
			CodigoTransacao:     row.CodigoTransacao,
			DataAtivacao:        timestamptzToTimePtr(row.DataAtivacao),
			DataVencimento:      timestamptzToTimePtr(row.DataVencimento),
			DataCancelamento:    timestamptzToTimePtr(row.DataCancelamento),
			CanceladoPor:        pgUUIDToUUIDPtr(row.CanceladoPor),
			ServicosUtilizados:  int(row.ServicosUtilizados),
			CreatedAt:           row.CreatedAt.Time,
			UpdatedAt:           row.UpdatedAt.Time,
			PlanoNome:           row.PlanoNome,
			ClienteNome:         row.ClienteNome,
			ClienteTelefone:     row.ClienteTelefone,
		})
	}
	return subs
}

func (r *SubscriptionRepositoryPG) mapOverdueRow(row *db.ListOverdueSubscriptionsRow) *entity.Subscription {
	return &entity.Subscription{
		ID:                  pgUUIDToUUID(row.ID),
		TenantID:            pgUUIDToUUID(row.TenantID),
		ClienteID:           pgUUIDToUUID(row.ClienteID),
		PlanoID:             pgUUIDToUUID(row.PlanoID),
		AsaasCustomerID:     row.AsaasCustomerID,
		AsaasSubscriptionID: row.AsaasSubscriptionID,
		FormaPagamento:      entity.PaymentMethod(row.FormaPagamento),
		Status:              entity.SubscriptionStatus(row.Status),
		Valor:               row.Valor,
		LinkPagamento:       row.LinkPagamento,
		CodigoTransacao:     row.CodigoTransacao,
		DataAtivacao:        timestamptzToTimePtr(row.DataAtivacao),
		DataVencimento:      timestamptzToTimePtr(row.DataVencimento),
		DataCancelamento:    timestamptzToTimePtr(row.DataCancelamento),
		CanceladoPor:        pgUUIDToUUIDPtr(row.CanceladoPor),
		ServicosUtilizados:  int(row.ServicosUtilizados),
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
		PlanoNome:           row.PlanoNome,
		ClienteNome:         row.ClienteNome,
	}
}

func (r *SubscriptionRepositoryPG) mapListExpiringRow(row *db.ListExpiringSoonRow) *entity.Subscription {
	return &entity.Subscription{
		ID:                  pgUUIDToUUID(row.ID),
		TenantID:            pgUUIDToUUID(row.TenantID),
		ClienteID:           pgUUIDToUUID(row.ClienteID),
		PlanoID:             pgUUIDToUUID(row.PlanoID),
		AsaasCustomerID:     row.AsaasCustomerID,
		AsaasSubscriptionID: row.AsaasSubscriptionID,
		FormaPagamento:      entity.PaymentMethod(row.FormaPagamento),
		Status:              entity.SubscriptionStatus(row.Status),
		Valor:               row.Valor,
		LinkPagamento:       row.LinkPagamento,
		CodigoTransacao:     row.CodigoTransacao,
		DataAtivacao:        timestamptzToTimePtr(row.DataAtivacao),
		DataVencimento:      timestamptzToTimePtr(row.DataVencimento),
		DataCancelamento:    timestamptzToTimePtr(row.DataCancelamento),
		CanceladoPor:        pgUUIDToUUIDPtr(row.CanceladoPor),
		ServicosUtilizados:  int(row.ServicosUtilizados),
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
		PlanoNome:           row.PlanoNome,
		ClienteNome:         row.ClienteNome,
		ClienteTelefone:     row.ClienteTelefone,
	}
}

func (r *SubscriptionRepositoryPG) mapSubscriptionRow(row db.GetSubscriptionByAsaasIDRow) *entity.Subscription {
	return &entity.Subscription{
		ID:                  pgUUIDToUUID(row.ID),
		TenantID:            pgUUIDToUUID(row.TenantID),
		ClienteID:           pgUUIDToUUID(row.ClienteID),
		PlanoID:             pgUUIDToUUID(row.PlanoID),
		AsaasCustomerID:     row.AsaasCustomerID,
		AsaasSubscriptionID: row.AsaasSubscriptionID,
		FormaPagamento:      entity.PaymentMethod(row.FormaPagamento),
		Status:              entity.SubscriptionStatus(row.Status),
		Valor:               row.Valor,
		LinkPagamento:       row.LinkPagamento,
		CodigoTransacao:     row.CodigoTransacao,
		DataAtivacao:        timestamptzToTimePtr(row.DataAtivacao),
		DataVencimento:      timestamptzToTimePtr(row.DataVencimento),
		DataCancelamento:    timestamptzToTimePtr(row.DataCancelamento),
		CanceladoPor:        pgUUIDToUUIDPtr(row.CanceladoPor),
		ServicosUtilizados:  int(row.ServicosUtilizados),
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
		PlanoNome:           row.PlanoNome,
	}
}

func (r *SubscriptionRepositoryPG) mapPaymentRow(row *db.SubscriptionPayment) *entity.SubscriptionPayment {
	return &entity.SubscriptionPayment{
		ID:              pgUUIDToUUID(row.ID),
		TenantID:        pgUUIDToUUID(row.TenantID),
		SubscriptionID:  pgUUIDToUUID(row.SubscriptionID),
		AsaasPaymentID:  row.AsaasPaymentID,
		Valor:           row.Valor,
		FormaPagamento:  entity.PaymentMethod(row.FormaPagamento),
		Status:          entity.PaymentStatus(row.Status),
		DataPagamento:   timestamptzToTimePtr(row.DataPagamento),
		CodigoTransacao: row.CodigoTransacao,
		Observacao:      row.Observacao,
		CreatedAt:       row.CreatedAt.Time,
	}
}

func (r *SubscriptionRepositoryPG) mapGetByIDRow(row *db.GetSubscriptionByIDRow) *entity.Subscription {
	return &entity.Subscription{
		ID:                  pgUUIDToUUID(row.ID),
		TenantID:            pgUUIDToUUID(row.TenantID),
		ClienteID:           pgUUIDToUUID(row.ClienteID),
		PlanoID:             pgUUIDToUUID(row.PlanoID),
		AsaasCustomerID:     row.AsaasCustomerID,
		AsaasSubscriptionID: row.AsaasSubscriptionID,
		FormaPagamento:      entity.PaymentMethod(row.FormaPagamento),
		Status:              entity.SubscriptionStatus(row.Status),
		Valor:               row.Valor,
		LinkPagamento:       row.LinkPagamento,
		CodigoTransacao:     row.CodigoTransacao,
		DataAtivacao:        timestamptzToTimePtr(row.DataAtivacao),
		DataVencimento:      timestamptzToTimePtr(row.DataVencimento),
		DataCancelamento:    timestamptzToTimePtr(row.DataCancelamento),
		CanceladoPor:        pgUUIDToUUIDPtr(row.CanceladoPor),
		ServicosUtilizados:  int(row.ServicosUtilizados),
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
		PlanoNome:           row.PlanoNome,
		PlanoQtdServicos:    int32PtrToIntPtr(row.PlanoQtdServicos),
		PlanoLimiteUso:      int32PtrToIntPtr(row.PlanoLimiteUsoMensal),
		ClienteNome:         row.ClienteNome,
		ClienteTelefone:     row.ClienteTelefone,
		ClienteEmail:        row.ClienteEmail,
	}
}

// timestampToTimestamptzPtr converte *time.Time para pgtype.Timestamptz
func timestampToTimestamptzPtr(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return timestampToTimestamptz(*t)
}
