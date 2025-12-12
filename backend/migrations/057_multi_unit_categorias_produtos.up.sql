-- ============================================================================
-- MIGRATION 057: Multi-Unidade - Adicionar unit_id em categorias_produtos
-- Objetivo:
-- - Permitir categorias de produtos isoladas por unidade (tenant + unit)
-- - Manter compatibilidade com dados legados (unit_id NULL)
--
-- Estratégia:
-- - Adiciona coluna unit_id (nullable)
-- - Substitui UNIQUE(tenant_id, nome) por índices únicos parciais:
--   * Para registros com unit_id IS NULL: único por tenant + nome
--   * Para registros com unit_id IS NOT NULL: único por tenant + unit + nome
-- ============================================================================

DO $$
BEGIN
    -- 1) Coluna unit_id (nullable para compatibilidade)
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'categorias_produtos'
          AND column_name = 'unit_id'
    ) THEN
        ALTER TABLE categorias_produtos
            ADD COLUMN unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 2) Remover constraint antiga (tenant_id, nome)
    BEGIN
        ALTER TABLE categorias_produtos
            DROP CONSTRAINT IF EXISTS categorias_produtos_tenant_nome_unique;
    EXCEPTION
        WHEN undefined_object THEN
            -- ignore
            NULL;
    END;

    -- 3) Índices de apoio
    CREATE INDEX IF NOT EXISTS idx_categorias_produtos_tenant_unit
        ON categorias_produtos(tenant_id, unit_id);

    -- 4) Unicidade por tenant para legados (unit_id NULL)
    CREATE UNIQUE INDEX IF NOT EXISTS uq_categorias_produtos_tenant_nome_unit_null
        ON categorias_produtos(tenant_id, LOWER(nome))
        WHERE unit_id IS NULL;

    -- 5) Unicidade por tenant+unit para registros novos (unit_id NOT NULL)
    CREATE UNIQUE INDEX IF NOT EXISTS uq_categorias_produtos_tenant_unit_nome
        ON categorias_produtos(tenant_id, unit_id, LOWER(nome))
        WHERE unit_id IS NOT NULL;
END $$;
