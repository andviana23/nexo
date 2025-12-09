package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
)

// AsaasWebhookLogRepository define operações de persistência para logs de webhooks Asaas
type AsaasWebhookLogRepository interface {
	// Create insere um novo log de webhook
	Create(ctx context.Context, log *entity.AsaasWebhookLog) error

	// MarkProcessed marca um webhook como processado
	MarkProcessed(ctx context.Context, id string) error

	// MarkFailed marca um webhook como falho com mensagem de erro
	MarkFailed(ctx context.Context, id string, errorMessage string) error

	// GetByPaymentID busca log por asaas_payment_id (para debug)
	GetByPaymentID(ctx context.Context, asaasPaymentID string) (*entity.AsaasWebhookLog, error)

	// ListUnprocessed lista webhooks não processados para retry
	ListUnprocessed(ctx context.Context, tenantID string, limit int) ([]*entity.AsaasWebhookLog, error)
}
