ALTER TABLE meios_pagamento ALTER COLUMN nome DROP NOT NULL;

DROP INDEX IF EXISTS idx_meios_pagamento_tenant_nome_tipo;

CREATE UNIQUE INDEX idx_meios_pagamento_tenant_tipo_bandeira 
ON meios_pagamento (tenant_id, tipo, COALESCE(bandeira, '')) 
WHERE ativo = true;
