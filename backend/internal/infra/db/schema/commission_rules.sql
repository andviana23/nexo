-- Schema: commission_rules
-- Regras de comissão

CREATE TABLE IF NOT EXISTS commission_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL DEFAULT 'PERCENTUAL' CHECK (type IN ('PERCENTUAL', 'FIXO')),
    default_rate NUMERIC(10,2) NOT NULL,
    min_amount NUMERIC(15,2),
    max_amount NUMERIC(15,2),
    calculation_base VARCHAR(20) DEFAULT 'BRUTO' CHECK (calculation_base IN ('BRUTO', 'LIQUIDO')),
    effective_from DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_to DATE,
    priority INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL
);
