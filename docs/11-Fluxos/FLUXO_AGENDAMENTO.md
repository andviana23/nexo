# Fluxo de Agendamento — NEXO v1.0

**Versão:** 1.2
**Última Atualização:** 27/11/2025
**Status:** 🟡 Parcialmente Implementado (75%)
**Responsável:** Product + Tech Lead

---

## 📊 Status de Implementação

| Área | Status | Progresso |
|------|--------|-----------|
| Backend (Go) | ✅ Completo | 100% |
| Frontend Base | ✅ Completo | 100% |
| Menu de Ações | ❌ Não iniciado | 0% |
| Comanda | ❌ Não iniciado | 0% |
| Pagamento Multi-Forma | ❌ Não iniciado | 0% |
| Status Extras (CHECKED_IN, etc) | ❌ Não iniciado | 0% |

### ✅ Implementado (27/11/2025)
- Backend completo: Entity, Repository, 6 Use Cases, 6 Handlers
- 36 testes (26 unit + 10 integration)
- Frontend: Calendário FullCalendar, Modal, Selectors conectados à API
- React Query hooks com Optimistic Updates
- 6 status: CREATED, CONFIRMED, IN_SERVICE, DONE, NO_SHOW, CANCELED

### ⏳ Pendente para MVP
- Status CHECKED_IN e AWAITING_PAYMENT
- Menu de ações (contexto/três pontinhos)
- Comanda estilo Trinks
- Pagamento multi-forma
- Drag & drop vertical (trocar profissional)

### 🚫 Movido para Futuro
- Fluxo 10: Consumo Interno (Produto do Estoque)
- Fluxo 11: Troca de Profissional com Split de Comissão
- Histórico de Edições

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

### RN-AGE-005: Status de Agendamento (Lifecycle Completo)

Status permitidos (ordem do fluxo):

- `CREATED` - Criado (pendente confirmação)
- `CONFIRMED` - Confirmado pelo cliente
- `CHECKED_IN` - Cliente chegou (marcou presença) ⭐ NOVO
- `IN_SERVICE` - Em atendimento
- `AWAITING_PAYMENT` - Aguardando pagamento ⭐ NOVO
- `DONE` - Finalizado com sucesso
- `NO_SHOW` - Cliente faltou
- `CANCELED` - Cancelado (cliente ou barbearia)

**Transições válidas:**
```
CREATED → CONFIRMED → CHECKED_IN → IN_SERVICE → AWAITING_PAYMENT → DONE
                  ↓         ↓           ↓              ↓
               CANCELED  NO_SHOW    CANCELED       CANCELED
```

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

### RN-AGE-008: Menu de Ações do Agendamento ⭐ NOVO

Ações disponíveis ao clicar com botão direito ou menu de três pontinhos:

**Ações de Status:**
- ✅ Confirmar agendamento (CREATED → CONFIRMED)
- ✅ Cliente chegou (CONFIRMED → CHECKED_IN)
- ✅ Iniciar atendimento (CHECKED_IN → IN_SERVICE)
- ✅ Finalizar atendimento (IN_SERVICE → AWAITING_PAYMENT)
- ✅ Marcar como concluído (AWAITING_PAYMENT → DONE)
- ✅ Cliente faltou (→ NO_SHOW)
- ✅ Cancelar agendamento (→ CANCELED)

**Ações de Edição:**
- ✅ Adicionar serviços extras
- ✅ Editar serviços existentes
- ✅ Trocar barbeiro/profissional
- ✅ Mover horário (drag & drop + opção manual)
- ✅ Reagendar para data futura

**Ações de Comanda:**
- ✅ Abrir comanda imediatamente
- ✅ Transformar em venda sem agendamento (check-in rápido/encaixe)

**Ações de Cliente:**
- ✅ Ver histórico do cliente
- ✅ Ver atendimentos anteriores
- ✅ Criar anotações internas do cliente

### RN-AGE-009: Drag & Drop Avançado ⭐ NOVO

