-- ============================================================================
-- SQLC Queries: User Units (Vínculo Usuário-Unidade)
-- ============================================================================

-- name: CreateUserUnit :one
INSERT INTO user_units (
    user_id,
    unit_id,
    is_default,
    role_override
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserUnit :one
SELECT * FROM user_units
WHERE user_id = $1 AND unit_id = $2;

-- name: GetUserDefaultUnit :one
SELECT uu.*, u.nome as unit_nome, u.apelido as unit_apelido, u.tenant_id
FROM user_units uu
JOIN units u ON u.id = uu.unit_id
WHERE uu.user_id = $1 AND uu.is_default = true
LIMIT 1;

-- name: ListUserUnits :many
SELECT 
    uu.*,
    u.nome as unit_nome,
    u.apelido as unit_apelido,
    u.is_matriz,
    u.ativa as unit_ativa,
    u.tenant_id
FROM user_units uu
JOIN units u ON u.id = uu.unit_id
WHERE uu.user_id = $1 AND u.ativa = true
ORDER BY uu.is_default DESC, u.is_matriz DESC, u.nome;

-- name: ListUnitUsers :many
SELECT 
    uu.*,
    usr.nome as user_nome,
    usr.email as user_email,
    usr.role as user_role
FROM user_units uu
JOIN users usr ON usr.id = uu.user_id
WHERE uu.unit_id = $1
ORDER BY usr.nome;

-- name: CheckUserUnitAccess :one
SELECT EXISTS(
    SELECT 1 FROM user_units uu
    JOIN units u ON u.id = uu.unit_id
    WHERE uu.user_id = $1 
      AND uu.unit_id = $2
      AND u.ativa = true
) as has_access;

-- name: SetUserDefaultUnit :exec
UPDATE user_units
SET is_default = CASE WHEN unit_id = $2 THEN true ELSE false END,
    atualizado_em = NOW()
WHERE user_id = $1;

-- name: UpdateUserUnitRole :one
UPDATE user_units
SET role_override = $3, atualizado_em = NOW()
WHERE user_id = $1 AND unit_id = $2
RETURNING *;

-- name: DeleteUserUnit :exec
DELETE FROM user_units
WHERE user_id = $1 AND unit_id = $2;

-- name: DeleteAllUserUnits :exec
DELETE FROM user_units
WHERE user_id = $1;

-- name: DeleteAllUnitUsers :exec
DELETE FROM user_units
WHERE unit_id = $1;

-- name: CountUserUnits :one
SELECT COUNT(*) FROM user_units
WHERE user_id = $1;

-- name: CountUnitUsers :one
SELECT COUNT(*) FROM user_units
WHERE unit_id = $1;

-- name: BulkCreateUserUnits :copyfrom
INSERT INTO user_units (
    user_id,
    unit_id,
    is_default,
    role_override
) VALUES (
    $1, $2, $3, $4
);
