-- Migration: 045_commission_items
-- Description: Cria tabela de itens de comissão (cada atendimento que gera comissão)
-- Author: Nexo Team
-- Date: 2025-12-05

-- ============================================================================
-- TABELA: commission_items
-- Itens individuais de comissão (vinculados a comandas/atendimentos)
-- ============================================================================

CREATE TABLE IF NOT EXISTS commission_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
    
    -- Profissional
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    
    -- Origem do item (comanda/command_item)
    command_id UUID REFERENCES commands(id) ON DELETE CASCADE,
    command_item_id UUID REFERENCES command_items(id) ON DELETE CASCADE,
    appointment_id UUID REFERENCES appointments(id) ON DELETE SET NULL,
    
    -- Serviço/Produto
    service_id UUID REFERENCES servicos(id) ON DELETE SET NULL,
    service_name VARCHAR(255),
    
    -- Valores
    gross_value NUMERIC(15,2) NOT NULL CHECK (gross_value >= 0),     -- Valor bruto do serviço
    commission_rate NUMERIC(10,2) NOT NULL,                           -- Taxa aplicada (% ou valor fixo)
    commission_type VARCHAR(20) NOT NULL DEFAULT 'PERCENTUAL'         -- Tipo: PERCENTUAL ou FIXO
        CHECK (commission_type IN ('PERCENTUAL', 'FIXO')),
    commission_value NUMERIC(15,2) NOT NULL CHECK (commission_value >= 0), -- Valor calculado
    
    -- Fonte da regra aplicada (para auditoria)
    commission_source VARCHAR(20) NOT NULL DEFAULT 'PROFISSIONAL'
        CHECK (commission_source IN ('SERVICO', 'PROFISSIONAL', 'REGRA', 'MANUAL')),
    rule_id UUID REFERENCES commission_rules(id) ON DELETE SET NULL,
    
    -- Data de competência
    reference_date DATE NOT NULL,
    
    -- Descrição/observação
    description TEXT,
    
    -- Status do item
    status VARCHAR(20) DEFAULT 'PENDENTE' NOT NULL
        CHECK (status IN ('PENDENTE', 'PROCESSADO', 'PAGO', 'CANCELADO', 'ESTORNADO')),
    
    -- Período de fechamento (quando processado)
    period_id UUID REFERENCES commission_periods(id) ON DELETE SET NULL,
    
    -- Auditoria
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    processed_at TIMESTAMPTZ,
    
    -- Constraints
    CONSTRAINT chk_commission_value_valid
        CHECK (
            (commission_type = 'PERCENTUAL' AND commission_rate >= 0 AND commission_rate <= 100) OR
            (commission_type = 'FIXO' AND commission_rate >= 0)
        )
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_commission_items_tenant 
    ON commission_items(tenant_id);
CREATE INDEX IF NOT EXISTS idx_commission_items_unit 
    ON commission_items(tenant_id, unit_id) WHERE unit_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_commission_items_professional 
    ON commission_items(tenant_id, professional_id);
CREATE INDEX IF NOT EXISTS idx_commission_items_period 
    ON commission_items(tenant_id, period_id) WHERE period_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_commission_items_status 
    ON commission_items(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_commission_items_pending 
    ON commission_items(tenant_id, professional_id, status) WHERE status = 'PENDENTE';
CREATE INDEX IF NOT EXISTS idx_commission_items_reference_date 
    ON commission_items(tenant_id, reference_date);
CREATE INDEX IF NOT EXISTS idx_commission_items_command 
    ON commission_items(tenant_id, command_id) WHERE command_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_commission_items_command_item 
    ON commission_items(command_item_id) WHERE command_item_id IS NOT NULL;

-- Índice único para evitar duplicação de comissão por item de comanda
CREATE UNIQUE INDEX IF NOT EXISTS idx_commission_items_command_item_unique
    ON commission_items(tenant_id, command_item_id) WHERE command_item_id IS NOT NULL;

-- Comentários
COMMENT ON TABLE commission_items IS 'Itens individuais de comissão gerados por atendimentos';
COMMENT ON COLUMN commission_items.commission_source IS 'Origem da taxa: SERVICO, PROFISSIONAL, REGRA ou MANUAL';
COMMENT ON COLUMN commission_items.reference_date IS 'Data de competência para fechamento';
COMMENT ON COLUMN commission_items.period_id IS 'Período de fechamento quando status = PROCESSADO ou PAGO';
