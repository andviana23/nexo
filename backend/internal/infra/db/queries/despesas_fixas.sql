-- ============================================================================
-- Queries sqlc: despesas_fixas
-- Módulo: Financeiro
-- Sprint: 2
-- ============================================================================

-- name: CreateDespesaFixa :one
-- Cria uma nova despesa fixa
INSERT INTO despesas_fixas (
    tenant_id,
    unidade_id,
    descricao,
    categoria_id,
    fornecedor,
    valor,
    dia_vencimento,
    ativo,
    observacoes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetDespesaFixaByID :one
-- Busca despesa fixa por ID (com isolamento de tenant)
SELECT * FROM despesas_fixas
WHERE id = $1 AND tenant_id = $2;

-- name: ListDespesasFixasByTenant :many
-- Lista todas as despesas fixas do tenant com paginação
SELECT * FROM despesas_fixas
WHERE tenant_id = $1
ORDER BY descricao ASC
LIMIT $2 OFFSET $3;

-- name: ListDespesasFixasAtivas :many
-- Lista apenas despesas fixas ativas (usado pelo cron job)
SELECT * FROM despesas_fixas
WHERE tenant_id = $1 AND ativo = true
ORDER BY dia_vencimento ASC;

-- name: ListDespesasFixasByUnidade :many
-- Lista despesas fixas de uma unidade específica
SELECT * FROM despesas_fixas
WHERE tenant_id = $1 AND unidade_id = $2
ORDER BY descricao ASC;

-- name: ListDespesasFixasByCategoria :many
-- Lista despesas fixas de uma categoria específica
SELECT * FROM despesas_fixas
WHERE tenant_id = $1 AND categoria_id = $2
ORDER BY descricao ASC;

-- name: UpdateDespesaFixa :one
-- Atualiza uma despesa fixa
UPDATE despesas_fixas
SET
    descricao = $3,
    categoria_id = $4,
    fornecedor = $5,
    valor = $6,
    dia_vencimento = $7,
    unidade_id = $8,
    observacoes = $9,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: ToggleDespesaFixa :one
-- Alterna o status ativo/inativo
UPDATE despesas_fixas
SET
    ativo = NOT ativo,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: ActivateDespesaFixa :one
-- Ativa uma despesa fixa
UPDATE despesas_fixas
SET
    ativo = true,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeactivateDespesaFixa :one
-- Desativa uma despesa fixa
UPDATE despesas_fixas
SET
    ativo = false,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteDespesaFixa :exec
-- Remove uma despesa fixa
DELETE FROM despesas_fixas
WHERE id = $1 AND tenant_id = $2;

-- name: SumDespesasFixasAtivas :one
-- Soma o valor total de despesas fixas ativas
SELECT COALESCE(SUM(valor), 0)::DECIMAL(15,2) as total
FROM despesas_fixas
WHERE tenant_id = $1 AND ativo = true;

-- name: SumDespesasFixasByUnidade :one
-- Soma o valor total de despesas fixas ativas por unidade
SELECT COALESCE(SUM(valor), 0)::DECIMAL(15,2) as total
FROM despesas_fixas
WHERE tenant_id = $1 AND unidade_id = $2 AND ativo = true;

-- name: CountDespesasFixas :one
-- Conta total de despesas fixas do tenant
SELECT COUNT(*) FROM despesas_fixas
WHERE tenant_id = $1;

-- name: CountDespesasFixasAtivas :one
-- Conta despesas fixas ativas do tenant
SELECT COUNT(*) FROM despesas_fixas
WHERE tenant_id = $1 AND ativo = true;

-- name: ListDespesasFixasAtivasPorTenants :many
-- Lista todas as despesas fixas ativas de todos os tenants
-- Usado pelo cron job para geração em massa
SELECT df.*, t.nome as tenant_nome
FROM despesas_fixas df
JOIN tenants t ON df.tenant_id = t.id
WHERE df.ativo = true AND t.ativo = true
ORDER BY df.tenant_id, df.dia_vencimento;

-- name: ExistsDespesaFixaByDescricao :one
-- Verifica se já existe despesa fixa com mesma descrição no tenant
SELECT EXISTS(
    SELECT 1 FROM despesas_fixas
    WHERE tenant_id = $1 
      AND LOWER(descricao) = LOWER($2)
      AND id != COALESCE($3, '00000000-0000-0000-0000-000000000000'::uuid)
) as exists;
