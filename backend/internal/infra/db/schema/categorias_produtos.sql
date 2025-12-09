-- ============================================================================
-- SCHEMA: categorias_produtos
-- Categorias customizáveis de produtos por tenant
-- ============================================================================

CREATE TABLE IF NOT EXISTS categorias_produtos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    cor VARCHAR(7) DEFAULT '#6B7280',
    icone VARCHAR(50) DEFAULT 'package',
    centro_custo VARCHAR(50) DEFAULT 'CMV'
        CHECK (centro_custo IN ('CUSTO_SERVICO', 'DESPESA_OPERACIONAL', 'CMV')),
    ativa BOOLEAN DEFAULT true NOT NULL,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    CONSTRAINT categorias_produtos_tenant_nome_unique UNIQUE (tenant_id, nome),
    CONSTRAINT chk_cor_hex CHECK (cor ~ '^#[0-9A-Fa-f]{6}$')
);

CREATE INDEX IF NOT EXISTS idx_categorias_produtos_tenant 
    ON categorias_produtos(tenant_id);
CREATE INDEX IF NOT EXISTS idx_categorias_produtos_tenant_ativa 
    ON categorias_produtos(tenant_id, ativa);
