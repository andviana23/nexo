-- Contas a receber de assinaturas e serviços com alertas de inadimplência
-- Tabela: contas_a_receber

CREATE TABLE IF NOT EXISTS contas_a_receber (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    origem VARCHAR(30) CHECK (origem IN ('ASSINATURA', 'SERVICO', 'OUTRO')),
    assinatura_id UUID,
    servico_id UUID,

    descricao VARCHAR(255) NOT NULL,
    valor NUMERIC(15,2) NOT NULL CHECK (valor > 0),
    valor_pago NUMERIC(15,2) DEFAULT 0 CHECK (valor_pago >= 0),

    data_vencimento DATE NOT NULL,
    data_recebimento DATE,
    status VARCHAR(20) DEFAULT 'PENDENTE' CHECK (status IN ('PENDENTE', 'RECEBIDO', 'ATRASADO', 'CANCELADO', 'CONFIRMADO', 'ESTORNADO')),

    observacoes TEXT,
    
    -- Campos Asaas v2 (Migration 041)
    subscription_id UUID REFERENCES subscriptions(id) ON DELETE SET NULL,
    asaas_payment_id VARCHAR(100),
    competencia_mes VARCHAR(7),
    confirmed_at TIMESTAMPTZ,
    received_at TIMESTAMPTZ,
    
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),

    CONSTRAINT contas_a_receber_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT contas_a_receber_assinatura_id_fkey FOREIGN KEY (assinatura_id) REFERENCES assinaturas(id) ON DELETE CASCADE,
    CONSTRAINT contas_a_receber_servico_id_fkey FOREIGN KEY (servico_id) REFERENCES servicos(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_contas_receber_tenant ON contas_a_receber(tenant_id);
CREATE INDEX IF NOT EXISTS idx_contas_receber_vencimento ON contas_a_receber(tenant_id, data_vencimento);
CREATE INDEX IF NOT EXISTS idx_contas_receber_status ON contas_a_receber(status, data_vencimento);
CREATE INDEX IF NOT EXISTS idx_contas_receber_assinatura ON contas_a_receber(assinatura_id) WHERE assinatura_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_contas_a_receber_asaas_payment_id ON contas_a_receber(tenant_id, asaas_payment_id) WHERE asaas_payment_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_contas_a_receber_competencia ON contas_a_receber(tenant_id, competencia_mes, status);

COMMENT ON TABLE contas_a_receber IS 'Contas a receber de assinaturas e serviços com alertas de inadimplência';
