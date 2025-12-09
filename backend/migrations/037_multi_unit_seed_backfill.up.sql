-- ============================================================================
-- MIGRATION 037: Multi-Unidade - Seed/Backfill Unidade Principal
-- Para cada tenant existente, cria unidade "Principal" e vincula usuários
-- SAFE: Verifica existência de cada tabela antes de alterar
-- ============================================================================

-- ============================================================================
-- 1. CRIAR UNIDADE "PRINCIPAL" PARA CADA TENANT EXISTENTE
-- ============================================================================
INSERT INTO units (tenant_id, nome, apelido, is_matriz, ativa, descricao)
SELECT 
    t.id as tenant_id,
    'Principal' as nome,
    'Matriz' as apelido,
    true as is_matriz,
    true as ativa,
    'Unidade principal criada automaticamente na migração multi-unidade' as descricao
FROM tenants t
WHERE NOT EXISTS (
    SELECT 1 FROM units u WHERE u.tenant_id = t.id
);

-- ============================================================================
-- 2. VINCULAR TODOS OS USUÁRIOS À UNIDADE PRINCIPAL DO SEU TENANT
-- ============================================================================
INSERT INTO user_units (user_id, unit_id, is_default)
SELECT 
    u.id as user_id,
    un.id as unit_id,
    true as is_default
FROM users u
JOIN units un ON un.tenant_id = u.tenant_id AND un.is_matriz = true
WHERE NOT EXISTS (
    SELECT 1 FROM user_units uu WHERE uu.user_id = u.id AND uu.unit_id = un.id
);

-- ============================================================================
-- 3. BACKFILL: Atualizar unit_id nas tabelas operacionais existentes
-- SAFE: Verifica existência da tabela e coluna antes de atualizar
-- ============================================================================

DO $$
BEGIN
    -- Appointments
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'appointments' AND column_name = 'unit_id') THEN
        UPDATE appointments a SET unit_id = (SELECT u.id FROM units u WHERE u.tenant_id = a.tenant_id AND u.is_matriz = true LIMIT 1) WHERE a.unit_id IS NULL;
    END IF;

    -- Commands
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'commands' AND column_name = 'unit_id') THEN
        UPDATE commands c SET unit_id = (SELECT u.id FROM units u WHERE u.tenant_id = c.tenant_id AND u.is_matriz = true LIMIT 1) WHERE c.unit_id IS NULL;
    END IF;

    -- Caixa Diário
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'caixa_diario' AND column_name = 'unit_id') THEN
        UPDATE caixa_diario cd SET unit_id = (SELECT u.id FROM units u WHERE u.tenant_id = cd.tenant_id AND u.is_matriz = true LIMIT 1) WHERE cd.unit_id IS NULL;
    END IF;

    -- Contas a Pagar
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'contas_a_pagar' AND column_name = 'unit_id') THEN
        UPDATE contas_a_pagar cp SET unit_id = (SELECT u.id FROM units u WHERE u.tenant_id = cp.tenant_id AND u.is_matriz = true LIMIT 1) WHERE cp.unit_id IS NULL;
    END IF;

    -- Contas a Receber
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'contas_a_receber' AND column_name = 'unit_id') THEN
        UPDATE contas_a_receber cr SET unit_id = (SELECT u.id FROM units u WHERE u.tenant_id = cr.tenant_id AND u.is_matriz = true LIMIT 1) WHERE cr.unit_id IS NULL;
    END IF;

    -- Compensações Bancárias
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'compensacoes_bancarias' AND column_name = 'unit_id') THEN
        UPDATE compensacoes_bancarias cb SET unit_id = (SELECT u.id FROM units u WHERE u.tenant_id = cb.tenant_id AND u.is_matriz = true LIMIT 1) WHERE cb.unit_id IS NULL;
    END IF;

    -- Fluxo Caixa Diário
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'fluxo_caixa_diario' AND column_name = 'unit_id') THEN
        UPDATE fluxo_caixa_diario fcd SET unit_id = (SELECT u.id FROM units u WHERE u.tenant_id = fcd.tenant_id AND u.is_matriz = true LIMIT 1) WHERE fcd.unit_id IS NULL;
    END IF;

    -- DRE Mensal
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'dre_mensal' AND column_name = 'unit_id') THEN
        UPDATE dre_mensal dm SET unit_id = (SELECT u.id FROM units u WHERE u.tenant_id = dm.tenant_id AND u.is_matriz = true LIMIT 1) WHERE dm.unit_id IS NULL;
    END IF;

    -- Movimentações Estoque
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'movimentacoes_estoque' AND column_name = 'unit_id') THEN
        UPDATE movimentacoes_estoque me SET unit_id = (SELECT u.id FROM units u WHERE u.tenant_id = me.tenant_id AND u.is_matriz = true LIMIT 1) WHERE me.unit_id IS NULL;
    END IF;

    -- Metas Mensais
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'metas_mensais' AND column_name = 'unit_id') THEN
        UPDATE metas_mensais mm SET unit_id = (SELECT u.id FROM units u WHERE u.tenant_id = mm.tenant_id AND u.is_matriz = true LIMIT 1) WHERE mm.unit_id IS NULL;
    END IF;

    -- Despesas Fixas
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'despesas_fixas' AND column_name = 'unit_id') THEN
        UPDATE despesas_fixas df SET unit_id = (SELECT u.id FROM units u WHERE u.tenant_id = df.tenant_id AND u.is_matriz = true LIMIT 1) WHERE df.unit_id IS NULL;
    END IF;
END $$;
