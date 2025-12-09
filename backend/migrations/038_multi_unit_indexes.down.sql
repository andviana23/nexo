-- ============================================================================
-- MIGRATION 038: Multi-Unidade - Rollback Índices
-- ============================================================================

DROP INDEX IF EXISTS idx_despesas_fixas_tenant_unit;
DROP INDEX IF EXISTS idx_metas_mensais_tenant_unit_periodo;
DROP INDEX IF EXISTS idx_blocked_times_tenant_unit;
DROP INDEX IF EXISTS idx_barber_turns_tenant_unit_date;
DROP INDEX IF EXISTS idx_mov_estoque_tenant_unit_data;
DROP INDEX IF EXISTS idx_dre_mensal_tenant_unit_periodo;
DROP INDEX IF EXISTS idx_fluxo_caixa_tenant_unit_data;
DROP INDEX IF EXISTS idx_compensacoes_tenant_unit_data;
DROP INDEX IF EXISTS idx_contas_receber_unit_status;
DROP INDEX IF EXISTS idx_contas_receber_tenant_unit_vencimento;
DROP INDEX IF EXISTS idx_contas_pagar_unit_status;
DROP INDEX IF EXISTS idx_contas_pagar_tenant_unit_vencimento;
DROP INDEX IF EXISTS idx_caixa_diario_tenant_unit_date;
DROP INDEX IF EXISTS idx_commands_unit_status;
DROP INDEX IF EXISTS idx_commands_tenant_unit_date;
DROP INDEX IF EXISTS idx_appointments_unit_status;
DROP INDEX IF EXISTS idx_appointments_tenant_unit_date;
DROP INDEX IF EXISTS idx_user_units_unit_user;
DROP INDEX IF EXISTS idx_units_tenant_matriz;
