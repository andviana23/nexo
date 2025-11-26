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
| `POST` | `/appointments` | Criar novo agendamento | Owner, Manager, Receptionist |
| `GET` | `/appointments` | Listar agendamentos (filtros) | Todos (Barbeiro vê próprios) |
| `GET` | `/appointments/:id` | Detalhes do agendamento | Todos (Barbeiro vê próprios) |
| `PUT` | `/appointments/:id` | Atualizar agendamento | Owner, Manager, Receptionist |
| `DELETE` | `/appointments/:id` | Cancelar agendamento | Owner, Manager, Receptionist |
| `GET` | `/appointments/availability` | Verificar disponibilidade | Todos |
| `PUT` | `/appointments/:id/status` | Alterar status (ex: No-Show) | Owner, Manager, Receptionist |

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
| `PROFESSIONAL_INACTIVE` | 400 | Barbeiro inativo | Escolher outro barbeiro |
| `CUSTOMER_NOT_FOUND` | 404 | Cliente não encontrado | Cadastrar cliente antes |
| `INSUFFICIENT_INTERVAL` | 400 | Intervalo mínimo não respeitado | Respeitar 10min entre agendamentos |

---

### 2.2 Listar Agendamentos

**Endpoint:** `GET /appointments`  
**Descrição:** Lista agendamentos com filtros.

#### Query Parameters

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `page` | Int | Não | Página atual (default: 1) |
| `page_size` | Int | Não | Itens por página (default: 20) |
| `professional_id` | UUID | Não | Filtrar por barbeiro |
| `customer_id` | UUID | Não | Filtrar por cliente |
| `date_from` | ISO8601 | Não | Data inicial (ex: 2025-12-01T00:00:00Z) |
| `date_to` | ISO8601 | Não | Data final |
| `status` | String | Não | Filtrar por status (ex: CONFIRMED) |

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
      "service_names": ["Corte", "Barba"]
    }
  ],
  "meta": {
    "page": 1,
    "page_size": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

---

### 2.3 Atualizar Agendamento

**Endpoint:** `PUT /appointments/:id`  
**Descrição:** Atualiza dados do agendamento (reagendamento, troca de serviço).

#### Request Body

```json
{
  "start_time": "2025-12-06T15:00:00Z",
  "service_ids": ["uuid1", "uuid2"],
  "notes": "Alterado a pedido do cliente"
}
```

#### Regras de Negócio
- ✅ Valida conflitos novamente se horário mudar.
- ✅ Recalcula `end_time` se serviços mudarem.
- ✅ Atualiza Google Agenda se status for `CONFIRMED`.

---

### 2.4 Cancelar Agendamento

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

---

### 2.5 Verificar Disponibilidade

**Endpoint:** `GET /appointments/availability`  
**Descrição:** Retorna slots disponíveis para um barbeiro em uma data.

#### Query Parameters

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `professional_id` | UUID | Sim | ID do barbeiro |
| `date` | Date | Sim | Data (YYYY-MM-DD) |

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

---

## 3. Códigos de Erro Padrão

| Código | HTTP Status | Descrição |
|--------|-------------|-----------|
| `INVALID_REQUEST` | 400 | Payload inválido (validação Zod falhou) |
| `UNAUTHORIZED` | 401 | Token ausente ou inválido |
| `FORBIDDEN` | 403 | Sem permissão para esta ação (RBAC) |
| `NOT_FOUND` | 404 | Recurso não encontrado |
| `CONFLICT` | 409 | Conflito de regra de negócio (horário ocupado) |
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

## 5. Interface Visual — FullCalendar Scheduler

### 5.1 Licença FullCalendar Premium (Modo Avaliação)

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
**Status:** ✅ COMPLETO
