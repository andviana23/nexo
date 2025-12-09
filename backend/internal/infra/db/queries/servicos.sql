-- ============================================================================
-- SERVIÇOS QUERIES (sqlc)
-- Módulo de Cadastro de Serviços — NEXO v1.0
-- Tabela: servicos (vinculada a categorias_servicos)
-- ============================================================================

-- ============================================================================
-- CREATE
-- ============================================================================

-- name: CreateServico :one
INSERT INTO servicos (
    id,
    tenant_id,
    categoria_id,
    nome,
    descricao,
    preco,
    duracao,
    comissao,
    cor,
    imagem,
    profissionais_ids,
    observacoes,
    tags,
    ativo,
    criado_em,
    atualizado_em
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW()
) RETURNING *;

-- ============================================================================
-- READ
-- ============================================================================

-- name: GetServicoByID :one
SELECT 
    s.id,
    s.tenant_id,
    s.categoria_id,
    s.nome,
    s.descricao,
    s.preco,
    s.duracao,
    s.comissao,
    s.cor,
    s.imagem,
    s.profissionais_ids,
    s.observacoes,
    s.tags,
    s.ativo,
    s.criado_em,
    s.atualizado_em,
    cs.nome AS categoria_nome,
    cs.cor AS categoria_cor
FROM servicos s
LEFT JOIN categorias_servicos cs ON cs.id = s.categoria_id AND cs.tenant_id = s.tenant_id
WHERE s.id = $1 AND s.tenant_id = $2;

-- name: ListServicos :many
SELECT 
    s.id,
    s.tenant_id,
    s.categoria_id,
    s.nome,
    s.descricao,
    s.preco,
    s.duracao,
    s.comissao,
    s.cor,
    s.imagem,
    s.profissionais_ids,
    s.observacoes,
    s.tags,
    s.ativo,
    s.criado_em,
    s.atualizado_em,
    cs.nome AS categoria_nome,
    cs.cor AS categoria_cor
FROM servicos s
LEFT JOIN categorias_servicos cs ON cs.id = s.categoria_id AND cs.tenant_id = s.tenant_id
WHERE s.tenant_id = $1
ORDER BY s.nome ASC;

-- name: ListServicosAtivos :many
SELECT 
    s.id,
    s.tenant_id,
    s.categoria_id,
    s.nome,
    s.descricao,
    s.preco,
    s.duracao,
    s.comissao,
    s.cor,
    s.imagem,
    s.profissionais_ids,
    s.observacoes,
    s.tags,
    s.ativo,
    s.criado_em,
    s.atualizado_em,
    cs.nome AS categoria_nome,
    cs.cor AS categoria_cor
FROM servicos s
LEFT JOIN categorias_servicos cs ON cs.id = s.categoria_id AND cs.tenant_id = s.tenant_id
WHERE s.tenant_id = $1 AND s.ativo = true
ORDER BY s.nome ASC;

-- name: ListServicosByCategoria :many
SELECT 
    s.id,
    s.tenant_id,
    s.categoria_id,
    s.nome,
    s.descricao,
    s.preco,
    s.duracao,
    s.comissao,
    s.cor,
    s.imagem,
    s.profissionais_ids,
    s.observacoes,
    s.tags,
    s.ativo,
    s.criado_em,
    s.atualizado_em,
    cs.nome AS categoria_nome,
    cs.cor AS categoria_cor
FROM servicos s
LEFT JOIN categorias_servicos cs ON cs.id = s.categoria_id AND cs.tenant_id = s.tenant_id
WHERE s.tenant_id = $1 AND s.categoria_id = $2
ORDER BY s.nome ASC;

-- name: ListServicosByProfissional :many
SELECT 
    s.id,
    s.tenant_id,
    s.categoria_id,
    s.nome,
    s.descricao,
    s.preco,
    s.duracao,
    s.comissao,
    s.cor,
    s.imagem,
    s.profissionais_ids,
    s.observacoes,
    s.tags,
    s.ativo,
    s.criado_em,
    s.atualizado_em,
    cs.nome AS categoria_nome,
    cs.cor AS categoria_cor
FROM servicos s
LEFT JOIN categorias_servicos cs ON cs.id = s.categoria_id AND cs.tenant_id = s.tenant_id
WHERE s.tenant_id = $1 
  AND s.ativo = true
  AND @profissional_id::uuid = ANY(s.profissionais_ids)
ORDER BY s.nome ASC;

-- name: CheckServicoNomeExists :one
SELECT EXISTS(
    SELECT 1 
    FROM servicos 
    WHERE tenant_id = $1 
      AND LOWER(nome) = LOWER($2)
      AND id != COALESCE($3, '00000000-0000-0000-0000-000000000000'::uuid)
) AS exists;

