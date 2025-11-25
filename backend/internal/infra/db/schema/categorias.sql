-- Tabela: categorias (schema mínimo para referência)
CREATE TABLE IF NOT EXISTS categorias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    nome VARCHAR(100) NOT NULL,
    tipo VARCHAR(20) NOT NULL CHECK (tipo IN ('RECEITA', 'DESPESA')),
    cor VARCHAR(7) DEFAULT '#000000',
    ativa BOOLEAN DEFAULT true,
    criado_em TIMESTAMPTZ DEFAULT NOW(),
    tipo_custo VARCHAR(20) DEFAULT 'FIXO' CHECK (tipo_custo IN ('FIXO', 'VARIAVEL')),
    CONSTRAINT categorias_tenant_id_nome_key UNIQUE (tenant_id, nome)
);

CREATE INDEX IF NOT EXISTS idx_categorias_tenant_id ON categorias(tenant_id);
