-- +goose Up
-- Migration consolidada para criar todas as tabelas restantes do sistema
-- Baseado no schema atual do banco de dados e arquivos em internal/infra/db/schema

-- 1. Tenants (Fundamental)
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome VARCHAR(255) NOT NULL UNIQUE,
    cnpj VARCHAR(14) UNIQUE,
    ativo BOOLEAN DEFAULT true,
    plano VARCHAR(50) DEFAULT 'free',
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    onboarding_completed BOOLEAN DEFAULT false,
    onboarding_step INTEGER DEFAULT 1
);

-- 2. Categorias
CREATE TABLE IF NOT EXISTS categorias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    tipo VARCHAR(20) NOT NULL CHECK (tipo IN ('RECEITA', 'DESPESA')),
    cor VARCHAR(7) DEFAULT '#000000',
    ativa BOOLEAN DEFAULT true,
    criado_em TIMESTAMPTZ DEFAULT now(),
    tipo_custo VARCHAR(20) DEFAULT 'FIXO' CHECK (tipo_custo IN ('FIXO', 'VARIAVEL')),
    CONSTRAINT categorias_tenant_id_nome_key UNIQUE (tenant_id, nome)
);

-- 3. Fornecedores
CREATE TABLE IF NOT EXISTS fornecedores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    razao_social VARCHAR(255) NOT NULL,
    nome_fantasia VARCHAR(255),
    cnpj VARCHAR(14) NOT NULL,
    email VARCHAR(255),
    telefone VARCHAR(20),
    celular VARCHAR(20),
    endereco_logradouro VARCHAR(255),
    endereco_numero VARCHAR(20),
    endereco_complemento VARCHAR(100),
    endereco_bairro VARCHAR(100),
    endereco_cidade VARCHAR(100),
    endereco_estado VARCHAR(2),
    endereco_cep VARCHAR(8),
    banco VARCHAR(100),
    agencia VARCHAR(20),
    conta VARCHAR(30),
    observacoes TEXT,
    ativo BOOLEAN DEFAULT true NOT NULL,
    criado_em TIMESTAMPTZ DEFAULT now() NOT NULL,
    atualizado_em TIMESTAMPTZ DEFAULT now() NOT NULL,
    CONSTRAINT chk_fornecedor_cnpj_valido CHECK (length(cnpj) = 14),
    CONSTRAINT fornecedores_cnpj_tenant_unique UNIQUE (tenant_id, cnpj)
);

-- 4. Produtos
CREATE TABLE IF NOT EXISTS produtos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    categoria_id UUID REFERENCES categorias(id) ON DELETE SET NULL,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    sku VARCHAR(100),
    codigo_barras VARCHAR(50),
    preco NUMERIC(10,2) NOT NULL CHECK (preco > 0),
    custo NUMERIC(10,2) CHECK (custo IS NULL OR custo >= 0),
    estoque INTEGER DEFAULT 0 CHECK (estoque >= 0),
    estoque_minimo INTEGER DEFAULT 0 CHECK (estoque_minimo >= 0),
    unidade VARCHAR(10) DEFAULT 'UN',
    fornecedor VARCHAR(255),
    imagem TEXT,
    observacoes TEXT,
    ativo BOOLEAN DEFAULT true,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    categoria_produto VARCHAR(30) DEFAULT 'REVENDA' NOT NULL CHECK (categoria_produto IN ('INSUMO', 'REVENDA', 'USO_INTERNO', 'PERMANENTE', 'PROMOCIONAL', 'KIT', 'SERVICO')),
    unidade_medida VARCHAR(20) DEFAULT 'UNIDADE' NOT NULL CHECK (unidade_medida IN ('UNIDADE', 'LITRO', 'MILILITRO', 'GRAMA', 'QUILOGRAMA')),
    quantidade_atual NUMERIC(15,3) DEFAULT 0 NOT NULL CHECK (quantidade_atual >= 0),
    quantidade_minima NUMERIC(15,3) DEFAULT 0 NOT NULL CHECK (quantidade_minima >= 0),
    localizacao VARCHAR(100),
    lote VARCHAR(50),
    data_validade DATE,
    ncm VARCHAR(8),
    permite_venda BOOLEAN DEFAULT true NOT NULL,
    CONSTRAINT idx_produtos_tenant_nome UNIQUE (tenant_id, nome),
    CONSTRAINT idx_produtos_tenant_sku UNIQUE (tenant_id, sku),
    CONSTRAINT idx_produtos_tenant_codigo_barras UNIQUE (tenant_id, codigo_barras)
);

