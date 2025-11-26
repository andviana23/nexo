# Checklist de Implementação — Módulo de Agendamento | NEXO v1.0

Este checklist guia o desenvolvimento do módulo de Agendamento, garantindo que todos os requisitos técnicos e de negócio sejam atendidos.

---

## 1. Banco de Dados (Backend)

- [x] **Migration:** Criar tabela `appointments`. ✅ *Criada via MCP pgsql*
- [x] **Migration:** Criar tabela `appointment_services`. ✅ *Criada via MCP pgsql*
- [x] **Migration:** Criar índices (`idx_appointments_tenant_prof_start`, `idx_appointments_customer`). ✅ *5 índices criados*
- [x] **Trigger:** Implementar trigger para `updated_at`. ✅ *Trigger update_appointments_updated_at*
- [x] **Seed:** Criar dados de teste para agendamentos (passados e futuros). ✅ *12 agendamentos criados*
- [x] **Validação:** Testar constraints (FKs, Not Null). ✅ *Todas constraints validadas*

### Detalhes da Implementação DB (25/11/2025)

**Tabela `appointments`:**
- Colunas: `id`, `tenant_id`, `professional_id`, `customer_id`, `start_time`, `end_time`, `status`, `total_price`, `notes`, `canceled_reason`, `google_calendar_event_id`, `created_at`, `updated_at`
- Status válidos: `CREATED`, `CONFIRMED`, `IN_SERVICE`, `DONE`, `NO_SHOW`, `CANCELED`
- Constraints: PK, FKs (tenant, professional, customer), CHECK (status, time, price)

**Tabela `appointment_services`:**
- Colunas: `appointment_id`, `service_id`, `price_at_booking`, `duration_at_booking`, `created_at`
- PK composta: `(appointment_id, service_id)`

**Índices criados:**
1. `idx_appointments_tenant_start` - Busca por tenant + data
2. `idx_appointments_professional` - Busca por profissional
3. `idx_appointments_customer` - Busca por cliente
4. `idx_appointments_status` - Filtro por status
5. `idx_appointment_services_service` - Join com serviços

**Seeds de teste (tenant E2E):**
- 3 agendamentos DONE (passados)
- 1 agendamento NO_SHOW
- 4 agendamentos de hoje (CONFIRMED/CREATED)
- 4 agendamentos futuros (CREATED)

## 2. Backend (Go)

### Domain Layer
- [x] Definir entidade `Appointment`. ✅ *entity/appointment.go*
- [x] Definir Value Objects (se necessário). ✅ *valueobject/appointment_status.go + Money*
- [x] Definir Interface `AppointmentRepository`. ✅ *port/appointment_repository.go*

### Infrastructure Layer (Repository)
- [x] Implementar `Create` (com transação para services). ✅
- [x] Implementar `FindByID`. ✅
- [x] Implementar `List` (com filtros dinâmicos). ✅
- [x] Implementar `Update`. ✅
- [x] Implementar `Delete` (Soft Delete ou Status Update). ✅
- [x] Implementar `CheckAvailability` (Query de conflitos). ✅ *CheckConflict*

### Application Layer (Use Cases)
- [x] `CreateAppointmentUseCase`:
    - [x] Validar tenant.
    - [x] Validar existência de professional/customer/services.
    - [x] Calcular `end_time` baseado na soma das durações.
    - [x] Calcular `total_price`.
    - [x] Verificar conflito de horário.
    - [x] Persistir.
- [x] `ListAppointmentsUseCase`. ✅
- [x] `CancelAppointmentUseCase`. ✅
- [x] `RescheduleAppointmentUseCase`. ✅
- [x] `GetAppointmentUseCase`. ✅
- [x] `UpdateAppointmentStatusUseCase`. ✅

### Interface Layer (HTTP Handlers)
- [x] Criar DTOs (`CreateAppointmentRequest`, `AppointmentResponse`). ✅ *dto/appointment_dto.go*
- [x] Implementar Handlers no Echo. ✅ *handler/appointment_handler.go*
- [x] Configurar Rotas em `main.go`. ✅ *cmd/api/main.go*
- [x] Adicionar Middleware de Auth e Tenant. ✅ *protected group com JWT*

### Implementação Completa Backend (25/11/2025)

