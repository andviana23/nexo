-- ============================================================================
-- SCHEMA: categorias_servicos
-- Categorias de Serviços (Cortes, Barba, Tratamentos, etc.)
-- IMPORTANTE: Separada de 'categorias' que é para receitas/despesas
-- ============================================================================

CREATE TABLE IF NOT EXISTS categorias_servicos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    cor VARCHAR(7) DEFAULT '#000000',
    icone VARCHAR(50),
    ativa BOOLEAN DEFAULT true,
    criado_em TIMESTAMPTZ DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT categorias_servicos_tenant_nome_unique UNIQUE (tenant_id, nome),
    CONSTRAINT chk_cor_hex CHECK (cor ~ '^#[0-9A-Fa-f]{6}$')
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_categorias_servicos_tenant 
ON categorias_servicos(tenant_id) 
WHERE ativa = true;

-- Comentários
COMMENT ON TABLE categorias_servicos IS 
'Categorias de serviços (Corte, Barba, Tratamentos, etc) - separadas das categorias financeiras';

COMMENT ON COLUMN categorias_servicos.tenant_id IS 
'Isolamento multi-tenant - OBRIGATÓRIO em todas as queries';

COMMENT ON COLUMN categorias_servicos.nome IS 
'Nome da categoria (ex: Cortes de Cabelo, Barba, Tratamentos Capilares)';

COMMENT ON COLUMN categorias_servicos.cor IS 
'Cor hexadecimal para exibição na UI (#RRGGBB)';

COMMENT ON COLUMN categorias_servicos.icone IS 
'Nome do ícone Material Icons (content_cut, straighten, spa, etc)';