-- 5. Serviços
CREATE TABLE IF NOT EXISTS servicos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    categoria_id UUID REFERENCES categorias(id) ON DELETE SET NULL,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    preco NUMERIC(10,2) NOT NULL CHECK (preco > 0),
    duracao INTEGER NOT NULL CHECK (duracao >= 5),
    comissao NUMERIC(5,2) DEFAULT 0.00 CHECK (comissao >= 0 AND comissao <= 100),
    cor VARCHAR(7),
    imagem TEXT,
    profissionais_ids UUID[],
    observacoes TEXT,
    tags TEXT[],
    ativo BOOLEAN DEFAULT true,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT idx_servicos_tenant_nome UNIQUE (tenant_id, nome)
);

-- 6. Profissionais
CREATE TABLE IF NOT EXISTS profissionais (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telefone VARCHAR(20) NOT NULL,
    cpf VARCHAR(11) NOT NULL,
    especialidades TEXT[],
    comissao NUMERIC(5,2) DEFAULT 0.00,
    tipo_comissao VARCHAR(20) DEFAULT 'PERCENTUAL',
    foto TEXT,
    data_admissao DATE DEFAULT CURRENT_DATE NOT NULL,
    data_demissao DATE,
    status VARCHAR(20) DEFAULT 'ATIVO' CHECK (status IN ('ATIVO', 'INATIVO', 'AFASTADO', 'DEMITIDO')),
    horario_trabalho JSONB,
    observacoes TEXT,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    tipo VARCHAR(30) DEFAULT 'BARBEIRO' NOT NULL CHECK (tipo IN ('BARBEIRO', 'MANICURE', 'RECEPCIONISTA', 'GERENTE', 'OUTRO')),
    CONSTRAINT chk_comissao_valida CHECK ((tipo_comissao = 'PERCENTUAL' AND comissao >= 0 AND comissao <= 100) OR (tipo_comissao = 'FIXO' AND comissao >= 0)),
    CONSTRAINT idx_profissionais_tenant_cpf UNIQUE (tenant_id, cpf),
    CONSTRAINT idx_profissionais_tenant_email UNIQUE (tenant_id, email),
    CONSTRAINT idx_profissionais_tenant_user UNIQUE (tenant_id, user_id)
);

-- 7. Clientes
CREATE TABLE IF NOT EXISTS clientes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    telefone VARCHAR(20) NOT NULL,
    cpf VARCHAR(11),
    data_nascimento DATE,
    genero VARCHAR(50),
    endereco_logradouro VARCHAR(255),
    endereco_numero VARCHAR(20),
    endereco_complemento VARCHAR(100),
    endereco_bairro VARCHAR(100),
    endereco_cidade VARCHAR(100),
    endereco_estado VARCHAR(2),
    endereco_cep VARCHAR(8),
    observacoes TEXT,
    tags TEXT[],
    ativo BOOLEAN DEFAULT true,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT idx_clientes_tenant_email_unique UNIQUE (tenant_id, email),
    CONSTRAINT idx_clientes_tenant_cpf_unique UNIQUE (tenant_id, cpf)
);

