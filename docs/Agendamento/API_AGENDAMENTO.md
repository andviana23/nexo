# API Contract — Módulo de Agendamento | NEXO v1.0

**Versão:** 1.0.0  
**Data:** 25/11/2025  
**Base URL:** `/api/v1`  
**Auth:** Bearer Token (JWT)  
**Content-Type:** `application/json`  

---

## 1. Visão Geral dos Endpoints

| Método | Endpoint | Descrição | Permissões |
|--------|----------|-----------|------------|
| `POST` | `/appointments` | Criar novo agendamento | Owner, Manager, Receptionist, **Employee (apenas para si)** |
| `GET` | `/appointments` | Listar agendamentos (filtros) | Owner, Manager, Receptionist, **Employee (somente próprios)** |
| `GET` | `/appointments/:id` | Detalhes do agendamento | Owner, Manager, Receptionist, **Employee (somente próprios)** |
| `PATCH` | `/appointments/:id/reschedule` | **Reagendar** (mudar data/horário) | Owner, Manager, Receptionist, **Employee (somente próprios, sem trocar profissional)** |
| `PUT` | `/appointments/:id` | Atualizar agendamento (serviços/notas) | Owner, Manager, Receptionist, **Employee (somente próprios)** |
| `DELETE` | `/appointments/:id` | Cancelar agendamento | Owner, Manager, Receptionist, **Employee (somente próprios)** |
| `GET` | `/appointments/availability` | Verificar disponibilidade | Todos (Employee vê apenas disponibilidade do próprio profissional) |
| `PATCH` | `/appointments/:id/status` | Alterar status (ex: No-Show) | Owner, Manager, Receptionist, **Employee (somente próprios)** |

---

## 2. Detalhamento dos Endpoints

### 2.1 Criar Agendamento

**Endpoint:** `POST /appointments`  
**Descrição:** Cria um novo agendamento validando conflitos e disponibilidade.

#### Request Body

```json
{
  "professional_id": "uuid",
  "customer_id": "uuid",
  "service_ids": ["uuid", "uuid"],
  "start_time": "2025-12-05T14:00:00Z",
  "notes": "Cliente prefere corte na tesoura"
}
```

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `professional_id` | UUID | Sim | ID do barbeiro |
| `customer_id` | UUID | Sim | ID do cliente |
| `service_ids` | Array<UUID> | Sim | Lista de serviços (mínimo 1) |
| `start_time` | ISO8601 | Sim | Data e hora de início (UTC) |
| `notes` | String | Não | Observações internas |

**Nota de RBAC:** usuários com role Employee/Barbeiro só podem criar agendamentos vinculando o próprio `professional_id`; caso contrário, recebem `403 Forbidden`.

#### Response (201 Created)

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "tenant_id": "uuid",
    "professional": {
      "id": "uuid",
      "name": "João Barbeiro"
    },
    "customer": {
      "id": "uuid",
      "name": "Carlos Cliente"
    },
    "services": [
      {
        "id": "uuid",
        "name": "Corte Masculino",
        "price": 50.00,
        "duration": 30
      }
    ],
    "start_time": "2025-12-05T14:00:00Z",
    "end_time": "2025-12-05T14:30:00Z",
    "status": "CREATED",
    "total_price": 50.00,
    "created_at": "2025-11-25T10:00:00Z"
  }
}
```

#### Erros Possíveis

| Código | Status | Mensagem | Solução |
|--------|--------|----------|---------|
| `TIME_SLOT_CONFLICT` | 409 | Horário já está ocupado | Escolher outro horário |
| `BLOCKED_TIME` | 409 | Horário está bloqueado para o profissional | Escolher outro horário |
| `INSUFFICIENT_INTERVAL` | 409 | Intervalo mínimo não respeitado | Respeitar 10min entre agendamentos |
| `PROFESSIONAL_NOT_FOUND` | 404 | Barbeiro não encontrado | Verificar profissional |
| `CUSTOMER_NOT_FOUND` | 404 | Cliente não encontrado | Cadastrar cliente antes |
| `SERVICE_NOT_FOUND` | 404 | Serviço não encontrado | Revisar lista de serviços |
| `FORBIDDEN_SCOPE` | 403 | Barbeiro tentando criar para outro profissional | Barbeiro só cria para si |

---

### 2.2 Listar Agendamentos

**Endpoint:** `GET /appointments`  
**Descrição:** Lista agendamentos com filtros.

#### Query Parameters

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `page` | Int | Não | Página atual (default: 1) |
| `page_size` | Int | Não | Itens por página (default: 20, máx: 100) |
| `professional_id` | UUID | Não | Filtrar por barbeiro |
| `customer_id` | UUID | Não | Filtrar por cliente |
| `start_date` | String | Não | Data inicial. **Aceita:** `YYYY-MM-DD` ou `ISO8601` (ex: `2025-12-01` ou `2025-12-01T00:00:00Z`) |
| `end_date` | String | Não | Data final. **Aceita:** `YYYY-MM-DD` ou `ISO8601` |
| `status` | String/Array | Não | Filtrar por status. **Aceita:** string única ou array via query params repetidos (ex: `?status=CREATED&status=CONFIRMED`) |

**Nota de RBAC:** se autenticado como Employee/Barbeiro, a listagem sempre retorna apenas agendamentos associados ao próprio `professional_id`, mesmo que outros IDs sejam passados nos filtros.

**Status válidos:** `CREATED`, `CONFIRMED`, `CHECKED_IN`, `IN_SERVICE`, `AWAITING_PAYMENT`, `DONE`, `NO_SHOW`, `CANCELED`

#### Exemplos de Filtros

```bash
# Filtro por data (YYYY-MM-DD)
GET /appointments?start_date=2025-12-01&end_date=2025-12-01

