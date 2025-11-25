# Fluxo de Agendamento — NEXO v1.0

**Versão:** 1.0
**Última Atualização:** 23/11/2025
**Status:** Planejado (v1.0.0 - Milestone 1.5)
**Responsável:** Product + Tech Lead

---

## 📋 Visão Geral

Módulo responsável pelo **agendamento de serviços** de forma visual (calendário estilo AppBarber/Trinks), com controle total de horários, barbeiros, unidades e integração com Google Agenda.

**Prioridade:** 🟡 MÉDIA (Milestone 1.5 - previsto para 10/12/2025)

---

## 🎯 Objetivos do Fluxo

1. ✅ Permitir criação de agendamentos por recepção/gerente
2. ✅ Validar disponibilidade de barbeiro e horário
3. ✅ Impedir conflitos de horário
4. ✅ Sincronizar com Google Agenda (barbeiro)
5. ✅ Notificar cliente via WhatsApp/SMS (futuro)
6. ✅ Controlar status do agendamento (lifecycle)
7. ✅ Respeitar isolamento multi-tenant
8. ✅ Suportar reagendamento e cancelamento

---

## 🔐 Regras de Negócio (RN)

### RN-AGE-001: Validação de Barbeiro

- ❌ Não pode agendar com barbeiro inativo
- ✅ Barbeiro deve pertencer ao mesmo `tenant_id`
- ✅ Barbeiro deve ter horário disponível no slot

### RN-AGE-002: Validação de Cliente

- ✅ Cliente deve existir no sistema antes do agendamento
- ✅ Se não existir, deve ser criado primeiro (tela "Novo Cliente")
- ✅ Cliente deve pertencer ao mesmo `tenant_id`

### RN-AGE-003: Intervalo Padrão

- ✅ Intervalo mínimo entre agendamentos: **10 minutos**
- ✅ Configurável por unidade (futuro)

### RN-AGE-004: Estrutura do Agendamento

Um agendamento sempre pertence a:

- 1 unidade (`unit_id`)
- 1 barbeiro (`professional_id`)
- 1 cliente (`customer_id`)
- 1 ou mais serviços (`services[]`)

### RN-AGE-005: Status de Agendamento

Status permitidos:

- `CREATED` - Criado (pendente confirmação)
- `CONFIRMED` - Confirmado pelo cliente
- `IN_SERVICE` - Em atendimento
- `DONE` - Finalizado com sucesso
- `NO_SHOW` - Cliente faltou
- `CANCELED` - Cancelado (cliente ou barbearia)

### RN-AGE-006: Permissões de Acesso

- **Recepção:** Pode criar, editar, mover, cancelar agendamentos
- **Gerente:** Idem recepção + visualizar todos os barbeiros
- **Barbeiro:** Apenas visualiza sua própria agenda (read-only)
- **Dono:** Acesso total

### RN-AGE-007: Google Agenda Integration

Sincronizar automaticamente:

- ✅ Agendamentos confirmados (status `CONFIRMED`)
- ✅ Cancelamentos (remover do Google Agenda)
- ✅ Alterações de horário (update no Google Agenda)

---

## 📊 Diagrama de Fluxo Principal

```mermaid
flowchart TD
    A[Início] --> B{Usuário autenticado?}
    B -->|Não| C[Redirecionar para Login]
    B -->|Sim| D[Extrair tenant_id do JWT]

    D --> E[Acessar Tela Agendamentos]
    E --> F[Carregar calendário semanal/mensal]
    F --> G[Clicar em Novo Agendamento]

    G --> H[Selecionar Cliente]
    H --> I{Cliente existe?}
    I -->|Não| J[Criar Cliente Primeiro]
    J --> K[Retornar ao Agendamento]
    I -->|Sim| K

    K --> L[Selecionar Serviço(s)]
    L --> M[Calcular duração total]
    M --> N[Selecionar Barbeiro]

    N --> O[Selecionar Data e Horário]
    O --> P{Horário disponível?}
    P -->|Não| Q[Exibir conflito + sugestões]
    Q --> O
    P -->|Sim| R[Validar tenant_id do barbeiro]

    R --> S{Barbeiro pertence ao tenant?}
    S -->|Não| T[Erro 403: Barbeiro inválido]
    S -->|Sim| U[Criar Agendamento no Backend]

    U --> V[POST /api/appointments]
    V --> W{Validação Backend OK?}
    W -->|Não| X[Retornar Erro de Validação]
    W -->|Sim| Y[Persistir no PostgreSQL]

    Y --> Z[Registrar Audit Log]
    Z --> AA[Status inicial: CREATED]
    AA --> AB{Confirmação automática?}
    AB -->|Sim| AC[Atualizar status: CONFIRMED]
    AB -->|Não| AD[Aguardar confirmação manual]

    AC --> AE[Sincronizar Google Agenda]
    AE --> AF[Enviar notificação ao cliente futuro]
    AF --> AG[Atualizar UI React]
    AG --> AH[Fim]

    AD --> AI[Cliente confirma depois]
    AI --> AC
```

````

---

