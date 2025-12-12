package subscription

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/andviana23/barber-analytics-backend/internal/infra/gateway/asaas"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// ProcessWebhookUseCaseV2 processa webhooks do Asaas com tratamento correto de eventos
// Reference: PLANO_AJUSTE_ASAAS.md
// Principais mudanças:
// - CONFIRMED → competência (DRE), gera conta a receber
// - RECEIVED → caixa, quita conta a receber
// - Idempotência via asaas_payment_id
// - Log de webhooks para auditoria
type ProcessWebhookUseCaseV2 struct {
	subRepo          port.SubscriptionRepository
	paymentRepo      port.SubscriptionPaymentRepository
	contaReceberRepo port.ContaReceberRepository
	webhookLogRepo   port.AsaasWebhookLogRepository
	caixaRepo        port.CaixaDiarioRepository // T-ASAAS-001: Lançar no caixa
	logger           *zap.Logger
}

// NewProcessWebhookUseCaseV2 cria o use case refatorado
func NewProcessWebhookUseCaseV2(
	subRepo port.SubscriptionRepository,
	paymentRepo port.SubscriptionPaymentRepository,
	contaReceberRepo port.ContaReceberRepository,
	webhookLogRepo port.AsaasWebhookLogRepository,
	caixaRepo port.CaixaDiarioRepository,
	logger *zap.Logger,
) *ProcessWebhookUseCaseV2 {
	return &ProcessWebhookUseCaseV2{
		subRepo:          subRepo,
		paymentRepo:      paymentRepo,
		contaReceberRepo: contaReceberRepo,
		webhookLogRepo:   webhookLogRepo,
		caixaRepo:        caixaRepo,
		logger:           logger,
	}
}

// Execute processa o webhook com logging e idempotência
func (uc *ProcessWebhookUseCaseV2) Execute(ctx context.Context, event asaas.WebhookEvent, rawPayload []byte, clientIP net.IP) error {
	uc.logger.Info("processing webhook event v2",
		zap.String("event", event.Event),
	)

	// 1. Encontrar assinatura pelo Asaas ID
	sub, err := uc.findSubscriptionByAsaasID(ctx, getSubscriptionIDFromEvent(event))
	if err != nil {
		return err
	}
	if sub == nil {
		uc.logger.Warn("subscription not found for webhook",
			zap.String("event", event.Event),
			zap.String("asaas_subscription_id", getSubscriptionIDFromEvent(event)),
		)
		return nil // Webhook órfão - não é erro
	}

	// 2. Criar log do webhook (para auditoria)
	webhookLog := uc.createWebhookLog(sub.TenantID, event, rawPayload, clientIP)
	if uc.webhookLogRepo != nil {
		if err := uc.webhookLogRepo.Create(ctx, webhookLog); err != nil {
			uc.logger.Error("failed to log webhook", zap.Error(err))
			// Não falha - log é auxiliar
		}
	}

	// 3. Processar evento específico
	var processErr error
	switch event.Event {
	case asaas.EventPaymentCreated:
		processErr = uc.handlePaymentCreated(ctx, sub, event)

	case asaas.EventPaymentConfirmed:
		processErr = uc.handlePaymentConfirmed(ctx, sub, event)

	case asaas.EventPaymentReceived, "PAYMENT_RECEIVED_IN_CASH":
		processErr = uc.handlePaymentReceived(ctx, sub, event)

	case asaas.EventPaymentOverdue:
		processErr = uc.handlePaymentOverdue(ctx, sub, event)

	case asaas.EventPaymentRefunded:
		processErr = uc.handlePaymentRefunded(ctx, sub, event)

	case asaas.EventSubscriptionDeleted, asaas.EventSubscriptionInactivated:
		processErr = uc.handleSubscriptionCanceled(ctx, sub, event)

	case asaas.EventSubscriptionRenewed:
		processErr = uc.handleSubscriptionRenewed(ctx, sub, event)

	default:
		uc.logger.Debug("ignoring webhook event", zap.String("event", event.Event))
	}

	// 4. Atualizar status do log
	if uc.webhookLogRepo != nil {
		if processErr != nil {
			webhookLog.MarkFailed(processErr)
		} else {
			webhookLog.MarkProcessed()
		}
		// Best effort update - não falhar se log não atualizar
		_ = uc.webhookLogRepo.MarkProcessed(ctx, webhookLog.ID.String())
	}

	return processErr
}

