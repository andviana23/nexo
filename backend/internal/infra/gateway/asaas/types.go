// Package asaas - Types and DTOs for Asaas API
// Documentation: https://docs.asaas.com/reference
package asaas

import "time"

// ============================================================================
// CUSTOMER TYPES
// Reference: https://docs.asaas.com/reference/criar-novo-cliente
// ============================================================================

// CustomerRequest represents the request to create a customer in Asaas
// Note: CPF is NOT required per REGRA AS-002, AS-003
type CustomerRequest struct {
	Name                 string `json:"name"`                           // Required
	Email                string `json:"email,omitempty"`                // Optional
	Phone                string `json:"phone,omitempty"`                // Optional (landline)
	MobilePhone          string `json:"mobilePhone,omitempty"`          // Optional (mobile)
	CpfCnpj              string `json:"cpfCnpj,omitempty"`              // Optional per REGRA AS-002
	PostalCode           string `json:"postalCode,omitempty"`           // Optional
	Address              string `json:"address,omitempty"`              // Optional
	AddressNumber        string `json:"addressNumber,omitempty"`        // Optional
	Complement           string `json:"complement,omitempty"`           // Optional
	Province             string `json:"province,omitempty"`             // Optional (bairro)
	ExternalReference    string `json:"externalReference,omitempty"`    // Our internal ID
	NotificationDisabled bool   `json:"notificationDisabled,omitempty"` // Disable Asaas notifications
	GroupName            string `json:"groupName,omitempty"`            // Optional group
}

// CustomerResponse represents the customer returned by Asaas
type CustomerResponse struct {
	ID                    string `json:"id"`
	DateCreated           string `json:"dateCreated"`
	Name                  string `json:"name"`
	Email                 string `json:"email"`
	Phone                 string `json:"phone"`
	MobilePhone           string `json:"mobilePhone"`
	CpfCnpj               string `json:"cpfCnpj"`
	PostalCode            string `json:"postalCode"`
	Address               string `json:"address"`
	AddressNumber         string `json:"addressNumber"`
	Complement            string `json:"complement"`
	Province              string `json:"province"`
	City                  *int64 `json:"city"`
	State                 string `json:"state"`
	Country               string `json:"country"`
	ExternalReference     string `json:"externalReference"`
	NotificationDisabled  bool   `json:"notificationDisabled"`
	AdditionalEmails      string `json:"additionalEmails"`
	MunicipalInscription  string `json:"municipalInscription"`
	StateInscription      string `json:"stateInscription"`
	Observations          string `json:"observations"`
	PersonType            string `json:"personType"` // FISICA or JURIDICA
	Deleted               bool   `json:"deleted"`
	CannotBeDeleted       bool   `json:"cannotBeDeleted"`
	CannotBeDeletedReason string `json:"cannotEditReason"`
}

// CustomerListResponse represents paginated customer list
type CustomerListResponse struct {
	Object     string             `json:"object"`
	HasMore    bool               `json:"hasMore"`
	TotalCount int                `json:"totalCount"`
	Limit      int                `json:"limit"`
	Offset     int                `json:"offset"`
	Data       []CustomerResponse `json:"data"`
}

// ============================================================================
// SUBSCRIPTION TYPES
// Reference: https://docs.asaas.com/reference/criar-nova-assinatura
// ============================================================================

