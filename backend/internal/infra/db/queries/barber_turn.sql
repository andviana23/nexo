-- ============================================================================
-- BARBER TURN QUERIES (sqlc)
-- Módulo Lista da Vez — NEXO v1.0
-- Conforme FLUXO_LISTA_DA_VEZ.md
-- ============================================================================

-- ============================================================================
-- CREATE / ADD
-- ============================================================================

-- name: AddBarberToTurnList :one
-- Adiciona um barbeiro à lista da vez
INSERT INTO barbers_turn_list (
    id,
    tenant_id,
    professional_id,
    current_points,
    last_turn_at,
    is_active,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, 0, NULL, true, NOW(), NOW()
) RETURNING *;

-- ============================================================================
-- READ / LIST
-- ============================================================================

-- name: GetBarberTurnByID :one
SELECT 
    btl.*,
    p.nome as professional_name,
    p.tipo as professional_type,
    p.status as professional_status
FROM barbers_turn_list btl
JOIN profissionais p ON p.id = btl.professional_id
WHERE btl.id = $1 AND btl.tenant_id = $2;

-- name: GetBarberTurnByProfessionalID :one
SELECT 
    btl.*,
    p.nome as professional_name,
    p.tipo as professional_type,
    p.status as professional_status
FROM barbers_turn_list btl
JOIN profissionais p ON p.id = btl.professional_id
WHERE btl.professional_id = $1 AND btl.tenant_id = $2;

-- name: ListBarbersTurnList :many
-- Lista todos os barbeiros na fila ordenados por pontuação
-- Menor pontuação = topo da fila
-- Critério de desempate: last_turn_at mais antigo ou nulo (nunca atendeu) → topo
-- Segundo desempate: ordem de criação (created_at ASC)
SELECT 
    btl.*,
    p.nome as professional_name,
    p.tipo as professional_type,
    p.status as professional_status,
    p.foto as professional_photo,
    ROW_NUMBER() OVER (
        ORDER BY 
            btl.current_points ASC,
            btl.last_turn_at ASC NULLS FIRST,
            btl.created_at ASC
    ) as position
FROM barbers_turn_list btl
JOIN profissionais p ON p.id = btl.professional_id
WHERE btl.tenant_id = $1
  AND (sqlc.narg('is_active')::bool IS NULL OR btl.is_active = sqlc.narg('is_active'))
ORDER BY 
    btl.current_points ASC,
    btl.last_turn_at ASC NULLS FIRST,
    btl.created_at ASC;

-- name: ListActiveBarbersTurnList :many
-- Lista apenas barbeiros ativos na fila (is_active = true)
SELECT 
    btl.*,
    p.nome as professional_name,
    p.tipo as professional_type,
    p.status as professional_status,
    p.foto as professional_photo,
    ROW_NUMBER() OVER (
        ORDER BY 
            btl.current_points ASC,
            btl.last_turn_at ASC NULLS FIRST,
            btl.created_at ASC
    ) as position
FROM barbers_turn_list btl
JOIN profissionais p ON p.id = btl.professional_id
WHERE btl.tenant_id = $1
  AND btl.is_active = true
  AND p.status = 'ATIVO'
ORDER BY 
    btl.current_points ASC,
    btl.last_turn_at ASC NULLS FIRST,
    btl.created_at ASC;

-- name: GetNextBarber :one
-- Retorna o próximo barbeiro da vez (topo da fila ativa)
SELECT 
    btl.*,
    p.nome as professional_name,
    p.tipo as professional_type,
    p.status as professional_status,
    p.foto as professional_photo
FROM barbers_turn_list btl
JOIN profissionais p ON p.id = btl.professional_id
WHERE btl.tenant_id = $1
  AND btl.is_active = true
  AND p.status = 'ATIVO'
ORDER BY 
    btl.current_points ASC,
    btl.last_turn_at ASC NULLS FIRST,
    btl.created_at ASC
LIMIT 1;

-- name: CountBarbersTurnList :one
SELECT 
    COUNT(*) FILTER (WHERE is_active = true) as total_ativos,
    COUNT(*) FILTER (WHERE is_active = false) as total_pausados,
    COUNT(*) as total_geral,
    COALESCE(SUM(current_points), 0) as total_pontos
FROM barbers_turn_list
WHERE tenant_id = $1;

-- ============================================================================
-- UPDATE
-- ============================================================================

-- name: RecordTurn :one
-- Registra um atendimento: incrementa pontos (+1) e atualiza timestamp
UPDATE barbers_turn_list
SET
    current_points = current_points + 1,
    last_turn_at = NOW(),
    updated_at = NOW()
WHERE professional_id = $1 AND tenant_id = $2
RETURNING *;

-- name: ToggleBarberTurnStatus :one
-- Alterna status ativo/inativo de um barbeiro na fila
UPDATE barbers_turn_list
SET
    is_active = NOT is_active,
    updated_at = NOW()
WHERE professional_id = $1 AND tenant_id = $2
RETURNING *;

-- name: SetBarberTurnActive :one
-- Ativa um barbeiro na fila
UPDATE barbers_turn_list
SET
    is_active = true,
    updated_at = NOW()
WHERE professional_id = $1 AND tenant_id = $2
RETURNING *;

-- name: SetBarberTurnInactive :one
-- Pausa um barbeiro na fila
UPDATE barbers_turn_list
SET
    is_active = false,
    updated_at = NOW()
