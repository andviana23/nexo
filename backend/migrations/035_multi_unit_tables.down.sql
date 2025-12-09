-- ============================================================================
-- MIGRATION 035: Multi-Unidade - Rollback
-- ============================================================================

-- Remover feature flag
ALTER TABLE tenants DROP COLUMN IF EXISTS multi_unit_enabled;

-- Remover trigger e função
DROP TRIGGER IF EXISTS trg_ensure_single_default_unit ON user_units;
DROP FUNCTION IF EXISTS ensure_single_default_unit();

-- Remover triggers de updated_at
DROP TRIGGER IF EXISTS trg_user_units_updated_at ON user_units;
DROP TRIGGER IF EXISTS trg_units_updated_at ON units;

-- Remover tabelas (ordem importa por FK)
DROP TABLE IF EXISTS user_units;
DROP TABLE IF EXISTS units;