// SubscriptionRequest represents the request to create a subscription
type SubscriptionRequest struct {
	Customer          string  `json:"customer"`                    // Asaas customer ID (required)
	BillingType       string  `json:"billingType"`                 // CREDIT_CARD, BOLETO, PIX, etc
	Value             float64 `json:"value"`                       // Subscription value
	NextDueDate       string  `json:"nextDueDate"`                 // First due date (YYYY-MM-DD)
	Cycle             string  `json:"cycle"`                       // MONTHLY, WEEKLY, etc
	Description       string  `json:"description,omitempty"`       // Subscription description
	ExternalReference string  `json:"externalReference,omitempty"` // Our internal subscription ID
	MaxPayments       *int    `json:"maxPayments,omitempty"`       // Max number of payments
	EndDate           string  `json:"endDate,omitempty"`           // End date (YYYY-MM-DD)

	// Discount (optional)
	Discount *DiscountRequest `json:"discount,omitempty"`

	// Fine and Interest (optional)
	Fine     *FineRequest     `json:"fine,omitempty"`
	Interest *InterestRequest `json:"interest,omitempty"`

	// Credit card specific
	CreditCard           *CreditCardRequest    `json:"creditCard,omitempty"`
	CreditCardHolderInfo *CreditCardHolderInfo `json:"creditCardHolderInfo,omitempty"`
	RemoteIp             string                `json:"remoteIp,omitempty"`
}

// DiscountRequest represents discount configuration
type DiscountRequest struct {
	Value            float64 `json:"value"`
	DueDateLimitDays int     `json:"dueDateLimitDays,omitempty"`
	Type             string  `json:"type"` // FIXED or PERCENTAGE
}

// FineRequest represents fine configuration
type FineRequest struct {
	Value float64 `json:"value"` // Percentage value
	Type  string  `json:"type"`  // FIXED or PERCENTAGE
}

// InterestRequest represents interest configuration
type InterestRequest struct {
	Value float64 `json:"value"` // Monthly percentage
	Type  string  `json:"type"`  // PERCENTAGE
}

// CreditCardRequest represents credit card data (tokenized or raw)
type CreditCardRequest struct {
	HolderName      string `json:"holderName,omitempty"`
	Number          string `json:"number,omitempty"`
	ExpiryMonth     string `json:"expiryMonth,omitempty"`
	ExpiryYear      string `json:"expiryYear,omitempty"`
	Ccv             string `json:"ccv,omitempty"`
	CreditCardToken string `json:"creditCardToken,omitempty"` // If using tokenization
}

// CreditCardHolderInfo represents the cardholder information
type CreditCardHolderInfo struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	CpfCnpj       string `json:"cpfCnpj"`
	PostalCode    string `json:"postalCode"`
	AddressNumber string `json:"addressNumber"`
	Phone         string `json:"phone,omitempty"`
	MobilePhone   string `json:"mobilePhone,omitempty"`
}

// SubscriptionResponse represents the subscription returned by Asaas
type SubscriptionResponse struct {
	ID                string  `json:"id"`
	DateCreated       string  `json:"dateCreated"`
	Customer          string  `json:"customer"`
	PaymentLink       string  `json:"paymentLink,omitempty"`
	BillingType       string  `json:"billingType"`
	Cycle             string  `json:"cycle"`
	Value             float64 `json:"value"`
	NextDueDate       string  `json:"nextDueDate"`
	Description       string  `json:"description"`
	Status            string  `json:"status"` // ACTIVE, INACTIVE, EXPIRED
	ExternalReference string  `json:"externalReference"`
	MaxPayments       *int    `json:"maxPayments"`
	EndDate           string  `json:"endDate"`
	Deleted           bool    `json:"deleted"`
}

// SubscriptionListResponse represents paginated subscription list
type SubscriptionListResponse struct {
	Object     string                 `json:"object"`
	HasMore    bool                   `json:"hasMore"`
	TotalCount int                    `json:"totalCount"`
	Limit      int                    `json:"limit"`
	Offset     int                    `json:"offset"`
	Data       []SubscriptionResponse `json:"data"`
}

// ============================================================================
// PAYMENT TYPES
// Reference: https://docs.asaas.com/reference/recuperar-uma-unica-cobranca
// ============================================================================

