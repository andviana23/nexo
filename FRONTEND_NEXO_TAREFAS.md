# 🚀 FRONTEND NEXO — Plano de Implementação MVP v1.0.0

**Data de Criação:** 25/11/2025
**Deadline MVP:** 05/12/2025 (10 dias restantes)
**Stack:** Next.js 16 + Tailwind CSS 4 + shadcn/ui + Framer Motion
**Status:** 🟡 EM ANDAMENTO

---

## 📊 Estado Atual do Frontend

### ✅ Concluído (Infraestrutura)

- [x] Next.js 16.0.4 + React 19.2.0 instalado
- [x] Tailwind CSS 4.1.17 configurado
- [x] 18 componentes shadcn/ui instalados
- [x] Framer Motion 12.23.24 instalado
- [x] Zustand 5.0.8 instalado
- [x] TanStack Query 5.90.11 instalado
- [x] React Hook Form 7.66.1 + Zod 4.1.13 instalados
- [x] Axios 1.13.2 instalado
- [x] CSS Variables Light/Dark configuradas
- [x] Função `cn()` em `src/lib/utils.ts`
- [x] Documentação Design System criada (`docs/03-frontend/`)

### ✅ Concluído (Estrutura Base)

- [x] Estrutura de pastas (`hooks/`, `store/`, `services/`, `types/`)
- [x] `src/lib/axios.ts` (instância configurada)
- [x] `src/lib/query-client.ts` (QueryClientProvider)
- [x] `src/store/auth-store.ts` (Zustand)
- [x] `src/store/ui-store.ts` (sidebar, theme)
- [x] `src/app/layout.tsx` com Providers
- [x] Route Groups (`(auth)/`, `(dashboard)/`)
- [x] `middleware.ts` (proteção de rotas)
- [x] Layout Dashboard (Sidebar, Header)

---

## 🗓️ Cronograma de Execução

### 📅 DIA 1 (25/11) — Fundação + Auth

**Meta:** Estrutura base + Login funcionando

#### FASE 0: Estrutura Base (2h) 🔴 BLOQUEADOR

| ID     | Tarefa                                                                                                                                               | Tempo | Status |
| ------ | ---------------------------------------------------------------------------------------------------------------------------------------------------- | ----- | ------ |
| F0-001 | Criar pastas: `src/hooks/`, `src/store/`, `src/services/`, `src/types/`, `src/components/layout/`, `src/components/shared/`, `src/components/forms/` | 15min | ✅     |
| F0-002 | Criar `src/lib/axios.ts` com interceptors de auth e erro                                                                                             | 30min | ✅     |
| F0-003 | Criar `src/lib/query-client.ts`                                                                                                                      | 15min | ✅     |
| F0-004 | Criar `src/types/index.ts` (User, Tenant, ApiError, PaginatedResponse)                                                                               | 30min | ✅     |
| F0-005 | Criar `src/store/auth-store.ts` (Zustand)                                                                                                            | 30min | ✅     |

#### FASE 1: Auth (4h) 🔴 BLOQUEADOR

| ID     | Tarefa                                                      | Tempo | Status |
| ------ | ----------------------------------------------------------- | ----- | ------ |
| F1-001 | Criar `src/services/auth-service.ts` (login, logout, me)    | 45min | ✅     |
| F1-002 | Criar `src/hooks/use-auth.ts`                               | 30min | ✅     |
| F1-003 | Criar `src/app/(auth)/layout.tsx` (layout centralizado)     | 20min | ✅     |
| F1-004 | Criar `src/app/(auth)/login/page.tsx` (formulário + Zod)    | 1h    | ✅     |
| F1-005 | Criar `src/middleware.ts` (proteção de rotas)               | 30min | ✅     |
| F1-006 | Atualizar `src/app/layout.tsx` com Providers                | 30min | ✅     |
| F1-007 | Criar `src/app/providers.tsx` (QueryClient, Theme, Toaster) | 30min | ✅     |

---

### 📅 DIA 2 (26/11) — Layout Dashboard

**Meta:** Sidebar + Header + Dashboard base

#### FASE 2: Layout (4h) 🔴 CRÍTICO

