-- ============================================================
-- QUERIES SQLC — MÓDULO ASSINATURAS DE CLIENTES
-- Referência: FLUXO_ASSINATURA.md
-- Data: 03/12/2025
-- ============================================================

-- ============================================================
-- PLANS (Planos de Assinatura de Clientes)
-- ============================================================

-- name: CreatePlan :one
-- Criar novo plano de assinatura
INSERT INTO plans (
    tenant_id, 
    nome, 
    descricao, 
    valor, 
    periodicidade, 
    qtd_servicos, 
    limite_uso_mensal, 
    ativo
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPlanByID :one
-- Buscar plano por ID (sempre com tenant_id)
SELECT * FROM plans 
WHERE id = $1 AND tenant_id = $2;

-- name: ListPlansByTenant :many
-- Listar todos os planos de um tenant
SELECT * FROM plans 
WHERE tenant_id = $1 
ORDER BY nome;

-- name: ListActivePlansByTenant :many
-- Listar apenas planos ativos (para seleção em nova assinatura - REGRA PL-002)
SELECT * FROM plans 
WHERE tenant_id = $1 AND ativo = true 
ORDER BY nome;

-- name: UpdatePlan :one
-- Atualizar plano existente
UPDATE plans SET 
    nome = $3, 
    descricao = $4, 
    valor = $5, 
    qtd_servicos = $6, 
    limite_uso_mensal = $7, 
    ativo = $8,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeactivatePlan :exec
-- Desativar plano (soft delete - REGRA PL-003)
UPDATE plans SET 
    ativo = false, 
    updated_at = NOW() 
WHERE id = $1 AND tenant_id = $2;

-- name: CountActiveSubscriptionsByPlan :one
-- Contar assinaturas ativas de um plano (para validar REGRA PL-003)
SELECT COUNT(*)::int FROM subscriptions 
WHERE plano_id = $1 AND tenant_id = $2 AND status = 'ATIVO';

-- name: CheckPlanNameExists :one
-- Verificar se nome de plano já existe no tenant (REGRA PL-005)
SELECT EXISTS(
    SELECT 1 FROM plans 
    WHERE tenant_id = $1 AND nome = $2 AND id != $3
) AS exists;

-- ============================================================
-- SUBSCRIPTIONS (Assinaturas de Clientes)
-- ============================================================

-- name: CreateSubscription :one
-- Criar nova assinatura
INSERT INTO subscriptions (
    tenant_id, 
    cliente_id, 
    plano_id, 
    asaas_customer_id, 
    asaas_subscription_id,
    forma_pagamento, 
    status, 
    valor, 
    link_pagamento, 
    codigo_transacao,
    data_ativacao, 
    data_vencimento
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetSubscriptionByID :one
-- Buscar assinatura por ID com dados de plano e cliente (JOIN)
SELECT 
    s.*,
    p.nome as plano_nome,
    p.qtd_servicos as plano_qtd_servicos,
    p.limite_uso_mensal as plano_limite_uso_mensal,
    c.nome as cliente_nome,
    c.telefone as cliente_telefone,
    c.email as cliente_email
FROM subscriptions s
JOIN plans p ON s.plano_id = p.id
JOIN clientes c ON s.cliente_id = c.id
WHERE s.id = $1 AND s.tenant_id = $2;

-- name: ListSubscriptionsByTenant :many
-- Listar todas as assinaturas de um tenant com dados de plano e cliente
SELECT 
    s.*,
    p.nome as plano_nome,
    c.nome as cliente_nome,
    c.telefone as cliente_telefone
FROM subscriptions s
JOIN plans p ON s.plano_id = p.id
JOIN clientes c ON s.cliente_id = c.id
WHERE s.tenant_id = $1
ORDER BY s.created_at DESC;

-- name: ListSubscriptionsByStatus :many
-- Listar assinaturas por status específico
SELECT 
    s.*,
    p.nome as plano_nome,
    c.nome as cliente_nome,
    c.telefone as cliente_telefone
FROM subscriptions s
JOIN plans p ON s.plano_id = p.id
JOIN clientes c ON s.cliente_id = c.id
WHERE s.tenant_id = $1 AND s.status = $2
ORDER BY s.created_at DESC;

-- name: ListSubscriptionsByCliente :many
-- Listar assinaturas de um cliente específico
SELECT 
    s.*,
    p.nome as plano_nome
FROM subscriptions s
JOIN plans p ON s.plano_id = p.id
WHERE s.cliente_id = $1 AND s.tenant_id = $2
ORDER BY s.created_at DESC;

-- name: ListOverdueSubscriptions :many
-- Buscar assinaturas vencidas para o cron job (RN-VENC-003, RN-VENC-004)
SELECT 
    s.*,
    p.nome as plano_nome,
    c.nome as cliente_nome,
    c.id as cliente_id
FROM subscriptions s
JOIN plans p ON s.plano_id = p.id
JOIN clientes c ON s.cliente_id = c.id
WHERE s.status = 'ATIVO' 
  AND s.data_vencimento < NOW() - INTERVAL '3 days';

-- name: ListExpiringSoon :many
-- Buscar assinaturas que vencem nos próximos N dias (para notificações)
SELECT 
    s.*,
    p.nome as plano_nome,
    c.nome as cliente_nome,
    c.telefone as cliente_telefone
FROM subscriptions s
JOIN plans p ON s.plano_id = p.id
JOIN clientes c ON s.cliente_id = c.id
WHERE s.tenant_id = $1
  AND s.status = 'ATIVO' 
  AND s.data_vencimento BETWEEN NOW() AND NOW() + ($2 || ' days')::INTERVAL;

-- name: UpdateSubscription :exec
-- Atualizar assinatura completa (usado pelo webhook após pagamento)
UPDATE subscriptions SET 
    plano_id = $3,
    forma_pagamento = $4,
    status = $5,
    valor = $6,
    data_ativacao = $7,
    data_vencimento = $8,
    asaas_customer_id = $9,
    asaas_subscription_id = $10,
    link_pagamento = $11,
    codigo_transacao = $12,
    servicos_utilizados = $13,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: UpdateSubscriptionStatus :exec
-- Atualizar apenas o status de uma assinatura
UPDATE subscriptions SET 
    status = $3, 
    updated_at = NOW() 
WHERE id = $1 AND tenant_id = $2;

-- name: ActivateSubscription :exec
-- Ativar assinatura (após pagamento confirmado)
UPDATE subscriptions SET 
    status = 'ATIVO', 
    data_ativacao = $3, 
    data_vencimento = $4,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: CancelSubscription :exec
-- Cancelar assinatura (RN-CANC-003)
UPDATE subscriptions SET 
    status = 'CANCELADO', 
    data_cancelamento = NOW(), 
    cancelado_por = $3,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: UpdateSubscriptionAsaasIDs :exec
-- Atualizar IDs do Asaas após criação no gateway
UPDATE subscriptions SET 
    asaas_customer_id = $3, 
    asaas_subscription_id = $4, 
    link_pagamento = $5,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: GetSubscriptionByAsaasID :one
-- Buscar assinatura pelo ID do Asaas (para webhooks)
SELECT s.*, p.nome as plano_nome
FROM subscriptions s
JOIN plans p ON s.plano_id = p.id
WHERE s.asaas_subscription_id = $1;

-- name: IncrementServicosUtilizados :exec
-- Incrementar contador de serviços utilizados (RN-BEN-002)
UPDATE subscriptions SET 
    servicos_utilizados = servicos_utilizados + 1,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: ResetServicosUtilizados :exec
-- Resetar contador de serviços na renovação (RN-BEN-004)
UPDATE subscriptions SET 
    servicos_utilizados = 0, 
    updated_at = NOW() 
WHERE id = $1 AND tenant_id = $2;

-- name: RenewSubscription :exec
-- Renovar assinatura manualmente (PIX/Dinheiro) - Seção 6.4
UPDATE subscriptions SET 
    status = 'ATIVO',
    data_ativacao = $3,
    data_vencimento = $4,
    servicos_utilizados = 0,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: CheckActiveSubscriptionExists :one
-- Verificar se já existe assinatura ativa do mesmo plano (RN-SUB-004)
SELECT EXISTS(
    SELECT 1 FROM subscriptions 
    WHERE cliente_id = $1 AND plano_id = $2 AND status = 'ATIVO'
) AS exists;

-- name: GetSubscriptionMetrics :one
-- Métricas para relatórios (Seção 5.1)
SELECT 
    COUNT(*) FILTER (WHERE s.status = 'ATIVO')::int as total_ativas,
    COUNT(*) FILTER (WHERE s.status IN ('INATIVO', 'CANCELADO'))::int as total_inativas,
    COUNT(*) FILTER (WHERE s.status = 'INADIMPLENTE')::int as total_inadimplentes,
    (SELECT COUNT(*) FROM plans p2 WHERE p2.tenant_id = $1 AND p2.ativo = true)::int as total_planos_ativos,
    COALESCE(SUM(s.valor) FILTER (WHERE s.status = 'ATIVO'), 0)::decimal(15,2) as receita_mensal,
    -- Taxa de renovação: assinaturas renovadas nos últimos 30 dias / total ativas * 100
    CASE 
        WHEN COUNT(*) FILTER (WHERE s.status = 'ATIVO') > 0 
        THEN ROUND(
            (COUNT(*) FILTER (WHERE s.data_ativacao >= NOW() - INTERVAL '30 days' AND s.status = 'ATIVO')::numeric 
            / NULLIF(COUNT(*) FILTER (WHERE s.status = 'ATIVO'), 0)::numeric) * 100, 1
        )
        ELSE 0
    END::float as taxa_renovacao,
    -- Renovações próximos 7 dias: assinaturas com vencimento próximo
    COUNT(*) FILTER (
        WHERE s.status = 'ATIVO' 
        AND s.data_vencimento IS NOT NULL 
        AND s.data_vencimento BETWEEN NOW() AND NOW() + INTERVAL '7 days'
    )::int as renovacoes_proximos_7_dias
FROM subscriptions s
WHERE s.tenant_id = $1;

-- name: GetSubscriptionsByPlanBreakdown :many
-- Breakdown por plano (Seção 5.1)
SELECT 
    p.id as plano_id,
    p.nome as plano_nome,
    COUNT(*)::int as total,
    COALESCE(SUM(s.valor), 0)::decimal(15,2) as receita
FROM subscriptions s
JOIN plans p ON s.plano_id = p.id
WHERE s.tenant_id = $1 AND s.status = 'ATIVO'
GROUP BY p.id, p.nome
ORDER BY total DESC;

-- name: GetSubscriptionsByPaymentMethodBreakdown :many
-- Breakdown por forma de pagamento (Seção 5.1)
SELECT 
    forma_pagamento,
    COUNT(*)::int as total,
    COALESCE(SUM(valor), 0)::decimal(15,2) as receita
FROM subscriptions
WHERE tenant_id = $1 AND status = 'ATIVO'
GROUP BY forma_pagamento
ORDER BY total DESC;

-- ============================================================
-- SUBSCRIPTION_PAYMENTS (Histórico de Pagamentos)
-- ============================================================

-- name: CreateSubscriptionPayment :one
-- Registrar novo pagamento
INSERT INTO subscription_payments (
    tenant_id, 
    subscription_id, 
    asaas_payment_id, 
    valor, 
    forma_pagamento, 
    status, 
    data_pagamento, 
    codigo_transacao, 
    observacao
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: ListPaymentsBySubscription :many
-- Listar histórico de pagamentos de uma assinatura
SELECT * FROM subscription_payments 
WHERE subscription_id = $1 AND tenant_id = $2
ORDER BY created_at DESC;

-- name: UpdatePaymentStatus :exec
-- Atualizar status de um pagamento
UPDATE subscription_payments SET 
    status = $3 
WHERE id = $1 AND tenant_id = $2;

-- name: GetPaymentByAsaasID :one
-- Buscar pagamento pelo ID do Asaas (para webhooks)
SELECT * FROM subscription_payments 
WHERE asaas_payment_id = $1;

-- name: ConfirmPayment :exec
-- Confirmar pagamento (atualizar status e data)
UPDATE subscription_payments SET 
    status = 'CONFIRMADO',
    data_pagamento = $3
WHERE id = $1 AND tenant_id = $2;

-- ============================================================
-- CLIENTES (Campos de Assinatura)
-- ============================================================

-- name: UpdateClienteAsaasID :exec
-- Atualizar ID do Asaas no cliente
UPDATE clientes
SET asaas_customer_id = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: GetClienteByAsaasID :one
-- Buscar cliente pelo ID do Asaas (para unificação - RN-CLI-002)
SELECT * FROM clientes
WHERE tenant_id = $1 AND asaas_customer_id = $2
LIMIT 1;

-- name: SetClienteAsSubscriber :exec
-- Marcar/desmarcar cliente como assinante (RN-CLI-003, RN-CLI-004)
UPDATE clientes
SET is_subscriber = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: ListSubscribers :many
-- Listar todos os clientes assinantes
SELECT * FROM clientes
WHERE tenant_id = $1 AND is_subscriber = true
ORDER BY nome;

-- name: CountActiveSubscriptionsByCliente :one
-- Contar assinaturas ativas de um cliente (para RN-CLI-004)
SELECT COUNT(*)::int FROM subscriptions 
WHERE cliente_id = $1 AND tenant_id = $2 AND status = 'ATIVO';

-- name: GetClienteByNameAndPhone :one
-- Buscar cliente por nome e telefone (para busca no Asaas - AS-001)
SELECT * FROM clientes
WHERE tenant_id = $1 
  AND nome ILIKE $2 
  AND telefone = $3
LIMIT 1;

-- ============================================================
-- SUBSCRIPTION_PAYMENTS - Queries v2 (Integração Asaas)
-- ============================================================

-- name: UpsertPaymentByAsaasID :one
-- Criar ou atualizar pagamento via webhook (idempotente)
INSERT INTO subscription_payments (
    tenant_id, 
    subscription_id, 
    asaas_payment_id, 
    valor, 
    forma_pagamento, 
    status,
    status_asaas,
    due_date,
    confirmed_date,
    client_payment_date,
    credit_date,
    estimated_credit_date,
    billing_type,
    net_value,
    invoice_url,
    bank_slip_url,
    pix_qrcode,
    data_pagamento, 
    codigo_transacao, 
    observacao
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
ON CONFLICT (asaas_payment_id) 
DO UPDATE SET
    status = EXCLUDED.status,
    status_asaas = EXCLUDED.status_asaas,
    due_date = COALESCE(EXCLUDED.due_date, subscription_payments.due_date),
    confirmed_date = COALESCE(EXCLUDED.confirmed_date, subscription_payments.confirmed_date),
    client_payment_date = COALESCE(EXCLUDED.client_payment_date, subscription_payments.client_payment_date),
    credit_date = COALESCE(EXCLUDED.credit_date, subscription_payments.credit_date),
    estimated_credit_date = COALESCE(EXCLUDED.estimated_credit_date, subscription_payments.estimated_credit_date),
    billing_type = COALESCE(EXCLUDED.billing_type, subscription_payments.billing_type),
    net_value = COALESCE(EXCLUDED.net_value, subscription_payments.net_value),
    invoice_url = COALESCE(EXCLUDED.invoice_url, subscription_payments.invoice_url),
    bank_slip_url = COALESCE(EXCLUDED.bank_slip_url, subscription_payments.bank_slip_url),
    pix_qrcode = COALESCE(EXCLUDED.pix_qrcode, subscription_payments.pix_qrcode),
    data_pagamento = COALESCE(EXCLUDED.data_pagamento, subscription_payments.data_pagamento),
    observacao = COALESCE(EXCLUDED.observacao, subscription_payments.observacao),
    updated_at = NOW()
RETURNING *;

-- name: UpdatePaymentStatusAsaas :exec
-- Atualizar status interno e status Asaas
UPDATE subscription_payments SET 
    status = $3,
    status_asaas = $4,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: UpdatePaymentConfirmed :exec
-- Marcar pagamento como confirmado (CONFIRMED webhook)
UPDATE subscription_payments SET 
    status = 'CONFIRMADO',
    status_asaas = $3,
    confirmed_date = $4,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: UpdatePaymentReceived :exec
-- Marcar pagamento como recebido (RECEIVED webhook)
UPDATE subscription_payments SET 
    status = 'RECEBIDO',
    status_asaas = $3,
    client_payment_date = $4,
    credit_date = $5,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: UpdatePaymentOverdue :exec
-- Marcar pagamento como vencido (OVERDUE webhook)
UPDATE subscription_payments SET 
    status = 'VENCIDO',
    status_asaas = $3,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: UpdatePaymentRefunded :exec
-- Marcar pagamento como estornado (REFUNDED webhook)
UPDATE subscription_payments SET 
    status = 'ESTORNADO',
    status_asaas = $3,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: ListPaymentsBySubscriptionPaginated :many
-- Listar pagamentos com paginação e filtros
SELECT * FROM subscription_payments 
WHERE subscription_id = $1 AND tenant_id = $2
ORDER BY due_date DESC NULLS LAST, created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListPaymentsPendingByTenant :many
-- Listar pagamentos pendentes por tenant
SELECT sp.*, s.asaas_subscription_id, c.nome as cliente_nome
FROM subscription_payments sp
JOIN subscriptions s ON sp.subscription_id = s.id
JOIN clientes c ON s.cliente_id = c.id
WHERE sp.tenant_id = $1 
  AND sp.status IN ('PENDENTE', 'CONFIRMADO')
ORDER BY sp.due_date ASC
LIMIT $2 OFFSET $3;

-- name: ListPaymentsOverdueByTenant :many
-- Listar pagamentos vencidos por tenant
SELECT sp.*, s.asaas_subscription_id, c.nome as cliente_nome
FROM subscription_payments sp
JOIN subscriptions s ON sp.subscription_id = s.id
JOIN clientes c ON s.cliente_id = c.id
WHERE sp.tenant_id = $1 
  AND sp.status = 'VENCIDO'
ORDER BY sp.due_date ASC
LIMIT $2 OFFSET $3;

-- name: CountPaymentsByStatus :one
-- Contar pagamentos por status
SELECT 
    COUNT(*) FILTER (WHERE status = 'PENDENTE')::int as pendentes,
    COUNT(*) FILTER (WHERE status = 'CONFIRMADO')::int as confirmados,
    COUNT(*) FILTER (WHERE status = 'RECEBIDO')::int as recebidos,
    COUNT(*) FILTER (WHERE status = 'VENCIDO')::int as vencidos,
    COUNT(*) FILTER (WHERE status = 'ESTORNADO')::int as estornados
FROM subscription_payments
WHERE tenant_id = $1;

-- name: SumPaymentsByStatusAndPeriod :one
-- Somar valores por status e período (para DRE)
SELECT 
    COALESCE(SUM(valor) FILTER (WHERE status = 'CONFIRMADO'), 0)::decimal(15,2) as confirmado_bruto,
    COALESCE(SUM(net_value) FILTER (WHERE status = 'CONFIRMADO'), 0)::decimal(15,2) as confirmado_liquido,
    COALESCE(SUM(valor) FILTER (WHERE status = 'RECEBIDO'), 0)::decimal(15,2) as recebido_bruto,
    COALESCE(SUM(net_value) FILTER (WHERE status = 'RECEBIDO'), 0)::decimal(15,2) as recebido_liquido
FROM subscription_payments
WHERE tenant_id = $1 
  AND due_date >= $2 
  AND due_date <= $3;

-- ============================================================
-- SUBSCRIPTIONS - Queries v2 (Campos Asaas)
-- ============================================================

-- name: UpdateSubscriptionAsaasFields :exec
-- Atualizar campos Asaas da assinatura (após webhook)
UPDATE subscriptions SET 
    next_due_date = $3,
    asaas_status = $4,
    last_confirmed_at = $5,
    last_sync_at = NOW(),
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: UpdateSubscriptionStatusWithAsaas :exec
-- Atualizar status interno e status Asaas juntos
UPDATE subscriptions SET 
    status = $3,
    asaas_status = $4,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: UpdateSubscriptionNextDueDate :exec
-- Atualizar próximo vencimento
UPDATE subscriptions SET 
    next_due_date = $3,
    data_vencimento = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: ListSubscriptionsNeedingSync :many
-- Listar assinaturas que precisam de sync (última sync > 24h)
SELECT * FROM subscriptions
WHERE tenant_id = $1 
  AND status = 'ATIVO'
  AND (last_sync_at IS NULL OR last_sync_at < NOW() - INTERVAL '24 hours')
ORDER BY last_sync_at ASC NULLS FIRST
LIMIT $2;

-- name: ListSubscriptionsByAsaasStatus :many
-- Listar assinaturas por status Asaas
SELECT s.*, p.nome as plano_nome, c.nome as cliente_nome
FROM subscriptions s
JOIN plans p ON s.plano_id = p.id
JOIN clientes c ON s.cliente_id = c.id
WHERE s.tenant_id = $1 AND s.asaas_status = $2
ORDER BY s.created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListPaymentsNeedingReconciliation :many
-- Listar pagamentos confirmados que podem precisar de conciliação
-- (status CONFIRMADO ou RECEBIDO com asaas_payment_id)
SELECT * FROM subscription_payments
WHERE tenant_id = $1 
  AND status IN ('CONFIRMADO', 'RECEBIDO')
  AND asaas_payment_id IS NOT NULL
ORDER BY created_at DESC
LIMIT 1000;
