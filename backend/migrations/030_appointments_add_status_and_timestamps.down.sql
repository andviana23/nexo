-- =====================================================
-- Rollback Migration 030: Remover novos status e timestamps
-- =====================================================

BEGIN;

-- 1. Remover índices
DROP INDEX IF EXISTS idx_appointments_timestamps;
DROP INDEX IF EXISTS idx_appointments_status_tenant;

-- 2. Remover constraint atual
ALTER TABLE appointments
DROP CONSTRAINT IF EXISTS appointments_status_check;

-- 3. Restaurar constraint antiga (apenas status originais)
ALTER TABLE appointments
ADD CONSTRAINT appointments_status_check 
CHECK (
  status IN (
    'CREATED',
    'CONFIRMED',
    'IN_SERVICE',
    'DONE',
    'NO_SHOW',
    'CANCELED'
  )
);

-- 4. Remover colunas de timestamp
ALTER TABLE appointments
DROP COLUMN IF EXISTS checked_in_at,
DROP COLUMN IF EXISTS started_at,
DROP COLUMN IF EXISTS finished_at;

-- 5. Restaurar comentário original
COMMENT ON COLUMN appointments.status IS 'Status: CREATED, CONFIRMED, IN_SERVICE, DONE, NO_SHOW, CANCELED';

COMMIT;
