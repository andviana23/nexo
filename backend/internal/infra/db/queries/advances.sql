-- ============================================================================
-- QUERIES: advances
-- Adiantamentos de profissionais
-- ============================================================================

-- name: CreateAdvance :one
INSERT INTO advances (
    tenant_id,
    unit_id,
    professional_id,
    amount,
    request_date,
    reason,
    status,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetAdvanceByID :one
SELECT a.*,
       p.nome as professional_name,
       u.nome as unit_name
FROM advances a
JOIN profissionais p ON a.professional_id = p.id
LEFT JOIN units u ON a.unit_id = u.id
WHERE a.id = $1 AND a.tenant_id = $2;

-- name: ListAdvancesByTenant :many
SELECT a.*,
       p.nome as professional_name,
       u.nome as unit_name
FROM advances a
JOIN profissionais p ON a.professional_id = p.id
LEFT JOIN units u ON a.unit_id = u.id
WHERE a.tenant_id = $1
ORDER BY a.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAdvancesByProfessional :many
SELECT a.*,
       u.nome as unit_name
FROM advances a
LEFT JOIN units u ON a.unit_id = u.id
WHERE a.tenant_id = $1
  AND a.professional_id = $2
ORDER BY a.created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListAdvancesByStatus :many
SELECT a.*,
       p.nome as professional_name,
       u.nome as unit_name
FROM advances a
JOIN profissionais p ON a.professional_id = p.id
LEFT JOIN units u ON a.unit_id = u.id
WHERE a.tenant_id = $1
  AND a.status = $2
ORDER BY a.created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListPendingAdvances :many
SELECT a.*,
       p.nome as professional_name,
       u.nome as unit_name
FROM advances a
JOIN profissionais p ON a.professional_id = p.id
LEFT JOIN units u ON a.unit_id = u.id
WHERE a.tenant_id = $1
  AND a.status = 'PENDING'
ORDER BY a.request_date ASC, a.created_at ASC;

-- name: ListApprovedAdvancesNotDeducted :many
SELECT a.*,
       p.nome as professional_name,
       u.nome as unit_name
FROM advances a
JOIN profissionais p ON a.professional_id = p.id
LEFT JOIN units u ON a.unit_id = u.id
WHERE a.tenant_id = $1
  AND a.status = 'APPROVED'
ORDER BY a.request_date ASC;

-- name: ListApprovedAdvancesForProfessional :many
SELECT a.*
FROM advances a
WHERE a.tenant_id = $1
  AND a.professional_id = $2
  AND a.status = 'APPROVED'
ORDER BY a.request_date ASC;

-- name: ApproveAdvance :one
UPDATE advances
SET
    status = 'APPROVED',
    approved_at = NOW(),
    approved_by = $3,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status = 'PENDING'
RETURNING *;

-- name: RejectAdvance :one
UPDATE advances
SET
    status = 'REJECTED',
    rejected_at = NOW(),
    rejected_by = $3,
    rejection_reason = $4,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status = 'PENDING'
RETURNING *;

-- name: DeductAdvance :one
UPDATE advances
SET
    status = 'DEDUCTED',
    deducted_at = NOW(),
    deduction_period_id = $3,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status = 'APPROVED'
RETURNING *;

-- name: CancelAdvance :one
UPDATE advances
SET
    status = 'CANCELLED',
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status IN ('PENDING', 'APPROVED')
RETURNING *;

-- name: DeleteAdvance :exec
DELETE FROM advances
WHERE id = $1 AND tenant_id = $2 AND status = 'PENDING';

-- name: SumPendingAdvancesByProfessional :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC(15,2) as total
FROM advances
WHERE tenant_id = $1
  AND professional_id = $2
  AND status = 'PENDING';

-- name: SumApprovedAdvancesByProfessional :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC(15,2) as total
FROM advances
WHERE tenant_id = $1
  AND professional_id = $2
  AND status = 'APPROVED';

-- name: SumAdvancesByPeriod :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC(15,2) as total
FROM advances
WHERE tenant_id = $1
  AND request_date >= $2
  AND request_date <= $3
  AND status IN ('APPROVED', 'DEDUCTED');

-- name: CountAdvancesByStatus :one
SELECT
    COUNT(*) FILTER (WHERE status = 'PENDING') as pending,
    COUNT(*) FILTER (WHERE status = 'APPROVED') as approved,
    COUNT(*) FILTER (WHERE status = 'REJECTED') as rejected,
    COUNT(*) FILTER (WHERE status = 'DEDUCTED') as deducted
FROM advances
WHERE tenant_id = $1;
