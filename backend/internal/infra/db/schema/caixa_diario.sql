-- ============================================================
-- Schema: caixa_diario + operacoes_caixa
-- Para uso com SQLC
-- ============================================================

-- ============================================================
-- TABELA: caixa_diario
-- Controle operacional da gaveta de dinheiro
-- ============================================================

CREATE TABLE IF NOT EXISTS caixa_diario (
    -- Identificadores
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    
    -- Operadores
    usuario_abertura_id UUID NOT NULL,
    usuario_fechamento_id UUID,
    
    -- Timestamps
    data_abertura TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    data_fechamento TIMESTAMPTZ,
    
    -- Valores Monetários (DECIMAL para precisão financeira)
    saldo_inicial DECIMAL(15,2) NOT NULL,
    total_entradas DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_saidas DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_sangrias DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_reforcos DECIMAL(15,2) NOT NULL DEFAULT 0,
    
    -- Saldo esperado (calculado pela aplicação, não GENERATED para compatibilidade sqlc)
    saldo_esperado DECIMAL(15,2) NOT NULL DEFAULT 0,
    
    -- Fechamento
    saldo_real DECIMAL(15,2),
    divergencia DECIMAL(15,2),
    
    -- Status e Justificativa
    status VARCHAR(20) NOT NULL DEFAULT 'ABERTO',
    justificativa_divergencia TEXT,
    
    -- Auditoria
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- TABELA: operacoes_caixa
-- Registro de todas as operações do caixa
-- ============================================================

CREATE TABLE IF NOT EXISTS operacoes_caixa (
    -- Identificadores
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    caixa_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    
    -- Tipo de Operação
    tipo VARCHAR(20) NOT NULL,
    
    -- Valor (sempre positivo)
    valor DECIMAL(15,2) NOT NULL,
    
    -- Detalhes
    descricao TEXT NOT NULL,
    
    -- Para SANGRIA: destino do dinheiro
    destino VARCHAR(50),
    
    -- Para REFORCO: origem do dinheiro
    origem VARCHAR(50),
    
    -- Operador
    usuario_id UUID NOT NULL,
    
    -- Auditoria
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
