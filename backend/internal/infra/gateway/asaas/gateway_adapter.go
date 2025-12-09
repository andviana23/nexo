// Package asaas - Gateway adapter implementing port.AsaasGateway
package asaas

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// GatewayAdapter implements port.AsaasGateway using the Asaas Client
type GatewayAdapter struct {
	client *Client
	logger *zap.Logger
}

// NewGatewayAdapter creates a new gateway adapter
func NewGatewayAdapter(client *Client, logger *zap.Logger) *GatewayAdapter {
	return &GatewayAdapter{
		client: client,
		logger: logger,
	}
}

// Ensure GatewayAdapter implements port.AsaasGateway
var _ port.AsaasGateway = (*GatewayAdapter)(nil)

// FindOrCreateCustomer finds a customer by name+phone or creates if not exists
// Reference: REGRA AS-001, AS-002, AS-003, RN-CLI-002
func (g *GatewayAdapter) FindOrCreateCustomer(ctx context.Context, params port.FindOrCreateCustomerParams) (*port.AsaasCustomerResult, error) {
	req := CustomerRequest{
		Name:                 params.Name,
		Email:                params.Email,
		Phone:                params.Phone,
		MobilePhone:          params.MobilePhone,
		CpfCnpj:              params.CpfCnpj, // Optional per REGRA AS-002
		ExternalReference:    params.ExternalReference,
		NotificationDisabled: false, // Allow Asaas notifications
	}

	customer, wasCreated, err := g.client.FindOrCreateCustomer(ctx, req)
	if err != nil {
		g.logger.Error("failed to find or create customer in Asaas",
			zap.String("name", params.Name),
			zap.Error(err),
		)
		return nil, err
	}

	return &port.AsaasCustomerResult{
		AsaasCustomerID: customer.ID,
		WasCreated:      wasCreated,
	}, nil
}

// CreateSubscription creates a subscription in Asaas
// Reference: REGRA AS-002
func (g *GatewayAdapter) CreateSubscription(ctx context.Context, params port.CreateAsaasSubscriptionParams) (*port.AsaasSubscriptionResult, error) {
	// If no next due date provided, use today + 1 day
	nextDueDate := params.NextDueDate
	if nextDueDate == "" {
		nextDueDate = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}

	sub, paymentLink, err := g.client.CreateSubscriptionWithPaymentLink(
		ctx,
		params.CustomerID,
		params.Value,
		params.Description,
		params.ExternalReference,
		nextDueDate,
	)
	if err != nil {
		g.logger.Error("failed to create subscription in Asaas",
			zap.String("customer_id", params.CustomerID),
			zap.Float64("value", params.Value),
			zap.Error(err),
		)
		return nil, err
	}

	return &port.AsaasSubscriptionResult{
		AsaasSubscriptionID: sub.ID,
		PaymentLink:         paymentLink,
		Status:              sub.Status,
	}, nil
}

// CancelSubscription cancels a subscription in Asaas
// Reference: REGRA AS-006
func (g *GatewayAdapter) CancelSubscription(ctx context.Context, subscriptionID string) error {
	if err := g.client.CancelSubscription(ctx, subscriptionID); err != nil {
		g.logger.Error("failed to cancel subscription in Asaas",
			zap.String("subscription_id", subscriptionID),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// GetPaymentLink gets the payment link for a subscription
// Reference: REGRA AS-004
func (g *GatewayAdapter) GetPaymentLink(ctx context.Context, subscriptionID string) (string, error) {
	link, err := g.client.GeneratePaymentLinkForSubscription(ctx, subscriptionID)
	if err != nil {
		g.logger.Error("failed to get payment link",
			zap.String("subscription_id", subscriptionID),
			zap.Error(err),
		)
		return "", err
	}

	return link, nil
}
