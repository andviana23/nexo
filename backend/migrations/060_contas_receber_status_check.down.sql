-- 060 - Rollback constraint de status de contas_a_receber (legado)

ALTER TABLE contas_a_receber
    DROP CONSTRAINT IF EXISTS contas_a_receber_status_check;

ALTER TABLE contas_a_receber
    ADD CONSTRAINT contas_a_receber_status_check
    CHECK (status IN (
        'PENDENTE',
        'RECEBIDO',
        'ATRASADO',
        'CANCELADO'
    ));

