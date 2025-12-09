-- ============================================================================
-- QUERIES: commission_rules
-- Regras de comissão
-- ============================================================================

-- name: CreateCommissionRule :one
INSERT INTO commission_rules (
    tenant_id,
    unit_id,
    name,
    description,
    type,
    default_rate,
    min_amount,
    max_amount,
    calculation_base,
    effective_from,
    effective_to,
    priority,
    is_active,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: GetCommissionRuleByID :one
SELECT * FROM commission_rules
WHERE id = $1 AND tenant_id = $2;

-- name: ListCommissionRulesByTenant :many
SELECT * FROM commission_rules
WHERE tenant_id = $1
ORDER BY priority DESC, created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListCommissionRulesActive :many
SELECT * FROM commission_rules
WHERE tenant_id = $1
  AND is_active = true
  AND effective_from <= CURRENT_DATE
  AND (effective_to IS NULL OR effective_to >= CURRENT_DATE)
ORDER BY priority DESC, created_at DESC;

-- name: ListCommissionRulesByUnit :many
SELECT * FROM commission_rules
WHERE tenant_id = $1
  AND (unit_id = $2 OR unit_id IS NULL)
  AND is_active = true
  AND effective_from <= CURRENT_DATE
  AND (effective_to IS NULL OR effective_to >= CURRENT_DATE)
ORDER BY 
    CASE WHEN unit_id IS NOT NULL THEN 0 ELSE 1 END,
    priority DESC;

-- name: UpdateCommissionRule :one
UPDATE commission_rules
SET
    name = $3,
    description = $4,
    type = $5,
    default_rate = $6,
    min_amount = $7,
    max_amount = $8,
    calculation_base = $9,
    effective_from = $10,
    effective_to = $11,
    priority = $12,
    is_active = $13,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeactivateCommissionRule :one
UPDATE commission_rules
SET
    is_active = false,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteCommissionRule :exec
DELETE FROM commission_rules
WHERE id = $1 AND tenant_id = $2;

-- name: CountCommissionRulesByTenant :one
SELECT COUNT(*) as total
FROM commission_rules
WHERE tenant_id = $1;

-- name: GetDefaultCommissionRule :one
SELECT * FROM commission_rules
WHERE tenant_id = $1
  AND unit_id IS NULL
  AND is_active = true
  AND effective_from <= CURRENT_DATE
  AND (effective_to IS NULL OR effective_to >= CURRENT_DATE)
ORDER BY priority DESC
LIMIT 1;

-- name: GetCommissionRuleByUnit :one
-- COM-001: Busca regra vigente específica de uma unidade
SELECT * FROM commission_rules
WHERE tenant_id = $1
  AND unit_id = $2
  AND is_active = true
  AND effective_from <= $3::date
  AND (effective_to IS NULL OR effective_to >= $3::date)
ORDER BY priority DESC
LIMIT 1;

-- name: GetGlobalCommissionRule :one
-- COM-001: Busca regra vigente global do tenant (unit_id IS NULL)
SELECT * FROM commission_rules
WHERE tenant_id = $1
  AND unit_id IS NULL
  AND is_active = true
  AND effective_from <= $2::date
  AND (effective_to IS NULL OR effective_to >= $2::date)
ORDER BY priority DESC
LIMIT 1;
