-- Schema: subscription_payments (Histórico de pagamentos de assinaturas)
-- Referência: FLUXO_ASSINATURA.md — Seção 6.4

CREATE TABLE subscription_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    
    asaas_payment_id VARCHAR(100),
    
    valor DECIMAL(10,2) NOT NULL CHECK (valor > 0),
    forma_pagamento VARCHAR(20) NOT NULL CHECK (
        forma_pagamento IN ('CARTAO', 'PIX', 'DINHEIRO')
    ),
    status VARCHAR(30) NOT NULL DEFAULT 'PENDENTE' CHECK (
        status IN ('PENDENTE', 'CONFIRMADO', 'RECEBIDO', 'VENCIDO', 'ESTORNADO', 'CANCELADO')
    ),
    data_pagamento TIMESTAMPTZ,
    codigo_transacao VARCHAR(100),
    observacao TEXT,
    
    -- Campos Asaas v2 (Migration 041)
    status_asaas VARCHAR(30),
    due_date DATE,
    confirmed_date TIMESTAMPTZ,
    client_payment_date DATE,
    credit_date DATE,
    estimated_credit_date DATE,
    billing_type VARCHAR(30),
    net_value DECIMAL(10,2),
    invoice_url TEXT,
    bank_slip_url TEXT,
    pix_qrcode TEXT,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscription_payments_subscription ON subscription_payments(subscription_id);
CREATE INDEX idx_subscription_payments_tenant ON subscription_payments(tenant_id);
CREATE INDEX idx_subscription_payments_status ON subscription_payments(status);
CREATE INDEX idx_subscription_payments_asaas ON subscription_payments(asaas_payment_id) WHERE asaas_payment_id IS NOT NULL;
CREATE UNIQUE INDEX idx_subscription_payments_asaas_payment_id ON subscription_payments(asaas_payment_id) WHERE asaas_payment_id IS NOT NULL;
