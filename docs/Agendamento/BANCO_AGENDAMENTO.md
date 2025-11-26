# Schema de Banco de Dados — Módulo de Agendamento | NEXO v1.0

**Versão:** 1.0.0  
**Data:** 25/11/2025  
**SGBD:** PostgreSQL 14+  
**ORM:** sqlc v1.30.0  

---

## 1. Diagrama ER (Entity-Relationship)

```
tenants (1) ──< (N) appointments
users/profissionais (1) ──< (N) appointments (as professional)
clientes (1) ──< (N) appointments (as customer)
appointments (1) ──< (N) appointment_services ──> (1) servicos
appointments (1) ──< (N) appointment_status_history
```

---

## 2. Tabelas Principais

### 2.1 `appointments`

**Descrição:** Tabela principal de agendamentos.

```sql
CREATE TABLE appointments (
    -- Identificação
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    
    -- Relacionamentos
    professional_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    unit_id UUID,
    
    -- Temporalidade
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- Status e Controle
    status VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    notes TEXT,
    canceled_reason TEXT,
    
    -- Integrações
    google_event_id VARCHAR(255),
    
    -- Auditoria
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    created_by UUID,
    updated_by UUID,
    
    -- Constraints
    CONSTRAINT chk_end_after_start 
        CHECK (end_time > start_time),
    
    CONSTRAINT chk_status_valid 
        CHECK (status IN ('CREATED', 'CONFIRMED', 'IN_SERVICE', 'DONE', 'NO_SHOW', 'CANCELED')),
    
    -- Foreign Keys
    CONSTRAINT fk_tenant 
        FOREIGN KEY (tenant_id) 
        REFERENCES tenants(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_professional 
        FOREIGN KEY (professional_id) 
        REFERENCES profissionais(id) 
        ON DELETE RESTRICT,
    
    CONSTRAINT fk_customer 
        FOREIGN KEY (customer_id) 
        REFERENCES clientes(id) 
        ON DELETE RESTRICT,
    
    CONSTRAINT fk_created_by 
        FOREIGN KEY (created_by) 
        REFERENCES users(id) 
        ON DELETE SET NULL,
    
    CONSTRAINT fk_updated_by 
        FOREIGN KEY (updated_by) 
        REFERENCES users(id) 
        ON DELETE SET NULL
);

-- Comentários
COMMENT ON TABLE appointments IS 'Agendamentos de serviços - core do sistema';
COMMENT ON COLUMN appointments.status IS 'CREATED → CONFIRMED → IN_SERVICE → DONE/NO_SHOW/CANCELED';
COMMENT ON COLUMN appointments.google_event_id IS 'ID do evento no Google Calendar (sync)';
COMMENT ON COLUMN appointments.canceled_reason IS 'Motivo do cancelamento (LGPD)';
```

### 2.2 `appointment_services`

**Descrição:** Tabela N:N entre agendamentos e serviços.

```sql
CREATE TABLE appointment_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL,
    service_id UUID NOT NULL,
    
    -- Snapshot de dados (para histórico)
    service_name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    duration_minutes INT NOT NULL,
    
    -- Auditoria
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- Constraints
    CONSTRAINT chk_price_positive CHECK (price >= 0),
    CONSTRAINT chk_duration_positive CHECK (duration_minutes > 0),
    
    -- Foreign Keys
    CONSTRAINT fk_appointment 
        FOREIGN KEY (appointment_id) 
        REFERENCES appointments(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_service 
        FOREIGN KEY (service_id) 
        REFERENCES servicos(id) 
        ON DELETE RESTRICT,
    
    -- Unique constraint (não pode duplicar serviço no mesmo agendamento)
    CONSTRAINT uq_appointment_service 
        UNIQUE (appointment_id, service_id)
);

COMMENT ON TABLE appointment_services IS 'Serviços associados a cada agendamento (N:N)';
COMMENT ON COLUMN appointment_services.service_name IS 'Snapshot do nome (imutável)';
COMMENT ON COLUMN appointment_services.price IS 'Snapshot do preço no momento do agendamento';
```

### 2.3 `appointment_status_history`

