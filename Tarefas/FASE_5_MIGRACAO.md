# 🟦 FASE 5 — Preparação para Produção (V2 Standalone)

**Objetivo:** Preparar V2 para rodar em produção de forma independente (sem MVP)
**Duração:** 7-14 dias
**Dependências:** ✅ Fase 3 + Fase 4 completas
**Sprint:** Sprint 7-8

---

## 📊 Progresso Geral

```
┌─────────────────────────────────────────────────────────────┐
│  FASE 5: PREPARAÇÃO PARA PRODUÇÃO V2                        │
├─────────────────────────────────────────────────────────────┤
│  Progresso:  ████████████████████████  100% (4/4 concluídas)│
│  Status:     ✅ CONCLUÍDO                                   │
│  Prioridade: 🔴 ALTA                                        │
│  Estimativa: 16 horas (8h gastas)                          │
│  Sprint:     Sprint 7-8                                     │
│  Abordagem:  🆕 V2 STANDALONE (sem migração de dados)      │
└─────────────────────────────────────────────────────────────┘
```

---

## ⚠️ **IMPORTANTE: ESTRATÉGIA SEM MIGRAÇÃO**

**Decisão:** O V2 **NÃO** migrará dados do MVP 1.0.

- ✅ V2 inicia com banco de dados limpo (apenas estrutura)
- ✅ Novos clientes começam direto no V2
- ✅ Clientes existentes continuam no MVP 1.0 (ou migram manualmente se desejarem)
- ❌ Sem dual-read (MVP + V2 ao mesmo tempo)
- ❌ Sem scripts de migração automática de dados

---

## ✅ Checklist de Tarefas

### ✅ T-PROD-001 — Seed de Dados Iniciais

- **Responsável:** Backend / DevOps
- **Prioridade:** 🔴 Alta
- **Estimativa:** 4h
- **Sprint:** Sprint 7
- **Status:** ✅ Concluído (17/11/2025)
- **Horas Gastas:** 4h
- **Deliverable:** Scripts de seed para dados essenciais do sistema

#### Critérios de Aceitação

- [x] Script `seed_categories.sql` - Categorias padrão de receitas e despesas
  - [x] Categorias de Receita: Serviços, Produtos, Assinaturas, Outros
  - [x] Categorias de Despesa: Salários, Aluguel, Fornecedores, Impostos, Marketing, Outros
- [x] Script `seed_plans.sql` - Planos de assinatura padrão (Clube do Trato)
  - [x] Plano Básico, Intermediário, Premium
- [x] Script `seed_demo_tenant.sql` - Tenant de demonstração com dados de exemplo
  - [x] 1 tenant demo
  - [x] 2 usuários (admin + barbeiro)
  - [x] 11 categorias (4 receita + 7 despesa)
  - [x] 3 planos de assinatura
  - [x] 10 receitas de exemplo
  - [x] 10 despesas de exemplo
  - [x] 3 assinaturas de exemplo
- [x] Documentação: `backend/scripts/SEED_GUIDE.md`
- [x] Programa Go: `backend/cmd/seed/main.go`
- [x] Comandos make: `make seed-demo`, `make seed-prod`, `make seed-clean`, `make seed-verify`

**Files Created:**

