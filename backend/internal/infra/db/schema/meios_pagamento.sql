-- Schema: meios_pagamento
-- Formas de pagamento aceitas pela barbearia

CREATE TABLE IF NOT EXISTS meios_pagamento (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome            VARCHAR(100) NOT NULL,
    tipo            VARCHAR(30) NOT NULL,
    bandeira        VARCHAR(50),
    taxa            NUMERIC(5,2) DEFAULT 0.00,
    taxa_fixa       NUMERIC(10,2) DEFAULT 0.00,
    d_mais          INTEGER DEFAULT 0,
    icone           VARCHAR(50),
    cor             VARCHAR(7),
    ordem_exibicao  INTEGER DEFAULT 0,
    observacoes     TEXT,
    ativo           BOOLEAN DEFAULT true,
    criado_em       TIMESTAMPTZ DEFAULT now(),
    atualizado_em   TIMESTAMPTZ DEFAULT now(),
    
    CONSTRAINT chk_taxa_valida CHECK (taxa >= 0 AND taxa <= 100),
    CONSTRAINT chk_taxa_fixa_valida CHECK (taxa_fixa >= 0),
    CONSTRAINT chk_tipo_valido CHECK (tipo IN ('DINHEIRO', 'PIX', 'CREDITO', 'DEBITO', 'TRANSFERENCIA', 'BOLETO', 'OUTRO')),
    CONSTRAINT chk_d_mais_valido CHECK (d_mais >= 0)
);

COMMENT ON TABLE meios_pagamento IS 'Formas de pagamento aceitas - isolamento por tenant_id';
COMMENT ON COLUMN meios_pagamento.bandeira IS 'Bandeira do cartão (Visa, Mastercard, Elo, etc.) - aplicável para CREDITO e DEBITO';
COMMENT ON COLUMN meios_pagamento.taxa IS 'Taxa percentual cobrada (0-100%)';
COMMENT ON COLUMN meios_pagamento.taxa_fixa IS 'Taxa fixa por transação em R$';
COMMENT ON COLUMN meios_pagamento.d_mais IS 'Dias para compensação bancária (D+0 para PIX/Dinheiro, D+1 para Débito, D+30 para Crédito)';
COMMENT ON COLUMN meios_pagamento.icone IS 'Nome do ícone Material Icons (ex: credit_card, pix, payments)';
