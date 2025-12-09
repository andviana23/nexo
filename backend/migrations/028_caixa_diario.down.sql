-- ============================================================
-- Migration: 028_caixa_diario (DOWN)
-- Descrição: Remove tabelas do módulo Caixa Diário
-- Data: 2025-11-29
-- ============================================================

-- Remover trigger primeiro
DROP TRIGGER IF EXISTS trg_caixa_diario_updated_at ON caixa_diario;

-- Remover índices de operacoes_caixa
DROP INDEX IF EXISTS idx_operacoes_caixa_created_at;
DROP INDEX IF EXISTS idx_operacoes_caixa_tenant_tipo;
DROP INDEX IF EXISTS idx_operacoes_caixa_caixa_id;

-- Remover tabela operacoes_caixa
DROP TABLE IF EXISTS operacoes_caixa;

-- Remover índices de caixa_diario
DROP INDEX IF EXISTS idx_caixa_diario_aberto_unico;
DROP INDEX IF EXISTS idx_caixa_diario_tenant_data_abertura;
DROP INDEX IF EXISTS idx_caixa_diario_tenant_status;

-- Remover tabela caixa_diario
DROP TABLE IF EXISTS caixa_diario;
