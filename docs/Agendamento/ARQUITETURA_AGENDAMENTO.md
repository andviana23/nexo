# Arquitetura — Módulo de Agendamento | NEXO v1.0

**Versão:** 1.0.0  
**Data:** 25/11/2025  
**Responsável:** Tech Lead  

---

## 1. Visão Arquitetural

### 1.1 Diagrama de Componentes

```
┌─────────────────────────────────────────────────────────────┐
│                     FRONTEND (Next.js 15)                    │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │  Calendar    │  │  Form Modal  │  │  Service        │   │
│  │  Component   │  │  (Create/    │  │  Layer          │   │
│  │ (FullCalendar│  │   Edit)      │  │ (API Calls)     │   │
│  └──────────────┘  └──────────────┘  └─────────────────┘   │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │        TanStack Query (State Management)             │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ HTTP/REST (JSON)
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                     BACKEND (Go + Echo)                      │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────┐   │
│  │            HTTP HANDLERS (Presentation)               │   │
│  │  - AppointmentHandler                                │   │
│  │  - Validation, RBAC, Error Handling                  │   │
│  └──────────────────────────────────────────────────────┘   │
│                            │                                  │
│                            ▼                                  │
│  ┌──────────────────────────────────────────────────────┐   │
│  │          APPLICATION LAYER (Use Cases)                │   │
│  │  - CreateAppointmentUseCase                          │   │
│  │  - UpdateAppointmentUseCase                          │   │
│  │  - CancelAppointmentUseCase                          │   │
│  │  - CheckAvailabilityUseCase                          │   │
│  │  - ListAppointmentsUseCase                           │   │
│  └──────────────────────────────────────────────────────┘   │
│                            │                                  │
│                            ▼                                  │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              DOMAIN LAYER (Business Logic)            │   │
│  │  - Appointment Entity                                │   │
│  │  - AppointmentRepository Interface                   │   │
│  │  - Value Objects (Status, TimeSlot)                  │   │
│  │  - Business Rules                                    │   │
│  └──────────────────────────────────────────────────────┘   │
│                            │                                  │
│                            ▼                                  │
│  ┌──────────────────────────────────────────────────────┐   │
│  │         INFRASTRUCTURE (Repositories, External)       │   │
│  │  - PostgresAppointmentRepository                     │   │
│  │  - GoogleCalendarService                             │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ SQL (sqlc)
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   POSTGRESQL (Neon Cloud)                    │
│  - appointments                                              │
│  - appointment_services                                      │
│  - profissionais                                             │
│  - clientes                                                  │
│  - servicos                                                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 2. Stack Tecnológica

### 2.1 Frontend

- **Framework:** Next.js 15.5.6 (App Router)
- **React:** 19.2.0
- **UI:** Tailwind CSS 4.1.17 + shadcn/ui
- **Calendário:** FullCalendar 6.x (ResourceTimeGrid)
- **State Management:** TanStack Query 5.90.11
- **Forms:** React Hook Form 7.66.1 + Zod 4.1.13

### ⚠️ Atenção Legal: Licença FullCalendar Premium

Durante o período de avaliação gratuita do **FullCalendar Premium (Scheduler)**, o NEXO está autorizado a utilizar a licença não-comercial fornecida pelo próprio FullCalendar, **exclusivamente para fins de desenvolvimento interno**, sem uso comercial e sem disponibilização aos clientes finais.

**Chave de Licença (Modo Avaliação):**

```javascript
var calendar = new Calendar(calendarEl, {
  schedulerLicenseKey: 'CC-Attribution-NonCommercial-NoDerivatives'
});
```

**Restrições:**

- ❌ **Proibido uso comercial** desta licença.
- ✅ **Permitido apenas** para testes internos, homologação e desenvolvimento.
- ⚠️ A versão final do NEXO que será usada por barbearias **exigirá a compra da licença oficial**.
- 🔄 **Substituir a chave** de desenvolvimento pela licença comercial **antes do lançamento em produção**.

**Dependências Externas Críticas:**

| Dependência | Versão | Status Licença | Próxima Ação |
|-------------|--------|----------------|---------------|
| FullCalendar Scheduler | 6.x | Modo Avaliação (CC-Attribution-NonCommercial-NoDerivatives) | Comprar licença comercial antes da produção |

### 2.2 Backend

- **Linguagem:** Go 1.24
- **Framework:** Echo v4
- **Arquitetura:** Clean Architecture + DDD
- **Banco:** PostgreSQL 14+ (Neon Cloud)
- **ORM:** sqlc v1.30.0

---

## 3. Frontend Architecture

### 3.1 Estrutura de Pastas

```
frontend/src/
├── app/
│   └── (dashboard)/
│       └── agendamentos/
│           ├── page.tsx                    # Calendário principal
│           ├── loading.tsx                 # Loading state
│           ├── error.tsx                   # Error boundary
│           └── novo/
│               └── page.tsx                # Form de criação
│
├── components/
│   └── appointments/
│       ├── AppointmentCalendar.tsx         # Wrapper FullCalendar
│       ├── AppointmentForm.tsx             # Formulário (Zod + RHF)
│       ├── AppointmentModal.tsx            # Modal de detalhes
│       ├── AppointmentCard.tsx             # Card mobile
│       ├── ProfessionalSelect.tsx          # Select barbeiro
│       ├── ServiceMultiSelect.tsx          # Multi-select serviços
│       ├── ClientAutocomplete.tsx          # Autocomplete cliente
│       └── StatusBadge.tsx                 # Badge de status
│
├── services/
│   └── appointment-service.ts              # API client
│       ├── create()
│       ├── update()
│       ├── cancel()
│       ├── list()
│       ├── checkAvailability()
│       └── getById()
│
├── hooks/
│   └── use-appointments.ts                 # React Query hooks
│       ├── useAppointments()               # List (GET)
│       ├── useAppointment()                # Get by ID
│       ├── useCreateAppointment()          # Create (POST)
│       ├── useUpdateAppointment()          # Update (PUT)
│       ├── useCancelAppointment()          # Cancel (DELETE)
│       └── useCheckAvailability()          # Check conflicts
│
└── types/
    └── appointment.ts                      # TypeScript interfaces
        ├── Appointment
        ├── AppointmentStatus
        ├── CreateAppointmentDTO
        ├── UpdateAppointmentDTO
        └── AvailabilityResponse
