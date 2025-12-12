-- ============================================================================
-- MIGRATION 057 DOWN: Reverter unit_id em categorias_produtos
-- ============================================================================

DO $$
BEGIN
    -- Remover índices novos
    DROP INDEX IF EXISTS uq_categorias_produtos_tenant_unit_nome;
    DROP INDEX IF EXISTS uq_categorias_produtos_tenant_nome_unit_null;
    DROP INDEX IF EXISTS idx_categorias_produtos_tenant_unit;

    -- Remover coluna unit_id
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'categorias_produtos'
          AND column_name = 'unit_id'
    ) THEN
        ALTER TABLE categorias_produtos
            DROP COLUMN unit_id;
    END IF;

    -- Restaurar constraint antiga (tenant_id, nome)
    -- (não recria se já existir; preserva idempotência)
    ALTER TABLE categorias_produtos
        ADD CONSTRAINT categorias_produtos_tenant_nome_unique UNIQUE (tenant_id, nome);
EXCEPTION
    WHEN duplicate_object THEN
        NULL;
END $$;
