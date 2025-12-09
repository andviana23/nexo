-- ============================================================================
-- CATEGORIAS DE PRODUTOS QUERIES (sqlc)
-- Módulo de Estoque — NEXO v1.0
-- Tabela: categorias_produtos (customizáveis por tenant)
-- ============================================================================

-- ============================================================================
-- CREATE
-- ============================================================================

-- name: CreateCategoriaProduto :one
INSERT INTO categorias_produtos (
    id,
    tenant_id,
    nome,
    descricao,
    cor,
    icone,
    centro_custo,
    ativa,
    criado_em,
    atualizado_em
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
) RETURNING *;

-- ============================================================================
-- READ
-- ============================================================================

-- name: GetCategoriaProdutoByID :one
SELECT 
    id,
    tenant_id,
    nome,
    descricao,
    cor,
    icone,
    centro_custo,
    ativa,
    criado_em,
    atualizado_em
FROM categorias_produtos
WHERE id = $1 AND tenant_id = $2;

-- name: GetCategoriaProdutoByNome :one
SELECT 
    id,
    tenant_id,
    nome,
    descricao,
    cor,
    icone,
    centro_custo,
    ativa,
    criado_em,
    atualizado_em
FROM categorias_produtos
WHERE tenant_id = $1 AND LOWER(nome) = LOWER($2);

-- name: ListCategoriasProdutos :many
SELECT 
    id,
    tenant_id,
    nome,
    descricao,
    cor,
    icone,
    centro_custo,
    ativa,
    criado_em,
    atualizado_em
FROM categorias_produtos
WHERE tenant_id = $1
ORDER BY nome ASC;

-- name: ListCategoriasProdutosAtivas :many
SELECT 
    id,
    tenant_id,
    nome,
    descricao,
    cor,
    icone,
    centro_custo,
    ativa,
    criado_em,
    atualizado_em
FROM categorias_produtos
WHERE tenant_id = $1 AND ativa = true
ORDER BY nome ASC;

-- name: CheckCategoriaProdutoNomeExists :one
SELECT EXISTS(
    SELECT 1 
    FROM categorias_produtos 
    WHERE tenant_id = $1 
      AND LOWER(nome) = LOWER($2)
      AND id != COALESCE($3, '00000000-0000-0000-0000-000000000000'::uuid)
) AS exists;

-- name: CountProdutosByCategoria :one
SELECT COUNT(*) AS count
FROM produtos
WHERE tenant_id = $1 AND categoria_produto_id = $2;

-- ============================================================================
-- UPDATE
-- ============================================================================

-- name: UpdateCategoriaProduto :one
UPDATE categorias_produtos
SET
    nome = $3,
    descricao = $4,
    cor = $5,
    icone = $6,
    centro_custo = $7,
    ativa = $8,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: ToggleCategoriaProdutoAtiva :one
UPDATE categorias_produtos
SET
    ativa = NOT ativa,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- ============================================================================
-- DELETE
-- ============================================================================

-- name: DeleteCategoriaProduto :exec
DELETE FROM categorias_produtos
WHERE id = $1 AND tenant_id = $2;
