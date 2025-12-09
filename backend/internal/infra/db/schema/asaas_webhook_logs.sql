-- Schema: asaas_webhook_logs (Auditoria de webhooks)
-- Referência: PLANO_AJUSTE_ASAAS.md — Migration 041

CREATE TABLE IF NOT EXISTS asaas_webhook_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Dados do evento
    event_type VARCHAR(50) NOT NULL,
    asaas_payment_id VARCHAR(100),
    asaas_subscription_id VARCHAR(100),
    
    -- Payload completo (para debug)
    payload JSONB NOT NULL,
    
    -- Resultado do processamento
    processed_at TIMESTAMPTZ,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_webhook_logs_tenant ON asaas_webhook_logs(tenant_id);
CREATE INDEX idx_webhook_logs_payment ON asaas_webhook_logs(asaas_payment_id) WHERE asaas_payment_id IS NOT NULL;
CREATE INDEX idx_webhook_logs_subscription ON asaas_webhook_logs(asaas_subscription_id) WHERE asaas_subscription_id IS NOT NULL;
CREATE INDEX idx_webhook_logs_unprocessed ON asaas_webhook_logs(created_at) WHERE processed_at IS NULL;

COMMENT ON TABLE asaas_webhook_logs IS 'Log de webhooks recebidos do Asaas para auditoria e retry';

-- Schema: asaas_reconciliation_logs (Auditoria de conciliação)
-- Referência: PLANO_AJUSTE_ASAAS.md — Migration 041

CREATE TABLE IF NOT EXISTS asaas_reconciliation_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Período conciliado
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    
    -- Resultado
    total_asaas INTEGER DEFAULT 0,
    total_nexo INTEGER DEFAULT 0,
    divergences INTEGER DEFAULT 0,
    auto_fixed INTEGER DEFAULT 0,
    pending_review INTEGER DEFAULT 0,
    
    -- Detalhes (JSON com lista de divergências)
    details JSONB,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reconciliation_logs_tenant ON asaas_reconciliation_logs(tenant_id);
CREATE INDEX idx_reconciliation_logs_period ON asaas_reconciliation_logs(tenant_id, period_start, period_end);

COMMENT ON TABLE asaas_reconciliation_logs IS 'Log de execuções de conciliação Asaas x NEXO';
