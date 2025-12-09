-- Migration: 042_commission_rules
-- Description: Cria tabela de regras de comissão flexíveis
-- Author: Nexo Team
-- Date: 2025-12-05

-- ============================================================================
-- TABELA: commission_rules
-- Regras globais/por serviço para cálculo de comissão
-- ============================================================================

CREATE TABLE IF NOT EXISTS commission_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
    
    -- Identificação da regra
    name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- Configuração de comissão
    type VARCHAR(20) NOT NULL DEFAULT 'PERCENTUAL' 
        CHECK (type IN ('PERCENTUAL', 'FIXO')),
    default_rate NUMERIC(10,2) NOT NULL 
        CHECK (
            (type = 'PERCENTUAL' AND default_rate >= 0 AND default_rate <= 100) OR
            (type = 'FIXO' AND default_rate >= 0)
        ),
    
    -- Limites opcionais
    min_amount NUMERIC(15,2),
    max_amount NUMERIC(15,2),
    
    -- Base de cálculo
    calculation_base VARCHAR(20) DEFAULT 'BRUTO'
        CHECK (calculation_base IN ('BRUTO', 'LIQUIDO')),
    
    -- Vigência
    effective_from DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_to DATE,
    
    -- Prioridade para resolução de conflitos
    priority INTEGER DEFAULT 0,
    
    -- Status
    is_active BOOLEAN DEFAULT true NOT NULL,
    
    -- Auditoria
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Constraints
    CONSTRAINT chk_effective_dates CHECK (effective_to IS NULL OR effective_to > effective_from),
    CONSTRAINT chk_min_max CHECK (max_amount IS NULL OR min_amount IS NULL OR max_amount >= min_amount)
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_commission_rules_tenant 
    ON commission_rules(tenant_id);
CREATE INDEX IF NOT EXISTS idx_commission_rules_unit 
    ON commission_rules(tenant_id, unit_id) WHERE unit_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_commission_rules_active 
    ON commission_rules(tenant_id, is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_commission_rules_effective 
    ON commission_rules(tenant_id, effective_from, effective_to);

-- Comentários
COMMENT ON TABLE commission_rules IS 'Regras de comissão globais ou por unidade';
COMMENT ON COLUMN commission_rules.type IS 'Tipo: PERCENTUAL ou FIXO';
COMMENT ON COLUMN commission_rules.default_rate IS 'Taxa padrão (% ou valor fixo)';
COMMENT ON COLUMN commission_rules.calculation_base IS 'Base: BRUTO ou LIQUIDO (desconsiderando taxas)';
COMMENT ON COLUMN commission_rules.priority IS 'Maior prioridade vence em conflitos';
