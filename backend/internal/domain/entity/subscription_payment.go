package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PaymentStatus enumera os status do histórico de pagamentos
type PaymentStatus string

const (
	PaymentStatusPendente   PaymentStatus = "PENDENTE"
	PaymentStatusConfirmado PaymentStatus = "CONFIRMADO"
	PaymentStatusRecebido   PaymentStatus = "RECEBIDO"
	PaymentStatusVencido    PaymentStatus = "VENCIDO"
	PaymentStatusEstornado  PaymentStatus = "ESTORNADO"
	PaymentStatusCancelado  PaymentStatus = "CANCELADO"
)

// SubscriptionPayment registra um pagamento relacionado a uma assinatura
type SubscriptionPayment struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	SubscriptionID  uuid.UUID
	AsaasPaymentID  *string
	Valor           decimal.Decimal
	FormaPagamento  PaymentMethod
	Status          PaymentStatus
	DataPagamento   *time.Time
	CodigoTransacao *string
	Observacao      *string
	CreatedAt       time.Time

	// Campos de integração Asaas (Migration 041)
	StatusAsaas         *string         // Status original Asaas: PENDING, CONFIRMED, RECEIVED, etc
	DueDate             *time.Time      // Data vencimento original
	ConfirmedDate       *time.Time      // Data confirmação (competência DRE)
	ClientPaymentDate   *time.Time      // Data que cliente pagou
	CreditDate          *time.Time      // Data que dinheiro caiu na conta (regime caixa)
	EstimatedCreditDate *time.Time      // Previsão de compensação
	BillingType         *string         // Tipo Asaas: CREDIT_CARD, PIX, BOLETO
	NetValue            decimal.Decimal // Valor líquido após taxas
	InvoiceURL          *string         // Link da fatura
	BankSlipURL         *string         // Link do boleto
	PixQRCode           *string         // Código PIX copia e cola
}