Suporte a arrastar agendamentos:

- ✅ **Horizontal:** Mover para outro horário (mesmo barbeiro)
- ✅ **Vertical:** Mover para outro barbeiro (mesmo horário)
- ✅ **Diagonal:** Mover para outro barbeiro E horário
- ✅ Validação de conflitos em tempo real durante o arraste
- ✅ Confirmação visual antes de soltar

### RN-AGE-010: Histórico de Edições ⭐ NOVO

Registrar todas as alterações do agendamento:

- ✅ Quem alterou (user_id)
- ✅ Quando alterou (timestamp)
- ✅ O que alterou (campo anterior → campo novo)
- ✅ Motivo da alteração (opcional)
- ✅ Visualização em timeline no modal do agendamento

### RN-AGE-011: Integração com Comanda ⭐ NOVO

Ao transitar para `IN_SERVICE`:

- ✅ Criar comanda automaticamente vinculada ao agendamento
- ✅ Popular serviços do agendamento na comanda
- ✅ Permitir adicionar/remover serviços na comanda
- ✅ Permitir adicionar produtos (consumo interno)
- ✅ Baixa automática de estoque ao adicionar produto
- ✅ Permitir trocar profissional durante atendimento (split de comissão)

### RN-AGE-012: Pagamento Multi-Forma ⭐ NOVO

Suporte a pagamento dividido:

- ✅ Múltiplas formas de pagamento na mesma comanda
  - Ex: R$ 50 no PIX + R$ 30 no cartão
- ✅ Pagamento de duas ou mais comandas juntas
  - Ex: Pai e filho pagam junto
- ✅ Desconto aplicado ao total
- ✅ Gorjeta opcional

### RN-AGE-013: Consumo Interno com Estoque ⭐ NOVO

Produtos consumidos durante atendimento:

- ✅ Adicionar produto à comanda (ex: pomada, cerveja)
- ✅ Baixa automática do estoque
- ✅ Produto não afeta comissão do barbeiro (configurável)
- ✅ Histórico de consumo por cliente

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

### Fluxo 5: Check-in do Cliente (Cliente Chegou) ⭐ NOVO

```
[Cliente chega na barbearia]
   ↓
[Recepção localiza agendamento no calendário]
   ↓
[Clica em "Cliente chegou" ou arrasta para coluna Check-in]
   ↓
[PUT /api/appointments/:id/check-in]
   ↓
[Atualizar status: CHECKED_IN]
   ↓
[Exibir na fila de espera do barbeiro]
   ↓
[Notificar barbeiro (push/som)]
   ↓
[Registrar hora de chegada]
   ↓
[Fim]
```

### Fluxo 6: Início do Atendimento ⭐ NOVO

```
[Barbeiro está livre]
   ↓
[Clica em "Iniciar atendimento" no cliente da fila]
   ↓
[PUT /api/appointments/:id/start-service]
   ↓
[Atualizar status: IN_SERVICE]
   ↓
[Criar comanda automaticamente]
   ↓
[POST /api/commands]
   ↓
[Popular serviços do agendamento na comanda]
   ↓
[Registrar hora de início]
   ↓
[Fim]
```

### Fluxo 7: Edição da Comanda Durante Atendimento ⭐ NOVO

```
[Durante o atendimento]
   ↓
[Barbeiro/Recepção abre comanda vinculada]
   ↓
[Opções disponíveis:]
   ├─→ [Adicionar serviço extra] → [Atualizar duração/preço]
   ├─→ [Remover serviço] → [Atualizar duração/preço]
   ├─→ [Adicionar produto] → [Baixar estoque automaticamente]
   ├─→ [Trocar profissional] → [Split de comissão se necessário]
   └─→ [Aplicar desconto] → [Registrar motivo]
   ↓
[PUT /api/commands/:id]
   ↓
[Registrar alteração no histórico]
   ↓
[Atualizar UI em tempo real]
   ↓
[Fim]
```

