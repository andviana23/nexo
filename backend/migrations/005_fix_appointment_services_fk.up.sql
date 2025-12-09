-- Migration: Corrigir FK constraint para permitir deleção de serviços
-- Problema: appointment_services_service_id_fkey usa RESTRICT, impedindo deleção de serviços vinculados
-- Nota: SET NULL não funciona porque service_id é parte da PK composta (appointment_id, service_id)
-- Solução: Alterar para ON DELETE CASCADE (remove registros de appointment_services ao deletar serviço)

-- Remover constraint antiga
ALTER TABLE appointment_services
DROP CONSTRAINT IF EXISTS appointment_services_service_id_fkey;

-- Recriar constraint com ON DELETE CASCADE
ALTER TABLE appointment_services
ADD CONSTRAINT appointment_services_service_id_fkey 
FOREIGN KEY (service_id) 
REFERENCES servicos(id) 
ON DELETE CASCADE;

-- Comentário explicativo
COMMENT ON CONSTRAINT appointment_services_service_id_fkey ON appointment_services IS 
'FK para servicos com CASCADE - ao deletar serviço, remove registros de appointment_services vinculados';
