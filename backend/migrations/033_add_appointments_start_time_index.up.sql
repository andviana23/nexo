-- =====================================================
-- Migration: Adicionar índice para filtro de appointments por data
-- Descrição: 
--   - Cria índice composto (tenant_id, start_time) para otimizar
--     consultas de listagem por range de datas no calendário
-- Data: 2025-12-01
-- =====================================================

BEGIN;

-- Índice composto para consultas de agendamentos por período
-- Usado principalmente pelo calendário semanal/mensal
CREATE INDEX IF NOT EXISTS idx_appointments_tenant_start_time 
ON appointments(tenant_id, start_time);

COMMENT ON INDEX idx_appointments_tenant_start_time IS 'Índice para consultas de agendamentos por tenant e período (start_date/end_date)';

-- Índice composto para consultas por profissional e data
-- Usado para verificar disponibilidade e listar agenda de um profissional
CREATE INDEX IF NOT EXISTS idx_appointments_professional_start_time 
ON appointments(professional_id, start_time)
WHERE status NOT IN ('CANCELED', 'NO_SHOW');

COMMENT ON INDEX idx_appointments_professional_start_time IS 'Índice para consultas de agendamentos ativos por profissional e data';

COMMIT;
