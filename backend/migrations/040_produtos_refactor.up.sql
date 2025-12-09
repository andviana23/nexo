-- ============================================================================
-- Migration: 040_produtos_refactor
-- Descrição: Refatoração da tabela produtos para novo modal de cadastro
-- Data: 2025-12-02
-- ============================================================================

-- Adicionar novos campos à tabela produtos
ALTER TABLE produtos
  -- Código de barras já existe, garantir que aceita NULL
  ALTER COLUMN codigo_barras DROP NOT NULL;

-- Adicionar estoque_maximo se não existir
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns 
    WHERE table_name = 'produtos' AND column_name = 'estoque_maximo'
  ) THEN
    ALTER TABLE produtos ADD COLUMN estoque_maximo NUMERIC(15,3) DEFAULT NULL;
  END IF;
END $$;

-- Adicionar valor_venda_profissional (preço de venda para profissionais)
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns 
    WHERE table_name = 'produtos' AND column_name = 'valor_venda_profissional'
  ) THEN
    ALTER TABLE produtos ADD COLUMN valor_venda_profissional NUMERIC(10,2) DEFAULT NULL;
  END IF;
END $$;

-- Adicionar valor_entrada (custo de entrada/compra)
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns 
    WHERE table_name = 'produtos' AND column_name = 'valor_entrada'
  ) THEN
    ALTER TABLE produtos ADD COLUMN valor_entrada NUMERIC(10,2) DEFAULT NULL;
  END IF;
END $$;

-- Adicionar FK para fornecedor (se não existir)
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns 
    WHERE table_name = 'produtos' AND column_name = 'fornecedor_id'
  ) THEN
    ALTER TABLE produtos ADD COLUMN fornecedor_id UUID REFERENCES fornecedores(id) ON DELETE SET NULL;
  END IF;
END $$;

-- Adicionar FK para categoria_produto (categorias customizadas)
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns 
    WHERE table_name = 'produtos' AND column_name = 'categoria_produto_id'
  ) THEN
    ALTER TABLE produtos ADD COLUMN categoria_produto_id UUID REFERENCES categorias_produtos(id) ON DELETE SET NULL;
  END IF;
END $$;

-- Criar índices para os novos campos
CREATE INDEX IF NOT EXISTS idx_produtos_fornecedor_id ON produtos(fornecedor_id);
CREATE INDEX IF NOT EXISTS idx_produtos_categoria_produto_id ON produtos(categoria_produto_id);
CREATE INDEX IF NOT EXISTS idx_produtos_codigo_barras ON produtos(tenant_id, codigo_barras) WHERE codigo_barras IS NOT NULL;

-- Comentários nos campos
COMMENT ON COLUMN produtos.codigo_barras IS 'Código de barras do produto (EAN-13, EAN-8, etc.) - Opcional';
COMMENT ON COLUMN produtos.estoque_maximo IS 'Quantidade máxima de estoque (para alertas de excesso)';
COMMENT ON COLUMN produtos.valor_venda_profissional IS 'Valor de venda para profissionais da barbearia';
COMMENT ON COLUMN produtos.valor_entrada IS 'Valor de entrada/custo de compra do produto';
COMMENT ON COLUMN produtos.fornecedor_id IS 'FK para fornecedor principal do produto';
COMMENT ON COLUMN produtos.categoria_produto_id IS 'FK para categoria customizada de produto';
