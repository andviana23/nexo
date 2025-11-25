-- Tabela: fornecedores
-- Cadastro de fornecedores com dados completos
CREATE TABLE IF NOT EXISTS fornecedores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,

    -- Dados principais
    razao_social VARCHAR(255) NOT NULL,
    nome_fantasia VARCHAR(255),
    cnpj VARCHAR(14) NOT NULL,

    -- Contato
    email VARCHAR(255),
    telefone VARCHAR(20),
    celular VARCHAR(20),

    -- Endereço
    endereco_logradouro VARCHAR(255),
    endereco_numero VARCHAR(20),
    endereco_complemento VARCHAR(100),
    endereco_bairro VARCHAR(100),
    endereco_cidade VARCHAR(100),
    endereco_estado VARCHAR(2),
    endereco_cep VARCHAR(8),

    -- Dados bancários
    banco VARCHAR(100),
    agencia VARCHAR(20),
    conta VARCHAR(30),

    -- Controle
    observacoes TEXT,
    ativo BOOLEAN DEFAULT true NOT NULL,

    -- Timestamps
    criado_em TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    atualizado_em TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    -- Constraints
    CONSTRAINT fornecedores_cnpj_tenant_unique UNIQUE(tenant_id, cnpj),
    CONSTRAINT chk_fornecedor_cnpj_valido CHECK (LENGTH(cnpj) = 14)
);

CREATE INDEX IF NOT EXISTS idx_fornecedores_tenant_id ON fornecedores(tenant_id);
CREATE INDEX IF NOT EXISTS idx_fornecedores_ativo ON fornecedores(tenant_id, ativo);
CREATE INDEX IF NOT EXISTS idx_fornecedores_cnpj ON fornecedores(tenant_id, cnpj);