```

### 2.2 Componentes Principais

#### AppointmentCalendar.tsx

```typescript
'use client';

import FullCalendar from '@fullcalendar/react';
import resourceTimeGridPlugin from '@fullcalendar/resource-timegrid';
import interactionPlugin from '@fullcalendar/interaction';
import { useAppointments } from '@/hooks/use-appointments';

export function AppointmentCalendar() {
  const { data: appointments, isLoading } = useAppointments({
    dateFrom: startOfWeek(new Date()),
    dateTo: endOfWeek(new Date()),
  });

  const resources = useProfessionals(); // Barbeiros como "recursos"

  const handleEventClick = (info) => {
    // Abrir modal de detalhes
  };

  const handleDateSelect = (selectInfo) => {
    // Abrir modal de criação
  };

  const handleEventDrop = (info) => {
    // Reagendar via drag & drop
    mutateUpdate.mutate({
      id: info.event.id,
      start_time: info.event.start,
      end_time: info.event.end,
    });
  };

  return (
    <FullCalendar
      plugins={[resourceTimeGridPlugin, interactionPlugin]}
      initialView="resourceTimeGridWeek"
      resources={resources}
      events={appointments}
      editable={hasPermission('edit')}
      selectable
      select={handleDateSelect}
      eventClick={handleEventClick}
      eventDrop={handleEventDrop}
      slotMinTime="08:00:00"
      slotMaxTime="20:00:00"
      slotDuration="00:15:00"
      locale="pt-br"
    />
  );
}
```

#### AppointmentForm.tsx

```typescript
'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

const appointmentSchema = z.object({
  professional_id: z.string().uuid(),
  customer_id: z.string().uuid(),
  service_ids: z.array(z.string().uuid()).min(1),
  start_time: z.date(),
  notes: z.string().optional(),
});

type FormData = z.infer<typeof appointmentSchema>;

export function AppointmentForm({ onSuccess }: Props) {
  const form = useForm<FormData>({
    resolver: zodResolver(appointmentSchema),
  });

  const { mutate, isPending } = useCreateAppointment();

  const onSubmit = (data: FormData) => {
    mutate(data, {
      onSuccess: () => {
        toast.success('Agendamento criado!');
        onSuccess?.();
      },
      onError: (error) => {
        if (error.code === 'TIME_SLOT_CONFLICT') {
          form.setError('start_time', {
            message: 'Horário já ocupado. Escolha outro.',
          });
        }
      },
    });
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <FormField name="customer_id" label="Cliente" required>
          <ClientAutocomplete />
        </FormField>

        <FormField name="service_ids" label="Serviços" required>
          <ServiceMultiSelect />
        </FormField>

        <FormField name="professional_id" label="Barbeiro" required>
          <ProfessionalSelect />
        </FormField>

        <FormField name="start_time" label="Data e Hora" required>
          <DateTimePicker />
        </FormField>

        <Button type="submit" disabled={isPending}>
          {isPending ? 'Criando...' : 'Agendar'}
        </Button>
      </form>
    </Form>
  );
}
```

### 2.3 React Query Hooks

```typescript
// hooks/use-appointments.ts

