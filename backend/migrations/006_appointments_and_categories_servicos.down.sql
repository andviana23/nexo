-- +goose Down
-- Rollback da migration 006

DROP TRIGGER IF EXISTS trg_check_appointment_conflict ON appointments;
DROP FUNCTION IF EXISTS check_appointment_conflict();
DROP TRIGGER IF EXISTS trg_categorias_servicos_updated_at ON categorias_servicos;
DROP TRIGGER IF EXISTS trg_appointments_updated_at ON appointments;
DROP FUNCTION IF EXISTS update_appointments_updated_at();
DROP TABLE IF EXISTS appointment_services;
DROP TABLE IF EXISTS appointments;
ALTER TABLE servicos DROP COLUMN IF EXISTS categoria_servico_id;
DROP TABLE IF EXISTS categorias_servicos;
