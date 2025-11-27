-- ============================================================================
-- CUSTOMERS QUERIES (sqlc)
-- Módulo de Cadastro de Clientes — NEXO v1.0
-- Conforme FLUXO_CADASTROS_CLIENTE.md
-- ============================================================================

-- ============================================================================
-- CREATE
-- ============================================================================

-- name: CreateCustomer :one
INSERT INTO clientes (
    id,
    tenant_id,
    nome,
    telefone,
    email,
    cpf,
    data_nascimento,
    genero,
    endereco_logradouro,
    endereco_numero,
    endereco_complemento,
    endereco_bairro,
    endereco_cidade,
    endereco_estado,
    endereco_cep,
    observacoes,
    tags,
    ativo,
    criado_em,
    atualizado_em
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, true, NOW(), NOW()
) RETURNING *;

-- ============================================================================
-- READ
-- ============================================================================

-- name: GetCustomerByID :one
SELECT * FROM clientes
WHERE id = $1 AND tenant_id = $2;

-- name: GetCustomerByPhone :one
SELECT * FROM clientes
WHERE tenant_id = $1 AND telefone = $2 AND ativo = true;

-- name: GetCustomerByCPF :one
SELECT * FROM clientes
WHERE tenant_id = $1 AND cpf = $2 AND ativo = true;

-- name: ListCustomers :many
SELECT * FROM clientes
WHERE tenant_id = @tenant_id
  AND (sqlc.narg(ativo)::bool IS NULL OR ativo = sqlc.narg(ativo))
  AND (
    sqlc.narg(search)::text IS NULL 
    OR nome ILIKE '%' || sqlc.narg(search) || '%'
    OR telefone ILIKE '%' || sqlc.narg(search) || '%'
    OR cpf ILIKE '%' || sqlc.narg(search) || '%'
    OR email ILIKE '%' || sqlc.narg(search) || '%'
  )
  AND (sqlc.narg(tags)::text[] IS NULL OR tags && sqlc.narg(tags))
ORDER BY 
  CASE WHEN sqlc.narg(order_by) = 'nome' THEN nome END ASC,
  CASE WHEN sqlc.narg(order_by) = 'criado_em' OR sqlc.narg(order_by) IS NULL THEN criado_em END DESC,
  CASE WHEN sqlc.narg(order_by) = 'atualizado_em' THEN atualizado_em END DESC
LIMIT @page_size OFFSET @page_offset;

-- name: CountCustomers :one
SELECT COUNT(*) FROM clientes
WHERE tenant_id = @tenant_id
  AND (sqlc.narg(ativo)::bool IS NULL OR ativo = sqlc.narg(ativo))
  AND (
    sqlc.narg(search)::text IS NULL 
    OR nome ILIKE '%' || sqlc.narg(search) || '%'
    OR telefone ILIKE '%' || sqlc.narg(search) || '%'
    OR cpf ILIKE '%' || sqlc.narg(search) || '%'
    OR email ILIKE '%' || sqlc.narg(search) || '%'
  )
  AND (sqlc.narg(tags)::text[] IS NULL OR tags && sqlc.narg(tags));

-- name: ListActiveCustomers :many
SELECT id, nome, telefone, email, tags
FROM clientes
WHERE tenant_id = $1 AND ativo = true
ORDER BY nome ASC
LIMIT 100;

-- name: SearchCustomers :many
SELECT id, nome, telefone, email, tags
FROM clientes
WHERE tenant_id = $1 
  AND ativo = true
  AND (
    nome ILIKE '%' || $2 || '%'
    OR telefone ILIKE '%' || $2 || '%'
    OR email ILIKE '%' || $2 || '%'
  )
ORDER BY nome ASC
LIMIT 20;

-- ============================================================================
-- UPDATE
-- ============================================================================

