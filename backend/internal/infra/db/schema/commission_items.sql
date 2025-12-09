-- Schema: commission_items
-- Itens individuais de comissão

CREATE TABLE IF NOT EXISTS commission_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    command_id UUID REFERENCES commands(id) ON DELETE CASCADE,
    command_item_id UUID REFERENCES command_items(id) ON DELETE CASCADE,
    appointment_id UUID REFERENCES appointments(id) ON DELETE SET NULL,
    service_id UUID REFERENCES servicos(id) ON DELETE SET NULL,
    service_name VARCHAR(255),
    gross_value NUMERIC(15,2) NOT NULL,
    commission_rate NUMERIC(10,2) NOT NULL,
    commission_type VARCHAR(20) NOT NULL DEFAULT 'PERCENTUAL' CHECK (commission_type IN ('PERCENTUAL', 'FIXO')),
    commission_value NUMERIC(15,2) NOT NULL,
    commission_source VARCHAR(20) NOT NULL DEFAULT 'PROFISSIONAL' CHECK (commission_source IN ('SERVICO', 'PROFISSIONAL', 'REGRA', 'MANUAL')),
    rule_id UUID REFERENCES commission_rules(id) ON DELETE SET NULL,
    reference_date DATE NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'PENDENTE' NOT NULL CHECK (status IN ('PENDENTE', 'PROCESSADO', 'PAGO', 'CANCELADO', 'ESTORNADO')),
    period_id UUID REFERENCES commission_periods(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    processed_at TIMESTAMPTZ
);
