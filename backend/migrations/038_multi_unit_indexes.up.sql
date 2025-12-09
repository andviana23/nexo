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
CREATE INDEX IF NOT EXISTS idx_appointments_tenant_unit_date 
    ON appointments(tenant_id, unit_id, start_time);
CREATE INDEX IF NOT EXISTS idx_appointments_unit_status 
    ON appointments(unit_id, status) WHERE status NOT IN ('CANCELADO', 'NO_SHOW');

-- ============================================================================
-- COMMANDS
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_commands_tenant_unit_date 
    ON commands(tenant_id, unit_id, created_at);
CREATE INDEX IF NOT EXISTS idx_commands_unit_status 
    ON commands(unit_id, status);

-- ============================================================================
-- CAIXA_DIARIO
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_caixa_diario_tenant_unit_date 
    ON caixa_diario(tenant_id, unit_id, data);

-- ============================================================================
-- CONTAS_A_PAGAR
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_contas_pagar_tenant_unit_vencimento 
    ON contas_a_pagar(tenant_id, unit_id, data_vencimento);
CREATE INDEX IF NOT EXISTS idx_contas_pagar_unit_status 
    ON contas_a_pagar(unit_id, status);

-- ============================================================================
-- CONTAS_A_RECEBER
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_contas_receber_tenant_unit_vencimento 
    ON contas_a_receber(tenant_id, unit_id, data_vencimento);
CREATE INDEX IF NOT EXISTS idx_contas_receber_unit_status 
    ON contas_a_receber(unit_id, status);

-- ============================================================================
-- COMPENSACOES_BANCARIAS
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_compensacoes_tenant_unit_data 
    ON compensacoes_bancarias(tenant_id, unit_id, data_compensacao);

-- ============================================================================
-- FLUXO_CAIXA_DIARIO
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_fluxo_caixa_tenant_unit_data 
    ON fluxo_caixa_diario(tenant_id, unit_id, data);

-- ============================================================================
-- DRE_MENSAL
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_dre_mensal_tenant_unit_periodo 
    ON dre_mensal(tenant_id, unit_id, ano, mes);

-- ============================================================================
-- MOVIMENTACOES_ESTOQUE
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_mov_estoque_tenant_unit_data 
    ON movimentacoes_estoque(tenant_id, unit_id, data_movimentacao);

-- ============================================================================
-- BARBER_TURNS
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_barber_turns_tenant_unit_date 
    ON barber_turns(tenant_id, unit_id, date);

-- ============================================================================
-- BLOCKED_TIMES
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_blocked_times_tenant_unit 
    ON blocked_times(tenant_id, unit_id);

-- ============================================================================
-- METAS_MENSAIS
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_metas_mensais_tenant_unit_periodo 
    ON metas_mensais(tenant_id, unit_id, ano, mes);

-- ============================================================================
-- DESPESAS_FIXAS
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_despesas_fixas_tenant_unit 
    ON despesas_fixas(tenant_id, unit_id);
