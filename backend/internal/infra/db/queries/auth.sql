-- name: GetUserByEmail :one
SELECT id, tenant_id, nome, email, password_hash, role, ativo, ultimo_login, criado_em, atualizado_em
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT id, tenant_id, nome, email, password_hash, role, ativo, ultimo_login, criado_em, atualizado_em
FROM users
WHERE id = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (tenant_id, nome, email, password_hash, role, ativo)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, tenant_id, nome, email, password_hash, role, ativo, ultimo_login, criado_em, atualizado_em;

-- name: UpdateLastLogin :exec
UPDATE users
SET ultimo_login = NOW(), atualizado_em = NOW()
WHERE id = $1;

-- name: SaveRefreshToken :exec
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3);

-- name: GetRefreshToken :one
SELECT id, user_id, token, expires_at, created_at
FROM refresh_tokens
WHERE token = $1 AND expires_at > NOW()
LIMIT 1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW();
