-- +goose Up
-- Seed de usuário de teste para desenvolvimento
-- Email: admin@teste.com
-- Senha: Admin123!
-- Role: owner
-- Tenant ID: 00000000-0000-0000-0000-000000000001

-- Usuário admin de teste
-- Senha: Admin123! (hash bcrypt com custo 10)
INSERT INTO users (id, tenant_id, nome, email, password_hash, role, ativo)
VALUES (
    '10000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    'Administrador Teste',
    'admin@teste.com',
    '$2a$10$v1RGWSKEVBgDKNKQuNfCpeanXjBAxHBgyU3PlGZeap2dyqIvWbgFO',  -- Admin123!
    'owner',
    true
)
ON CONFLICT (tenant_id, email) DO NOTHING;

-- +goose Down
DELETE FROM users WHERE email = 'admin@teste.com';
