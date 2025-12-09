-- Migration 041: Campos para integração completa Asaas x NEXO
-- Objetivo: alinhar modelagem ao comportamento oficial do Asaas
-- Ref: PLANO_AJUSTE_ASAAS.md

-- ============================================================================
-- 1) SUBSCRIPTIONS: campos de sincronização com Asaas
-- ============================================================================

-- next_due_date: próxima data de vencimento (vem do Asaas nextDueDate)
ALTER TABLE subscriptions 
ADD COLUMN IF NOT EXISTS next_due_date DATE;

-- cycle: ciclo de cobrança (MONTHLY, WEEKLY, BIWEEKLY, QUARTERLY, etc)
ALTER TABLE subscriptions 
ADD COLUMN IF NOT EXISTS cycle VARCHAR(20) DEFAULT 'MONTHLY';

-- asaas_status: status real no Asaas (ACTIVE, INACTIVE, EXPIRED, etc)
ALTER TABLE subscriptions 
ADD COLUMN IF NOT EXISTS asaas_status VARCHAR(30);

-- last_confirmed_at: data do último pagamento confirmado
ALTER TABLE subscriptions 
ADD COLUMN IF NOT EXISTS last_confirmed_at TIMESTAMP WITH TIME ZONE;

-- last_sync_at: última sincronização com Asaas
ALTER TABLE subscriptions 
ADD COLUMN IF NOT EXISTS last_sync_at TIMESTAMP WITH TIME ZONE;

COMMENT ON COLUMN subscriptions.next_due_date IS 'Próxima data de vencimento (Asaas nextDueDate)';
COMMENT ON COLUMN subscriptions.cycle IS 'Ciclo: MONTHLY, WEEKLY, BIWEEKLY, QUARTERLY, SEMIANNUALLY, YEARLY';
COMMENT ON COLUMN subscriptions.asaas_status IS 'Status no Asaas: ACTIVE, INACTIVE, EXPIRED';
COMMENT ON COLUMN subscriptions.last_confirmed_at IS 'Data do último pagamento CONFIRMED';
COMMENT ON COLUMN subscriptions.last_sync_at IS 'Última sincronização com API Asaas';

-- ============================================================================
-- 2) SUBSCRIPTION_PAYMENTS: campos completos de pagamento
-- ============================================================================

-- status_asaas: status original do Asaas (preservar para auditoria)
ALTER TABLE subscription_payments 
ADD COLUMN IF NOT EXISTS status_asaas VARCHAR(30);

-- due_date: data de vencimento original da cobrança
ALTER TABLE subscription_payments 
ADD COLUMN IF NOT EXISTS due_date DATE;

-- confirmed_date: data em que foi confirmado (competência para DRE)
ALTER TABLE subscription_payments 
ADD COLUMN IF NOT EXISTS confirmed_date TIMESTAMP WITH TIME ZONE;

-- client_payment_date: data em que o cliente efetuou o pagamento
ALTER TABLE subscription_payments 
ADD COLUMN IF NOT EXISTS client_payment_date DATE;

-- credit_date: data em que o dinheiro caiu na conta (regime de caixa)
ALTER TABLE subscription_payments 
ADD COLUMN IF NOT EXISTS credit_date DATE;

-- estimated_credit_date: previsão de compensação (D+N)
ALTER TABLE subscription_payments 
ADD COLUMN IF NOT EXISTS estimated_credit_date DATE;

-- billing_type: tipo de cobrança Asaas (CREDIT_CARD, PIX, BOLETO, etc)
ALTER TABLE subscription_payments 
ADD COLUMN IF NOT EXISTS billing_type VARCHAR(30);

-- net_value: valor líquido após taxas Asaas
ALTER TABLE subscription_payments 
ADD COLUMN IF NOT EXISTS net_value NUMERIC(10,2);

-- invoice_url: link da fatura/boleto
ALTER TABLE subscription_payments 
ADD COLUMN IF NOT EXISTS invoice_url TEXT;

-- bank_slip_url: link do boleto (se aplicável)
ALTER TABLE subscription_payments 
ADD COLUMN IF NOT EXISTS bank_slip_url TEXT;

-- pix_qr_code: código PIX copia e cola
ALTER TABLE subscription_payments 
ADD COLUMN IF NOT EXISTS pix_qr_code TEXT;

-- Índice único para idempotência de webhooks
CREATE UNIQUE INDEX IF NOT EXISTS idx_subscription_payments_asaas_payment_id 
ON subscription_payments(asaas_payment_id) 
WHERE asaas_payment_id IS NOT NULL;

