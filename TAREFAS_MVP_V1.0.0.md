# 🎯 Tarefas para Concluir v1.0.0 — MVP CORE

**Data de Emissão:** 27/11/2025  
**Última Atualização:** 06/12/2025 10:00  
**Deadline:** 05/12/2025 (2 dias)  
**Progresso Atual:** 100% ✅  
**Status:** 🟢 **CORE COMPLETO** — P1-P5 finalizados. MVP pronto para release!

---

## 📊 Status por Módulo (06/12/2025)

| Módulo | Backend | Frontend | Status Geral |
|--------|---------|----------|--------------|
| Autenticação | ✅ 100% | ✅ 100% | ✅ PRONTO |
| Serviços | ✅ 100% | ✅ 100% | ✅ PRONTO |
| Categorias | ✅ 100% | ✅ 100% | ✅ PRONTO |
| Profissionais | ✅ 100% | ✅ 100% | ✅ PRONTO |
| Lista da Vez | ✅ 100% | ✅ 100% | ✅ PRONTO |
| Clientes (CRM) | ✅ 100% | ✅ 100% | ✅ **PRONTO** (Histórico de atendimentos) |
| **Agendamento** | ✅ 100% | ✅ 100% | ✅ **PRONTO** (01/12 - 8 bugs + 6 features) |
| **Estoque** | ✅ 100% | ✅ 100% | ✅ **PRONTO** (Entrada/Saída implementadas) |
| **Financeiro** | ✅ 100% | ✅ 100% | ✅ **PRONTO** (Dashboard + Projections funcionando) |
| **Metas** | ✅ 100% | ✅ 100% | ✅ **PRONTO** |
| **Comissões** | ✅ 100% | ✅ 100% | ✅ **PRONTO** (06/12 - Backend 35+ endpoints + Frontend completo) |
| Precificação | ✅ 100% | ❌ 0% | 🔴 Sem UI |
| Relatórios | ✅ 100% | ✅ 100% | ✅ **PRONTO** (DRE + Fluxo + Faturamento + Despesas) |
| **RBAC** | ✅ 100% | ✅ 100% | ✅ **PRONTO** (Middleware + Sidebar Filter) |
| Assinaturas Asaas | ❌ 0% | ❌ 0% | ⚪ Última prioridade |

---

## 🔥 TAREFAS PENDENTES — ORDEM DE PRIORIDADE

### ✅ PRIORIDADE 1: BLOQUEADORES CRÍTICOS — CONCLUÍDO

| # | Tarefa | Esforço | Status | Arquivos Afetados |
|---|--------|---------|--------|-------------------|
| ~~**P1.1**~~ | ~~Criar migration `appointments` + `appointment_services`~~ | ~~3h~~ | ✅ PRONTO | Tabelas já existem no banco |
| ~~**P1.2**~~ | ~~Mover rotas `/stock/*` para grupo JWT protegido~~ | ~~1h~~ | ✅ PRONTO | `backend/cmd/api/main.go` |
| ~~**P1.3**~~ | ~~Alinhar contrato frontend agendamento com backend~~ | ~~2h~~ | ✅ PRONTO | `frontend/src/services/appointment-service.ts`, `use-appointments.ts` |
| ~~**P1.4**~~ | ~~Validar sqlc generate após migrations~~ | ~~1h~~ | ✅ PRONTO | `backend/internal/infra/db/` |

**Subtotal P1:** ~~7h~~ → ✅ CONCLUÍDO

---

### ✅ PRIORIDADE 2: FRONTEND FINANCEIRO — CONCLUÍDO (29/11)