// handlePaymentCreated processa PAYMENT_CREATED
// Ações: criar payment PENDENTE com due_date
func (uc *ProcessWebhookUseCaseV2) handlePaymentCreated(ctx context.Context, sub *entity.Subscription, event asaas.WebhookEvent) error {
	if event.Payment == nil {
		return nil
	}

	// Parse due date
	var dueDate *time.Time
	if event.Payment.DueDate != "" {
		if d, err := asaas.ParseDate(event.Payment.DueDate); err == nil {
			dueDate = &d
		}
	}

	payment := &entity.SubscriptionPayment{
		ID:             uuid.New(),
		TenantID:       sub.TenantID,
		SubscriptionID: sub.ID,
		AsaasPaymentID: &event.Payment.ID,
		Valor:          decimal.NewFromFloat(event.Payment.Value),
		FormaPagamento: sub.FormaPagamento,
		Status:         entity.PaymentStatusPendente,
		DueDate:        dueDate,
		StatusAsaas:    strPtr("PENDING"),
		BillingType:    &event.Payment.BillingType,
		NetValue:       decimal.NewFromFloat(event.Payment.NetValue),
		InvoiceURL:     strPtrIfNotEmpty(event.Payment.InvoiceUrl),
		BankSlipURL:    strPtrIfNotEmpty(event.Payment.BankSlipUrl),
		CreatedAt:      time.Now(),
	}

	// Upsert para idempotência
	if err := uc.paymentRepo.UpsertByAsaasID(ctx, payment); err != nil {
		uc.logger.Error("failed to create payment record", zap.Error(err))
		return err
	}

	uc.logger.Info("payment created via webhook",
		zap.String("subscription_id", sub.ID.String()),
		zap.String("payment_id", event.Payment.ID),
	)

	return nil
}