export const appointmentKeys = {
  all: ['appointments'] as const,
  lists: () => [...appointmentKeys.all, 'list'] as const,
  list: (filters: AppointmentFilters) =>
    [...appointmentKeys.lists(), filters] as const,
  details: () => [...appointmentKeys.all, 'detail'] as const,
  detail: (id: string) => [...appointmentKeys.details(), id] as const,
};

export function useAppointments(filters: AppointmentFilters) {
  return useQuery({
    queryKey: appointmentKeys.list(filters),
    queryFn: () => appointmentService.list(filters),
    staleTime: 30_000, // 30s
  });
}

export function useCreateAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: appointmentService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
    },
  });
}

export function useUpdateAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateAppointmentDTO }) =>
      appointmentService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: appointmentKeys.detail(variables.id),
      });
    },
  });
}

export function useCancelAppointment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: appointmentService.cancel,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
    },
  });
}
```

### 2.4 Service Layer

```typescript
// services/appointment-service.ts

import { apiClient } from '@/lib/axios';
import {
  Appointment,
  CreateAppointmentDTO,
  UpdateAppointmentDTO,
  AppointmentFilters,
} from '@/types/appointment';

export const appointmentService = {
  list: async (filters: AppointmentFilters): Promise<Appointment[]> => {
    const { data } = await apiClient.get('/api/v1/appointments', {
      params: filters,
    });
    return data.data;
  },

  getById: async (id: string): Promise<Appointment> => {
    const { data } = await apiClient.get(`/api/v1/appointments/${id}`);
    return data.data;
  },

  create: async (dto: CreateAppointmentDTO): Promise<Appointment> => {
    const { data } = await apiClient.post('/api/v1/appointments', dto);
    return data.data;
  },

  update: async (
    id: string,
    dto: UpdateAppointmentDTO
  ): Promise<Appointment> => {
    const { data } = await apiClient.put(`/api/v1/appointments/${id}`, dto);
    return data.data;
  },

  cancel: async (id: string): Promise<void> => {
    await apiClient.delete(`/api/v1/appointments/${id}`);
  },

  checkAvailability: async (
    professionalId: string,
    date: string
  ): Promise<AvailabilitySlot[]> => {
    const { data } = await apiClient.get('/api/v1/appointments/availability', {
      params: { professional_id: professionalId, date },
    });
    return data.data;
  },
};
```

---

## 3. Backend Architecture (Go + Clean Architecture)

### 3.1 Estrutura de Pastas

```
backend/
├── internal/
│   ├── domain/
│   │   └── appointment/
│   │       ├── appointment.go          # Entity
│   │       ├── repository.go           # Interface
│   │       ├── status.go               # Value Object
│   │       └── errors.go               # Domain errors
│   │
│   ├── application/
│   │   ├── dto/
│   │   │   └── appointment_dto.go
│   │   ├── mapper/
│   │   │   └── appointment_mapper.go
│   │   └── usecase/
│   │       └── appointment/
│   │           ├── create_appointment.go
│   │           ├── update_appointment.go
│   │           ├── cancel_appointment.go
│   │           ├── list_appointments.go
│   │           └── check_availability.go
│   │
│   └── infrastructure/
│       ├── repository/
│       │   └── postgres/
│       │       ├── appointment_repository.go
│       │       └── queries.sql (sqlc)
│       ├── http/
│       │   └── handler/
│       │       └── appointment_handler.go
│       └── external/
│           └── google/
│               └── calendar_service.go
│
└── migrations/
    ├── 00XX_create_appointments.up.sql
    └── 00XX_create_appointments.down.sql
```

### 3.2 Domain Layer

#### appointment.go (Entity)

```go
package appointment

import (
    "time"
    "github.com/google/uuid"
)

