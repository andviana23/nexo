-- =============================================
-- CAIXA DIÁRIO - Queries SQLC
-- Módulo de controle operacional da gaveta
-- =============================================

-- ========== CREATE ==========

-- name: CreateCaixaDiario :one
INSERT INTO caixa_diario (
    id,
    tenant_id,
    usuario_abertura_id,
    data_abertura,
    saldo_inicial,
    total_entradas,
    total_saidas,
    total_sangrias,
    total_reforcos,
    saldo_esperado,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- ========== READ ==========

-- name: GetCaixaDiarioByID :one
SELECT 
    c.*,
    ua.nome as usuario_abertura_nome,
    COALESCE(uf.nome, '') as usuario_fechamento_nome
FROM caixa_diario c
LEFT JOIN users ua ON ua.id = c.usuario_abertura_id
LEFT JOIN users uf ON uf.id = c.usuario_fechamento_id
WHERE c.id = $1 AND c.tenant_id = $2;

-- name: GetCaixaDiarioAberto :one
SELECT 
    c.*,
    ua.nome as usuario_abertura_nome,
    '' as usuario_fechamento_nome
FROM caixa_diario c
LEFT JOIN users ua ON ua.id = c.usuario_abertura_id
WHERE c.tenant_id = $1 AND c.status = 'ABERTO'
LIMIT 1;

-- name: ExistsCaixaAberto :one
SELECT EXISTS(
    SELECT 1 FROM caixa_diario 
    WHERE tenant_id = $1 AND status = 'ABERTO'
) as exists;

-- ========== UPDATE ==========

-- name: UpdateCaixaDiario :one
UPDATE caixa_diario
SET
    total_entradas = $3,
    total_saidas = $4,
    total_sangrias = $5,
    total_reforcos = $6,
    saldo_esperado = $7,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: UpdateCaixaDiarioTotais :exec
UPDATE caixa_diario
SET
    total_sangrias = $3,
    total_reforcos = $4,
    total_entradas = $5,
    saldo_esperado = saldo_inicial + $5 - $3 + $4,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: FecharCaixaDiario :one
UPDATE caixa_diario
SET
    usuario_fechamento_id = $3,
    data_fechamento = $4,
    saldo_real = $5,
    divergencia = $6,
    status = 'FECHADO',
    justificativa_divergencia = $7,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status = 'ABERTO'
RETURNING *;

-- ========== LIST ==========

-- name: ListCaixaDiarioHistorico :many
SELECT 
    c.*,
    ua.nome as usuario_abertura_nome,
    COALESCE(uf.nome, '') as usuario_fechamento_nome
FROM caixa_diario c
LEFT JOIN users ua ON ua.id = c.usuario_abertura_id
LEFT JOIN users uf ON uf.id = c.usuario_fechamento_id
WHERE c.tenant_id = $1 
    AND c.status = 'FECHADO'
    AND ($2::date IS NULL OR c.data_abertura >= $2)
    AND ($3::date IS NULL OR c.data_abertura <= $3)
    AND ($4::uuid IS NULL OR c.usuario_abertura_id = $4 OR c.usuario_fechamento_id = $4)
ORDER BY c.data_abertura DESC
LIMIT $5 OFFSET $6;

-- name: CountCaixaDiarioHistorico :one
SELECT COUNT(*) 
FROM caixa_diario
WHERE tenant_id = $1 
    AND status = 'FECHADO'
    AND ($2::date IS NULL OR data_abertura >= $2)
    AND ($3::date IS NULL OR data_abertura <= $3)
    AND ($4::uuid IS NULL OR usuario_abertura_id = $4 OR usuario_fechamento_id = $4);

-- =============================================
-- OPERAÇÕES DO CAIXA
-- =============================================

-- name: CreateOperacaoCaixa :one
INSERT INTO operacoes_caixa (
    id,
    caixa_id,
    tenant_id,
    tipo,
    valor,
    descricao,
    destino,
    origem,
    usuario_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: ListOperacoesByCaixa :many
SELECT 
    o.*,
    u.nome as usuario_nome
FROM operacoes_caixa o
LEFT JOIN users u ON u.id = o.usuario_id
WHERE o.caixa_id = $1 AND o.tenant_id = $2
ORDER BY o.created_at ASC;

-- name: ListOperacoesByCaixaAndTipo :many
SELECT 
    o.*,
    u.nome as usuario_nome
FROM operacoes_caixa o
LEFT JOIN users u ON u.id = o.usuario_id
WHERE o.caixa_id = $1 AND o.tenant_id = $2 AND o.tipo = $3
ORDER BY o.created_at ASC;

-- name: SumOperacoesByTipo :many
SELECT 
    tipo,
    COALESCE(SUM(valor), 0) as total
FROM operacoes_caixa
WHERE caixa_id = $1 AND tenant_id = $2
GROUP BY tipo;

-- name: GetLastOperacao :one
SELECT 
    o.*,
    u.nome as usuario_nome
FROM operacoes_caixa o
LEFT JOIN users u ON u.id = o.usuario_id
WHERE o.caixa_id = $1 AND o.tenant_id = $2
ORDER BY o.created_at DESC
LIMIT 1;
