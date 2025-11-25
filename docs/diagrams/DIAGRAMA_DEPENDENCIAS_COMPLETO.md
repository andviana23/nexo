> Criado em: 21/11/2025 19:00 (America/Sao_Paulo)

# 🔍 Diagrama de Dependências Completo - Barber Analytics Pro v2.0

**Data:** 22/11/2025  
**Tipo:** Auditoria Arquitetural (estado atual vs planejado)  
**Autor:** Auditor Técnico / Arquiteto Sênior

---

## 🗺️ Diagrama Mestre de Dependências (Estado Atual)

```mermaid
flowchart TB
    subgraph EXTERNAL["🌐 SISTEMAS EXTERNOS"]
        NEON[Neon PostgreSQL<br/>Banco Principal]
    end

    subgraph INFRA["⚙️ INFRAESTRUTURA"]
        NGINX[NGINX<br/>Reverse Proxy]
        PROMETHEUS[Prometheus<br/>Métricas]
        GRAFANA[Grafana<br/>Dashboards]
        SENTRY[Sentry<br/>APM/Errors]
    end

    subgraph BACKEND["🔧 BACKEND GO"]
        direction TB

        subgraph HTTP_LAYER["📡 HTTP LAYER"]
            ECHO[Echo v4]
            LOGGER_MW[Logger Middleware]
            RECOVERY_MW[Recovery Middleware]
            CORS_MW[CORS Middleware]
            TENANT_MW[Tenant Mock Middleware<br/>Header X-Tenant-ID]
            HANDLERS[Handlers HTTP<br/>financeiro / metas / precificação]
            ROUTES["/api/v1/*"]
        end

        subgraph APP_LAYER["🎯 APPLICATION LAYER"]
            DTOS[DTOs Request/Response]
            MAPPERS[Mappers]
            USECASES[Use Cases]
            UC_FINANCIAL[Financial Use Cases<br/>payables/receivables/fluxo/DRE]
            UC_METAS[Metas Use Cases]
            UC_PRICING[Pricing Use Cases]
            UC_LGPD[User Preferences Use Cases]
        end

        subgraph DOMAIN_LAYER["💎 DOMAIN LAYER"]
            ENTITIES[Entities/Aggregates<br/>ContaPagar, ContaReceber,<br/>FluxoCaixa, DRE, Metas,<br/>Precificação, UserPreferences]
            VALUE_OBJECTS[Value Objects<br/>Money, Percentual, MesAno, etc]
            PORTS[Repository Interfaces]
        end

        subgraph INFRASTRUCTURE_LAYER["🏗️ INFRASTRUCTURE LAYER"]
            REPOSITORIES[PostgreSQL Repositories<br/>SQLC]
            REPO_FINANCIAL[ContaPagar/Receber/Fluxo/DRE/Compensação]
            REPO_METAS[Metas Mensais/Barbeiro/Ticket]
            REPO_PRICING[Precificação Config/Simulação]
            REPO_PREFS[User Preferences]

            SCHEDULER[Cron Scheduler<br/>robfig/cron/v3]
            CRON_DRE[GenerateDREMonthly]
            CRON_FLUXO[GenerateFluxoDiario]
            CRON_COMP[MarcarCompensacoes]
        end
    end

    subgraph FRONTEND["🎨 FRONTEND"]
        NEXTJS[Next.js App Router]
        TANSTACK[TanStack Query]
        MUI_SHADCN[MUI + shadcn/ui]
        HOOKS[Hooks/Services TS<br/>financeiro/metas/precificação]
    end

    subgraph PERSISTENCE["💾 PERSISTENCE"]
        DB_PAYABLES[(contas_a_pagar)]
        DB_RECEIVABLES[(contas_a_receber)]
        DB_COMP[(compensacoes_bancarias)]
        DB_FLUXO[(fluxo_caixa_diario)]
        DB_DRE[(dre_mensal)]
        DB_METAS[(metas_mensais<br/>metas_barbeiro<br/>metas_ticket_medio)]
        DB_PRICING[(precificacao_config<br/>precificacao_simulacoes)]
        DB_PREFS[(user_preferences)]
    end

    %% Connections
    NGINX --> ECHO
    PROMETHEUS -.-> ECHO
    GRAFANA -.-> PROMETHEUS
    SENTRY -.-> ECHO

    ECHO --> LOGGER_MW --> RECOVERY_MW --> CORS_MW --> TENANT_MW --> HANDLERS --> ROUTES
    ROUTES --> USECASES
    HANDLERS --> DTOS --> MAPPERS --> USECASES
    USECASES --> UC_FINANCIAL
    USECASES --> UC_METAS
    USECASES --> UC_PRICING
    USECASES --> UC_LGPD

    UC_FINANCIAL --> ENTITIES
    UC_METAS --> ENTITIES
    UC_PRICING --> ENTITIES
    UC_LGPD --> ENTITIES
    ENTITIES --> VALUE_OBJECTS --> PORTS

    PORTS -.-> REPOSITORIES
    REPOSITORIES --> REPO_FINANCIAL
    REPOSITORIES --> REPO_METAS
    REPOSITORIES --> REPO_PRICING
    REPOSITORIES --> REPO_PREFS

    REPO_FINANCIAL --> DB_PAYABLES
    REPO_FINANCIAL --> DB_RECEIVABLES
    REPO_FINANCIAL --> DB_COMP
    REPO_FINANCIAL --> DB_FLUXO
    REPO_FINANCIAL --> DB_DRE
    REPO_METAS --> DB_METAS
    REPO_PRICING --> DB_PRICING
    REPO_PREFS --> DB_PREFS

    SCHEDULER --> CRON_DRE --> UC_FINANCIAL
    SCHEDULER --> CRON_FLUXO --> UC_FINANCIAL
    SCHEDULER --> CRON_COMP --> UC_FINANCIAL

    FRONTEND -.->|HTTP| ROUTES
    FRONTEND --> HOOKS
    HOOKS -.->|REST| ROUTES
```