| ID     | Tarefa                                                                   | Tempo | Status |
| ------ | ------------------------------------------------------------------------ | ----- | ------ |
| F2-001 | Criar `src/store/ui-store.ts` (sidebar open, theme)                      | 20min | ✅     |
| F2-002 | Criar `src/app/(dashboard)/layout.tsx`                                   | 30min | ✅     |
| F2-003 | Criar `src/components/layout/Sidebar.tsx` (navegação, collapse, mobile)  | 2h    | ✅     |
| F2-004 | Criar `src/components/layout/Header.tsx` (user menu, breadcrumb, mobile) | 1h    | ✅     |
| F2-005 | Criar `src/components/layout/UserNav.tsx` (dropdown user)                | 30min | ✅     |

---

### 📅 DIA 3-4 (27-28/11) — Módulo Estoque 🔴 CRÍTICO

**Meta:** Telas de Estoque completas (entrada, saída, inventário)

#### FASE 3: Estoque Frontend (14h)

| ID     | Tarefa                                                          | Tempo | Status |
| ------ | --------------------------------------------------------------- | ----- | ------ |
| F3-001 | Criar `src/types/stock.ts` (StockEntry, StockExit, Inventory)   | 30min | ✅     |
| F3-002 | Criar `src/services/stock-service.ts`                           | 45min | ✅     |
| F3-003 | Criar `src/hooks/use-stock.ts`                                  | 45min | ✅     |
| F3-004 | Criar `src/components/stock/EntryForm.tsx`                      | 1h30  | ✅     |
| F3-005 | Criar `src/app/(dashboard)/estoque/entrada/page.tsx`            | 1h    | ✅     |
| F3-006 | Criar `src/components/stock/ExitForm.tsx`                       | 1h    | ✅     |
| F3-007 | Criar `src/app/(dashboard)/estoque/saida/page.tsx`              | 45min | ✅     |
| F3-008 | Criar `src/components/shared/DataTable.tsx` (REUTILIZÁVEL)      | 2h    | ✅     |
| F3-009 | Criar `src/components/stock/InventoryTable.tsx`                 | 1h30  | ✅     |
| F3-010 | Criar `src/app/(dashboard)/estoque/page.tsx` (inventário)       | 1h30  | ✅     |
| F3-011 | Criar `src/app/(dashboard)/estoque/layout.tsx` (tabs navegação) | 30min | ✅     |

---

### 📅 DIA 5-6 (29-30/11) — Módulo Agendamento 🔴 BLOQUEADOR

**Meta:** Calendário visual + CRUD de agendamentos

#### FASE 4: Agendamento Frontend (18h)

| ID     | Tarefa                                                                 | Tempo | Status |
| ------ | ---------------------------------------------------------------------- | ----- | ------ |
| F4-001 | Instalar biblioteca de calendário (FullCalendar ou react-big-calendar) | 30min | ⬜     |
| F4-002 | Criar `src/types/appointment.ts`                                       | 30min | ⬜     |
| F4-003 | Criar `src/services/appointment-service.ts`                            | 1h    | ⬜     |
| F4-004 | Criar `src/hooks/use-appointments.ts`                                  | 45min | ⬜     |
| F4-005 | Criar `src/components/appointments/Calendar.tsx` (wrapper)             | 3h    | ⬜     |
| F4-006 | Criar `src/components/appointments/AppointmentForm.tsx` (modal)        | 2h    | ⬜     |
| F4-007 | Criar `src/components/appointments/AppointmentCard.tsx`                | 1h    | ⬜     |
| F4-008 | Criar `src/app/(dashboard)/agenda/page.tsx`                            | 2h    | ⬜     |
| F4-009 | Integrar validação de conflitos no form                                | 1h    | ⬜     |
| F4-010 | Criar `src/components/shared/ClientSelect.tsx` (busca async)           | 1h    | ⬜     |
| F4-011 | Criar `src/components/shared/ServiceMultiSelect.tsx`                   | 1h    | ⬜     |

---

### 📅 DIA 7 (01/12) — Lista da Vez 🔴 CRÍTICO

**Meta:** Fila de espera funcional

#### FASE 5: Lista da Vez Frontend (10h)