# Filtro por status único
GET /appointments?status=AWAITING_PAYMENT

# Filtro por múltiplos status
GET /appointments?status=CREATED&status=CONFIRMED

# Combinação de filtros
GET /appointments?start_date=2025-12-01&end_date=2025-12-07&status=AWAITING_PAYMENT&professional_id=uuid
```

#### Response (200 OK)

```json
{
  "data": [
    {
      "id": "uuid",
      "professional_name": "João Barbeiro",
      "customer_name": "Carlos Cliente",
      "start_time": "2025-12-05T14:00:00Z",
      "end_time": "2025-12-05T14:30:00Z",
      "status": "CONFIRMED",
      "services": [
        {
          "service_id": "uuid",
          "service_name": "Corte",
          "price": "50.00",
          "duration": 30
        }
      ]
    }
  ],
  "page": 1,
  "page_size": 20,
  "total": 150
}
```

---

### 2.3 Reagendar Agendamento

**Endpoint:** `PATCH /appointments/:id/reschedule`  
**Descrição:** Reagenda um agendamento existente para nova data/horário e/ou novo profissional.

#### Request Body

```json
{
  "new_start_time": "2025-12-06T15:00:00Z",
  "professional_id": "uuid-opcional"
}
```

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `new_start_time` | ISO8601 | Sim | Nova data e hora de início (UTC) |
| `professional_id` | UUID | Não | Novo profissional (opcional, mantém o mesmo se omitido) |

#### Response (200 OK)

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "tenant_id": "uuid",
  "professional": {
    "id": "uuid",
    "name": "João Barbeiro"
  },
  "customer": {
    "id": "uuid",
    "name": "Carlos Cliente"
  },
  "services": [
    {
      "id": "uuid",
      "name": "Corte Masculino",
      "price": 50.00,
      "duration": 30
    }
  ],
  "start_time": "2025-12-06T15:00:00Z",
  "end_time": "2025-12-06T15:30:00Z",
  "status": "CREATED",
  "total_price": 50.00,
  "updated_at": "2025-11-25T10:30:00Z"
}
```

#### Regras de Negócio
- ✅ Mantém a duração original do agendamento (recalcula `end_time` automaticamente).
- ✅ Valida conflitos novamente no novo horário.
- ✅ Permite trocar de profissional durante o reagendamento.
- ✅ Verifica bloqueios de horário do profissional.
- ✅ Respeita intervalo mínimo de 10 minutos entre agendamentos.
- ✅ Barbeiro só reage agenda os próprios atendimentos e não pode trocar o profissional.

#### Erros Possíveis

