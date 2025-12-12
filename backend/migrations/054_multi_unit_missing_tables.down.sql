-- ============================================================================
-- MIGRATION 054 DOWN: Remove unit_id das tabelas faltantes
-- ============================================================================

ALTER TABLE IF EXISTS profissionais DROP COLUMN IF EXISTS unit_id;
ALTER TABLE IF EXISTS servicos DROP COLUMN IF EXISTS unit_id;
ALTER TABLE IF EXISTS categorias_servicos DROP COLUMN IF EXISTS unit_id;
ALTER TABLE IF EXISTS subscriptions DROP COLUMN IF EXISTS unit_id;
ALTER TABLE IF EXISTS assinaturas DROP COLUMN IF EXISTS unit_id;