-- 8. Meios de Pagamento
CREATE TABLE IF NOT EXISTS meios_pagamento (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    tipo VARCHAR(30) NOT NULL CHECK (tipo IN ('DINHEIRO', 'PIX', 'CREDITO', 'DEBITO', 'TRANSFERENCIA')),
    taxa NUMERIC(5,2) DEFAULT 0.00 CHECK (taxa >= 0 AND taxa <= 100),
    taxa_fixa NUMERIC(10,2) DEFAULT 0.00 CHECK (taxa_fixa >= 0),
    icone VARCHAR(50),
    cor VARCHAR(7),
    ordem_exibicao INTEGER DEFAULT 0,
    observacoes TEXT,
    ativo BOOLEAN DEFAULT true,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    d_mais INTEGER DEFAULT 0 CHECK (d_mais >= 0),
    CONSTRAINT idx_meios_pagamento_tenant_nome_tipo UNIQUE (tenant_id, nome, tipo)
);

-- 9. Planos de Assinatura
CREATE TABLE IF NOT EXISTS planos_assinatura (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    valor NUMERIC(10,2) NOT NULL CHECK (valor > 0),
    periodicidade VARCHAR(50) NOT NULL,
    quantidade_servicos INTEGER DEFAULT 0,
    ativa BOOLEAN DEFAULT true,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    beneficios_json TEXT,
    limite_servicos INTEGER,
    desconto NUMERIC(5,2) DEFAULT 0.00 CHECK (desconto >= 0 AND desconto <= 100),
    cor VARCHAR(7),
    popular BOOLEAN DEFAULT false,
    ordem_exibicao INTEGER DEFAULT 0,
    imagem TEXT,
    CONSTRAINT planos_assinatura_tenant_id_nome_key UNIQUE (tenant_id, nome)
);

-- 10. Assinaturas
CREATE TABLE IF NOT EXISTS assinaturas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES planos_assinatura(id) ON DELETE RESTRICT,
    barbeiro_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    asaas_subscription_id VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) DEFAULT 'ATIVA' NOT NULL,
    data_inicio DATE NOT NULL,
    data_fim DATE,
    proxima_fatura_data DATE NOT NULL,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    data_proximo_pagamento DATE,
    origem_dado VARCHAR(100) DEFAULT 'asaas' NOT NULL
);

-- 11. Assinatura Invoices
CREATE TABLE IF NOT EXISTS assinatura_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    assinatura_id UUID NOT NULL REFERENCES assinaturas(id) ON DELETE CASCADE,
    asaas_invoice_id VARCHAR(255) NOT NULL UNIQUE,
    valor NUMERIC(10,2) NOT NULL CHECK (valor > 0),
    status VARCHAR(50) DEFAULT 'PENDENTE' NOT NULL,
    data_vencimento DATE NOT NULL,
    data_pagamento DATE,
    processada BOOLEAN DEFAULT false,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now()
);

-- 12. Receitas
CREATE TABLE IF NOT EXISTS receitas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    usuario_id UUID REFERENCES users(id) ON DELETE SET NULL,
    descricao VARCHAR(255) NOT NULL,
    valor NUMERIC(15,2) NOT NULL CHECK (valor > 0),
    categoria_id UUID NOT NULL REFERENCES categorias(id) ON DELETE RESTRICT,
    metodo_pagamento VARCHAR(50) NOT NULL,
    data DATE DEFAULT CURRENT_DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'CONFIRMADO',
    observacoes TEXT,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    manual BOOLEAN DEFAULT false NOT NULL,
    origem_dado VARCHAR(100) DEFAULT 'imported' NOT NULL,
    subtipo VARCHAR(30) DEFAULT 'SERVICO' CHECK (subtipo IN ('SERVICO', 'PRODUTO', 'PLANO'))
);

-- 13. Despesas
CREATE TABLE IF NOT EXISTS despesas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    usuario_id UUID REFERENCES users(id) ON DELETE SET NULL,
    descricao VARCHAR(255) NOT NULL,
    valor NUMERIC(15,2) NOT NULL CHECK (valor > 0),
    categoria_id UUID NOT NULL REFERENCES categorias(id) ON DELETE RESTRICT,
    fornecedor VARCHAR(255),
    metodo_pagamento VARCHAR(50) NOT NULL,
    data DATE DEFAULT CURRENT_DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDENTE',
    observacoes TEXT,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    manual BOOLEAN DEFAULT false NOT NULL,
    origem_dado VARCHAR(100) DEFAULT 'imported' NOT NULL
);