| Código | Status | Mensagem | Solução |
|--------|--------|----------|---------|
| `TIME_SLOT_CONFLICT` | 409 | Novo horário já está ocupado | Escolher outro horário |
| `BLOCKED_TIME` | 409 | Horário bloqueado para o profissional | Escolher outro horário |
| `INSUFFICIENT_INTERVAL` | 409 | Intervalo mínimo não respeitado | Respeitar 10min entre agendamentos |
| `INVALID_TRANSITION` | 409 | Agendamento não pode ser reagendado neste status | Ajustar fluxo antes de reagendar |
| `PROFESSIONAL_NOT_FOUND` | 404 | Novo barbeiro não existe | Escolher barbeiro válido |
| `APPOINTMENT_NOT_FOUND` | 404 | Agendamento não existe | Verificar ID |
| `FORBIDDEN_SCOPE` | 403 | Barbeiro tentou mover para outro profissional | Barbeiro só reage agenda para si |

---

### 2.4 Atualizar Agendamento (Geral)

**Endpoint:** `PUT /appointments/:id`  
**Descrição:** Atualiza dados gerais do agendamento (serviços, notas).  
**⚠️ Nota:** Para reagendamento de data/horário, use `PATCH /appointments/:id/reschedule` ao invés deste endpoint.

#### Request Body

```json
{
  "service_ids": ["uuid1", "uuid2"],
  "notes": "Alterado a pedido do cliente"
}
```

#### Regras de Negócio
- ✅ Recalcula `end_time` se serviços mudarem.
- ✅ Atualiza Google Agenda se status for `CONFIRMED`.
- ❌ **Não permite alterar** `start_time` diretamente (use `/reschedule`).

---

### 2.5 Cancelar Agendamento

**Endpoint:** `DELETE /appointments/:id`  
**Descrição:** Cancela um agendamento (Soft Delete ou Status Update).

#### Request Body (Opcional)

```json
{
  "reason": "Cliente desistiu"
}
```

#### Response (200 OK)

```json
{
  "message": "Agendamento cancelado com sucesso",
  "id": "uuid",
  "status": "CANCELED"
}
```

#### Regras e Erros
- ✅ Barbeiro só pode cancelar agendamentos associados ao próprio `professional_id`.
- ✅ Cancelamento em status não permitido retorna `409 Conflict` (`INVALID_TRANSITION`).
- ✅ Agendamento inexistente retorna `404 Not Found`.

---

### 2.5 Verificar Disponibilidade

**Endpoint:** `GET /appointments/availability`  
**Descrição:** Retorna slots disponíveis para um barbeiro em uma data.

#### Query Parameters

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `professional_id` | UUID | Sim | ID do barbeiro |
| `date` | Date | Sim | Data (YYYY-MM-DD) |

**Nota de RBAC:** quando autenticado como Employee/Barbeiro, a API força `professional_id` para o ID do próprio barbeiro; consultas de disponibilidade de outros profissionais retornam `403 Forbidden`.

#### Response (200 OK)

```json
{
  "data": [
    { "time": "08:00", "available": true },
    { "time": "08:15", "available": true },
    { "time": "08:30", "available": false, "reason": "BOOKED" },
    { "time": "08:45", "available": false, "reason": "BLOCKED" }
  ]
}
```

---

### 2.6 Alterar Status

**Endpoint:** `PUT /appointments/:id/status`  
**Descrição:** Transição de status específica (ex: Marcar No-Show).

#### Request Body

```json
{
  "status": "NO_SHOW",
  "reason": "Cliente não compareceu e não avisou"
}
```

#### Response (200 OK)

```json
{
  "id": "uuid",
  "previous_status": "CONFIRMED",
  "new_status": "NO_SHOW",
  "updated_at": "2025-12-05T15:00:00Z"
}
```

#### Regras e Erros
- ✅ Barbeiro só pode alterar status dos próprios agendamentos.
- ✅ Transições inválidas retornam `409 Conflict` (`INVALID_TRANSITION`).
- ✅ Agendamento inexistente retorna `404 Not Found`.
- ✅ Tentativa de acessar agendamento de outro profissional (barbeiro) retorna `403 Forbidden`.

---

## 3. Códigos de Erro Padrão

