-- Tabela: tenants (schema mínimo para referência)
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome VARCHAR(255) NOT NULL UNIQUE,
    cnpj VARCHAR(14) UNIQUE,
    ativo BOOLEAN DEFAULT true,
    plano VARCHAR(50) DEFAULT 'free',
    criado_em TIMESTAMPTZ DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ DEFAULT NOW(),
    onboarding_completed BOOLEAN DEFAULT false,
    onboarding_step INTEGER DEFAULT 1
);