// handlePaymentConfirmed processa PAYMENT_CONFIRMED
// Ações:
// - Atualizar payment para CONFIRMADO com confirmed_date
// - Criar/atualizar conta a receber (competência do mês)
// - Atualizar subscription para ATIVO com last_confirmed_at
// - Marcar cliente como assinante
func (uc *ProcessWebhookUseCaseV2) handlePaymentConfirmed(ctx context.Context, sub *entity.Subscription, event asaas.WebhookEvent) error {
	if event.Payment == nil {
		return nil
	}

	// Parse dates
	var confirmedDate time.Time
	if event.Payment.ConfirmedDate != "" {
		if d, err := asaas.ParseDate(event.Payment.ConfirmedDate); err == nil {
			confirmedDate = d
		}
	}
	if confirmedDate.IsZero() {
		confirmedDate = time.Now()
	}

	var dueDate, estimatedCreditDate *time.Time
	if event.Payment.DueDate != "" {
		if d, err := asaas.ParseDate(event.Payment.DueDate); err == nil {
			dueDate = &d
		}
	}
	if event.Payment.EstimatedCreditDate != "" {
		if d, err := asaas.ParseDate(event.Payment.EstimatedCreditDate); err == nil {
			estimatedCreditDate = &d
		}
	}

	// 1. Upsert payment com status CONFIRMADO
	payment := &entity.SubscriptionPayment{
		ID:                  uuid.New(),
		TenantID:            sub.TenantID,
		SubscriptionID:      sub.ID,
		AsaasPaymentID:      &event.Payment.ID,
		Valor:               decimal.NewFromFloat(event.Payment.Value),
		FormaPagamento:      sub.FormaPagamento,
		Status:              entity.PaymentStatusConfirmado,
		DataPagamento:       &confirmedDate,
		DueDate:             dueDate,
		ConfirmedDate:       &confirmedDate,
		EstimatedCreditDate: estimatedCreditDate,
		StatusAsaas:         strPtr("CONFIRMED"),
		BillingType:         &event.Payment.BillingType,
		NetValue:            decimal.NewFromFloat(event.Payment.NetValue),
		InvoiceURL:          strPtrIfNotEmpty(event.Payment.InvoiceUrl),
		CreatedAt:           time.Now(),
	}

	if err := uc.paymentRepo.UpsertByAsaasID(ctx, payment); err != nil {
		uc.logger.Error("failed to update payment to confirmed", zap.Error(err))
		return err
	}

	// 2. Criar/atualizar conta a receber (competência)
	if uc.contaReceberRepo != nil {
		competenciaMes := confirmedDate.Format("2006-01") // YYYY-MM
		descricao := fmt.Sprintf("Assinatura - %s", sub.ClienteNome)
		if sub.ClienteNome == "" {
			descricao = fmt.Sprintf("Assinatura #%s", sub.ID.String()[:8])
		}

		conta := &entity.ContaReceber{
			ID:              uuid.NewString(),
			TenantID:        sub.TenantID,
			Origem:          "ASSINATURA",
			SubscriptionID:  strPtr(sub.ID.String()),
			AsaasPaymentID:  &event.Payment.ID,
			DescricaoOrigem: descricao,
			Valor:           valueobject.NewMoneyFromFloat(event.Payment.Value),
			ValorPago:       valueobject.Zero(),
			ValorAberto:     valueobject.NewMoneyFromFloat(event.Payment.Value),
			DataVencimento:  confirmedDate,
			Status:          valueobject.StatusContaConfirmado, // Confirmado mas não recebido ainda
			CompetenciaMes:  &competenciaMes,
			ConfirmedAt:     &confirmedDate,
			CriadoEm:        time.Now(),
			AtualizadoEm:    time.Now(),
		}

		if err := uc.contaReceberRepo.UpsertByAsaasPaymentID(ctx, conta); err != nil {
			uc.logger.Error("failed to create conta a receber", zap.Error(err))
			// Não falha - conta a receber é complementar
		}
	}

	// 3. Atualizar subscription
	expirationDate := confirmedDate.AddDate(0, 0, 30)
	sub.Status = entity.StatusAtivo
	sub.DataAtivacao = &confirmedDate
	sub.DataVencimento = &expirationDate
	sub.LastConfirmedAt = &confirmedDate
	sub.ServicosUtilizados = 0

	if err := uc.subRepo.Update(ctx, sub); err != nil {
		uc.logger.Error("failed to update subscription after confirmed", zap.Error(err))
		return err
	}

	// 4. Marcar cliente como assinante
	if err := uc.subRepo.SetClienteAsSubscriber(ctx, sub.ClienteID, sub.TenantID, true); err != nil {
		uc.logger.Error("failed to mark client as subscriber", zap.Error(err))
	}

	uc.logger.Info("payment confirmed via webhook",
		zap.String("subscription_id", sub.ID.String()),
		zap.String("payment_id", event.Payment.ID),
		zap.Time("confirmed_date", confirmedDate),
	)

	return nil
}