| ID     | Tarefa                                            | Tempo | Status |
| ------ | ------------------------------------------------- | ----- | ------ |
| F5-001 | Criar `src/types/queue.ts`                        | 20min | ⬜     |
| F5-002 | Criar `src/services/queue-service.ts`             | 45min | ⬜     |
| F5-003 | Criar `src/hooks/use-queue.ts`                    | 30min | ⬜     |
| F5-004 | Criar `src/components/queue/QueueList.tsx`        | 2h    | ⬜     |
| F5-005 | Criar `src/components/queue/QueueCard.tsx`        | 1h    | ⬜     |
| F5-006 | Criar `src/components/queue/AddToQueueForm.tsx`   | 1h30  | ⬜     |
| F5-007 | Criar `src/components/queue/CallNextButton.tsx`   | 30min | ⬜     |
| F5-008 | Criar `src/app/(dashboard)/lista-da-vez/page.tsx` | 2h    | ⬜     |
| F5-009 | Adicionar animações Framer Motion na fila         | 1h    | ⬜     |

---

### 📅 DIA 8 (02/12) — Assinaturas Asaas 🔴 BLOQUEADOR

**Meta:** Checkout + Gerenciamento de assinatura

#### FASE 6: Assinaturas Frontend (11h)

| ID     | Tarefa                                                     | Tempo | Status |
| ------ | ---------------------------------------------------------- | ----- | ------ |
| F6-001 | Criar `src/types/subscription.ts`                          | 30min | ⬜     |
| F6-002 | Criar `src/services/subscription-service.ts`               | 45min | ⬜     |
| F6-003 | Criar `src/hooks/use-subscription.ts`                      | 30min | ⬜     |
| F6-004 | Criar `src/components/subscription/PlanCard.tsx`           | 1h    | ⬜     |
| F6-005 | Criar `src/app/(public)/planos/page.tsx`                   | 1h30  | ⬜     |
| F6-006 | Criar `src/components/subscription/CheckoutForm.tsx`       | 2h    | ⬜     |
| F6-007 | Criar `src/app/(public)/checkout/page.tsx`                 | 1h30  | ⬜     |
| F6-008 | Criar `src/components/subscription/ManageSubscription.tsx` | 1h    | ⬜     |
| F6-009 | Criar `src/app/(dashboard)/assinatura/page.tsx`            | 1h    | ⬜     |

---

### 📅 DIA 9 (03/12) — CRM + Dashboard 🟡 ALTA

**Meta:** Clientes + KPIs Dashboard

#### FASE 7: CRM Frontend (8h)

| ID     | Tarefa                                             | Tempo | Status |
| ------ | -------------------------------------------------- | ----- | ------ |
| F7-001 | Criar `src/types/client.ts`                        | 20min | ⬜     |
| F7-002 | Criar `src/services/client-service.ts`             | 45min | ⬜     |
| F7-003 | Criar `src/hooks/use-clients.ts`                   | 30min | ⬜     |
| F7-004 | Criar `src/components/clients/ClientForm.tsx`      | 1h30  | ⬜     |
| F7-005 | Criar `src/app/(dashboard)/clientes/page.tsx`      | 1h30  | ⬜     |
| F7-006 | Criar `src/app/(dashboard)/clientes/novo/page.tsx` | 30min | ⬜     |
| F7-007 | Criar `src/components/clients/ClientHistory.tsx`   | 1h    | ⬜     |
| F7-008 | Criar `src/app/(dashboard)/clientes/[id]/page.tsx` | 1h    | ⬜     |

#### FASE 8: Dashboard (4h)

| ID     | Tarefa                                                     | Tempo | Status |
| ------ | ---------------------------------------------------------- | ----- | ------ |
| F8-001 | Instalar Recharts                                          | 15min | ⬜     |
| F8-002 | Criar `src/types/dashboard.ts`                             | 20min | ⬜     |
| F8-003 | Criar `src/services/dashboard-service.ts`                  | 30min | ⬜     |
| F8-004 | Criar `src/hooks/use-dashboard.ts`                         | 30min | ⬜     |
| F8-005 | Criar `src/components/shared/MetricCard.tsx`               | 45min | ⬜     |
| F8-006 | Criar `src/app/(dashboard)/page.tsx` (Dashboard principal) | 2h    | ⬜     |