### Fluxo 8: Finalização e Pagamento ⭐ NOVO

```
[Atendimento concluído]
   ↓
[Clica em "Finalizar atendimento"]
   ↓
[PUT /api/appointments/:id/finish]
   ↓
[Atualizar status: AWAITING_PAYMENT]
   ↓
[Exibir resumo da comanda]
   ↓
[Selecionar forma(s) de pagamento]
   ├─→ [Pagamento único] → [100% em uma forma]
   └─→ [Pagamento dividido] → [Distribuir valor entre formas]
   ↓
[POST /api/payments]
   ↓
[Registrar pagamento]
   ↓
[Atualizar status: DONE]
   ↓
[Calcular comissão do barbeiro]
   ↓
[Emitir recibo/NF (se configurado)]
   ↓
[Fim]
```

### Fluxo 9: Pagamento de Múltiplas Comandas ⭐ NOVO

```
[Cliente quer pagar duas ou mais comandas juntas]
   ↓
[Recepção seleciona comandas do mesmo cliente/grupo]
   ↓
[Sistema agrupa em um único pagamento]
   ↓
[Exibir total consolidado]
   ↓
[Selecionar forma(s) de pagamento]
   ↓
[POST /api/payments/bulk]
   ↓
[Distribuir pagamento entre comandas]
   ↓
[Atualizar status de todos: DONE]
   ↓
[Fim]
```

### Fluxo 10: Consumo Interno (Produto do Estoque) ⭐ NOVO

```
[Durante atendimento, cliente consome produto]
   ↓
[Ex: Cerveja, pomada, gel]
   ↓
[Adicionar produto à comanda]
   ↓
[POST /api/commands/:id/items]
   ↓
[Verificar estoque disponível]
   ↓
{Estoque suficiente?}
   ├─→ [Não] → [Exibir alerta de estoque baixo]
   └─→ [Sim] → [Continuar]
   ↓
[Baixar quantidade do estoque]
   ↓
[PUT /api/inventory/:product_id/decrement]
   ↓
[Adicionar ao total da comanda]
   ↓
[Produto NÃO entra na comissão (configurável)]
   ↓
[Fim]
```

### Fluxo 11: Troca de Profissional Durante Atendimento ⭐ NOVO

```
[Necessidade de trocar barbeiro durante atendimento]
   ↓
[Abrir comanda vinculada]
   ↓
[Clicar em "Trocar profissional"]
   ↓
[Selecionar novo profissional]
   ↓
[Definir split de comissão:]
   ├─→ [100% novo] → [Comissão toda para novo barbeiro]
   ├─→ [50/50] → [Dividir entre ambos]
   └─→ [Proporcional] → [Baseado no tempo de cada um]
   ↓
[PUT /api/commands/:id/transfer]
   ↓
[Registrar histórico de transferência]
   ↓
[Notificar ambos profissionais]
   ↓
[Fim]
```

### Fluxo 12: Check-in Rápido (Encaixe/Venda sem Agendamento) ⭐ NOVO

