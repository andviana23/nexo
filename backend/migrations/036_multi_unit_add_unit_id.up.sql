-- ============================================================================
-- MIGRATION 036: Multi-Unidade - Propagação de unit_id nas tabelas operacionais
-- Adiciona coluna unit_id em todas as tabelas que precisam de isolamento por unidade
-- SAFE: Verifica existência de cada tabela antes de alterar
-- ============================================================================

DO $$
BEGIN
    -- 1. APPOINTMENTS (Agendamentos)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'appointments') THEN
        ALTER TABLE appointments ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 2. APPOINTMENT_STATUS_HISTORY (Histórico de Status)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'appointment_status_history') THEN
        ALTER TABLE appointment_status_history ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 3. COMMANDS (Comandas)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'commands') THEN
        ALTER TABLE commands ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 4. CAIXA_DIARIO (Caixa Diário)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'caixa_diario') THEN
        ALTER TABLE caixa_diario ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 5. CONTAS_A_PAGAR (Contas a Pagar)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'contas_a_pagar') THEN
        ALTER TABLE contas_a_pagar ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 6. CONTAS_A_RECEBER (Contas a Receber)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'contas_a_receber') THEN
        ALTER TABLE contas_a_receber ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 7. COMPENSACOES_BANCARIAS (Compensações Bancárias)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'compensacoes_bancarias') THEN
        ALTER TABLE compensacoes_bancarias ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 8. FLUXO_CAIXA_DIARIO (Fluxo de Caixa Diário)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'fluxo_caixa_diario') THEN
        ALTER TABLE fluxo_caixa_diario ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 9. DRE_MENSAL (DRE Mensal)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'dre_mensal') THEN
        ALTER TABLE dre_mensal ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 10. MOVIMENTACOES_ESTOQUE (Movimentações de Estoque)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'movimentacoes_estoque') THEN
        ALTER TABLE movimentacoes_estoque ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 11. PRODUTOS (Estoque)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'produtos') THEN
        ALTER TABLE produtos ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 12. BARBER_TURNS (Turnos de Barbeiros)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'barber_turns') THEN
        ALTER TABLE barber_turns ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 13. BLOCKED_TIMES (Horários Bloqueados)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'blocked_times') THEN
        ALTER TABLE blocked_times ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 14. METAS_MENSAIS (Metas Mensais)
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'metas_mensais') THEN
        ALTER TABLE metas_mensais ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;

    -- 15. DESPESAS_FIXAS
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'despesas_fixas') THEN
        ALTER TABLE despesas_fixas ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id) ON DELETE RESTRICT;
    END IF;
END $$;
