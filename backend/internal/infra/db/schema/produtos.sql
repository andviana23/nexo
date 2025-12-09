-- Tabela: produtos (schema base)
CREATE TABLE IF NOT EXISTS produtos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    categoria_id UUID,
    categoria_produto_id UUID,  -- FK para categorias_produtos customizadas
    fornecedor_id UUID,         -- FK para fornecedores

    -- Dados básicos
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    sku VARCHAR(100),
    codigo_barras VARCHAR(50),

    -- Pricing
    preco NUMERIC(10,2) NOT NULL CHECK (preco > 0),
    custo NUMERIC(10,2) CHECK (custo IS NULL OR custo >= 0),
    valor_venda_profissional NUMERIC(10,2) DEFAULT 0,  -- Preço de venda para profissionais
    valor_entrada NUMERIC(10,2) DEFAULT 0,             -- Custo de aquisição

    -- Estoque (colunas legadas - manter compatibilidade)
    estoque INTEGER DEFAULT 0 CHECK (estoque >= 0),
    estoque_minimo INTEGER DEFAULT 0 CHECK (estoque_minimo >= 0),
    estoque_maximo INTEGER DEFAULT 0 CHECK (estoque_maximo >= 0),
    unidade VARCHAR(10) DEFAULT 'UN',
    fornecedor VARCHAR(255),

    -- Módulo de Estoque (novas colunas v2.0)
    categoria_produto VARCHAR(30) DEFAULT 'REVENDA' NOT NULL
        CHECK (categoria_produto IN ('POMADA', 'SHAMPOO', 'CREME', 'LAMINA', 'TOALHA', 'LIMPEZA', 'ESCRITORIO', 'BEBIDA', 'REVENDA', 'INSUMO', 'USO_INTERNO', 'PERMANENTE', 'PROMOCIONAL', 'KIT', 'SERVICO')),
    centro_custo VARCHAR(50)
        CHECK (centro_custo IN ('CUSTO_SERVICO', 'DESPESA_OPERACIONAL', 'CMV')),
    unidade_medida VARCHAR(20) DEFAULT 'UN' NOT NULL
        CHECK (unidade_medida IN ('UN', 'KG', 'G', 'ML', 'L')),
    quantidade_atual NUMERIC(15,3) DEFAULT 0 NOT NULL CHECK (quantidade_atual >= 0),
    quantidade_minima NUMERIC(15,3) DEFAULT 0 NOT NULL CHECK (quantidade_minima >= 0),
    localizacao VARCHAR(100),
    lote VARCHAR(50),
    data_validade DATE,
    ncm VARCHAR(8),
    permite_venda BOOLEAN DEFAULT true NOT NULL,
    controla_validade BOOLEAN DEFAULT false,
    lead_time_dias INT DEFAULT 7,

    -- Metadados
    imagem TEXT,
    observacoes TEXT,
    ativo BOOLEAN DEFAULT true,
    unit_id UUID,

    -- Timestamps
    criado_em TIMESTAMPTZ DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_produtos_tenant_id ON produtos(tenant_id);
CREATE INDEX IF NOT EXISTS idx_produtos_categoria_id ON produtos(categoria_id);
CREATE INDEX IF NOT EXISTS idx_produtos_categoria_estoque ON produtos(tenant_id, categoria_produto);
CREATE INDEX IF NOT EXISTS idx_produtos_sku_tenant ON produtos(tenant_id, sku) WHERE sku IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_produtos_ativo ON produtos(tenant_id, ativo);
