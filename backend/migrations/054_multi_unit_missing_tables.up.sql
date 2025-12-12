-- ============================================================================
-- MIGRATION 054: Multi-Unidade - Tabelas Faltantes
-- Adiciona unit_id em tabelas que não foram cobertas na migration 036
-- ============================================================================

DO $$
BEGIN
    -- 1. PROFISSIONAIS (Professionals)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'profissionais') THEN
        ALTER TABLE profissionais ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
        -- Cria índice se não existir
        IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_profissionais_unit') THEN
             CREATE INDEX idx_profissionais_unit ON profissionais(tenant_id, unit_id);
        END IF;
    END IF;

    -- 2. SERVICOS (Services)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'servicos') THEN
        ALTER TABLE servicos ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
        IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_servicos_unit') THEN
             CREATE INDEX idx_servicos_unit ON servicos(tenant_id, unit_id);
        END IF;
    END IF;

    -- 3. CATEGORIAS_SERVICOS (Service Categories)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'categorias_servicos') THEN
        ALTER TABLE categorias_servicos ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
        IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_categorias_servicos_unit') THEN
             CREATE INDEX idx_categorias_servicos_unit ON categorias_servicos(tenant_id, unit_id);
        END IF;
    END IF;

    -- 4. ASSINATURAS (Subscriptions)
    -- Verificar nome da tabela (subscriptions ou assinaturas)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'subscriptions') THEN
        ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
        IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_subscriptions_unit') THEN
             CREATE INDEX idx_subscriptions_unit ON subscriptions(tenant_id, unit_id);
        END IF;
    ELSIF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'assinaturas') THEN
        ALTER TABLE assinaturas ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
        IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_assinaturas_unit') THEN
             CREATE INDEX idx_assinaturas_unit ON assinaturas(tenant_id, unit_id);
        END IF;
    END IF;

END $$;
