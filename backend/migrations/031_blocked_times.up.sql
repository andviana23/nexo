-- Migration: Create blocked_times table for schedule blocking feature
-- Created: 2025-11-30
-- Description: Allows blocking time slots in the schedule (vacation, breaks, etc.)

BEGIN;

-- Create blocked_times table
CREATE TABLE IF NOT EXISTS blocked_times (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    
    -- Time range
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    
    -- Metadata
    reason TEXT NOT NULL,
    is_recurring BOOLEAN DEFAULT FALSE NOT NULL,
    recurrence_rule TEXT, -- iCal RRULE format (for future use)
    
    -- Audit
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Constraints
    CONSTRAINT blocked_times_time_check CHECK (end_time > start_time)
);

-- Indexes for efficient querying
CREATE INDEX idx_blocked_times_tenant_professional 
    ON blocked_times(tenant_id, professional_id, start_time);

CREATE INDEX idx_blocked_times_time_range 
    ON blocked_times(tenant_id, professional_id, start_time, end_time);

-- Comments
COMMENT ON TABLE blocked_times IS 'Horários bloqueados na agenda (férias, intervalos, etc.)';
COMMENT ON COLUMN blocked_times.tenant_id IS 'ID do tenant (multi-tenancy)';
COMMENT ON COLUMN blocked_times.professional_id IS 'ID do profissional (barbeiro)';
COMMENT ON COLUMN blocked_times.start_time IS 'Início do bloqueio';
COMMENT ON COLUMN blocked_times.end_time IS 'Fim do bloqueio';
COMMENT ON COLUMN blocked_times.reason IS 'Motivo do bloqueio (obrigatório)';
COMMENT ON COLUMN blocked_times.is_recurring IS 'Se o bloqueio se repete';
COMMENT ON COLUMN blocked_times.recurrence_rule IS 'Regra de recorrência no formato iCal RRULE';

-- Enable RLS
ALTER TABLE blocked_times ENABLE ROW LEVEL SECURITY;

-- RLS Policy: Users can only see blocked_times from their tenant
CREATE POLICY blocked_times_tenant_isolation ON blocked_times
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id', true)::uuid);

COMMIT;
