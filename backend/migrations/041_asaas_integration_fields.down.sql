-- Rollback Migration 041: Campos para integração Asaas x NEXO

-- ============================================================================
-- 1) Remover tabelas de log
-- ============================================================================
DROP TABLE IF EXISTS asaas_reconciliation_logs;
DROP TABLE IF EXISTS asaas_webhook_logs;

-- ============================================================================
-- 2) Remover campos de fluxo_caixa_diario
-- ============================================================================
ALTER TABLE fluxo_caixa_diario DROP COLUMN IF EXISTS asaas_payments_count;
ALTER TABLE fluxo_caixa_diario DROP COLUMN IF EXISTS asaas_payments_total;

-- ============================================================================
-- 3) Remover campos e índices de contas_a_receber
-- ============================================================================
DROP INDEX IF EXISTS idx_contas_a_receber_competencia;
DROP INDEX IF EXISTS idx_contas_a_receber_asaas_payment_id;

ALTER TABLE contas_a_receber DROP COLUMN IF EXISTS received_at;
ALTER TABLE contas_a_receber DROP COLUMN IF EXISTS confirmed_at;
ALTER TABLE contas_a_receber DROP COLUMN IF EXISTS competencia_mes;
ALTER TABLE contas_a_receber DROP COLUMN IF EXISTS asaas_payment_id;
ALTER TABLE contas_a_receber DROP COLUMN IF EXISTS subscription_id;

-- ============================================================================
-- 4) Remover campos e índices de subscription_payments
-- ============================================================================
DROP INDEX IF EXISTS idx_subscription_payments_asaas_payment_id;

ALTER TABLE subscription_payments DROP COLUMN IF EXISTS pix_qr_code;
ALTER TABLE subscription_payments DROP COLUMN IF EXISTS bank_slip_url;
ALTER TABLE subscription_payments DROP COLUMN IF EXISTS invoice_url;
ALTER TABLE subscription_payments DROP COLUMN IF EXISTS net_value;
ALTER TABLE subscription_payments DROP COLUMN IF EXISTS billing_type;
ALTER TABLE subscription_payments DROP COLUMN IF EXISTS estimated_credit_date;
ALTER TABLE subscription_payments DROP COLUMN IF EXISTS credit_date;
ALTER TABLE subscription_payments DROP COLUMN IF EXISTS client_payment_date;
ALTER TABLE subscription_payments DROP COLUMN IF EXISTS confirmed_date;
ALTER TABLE subscription_payments DROP COLUMN IF EXISTS due_date;
ALTER TABLE subscription_payments DROP COLUMN IF EXISTS status_asaas;

-- Restaurar constraint original
ALTER TABLE subscription_payments DROP CONSTRAINT IF EXISTS subscription_payments_status_check;
ALTER TABLE subscription_payments ADD CONSTRAINT subscription_payments_status_check 
CHECK (status IN ('PENDENTE', 'CONFIRMADO', 'ESTORNADO', 'CANCELADO'));

-- ============================================================================
-- 5) Remover campos de subscriptions
-- ============================================================================
ALTER TABLE subscriptions DROP COLUMN IF EXISTS last_sync_at;
ALTER TABLE subscriptions DROP COLUMN IF EXISTS last_confirmed_at;
ALTER TABLE subscriptions DROP COLUMN IF EXISTS asaas_status;
ALTER TABLE subscriptions DROP COLUMN IF EXISTS cycle;
ALTER TABLE subscriptions DROP COLUMN IF EXISTS next_due_date;

-- Restaurar constraint original
ALTER TABLE subscriptions DROP CONSTRAINT IF EXISTS subscriptions_status_check;
ALTER TABLE subscriptions ADD CONSTRAINT subscriptions_status_check 
CHECK (status IN ('AGUARDANDO_PAGAMENTO', 'ATIVO', 'INADIMPLENTE', 'INATIVO', 'CANCELADO'));
