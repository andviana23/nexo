-- ============================================================================
-- QUERIES: commission_items
-- Itens individuais de comissão
-- ============================================================================

-- name: CreateCommissionItem :one
INSERT INTO commission_items (
    tenant_id,
    unit_id,
    professional_id,
    command_id,
    command_item_id,
    appointment_id,
    service_id,
    service_name,
    gross_value,
    commission_rate,
    commission_type,
    commission_value,
    commission_source,
    rule_id,
    reference_date,
    description,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
) RETURNING *;

-- name: GetCommissionItemByID :one
SELECT ci.*,
       p.nome as professional_name,
       s.nome as service_display_name,
       u.nome as unit_name
FROM commission_items ci
JOIN profissionais p ON ci.professional_id = p.id
LEFT JOIN servicos s ON ci.service_id = s.id
LEFT JOIN units u ON ci.unit_id = u.id
WHERE ci.id = $1 AND ci.tenant_id = $2;

-- name: GetCommissionItemByCommandItem :one
SELECT * FROM commission_items
WHERE tenant_id = $1
  AND command_item_id = $2;

-- name: ListCommissionItemsByTenant :many
SELECT ci.*,
       p.nome as professional_name,
       s.nome as service_display_name,
       u.nome as unit_name
FROM commission_items ci
JOIN profissionais p ON ci.professional_id = p.id
LEFT JOIN servicos s ON ci.service_id = s.id
LEFT JOIN units u ON ci.unit_id = u.id
WHERE ci.tenant_id = $1
ORDER BY ci.reference_date DESC, ci.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListCommissionItemsByProfessional :many
SELECT ci.*,
       s.nome as service_display_name,
       u.nome as unit_name
FROM commission_items ci
LEFT JOIN servicos s ON ci.service_id = s.id
LEFT JOIN units u ON ci.unit_id = u.id
WHERE ci.tenant_id = $1
  AND ci.professional_id = $2
ORDER BY ci.reference_date DESC, ci.created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListCommissionItemsByStatus :many
SELECT ci.*,
       p.nome as professional_name,
       s.nome as service_display_name,
       u.nome as unit_name
FROM commission_items ci
JOIN profissionais p ON ci.professional_id = p.id
LEFT JOIN servicos s ON ci.service_id = s.id
LEFT JOIN units u ON ci.unit_id = u.id
WHERE ci.tenant_id = $1
  AND ci.status = $2
ORDER BY ci.reference_date DESC, ci.created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListPendingCommissionItems :many
SELECT ci.*,
       p.nome as professional_name,
       s.nome as service_display_name,
       u.nome as unit_name
FROM commission_items ci
JOIN profissionais p ON ci.professional_id = p.id
LEFT JOIN servicos s ON ci.service_id = s.id
LEFT JOIN units u ON ci.unit_id = u.id
WHERE ci.tenant_id = $1
  AND ci.status = 'PENDENTE'
ORDER BY ci.reference_date ASC, ci.created_at ASC;

-- name: ListPendingCommissionItemsByProfessional :many
SELECT ci.*,
       s.nome as service_display_name
FROM commission_items ci
LEFT JOIN servicos s ON ci.service_id = s.id
WHERE ci.tenant_id = $1
  AND ci.professional_id = $2
  AND ci.status = 'PENDENTE'
ORDER BY ci.reference_date ASC, ci.created_at ASC;

-- name: ListCommissionItemsByPeriod :many
SELECT ci.*,
       p.nome as professional_name,
       s.nome as service_display_name
FROM commission_items ci
JOIN profissionais p ON ci.professional_id = p.id
LEFT JOIN servicos s ON ci.service_id = s.id
WHERE ci.tenant_id = $1
  AND ci.period_id = $2
ORDER BY p.nome ASC, ci.reference_date ASC;

-- name: ListCommissionItemsByDateRange :many
SELECT ci.*,
       p.nome as professional_name,
       s.nome as service_display_name,
       u.nome as unit_name
FROM commission_items ci
JOIN profissionais p ON ci.professional_id = p.id
LEFT JOIN servicos s ON ci.service_id = s.id
LEFT JOIN units u ON ci.unit_id = u.id
WHERE ci.tenant_id = $1
  AND ci.reference_date >= $2
  AND ci.reference_date <= $3
ORDER BY ci.reference_date ASC, p.nome ASC;

-- name: ListCommissionItemsByProfessionalAndDateRange :many
SELECT ci.*,
       s.nome as service_display_name
FROM commission_items ci
LEFT JOIN servicos s ON ci.service_id = s.id
WHERE ci.tenant_id = $1
  AND ci.professional_id = $2
  AND ci.reference_date >= $3
  AND ci.reference_date <= $4
ORDER BY ci.reference_date ASC, ci.created_at ASC;