// handlePaymentReceived processa PAYMENT_RECEIVED / RECEIVED_IN_CASH
// Ações:
// - Atualizar payment para RECEBIDO com credit_date
// - Quitar conta a receber
// - (TODO) Lançar no fluxo de caixa diário
func (uc *ProcessWebhookUseCaseV2) handlePaymentReceived(ctx context.Context, sub *entity.Subscription, event asaas.WebhookEvent) error {
	if event.Payment == nil {
		return nil
	}

	// Parse dates
	var creditDate time.Time
	if event.Payment.CreditDate != "" {
		if d, err := asaas.ParseDate(event.Payment.CreditDate); err == nil {
			creditDate = d
		}
	}
	if creditDate.IsZero() {
		if event.Payment.PaymentDate != "" {
			if d, err := asaas.ParseDate(event.Payment.PaymentDate); err == nil {
				creditDate = d
			}
		}
	}
	if creditDate.IsZero() {
		creditDate = time.Now()
	}

	var clientPaymentDate *time.Time
	if event.Payment.ClientPaymentDate != "" {
		if d, err := asaas.ParseDate(event.Payment.ClientPaymentDate); err == nil {
			clientPaymentDate = &d
		}
	}

	// 1. Upsert payment com status RECEBIDO
	payment := &entity.SubscriptionPayment{
		ID:                uuid.New(),
		TenantID:          sub.TenantID,
		SubscriptionID:    sub.ID,
		AsaasPaymentID:    &event.Payment.ID,
		Valor:             decimal.NewFromFloat(event.Payment.Value),
		FormaPagamento:    sub.FormaPagamento,
		Status:            entity.PaymentStatusRecebido,
		DataPagamento:     &creditDate,
		CreditDate:        &creditDate,
		ClientPaymentDate: clientPaymentDate,
		StatusAsaas:       strPtr("RECEIVED"),
		BillingType:       &event.Payment.BillingType,
		NetValue:          decimal.NewFromFloat(event.Payment.NetValue),
		CreatedAt:         time.Now(),
	}

	if err := uc.paymentRepo.UpsertByAsaasID(ctx, payment); err != nil {
		uc.logger.Error("failed to update payment to received", zap.Error(err))
		return err
	}

	// 2. Quitar conta a receber
	if uc.contaReceberRepo != nil {
		conta, err := uc.contaReceberRepo.GetByAsaasPaymentID(ctx, sub.TenantID.String(), event.Payment.ID)
		if err != nil {
			uc.logger.Warn("conta a receber not found for payment", zap.String("payment_id", event.Payment.ID))
		} else if conta != nil {
			conta.Status = valueobject.StatusContaRecebido
			conta.ValorPago = conta.Valor
			conta.ValorAberto = valueobject.Zero()
			conta.DataRecebimento = &creditDate
			conta.ReceivedAt = &creditDate
			conta.AtualizadoEm = time.Now()

			if err := uc.contaReceberRepo.Update(ctx, conta); err != nil {
				uc.logger.Error("failed to update conta a receber", zap.Error(err))
			}
		}
	}

	// 3. T-ASAAS-001: Lançar no fluxo de caixa diário
	if uc.caixaRepo != nil {
		// Buscar caixa aberto do dia
		caixaAberto, err := uc.caixaRepo.FindAberto(ctx, sub.TenantID)
		if err != nil {
			uc.logger.Warn("erro ao buscar caixa aberto para pagamento Asaas",
				zap.String("tenant_id", sub.TenantID.String()),
				zap.Error(err))
		} else if caixaAberto != nil {
			// Criar operação de entrada no caixa
			valor := decimal.NewFromFloat(event.Payment.NetValue) // Usar valor líquido (após taxas Asaas)
			descricao := fmt.Sprintf("Assinatura Asaas - %s #%s", sub.ClienteNome, event.Payment.ID[:8])

			// Identificar tipo de pagamento
			var tipoOperacao string
			switch event.Payment.BillingType {
			case "PIX":
				tipoOperacao = "PIX"
			case "BOLETO":
				tipoOperacao = "BOLETO"
			case "CREDIT_CARD":
				tipoOperacao = "CARTAO_CREDITO"
			default:
				tipoOperacao = "ASAAS"
			}

			operacao, err := entity.NewOperacaoVenda(
				caixaAberto.ID,
				sub.TenantID,
				uuid.Nil, // UserID do sistema (webhook)
				valor,
				fmt.Sprintf("%s - %s", tipoOperacao, descricao),
			)
			if err != nil {
				uc.logger.Error("erro ao criar operação de caixa para webhook Asaas", zap.Error(err))
			} else {
				if err := uc.caixaRepo.CreateOperacao(ctx, operacao); err != nil {
					uc.logger.Error("erro ao registrar operação no caixa", zap.Error(err))
				} else {
					// Atualizar totais do caixa
					novoTotal := caixaAberto.TotalEntradas.Add(valor)
					caixaAberto.TotalEntradas = novoTotal
					if err := uc.caixaRepo.UpdateTotais(ctx, caixaAberto.ID, sub.TenantID, caixaAberto.TotalSangrias, caixaAberto.TotalReforcos, novoTotal); err != nil {
						uc.logger.Error("erro ao atualizar totais do caixa", zap.Error(err))
					}

					uc.logger.Info("pagamento Asaas lançado no caixa",
						zap.String("operacao_id", operacao.ID.String()),
						zap.String("valor", valor.String()),
						zap.String("tipo", tipoOperacao))
				}
			}
		} else {
			uc.logger.Warn("nenhum caixa aberto para receber pagamento Asaas",
				zap.String("tenant_id", sub.TenantID.String()),
				zap.String("payment_id", event.Payment.ID))
		}
	}

	uc.logger.Info("payment received via webhook",
		zap.String("subscription_id", sub.ID.String()),
		zap.String("payment_id", event.Payment.ID),
		zap.Time("credit_date", creditDate),
	)

	return nil
}

