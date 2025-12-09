-- Migration: 044_advances
-- Description: Cria tabela de adiantamentos/vales para profissionais
-- Author: Nexo Team
-- Date: 2025-12-05

-- ============================================================================
-- TABELA: advances
-- Adiantamentos (vales) de profissionais
-- ============================================================================

CREATE TABLE IF NOT EXISTS advances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
    
    -- Profissional
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    
    -- Valores
    amount NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    
    -- Data da solicitação/efetivação
    request_date DATE NOT NULL DEFAULT CURRENT_DATE,
    
    -- Motivo opcional
    reason TEXT,
    
    -- Status do adiantamento
    status VARCHAR(20) DEFAULT 'PENDING' NOT NULL
        CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED', 'DEDUCTED', 'CANCELLED')),
    
    -- Aprovação
    approved_at TIMESTAMPTZ,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Rejeição
    rejected_at TIMESTAMPTZ,
    rejected_by UUID REFERENCES users(id) ON DELETE SET NULL,
    rejection_reason TEXT,
    
    -- Dedução (quando for descontado da comissão)
    deducted_at TIMESTAMPTZ,
    deduction_period_id UUID REFERENCES commission_periods(id) ON DELETE SET NULL,
    
    -- Auditoria
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Constraints
    CONSTRAINT chk_approval_data 
        CHECK (
            (status != 'APPROVED' OR approved_at IS NOT NULL) AND
            (status != 'REJECTED' OR rejected_at IS NOT NULL)
        )
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_advances_tenant 
    ON advances(tenant_id);
CREATE INDEX IF NOT EXISTS idx_advances_unit 
    ON advances(tenant_id, unit_id) WHERE unit_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_advances_professional 
    ON advances(tenant_id, professional_id);
CREATE INDEX IF NOT EXISTS idx_advances_status 
    ON advances(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_advances_pending 
    ON advances(tenant_id, status) WHERE status = 'PENDING';
CREATE INDEX IF NOT EXISTS idx_advances_approved_not_deducted 
    ON advances(tenant_id, professional_id, status) WHERE status = 'APPROVED';
CREATE INDEX IF NOT EXISTS idx_advances_request_date 
    ON advances(tenant_id, request_date);

-- Comentários
COMMENT ON TABLE advances IS 'Adiantamentos (vales) solicitados por profissionais';
COMMENT ON COLUMN advances.status IS 'PENDING=aguardando, APPROVED=aprovado, REJECTED=rejeitado, DEDUCTED=descontado, CANCELLED=cancelado';
COMMENT ON COLUMN advances.deduction_period_id IS 'Período de comissão onde foi deduzido';
