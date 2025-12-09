# Pendências Backend Resolvidas — CommandModal Integration

**Data:** 30/11/2025  
**Versão:** 1.0.0  
**Status:** ✅ Completo (aguarda build)

---

## 📋 Pendências Identificadas

| # | Pendência | Status |
|---|-----------|--------|
| 1 | Endpoint GET /api/v1/appointments/:id | ✅ Já existia |
| 2 | Campo `command_id` em `appointments` table | ✅ Migration criada |
| 3 | Campo `command_id` retornado em responses | ✅ Implementado |

---

## ✅ Alterações Implementadas

### 1️⃣ Migration: Adicionar command_id em appointments

**Arquivo:** `backend/migrations/032_add_command_id_to_appointments.up.sql`

```sql
ALTER TABLE appointments
ADD COLUMN IF NOT EXISTS command_id UUID REFERENCES commands(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_appointments_command_id 
ON appointments(command_id) 
WHERE command_id IS NOT NULL;
```

**Reversão:** `032_add_command_id_to_appointments.down.sql`

```sql
DROP INDEX IF EXISTS idx_appointments_command_id;
ALTER TABLE appointments DROP COLUMN IF EXISTS command_id;
```

---

### 2️⃣ DTO: AppointmentResponse

**Arquivo:** `backend/internal/application/dto/appointment_dto.go`

**Alteração:**
```go
type AppointmentResponse struct {
    // ... campos existentes
    CommandID string `json:"command_id,omitempty"` // ← NOVO
    // ... demais campos
}
```

**Impacto:**
- ✅ Todas as responses de appointments agora incluem `command_id`
- ✅ Campo é `omitempty` (não aparece se vazio)
- ✅ Frontend pode acessar via `appointment.command_id`

---

### 3️⃣ Entity: Appointment

**Arquivo:** `backend/internal/domain/entity/appointment.go`

**Alteração:**
```go
type Appointment struct {
    // ... campos existentes
    CommandID string // ← NOVO: Comanda vinculada ao agendamento
    // ... demais campos
}
```

---

### 4️⃣ Mapper: AppointmentToResponse

**Arquivo:** `backend/internal/application/mapper/appointment_mapper.go`

**Alteração:**
```go
func AppointmentToResponse(a *entity.Appointment) dto.AppointmentResponse {
    return dto.AppointmentResponse{
        // ... campos existentes
        CommandID: a.CommandID, // ← NOVO: Campo mapeado
        // ... demais campos
    }
}
```

---

### 5️⃣ SQLC Queries: appointments.sql

**Arquivo:** `backend/internal/infra/db/queries/appointments.sql`

**Alterações:**

**a) CreateAppointment:**
```sql
-- Adicionado command_id na lista de colunas
INSERT INTO appointments (
    id, tenant_id, professional_id, customer_id,
    start_time, end_time, status, total_price,
    notes, canceled_reason, google_calendar_event_id,
    command_id  -- ← NOVO
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12  -- ← $12 = command_id
) RETURNING *;
```

**b) UpdateAppointment:**
```sql
UPDATE appointments
SET
    professional_id = $3,
    start_time = $4,
    end_time = $5,
    status = $6,
    total_price = $7,
    notes = $8,
    canceled_reason = $9,
    google_calendar_event_id = $10,
    checked_in_at = $11,
    started_at = $12,
    finished_at = $13,
    command_id = $14,  -- ← NOVO
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;
```

**GetAppointmentByID:**
- ✅ Já retorna todos os campos (incluirá `command_id` após regenerar SQLC)

---

### 6️⃣ Repository: AppointmentRepository

**Arquivo:** `backend/internal/infra/repository/postgres/appointment_repository.go`

**Alterações:**

**a) Create:**
```go
params := db.CreateAppointmentParams{
    // ... campos existentes
    CommandID: uuidStrPtrToPgtype(appointment.CommandID), // ← NOVO
}
```

**b) Update:**
```go
params := db.UpdateAppointmentParams{
    // ... campos existentes
    CommandID: uuidStrPtrToPgtype(appointment.CommandID), // ← NOVO
}
```

**c) rowToDomain (todos os métodos):**
```go
return &entity.Appointment{
    // ... campos existentes
    CommandID: pgUUIDPtrToString(row.CommandID), // ← NOVO
    // ... demais campos
}
```

Métodos atualizados:
- `rowToDomain` (GetByID)
- `listRowToDomain` (List)
- `professionalRangeRowToDomain` (ListByProfessional)
- `customerRowToDomain` (ListByCustomer)

---

### 7️⃣ Helpers: Funções de conversão UUID nullable

**Arquivo:** `backend/internal/infra/repository/postgres/helpers.go`