**Descrição:** Histórico de mudanças de status (auditoria).

```sql
CREATE TABLE appointment_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL,
    
    -- Transição
    from_status VARCHAR(20),
    to_status VARCHAR(20) NOT NULL,
    
    -- Contexto
    changed_by UUID,
    reason TEXT,
    
    -- Auditoria
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    -- Foreign Keys
    CONSTRAINT fk_appointment 
        FOREIGN KEY (appointment_id) 
        REFERENCES appointments(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_changed_by 
        FOREIGN KEY (changed_by) 
        REFERENCES users(id) 
        ON DELETE SET NULL
);

CREATE INDEX idx_status_history_appointment 
ON appointment_status_history(appointment_id, changed_at DESC);

COMMENT ON TABLE appointment_status_history IS 'Rastreabilidade de mudanças de status';
```

### 2.4 `blocked_times` (Futuro)

**Descrição:** Horários bloqueados (férias, almoço, etc).

```sql
CREATE TABLE blocked_times (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    professional_id UUID NOT NULL,
    
    -- Temporalidade
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- Metadados
    reason VARCHAR(255) NOT NULL,
    is_recurring BOOLEAN DEFAULT FALSE,
    recurrence_rule TEXT, -- iCal RRULE format
    
    -- Auditoria
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    created_by UUID,
    
    -- Constraints
    CONSTRAINT chk_end_after_start 
        CHECK (end_time > start_time),
    
    -- Foreign Keys
    CONSTRAINT fk_tenant 
        FOREIGN KEY (tenant_id) 
        REFERENCES tenants(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_professional 
        FOREIGN KEY (professional_id) 
        REFERENCES profissionais(id) 
        ON DELETE CASCADE
);

CREATE INDEX idx_blocked_times_professional 
ON blocked_times(tenant_id, professional_id, start_time);

COMMENT ON TABLE blocked_times IS 'Horários bloqueados para agendamento (férias, almoço)';
```

---

## 3. Índices de Performance

### 3.1 Índices Essenciais

```sql
-- Query principal: buscar agendamentos por barbeiro e data
CREATE INDEX idx_appointments_tenant_professional_date 
ON appointments(tenant_id, professional_id, start_time DESC)
WHERE status NOT IN ('CANCELED');

-- Query de conflitos (overlapping)
CREATE INDEX idx_appointments_time_range 
ON appointments(tenant_id, professional_id, start_time, end_time)
WHERE status NOT IN ('CANCELED', 'NO_SHOW');

-- Query por cliente (histórico)
CREATE INDEX idx_appointments_customer 
ON appointments(tenant_id, customer_id, start_time DESC);

-- Query por status
CREATE INDEX idx_appointments_status 
ON appointments(tenant_id, status, start_time DESC);

-- Query por data (calendário)
CREATE INDEX idx_appointments_date_range 
ON appointments(tenant_id, start_time, end_time)
WHERE status NOT IN ('CANCELED');

-- Google Calendar sync
CREATE INDEX idx_appointments_google_event 
ON appointments(google_event_id)
WHERE google_event_id IS NOT NULL;
```

### 3.2 Justificativa dos Índices

| Índice | Query Beneficiada | Frequência |
|--------|-------------------|------------|
| `idx_appointments_tenant_professional_date` | Calendário por barbeiro | Muito Alta |
| `idx_appointments_time_range` | Validação de conflitos | Alta |
| `idx_appointments_customer` | Histórico do cliente | Média |
| `idx_appointments_status` | Dashboard/filtros | Média |
| `idx_appointments_date_range` | Calendário geral | Alta |

---

## 4. Triggers e Funções

### 4.1 Trigger de Auditoria (Status History)

```sql
CREATE OR REPLACE FUNCTION fn_log_appointment_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        INSERT INTO appointment_status_history (
            appointment_id,
            from_status,
            to_status,
            changed_by
        ) VALUES (
            NEW.id,
            OLD.status,
            NEW.status,
            NEW.updated_by
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_appointment_status_change
    AFTER UPDATE ON appointments
    FOR EACH ROW
    WHEN (OLD.status IS DISTINCT FROM NEW.status)
    EXECUTE FUNCTION fn_log_appointment_status_change();

COMMENT ON FUNCTION fn_log_appointment_status_change() IS 'Registra mudanças de status automaticamente';
```

