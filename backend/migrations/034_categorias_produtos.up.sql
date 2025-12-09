-- ============================================================================
-- MIGRATION: Criar tabela categorias_produtos
-- Permite cadastro customizado de categorias de produtos por tenant
-- ============================================================================

-- Tabela de categorias de produtos
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
    
    -- Nome único por tenant
    CONSTRAINT categorias_produtos_tenant_nome_unique UNIQUE (tenant_id, nome),
    -- Validação cor hexadecimal
    CONSTRAINT chk_cor_hex CHECK (cor ~ '^#[0-9A-Fa-f]{6}$')
);

-- Comentários
COMMENT ON TABLE categorias_produtos IS 'Categorias customizáveis de produtos por tenant - isolamento multi-tenant obrigatório';
COMMENT ON COLUMN categorias_produtos.nome IS 'Nome da categoria (ex: Pomadas, Shampoos, Bebidas, Material Escritório)';
COMMENT ON COLUMN categorias_produtos.cor IS 'Cor hexadecimal para exibição na UI (#RRGGBB)';
COMMENT ON COLUMN categorias_produtos.icone IS 'Nome do ícone Lucide (package, droplet, scissors, etc)';
COMMENT ON COLUMN categorias_produtos.centro_custo IS 'Classificação DRE: CMV (revenda), CUSTO_SERVICO (insumos), DESPESA_OPERACIONAL (uso interno)';

-- Índices
CREATE INDEX IF NOT EXISTS idx_categorias_produtos_tenant 
    ON categorias_produtos(tenant_id);
CREATE INDEX IF NOT EXISTS idx_categorias_produtos_tenant_ativa 
    ON categorias_produtos(tenant_id, ativa);

-- Trigger para atualizar updated_at (drop first to make idempotent)
DROP TRIGGER IF EXISTS trg_categorias_produtos_updated_at ON categorias_produtos;
CREATE TRIGGER trg_categorias_produtos_updated_at
    BEFORE UPDATE ON categorias_produtos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- Adicionar FK opcional na tabela produtos (mantendo compatibilidade)
-- ============================================================================

-- Adiciona coluna de FK para categoria customizada
ALTER TABLE produtos 
ADD COLUMN IF NOT EXISTS categoria_produto_id UUID REFERENCES categorias_produtos(id) ON DELETE SET NULL;

COMMENT ON COLUMN produtos.categoria_produto_id IS 'FK para categoria customizada (substitui enum categoria_produto)';

CREATE INDEX IF NOT EXISTS idx_produtos_categoria_produto_id 
    ON produtos(tenant_id, categoria_produto_id);

-- ============================================================================
-- Seed: Categorias padrão (serão criadas por tenant no onboarding)
-- ============================================================================
-- Nota: As categorias padrão serão inseridas via seed ou no primeiro acesso
-- do tenant, não diretamente na migration para respeitar multi-tenancy.
