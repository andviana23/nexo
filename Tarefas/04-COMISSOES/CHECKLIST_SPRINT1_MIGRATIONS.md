# ✅ CHECKLIST — SPRINT 1: MIGRATIONS + QUERIES

> **Status:** ✅ Concluído  
> **Dependência:** Pacote 03-FINANCEIRO Sprint 1 (✅)  
> **Esforço Estimado:** 15 horas  
> **Prioridade:** P0 — Bloqueia todo o módulo

---

## 📊 RESUMO

```
████████████████████ 100% COMPLETO
```

| Categoria | Completo | Pendente |
|-----------|:--------:|:--------:|
| Migrations | 4/4 | 0 |
| Schema sqlc | 4/4 | 0 |
| Queries sqlc | 4/4 | 0 |
| sqlc generate | 1/1 | 0 |
| Testes migration | 1/1 | 0 |

---

## 1️⃣ DATABASE — MIGRATIONS

### 1.1 Migration: `commission_rules` (Esforço: 3h)

- [ ] Criar arquivo `backend/migrations/XXX_commission_rules.up.sql`
- [ ] Criar arquivo `backend/migrations/XXX_commission_rules.down.sql`

#### SQL UP

```sql
-- XXX_commission_rules.up.sql

CREATE TABLE IF NOT EXISTS commission_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE CASCADE,
    professional_id UUID REFERENCES profissionais(id) ON DELETE CASCADE,
    service_id UUID REFERENCES servicos(id) ON DELETE CASCADE,
    
    -- Tipo de comissão
    type VARCHAR(20) NOT NULL DEFAULT 'PERCENTAGE',
    -- Valores: PERCENTAGE, FIXED, HYBRID, PROGRESSIVE
    
    -- Valor principal (% para PERCENTAGE, R$ para FIXED)
    value NUMERIC(10,2) NOT NULL DEFAULT 0,
    
    -- Valor fixo (para HYBRID)
    fixed_value NUMERIC(10,2) DEFAULT 0,
    
    -- Faixas progressivas (JSON)
    -- Exemplo: [{"min": 0, "max": 5000, "pct": 40}, {"min": 5000, "pct": 50}]
    tiers JSONB,
    
    -- Prioridade (quanto maior, mais prioritário)
    priority INTEGER NOT NULL DEFAULT 0,
    
    -- Status
    active BOOLEAN NOT NULL DEFAULT true,
    
    -- Auditoria
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    
    -- Constraints
    CONSTRAINT valid_commission_type CHECK (type IN ('PERCENTAGE', 'FIXED', 'HYBRID', 'PROGRESSIVE')),
    CONSTRAINT valid_percentage CHECK (type != 'PERCENTAGE' OR (value >= 0 AND value <= 100)),
    CONSTRAINT valid_fixed_value CHECK (fixed_value >= 0),
    CONSTRAINT valid_priority CHECK (priority >= 0)
);

-- Índices
CREATE INDEX idx_commission_rules_tenant ON commission_rules(tenant_id);
CREATE INDEX idx_commission_rules_unit ON commission_rules(tenant_id, unit_id);
CREATE INDEX idx_commission_rules_professional ON commission_rules(tenant_id, professional_id);
CREATE INDEX idx_commission_rules_service ON commission_rules(tenant_id, service_id);
CREATE INDEX idx_commission_rules_lookup ON commission_rules(tenant_id, unit_id, professional_id, service_id) WHERE active = true;
CREATE INDEX idx_commission_rules_priority ON commission_rules(tenant_id, priority DESC) WHERE active = true;

-- Trigger updated_at
CREATE TRIGGER update_commission_rules_updated_at
    BEFORE UPDATE ON commission_rules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentários
COMMENT ON TABLE commission_rules IS 'Regras de comissão configuráveis por tenant/unidade/profissional/serviço';
COMMENT ON COLUMN commission_rules.type IS 'Tipo: PERCENTAGE (%), FIXED (R$), HYBRID (fixo + %), PROGRESSIVE (faixas)';
COMMENT ON COLUMN commission_rules.tiers IS 'Faixas para PROGRESSIVE: [{"min": 0, "max": 5000, "pct": 40}]';
```

