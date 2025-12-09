-- Schema: plans (Modelos de planos de assinatura)
-- Referência: FLUXO_ASSINATURA.md — Seção 2.2

CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    valor DECIMAL(10,2) NOT NULL CHECK (valor > 0),
    periodicidade VARCHAR(20) NOT NULL DEFAULT 'MENSAL',
    qtd_servicos INTEGER,
    limite_uso_mensal INTEGER,
    ativo BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(tenant_id, nome)
);

CREATE INDEX idx_plans_tenant ON plans(tenant_id);
CREATE INDEX idx_plans_tenant_ativo ON plans(tenant_id, ativo);
