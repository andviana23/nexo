-- ============================================================================
-- SQLC Schema: Units (Unidades)
-- ============================================================================

CREATE TABLE IF NOT EXISTS units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    
    -- Dados básicos
    nome VARCHAR(100) NOT NULL,
    apelido VARCHAR(50),
    descricao TEXT,
    
    -- Endereço resumido
    endereco_resumo VARCHAR(255),
    cidade VARCHAR(100),
    estado VARCHAR(2),
    
    -- Configurações
    timezone VARCHAR(50) DEFAULT 'America/Sao_Paulo' NOT NULL,
    ativa BOOLEAN DEFAULT true NOT NULL,
    is_matriz BOOLEAN DEFAULT false NOT NULL,
    
    -- Timestamps
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- Constraints
    CONSTRAINT units_tenant_nome_unique UNIQUE (tenant_id, nome)
);
