-- Schema: subscriptions (Assinaturas de clientes)
-- Referência: FLUXO_ASSINATURA.md — Seção 3

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    cliente_id UUID NOT NULL REFERENCES clientes(id) ON DELETE RESTRICT,
    plano_id UUID NOT NULL REFERENCES plans(id) ON DELETE RESTRICT,
    
    asaas_customer_id VARCHAR(100),
    asaas_subscription_id VARCHAR(100),
    
    forma_pagamento VARCHAR(20) NOT NULL CHECK (
        forma_pagamento IN ('CARTAO', 'PIX', 'DINHEIRO')
    ),
    status VARCHAR(30) NOT NULL DEFAULT 'AGUARDANDO_PAGAMENTO' CHECK (
        status IN ('AGUARDANDO_PAGAMENTO', 'ATIVO', 'INADIMPLENTE', 'INATIVO', 'CANCELADO')
    ),
    valor DECIMAL(10,2) NOT NULL CHECK (valor >= 1.00),
    link_pagamento TEXT,
    codigo_transacao VARCHAR(100),
    
    data_ativacao TIMESTAMPTZ,
    data_vencimento TIMESTAMPTZ,
    data_cancelamento TIMESTAMPTZ,
    cancelado_por UUID REFERENCES users(id),
    
    servicos_utilizados INTEGER NOT NULL DEFAULT 0,
    
    -- Campos Asaas v2 (Migration 041)
    next_due_date DATE,
    cycle VARCHAR(20) DEFAULT 'MONTHLY',
    asaas_status VARCHAR(30),
    last_confirmed_at TIMESTAMPTZ,
    last_sync_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_tenant ON subscriptions(tenant_id);
CREATE INDEX idx_subscriptions_cliente ON subscriptions(cliente_id);
CREATE INDEX idx_subscriptions_plano ON subscriptions(plano_id);
CREATE INDEX idx_subscriptions_tenant_status ON subscriptions(tenant_id, status);
CREATE INDEX idx_subscriptions_vencimento ON subscriptions(data_vencimento) WHERE status = 'ATIVO';
CREATE INDEX idx_subscriptions_asaas ON subscriptions(asaas_subscription_id) WHERE asaas_subscription_id IS NOT NULL;
