-- name: CreateContaPagar :one
INSERT INTO contas_a_pagar (
    tenant_id,
    unit_id,
    descricao,
    categoria_id,
    fornecedor,
    valor,
    tipo,
    recorrente,
    periodicidade,
    data_vencimento,
    data_pagamento,
    status,
    comprovante_url,
    pix_code,
    observacoes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING *;

-- name: GetContaPagarByID :one
SELECT * FROM contas_a_pagar
WHERE id = $1 AND tenant_id = $2;

-- name: ListContasPagarByTenant :many
SELECT * FROM contas_a_pagar
WHERE tenant_id = $1
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
ORDER BY data_vencimento DESC
LIMIT $2 OFFSET $3;

-- name: ListContasPagarByStatus :many
SELECT * FROM contas_a_pagar
WHERE tenant_id = $1
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
  AND status = $2
ORDER BY data_vencimento ASC
LIMIT $3 OFFSET $4;

-- name: ListContasPagarByPeriod :many
SELECT * FROM contas_a_pagar
WHERE tenant_id = $1
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
  AND data_vencimento >= $2
  AND data_vencimento <= $3
ORDER BY data_vencimento ASC;

-- name: ListContasPagarFiltered :many
SELECT * FROM contas_a_pagar
WHERE tenant_id = $1
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
  AND (sqlc.narg(status)::text IS NULL OR status = sqlc.narg(status))
  AND (sqlc.narg(tipo)::text IS NULL OR tipo = sqlc.narg(tipo))
  AND (sqlc.narg(categoria_id)::uuid IS NULL OR categoria_id = sqlc.narg(categoria_id))
  AND (sqlc.narg(data_inicio)::date IS NULL OR data_vencimento >= sqlc.narg(data_inicio))
  AND (sqlc.narg(data_fim)::date IS NULL OR data_vencimento <= sqlc.narg(data_fim))
ORDER BY data_vencimento DESC
LIMIT $2 OFFSET $3;

-- name: ListContasPagarVencidas :many
SELECT * FROM contas_a_pagar
WHERE tenant_id = $1
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
  AND status IN ('ABERTO', 'ATRASADO')
  AND data_vencimento < $2
ORDER BY data_vencimento ASC;

-- name: ListContasPagarRecorrentes :many
SELECT * FROM contas_a_pagar
WHERE tenant_id = $1 AND recorrente = true
ORDER BY data_vencimento DESC;

-- name: UpdateContaPagar :one
UPDATE contas_a_pagar
SET
    descricao = $3,
    categoria_id = $4,
    fornecedor = $5,
    valor = $6,
    tipo = $7,
    recorrente = $8,
    periodicidade = $9,
    data_vencimento = $10,
    data_pagamento = $11,
    status = $12,
    comprovante_url = $13,
    pix_code = $14,
    observacoes = $15,
    unit_id = COALESCE(sqlc.narg(unit_id), unit_id),
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: MarcarContaPagarComoPaga :one
UPDATE contas_a_pagar
SET
    status = 'PAGO',
    data_pagamento = $3,
    atualizado_em = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: MarcarContaPagarComoAtrasada :exec
UPDATE contas_a_pagar
SET
    status = 'ATRASADO',
    atualizado_em = NOW()
WHERE tenant_id = $1
  AND status = 'ABERTO'
  AND data_vencimento < $2;

-- name: DeleteContaPagar :exec
DELETE FROM contas_a_pagar
WHERE id = $1 AND tenant_id = $2;

-- name: SumContasPagarByPeriod :one
SELECT
    COALESCE(SUM(valor), 0) as total_a_pagar
FROM contas_a_pagar
WHERE tenant_id = $1
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
  AND data_vencimento >= $2
  AND data_vencimento <= $3
  AND status != 'CANCELADO';

-- name: SumContasPagasByPeriod :one
SELECT
    COALESCE(SUM(valor), 0) as total_pago
FROM contas_a_pagar
WHERE tenant_id = $1
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
  AND data_pagamento >= $2
  AND data_pagamento <= $3
  AND status = 'PAGO';

-- name: CountContasPagarByStatus :one
SELECT COUNT(*) FROM contas_a_pagar
WHERE tenant_id = $1 AND status = $2;

-- name: CountContasPagarByTenant :one
SELECT COUNT(*) FROM contas_a_pagar
WHERE tenant_id = $1;
