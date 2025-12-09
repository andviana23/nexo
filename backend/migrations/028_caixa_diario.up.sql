-- ============================================================
-- Migration: 028_caixa_diario
-- Descrição: Cria tabelas para o módulo Caixa Diário
-- Data: 2025-11-29
-- Autor: NEXO Team
-- ============================================================

-- ============================================================
-- TABELA: caixa_diario
-- Controle operacional da gaveta de dinheiro
-- ============================================================

CREATE TABLE IF NOT EXISTS caixa_diario (
    -- Identificadores
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Operadores
    usuario_abertura_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    usuario_fechamento_id UUID REFERENCES users(id) ON DELETE RESTRICT,
    
    -- Timestamps
    data_abertura TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    data_fechamento TIMESTAMPTZ,
    
    -- Valores Monetários (DECIMAL para precisão financeira)
    saldo_inicial DECIMAL(15,2) NOT NULL CHECK (saldo_inicial >= 0),
    total_entradas DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (total_entradas >= 0),
    total_saidas DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (total_saidas >= 0),
    total_sangrias DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (total_sangrias >= 0),
    total_reforcos DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (total_reforcos >= 0),
    
    -- Saldo esperado = inicial + entradas - sangrias + reforços
    saldo_esperado DECIMAL(15,2) GENERATED ALWAYS AS (
        saldo_inicial + total_entradas - total_sangrias + total_reforcos
    ) STORED,
    
    -- Fechamento
    saldo_real DECIMAL(15,2),
    divergencia DECIMAL(15,2),
    
    -- Status e Justificativa
    status VARCHAR(20) NOT NULL DEFAULT 'ABERTO' 
        CHECK (status IN ('ABERTO', 'FECHADO')),
    justificativa_divergencia TEXT,
    
    -- Auditoria
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Comentários para documentação
COMMENT ON TABLE caixa_diario IS 'Caixa Diário - Controle operacional da gaveta de dinheiro. Isolamento multi-tenant obrigatório.';
COMMENT ON COLUMN caixa_diario.tenant_id IS 'Isolamento multi-tenant - OBRIGATÓRIO em todas as queries';
COMMENT ON COLUMN caixa_diario.usuario_abertura_id IS 'Usuário que abriu o caixa';
COMMENT ON COLUMN caixa_diario.usuario_fechamento_id IS 'Usuário que fechou o caixa (null se ainda aberto)';
COMMENT ON COLUMN caixa_diario.saldo_inicial IS 'Saldo de abertura informado pelo operador (conferência física)';
COMMENT ON COLUMN caixa_diario.total_entradas IS 'Total de vendas em dinheiro registradas';
COMMENT ON COLUMN caixa_diario.total_sangrias IS 'Total de retiradas (depósitos, pagamentos, etc)';
COMMENT ON COLUMN caixa_diario.total_reforcos IS 'Total de adições (troco, capital de giro, etc)';
COMMENT ON COLUMN caixa_diario.saldo_esperado IS 'Calculado automaticamente: inicial + entradas - sangrias + reforços';
COMMENT ON COLUMN caixa_diario.saldo_real IS 'Valor contado fisicamente no fechamento';
COMMENT ON COLUMN caixa_diario.divergencia IS 'saldo_real - saldo_esperado (negativo = falta, positivo = sobra)';
COMMENT ON COLUMN caixa_diario.status IS 'ABERTO ou FECHADO';
COMMENT ON COLUMN caixa_diario.justificativa_divergencia IS 'Obrigatória se divergência > R$ 5,00';

-- ============================================================
-- ÍNDICES para caixa_diario
-- ============================================================

-- Índice para buscar caixas por tenant e status (muito usado)
CREATE INDEX IF NOT EXISTS idx_caixa_diario_tenant_status 
    ON caixa_diario(tenant_id, status);

-- Índice para histórico ordenado por data de abertura
CREATE INDEX IF NOT EXISTS idx_caixa_diario_tenant_data_abertura 
    ON caixa_diario(tenant_id, data_abertura DESC);

-- Constraint ÚNICA: apenas 1 caixa ABERTO por tenant
-- Isso garante RN-CAI-001: Abertura Única
CREATE UNIQUE INDEX IF NOT EXISTS idx_caixa_diario_aberto_unico 
    ON caixa_diario(tenant_id) 
    WHERE status = 'ABERTO';

-- ============================================================
-- TABELA: operacoes_caixa
-- Registro de todas as operações do caixa (sangrias, reforços, vendas)
-- ============================================================

CREATE TABLE IF NOT EXISTS operacoes_caixa (
    -- Identificadores
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    caixa_id UUID NOT NULL REFERENCES caixa_diario(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Tipo de Operação
    tipo VARCHAR(20) NOT NULL CHECK (tipo IN ('VENDA', 'SANGRIA', 'REFORCO', 'DESPESA')),
    
    -- Valor (sempre positivo, tipo define se é entrada ou saída)
    valor DECIMAL(15,2) NOT NULL CHECK (valor > 0),
    
    -- Detalhes
    descricao TEXT NOT NULL,
    
    -- Para SANGRIA: destino do dinheiro
    destino VARCHAR(50) CHECK (
        tipo != 'SANGRIA' OR destino IN ('DEPOSITO', 'PAGAMENTO', 'COFRE', 'OUTROS')
    ),
    
    -- Para REFORCO: origem do dinheiro
    origem VARCHAR(50) CHECK (
        tipo != 'REFORCO' OR origem IN ('TROCO', 'CAPITAL_GIRO', 'TRANSFERENCIA', 'OUTROS')
    ),
    
    -- Operador que realizou a operação
    usuario_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    
    -- Auditoria
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Comentários para documentação
COMMENT ON TABLE operacoes_caixa IS 'Operações do Caixa Diário - vendas, sangrias, reforços. Isolamento multi-tenant obrigatório.';
COMMENT ON COLUMN operacoes_caixa.caixa_id IS 'Referência ao caixa_diario pai';
COMMENT ON COLUMN operacoes_caixa.tenant_id IS 'Isolamento multi-tenant - OBRIGATÓRIO em todas as queries';
COMMENT ON COLUMN operacoes_caixa.tipo IS 'VENDA (entrada), SANGRIA (retirada), REFORCO (adição), DESPESA (pagamento)';
COMMENT ON COLUMN operacoes_caixa.valor IS 'Valor da operação (sempre positivo)';
COMMENT ON COLUMN operacoes_caixa.descricao IS 'Descrição obrigatória da operação';
COMMENT ON COLUMN operacoes_caixa.destino IS 'Para SANGRIA: DEPOSITO, PAGAMENTO, COFRE, OUTROS';
COMMENT ON COLUMN operacoes_caixa.origem IS 'Para REFORCO: TROCO, CAPITAL_GIRO, TRANSFERENCIA, OUTROS';
COMMENT ON COLUMN operacoes_caixa.usuario_id IS 'Operador que realizou a operação';

-- ============================================================
-- ÍNDICES para operacoes_caixa
-- ============================================================

-- Índice para listar operações de um caixa específico
CREATE INDEX IF NOT EXISTS idx_operacoes_caixa_caixa_id 
    ON operacoes_caixa(caixa_id);

-- Índice para buscar operações por tenant e tipo
CREATE INDEX IF NOT EXISTS idx_operacoes_caixa_tenant_tipo 
    ON operacoes_caixa(tenant_id, tipo);

-- Índice para ordenação cronológica
CREATE INDEX IF NOT EXISTS idx_operacoes_caixa_created_at 
    ON operacoes_caixa(caixa_id, created_at DESC);

-- ============================================================
-- TRIGGER: Atualizar updated_at automaticamente
-- ============================================================

-- Trigger para caixa_diario
CREATE TRIGGER trg_caixa_diario_updated_at
    BEFORE UPDATE ON caixa_diario
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
