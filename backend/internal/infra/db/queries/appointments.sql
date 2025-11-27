-- ============================================================================
-- APPOINTMENTS QUERIES (sqlc)
-- Módulo de Agendamento — NEXO v1.0
-- ============================================================================

-- name: CreateAppointment :one
INSERT INTO appointments (
    id,
    tenant_id,
    professional_id,
    customer_id,
    start_time,
    end_time,
    status,
    total_price,
    notes,
    canceled_reason,
    google_calendar_event_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
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
WHERE a.id = $1 AND a.tenant_id = $2;

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

-- name: DeleteAppointment :exec
UPDATE appointments
SET
    status = 'CANCELED',
    canceled_reason = 'Removido pelo sistema',
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2;

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
  AND ($2::uuid IS NULL OR a.professional_id = $2)
  AND ($3::uuid IS NULL OR a.customer_id = $3)
  AND ($4::text IS NULL OR a.status = $4)
  AND ($5::timestamptz IS NULL OR a.start_time >= $5)
  AND ($6::timestamptz IS NULL OR a.start_time < $6)
ORDER BY a.start_time DESC
LIMIT $7 OFFSET $8;

-- name: CountAppointments :one
SELECT COUNT(*)
FROM appointments a
WHERE a.tenant_id = $1
  AND ($2::uuid IS NULL OR a.professional_id = $2)
  AND ($3::uuid IS NULL OR a.customer_id = $3)
  AND ($4::text IS NULL OR a.status = $4)
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
  AND a.customer_id = $2
ORDER BY a.start_time DESC
LIMIT 50;

-- name: CheckAppointmentConflict :one
SELECT EXISTS (
    SELECT 1 FROM appointments
    WHERE tenant_id = $1
      AND professional_id = $2
      AND id != $3
      AND status NOT IN ('CANCELED', 'NO_SHOW')
      AND start_time < $5
      AND end_time > $4
) as has_conflict;

-- name: CountAppointmentsByStatus :one
SELECT COUNT(*)
FROM appointments
WHERE tenant_id = $1 AND status = $2;

-- name: GetDailyAppointmentStats :one
SELECT 
    COUNT(*) as total_appointments,
    COUNT(*) FILTER (WHERE status = 'DONE') as completed_count,
    COUNT(*) FILTER (WHERE status = 'CANCELED') as canceled_count,
    COUNT(*) FILTER (WHERE status = 'NO_SHOW') as no_show_count,
    COALESCE(SUM(total_price) FILTER (WHERE status = 'DONE'), 0) as total_revenue
FROM appointments
WHERE tenant_id = $1
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
SELECT id, nome, status, NULL::text as cor
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
SELECT id, nome, preco, duracao, ativo
FROM servicos
WHERE id = $1 AND tenant_id = $2;

-- name: GetServicesByIDs :many
SELECT id, nome, preco, duracao, ativo
FROM servicos
WHERE tenant_id = $1 AND id = ANY($2::uuid[])
ORDER BY nome ASC;