| # | Tarefa | Esforço | Status | Arquivos Criados |
|---|--------|---------|--------|------------------|
| ~~**P2.1**~~ | ~~Dashboard Financeiro consolidado~~ | ~~4h~~ | ✅ PRONTO | `financeiro/page.tsx` |
| ~~**P2.2**~~ | ~~Página DRE Mensal~~ | ~~4h~~ | ✅ PRONTO | `financeiro/dre/page.tsx` |
| ~~**P2.3**~~ | ~~Página Fluxo de Caixa~~ | ~~4h~~ | ✅ PRONTO | `financeiro/fluxo-caixa/page.tsx` |
| ~~**P2.4**~~ | ~~Página Contas a Pagar~~ | ~~3h~~ | ✅ PRONTO | `financeiro/contas-pagar/page.tsx` |
| ~~**P2.5**~~ | ~~Página Contas a Receber~~ | ~~3h~~ | ✅ PRONTO | `financeiro/contas-receber/page.tsx` |
| ~~**P2.6**~~ | ~~Endpoints Backend Dashboard + Projections~~ | ~~2h~~ | ✅ PRONTO | `backend/cmd/api/main.go` (rotas registradas) |

**Subtotal P2:** ~~20h~~ → ✅ CONCLUÍDO (29/11/2025)

---

### ✅ PRIORIDADE 3: METAS + PRECIFICAÇÃO UI — CONCLUÍDO

| # | Tarefa | Esforço | Arquivos Criados | Status |
|---|--------|---------|------------------|--------|
| ~~**P3.1**~~ | ~~Página Metas Mensais (CRUD + progresso)~~ | ~~4h~~ | `metas/mensais/page.tsx` | ✅ PRONTO |
| ~~**P3.2**~~ | ~~Página Metas por Barbeiro~~ | ~~3h~~ | `metas/barbeiros/page.tsx`, `metas/ticket/page.tsx` | ✅ PRONTO |
| ~~**P3.3**~~ | ~~Página Simulador Precificação~~ | ~~4h~~ | `frontend/src/app/(dashboard)/precificacao/page.tsx` | ✅ PRONTO |
| ~~**P3.4**~~ | ~~Relatórios (DRE + Fluxo + Faturamento + Despesas)~~ | ~~4h~~ | `frontend/src/app/(dashboard)/relatorios/page.tsx` | ✅ PRONTO |

**Subtotal P3:** ~~15h~~ → ✅ CONCLUÍDO (02/12/2025)

---

### ✅ PRIORIDADE 4: ESTOQUE + CRM — CONCLUÍDO

| # | Tarefa | Esforço | Arquivos Afetados | Dependências |
|---|--------|---------|-------------------|--------------|
| ~~**P4.1**~~ | ~~Alinhar contrato estoque frontend/backend~~ | ~~2h~~ | `frontend/src/services/stock-service.ts` | ✅ PRONTO |
| ~~**P4.2**~~ | ~~Tela Entrada de Estoque~~ | ~~3h~~ | `frontend/src/app/(dashboard)/estoque/entrada/page.tsx` | ✅ PRONTO |
| ~~**P4.3**~~ | ~~Tela Saída de Estoque~~ | ~~2h~~ | `frontend/src/app/(dashboard)/estoque/saida/page.tsx` | ✅ PRONTO |
| ~~**P4.4**~~ | ~~CRM: histórico de atendimentos~~ | ~~3h~~ | `frontend/src/app/(dashboard)/clientes/[id]/page.tsx` | ✅ PRONTO |

**Subtotal P4:** ~~10h~~ → ✅ CONCLUÍDO (02/12/2025)

---

### ✅ PRIORIDADE 5: QUALIDADE + DEPLOY — CONCLUÍDO

