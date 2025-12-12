-- Fase 1: vínculos explícitos de contas a receber com comanda/pagamento
-- + normalização de origem para incluir PRODUTO

ALTER TABLE contas_a_receber
    ADD COLUMN IF NOT EXISTS command_id UUID REFERENCES commands(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS command_payment_id UUID REFERENCES command_payments(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_contas_a_receber_command
    ON contas_a_receber(tenant_id, command_id)
    WHERE command_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_contas_a_receber_command_payment
    ON contas_a_receber(tenant_id, command_payment_id)
    WHERE command_payment_id IS NOT NULL;

-- Atualizar constraint de origem para permitir PRODUTO
ALTER TABLE contas_a_receber DROP CONSTRAINT IF EXISTS contas_a_receber_origem_check;
ALTER TABLE contas_a_receber
    ADD CONSTRAINT contas_a_receber_origem_check
    CHECK (origem IN ('ASSINATURA', 'SERVICO', 'PRODUTO', 'OUTRO'));

