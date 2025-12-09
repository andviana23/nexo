-- Migration: 008_despesas_fixas
-- Descrição: Cria tabela de despesas fixas (recorrentes) para o módulo financeiro
-- Data: 2025-11-29
-- Sprint: 2 - Módulo Financeiro

-- ============================================================================
-- TABELA: despesas_fixas
-- Gerencia despesas recorrentes que serão automaticamente convertidas em
-- contas a pagar no início de cada mês via cron job.
-- ============================================================================

CREATE TABLE IF NOT EXISTS despesas_fixas (
    -- Identificação
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    -- unidade_id será adicionado quando tabela unidades existir
    unidade_id UUID,
    
    -- Dados da despesa
    descricao VARCHAR(255) NOT NULL,
    categoria_id UUID REFERENCES categorias(id) ON DELETE SET NULL,
    fornecedor VARCHAR(255),
    valor DECIMAL(15,2) NOT NULL CHECK (valor > 0),
    
    -- Configuração de recorrência
    dia_vencimento INTEGER NOT NULL CHECK (dia_vencimento BETWEEN 1 AND 31),
    ativo BOOLEAN NOT NULL DEFAULT true,
    
    -- Metadados
    observacoes TEXT,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- ÍNDICES
-- ============================================================================

-- Índice principal para isolamento multi-tenant
CREATE INDEX idx_despesas_fixas_tenant ON despesas_fixas(tenant_id);

-- Índice para buscar despesas ativas (usado pelo cron job)
CREATE INDEX idx_despesas_fixas_ativo ON despesas_fixas(tenant_id, ativo) 
    WHERE ativo = true;

-- Índice para filtrar por unidade
CREATE INDEX idx_despesas_fixas_unidade ON despesas_fixas(unidade_id) 
    WHERE unidade_id IS NOT NULL;

-- Índice para filtrar por categoria
CREATE INDEX idx_despesas_fixas_categoria ON despesas_fixas(categoria_id) 
    WHERE categoria_id IS NOT NULL;

-- ============================================================================
-- ROW LEVEL SECURITY (RLS)
-- ============================================================================

ALTER TABLE despesas_fixas ENABLE ROW LEVEL SECURITY;

-- Policy para isolamento multi-tenant
CREATE POLICY despesas_fixas_tenant_isolation ON despesas_fixas
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant', true)::uuid);

-- ============================================================================
-- TRIGGER para atualizar updated_at automaticamente
-- ============================================================================

CREATE TRIGGER update_despesas_fixas_updated_at
    BEFORE UPDATE ON despesas_fixas
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- COMENTÁRIOS
-- ============================================================================

COMMENT ON TABLE despesas_fixas IS 'Despesas recorrentes que geram contas a pagar automaticamente no início de cada mês';
COMMENT ON COLUMN despesas_fixas.dia_vencimento IS 'Dia do mês para vencimento (1-31). Para meses com menos dias, ajusta para o último dia';
COMMENT ON COLUMN despesas_fixas.ativo IS 'Se false, não gera conta no próximo mês';