#### SQL DOWN

```sql
-- XXX_commission_rules.down.sql

DROP TRIGGER IF EXISTS update_commission_rules_updated_at ON commission_rules;
DROP TABLE IF EXISTS commission_rules;
```

#### Checklist

- [ ] Arquivo UP criado
- [ ] Arquivo DOWN criado
- [ ] Constraints validados
- [ ] Índices criados
- [ ] Trigger updated_at
- [ ] Comentários adicionados

---

### 1.2 Migration: `commission_periods` (Esforço: 3h)

- [ ] Criar arquivo `backend/migrations/XXX_commission_periods.up.sql`
- [ ] Criar arquivo `backend/migrations/XXX_commission_periods.down.sql`

#### SQL UP

```sql
-- XXX_commission_periods.up.sql

CREATE TABLE IF NOT EXISTS commission_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE RESTRICT,
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE RESTRICT,
    
    -- Período
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    
    -- Totais
    total_services NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_products NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_commission NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_bonus NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_deductions NUMERIC(15,2) NOT NULL DEFAULT 0,
    net_value NUMERIC(15,2) NOT NULL DEFAULT 0,
    
    -- Quantidades
    qty_services INTEGER NOT NULL DEFAULT 0,
    qty_products INTEGER NOT NULL DEFAULT 0,
    
    -- Status: DRAFT, CLOSED, PAID
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    
    -- Referência ao título gerado
    bill_id UUID REFERENCES contas_a_pagar(id),
    
    -- Observações
    notes TEXT,
    
    -- Auditoria
    closed_at TIMESTAMPTZ,
    closed_by UUID REFERENCES users(id),
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT valid_period_status CHECK (status IN ('DRAFT', 'CLOSED', 'PAID')),
    CONSTRAINT valid_dates CHECK (end_date >= start_date),
    CONSTRAINT valid_totals CHECK (total_services >= 0 AND total_products >= 0 AND total_commission >= 0),
    CONSTRAINT closed_must_have_closer CHECK (status != 'CLOSED' OR closed_by IS NOT NULL)
);

-- Índices
CREATE INDEX idx_commission_periods_tenant ON commission_periods(tenant_id);
CREATE INDEX idx_commission_periods_professional ON commission_periods(tenant_id, professional_id);
CREATE INDEX idx_commission_periods_unit ON commission_periods(tenant_id, unit_id);
CREATE INDEX idx_commission_periods_dates ON commission_periods(tenant_id, start_date, end_date);
CREATE INDEX idx_commission_periods_status ON commission_periods(tenant_id, status);
CREATE INDEX idx_commission_periods_bill ON commission_periods(bill_id) WHERE bill_id IS NOT NULL;

-- Unique: Apenas um período DRAFT por profissional
CREATE UNIQUE INDEX idx_commission_periods_unique_draft 
    ON commission_periods(tenant_id, professional_id, start_date, end_date) 
    WHERE status = 'DRAFT';

-- Trigger updated_at
CREATE TRIGGER update_commission_periods_updated_at
    BEFORE UPDATE ON commission_periods
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentários
COMMENT ON TABLE commission_periods IS 'Períodos de fechamento de comissão (folha)';
COMMENT ON COLUMN commission_periods.status IS 'Status: DRAFT (rascunho), CLOSED (fechado), PAID (pago)';
COMMENT ON COLUMN commission_periods.bill_id IS 'Referência à conta a pagar gerada no fechamento';
```

#### SQL DOWN

```sql
-- XXX_commission_periods.down.sql

DROP TRIGGER IF EXISTS update_commission_periods_updated_at ON commission_periods;
DROP TABLE IF EXISTS commission_periods;
```

#### Checklist

- [ ] Arquivo UP criado
- [ ] Arquivo DOWN criado
- [ ] FK para contas_a_pagar
- [ ] Unique index para evitar duplicatas
- [ ] Status workflow válido

---

### 1.3 Migration: `advances` (Esforço: 2h)

