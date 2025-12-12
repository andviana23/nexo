-- name: CreateContaReceber :one
INSERT INTO contas_a_receber (
    tenant_id,
    origem,
    assinatura_id,
    servico_id,
    command_id,
    command_payment_id,
    descricao,
    valor,
    valor_pago,
    data_vencimento,
    data_recebimento,
    status,
    observacoes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING *;

-- name: GetContaReceberByID :one
SELECT * FROM contas_a_receber
WHERE id = $1 AND tenant_id = $2;

-- name: ListContasReceberByTenant :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1
ORDER BY data_vencimento DESC
LIMIT $2 OFFSET $3;

-- name: ListContasReceberByStatus :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1 AND status = $2
ORDER BY data_vencimento ASC
LIMIT $3 OFFSET $4;

-- name: ListContasReceberByPeriod :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1
  AND data_vencimento >= $2
  AND data_vencimento <= $3
ORDER BY data_vencimento ASC;

-- name: ListContasReceberFiltered :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1
  AND (sqlc.narg(status)::text IS NULL OR status = sqlc.narg(status))
  AND (sqlc.narg(origem)::text IS NULL OR origem = sqlc.narg(origem))
  AND (sqlc.narg(assinatura_id)::uuid IS NULL OR assinatura_id = sqlc.narg(assinatura_id))
  AND (sqlc.narg(data_inicio)::date IS NULL OR data_vencimento >= sqlc.narg(data_inicio))
  AND (sqlc.narg(data_fim)::date IS NULL OR data_vencimento <= sqlc.narg(data_fim))
ORDER BY data_vencimento DESC
LIMIT $2 OFFSET $3;

-- name: ListContasReceberVencidas :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1
  AND status IN ('PENDENTE', 'ATRASADO')
  AND data_vencimento < $2
ORDER BY data_vencimento ASC;

-- name: ListContasReceberByAssinatura :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1 AND assinatura_id = $2
ORDER BY data_vencimento DESC;

-- name: ListContasReceberByOrigem :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1 AND origem = $2
ORDER BY data_vencimento DESC
LIMIT $3 OFFSET $4;

-- name: ListContasReceberByCommandID :many
SELECT * FROM contas_a_receber
WHERE tenant_id = $1 AND command_id = $2
ORDER BY criado_em ASC;

-- name: UpdateContaReceber :one
UPDATE contas_a_receber
SET
    descricao = $3,
    valor = $4,
    valor_pago = $5,
    data_vencimento = $6,
    data_recebimento = $7,
    status = $8,
    observacoes = $9,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: MarcarContaReceberComoRecebida :one
UPDATE contas_a_receber
SET
    status = 'RECEBIDO',
    data_recebimento = $3,
    valor_pago = $4,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: MarcarContaReceberComoAtrasada :exec
UPDATE contas_a_receber
SET
    status = 'ATRASADO',
    atualizado_em = NOW()
WHERE tenant_id = $1
  AND status = 'PENDENTE'
  AND data_vencimento < $2;

-- name: DeleteContaReceber :exec
DELETE FROM contas_a_receber
WHERE id = $1 AND tenant_id = $2;

-- name: SumContasReceberByPeriod :one
SELECT
    COALESCE(SUM(valor), 0) as total_a_receber
FROM contas_a_receber
WHERE tenant_id = $1
  AND data_vencimento >= $2
  AND data_vencimento <= $3
  AND status NOT IN ('CANCELADO', 'ESTORNADO');

-- name: SumContasReceberByOrigem :one
SELECT
    COALESCE(SUM(valor), 0) as total_por_origem
FROM contas_a_receber
WHERE tenant_id = $1
  AND origem = $2
  AND data_vencimento >= $3
  AND data_vencimento <= $4
  AND status NOT IN ('CANCELADO', 'ESTORNADO');

-- name: SumContasRecebidasByPeriod :one
SELECT
    COALESCE(SUM(valor_pago), 0) as total_recebido
FROM contas_a_receber
WHERE tenant_id = $1
  AND data_recebimento >= $2
  AND data_recebimento <= $3
  AND status = 'RECEBIDO';

-- name: CountContasReceberByStatus :one
SELECT COUNT(*) FROM contas_a_receber
WHERE tenant_id = $1 AND status = $2;

-- name: CountContasReceberByTenant :one
SELECT COUNT(*) FROM contas_a_receber
WHERE tenant_id = $1;

-- ============================================================
-- CONTAS_A_RECEBER - Queries v2 (Integração Asaas)
-- ============================================================

-- name: UpsertContaReceberByAsaasPaymentID :one
-- Criar ou atualizar conta a receber via webhook (idempotente)
-- Nota: índice único é (tenant_id, asaas_payment_id)
INSERT INTO contas_a_receber (
    tenant_id,
    origem,
    assinatura_id,
    subscription_id,
    asaas_payment_id,
    servico_id,
    descricao,
    valor,
    valor_pago,
    data_vencimento,
    data_recebimento,
    competencia_mes,
    confirmed_at,
    received_at,
    status,
    observacoes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
ON CONFLICT (tenant_id, asaas_payment_id) WHERE asaas_payment_id IS NOT NULL
DO UPDATE SET
    valor = EXCLUDED.valor,
    valor_pago = COALESCE(EXCLUDED.valor_pago, contas_a_receber.valor_pago),
    data_vencimento = COALESCE(EXCLUDED.data_vencimento, contas_a_receber.data_vencimento),
    data_recebimento = COALESCE(EXCLUDED.data_recebimento, contas_a_receber.data_recebimento),
    confirmed_at = COALESCE(EXCLUDED.confirmed_at, contas_a_receber.confirmed_at),
    received_at = COALESCE(EXCLUDED.received_at, contas_a_receber.received_at),
    status = EXCLUDED.status,
    observacoes = COALESCE(EXCLUDED.observacoes, contas_a_receber.observacoes),
    atualizado_em = NOW()
RETURNING *;

-- name: GetContaReceberByAsaasPaymentID :one
-- Buscar conta pelo payment ID do Asaas
SELECT * FROM contas_a_receber
WHERE tenant_id = $1 AND asaas_payment_id = $2;

-- name: GetContaReceberBySubscriptionID :many
-- Listar contas por subscription
SELECT * FROM contas_a_receber
WHERE tenant_id = $1 AND subscription_id = $2
ORDER BY data_vencimento DESC;

-- name: MarcarContaReceberRecebidaViaAsaas :one
-- Quitar conta quando webhook RECEIVED chegar
UPDATE contas_a_receber
SET
    status = 'RECEBIDO',
    data_recebimento = $3,
    received_at = $3,
    valor_pago = $4,
    atualizado_em = NOW()
WHERE tenant_id = $1 AND asaas_payment_id = $2
RETURNING *;

-- name: EstornarContaReceberViaAsaas :one
-- Estornar conta quando webhook REFUNDED chegar
UPDATE contas_a_receber
SET
    status = 'ESTORNADO',
    observacoes = COALESCE(observacoes || ' | ', '') || $3,
    atualizado_em = NOW()
WHERE tenant_id = $1 AND asaas_payment_id = $2
RETURNING *;

-- name: SumContasReceberByCompetencia :one
-- Somar contas por competência (para DRE por competência)
SELECT 
    COALESCE(SUM(valor), 0)::decimal(15,2) as total_bruto,
    COALESCE(SUM(valor_pago), 0)::decimal(15,2) as total_pago,
    COUNT(*)::int as quantidade
FROM contas_a_receber
WHERE tenant_id = $1 
  AND competencia_mes = $2
  AND status NOT IN ('CANCELADO', 'ESTORNADO');

-- name: SumContasReceberByCompetenciaAndStatus :one
-- Somar por competência e status específico
SELECT 
    COALESCE(SUM(valor), 0)::decimal(15,2) as total
FROM contas_a_receber
WHERE tenant_id = $1 
  AND competencia_mes = $2
  AND status = $3;

-- name: ListContasReceberByCompetencia :many
-- Listar contas por competência
SELECT * FROM contas_a_receber
WHERE tenant_id = $1 AND competencia_mes = $2
ORDER BY data_vencimento ASC;

-- name: ListContasReceberPendentesAsaas :many
-- Listar contas pendentes de assinaturas (para conciliação)
SELECT cr.*, s.asaas_subscription_id
FROM contas_a_receber cr
LEFT JOIN subscriptions s ON cr.subscription_id = s.id
WHERE cr.tenant_id = $1 
  AND cr.origem = 'ASSINATURA'
  AND cr.status IN ('PENDENTE', 'CONFIRMADO')
ORDER BY cr.data_vencimento ASC;

-- name: GetContasReceberResumoMensal :one
-- Resumo mensal para DRE
SELECT 
    COUNT(*) FILTER (WHERE status = 'PENDENTE')::int as pendentes,
    COUNT(*) FILTER (WHERE status = 'CONFIRMADO')::int as confirmados,
    COUNT(*) FILTER (WHERE status = 'RECEBIDO')::int as recebidos,
    COUNT(*) FILTER (WHERE status = 'ATRASADO')::int as atrasados,
    COUNT(*) FILTER (WHERE status = 'ESTORNADO')::int as estornados,
    COALESCE(SUM(valor) FILTER (WHERE status IN ('PENDENTE', 'CONFIRMADO', 'ATRASADO')), 0)::decimal(15,2) as valor_a_receber,
    COALESCE(SUM(valor_pago) FILTER (WHERE status = 'RECEBIDO'), 0)::decimal(15,2) as valor_recebido
FROM contas_a_receber
WHERE tenant_id = $1 
  AND competencia_mes = $2;

-- name: SumContasReceberByReceivedDate :one
-- Somar por data de recebimento (para fluxo de caixa)
SELECT 
    COALESCE(SUM(valor_pago), 0)::decimal(15,2) as total
FROM contas_a_receber
WHERE tenant_id = $1 
  AND received_at >= $2
  AND received_at < $3
  AND status = 'RECEBIDO';

-- name: SumContasReceberByConfirmedDate :one
-- Somar por data de confirmação (para DRE regime competência)
SELECT 
    COALESCE(SUM(valor), 0)::decimal(15,2) as total
FROM contas_a_receber
WHERE tenant_id = $1 
  AND confirmed_at >= $2
  AND confirmed_at < $3
  AND status IN ('CONFIRMADO', 'RECEBIDO');