| # | Tarefa | Esforço | Arquivos Afetados | Dependências |
|---|--------|---------|-------------------|--------------|
| ~~**P5.1**~~ | ~~RBAC: middleware de permissões~~ | ~~3h~~ | `backend/internal/infra/http/middleware/` | ✅ PRONTO |
| ~~**P5.2**~~ | ~~RBAC: filtro sidebar por role~~ | ~~2h~~ | `frontend/src/components/layout/Sidebar.tsx` | ✅ PRONTO |
| ~~**P5.3**~~ | ~~Testes E2E: agendamento + financeiro~~ | ~~4h~~ | `frontend/tests/e2e/financeiro.spec.ts`, `metas.spec.ts`, `estoque.spec.ts`, `crm.spec.ts` | ✅ PRONTO |
| ~~**P5.4**~~ | ~~Smoke tests backend~~ | ~~2h~~ | `scripts/smoke_tests_complete.sh` | ✅ PRONTO |
| ~~**P5.5**~~ | ~~Deploy staging + validação~~ | ~~3h~~ | `scripts/validate-staging.sh` | ✅ PRONTO |
| ~~**P5.6**~~ | ~~Documentação Swagger atualizada~~ | ~~2h~~ | `docs/04-backend/API_REFERENCE.md` | ✅ PRONTO |

**Subtotal P5:** ~~16h~~ → ✅ CONCLUÍDO (03/12/2025)

---

### ⚪ PRIORIDADE 6: INTEGRAÇÃO ASAAS (04-05/12) — ÚLTIMA PRIORIDADE

> ⚠️ **Nota:** Assinaturas Asaas foi movida para última prioridade. O MVP pode lançar com cobrança manual (PIX/dinheiro) e integração Asaas pode ser entregue em v1.0.1 se necessário.

| # | Tarefa | Esforço | Arquivos Afetados | Dependências |
|---|--------|---------|-------------------|--------------|
| **P6.1** | Criar gateway Asaas (client HTTP) | 4h | `backend/internal/infra/gateway/asaas/` | Nenhuma |
| **P6.2** | Implementar use cases assinatura | 4h | `backend/internal/application/usecase/subscription/` | P6.1 |
| **P6.3** | Criar handlers HTTP assinaturas | 3h | `backend/internal/infra/http/handler/subscription_handler.go` | P6.2 |
| **P6.4** | Webhook Asaas (receber eventos) | 2h | `backend/internal/infra/http/handler/webhook_handler.go` | P6.3 |
| **P6.5** | Frontend: página de planos/checkout | 6h | `frontend/src/app/(dashboard)/assinaturas/` | P6.3 |

**Subtotal P6:** 19h (pode ser postergado para v1.0.1)

---

## 📅 CRONOGRAMA ATUALIZADO

| Data | Sprint | Tarefas | Horas | Status |
|------|--------|---------|-------|--------|
| **27/11 (Qui)** | Sprint 1 | P1.1-P1.4 | 7h | ✅ CONCLUÍDO |
| **27/11 (Qui)** | Sprint 2.1 | P2.1-P2.5 (Financeiro UI) | 18h | ✅ CONCLUÍDO |
| **27/11 (Qui)** | Sprint 3.1 | P3.1, P3.2 (Metas) | 7h | ✅ CONCLUÍDO |
| **29/11 (Sáb)** | Sprint 2.2 | P2.6 (Endpoints Backend) | 2h | ✅ CONCLUÍDO |
| **30/11 (Dom)** | Sprint 3.2 | P3.3, P3.4 (Precificação + Relatórios) | 8h | ✅ CONCLUÍDO |
| **29/11 (Sáb)** | Sprint 4 | P4.1-P4.4 (Estoque + CRM) | 10h | ✅ CONCLUÍDO |
| **01/12 (Seg)** | Sprint 5.1 | P5.1, P5.2 (RBAC) | 5h | ✅ CONCLUÍDO |
| **02-03/12 (Ter-Qua)** | Sprint 5.2 | P5.3-P5.6 (Testes + Deploy) | 11h | ✅ CONCLUÍDO |
| **03/12 (Qua)** | Buffer | Correções e ajustes | 6h | Disponível |
| **04/12 (Qui)** | Sprint 6.1 | P6.1, P6.2, P6.3 (Asaas) | 11h | Opcional |
| **05/12 (Sex)** | Sprint 6.2 | P6.4, P6.5 (Asaas) | 8h | Opcional |

