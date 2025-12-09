# ✅ CHECKLIST — SPRINT 1: INFRAESTRUTURA BASE

> **Status:** 🟢 90% Completo  
> **Período:** Sprints 10-12 (Concluído)  
> **Próximo:** Sprint 2 (Despesas Fixas)

---

## 📊 RESUMO

```
████████████████████ 90% COMPLETO
```

| Categoria | Completo | Pendente |
|-----------|:--------:|:--------:|
| Database | 7/8 | 1 |
| Queries sqlc | 5/6 | 1 |
| Domain | 5/6 | 1 |
| Repository | 5/6 | 1 |
| Use Cases | 21/24 | 3 |
| Handlers | 20/26 | 6 |
| DTOs | 9/12 | 3 |

---

## 1️⃣ DATABASE — MIGRATIONS

### Tabelas Principais

- [x] `contas_a_pagar` — `migrations/003_full_schema.sql:280-298`
  - [x] Campos: id, tenant_id, descricao, categoria_id, fornecedor, valor, tipo, recorrente, periodicidade, data_vencimento, data_pagamento, status, comprovante_url, pix_code, observacoes, criado_em, atualizado_em
  - [x] FK para tenants
  - [x] FK para categorias
  - [x] Índice por tenant_id
  - [x] Índice por data_vencimento
  - [x] RLS habilitado

- [x] `contas_a_receber` — `migrations/003_full_schema.sql:300-316`
  - [x] Campos: id, tenant_id, origem, assinatura_id, servico_id, descricao, valor, valor_pago, data_vencimento, data_recebimento, status, observacoes, criado_em, atualizado_em
  - [x] FK para tenants
  - [x] FK para assinaturas
  - [x] FK para servicos
  - [x] Índice por tenant_id
  - [x] RLS habilitado

- [x] `compensacoes_bancarias` — `migrations/003_full_schema.sql:318-335`
  - [x] Campos: id, tenant_id, conta_receber_id, data_compensacao, valor_compensado, banco, agencia, conta, observacoes, criado_em
  - [x] FK para contas_a_receber
  - [x] Índice por tenant_id
  - [x] RLS habilitado

- [x] `fluxo_caixa_diario` — `migrations/003_full_schema.sql:337-352`
  - [x] Campos: id, tenant_id, unidade_id, data, abertura, entradas_dinheiro, entradas_cartao, entradas_pix, saidas, sangrias, suprimentos, fechamento, diferenca, observacoes
  - [x] FK para tenants
  - [x] FK para unidades
  - [x] Índice por (tenant_id, data)
  - [x] RLS habilitado

- [x] `dre_mensal` — `migrations/003_full_schema.sql:354-378`
  - [x] Campos: id, tenant_id, unidade_id, ano, mes, receita_bruta, deducoes, receita_liquida, custos_servicos, lucro_bruto, despesas_operacionais, despesas_fixas, lucro_operacional, resultado_financeiro, lucro_antes_ir, provisao_ir, lucro_liquido, gerado_em
  - [x] FK para tenants
  - [x] FK para unidades
  - [x] Unique constraint (tenant_id, unidade_id, ano, mes)
  - [x] RLS habilitado

- [x] `metas_mensais` — `migrations/003_full_schema.sql:380-392`
  - [x] Campos padrão + meta_valor
  - [x] FK para tenants
  - [x] RLS habilitado

- [x] `metas_barbeiro` — `migrations/003_full_schema.sql:394+`
  - [x] Campos para metas individuais
  - [x] FK para barbeiros
  - [x] RLS habilitado

- [ ] ❌ `despesas_fixas` — **NÃO EXISTE** (Sprint 2)

---

## 2️⃣ SQL QUERIES (sqlc)

### ✅ `contas_a_pagar.sql` (12 queries)

- [x] `CreateContaPagar` — INSERT com todos os campos
- [x] `GetContaPagarByID` — SELECT por id + tenant_id
- [x] `ListContasPagarByTenant` — SELECT com paginação
- [x] `ListContasPagarByStatus` — Filtro por status
- [x] `ListContasPagarByPeriod` — Filtro por período
- [x] `ListContasPagarVencidas` — Contas vencidas
- [x] `ListContasPagarRecorrentes` — Apenas recorrentes
- [x] `UpdateContaPagar` — UPDATE completo
- [x] `MarcarContaPagarComoPaga` — Atualiza status para PAGO
- [x] `MarcarContaPagarComoAtrasada` — Batch update para atrasadas
- [x] `DeleteContaPagar` — DELETE por id + tenant_id
- [x] `SumContasPagarByPeriod` — Soma total do período
- [x] `SumContasPagasByPeriod` — Soma das pagas
- [x] `CountContasPagarByStatus` — Contagem por status
- [x] `CountContasPagarByTenant` — Total do tenant

