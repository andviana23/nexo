-- =====================================================
-- Migration: Adicionar novos status e timestamps para appointments
-- Descrição: 
--   - Adiciona status CHECKED_IN e AWAITING_PAYMENT
--   - Adiciona timestamps: checked_in_at, started_at, finished_at
--   - Permite rastreamento completo do fluxo de agendamento
-- Data: 2025-01-XX
-- =====================================================

BEGIN;

-- 1. Adicionar novas colunas de timestamp
ALTER TABLE appointments
ADD COLUMN IF NOT EXISTS checked_in_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS finished_at TIMESTAMPTZ;

COMMENT ON COLUMN appointments.checked_in_at IS 'Momento em que o cliente fez check-in (chegou para o atendimento)';
COMMENT ON COLUMN appointments.started_at IS 'Momento em que o atendimento iniciou (profissional começou o serviço)';
COMMENT ON COLUMN appointments.finished_at IS 'Momento em que o atendimento foi concluído (serviços finalizados)';

-- 2. Remover constraint antiga de status
ALTER TABLE appointments
DROP CONSTRAINT IF EXISTS appointments_status_check;

-- 3. Adicionar nova constraint com os novos status
ALTER TABLE appointments
ADD CONSTRAINT appointments_status_check 
CHECK (
  status IN (
    'CREATED',           -- Agendamento criado
    'CONFIRMED',         -- Agendamento confirmado
    'CHECKED_IN',        -- Cliente chegou (novo)
    'IN_SERVICE',        -- Atendimento em andamento
    'AWAITING_PAYMENT',  -- Aguardando pagamento (novo)
    'DONE',              -- Concluído e pago
    'NO_SHOW',           -- Cliente não compareceu
    'CANCELED'           -- Cancelado
  )
);

-- 4. Atualizar comentário da coluna status com os novos valores
COMMENT ON COLUMN appointments.status IS 'Status: CREATED, CONFIRMED, CHECKED_IN, IN_SERVICE, AWAITING_PAYMENT, DONE, NO_SHOW, CANCELED';

-- 5. Criar índice para melhorar performance de consultas por status
CREATE INDEX IF NOT EXISTS idx_appointments_status_tenant 
ON appointments(tenant_id, status) 
WHERE status NOT IN ('DONE', 'CANCELED', 'NO_SHOW');

COMMENT ON INDEX idx_appointments_status_tenant IS 'Índice para buscar agendamentos ativos por tenant e status';

-- 6. Criar índice para timestamps (usado em relatórios e métricas)
CREATE INDEX IF NOT EXISTS idx_appointments_timestamps 
ON appointments(tenant_id, started_at, finished_at) 
WHERE started_at IS NOT NULL;

COMMENT ON INDEX idx_appointments_timestamps IS 'Índice para consultas de agendamentos por período de execução';

COMMIT;
