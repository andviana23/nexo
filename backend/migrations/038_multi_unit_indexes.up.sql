-- ============================================================================
-- MIGRATION 038: Multi-Unidade - Índices Compostos
-- Índices otimizados para queries com tenant_id + unit_id
-- ============================================================================

-- ============================================================================
-- UNITS
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_units_tenant_matriz 
    ON units(tenant_id, is_matriz) WHERE is_matriz = true;

-- ============================================================================
-- USER_UNITS
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_user_units_unit_user 
    ON user_units(unit_id, user_id);

-- ============================================================================
-- APPOINTMENTS
-- ============================================================================
DO $$
BEGIN
    IF to_regclass('public.appointments') IS NOT NULL
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'appointments'
              AND column_name = 'unit_id'
        ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_appointments_tenant_unit_date ON appointments(tenant_id, unit_id, start_time)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_appointments_unit_status ON appointments(unit_id, status) WHERE status NOT IN (''CANCELADO'', ''NO_SHOW'')';
    END IF;
END $$;

-- ============================================================================
-- COMMANDS
-- ============================================================================
-- Compatibilidade: alguns schemas usam `created_at`, outros usam `criado_em`.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'commands'
          AND column_name = 'unit_id'
    ) THEN
        IF EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'commands'
              AND column_name = 'created_at'
        ) THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_commands_tenant_unit_date ON commands(tenant_id, unit_id, created_at)';
        ELSIF EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'commands'
              AND column_name = 'criado_em'
        ) THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_commands_tenant_unit_date ON commands(tenant_id, unit_id, criado_em)';
        END IF;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_commands_unit_status 
    ON commands(unit_id, status);

-- ============================================================================
-- CAIXA_DIARIO
-- ============================================================================
DO $$
BEGIN
    IF to_regclass('public.caixa_diario') IS NOT NULL
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'caixa_diario'
              AND column_name = 'unit_id'
        ) THEN
        IF EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'caixa_diario'
              AND column_name = 'data'
        ) THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_caixa_diario_tenant_unit_date ON caixa_diario(tenant_id, unit_id, data)';
        ELSIF EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'caixa_diario'
              AND column_name = 'data_abertura'
        ) THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_caixa_diario_tenant_unit_date ON caixa_diario(tenant_id, unit_id, data_abertura)';
        ELSIF EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'caixa_diario'
              AND column_name = 'created_at'
        ) THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_caixa_diario_tenant_unit_date ON caixa_diario(tenant_id, unit_id, created_at)';
        END IF;
    END IF;
END $$;

-- ============================================================================
-- CONTAS_A_PAGAR
-- ============================================================================
DO $$
BEGIN
    IF to_regclass('public.contas_a_pagar') IS NOT NULL
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'contas_a_pagar'
              AND column_name = 'unit_id'
        ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_contas_pagar_tenant_unit_vencimento ON contas_a_pagar(tenant_id, unit_id, data_vencimento)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_contas_pagar_unit_status ON contas_a_pagar(unit_id, status)';
    END IF;
END $$;

-- ============================================================================
-- CONTAS_A_RECEBER
-- ============================================================================
DO $$
BEGIN
    IF to_regclass('public.contas_a_receber') IS NOT NULL
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'contas_a_receber'
              AND column_name = 'unit_id'
        ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_contas_receber_tenant_unit_vencimento ON contas_a_receber(tenant_id, unit_id, data_vencimento)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_contas_receber_unit_status ON contas_a_receber(unit_id, status)';
    END IF;
END $$;

-- ============================================================================
-- COMPENSACOES_BANCARIAS
-- ============================================================================
DO $$
BEGIN
    IF to_regclass('public.compensacoes_bancarias') IS NOT NULL
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'compensacoes_bancarias'
              AND column_name = 'unit_id'
        ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_compensacoes_tenant_unit_data ON compensacoes_bancarias(tenant_id, unit_id, data_compensacao)';
    END IF;
END $$;

-- ============================================================================
-- FLUXO_CAIXA_DIARIO
-- ============================================================================
DO $$
BEGIN
    IF to_regclass('public.fluxo_caixa_diario') IS NOT NULL
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'fluxo_caixa_diario'
              AND column_name = 'unit_id'
        ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_fluxo_caixa_tenant_unit_data ON fluxo_caixa_diario(tenant_id, unit_id, data)';
    END IF;
END $$;

-- ============================================================================
-- DRE_MENSAL
-- ============================================================================
DO $$
BEGIN
    IF to_regclass('public.dre_mensal') IS NOT NULL
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'dre_mensal'
              AND column_name = 'unit_id'
        )
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'dre_mensal'
              AND column_name = 'ano'
        )
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'dre_mensal'
              AND column_name = 'mes'
        ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_dre_mensal_tenant_unit_periodo ON dre_mensal(tenant_id, unit_id, ano, mes)';
    END IF;
END $$;

-- ============================================================================
-- MOVIMENTACOES_ESTOQUE
-- ============================================================================
DO $$
BEGIN
    IF to_regclass('public.movimentacoes_estoque') IS NOT NULL
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'movimentacoes_estoque'
              AND column_name = 'unit_id'
        ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_mov_estoque_tenant_unit_data ON movimentacoes_estoque(tenant_id, unit_id, data_movimentacao)';
    END IF;
END $$;

-- ============================================================================
-- BARBER_TURNS
-- ============================================================================
DO $$
BEGIN
    IF to_regclass('public.barber_turns') IS NOT NULL
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'barber_turns'
              AND column_name = 'unit_id'
        )
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'barber_turns'
              AND column_name = 'date'
        ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_barber_turns_tenant_unit_date ON barber_turns(tenant_id, unit_id, date)';
    END IF;
END $$;

-- ============================================================================
-- BLOCKED_TIMES
-- ============================================================================
DO $$
BEGIN
    IF to_regclass('public.blocked_times') IS NOT NULL
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'blocked_times'
              AND column_name = 'unit_id'
        ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_blocked_times_tenant_unit ON blocked_times(tenant_id, unit_id)';
    END IF;
END $$;

-- ============================================================================
-- METAS_MENSAIS
-- ============================================================================
DO $$
BEGIN
    IF to_regclass('public.metas_mensais') IS NOT NULL
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'metas_mensais'
              AND column_name = 'unit_id'
        )
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'metas_mensais'
              AND column_name = 'ano'
        )
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'metas_mensais'
              AND column_name = 'mes'
        ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_metas_mensais_tenant_unit_periodo ON metas_mensais(tenant_id, unit_id, ano, mes)';
    END IF;
END $$;

-- ============================================================================
-- DESPESAS_FIXAS
-- ============================================================================
DO $$
BEGIN
    IF to_regclass('public.despesas_fixas') IS NOT NULL
       AND EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'despesas_fixas'
              AND column_name = 'unit_id'
        ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_despesas_fixas_tenant_unit ON despesas_fixas(tenant_id, unit_id)';
    END IF;
END $$;