## 🏗️ Arquitetura Técnica

### Backend (Go - Clean Architecture)

**Domain Layer:**

```go
// internal/domain/appointment/appointment.go
type Appointment struct {
    ID             string
    TenantID       string
    UnitID         string
    ProfessionalID string
    CustomerID     string
    ServiceIDs     []string
    StartTime      time.Time
    EndTime        time.Time
    Status         AppointmentStatus
    Notes          string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

type AppointmentStatus string
const (
    StatusCreated    AppointmentStatus = "CREATED"
    StatusConfirmed  AppointmentStatus = "CONFIRMED"
    StatusInService  AppointmentStatus = "IN_SERVICE"
    StatusDone       AppointmentStatus = "DONE"
    StatusNoShow     AppointmentStatus = "NO_SHOW"
    StatusCanceled   AppointmentStatus = "CANCELED"
)
````

**Application Layer:**

```go
// internal/application/usecase/appointment/create_appointment.go
type CreateAppointmentUseCase struct {
    appointmentRepo domain.AppointmentRepository
    professionalRepo domain.ProfessionalRepository
    customerRepo domain.CustomerRepository
    googleCalendar external.GoogleCalendarService
}

func (uc *CreateAppointmentUseCase) Execute(
    ctx context.Context,
    tenantID string,
    req *dto.CreateAppointmentRequest,
) (*dto.CreateAppointmentResponse, error) {
    // 1. Validar tenant do barbeiro
    professional, err := uc.professionalRepo.FindByID(ctx, tenantID, req.ProfessionalID)
    if err != nil {
        return nil, ErrProfessionalNotFound
    }

    // 2. Validar disponibilidade
    conflicts, err := uc.appointmentRepo.CheckConflicts(
        ctx, tenantID, req.ProfessionalID, req.StartTime, req.EndTime,
    )
    if len(conflicts) > 0 {
        return nil, ErrTimeSlotUnavailable
    }

    // 3. Criar agendamento
    appointment := domain.NewAppointment(...)

    // 4. Persistir
    if err := uc.appointmentRepo.Save(ctx, tenantID, appointment); err != nil {
        return nil, err
    }

    // 5. Sincronizar Google Agenda (async)
    go uc.googleCalendar.CreateEvent(appointment)

    return mapper.ToAppointmentResponse(appointment), nil
}
```

**HTTP Handler:**

```go
// internal/infrastructure/http/handler/appointment_handler.go
func (h *AppointmentHandler) Create(c echo.Context) error {
    tenantID := c.Get("tenant_id").(string) // Middleware

    var req dto.CreateAppointmentRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(400, ErrorResponse{Message: "Invalid request"})
    }

    resp, err := h.createUC.Execute(c.Request().Context(), tenantID, &req)
    if err != nil {
        return handleError(c, err)
    }

    return c.JSON(201, resp)
}
```

### Frontend (Next.js + React Query)

**Service:**

```typescript
// frontend/app/lib/services/appointmentService.ts
export const appointmentService = {
  create: async (data: CreateAppointmentDTO) => {
    const response = await apiClient.post('/api/appointments', data);
    return CreateAppointmentResponseSchema.parse(response.data);
  },

  checkAvailability: async (professionalId: string, date: string) => {
    const response = await apiClient.get(`/api/appointments/availability`, {
      params: { professional_id: professionalId, date },
    });
    return AvailabilityResponseSchema.parse(response.data);
  },
};
```

**Hook:**

```typescript
// frontend/app/hooks/useAppointments.ts
export function useCreateAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: appointmentService.create,
    onSuccess: () => {
      toast.success('Agendamento criado com sucesso!');
      queryClient.invalidateQueries(['appointments']);
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
}
```

---

## 🔄 Fluxos Alternativos

### Fluxo 2: Reagendamento

```
[Cliente solicita reagendamento]
   ↓
[Recepção acessa agendamento existente]
   ↓
[Clica em Reagendar]
   ↓
[Seleciona nova data/hora]
   ↓
[Validar disponibilidade]
   ↓
[PUT /api/appointments/:id]
   ↓
[Atualizar no Google Agenda]
   ↓
[Notificar cliente]
   ↓
[Registrar histórico de mudanças]
   ↓
[Fim]
```

### Fluxo 3: Cancelamento

```
[Usuário clica em Cancelar]
   ↓
[Confirmar ação: Sim/Não]
   ↓
[DELETE /api/appointments/:id OU PUT status=CANCELED]
   ↓
[Remover do Google Agenda]
   ↓
[Notificar cliente]
   ↓
[Liberar horário do barbeiro]
   ↓
[Registrar motivo cancelamento (opcional)]
   ↓
[Fim]
```

### Fluxo 4: Marcar No-Show

```
[Barbeiro/Recepção marca cliente faltou]
   ↓
[PUT /api/appointments/:id/no-show]
   ↓
[Atualizar status: NO_SHOW]
   ↓
[Incrementar contador no-show do cliente]
   ↓
[Liberar horário]
   ↓
[Notificar gerente (se configurado)]
   ↓
