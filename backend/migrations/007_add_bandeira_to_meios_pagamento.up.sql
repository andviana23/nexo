-- Migration: 007_add_bandeira_to_meios_pagamento
-- Adiciona coluna bandeira à tabela meios_pagamento para suportar bandeiras de cartão

-- Adiciona coluna bandeira
ALTER TABLE meios_pagamento 
ADD COLUMN IF NOT EXISTS bandeira VARCHAR(50);

-- Adiciona comentário
COMMENT ON COLUMN meios_pagamento.bandeira IS 'Bandeira do cartão (Visa, Mastercard, Elo, etc.) - aplicável para CREDITO e DEBITO';

-- Atualiza constraint de tipo para incluir BOLETO e OUTRO
ALTER TABLE meios_pagamento 
DROP CONSTRAINT IF EXISTS chk_tipo_valido;

ALTER TABLE meios_pagamento 
ADD CONSTRAINT chk_tipo_valido 
CHECK (tipo IN ('DINHEIRO', 'PIX', 'CREDITO', 'DEBITO', 'TRANSFERENCIA', 'BOLETO', 'OUTRO'));