---

## ⚠️ Gaps e Alertas (Estado Atual)

- **Auth/RBAC ausente:** TENANT_MW usa header mock; falta JWT RS256 e roles.
- **Validator não registrado:** handlers chamam `c.Validate`, mas o Echo não tem validator global configurado.
- **Repos financeiros incompletos:** `SumByPeriod` e agregações retornam zero, impactando Fluxo/DRE.
- **LGPD parcial:** `user_preferences` tem repo, mas handlers `/me/*` não expostos/completos.
- **Futuros não implementados:** Assinaturas/Asaas, Agenda/Lista da Vez, Comissões, Estoque, CRM.
- **RLS/auditoria:** sem RLS no Postgres; sem audit logs.

---

## 🧭 Diagrama Planejado (quando módulos forem adicionados)

```mermaid
flowchart TB
    subgraph EXTERNAL["🌐 EXTERNOS"]
        ASAAS[Asaas API v3]
        GOOGLE[Google Calendar]
        NEON[Neon PostgreSQL]
    end

    subgraph BACKEND["🔧 BACKEND"]
        AUTH_MW[Auth JWT RS256 + RBAC]
        TENANT_MW[Tenant Middleware]
        LGPD[LGPD Handlers]
        AGENDA[Agenda/Lista da Vez UCs]
        SUBS[Assinaturas/Asaas UCs]
        COMISSOES[Comissões UCs]
        ESTOQUE[Estoque UCs]
        CRM[CRM UCs]
        AUDIT[Audit Log]
        RLS[RLS Policies]
        SCHED[Scheduler Jobs]
    end

    subgraph PERSISTENCE["💾 DB FUTURO"]
        DB_AGENDA[(agendamentos<br/>turns)]
        DB_ASAAS[(assinaturas<br/>faturas)]
        DB_ESTOQUE[(produtos<br/>movimentacoes)]
        DB_CRM[(clientes<br/>contatos)]
        DB_AUDIT[(audit_logs)]
    end

    SUBS --> ASAAS
    AGENDA --> GOOGLE
    RLS -.-> DB_AGENDA
    RLS -.-> DB_ASAAS
    RLS -.-> DB_ESTOQUE
    RLS -.-> DB_CRM
    AUDIT --> DB_AUDIT
```

---

## ✅ Próximos Passos Sugeridos

1. **Auth/RBAC + Tenant real:** adicionar middleware JWT RS256 e popular `tenant_id` a partir do token.
2. **Validator:** registrar `validator/v10` no Echo antes dos handlers.
3. **Repos financeiros:** implementar agregações/filtros em `conta_pagar/receber` para fluxos/DRE corretos.
4. **LGPD:** expor rotas `/me/preferences|export|delete` e completar use cases.
5. **Planejar módulos futuros:** agenda/lista da vez, assinaturas/Asaas, comissões, estoque, CRM — alinhar contratos e schemas.
6. **RLS/Auditoria:** ativar RLS por tabela e logar `tenant_id`/`user_id` em operações sensíveis.
```