type Appointment struct {
    ID             uuid.UUID
    TenantID       uuid.UUID
    ProfessionalID uuid.UUID
    CustomerID     uuid.UUID
    ServiceIDs     []uuid.UUID
    StartTime      time.Time
    EndTime        time.Time
    Status         Status
    Notes          string
    GoogleEventID  *string
    CanceledReason *string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

// NewAppointment cria novo agendamento com validações
func NewAppointment(
    tenantID uuid.UUID,
    professionalID uuid.UUID,
    customerID uuid.UUID,
    serviceIDs []uuid.UUID,
    startTime time.Time,
    endTime time.Time,
    notes string,
) (*Appointment, error) {
    if err := validateTimeRange(startTime, endTime); err != nil {
        return nil, err
    }

    if len(serviceIDs) == 0 {
        return nil, ErrNoServicesProvided
    }

    return &Appointment{
        ID:             uuid.New(),
        TenantID:       tenantID,
        ProfessionalID: professionalID,
        CustomerID:     customerID,
        ServiceIDs:     serviceIDs,
        StartTime:      startTime,
        EndTime:        endTime,
        Status:         StatusCreated,
        Notes:          notes,
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }, nil
}

// CanTransitionTo valida se pode mudar para novo status
func (a *Appointment) CanTransitionTo(newStatus Status) bool {
    validTransitions := map[Status][]Status{
        StatusCreated:   {StatusConfirmed, StatusCanceled},
        StatusConfirmed: {StatusInService, StatusNoShow, StatusCanceled},
        StatusInService: {StatusDone, StatusCanceled},
        StatusDone:      {},
        StatusNoShow:    {},
        StatusCanceled:  {},
    }

    allowed, ok := validTransitions[a.Status]
    if !ok {
        return false
    }

    for _, s := range allowed {
        if s == newStatus {
            return true
        }
    }
    return false
}

// UpdateStatus atualiza status com validação
func (a *Appointment) UpdateStatus(newStatus Status) error {
    if !a.CanTransitionTo(newStatus) {
        return ErrInvalidStatusTransition
    }

    a.Status = newStatus
    a.UpdatedAt = time.Now()
    return nil
}
```

#### status.go (Value Object)

```go
package appointment

type Status string

const (
    StatusCreated   Status = "CREATED"
    StatusConfirmed Status = "CONFIRMED"
    StatusInService Status = "IN_SERVICE"
    StatusDone      Status = "DONE"
    StatusNoShow    Status = "NO_SHOW"
    StatusCanceled  Status = "CANCELED"
)

func (s Status) String() string {
    return string(s)
}

func (s Status) IsValid() bool {
    switch s {
    case StatusCreated, StatusConfirmed, StatusInService,
        StatusDone, StatusNoShow, StatusCanceled:
        return true
    }
    return false
}
```

#### repository.go (Interface)

```go
package appointment

import (
    "context"
    "time"
    "github.com/google/uuid"
)

type Repository interface {
    // CRUD
    Save(ctx context.Context, appointment *Appointment) error
    FindByID(ctx context.Context, tenantID, id uuid.UUID) (*Appointment, error)
    Update(ctx context.Context, appointment *Appointment) error
    Delete(ctx context.Context, tenantID, id uuid.UUID) error

    // Queries
    List(ctx context.Context, tenantID uuid.UUID, filters ListFilters) ([]*Appointment, error)
    CheckConflicts(
        ctx context.Context,
        tenantID uuid.UUID,
        professionalID uuid.UUID,
        startTime time.Time,
        endTime time.Time,
        excludeID *uuid.UUID,
    ) ([]*Appointment, error)

    // Stats
    CountByStatus(ctx context.Context, tenantID uuid.UUID, status Status) (int, error)
}
```

### 3.3 Application Layer (Use Cases)

#### create_appointment.go

```go
package appointment

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "nexo/internal/domain/appointment"
    "nexo/internal/application/dto"
)

type CreateAppointmentUseCase struct {
    appointmentRepo  appointment.Repository
    professionalRepo domain.ProfessionalRepository
    customerRepo     domain.CustomerRepository
    serviceRepo      domain.ServiceRepository
    googleCalendar   external.GoogleCalendarService
}

