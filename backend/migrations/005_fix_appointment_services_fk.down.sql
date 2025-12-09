-- Rollback: Restaurar FK constraint original com RESTRICT

-- Remover constraint alterada
ALTER TABLE appointment_services
DROP CONSTRAINT IF EXISTS appointment_services_service_id_fkey;

-- Recriar constraint original com RESTRICT (comportamento padrão)
ALTER TABLE appointment_services
ADD CONSTRAINT appointment_services_service_id_fkey 
FOREIGN KEY (service_id) 
REFERENCES servicos(id);