-- 14. Contas a Pagar
CREATE TABLE IF NOT EXISTS contas_a_pagar (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    descricao VARCHAR(255) NOT NULL,
    categoria_id UUID REFERENCES categorias(id) ON DELETE SET NULL,
    fornecedor VARCHAR(255),
    valor NUMERIC(15,2) NOT NULL CHECK (valor > 0),
    tipo VARCHAR(20) DEFAULT 'VARIAVEL' CHECK (tipo IN ('FIXA', 'VARIAVEL')),
    recorrente BOOLEAN DEFAULT false,
    periodicidade VARCHAR(20) CHECK (periodicidade IN ('MENSAL', 'TRIMESTRAL', 'ANUAL')),
    data_vencimento DATE NOT NULL,
    data_pagamento DATE,
    status VARCHAR(20) DEFAULT 'ABERTO' CHECK (status IN ('ABERTO', 'PAGO', 'ATRASADO', 'CANCELADO')),
    comprovante_url TEXT,
    pix_code TEXT,
    observacoes TEXT,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now()
);

-- 15. Contas a Receber
CREATE TABLE IF NOT EXISTS contas_a_receber (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    origem VARCHAR(30) CHECK (origem IN ('ASSINATURA', 'SERVICO', 'OUTRO')),
    assinatura_id UUID REFERENCES assinaturas(id) ON DELETE CASCADE,
    servico_id UUID REFERENCES servicos(id) ON DELETE SET NULL,
    descricao VARCHAR(255) NOT NULL,
    valor NUMERIC(15,2) NOT NULL CHECK (valor > 0),
    valor_pago NUMERIC(15,2) DEFAULT 0 CHECK (valor_pago >= 0),
    data_vencimento DATE NOT NULL,
    data_recebimento DATE,
    status VARCHAR(20) DEFAULT 'PENDENTE' CHECK (status IN ('PENDENTE', 'RECEBIDO', 'ATRASADO', 'CANCELADO')),
    observacoes TEXT,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now()
);

-- 16. Compensações Bancárias
CREATE TABLE IF NOT EXISTS compensacoes_bancarias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    receita_id UUID REFERENCES receitas(id) ON DELETE CASCADE,
    data_transacao DATE NOT NULL,
    data_compensacao DATE NOT NULL,
    data_compensado DATE,
    valor_bruto NUMERIC(15,2) NOT NULL,
    taxa_percentual NUMERIC(5,2) DEFAULT 0,
    taxa_fixa NUMERIC(10,2) DEFAULT 0,
    valor_liquido NUMERIC(15,2) NOT NULL,
    meio_pagamento_id UUID REFERENCES meios_pagamento(id),
    d_mais INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'PREVISTO' CHECK (status IN ('PREVISTO', 'CONFIRMADO', 'COMPENSADO', 'CANCELADO')),
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now()
);

-- 17. Fluxo de Caixa Diário
CREATE TABLE IF NOT EXISTS fluxo_caixa_diario (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    data DATE NOT NULL,
    saldo_inicial NUMERIC(15,2) DEFAULT 0,
    saldo_final NUMERIC(15,2) DEFAULT 0,
    entradas_confirmadas NUMERIC(15,2) DEFAULT 0,
    entradas_previstas NUMERIC(15,2) DEFAULT 0,
    saidas_pagas NUMERIC(15,2) DEFAULT 0,
    saidas_previstas NUMERIC(15,2) DEFAULT 0,
    processado_em TIMESTAMPTZ DEFAULT now(),
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT fluxo_caixa_diario_tenant_id_data_key UNIQUE (tenant_id, data)
);

