package subscription

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateSubscriptionUseCase cria uma nova assinatura seguindo as regras do fluxo
// Reference: FLUXO_ASSINATURA.md — Seção 6.1 (Cartão), 6.2 (PIX), 6.3 (Dinheiro)
type CreateSubscriptionUseCase struct {
	planRepo     port.PlanRepository
	subRepo      port.SubscriptionRepository
	customerRepo port.CustomerRepository
	asaasGateway port.AsaasGateway // Integração Asaas (pode ser nil para PIX/Dinheiro)
	logger       *zap.Logger
}

// NewCreateSubscriptionUseCase constrói o use case
func NewCreateSubscriptionUseCase(
	planRepo port.PlanRepository,
	subRepo port.SubscriptionRepository,
	customerRepo port.CustomerRepository,
	asaasGateway port.AsaasGateway,
	logger *zap.Logger,
) *CreateSubscriptionUseCase {
	return &CreateSubscriptionUseCase{
		planRepo:     planRepo,
		subRepo:      subRepo,
		customerRepo: customerRepo,
		asaasGateway: asaasGateway,
		logger:       logger,
	}
}

// Execute cria a assinatura
func (uc *CreateSubscriptionUseCase) Execute(ctx context.Context, tenantID string, req dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidTenantID
	}
	clienteUUID, err := uuid.Parse(req.ClienteID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}
	planoUUID, err := uuid.Parse(req.PlanoID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	// Validar cliente
	cliente, err := uc.customerRepo.FindByID(ctx, tenantID, req.ClienteID)
	if err != nil || cliente == nil {
		return nil, domain.ErrCustomerNotFound
	}

	// Validar plano
	plano, err := uc.planRepo.GetByID(ctx, planoUUID, tenantUUID)
	if err != nil {
		return nil, err
	}
	if plano == nil {
		return nil, domain.ErrPlanNotFound
	}
	if !plano.Ativo {
		return nil, domain.ErrPlanInactive
	}

	// Verificar duplicidade
	exists, err := uc.subRepo.CheckActiveExists(ctx, clienteUUID, planoUUID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrSubscriptionDuplicateActive
	}

	forma := entity.PaymentMethod(req.FormaPagamento)
	if forma != entity.PaymentMethodCartao && forma != entity.PaymentMethodPix && forma != entity.PaymentMethodDinheiro {
		return nil, domain.ErrSubscriptionPaymentMethodInvalid
	}

	// Criar entidade base
	sub, err := entity.NewSubscription(tenantUUID, clienteUUID, planoUUID, forma, plano.Valor)
	if err != nil {
		return nil, err
	}
	sub.PlanoNome = plano.Nome
	sub.ClienteNome = cliente.Nome
	sub.ClienteTelefone = cliente.Telefone

	now := time.Now()
	switch forma {
	case entity.PaymentMethodCartao:
		sub.Status = entity.StatusAguardandoPagamento
		sub.DataAtivacao = nil
		sub.DataVencimento = nil

		// Integração Asaas para pagamento via cartão
		if uc.asaasGateway != nil {
			asaasResult, err := uc.integrateWithAsaas(ctx, sub, cliente, plano)
			if err != nil {
				// Log error but allow fallback to manual (AS-013)
				uc.logger.Warn("asaas integration failed, subscription created without Asaas",
					zap.String("cliente_id", clienteUUID.String()),
					zap.Error(err),
				)
				// Subscription will be created without Asaas IDs
				// Can be synced later via cron or manual action
			} else {
				// Update subscription with Asaas data
				sub.AsaasCustomerID = &asaasResult.asaasCustomerID
				sub.AsaasSubscriptionID = &asaasResult.asaasSubscriptionID
				sub.LinkPagamento = &asaasResult.paymentLink

				// Update cliente with Asaas customer ID if needed (RN-CLI-002)
				if asaasResult.asaasCustomerID != "" {
					_ = uc.subRepo.UpdateClienteAsaasID(ctx, clienteUUID, tenantUUID, &asaasResult.asaasCustomerID)
				}
			}
		}

	case entity.PaymentMethodPix, entity.PaymentMethodDinheiro:
		sub.Status = entity.StatusAtivo
		sub.DataAtivacao = &now
		venc := now.AddDate(0, 0, 30)
		sub.DataVencimento = &venc
	}

	// Salvar
	if err := uc.subRepo.Create(ctx, sub); err != nil {
		uc.logger.Error("erro ao criar assinatura", zap.Error(err))
		return nil, err
	}

	// Marcar cliente como assinante se aplicável
	if sub.ShouldMarkClientAsSubscriber() {
		_ = uc.subRepo.SetClienteAsSubscriber(ctx, clienteUUID, tenantUUID, true)
	}

	return mapper.SubscriptionToResponse(sub), nil
}

// asaasIntegrationResult holds the result of Asaas integration
type asaasIntegrationResult struct {
	asaasCustomerID     string
	asaasSubscriptionID string
	paymentLink         string
}

// integrateWithAsaas handles the Asaas integration for credit card subscriptions
// Reference: FLUXO_ASSINATURA.md — Seção 6.1
func (uc *CreateSubscriptionUseCase) integrateWithAsaas(
	ctx context.Context,
	sub *entity.Subscription,
	cliente *entity.Customer,
	plano *entity.Plan,
) (*asaasIntegrationResult, error) {
	// Step 1: Find or create customer in Asaas
	// Reference: REGRA AS-001, AS-002, AS-003, RN-CLI-002
	email := ""
	if cliente.Email != nil {
		email = *cliente.Email
	}
	cpf := ""
	if cliente.CPF != nil {
		cpf = *cliente.CPF
	}

	customerParams := port.FindOrCreateCustomerParams{
		Name:              cliente.Nome,
		Email:             email,
		MobilePhone:       cliente.Telefone,
		CpfCnpj:           cpf, // Optional per REGRA AS-002
		ExternalReference: cliente.ID,
	}

	customerResult, err := uc.asaasGateway.FindOrCreateCustomer(ctx, customerParams)
	if err != nil {
		uc.logger.Error("failed to find or create customer in Asaas",
			zap.String("cliente_id", cliente.ID),
			zap.Error(err),
		)
		return nil, err
	}

	uc.logger.Info("asaas customer ready",
		zap.String("asaas_customer_id", customerResult.AsaasCustomerID),
		zap.Bool("was_created", customerResult.WasCreated),
	)

	// Step 2: Create subscription in Asaas
	// Reference: REGRA AS-002
	subscriptionParams := port.CreateAsaasSubscriptionParams{
		CustomerID:        customerResult.AsaasCustomerID,
		Value:             plano.Valor.InexactFloat64(),
		Description:       plano.Nome,
		ExternalReference: sub.ID.String(),
		NextDueDate:       "", // Use default (tomorrow)
	}

	subscriptionResult, err := uc.asaasGateway.CreateSubscription(ctx, subscriptionParams)
	if err != nil {
		uc.logger.Error("failed to create subscription in Asaas",
			zap.String("asaas_customer_id", customerResult.AsaasCustomerID),
			zap.Error(err),
		)
		return nil, err
	}

	uc.logger.Info("asaas subscription created",
		zap.String("asaas_subscription_id", subscriptionResult.AsaasSubscriptionID),
		zap.String("payment_link", subscriptionResult.PaymentLink),
	)

	return &asaasIntegrationResult{
		asaasCustomerID:     customerResult.AsaasCustomerID,
		asaasSubscriptionID: subscriptionResult.AsaasSubscriptionID,
		paymentLink:         subscriptionResult.PaymentLink,
	}, nil
}