```
[Cliente chega sem agendamento]
   ↓
[Recepção clica em "Encaixe" ou "Venda rápida"]
   ↓
[Selecionar/criar cliente]
   ↓
[Selecionar serviço(s)]
   ↓
[Selecionar barbeiro disponível]
   ↓
[Criar agendamento com status CHECKED_IN diretamente]
   ↓
[POST /api/appointments (status: CHECKED_IN)]
   ↓
[Seguir fluxo normal de atendimento]
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
    command_id UUID REFERENCES commands(id) ON DELETE SET NULL, -- ⭐ NOVO: Vinculo com comanda
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    checked_in_at TIMESTAMP, -- ⭐ NOVO: Hora que cliente chegou
    service_started_at TIMESTAMP, -- ⭐ NOVO: Hora que iniciou atendimento
    service_finished_at TIMESTAMP, -- ⭐ NOVO: Hora que finalizou
    status VARCHAR(50) NOT NULL DEFAULT 'CREATED',
    notes TEXT,
    google_event_id VARCHAR(255),
    canceled_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_end_after_start CHECK (end_time > start_time),
    CONSTRAINT chk_valid_status CHECK (status IN (
        'CREATED', 'CONFIRMED', 'CHECKED_IN', 'IN_SERVICE', 
        'AWAITING_PAYMENT', 'DONE', 'NO_SHOW', 'CANCELED'
    ))
);

CREATE INDEX idx_appointments_tenant_professional_date
  ON appointments(tenant_id, professional_id, start_time DESC);

CREATE INDEX idx_appointments_tenant_customer
  ON appointments(tenant_id, customer_id);

CREATE INDEX idx_appointments_tenant_status
  ON appointments(tenant_id, status);

CREATE INDEX idx_appointments_command
  ON appointments(command_id);
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

### Tabela: `appointment_history` ⭐ NOVO (Histórico de Alterações)

```sql
CREATE TABLE appointment_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    appointment_id UUID NOT NULL REFERENCES appointments(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL, -- 'CREATED', 'STATUS_CHANGED', 'RESCHEDULED', 'PROFESSIONAL_CHANGED', etc.
    field_changed VARCHAR(100), -- 'status', 'start_time', 'professional_id', etc.
    old_value TEXT,
    new_value TEXT,
    reason TEXT, -- Motivo da alteração (opcional)
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_appointment_history_appointment
  ON appointment_history(appointment_id, created_at DESC);
```

### Tabela: `commands` ⭐ NOVO (Comandas)

```sql
CREATE TABLE commands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    appointment_id UUID REFERENCES appointments(id) ON DELETE SET NULL,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    professional_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN', -- 'OPEN', 'AWAITING_PAYMENT', 'PAID', 'CANCELED'
    subtotal DECIMAL(10,2) NOT NULL DEFAULT 0,
    discount DECIMAL(10,2) DEFAULT 0,
    discount_reason TEXT,
    tip DECIMAL(10,2) DEFAULT 0,
    total DECIMAL(10,2) NOT NULL DEFAULT 0,
    opened_at TIMESTAMP DEFAULT NOW(),
    closed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_commands_tenant_status
  ON commands(tenant_id, status);

CREATE INDEX idx_commands_tenant_customer
  ON commands(tenant_id, customer_id);

CREATE INDEX idx_commands_appointment
  ON commands(appointment_id);
```

### Tabela: `command_items` ⭐ NOVO (Itens da Comanda)

```sql
CREATE TABLE command_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    command_id UUID NOT NULL REFERENCES commands(id) ON DELETE CASCADE,
    item_type VARCHAR(20) NOT NULL, -- 'SERVICE' ou 'PRODUCT'
    service_id UUID REFERENCES services(id) ON DELETE SET NULL,
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,
    professional_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Quem executou o serviço
    name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    affects_commission BOOLEAN DEFAULT TRUE, -- Produtos geralmente não afetam
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_item_type CHECK (
        (item_type = 'SERVICE' AND service_id IS NOT NULL) OR
        (item_type = 'PRODUCT' AND product_id IS NOT NULL)
    )
);

CREATE INDEX idx_command_items_command
  ON command_items(command_id);
```

### Tabela: `command_payments` ⭐ NOVO (Pagamentos Multi-Forma)

```sql
CREATE TABLE command_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    command_id UUID NOT NULL REFERENCES commands(id) ON DELETE CASCADE,
    payment_method VARCHAR(50) NOT NULL, -- 'PIX', 'CREDIT_CARD', 'DEBIT_CARD', 'CASH', 'TRANSFER'
    amount DECIMAL(10,2) NOT NULL,
    reference VARCHAR(255), -- ID da transação, NSU, etc.
    paid_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_command_payments_command
  ON command_payments(command_id);
