-- =============================================================================
-- Schema: blocked_times (para sqlc)
-- Bloqueios de horário na agenda
-- =============================================================================

CREATE TABLE IF NOT EXISTS blocked_times (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    reason TEXT NOT NULL,
    is_recurring BOOLEAN DEFAULT FALSE NOT NULL,
    recurrence_rule TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    
    CONSTRAINT blocked_times_time_check CHECK (end_time > start_time)
);

CREATE INDEX IF NOT EXISTS idx_blocked_times_tenant_professional 
    ON blocked_times(tenant_id, professional_id, start_time);

CREATE INDEX IF NOT EXISTS idx_blocked_times_time_range 
    ON blocked_times(tenant_id, professional_id, start_time, end_time);
