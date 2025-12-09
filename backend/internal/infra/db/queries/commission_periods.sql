-- ============================================================================
-- QUERIES: commission_periods
-- Períodos de fechamento de comissões
-- ============================================================================

-- name: CreateCommissionPeriod :one
INSERT INTO commission_periods (
    tenant_id,
    unit_id,
    reference_month,
    professional_id,
    total_gross,
    total_commission,
    total_advances,
    total_adjustments,
    total_net,
    items_count,
    status,
    period_start,
    period_end,
    notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: GetCommissionPeriodByID :one
SELECT * FROM commission_periods
WHERE id = $1 AND tenant_id = $2;

-- name: GetCommissionPeriodByProfessionalAndMonth :one
SELECT * FROM commission_periods
WHERE tenant_id = $1
  AND professional_id = $2
  AND reference_month = $3
  AND (unit_id = $4 OR ($4 IS NULL AND unit_id IS NULL));

-- name: ListCommissionPeriodsByTenant :many
SELECT cp.*,
       p.nome as professional_name,
       u.nome as unit_name
FROM commission_periods cp
LEFT JOIN profissionais p ON cp.professional_id = p.id
LEFT JOIN units u ON cp.unit_id = u.id
WHERE cp.tenant_id = $1
ORDER BY cp.reference_month DESC, p.nome ASC
LIMIT $2 OFFSET $3;

-- name: ListCommissionPeriodsByMonth :many
SELECT cp.*,
       p.nome as professional_name,
       u.nome as unit_name
FROM commission_periods cp
LEFT JOIN profissionais p ON cp.professional_id = p.id
LEFT JOIN units u ON cp.unit_id = u.id
WHERE cp.tenant_id = $1
  AND cp.reference_month = $2
ORDER BY p.nome ASC;

-- name: ListCommissionPeriodsByProfessional :many
SELECT cp.*,
       u.nome as unit_name
FROM commission_periods cp
LEFT JOIN units u ON cp.unit_id = u.id
WHERE cp.tenant_id = $1
  AND cp.professional_id = $2
ORDER BY cp.reference_month DESC
LIMIT $3 OFFSET $4;

-- name: ListCommissionPeriodsByStatus :many
SELECT cp.*,
       p.nome as professional_name,
       u.nome as unit_name
FROM commission_periods cp
LEFT JOIN profissionais p ON cp.professional_id = p.id
LEFT JOIN units u ON cp.unit_id = u.id
WHERE cp.tenant_id = $1
  AND cp.status = $2
ORDER BY cp.reference_month DESC, p.nome ASC
LIMIT $3 OFFSET $4;

-- name: ListOpenCommissionPeriods :many
SELECT cp.*,
       p.nome as professional_name,
       u.nome as unit_name
FROM commission_periods cp
LEFT JOIN profissionais p ON cp.professional_id = p.id
LEFT JOIN units u ON cp.unit_id = u.id
WHERE cp.tenant_id = $1
  AND cp.status = 'ABERTO'
ORDER BY cp.reference_month DESC, p.nome ASC;

-- name: UpdateCommissionPeriodTotals :one
UPDATE commission_periods
SET
    total_gross = $3,
    total_commission = $4,
    total_advances = $5,
    total_adjustments = $6,
    total_net = $7,
    items_count = $8,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: CloseCommissionPeriod :one
UPDATE commission_periods
SET
    status = 'FECHADO',
    closed_at = NOW(),
    closed_by = $3,
    conta_pagar_id = $4,
    notes = COALESCE($5, notes),
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: MarkCommissionPeriodAsPaid :one
UPDATE commission_periods
SET
    status = 'PAGO',
    paid_at = NOW(),
    paid_by = $3,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: CancelCommissionPeriod :one
UPDATE commission_periods
SET
    status = 'CANCELADO',
    notes = $3,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteCommissionPeriod :exec
DELETE FROM commission_periods
WHERE id = $1 AND tenant_id = $2 AND status = 'ABERTO';

-- name: SumCommissionsByMonth :one
SELECT
    COALESCE(SUM(total_commission), 0)::NUMERIC(15,2) as total_commission,
    COALESCE(SUM(total_net), 0)::NUMERIC(15,2) as total_net,
    COUNT(*) as periods_count
FROM commission_periods
WHERE tenant_id = $1
  AND reference_month = $2;

-- name: SumCommissionsByPeriodRange :one
SELECT
    COALESCE(SUM(total_commission), 0)::NUMERIC(15,2) as total_commission,
    COALESCE(SUM(total_net), 0)::NUMERIC(15,2) as total_net,
    COUNT(*) as periods_count
FROM commission_periods
WHERE tenant_id = $1
  AND reference_month >= $2
  AND reference_month <= $3;

-- name: GetCommissionSummaryByProfessional :many
SELECT
    professional_id,
    p.nome as professional_name,
    COUNT(*) as periods_count,
    SUM(total_gross)::NUMERIC(15,2) as sum_gross,
    SUM(total_commission)::NUMERIC(15,2) as sum_commission,
    SUM(total_advances)::NUMERIC(15,2) as sum_advances,
    SUM(total_net)::NUMERIC(15,2) as sum_net
FROM commission_periods cp
JOIN profissionais p ON cp.professional_id = p.id
WHERE cp.tenant_id = $1
  AND cp.reference_month >= $2
  AND cp.reference_month <= $3
GROUP BY professional_id, p.nome
ORDER BY sum_commission DESC;