- [ ] Criar arquivo `backend/migrations/XXX_advances.up.sql`
- [ ] Criar arquivo `backend/migrations/XXX_advances.down.sql`

#### SQL UP

```sql
-- XXX_advances.up.sql

CREATE TABLE IF NOT EXISTS advances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE RESTRICT,
    professional_id UUID NOT NULL REFERENCES profissionais(id) ON DELETE RESTRICT,
    
    -- Valor
    amount NUMERIC(15,2) NOT NULL,
    
    -- Datas
    request_date DATE NOT NULL DEFAULT CURRENT_DATE,
    
    -- Motivo
    reason TEXT,
    
    -- Status: PENDING, APPROVED, REJECTED, DEDUCTED
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    
    -- Aprovação
    approved_at TIMESTAMPTZ,
    approved_by UUID REFERENCES users(id),
    
    -- Rejeição
    rejected_at TIMESTAMPTZ,
    rejected_by UUID REFERENCES users(id),
    rejection_reason TEXT,
    
    -- Dedução
    deducted_in UUID REFERENCES commission_periods(id),
    deducted_at TIMESTAMPTZ,
    
    -- Auditoria
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    
    -- Constraints
    CONSTRAINT valid_advance_status CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED', 'DEDUCTED')),
    CONSTRAINT positive_amount CHECK (amount > 0),
    CONSTRAINT approved_must_have_approver CHECK (status != 'APPROVED' OR approved_by IS NOT NULL),
    CONSTRAINT rejected_must_have_rejector CHECK (status != 'REJECTED' OR rejected_by IS NOT NULL),
    CONSTRAINT deducted_must_have_period CHECK (status != 'DEDUCTED' OR deducted_in IS NOT NULL)
);

-- Índices
CREATE INDEX idx_advances_tenant ON advances(tenant_id);
CREATE INDEX idx_advances_professional ON advances(tenant_id, professional_id);
CREATE INDEX idx_advances_unit ON advances(tenant_id, unit_id);
CREATE INDEX idx_advances_status ON advances(tenant_id, status);
CREATE INDEX idx_advances_pending ON advances(tenant_id, professional_id) WHERE status = 'PENDING';
CREATE INDEX idx_advances_approved ON advances(tenant_id, professional_id) WHERE status = 'APPROVED';
CREATE INDEX idx_advances_deducted_in ON advances(deducted_in) WHERE deducted_in IS NOT NULL;

-- Trigger updated_at
CREATE TRIGGER update_advances_updated_at
    BEFORE UPDATE ON advances
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentários
COMMENT ON TABLE advances IS 'Adiantamentos/vales solicitados por profissionais';
COMMENT ON COLUMN advances.status IS 'Status: PENDING, APPROVED, REJECTED, DEDUCTED';
COMMENT ON COLUMN advances.deducted_in IS 'Período onde o adiantamento foi deduzido';
```

#### SQL DOWN

```sql
-- XXX_advances.down.sql

DROP TRIGGER IF EXISTS update_advances_updated_at ON advances;
DROP TABLE IF EXISTS advances;
```

#### Checklist

- [ ] Arquivo UP criado
- [ ] Arquivo DOWN criado
- [ ] FK para commission_periods
- [ ] Constraints de status
- [ ] Índices otimizados

---

### 1.4 Migration: Alter `barber_commissions` (Esforço: 1h)

- [ ] Criar arquivo `backend/migrations/XXX_alter_barber_commissions.up.sql`
- [ ] Criar arquivo `backend/migrations/XXX_alter_barber_commissions.down.sql`

#### SQL UP

