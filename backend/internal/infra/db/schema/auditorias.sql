-- Tabela Auditorias
CREATE TABLE IF NOT EXISTS auditorias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    responsavel_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'ABERTA',
    data_inicio TIMESTAMP DEFAULT NOW(),
    data_fim TIMESTAMP
);

-- Itens Auditoria
CREATE TABLE IF NOT EXISTS itens_auditoria (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auditoria_id UUID NOT NULL REFERENCES auditorias(id),
    produto_id UUID NOT NULL REFERENCES produtos(id),
    qtd_sistema INT NOT NULL,
    qtd_contada INT NOT NULL,
    divergencia INT GENERATED ALWAYS AS (qtd_contada - qtd_sistema) STORED,
    justificativa TEXT
);
