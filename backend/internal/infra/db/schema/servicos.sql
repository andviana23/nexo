-- =============================================================================
-- Schema: servicos (para sqlc)
-- =============================================================================
CREATE TABLE IF NOT EXISTS servicos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE RESTRICT,
    categoria_id UUID,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    preco NUMERIC(10,2) NOT NULL,
    duracao INTEGER NOT NULL,
    comissao NUMERIC(5,2) DEFAULT 0.00,
    cor VARCHAR(7),
    imagem TEXT,
    profissionais_ids UUID[],
    observacoes TEXT,
    tags TEXT[],
    ativo BOOLEAN DEFAULT true,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now()
);
