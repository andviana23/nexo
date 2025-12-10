-- Reverte as mudanças da migration UP
DROP TRIGGER IF EXISTS trg_user_units_atualizado_em ON user_units;
DROP FUNCTION IF EXISTS update_atualizado_em_column();

-- Recria o trigger incorreto (apenas para rollback exato, se necessário, embora seja "bugado")
-- Idealmente, não queremos voltar para o estado bugado, mas para rollback fiel:
-- CREATE TRIGGER trg_user_units_updated_at
--    BEFORE UPDATE ON user_units
--    FOR EACH ROW
--    EXECUTE FUNCTION update_updated_at_column();
