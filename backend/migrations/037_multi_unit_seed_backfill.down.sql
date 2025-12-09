-- ============================================================================
-- MIGRATION 037: Multi-Unidade - Rollback Seed/Backfill
-- ATENÇÃO: Este rollback remove dados! Use com cuidado.
-- ============================================================================

-- Limpar unit_id das tabelas (setar para NULL)
UPDATE despesas_fixas SET unit_id = NULL;
UPDATE metas_mensais SET unit_id = NULL;
UPDATE blocked_times SET unit_id = NULL;
UPDATE barber_turns SET unit_id = NULL;
UPDATE movimentacoes_estoque SET unit_id = NULL;
UPDATE dre_mensal SET unit_id = NULL;
UPDATE fluxo_caixa_diario SET unit_id = NULL;
UPDATE compensacoes_bancarias SET unit_id = NULL;
UPDATE contas_a_receber SET unit_id = NULL;
UPDATE contas_a_pagar SET unit_id = NULL;
UPDATE caixa_diario SET unit_id = NULL;
UPDATE commands SET unit_id = NULL;
UPDATE appointment_status_history SET unit_id = NULL;
UPDATE appointments SET unit_id = NULL;

-- Remover vínculos usuário-unidade
DELETE FROM user_units;

-- Remover unidades "Principal" criadas automaticamente
DELETE FROM units WHERE descricao LIKE '%migração multi-unidade%';
