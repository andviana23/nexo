package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// SubscriptionStatus representa o status de uma assinatura
type SubscriptionStatus string

const (
	StatusAguardandoPagamento SubscriptionStatus = "AGUARDANDO_PAGAMENTO"
	StatusAtivo               SubscriptionStatus = "ATIVO"
	StatusInadimplente        SubscriptionStatus = "INADIMPLENTE"
	StatusInativo             SubscriptionStatus = "INATIVO"
	StatusCancelado           SubscriptionStatus = "CANCELADO"
)

// PaymentMethod representa a forma de pagamento da assinatura
type PaymentMethod string

const (
	PaymentMethodCartao   PaymentMethod = "CARTAO"
	PaymentMethodPix      PaymentMethod = "PIX"
	PaymentMethodDinheiro PaymentMethod = "DINHEIRO"
)

// Subscription representa uma assinatura de cliente
type Subscription struct {
	ID                  uuid.UUID
	TenantID            uuid.UUID
	ClienteID           uuid.UUID
	PlanoID             uuid.UUID
	AsaasCustomerID     *string
	AsaasSubscriptionID *string
	FormaPagamento      PaymentMethod
	Status              SubscriptionStatus
	Valor               decimal.Decimal
	LinkPagamento       *string
	CodigoTransacao     *string
	DataAtivacao        *time.Time
	DataVencimento      *time.Time
	DataCancelamento    *time.Time
	CanceladoPor        *uuid.UUID
	ServicosUtilizados  int
	CreatedAt           time.Time
	UpdatedAt           time.Time

	// Campos de integração Asaas (Migration 041)
	NextDueDate     *time.Time // Próxima data vencimento (Asaas nextDueDate)
	Cycle           string     // Ciclo: MONTHLY, WEEKLY, etc
	AsaasStatus     *string    // Status no Asaas: ACTIVE, INACTIVE, EXPIRED
	LastConfirmedAt *time.Time // Último pagamento CONFIRMED
	LastSyncAt      *time.Time // Última sincronização com Asaas

	// Campos derivados (JOIN)
	PlanoNome        string
	PlanoQtdServicos *int
	PlanoLimiteUso   *int
	ClienteNome      string
	ClienteTelefone  string
	ClienteEmail     *string
}

// SubscriptionMetrics representa métricas agregadas do módulo
type SubscriptionMetrics struct {
	TotalAtivas             int
	TotalInativas           int
	TotalInadimplentes      int
	TotalPlanosAtivos       int
	ReceitaMensal           decimal.Decimal
	TaxaRenovacao           float64
	RenovacoesProximos7Dias int
}

// NewSubscription cria uma nova assinatura com validação mínima.
func NewSubscription(tenantID, clienteID, planoID uuid.UUID, forma PaymentMethod, valor decimal.Decimal) (*Subscription, error) {
	if valor.LessThanOrEqual(decimal.Zero) {
		return nil, domain.ErrValorInvalido
	}

	if forma != PaymentMethodCartao && forma != PaymentMethodPix && forma != PaymentMethodDinheiro {
		return nil, domain.ErrSubscriptionPaymentMethodInvalid
	}

	now := time.Now()
	return &Subscription{
		ID:             uuid.New(),
		TenantID:       tenantID,
		ClienteID:      clienteID,
		PlanoID:        planoID,
		FormaPagamento: forma,
		Status:         StatusAguardandoPagamento,
		Valor:          valor,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// CanUseServices aplica RN-BEN-001
func (s *Subscription) CanUseServices() bool {
	return s.Status == StatusAtivo
}

// HasReachedServiceLimit aplica RN-BEN-003
func (s *Subscription) HasReachedServiceLimit(planoLimite *int) bool {
	if planoLimite == nil {
		return false
	}
	return s.ServicosUtilizados >= *planoLimite
}

// ShouldBecomeInadimplente aplica RN-VENC-004 (3 dias após vencimento)
func (s *Subscription) ShouldBecomeInadimplente(now time.Time) bool {
	if s.DataVencimento == nil || s.Status != StatusAtivo {
		return false
	}
	limite := s.DataVencimento.AddDate(0, 0, 3)
	return now.After(limite)
}

// CanBeReactivated aplica RN-CANC-004
func (s *Subscription) CanBeReactivated() bool {
	return s.Status != StatusCancelado
}

// RequiresAsaasIntegration indica se deve integrar com Asaas
func (s *Subscription) RequiresAsaasIntegration() bool {
	return s.FormaPagamento == PaymentMethodCartao
}

// ShouldMarkClientAsSubscriber aplica RN-CLI-003
func (s *Subscription) ShouldMarkClientAsSubscriber() bool {
	return s.Status == StatusAtivo
}

// SetCanceled atualiza dados de cancelamento
func (s *Subscription) SetCanceled(by uuid.UUID) {
	s.Status = StatusCancelado
	s.CanceladoPor = &by
	now := time.Now()
	s.DataCancelamento = &now
	s.UpdatedAt = now
}
