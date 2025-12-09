-- ============================================================================
-- BLOCKED_TIMES QUERIES (sqlc)
-- Bloqueios de horário na agenda
-- ============================================================================

-- name: CreateBlockedTime :one
INSERT INTO blocked_times (
    id,
    tenant_id,
    professional_id,
    start_time,
    end_time,
    reason,
    is_recurring,
    recurrence_rule,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetBlockedTimeByID :one
SELECT * FROM blocked_times
WHERE id = $1 AND tenant_id = $2;

-- name: ListBlockedTimes :many
-- Lista bloqueios com filtros opcionais
SELECT * FROM blocked_times
WHERE tenant_id = sqlc.arg(tenant_id)::uuid
  AND (sqlc.narg(professional_id)::uuid IS NULL OR professional_id = sqlc.narg(professional_id)::uuid)
  AND (sqlc.narg(start_date)::timestamptz IS NULL OR start_time >= sqlc.narg(start_date)::timestamptz)
  AND (sqlc.narg(end_date)::timestamptz IS NULL OR end_time <= sqlc.narg(end_date)::timestamptz)
ORDER BY start_time ASC;

-- name: CheckBlockedTimeConflict :one
-- Verifica se há conflito com bloqueios existentes
SELECT EXISTS (
    SELECT 1 FROM blocked_times
    WHERE tenant_id = sqlc.arg(tenant_id)::uuid
      AND professional_id = sqlc.arg(professional_id)::uuid
      AND (sqlc.narg(exclude_id)::uuid IS NULL OR id != sqlc.narg(exclude_id)::uuid)
      AND start_time < sqlc.arg(end_time)::timestamptz
      AND end_time > sqlc.arg(start_time)::timestamptz
) as has_conflict;

-- name: GetBlockedTimesInRange :many
-- Busca bloqueios em um intervalo de tempo (para validação de agendamentos)
SELECT * FROM blocked_times
WHERE tenant_id = $1
  AND professional_id = $2
  AND start_time < $4
  AND end_time > $3
ORDER BY start_time ASC;

-- name: UpdateBlockedTime :one
UPDATE blocked_times
SET
    start_time = $3,
    end_time = $4,
    reason = $5,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteBlockedTime :exec
DELETE FROM blocked_times
WHERE id = $1 AND tenant_id = $2;
