package subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"go.uber.org/zap"
)

// ReconcileAsaasInput define os parâmetros de entrada para conciliação
type ReconcileAsaasInput struct {
	TenantID string
	// DataInicio é a data de início do período de conciliação (opcional)
	DataInicio *time.Time
	// DataFim é a data de fim do período de conciliação (opcional)
	DataFim *time.Time
	// FullSync força sincronização completa (ignora período)
	FullSync bool
	// AutoFix se true, cria automaticamente ContasReceber faltantes
	AutoFix bool
}

// ReconcileAsaasOutput retorna o resultado da conciliação
type ReconcileAsaasOutput struct {
	TotalProcessed    int
	TotalMatched      int
	TotalMissingNexo  int // Existe no Asaas mas não no NEXO
	TotalMissingAsaas int // Existe no NEXO mas não no Asaas
	TotalDivergent    int // Valores diferentes
	TotalAutoFixed    int // Contas criadas automaticamente
	Errors            []string
	ReconciliationID  string
}

// ReconcileAsaasUseCase implementa a conciliação entre Asaas e NEXO
// Alinhado com PLANO_AJUSTE_ASAAS.md - Sprint 4
type ReconcileAsaasUseCase struct {
	subscriptionRepo   port.SubscriptionPaymentRepository
	contasReceberRepo  port.ContaReceberRepository
	webhookLogRepo     port.AsaasWebhookLogRepository
	reconciliationRepo port.AsaasReconciliationLogRepository
	logger             *zap.Logger
}

// NewReconcileAsaasUseCase cria nova instância do use case
func NewReconcileAsaasUseCase(
	subscriptionRepo port.SubscriptionPaymentRepository,
	contasReceberRepo port.ContaReceberRepository,
	webhookLogRepo port.AsaasWebhookLogRepository,
	reconciliationRepo port.AsaasReconciliationLogRepository,
	logger *zap.Logger,
) *ReconcileAsaasUseCase {
	return &ReconcileAsaasUseCase{
		subscriptionRepo:   subscriptionRepo,
		contasReceberRepo:  contasReceberRepo,
		webhookLogRepo:     webhookLogRepo,
		reconciliationRepo: reconciliationRepo,
		logger:             logger,
	}
}

// Execute executa a conciliação entre Asaas e NEXO
func (uc *ReconcileAsaasUseCase) Execute(ctx context.Context, input ReconcileAsaasInput) (*ReconcileAsaasOutput, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	// Default: últimas 24h se não informado
	periodStart := time.Now().AddDate(0, 0, -1)
	periodEnd := time.Now()
	if input.DataInicio != nil {
		periodStart = *input.DataInicio
	}
	if input.DataFim != nil {
		periodEnd = *input.DataFim
	}

	output := &ReconcileAsaasOutput{
		Errors: make([]string, 0),
	}

	// 1. Verificar pagamentos pendentes de processamento nos webhooks
	unprocessedWebhooks, err := uc.webhookLogRepo.ListUnprocessed(ctx, input.TenantID, 1000)
	if err != nil {
		output.Errors = append(output.Errors, fmt.Sprintf("erro ao listar webhooks não processados: %v", err))
		uc.logger.Error("erro ao listar webhooks não processados", zap.Error(err))
	} else {
		for _, wh := range unprocessedWebhooks {
			output.TotalProcessed++
			uc.logger.Warn("webhook não processado encontrado",
				zap.String("webhook_id", wh.ID.String()),
				zap.String("event_type", wh.EventType),
			)
		}
	}

	// 2. Verificar subscription_payments sem conta_receber correspondente
	// (CONFIRMED mas sem conta criada)
	paymentsToReconcile, err := uc.subscriptionRepo.ListNeedingReconciliation(ctx, input.TenantID)
	if err != nil {
		output.Errors = append(output.Errors, fmt.Sprintf("erro ao listar pagamentos para conciliação: %v", err))
		uc.logger.Error("erro ao listar pagamentos para conciliação", zap.Error(err))
	} else {
		for _, payment := range paymentsToReconcile {
			output.TotalProcessed++

			// Pular se não tiver asaas_payment_id
			if payment.AsaasPaymentID == nil || *payment.AsaasPaymentID == "" {
				continue
			}

			// Verificar se existe conta a receber correspondente
			conta, err := uc.contasReceberRepo.GetByAsaasPaymentID(ctx, input.TenantID, *payment.AsaasPaymentID)
			if err != nil {
				output.TotalMissingNexo++

				// Auto-fix: criar ContaReceber se solicitado
				if input.AutoFix {
					if createErr := uc.createContaReceberFromPayment(ctx, input.TenantID, payment); createErr != nil {
						output.Errors = append(output.Errors, fmt.Sprintf("erro ao criar conta para payment %s: %v", *payment.AsaasPaymentID, createErr))
						uc.logger.Error("erro ao criar conta a receber automaticamente",
							zap.String("payment_id", *payment.AsaasPaymentID),
							zap.Error(createErr),
						)
					} else {
						output.TotalAutoFixed++
						uc.logger.Info("conta a receber criada automaticamente",
							zap.String("payment_id", *payment.AsaasPaymentID),
							zap.String("valor", payment.Valor.String()),
						)
					}
				} else {
					output.Errors = append(output.Errors, fmt.Sprintf("conta não encontrada para payment %s", *payment.AsaasPaymentID))
					uc.logger.Warn("conta a receber não encontrada para payment",
						zap.String("payment_id", *payment.AsaasPaymentID),
						zap.String("status", string(payment.Status)),
					)
				}
			} else {
				// Verificar se valores batem
				contaDecimal := conta.Valor.Value()
				if !payment.Valor.Equal(contaDecimal) {
					output.TotalDivergent++
					output.Errors = append(output.Errors, fmt.Sprintf(
						"valor divergente para payment %s: Payment=%s, Conta=%s",
						*payment.AsaasPaymentID,
						payment.Valor.String(),
						conta.Valor.String(),
					))
				} else {
					output.TotalMatched++
				}
			}
		}
	}

	// 3. Verificar contas a receber órfãs (sem subscription_payment correspondente)
	// Filtrar por origem ASSINATURA
	contasOrfas, err := uc.listContasAssinaturaOrfas(ctx, input.TenantID)
	if err != nil {
		output.Errors = append(output.Errors, fmt.Sprintf("erro ao listar contas órfãs: %v", err))
		uc.logger.Error("erro ao listar contas órfãs", zap.Error(err))
	} else {
		output.TotalMissingAsaas = len(contasOrfas)
		for _, conta := range contasOrfas {
			asaasID := ""
			if conta.AsaasPaymentID != nil {
				asaasID = *conta.AsaasPaymentID
			}
			uc.logger.Warn("conta a receber sem payment correspondente",
				zap.String("conta_id", conta.ID),
				zap.String("asaas_payment_id", asaasID),
			)
		}
	}

	// 4. Criar log de conciliação
	if uc.reconciliationRepo != nil {
		reconciliationLog := entity.NewAsaasReconciliationLog(input.TenantID, periodStart, periodEnd)

		// Converter erros para JSON
		var detailsJSON *string
		if len(output.Errors) > 0 {
			detailsStr := stringSliceToJSON(output.Errors)
			detailsJSON = &detailsStr
		}

		reconciliationLog.SetResults(
			output.TotalProcessed,                            // totalAsaas (simplificado)
			output.TotalProcessed,                            // totalNexo
			output.TotalMissingNexo+output.TotalMissingAsaas, // divergences
			output.TotalAutoFixed,                            // autoFixed
			output.TotalDivergent,                            // pendingReview
			detailsJSON,
		)

		if err := uc.reconciliationRepo.Create(ctx, reconciliationLog); err != nil {
			uc.logger.Warn("erro ao criar log de conciliação", zap.Error(err))
		} else {
			output.ReconciliationID = reconciliationLog.ID
		}
	}

	uc.logger.Info("Conciliação Asaas concluída",
		zap.String("tenant_id", input.TenantID),
		zap.Int("total_processed", output.TotalProcessed),
		zap.Int("total_matched", output.TotalMatched),
		zap.Int("total_missing_nexo", output.TotalMissingNexo),
		zap.Int("total_missing_asaas", output.TotalMissingAsaas),
		zap.Int("total_divergent", output.TotalDivergent),
		zap.Int("total_auto_fixed", output.TotalAutoFixed),
	)

	return output, nil
}

