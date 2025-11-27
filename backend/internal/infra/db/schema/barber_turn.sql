-- =============================================================================
-- Schema: barbers_turn_list e barber_turn_history (para sqlc)
-- Lista da Vez — NEXO v1.0
-- =============================================================================

-- Tabela principal: Estado atual da fila
CREATE TABLE IF NOT EXISTS barbers_turn_list (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    current_points INTEGER DEFAULT 0 NOT NULL CHECK (current_points >= 0),
    last_turn_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    CONSTRAINT uq_barber_turn_professional_tenant UNIQUE (professional_id, tenant_id)
);

-- Trigger para validar que profissional é do tipo BARBEIRO
CREATE TRIGGER trg_validate_barber_type 
    BEFORE INSERT OR UPDATE ON barbers_turn_list 
    FOR EACH ROW 
    EXECUTE FUNCTION check_professional_is_barber();

-- Tabela de histórico: Snapshot mensal
CREATE TABLE IF NOT EXISTS barber_turn_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    month_year VARCHAR(7) NOT NULL,
    total_turns INTEGER DEFAULT 0 NOT NULL,
    final_points INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    CONSTRAINT uq_history_professional_month UNIQUE (professional_id, tenant_id, month_year)
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_barbers_turn_tenant ON barbers_turn_list(tenant_id);
CREATE INDEX IF NOT EXISTS idx_barbers_turn_points ON barbers_turn_list(current_points);
CREATE INDEX IF NOT EXISTS idx_barbers_turn_active ON barbers_turn_list(is_active);
CREATE INDEX IF NOT EXISTS idx_barbers_turn_last_turn ON barbers_turn_list(last_turn_at);

CREATE INDEX IF NOT EXISTS idx_turn_history_tenant ON barber_turn_history(tenant_id);
CREATE INDEX IF NOT EXISTS idx_turn_history_month ON barber_turn_history(month_year);
CREATE INDEX IF NOT EXISTS idx_turn_history_professional ON barber_turn_history(professional_id);
