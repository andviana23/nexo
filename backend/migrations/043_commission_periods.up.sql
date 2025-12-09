-- Migration: 043_commission_periods
-- Description: Cria tabela de períodos/folhas de comissão para fechamento mensal
-- Author: Nexo Team
-- Date: 2025-12-05

-- ============================================================================
-- TABELA: commission_periods
-- Períodos de fechamento de comissões (folhas mensais)
-- ============================================================================

CREATE TABLE IF NOT EXISTS commission_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
    
    -- Período
    reference_month VARCHAR(7) NOT NULL, -- formato: YYYY-MM
    
    -- Profissional (se NULL, é período global)
    professional_id UUID REFERENCES profissionais(id) ON DELETE CASCADE,
    
    -- Totais calculados
    total_gross NUMERIC(15,2) DEFAULT 0 NOT NULL,          -- Total bruto de atendimentos
    total_commission NUMERIC(15,2) DEFAULT 0 NOT NULL,     -- Total de comissões
    total_advances NUMERIC(15,2) DEFAULT 0 NOT NULL,       -- Total de adiantamentos deduzidos
    total_adjustments NUMERIC(15,2) DEFAULT 0 NOT NULL,    -- Ajustes manuais (+/-)
    total_net NUMERIC(15,2) DEFAULT 0 NOT NULL,            -- Líquido a pagar
    
    -- Contadores
    items_count INTEGER DEFAULT 0 NOT NULL,                -- Qtd de itens de comissão
    
    -- Status do período
    status VARCHAR(20) DEFAULT 'ABERTO' NOT NULL
        CHECK (status IN ('ABERTO', 'PROCESSANDO', 'FECHADO', 'PAGO', 'CANCELADO')),
    
    -- Datas importantes
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    closed_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    
    -- Integração com contas a pagar
    conta_pagar_id UUID REFERENCES contas_a_pagar(id) ON DELETE SET NULL,
    
    -- Quem fechou/pagou
    closed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    paid_by UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Observações
    notes TEXT,
    
    -- Auditoria
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    
    -- Constraints
    CONSTRAINT uq_commission_period_professional_month 
        UNIQUE (tenant_id, unit_id, professional_id, reference_month),
    CONSTRAINT chk_period_dates 
        CHECK (period_end >= period_start),
    CONSTRAINT chk_totals_positive
        CHECK (total_gross >= 0 AND total_commission >= 0 AND total_advances >= 0)
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_commission_periods_tenant 
    ON commission_periods(tenant_id);
CREATE INDEX IF NOT EXISTS idx_commission_periods_unit 
    ON commission_periods(tenant_id, unit_id) WHERE unit_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_commission_periods_professional 
    ON commission_periods(tenant_id, professional_id) WHERE professional_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_commission_periods_month 
    ON commission_periods(tenant_id, reference_month);
CREATE INDEX IF NOT EXISTS idx_commission_periods_status 
    ON commission_periods(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_commission_periods_open 
    ON commission_periods(tenant_id, status) WHERE status = 'ABERTO';

-- Comentários
COMMENT ON TABLE commission_periods IS 'Períodos mensais de fechamento de comissões por profissional';
COMMENT ON COLUMN commission_periods.reference_month IS 'Mês de referência no formato YYYY-MM';
COMMENT ON COLUMN commission_periods.total_net IS 'Valor líquido = comissão - adiantamentos + ajustes';
COMMENT ON COLUMN commission_periods.conta_pagar_id IS 'FK para conta a pagar gerada no fechamento';