---

### 📅 DIA 10 (04/12) — Relatórios + RBAC 🟡 MÉDIA

**Meta:** DRE, Fluxo de Caixa, Permissões básicas

#### FASE 9: Relatórios Frontend (8h)

| ID     | Tarefa                                                | Tempo | Status |
| ------ | ----------------------------------------------------- | ----- | ------ |
| F9-001 | Criar `src/types/financial.ts`                        | 30min | ⬜     |
| F9-002 | Criar `src/services/financial-service.ts`             | 45min | ⬜     |
| F9-003 | Criar `src/hooks/use-financial.ts`                    | 30min | ⬜     |
| F9-004 | Criar `src/components/financial/DREChart.tsx`         | 1h30  | ⬜     |
| F9-005 | Criar `src/app/(dashboard)/financeiro/dre/page.tsx`   | 1h30  | ⬜     |
| F9-006 | Criar `src/components/financial/CashflowChart.tsx`    | 1h30  | ⬜     |
| F9-007 | Criar `src/app/(dashboard)/financeiro/fluxo/page.tsx` | 1h    | ⬜     |

#### FASE 10: RBAC Frontend (4h)

| ID      | Tarefa                                                    | Tempo | Status |
| ------- | --------------------------------------------------------- | ----- | ------ |
| F10-001 | Criar `src/hooks/use-permissions.ts`                      | 30min | ⬜     |
| F10-002 | Atualizar `Sidebar.tsx` com permissões                    | 1h    | ⬜     |
| F10-003 | Atualizar `middleware.ts` com RBAC                        | 1h    | ⬜     |
| F10-004 | Criar `src/components/shared/PermissionGate.tsx`          | 30min | ⬜     |
| F10-005 | Criar `src/app/(dashboard)/configuracoes/equipe/page.tsx` | 1h    | ⬜     |

---

### 📅 DIA 11 (05/12) — Deploy + Testes 🔴 CRÍTICO

**Meta:** Deploy Staging + Smoke Tests

#### FASE 11: Deploy (4h)

| ID      | Tarefa                                             | Tempo | Status |
| ------- | -------------------------------------------------- | ----- | ------ |
| F11-001 | Configurar variáveis de ambiente (.env.production) | 30min | ⬜     |
| F11-002 | Build de produção (`pnpm build`)                   | 30min | ⬜     |
| F11-003 | Deploy Vercel (staging)                            | 1h    | ⬜     |
| F11-004 | Testar todas as rotas em staging                   | 1h    | ⬜     |
| F11-005 | Corrigir bugs encontrados                          | 1h    | ⬜     |

---

## 📊 Resumo de Esforço por Fase

| Fase      | Descrição        | Horas   | Dia         | Prioridade    |
| --------- | ---------------- | ------- | ----------- | ------------- |
| F0        | Estrutura Base   | 2h      | 25/11       | 🔴 BLOQUEADOR |
| F1        | Autenticação     | 4h      | 25/11       | 🔴 BLOQUEADOR |
| F2        | Layout Dashboard | 4h      | 26/11       | 🔴 CRÍTICO    |
| F3        | Estoque          | 14h     | 27-28/11    | 🔴 CRÍTICO    |
| F4        | Agendamento      | 18h     | 29-30/11    | 🔴 BLOQUEADOR |
| F5        | Lista da Vez     | 10h     | 01/12       | 🔴 CRÍTICO    |
| F6        | Assinaturas      | 11h     | 02/12       | 🔴 BLOQUEADOR |
| F7        | CRM              | 8h      | 03/12       | 🟡 ALTA       |
| F8        | Dashboard        | 4h      | 03/12       | 🟡 ALTA       |
| F9        | Relatórios       | 8h      | 04/12       | 🟡 MÉDIA      |
| F10       | RBAC             | 4h      | 04/12       | 🟡 MÉDIA      |
| F11       | Deploy           | 4h      | 05/12       | 🔴 CRÍTICO    |
| **TOTAL** |                  | **91h** | **11 dias** |               |

---

## 🎯 Componentes shadcn/ui Já Instalados

✅ **Disponíveis para uso imediato:**