```sql
-- XXX_alter_barber_commissions.up.sql

-- Adicionar coluna para rastrear item da comanda
ALTER TABLE barber_commissions 
ADD COLUMN IF NOT EXISTS command_item_id UUID REFERENCES command_items(id);

-- Adicionar coluna para rastrear período de fechamento
ALTER TABLE barber_commissions 
ADD COLUMN IF NOT EXISTS commission_period_id UUID REFERENCES commission_periods(id);

-- Adicionar coluna unit_id se não existir
ALTER TABLE barber_commissions 
ADD COLUMN IF NOT EXISTS unit_id UUID REFERENCES units(id);

-- Índices
CREATE INDEX IF NOT EXISTS idx_barber_commissions_command_item ON barber_commissions(command_item_id);
CREATE INDEX IF NOT EXISTS idx_barber_commissions_period ON barber_commissions(commission_period_id);
CREATE INDEX IF NOT EXISTS idx_barber_commissions_unit ON barber_commissions(tenant_id, unit_id);
CREATE INDEX IF NOT EXISTS idx_barber_commissions_pending ON barber_commissions(tenant_id, barbeiro_id) WHERE status = 'PENDENTE';
```

#### SQL DOWN

```sql
-- XXX_alter_barber_commissions.down.sql

DROP INDEX IF EXISTS idx_barber_commissions_pending;
DROP INDEX IF EXISTS idx_barber_commissions_unit;
DROP INDEX IF EXISTS idx_barber_commissions_period;
DROP INDEX IF EXISTS idx_barber_commissions_command_item;

ALTER TABLE barber_commissions DROP COLUMN IF EXISTS unit_id;
ALTER TABLE barber_commissions DROP COLUMN IF EXISTS commission_period_id;
ALTER TABLE barber_commissions DROP COLUMN IF EXISTS command_item_id;
```

#### Checklist

- [ ] Arquivo UP criado
- [ ] Arquivo DOWN criado
- [ ] Colunas adicionadas
- [ ] Índices otimizados

---

## 2️⃣ SQL QUERIES — sqlc

### 2.1 Arquivo: `commission_rules.sql` (Esforço: 2h)

- [ ] Criar `backend/internal/infra/db/queries/commission_rules.sql`

#### Queries a Implementar

```sql
-- name: CreateCommissionRule :one
INSERT INTO commission_rules (
    tenant_id, unit_id, professional_id, service_id,
    type, value, fixed_value, tiers, priority, active,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetCommissionRuleByID :one
SELECT * FROM commission_rules
WHERE id = $1 AND tenant_id = $2;

-- name: ListCommissionRulesByTenant :many
SELECT * FROM commission_rules
WHERE tenant_id = $1
ORDER BY priority DESC, created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListActiveCommissionRules :many
SELECT * FROM commission_rules
WHERE tenant_id = $1 AND active = true
ORDER BY priority DESC;

-- name: ListCommissionRulesByUnit :many
SELECT * FROM commission_rules
WHERE tenant_id = $1 AND unit_id = $2 AND active = true
ORDER BY priority DESC;

-- name: ListCommissionRulesByProfessional :many
SELECT * FROM commission_rules
WHERE tenant_id = $1 AND professional_id = $2 AND active = true
ORDER BY priority DESC;

-- name: ListCommissionRulesByService :many
SELECT * FROM commission_rules
WHERE tenant_id = $1 AND service_id = $2 AND active = true
ORDER BY priority DESC;

-- name: FindApplicableRule :one
-- Busca a regra mais específica aplicável
SELECT * FROM commission_rules
WHERE tenant_id = $1
  AND active = true
  AND (unit_id = $2 OR unit_id IS NULL)
  AND (professional_id = $3 OR professional_id IS NULL)
  AND (service_id = $4 OR service_id IS NULL)
ORDER BY 
    CASE WHEN service_id IS NOT NULL THEN 4 ELSE 0 END +
    CASE WHEN professional_id IS NOT NULL THEN 2 ELSE 0 END +
    CASE WHEN unit_id IS NOT NULL THEN 1 ELSE 0 END DESC,
    priority DESC
LIMIT 1;

-- name: UpdateCommissionRule :one
UPDATE commission_rules
SET
    unit_id = $3,
    professional_id = $4,
    service_id = $5,
    type = $6,
    value = $7,
    fixed_value = $8,
    tiers = $9,
    priority = $10,
    active = $11,
    updated_by = $12,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: ToggleCommissionRule :one
UPDATE commission_rules
SET active = NOT active, updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteCommissionRule :exec
DELETE FROM commission_rules
WHERE id = $1 AND tenant_id = $2;

-- name: CountCommissionRules :one
SELECT COUNT(*) FROM commission_rules
WHERE tenant_id = $1;

-- name: CountActiveCommissionRules :one
SELECT COUNT(*) FROM commission_rules
WHERE tenant_id = $1 AND active = true;
```

