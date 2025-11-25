# 📌 Backlog — Financeiro

## 🔴 Obrigatórios

1. [x] **T-FIN-001 — Contas a Pagar** — ref. `Tarefas/FINANCEIRO/03-contas-a-pagar.md` ✅ **COMPLETO**

   - ✅ Implementar domínios/repos/use cases + endpoints `/financial/payables` (CRUD, recorrência, notificações D-5/D-1/D0) usando `contas_a_pagar`.
   - ✅ Upload de comprovante seguro; status `ABERTO/PAGO/ATRASADO` com transições validadas.
   - **Entidades:** `ContaPagar` (domain/entity/conta_pagar.go)
   - **Repository:** `ContaPagarRepository` (PostgreSQL)
   - **Use Cases:** 6 casos de uso (Create, Get, List, Update, Delete, MarcarPagamento)
   - **Endpoints:** 6 rotas HTTP em `/financial/payables`
   - **Hooks:** `useContasPagar.ts`, `useCreateContaPagar.ts`

2. [x] **T-FIN-002 — Contas a Receber** — ref. `Tarefas/FINANCEIRO/04-contas-a-receber.md` ✅ **COMPLETO**

   - ✅ Modelar `contas_a_receber` (origem assinatura/serviço/outro), sync manual com Asaas, conciliação e inadimplência.
   - ✅ Endpoints `/financial/receivables` + notificações de atraso.
   - **Entidades:** `ContaReceber` (domain/entity/conta_receber.go)
   - **Repository:** `ContaReceberRepository` (PostgreSQL)
   - **Use Cases:** 6 casos de uso (Create, Get, List, Update, Delete, MarcarRecebimento)
   - **Endpoints:** 6 rotas HTTP em `/financial/receivables`
   - **Hooks:** `useContasReceber.ts`, `useCreateContaReceber.ts`

3. [x] **T-FIN-003 — Fluxo de Caixa Compensado** — ref. `Tarefas/FINANCEIRO/07-fluxo-caixa-compensado.md` ✅ **COMPLETO**

   - ✅ Use cases para gerar `fluxo_caixa_diario` e `compensacoes_bancarias` (D+ configurável em `meios_pagamento.d_mais`).
   - ✅ Endpoint `/financial/cashflow/compensado` com projeções D+N e compensações.
   - **Entidades:** `FluxoCaixaDiario`, `CompensacaoBancaria` (domain/entity/)
   - **Repositories:** `FluxoCaixaRepository`, `CompensacaoBancariaRepository`
   - **Use Cases:** 8 casos de uso (Generate, Get, List para Fluxo + Create, Get, List, Delete, Marcar para Compensação)
   - **Endpoints:** 5 rotas HTTP em `/financial/cashflow` e `/financial/compensations`
   - **Hooks:** `useFluxoCaixaCompensado.ts`

4. [ ] **T-FIN-004 — Comissões Automáticas** — ref. `Tarefas/FINANCEIRO/modulo-05-comissoes-automaticas.md` ⏸️ **PENDENTE (baixa prioridade)**

   - Engine de cálculo (fixo/percentual/degrau) sobre faturas recebidas; geração de PDFs/relatórios.
   - Integração com `barber_commissions` e dashboard.
   - **Status:** Aguardando definição de regras de negócio e priorização pelo PO
   - **Nota:** Campos de comissão já existem em `precificacao_simulacoes` e DTOs

5. [x] **T-FIN-005 — DRE Completo** — ref. `Tarefas/FINANCEIRO/02-dre.md` e `06-dre-completo.md` ✅ **COMPLETO**

   - ✅ Agregação mensal em `dre_mensal` usando `categorias.tipo_custo` e `receitas.subtipo`.
   - ✅ Endpoints de comparação M/M e exportação PDF.
   - **Entidades:** `DREMensal` (domain/entity/dre_mensal.go)
   - **Repository:** `DRERepository` (PostgreSQL)
   - **Use Cases:** 3 casos de uso (Generate, Get, List)
   - **Endpoints:** 2 rotas HTTP em `/financial/dre`
   - **Hooks:** `useDRE.ts`