-- 18. DRE Mensal
CREATE TABLE IF NOT EXISTS dre_mensal (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    mes_ano VARCHAR(7) NOT NULL,
    receita_servicos NUMERIC(15,2) DEFAULT 0,
    receita_produtos NUMERIC(15,2) DEFAULT 0,
    receita_planos NUMERIC(15,2) DEFAULT 0,
    receita_total NUMERIC(15,2) DEFAULT 0,
    custo_comissoes NUMERIC(15,2) DEFAULT 0,
    custo_insumos NUMERIC(15,2) DEFAULT 0,
    custo_variavel_total NUMERIC(15,2) DEFAULT 0,
    despesa_fixa NUMERIC(15,2) DEFAULT 0,
    despesa_variavel NUMERIC(15,2) DEFAULT 0,
    despesa_total NUMERIC(15,2) DEFAULT 0,
    resultado_bruto NUMERIC(15,2) DEFAULT 0,
    resultado_operacional NUMERIC(15,2) DEFAULT 0,
    margem_bruta NUMERIC(5,2) DEFAULT 0,
    margem_operacional NUMERIC(5,2) DEFAULT 0,
    lucro_liquido NUMERIC(15,2) DEFAULT 0,
    processado_em TIMESTAMPTZ DEFAULT now(),
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT dre_mensal_tenant_id_mes_ano_key UNIQUE (tenant_id, mes_ano)
);

-- 19. Metas Mensais
CREATE TABLE IF NOT EXISTS metas_mensais (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    mes_ano VARCHAR(7) NOT NULL,
    meta_faturamento NUMERIC(15,2) NOT NULL CHECK (meta_faturamento > 0),
    origem VARCHAR(20) DEFAULT 'MANUAL' CHECK (origem IN ('MANUAL', 'AUTOMATICA')),
    status VARCHAR(20) DEFAULT 'PENDENTE' CHECK (status IN ('PENDENTE', 'ACEITA', 'REJEITADA')),
    criado_por UUID REFERENCES users(id) ON DELETE SET NULL,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT metas_mensais_tenant_id_mes_ano_key UNIQUE (tenant_id, mes_ano)
);

-- 20. Metas Barbeiro
CREATE TABLE IF NOT EXISTS metas_barbeiro (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    barbeiro_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mes_ano VARCHAR(7) NOT NULL,
    meta_servicos_gerais NUMERIC(15,2) DEFAULT 0,
    meta_servicos_extras NUMERIC(15,2) DEFAULT 0,
    meta_produtos NUMERIC(15,2) DEFAULT 0,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT metas_barbeiro_tenant_id_barbeiro_id_mes_ano_key UNIQUE (tenant_id, barbeiro_id, mes_ano)
);

-- 21. Metas Ticket Médio
CREATE TABLE IF NOT EXISTS metas_ticket_medio (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    mes_ano VARCHAR(7) NOT NULL,
    meta_valor NUMERIC(10,2) NOT NULL CHECK (meta_valor > 0),
    tipo VARCHAR(20) DEFAULT 'GERAL' CHECK (tipo IN ('GERAL', 'BARBEIRO')),
    barbeiro_id UUID REFERENCES users(id) ON DELETE CASCADE,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT chk_barbeiro_tipo CHECK ((tipo = 'GERAL' AND barbeiro_id IS NULL) OR (tipo = 'BARBEIRO' AND barbeiro_id IS NOT NULL))
);

-- 22. Precificação Config
CREATE TABLE IF NOT EXISTS precificacao_config (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    margem_desejada NUMERIC(5,2) DEFAULT 30.00 CHECK (margem_desejada >= 5 AND margem_desejada <= 100),
    markup_alvo NUMERIC(5,2) DEFAULT 1.5 CHECK (markup_alvo >= 1),
    imposto_percentual NUMERIC(5,2) DEFAULT 0.00 CHECK (imposto_percentual >= 0 AND imposto_percentual <= 100),
    comissao_percentual_default NUMERIC(5,2) DEFAULT 30.00 CHECK (comissao_percentual_default >= 0 AND comissao_percentual_default <= 100),
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT precificacao_config_tenant_id_key UNIQUE (tenant_id)
);

