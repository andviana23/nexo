-- ========================================
-- QUERIES SQL - MÓDULO DE ESTOQUE
-- Geradas via sqlc para type-safety
-- ========================================

-- ========================================
-- FORNECEDORES
-- ========================================

-- name: CreateFornecedor :one
INSERT INTO fornecedores (
    tenant_id,
    razao_social,
    nome_fantasia,
    cnpj,
    email,
    telefone,
    celular,
    endereco_logradouro,
    endereco_numero,
    endereco_complemento,
    endereco_bairro,
    endereco_cidade,
    endereco_estado,
    endereco_cep,
    banco,
    agencia,
    conta,
    observacoes,
    ativo
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19
)
RETURNING *;

-- name: GetFornecedorByID :one
SELECT * FROM fornecedores
WHERE id = $1 AND tenant_id = $2;

-- name: GetFornecedorByCNPJ :one
SELECT * FROM fornecedores
WHERE cnpj = $1 AND tenant_id = $2;

-- name: ListFornecedores :many
SELECT * FROM fornecedores
WHERE tenant_id = $1
ORDER BY razao_social;

-- name: ListFornecedoresAtivos :many
SELECT * FROM fornecedores
WHERE tenant_id = $1 AND ativo = true
ORDER BY razao_social;

-- name: UpdateFornecedor :one
UPDATE fornecedores
SET
    razao_social = $3,
    nome_fantasia = $4,
    email = $5,
    telefone = $6,
    celular = $7,
    endereco_logradouro = $8,
    endereco_numero = $9,
    endereco_complemento = $10,
    endereco_bairro = $11,
    endereco_cidade = $12,
    endereco_estado = $13,
    endereco_cep = $14,
    banco = $15,
    agencia = $16,
    conta = $17,
    observacoes = $18,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteFornecedor :exec
UPDATE fornecedores
SET ativo = false, atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: ReativarFornecedor :exec
UPDATE fornecedores
SET ativo = true, atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;


-- ========================================
-- PRODUTOS (operações de estoque)
-- ========================================

-- name: CreateProduto :one
INSERT INTO produtos (
    tenant_id,
    categoria_id,
    nome,
    descricao,
    sku,
    codigo_barras,
    preco,
    custo,
    categoria_produto,
    unidade_medida,
    quantidade_atual,
    quantidade_minima,
    localizacao,
    lote,
    data_validade,
    ncm,
    permite_venda,
    ativo
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18
)
RETURNING *;

-- name: GetProdutoByID :one
SELECT * FROM produtos
WHERE id = $1 AND tenant_id = $2;

-- name: GetProdutoBySKU :one
SELECT * FROM produtos
WHERE sku = $1 AND tenant_id = $2;

-- name: ListProdutos :many
SELECT * FROM produtos
WHERE tenant_id = $1
ORDER BY nome;

-- name: ListProdutosAtivos :many
SELECT * FROM produtos
WHERE tenant_id = $1 AND ativo = true
ORDER BY nome;

-- name: ListProdutosByCategoria :many
SELECT * FROM produtos
WHERE tenant_id = $1 AND categoria_produto = $2 AND ativo = true
ORDER BY nome;

-- name: ListProdutosAbaixoDoMinimo :many
SELECT * FROM produtos
WHERE tenant_id = $1
  AND ativo = true
  AND quantidade_atual <= quantidade_minima
ORDER BY (quantidade_atual / NULLIF(quantidade_minima, 0)) ASC;

-- name: UpdateProduto :one
UPDATE produtos
SET
    categoria_id = $3,
    nome = $4,
    descricao = $5,
    sku = $6,
    codigo_barras = $7,
    preco = $8,
    custo = $9,
    categoria_produto = $10,
    unidade_medida = $11,
    quantidade_minima = $12,
    localizacao = $13,
    lote = $14,
    data_validade = $15,
    ncm = $16,
    permite_venda = $17,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: AtualizarQuantidadeProduto :one
UPDATE produtos
SET
    quantidade_atual = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteProduto :exec
UPDATE produtos
SET ativo = false, atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;


-- ========================================
-- MOVIMENTAÇÕES DE ESTOQUE
-- ========================================