- ✅ `backend/scripts/sql/seed_categories.sql` (11 categorias com cores do Design System)
- ✅ `backend/scripts/sql/seed_plans.sql` (3 planos: Básico R$59.90, Intermediário R$89.90, Premium R$129.90)
- ✅ `backend/scripts/sql/seed_demo_tenant.sql` (tenant completo + usuários + dados de exemplo)
- ✅ `backend/scripts/SEED_GUIDE.md` (documentação completa com troubleshooting)
- ✅ `backend/cmd/seed/main.go` (programa Go com flags --mode e --tenant-id)
- ✅ `backend/Makefile` (seção ##@ Seeds com 6 comandos)

---

### ✅ T-PROD-002 — Validação de Integridade

- **Responsável:** QA / Backend
- **Prioridade:** 🔴 Alta
- **Estimativa:** 4h
- **Sprint:** Sprint 7
- **Status:** ✅ Concluído (17/11/2025)
- **Deliverable:** Suite de validação de banco e APIs

#### Critérios de Aceitação

- [x] Script de validação de schema (`scripts/validate_schema.sh`)
  - [x] Verifica tabelas core existem
  - [x] Verifica índices essenciais
  - [x] Verifica RLS em tabelas sensíveis (warning se ausente)
  - [x] Verifica constraints críticas e migrations
- [x] Health check endpoint completo (`GET /health`)
  - [x] Database connection OK (ping + pool stats)
  - [x] Migrations: versão mais recente (`schema_migrations`)
  - [x] Redis: suporte previsto (retorna `not_configured` se indisponível)
  - [x] External APIs: Asaas reachability
- [x] Smoke tests E2E (`scripts/smoke_tests.sh`)
  - [x] Criar tenant → OK
  - [x] Criar usuário → OK
  - [x] Login → OK
  - [x] Criar receita → OK (com fallback de aviso se categoria não existir)
  - [x] Listar receitas → OK
- [x] Documentação: `VALIDATION_GUIDE.md` atualizada

**Deliverables criados/ajustados:**

- `scripts/validate_schema.sh`
- `scripts/smoke_tests.sh`
- `VALIDATION_GUIDE.md`
- `backend/internal/infrastructure/http/handler/health.go` (melhorado)

---

### ✅ T-PROD-003 — Onboarding Flow

- **Responsável:** Frontend / Backend
- **Prioridade:** 🟡 Média
- **Estimativa:** 6h
- **Sprint:** Sprint 8
- **Status:** ✅ Concluído (signup + onboarding + tutorial)
- **Deliverable:** Fluxo de cadastro de novo tenant

#### Critérios de Aceitação

- [x] Página `/signup` (cadastro de novo tenant)
  - [x] Form: Nome da barbearia, CNPJ, Email, Senha
  - [x] Validação: CNPJ válido, email único, senha forte
  - [x] Criação de tenant + primeiro usuário (OWNER)
- [x] Endpoint `POST /auth/signup`
  - [x] Cria tenant
  - [x] Cria primeiro usuário (role: OWNER)
  - [x] Envia email de boas-vindas (opcional)
  - [x] Retorna access_token e refresh_token
- [x] Página `/onboarding` (primeiro acesso)
  - [x] Tour guiado (opcional)
  - [x] Configurar categorias personalizadas
  - [x] Configurar planos de assinatura (se usar Clube do Trato)
- [x] Documentação: Tutorial de primeiro acesso

**Notas de Progresso (20/11/2025):**

- ✅ Wizard finalizado (salva preferências + conclui onboarding com cookie de bloqueio até completar).
- ✅ Backend `/auth/signup` com validação de CNPJ, senha forte, tokens (access + refresh) e retorno do tenant.
- ✅ `/auth/me` inclui dados do tenant (`onboarding_completed`) para redirecionamento automático.
- ✅ Guia de primeiro acesso: `docs/ONBOARDING_GUIDE.md`.

**Files to Create:**

- `frontend/app/(auth)/signup/page.tsx`
- `frontend/app/(private)/onboarding/page.tsx`
- `backend/internal/application/usecase/auth/signup_usecase.go`
- `backend/internal/infrastructure/http/handler/auth_handler.go` (adicionar signup)
- `docs/ONBOARDING_GUIDE.md`

---

### ✅ T-PROD-004 — Documentação de Deploy

- **Responsável:** DevOps
- **Prioridade:** 🟡 Média
- **Estimativa:** 2h
- **Sprint:** Sprint 8
- **Status:** ✅ Concluído
- **Deliverable:** Guia completo de deploy em produção

#### Critérios de Aceitação

- [x] `docs/DEPLOY_PRODUCTION.md` criado com:
  - [x] Checklist pré-deploy
  - [x] Variáveis de ambiente obrigatórias
  - [x] Comandos de deploy (backend + frontend)
  - [x] Verificação pós-deploy
  - [x] Procedimentos de rollback
  - [x] Monitoramento inicial (logs, métricas)
- [x] Scripts de deploy atualizados
  - [x] `scripts/deploy-backend.sh`
  - [x] `scripts/deploy-frontend.sh`
- [x] CI/CD pipeline validado
  - [x] GitHub Actions roda testes
  - [x] Deploy manual em produção (aprovação)

**Files to Create:**

- `docs/DEPLOY_PRODUCTION.md`
- `scripts/deploy-backend.sh`
- `scripts/deploy-frontend.sh`
- `.github/workflows/deploy-production.yml`

---

## 📈 Métricas de Sucesso

### Fase 5 completa quando:

- [x] ✅ Todos os 4 tasks concluídos (100%)
- [x] ✅ Seeds de dados essenciais criados
- [x] ✅ Validação de integridade passando
- [x] ✅ Onboarding flow funcional
- [x] ✅ Documentação de deploy completa
- [x] ✅ V2 pronto para receber primeiros clientes em produção

---

## 🎯 Deliverables da Fase 5

| #   | Deliverable                                              | Status                    |
| --- | -------------------------------------------------------- | ------------------------- |
| 1   | Seeds de dados iniciais (categorias, planos, demo)       | ✅ Concluído (17/11/2025) |
| 2   | Validação de integridade (schema + health + smoke tests) | ✅ Concluído (17/11/2025) |
| 3   | Onboarding flow (signup + primeiro acesso)               | ✅ Concluído (20/11/2025) |
| 4   | Documentação de deploy em produção                       | ✅ Concluído (20/11/2025) |

---

## 🚀 Próximos Passos

Após completar **100%** da Fase 5:

👉 **Iniciar FASE 6 — Hardening** (`Tarefas/FASE_6_HARDENING.md`)

**Resumo Fase 6:**

- Segurança (rate limiting avançado, auditoria, RBAC completo)
- Observabilidade (Prometheus, Grafana, Sentry)
- Performance (query optimization, caching Redis)
- Compliance (LGPD, backup, DR)
- Load testing e otimização

---

## 📝 Notas de Implementação

### Seed de Categorias Padrão

As categorias padrão devem cobrir os casos mais comuns de barbearias:

**Categorias de Receita:**

- Serviços (corte, barba, coloração, etc.)
- Produtos (pomadas, shampoos, etc.)
- Assinaturas (Clube do Trato)
- Outros

**Categorias de Despesa:**

- Salários (barbeiros, recepcionista)
- Aluguel (espaço físico)
- Fornecedores (produtos para revenda)
- Impostos (MEI, SIMPLES)
- Marketing (redes sociais, anúncios)
- Utilidades (água, luz, internet)
- Outros

### Planos de Assinatura Padrão

Sugestão de planos iniciais para o Clube do Trato:

- **Básico** (R$ 59,90/mês): 2 cortes/mês
- **Intermediário** (R$ 89,90/mês): 4 cortes/mês + 1 barba
- **Premium** (R$ 129,90/mês): Ilimitado cortes + barbas

### Tenant Demo

O tenant demo deve ter dados realistas para:

- Demonstrações comerciais
- Testes de integração
- Validação visual do sistema

**Dados sugeridos:**

- Nome: "Barbearia Demo"
- CNPJ: 00.000.000/0001-00 (fictício)
- 10 receitas nos últimos 30 dias
- 10 despesas nos últimos 30 dias
- 3 assinaturas ativas
- 1 usuário admin (demo@barberpro.dev / Demo@1234)

### Health Check Completo

O endpoint `/health` deve retornar:

```json
{
  "status": "healthy",
  "timestamp": "2025-11-17T10:00:00Z",
  "version": "2.0.0",
  "checks": {
    "database": {
      "status": "up",
      "latency_ms": 12
    },
    "migrations": {
      "status": "up_to_date",
      "applied": 15,
      "pending": 0
    },
    "redis": {
      "status": "up",
      "latency_ms": 3
    },
    "external_apis": {
      "asaas": {
        "status": "up",
        "latency_ms": 150
      }
    }
  }
}
```

---

## 📝 Changelog

### 21/11/2025

- ✅ **T-PROD-003 Concluído** — Signup + onboarding guiado finalizados
  - Backend: `/auth/signup` com validação de CNPJ, senha forte, tokens completos e retorno do tenant em `/auth/me`.
  - Frontend: validações fortes em `/signup`, guarda de onboarding via middleware/cookies e wizard ajustado (config + conclusão).
  - Documentação: `docs/ONBOARDING_GUIDE.md` e ajustes de testes (unit + e2e).
- ✅ **T-PROD-004 Concluído** — Guia e pipeline de deploy
  - Scripts `scripts/deploy-backend.sh` e `scripts/deploy-frontend.sh` com backup, owner correto e restart seguro.
  - Workflow GitHub Actions `deploy-production.yml` com aprovação de ambiente `production`.
  - Documentação `docs/DEPLOY_PRODUCTION.md` com checklist, rollback e monitoramento pós-deploy.

### 20/11/2025

- ✅ **T-PROD-002 Concluído** — Validação de integridade completa
  - Scripts de validação de schema e smoke tests criados
  - Health check endpoint aprimorado
  - Documentação VALIDATION_GUIDE.md atualizada
  - Progresso: 25% → 50%
- 🟡 **T-PROD-003 em andamento** — Step 2 (configurações iniciais) finalizado no frontend
  - `tenantConfigService` atualizado para persistir preferências pós-login
  - Formulário do wizard conectado ao backend com validações
  - Faltam `/signup`, endpoint `POST /auth/signup` e tutorial de primeiro acesso

### 17/11/2025

- ✅ **T-PROD-001 Concluído** — Seeds de dados iniciais implementados
  - Criados 3 scripts SQL (categories, plans, demo_tenant)
  - Programa Go com suporte a --mode=demo e --mode=prod
  - 6 comandos make adicionados (seed-demo, seed-prod, seed-clean, etc)
  - Documentação completa em SEED_GUIDE.md
  - Progresso: 0% → 25%

---

**Última Atualização:** 21/11/2025
**Status:** ✅ Concluída (100% - 4/4 tarefas concluídas)
**Próxima Tarefa:** Abrir checklist da Fase 6 (LGPD/Backup)
