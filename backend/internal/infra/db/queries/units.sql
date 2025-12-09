-- ============================================================================
-- SQLC Queries: Units (Unidades)
-- ============================================================================

-- name: CreateUnit :one
INSERT INTO units (
    tenant_id,
    nome,
    apelido,
    descricao,
    endereco_resumo,
    cidade,
    estado,
    timezone,
    ativa,
    is_matriz
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetUnitByID :one
SELECT * FROM units
WHERE id = $1 AND tenant_id = $2;

-- name: GetUnitByName :one
SELECT * FROM units
WHERE tenant_id = $1 AND nome = $2;

-- name: GetMatrizUnit :one
SELECT * FROM units
WHERE tenant_id = $1 AND is_matriz = true
LIMIT 1;

-- name: ListUnitsByTenant :many
SELECT * FROM units
WHERE tenant_id = $1
ORDER BY is_matriz DESC, nome;

-- name: ListActiveUnitsByTenant :many
SELECT * FROM units
WHERE tenant_id = $1 AND ativa = true
ORDER BY is_matriz DESC, nome;

-- name: UpdateUnit :one
UPDATE units
SET
    nome = COALESCE(NULLIF($3, ''), nome),
    apelido = $4,
    descricao = $5,
    endereco_resumo = $6,
    cidade = $7,
    estado = $8,
    timezone = COALESCE(NULLIF($9, ''), timezone),
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: ToggleUnitStatus :one
UPDATE units
SET ativa = NOT ativa, atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: SetMatrizUnit :exec
UPDATE units
SET is_matriz = CASE WHEN id = $1 THEN true ELSE false END,
    atualizado_em = NOW()
WHERE tenant_id = $2;

-- name: DeleteUnit :exec
DELETE FROM units
WHERE id = $1 AND tenant_id = $2 AND is_matriz = false;

-- name: CountUnitsByTenant :one
SELECT COUNT(*) FROM units
WHERE tenant_id = $1;