func (uc *CreateAppointmentUseCase) Execute(
    ctx context.Context,
    tenantID uuid.UUID,
    req *dto.CreateAppointmentRequest,
) (*dto.AppointmentResponse, error) {
    // 1. Validar barbeiro pertence ao tenant
    professional, err := uc.professionalRepo.FindByID(ctx, tenantID, req.ProfessionalID)
    if err != nil {
        return nil, fmt.Errorf("professional not found: %w", err)
    }

    if !professional.IsActive {
        return nil, appointment.ErrProfessionalInactive
    }

    // 2. Validar cliente pertence ao tenant
    customer, err := uc.customerRepo.FindByID(ctx, tenantID, req.CustomerID)
    if err != nil {
        return nil, fmt.Errorf("customer not found: %w", err)
    }

    // 3. Buscar serviços e calcular duração total
    services, err := uc.serviceRepo.FindByIDs(ctx, tenantID, req.ServiceIDs)
    if err != nil {
        return nil, err
    }

    totalDuration := calculateTotalDuration(services)
    endTime := req.StartTime.Add(totalDuration)

    // 4. Verificar conflitos de horário
    conflicts, err := uc.appointmentRepo.CheckConflicts(
        ctx,
        tenantID,
        req.ProfessionalID,
        req.StartTime,
        endTime,
        nil, // Não excluir nenhum ID (criação nova)
    )
    if err != nil {
        return nil, err
    }

    if len(conflicts) > 0 {
        return nil, appointment.ErrTimeSlotConflict
    }

    // 5. Criar entidade de domínio
    app, err := appointment.NewAppointment(
        tenantID,
        req.ProfessionalID,
        req.CustomerID,
        req.ServiceIDs,
        req.StartTime,
        endTime,
        req.Notes,
    )
    if err != nil {
        return nil, err
    }

    // 6. Persistir
    if err := uc.appointmentRepo.Save(ctx, app); err != nil {
        return nil, fmt.Errorf("failed to save appointment: %w", err)
    }

    // 7. Sincronizar Google Agenda (async, se configurado)
    if professional.HasGoogleCalendar() && app.Status == appointment.StatusConfirmed {
        go uc.syncToGoogleCalendar(app, professional)
    }

    // 8. Retornar DTO
    return mapper.ToAppointmentResponse(app, services, professional, customer), nil
}

func calculateTotalDuration(services []*domain.Service) time.Duration {
    total := 0
    for _, s := range services {
        total += s.DurationMinutes
    }
    return time.Duration(total) * time.Minute
}
```

### 3.4 Infrastructure Layer

#### appointment_repository.go (PostgreSQL + sqlc)

```go
package postgres

import (
    "context"
    "database/sql"
    "time"

    "github.com/google/uuid"
    "nexo/internal/domain/appointment"
)

type AppointmentRepository struct {
    db      *sql.DB
    queries *sqlc.Queries
}

func (r *AppointmentRepository) Save(ctx context.Context, app *appointment.Appointment) error {
    return r.queries.CreateAppointment(ctx, sqlc.CreateAppointmentParams{
        ID:             app.ID,
        TenantID:       app.TenantID,
        ProfessionalID: app.ProfessionalID,
        CustomerID:     app.CustomerID,
        StartTime:      app.StartTime,
        EndTime:        app.EndTime,
        Status:         string(app.Status),
        Notes:          sql.NullString{String: app.Notes, Valid: app.Notes != ""},
    })
}

func (r *AppointmentRepository) CheckConflicts(
    ctx context.Context,
    tenantID uuid.UUID,
    professionalID uuid.UUID,
    startTime time.Time,
    endTime time.Time,
    excludeID *uuid.UUID,
) ([]*appointment.Appointment, error) {
    rows, err := r.queries.CheckAppointmentConflicts(ctx, sqlc.CheckAppointmentConflictsParams{
        TenantID:       tenantID,
        ProfessionalID: professionalID,
        StartTime:      startTime,
        EndTime:        endTime,
        ExcludeID:      excludeID,
    })
    if err != nil {
        return nil, err
    }

    appointments := make([]*appointment.Appointment, len(rows))
    for i, row := range rows {
        appointments[i] = r.rowToAppointment(row)
    }

    return appointments, nil
}
```

#### queries.sql (sqlc)

```sql
-- name: CreateAppointment :exec
INSERT INTO appointments (
    id, tenant_id, professional_id, customer_id,
    start_time, end_time, status, notes,
    created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
);

-- name: CheckAppointmentConflicts :many
SELECT * FROM appointments
WHERE tenant_id = $1
  AND professional_id = $2
  AND start_time < $4  -- novo.end_time
  AND end_time > $3    -- novo.start_time
  AND status NOT IN ('CANCELED', 'NO_SHOW')
  AND ($5::uuid IS NULL OR id != $5); -- Excluir ID (para updates)