**Novas funções:**

```go
// pgUUIDPtrToString converte pgtype.UUID (nullable) para string
func pgUUIDPtrToString(u pgtype.UUID) string {
    if !u.Valid {
        return ""
    }
    return uuid.UUID(u.Bytes).String()
}

// uuidStrPtrToPgtype converte string (possivelmente vazia) para pgtype.UUID nullable
func uuidStrPtrToPgtype(s string) pgtype.UUID {
    if s == "" {
        return pgtype.UUID{Valid: false}
    }
    var pguuid pgtype.UUID
    _ = pguuid.Scan(s)
    return pguuid
}
```

**Uso:**
- `uuidStrPtrToPgtype("")` → `pgtype.UUID{Valid: false}` (NULL no DB)
- `pgUUIDPtrToString(pgtype.UUID{Valid: false})` → `""` (string vazia)

---

## 🔧 Como Aplicar as Mudanças

### Passo 1: Aplicar Migration

```bash
cd backend
make migrate-up
```

**Ou manualmente:**
```bash
psql -U postgres -d nexo_dev -f migrations/032_add_command_id_to_appointments.up.sql
```

**Verificar:**
```sql
SELECT column_name, data_type, is_nullable
FROM information_schema.columns
WHERE table_name = 'appointments' AND column_name = 'command_id';
```

Deve retornar:
```
 column_name | data_type |  is_nullable
-------------+-----------+--------------
 command_id  | uuid      | YES
```

---

### Passo 2: Regenerar SQLC

```bash
cd backend
make sqlc-generate
```

**Ou:**
```bash
sqlc generate
```

**Verificar:**
Arquivos gerados em `backend/internal/infra/db/sqlc/` devem incluir campo `CommandID` em:
- `GetAppointmentByIDRow`
- `ListAppointmentsRow`
- `CreateAppointmentParams`
- `UpdateAppointmentParams`

---

### Passo 3: Compilar Backend

```bash
cd backend
go build ./...
```

**Verificar erros de compilação:**
- ✅ Se nenhum erro → tudo certo
- ❌ Se erros → verificar se SQLC foi regenerado corretamente

---

### Passo 4: Rodar Testes

```bash
cd backend
make test
```

**Ou:**
```bash
go test ./...
```

---

## 🧪 Como Testar

### Teste 1: Verificar Endpoint GET /api/v1/appointments/:id

```bash
# Login
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@teste.com","password":"senha123"}' \
  | jq -r '.token')

# Buscar appointment
curl -X GET http://localhost:8080/api/v1/appointments/<UUID> \
  -H "Authorization: Bearer $TOKEN" \
  | jq '.'
```

**Response esperada:**
```json
{
  "id": "...",
  "tenant_id": "...",
  "customer_name": "João Silva",
  "status": "AWAITING_PAYMENT",
  "command_id": "abc-123-def-456",  ← DEVE APARECER
  "total_price": "80.00",
  ...
}
```

Se `command_id` for `null` ou não existir, deve retornar:
```json
{
  "command_id": null
}
```

Ou simplesmente omitir o campo (devido ao `omitempty`).

---

### Teste 2: Verificar Listagem

```bash
curl -X GET http://localhost:8080/api/v1/appointments \
  -H "Authorization: Bearer $TOKEN" \
  | jq '.data[] | {id, status, command_id}'
```

**Response esperada:**
```json
[
  {
    "id": "...",
    "status": "AWAITING_PAYMENT",
    "command_id": "abc-123..."
  },
  {
    "id": "...",
    "status": "CONFIRMED",
    "command_id": null
  }
]
```

---

### Teste 3: Frontend Routing Inteligente

1. **Reiniciar Next.js:**
```bash
cd frontend
pnpm run dev
```

2. **Cenário:** Agendamento com `status = AWAITING_PAYMENT` e `command_id` preenchido

3. **Ação:** Clicar no agendamento no calendário

4. **Resultado Esperado:**
   - ✅ CommandModal abre
   - ❌ AppointmentModal **NÃO** abre

5. **Verificar no DevTools:**
```javascript
// Console do navegador
appointment.command_id // Deve retornar UUID da comanda
```

---

## 📊 Impacto das Mudanças

### Backend

| Componente | Alterações | Status |
|------------|------------|--------|
| Database | +1 coluna, +1 índice | ✅ Migration OK |
| DTO | +1 campo (JSON) | ✅ Implementado |
| Entity | +1 campo (string) | ✅ Implementado |
| Mapper | +1 mapeamento | ✅ Implementado |
| Repository | +3 conversões | ✅ Implementado |
| SQLC Queries | +2 parâmetros | ✅ Implementado |
| Helpers | +2 funções | ✅ Implementado |

