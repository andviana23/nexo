# 🎫 Módulo de Assinaturas — NEXO

**Versão:** 1.0  
**Data de Criação:** 03/12/2025  
**Status:** 🚧 EM PROGRESSO  
**Responsável:** Equipe de Desenvolvimento  
**Estimativa Total:** 70-90 horas (~15-18 dias úteis)

---

## 📊 Progresso Geral

| Sprint | Componente | Arquivo | Status | Progresso |
|--------|-----------|---------|--------|-----------|
| Sprint 1 | Banco de Dados | [01-BANCO-DE-DADOS.md](./01-BANCO-DE-DADOS.md) | ✅ Concluído | 100% |
| Sprint 2 | Backend Core | [02-BACKEND.md](./02-BACKEND.md) | ⬜ Não Iniciado | 0% |
| Sprint 3 | Integração Asaas | [03-INTEGRACAO-ASAAS.md](./03-INTEGRACAO-ASAAS.md) | ⬜ Não Iniciado | 0% |
| Sprint 4 | Frontend | [04-FRONTEND.md](./04-FRONTEND.md) | ⬜ Não Iniciado | 0% |
| Sprint 5 | Testes & QA | [05-TESTES-QA.md](./05-TESTES-QA.md) | ⬜ Não Iniciado | 0% |

**📈 PROGRESSO TOTAL: 20% (1/5 Sprints)**

---

## 📚 Documentação de Referência

> ⚠️ **OBRIGATÓRIO:** Antes de iniciar qualquer tarefa, consultar:
> 
> - **[FLUXO_ASSINATURA.md](../../docs/11-Fluxos/FLUXO_ASSINATURA.md)** — Fonte da verdade do módulo
> - **[PRD-NEXO.md](../../docs/07-produto-e-funcionalidades/PRD-NEXO.md)** — Requisitos de produto
> - **[RBAC.md](../../docs/06-seguranca/RBAC.md)** — Permissões por role
> - **[ARQUITETURA.md](../../docs/02-arquitetura/ARQUITETURA.md)** — Padrões arquiteturais

---

## 🎯 Objetivo do Módulo

Implementar sistema completo de **assinaturas recorrentes** para barbearias, com:

1. **Planos** — CRUD de modelos de assinatura (templates internos)
2. **Assinantes** — Gestão de assinaturas ativas com 3 formas de pagamento
3. **Integração Asaas** — Cobranças via cartão de crédito com renovação automática
4. **Pagamentos Manuais** — PIX e Dinheiro com controle de vencimento
5. **Relatórios** — Métricas de receita, churn, e breakdown por plano/forma

---

## 🗓️ Cronograma de Sprints

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    ROADMAP DE IMPLEMENTAÇÃO                              │
└─────────────────────────────────────────────────────────────────────────┘

Semana 1 (Dias 1-3):
┌──────────────────────────────────────────────────────────────────────────┐
│ SPRINT 1: BANCO DE DADOS                                                 │
│ 📂 01-BANCO-DE-DADOS.md                                                  │
│ ⏱️ Estimativa: 4-6h                                                      │
│ 🎯 Entregável: 3 tabelas + índices + migrations aplicadas               │
└──────────────────────────────────────────────────────────────────────────┘
                                    ↓
Semana 1-2 (Dias 3-7):
┌──────────────────────────────────────────────────────────────────────────┐
│ SPRINT 2: BACKEND CORE                                                   │
│ 📂 02-BACKEND.md                                                         │
│ ⏱️ Estimativa: 20-25h                                                    │
│ 🎯 Entregável: Entidades, Repos, Use Cases, Handlers REST               │
└──────────────────────────────────────────────────────────────────────────┘
                                    ↓
Semana 2 (Dias 6-9):
┌──────────────────────────────────────────────────────────────────────────┐
│ SPRINT 3: INTEGRAÇÃO ASAAS                                               │
│ 📂 03-INTEGRACAO-ASAAS.md                                               │
│ ⏱️ Estimativa: 15-20h                                                    │
│ 🎯 Entregável: Gateway HTTP, Webhooks, Sync status                      │
└──────────────────────────────────────────────────────────────────────────┘
                                    ↓
