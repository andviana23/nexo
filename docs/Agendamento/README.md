# 📅 Módulo de Agendamento — NEXO v1.0

**Versão:** 1.0.0  
**Status:** 🟡 Em Desenvolvimento (Milestone 1.5)  
**Prioridade:** 🔴 CRÍTICA  
**Data de Criação:** 25/11/2025  
**Última Atualização:** 25/11/2025  
**Responsável:** Product + Tech Lead  

---

## 📋 Índice da Documentação

Este diretório contém toda a documentação técnica, arquitetural e funcional do **Módulo de Agendamento** do NEXO.

### 📚 Documentos Disponíveis

| Documento | Descrição | Público-Alvo |
|-----------|-----------|--------------|
| **[PRD_AGENDAMENTO.md](./PRD_AGENDAMENTO.md)** | Product Requirements Document completo | Product Manager, Stakeholders |
| **[ARQUITETURA_AGENDAMENTO.md](./ARQUITETURA_AGENDAMENTO.md)** | Arquitetura completa (Frontend + Backend) | Tech Lead, Desenvolvedores |
| **[BANCO_AGENDAMENTO.md](./BANCO_AGENDAMENTO.md)** | Schema de banco de dados completo | DBA, Backend Devs |
| **[API_AGENDAMENTO.md](./API_AGENDAMENTO.md)** | Contrato completo da API REST | Frontend + Backend Devs |
| **[DIAGRAMAS_AGENDAMENTO.md](./DIAGRAMAS_AGENDAMENTO.md)** | Fluxogramas e diagramas técnicos | Todos os times |
| **[CHECKLIST_IMPLEMENTACAO.md](./CHECKLIST_IMPLEMENTACAO.md)** | Checklist completo de implementação | Todos os desenvolvedores |

---

## 🎯 Visão Geral do Módulo

### Objetivo

Permitir o **agendamento visual e intuitivo** de serviços de barbearia, com:

- ✅ Calendário visual (estilo AppBarber/Trinks)
- ✅ Validação de conflitos em tempo real
- ✅ Isolamento multi-tenant completo
- ✅ Sincronização com Google Agenda
- ✅ CRUD completo (criar, editar, cancelar, reagendar)
- ✅ Status lifecycle (CREATED → CONFIRMED → IN_SERVICE → DONE)
- ✅ Controle de permissões por role (RBAC)

---

## 🏗️ Stack Tecnológica

### Frontend
- **Framework:** Next.js 15.5.6 (App Router)
- **React:** 19.2.0
- **UI:** Tailwind CSS 4.1.17 + shadcn/ui
- **Calendário:** FullCalendar 6.x (ResourceTimeGrid) ⚠️ **Licença Avaliação**
- **State Management:** TanStack Query 5.90.11
- **Forms:** React Hook Form 7.66.1 + Zod 4.1.13

### ⚠️ Atenção: Licença FullCalendar Scheduler – Modo Avaliação

O NEXO utiliza o **FullCalendar Premium (Scheduler)** durante o período de **avaliação gratuita** com a seguinte chave:

```javascript
schedulerLicenseKey: 'CC-Attribution-NonCommercial-NoDerivatives'
```

**Restrições Legais:**

- ❌ **Proibido uso comercial** neste modo.
- ✅ **Permitido apenas** para desenvolvimento interno, testes e homologação.
- ⚠️ **Antes do lançamento em produção**, será necessário **adquirir a licença comercial oficial** do FullCalendar.