- `button`, `input`, `label`, `textarea`, `checkbox`, `select`
- `card`, `dialog`, `sheet`, `dropdown-menu`
- `table`, `badge`, `alert`, `skeleton`
- `form`, `separator`, `avatar`, `sonner`

❌ **Instalar quando necessário:**

```bash
npx shadcn@latest add calendar tabs command popover progress tooltip
```

---

## 🔗 Dependência entre Tarefas

```
F0 (Estrutura) ─── F1 (Auth) ─── F2 (Layout) ───┬── F3 (Estoque)
                                                 ├── F4 (Agendamento)
                                                 ├── F5 (Lista da Vez)
                                                 ├── F6 (Assinaturas)
                                                 ├── F7 (CRM)
                                                 ├── F8 (Dashboard)
                                                 ├── F9 (Relatórios)
                                                 └── F10 (RBAC)
                                                      │
                                                      └── F11 (Deploy)
```

---

## ✅ Checklist de Conclusão MVP Frontend

### Estrutura Base

- [x] Todas as pastas criadas (`hooks/`, `store/`, `services/`, `types/`)
- [x] Axios configurado com interceptors
- [x] React Query configurado (QueryClient + queryKeys)
- [x] Auth Store funcionando (Zustand)
- [x] UI Store criado (sidebar, theme, breadcrumbs)
- [x] Middleware protegendo rotas
- [x] Providers configurados (QueryClient, Theme, Toaster)
- [x] Layout root atualizado com Providers

### Módulos Core

- [x] **Login:** Formulário funcionando + redirect
- [x] **Layout:** Sidebar + Header responsivos
- [x] **Estoque:** Entrada, Saída, Inventário
- [ ] **Agendamento:** Calendário + CRUD
- [ ] **Lista da Vez:** Fila completa
- [ ] **Assinaturas:** Checkout + Gerenciamento
- [ ] **CRM:** Clientes + Histórico
- [ ] **Dashboard:** KPIs + Gráficos
- [ ] **Relatórios:** DRE + Fluxo de Caixa
- [ ] **RBAC:** Menus filtrados por permissão

### Qualidade

- [ ] Todas as telas responsivas (375px → 1920px)
- [ ] Dark Mode funcionando
- [ ] Validação Zod em todos os formulários
- [ ] Estados de Loading em todas as telas
- [ ] Estados de Error tratados
- [ ] Build sem erros (`pnpm build`)
- [ ] Deploy em Staging funcionando

---

## 🚀 Próximos Passos Imediatos

### AGORA — Executar Fase F0:

1. Criar estrutura de pastas
2. Criar `src/lib/axios.ts`
3. Criar `src/lib/query-client.ts`
4. Criar `src/types/index.ts`
5. Criar `src/store/auth-store.ts`

### HOJE (25/11) — Executar Fase F1:

6. Criar auth service + hook
7. Criar layout de auth
8. Criar página de login
9. Configurar middleware
10. Configurar providers

**Meta do dia:** Login funcionando até 23:59

---

## 📝 Log de Progresso

| Data  | Fase | Tarefas Concluídas                          | Observações                                                           |
| ----- | ---- | ------------------------------------------- | --------------------------------------------------------------------- |
| 25/11 | -    | Arquivo de tarefas criado                   | Iniciando implementação                                               |
| 25/11 | F0   | F0-001 a F0-005 (Estrutura Base)            | ✅ FASE F0 COMPLETA                                                   |
| 25/11 | F1   | F1-001 a F1-007 (Auth completa)             | ✅ FASE F1 COMPLETA                                                   |
| 25/11 | F2   | F2-001 a F2-005 (Layout Dashboard completo) | ✅ FASE F2 COMPLETA - Sidebar, Header, UserNav                        |
| 25/11 | F3   | F3-001 a F3-011 (Estoque completo)          | ✅ FASE F3 COMPLETA - Types, Service, Hooks, Pages (inventário total) |

---

**Última Atualização:** 25/11/2025 23:45
**Responsável:** Andrey Viana + GitHub Copilot
**Próxima Revisão:** 26/11/2025 09:00
**Progresso:** 4/11 fases concluídas (36%)

---

**🚀 VAMOS ENTREGAR ESSE FRONTEND! 🚀**