### ✅ `contas_a_receber.sql` (11 queries)

- [x] `CreateContaReceber` — INSERT
- [x] `GetContaReceberByID` — SELECT por id + tenant_id
- [x] `ListContasReceberByTenant` — Com paginação
- [x] `ListContasReceberByStatus` — Filtro status
- [x] `ListContasReceberByPeriod` — Filtro período
- [x] `ListContasReceberVencidas` — Vencidas
- [x] `ListContasReceberByAssinatura` — Por assinatura
- [x] `ListContasReceberByOrigem` — Por origem
- [x] `UpdateContaReceber` — UPDATE completo
- [x] `MarcarContaReceberComoRecebida` — Atualiza para RECEBIDO
- [x] `DeleteContaReceber` — DELETE

### ✅ `compensacoes_bancarias.sql`

- [x] `CreateCompensacao`
- [x] `GetCompensacaoByID`
- [x] `ListCompensacoesByTenant`
- [x] `ListCompensacoesByContaReceber`
- [x] `DeleteCompensacao`

### ✅ `fluxo_caixa_diario.sql`

- [x] `CreateFluxoCaixa`
- [x] `GetFluxoCaixaByDate`
- [x] `ListFluxoCaixaByPeriod`
- [x] `UpdateFluxoCaixa`
- [x] Queries de agregação

### ✅ `dre_mensal.sql`

- [x] `CreateDRE`
- [x] `GetDREByMonth`
- [x] `ListDREByYear`
- [x] `UpdateDRE`
- [x] `UpsertDRE`

### ❌ `despesas_fixas.sql` — **NÃO EXISTE** (Sprint 2)

---

## 3️⃣ DOMAIN LAYER

### Entities

- [x] `ContaPagar` — `internal/domain/entity/conta_pagar.go`
  - [x] Struct com todos os campos
  - [x] Método `Validate()`
  - [x] Método `MarcarComoPaga()`
  - [x] Método `EstaVencida()`
  - [x] Método `IsRecorrente()`

- [x] `ContaReceber` — `internal/domain/entity/conta_receber.go`
  - [x] Struct com todos os campos
  - [x] Método `Validate()`
  - [x] Método `MarcarComoRecebida()`
  - [x] Método `EstaVencida()`

- [x] `CompensacaoBancaria` — `internal/domain/entity/compensacao_bancaria.go`
  - [x] Struct
  - [x] Métodos de validação

- [x] `FluxoCaixaDiario` — `internal/domain/entity/fluxo_caixa.go`
  - [x] Struct
  - [x] Cálculo de saldo

- [x] `DREMensal` — `internal/domain/entity/dre_mensal.go`
  - [x] Struct
  - [x] Cálculo de margens

- [ ] ❌ `DespesaFixa` — **NÃO EXISTE** (Sprint 2)

### Value Objects

- [x] `StatusConta` — Enum: ABERTO, PAGO, ATRASADO, CANCELADO
- [x] `StatusRecebimento` — Enum: PENDENTE, RECEBIDO, ATRASADO
- [x] `TipoConta` — Enum: DESPESA_FIXA, DESPESA_VARIAVEL, etc.
- [x] `Periodicidade` — Enum: MENSAL, SEMANAL, QUINZENAL, etc.

---

## 4️⃣ REPOSITORY LAYER

### Interfaces

- [x] `ContaPagarRepository` — `internal/domain/repository/conta_pagar_repository.go`
- [x] `ContaReceberRepository`
- [x] `CompensacaoBancariaRepository`
- [x] `FluxoCaixaDiarioRepository`
- [x] `DREMensalRepository`
- [ ] ❌ `DespesaFixaRepository` — **NÃO EXISTE** (Sprint 2)

### Implementações PostgreSQL

- [x] `PGContaPagarRepository`
- [x] `PGContaReceberRepository`
- [x] `PGCompensacaoBancariaRepository`
- [x] `PGFluxoCaixaDiarioRepository`
- [x] `PGDREMensalRepository`
- [ ] ❌ `PGDespesaFixaRepository` — **NÃO EXISTE** (Sprint 2)

