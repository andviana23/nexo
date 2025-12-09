package subscription

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/infra/gateway/asaas"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// ProcessWebhookUseCase processes incoming webhooks from Asaas
// Reference: FLUXO_ASSINATURA.md — Seção 6.6
type ProcessWebhookUseCase struct {
	subRepo     port.SubscriptionRepository
	paymentRepo port.SubscriptionPaymentRepository
	logger      *zap.Logger
}

// NewProcessWebhookUseCase creates a new process webhook use case
func NewProcessWebhookUseCase(
	subRepo port.SubscriptionRepository,
	paymentRepo port.SubscriptionPaymentRepository,
	logger *zap.Logger,
) *ProcessWebhookUseCase {
	return &ProcessWebhookUseCase{
		subRepo:     subRepo,
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

// Execute processes the webhook event
func (uc *ProcessWebhookUseCase) Execute(ctx context.Context, event asaas.WebhookEvent) error {
	uc.logger.Info("processing webhook event",
		zap.String("event", event.Event),
	)

	switch event.Event {
	case asaas.EventPaymentConfirmed, asaas.EventPaymentReceived:
		return uc.handlePaymentConfirmed(ctx, event)

	case asaas.EventPaymentOverdue:
		return uc.handlePaymentOverdue(ctx, event)

	case asaas.EventSubscriptionDeleted, asaas.EventSubscriptionInactivated:
		return uc.handleSubscriptionCanceled(ctx, event)

	case asaas.EventSubscriptionRenewed:
		return uc.handleSubscriptionRenewed(ctx, event)

	case asaas.EventPaymentRefunded:
		return uc.handlePaymentRefunded(ctx, event)

	default:
		uc.logger.Debug("ignoring webhook event",
			zap.String("event", event.Event),
		)
		return nil
	}
}

// handlePaymentConfirmed processes PAYMENT_CONFIRMED and PAYMENT_RECEIVED events
// Actions:
// - Update subscription status to ATIVO
// - Update data_ativacao and data_vencimento
// - Register payment in subscription_payments
// - Mark client as is_subscriber = true (RN-CLI-003)
func (uc *ProcessWebhookUseCase) handlePaymentConfirmed(ctx context.Context, event asaas.WebhookEvent) error {
	if event.Payment == nil {
		uc.logger.Warn("payment confirmed event without payment data")
		return nil
	}

	// Find subscription by Asaas subscription ID
	sub, err := uc.findSubscriptionByAsaasID(ctx, event.Payment.Subscription)
	if err != nil {
		return err
	}
	if sub == nil {
		uc.logger.Warn("subscription not found for payment",
			zap.String("asaas_subscription_id", event.Payment.Subscription),
			zap.String("payment_id", event.Payment.ID),
		)
		return nil // Don't error - orphan webhook
	}

	// Parse confirmed date
	var activationDate time.Time
	if event.Payment.ConfirmedDate != "" {
		activationDate, _ = asaas.ParseDate(event.Payment.ConfirmedDate)
	}
	if activationDate.IsZero() {
		activationDate = time.Now()
	}

	// Calculate expiration date (30 days)
	expirationDate := activationDate.AddDate(0, 0, 30)

	// Update subscription status
	sub.Status = entity.StatusAtivo
	sub.DataAtivacao = &activationDate
	sub.DataVencimento = &expirationDate
	sub.ServicosUtilizados = 0 // Reset usage counter

	if err := uc.subRepo.Update(ctx, sub); err != nil {
		uc.logger.Error("failed to update subscription after payment confirmed",
			zap.String("subscription_id", sub.ID.String()),
			zap.Error(err),
		)
		return err
	}

	// Register payment
	payment := &entity.SubscriptionPayment{
		ID:              uuid.New(),
		TenantID:        sub.TenantID,
		SubscriptionID:  sub.ID,
		AsaasPaymentID:  &event.Payment.ID,
		Valor:           decimal.NewFromFloat(event.Payment.Value),
		FormaPagamento:  sub.FormaPagamento,
		Status:          entity.PaymentStatusConfirmado,
		DataPagamento:   &activationDate,
		CodigoTransacao: nil,
		CreatedAt:       time.Now(),
	}

	if err := uc.paymentRepo.Create(ctx, payment); err != nil {
		uc.logger.Error("failed to register payment",
			zap.String("subscription_id", sub.ID.String()),
			zap.Error(err),
		)
		// Don't fail - subscription was updated
	}

	// Mark client as subscriber (RN-CLI-003)
	if err := uc.subRepo.SetClienteAsSubscriber(ctx, sub.ClienteID, sub.TenantID, true); err != nil {
		uc.logger.Error("failed to mark client as subscriber",
			zap.String("cliente_id", sub.ClienteID.String()),
			zap.Error(err),
		)
	}

	uc.logger.Info("subscription activated via webhook",
		zap.String("subscription_id", sub.ID.String()),
		zap.String("status", string(sub.Status)),
		zap.Time("activation_date", activationDate),
		zap.Time("expiration_date", expirationDate),
	)

	return nil
}

// handlePaymentOverdue processes PAYMENT_OVERDUE event
// Actions:
// - Update subscription status to INADIMPLENTE
// - Check if client has other active subscriptions (RN-CLI-004)
// - If not, update is_subscriber = false
func (uc *ProcessWebhookUseCase) handlePaymentOverdue(ctx context.Context, event asaas.WebhookEvent) error {
	if event.Payment == nil {
		return nil
	}

	sub, err := uc.findSubscriptionByAsaasID(ctx, event.Payment.Subscription)
	if err != nil {
		return err
	}
	if sub == nil {
		uc.logger.Warn("subscription not found for overdue payment",
			zap.String("asaas_subscription_id", event.Payment.Subscription),
		)
		return nil
	}

	// Update status
	sub.Status = entity.StatusInadimplente

	if err := uc.subRepo.Update(ctx, sub); err != nil {
		uc.logger.Error("failed to update subscription to inadimplente",
			zap.String("subscription_id", sub.ID.String()),
			zap.Error(err),
		)
		return err
	}

	// Check if client has other active subscriptions (RN-CLI-004)
	count, err := uc.subRepo.CountActiveSubscriptionsByCliente(ctx, sub.ClienteID, sub.TenantID)
	if err == nil && count == 0 {
		if err := uc.subRepo.SetClienteAsSubscriber(ctx, sub.ClienteID, sub.TenantID, false); err != nil {
			uc.logger.Error("failed to remove subscriber flag",
				zap.String("cliente_id", sub.ClienteID.String()),
				zap.Error(err),
			)
		}
	}

	uc.logger.Info("subscription marked as inadimplente via webhook",
		zap.String("subscription_id", sub.ID.String()),
	)

	return nil
}

// handleSubscriptionCanceled processes SUBSCRIPTION_DELETED and SUBSCRIPTION_INACTIVATED events
// Actions:
// - Update subscription status to CANCELADO
// - Check if client has other active subscriptions (RN-CLI-004)
// - If not, update is_subscriber = false
func (uc *ProcessWebhookUseCase) handleSubscriptionCanceled(ctx context.Context, event asaas.WebhookEvent) error {
	if event.Payment == nil {
		return nil
	}

	sub, err := uc.findSubscriptionByAsaasID(ctx, event.Payment.Subscription)
	if err != nil {
		return err
	}
	if sub == nil {
		uc.logger.Warn("subscription not found for cancel event",
			zap.String("asaas_subscription_id", event.Payment.Subscription),
		)
		return nil
	}

	// Update status
	now := time.Now()
	sub.Status = entity.StatusCancelado
	sub.DataCancelamento = &now

	if err := uc.subRepo.Update(ctx, sub); err != nil {
		uc.logger.Error("failed to update subscription to cancelado",
			zap.String("subscription_id", sub.ID.String()),
			zap.Error(err),
		)
		return err
	}

	// Check if client has other active subscriptions (RN-CLI-004)
	count, err := uc.subRepo.CountActiveSubscriptionsByCliente(ctx, sub.ClienteID, sub.TenantID)
	if err == nil && count == 0 {
		if err := uc.subRepo.SetClienteAsSubscriber(ctx, sub.ClienteID, sub.TenantID, false); err != nil {
			uc.logger.Error("failed to remove subscriber flag",
				zap.String("cliente_id", sub.ClienteID.String()),
				zap.Error(err),
			)
		}
	}

	uc.logger.Info("subscription canceled via webhook",
		zap.String("subscription_id", sub.ID.String()),
	)

	return nil
}

// handleSubscriptionRenewed processes SUBSCRIPTION_RENEWED event
// Actions:
// - Update data_vencimento to new due date
// - Reset servicos_utilizados (RN-BEN-004)
func (uc *ProcessWebhookUseCase) handleSubscriptionRenewed(ctx context.Context, event asaas.WebhookEvent) error {
	if event.Payment == nil {
		return nil
	}

	sub, err := uc.findSubscriptionByAsaasID(ctx, event.Payment.Subscription)
	if err != nil {
		return err
	}
	if sub == nil {
		uc.logger.Warn("subscription not found for renewal event",
			zap.String("asaas_subscription_id", event.Payment.Subscription),
		)
		return nil
	}

	// Parse new due date
	var newDueDate time.Time
	if event.Payment.DueDate != "" {
		newDueDate, _ = asaas.ParseDate(event.Payment.DueDate)
	}
	if newDueDate.IsZero() {
		// If no due date, add 30 days from current expiration
		if sub.DataVencimento != nil {
			newDueDate = sub.DataVencimento.AddDate(0, 0, 30)
		} else {
			newDueDate = time.Now().AddDate(0, 0, 30)
		}
	}

	// Update subscription
	sub.DataVencimento = &newDueDate
	sub.ServicosUtilizados = 0 // Reset usage (RN-BEN-004)
	sub.Status = entity.StatusAtivo

	if err := uc.subRepo.Update(ctx, sub); err != nil {
		uc.logger.Error("failed to update subscription after renewal",
			zap.String("subscription_id", sub.ID.String()),
			zap.Error(err),
		)
		return err
	}

	uc.logger.Info("subscription renewed via webhook",
		zap.String("subscription_id", sub.ID.String()),
		zap.Time("new_due_date", newDueDate),
	)

	return nil
}

// handlePaymentRefunded processes PAYMENT_REFUNDED event
// Actions:
// - Update subscription status to INATIVO
// - Register refund in payments
func (uc *ProcessWebhookUseCase) handlePaymentRefunded(ctx context.Context, event asaas.WebhookEvent) error {
	if event.Payment == nil {
		return nil
	}

	sub, err := uc.findSubscriptionByAsaasID(ctx, event.Payment.Subscription)
	if err != nil {
		return err
	}
	if sub == nil {
		uc.logger.Warn("subscription not found for refund event",
			zap.String("asaas_subscription_id", event.Payment.Subscription),
		)
		return nil
	}

	// Update status
	sub.Status = entity.StatusInativo

	if err := uc.subRepo.Update(ctx, sub); err != nil {
		uc.logger.Error("failed to update subscription after refund",
			zap.String("subscription_id", sub.ID.String()),
			zap.Error(err),
		)
		return err
	}

	// Register refund payment
	now := time.Now()
	refundPayment := &entity.SubscriptionPayment{
		ID:             uuid.New(),
		TenantID:       sub.TenantID,
		SubscriptionID: sub.ID,
		AsaasPaymentID: &event.Payment.ID,
		Valor:          decimal.NewFromFloat(event.Payment.Value),
		FormaPagamento: sub.FormaPagamento,
		Status:         entity.PaymentStatusEstornado,
		DataPagamento:  &now,
		CreatedAt:      time.Now(),
	}

	if err := uc.paymentRepo.Create(ctx, refundPayment); err != nil {
		uc.logger.Error("failed to register refund payment",
			zap.String("subscription_id", sub.ID.String()),
			zap.Error(err),
		)
	}

	// Check subscriber flag
	count, err := uc.subRepo.CountActiveSubscriptionsByCliente(ctx, sub.ClienteID, sub.TenantID)
	if err == nil && count == 0 {
		_ = uc.subRepo.SetClienteAsSubscriber(ctx, sub.ClienteID, sub.TenantID, false)
	}

	uc.logger.Info("subscription refund processed via webhook",
		zap.String("subscription_id", sub.ID.String()),
	)

	return nil
}

// findSubscriptionByAsaasID finds a subscription by its Asaas ID
func (uc *ProcessWebhookUseCase) findSubscriptionByAsaasID(ctx context.Context, asaasSubID string) (*entity.Subscription, error) {
	if asaasSubID == "" {
		return nil, nil
	}
	return uc.subRepo.GetByAsaasSubscriptionID(ctx, &asaasSubID)
}