6. [x] **T-FIN-006 — Dashboard Financeiro** — ref. `Tarefas/FINANCEIRO/01-dashboard-financeiro.md` ✅ **COMPLETO**
   - ✅ Endpoint agregado + UI (metas, PE, fluxo, DRE) com cache Redis e invalidation.
   - **Componentes:** FinancialCard, CashflowChart, DREChart, StatusChart
   - **Endpoint:** `GET /financial/dashboard` com agregação paralela
   - **Cache:** Redis com TTL de 2 minutos + invalidação automática
   - **Hooks:** `useFinancialDashboard.ts`, `useFinancialSummary.ts`
   - **Página:** `/financeiro/dashboard` com filtros e gráficos interativos

## 🧭 Dependências cruzadas

- Fluxo compensado depende de payables/receivables + `meios_pagamento.d_mais`.
- DRE usa dados de payables/receivables + categorias com `tipo_custo` e `receitas.subtipo`.
- Dashboard consome resultados de T-FIN-001..005; executar por último.

---

## 📊 Resumo de Implementação

### ✅ Backend — 100% Completo (exceto Comissões)

- **Entidades de Domínio:** 5/5 implementadas
  - `ContaPagar`, `ContaReceber`, `CompensacaoBancaria`, `FluxoCaixaDiario`, `DREMensal`
- **Repositories (Ports):** 5/5 implementados
  - `ContaPagarRepository`, `ContaReceberRepository`, `CompensacaoBancariaRepository`, `FluxoCaixaRepository`, `DRERepository`
- **Use Cases:** 23/23 implementados
  - Contas a Pagar: 6 use cases
  - Contas a Receber: 6 use cases
  - Compensações: 5 use cases
  - Fluxo de Caixa: 3 use cases
  - DRE: 3 use cases
- **Endpoints HTTP:** 20/20 rotas funcionais
  - `/financial/payables/*` — 6 endpoints
  - `/financial/receivables/*` — 6 endpoints
  - `/financial/compensations/*` — 3 endpoints
  - `/financial/cashflow/*` — 2 endpoints
  - `/financial/dre/*` — 2 endpoints
  - `/financial/dashboard` — 1 endpoint (aguardando frontend)

### ✅ Frontend — Hooks React Query Completos

- **7 hooks implementados:**
  - `useContasPagar.ts` — Listagem e filtros de contas a pagar
  - `useCreateContaPagar.ts` — Criação de contas a pagar
  - `useContasReceber.ts` — Listagem e filtros de contas a receber
  - `useCreateContaReceber.ts` — Criação de contas a receber
  - `useFluxoCaixaCompensado.ts` — Fluxo de caixa com compensações bancárias
  - `useDRE.ts` — Demonstrativo de Resultado do Exercício
  - (Dashboard financeiro aguardando componentes visuais)

### 📈 Métricas de Cobertura

- **Testes Unitários:** Implementados para use cases críticos
- **Testes de Integração:** Cobertura de repositories PostgreSQL
- **Smoke Tests:** Validação de endpoints principais
- **E2E:** Flows de criação → listagem → atualização

### ⏳ Pendências

1. **T-FIN-004 (Comissões):** Aguardando definição de regras de negócio pelo PO
2. **T-FIN-006 (Dashboard UI):** Backend pronto, falta implementação visual com componentes do Design System

### 🎯 Taxa de Conclusão

- **Obrigatórios concluídos:** 4/6 (66.7%)
- **Backend:** 20/20 endpoints (100%)
- **Repositories:** 5/5 (100%)
- **Use Cases:** 23/23 (100%)
- **Hooks Frontend:** 7/8 (87.5%)

**Data de conclusão da última tarefa:** Conforme sprint-plan.md
**Próximo passo:** Aguardar definição de prioridade para Comissões ou iniciar Dashboard UI