-- name: SearchServicos :many
SELECT 
    s.id,
    s.tenant_id,
    s.categoria_id,
    s.nome,
    s.descricao,
    s.preco,
    s.duracao,
    s.comissao,
    s.cor,
    s.imagem,
    s.profissionais_ids,
    s.observacoes,
    s.tags,
    s.ativo,
    s.criado_em,
    s.atualizado_em,
    cs.nome AS categoria_nome,
    cs.cor AS categoria_cor
FROM servicos s
LEFT JOIN categorias_servicos cs ON cs.id = s.categoria_id AND cs.tenant_id = s.tenant_id
WHERE s.tenant_id = $1 
  AND (
      LOWER(s.nome) LIKE LOWER('%' || @search_term::text || '%')
      OR LOWER(s.descricao) LIKE LOWER('%' || @search_term::text || '%')
      OR @search_term::text = ANY(s.tags)
  )
ORDER BY s.nome ASC;

-- ============================================================================
-- UPDATE
-- ============================================================================

-- name: UpdateServico :one
UPDATE servicos SET
    categoria_id = $3,
    nome = $4,
    descricao = $5,
    preco = $6,
    duracao = $7,
    comissao = $8,
    cor = $9,
    imagem = $10,
    profissionais_ids = $11,
    observacoes = $12,
    tags = $13,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: ToggleServicoStatus :one
UPDATE servicos SET
    ativo = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: UpdateServicoCategoria :one
UPDATE servicos SET
    categoria_id = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: UpdateServicoProfissionais :one
UPDATE servicos SET
    profissionais_ids = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- ============================================================================
-- DELETE
-- ============================================================================

-- name: DeleteServico :exec
DELETE FROM servicos
WHERE id = $1 AND tenant_id = $2;

-- name: DeleteServicosByCategoria :exec
DELETE FROM servicos
WHERE categoria_id = $1 AND tenant_id = $2;

-- ============================================================================
-- QUERIES AUXILIARES
-- ============================================================================

-- name: CountServicosByTenant :one
SELECT COUNT(*) AS total
FROM servicos
WHERE tenant_id = $1;

-- name: CountServicosAtivosByTenant :one
SELECT COUNT(*) AS total
FROM servicos
WHERE tenant_id = $1 AND ativo = true;

-- name: GetServicosStats :one
SELECT 
    COUNT(*) AS total_servicos,
    COUNT(*) FILTER (WHERE ativo = true) AS servicos_ativos,
    COUNT(*) FILTER (WHERE ativo = false) AS servicos_inativos,
    COALESCE(AVG(preco), 0) AS preco_medio,
    COALESCE(AVG(duracao), 0) AS duracao_media,
    COALESCE(AVG(comissao), 0) AS comissao_media
FROM servicos
WHERE tenant_id = $1;

-- name: ListServicosComCategoria :many
SELECT 
    s.id,
    s.tenant_id,
    s.categoria_id,
    s.nome,
    s.descricao,
    s.preco,
    s.duracao,
    s.comissao,
    s.cor,
    s.imagem,
    s.profissionais_ids,
    s.observacoes,
    s.tags,
    s.ativo,
    s.criado_em,
    s.atualizado_em,
    cs.nome AS categoria_nome,
    cs.cor AS categoria_cor,
    cs.icone AS categoria_icone
FROM servicos s
INNER JOIN categorias_servicos cs ON cs.id = s.categoria_id AND cs.tenant_id = s.tenant_id
WHERE s.tenant_id = $1
ORDER BY cs.nome ASC, s.nome ASC;

-- name: ListServicosSemCategoria :many
SELECT 
    s.id,
    s.tenant_id,
    s.categoria_id,
    s.nome,
    s.descricao,
    s.preco,
    s.duracao,
    s.comissao,
    s.cor,
    s.imagem,
    s.profissionais_ids,
    s.observacoes,
    s.tags,
    s.ativo,
    s.criado_em,
    s.atualizado_em
FROM servicos s
WHERE s.tenant_id = $1 AND s.categoria_id IS NULL
ORDER BY s.nome ASC;

-- name: GetServicosByIDs :many
SELECT 
    s.id,
    s.tenant_id,
    s.categoria_id,
    s.nome,
    s.descricao,
    s.preco,
    s.duracao,
    s.comissao,
    s.cor,
    s.imagem,
    s.profissionais_ids,
    s.observacoes,
    s.tags,
    s.ativo,
    s.criado_em,
    s.atualizado_em,
    cs.nome AS categoria_nome,
    cs.cor AS categoria_cor
FROM servicos s
LEFT JOIN categorias_servicos cs ON cs.id = s.categoria_id AND cs.tenant_id = s.tenant_id
WHERE s.tenant_id = $1 AND s.id = ANY(@servico_ids::uuid[])
ORDER BY s.nome ASC;