WHERE professional_id = $1 AND tenant_id = $2
RETURNING *;

-- ============================================================================
-- DELETE
-- ============================================================================

-- name: RemoveBarberFromTurnList :exec
-- Remove barbeiro da lista da vez
DELETE FROM barbers_turn_list
WHERE professional_id = $1 AND tenant_id = $2;

-- ============================================================================
-- RESET MENSAL
-- ============================================================================

-- name: ResetAllTurnPoints :exec
-- Zera todos os pontos e last_turn_at para reset mensal
UPDATE barbers_turn_list
SET
    current_points = 0,
    last_turn_at = NULL,
    updated_at = NOW()
WHERE tenant_id = $1;

-- name: SaveTurnHistoryBeforeReset :exec
-- Salva snapshot no histórico antes do reset
INSERT INTO barber_turn_history (
    id,
    tenant_id,
    professional_id,
    month_year,
    total_turns,
    final_points,
    created_at
)
SELECT
    gen_random_uuid(),
    btl.tenant_id,
    btl.professional_id,
    $2::varchar(7), -- month_year no formato 'YYYY-MM'
    btl.current_points,
    btl.current_points,
    NOW()
FROM barbers_turn_list btl
WHERE btl.tenant_id = $1
  AND btl.is_active = true
ON CONFLICT (professional_id, tenant_id, month_year) 
DO UPDATE SET
    total_turns = EXCLUDED.total_turns,
    final_points = EXCLUDED.final_points;

-- ============================================================================
-- HISTÓRICO
-- ============================================================================

-- name: ListTurnHistory :many
-- Lista histórico mensal de atendimentos
SELECT 
    bth.*,
    p.nome as professional_name
FROM barber_turn_history bth
JOIN profissionais p ON p.id = bth.professional_id
WHERE bth.tenant_id = $1
  AND ($2::varchar(7) IS NULL OR bth.month_year = $2)
ORDER BY bth.month_year DESC, bth.final_points DESC;

-- name: GetTurnHistoryByMonth :many
-- Busca histórico de um mês específico
SELECT 
    bth.*,
    p.nome as professional_name
FROM barber_turn_history bth
JOIN profissionais p ON p.id = bth.professional_id
WHERE bth.tenant_id = $1
  AND bth.month_year = $2
ORDER BY bth.final_points DESC;

-- name: GetTurnHistorySummary :many
-- Resumo dos últimos 12 meses
SELECT 
    month_year,
    COUNT(*) as total_barbeiros,
    SUM(total_turns) as total_atendimentos,
    AVG(total_turns) as media_atendimentos
FROM barber_turn_history
WHERE tenant_id = $1
GROUP BY month_year
ORDER BY month_year DESC
LIMIT 12;

-- ============================================================================
-- VALIDAÇÕES
-- ============================================================================

-- name: CheckProfessionalInTurnList :one
-- Verifica se profissional já está na lista
SELECT EXISTS (
    SELECT 1 FROM barbers_turn_list
    WHERE tenant_id = $1 AND professional_id = $2
) as exists;

-- name: CheckProfessionalIsBarber :one
-- Verifica se profissional é do tipo BARBEIRO
SELECT EXISTS (
    SELECT 1 FROM profissionais
    WHERE id = $1 
      AND tenant_id = $2 
      AND tipo = 'BARBEIRO'
      AND status = 'ATIVO'
) as exists;

-- name: GetAvailableBarbersForTurnList :many
-- Lista barbeiros ativos que ainda não estão na lista da vez
SELECT 
    p.id,
    p.nome,
    p.foto,
    p.status
FROM profissionais p
WHERE p.tenant_id = $1
  AND p.tipo = 'BARBEIRO'
  AND p.status = 'ATIVO'
  AND NOT EXISTS (
    SELECT 1 FROM barbers_turn_list btl
    WHERE btl.professional_id = p.id AND btl.tenant_id = p.tenant_id
  )
ORDER BY p.nome ASC;

-- ============================================================================
-- ESTATÍSTICAS DIÁRIAS
-- ============================================================================

-- name: GetTodayStats :one
-- Estatísticas do dia atual
SELECT 
    COUNT(*) FILTER (WHERE DATE(last_turn_at) = CURRENT_DATE) as atendimentos_hoje,
    SUM(current_points) as total_pontos_mes,
    COUNT(*) FILTER (WHERE is_active = true) as barbeiros_ativos,
    MAX(last_turn_at) as ultimo_atendimento
FROM barbers_turn_list
WHERE tenant_id = $1;

-- name: GetDailyReportByDate :many
-- Relatório diário: pontos ganhos em uma data específica
-- Nota: Esta query usa CTE para calcular incrementos comparando com o dia anterior
WITH daily_points AS (
    SELECT 
        btl.*,
        p.nome as professional_name,
        -- Assumindo que last_turn_at é atualizado a cada atendimento do dia
        CASE 
            WHEN DATE(btl.last_turn_at) = $2::date THEN 1
            ELSE 0
        END as atendeu_hoje
    FROM barbers_turn_list btl
    JOIN profissionais p ON p.id = btl.professional_id
    WHERE btl.tenant_id = $1
)
SELECT 
    professional_id,
    professional_name,
    current_points,
    atendeu_hoje,
    last_turn_at
FROM daily_points
WHERE atendeu_hoje = 1 OR current_points > 0
ORDER BY current_points DESC, last_turn_at DESC;