-- 23. Precificação Simulações
CREATE TABLE IF NOT EXISTS precificacao_simulacoes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    item_id UUID,
    tipo_item VARCHAR(20) CHECK (tipo_item IN ('SERVICO', 'PRODUTO')),
    custo_insumos NUMERIC(15,2) NOT NULL,
    comissao_percentual NUMERIC(5,2) NOT NULL,
    imposto_percentual NUMERIC(5,2) NOT NULL,
    margem_desejada NUMERIC(5,2) NOT NULL,
    markup_aplicado NUMERIC(5,2) NOT NULL,
    preco_sugerido NUMERIC(15,2) NOT NULL,
    margem_resultante NUMERIC(5,2) NOT NULL,
    parametros_json JSONB,
    criado_por UUID REFERENCES users(id) ON DELETE SET NULL,
    criado_em TIMESTAMPTZ DEFAULT now(),
    custo_materiais NUMERIC(15,2) DEFAULT 0.00,
    custo_mao_de_obra NUMERIC(15,2) DEFAULT 0.00,
    custo_total NUMERIC(15,2) DEFAULT 0.00,
    preco_atual NUMERIC(15,2) DEFAULT 0.00,
    diferenca_percentual NUMERIC(5,2) DEFAULT 0.00,
    lucro_estimado NUMERIC(15,2) DEFAULT 0.00,
    margem_final NUMERIC(5,2) DEFAULT 0.00
);

-- 24. Movimentações Estoque
CREATE TABLE IF NOT EXISTS movimentacoes_estoque (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    produto_id UUID NOT NULL REFERENCES produtos(id) ON DELETE RESTRICT,
    tipo_movimentacao VARCHAR(30) NOT NULL CHECK (tipo_movimentacao IN ('ENTRADA', 'SAIDA', 'CONSUMO_INTERNO', 'AJUSTE', 'DEVOLUCAO', 'PERDA')),
    quantidade NUMERIC(15,3) NOT NULL CHECK (quantidade > 0),
    valor_unitario NUMERIC(15,2) DEFAULT 0 NOT NULL CHECK (valor_unitario >= 0),
    valor_total NUMERIC(15,2) DEFAULT 0 NOT NULL CHECK (valor_total >= 0),
    fornecedor_id UUID REFERENCES fornecedores(id) ON DELETE SET NULL,
    usuario_id UUID REFERENCES users(id) ON DELETE SET NULL,
    data_movimentacao TIMESTAMPTZ DEFAULT now() NOT NULL,
    observacoes TEXT,
    documento VARCHAR(100),
    criado_em TIMESTAMPTZ DEFAULT now() NOT NULL,
    atualizado_em TIMESTAMPTZ DEFAULT now() NOT NULL,
    CONSTRAINT chk_observacoes_obrigatorias_ajuste_perda CHECK ((tipo_movimentacao NOT IN ('AJUSTE', 'PERDA')) OR (observacoes IS NOT NULL AND length(TRIM(observacoes)) > 0))
);

-- 25. Produto Fornecedor
CREATE TABLE IF NOT EXISTS produto_fornecedor (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    produto_id UUID NOT NULL REFERENCES produtos(id) ON DELETE CASCADE,
    fornecedor_id UUID NOT NULL REFERENCES fornecedores(id) ON DELETE CASCADE,
    codigo_fornecedor VARCHAR(100),
    preco_compra NUMERIC(15,2) CHECK (preco_compra IS NULL OR preco_compra >= 0),
    prazo_entrega_dias INTEGER DEFAULT 0 CHECK (prazo_entrega_dias >= 0),
    fornecedor_preferencial BOOLEAN DEFAULT false,
    ativo BOOLEAN DEFAULT true NOT NULL,
    criado_em TIMESTAMPTZ DEFAULT now() NOT NULL,
    atualizado_em TIMESTAMPTZ DEFAULT now() NOT NULL,
    CONSTRAINT produto_fornecedor_unique UNIQUE (tenant_id, produto_id, fornecedor_id)
);

