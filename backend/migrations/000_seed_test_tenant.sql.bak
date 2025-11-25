-- Seed data para tenant de testes
-- Uso: psql $DATABASE_URL < migrations/seed_test_tenant.sql

-- Criar tenant de teste se não existir
INSERT INTO tenants (id, nome, cnpj, ativo, plano, criado_em, atualizado_em, onboarding_completed, onboarding_step)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Barbearia Teste Smoke',
    '11111111111111',
    true,
    'premium',
    NOW(),
    NOW(),
    true,
    5
)
ON CONFLICT (id) DO UPDATE SET
    atualizado_em = NOW(),
    ativo = true;

-- Criar configuração de precificação padrão para o tenant de teste
INSERT INTO precificacao_config (
    id,
    tenant_id,
    margem_desejada,
    markup_alvo,
    imposto_percentual,
    comissao_percentual_default,
    criado_em,
    atualizado_em
)
VALUES (
    '00000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000001',
    30.00,  -- 30% margem desejada
    1.43,   -- markup alvo
    8.50,   -- 8.5% imposto
    40.00,  -- 40% comissão default
    NOW(),
    NOW()
)
ON CONFLICT (tenant_id) DO UPDATE SET
    margem_desejada = EXCLUDED.margem_desejada,
    markup_alvo = EXCLUDED.markup_alvo,
    imposto_percentual = EXCLUDED.imposto_percentual,
    comissao_percentual_default = EXCLUDED.comissao_percentual_default,
    atualizado_em = NOW();

SELECT 'Seed data criado com sucesso!' AS resultado;
SELECT 'Tenant ID: 00000000-0000-0000-0000-000000000001' AS tenant_info;
SELECT 'Pricing Config criada!' AS pricing_info;
