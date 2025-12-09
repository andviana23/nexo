-- =============================================
-- Commands Queries (sqlc)
-- =============================================

-- name: CreateCommand :one
INSERT INTO commands (
    id,
    tenant_id,
    appointment_id,
    customer_id,
    numero,
    status,
    subtotal,
    desconto,
    total,
    total_recebido,
    troco,
    saldo_devedor,
    observacoes,
    deixar_troco_gorjeta,
    deixar_saldo_divida,
    criado_em,
    atualizado_em
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
) RETURNING *;

-- name: GetCommandByID :one
SELECT * FROM commands
WHERE id = $1 AND tenant_id = $2;

-- name: GetCommandByAppointmentID :one
SELECT * FROM commands
WHERE appointment_id = $1 AND tenant_id = $2
ORDER BY criado_em DESC
LIMIT 1;

-- name: UpdateCommand :one
UPDATE commands SET
    status = $3,
    subtotal = $4,
    desconto = $5,
    total = $6,
    total_recebido = $7,
    troco = $8,
    saldo_devedor = $9,
    observacoes = $10,
    deixar_troco_gorjeta = $11,
    deixar_saldo_divida = $12,
    fechado_em = $13,
    fechado_por = $14,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteCommand :exec
UPDATE commands SET
    status = 'CANCELED',
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2;

-- name: ListCommands :many
SELECT * FROM commands
WHERE tenant_id = $1
    AND ($2::VARCHAR IS NULL OR status = $2)
    AND ($3::UUID IS NULL OR customer_id = $3)
    AND ($4::DATE IS NULL OR DATE(criado_em) >= $4)
    AND ($5::DATE IS NULL OR DATE(criado_em) <= $5)
ORDER BY criado_em DESC
LIMIT $6 OFFSET $7;

-- name: CountCommands :one
SELECT COUNT(*) FROM commands
WHERE tenant_id = $1
    AND ($2::VARCHAR IS NULL OR status = $2)
    AND ($3::UUID IS NULL OR customer_id = $3)
    AND ($4::DATE IS NULL OR DATE(criado_em) >= $4)
    AND ($5::DATE IS NULL OR DATE(criado_em) <= $5);

-- name: GetNextCommandNumber :one
-- Retorna o próximo número sequencial para comandas do tenant no ano atual
SELECT COALESCE(MAX(
    CAST(
        NULLIF(REGEXP_REPLACE(numero, '[^0-9]', '', 'g'), '') 
        AS INTEGER
    )
), 0) + 1 as next_number
FROM commands
WHERE tenant_id = $1
    AND EXTRACT(YEAR FROM criado_em) = EXTRACT(YEAR FROM NOW());

-- =============================================
-- Command Items Queries
-- =============================================

-- name: CreateCommandItem :one
INSERT INTO command_items (
    id,
    command_id,
    tipo,
    item_id,
    descricao,
    preco_unitario,
    quantidade,
    desconto_valor,
    desconto_percentual,
    preco_final,
    observacoes,
    criado_em
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: GetCommandItems :many
SELECT ci.* FROM command_items ci
INNER JOIN commands c ON c.id = ci.command_id
WHERE ci.command_id = $1 AND c.tenant_id = $2
ORDER BY ci.criado_em ASC;

-- name: GetCommandItemByID :one
SELECT ci.* FROM command_items ci
INNER JOIN commands c ON c.id = ci.command_id
WHERE ci.id = $1 AND c.tenant_id = $2;

-- name: UpdateCommandItem :one
UPDATE command_items SET
    preco_unitario = $3,
    quantidade = $4,
    desconto_valor = $5,
    desconto_percentual = $6,
    preco_final = $7,
    observacoes = $8
WHERE command_items.id = $1
    AND EXISTS (
        SELECT 1 FROM commands 
        WHERE commands.id = command_items.command_id 
        AND commands.tenant_id = $2
    )
RETURNING *;

-- name: DeleteCommandItem :exec
DELETE FROM command_items
WHERE command_items.id = $1
    AND EXISTS (
        SELECT 1 FROM commands 
        WHERE commands.id = command_items.command_id 
        AND commands.tenant_id = $2
    );

-- =============================================
-- Command Payments Queries
-- =============================================

-- name: CreateCommandPayment :one
INSERT INTO command_payments (
    id,
    command_id,
    meio_pagamento_id,
    valor_recebido,
    taxa_percentual,
    taxa_fixa,
    valor_liquido,
    observacoes,
    criado_em,
    criado_por
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetCommandPayments :many
SELECT cp.* FROM command_payments cp
INNER JOIN commands c ON c.id = cp.command_id
WHERE cp.command_id = $1 AND c.tenant_id = $2
ORDER BY cp.criado_em ASC;

-- name: GetCommandPaymentByID :one
SELECT cp.* FROM command_payments cp
INNER JOIN commands c ON c.id = cp.command_id
WHERE cp.id = $1 AND c.tenant_id = $2;

-- name: DeleteCommandPayment :exec
DELETE FROM command_payments
WHERE command_payments.id = $1
    AND EXISTS (
        SELECT 1 FROM commands 
        WHERE commands.id = command_payments.command_id 
        AND commands.tenant_id = $2
    );

-- name: GetCommandPaymentsSummary :one
SELECT 
    COALESCE(SUM(valor_recebido), 0) as total_recebido,
    COALESCE(SUM(valor_liquido), 0) as total_liquido,
    COALESCE(SUM(valor_recebido - valor_liquido), 0) as total_taxas,
    COUNT(*) as total_pagamentos
FROM command_payments cp
INNER JOIN commands c ON c.id = cp.command_id
WHERE cp.command_id = $1 AND c.tenant_id = $2;
