-- +goose Up
-- Migration 006: Verificação de tabelas de agendamentos e categorias de serviço
-- Data: 27/11/2025
-- Nota: Tabelas já existem no banco (criadas via SQL direto ou migration anterior)
-- Esta migration garante consistência e adiciona índices faltantes

-- ============================================
-- 1. CATEGORIAS DE SERVIÇO (já existe)
-- ============================================
CREATE TABLE IF NOT EXISTS categorias_servicos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    cor VARCHAR(7) DEFAULT '#6366f1',
    icone VARCHAR(50),
    ordem INTEGER DEFAULT 0,
    ativo BOOLEAN DEFAULT true,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT categorias_servicos_tenant_nome_unique UNIQUE (tenant_id, nome)
);

-- Índices para categorias_servicos
CREATE INDEX IF NOT EXISTS idx_categorias_servicos_tenant ON categorias_servicos(tenant_id);
CREATE INDEX IF NOT EXISTS idx_categorias_servicos_ativo ON categorias_servicos(tenant_id, ativo);

-- ============================================
-- 2. ADICIONAR COLUNA categoria_servico_id EM SERVICOS
-- ============================================
ALTER TABLE servicos 
ADD COLUMN IF NOT EXISTS categoria_servico_id UUID REFERENCES categorias_servicos(id) ON DELETE SET NULL;

-- Índice para busca por categoria de serviço
CREATE INDEX IF NOT EXISTS idx_servicos_categoria_servico ON servicos(categoria_servico_id);

-- ============================================
-- 3. AGENDAMENTOS (APPOINTMENTS)
-- ============================================
CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    cliente_id UUID NOT NULL REFERENCES clientes(id) ON DELETE RESTRICT,
    profissional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE RESTRICT,
    data_hora_inicio TIMESTAMPTZ NOT NULL,
    data_hora_fim TIMESTAMPTZ NOT NULL,
    status VARCHAR(30) DEFAULT 'AGENDADO' NOT NULL CHECK (status IN (
        'AGENDADO',
        'CONFIRMADO', 
        'EM_ATENDIMENTO',
        'CONCLUIDO',
        'CANCELADO',
        'NAO_COMPARECEU'
    )),
    observacoes TEXT,
    valor_total NUMERIC(15,2) DEFAULT 0,
    valor_desconto NUMERIC(15,2) DEFAULT 0,
    valor_final NUMERIC(15,2) DEFAULT 0,
    cupom_id UUID REFERENCES cupons_desconto(id) ON DELETE SET NULL,
    meio_pagamento_id UUID REFERENCES meios_pagamento(id) ON DELETE SET NULL,
    origem VARCHAR(30) DEFAULT 'SISTEMA' CHECK (origem IN ('SISTEMA', 'WHATSAPP', 'APP_CLIENTE', 'GOOGLE_CALENDAR')),
    google_event_id VARCHAR(255),
    lembrete_enviado BOOLEAN DEFAULT false,
    cancelado_por UUID REFERENCES users(id) ON DELETE SET NULL,
    motivo_cancelamento TEXT,
    criado_em TIMESTAMPTZ DEFAULT now(),
    atualizado_em TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT chk_data_hora CHECK (data_hora_fim > data_hora_inicio),
    CONSTRAINT chk_valores CHECK (valor_final >= 0 AND valor_desconto >= 0 AND valor_total >= 0)
);

-- Índices para appointments
CREATE INDEX IF NOT EXISTS idx_appointments_tenant ON appointments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_appointments_cliente ON appointments(cliente_id);
CREATE INDEX IF NOT EXISTS idx_appointments_profissional ON appointments(profissional_id);
CREATE INDEX IF NOT EXISTS idx_appointments_data ON appointments(tenant_id, data_hora_inicio);
CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_appointments_profissional_data ON appointments(profissional_id, data_hora_inicio);

-- ============================================
-- 4. SERVIÇOS DO AGENDAMENTO (APPOINTMENT_SERVICES)
-- ============================================
CREATE TABLE IF NOT EXISTS appointment_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    appointment_id UUID NOT NULL REFERENCES appointments(id) ON DELETE CASCADE,
    service_id UUID NOT NULL REFERENCES servicos(id) ON DELETE RESTRICT,
    preco NUMERIC(15,2) NOT NULL CHECK (preco >= 0),
    duracao INTEGER NOT NULL CHECK (duracao > 0),
    ordem INTEGER DEFAULT 0,
    criado_em TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT appointment_services_unique UNIQUE (appointment_id, service_id)
);

-- Índices para appointment_services
CREATE INDEX IF NOT EXISTS idx_appointment_services_appointment ON appointment_services(appointment_id);
CREATE INDEX IF NOT EXISTS idx_appointment_services_service ON appointment_services(service_id);
CREATE INDEX IF NOT EXISTS idx_appointment_services_tenant ON appointment_services(tenant_id);

-- ============================================
-- 5. FUNÇÃO PARA ATUALIZAR updated_at
-- ============================================
CREATE OR REPLACE FUNCTION update_appointments_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.atualizado_em = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para appointments
DROP TRIGGER IF EXISTS trg_appointments_updated_at ON appointments;
CREATE TRIGGER trg_appointments_updated_at
    BEFORE UPDATE ON appointments
    FOR EACH ROW
    EXECUTE FUNCTION update_appointments_updated_at();

-- Trigger para categorias_servicos
DROP TRIGGER IF EXISTS trg_categorias_servicos_updated_at ON categorias_servicos;
CREATE TRIGGER trg_categorias_servicos_updated_at
    BEFORE UPDATE ON categorias_servicos
    FOR EACH ROW
    EXECUTE FUNCTION update_appointments_updated_at();

-- ============================================
-- 6. FUNÇÃO PARA VALIDAR CONFLITO DE HORÁRIO
-- ============================================
CREATE OR REPLACE FUNCTION check_appointment_conflict()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM appointments
        WHERE tenant_id = NEW.tenant_id
        AND profissional_id = NEW.profissional_id
        AND id != COALESCE(NEW.id, '00000000-0000-0000-0000-000000000000'::uuid)
        AND status NOT IN ('CANCELADO', 'NAO_COMPARECEU')
        AND (
            (NEW.data_hora_inicio >= data_hora_inicio AND NEW.data_hora_inicio < data_hora_fim)
            OR (NEW.data_hora_fim > data_hora_inicio AND NEW.data_hora_fim <= data_hora_fim)
            OR (NEW.data_hora_inicio <= data_hora_inicio AND NEW.data_hora_fim >= data_hora_fim)
        )
    ) THEN
        RAISE EXCEPTION 'Conflito de horário: profissional já possui agendamento neste período';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para validar conflitos
DROP TRIGGER IF EXISTS trg_check_appointment_conflict ON appointments;
CREATE TRIGGER trg_check_appointment_conflict
    BEFORE INSERT OR UPDATE ON appointments
    FOR EACH ROW
    EXECUTE FUNCTION check_appointment_conflict();

-- +goose Down
DROP TRIGGER IF EXISTS trg_check_appointment_conflict ON appointments;
DROP FUNCTION IF EXISTS check_appointment_conflict();
DROP TRIGGER IF EXISTS trg_categorias_servicos_updated_at ON categorias_servicos;
DROP TRIGGER IF EXISTS trg_appointments_updated_at ON appointments;
DROP FUNCTION IF EXISTS update_appointments_updated_at();
DROP TABLE IF EXISTS appointment_services;
DROP TABLE IF EXISTS appointments;
ALTER TABLE servicos DROP COLUMN IF EXISTS categoria_servico_id;
DROP TABLE IF EXISTS categorias_servicos;
