-- ============================================================================
-- CATEGORIAS DE SERVIÇOS QUERIES (sqlc)
-- Módulo de Cadastro de Serviços — NEXO v1.0
-- Tabela: categorias_servicos (separada de categorias financeiras)
-- ============================================================================

-- ============================================================================
-- CREATE
-- ============================================================================

-- name: CreateCategoriaServico :one
INSERT INTO categorias_servicos (
    id,
    tenant_id,
    nome,
    descricao,
    cor,
    icone,
    ativa,
    criado_em,
    atualizado_em
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
) RETURNING *;

-- ============================================================================
-- READ
-- ============================================================================

-- name: GetCategoriaServicoByID :one
SELECT 
    id,
    tenant_id,
    nome,
    descricao,
    cor,
    icone,
    ativa,
    criado_em,
    atualizado_em
FROM categorias_servicos
WHERE id = $1 AND tenant_id = $2;

-- name: ListCategoriasServicos :many
SELECT 
    id,
    tenant_id,
    nome,
    descricao,
    cor,
    icone,
    ativa,
    criado_em,
    atualizado_em
FROM categorias_servicos
WHERE tenant_id = $1
ORDER BY nome ASC;

-- name: ListCategoriasServicosAtivas :many
SELECT 
    id,
    tenant_id,
    nome,
    descricao,
    cor,
    icone,
    ativa,
    criado_em,
    atualizado_em
FROM categorias_servicos
WHERE tenant_id = $1 AND ativa = true
ORDER BY nome ASC;

-- name: CheckCategoriaServicoNomeExists :one
SELECT EXISTS(
    SELECT 1 
    FROM categorias_servicos 
    WHERE tenant_id = $1 
      AND LOWER(nome) = LOWER($2)
      AND id != COALESCE($3, '00000000-0000-0000-0000-000000000000'::uuid)
) AS exists;

-- ============================================================================
-- UPDATE
-- ============================================================================

-- name: UpdateCategoriaServico :one
UPDATE categorias_servicos SET
    nome = $3,
    descricao = $4,
    cor = $5,
    icone = $6,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: ToggleCategoriaServicoStatus :one
UPDATE categorias_servicos SET
    ativa = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- ============================================================================
-- DELETE
-- ============================================================================

-- name: DeleteCategoriaServico :exec
DELETE FROM categorias_servicos
WHERE id = $1 AND tenant_id = $2;

-- ============================================================================
-- QUERIES AUXILIARES
-- ============================================================================

-- name: CountServicosInCategoria :one
SELECT COUNT(*) AS total
FROM servicos
WHERE categoria_id = $1 AND tenant_id = $2;

-- name: CountCategoriasServicosByTenant :one
SELECT COUNT(*) AS total
FROM categorias_servicos
WHERE tenant_id = $1;

-- name: GetCategoriasServicosComServicos :many
SELECT 
    cs.id,
    cs.tenant_id,
    cs.nome,
    cs.descricao,
    cs.cor,
    cs.icone,
    cs.ativa,
    cs.criado_em,
    cs.atualizado_em,
    COUNT(s.id) AS total_servicos
FROM categorias_servicos cs
LEFT JOIN servicos s ON s.categoria_id = cs.id AND s.tenant_id = cs.tenant_id
WHERE cs.tenant_id = $1
GROUP BY cs.id, cs.tenant_id, cs.nome, cs.descricao, cs.cor, cs.icone, cs.ativa, cs.criado_em, cs.atualizado_em
ORDER BY cs.nome ASC;
