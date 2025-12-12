-- =============================================================================
-- QUERIES PARA PROFISSIONAIS
-- =============================================================================

-- name: ListProfessionals :many
SELECT id, tenant_id, unit_id, user_id, nome, email, telefone, cpf, especialidades, 
       comissao, tipo_comissao, foto, data_admissao, data_demissao, status, 
       horario_trabalho, observacoes, criado_em, atualizado_em, tipo
FROM profissionais
WHERE tenant_id = @tenant_id
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
  AND (sqlc.narg(status)::text IS NULL OR status = sqlc.narg(status))
  AND (sqlc.narg(tipo)::text IS NULL OR tipo = sqlc.narg(tipo))
  AND (
    sqlc.narg(search)::text IS NULL 
    OR nome ILIKE '%' || sqlc.narg(search) || '%'
    OR email ILIKE '%' || sqlc.narg(search) || '%'
    OR cpf ILIKE '%' || sqlc.narg(search) || '%'
  )
ORDER BY 
  CASE WHEN sqlc.narg(order_by)::text = 'nome' THEN nome END ASC,
  CASE WHEN sqlc.narg(order_by)::text = 'criado_em' OR sqlc.narg(order_by)::text IS NULL THEN criado_em END DESC,
  CASE WHEN sqlc.narg(order_by)::text = 'data_admissao' THEN data_admissao END DESC
LIMIT @page_size OFFSET @page_offset;

-- name: CountProfessionals :one
SELECT COUNT(*) FROM profissionais
WHERE tenant_id = @tenant_id
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
  AND (sqlc.narg(status)::text IS NULL OR status = sqlc.narg(status))
  AND (sqlc.narg(tipo)::text IS NULL OR tipo = sqlc.narg(tipo))
  AND (
    sqlc.narg(search)::text IS NULL 
    OR nome ILIKE '%' || sqlc.narg(search) || '%'
    OR email ILIKE '%' || sqlc.narg(search) || '%'
    OR cpf ILIKE '%' || sqlc.narg(search) || '%'
  );

-- name: GetProfessionalByID :one
SELECT id, tenant_id, unit_id, user_id, nome, email, telefone, cpf, especialidades, 
       comissao, tipo_comissao, foto, data_admissao, data_demissao, status, 
       horario_trabalho, observacoes, criado_em, atualizado_em, tipo
FROM profissionais
WHERE id = @id AND tenant_id = @tenant_id
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id));

-- name: CreateProfessional :one
INSERT INTO profissionais (
    tenant_id, unit_id, user_id, nome, email, telefone, cpf, especialidades,
    comissao, tipo_comissao, foto, data_admissao, status, horario_trabalho,
    observacoes, tipo
) VALUES (
    @tenant_id, sqlc.narg(unit_id), sqlc.narg(user_id), @nome, @email, @telefone, @cpf, @especialidades,
    @comissao, @tipo_comissao, sqlc.narg(foto), @data_admissao, @status, sqlc.narg(horario_trabalho),
    sqlc.narg(observacoes), @tipo
)
RETURNING id, tenant_id, unit_id, user_id, nome, email, telefone, cpf, especialidades, 
          comissao, tipo_comissao, foto, data_admissao, data_demissao, status, 
          horario_trabalho, observacoes, criado_em, atualizado_em, tipo;

-- name: UpdateProfessional :one
UPDATE profissionais SET
    nome = @nome,
    email = @email,
    telefone = @telefone,
    cpf = @cpf,
    especialidades = @especialidades,
    comissao = @comissao,
    tipo_comissao = @tipo_comissao,
    foto = sqlc.narg(foto),
    data_admissao = @data_admissao,
    data_demissao = sqlc.narg(data_demissao),
    status = @status,
    horario_trabalho = sqlc.narg(horario_trabalho),
    observacoes = sqlc.narg(observacoes),
    tipo = @tipo,
    atualizado_em = NOW()
WHERE id = @id AND tenant_id = @tenant_id
  AND unit_id = @unit_id
RETURNING id, tenant_id, unit_id, user_id, nome, email, telefone, cpf, especialidades, 
          comissao, tipo_comissao, foto, data_admissao, data_demissao, status, 
          horario_trabalho, observacoes, criado_em, atualizado_em, tipo;

-- name: UpdateProfessionalStatus :one
UPDATE profissionais SET
    status = @status,
    data_demissao = sqlc.narg(data_demissao),
    atualizado_em = NOW()
WHERE id = @id AND tenant_id = @tenant_id
  AND unit_id = @unit_id
RETURNING id, tenant_id, unit_id, user_id, nome, email, telefone, cpf, especialidades, 
          comissao, tipo_comissao, foto, data_admissao, data_demissao, status, 
          horario_trabalho, observacoes, criado_em, atualizado_em, tipo;

-- name: DeleteProfessional :exec
DELETE FROM profissionais
WHERE id = @id AND tenant_id = @tenant_id
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id));

-- name: CheckEmailExistsProfessional :one
SELECT EXISTS (
    SELECT 1 FROM profissionais
    WHERE tenant_id = @tenant_id 
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
      AND email = @email 
      AND (sqlc.narg(exclude_id)::uuid IS NULL OR id != sqlc.narg(exclude_id))
) as exists;

-- name: CheckCpfExistsProfessional :one
SELECT EXISTS (
    SELECT 1 FROM profissionais
    WHERE tenant_id = @tenant_id 
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
      AND cpf = @cpf 
      AND (sqlc.narg(exclude_id)::uuid IS NULL OR id != sqlc.narg(exclude_id))
) as exists;

-- name: ListBarbers :many
SELECT id, tenant_id, unit_id, user_id, nome, email, telefone, cpf, especialidades, 
       comissao, tipo_comissao, foto, data_admissao, data_demissao, status, 
       horario_trabalho, observacoes, criado_em, atualizado_em, tipo
FROM profissionais
WHERE tenant_id = @tenant_id 
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id)) 
  AND tipo = 'BARBEIRO' 
  AND status = 'ATIVO'
ORDER BY nome;

-- name: ListProfessionalCategoryCommissions :many
SELECT * FROM comissoes_categoria_profissional
WHERE tenant_id = @tenant_id AND profissional_id = @profissional_id;

-- name: CreateProfessionalCategoryCommission :one
INSERT INTO comissoes_categoria_profissional (
  tenant_id, profissional_id, categoria_id, comissao
) VALUES (
  @tenant_id, @profissional_id, @categoria_id, @comissao
) RETURNING *;

-- name: DeleteProfessionalCategoryCommissionsByProfessional :exec
DELETE FROM comissoes_categoria_profissional
WHERE tenant_id = @tenant_id AND profissional_id = @profissional_id;

-- name: GetProfessionalCategoryCommission :one
-- Busca comissão específica de um profissional para uma categoria de serviço
SELECT comissao::text as comissao
FROM comissoes_categoria_profissional
WHERE tenant_id = @tenant_id 
  AND profissional_id = @profissional_id 
  AND categoria_id = @categoria_id;
