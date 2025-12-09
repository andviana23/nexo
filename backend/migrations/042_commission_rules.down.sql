-- Migration: 042_commission_rules (rollback)
-- Description: Remove tabela de regras de comissão

DROP INDEX IF EXISTS idx_commission_rules_effective;
DROP INDEX IF EXISTS idx_commission_rules_active;
DROP INDEX IF EXISTS idx_commission_rules_unit;
DROP INDEX IF EXISTS idx_commission_rules_tenant;
DROP TABLE IF EXISTS commission_rules;
