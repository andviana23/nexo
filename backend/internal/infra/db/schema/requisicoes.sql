-- Requisições de Compra
CREATE TABLE IF NOT EXISTS requisicoes_compra (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    solicitante_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'PENDENTE', -- PENDENTE, APROVADA, COMPRADA, CANCELADA
    observacoes TEXT,
    valor_estimado NUMERIC(10,2),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS itens_requisicao (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requisicao_id UUID NOT NULL REFERENCES requisicoes_compra(id),
    produto_id UUID NOT NULL REFERENCES produtos(id),
    qtd_sugerida INT,
    qtd_aprovada INT
);
