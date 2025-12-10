-- 1. Criar função genérica para atualizar 'atualizado_em'
CREATE OR REPLACE FUNCTION update_atualizado_em_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.atualizado_em = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 2. Remover o trigger incorreto (se existir)
DROP TRIGGER IF EXISTS trg_user_units_updated_at ON user_units;

-- 3. Criar o trigger correto
DROP TRIGGER IF EXISTS trg_user_units_atualizado_em ON user_units;

CREATE TRIGGER trg_user_units_atualizado_em
    BEFORE UPDATE ON user_units
    FOR EACH ROW
    EXECUTE FUNCTION update_atualizado_em_column();
