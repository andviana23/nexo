-- ============================================================================
-- SQLC Schema: User Units (Vínculo Usuário-Unidade)
-- ============================================================================

CREATE TABLE IF NOT EXISTS user_units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    unit_id UUID NOT NULL,
    
    -- Configurações do vínculo
    is_default BOOLEAN DEFAULT false NOT NULL,
    role_override VARCHAR(50),
    
    -- Timestamps
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- Constraints
    CONSTRAINT user_units_unique UNIQUE (user_id, unit_id)
);
