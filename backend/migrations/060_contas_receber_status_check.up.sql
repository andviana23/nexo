-- 060 - Atualizar constraint de status de contas_a_receber
-- Motivo: Fase 1/2 introduz status canônicos CONFIRMADO/ESTORNADO e compat PAGO.

ALTER TABLE contas_a_receber
    DROP CONSTRAINT IF EXISTS contas_a_receber_status_check;

ALTER TABLE contas_a_receber
    ADD CONSTRAINT contas_a_receber_status_check
    CHECK (status IN (
        'PENDENTE',
        'CONFIRMADO',
        'RECEBIDO',
        'PAGO',       -- legado (lido como RECEBIDO)
        'ATRASADO',
        'ESTORNADO',
        'CANCELADO'
    ));