COMMENT ON COLUMN subscription_payments.status_asaas IS 'Status original Asaas: PENDING, CONFIRMED, RECEIVED, OVERDUE, REFUNDED, etc';
COMMENT ON COLUMN subscription_payments.due_date IS 'Data de vencimento original da cobrança';
COMMENT ON COLUMN subscription_payments.confirmed_date IS 'Data de confirmação (competência DRE)';
COMMENT ON COLUMN subscription_payments.client_payment_date IS 'Data em que cliente pagou';
COMMENT ON COLUMN subscription_payments.credit_date IS 'Data que o dinheiro caiu na conta (regime caixa)';
COMMENT ON COLUMN subscription_payments.estimated_credit_date IS 'Previsão de compensação (D+N conforme meio)';
COMMENT ON COLUMN subscription_payments.billing_type IS 'Tipo Asaas: CREDIT_CARD, PIX, BOLETO, UNDEFINED';
COMMENT ON COLUMN subscription_payments.net_value IS 'Valor líquido após taxas Asaas';

-- ============================================================================
-- 3) CONTAS_A_RECEBER: vincular a subscriptions e asaas_payment_id
-- ============================================================================

-- subscription_id: FK para subscriptions (nova tabela)
ALTER TABLE contas_a_receber 
ADD COLUMN IF NOT EXISTS subscription_id UUID REFERENCES subscriptions(id) ON DELETE SET NULL;

-- asaas_payment_id: vincular à cobrança específica do Asaas
ALTER TABLE contas_a_receber 
ADD COLUMN IF NOT EXISTS asaas_payment_id VARCHAR(100);

-- competencia_mes: mês de competência (YYYY-MM) para DRE
ALTER TABLE contas_a_receber 
ADD COLUMN IF NOT EXISTS competencia_mes VARCHAR(7);

-- confirmed_at: quando foi confirmado
ALTER TABLE contas_a_receber 
ADD COLUMN IF NOT EXISTS confirmed_at TIMESTAMP WITH TIME ZONE;

-- received_at: quando foi recebido (creditado)
ALTER TABLE contas_a_receber 
ADD COLUMN IF NOT EXISTS received_at TIMESTAMP WITH TIME ZONE;

-- Índice para evitar duplicatas de cobrança Asaas
CREATE UNIQUE INDEX IF NOT EXISTS idx_contas_a_receber_asaas_payment_id 
ON contas_a_receber(tenant_id, asaas_payment_id) 
WHERE asaas_payment_id IS NOT NULL;

-- Índice para consulta por competência (DRE)
CREATE INDEX IF NOT EXISTS idx_contas_a_receber_competencia 
ON contas_a_receber(tenant_id, competencia_mes, status);

COMMENT ON COLUMN contas_a_receber.subscription_id IS 'FK para subscriptions (assinaturas de clientes)';
COMMENT ON COLUMN contas_a_receber.asaas_payment_id IS 'ID da cobrança no Asaas (idempotência)';
COMMENT ON COLUMN contas_a_receber.competencia_mes IS 'Mês de competência YYYY-MM para DRE';
COMMENT ON COLUMN contas_a_receber.confirmed_at IS 'Timestamp de confirmação (CONFIRMED)';
COMMENT ON COLUMN contas_a_receber.received_at IS 'Timestamp de recebimento (RECEIVED)';

-- ============================================================================
-- 4) FLUXO_CAIXA_DIARIO: campos para rastreabilidade
-- ============================================================================

-- asaas_payments_count: quantidade de pagamentos Asaas no dia
ALTER TABLE fluxo_caixa_diario 
ADD COLUMN IF NOT EXISTS asaas_payments_count INTEGER DEFAULT 0;

-- asaas_payments_total: valor total de pagamentos Asaas no dia
ALTER TABLE fluxo_caixa_diario 
ADD COLUMN IF NOT EXISTS asaas_payments_total NUMERIC(15,2) DEFAULT 0;

COMMENT ON COLUMN fluxo_caixa_diario.asaas_payments_count IS 'Qtd de pagamentos Asaas RECEIVED no dia';
COMMENT ON COLUMN fluxo_caixa_diario.asaas_payments_total IS 'Total de pagamentos Asaas RECEIVED no dia';

-- ============================================================================
-- 5) TABELA DE LOG DE WEBHOOKS (auditoria e debug)
-- ============================================================================

