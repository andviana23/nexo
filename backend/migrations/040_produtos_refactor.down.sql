-- ============================================================================
-- Migration: 040_produtos_refactor (ROLLBACK)
-- Descrição: Reverte a refatoração da tabela produtos
-- Data: 2025-12-02
-- ============================================================================

-- Remover índices
DROP INDEX IF EXISTS idx_produtos_fornecedor_id;
DROP INDEX IF EXISTS idx_produtos_categoria_produto_id;
DROP INDEX IF EXISTS idx_produtos_codigo_barras;

-- Remover colunas adicionadas
ALTER TABLE produtos DROP COLUMN IF EXISTS estoque_maximo;
ALTER TABLE produtos DROP COLUMN IF EXISTS valor_venda_profissional;
ALTER TABLE produtos DROP COLUMN IF EXISTS valor_entrada;
-- Nota: fornecedor_id e categoria_produto_id podem já existir, não removemos