// handlePaymentOverdue processa PAYMENT_OVERDUE
// Ações:
// - Atualizar payment para VENCIDO
// - Atualizar subscription para INADIMPLENTE
// - Verificar se cliente tem outras assinaturas ativas
func (uc *ProcessWebhookUseCaseV2) handlePaymentOverdue(ctx context.Context, sub *entity.Subscription, event asaas.WebhookEvent) error {
	if event.Payment == nil {
		return nil
	}

	// 1. Atualizar/criar payment com status VENCIDO
	payment := &entity.SubscriptionPayment{
		ID:             uuid.New(),
		TenantID:       sub.TenantID,
		SubscriptionID: sub.ID,
		AsaasPaymentID: &event.Payment.ID,
		Valor:          decimal.NewFromFloat(event.Payment.Value),
		FormaPagamento: sub.FormaPagamento,
		Status:         entity.PaymentStatusVencido,
		StatusAsaas:    strPtr("OVERDUE"),
		BillingType:    &event.Payment.BillingType,
		CreatedAt:      time.Now(),
	}

	if err := uc.paymentRepo.UpsertByAsaasID(ctx, payment); err != nil {
		uc.logger.Error("failed to update payment to overdue", zap.Error(err))
	}

	// 2. Atualizar subscription para INADIMPLENTE
	sub.Status = entity.StatusInadimplente
	if err := uc.subRepo.Update(ctx, sub); err != nil {
		uc.logger.Error("failed to update subscription to inadimplente", zap.Error(err))
		return err
	}

	// 3. Verificar se cliente tem outras assinaturas ativas
	count, err := uc.subRepo.CountActiveSubscriptionsByCliente(ctx, sub.ClienteID, sub.TenantID)
	if err == nil && count == 0 {
		if err := uc.subRepo.SetClienteAsSubscriber(ctx, sub.ClienteID, sub.TenantID, false); err != nil {
			uc.logger.Error("failed to remove subscriber flag", zap.Error(err))
		}
	}

	uc.logger.Info("payment marked as overdue via webhook",
		zap.String("subscription_id", sub.ID.String()),
		zap.String("payment_id", event.Payment.ID),
	)

	return nil
}

// handlePaymentRefunded processa PAYMENT_REFUNDED
// Ações:
// - Atualizar payment para ESTORNADO
// - Atualizar subscription para INATIVO
// - Estornar conta a receber
func (uc *ProcessWebhookUseCaseV2) handlePaymentRefunded(ctx context.Context, sub *entity.Subscription, event asaas.WebhookEvent) error {
	if event.Payment == nil {
		return nil
	}

	now := time.Now()

	// 1. Atualizar payment para ESTORNADO
	payment := &entity.SubscriptionPayment{
		ID:             uuid.New(),
		TenantID:       sub.TenantID,
		SubscriptionID: sub.ID,
		AsaasPaymentID: &event.Payment.ID,
		Valor:          decimal.NewFromFloat(event.Payment.Value),
		FormaPagamento: sub.FormaPagamento,
		Status:         entity.PaymentStatusEstornado,
		DataPagamento:  &now,
		StatusAsaas:    strPtr("REFUNDED"),
		CreatedAt:      time.Now(),
	}

	if err := uc.paymentRepo.UpsertByAsaasID(ctx, payment); err != nil {
		uc.logger.Error("failed to update payment to refunded", zap.Error(err))
	}

	// 2. Atualizar subscription para INATIVO
	sub.Status = entity.StatusInativo
	if err := uc.subRepo.Update(ctx, sub); err != nil {
		uc.logger.Error("failed to update subscription after refund", zap.Error(err))
		return err
	}

	// 3. Estornar conta a receber
	if uc.contaReceberRepo != nil {
		conta, _ := uc.contaReceberRepo.GetByAsaasPaymentID(ctx, sub.TenantID.String(), event.Payment.ID)
		if conta != nil {
			conta.Status = valueobject.StatusContaEstornado
			conta.AtualizadoEm = time.Now()
			_ = uc.contaReceberRepo.Update(ctx, conta)
		}
	}

	// 4. Verificar subscriber flag
	count, err := uc.subRepo.CountActiveSubscriptionsByCliente(ctx, sub.ClienteID, sub.TenantID)
	if err == nil && count == 0 {
		_ = uc.subRepo.SetClienteAsSubscriber(ctx, sub.ClienteID, sub.TenantID, false)
	}

	uc.logger.Info("payment refunded via webhook",
		zap.String("subscription_id", sub.ID.String()),
		zap.String("payment_id", event.Payment.ID),
	)

	return nil
}

