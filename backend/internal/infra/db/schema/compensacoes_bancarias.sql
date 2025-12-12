-- Compensações bancárias com D+ para fluxo de caixa compensado
-- Tabela: compensacoes_bancarias

CREATE TABLE IF NOT EXISTS compensacoes_bancarias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    receita_id UUID,

    data_transacao DATE NOT NULL,
    data_compensacao DATE NOT NULL,
    data_compensado DATE,

    valor_bruto NUMERIC(15,2) NOT NULL,
    taxa_percentual NUMERIC(5,2) DEFAULT 0,
    taxa_fixa NUMERIC(10,2) DEFAULT 0,
    valor_liquido NUMERIC(15,2) NOT NULL,

    meio_pagamento_id UUID,
    d_mais INTEGER NOT NULL,

    status VARCHAR(20) DEFAULT 'PREVISTO',

    criado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT now(),

    CONSTRAINT compensacoes_bancarias_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    -- receita_id referencia a receita canônica (contas_a_receber) nos fluxos V2
    CONSTRAINT compensacoes_bancarias_meio_pagamento_id_fkey FOREIGN KEY (meio_pagamento_id) REFERENCES meios_pagamento(id),
    CONSTRAINT chk_status_compensacao CHECK (status IN ('PREVISTO', 'CONFIRMADO', 'COMPENSADO', 'CANCELADO'))
);

CREATE INDEX IF NOT EXISTS idx_compensacoes_tenant ON compensacoes_bancarias(tenant_id);
CREATE INDEX IF NOT EXISTS idx_compensacoes_data_compensacao ON compensacoes_bancarias(tenant_id, data_compensacao);
CREATE INDEX IF NOT EXISTS idx_compensacoes_status ON compensacoes_bancarias(status, data_compensacao);
CREATE INDEX IF NOT EXISTS idx_compensacoes_receita ON compensacoes_bancarias(receita_id);

COMMENT ON TABLE compensacoes_bancarias IS 'Compensações bancárias com D+ para fluxo de caixa compensado';
