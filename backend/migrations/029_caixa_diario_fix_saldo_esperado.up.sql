-- ============================================================
-- Migration: 029_caixa_diario_fix_saldo_esperado
-- Descrição: Remove GENERATED ALWAYS de saldo_esperado para
--            permitir atualização direta via aplicação
-- Data: 2025-11-29
-- Autor: NEXO Team
-- ============================================================

-- Remover a coluna GENERATED e recriar como coluna normal
ALTER TABLE caixa_diario DROP COLUMN saldo_esperado;

ALTER TABLE caixa_diario 
    ADD COLUMN saldo_esperado DECIMAL(15,2) NOT NULL DEFAULT 0;

-- Atualizar valores existentes (recalcular)
UPDATE caixa_diario
SET saldo_esperado = saldo_inicial + total_entradas - total_sangrias + total_reforcos;
