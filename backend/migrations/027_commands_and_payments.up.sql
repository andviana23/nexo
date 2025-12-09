-- =============================================
-- Migration: 027_commands_and_payments
-- Descrição: Criação das tabelas para comandas e pagamentos
-- Data: 29/11/2025
-- =============================================

-- Tabela principal de comandas
CREATE TABLE IF NOT EXISTS commands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    appointment_id UUID REFERENCES appointments(id) ON DELETE SET NULL,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    numero VARCHAR(50), -- Número sequencial da comanda (ex: "CMD-2025-00123")
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' 
        CHECK (status IN ('OPEN', 'CLOSED', 'CANCELED')),
    
    -- Valores financeiros
    subtotal NUMERIC(10,2) NOT NULL DEFAULT 0,
    desconto NUMERIC(10,2) NOT NULL DEFAULT 0,
    total NUMERIC(10,2) NOT NULL DEFAULT 0,
    total_recebido NUMERIC(10,2) NOT NULL DEFAULT 0,
    troco NUMERIC(10,2) NOT NULL DEFAULT 0,
    saldo_devedor NUMERIC(10,2) NOT NULL DEFAULT 0,
    
    -- Opções de fechamento
    observacoes TEXT,
    deixar_troco_gorjeta BOOLEAN DEFAULT false,
    deixar_saldo_divida BOOLEAN DEFAULT false,
    
    -- Auditoria
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    fechado_em TIMESTAMPTZ,
    fechado_por UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Constraints
    CONSTRAINT commands_total_check CHECK (total >= 0),
    CONSTRAINT commands_subtotal_check CHECK (subtotal >= 0),
    CONSTRAINT commands_desconto_check CHECK (desconto >= 0),
    CONSTRAINT commands_tenant_customer_check 
        CHECK (tenant_id IS NOT NULL AND customer_id IS NOT NULL)
);

-- Itens da comanda (serviços, produtos, pacotes)
CREATE TABLE IF NOT EXISTS command_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    command_id UUID NOT NULL REFERENCES commands(id) ON DELETE CASCADE,
    
    -- Tipo e referência
    tipo VARCHAR(20) NOT NULL 
        CHECK (tipo IN ('SERVICO', 'PRODUTO', 'PACOTE')),
    item_id UUID NOT NULL, -- ID do serviço/produto/pacote
    descricao VARCHAR(200) NOT NULL,
    
    -- Preços
    preco_unitario NUMERIC(10,2) NOT NULL,
    quantidade INTEGER NOT NULL DEFAULT 1,
    desconto_valor NUMERIC(10,2) NOT NULL DEFAULT 0,
    desconto_percentual NUMERIC(5,2) NOT NULL DEFAULT 0,
    preco_final NUMERIC(10,2) NOT NULL,
    
    -- Auditoria
    observacoes TEXT,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    -- Constraints
    CONSTRAINT command_items_preco_check CHECK (preco_unitario >= 0),
    CONSTRAINT command_items_quantidade_check CHECK (quantidade > 0),
    CONSTRAINT command_items_desconto_valor_check CHECK (desconto_valor >= 0),
    CONSTRAINT command_items_desconto_percentual_check 
        CHECK (desconto_percentual >= 0 AND desconto_percentual <= 100),
    CONSTRAINT command_items_preco_final_check CHECK (preco_final >= 0)
);

-- Pagamentos da comanda (múltiplas formas de pagamento)
CREATE TABLE IF NOT EXISTS command_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    command_id UUID NOT NULL REFERENCES commands(id) ON DELETE CASCADE,
    meio_pagamento_id UUID NOT NULL REFERENCES meios_pagamento(id) ON DELETE RESTRICT,
    
    -- Valores
    valor_recebido NUMERIC(10,2) NOT NULL,
    taxa_percentual NUMERIC(5,2) NOT NULL DEFAULT 0,
    taxa_fixa NUMERIC(10,2) NOT NULL DEFAULT 0,
    valor_liquido NUMERIC(10,2) NOT NULL, -- Após taxas
    
    -- Auditoria
    observacoes TEXT,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    criado_por UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Constraints
    CONSTRAINT command_payments_valor_check CHECK (valor_recebido > 0),
    CONSTRAINT command_payments_taxa_percentual_check 
        CHECK (taxa_percentual >= 0 AND taxa_percentual <= 100),
    CONSTRAINT command_payments_taxa_fixa_check CHECK (taxa_fixa >= 0),
    CONSTRAINT command_payments_liquido_check CHECK (valor_liquido >= 0)
);

