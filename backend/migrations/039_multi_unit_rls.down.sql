-- ============================================================================
-- MIGRATION 039: Multi-Unidade - Rollback RLS
-- ============================================================================

-- Remover policies
DROP POLICY IF EXISTS rls_contas_receber_tenant_unit ON contas_a_receber;
DROP POLICY IF EXISTS rls_contas_pagar_tenant_unit ON contas_a_pagar;
DROP POLICY IF EXISTS rls_caixa_diario_tenant_unit ON caixa_diario;
DROP POLICY IF EXISTS rls_commands_tenant_unit ON commands;
DROP POLICY IF EXISTS rls_appointments_tenant_unit ON appointments;
DROP POLICY IF EXISTS rls_units_tenant ON units;

-- Remover funções
DROP FUNCTION IF EXISTS check_tenant_unit_access(UUID, UUID);
DROP FUNCTION IF EXISTS get_current_unit();
DROP FUNCTION IF EXISTS get_current_tenant();
DROP FUNCTION IF EXISTS set_app_context(UUID, UUID);