-- 26. Cupons Desconto
CREATE TABLE IF NOT EXISTS cupons_desconto (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    codigo VARCHAR(50) NOT NULL,
    descricao TEXT,
    tipo_desconto VARCHAR(20) NOT NULL CHECK (tipo_desconto IN ('PERCENTUAL', 'FIXO')),
    valor NUMERIC(10,2) NOT NULL CHECK (valor > 0),
    valor_minimo NUMERIC(10,2) CHECK (valor_minimo IS NULL OR valor_minimo >= 0),
    data_inicio TIMESTAMPTZ DEFAULT now() NOT NULL,
    data_fim TIMESTAMPTZ NOT NULL CHECK (data_fim > data_inicio),
    limite_uso INTEGER CHECK (limite_uso IS NULL OR limite_uso > 0),
    usos_realizados INTEGER DEFAULT 0 CHECK (usos_realizados >= 0),
    limite_por_cliente INTEGER CHECK (limite_por_cliente IS NULL OR limite_por_cliente > 0),
    apenas_novo_cliente BOOLEAN DEFAULT false,
    servico_ids UUID[],
    ativo BOOLEAN DEFAULT true,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT idx_cupons_tenant_codigo UNIQUE (tenant_id, codigo)
);

-- 27. Barbers Turn List
CREATE TABLE IF NOT EXISTS barbers_turn_list (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    current_points INTEGER DEFAULT 0 NOT NULL CHECK (current_points >= 0),
    last_turn_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    CONSTRAINT uq_barber_turn_professional_tenant UNIQUE (professional_id, tenant_id)
);

-- 28. Barber Turn History
CREATE TABLE IF NOT EXISTS barber_turn_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE CASCADE,
    month_year VARCHAR(7) NOT NULL,
    total_turns INTEGER DEFAULT 0 NOT NULL,
    final_points INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    CONSTRAINT uq_history_professional_month UNIQUE (professional_id, tenant_id, month_year)
);

-- 29. Barber Commissions
CREATE TABLE IF NOT EXISTS barber_commissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    barbeiro_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    receita_id UUID REFERENCES receitas(id) ON DELETE SET NULL,
    valor NUMERIC(18,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDENTE' NOT NULL,
    data_competencia DATE NOT NULL,
    manual BOOLEAN DEFAULT false NOT NULL,
    criado_em TIMESTAMPTZ DEFAULT now() NOT NULL,
    atualizado_em TIMESTAMPTZ DEFAULT now() NOT NULL
);

-- 30. Audit Logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    resource_name VARCHAR(100) NOT NULL,
    resource_id VARCHAR(255),
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    timestamp TIMESTAMPTZ DEFAULT now(),
    resource_type VARCHAR(100),
    user_agent TEXT,
    deleted_at TIMESTAMPTZ
);

-- 31. Cron Run Logs
CREATE TABLE IF NOT EXISTS cron_run_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    job_name VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    started_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    finished_at TIMESTAMPTZ,
    details JSONB,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

-- 32. Feature Flags
CREATE TABLE IF NOT EXISTS feature_flags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    feature VARCHAR(120) NOT NULL CHECK (feature <> ''),
    enabled BOOLEAN DEFAULT false NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

-- 33. Financial Snapshots
CREATE TABLE IF NOT EXISTS financial_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    periodo_inicio DATE NOT NULL,
    periodo_fim DATE NOT NULL,
    entradas NUMERIC(18,2) DEFAULT 0 NOT NULL,
    saidas NUMERIC(18,2) DEFAULT 0 NOT NULL,
    saldo NUMERIC(18,2) DEFAULT 0 NOT NULL,
    origem_dado VARCHAR(100) DEFAULT 'cron-snapshot' NOT NULL,
    criado_em TIMESTAMPTZ DEFAULT now() NOT NULL,
    atualizado_em TIMESTAMPTZ DEFAULT now() NOT NULL
);

