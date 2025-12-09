// Package asaas - Subscription API methods
// Reference: https://docs.asaas.com/reference/criar-nova-assinatura
// REGRAS: AS-002, AS-004, AS-006
package asaas

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

// ============================================================================
// SUBSCRIPTION METHODS
// ============================================================================

// CreateSubscription creates a new subscription in Asaas
// Reference: REGRA AS-002 - Assinatura com billingType CREDIT_CARD
func (c *Client) CreateSubscription(ctx context.Context, req SubscriptionRequest) (*SubscriptionResponse, error) {
	c.logger.Debug("creating subscription in Asaas",
		zap.String("customer", req.Customer),
		zap.String("billingType", req.BillingType),
		zap.Float64("value", req.Value),
		zap.String("cycle", req.Cycle),
		zap.String("externalReference", req.ExternalReference),
	)

	body, err := c.doRequest(ctx, "POST", "/subscriptions", req)
	if err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	var resp SubscriptionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal subscription response: %w", err)
	}

	c.logger.Info("subscription created in Asaas",
		zap.String("asaas_subscription_id", resp.ID),
		zap.String("status", resp.Status),
		zap.String("paymentLink", resp.PaymentLink),
	)

	return &resp, nil
}

// GetSubscription retrieves a subscription by ID
func (c *Client) GetSubscription(ctx context.Context, subscriptionID string) (*SubscriptionResponse, error) {
	path := fmt.Sprintf("/subscriptions/%s", subscriptionID)

	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}

	var resp SubscriptionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal subscription: %w", err)
	}

	return &resp, nil
}

// CancelSubscription cancels/deletes a subscription in Asaas
// Reference: REGRA AS-006 - DELETE /subscriptions/{id}
func (c *Client) CancelSubscription(ctx context.Context, subscriptionID string) error {
	c.logger.Debug("canceling subscription in Asaas",
		zap.String("subscriptionID", subscriptionID),
	)

	path := fmt.Sprintf("/subscriptions/%s", subscriptionID)

	body, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		// Check if it's already deleted/not found
		if apiErr, ok := err.(*APIError); ok && apiErr.IsNotFound() {
			c.logger.Warn("subscription not found in Asaas (may already be deleted)",
				zap.String("subscriptionID", subscriptionID),
			)
			return nil // Treat as success
		}
		return fmt.Errorf("cancel subscription: %w", err)
	}

	var resp SubscriptionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("unmarshal cancel response: %w", err)
	}

	if !resp.Deleted {
		c.logger.Warn("subscription delete returned but deleted flag is false",
			zap.String("subscriptionID", subscriptionID),
		)
	}

	c.logger.Info("subscription canceled in Asaas",
		zap.String("subscriptionID", subscriptionID),
		zap.Bool("deleted", resp.Deleted),
	)

	return nil
}

// ListSubscriptionPayments lists all payments for a subscription
// Used to find pending payments for generating payment links
func (c *Client) ListSubscriptionPayments(ctx context.Context, subscriptionID string) ([]PaymentResponse, error) {
	path := fmt.Sprintf("/subscriptions/%s/payments", subscriptionID)

	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("list subscription payments: %w", err)
	}

	var resp PaymentListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal payments list: %w", err)
	}

	return resp.Data, nil
}

// GetPayment retrieves a payment by ID
func (c *Client) GetPayment(ctx context.Context, paymentID string) (*PaymentResponse, error) {
	path := fmt.Sprintf("/payments/%s", paymentID)

	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get payment: %w", err)
	}

	var resp PaymentResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal payment: %w", err)
	}

	return &resp, nil
}

// GetPaymentLink generates a payment link for an existing payment
// Reference: REGRA AS-004 - Link de pagamento expira em 24h
func (c *Client) GetPaymentLink(ctx context.Context, subscriptionID string) (string, error) {
	// First, get the pending payment for this subscription
	payments, err := c.ListSubscriptionPayments(ctx, subscriptionID)
	if err != nil {
		return "", fmt.Errorf("list payments for link: %w", err)
	}

	// Find first pending payment
	for _, payment := range payments {
		if payment.Status == PaymentStatusPending {
			// Use invoiceUrl if available, otherwise paymentLink
			if payment.InvoiceUrl != "" {
				return payment.InvoiceUrl, nil
			}
			if payment.PaymentLink != "" {
				return payment.PaymentLink, nil
			}
		}
	}

	return "", fmt.Errorf("no pending payment found for subscription %s", subscriptionID)
}

// GeneratePaymentLinkForSubscription gets or generates a payment link for a subscription
// This is used when creating a new subscription with CREDIT_CARD
func (c *Client) GeneratePaymentLinkForSubscription(ctx context.Context, subscriptionID string) (string, error) {
	// Get subscription first to check if it has a payment link
	sub, err := c.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return "", fmt.Errorf("get subscription for payment link: %w", err)
	}

	// If subscription already has a payment link, return it
	if sub.PaymentLink != "" {
		return sub.PaymentLink, nil
	}

	// Otherwise, try to get from pending payment
	return c.GetPaymentLink(ctx, subscriptionID)
}

// ============================================================================
// SUBSCRIPTION HELPER METHODS
// ============================================================================

// CreateCreditCardSubscription creates a subscription for credit card payments
// This is a convenience method that sets the correct billing type and cycle
func (c *Client) CreateCreditCardSubscription(
	ctx context.Context,
	customerID string,
	value float64,
	description string,
	externalReference string,
	nextDueDate string,
) (*SubscriptionResponse, error) {
	req := SubscriptionRequest{
		Customer:          customerID,
		BillingType:       BillingTypeCreditCard,
		Value:             value,
		NextDueDate:       nextDueDate,
		Cycle:             CycleMonthly,
		Description:       description,
		ExternalReference: externalReference,
	}

	return c.CreateSubscription(ctx, req)
}

// CreateSubscriptionWithPaymentLink creates a subscription and returns the payment link
// This combines CreateSubscription + GeneratePaymentLinkForSubscription
func (c *Client) CreateSubscriptionWithPaymentLink(
	ctx context.Context,
	customerID string,
	value float64,
	description string,
	externalReference string,
	nextDueDate string,
) (*SubscriptionResponse, string, error) {
	// Create subscription
	sub, err := c.CreateCreditCardSubscription(
		ctx,
		customerID,
		value,
		description,
		externalReference,
		nextDueDate,
	)
	if err != nil {
		return nil, "", err
	}

	// Get payment link
	var paymentLink string
	if sub.PaymentLink != "" {
		paymentLink = sub.PaymentLink
	} else {
		paymentLink, err = c.GeneratePaymentLinkForSubscription(ctx, sub.ID)
		if err != nil {
			c.logger.Warn("could not get payment link, subscription created without link",
				zap.String("subscriptionID", sub.ID),
				zap.Error(err),
			)
			// Don't fail - subscription was created, link can be retrieved later
		}
	}

	return sub, paymentLink, nil
}
