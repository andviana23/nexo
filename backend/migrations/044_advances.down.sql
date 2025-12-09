-- Migration: 044_advances (rollback)
-- Description: Remove tabela de adiantamentos

DROP INDEX IF EXISTS idx_advances_request_date;
DROP INDEX IF EXISTS idx_advances_approved_not_deducted;
DROP INDEX IF EXISTS idx_advances_pending;
DROP INDEX IF EXISTS idx_advances_status;
DROP INDEX IF EXISTS idx_advances_professional;
DROP INDEX IF EXISTS idx_advances_unit;
DROP INDEX IF EXISTS idx_advances_tenant;
DROP TABLE IF EXISTS advances;