-- name: CreateMovimentacaoEstoque :one
INSERT INTO movimentacoes_estoque (
    tenant_id,
    produto_id,
    tipo_movimentacao,
    quantidade,
    valor_unitario,
    valor_total,
    fornecedor_id,
    usuario_id,
    data_movimentacao,
    observacoes,
    documento
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: GetMovimentacaoByID :one
SELECT * FROM movimentacoes_estoque
WHERE id = $1 AND tenant_id = $2;

-- name: ListMovimentacoesByProduto :many
SELECT m.*
FROM movimentacoes_estoque m
WHERE m.tenant_id = $1 AND m.produto_id = $2
ORDER BY m.data_movimentacao DESC
LIMIT $3 OFFSET $4;

-- name: ListMovimentacoesByTipo :many
SELECT m.*
FROM movimentacoes_estoque m
WHERE m.tenant_id = $1 AND m.tipo_movimentacao = $2
ORDER BY m.data_movimentacao DESC
LIMIT $3 OFFSET $4;

-- name: ListMovimentacoesByPeriodo :many
SELECT m.*
FROM movimentacoes_estoque m
WHERE m.tenant_id = $1
  AND m.data_movimentacao >= $2
  AND m.data_movimentacao <= $3
ORDER BY m.data_movimentacao DESC;

-- name: ListMovimentacoesByFornecedor :many
SELECT m.*
FROM movimentacoes_estoque m
WHERE m.tenant_id = $1 AND m.fornecedor_id = $2
ORDER BY m.data_movimentacao DESC
LIMIT $3 OFFSET $4;

-- name: GetTotalPorTipo :one
SELECT
    tipo_movimentacao,
    SUM(quantidade) as total_quantidade,
    SUM(valor_total) as total_valor
FROM movimentacoes_estoque
WHERE tenant_id = $1
  AND tipo_movimentacao = $2
  AND data_movimentacao >= $3
  AND data_movimentacao <= $4
GROUP BY tipo_movimentacao;


-- ========================================
-- PRODUTO-FORNECEDOR (relacionamento)
-- ========================================

-- name: CreateProdutoFornecedor :one
INSERT INTO produto_fornecedor (
    tenant_id,
    produto_id,
    fornecedor_id,
    codigo_fornecedor,
    preco_compra,
    prazo_entrega_dias,
    fornecedor_preferencial,
    ativo
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetProdutoFornecedor :one
SELECT * FROM produto_fornecedor
WHERE produto_id = $1
  AND fornecedor_id = $2
  AND tenant_id = $3;

-- name: ListFornecedoresByProduto :many
SELECT
    pf.*,
    f.razao_social,
    f.nome_fantasia,
    f.cnpj,
    f.email,
    f.telefone
FROM produto_fornecedor pf
JOIN fornecedores f ON pf.fornecedor_id = f.id
WHERE pf.produto_id = $1
  AND pf.tenant_id = $2
  AND pf.ativo = true
ORDER BY pf.fornecedor_preferencial DESC, f.razao_social;

-- name: ListProdutosByFornecedor :many
SELECT
    pf.*,
    p.nome,
    p.sku,
    p.categoria_produto,
    p.unidade_medida
FROM produto_fornecedor pf
JOIN produtos p ON pf.produto_id = p.id
WHERE pf.fornecedor_id = $1
  AND pf.tenant_id = $2
  AND pf.ativo = true
ORDER BY p.nome;

-- name: UpdateProdutoFornecedor :one
UPDATE produto_fornecedor
SET
    codigo_fornecedor = $4,
    preco_compra = $5,
    prazo_entrega_dias = $6,
    fornecedor_preferencial = $7,
    atualizado_em = NOW()
WHERE produto_id = $1
  AND fornecedor_id = $2
  AND tenant_id = $3
RETURNING *;

-- name: DeleteProdutoFornecedor :exec
UPDATE produto_fornecedor
SET ativo = false, atualizado_em = NOW()
WHERE produto_id = $1
  AND fornecedor_id = $2
  AND tenant_id = $3;


-- ========================================
-- RELATÓRIOS E CONSULTAS ESPECIAIS
-- ========================================

-- name: GetProdutosComEstoqueBaixo :many
SELECT
    p.*,
    (p.quantidade_atual / NULLIF(p.quantidade_minima, 0)) * 100 as percentual_estoque
FROM produtos p
WHERE p.tenant_id = $1
  AND p.ativo = true
  AND p.quantidade_atual <= p.quantidade_minima
ORDER BY percentual_estoque ASC;

-- name: GetMovimentacoesPorPeriodoComDetalhes :many
SELECT
    m.*,
    p.nome as produto_nome,
    p.sku as produto_sku,
    f.razao_social as fornecedor_nome,
    u.nome as usuario_nome
FROM movimentacoes_estoque m
JOIN produtos p ON m.produto_id = p.id
LEFT JOIN fornecedores f ON m.fornecedor_id = f.id
LEFT JOIN users u ON m.usuario_id = u.id
WHERE m.tenant_id = $1
  AND m.data_movimentacao >= $2
  AND m.data_movimentacao <= $3
ORDER BY m.data_movimentacao DESC;

-- name: GetValorTotalEstoque :one
SELECT
    SUM(quantidade_atual * COALESCE(custo, 0)) as valor_total_estoque,
    COUNT(*) as total_produtos
FROM produtos
WHERE tenant_id = $1 AND ativo = true;

-- name: GetCurvaABC :many
SELECT
    p.id,
    p.nome,
    p.sku,
    p.categoria_produto,
    p.quantidade_atual,
    p.custo,
    (p.quantidade_atual * COALESCE(p.custo, 0)) as valor_estoque,
    CASE
        WHEN ROW_NUMBER() OVER (ORDER BY (p.quantidade_atual * COALESCE(p.custo, 0)) DESC) <= (COUNT(*) OVER () * 0.2) THEN 'A'
        WHEN ROW_NUMBER() OVER (ORDER BY (p.quantidade_atual * COALESCE(p.custo, 0)) DESC) <= (COUNT(*) OVER () * 0.5) THEN 'B'
        ELSE 'C'
    END as classe_abc
FROM produtos p
WHERE p.tenant_id = $1 AND p.ativo = true
ORDER BY valor_estoque DESC;