**Documentação oficial:** [FullCalendar Scheduler License](https://fullcalendar.io/docs/schedulerLicenseKey)

### Backend
- **Linguagem:** Go 1.24
- **Framework:** Echo v4
- **Arquitetura:** Clean Architecture + DDD
- **Banco:** PostgreSQL 14+ (Neon Cloud)
- **ORM:** sqlc v1.30.0

---

## 📊 Fluxo Simplificado

```
[Usuário Autenticado]
    ↓
[Acessar Tela de Agendamentos]
    ↓
[Visualizar Calendário com Barbeiros]
    ↓
[Clicar em Novo Agendamento]
    ↓
[Selecionar: Cliente + Serviço(s) + Barbeiro + Data/Hora]
    ↓
[Validar Disponibilidade (Backend)]
    ↓
[Criar Agendamento]
    ↓
[Sincronizar Google Agenda (Async)]
    ↓
[Notificar Cliente (Futuro)]
```

---

## 🔐 Regras de Negócio Críticas

| ID | Regra | Descrição |
|----|-------|-----------|
| **RN-AGE-001** | Validação de Barbeiro | Barbeiro deve estar ativo e pertencer ao mesmo tenant |
| **RN-AGE-002** | Validação de Cliente | Cliente deve existir antes de agendar |
| **RN-AGE-003** | Intervalo Mínimo | 10 minutos entre agendamentos (configurável) |
| **RN-AGE-004** | Multi-Tenant | TODOS os dados isolados por `tenant_id` |
| **RN-AGE-005** | Status Lifecycle | CREATED → CONFIRMED → IN_SERVICE → DONE/CANCELED/NO_SHOW |
| **RN-AGE-006** | Conflitos | Sistema DEVE impedir conflitos de horário |
| **RN-AGE-007** | Google Sync | Sincronizar apenas agendamentos CONFIRMED |

---

## 👥 Personas e Permissões

| Persona | Ver Agenda | Criar | Editar | Cancelar | Reagendar | Visualizar Todos |
|---------|------------|-------|--------|----------|-----------|------------------|
| **Dono** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Gerente** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Recepção** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ (unidade) |
| **Barbeiro** | ✅ (própria) | ❌ | ❌ | ❌ | ❌ | ❌ |
| **Cliente** | ✅ (próprios) | ✅ (app) | ❌ | ✅ (app) | ✅ (app) | ❌ |

---

## 📦 Estrutura de Pastas

### Frontend
```
frontend/src/
├── app/
│   └── (dashboard)/
│       └── agendamentos/
│           ├── page.tsx              # Página principal (calendário)
│           ├── novo/
│           │   └── page.tsx          # Formulário de novo agendamento
│           └── [id]/
│               ├── page.tsx          # Detalhes do agendamento
│               └── editar/
│                   └── page.tsx      # Editar agendamento
├── components/
│   └── appointments/
│       ├── AppointmentCalendar.tsx   # Wrapper FullCalendar
│       ├── AppointmentForm.tsx       # Formulário (criar/editar)
│       ├── AppointmentCard.tsx       # Card de agendamento
│       ├── AppointmentModal.tsx      # Modal de detalhes
│       ├── ProfessionalSelect.tsx    # Seletor de barbeiro
│       └── ServiceMultiSelect.tsx    # Multi-select de serviços
├── services/
│   └── appointment-service.ts        # API calls
├── hooks/
│   └── use-appointments.ts           # React Query hooks
└── types/
    └── appointment.ts                # TypeScript types
```

### Backend
```
backend/
├── internal/
│   ├── domain/
│   │   └── appointment/
│   │       ├── appointment.go        # Entity
│   │       ├── repository.go         # Interface
│   │       └── value_objects.go      # VOs (Status, TimeSlot)
│   ├── application/
│   │   ├── dto/
│   │   │   └── appointment_dto.go    # DTOs
│   │   ├── mapper/
│   │   │   └── appointment_mapper.go # Mappers
│   │   └── usecase/
│   │       └── appointment/
│   │           ├── create_appointment.go
│   │           ├── update_appointment.go
│   │           ├── cancel_appointment.go
│   │           ├── check_availability.go
│   │           └── list_appointments.go
│   └── infrastructure/
│       ├── repository/
│       │   └── postgres/
│       │       └── appointment_repository.go
│       ├── http/
│       │   └── handler/
│       │       └── appointment_handler.go
│       └── external/
│           └── google/
│               └── calendar_service.go
└── migrations/
    ├── 00XX_create_appointments_table.up.sql
    ├── 00XX_create_appointments_table.down.sql
    ├── 00XX_create_appointment_services_table.up.sql
    └── 00XX_create_appointment_services_table.down.sql
```

---

## 🚀 Início Rápido

### 1. Leia a Documentação na Ordem

1. **[PRD_AGENDAMENTO.md](./PRD_AGENDAMENTO.md)** - Entenda O QUE será feito
2. **[ARQUITETURA_AGENDAMENTO.md](./ARQUITETURA_AGENDAMENTO.md)** - Entenda COMO será feito
3. **[BANCO_AGENDAMENTO.md](./BANCO_AGENDAMENTO.md)** - Entenda o SCHEMA de dados
4. **[API_AGENDAMENTO.md](./API_AGENDAMENTO.md)** - Entenda os ENDPOINTS
5. **[DIAGRAMAS_AGENDAMENTO.md](./DIAGRAMAS_AGENDAMENTO.md)** - Visualize os FLUXOS
6. **[CHECKLIST_IMPLEMENTACAO.md](./CHECKLIST_IMPLEMENTACAO.md)** - Execute a implementação

### 2. Requisitos Técnicos

**Backend:**
- Go 1.24+
- PostgreSQL 14+
- Echo v4
- sqlc v1.30.0

**Frontend:**
- Node.js 20+
- pnpm 10+
- Next.js 15.5.6+
- FullCalendar 6.x

### 3. Comandos Úteis

```bash
# Backend
cd backend
make migrate-up        # Rodar migrations
make dev              # Iniciar servidor dev
make test             # Rodar testes

# Frontend
cd frontend
pnpm install          # Instalar dependências
pnpm dev             # Iniciar dev server
pnpm build           # Build de produção
```

---

## 📊 Métricas de Sucesso

### Operacionais
- Taxa de no-show < 10%
- Tempo médio de criação de agendamento < 30 segundos
- Conflitos de horário = 0

### Técnicas
- Latência da API < 150ms
- Uptime > 99.5%
- Sincronização Google Calendar < 500ms

---

## 🔗 Referências Externas

- [FullCalendar Docs](https://fullcalendar.io/docs)
- [Google Calendar API](https://developers.google.com/calendar/api)
- [Clean Architecture (Go)](https://github.com/bxcodec/go-clean-arch)
- [Next.js App Router](https://nextjs.org/docs/app)

---

## 📞 Suporte

- **Tech Lead:** Andrey Viana
- **Product Owner:** Andrey Viana
- **Documentação:** Este diretório (`docs/Agendamento/`)

---

**🚀 Vamos construir o melhor sistema de agendamento para barbearias do Brasil! 🚀**
