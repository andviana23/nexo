-- ============================================================================
-- MIGRATION 035: Multi-Unidade - Tabelas Base
-- Suporte a múltiplas unidades (filiais) por tenant
-- ============================================================================

-- ============================================================================
-- 1. TABELA: units (Unidades/Filiais)
-- ============================================================================
CREATE TABLE IF NOT EXISTS units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Dados básicos
    nome VARCHAR(100) NOT NULL,
    apelido VARCHAR(50),
    descricao TEXT,
    
    -- Endereço resumido
    endereco_resumo VARCHAR(255),
    cidade VARCHAR(100),
    estado VARCHAR(2),
    
    -- Configurações
    timezone VARCHAR(50) DEFAULT 'America/Sao_Paulo' NOT NULL,
    ativa BOOLEAN DEFAULT true NOT NULL,
    is_matriz BOOLEAN DEFAULT false NOT NULL,
    
    -- Timestamps
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- Constraints
    CONSTRAINT units_tenant_nome_unique UNIQUE (tenant_id, nome)
);

-- Comentários
COMMENT ON TABLE units IS 'Unidades/Filiais do tenant - cada barbearia pode ter múltiplas unidades';
COMMENT ON COLUMN units.nome IS 'Nome oficial da unidade (único por tenant)';
COMMENT ON COLUMN units.apelido IS 'Nome curto para exibição (ex: "Centro", "Shopping")';
COMMENT ON COLUMN units.is_matriz IS 'Se true, é a unidade principal/matriz do tenant';
COMMENT ON COLUMN units.timezone IS 'Fuso horário da unidade para agendamentos';

-- Índices
CREATE INDEX IF NOT EXISTS idx_units_tenant ON units(tenant_id);
CREATE INDEX IF NOT EXISTS idx_units_tenant_ativa ON units(tenant_id, ativa);

-- Trigger para atualizar updated_at (drop first for idempotency)
DROP TRIGGER IF EXISTS trg_units_updated_at ON units;
CREATE TRIGGER trg_units_updated_at
    BEFORE UPDATE ON units
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- 2. TABELA: user_units (Vínculo Usuário-Unidade)
-- ============================================================================
CREATE TABLE IF NOT EXISTS user_units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    
    -- Configurações do vínculo
    is_default BOOLEAN DEFAULT false NOT NULL,
    role_override VARCHAR(50), -- Permite papel diferente por unidade (null = usa role padrão)
    
    -- Timestamps
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- Constraints
    CONSTRAINT user_units_unique UNIQUE (user_id, unit_id)
);

-- Comentários
COMMENT ON TABLE user_units IS 'Vínculo entre usuários e unidades - define acesso e unidade padrão';
COMMENT ON COLUMN user_units.is_default IS 'Se true, é a unidade padrão do usuário no login';
COMMENT ON COLUMN user_units.role_override IS 'Papel específico nesta unidade (null = usa role do user)';

-- Índices
CREATE INDEX IF NOT EXISTS idx_user_units_user ON user_units(user_id);
CREATE INDEX IF NOT EXISTS idx_user_units_unit ON user_units(unit_id);
CREATE INDEX IF NOT EXISTS idx_user_units_user_default ON user_units(user_id, is_default) WHERE is_default = true;

-- Trigger para atualizar updated_at (drop first for idempotency)
DROP TRIGGER IF EXISTS trg_user_units_updated_at ON user_units;
CREATE TRIGGER trg_user_units_updated_at
    BEFORE UPDATE ON user_units
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- 3. FUNÇÃO: Garantir apenas uma unidade padrão por usuário
-- ============================================================================
CREATE OR REPLACE FUNCTION ensure_single_default_unit()
RETURNS TRIGGER AS $$
BEGIN
    -- Se está definindo como padrão, remove padrão das outras
    IF NEW.is_default = true THEN
        UPDATE user_units 
        SET is_default = false
        WHERE user_id = NEW.user_id 
          AND id != NEW.id 
          AND is_default = true;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_ensure_single_default_unit ON user_units;
CREATE TRIGGER trg_ensure_single_default_unit
    BEFORE INSERT OR UPDATE OF is_default ON user_units
    FOR EACH ROW
    WHEN (NEW.is_default = true)
    EXECUTE FUNCTION ensure_single_default_unit();

-- ============================================================================
-- 4. ADICIONAR COLUNA multi_unit_enabled EM tenants (feature flag)
-- ============================================================================
ALTER TABLE tenants 
ADD COLUMN IF NOT EXISTS multi_unit_enabled BOOLEAN DEFAULT false NOT NULL;

COMMENT ON COLUMN tenants.multi_unit_enabled IS 'Feature flag: se true, tenant tem acesso a multi-unidade';