### 4.2 Trigger de Updated At

```sql
CREATE OR REPLACE FUNCTION fn_update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_appointments_updated_at
    BEFORE UPDATE ON appointments
    FOR EACH ROW
    EXECUTE FUNCTION fn_update_updated_at();
```

---

## 5. Views Úteis

### 5.1 `v_appointments_full`

**Descrição:** View com dados completos (JOIN de todas as tabelas).

```sql
CREATE OR REPLACE VIEW v_appointments_full AS
SELECT 
    a.id,
    a.tenant_id,
    a.start_time,
    a.end_time,
    a.status,
    a.notes,
    a.created_at,
    a.updated_at,
    
    -- Professional
    p.nome AS professional_name,
    p.email AS professional_email,
    
    -- Customer
    c.nome AS customer_name,
    c.telefone AS customer_phone,
    c.email AS customer_email,
    
    -- Services (array aggregation)
    ARRAY_AGG(DISTINCT s.nome) AS service_names,
    SUM(aps.price) AS total_price,
    SUM(aps.duration_minutes) AS total_duration,
    
    -- Calculated
    CASE 
        WHEN a.start_time < NOW() AND a.status = 'CONFIRMED' THEN TRUE
        ELSE FALSE
    END AS is_late,
    
    EXTRACT(EPOCH FROM (a.end_time - a.start_time))/60 AS duration_minutes
    
FROM appointments a
INNER JOIN profissionais p ON a.professional_id = p.id
INNER JOIN clientes c ON a.customer_id = c.id
LEFT JOIN appointment_services aps ON a.id = aps.appointment_id
LEFT JOIN servicos s ON aps.service_id = s.id
GROUP BY a.id, p.id, c.id;

COMMENT ON VIEW v_appointments_full IS 'View com dados completos de agendamentos (JOIN)';
```

### 5.2 `v_daily_schedule`

**Descrição:** Agenda diária por barbeiro.

```sql
CREATE OR REPLACE VIEW v_daily_schedule AS
SELECT 
    a.tenant_id,
    DATE(a.start_time) AS date,
    a.professional_id,
    p.nome AS professional_name,
    
    COUNT(*) AS total_appointments,
    COUNT(*) FILTER (WHERE a.status = 'CONFIRMED') AS confirmed_count,
    COUNT(*) FILTER (WHERE a.status = 'DONE') AS done_count,
    COUNT(*) FILTER (WHERE a.status = 'NO_SHOW') AS no_show_count,
    
    SUM(aps.price) AS total_revenue,
    
    MIN(a.start_time) AS first_appointment,
    MAX(a.end_time) AS last_appointment
    
FROM appointments a
INNER JOIN profissionais p ON a.professional_id = p.id
LEFT JOIN appointment_services aps ON a.id = aps.appointment_id
WHERE a.status NOT IN ('CANCELED')
GROUP BY a.tenant_id, DATE(a.start_time), a.professional_id, p.nome;

COMMENT ON VIEW v_daily_schedule IS 'Resumo diário da agenda por barbeiro';
```

---

## 6. Queries SQL Essenciais

### 6.1 Validar Conflitos (Overlapping)

```sql
-- Query usada no CheckConflicts do repository
SELECT * FROM appointments
WHERE tenant_id = $1
  AND professional_id = $2
  AND start_time < $4  -- novo.end_time
  AND end_time > $3    -- novo.start_time
  AND status NOT IN ('CANCELED', 'NO_SHOW')
  AND ($5::uuid IS NULL OR id != $5); -- Excluir ID próprio (updates)
```

**Explicação:**
- Conflito existe quando há sobreposição de tempo
- Ignora agendamentos cancelados/no-show
- Permite excluir o próprio ID (para updates)

### 6.2 Listar Agendamentos com Filtros