```

### Tabela: `payment_groups` ⭐ NOVO (Pagamento de Múltiplas Comandas)

```sql
CREATE TABLE payment_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    total_amount DECIMAL(10,2) NOT NULL,
    paid_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE payment_group_commands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_group_id UUID NOT NULL REFERENCES payment_groups(id) ON DELETE CASCADE,
    command_id UUID NOT NULL REFERENCES commands(id) ON DELETE CASCADE,
    amount_from_group DECIMAL(10,2) NOT NULL, -- Quanto deste pagamento foi para esta comanda
    UNIQUE(payment_group_id, command_id)
);
```

### Tabela: `command_professional_splits` ⭐ NOVO (Split de Comissão)

```sql
CREATE TABLE command_professional_splits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    command_id UUID NOT NULL REFERENCES commands(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    percentage DECIMAL(5,2) NOT NULL, -- Ex: 50.00 para 50%
    amount DECIMAL(10,2) NOT NULL,
    reason TEXT, -- 'TRANSFER', 'ASSISTANCE', etc.
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_percentage CHECK (percentage >= 0 AND percentage <= 100)
);

CREATE INDEX idx_command_splits_command
  ON command_professional_splits(command_id);
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

### PUT `/api/appointments/:id/confirm` ⭐ NOVO

**Descrição:** Confirmar agendamento
**Status:** CREATED → CONFIRMED

### PUT `/api/appointments/:id/check-in` ⭐ NOVO

**Descrição:** Marcar cliente como chegou
**Status:** CONFIRMED → CHECKED_IN
**Ação:** Registra `checked_in_at`

### PUT `/api/appointments/:id/start-service` ⭐ NOVO

**Descrição:** Iniciar atendimento
**Status:** CHECKED_IN → IN_SERVICE
**Ação:** Cria comanda automaticamente, registra `service_started_at`

### PUT `/api/appointments/:id/finish` ⭐ NOVO

**Descrição:** Finalizar atendimento
**Status:** IN_SERVICE → AWAITING_PAYMENT
**Ação:** Registra `service_finished_at`

### GET `/api/appointments/:id/history` ⭐ NOVO

**Descrição:** Histórico de alterações do agendamento

---

## 📡 Endpoints de Comanda ⭐ NOVO

### POST `/api/commands`

**Descrição:** Criar nova comanda
**Body:**
```json
{
  "appointment_id": "uuid (opcional)",
  "customer_id": "uuid",
  "professional_id": "uuid"
}
```

### GET `/api/commands/:id`

**Descrição:** Detalhes da comanda com itens

### POST `/api/commands/:id/items`

**Descrição:** Adicionar item à comanda (serviço ou produto)
**Body:**
```json
{
  "item_type": "SERVICE | PRODUCT",
  "service_id": "uuid (se SERVICE)",
  "product_id": "uuid (se PRODUCT)",
  "quantity": 1,
  "unit_price": "45.00"
}
```

### DELETE `/api/commands/:id/items/:item_id`

**Descrição:** Remover item da comanda

### PUT `/api/commands/:id/transfer`

**Descrição:** Trocar profissional (com split de comissão)
**Body:**
```json
{
  "new_professional_id": "uuid",
  "split_type": "FULL | HALF | PROPORTIONAL",
  "reason": "Troca durante atendimento"
}
```

### PUT `/api/commands/:id/discount`

**Descrição:** Aplicar desconto
**Body:**
```json
{
  "amount": "10.00",
  "reason": "Desconto fidelidade"
}
```

---

## 📡 Endpoints de Pagamento ⭐ NOVO

### POST `/api/payments`

**Descrição:** Registrar pagamento de uma comanda
**Body:**
```json
{
  "command_id": "uuid",
  "payments": [
    { "method": "PIX", "amount": "50.00" },
    { "method": "CREDIT_CARD", "amount": "30.00" }
  ],
  "tip": "5.00"
}
```

### POST `/api/payments/bulk`

