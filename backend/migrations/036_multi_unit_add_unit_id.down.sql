-- ============================================================================
-- MIGRATION 036: Multi-Unidade - Rollback unit_id das tabelas operacionais
-- ============================================================================

ALTER TABLE despesas_fixas DROP COLUMN IF EXISTS unit_id;
ALTER TABLE metas_mensais DROP COLUMN IF EXISTS unit_id;
ALTER TABLE blocked_times DROP COLUMN IF EXISTS unit_id;
ALTER TABLE barber_turns DROP COLUMN IF EXISTS unit_id;
ALTER TABLE produtos DROP COLUMN IF EXISTS unit_id;
ALTER TABLE movimentacoes_estoque DROP COLUMN IF EXISTS unit_id;
ALTER TABLE dre_mensal DROP COLUMN IF EXISTS unit_id;
ALTER TABLE fluxo_caixa_diario DROP COLUMN IF EXISTS unit_id;
ALTER TABLE compensacoes_bancarias DROP COLUMN IF EXISTS unit_id;
ALTER TABLE contas_a_receber DROP COLUMN IF EXISTS unit_id;
ALTER TABLE contas_a_pagar DROP COLUMN IF EXISTS unit_id;
ALTER TABLE caixa_diario DROP COLUMN IF EXISTS unit_id;
ALTER TABLE commands DROP COLUMN IF EXISTS unit_id;
ALTER TABLE appointment_status_history DROP COLUMN IF EXISTS unit_id;
ALTER TABLE appointments DROP COLUMN IF EXISTS unit_id;