#### Checklist Queries

- [ ] CreateCommissionRule
- [ ] GetCommissionRuleByID
- [ ] ListCommissionRulesByTenant
- [ ] ListActiveCommissionRules
- [ ] ListCommissionRulesByUnit
- [ ] ListCommissionRulesByProfessional
- [ ] ListCommissionRulesByService
- [ ] FindApplicableRule (hierarquia)
- [ ] UpdateCommissionRule
- [ ] ToggleCommissionRule
- [ ] DeleteCommissionRule
- [ ] CountCommissionRules
- [ ] CountActiveCommissionRules

---

### 2.2 Arquivo: `commission_periods.sql` (Esforço: 2h)

- [ ] Criar `backend/internal/infra/db/queries/commission_periods.sql`

#### Queries a Implementar

```sql
-- name: CreateCommissionPeriod :one
INSERT INTO commission_periods (
    tenant_id, unit_id, professional_id,
    start_date, end_date,
    total_services, total_products, total_commission,
    total_bonus, total_deductions, net_value,
    qty_services, qty_products, status, notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING *;

-- name: GetCommissionPeriodByID :one
SELECT * FROM commission_periods
WHERE id = $1 AND tenant_id = $2;

-- name: ListCommissionPeriodsByTenant :many
SELECT * FROM commission_periods
WHERE tenant_id = $1
ORDER BY start_date DESC
LIMIT $2 OFFSET $3;

-- name: ListCommissionPeriodsByProfessional :many
SELECT * FROM commission_periods
WHERE tenant_id = $1 AND professional_id = $2
ORDER BY start_date DESC
LIMIT $3 OFFSET $4;

-- name: ListCommissionPeriodsByUnit :many
SELECT * FROM commission_periods
WHERE tenant_id = $1 AND unit_id = $2
ORDER BY start_date DESC
LIMIT $3 OFFSET $4;

-- name: ListCommissionPeriodsByStatus :many
SELECT * FROM commission_periods
WHERE tenant_id = $1 AND status = $2
ORDER BY start_date DESC
LIMIT $3 OFFSET $4;

-- name: ListCommissionPeriodsByDateRange :many
SELECT * FROM commission_periods
WHERE tenant_id = $1
  AND start_date >= $2
  AND end_date <= $3
ORDER BY start_date DESC;

-- name: GetDraftPeriod :one
SELECT * FROM commission_periods
WHERE tenant_id = $1 
  AND professional_id = $2
  AND status = 'DRAFT'
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateCommissionPeriod :one
UPDATE commission_periods
SET
    total_services = $3,
    total_products = $4,
    total_commission = $5,
    total_bonus = $6,
    total_deductions = $7,
    net_value = $8,
    qty_services = $9,
    qty_products = $10,
    notes = $11,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: ClosePeriod :one
UPDATE commission_periods
SET
    status = 'CLOSED',
    closed_at = NOW(),
    closed_by = $3,
    bill_id = $4,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status = 'DRAFT'
RETURNING *;

-- name: MarkPeriodAsPaid :one
UPDATE commission_periods
SET
    status = 'PAID',
    paid_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status = 'CLOSED'
RETURNING *;

-- name: DeleteCommissionPeriod :exec
DELETE FROM commission_periods
WHERE id = $1 AND tenant_id = $2 AND status = 'DRAFT';

-- name: SumCommissionsByProfessionalAndPeriod :one
SELECT 
    COALESCE(SUM(valor), 0) as total_commission,
    COUNT(*) as qty_items
FROM barber_commissions
WHERE tenant_id = $1
  AND barbeiro_id = $2
  AND data_competencia >= $3
  AND data_competencia <= $4
  AND status = 'PENDENTE';

-- name: CountPeriodsByStatus :one
SELECT COUNT(*) FROM commission_periods
WHERE tenant_id = $1 AND status = $2;
```

