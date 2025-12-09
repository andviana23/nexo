-- Migration: 007_add_bandeira_to_meios_pagamento (DOWN)
-- Remove coluna bandeira da tabela meios_pagamento

-- Remove coluna bandeira
ALTER TABLE meios_pagamento 
DROP COLUMN IF EXISTS bandeira;

-- Reverte constraint de tipo para versão anterior
ALTER TABLE meios_pagamento 
DROP CONSTRAINT IF EXISTS chk_tipo_valido;

ALTER TABLE meios_pagamento 
ADD CONSTRAINT chk_tipo_valido 
CHECK (tipo IN ('DINHEIRO', 'PIX', 'CREDITO', 'DEBITO', 'TRANSFERENCIA'));
