-- ============================================================================
-- MIGRATION 056: Enforce Unit Isolation
-- Objective: Add NOT NULL constraint to unit_id columns on critical tables
-- ============================================================================

DO $$
BEGIN
    -- 1. PROFISSIONAIS
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'profissionais' AND column_name = 'unit_id') THEN
        ALTER TABLE profissionais ALTER COLUMN unit_id SET NOT NULL;
    END IF;

    -- 2. SERVICOS
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'servicos' AND column_name = 'unit_id') THEN
        ALTER TABLE servicos ALTER COLUMN unit_id SET NOT NULL;
    END IF;

    -- 3. CATEGORIAS
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'categorias_servicos' AND column_name = 'unit_id') THEN
        ALTER TABLE categorias_servicos ALTER COLUMN unit_id SET NOT NULL;
    END IF;

    -- 4. APPOINTMENTS
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'appointments' AND column_name = 'unit_id') THEN
        ALTER TABLE appointments ALTER COLUMN unit_id SET NOT NULL;
    END IF;

    -- 5. COMMANDS
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'commands' AND column_name = 'unit_id') THEN
        ALTER TABLE commands ALTER COLUMN unit_id SET NOT NULL;
    END IF;

    -- 6. FINANCEIRO (Exemplo: dre_mensal)
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'dre_mensal' AND column_name = 'unit_id') THEN
        ALTER TABLE dre_mensal ALTER COLUMN unit_id SET NOT NULL;
    END IF;
    
    -- 7. CAIXA DIARIO
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'caixa_diario' AND column_name = 'unit_id') THEN
        ALTER TABLE caixa_diario ALTER COLUMN unit_id SET NOT NULL;
    END IF;

    -- 8. CONTAS A PAGAR/RECEBER
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'contas_a_pagar' AND column_name = 'unit_id') THEN
        ALTER TABLE contas_a_pagar ALTER COLUMN unit_id SET NOT NULL;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'contas_a_receber' AND column_name = 'unit_id') THEN
        ALTER TABLE contas_a_receber ALTER COLUMN unit_id SET NOT NULL;
    END IF;

END $$;
