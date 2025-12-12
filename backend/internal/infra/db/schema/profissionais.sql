-- =============================================================================
-- Schema: profissionais (para sqlc)
-- =============================================================================
CREATE TABLE IF NOT EXISTS profissionais (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE RESTRICT,
    user_id UUID,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telefone VARCHAR(20) NOT NULL,
    cpf VARCHAR(11) NOT NULL,
    especialidades TEXT[],
    comissao NUMERIC(5,2) DEFAULT 0.00,
    tipo_comissao VARCHAR(20) DEFAULT 'PERCENTUAL',
    foto TEXT,
    data_admissao DATE DEFAULT CURRENT_DATE NOT NULL,
    data_demissao DATE,
    status VARCHAR(20) DEFAULT 'ATIVO' CHECK (status IN ('ATIVO', 'INATIVO', 'AFASTADO', 'DEMITIDO')),
    horario_trabalho JSONB,
    observacoes TEXT,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    tipo VARCHAR(30) DEFAULT 'BARBEIRO' NOT NULL CHECK (tipo IN ('BARBEIRO', 'MANICURE', 'RECEPCIONISTA', 'GERENTE', 'OUTRO')),
    cor VARCHAR(7)
);