| Código | HTTP Status | Descrição |
|--------|-------------|-----------|
| `INVALID_REQUEST` | 400 | Payload inválido (validação Zod falhou) |
| `UNAUTHORIZED` | 401 | Token ausente ou inválido |
| `FORBIDDEN` | 403 | Sem permissão para esta ação (RBAC) |
| `NOT_FOUND` | 404 | Recurso não encontrado |
| `CONFLICT` | 409 | Conflito de regra de negócio (horário ocupado/bloqueado, intervalo mínimo) |
| `INVALID_TRANSITION` | 409 | Transição de status/fluxo não permitida |
| `INTERNAL_ERROR` | 500 | Erro inesperado no servidor |

---

## 4. Exemplo de Fluxo Completo (Frontend → Backend)

### Cenário: Criar Agendamento

1. **Frontend** envia `POST /appointments`
   ```json
   {
     "professional_id": "123",
     "customer_id": "456",
     "service_ids": ["789"],
     "start_time": "2025-12-10T10:00:00Z"
   }
   ```

2. **Backend** valida token JWT.
3. **Backend** extrai `tenant_id` do token.
4. **Backend** verifica se `professional_id` pertence ao tenant.
5. **Backend** verifica conflitos no banco:
   ```sql
   SELECT count(*) FROM appointments 
   WHERE professional_id = '123' 
   AND start_time < '2025-12-10T10:30:00Z' 
   AND end_time > '2025-12-10T10:00:00Z'
   AND status != 'CANCELED'
   ```
6. **Backend** insere no banco.
7. **Backend** retorna `201 Created`.

---

## 5. Formato de Valores Monetários

### 5.1 Convenção de Preços na API

Todos os valores monetários (`total_price`, `price` de serviços) são retornados pela API como **strings numéricas** no formato americano (ponto decimal).

**Formato:** `"XX.XX"` (string numérica com 2 casas decimais)

| Campo | Tipo | Exemplo | Descrição |
|-------|------|---------|-----------|
| `total_price` | String | `"50.00"` | Valor total do agendamento |
| `services[].price` | String | `"35.50"` | Preço unitário do serviço |

### 5.2 Exemplos

```json
{
  "id": "uuid",
  "total_price": "85.50",
  "services": [
    {
      "service_id": "uuid",
      "service_name": "Corte Masculino",
      "price": "50.00",
      "duration": 30
    },
    {
      "service_id": "uuid",
      "service_name": "Barba",
      "price": "35.50",
      "duration": 20
    }
  ]
}
```

### 5.3 Formatação no Frontend

O frontend deve formatar os valores para exibição usando a função `formatCurrency()`:

```typescript
import { formatCurrency } from '@/types/appointment';

// Uso
formatCurrency("50.00")   // "R$ 50,00"
formatCurrency(50)        // "R$ 50,00"
formatCurrency("85.50")   // "R$ 85,50"
```

### 5.4 Histórico de Alterações

| Data | Versão | Alteração |
|------|--------|-----------|
| 01/12/2025 | 1.0.1 | **BUG-004:** Alterado formato de `total_price` e `price` de string formatada (`"R$ 50,00"`) para string numérica (`"50.00"`) para evitar NaN no frontend |

---

## 6. Interface Visual — FullCalendar Scheduler

### 6.1 Licença FullCalendar Premium (Modo Avaliação)

A interface visual do calendário utiliza o **FullCalendar Scheduler (Premium)**, que está em período de avaliação gratuita.

**Chave de Licença (Desenvolvimento):**

```javascript
var calendar = new Calendar(calendarEl, {
  schedulerLicenseKey: 'CC-Attribution-NonCommercial-NoDerivatives'
});
```

### ⚠️ Atenção Legal

- ❌ **Proibido uso comercial** desta licença.
- ✅ **Permitido apenas** para:
  - Desenvolvimento interno
  - Testes de integração
  - Homologação
  - Demonstrações internas
- ⚠️ **Antes do lançamento em produção**, será necessário:
  - Adquirir a **licença comercial oficial** do FullCalendar.
  - Substituir a chave de avaliação pela chave comercial.
  - Validar conformidade legal.

**Documentação oficial:** [FullCalendar Pricing](https://fullcalendar.io/pricing)

---

**Responsável:** Tech Lead  
**Data:** 25/11/2025  
**Atualizado:** 01/12/2025 (BUG-004)  
**Status:** ✅ COMPLETO