-- name: UpdateCustomer :one
UPDATE clientes
SET
    nome = COALESCE($3, nome),
    telefone = COALESCE($4, telefone),
    email = COALESCE($5, email),
    cpf = COALESCE($6, cpf),
    data_nascimento = COALESCE($7, data_nascimento),
    genero = COALESCE($8, genero),
    endereco_logradouro = COALESCE($9, endereco_logradouro),
    endereco_numero = COALESCE($10, endereco_numero),
    endereco_complemento = COALESCE($11, endereco_complemento),
    endereco_bairro = COALESCE($12, endereco_bairro),
    endereco_cidade = COALESCE($13, endereco_cidade),
    endereco_estado = COALESCE($14, endereco_estado),
    endereco_cep = COALESCE($15, endereco_cep),
    observacoes = COALESCE($16, observacoes),
    tags = COALESCE($17, tags),
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: UpdateCustomerTags :one
UPDATE clientes
SET
    tags = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- ============================================================================
-- DELETE (Soft Delete)
-- ============================================================================

-- name: InactivateCustomer :exec
UPDATE clientes
SET
    ativo = false,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: ReactivateCustomer :exec
UPDATE clientes
SET
    ativo = true,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;

-- ============================================================================
-- VALIDAÇÕES
-- ============================================================================

-- name: CheckPhoneExists :one
SELECT EXISTS (
    SELECT 1 FROM clientes
    WHERE tenant_id = $1 
      AND telefone = $2 
      AND ativo = true
      AND ($3::uuid IS NULL OR id != $3)
) as exists;

-- name: CheckCPFExists :one
SELECT EXISTS (
    SELECT 1 FROM clientes
    WHERE tenant_id = $1 
      AND cpf = $2 
      AND ativo = true
      AND ($3::uuid IS NULL OR id != $3)
) as exists;

-- name: CheckEmailExists :one
SELECT EXISTS (
    SELECT 1 FROM clientes
    WHERE tenant_id = $1 
      AND email = $2 
      AND ativo = true
      AND ($3::uuid IS NULL OR id != $3)
) as exists;

-- ============================================================================
-- ESTATÍSTICAS E RELATÓRIOS
-- ============================================================================

-- name: GetCustomerStats :one
SELECT 
    COUNT(*) FILTER (WHERE ativo = true) as total_ativos,
    COUNT(*) FILTER (WHERE ativo = false) as total_inativos,
    COUNT(*) FILTER (WHERE criado_em >= NOW() - INTERVAL '30 days') as novos_ultimos_30_dias,
    COUNT(*) as total_geral
FROM clientes
WHERE tenant_id = $1;

-- name: GetCustomerWithHistory :one
SELECT 
    c.*,
    COALESCE(stats.total_atendimentos, 0) as total_atendimentos,
    COALESCE(stats.total_gasto, 0) as total_gasto,
    stats.ultimo_atendimento
FROM clientes c
LEFT JOIN LATERAL (
    SELECT 
        COUNT(*) as total_atendimentos,
        SUM(total_price) as total_gasto,
        MAX(start_time) as ultimo_atendimento
    FROM appointments
    WHERE customer_id = c.id 
      AND tenant_id = c.tenant_id
      AND status = 'DONE'
) stats ON true
WHERE c.id = $1 AND c.tenant_id = $2;

-- name: ListCustomersWithoutAppointments :many
SELECT c.*
FROM clientes c
WHERE c.tenant_id = $1 
  AND c.ativo = true
  AND NOT EXISTS (
    SELECT 1 FROM appointments a
    WHERE a.customer_id = c.id
      AND a.start_time >= NOW() - INTERVAL '90 days'
  )
ORDER BY c.criado_em DESC
LIMIT $2 OFFSET $3;

-- ============================================================================
-- EXPORTAÇÃO LGPD
-- ============================================================================

-- name: GetCustomerDataForExport :one
SELECT 
    c.*,
    COALESCE(
        (SELECT json_agg(json_build_object(
            'data', a.start_time,
            'status', a.status,
            'valor_total', a.total_price,
            'profissional', p.nome
        ) ORDER BY a.start_time DESC)
        FROM appointments a
        JOIN profissionais p ON p.id = a.professional_id
        WHERE a.customer_id = c.id AND a.tenant_id = c.tenant_id
        ), '[]'::json
    ) as historico_atendimentos
FROM clientes c
WHERE c.id = $1 AND c.tenant_id = $2;