-- name: ListAppointments :many
SELECT * FROM appointments
WHERE tenant_id = $1
  AND ($2::uuid IS NULL OR professional_id = $2)
  AND ($3::timestamp IS NULL OR start_time >= $3)
  AND ($4::timestamp IS NULL OR start_time <= $4)
  AND ($5::text IS NULL OR status = $5)
ORDER BY start_time ASC;
```

---

## 4. Fluxo de Dados Completo

### 4.1 Criação de Agendamento

```
[Frontend] AppointmentForm
    │
    ├─> onSubmit(formData)
    │
    ▼
[React Query] useCreateAppointment()
    │
    ├─> mutationFn: appointmentService.create(data)
    │
    ▼
[HTTP] POST /api/v1/appointments
    │   Headers: Authorization: Bearer <JWT>
    │   Body: { professional_id, customer_id, service_ids, start_time, notes }
    │
    ▼
[Backend] AppointmentHandler.Create()
    │
    ├─> Extract tenant_id from JWT (middleware)
    ├─> Bind & Validate DTO
    │
    ▼
[Use Case] CreateAppointmentUseCase.Execute()
    │
    ├─> Validate Professional (tenant_id, active)
    ├─> Validate Customer (tenant_id)
    ├─> Calculate end_time (services duration)
    ├─> Check Conflicts (SQL query)
    │   └─> If conflicts → return ErrTimeSlotConflict (409)
    ├─> Create Domain Entity
    ├─> Save to DB (transaction)
    ├─> Sync Google Calendar (async goroutine)
    │
    ▼
[Repository] appointmentRepo.Save()
    │
    ├─> sqlc.CreateAppointment()
    ├─> INSERT INTO appointments (...)
    │
    ▼
[Database] PostgreSQL
    │
    ▼
[Response] 201 Created
    │   { "data": { "id": "uuid", ... } }
    │
    ▼
[Frontend] onSuccess()
    │
    ├─> Invalidate queries (React Query)
    ├─> Close modal
    ├─> Toast: "Agendamento criado!"
    ├─> Refetch calendar
```

---

## 5. Estratégias de Cache

### 5.1 Frontend (TanStack Query)

```typescript
// Configuração global
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000, // 30s - dados considerados frescos
      gcTime: 5 * 60 * 1000, // 5min - garbage collection
      refetchOnWindowFocus: true,
      retry: 1,
    },
  },
});

// Cache keys organizados
const appointmentKeys = {
  all: ['appointments'],
  lists: () => [...appointmentKeys.all, 'list'],
  list: (filters) => [...appointmentKeys.lists(), filters],
  details: () => [...appointmentKeys.all, 'detail'],
  detail: (id) => [...appointmentKeys.details(), id],
};

// Invalidação granular
queryClient.invalidateQueries({ queryKey: appointmentKeys.lists() });
```

### 5.2 Backend (Redis - Futuro)

```go
// Cache de disponibilidade (query pesada)
func (r *CachedAppointmentRepository) CheckConflicts(...) {
    cacheKey := fmt.Sprintf("availability:%s:%s:%s",
        tenantID, professionalID, date.Format("2006-01-02"))

    // Tentar cache primeiro
    if cached, err := r.redis.Get(ctx, cacheKey); err == nil {
        return parseFromCache(cached)
    }

    // Buscar do banco
    conflicts, err := r.postgres.CheckConflicts(...)

    // Cachear por 1 minuto
    r.redis.Set(ctx, cacheKey, conflicts, 1*time.Minute)

    return conflicts, err
}
```

---

## 6. Tratamento de Erros

### 6.1 Errors de Domínio

```go
// domain/appointment/errors.go

var (
    ErrNotFound              = errors.New("appointment not found")
    ErrProfessionalInactive  = errors.New("professional is inactive")
    ErrTimeSlotConflict      = errors.New("time slot conflict")
    ErrInvalidStatusTransition = errors.New("invalid status transition")
    ErrPastAppointment       = errors.New("cannot modify past appointment")
)
```

### 6.2 HTTP Error Handling

```go
// handler/appointment_handler.go

func (h *AppointmentHandler) handleError(c echo.Context, err error) error {
    switch {
    case errors.Is(err, appointment.ErrNotFound):
        return c.JSON(404, ErrorResponse{
            Code:    "NOT_FOUND",
            Message: "Agendamento não encontrado",
        })

    case errors.Is(err, appointment.ErrTimeSlotConflict):
        return c.JSON(409, ErrorResponse{
            Code:    "TIME_SLOT_CONFLICT",
            Message: "Horário já está ocupado",
        })

    case errors.Is(err, appointment.ErrProfessionalInactive):
        return c.JSON(400, ErrorResponse{
            Code:    "PROFESSIONAL_INACTIVE",
            Message: "Barbeiro está inativo",
        })

    default:
        log.Error("unexpected error", "error", err)
        return c.JSON(500, ErrorResponse{
            Code:    "INTERNAL_ERROR",
            Message: "Erro interno do servidor",
        })
    }
}
```

### 6.3 Frontend Error Handling

```typescript
const { mutate } = useCreateAppointment();

