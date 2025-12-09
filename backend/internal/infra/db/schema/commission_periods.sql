-- Schema: commission_periods
-- Períodos de fechamento de comissões

CREATE TABLE IF NOT EXISTS commission_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
    reference_month VARCHAR(7) NOT NULL,
    professional_id UUID REFERENCES profissionais(id) ON DELETE CASCADE,
    total_gross NUMERIC(15,2) DEFAULT 0 NOT NULL,
    total_commission NUMERIC(15,2) DEFAULT 0 NOT NULL,
    total_advances NUMERIC(15,2) DEFAULT 0 NOT NULL,
    total_adjustments NUMERIC(15,2) DEFAULT 0 NOT NULL,
    total_net NUMERIC(15,2) DEFAULT 0 NOT NULL,
    items_count INTEGER DEFAULT 0 NOT NULL,
    status VARCHAR(20) DEFAULT 'ABERTO' NOT NULL CHECK (status IN ('ABERTO', 'PROCESSANDO', 'FECHADO', 'PAGO', 'CANCELADO')),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    closed_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    conta_pagar_id UUID REFERENCES contas_a_pagar(id) ON DELETE SET NULL,
    closed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    paid_by UUID REFERENCES users(id) ON DELETE SET NULL,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);
