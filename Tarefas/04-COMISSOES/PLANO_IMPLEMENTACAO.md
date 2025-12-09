# 📊 PLANO DE IMPLEMENTAÇÃO — MÓDULO DE COMISSÕES

> **Versão:** 1.0.0  
> **Data:** Dezembro 2024  
> **Status:** NÃO INICIADO  
> **Sprints Alvo:** 15-17  
> **Dependências:** ✅ Pacote 03-FINANCEIRO (Sprint 1)

---

## 📋 ÍNDICE

1. [Resumo Executivo](#resumo-executivo)
2. [Status Global](#status-global)
3. [Fases de Implementação](#fases-de-implementação)
4. [Sprint 1: Migrations + Queries](#sprint-1-migrations--queries)
5. [Sprint 2: Domain + Repository + UseCases](#sprint-2-domain--repository--usecases)
6. [Sprint 3: Handlers + Motor de Cálculo](#sprint-3-handlers--motor-de-cálculo)
7. [Sprint 4: Frontend Config + Fechamento](#sprint-4-frontend-config--fechamento)
8. [Sprint 5: Frontend Dashboard Barbeiro](#sprint-5-frontend-dashboard-barbeiro)
9. [Sprint 6: Testes E2E + QA](#sprint-6-testes-e2e--qa)
10. [Dependências Críticas](#dependências-críticas)
11. [Riscos e Mitigações](#riscos-e-mitigações)

---

## 📌 RESUMO EXECUTIVO

O Módulo de Comissões automatiza todo o ciclo de pagamento de profissionais:

- Cálculo automático baseado em regras flexíveis
- Fechamento de período com consolidação
- Integração com Contas a Pagar
- Dashboard individual do barbeiro
- Gestão de adiantamentos

### Progresso Atual

```
░░░░░░░░░░░░░░░░░░░░ 0% Completo
```

| Componente | Backend | Frontend | Testes | Docs |
|------------|:-------:|:--------:|:------:|:----:|
| commission_rules | ❌ | ❌ | ❌ | ✅ |
| commission_periods | ❌ | ❌ | ❌ | ✅ |
| advances | ❌ | ❌ | ❌ | ✅ |
| Motor de Cálculo | ❌ | N/A | ❌ | ✅ |
| Fechamento | ❌ | ❌ | ❌ | ✅ |
| Dashboard Barbeiro | N/A | ❌ | ❌ | ✅ |

---

## 📊 STATUS GLOBAL

### ✅ EXISTENTE (Aproveitável)

| Item | Localização | Status |
|------|-------------|--------|
| Tabela `profissionais` | `migrations/003_full_schema.sql` | ✅ Tem comissao + tipo_comissao |
| Tabela `servicos` | `migrations/003_full_schema.sql` | ✅ Tem comissao |
| Tabela `barber_commissions` | `migrations/003_full_schema.sql` | ✅ Precisa ajuste |
| Tabela `contas_a_pagar` | `migrations/003_full_schema.sql` | ✅ Pronto |
| Tabela `dre_mensal` | `migrations/003_full_schema.sql` | ✅ Tem custo_comissoes |
| Tabela `metas_mensais` | `migrations/003_full_schema.sql` | ✅ Para bônus |

### ❌ PENDENTE (Bloqueia MVP)

| Item | Prioridade | Sprint | Esforço |
|------|:----------:|:------:|:-------:|
| Migration `commission_rules` | 🔴 P0 | Sprint 1 | 3h |
| Migration `commission_periods` | 🔴 P0 | Sprint 1 | 3h |
| Migration `advances` | 🔴 P0 | Sprint 1 | 2h |
| Alter `barber_commissions` | 🔴 P0 | Sprint 1 | 1h |
| Queries sqlc (4 tabelas) | 🔴 P0 | Sprint 1 | 6h |
| Domain Entities | 🔴 P0 | Sprint 2 | 4h |
| Repositories | 🔴 P0 | Sprint 2 | 6h |
| UseCases (CRUD + Cálculo) | 🔴 P0 | Sprint 2 | 10h |
| Motor de Cálculo | 🔴 P0 | Sprint 3 | 8h |
| Handlers API | 🟡 P1 | Sprint 3 | 6h |
| Tela Config Regras | 🟡 P1 | Sprint 4 | 8h |
| Tela Fechamento | 🟡 P1 | Sprint 4 | 10h |
| Dashboard Barbeiro | 🟡 P1 | Sprint 5 | 12h |
| Tela Adiantamentos | 🟡 P1 | Sprint 5 | 6h |
| Testes E2E | 🟢 P2 | Sprint 6 | 8h |

---

## 🏗 FASES DE IMPLEMENTAÇÃO

```
┌─────────────────────────────────────────────────────────────────────┐
│                     ROADMAP COMISSÕES                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  SPRINT 1 ❌        SPRINT 2 ❌        SPRINT 3 ❌                   │
│  ┌───────────┐     ┌───────────┐      ┌───────────┐                 │
│  │ Migrations│────▶│ Domain    │─────▶│ Handlers  │                 │
│  │ + Queries │     │ Repository│      │ Motor     │                 │
│  │ sqlc gen  │     │ UseCases  │      │ Cálculo   │                 │
│  └───────────┘     └───────────┘      └───────────┘                 │
│       │                 │                  │                         │
│       ▼                 ▼                  ▼                         │
│  SPRINT 4 ❌        SPRINT 5 ❌        SPRINT 6 ❌                   │
│  ┌───────────┐     ┌───────────┐      ┌───────────┐                 │
│  │ Frontend  │────▶│ Dashboard │─────▶│ Testes E2E│                 │
│  │ Config    │     │ Barbeiro  │      │ QA Final  │                 │
│  │ Fechamento│     │ Adianta.  │      │ Deploy    │                 │
│  └───────────┘     └───────────┘      └───────────┘                 │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 🔵 SPRINT 1: MIGRATIONS + QUERIES

### 1.1 Overview

**Objetivo:** Criar toda a estrutura de banco de dados  
**Duração:** 1 semana  
**Esforço:** ~15 horas  
**Checklist:** `CHECKLIST_SPRINT1_MIGRATIONS.md`

### 1.2 Entregas

| Entrega | Descrição |
|---------|-----------|
| Migration `XXX_commission_rules` | Tabela de regras flexíveis |
| Migration `XXX_commission_periods` | Tabela de períodos/folhas |
| Migration `XXX_advances` | Tabela de adiantamentos |
| Migration `XXX_alter_barber_commissions` | Adicionar `command_item_id` |
| Queries sqlc | CRUD completo para todas as tabelas |
| `sqlc generate` | Gerar código Go |

### 1.3 Dependências

- ✅ Migration 003 (schema base)
- ✅ Tabela `profissionais`
- ✅ Tabela `servicos`
- ✅ Tabela `barber_commissions`
- ✅ Tabela `contas_a_pagar`

---

## 🟢 SPRINT 2: DOMAIN + REPOSITORY + USECASES

### 2.1 Overview

**Objetivo:** Criar toda a camada de domínio e aplicação  
**Duração:** 1 semana  
**Esforço:** ~20 horas  
**Checklist:** `CHECKLIST_SPRINT2_BACKEND.md`

### 2.2 Entregas

| Entrega | Descrição |
|---------|-----------|
| Entity `CommissionRule` | Regra de comissão |
| Entity `CommissionPeriod` | Período/Folha |
| Entity `Advance` | Adiantamento |
| Value Objects | Enums de status e tipos |
| Repository Interfaces | Contratos |
| Repository Implementations | PostgreSQL |
| UseCases CRUD | Create, Get, List, Update, Delete |
| UseCases Específicos | CalculateCommission, ClosePeriod |

### 2.3 Dependências

- ✅ Sprint 1 completo
- ✅ Queries sqlc geradas

---

## 🟡 SPRINT 3: HANDLERS + MOTOR DE CÁLCULO

### 3.1 Overview

**Objetivo:** Expor API e implementar motor de cálculo  
**Duração:** 1 semana  
**Esforço:** ~14 horas  
**Checklist:** `CHECKLIST_SPRINT3_HANDLERS.md`

### 3.2 Entregas

| Entrega | Descrição |
|---------|-----------|
| CommissionRulesHandler | CRUD de regras |
| CommissionsHandler | Listagem e resumo |
| CommissionPeriodsHandler | Preview e fechamento |
| AdvancesHandler | CRUD + aprovação |
| Motor de Cálculo | Trigger no fechamento de comanda |
| Integração Contas a Pagar | Geração automática |

### 3.3 Dependências

- ✅ Sprint 2 completo
- ✅ UseCases implementados

---

## 🟠 SPRINT 4: FRONTEND CONFIG + FECHAMENTO

### 4.1 Overview

**Objetivo:** Telas de configuração e fechamento  
**Duração:** 1 semana  
**Esforço:** ~18 horas  
**Checklist:** `CHECKLIST_SPRINT4_FRONTEND_CONFIG.md`

### 4.2 Entregas

| Entrega | Descrição |
|---------|-----------|
| Página `/admin/comissoes/config` | Config global + por serviço |
| Página `/financeiro/comissoes` | Fechamento de período |
| Componente RegrasComissaoForm | Form de regras |
| Componente FechamentoTable | Tabela de fechamento |
| Componente PreviewModal | Prévia antes de fechar |
| Services/Hooks | Integração com API |

### 4.3 Dependências

- ✅ Sprint 3 completo
- ✅ API disponível

---

## 🟣 SPRINT 5: FRONTEND DASHBOARD BARBEIRO

### 5.1 Overview

**Objetivo:** Dashboard individual e adiantamentos  
**Duração:** 1 semana  
**Esforço:** ~18 horas  
**Checklist:** `CHECKLIST_SPRINT5_FRONTEND_DASHBOARD.md`

### 5.2 Entregas

| Entrega | Descrição |
|---------|-----------|
| Página `/barbeiro/painel` | Dashboard individual |
| Página `/financeiro/adiantamentos` | Gestão de vales |
| Componente ComissaoCard | Card resumo |
| Componente ComissaoChart | Gráfico evolução |
| Componente ExtratoList | Lista de atendimentos |
| Componente AdiantamentoForm | Solicitação |

### 5.3 Dependências

- ✅ Sprint 4 completo
- ✅ RBAC configurado (barbeiro só vê seus dados)

---

## ⚫ SPRINT 6: TESTES E2E + QA

### 6.1 Overview

**Objetivo:** Garantir qualidade e estabilidade  
**Duração:** 1 semana  
**Esforço:** ~10 horas  
**Checklist:** `CHECKLIST_SPRINT6_TESTES.md`

### 6.2 Entregas

| Entrega | Descrição |
|---------|-----------|
| Testes unitários | Motor de cálculo |
| Testes integração | Fechamento + Contas a Pagar |
| Testes E2E | Fluxo completo |
| Testes RBAC | Isolamento barbeiro |
| Smoke tests | Regressão |
| Documentação | Atualizar docs |

---

## ⚠️ DEPENDÊNCIAS CRÍTICAS

### Internas

| Dependência | Status | Impacto |
|-------------|--------|---------|
| `contas_a_pagar` | ✅ Pronto | Bloqueia fechamento |
| `dre_mensal` | ✅ Pronto | Bloqueia relatório |
| `profissionais` | ✅ Pronto | Bloqueia cálculo |
| `servicos` | ✅ Pronto | Bloqueia cálculo |
| `commands` | ✅ Pronto | Bloqueia trigger |
| RBAC barbeiro | ✅ Pronto | Bloqueia dashboard |

### Externas

| Dependência | Status | Impacto |
|-------------|--------|---------|
| sqlc instalado | ✅ | Bloqueia queries |
| Node.js/pnpm | ✅ | Bloqueia frontend |
| Design System | ✅ | Bloqueia telas |

---

## 🚨 RISCOS E MITIGAÇÕES

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Erro de cálculo complexo | Média | Alto | Testes extensivos + log |
| Performance no fechamento | Baixa | Médio | Batch processing |
| Conflito de regras | Baixa | Médio | Prioridade explícita |
| RBAC incorreto | Média | Alto | Testes de segurança |
| Integração DRE falhar | Baixa | Médio | Transaction rollback |

---

## 📁 ESTRUTURA DE ARQUIVOS (Previsão)

```
backend/
├── migrations/
│   ├── XXX_commission_rules.up.sql
│   ├── XXX_commission_rules.down.sql
│   ├── XXX_commission_periods.up.sql
│   ├── XXX_commission_periods.down.sql
│   ├── XXX_advances.up.sql
│   ├── XXX_advances.down.sql
│   └── XXX_alter_barber_commissions.up.sql
│
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── commission_rule.go
│   │   │   ├── commission_period.go
│   │   │   └── advance.go
│   │   └── repository/
│   │       ├── commission_rule_repository.go
│   │       ├── commission_period_repository.go
│   │       └── advance_repository.go
│   │
│   ├── application/
│   │   ├── dto/
│   │   │   ├── commission_rule_dto.go
│   │   │   ├── commission_period_dto.go
│   │   │   └── advance_dto.go
│   │   └── usecase/
│   │       ├── commission/
│   │       │   ├── calculate_commission.go
│   │       │   ├── create_rule.go
│   │       │   └── close_period.go
│   │       └── advance/
│   │           ├── create_advance.go
│   │           └── approve_advance.go
│   │
│   ├── infra/
│   │   ├── db/
│   │   │   └── queries/
│   │   │       ├── commission_rules.sql
│   │   │       ├── commission_periods.sql
│   │   │       └── advances.sql
│   │   └── repository/
│   │       ├── pg_commission_rule_repository.go
│   │       ├── pg_commission_period_repository.go
│   │       └── pg_advance_repository.go
│   │
│   └── interfaces/
│       └── http/
│           └── handler/
│               ├── commission_rules_handler.go
│               ├── commissions_handler.go
│               ├── commission_periods_handler.go
│               └── advances_handler.go

frontend/
└── src/
    ├── app/
    │   ├── (authenticated)/
    │   │   ├── admin/
    │   │   │   └── comissoes/
    │   │   │       └── config/
    │   │   │           └── page.tsx
    │   │   ├── financeiro/
    │   │   │   ├── comissoes/
    │   │   │   │   └── page.tsx
    │   │   │   └── adiantamentos/
    │   │   │       └── page.tsx
    │   │   └── barbeiro/
    │   │       └── painel/
    │   │           └── page.tsx
    │
    ├── components/
    │   └── comissoes/
    │       ├── RegrasComissaoForm.tsx
    │       ├── FechamentoTable.tsx
    │       ├── PreviewModal.tsx
    │       ├── ComissaoCard.tsx
    │       ├── ComissaoChart.tsx
    │       └── ExtratoList.tsx
    │
    ├── services/
    │   ├── commissionRulesService.ts
    │   ├── commissionsService.ts
    │   ├── commissionPeriodsService.ts
    │   └── advancesService.ts
    │
    └── hooks/
        ├── useCommissionRules.ts
        ├── useCommissions.ts
        ├── useCommissionPeriods.ts
        └── useAdvances.ts
```

---

*Documento criado em: 05/12/2025*
