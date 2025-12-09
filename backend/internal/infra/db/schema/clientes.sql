-- =============================================================================
-- Schema: clientes (para sqlc)
-- =============================================================================
CREATE TABLE IF NOT EXISTS clientes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    telefone VARCHAR(20) NOT NULL,
    cpf VARCHAR(11),
    data_nascimento DATE,
    genero VARCHAR(50),
    endereco_logradouro VARCHAR(255),
    endereco_numero VARCHAR(20),
    endereco_complemento VARCHAR(100),
    endereco_bairro VARCHAR(100),
    endereco_cidade VARCHAR(100),
    endereco_estado VARCHAR(2),
    endereco_cep VARCHAR(8),
    observacoes TEXT,
    tags TEXT[],
    ativo BOOLEAN DEFAULT true,
    -- Campos de integração Asaas (FLUXO_ASSINATURA.md — Seção 3.5)
    asaas_customer_id VARCHAR(100),
    is_subscriber BOOLEAN NOT NULL DEFAULT false,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    
    CONSTRAINT clientes_asaas_customer_unique UNIQUE(tenant_id, asaas_customer_id)
);

CREATE INDEX IF NOT EXISTS idx_clientes_subscriber ON clientes(tenant_id, is_subscriber);
CREATE INDEX IF NOT EXISTS idx_clientes_asaas_customer ON clientes(asaas_customer_id) WHERE asaas_customer_id IS NOT NULL;
