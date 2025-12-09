-- ============================================================================
-- ROLLBACK: Remover tabela categorias_produtos
-- ============================================================================

-- Remove FK da tabela produtos
ALTER TABLE produtos DROP COLUMN IF EXISTS categoria_produto_id;

-- Remove trigger
DROP TRIGGER IF EXISTS trg_categorias_produtos_updated_at ON categorias_produtos;

-- Remove tabela
DROP TABLE IF EXISTS categorias_produtos;