Semana 2-3 (Dias 8-14):
┌──────────────────────────────────────────────────────────────────────────┐
│ SPRINT 4: FRONTEND                                                       │
│ 📂 04-FRONTEND.md                                                        │
│ ⏱️ Estimativa: 25-30h                                                    │
│ 🎯 Entregável: 4 páginas, componentes, hooks, services                  │
└──────────────────────────────────────────────────────────────────────────┘
                                    ↓
Semana 3 (Dias 14-18):
┌──────────────────────────────────────────────────────────────────────────┐
│ SPRINT 5: TESTES & QA                                                    │
│ 📂 05-TESTES-QA.md                                                       │
│ ⏱️ Estimativa: 6-10h                                                     │
│ 🎯 Entregável: E2E tests, Smoke tests, Validação completa               │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## 📋 Dependências entre Sprints

```
Sprint 1 (DB) ─────► Sprint 2 (Backend) ─────► Sprint 3 (Asaas)
                           │                        │
                           │                        ▼
                           └──────────────► Sprint 4 (Frontend)
                                                    │
                                                    ▼
                                            Sprint 5 (Testes)
```

**Regras:**
- Sprint 2 só inicia após Sprint 1 concluída
- Sprint 3 e 4 podem iniciar em paralelo após Sprint 2
- Sprint 5 só inicia após Sprint 3 e 4 concluídas

---

## 🔐 Permissões RBAC (Referência)

| Página | Administrador | Gerente | Recepção | Barbeiro |
|--------|---------------|---------|----------|----------|
| Planos | CRUD completo | CRUD completo | Visualizar | ❌ |
| Assinantes | CRUD + Cancelar | CRUD + Cancelar | Criar + Visualizar | ❌ |
| Relatórios | Visualizar | Visualizar | Visualizar | ❌ |

---

## 📦 Entregáveis por Sprint

### Sprint 1: Banco de Dados
- [ ] Migration: tabela `plans`
- [ ] Migration: tabela `subscriptions`  
- [ ] Migration: tabela `subscription_payments`
- [ ] Índices de performance
- [ ] Queries sqlc

### Sprint 2: Backend Core
- [ ] Entidades de domínio
- [ ] Interfaces de repositório
- [ ] Implementações sqlc
- [ ] DTOs Request/Response
- [ ] Use Cases (CRUD + Actions)
- [ ] Handlers HTTP
- [ ] Cron Job de vencimentos

### Sprint 3: Integração Asaas
- [ ] Gateway HTTP com retry
- [ ] Métodos: Customer, Subscription, PaymentLink
- [ ] Webhook handler
- [ ] Validação de signature
- [ ] Processamento de eventos
- [ ] Fallback para manual

### Sprint 4: Frontend
- [ ] Página: Lista de Planos
- [ ] Página: Lista de Assinantes
- [ ] Página: Relatórios
- [ ] Modal: Novo Plano
- [ ] Wizard: Nova Assinatura
- [ ] Modal: Renovar/Cancelar
- [ ] Hooks e Services

### Sprint 5: Testes & QA
- [ ] Smoke tests backend
- [ ] Testes E2E Playwright
- [ ] Validação de RBAC
- [ ] Teste de integração Asaas (sandbox)

---

## 🚀 Como Iniciar

1. Ler completamente o [FLUXO_ASSINATURA.md](../../docs/11-Fluxos/FLUXO_ASSINATURA.md)
2. Iniciar pela Sprint 1: [01-BANCO-DE-DADOS.md](./01-BANCO-DE-DADOS.md)
3. Marcar tarefas como ✅ conforme conclusão
4. Atualizar progresso neste arquivo

---

## 📝 Histórico de Alterações

| Data | Versão | Alteração |
|------|--------|-----------|
| 03/12/2025 | 1.0 | Criação do plano de implementação |

---

**FIM DO DOCUMENTO OVERVIEW**