**Total:** 7 arquivos modificados, 2 arquivos criados (migrations)

---

### Frontend

| Componente | Impacto |
|------------|---------|
| TypeScript Types | ✅ Já tem `command_id?: string` |
| handleEventClick | ✅ Já usa `appointment.command_id` |
| AppointmentCard | ✅ Já preparado |
| AgendaCalendar | ✅ Já passa appointment completo |

**Impacto:** ZERO mudanças necessárias no frontend (já estava preparado)

---

## 🔍 Checklist de Verificação

### Backend

- [x] Migration 032 criada (up + down)
- [x] Campo `command_id` adicionado em DTO
- [x] Campo `command_id` adicionado em Entity
- [x] Mapper atualizado para incluir `command_id`
- [x] SQLC queries atualizadas (CreateAppointment, UpdateAppointment)
- [x] Repository: Create atualizado
- [x] Repository: Update atualizado
- [x] Repository: rowToDomain (4 métodos) atualizados
- [x] Helpers: funções UUID nullable criadas
- [ ] Migration aplicada no banco (aguarda execução)
- [ ] SQLC regenerado (aguarda execução)
- [ ] Backend compilado sem erros (aguarda execução)
- [ ] Testes rodados (aguarda execução)

### Frontend

- [x] Type `AppointmentResponse` tem `command_id?`
- [x] handleEventClick usa `appointment.command_id`
- [x] CommandModal integrado
- [x] Menu de contexto implementado

---

## 🚀 Próximos Passos

1. **Aplicar Migration:**
   ```bash
   cd backend
   bash apply-command-id-migration.sh
   ```

2. **Verificar compilação:**
   ```bash
   go build ./...
   ```

3. **Rodar testes:**
   ```bash
   make test
   ```

4. **Testar endpoint:**
   ```bash
   # Buscar appointment e verificar se command_id aparece
   ```

5. **Testar frontend:**
   ```bash
   cd frontend
   pnpm run dev
   # Clicar em appointment AWAITING_PAYMENT → CommandModal deve abrir
   ```

---

## 🎯 Critérios de Sucesso

### ✅ Backend

- [ ] Migration 032 aplicada sem erros
- [ ] SQLC regenerado sem erros
- [ ] Backend compila sem erros
- [ ] Endpoint `GET /api/v1/appointments/:id` retorna `command_id`
- [ ] Endpoint `GET /api/v1/appointments` retorna `command_id` em cada item
- [ ] Campo `command_id` é `null` quando não há comanda vinculada
- [ ] Campo `command_id` contém UUID quando há comanda vinculada

### ✅ Frontend

- [ ] Clicar em appointment AWAITING_PAYMENT (com command_id) → CommandModal abre
- [ ] Clicar em appointment CONFIRMED (sem command_id) → AppointmentModal abre
- [ ] Menu de contexto (botão direito) funciona
- [ ] "Fechar Comanda" no menu abre CommandModal

---

## 📚 Arquivos Criados/Modificados

### Novos Arquivos (2)

| Arquivo | Descrição |
|---------|-----------|
| `backend/migrations/032_add_command_id_to_appointments.up.sql` | Migration para adicionar command_id |
| `backend/migrations/032_add_command_id_to_appointments.down.sql` | Reversão da migration |

### Arquivos Modificados (7)

| Arquivo | Linhas Alteradas | Mudança |
|---------|------------------|---------|
| `backend/internal/application/dto/appointment_dto.go` | +1 | Adicionar CommandID ao DTO |
| `backend/internal/domain/entity/appointment.go` | +1 | Adicionar CommandID à entity |
| `backend/internal/application/mapper/appointment_mapper.go` | +1 | Mapear CommandID |
| `backend/internal/infra/db/queries/appointments.sql` | +3 | Adicionar command_id em queries |
| `backend/internal/infra/repository/postgres/appointment_repository.go` | +8 | Mapear command_id em conversões |
| `backend/internal/infra/repository/postgres/helpers.go` | +16 | Funções UUID nullable |
| `backend/apply-command-id-migration.sh` | NOVO | Script automático |

**Total:** 9 arquivos (2 novos + 7 modificados)

---

## ⚠️ Notas Importantes

1. **SQLC deve ser regenerado** após aplicar a migration
2. **Ordem correta:** Migration → SQLC → Build → Test
3. **Rollback disponível:** Se necessário, rodar `032_add_command_id_to_appointments.down.sql`
4. **Compatibilidade:** Campo é nullable (não quebra dados existentes)
5. **Performance:** Índice criado para otimizar consultas por command_id

---

**Status Final:** ✅ Todas as pendências backend foram resolvidas. Aguarda execução da migration e regeneração do SQLC.
