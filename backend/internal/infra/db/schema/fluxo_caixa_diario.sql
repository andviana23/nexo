-- Fluxo de caixa diário com previsões e compensações bancárias
-- Tabela: fluxo_caixa_diario

CREATE TABLE IF NOT EXISTS fluxo_caixa_diario (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    data DATE NOT NULL,

    saldo_inicial NUMERIC(15,2) DEFAULT 0,
    saldo_final NUMERIC(15,2) DEFAULT 0,

    entradas_confirmadas NUMERIC(15,2) DEFAULT 0,
    entradas_previstas NUMERIC(15,2) DEFAULT 0,

    saidas_pagas NUMERIC(15,2) DEFAULT 0,
    saidas_previstas NUMERIC(15,2) DEFAULT 0,
    
    -- Campos Asaas v2 (Migration 041)
    asaas_payments_count INTEGER DEFAULT 0,
    asaas_payments_total NUMERIC(15,2) DEFAULT 0,

    processado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),

    CONSTRAINT fluxo_caixa_diario_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fluxo_caixa_diario_tenant_id_data_key UNIQUE (tenant_id, data)
);

CREATE UNIQUE INDEX IF NOT EXISTS fluxo_caixa_diario_tenant_id_data_key ON fluxo_caixa_diario(tenant_id, data);
CREATE INDEX IF NOT EXISTS idx_fluxo_caixa_diario_tenant ON fluxo_caixa_diario(tenant_id);
CREATE INDEX IF NOT EXISTS idx_fluxo_caixa_diario_data ON fluxo_caixa_diario(tenant_id, data DESC);

COMMENT ON TABLE fluxo_caixa_diario IS 'Fluxo de caixa diário com previsões e compensações bancárias';
