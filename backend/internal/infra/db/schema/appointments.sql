-- =============================================================================
-- Schema: appointments (para sqlc)
-- Módulo de Agendamento — NEXO v1.0
-- =============================================================================

CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE RESTRICT,
    customer_id UUID NOT NULL REFERENCES clientes(id) ON DELETE RESTRICT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'CREATED' 
        CHECK (status IN ('CREATED', 'CONFIRMED', 'CHECKED_IN', 'IN_SERVICE', 'AWAITING_PAYMENT', 'DONE', 'NO_SHOW', 'CANCELED')),
    total_price NUMERIC(10,2) NOT NULL CHECK (total_price >= 0),
    notes TEXT,
    canceled_reason TEXT,
    google_calendar_event_id VARCHAR(255),
    command_id UUID REFERENCES commands(id) ON DELETE SET NULL,
    unit_id UUID REFERENCES units(id) ON DELETE RESTRICT,
    checked_in_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT appointments_time_check CHECK (end_time > start_time)
);

-- Índices para busca eficiente
CREATE INDEX IF NOT EXISTS idx_appointments_tenant_start ON appointments(tenant_id, start_time);
CREATE INDEX IF NOT EXISTS idx_appointments_professional ON appointments(professional_id);
CREATE INDEX IF NOT EXISTS idx_appointments_customer ON appointments(customer_id);
CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(status);
CREATE INDEX IF NOT EXISTS idx_appointments_command_id ON appointments(command_id) WHERE command_id IS NOT NULL;

-- =============================================================================
-- Schema: appointment_services (para sqlc)
-- Relacionamento N:N entre agendamentos e serviços
-- =============================================================================

CREATE TABLE IF NOT EXISTS appointment_services (
    appointment_id UUID NOT NULL REFERENCES appointments(id) ON DELETE CASCADE,
    service_id UUID NOT NULL REFERENCES servicos(id) ON DELETE RESTRICT,
    price_at_booking NUMERIC(10,2) NOT NULL CHECK (price_at_booking >= 0),
    duration_at_booking INTEGER NOT NULL CHECK (duration_at_booking > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (appointment_id, service_id)
);

-- Índice para busca por serviço
CREATE INDEX IF NOT EXISTS idx_appointment_services_service ON appointment_services(service_id);
