-- ============================================================================
-- Schema: despesas_fixas
-- Módulo: Financeiro - Sprint 2 (Despesas Fixas)
-- Descrição: Despesas recorrentes que geram contas a pagar automaticamente
-- ============================================================================

CREATE TABLE IF NOT EXISTS despesas_fixas (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unidade_id    UUID REFERENCES tenants(id) ON DELETE SET NULL,  -- Opcional
    
    descricao     VARCHAR(255) NOT NULL,
    categoria_id  UUID REFERENCES categorias(id) ON DELETE SET NULL,
    fornecedor    VARCHAR(255),
    valor         DECIMAL(12, 2) NOT NULL CHECK (valor > 0),
    
    dia_vencimento INTEGER NOT NULL CHECK (dia_vencimento BETWEEN 1 AND 31),
    ativo          BOOLEAN DEFAULT true,
    
    observacoes   TEXT,
    
    criado_em     TIMESTAMPTZ DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ DEFAULT NOW()
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_despesas_fixas_tenant ON despesas_fixas(tenant_id);
CREATE INDEX IF NOT EXISTS idx_despesas_fixas_ativo ON despesas_fixas(tenant_id, ativo);
CREATE INDEX IF NOT EXISTS idx_despesas_fixas_unidade ON despesas_fixas(unidade_id) WHERE unidade_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_despesas_fixas_categoria ON despesas_fixas(categoria_id) WHERE categoria_id IS NOT NULL;

-- Comentários
COMMENT ON TABLE despesas_fixas IS 'Despesas fixas recorrentes que geram contas a pagar mensalmente';
COMMENT ON COLUMN despesas_fixas.dia_vencimento IS 'Dia do mês (1-31) para vencimento. Ajustado automaticamente para meses com menos dias';
COMMENT ON COLUMN despesas_fixas.ativo IS 'Se false, não gera conta no próximo ciclo';
