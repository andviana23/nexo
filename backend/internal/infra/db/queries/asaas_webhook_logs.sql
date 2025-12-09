-- ============================================================
-- ASAAS_WEBHOOK_LOGS (Auditoria de Webhooks)
-- ============================================================

-- name: CreateWebhookLog :one
-- Registrar webhook recebido
INSERT INTO asaas_webhook_logs (
    tenant_id,
    event_type,
    asaas_payment_id,
    asaas_subscription_id,
    payload
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetWebhookLogByID :one
SELECT * FROM asaas_webhook_logs
WHERE id = $1;

-- name: GetWebhookLogByPaymentID :one
-- Buscar log por payment ID (para verificar duplicatas)
SELECT * FROM asaas_webhook_logs
WHERE asaas_payment_id = $1
  AND event_type = $2
  AND processed_at IS NOT NULL
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkWebhookProcessed :exec
-- Marcar webhook como processado com sucesso
UPDATE asaas_webhook_logs SET
    processed_at = NOW()
WHERE id = $1;

-- name: MarkWebhookFailed :exec
-- Marcar webhook como falha
UPDATE asaas_webhook_logs SET
    error_message = $2,
    retry_count = retry_count + 1
WHERE id = $1;

-- name: ListUnprocessedWebhooks :many
-- Listar webhooks não processados (para retry)
SELECT * FROM asaas_webhook_logs
WHERE processed_at IS NULL
  AND retry_count < 5
  AND created_at > NOW() - INTERVAL '7 days'
ORDER BY created_at ASC
LIMIT $1;

-- name: ListWebhooksByTenant :many
-- Listar webhooks por tenant (para auditoria)
SELECT * FROM asaas_webhook_logs
WHERE tenant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListWebhooksByPaymentID :many
-- Histórico de webhooks de um payment
SELECT * FROM asaas_webhook_logs
WHERE asaas_payment_id = $1
ORDER BY created_at DESC;

-- name: ListWebhooksBySubscriptionID :many
-- Histórico de webhooks de uma subscription
SELECT * FROM asaas_webhook_logs
WHERE asaas_subscription_id = $1
ORDER BY created_at DESC;

-- name: CountWebhooksByEventType :many
-- Estatísticas de webhooks por tipo (últimos 30 dias)
SELECT 
    event_type,
    COUNT(*)::int as total,
    COUNT(*) FILTER (WHERE processed_at IS NOT NULL)::int as processed,
    COUNT(*) FILTER (WHERE processed_at IS NULL AND retry_count >= 5)::int as failed
FROM asaas_webhook_logs
WHERE created_at > NOW() - INTERVAL '30 days'
GROUP BY event_type
ORDER BY total DESC;

-- name: CleanupOldWebhookLogs :exec
-- Limpar logs antigos (manter últimos 90 dias)
DELETE FROM asaas_webhook_logs
WHERE created_at < NOW() - INTERVAL '90 days'
  AND processed_at IS NOT NULL;

-- ============================================================
-- ASAAS_RECONCILIATION_LOGS (Auditoria de Conciliação)
-- ============================================================

-- name: CreateReconciliationLog :one
-- Registrar execução de conciliação
INSERT INTO asaas_reconciliation_logs (
    tenant_id,
    period_start,
    period_end,
    total_asaas,
    total_nexo,
    divergences,
    auto_fixed,
    pending_review,
    details
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetReconciliationLogByID :one
SELECT * FROM asaas_reconciliation_logs
WHERE id = $1 AND tenant_id = $2;

-- name: ListReconciliationLogs :many
-- Listar logs de conciliação
SELECT * FROM asaas_reconciliation_logs
WHERE tenant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetLastReconciliation :one
-- Última conciliação executada
SELECT * FROM asaas_reconciliation_logs
WHERE tenant_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: SumReconciliationDivergences :one
-- Somar divergências do período
SELECT 
    COALESCE(SUM(divergences), 0)::int as total_divergences,
    COALESCE(SUM(auto_fixed), 0)::int as total_auto_fixed,
    COALESCE(SUM(pending_review), 0)::int as total_pending_review
FROM asaas_reconciliation_logs
WHERE tenant_id = $1
  AND created_at >= $2
  AND created_at <= $3;
