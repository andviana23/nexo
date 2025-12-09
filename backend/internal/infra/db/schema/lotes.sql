-- Tabela Lotes
CREATE TABLE IF NOT EXISTS lotes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    produto_id UUID NOT NULL REFERENCES produtos(id),
    codigo_lote VARCHAR(50),
    data_validade DATE NOT NULL,
    quantidade_inicial INT NOT NULL,
    quantidade_atual INT NOT NULL,
    ativo BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_lotes_validade ON lotes(data_validade);
CREATE INDEX IF NOT EXISTS idx_lotes_produto ON lotes(produto_id);
