package entity

import (
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"
)

// AsaasWebhookLog registra webhooks recebidos do Asaas para auditoria
type AsaasWebhookLog struct {
	ID                  uuid.UUID
	TenantID            uuid.UUID
	EventType           string
	AsaasPaymentID      *string
	AsaasSubscriptionID *string
	Payload             json.RawMessage
	Processed           bool
	ProcessedAt         *time.Time
	ErrorMessage        *string
	ReceivedAt          time.Time
	IPAddress           net.IP
	CreatedAt           time.Time
}

// NewAsaasWebhookLog cria um novo log de webhook
func NewAsaasWebhookLog(
	tenantID uuid.UUID,
	eventType string,
	paymentID *string,
	subscriptionID *string,
	payload json.RawMessage,
	ipAddress net.IP,
) *AsaasWebhookLog {
	now := time.Now()
	return &AsaasWebhookLog{
		ID:                  uuid.New(),
		TenantID:            tenantID,
		EventType:           eventType,
		AsaasPaymentID:      paymentID,
		AsaasSubscriptionID: subscriptionID,
		Payload:             payload,
		Processed:           false,
		ReceivedAt:          now,
		IPAddress:           ipAddress,
		CreatedAt:           now,
	}
}

// MarkProcessed marca o webhook como processado com sucesso
func (w *AsaasWebhookLog) MarkProcessed() {
	now := time.Now()
	w.Processed = true
	w.ProcessedAt = &now
}

// MarkFailed marca o webhook como falho
func (w *AsaasWebhookLog) MarkFailed(err error) {
	now := time.Now()
	w.Processed = false
	w.ProcessedAt = &now
	errMsg := err.Error()
	w.ErrorMessage = &errMsg
}