mutate(data, {
  onError: (error: ApiError) => {
    switch (error.code) {
      case 'TIME_SLOT_CONFLICT':
        form.setError('start_time', {
          message: 'Horário já ocupado. Escolha outro horário.',
        });
        break;

      case 'PROFESSIONAL_INACTIVE':
        toast.error('Barbeiro não está disponível');
        break;

      default:
        toast.error('Erro ao criar agendamento. Tente novamente.');
    }
  },
});
```

---

## 7. Segurança

### 7.1 Multi-Tenant Isolation

```go
// TODAS as queries DEVEM incluir tenant_id
func (r *AppointmentRepository) FindByID(
    ctx context.Context,
    tenantID uuid.UUID,
    id uuid.UUID,
) (*Appointment, error) {
    row := r.queries.GetAppointment(ctx, sqlc.GetAppointmentParams{
        TenantID: tenantID, // OBRIGATÓRIO
        ID:       id,
    })
    // ...
}
```

### 7.2 RBAC Middleware

```go
// middleware/rbac.go

func RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            userRole := c.Get("user_role").(string)

            for _, role := range allowedRoles {
                if role == userRole {
                    return next(c)
                }
            }

            return c.JSON(403, ErrorResponse{
                Code:    "FORBIDDEN",
                Message: "Você não tem permissão para esta ação",
            })
        }
    }
}

// Uso:
e.POST("/appointments", handler.Create, RequireRole("owner", "manager", "receptionist"))
e.GET("/appointments", handler.List, RequireRole("owner", "manager", "receptionist", "barbeiro"))
```

### 7.3 Validação de Ownership

```go
// Barbeiro só pode ver seus próprios agendamentos
func (uc *ListAppointmentsUseCase) Execute(...) {
    if userRole == "barbeiro" {
        filters.ProfessionalID = &userID // Forçar filtro
    }

    return uc.appointmentRepo.List(ctx, tenantID, filters)
}
```

---

## 8. Performance

### 8.1 Índices de Banco de Dados

```sql
-- Query principal: buscar agendamentos por barbeiro e data
CREATE INDEX idx_appointments_tenant_professional_date
ON appointments(tenant_id, professional_id, start_time DESC);

-- Query de conflitos
CREATE INDEX idx_appointments_time_range
ON appointments(tenant_id, professional_id, start_time, end_time)
WHERE status NOT IN ('CANCELED', 'NO_SHOW');

-- Query por cliente
CREATE INDEX idx_appointments_customer
ON appointments(tenant_id, customer_id, start_time DESC);

-- Query por status
CREATE INDEX idx_appointments_status
ON appointments(tenant_id, status, start_time DESC);
```

### 8.2 Paginação

```go
// DTO
type ListAppointmentsRequest struct {
    Page     int       `query:"page" validate:"min=1"`
    PageSize int       `query:"page_size" validate:"min=1,max=100"`
    // ... filters
}

// Repository
func (r *AppointmentRepository) List(...) {
    offset := (filters.Page - 1) * filters.PageSize

    return r.queries.ListAppointments(ctx, sqlc.ListAppointmentsParams{
        Limit:  filters.PageSize,
        Offset: offset,
        // ... other params
    })
}
```

### 8.3 Lazy Loading (Frontend)

```typescript
// Carregar apenas agendamentos visíveis no calendário
const { data } = useAppointments({
  dateFrom: startOfWeek(currentDate),
  dateTo: endOfWeek(currentDate),
  professionalIds: selectedProfessionals, // Filtrar por barbeiros selecionados
});
```

---

## 9. Testes

### 9.1 Backend (Go)

```go
// usecase/appointment/create_appointment_test.go

