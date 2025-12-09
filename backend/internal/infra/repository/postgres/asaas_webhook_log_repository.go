package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// AsaasWebhookLogRepositoryPG implementa port.AsaasWebhookLogRepository
type AsaasWebhookLogRepositoryPG struct {
	queries *db.Queries
}

// NewAsaasWebhookLogRepository cria uma instância do repositório de webhook logs
func NewAsaasWebhookLogRepository(queries *db.Queries) port.AsaasWebhookLogRepository {
	return &AsaasWebhookLogRepositoryPG{queries: queries}
}

// Create registra um novo webhook log
func (r *AsaasWebhookLogRepositoryPG) Create(ctx context.Context, log *entity.AsaasWebhookLog) error {
	params := db.CreateWebhookLogParams{
		TenantID:            uuidToPgUUID(log.TenantID),
		EventType:           log.EventType,
		AsaasPaymentID:      log.AsaasPaymentID,
		AsaasSubscriptionID: log.AsaasSubscriptionID,
		Payload:             log.Payload,
	}

	row, err := r.queries.CreateWebhookLog(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar webhook log: %w", err)
	}

	log.ID = pgUUIDToUUID(row.ID)
	log.CreatedAt = row.CreatedAt.Time
	return nil
}

// MarkProcessed marca webhook como processado com sucesso
func (r *AsaasWebhookLogRepositoryPG) MarkProcessed(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("ID inválido: %w", err)
	}
	return r.queries.MarkWebhookProcessed(ctx, uuidToPgUUID(parsedID))
}

// MarkFailed marca webhook como falha com mensagem de erro
func (r *AsaasWebhookLogRepositoryPG) MarkFailed(ctx context.Context, id string, errorMsg string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("ID inválido: %w", err)
	}
	return r.queries.MarkWebhookFailed(ctx, db.MarkWebhookFailedParams{
		ID:           uuidToPgUUID(parsedID),
		ErrorMessage: &errorMsg,
	})
}

// GetByPaymentID busca log por payment ID (para debug)
func (r *AsaasWebhookLogRepositoryPG) GetByPaymentID(ctx context.Context, asaasPaymentID string) (*entity.AsaasWebhookLog, error) {
	rows, err := r.queries.ListWebhooksByPaymentID(ctx, &asaasPaymentID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar webhook log: %w", err)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return r.mapWebhookLogRow(&rows[0]), nil
}

// ListUnprocessed lista webhooks não processados para retry
func (r *AsaasWebhookLogRepositoryPG) ListUnprocessed(ctx context.Context, tenantID string, limit int) ([]*entity.AsaasWebhookLog, error) {
	rows, err := r.queries.ListUnprocessedWebhooks(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar webhooks não processados: %w", err)
	}

	var parsedTenantID uuid.UUID
	if tenantID != "" {
		parsedTenantID, _ = uuid.Parse(tenantID)
	}

	result := make([]*entity.AsaasWebhookLog, 0, len(rows))
	for i := range rows {
		// Filtrar por tenant se necessário
		if parsedTenantID != uuid.Nil && pgUUIDToUUID(rows[i].TenantID) != parsedTenantID {
			continue
		}
		result = append(result, r.mapWebhookLogRow(&rows[i]))
	}
	return result, nil
}

// mapWebhookLogRow mapeia row para entity
func (r *AsaasWebhookLogRepositoryPG) mapWebhookLogRow(row *db.AsaasWebhookLog) *entity.AsaasWebhookLog {
	return &entity.AsaasWebhookLog{
		ID:                  pgUUIDToUUID(row.ID),
		TenantID:            pgUUIDToUUID(row.TenantID),
		EventType:           row.EventType,
		AsaasPaymentID:      row.AsaasPaymentID,
		AsaasSubscriptionID: row.AsaasSubscriptionID,
		Payload:             json.RawMessage(row.Payload),
		ProcessedAt:         timestamptzToTimePtr(row.ProcessedAt),
		ErrorMessage:        row.ErrorMessage,
		CreatedAt:           row.CreatedAt.Time,
	}
}