// listContasAssinaturaOrfas lista contas de assinatura sem payment correspondente
func (uc *ReconcileAsaasUseCase) listContasAssinaturaOrfas(ctx context.Context, tenantID string) ([]*entity.ContaReceber, error) {
	// TODO: Implementar query específica para contas órfãs
	// Por enquanto, retornar vazio para não quebrar build
	return []*entity.ContaReceber{}, nil
}

// createContaReceberFromPayment cria uma ContaReceber a partir de um SubscriptionPayment
func (uc *ReconcileAsaasUseCase) createContaReceberFromPayment(ctx context.Context, tenantID string, payment *entity.SubscriptionPayment) error {
	// Determinar data de vencimento
	dataVencimento := time.Now()
	if payment.DueDate != nil {
		dataVencimento = *payment.DueDate
	}

	// Criar valor como Money (decimal.Decimal -> Money)
	valor := valueobject.NewMoneyFromDecimal(payment.Valor)

	// Criar descrição
	descricao := "Assinatura"
	if payment.BillingType != nil {
		descricao = fmt.Sprintf("Assinatura - %s", *payment.BillingType)
	}

	// Criar ContaReceber
	subscriptionIDStr := payment.SubscriptionID.String()
	conta, err := entity.NewContaReceber(
		tenantID,
		"ASSINATURA",
		nil, // assinaturaID (tabela legada)
		descricao,
		valor,
		dataVencimento,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar entidade ContaReceber: %w", err)
	}

	// Definir campos de integração Asaas
	conta.SubscriptionID = &subscriptionIDStr
	conta.AsaasPaymentID = payment.AsaasPaymentID

	// Determinar competência (mês do vencimento)
	competencia := dataVencimento.Format("2006-01")
	conta.CompetenciaMes = &competencia

	// Se já confirmado, atualizar status e data
	if payment.Status == entity.PaymentStatusConfirmado || payment.Status == entity.PaymentStatusRecebido {
		conta.Status = valueobject.StatusContaPago
		if payment.ConfirmedDate != nil {
			conta.ConfirmedAt = payment.ConfirmedDate
		}
		if payment.DataPagamento != nil {
			conta.DataRecebimento = payment.DataPagamento
			conta.ReceivedAt = payment.DataPagamento
		}
		conta.ValorPago = valor
		conta.ValorAberto = valueobject.Zero()
	}

	// Usar UpsertByAsaasPaymentID para idempotência
	if err := uc.contasReceberRepo.UpsertByAsaasPaymentID(ctx, conta); err != nil {
		return fmt.Errorf("erro ao persistir ContaReceber: %w", err)
	}

	return nil
}

// stringSliceToJSON converte slice de strings para JSON string
func stringSliceToJSON(s []string) string {
	if len(s) == 0 {
		return "[]"
	}
	result := "["
	for i, err := range s {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf(`"%s"`, err)
	}
	result += "]"
	return result
}
