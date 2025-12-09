-- =============================================
-- Migration Rollback: 027_commands_and_payments
-- Descrição: Remove tabelas de comandas e pagamentos
-- Data: 29/11/2025
-- =============================================

-- Remove policies RLS
DROP POLICY IF EXISTS command_payments_tenant_isolation ON command_payments;
DROP POLICY IF EXISTS command_items_tenant_isolation ON command_items;
DROP POLICY IF EXISTS commands_tenant_isolation ON commands;

-- Remove trigger
DROP TRIGGER IF EXISTS trigger_commands_updated_at ON commands;
DROP FUNCTION IF EXISTS update_commands_updated_at();

-- Remove índices (serão removidos automaticamente com DROP TABLE CASCADE)
-- Mas listamos para documentação
-- DROP INDEX IF EXISTS idx_commands_tenant_id;
-- DROP INDEX IF EXISTS idx_commands_customer_id;
-- DROP INDEX IF EXISTS idx_commands_appointment_id;
-- DROP INDEX IF EXISTS idx_commands_status;
-- DROP INDEX IF EXISTS idx_commands_criado_em;
-- DROP INDEX IF EXISTS idx_commands_numero;
-- DROP INDEX IF EXISTS idx_command_items_command_id;
-- DROP INDEX IF EXISTS idx_command_items_tipo;
-- DROP INDEX IF EXISTS idx_command_items_item_id;
-- DROP INDEX IF EXISTS idx_command_payments_command_id;
-- DROP INDEX IF EXISTS idx_command_payments_meio_pagamento;
-- DROP INDEX IF EXISTS idx_command_payments_criado_em;

-- Remove tabelas (ordem reversa de criação)
DROP TABLE IF EXISTS command_payments CASCADE;
DROP TABLE IF EXISTS command_items CASCADE;
DROP TABLE IF EXISTS commands CASCADE;