-- 34. Tenant Settings
CREATE TABLE IF NOT EXISTS tenant_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    business_hours JSONB DEFAULT '{"days_open": ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday"], "closing_time": "18:00", "opening_time": "08:00"}' NOT NULL,
    financial_settings JSONB DEFAULT '{"default_commission_rate": 30, "accepted_payment_methods": ["PIX", "DINHEIRO", "DEBITO", "CREDITO"]}' NOT NULL,
    preferences JSONB DEFAULT '{"timezone": "America/Sao_Paulo", "default_service_duration": 30}' NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    CONSTRAINT tenant_settings_tenant_id_key UNIQUE (tenant_id)
);

-- 35. User Preferences
CREATE TABLE IF NOT EXISTS user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    analytics_enabled BOOLEAN DEFAULT false,
    error_tracking_enabled BOOLEAN DEFAULT false,
    marketing_enabled BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    CONSTRAINT user_preferences_user_id_key UNIQUE (user_id)
);

-- Funções auxiliares
CREATE OR REPLACE FUNCTION check_professional_is_barber()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM profissionais
        WHERE id = NEW.professional_id
        AND tenant_id = NEW.tenant_id
        AND tipo = 'BARBEIRO'
    ) THEN
        RAISE EXCEPTION 'Apenas profissionais do tipo BARBEIRO podem ser adicionados à lista da vez';
    END IF;
    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION update_updated_at_column()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    -- Compatível com tabelas que usam updated_at e/ou atualizado_em
    IF to_jsonb(NEW) ? 'updated_at' THEN
        NEW.updated_at = NOW();
    END IF;
    IF to_jsonb(NEW) ? 'atualizado_em' THEN
        NEW.atualizado_em = NOW();
    END IF;
    RETURN NEW;
END;
$function$;

-- Triggers
DROP TRIGGER IF EXISTS trg_validate_barber_type ON barbers_turn_list;
CREATE TRIGGER trg_validate_barber_type BEFORE INSERT OR UPDATE ON barbers_turn_list FOR EACH ROW EXECUTE FUNCTION check_professional_is_barber();

DROP TRIGGER IF EXISTS update_user_preferences_updated_at ON user_preferences;
CREATE TRIGGER update_user_preferences_updated_at BEFORE UPDATE ON user_preferences FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
-- Remoção das tabelas (ordem inversa)
DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS tenant_settings;
DROP TABLE IF EXISTS financial_snapshots;
DROP TABLE IF EXISTS feature_flags;
DROP TABLE IF EXISTS cron_run_logs;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS barber_commissions;
DROP TABLE IF EXISTS barber_turn_history;
DROP TABLE IF EXISTS barbers_turn_list;
DROP TABLE IF EXISTS cupons_desconto;
DROP TABLE IF EXISTS produto_fornecedor;
DROP TABLE IF EXISTS movimentacoes_estoque;
DROP TABLE IF EXISTS precificacao_simulacoes;
DROP TABLE IF EXISTS precificacao_config;
DROP TABLE IF EXISTS metas_ticket_medio;
DROP TABLE IF EXISTS metas_barbeiro;
DROP TABLE IF EXISTS metas_mensais;
DROP TABLE IF EXISTS dre_mensal;
DROP TABLE IF EXISTS fluxo_caixa_diario;
DROP TABLE IF EXISTS compensacoes_bancarias;
DROP TABLE IF EXISTS contas_a_receber;
DROP TABLE IF EXISTS contas_a_pagar;
DROP TABLE IF EXISTS despesas;
DROP TABLE IF EXISTS receitas;
DROP TABLE IF EXISTS assinatura_invoices;
DROP TABLE IF EXISTS assinaturas;
DROP TABLE IF EXISTS planos_assinatura;
DROP TABLE IF EXISTS meios_pagamento;
DROP TABLE IF EXISTS clientes;
DROP TABLE IF EXISTS profissionais;
DROP TABLE IF EXISTS servicos;
DROP TABLE IF EXISTS produtos;
DROP TABLE IF EXISTS fornecedores;
DROP TABLE IF EXISTS categorias;
DROP TABLE IF EXISTS tenants;