```sql
-- Query com todos os filtros possíveis
SELECT a.*, 
       p.nome AS professional_name,
       c.nome AS customer_name,
       ARRAY_AGG(s.nome) AS service_names
FROM appointments a
INNER JOIN profissionais p ON a.professional_id = p.id
INNER JOIN clientes c ON a.customer_id = c.id
LEFT JOIN appointment_services aps ON a.id = aps.appointment_id
LEFT JOIN servicos s ON aps.service_id = s.id
WHERE a.tenant_id = $1
  AND ($2::uuid IS NULL OR a.professional_id = $2)
  AND ($3::uuid IS NULL OR a.customer_id = $3)
  AND ($4::timestamp IS NULL OR a.start_time >= $4)
  AND ($5::timestamp IS NULL OR a.start_time <= $5)
  AND ($6::text IS NULL OR a.status = $6)
GROUP BY a.id, p.nome, c.nome
ORDER BY a.start_time ASC
LIMIT $7 OFFSET $8;
```

### 6.3 Slots Disponíveis (Complexo)

```sql
-- Gerar slots de 15 em 15 minutos e verificar disponibilidade
WITH time_slots AS (
    SELECT 
        generate_series(
            DATE($1) + TIME '08:00',
            DATE($1) + TIME '20:00',
            INTERVAL '15 minutes'
        ) AS slot_time
),
occupied_slots AS (
    SELECT 
        generate_series(
            DATE_TRUNC('minute', start_time),
            DATE_TRUNC('minute', end_time) - INTERVAL '1 minute',
            INTERVAL '15 minutes'
        ) AS slot_time
    FROM appointments
    WHERE tenant_id = $2
      AND professional_id = $3
      AND DATE(start_time) = $1
      AND status NOT IN ('CANCELED', 'NO_SHOW')
)
SELECT ts.slot_time, 
       CASE WHEN os.slot_time IS NULL THEN TRUE ELSE FALSE END AS is_available
FROM time_slots ts
LEFT JOIN occupied_slots os ON ts.slot_time = os.slot_time
ORDER BY ts.slot_time;
```

### 6.4 Dashboard: Ocupação do Barbeiro

```sql
-- Calcular taxa de ocupação (% do tempo com agendamentos)
SELECT 
    p.nome AS professional_name,
    COUNT(*) AS total_appointments,
    SUM(EXTRACT(EPOCH FROM (a.end_time - a.start_time))/3600) AS total_hours,
    
    -- Assumindo 8h/dia útil
    SUM(EXTRACT(EPOCH FROM (a.end_time - a.start_time))/3600) / (8.0 * COUNT(DISTINCT DATE(a.start_time))) * 100 AS occupation_rate
    
FROM appointments a
INNER JOIN profissionais p ON a.professional_id = p.id
WHERE a.tenant_id = $1
  AND a.start_time >= $2
  AND a.start_time <= $3
  AND a.status IN ('CONFIRMED', 'DONE')
GROUP BY p.id, p.nome
ORDER BY occupation_rate DESC;
```

---

## 7. Migrations (golang-migrate)

### 7.1 Migration Up