[Fim]
```

---

## 🗄️ Modelo de Dados

### Tabela: `appointments`

```sql
CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE SET NULL,
    professional_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'CREATED',
    notes TEXT,
    google_event_id VARCHAR(255), -- ID do evento no Google Calendar
    canceled_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_end_after_start CHECK (end_time > start_time)
);

CREATE INDEX idx_appointments_tenant_professional_date
  ON appointments(tenant_id, professional_id, start_time DESC);

CREATE INDEX idx_appointments_tenant_customer
  ON appointments(tenant_id, customer_id);

CREATE INDEX idx_appointments_tenant_status
  ON appointments(tenant_id, status);
```

### Tabela: `appointment_services` (Many-to-Many)

```sql
CREATE TABLE appointment_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL REFERENCES appointments(id) ON DELETE CASCADE,
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE RESTRICT,
    price DECIMAL(10,2) NOT NULL,
    duration_minutes INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(appointment_id, service_id)
);
```

---

## 📡 Endpoints da API

### POST `/api/appointments`

**Descrição:** Criar novo agendamento
**Auth:** JWT (recepção/gerente/dono)
**Body:**

```json
{
  "professional_id": "uuid",
  "customer_id": "uuid",
  "service_ids": ["uuid1", "uuid2"],
  "start_time": "2024-12-05T14:00:00Z",
  "notes": "Cliente prefere corte mais curto"
}
```

**Response:** `201 Created`

### GET `/api/appointments`

**Descrição:** Listar agendamentos (com filtros)
**Query Params:**

- `professional_id` (opcional)
- `customer_id` (opcional)
- `date_from` (opcional)
- `date_to` (opcional)
- `status` (opcional)

### GET `/api/appointments/:id`

**Descrição:** Detalhes de um agendamento

### PUT `/api/appointments/:id`

**Descrição:** Atualizar agendamento (reagendar)

### DELETE `/api/appointments/:id`

**Descrição:** Cancelar agendamento

### PUT `/api/appointments/:id/no-show`

**Descrição:** Marcar cliente como faltante

### GET `/api/appointments/availability`

**Descrição:** Verificar slots disponíveis
**Query:** `professional_id`, `date`

---

## 🔗 Integrações

### Google Calendar API

**Objetivo:** Sincronizar agendamentos com Google Agenda do barbeiro

**Fluxo:**

1. Barbeiro conecta conta Google (OAuth 2.0)
2. Sistema armazena `refresh_token` do barbeiro
3. A cada agendamento CONFIRMED:
   - Backend chama Google Calendar API
   - Cria evento com título, horário, cliente
   - Armazena `google_event_id` no agendamento
4. Em alterações:
   - Update do evento via `google_event_id`
5. Em cancelamentos:
   - Delete do evento

**Configuração:**

```go
// internal/infrastructure/external/google/calendar.go
type GoogleCalendarService struct {
    client *calendar.Service
}

func (g *GoogleCalendarService) CreateEvent(
    appointment *domain.Appointment,
    accessToken string,
) (string, error) {
    event := &calendar.Event{
        Summary: "Agendamento - " + appointment.CustomerName,
        Start:   &calendar.EventDateTime{DateTime: appointment.StartTime.Format(time.RFC3339)},
        End:     &calendar.EventDateTime{DateTime: appointment.EndTime.Format(time.RFC3339)},
    }

    result, err := g.client.Events.Insert("primary", event).Do()
    return result.Id, err
}
```

---

## ✅ Critérios de Aceite

Para considerar o módulo **PRONTO** na v1.0:

- [ ] ✅ Backend implementado (Domain + Use Cases + Handlers)
- [ ] ✅ Frontend com calendário visual (FullCalendar ou similar)
- [ ] ✅ Validação de conflitos de horário funcionando
- [ ] ✅ Multi-tenant isolamento garantido
- [ ] ✅ Integração Google Agenda ativa
- [ ] ✅ RBAC respeitado (barbeiro read-only na própria agenda)
- [ ] ✅ Testes E2E cobrindo fluxo principal
- [ ] ✅ Drag & drop para reagendamento (UX)
- [ ] ✅ Notificações de confirmação (WhatsApp/SMS ou email)

---

## 📊 Métricas de Sucesso

**Operacionais:**

- Taxa de no-show < 10%
- Tempo médio para criar agendamento < 30 segundos
- Conflitos de horário zerados

**Técnicas:**

- Latência API < 150ms
- Sincronização Google < 500ms
- Uptime > 99.5%

---

## 📚 Referências

- `docs/02-arquitetura/ARQUITETURA.md` - Clean Architecture
- `docs/02-arquitetura/MODELO_MULTI_TENANT.md` - Isolamento de dados
- `docs/04-backend/GUIA_DEV_BACKEND.md` - Padrões Go
- `docs/03-frontend/GUIA_FRONTEND.md` - Padrões React/Next.js
- `PRD-NEXO.md` - Requisitos de produto

---

**Status:** 🟡 Planejado
**Próximo Marco:** Milestone 1.5 (10/12/2025)
**Última Revisão:** 23/11/2025
