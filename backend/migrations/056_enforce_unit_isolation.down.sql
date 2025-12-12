-- ============================================================================
-- MIGRATION 056: Enforce Unit Isolation (Revert)
-- Objective: Remove NOT NULL constraint from unit_id columns
-- ============================================================================

DO $$
BEGIN
    ALTER TABLE IF EXISTS profissionais ALTER COLUMN unit_id DROP NOT NULL;
    ALTER TABLE IF EXISTS servicos ALTER COLUMN unit_id DROP NOT NULL;
    ALTER TABLE IF EXISTS categorias_servicos ALTER COLUMN unit_id DROP NOT NULL;
    ALTER TABLE IF EXISTS appointments ALTER COLUMN unit_id DROP NOT NULL;
    ALTER TABLE IF EXISTS commands ALTER COLUMN unit_id DROP NOT NULL;
    ALTER TABLE IF EXISTS dre_mensal ALTER COLUMN unit_id DROP NOT NULL;
    ALTER TABLE IF EXISTS caixa_diario ALTER COLUMN unit_id DROP NOT NULL;
    ALTER TABLE IF EXISTS contas_a_pagar ALTER COLUMN unit_id DROP NOT NULL;
    ALTER TABLE IF EXISTS contas_a_receber ALTER COLUMN unit_id DROP NOT NULL;
END $$;
