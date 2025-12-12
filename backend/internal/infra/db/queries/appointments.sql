-- ============================================================================
-- APPOINTMENTS QUERIES (sqlc)
-- Módulo de Agendamento — NEXO v1.0
-- ============================================================================

-- name: CreateAppointment :one
INSERT INTO appointments (
    id,
    tenant_id,
    unit_id,
    professional_id,
    customer_id,
    start_time,
    end_time,
    status,
    total_price,
    notes,
    canceled_reason,
    google_calendar_event_id,
    command_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING *;

-- name: CreateAppointmentService :exec
INSERT INTO appointment_services (
    appointment_id,
    service_id,
    price_at_booking,
    duration_at_booking
) VALUES (
    $1, $2, $3, $4
);

-- name: GetAppointmentByID :one
SELECT 
    a.*,
    p.nome as professional_name,
    c.nome as customer_name,
    c.telefone as customer_phone
FROM appointments a
JOIN profissionais p ON p.id = a.professional_id
JOIN clientes c ON c.id = a.customer_id
WHERE a.id = $1 AND a.tenant_id = $2
  AND (sqlc.narg(unit_id)::uuid IS NULL OR a.unit_id = sqlc.narg(unit_id));

-- name: GetAppointmentServices :many
SELECT 
    aps.*,
    s.nome as service_name
FROM appointment_services aps
JOIN servicos s ON s.id = aps.service_id
WHERE aps.appointment_id = $1;

-- name: UpdateAppointment :one
UPDATE appointments
SET
    professional_id = $3,
    start_time = $4,
    end_time = $5,
    status = $6,
    total_price = $7,
    notes = $8,
    canceled_reason = $9,
    google_calendar_event_id = $10,
    checked_in_at = $11,
    started_at = $12,
    finished_at = $13,
    command_id = $14,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: UpdateAppointmentStatus :one
UPDATE appointments
SET
    status = $3,
    canceled_reason = $4,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: CheckInAppointment :one
-- Marca que o cliente chegou para o atendimento
UPDATE appointments
SET
    status = 'CHECKED_IN',
    checked_in_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
  AND status IN ('CREATED', 'CONFIRMED')
RETURNING *;

-- name: StartAppointment :one
-- Inicia o atendimento (profissional começou os serviços)
UPDATE appointments
SET
    status = 'IN_SERVICE',
    started_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
  AND status = 'CHECKED_IN'
RETURNING *;

-- name: FinishAppointment :one
-- Finaliza o atendimento (serviços concluídos, aguardando pagamento)
UPDATE appointments
SET
    status = 'AWAITING_PAYMENT',
    finished_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
  AND status = 'IN_SERVICE'
RETURNING *;

-- name: CompleteAppointment :one
-- Completa o agendamento após pagamento confirmado
UPDATE appointments
SET
    status = 'DONE',
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
  AND status = 'AWAITING_PAYMENT'
RETURNING *;

-- name: DeleteAppointment :exec
UPDATE appointments
SET
    status = 'CANCELED',
    canceled_reason = 'Removido pelo sistema',
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id));

-- name: DeleteAppointmentServices :exec
DELETE FROM appointment_services
WHERE appointment_id = $1;

-- name: ListAppointments :many
SELECT 
    a.*,
    p.nome as professional_name,
    c.nome as customer_name,
    c.telefone as customer_phone
FROM appointments a
JOIN profissionais p ON p.id = a.professional_id
JOIN clientes c ON c.id = a.customer_id
WHERE a.tenant_id = $1
  AND (sqlc.narg(unit_id)::uuid IS NULL OR a.unit_id = sqlc.narg(unit_id))
  AND ($2::uuid IS NULL OR a.professional_id = $2)
  AND ($3::uuid IS NULL OR a.customer_id = $3)
  AND (COALESCE(array_length($4::text[], 1), 0) = 0 OR a.status = ANY($4::text[]))
  AND ($5::timestamptz IS NULL OR a.start_time >= $5)
  AND ($6::timestamptz IS NULL OR a.start_time < $6)
ORDER BY a.start_time DESC
LIMIT $7 OFFSET $8;

-- name: CountAppointments :one
SELECT COUNT(*)
FROM appointments a
WHERE a.tenant_id = $1
  AND (sqlc.narg(unit_id)::uuid IS NULL OR a.unit_id = sqlc.narg(unit_id))
  AND ($2::uuid IS NULL OR a.professional_id = $2)
  AND ($3::uuid IS NULL OR a.customer_id = $3)
  AND (COALESCE(array_length($4::text[], 1), 0) = 0 OR a.status = ANY($4::text[]))
  AND ($5::timestamptz IS NULL OR a.start_time >= $5)
  AND ($6::timestamptz IS NULL OR a.start_time < $6);

-- name: ListAppointmentsByProfessionalAndDateRange :many
SELECT 
    a.*,
    p.nome as professional_name,
    c.nome as customer_name,
    c.telefone as customer_phone
FROM appointments a
JOIN profissionais p ON p.id = a.professional_id
JOIN clientes c ON c.id = a.customer_id
WHERE a.tenant_id = $1
  AND (sqlc.narg(unit_id)::uuid IS NULL OR a.unit_id = sqlc.narg(unit_id))
  AND a.professional_id = $2
  AND a.start_time >= $3
  AND a.start_time < $4
  AND a.status NOT IN ('CANCELED')
ORDER BY a.start_time ASC;

-- name: ListAppointmentsByCustomer :many
SELECT 
    a.*,
    p.nome as professional_name,
    c.nome as customer_name,
    c.telefone as customer_phone