-- name: UpdateCommissionItem :one
UPDATE commission_items
SET
    commission_rate = $3,
    commission_type = $4,
    commission_value = $5,
    commission_source = $6,
    description = $7,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: ProcessCommissionItem :one
UPDATE commission_items
SET
    status = 'PROCESSADO',
    period_id = $3,
    processed_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status = 'PENDENTE'
RETURNING *;

-- name: MarkCommissionItemAsPaid :one
UPDATE commission_items
SET
    status = 'PAGO',
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status = 'PROCESSADO'
RETURNING *;

-- name: CancelCommissionItem :one
UPDATE commission_items
SET
    status = 'CANCELADO',
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status = 'PENDENTE'
RETURNING *;

-- name: ReverseCommissionItem :one
UPDATE commission_items
SET
    status = 'ESTORNADO',
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status IN ('PENDENTE', 'PROCESSADO')
RETURNING *;

-- name: BulkProcessCommissionItems :exec
UPDATE commission_items
SET
    status = 'PROCESSADO',
    period_id = $3,
    processed_at = NOW(),
    updated_at = NOW()
WHERE tenant_id = $1
  AND professional_id = $2
  AND status = 'PENDENTE'
  AND reference_date >= $4
  AND reference_date <= $5;

-- name: DeleteCommissionItem :exec
DELETE FROM commission_items
WHERE id = $1 AND tenant_id = $2 AND status = 'PENDENTE';

-- name: SumPendingCommissionsByProfessional :one
SELECT
    COALESCE(SUM(gross_value), 0)::NUMERIC(15,2) as total_gross,
    COALESCE(SUM(commission_value), 0)::NUMERIC(15,2) as total_commission,
    COUNT(*) as items_count
FROM commission_items
WHERE tenant_id = $1
  AND professional_id = $2
  AND status = 'PENDENTE';

-- name: SumCommissionsByProfessionalAndDateRange :one
SELECT
    COALESCE(SUM(gross_value), 0)::NUMERIC(15,2) as total_gross,
    COALESCE(SUM(commission_value), 0)::NUMERIC(15,2) as total_commission,
    COUNT(*) as items_count
FROM commission_items
WHERE tenant_id = $1
  AND professional_id = $2
  AND reference_date >= $3
  AND reference_date <= $4
  AND status != 'CANCELADO'
  AND status != 'ESTORNADO';

-- name: SumCommissionsByDateRange :one
SELECT
    COALESCE(SUM(gross_value), 0)::NUMERIC(15,2) as total_gross,
    COALESCE(SUM(commission_value), 0)::NUMERIC(15,2) as total_commission,
    COUNT(*) as items_count
FROM commission_items
WHERE tenant_id = $1
  AND reference_date >= $2
  AND reference_date <= $3
  AND status != 'CANCELADO'
  AND status != 'ESTORNADO';

-- name: GetCommissionSummaryByService :many
SELECT
    service_id,
    service_name,
    COUNT(*) as items_count,
    SUM(gross_value)::NUMERIC(15,2) as total_gross,
    SUM(commission_value)::NUMERIC(15,2) as total_commission,
    AVG(commission_rate)::NUMERIC(5,2) as avg_rate
FROM commission_items
WHERE tenant_id = $1
  AND reference_date >= $2
  AND reference_date <= $3
  AND status != 'CANCELADO'
  AND status != 'ESTORNADO'
GROUP BY service_id, service_name
ORDER BY total_commission DESC;

-- name: GetCommissionSummaryByProfessionalAndMonth :one
SELECT
    professional_id,
    COALESCE(SUM(gross_value), 0)::NUMERIC(15,2) as total_gross,
    COALESCE(SUM(commission_value), 0)::NUMERIC(15,2) as total_commission,
    COUNT(*) as items_count,
    COUNT(*) FILTER (WHERE status = 'PENDENTE') as pending_count,
    COUNT(*) FILTER (WHERE status = 'PROCESSADO') as processed_count,
    COUNT(*) FILTER (WHERE status = 'PAGO') as paid_count
FROM commission_items
WHERE tenant_id = $1
  AND professional_id = $2
  AND to_char(reference_date, 'YYYY-MM') = $3
  AND status != 'CANCELADO'
  AND status != 'ESTORNADO'
GROUP BY professional_id;

-- name: CountCommissionItemsByStatus :one
SELECT
    COUNT(*) FILTER (WHERE status = 'PENDENTE') as pending,
    COUNT(*) FILTER (WHERE status = 'PROCESSADO') as processed,
    COUNT(*) FILTER (WHERE status = 'PAGO') as paid,
    COUNT(*) FILTER (WHERE status = 'CANCELADO') as cancelled,
    COUNT(*) FILTER (WHERE status = 'ESTORNADO') as reversed
FROM commission_items
WHERE tenant_id = $1;
