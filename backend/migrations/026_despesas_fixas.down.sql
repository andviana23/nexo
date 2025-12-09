-- Migration: 008_despesas_fixas (DOWN)
-- Descrição: Remove tabela de despesas fixas
-- Data: 2025-11-29

-- Remove RLS policy
DROP POLICY IF EXISTS despesas_fixas_tenant_isolation ON despesas_fixas;

-- Remove trigger
DROP TRIGGER IF EXISTS update_despesas_fixas_updated_at ON despesas_fixas;

-- Remove índices (serão removidos automaticamente com a tabela, mas explicitamos)
DROP INDEX IF EXISTS idx_despesas_fixas_tenant;
DROP INDEX IF EXISTS idx_despesas_fixas_ativo;
DROP INDEX IF EXISTS idx_despesas_fixas_unidade;
DROP INDEX IF EXISTS idx_despesas_fixas_categoria;

-- Remove tabela
DROP TABLE IF EXISTS despesas_fixas;
