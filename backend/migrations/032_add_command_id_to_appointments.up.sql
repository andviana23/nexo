-- =====================================================
-- Migration: Adicionar command_id em appointments
-- Descrição: 
--   - Adiciona coluna command_id em appointments
--   - Permite vincular agendamento com comanda
--   - Necessário para fluxo de pagamento Trinks
-- Data: 30/11/2025
-- =====================================================

BEGIN;

-- 1. Adicionar coluna command_id
ALTER TABLE appointments
ADD COLUMN IF NOT EXISTS command_id UUID REFERENCES commands(id) ON DELETE SET NULL;

-- 2. Comentário na coluna
COMMENT ON COLUMN appointments.command_id IS 'Referência para a comanda vinculada ao agendamento (quando status = AWAITING_PAYMENT)';

-- 3. Criar índice para melhorar performance de consultas
CREATE INDEX IF NOT EXISTS idx_appointments_command_id 
ON appointments(command_id) 
WHERE command_id IS NOT NULL;

COMMENT ON INDEX idx_appointments_command_id IS 'Índice para buscar agendamentos por comanda';

COMMIT;