```sql
-- migrations/XXXX_create_appointments.up.sql

BEGIN;

-- Tabela principal
CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    professional_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    unit_id UUID,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    notes TEXT,
    canceled_reason TEXT,
    google_event_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    created_by UUID,
    updated_by UUID,
    
    CONSTRAINT chk_end_after_start CHECK (end_time > start_time),
    CONSTRAINT chk_status_valid CHECK (status IN ('CREATED', 'CONFIRMED', 'IN_SERVICE', 'DONE', 'NO_SHOW', 'CANCELED')),
    CONSTRAINT fk_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_professional FOREIGN KEY (professional_id) REFERENCES profissionais(id) ON DELETE RESTRICT,
    CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES clientes(id) ON DELETE RESTRICT
);

-- Tabela N:N
CREATE TABLE appointment_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL,
    service_id UUID NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    duration_minutes INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    CONSTRAINT chk_price_positive CHECK (price >= 0),
    CONSTRAINT chk_duration_positive CHECK (duration_minutes > 0),
    CONSTRAINT fk_appointment FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE,
    CONSTRAINT fk_service FOREIGN KEY (service_id) REFERENCES servicos(id) ON DELETE RESTRICT,
    CONSTRAINT uq_appointment_service UNIQUE (appointment_id, service_id)
);

-- Histórico de status
CREATE TABLE appointment_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL,
    from_status VARCHAR(20),
    to_status VARCHAR(20) NOT NULL,
    changed_by UUID,
    reason TEXT,
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    
    CONSTRAINT fk_appointment FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE,
    CONSTRAINT fk_changed_by FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Índices
CREATE INDEX idx_appointments_tenant_professional_date ON appointments(tenant_id, professional_id, start_time DESC);
CREATE INDEX idx_appointments_time_range ON appointments(tenant_id, professional_id, start_time, end_time) WHERE status NOT IN ('CANCELED', 'NO_SHOW');
CREATE INDEX idx_appointments_customer ON appointments(tenant_id, customer_id, start_time DESC);
CREATE INDEX idx_appointments_status ON appointments(tenant_id, status, start_time DESC);
CREATE INDEX idx_status_history_appointment ON appointment_status_history(appointment_id, changed_at DESC);

-- Triggers
CREATE OR REPLACE FUNCTION fn_log_appointment_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        INSERT INTO appointment_status_history (appointment_id, from_status, to_status, changed_by)
        VALUES (NEW.id, OLD.status, NEW.status, NEW.updated_by);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_appointment_status_change
    AFTER UPDATE ON appointments
    FOR EACH ROW
    WHEN (OLD.status IS DISTINCT FROM NEW.status)
    EXECUTE FUNCTION fn_log_appointment_status_change();

CREATE OR REPLACE FUNCTION fn_update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_appointments_updated_at
    BEFORE UPDATE ON appointments
    FOR EACH ROW
    EXECUTE FUNCTION fn_update_updated_at();

COMMIT;
```

### 7.2 Migration Down

```sql
-- migrations/XXXX_create_appointments.down.sql

BEGIN;

DROP TRIGGER IF EXISTS trg_appointments_updated_at ON appointments;
DROP TRIGGER IF EXISTS trg_appointment_status_change ON appointments;

DROP FUNCTION IF EXISTS fn_update_updated_at();
DROP FUNCTION IF EXISTS fn_log_appointment_status_change();

DROP TABLE IF EXISTS appointment_status_history CASCADE;
DROP TABLE IF EXISTS appointment_services CASCADE;
DROP TABLE IF EXISTS appointments CASCADE;

COMMIT;
```

---

## 8. Particionamento (Futuro - Alta Escala)

```sql
-- Particionar appointments por mês (para tenants grandes)
CREATE TABLE appointments_2025_01 PARTITION OF appointments
FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE appointments_2025_02 PARTITION OF appointments
FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');

-- Continuar para demais meses...
```

---

## 9. Backup e Retenção

### 9.1 Política de Retenção

| Tabela | Retenção | Soft Delete | Observação |
|--------|----------|-------------|------------|
| `appointments` | 2 anos | Não | Após 2 anos, arquivar em cold storage |
| `appointment_services` | 2 anos | Não | Manter snapshot histórico |
| `appointment_status_history` | Permanente | Não | Auditoria LGPD |

### 9.2 Arquivamento

```sql
-- Mover agendamentos antigos para tabela de arquivo
CREATE TABLE appointments_archive (LIKE appointments INCLUDING ALL);

INSERT INTO appointments_archive
SELECT * FROM appointments
WHERE start_time < NOW() - INTERVAL '2 years';

DELETE FROM appointments
WHERE start_time < NOW() - INTERVAL '2 years';
```

---

## 10. Checklist de Implementação

### Database Setup

- [ ] Criar migrations (up/down)
- [ ] Rodar migrations em dev
- [ ] Rodar migrations em staging
- [ ] Validar índices criados
- [ ] Validar triggers funcionando
- [ ] Popular tabela de teste com seed data
- [ ] Testar queries de conflito
- [ ] Testar performance com 10k registros
- [ ] Configurar backup automático
- [ ] Configurar retenção de logs

---

**Responsável:** DBA + Backend Dev  
**Data:** 25/11/2025  
**Status:** ✅ COMPLETO