**Total Core (P1-P5):** 68h → ✅ **100% CONCLUÍDO**  
**Total com Asaas (P6):** 87h (opcional para v1.0.1)

**Horas Trabalhadas:** 68h  
**Progresso:** 100% do core ✅

---

## ✅ CHECKLIST GATE v1.0.0 (05/12/2025 - 19:00)

### Funcionalidades Core (100% obrigatório)
- [x] Agendamento funcional end-to-end ✅ (01/12 - 8 bugs + 6 features)
- [x] Lista da Vez operacional
- [x] Financeiro (DRE + Fluxo + Contas) com UI ✅
- [x] Metas (Mensais + Barbeiro + Ticket) com UI ✅
- [x] Comissões calculando corretamente ✅ (05/12 - Backend completo: 35+ endpoints)
- [x] Estoque CRUD operacional (UI entrada/saída concluída) ✅
- [x] CRM básico funcional (histórico atendimentos) ✅ (02/12)
- [x] Relatórios mensais gerados ✅ (02/12)
- [x] Permissões (RBAC) funcionando ✅ (02/12 - Middleware + Sidebar)

### Funcionalidades Opcionais (pode ir para v1.0.1)
- [ ] Assinaturas Asaas integradas (cobrança manual ok para MVP)
- [x] Precificação UI (simulador) ✅ (02/12)

### Qualidade
- [x] Testes backend ≥70%
- [x] Testes frontend ≥60%
- [x] E2E ≥80% passando ✅ (33+ testes - financeiro, metas, estoque, crm)
- [ ] Performance p95 <300ms

### Compliance
- [x] LGPD endpoints funcionais
- [x] Backup automático configurado
- [ ] Multi-tenant validado

### Operacional
- [x] Deploy staging validado ✅ (03/12 - validate-staging.sh 6/6)
- [x] Smoke tests passando ✅ (03/12 - 29/29 testes)
- [x] Documentação atualizada ✅ (03/12 - API_REFERENCE.md)

---

## 🚨 RISCOS E MITIGAÇÕES

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Prazo insuficiente | ~~Baixa~~ ✅ Mitigado | Alto | Core 100% completo |
| UI Financeiro complexa | ~~Baixa~~ ✅ Mitigado | Médio | Componentes implementados |
| Bugs pós-deploy | Média | Médio | Smoke tests + staging validados |

---

## 📝 DECISÕES TOMADAS

1. ✅ **Asaas:** Postergado para última prioridade (v1.0.1 se necessário)
2. ✅ **Financeiro Frontend:** Concluído em 27/11 (P2.1-P2.5)
3. ✅ **Metas Frontend:** Concluído em 27/11 (P3.1-P3.2) — Dashboard, Mensais, Barbeiros, Ticket
4. ✅ **Testes E2E Serviços:** 9/9 passando no Chromium (27/11)
5. ✅ **Correções Backend Metas:** Status constraint e FK criado_por corrigidos
6. ✅ **Agendamento Completo:** 8 bugs + 6 features concluídos em 01/12
7. ✅ **RBAC Agendamento:** Middleware implementado em todas as rotas
8. ✅ **Estoque UI:** Entrada e Saída implementadas (02/12)
9. ✅ **RBAC Sidebar:** Filtro de menu por role implementado (02/12)
10. ✅ **Comissões Backend:** Módulo completo com 35+ endpoints (05/12)
11. ✅ **Comissões Frontend:** Módulo completo (06/12)
    - Dashboard com KPIs e resumos
    - Página de Regras (CRUD)
    - Página de Períodos (CRUD + Close + Pay)
    - Página de Adiantamentos (CRUD + Approve + Reject)
    - Página de Itens de Comissão
    - Menu na Sidebar com RBAC
12. **Relatórios PDF:** Deixar para v1.1 (visualização web é suficiente)
13. **Deploy:** Definir até 02/12

---

**Responsável:** Equipe NEXO  
**Próxima revisão:** 06/12/2025  
**Contato:** CEO - Andrey Viana
