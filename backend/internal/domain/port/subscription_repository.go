package port

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
)

// SubscriptionRepository define operações de persistência para assinaturas e pagamentos
type SubscriptionRepository interface {
	// CRUD básico
	Create(ctx context.Context, sub *entity.Subscription) error
	Update(ctx context.Context, sub *entity.Subscription) error
	GetByID(ctx context.Context, id, tenantID uuid.UUID) (*entity.Subscription, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.Subscription, error)
	ListByStatus(ctx context.Context, tenantID uuid.UUID, status entity.SubscriptionStatus) ([]*entity.Subscription, error)

	// Status management
	UpdateStatus(ctx context.Context, id, tenantID uuid.UUID, status entity.SubscriptionStatus) error
	Activate(ctx context.Context, id, tenantID uuid.UUID, dataAtivacao, dataVencimento time.Time) error
	Cancel(ctx context.Context, id, tenantID uuid.UUID, canceladoPor uuid.UUID) error

	// Regras de duplicidade e limites
	CheckActiveExists(ctx context.Context, clienteID, planoID uuid.UUID) (bool, error)

	// Integração Asaas / IDs externos
	UpdateAsaasIDs(ctx context.Context, id, tenantID uuid.UUID, customerID, subscriptionID, linkPagamento *string) error
	GetByAsaasSubscriptionID(ctx context.Context, asaasSubscriptionID *string) (*entity.Subscription, error)

	// Uso de serviços
	IncrementServicosUtilizados(ctx context.Context, id, tenantID uuid.UUID) error
	ResetServicosUtilizados(ctx context.Context, id, tenantID uuid.UUID) error

	// Cron / vencimentos
	ListOverdue(ctx context.Context) ([]*entity.Subscription, error)
	ListExpiringSoon(ctx context.Context, tenantID uuid.UUID, days int) ([]*entity.Subscription, error)

	// Métricas
	GetMetrics(ctx context.Context, tenantID uuid.UUID) (*entity.SubscriptionMetrics, error)

	// Pagamentos
	CreatePayment(ctx context.Context, payment *entity.SubscriptionPayment) error
	ListPaymentsBySubscription(ctx context.Context, subscriptionID, tenantID uuid.UUID) ([]*entity.SubscriptionPayment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID, tenantID uuid.UUID, status entity.PaymentStatus, dataPagamento *time.Time) error
	GetPaymentByAsaasID(ctx context.Context, asaasPaymentID *string) (*entity.SubscriptionPayment, error)

	// Cliente (flags de assinante e IDs externos)
	UpdateClienteAsaasID(ctx context.Context, clienteID, tenantID uuid.UUID, asaasCustomerID *string) error
	SetClienteAsSubscriber(ctx context.Context, clienteID, tenantID uuid.UUID, isSubscriber bool) error
	CountActiveSubscriptionsByCliente(ctx context.Context, clienteID, tenantID uuid.UUID) (int, error)

	// Campos Asaas (Migration 041)
	UpdateAsaasFields(ctx context.Context, id, tenantID uuid.UUID, nextDueDate *time.Time, asaasStatus *string, lastConfirmedAt *time.Time) error
	UpdateStatusWithAsaas(ctx context.Context, id, tenantID uuid.UUID, status entity.SubscriptionStatus, asaasStatus *string) error
}

// SubscriptionPaymentRepository define operações de persistência para pagamentos de assinatura
// Separado para uso no ProcessWebhookUseCase
type SubscriptionPaymentRepository interface {
	Create(ctx context.Context, payment *entity.SubscriptionPayment) error
	ListBySubscription(ctx context.Context, subscriptionID, tenantID uuid.UUID) ([]*entity.SubscriptionPayment, error)
	UpdateStatus(ctx context.Context, paymentID, tenantID uuid.UUID, status entity.PaymentStatus, dataPagamento *time.Time) error
	GetByAsaasID(ctx context.Context, asaasPaymentID string) (*entity.SubscriptionPayment, error)

	// Idempotência para webhooks (Migration 041)
	UpsertByAsaasID(ctx context.Context, payment *entity.SubscriptionPayment) error

	// Métodos de atualização por status específico
	UpdatePaymentConfirmed(ctx context.Context, paymentID, tenantID uuid.UUID, statusAsaas string, confirmedDate time.Time) error
	UpdatePaymentReceived(ctx context.Context, paymentID, tenantID uuid.UUID, statusAsaas string, clientPaymentDate, creditDate time.Time) error
	UpdatePaymentOverdue(ctx context.Context, paymentID, tenantID uuid.UUID, statusAsaas string) error
	UpdatePaymentRefunded(ctx context.Context, paymentID, tenantID uuid.UUID, statusAsaas string) error

	// Conciliação (Sprint 4)
	ListNeedingReconciliation(ctx context.Context, tenantID string) ([]*entity.SubscriptionPayment, error)
}