CREATE TABLE IF NOT EXISTS asaas_webhook_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Dados do evento
    event_type VARCHAR(50) NOT NULL,
    asaas_payment_id VARCHAR(100),
    asaas_subscription_id VARCHAR(100),
    
    -- Payload completo (para debug)
    payload JSONB NOT NULL,
    
    -- Resultado do processamento
    processed BOOLEAN DEFAULT FALSE,
    processed_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    
    -- Rastreabilidade
    received_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address INET,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_asaas_webhook_logs_payment 
ON asaas_webhook_logs(asaas_payment_id);

CREATE INDEX IF NOT EXISTS idx_asaas_webhook_logs_subscription 
ON asaas_webhook_logs(asaas_subscription_id);

CREATE INDEX IF NOT EXISTS idx_asaas_webhook_logs_event 
ON asaas_webhook_logs(tenant_id, event_type, received_at DESC);

COMMENT ON TABLE asaas_webhook_logs IS 'Log de webhooks Asaas para auditoria e debug';
COMMENT ON COLUMN asaas_webhook_logs.event_type IS 'Tipo: PAYMENT_CREATED, PAYMENT_CONFIRMED, PAYMENT_RECEIVED, etc';
COMMENT ON COLUMN asaas_webhook_logs.payload IS 'Payload JSON completo do webhook';
COMMENT ON COLUMN asaas_webhook_logs.processed IS 'Se foi processado com sucesso';

-- ============================================================================
-- 6) TABELA DE CONCILIAÇÃO (divergências)
-- ============================================================================

CREATE TABLE IF NOT EXISTS asaas_reconciliation_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Período da conciliação
    reference_date DATE NOT NULL,
    
    -- Tipo de divergência
    divergence_type VARCHAR(50) NOT NULL,
    -- MISSING_PAYMENT: payment no Asaas mas não no NEXO
    -- MISSING_SUBSCRIPTION: subscription no Asaas mas não no NEXO  
    -- STATUS_MISMATCH: status diferente entre sistemas
    -- VALUE_MISMATCH: valor diferente entre sistemas
    
    -- Referências
    asaas_payment_id VARCHAR(100),
    asaas_subscription_id VARCHAR(100),
    nexo_payment_id UUID,
    nexo_subscription_id UUID,
    
    -- Detalhes
    asaas_value NUMERIC(10,2),
    nexo_value NUMERIC(10,2),
    asaas_status VARCHAR(30),
    nexo_status VARCHAR(30),
    
    -- Resolução
    resolved BOOLEAN DEFAULT FALSE,
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolved_by UUID REFERENCES users(id),
    resolution_notes TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_asaas_reconciliation_pending 
ON asaas_reconciliation_logs(tenant_id, resolved, created_at DESC) 
WHERE resolved = FALSE;

COMMENT ON TABLE asaas_reconciliation_logs IS 'Log de divergências entre Asaas e NEXO';
COMMENT ON COLUMN asaas_reconciliation_logs.divergence_type IS 'Tipo: MISSING_PAYMENT, MISSING_SUBSCRIPTION, STATUS_MISMATCH, VALUE_MISMATCH';

-- ============================================================================
-- 7) ENUM DE STATUS ASAAS (para consistência)
-- ============================================================================

-- Atualizar constraint de status em subscription_payments
ALTER TABLE subscription_payments 
DROP CONSTRAINT IF EXISTS subscription_payments_status_check;

ALTER TABLE subscription_payments 
ADD CONSTRAINT subscription_payments_status_check 
CHECK (status IN (
    'PENDENTE',      -- Aguardando pagamento
    'CONFIRMADO',    -- Pagamento confirmado (competência)
    'RECEBIDO',      -- Dinheiro na conta (caixa)
    'VENCIDO',       -- Passou do vencimento sem pagar
    'ESTORNADO',     -- Estornado/Refundido
    'CANCELADO'      -- Cancelado
));

-- Atualizar constraint de status em subscriptions
ALTER TABLE subscriptions 
DROP CONSTRAINT IF EXISTS subscriptions_status_check;

ALTER TABLE subscriptions 
ADD CONSTRAINT subscriptions_status_check 
CHECK (status IN (
    'AGUARDANDO_PAGAMENTO',  -- Criada, aguardando 1º pagamento
    'ATIVO',                 -- Pagamentos em dia
    'INADIMPLENTE',          -- Com pagamento vencido
    'INATIVO',               -- Pausada/suspensa
    'CANCELADO'              -- Cancelada definitivamente
));
