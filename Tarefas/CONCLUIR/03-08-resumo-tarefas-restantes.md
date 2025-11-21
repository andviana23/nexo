# 03-08 - Tarefas Backend e Frontend Restantes

**Nota:** Arquivos resumidos. Detalhamento completo será fornecido quando necessário.

---

## 03 - Backend: Repository Implementations (PostgreSQL)

**Prioridade:** 🔴 CRÍTICA
**Estimativa:** 5 dias
**Dependências:** 01, 02

Implementar todos os repositórios PostgreSQL para as 19 entidades usando sqlc.

---

## 04 - Backend: Use Cases Base

**Prioridade:** 🔴 CRÍTICA
**Estimativa:** 4 dias
**Dependências:** 01, 02, 03

Use cases essenciais:

- DRE: GenerateDREUseCase
- Fluxo: GenerateFluxoDiarioUseCase, CreateCompensacaoUseCase
- Metas: SetMetaMensalUseCase, CalculateProgressUseCase
- Precificação: CalculatePrecoUseCase
- Contas: CreateContaPagarUseCase, CreateContaReceberUseCase

---

## 05 - Backend: HTTP Handlers

**Prioridade:** 🔴 CRÍTICA
**Estimativa:** 3 dias
**Dependências:** 01-04

Handlers REST para todos os módulos:

- DTOs (Request/Response)
- Mappers (Domain ↔ DTO)
- Handlers
- Rotas

---

## 06 - Backend: Cron Jobs

**Prioridade:** 🟡 MÉDIA
**Estimativa:** 2 dias
**Dependências:** 01-05

Jobs agendados:

- GenerateDREJob (dia 1º 05:00)
- GenerateFluxoDiarioJob (06:00)
- MarcarCompensacoesJob (07:00)

---

## 07 - Frontend: Service Layer

**Prioridade:** 🔴 CRÍTICA
**Estimativa:** 2 dias
**Dependências:** 05 (handlers prontos)

Services API:

- `api/dre.ts`
- `api/fluxo-caixa.ts`
- `api/metas.ts`
- `api/precificacao.ts`
- `api/contas.ts`
- `api/comissoes.ts`
- `api/estoque.ts`

---

## 08 - Frontend: Hooks Base

**Prioridade:** 🔴 CRÍTICA
**Estimativa:** 2 dias
**Dependências:** 07

Hooks React Query:

- `useDRE`, `useDREComparison`
- `useFluxoCaixa`, `useCompensacoes`
- `useMetas`, `useMetasBarbeiro`
- `usePrecificacao`, `useSimularPreco`
- `useContasPagar`, `useContasReceber`
- `useComissoes`
- `useEstoque`

---

**Total:** ~23 dias (3 semanas full-time) para completar a base do sistema.
