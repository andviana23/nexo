CREATE TABLE IF NOT EXISTS comissoes_categoria_profissional (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    profissional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    categoria_id UUID NOT NULL REFERENCES categorias_servicos(id) ON DELETE CASCADE,
    comissao NUMERIC(5,2) NOT NULL CHECK (comissao >= 0 AND comissao <= 100),
    criado_em TIMESTAMPTZ DEFAULT now() NOT NULL,
    atualizado_em TIMESTAMPTZ DEFAULT now() NOT NULL,

    CONSTRAINT uq_comissao_profissional_categoria UNIQUE (tenant_id, profissional_id, categoria_id)
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_comissoes_cat_prof_tenant ON comissoes_categoria_profissional(tenant_id);
CREATE INDEX IF NOT EXISTS idx_comissoes_cat_prof_profissional ON comissoes_categoria_profissional(tenant_id, profissional_id);
