-- Tabela: movimentacoes_estoque
-- Histórico completo de movimentações de estoque
CREATE TABLE IF NOT EXISTS movimentacoes_estoque (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,

    -- Produto e tipo
    produto_id UUID NOT NULL,
    tipo_movimentacao VARCHAR(30) NOT NULL CHECK (
        tipo_movimentacao IN ('ENTRADA', 'SAIDA', 'CONSUMO_INTERNO', 'AJUSTE', 'DEVOLUCAO', 'PERDA')
    ),

    -- Quantidades e valores
    quantidade NUMERIC(15,3) NOT NULL CHECK (quantidade > 0),
    valor_unitario NUMERIC(15,2) DEFAULT 0 NOT NULL CHECK (valor_unitario >= 0),
    valor_total NUMERIC(15,2) DEFAULT 0 NOT NULL CHECK (valor_total >= 0),

    -- Referências opcionais
    fornecedor_id UUID,
    usuario_id UUID,

    -- Metadados
    data_movimentacao TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    observacoes TEXT,
    documento VARCHAR(100),

    -- Timestamps
    criado_em TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    atualizado_em TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    -- Observações obrigatórias para AJUSTE e PERDA
    CONSTRAINT chk_observacoes_obrigatorias_ajuste_perda CHECK (
        tipo_movimentacao NOT IN ('AJUSTE', 'PERDA')
        OR (observacoes IS NOT NULL AND LENGTH(TRIM(observacoes)) > 0)
    )
);

CREATE INDEX IF NOT EXISTS idx_movimentacoes_tenant_produto ON movimentacoes_estoque(tenant_id, produto_id);
CREATE INDEX IF NOT EXISTS idx_movimentacoes_tenant_tipo ON movimentacoes_estoque(tenant_id, tipo_movimentacao);
CREATE INDEX IF NOT EXISTS idx_movimentacoes_tenant_data ON movimentacoes_estoque(tenant_id, data_movimentacao DESC);
CREATE INDEX IF NOT EXISTS idx_movimentacoes_fornecedor ON movimentacoes_estoque(fornecedor_id) WHERE fornecedor_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_movimentacoes_usuario ON movimentacoes_estoque(usuario_id) WHERE usuario_id IS NOT NULL;