// handleSubscriptionCanceled processa SUBSCRIPTION_DELETED/INACTIVATED
func (uc *ProcessWebhookUseCaseV2) handleSubscriptionCanceled(ctx context.Context, sub *entity.Subscription, event asaas.WebhookEvent) error {
	now := time.Now()
	sub.Status = entity.StatusCancelado
	sub.DataCancelamento = &now

	if err := uc.subRepo.Update(ctx, sub); err != nil {
		return err
	}

	// Verificar subscriber flag
	count, _ := uc.subRepo.CountActiveSubscriptionsByCliente(ctx, sub.ClienteID, sub.TenantID)
	if count == 0 {
		_ = uc.subRepo.SetClienteAsSubscriber(ctx, sub.ClienteID, sub.TenantID, false)
	}

	uc.logger.Info("subscription canceled via webhook", zap.String("subscription_id", sub.ID.String()))
	return nil
}

// handleSubscriptionRenewed processa SUBSCRIPTION_RENEWED
func (uc *ProcessWebhookUseCaseV2) handleSubscriptionRenewed(ctx context.Context, sub *entity.Subscription, event asaas.WebhookEvent) error {
	// Parse new due date
	var newDueDate time.Time
	if event.Payment != nil && event.Payment.DueDate != "" {
		newDueDate, _ = asaas.ParseDate(event.Payment.DueDate)
	}
	if newDueDate.IsZero() {
		if sub.DataVencimento != nil {
			newDueDate = sub.DataVencimento.AddDate(0, 0, 30)
		} else {
			newDueDate = time.Now().AddDate(0, 0, 30)
		}
	}

	sub.DataVencimento = &newDueDate
	sub.NextDueDate = &newDueDate
	sub.ServicosUtilizados = 0
	sub.Status = entity.StatusAtivo

	if err := uc.subRepo.Update(ctx, sub); err != nil {
		return err
	}

	uc.logger.Info("subscription renewed via webhook",
		zap.String("subscription_id", sub.ID.String()),
		zap.Time("new_due_date", newDueDate),
	)

	return nil
}

// Helper functions

func (uc *ProcessWebhookUseCaseV2) findSubscriptionByAsaasID(ctx context.Context, asaasSubID string) (*entity.Subscription, error) {
	if asaasSubID == "" {
		return nil, nil
	}
	return uc.subRepo.GetByAsaasSubscriptionID(ctx, &asaasSubID)
}

func (uc *ProcessWebhookUseCaseV2) createWebhookLog(tenantID uuid.UUID, event asaas.WebhookEvent, rawPayload []byte, clientIP net.IP) *entity.AsaasWebhookLog {
	var paymentID, subscriptionID *string
	if event.Payment != nil {
		paymentID = &event.Payment.ID
		if event.Payment.Subscription != "" {
			subscriptionID = &event.Payment.Subscription
		}
	}

	var payload json.RawMessage
	if len(rawPayload) > 0 {
		payload = rawPayload
	} else {
		payload, _ = json.Marshal(event)
	}

	return entity.NewAsaasWebhookLog(tenantID, event.Event, paymentID, subscriptionID, payload, clientIP)
}

func getSubscriptionIDFromEvent(event asaas.WebhookEvent) string {
	if event.Payment != nil {
		return event.Payment.Subscription
	}
	return ""
}

func strPtr(s string) *string {
	return &s
}

func strPtrIfNotEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
