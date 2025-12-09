-- Migration: 045_commission_items (rollback)
-- Description: Remove tabela de itens de comissão

DROP INDEX IF EXISTS idx_commission_items_command_item_unique;
DROP INDEX IF EXISTS idx_commission_items_command_item;
DROP INDEX IF EXISTS idx_commission_items_command;
DROP INDEX IF EXISTS idx_commission_items_reference_date;
DROP INDEX IF EXISTS idx_commission_items_pending;
DROP INDEX IF EXISTS idx_commission_items_status;
DROP INDEX IF EXISTS idx_commission_items_period;
DROP INDEX IF EXISTS idx_commission_items_professional;
DROP INDEX IF EXISTS idx_commission_items_unit;
DROP INDEX IF EXISTS idx_commission_items_tenant;
DROP TABLE IF EXISTS commission_items;
