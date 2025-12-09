-- =====================================================
-- Migration: Reverter command_id em appointments
-- Descrição: Remove coluna command_id de appointments
-- Data: 30/11/2025
-- =====================================================

BEGIN;

-- 1. Remover índice
DROP INDEX IF EXISTS idx_appointments_command_id;

-- 2. Remover coluna command_id
ALTER TABLE appointments
DROP COLUMN IF EXISTS command_id;

COMMIT;