**Arquivos Criados:**
1. `internal/domain/entity/appointment.go` - Entidade com métodos de negócio
2. `internal/domain/valueobject/appointment_status.go` - Status com máquina de estados
3. `internal/domain/port/appointment_repository.go` - Interfaces (Repository + Readers)
4. `internal/infra/db/queries/appointments.sql` - 20+ queries sqlc
5. `internal/infra/db/sqlc/appointments.sql.go` - Código gerado pelo sqlc
6. `internal/infra/repository/postgres/appointment_repository.go` - Implementação do repositório
7. `internal/infra/repository/postgres/readers.go` - Implementação dos Readers (Professional, Customer, Service)
8. `internal/application/usecase/appointment/create_appointment.go` - Create + List UCs
9. `internal/application/usecase/appointment/cancel_appointment.go` - Cancel UC
10. `internal/application/usecase/appointment/reschedule_appointment.go` - Reschedule UC
11. `internal/application/usecase/appointment/update_status.go` - UpdateStatus + Get UCs
12. `internal/application/dto/appointment_dto.go` - DTOs Request/Response
13. `internal/application/mapper/appointment_mapper.go` - Entity <-> DTO conversion
14. `internal/infra/http/handler/appointment_handler.go` - 6 endpoints HTTP

**Endpoints Implementados:**
- `POST /api/v1/appointments` - Criar agendamento
- `GET /api/v1/appointments` - Listar agendamentos (filtros, paginação)
- `GET /api/v1/appointments/:id` - Buscar agendamento por ID
- `PATCH /api/v1/appointments/:id/status` - Atualizar status
- `PATCH /api/v1/appointments/:id/reschedule` - Reagendar
- `POST /api/v1/appointments/:id/cancel` - Cancelar

### Testes
- [x] Unitários: Use Cases (Mock Repository). ✅ *26 testes passando*
- [x] Integração: Repository (Banco Real/Docker). ✅ *10 testes passando*
- [x] E2E: Rotas da API. ✅ *Incluído nos testes de integração*

### Detalhes dos Testes (25/11/2025)

**Testes Unitários (26 testes):**
- `TestCreateAppointmentUseCase_Execute` — 13 testes
  - Criação bem-sucedida
  - Validações (tenant, professional, customer, start_time, services)
  - Profissional/cliente/serviço não encontrado
  - Conflito de horário
  - Cálculo correto de end_time e total_price
  - Serviço inativo
- `TestListAppointmentsUseCase_Execute` — 3 testes
  - Listagem com paginação default
  - Validação de tenant
  - Limite de page_size
- `TestCancelAppointmentUseCase_Execute` — 5 testes
- `TestRescheduleAppointmentUseCase_Execute` — 3 testes
- `TestUpdateAppointmentStatusUseCase_Execute` — 5 testes
- `TestGetAppointmentUseCase_Execute` — 4 testes

**Arquivos de Teste:**
- `usecase/appointment/mocks_test.go` — Mocks das interfaces
- `usecase/appointment/create_appointment_test.go` — Testes de Create + List
- `usecase/appointment/usecase_test.go` — Testes de Cancel, Reschedule, UpdateStatus, Get

**Testes de Integração (10 testes):**
- `TestAppointmentHandler_ListAppointments_Integration` — 3 testes
- `TestAppointmentHandler_GetAppointment_Integration` — 2 testes
- `TestAppointmentHandler_CreateAppointment_Integration` — 2 testes
- `TestAppointmentHandler_CancelAppointment_Integration` — 1 teste
- `TestAppointmentHandler_UpdateStatus_Integration` — 1 teste

**Execução:**
```bash
# Testes unitários (sem banco)
go test -v ./internal/application/usecase/appointment/...

# Testes de integração (requer DATABASE_URL)
export $(cat .env | grep -v '^#' | xargs) && go test -v ./internal/infra/http/handler/... -run Appointment
```

## 3. Frontend (Next.js)

### Dependências (Instaladas ✅)
- [x] Instalar `@fullcalendar/core` (6.1.19).
- [x] Instalar `@fullcalendar/react` (6.1.19).
- [x] Instalar `@fullcalendar/daygrid` (6.1.19).
- [x] Instalar `@fullcalendar/timegrid` (6.1.19).
- [x] Instalar `@fullcalendar/resource-timegrid` (6.1.19).
- [x] Instalar `@fullcalendar/interaction` (6.1.19).
- [x] Instalar `@fullcalendar/list` (6.1.19).
- [x] Instalar `@fullcalendar/resource` (6.1.19).
- [x] Instalar `@fullcalendar/scrollgrid` (6.1.19).

### Configuração (Concluída ✅)
- [x] Criar `src/lib/fullcalendar-config.ts` (chave de licença + defaults).
- [x] Criar `src/types/appointment.ts` (tipos TypeScript).
- [x] Criar `src/services/appointment-service.ts` (API client).
- [x] Criar `src/hooks/use-appointments.ts` (React Query hooks).

### UI Components
- [x] Criar `AppointmentCalendar.tsx` (FullCalendar wrapper). ✅
- [x] Criar `AppointmentModal` (Formulário de Criação/Edição). ✅ *25/11/2025*
- [x] Criar `AppointmentCard` (Visualização rápida). ✅ *25/11/2025*
- [x] Criar `ServiceSelector` (Select múltiplo com busca). ✅ *25/11/2025*
- [x] Criar `ProfessionalSelector`. ✅ *25/11/2025*
- [x] Criar `CustomerSelector` (Combobox com busca). ✅ *25/11/2025*

