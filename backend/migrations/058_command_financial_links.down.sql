-- Reverter vínculos explícitos e constraint de origem

ALTER TABLE contas_a_receber DROP CONSTRAINT IF EXISTS contas_a_receber_origem_check;
ALTER TABLE contas_a_receber
    ADD CONSTRAINT contas_a_receber_origem_check
    CHECK (origem IN ('ASSINATURA', 'SERVICO', 'OUTRO'));

DROP INDEX IF EXISTS idx_contas_a_receber_command_payment;
DROP INDEX IF EXISTS idx_contas_a_receber_command;

ALTER TABLE contas_a_receber
    DROP COLUMN IF EXISTS command_payment_id,
    DROP COLUMN IF EXISTS command_id;

