-- Migration: 052_create_professional_category_commissions
-- Description: Cria tabela de comissões personalizadas por categoria para profissionais
-- Author: Antigravity
-- Date: 2025-12-09

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

-- Comentários
COMMENT ON TABLE comissoes_categoria_profissional IS 'Comissões personalizadas por categoria para profissionais';
COMMENT ON COLUMN comissoes_categoria_profissional.comissao IS 'Percentual de comissão (0-100)';
