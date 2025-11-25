-- Tabela: produto_fornecedor (relacionamento N:N)
CREATE TABLE IF NOT EXISTS produto_fornecedor (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    produto_id UUID NOT NULL,
    fornecedor_id UUID NOT NULL,

    -- Dados específicos
    codigo_fornecedor VARCHAR(100),
    preco_compra NUMERIC(15,2) CHECK (preco_compra IS NULL OR preco_compra >= 0),
    prazo_entrega_dias INTEGER DEFAULT 0 CHECK (prazo_entrega_dias >= 0),

    -- Controle
    fornecedor_preferencial BOOLEAN DEFAULT false,
    ativo BOOLEAN DEFAULT true NOT NULL,

    -- Timestamps
    criado_em TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    atualizado_em TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    -- Unique constraint
    CONSTRAINT produto_fornecedor_unique UNIQUE(tenant_id, produto_id, fornecedor_id)
);

CREATE INDEX IF NOT EXISTS idx_produto_fornecedor_tenant_produto ON produto_fornecedor(tenant_id, produto_id);
CREATE INDEX IF NOT EXISTS idx_produto_fornecedor_tenant_fornecedor ON produto_fornecedor(tenant_id, fornecedor_id);
CREATE INDEX IF NOT EXISTS idx_produto_fornecedor_preferencial ON produto_fornecedor(tenant_id, produto_id)
    WHERE fornecedor_preferencial = true;