---

## 5️⃣ USE CASES

### Contas a Pagar (6/6) ✅

- [x] `CreateContaPagarUseCase`
- [x] `GetContaPagarUseCase`
- [x] `ListContasPagarUseCase`
- [x] `UpdateContaPagarUseCase`
- [x] `DeleteContaPagarUseCase`
- [x] `MarcarPagamentoUseCase`

### Contas a Receber (6/6) ✅

- [x] `CreateContaReceberUseCase`
- [x] `GetContaReceberUseCase`
- [x] `ListContasReceberUseCase`
- [x] `UpdateContaReceberUseCase`
- [x] `DeleteContaReceberUseCase`
- [x] `MarcarRecebimentoUseCase`

### Compensações (5/5) ✅

- [x] `CreateCompensacaoUseCase`
- [x] `GetCompensacaoUseCase`
- [x] `ListCompensacoesUseCase`
- [x] `DeleteCompensacaoUseCase`
- [x] `MarcarCompensacaoUseCase`

### Fluxo de Caixa (3/3) ✅

- [x] `GenerateFluxoDiarioUseCase`
- [x] `GetFluxoCaixaUseCase`
- [x] `ListFluxoCaixaUseCase`

### DRE (3/3) ✅

- [x] `GenerateDREUseCase`
- [x] `GetDREUseCase`
- [x] `ListDREUseCase`

### Despesas Fixas (0/6) ❌

- [ ] `CreateDespesaFixaUseCase` — Sprint 2
- [ ] `GetDespesaFixaUseCase` — Sprint 2
- [ ] `ListDespesasFixasUseCase` — Sprint 2
- [ ] `UpdateDespesaFixaUseCase` — Sprint 2
- [ ] `ToggleDespesaFixaUseCase` — Sprint 2
- [ ] `DeleteDespesaFixaUseCase` — Sprint 2

### Painel Mensal (0/2) ❌

- [ ] `GetPainelMensalUseCase` — Sprint 3
- [ ] `GetProjecoesUseCase` — Sprint 3

---

## 6️⃣ HTTP HANDLERS

### FinancialHandler — `internal/infra/http/handler/financial_handler.go`

**Arquivo:** 1342 linhas ✅

### Endpoints Contas a Pagar (6/6) ✅

- [x] `POST /financial/payables` → `CreateContaPagar()`
- [x] `GET /financial/payables` → `ListContasPagar()`
- [x] `GET /financial/payables/:id` → `GetContaPagar()`
- [x] `PUT /financial/payables/:id` → `UpdateContaPagar()`
- [x] `DELETE /financial/payables/:id` → `DeleteContaPagar()`
- [x] `POST /financial/payables/:id/payment` → `MarcarPagamento()`

### Endpoints Contas a Receber (6/6) ✅

- [x] `POST /financial/receivables` → `CreateContaReceber()`
- [x] `GET /financial/receivables` → `ListContasReceber()`
- [x] `GET /financial/receivables/:id` → `GetContaReceber()`
- [x] `PUT /financial/receivables/:id` → `UpdateContaReceber()`
- [x] `DELETE /financial/receivables/:id` → `DeleteContaReceber()`
- [x] `POST /financial/receivables/:id/receipt` → `MarcarRecebimento()`

### Endpoints Compensações (3/3) ✅

- [x] `GET /financial/compensations` → `ListCompensacoes()`
- [x] `GET /financial/compensations/:id` → `GetCompensacao()`
- [x] `DELETE /financial/compensations/:id` → `DeleteCompensacao()`

### Endpoints Fluxo de Caixa (2/2) ✅

- [x] `GET /financial/cashflow` → `ListFluxoCaixa()`
- [x] `GET /financial/cashflow/:date` → `GetFluxoCaixa()`

### Endpoints DRE (2/2) ✅

- [x] `GET /financial/dre` → `ListDRE()`
- [x] `GET /financial/dre/:year/:month` → `GetDRE()`

### Endpoints Despesas Fixas (0/6) ❌

- [ ] `POST /financial/fixed-expenses` — Sprint 2
- [ ] `GET /financial/fixed-expenses` — Sprint 2
- [ ] `GET /financial/fixed-expenses/:id` — Sprint 2
- [ ] `PUT /financial/fixed-expenses/:id` — Sprint 2
- [ ] `POST /financial/fixed-expenses/:id/toggle` — Sprint 2
- [ ] `DELETE /financial/fixed-expenses/:id` — Sprint 2

