package port

import "context"

// AsaasGateway define a interface para integração com o Asaas
// Referência: FLUXO_ASSINATURA.md — Seção 4
type AsaasGateway interface {
	// Customer operations
	// Reference: REGRA AS-001, AS-002, AS-003, RN-CLI-002
	FindOrCreateCustomer(ctx context.Context, params FindOrCreateCustomerParams) (*AsaasCustomerResult, error)

	// Subscription operations
	// Reference: REGRA AS-002
	CreateSubscription(ctx context.Context, params CreateAsaasSubscriptionParams) (*AsaasSubscriptionResult, error)

	// Reference: REGRA AS-006
	CancelSubscription(ctx context.Context, subscriptionID string) error

	// Reference: REGRA AS-004
	GetPaymentLink(ctx context.Context, subscriptionID string) (string, error)
}

// FindOrCreateCustomerParams holds the parameters for finding or creating a customer
type FindOrCreateCustomerParams struct {
	Name              string
	Email             string
	Phone             string
	MobilePhone       string
	CpfCnpj           string // Optional per REGRA AS-002
	ExternalReference string // Our internal client ID
}

// AsaasCustomerResult holds the result of finding or creating a customer
type AsaasCustomerResult struct {
	AsaasCustomerID string
	WasCreated      bool // true if customer was created, false if already existed
}

// CreateAsaasSubscriptionParams holds the parameters for creating a subscription
type CreateAsaasSubscriptionParams struct {
	CustomerID        string // Asaas customer ID
	Value             float64
	Description       string
	ExternalReference string // Our internal subscription ID
	NextDueDate       string // YYYY-MM-DD format
}

// AsaasSubscriptionResult holds the result of creating a subscription
type AsaasSubscriptionResult struct {
	AsaasSubscriptionID string
	PaymentLink         string
	Status              string
}