#### Checklist Queries

- [ ] CreateCommissionPeriod
- [ ] GetCommissionPeriodByID
- [ ] ListCommissionPeriodsByTenant
- [ ] ListCommissionPeriodsByProfessional
- [ ] ListCommissionPeriodsByUnit
- [ ] ListCommissionPeriodsByStatus
- [ ] ListCommissionPeriodsByDateRange
- [ ] GetDraftPeriod
- [ ] UpdateCommissionPeriod
- [ ] ClosePeriod
- [ ] MarkPeriodAsPaid
- [ ] DeleteCommissionPeriod
- [ ] SumCommissionsByProfessionalAndPeriod
- [ ] CountPeriodsByStatus

---

### 2.3 Arquivo: `advances.sql` (Esforço: 1.5h)

- [ ] Criar `backend/internal/infra/db/queries/advances.sql`

#### Queries a Implementar

```sql
-- name: CreateAdvance :one
INSERT INTO advances (
    tenant_id, unit_id, professional_id,
    amount, request_date, reason, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetAdvanceByID :one
SELECT * FROM advances
WHERE id = $1 AND tenant_id = $2;

-- name: ListAdvancesByTenant :many
SELECT * FROM advances
WHERE tenant_id = $1
ORDER BY request_date DESC
LIMIT $2 OFFSET $3;

-- name: ListAdvancesByProfessional :many
SELECT * FROM advances
WHERE tenant_id = $1 AND professional_id = $2
ORDER BY request_date DESC
LIMIT $3 OFFSET $4;

-- name: ListAdvancesByStatus :many
SELECT * FROM advances
WHERE tenant_id = $1 AND status = $2
ORDER BY request_date DESC
LIMIT $3 OFFSET $4;

-- name: ListPendingAdvances :many
SELECT * FROM advances
WHERE tenant_id = $1 AND status = 'PENDING'
ORDER BY request_date ASC;

-- name: ListApprovedAdvancesNotDeducted :many
SELECT * FROM advances
WHERE tenant_id = $1 
  AND professional_id = $2
  AND status = 'APPROVED'
  AND deducted_in IS NULL
ORDER BY request_date ASC;

-- name: ApproveAdvance :one
UPDATE advances
SET
    status = 'APPROVED',
    approved_at = NOW(),
    approved_by = $3,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status = 'PENDING'
RETURNING *;

-- name: RejectAdvance :one
UPDATE advances
SET
    status = 'REJECTED',
    rejected_at = NOW(),
    rejected_by = $3,
    rejection_reason = $4,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status = 'PENDING'
RETURNING *;

-- name: DeductAdvance :one
UPDATE advances
SET
    status = 'DEDUCTED',
    deducted_in = $3,
    deducted_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND status = 'APPROVED'
RETURNING *;

-- name: DeleteAdvance :exec
DELETE FROM advances
WHERE id = $1 AND tenant_id = $2 AND status = 'PENDING';

-- name: SumApprovedAdvancesNotDeducted :one
SELECT COALESCE(SUM(amount), 0) as total
FROM advances
WHERE tenant_id = $1 
  AND professional_id = $2
  AND status = 'APPROVED'
  AND deducted_in IS NULL;

-- name: CountAdvancesByStatus :one
SELECT COUNT(*) FROM advances
WHERE tenant_id = $1 AND status = $2;
```

#### Checklist Queries

- [ ] CreateAdvance
- [ ] GetAdvanceByID
- [ ] ListAdvancesByTenant
- [ ] ListAdvancesByProfessional
- [ ] ListAdvancesByStatus
- [ ] ListPendingAdvances
- [ ] ListApprovedAdvancesNotDeducted
- [ ] ApproveAdvance
- [ ] RejectAdvance
- [ ] DeductAdvance
- [ ] DeleteAdvance
- [ ] SumApprovedAdvancesNotDeducted
- [ ] CountAdvancesByStatus

---

### 2.4 Arquivo: `barber_commissions.sql` — Ajustes (Esforço: 0.5h)

