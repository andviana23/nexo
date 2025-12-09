-- ============================================================================
-- MIGRATION 039: Multi-Unidade - RLS Policies e Função de Contexto
-- Segurança em nível de linha para isolamento multi-tenant e multi-unit
-- ============================================================================

-- ============================================================================
-- 1. FUNÇÃO: set_app_context
-- Define tenant_id e unit_id no contexto da sessão
-- ============================================================================
CREATE OR REPLACE FUNCTION set_app_context(
    p_tenant_id UUID,
    p_unit_id UUID DEFAULT NULL
) RETURNS VOID AS $$
BEGIN
    -- Definir tenant
    PERFORM set_config('app.current_tenant', p_tenant_id::text, true);
    
    -- Definir unit (pode ser NULL para operações tenant-wide)
    IF p_unit_id IS NOT NULL THEN
        PERFORM set_config('app.current_unit', p_unit_id::text, true);
    ELSE
        PERFORM set_config('app.current_unit', '', true);
    END IF;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

COMMENT ON FUNCTION set_app_context IS 'Define tenant_id e unit_id no contexto da sessão para RLS';

-- ============================================================================
-- 2. FUNÇÃO: get_current_tenant
-- Obtém tenant_id do contexto com fallback seguro
-- ============================================================================
CREATE OR REPLACE FUNCTION get_current_tenant() RETURNS UUID AS $$
DECLARE
    v_tenant TEXT;
BEGIN
    v_tenant := current_setting('app.current_tenant', true);
    IF v_tenant IS NULL OR v_tenant = '' THEN
        RETURN NULL;
    END IF;
    RETURN v_tenant::UUID;
EXCEPTION WHEN OTHERS THEN
    RETURN NULL;
END;
$$ LANGUAGE plpgsql STABLE SECURITY DEFINER;

-- ============================================================================
-- 3. FUNÇÃO: get_current_unit
-- Obtém unit_id do contexto com fallback seguro
-- ============================================================================
CREATE OR REPLACE FUNCTION get_current_unit() RETURNS UUID AS $$
DECLARE
    v_unit TEXT;
BEGIN
    v_unit := current_setting('app.current_unit', true);
    IF v_unit IS NULL OR v_unit = '' THEN
        RETURN NULL;
    END IF;
    RETURN v_unit::UUID;
EXCEPTION WHEN OTHERS THEN
    RETURN NULL;
END;
$$ LANGUAGE plpgsql STABLE SECURITY DEFINER;

-- ============================================================================
-- 4. FUNÇÃO: check_tenant_unit_access
-- Verifica se o contexto tem acesso à linha (tenant + unit opcional)
-- ============================================================================
CREATE OR REPLACE FUNCTION check_tenant_unit_access(
    row_tenant_id UUID,
    row_unit_id UUID
) RETURNS BOOLEAN AS $$
DECLARE
    v_current_tenant UUID;
    v_current_unit UUID;
BEGIN
    v_current_tenant := get_current_tenant();
    v_current_unit := get_current_unit();
    
    -- Tenant deve sempre coincidir
    IF row_tenant_id IS DISTINCT FROM v_current_tenant THEN
        RETURN FALSE;
    END IF;
    
    -- Se não há unit no contexto, permite acesso (operação tenant-wide)
    IF v_current_unit IS NULL THEN
        RETURN TRUE;
    END IF;
    
    -- Se a linha não tem unit, permite acesso (recurso compartilhado)
    IF row_unit_id IS NULL THEN
        RETURN TRUE;
    END IF;
    
    -- Unit deve coincidir
    RETURN row_unit_id = v_current_unit;
END;
$$ LANGUAGE plpgsql STABLE SECURITY DEFINER;

-- ============================================================================
-- NOTA SOBRE RLS:
-- A ativação de RLS nas tabelas é opcional e deve ser feita apenas
-- quando o middleware estiver configurado para definir o contexto.
-- Por ora, deixamos as políticas criadas mas não ativamos RLS,
-- permitindo que o código da aplicação faça o filtro explícito.
-- 
-- Para ativar RLS em uma tabela:
--   ALTER TABLE <table> ENABLE ROW LEVEL SECURITY;
--   ALTER TABLE <table> FORCE ROW LEVEL SECURITY;
--
-- Exemplo de policy:
--   CREATE POLICY tenant_unit_policy ON appointments
--     USING (check_tenant_unit_access(tenant_id, unit_id))
--     WITH CHECK (check_tenant_unit_access(tenant_id, unit_id));
-- ============================================================================

-- ============================================================================
-- 5. CRIAR POLICIES (desabilitadas por padrão - ativar quando pronto)
-- ============================================================================

-- Units policy (apenas tenant, sem unit)
DROP POLICY IF EXISTS rls_units_tenant ON units;
CREATE POLICY rls_units_tenant ON units
    USING (tenant_id = get_current_tenant())
    WITH CHECK (tenant_id = get_current_tenant());

-- Appointments policy
DROP POLICY IF EXISTS rls_appointments_tenant_unit ON appointments;
CREATE POLICY rls_appointments_tenant_unit ON appointments
    USING (check_tenant_unit_access(tenant_id, unit_id))
    WITH CHECK (check_tenant_unit_access(tenant_id, unit_id));

-- Commands policy
DROP POLICY IF EXISTS rls_commands_tenant_unit ON commands;
CREATE POLICY rls_commands_tenant_unit ON commands
    USING (check_tenant_unit_access(tenant_id, unit_id))
    WITH CHECK (check_tenant_unit_access(tenant_id, unit_id));

-- Caixa Diário policy
DROP POLICY IF EXISTS rls_caixa_diario_tenant_unit ON caixa_diario;
CREATE POLICY rls_caixa_diario_tenant_unit ON caixa_diario
    USING (check_tenant_unit_access(tenant_id, unit_id))
    WITH CHECK (check_tenant_unit_access(tenant_id, unit_id));

-- Contas a Pagar policy
DROP POLICY IF EXISTS rls_contas_pagar_tenant_unit ON contas_a_pagar;
CREATE POLICY rls_contas_pagar_tenant_unit ON contas_a_pagar
    USING (check_tenant_unit_access(tenant_id, unit_id))
    WITH CHECK (check_tenant_unit_access(tenant_id, unit_id));

-- Contas a Receber policy
DROP POLICY IF EXISTS rls_contas_receber_tenant_unit ON contas_a_receber;
CREATE POLICY rls_contas_receber_tenant_unit ON contas_a_receber
    USING (check_tenant_unit_access(tenant_id, unit_id))
    WITH CHECK (check_tenant_unit_access(tenant_id, unit_id));

-- ============================================================================
-- LOG: Registro de criação das policies
-- ============================================================================
DO $$
BEGIN
    RAISE NOTICE 'RLS policies criadas. Para ativar, execute:';
    RAISE NOTICE 'ALTER TABLE <table> ENABLE ROW LEVEL SECURITY;';
    RAISE NOTICE 'ALTER TABLE <table> FORCE ROW LEVEL SECURITY;';
END $$;
