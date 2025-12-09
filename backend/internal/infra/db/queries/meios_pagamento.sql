-- ============================================================================
-- MEIOS DE PAGAMENTO QUERIES (sqlc)
-- Módulo de Cadastro de Tipos de Recebimento — NEXO v1.0
-- Tabela: meios_pagamento
-- ============================================================================

-- ============================================================================
-- CREATE
-- ============================================================================

-- name: CreateMeioPagamento :one
INSERT INTO meios_pagamento (
    id,
    tenant_id,
    nome,
    tipo,
    bandeira,
    taxa,
    taxa_fixa,
    d_mais,
    icone,
    cor,
    ordem_exibicao,
    observacoes,
    ativo,
    criado_em,
    atualizado_em
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW()
) RETURNING *;

-- ============================================================================
-- READ
-- ============================================================================

-- name: GetMeioPagamentoByID :one
SELECT 
    id,
    tenant_id,
    nome,
    tipo,
    bandeira,
    taxa,
    taxa_fixa,
    d_mais,
    icone,
    cor,
    ordem_exibicao,
    observacoes,
    ativo,
    criado_em,
    atualizado_em
FROM meios_pagamento
WHERE id = $1 AND tenant_id = $2;

-- name: ListMeiosPagamento :many
SELECT 
    id,
    tenant_id,
    nome,
    tipo,
    bandeira,
    taxa,
    taxa_fixa,
    d_mais,
    icone,
    cor,
    ordem_exibicao,
    observacoes,
    ativo,
    criado_em,
    atualizado_em
FROM meios_pagamento
WHERE tenant_id = $1
ORDER BY ordem_exibicao ASC, nome ASC;

-- name: ListMeiosPagamentoAtivos :many
SELECT 
    id,
    tenant_id,
    nome,
    tipo,
    bandeira,
    taxa,
    taxa_fixa,
    d_mais,
    icone,
    cor,
    ordem_exibicao,
    observacoes,
    ativo,
    criado_em,
    atualizado_em
FROM meios_pagamento
WHERE tenant_id = $1 AND ativo = true
ORDER BY ordem_exibicao ASC, nome ASC;

-- name: ListMeiosPagamentoPorTipo :many
SELECT 
    id,
    tenant_id,
    nome,
    tipo,
    bandeira,
    taxa,
    taxa_fixa,
    d_mais,
    icone,
    cor,
    ordem_exibicao,
    observacoes,
    ativo,
    criado_em,
    atualizado_em
FROM meios_pagamento
WHERE tenant_id = $1 AND tipo = $2
ORDER BY ordem_exibicao ASC, nome ASC;

-- name: CountMeiosPagamento :one
SELECT COUNT(*) FROM meios_pagamento
WHERE tenant_id = $1;

-- name: CountMeiosPagamentoAtivos :one
SELECT COUNT(*) FROM meios_pagamento
WHERE tenant_id = $1 AND ativo = true;

-- name: ExistsMeioPagamentoByNome :one
SELECT EXISTS(
    SELECT 1 FROM meios_pagamento
    WHERE tenant_id = $1 AND LOWER(nome) = LOWER($2)
) AS exists;

-- ============================================================================
-- UPDATE
-- ============================================================================

-- name: UpdateMeioPagamento :one
UPDATE meios_pagamento SET
    nome = $3,
    tipo = $4,
    bandeira = $5,
    taxa = $6,
    taxa_fixa = $7,
    d_mais = $8,
    icone = $9,
    cor = $10,
    ordem_exibicao = $11,
    observacoes = $12,
    ativo = $13,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: ToggleMeioPagamentoAtivo :one
UPDATE meios_pagamento SET
    ativo = NOT ativo,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: ActivateMeioPagamento :exec
UPDATE meios_pagamento SET
    ativo = true,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: DeactivateMeioPagamento :exec
UPDATE meios_pagamento SET
    ativo = false,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: UpdateOrdemExibicao :exec
UPDATE meios_pagamento SET
    ordem_exibicao = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;

-- ============================================================================
-- DELETE
-- ============================================================================

-- name: DeleteMeioPagamento :exec
DELETE FROM meios_pagamento
WHERE id = $1 AND tenant_id = $2;

-- name: DeleteMeiosPagamentoByTenant :exec
DELETE FROM meios_pagamento
WHERE tenant_id = $1;