- [ ] Atualizar `backend/internal/infra/db/queries/barber_commissions.sql`

#### Queries a Adicionar

```sql
-- name: CreateBarberCommissionFromCommand :one
INSERT INTO barber_commissions (
    tenant_id, barbeiro_id, receita_id, command_item_id,
    unit_id, valor, status, data_competencia, manual
) VALUES (
    $1, $2, $3, $4, $5, $6, 'PENDENTE', $7, false
) RETURNING *;

-- name: ListPendingCommissionsByProfessional :many
SELECT * FROM barber_commissions
WHERE tenant_id = $1 
  AND barbeiro_id = $2
  AND status = 'PENDENTE'
ORDER BY data_competencia ASC;

-- name: ListPendingCommissionsByPeriod :many
SELECT * FROM barber_commissions
WHERE tenant_id = $1 
  AND barbeiro_id = $2
  AND data_competencia >= $3
  AND data_competencia <= $4
  AND status = 'PENDENTE'
ORDER BY data_competencia ASC;

-- name: MarkCommissionsAsProcessed :exec
UPDATE barber_commissions
SET 
    status = 'PROCESSADO',
    commission_period_id = $3
WHERE tenant_id = $1
  AND barbeiro_id = $2
  AND status = 'PENDENTE'
  AND data_competencia >= $4
  AND data_competencia <= $5;

-- name: MarkCommissionsAsPaid :exec
UPDATE barber_commissions
SET status = 'PAGO'
WHERE commission_period_id = $1;
```

#### Checklist Queries Adicionais

- [ ] CreateBarberCommissionFromCommand
- [ ] ListPendingCommissionsByProfessional
- [ ] ListPendingCommissionsByPeriod
- [ ] MarkCommissionsAsProcessed
- [ ] MarkCommissionsAsPaid

---

## 3️⃣ GERAÇÃO sqlc

### 3.1 Gerar Código

- [ ] Executar `sqlc generate`
- [ ] Verificar arquivos gerados em `internal/infra/db/sqlc/`
- [ ] Corrigir erros de compilação (se houver)
- [ ] Commit dos arquivos gerados

### 3.2 Validação

- [ ] Compilar projeto: `go build ./...`
- [ ] Sem erros de tipo
- [ ] Sem imports faltando

---

## 4️⃣ TESTES DE MIGRATION

### 4.1 Testar UP

```bash
make migrate-up
# ou
cd backend && go run cmd/migrate/main.go up
```

- [ ] Migration executada sem erros
- [ ] Tabelas criadas corretamente
- [ ] Índices criados
- [ ] FKs válidas

### 4.2 Testar DOWN

```bash
make migrate-down
# ou
cd backend && go run cmd/migrate/main.go down 4
```

- [ ] Rollback executado sem erros
- [ ] Tabelas removidas
- [ ] Sem resíduos

### 4.3 Verificar Schema

```sql
-- Verificar tabelas criadas
\dt commission*
\dt advances

-- Verificar índices
\di+ idx_commission*
\di+ idx_advances*

-- Verificar constraints
\d commission_rules
\d commission_periods
\d advances
```

---

## 📝 NOTAS

### Próximos Passos

Após completar esta sprint:
1. Iniciar Sprint 2 (Domain + Repository + UseCases)
2. Checklist: `CHECKLIST_SPRINT2_BACKEND.md`

### Arquivos Criados

| Arquivo | Status |
|---------|--------|
| `XXX_commission_rules.up.sql` | ❌ |
| `XXX_commission_rules.down.sql` | ❌ |
| `XXX_commission_periods.up.sql` | ❌ |
| `XXX_commission_periods.down.sql` | ❌ |
| `XXX_advances.up.sql` | ❌ |
| `XXX_advances.down.sql` | ❌ |
| `XXX_alter_barber_commissions.up.sql` | ❌ |
| `XXX_alter_barber_commissions.down.sql` | ❌ |
| `queries/commission_rules.sql` | ❌ |
| `queries/commission_periods.sql` | ❌ |
| `queries/advances.sql` | ❌ |

---

*Checklist criado em: 05/12/2025*