// PaymentResponse represents a payment/charge in Asaas
type PaymentResponse struct {
	ID                    string   `json:"id"`
	DateCreated           string   `json:"dateCreated"`
	Customer              string   `json:"customer"`
	PaymentLink           string   `json:"paymentLink,omitempty"`
	DueDate               string   `json:"dueDate"`
	Value                 float64  `json:"value"`
	NetValue              float64  `json:"netValue"`
	BillingType           string   `json:"billingType"`
	CanBePaidAfterDueDate bool     `json:"canBePaidAfterDueDate"`
	PixTransaction        string   `json:"pixTransaction,omitempty"`
	Status                string   `json:"status"` // PENDING, RECEIVED, CONFIRMED, OVERDUE, etc
	Description           string   `json:"description"`
	ExternalReference     string   `json:"externalReference"`
	OriginalValue         *float64 `json:"originalValue,omitempty"`
	InterestValue         *float64 `json:"interestValue,omitempty"`
	OriginalDueDate       string   `json:"originalDueDate,omitempty"`
	ClientPaymentDate     string   `json:"clientPaymentDate,omitempty"`
	PaymentDate           string   `json:"paymentDate,omitempty"`
	InvoiceUrl            string   `json:"invoiceUrl,omitempty"`
	BankSlipUrl           string   `json:"bankSlipUrl,omitempty"`
	TransactionReceiptUrl string   `json:"transactionReceiptUrl,omitempty"`
	InvoiceNumber         string   `json:"invoiceNumber,omitempty"`
	Deleted               bool     `json:"deleted"`
	Anticipated           bool     `json:"anticipated"`
	Subscription          string   `json:"subscription,omitempty"` // Subscription ID if from subscription
	ConfirmedDate         string   `json:"confirmedDate,omitempty"`
	CreditDate            string   `json:"creditDate,omitempty"`
	EstimatedCreditDate   string   `json:"estimatedCreditDate,omitempty"`
}

// PaymentListResponse represents paginated payment list
type PaymentListResponse struct {
	Object     string            `json:"object"`
	HasMore    bool              `json:"hasMore"`
	TotalCount int               `json:"totalCount"`
	Limit      int               `json:"limit"`
	Offset     int               `json:"offset"`
	Data       []PaymentResponse `json:"data"`
}

// ============================================================================
// PAYMENT LINK TYPES
// Reference: https://docs.asaas.com/reference/gerar-link-pagamento
// ============================================================================

// PaymentLinkRequest represents the request to generate a payment link
type PaymentLinkRequest struct {
	Name                string  `json:"name"`                       // Link name/description
	Description         string  `json:"description,omitempty"`      // Detailed description
	EndDate             string  `json:"endDate,omitempty"`          // Expiration date
	Value               float64 `json:"value,omitempty"`            // Fixed value (optional)
	BillingType         string  `json:"billingType"`                // UNDEFINED, BOLETO, CREDIT_CARD, PIX
	ChargeType          string  `json:"chargeType"`                 // DETACHED, RECURRENT, INSTALLMENT
	MaxInstallments     int     `json:"maxInstallments,omitempty"`  // Max installments
	DueDateLimitDays    int     `json:"dueDateLimitDays,omitempty"` // Days until due date
	NotificationEnabled bool    `json:"notificationEnabled,omitempty"`
}

// PaymentLinkResponse represents the payment link returned by Asaas
type PaymentLinkResponse struct {
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	URL                 string  `json:"url"`
	Description         string  `json:"description"`
	EndDate             string  `json:"endDate"`
	Value               float64 `json:"value"`
	BillingType         string  `json:"billingType"`
	ChargeType          string  `json:"chargeType"`
	SubscriptionCycle   string  `json:"subscriptionCycle,omitempty"`
	MaxInstallmentCount int     `json:"maxInstallmentCount"`
	Active              bool    `json:"active"`
	NotificationEnabled bool    `json:"notificationEnabled"`
	ViewCount           int     `json:"viewCount"`
	Deleted             bool    `json:"deleted"`
}

// ============================================================================
// WEBHOOK TYPES
// Reference: https://docs.asaas.com/reference/webhooks
// ============================================================================

