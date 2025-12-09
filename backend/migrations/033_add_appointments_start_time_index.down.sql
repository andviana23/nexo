-- =====================================================
-- Migration: Rollback do índice de appointments por data
-- =====================================================

BEGIN;

DROP INDEX IF EXISTS idx_appointments_tenant_start_time;
DROP INDEX IF EXISTS idx_appointments_professional_start_time;

COMMIT;