**Descrição:** Pagar múltiplas comandas juntas
**Body:**
```json
{
  "command_ids": ["uuid1", "uuid2"],
  "payments": [
    { "method": "PIX", "amount": "150.00" }
  ]
}
```

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

### Agendamento Básico ✅ IMPLEMENTADO
- [x] Backend implementado (Domain + Use Cases + Handlers)
- [x] Frontend com calendário visual (FullCalendar)
- [x] Validação de conflitos de horário funcionando
- [x] Multi-tenant isolamento garantido
- [ ] Integração Google Agenda ativa (mock apenas)
- [x] RBAC respeitado (barbeiro read-only na própria agenda)
- [x] Testes cobrindo fluxo principal (36 testes)

### Status Lifecycle ⏳ PENDENTE
- [ ] Status CHECKED_IN (Cliente chegou) implementado
- [x] Status IN_SERVICE (Em atendimento) implementado
- [ ] Status AWAITING_PAYMENT (Aguardando pagamento) implementado
- [x] Transições de status validadas no backend
- [x] Indicadores visuais por status no calendário

### Drag & Drop ⏳ PARCIAL
- [x] Drag & drop horizontal (mudar horário)
- [ ] Drag & drop vertical (mudar profissional)
- [ ] Validação de conflitos durante arraste
- [ ] Confirmação visual antes de soltar

### Menu de Ações ⏳ PENDENTE
- [ ] Menu de contexto (botão direito / três pontinhos)
- [ ] Todas as ações de status disponíveis
- [ ] Ações de edição (adicionar/remover serviços)
- [ ] Ações de comanda (abrir, editar)
- [ ] Ações de cliente (histórico, anotações)

### Comanda ⏳ PENDENTE
- [ ] Criação automática ao iniciar atendimento
- [ ] Adicionar/remover serviços durante atendimento
- [ ] Adicionar produtos (consumo interno) 🚫 FUTURO
- [ ] Troca de profissional com split de comissão 🚫 FUTURO
- [ ] Aplicar desconto com motivo

### Pagamento ⏳ PENDENTE
- [ ] Pagamento único (uma forma)
- [ ] Pagamento dividido (múltiplas formas)
- [ ] Pagamento de múltiplas comandas juntas
- [ ] Gorjeta opcional
- [ ] Cálculo automático de comissão

### Estoque 🚫 FUTURO (v1.2.0)
- [ ] Baixa automática ao adicionar produto à comanda
- [ ] Alerta de estoque baixo
- [ ] Produto não afeta comissão (configurável)

### Histórico 🚫 FUTURO
- [ ] Registro de todas as alterações
- [ ] Timeline visual no modal
- [ ] Identificação de quem/quando alterou

---

## 📊 Métricas de Sucesso

**Operacionais:**

- Taxa de no-show < 10%
- Tempo médio para criar agendamento < 30 segundos
- Conflitos de horário zerados
- Tempo médio de check-in a pagamento < 45 minutos ⭐ NOVO

**Técnicas:**

- Latência API < 150ms
- Sincronização Google < 500ms
- Uptime > 99.5%
- Baixa de estoque em tempo real ⭐ NOVO

---

## 📚 Referências

- `docs/02-arquitetura/ARQUITETURA.md` - Clean Architecture
- `docs/02-arquitetura/MODELO_MULTI_TENANT.md` - Isolamento de dados
- `docs/04-backend/GUIA_DEV_BACKEND.md` - Padrões Go
- `docs/03-frontend/GUIA_FRONTEND.md` - Padrões React/Next.js
- `docs/11-Fluxos/FLUXO_COMANDA.md` - Fluxo detalhado de comandas ⭐ NOVO
- `docs/11-Fluxos/FLUXO_PAGAMENTO.md` - Fluxo de pagamentos ⭐ NOVO
- `PRD-NEXO.md` - Requisitos de produto

---

**Status:** 🟡 Parcialmente Implementado (75%)
**Próximo Marco:** Comanda e Menu de Ações
**Última Revisão:** 27/11/2025