// WebhookEvent represents an incoming webhook from Asaas
type WebhookEvent struct {
	Event   string           `json:"event"`
	Payment *PaymentResponse `json:"payment,omitempty"`
}

// Webhook event types
const (
	EventPaymentConfirmed                = "PAYMENT_CONFIRMED"
	EventPaymentReceived                 = "PAYMENT_RECEIVED"
	EventPaymentOverdue                  = "PAYMENT_OVERDUE"
	EventPaymentDeleted                  = "PAYMENT_DELETED"
	EventPaymentRestored                 = "PAYMENT_RESTORED"
	EventPaymentRefunded                 = "PAYMENT_REFUNDED"
	EventPaymentCreated                  = "PAYMENT_CREATED"
	EventPaymentUpdated                  = "PAYMENT_UPDATED"
	EventPaymentDueDateChanged           = "PAYMENT_DUEDATE_WARNING"
	EventPaymentCreditCardCaptureRefused = "PAYMENT_CREDIT_CARD_CAPTURE_REFUSED"

	EventSubscriptionCreated     = "SUBSCRIPTION_CREATED"
	EventSubscriptionUpdated     = "SUBSCRIPTION_UPDATED"
	EventSubscriptionDeleted     = "SUBSCRIPTION_DELETED"
	EventSubscriptionRenewed     = "SUBSCRIPTION_RENEWED"
	EventSubscriptionActivated   = "SUBSCRIPTION_ACTIVATED"
	EventSubscriptionInactivated = "SUBSCRIPTION_INACTIVATED"
)

// ============================================================================
// CONSTANTS
// ============================================================================

// BillingType constants
const (
	BillingTypeBoleto     = "BOLETO"
	BillingTypeCreditCard = "CREDIT_CARD"
	BillingTypeDebitCard  = "DEBIT_CARD"
	BillingTypePix        = "PIX"
	BillingTypeUndefined  = "UNDEFINED"
)

// Cycle constants
const (
	CycleWeekly       = "WEEKLY"
	CycleBiweekly     = "BIWEEKLY"
	CycleMonthly      = "MONTHLY"
	CycleQuarterly    = "QUARTERLY"
	CycleSemiannually = "SEMIANNUALLY"
	CycleYearly       = "YEARLY"
)

// PaymentStatus constants
const (
	PaymentStatusPending              = "PENDING"
	PaymentStatusReceived             = "RECEIVED"
	PaymentStatusConfirmed            = "CONFIRMED"
	PaymentStatusOverdue              = "OVERDUE"
	PaymentStatusRefunded             = "REFUNDED"
	PaymentStatusReceivedInCash       = "RECEIVED_IN_CASH"
	PaymentStatusRefundRequested      = "REFUND_REQUESTED"
	PaymentStatusChargeback           = "CHARGEBACK_REQUESTED"
	PaymentStatusChargebackDispute    = "CHARGEBACK_DISPUTE"
	PaymentStatusAwaitingChargeback   = "AWAITING_CHARGEBACK_REVERSAL"
	PaymentStatusDunningRequested     = "DUNNING_REQUESTED"
	PaymentStatusDunningReceived      = "DUNNING_RECEIVED"
	PaymentStatusAwaitingRiskAnalysis = "AWAITING_RISK_ANALYSIS"
)

// SubscriptionStatus constants
const (
	SubscriptionStatusActive   = "ACTIVE"
	SubscriptionStatusInactive = "INACTIVE"
	SubscriptionStatusExpired  = "EXPIRED"
)

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// FormatDate formats a time.Time to Asaas date format (YYYY-MM-DD)
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// ParseDate parses an Asaas date string to time.Time
func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// ParseDateTime parses an Asaas datetime string to time.Time
func ParseDateTime(s string) (time.Time, error) {
	// Asaas returns dates in ISO 8601 format
	return time.Parse(time.RFC3339, s)
}