FROM appointments a
JOIN profissionais p ON p.id = a.professional_id
JOIN clientes c ON c.id = a.customer_id
WHERE a.tenant_id = $1
  AND (sqlc.narg(unit_id)::uuid IS NULL OR a.unit_id = sqlc.narg(unit_id))
  AND a.customer_id = $2
ORDER BY a.start_time DESC
LIMIT 50;

-- name: CheckAppointmentConflict :one
-- Verifica conflito com agendamentos existentes
-- Parâmetros: tenant_id, professional_id, exclude_id, start_time, end_time
SELECT EXISTS (
    SELECT 1 FROM appointments
    WHERE tenant_id = $1
      AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
      AND professional_id = $2
      AND id != $3
      AND status NOT IN ('CANCELED', 'NO_SHOW')
      AND start_time < $5
      AND end_time > $4
) as has_conflict;

-- name: CheckBlockedTimeConflictForAppointment :one
-- Verifica se há conflito com horários bloqueados (blocked_times)
SELECT EXISTS (
    SELECT 1 FROM blocked_times
    WHERE tenant_id = sqlc.arg(tenant_id)::uuid
      AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
      AND professional_id = sqlc.arg(professional_id)::uuid
      AND start_time < sqlc.arg(end_time)::timestamptz
      AND end_time > sqlc.arg(start_time)::timestamptz
) as has_blocked_conflict;

-- name: CheckMinimumIntervalConflict :one
-- Verifica se há conflito de intervalo mínimo (10 minutos entre agendamentos)
-- Um agendamento que termina exatamente quando outro começa não é conflito,
-- mas se o intervalo for menor que 10 minutos, é conflito.
SELECT EXISTS (
    SELECT 1 FROM appointments
    WHERE tenant_id = sqlc.arg(tenant_id)::uuid
      AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
      AND professional_id = sqlc.arg(professional_id)::uuid
      AND id != sqlc.arg(exclude_id)::uuid
      AND status NOT IN ('CANCELED', 'NO_SHOW')
      AND (
          -- Agendamento existente termina menos de X minutos antes do novo início
          (end_time > sqlc.arg(start_time)::timestamptz - (sqlc.arg(interval_minutes)::int * interval '1 minute') AND end_time <= sqlc.arg(start_time)::timestamptz)
          OR
          -- Novo agendamento termina menos de X minutos antes do início existente
          (sqlc.arg(end_time)::timestamptz > start_time - (sqlc.arg(interval_minutes)::int * interval '1 minute') AND sqlc.arg(end_time)::timestamptz <= start_time)
      )
) as has_interval_conflict;

-- name: CountAppointmentsByStatus :one
SELECT COUNT(*)
FROM appointments
WHERE tenant_id = $1 
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
  AND status = $2;

-- name: GetDailyAppointmentStats :one
SELECT 
    COUNT(*) as total_appointments,
    COUNT(*) FILTER (WHERE status = 'DONE') as completed_count,
    COUNT(*) FILTER (WHERE status = 'CANCELED') as canceled_count,
    COUNT(*) FILTER (WHERE status = 'NO_SHOW') as no_show_count,
    COALESCE(SUM(total_price) FILTER (WHERE status = 'DONE'), 0) as total_revenue
FROM appointments
WHERE tenant_id = $1
  AND (sqlc.narg(unit_id)::uuid IS NULL OR unit_id = sqlc.narg(unit_id))
  AND start_time >= $2
  AND start_time < $3;

-- ============================================================================
-- QUERIES AUXILIARES: Profissionais, Clientes, Serviços (Read-Only)
-- ============================================================================

-- name: ProfessionalExists :one
SELECT EXISTS (
    SELECT 1 FROM profissionais
    WHERE id = $1 AND tenant_id = $2 AND status = 'ATIVO'
) as exists;

-- name: GetProfessionalInfo :one
SELECT id, nome, status, NULL::text as cor, comissao::text as comissao, tipo_comissao
FROM profissionais
WHERE id = $1 AND tenant_id = $2;

-- name: ListActiveProfessionals :many
SELECT id, nome, status, NULL::text as cor
FROM profissionais
WHERE tenant_id = $1 AND status = 'ATIVO'
ORDER BY nome ASC;

-- name: CustomerExists :one
SELECT EXISTS (
    SELECT 1 FROM clientes
    WHERE id = $1 AND tenant_id = $2 AND ativo = true
) as exists;

-- name: GetCustomerInfo :one
SELECT id, nome, telefone, email
FROM clientes
WHERE id = $1 AND tenant_id = $2;

-- name: ServiceExists :one
SELECT EXISTS (
    SELECT 1 FROM servicos
    WHERE id = $1 AND tenant_id = $2 AND ativo = true
) as exists;

-- name: GetServiceInfo :one
SELECT id, nome, preco, duracao, ativo, comissao::text as comissao, categoria_id
FROM servicos
WHERE id = $1 AND tenant_id = $2;

-- name: GetServicesByIDs :many
SELECT id, nome, preco, duracao, ativo, comissao::text as comissao
FROM servicos
WHERE tenant_id = $1 AND id = ANY($2::uuid[])
ORDER BY nome ASC;

-- name: GetServicesForAppointments :many
-- Busca todos os serviços de múltiplos agendamentos de uma vez (evita N+1)
SELECT 
    aps.appointment_id,
    aps.service_id,
    aps.price_at_booking,
    aps.duration_at_booking,
    aps.created_at,
    s.nome as service_name
FROM appointment_services aps
JOIN servicos s ON s.id = aps.service_id
WHERE aps.appointment_id = ANY($1::uuid[])
ORDER BY aps.appointment_id, s.nome;