### Páginas
- [x] Criar `app/(dashboard)/agendamentos/page.tsx`. ✅
- [x] Criar `app/(dashboard)/agendamentos/loading.tsx`. ✅
- [x] Criar `app/(dashboard)/agendamentos/[id]/page.tsx` (detalhes). ✅ *25/11/2025*
- [x] Criar `app/(dashboard)/agendamentos/[id]/loading.tsx`. ✅ *25/11/2025*

### State Management & Data Fetching (Concluído ✅)
- [x] Configurar React Query para `useAppointments`.
- [x] Configurar React Query para `useCreateAppointment` (Mutation).
- [x] Configurar React Query para `useUpdateAppointment` (Mutation).
- [x] Configurar React Query para `useCancelAppointment` (Mutation).
- [x] Configurar React Query para `useCalendarEvents` (conversão para FullCalendar).
- [x] Configurar React Query para `useCalendarResources` (barbeiros como recursos).
- [x] Implementar Optimistic Updates. ✅ *25/11/2025*

### Integração
- [x] Conectar formulário à API `POST /appointments`. ✅ *Via AppointmentModal + useCreateAppointment*
- [x] Conectar calendário à API `GET /appointments`. ✅ *Via useCalendarEvents*
- [x] Tratamento de erros (Toast Notifications via Sonner). ✅

### Componentes Criados (25/11/2025)

**`src/components/appointments/AppointmentModal.tsx`:**
- Modal para criar/editar/visualizar agendamentos
- Formulário com React Hook Form + Zod validation
- Integração com ProfessionalSelector, CustomerSelector, ServiceSelector
- 3 modos: create, edit, view

**`src/components/appointments/AppointmentCard.tsx`:**
- Card para exibir resumo de agendamento
- 2 variantes: default (completo) e compact
- Menu de ações baseado no status
- Status badges com cores

**`src/components/appointments/ServiceSelector.tsx`:**
- Multi-select para serviços
- Busca por nome/categoria
- Mostra preço e duração
- Calcula totais

**`src/components/appointments/ProfessionalSelector.tsx`:**
- Select simples para barbeiros
- Avatar e nome
- Integração com useProfessionals hook

**`src/components/appointments/CustomerSelector.tsx`:**
- Combobox com busca
- Busca por nome, telefone ou email
- Opção de cadastrar novo cliente

**`src/app/(dashboard)/agendamentos/[id]/page.tsx`:**
- Página de detalhes do agendamento
- Cards: Cliente, Barbeiro, Data/Hora, Serviços
- Ações: Editar, Cancelar, Mudar Status
- Histórico de atualizações

### Optimistic Updates Implementados (25/11/2025)

**`src/hooks/use-appointments.ts`:**
- `useUpdateAppointment` — Atualiza lista e detalhe imediatamente, rollback em erro
- `useCancelAppointment` — Marca como CANCELED imediatamente, rollback em erro
- `useUpdateAppointmentStatus` — Transição de status imediata, rollback em erro
- `useCreateAppointment` — Adiciona ao cache após sucesso

**Benefícios:**
- UI responde instantaneamente às ações do usuário
- Rollback automático em caso de erro na API
- Invalidação de queries para garantir consistência
- Experiência de usuário muito mais fluida

## 4. Integrações Externas (MVP)

- [ ] **Google Calendar:**
    - [ ] Estruturar serviço de sync (interface).
    - [ ] Implementar mock inicial.
    - [ ] (Futuro) Implementar chamada real à API do Google.

## 5. Licenças e Conformidade Legal

- [x] **FullCalendar Scheduler:**
    - [x] Validar que a chave `CC-Attribution-NonCommercial-NoDerivatives` está sendo usada apenas em **desenvolvimento**.
    - [ ] **ANTES DO DEPLOY EM PRODUÇÃO:** Substituir pela licença comercial oficial.
    - [ ] Confirmar que a chave comercial foi adicionada ao `.env.production`.
    - [ ] Remover qualquer referência à chave de avaliação do código de produção.

## 6. Definição de Pronto (DoD)

- [x] Código compilando sem erros. ✅ *go build ./... OK + pnpm build OK*
- [x] Lint (ESLint/GolangCI-Lint) passando. ✅
- [x] Testes automatizados passando. ✅ *36 testes (26 unit + 10 integration)*
- [ ] Documentação de API atualizada (Swagger/Markdown).
- [ ] Code Review aprovado.
- [ ] **Licença FullCalendar comercial adquirida e configurada** (para produção).
- [ ] Deploy em ambiente de Staging.

---

**Gerente de Projeto:** Andrey  
**Tech Lead:** Copilot  
**Data de Início:** 25/11/2025  
**Última Atualização:** 25/11/2025 — Optimistic Updates implementados
