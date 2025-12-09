-- ============================================================
-- Migration: 029_caixa_diario_fix_saldo_esperado (DOWN)
-- ============================================================

-- Reverter para GENERATED ALWAYS
ALTER TABLE caixa_diario DROP COLUMN saldo_esperado;

ALTER TABLE caixa_diario 
    ADD COLUMN saldo_esperado DECIMAL(15,2) GENERATED ALWAYS AS (
        saldo_inicial + total_entradas - total_sangrias + total_reforcos
    ) STORED;
