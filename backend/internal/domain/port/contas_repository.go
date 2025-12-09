package port

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
)

// ContaPagarRepository define operações para Contas a Pagar
type ContaPagarRepository interface {
	// Create cria uma nova conta a pagar
	Create(ctx context.Context, conta *entity.ContaPagar) error

	// FindByID busca uma conta por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.ContaPagar, error)

	// Update atualiza uma conta existente
	Update(ctx context.Context, conta *entity.ContaPagar) error

	// Delete remove uma conta
	Delete(ctx context.Context, tenantID, id string) error

	// List lista contas com filtros
	List(ctx context.Context, tenantID string, filters ContaPagarListFilters) ([]*entity.ContaPagar, error)

	// ListByStatus lista contas por status
	ListByStatus(ctx context.Context, tenantID string, status valueobject.StatusConta) ([]*entity.ContaPagar, error)

	// ListVencendoEm lista contas que vencem em até N dias
	ListVencendoEm(ctx context.Context, tenantID string, dias int) ([]*entity.ContaPagar, error)

	// ListAtrasadas lista contas atrasadas
	ListAtrasadas(ctx context.Context, tenantID string) ([]*entity.ContaPagar, error)

	// ListByDateRange lista contas em um período (data vencimento)
	ListByDateRange(ctx context.Context, tenantID string, inicio, fim time.Time) ([]*entity.ContaPagar, error)

	// SumByPeriod soma valores de contas em um período
	SumByPeriod(ctx context.Context, tenantID string, inicio, fim time.Time, status *valueobject.StatusConta) (valueobject.Money, error)

	// SumByCategoria soma valores por categoria
	SumByCategoria(ctx context.Context, tenantID, categoriaID string, inicio, fim time.Time) (valueobject.Money, error)
}

// ContaPagarListFilters filtros para listagem de contas a pagar
type ContaPagarListFilters struct {
	Status      *valueobject.StatusConta
	Tipo        *valueobject.TipoCusto
	CategoriaID *string
	Fornecedor  *string
	Recorrente  *bool
	Page        int
	PageSize    int
	OrderBy     string
}

// ContaReceberRepository define operações para Contas a Receber
type ContaReceberRepository interface {
	// Create cria uma nova conta a receber
	Create(ctx context.Context, conta *entity.ContaReceber) error

	// FindByID busca uma conta por ID
	FindByID(ctx context.Context, tenantID, id string) (*entity.ContaReceber, error)

	// Update atualiza uma conta existente
	Update(ctx context.Context, conta *entity.ContaReceber) error

	// Delete remove uma conta
	Delete(ctx context.Context, tenantID, id string) error

	// List lista contas com filtros
	List(ctx context.Context, tenantID string, filters ContaReceberListFilters) ([]*entity.ContaReceber, error)

	// ListByStatus lista contas por status
	ListByStatus(ctx context.Context, tenantID string, status valueobject.StatusConta) ([]*entity.ContaReceber, error)

	// ListByAssinatura lista contas de uma assinatura
	ListByAssinatura(ctx context.Context, tenantID, assinaturaID string) ([]*entity.ContaReceber, error)

	// ListVencendoEm lista contas que vencem em até N dias
	ListVencendoEm(ctx context.Context, tenantID string, dias int) ([]*entity.ContaReceber, error)

	// ListAtrasadas lista contas atrasadas
	ListAtrasadas(ctx context.Context, tenantID string) ([]*entity.ContaReceber, error)

	// ListByDateRange lista contas em um período (data vencimento)
	ListByDateRange(ctx context.Context, tenantID string, inicio, fim time.Time) ([]*entity.ContaReceber, error)

	// SumByPeriod soma valores de contas em um período
	SumByPeriod(ctx context.Context, tenantID string, inicio, fim time.Time, status *valueobject.StatusConta) (valueobject.Money, error)

	// SumByOrigem soma valores por origem
	SumByOrigem(ctx context.Context, tenantID, origem string, inicio, fim time.Time) (valueobject.Money, error)

	// Integração Asaas (Migration 041)
	// UpsertByAsaasPaymentID cria ou atualiza conta a receber pela cobrança Asaas (idempotência)
	UpsertByAsaasPaymentID(ctx context.Context, conta *entity.ContaReceber) error

	// GetByAsaasPaymentID busca conta por ID de pagamento Asaas
	GetByAsaasPaymentID(ctx context.Context, tenantID, asaasPaymentID string) (*entity.ContaReceber, error)

	// SumByCompetencia soma valores por mês de competência (DRE)
	SumByCompetencia(ctx context.Context, tenantID, competenciaMes string, status *valueobject.StatusConta) (valueobject.Money, error)

	// ListBySubscriptionID lista contas de uma subscription (nova tabela)
	ListBySubscriptionID(ctx context.Context, tenantID, subscriptionID string) ([]*entity.ContaReceber, error)

	// SumByReceivedDate soma valores recebidos em um período (para fluxo de caixa)
	SumByReceivedDate(ctx context.Context, tenantID string, inicio, fim time.Time) (valueobject.Money, error)

	// SumByConfirmedDate soma valores confirmados em um período (para DRE regime competência)
	SumByConfirmedDate(ctx context.Context, tenantID string, inicio, fim time.Time) (valueobject.Money, error)

	// MarcarRecebidaViaAsaas marca conta como recebida via webhook Asaas
	MarcarRecebidaViaAsaas(ctx context.Context, tenantID, asaasPaymentID string, dataRecebimento time.Time, valorPago valueobject.Money) error

	// EstornarViaAsaas estorna conta via webhook Asaas
	EstornarViaAsaas(ctx context.Context, tenantID, asaasPaymentID, observacao string) error
}

// ContaReceberListFilters filtros para listagem de contas a receber
type ContaReceberListFilters struct {
	Status       *valueobject.StatusConta
	Origem       *string
	AssinaturaID *string
	Page         int
	PageSize     int
	OrderBy      string
}