### Endpoints Dashboard (0/2) ❌

- [ ] `GET /financial/dashboard` — Sprint 3
- [ ] `GET /financial/projections` — Sprint 3

---

## 7️⃣ DTOs

### Contas a Pagar ✅

- [x] `ContaPagarCreateRequest`
- [x] `ContaPagarUpdateRequest`
- [x] `ContaPagarResponse`

### Contas a Receber ✅

- [x] `ContaReceberCreateRequest`
- [x] `ContaReceberUpdateRequest`
- [x] `ContaReceberResponse`

### Compensações ✅

- [x] `CompensacaoResponse`

### Fluxo de Caixa ✅

- [x] `FluxoCaixaResponse`

### DRE ✅

- [x] `DREMensalResponse`

### Despesas Fixas ❌

- [ ] `DespesaFixaCreateRequest` — Sprint 2
- [ ] `DespesaFixaUpdateRequest` — Sprint 2
- [ ] `DespesaFixaResponse` — Sprint 2

### Painel Mensal ❌

- [ ] `PainelMensalResponse` — Sprint 3
- [ ] `ProjecaoResponse` — Sprint 3

---

## 8️⃣ ROTAS REGISTRADAS

**Arquivo:** `cmd/api/main.go` linhas 568-594

```go
// ✅ Registradas
financial.POST("/payables", financialHandler.CreateContaPagar)
financial.GET("/payables", financialHandler.ListContasPagar)
financial.GET("/payables/:id", financialHandler.GetContaPagar)
financial.PUT("/payables/:id", financialHandler.UpdateContaPagar)
financial.DELETE("/payables/:id", financialHandler.DeleteContaPagar)
financial.POST("/payables/:id/payment", financialHandler.MarcarPagamento)

financial.POST("/receivables", financialHandler.CreateContaReceber)
financial.GET("/receivables", financialHandler.ListContasReceber)
financial.GET("/receivables/:id", financialHandler.GetContaReceber)
financial.PUT("/receivables/:id", financialHandler.UpdateContaReceber)
financial.DELETE("/receivables/:id", financialHandler.DeleteContaReceber)
financial.POST("/receivables/:id/receipt", financialHandler.MarcarRecebimento)

financial.GET("/compensations", financialHandler.ListCompensacoes)
financial.GET("/compensations/:id", financialHandler.GetCompensacao)
financial.DELETE("/compensations/:id", financialHandler.DeleteCompensacao)

financial.GET("/cashflow", financialHandler.ListFluxoCaixa)
financial.GET("/cashflow/:date", financialHandler.GetFluxoCaixa)

financial.GET("/dre", financialHandler.ListDRE)
financial.GET("/dre/:year/:month", financialHandler.GetDRE)

// ❌ Pendentes (Sprint 2-3)
// financial.POST("/fixed-expenses", ...)
// financial.GET("/fixed-expenses", ...)
// financial.GET("/dashboard", ...)
// financial.GET("/projections", ...)
```

---

## 9️⃣ TESTES

### Testes Unitários

- [x] Domain entities
- [x] Value objects
- [ ] 🔄 Use cases (parcial)

### Testes de Integração

- [x] Repository tests
- [x] Handler tests básicos
- [ ] 🔄 Fluxos completos

### Testes E2E

- [ ] ❌ Pendente Sprint 5

---

## 🎯 PRÓXIMOS PASSOS

1. **✅ Sprint 1 está 90% completo**
2. **➡️ Iniciar Sprint 2: Despesas Fixas**
   - Criar migration da tabela
   - Implementar queries sqlc
   - Criar domain entity
   - Implementar repository
   - Criar use cases
   - Adicionar handlers
3. **➡️ Configurar Cron Job para geração automática**

---

## 📎 ARQUIVOS REFERÊNCIA

| Componente | Caminho |
|------------|---------|
| Migration | `backend/migrations/003_full_schema.sql` |
| Queries | `backend/internal/infra/db/queries/` |
| Handler | `backend/internal/infra/http/handler/financial_handler.go` |
| Use Cases | `backend/internal/application/usecase/financial/` |
| Rotas | `backend/cmd/api/main.go:568-594` |

---

*Última atualização: Dezembro 2024*
