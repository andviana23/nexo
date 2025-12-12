-- Contas a pagar com suporte a recorrência e notificações
-- Tabela: contas_a_pagar

CREATE TABLE IF NOT EXISTS contas_a_pagar (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    descricao VARCHAR(255) NOT NULL,
    categoria_id UUID,
    fornecedor VARCHAR(255),

    valor NUMERIC(15,2) NOT NULL CHECK (valor > 0),
    tipo VARCHAR(20) DEFAULT 'VARIAVEL' CHECK (tipo IN ('FIXA', 'VARIAVEL')),

    recorrente BOOLEAN DEFAULT false,
    periodicidade VARCHAR(20) CHECK (periodicidade IN ('MENSAL', 'TRIMESTRAL', 'ANUAL')),

    data_vencimento DATE NOT NULL,
    data_pagamento DATE,
    status VARCHAR(20) DEFAULT 'ABERTO' CHECK (status IN ('ABERTO', 'PAGO', 'ATRASADO', 'CANCELADO')),
    unit_id UUID REFERENCES units(id) ON DELETE RESTRICT,

    comprovante_url TEXT,
    pix_code TEXT,
    observacoes TEXT,

    criado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),

    CONSTRAINT contas_a_pagar_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT contas_a_pagar_categoria_id_fkey FOREIGN KEY (categoria_id) REFERENCES categorias(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_contas_pagar_tenant ON contas_a_pagar(tenant_id);
CREATE INDEX IF NOT EXISTS idx_contas_pagar_vencimento ON contas_a_pagar(tenant_id, data_vencimento);
CREATE INDEX IF NOT EXISTS idx_contas_pagar_status ON contas_a_pagar(status, data_vencimento);

COMMENT ON TABLE contas_a_pagar IS 'Contas a pagar com suporte a recorrência e notificações';
