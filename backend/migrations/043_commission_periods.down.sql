-- Migration: 043_commission_periods (rollback)
-- Description: Remove tabela de períodos de comissão

DROP INDEX IF EXISTS idx_commission_periods_open;
DROP INDEX IF EXISTS idx_commission_periods_status;
DROP INDEX IF EXISTS idx_commission_periods_month;
DROP INDEX IF EXISTS idx_commission_periods_professional;
DROP INDEX IF EXISTS idx_commission_periods_unit;
DROP INDEX IF EXISTS idx_commission_periods_tenant;
DROP TABLE IF EXISTS commission_periods;
