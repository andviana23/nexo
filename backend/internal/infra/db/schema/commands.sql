-- Commands schema for sqlc
CREATE TABLE commands (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    appointment_id UUID,
    customer_id UUID NOT NULL,
    numero VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    
    subtotal NUMERIC(10,2) NOT NULL DEFAULT 0,
    desconto NUMERIC(10,2) NOT NULL DEFAULT 0,
    total NUMERIC(10,2) NOT NULL DEFAULT 0,
    total_recebido NUMERIC(10,2) NOT NULL DEFAULT 0,
    troco NUMERIC(10,2) NOT NULL DEFAULT 0,
    saldo_devedor NUMERIC(10,2) NOT NULL DEFAULT 0,
    
    observacoes TEXT,
    deixar_troco_gorjeta BOOLEAN DEFAULT false,
    deixar_saldo_divida BOOLEAN DEFAULT false,
    
    criado_em TIMESTAMPTZ NOT NULL,
    atualizado_em TIMESTAMPTZ NOT NULL,
    fechado_em TIMESTAMPTZ,
    fechado_por UUID
);

CREATE TABLE command_items (
    id UUID PRIMARY KEY,
    command_id UUID NOT NULL,
    
    tipo VARCHAR(20) NOT NULL,
    item_id UUID NOT NULL,
    descricao VARCHAR(200) NOT NULL,
    
    preco_unitario NUMERIC(10,2) NOT NULL,
    quantidade INTEGER NOT NULL DEFAULT 1,
    desconto_valor NUMERIC(10,2) NOT NULL DEFAULT 0,
    desconto_percentual NUMERIC(5,2) NOT NULL DEFAULT 0,
    preco_final NUMERIC(10,2) NOT NULL,
    
    observacoes TEXT,
    criado_em TIMESTAMPTZ NOT NULL
);

CREATE TABLE command_payments (
    id UUID PRIMARY KEY,
    command_id UUID NOT NULL,
    meio_pagamento_id UUID NOT NULL,
    
    valor_recebido NUMERIC(10,2) NOT NULL,
    taxa_percentual NUMERIC(5,2) NOT NULL DEFAULT 0,
    taxa_fixa NUMERIC(10,2) NOT NULL DEFAULT 0,
    valor_liquido NUMERIC(10,2) NOT NULL,
    
    observacoes TEXT,
    criado_em TIMESTAMPTZ NOT NULL,
    criado_por UUID
);