-- =============================================
-- ÍNDICES
-- =============================================

-- Commands
CREATE INDEX idx_commands_tenant_id ON commands(tenant_id);
CREATE INDEX idx_commands_customer_id ON commands(customer_id);
CREATE INDEX idx_commands_appointment_id ON commands(appointment_id) WHERE appointment_id IS NOT NULL;
CREATE INDEX idx_commands_status ON commands(status);
CREATE INDEX idx_commands_criado_em ON commands(criado_em DESC);
CREATE INDEX idx_commands_numero ON commands(tenant_id, numero);

-- Command Items
CREATE INDEX idx_command_items_command_id ON command_items(command_id);
CREATE INDEX idx_command_items_tipo ON command_items(tipo);
CREATE INDEX idx_command_items_item_id ON command_items(item_id);

-- Command Payments
CREATE INDEX idx_command_payments_command_id ON command_payments(command_id);
CREATE INDEX idx_command_payments_meio_pagamento ON command_payments(meio_pagamento_id);
CREATE INDEX idx_command_payments_criado_em ON command_payments(criado_em DESC);

-- =============================================
-- TRIGGERS
-- =============================================

-- Trigger para atualizar updated_at em commands
CREATE OR REPLACE FUNCTION update_commands_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.atualizado_em = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_commands_updated_at
    BEFORE UPDATE ON commands
    FOR EACH ROW
    EXECUTE FUNCTION update_commands_updated_at();

-- =============================================
-- RLS (Row Level Security)
-- =============================================

ALTER TABLE commands ENABLE ROW LEVEL SECURITY;
ALTER TABLE command_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE command_payments ENABLE ROW LEVEL SECURITY;

-- Policy para commands
CREATE POLICY commands_tenant_isolation ON commands
    USING (tenant_id = current_setting('app.current_tenant_id', TRUE)::UUID);

-- Policy para command_items (via command)
CREATE POLICY command_items_tenant_isolation ON command_items
    USING (
        EXISTS (
            SELECT 1 FROM commands 
            WHERE commands.id = command_items.command_id 
            AND commands.tenant_id = current_setting('app.current_tenant_id', TRUE)::UUID
        )
    );

-- Policy para command_payments (via command)
CREATE POLICY command_payments_tenant_isolation ON command_payments
    USING (
        EXISTS (
            SELECT 1 FROM commands 
            WHERE commands.id = command_payments.command_id 
            AND commands.tenant_id = current_setting('app.current_tenant_id', TRUE)::UUID
        )
    );

-- =============================================
-- COMENTÁRIOS
-- =============================================

COMMENT ON TABLE commands IS 'Comandas de atendimento - controle de serviços e pagamentos';
COMMENT ON TABLE command_items IS 'Itens da comanda (serviços, produtos, pacotes)';
COMMENT ON TABLE command_payments IS 'Pagamentos da comanda (múltiplas formas aceitas)';

COMMENT ON COLUMN commands.numero IS 'Número sequencial da comanda para exibição';
COMMENT ON COLUMN commands.deixar_troco_gorjeta IS 'Se true, troco vira gorjeta do profissional';
COMMENT ON COLUMN commands.deixar_saldo_divida IS 'Se true, permite fechar com saldo devedor';
COMMENT ON COLUMN command_items.tipo IS 'Tipo do item: SERVICO, PRODUTO ou PACOTE';
COMMENT ON COLUMN command_payments.valor_liquido IS 'Valor após dedução de taxas';