func TestCreateAppointmentUseCase_Execute_Success(t *testing.T) {
    // Arrange
    mockRepo := &MockAppointmentRepository{}
    mockProfRepo := &MockProfessionalRepository{}

    uc := NewCreateAppointmentUseCase(mockRepo, mockProfRepo, ...)

    req := &dto.CreateAppointmentRequest{
        ProfessionalID: uuid.New(),
        CustomerID:     uuid.New(),
        ServiceIDs:     []uuid.UUID{uuid.New()},
        StartTime:      time.Now().Add(24 * time.Hour),
    }

    // Mock: professional exists and is active
    mockProfRepo.On("FindByID", mock.Anything, mock.Anything, req.ProfessionalID).
        Return(&domain.Professional{IsActive: true}, nil)

    // Mock: no conflicts
    mockRepo.On("CheckConflicts", mock.Anything, mock.Anything, mock.Anything).
        Return([]*appointment.Appointment{}, nil)

    // Mock: save succeeds
    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

    // Act
    resp, err := uc.Execute(context.Background(), uuid.New(), req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, appointment.StatusCreated, resp.Status)
}

func TestCreateAppointmentUseCase_Execute_Conflict(t *testing.T) {
    // ... setup

    // Mock: conflict exists
    mockRepo.On("CheckConflicts", mock.Anything, mock.Anything, mock.Anything).
        Return([]*appointment.Appointment{{ID: uuid.New()}}, nil)

    // Act
    _, err := uc.Execute(context.Background(), uuid.New(), req)

    // Assert
    assert.ErrorIs(t, err, appointment.ErrTimeSlotConflict)
}
```

### 9.2 Frontend (Jest + React Testing Library)

```typescript
// components/appointments/__tests__/AppointmentForm.test.tsx

describe('AppointmentForm', () => {
  it('should create appointment successfully', async () => {
    const mockMutate = jest.fn();
    (useCreateAppointment as jest.Mock).mockReturnValue({
      mutate: mockMutate,
      isPending: false,
    });

    render(<AppointmentForm />);

    // Fill form
    await userEvent.selectOptions(screen.getByLabelText('Cliente'), 'cliente-1');
    await userEvent.selectOptions(
      screen.getByLabelText('Barbeiro'),
      'barbeiro-1'
    );
    await userEvent.click(screen.getByLabelText('Serviços'));
    await userEvent.click(screen.getByText('Corte'));

    // Submit
    await userEvent.click(screen.getByRole('button', { name: /agendar/i }));

    expect(mockMutate).toHaveBeenCalledWith(
      expect.objectContaining({
        customer_id: 'cliente-1',
        professional_id: 'barbeiro-1',
        service_ids: expect.arrayContaining([expect.any(String)]),
      }),
      expect.any(Object)
    );
  });

  it('should show error on conflict', async () => {
    const mockMutate = jest.fn((_, { onError }) => {
      onError({ code: 'TIME_SLOT_CONFLICT' });
    });

    (useCreateAppointment as jest.Mock).mockReturnValue({
      mutate: mockMutate,
      isPending: false,
    });

    render(<AppointmentForm />);

    // ... fill and submit

    await waitFor(() => {
      expect(screen.getByText(/horário já ocupado/i)).toBeInTheDocument();
    });
  });
});
```

---

## 10. Monitoramento

### 10.1 Métricas (Prometheus)

```go
var (
    appointmentCreationDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "appointment_creation_duration_seconds",
            Help: "Tempo de criação de agendamento",
        },
        []string{"tenant_id"},
    )

    appointmentConflicts = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "appointment_conflicts_total",
            Help: "Total de conflitos de horário detectados",
        },
        []string{"tenant_id"},
    )
)

// No use case:
func (uc *CreateAppointmentUseCase) Execute(...) {
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        appointmentCreationDuration.WithLabelValues(tenantID.String()).Observe(duration)
    }()

    // ... lógica
}
```

### 10.2 Logs Estruturados

```go
import "go.uber.org/zap"

logger.Info("creating appointment",
    zap.String("tenant_id", tenantID.String()),
    zap.String("professional_id", req.ProfessionalID.String()),
    zap.Time("start_time", req.StartTime),
)

logger.Error("conflict detected",
    zap.String("tenant_id", tenantID.String()),
    zap.String("professional_id", req.ProfessionalID.String()),
    zap.Int("conflicts_count", len(conflicts)),
)
```

---

**Próximos Documentos:**
- `BANCO_AGENDAMENTO.md` - Schema de banco completo
- `API_AGENDAMENTO.md` - Contrato da API REST
- `DIAGRAMAS_AGENDAMENTO.md` - Fluxogramas detalhados
- `CHECKLIST_IMPLEMENTACAO.md` - Checklist de desenvolvimento

---

**Responsável:** Tech Lead  
**Data:** 25/11/2025  
**Status:** ✅ COMPLETO
