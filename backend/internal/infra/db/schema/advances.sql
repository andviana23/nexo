-- Schema: advances
-- Adiantamentos de profissionais

CREATE TABLE IF NOT EXISTS advances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    amount NUMERIC(15,2) NOT NULL,
    request_date DATE NOT NULL DEFAULT CURRENT_DATE,
    reason TEXT,
    status VARCHAR(20) DEFAULT 'PENDING' NOT NULL CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED', 'DEDUCTED', 'CANCELLED')),
    approved_at TIMESTAMPTZ,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    rejected_at TIMESTAMPTZ,
    rejected_by UUID REFERENCES users(id) ON DELETE SET NULL,
    rejection_reason TEXT,
    deducted_at TIMESTAMPTZ,
    deduction_period_id UUID REFERENCES commission_periods(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL
);
