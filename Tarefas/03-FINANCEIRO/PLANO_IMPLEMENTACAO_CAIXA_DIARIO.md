# Plano de Implementação — Caixa Diário

**Versão:** 2.0  
**Data:** 29/11/2025  
**Baseado em:** `docs/11-Fluxos/Fluxo_Financeiro/FLUXO_CAIXA.md`  
**Status:** ✅ CONCLUÍDO (Todas as 7 Fases implementadas)

---

## 📋 Sumário Executivo

### O que é o Caixa Diário?
O Caixa Diário é o **ponto operacional** onde acontece o controle físico de numerário (gaveta de dinheiro). É diferente do Fluxo de Caixa (visão estratégica) — aqui tratamos do operacional: abertura, sangrias, reforços e fechamento.

### Escopo da Implementação
- **2 tabelas novas:** `caixa_diario`, `operacoes_caixa` ✅
- **9 endpoints de API** ✅
- **10 DTOs** (Request/Response) ✅
- **7 Use Cases** ✅
- **5 regras de negócio** (RN-CAI-001 a 005) ✅
- **RBAC** por papel ✅
- **Frontend completo** (páginas, hooks, componentes) ✅

### Progresso Atual

| Fase | Status | Arquivos |
|------|--------|----------|
| 1. Database | ✅ Concluída | 4 migrations |
| 2. Domain | ✅ Concluída | 4 arquivos |
| 3. Infrastructure | ✅ Concluída | 4 arquivos |
| 4. Application | ✅ Concluída | 9 arquivos |
| 5. Interface | ✅ Concluída | 1 handler |
| 6. Integração | ✅ Concluída | Wire + Routes |
| 7. Frontend | ✅ Concluída | 13 arquivos |

---

## 🏗️ Fases de Implementação

### FASE 1: Database (Migrations)
**Estimativa:** 2-3 horas  
**Status:** ✅ CONCLUÍDA  
**Arquivos criados:**
- `backend/migrations/028_caixa_diario.up.sql`
- `backend/migrations/028_caixa_diario.down.sql`
- `backend/migrations/029_caixa_diario_fix_saldo_esperado.up.sql`
- `backend/migrations/029_caixa_diario_fix_saldo_esperado.down.sql`

#### 1.1 Tabela `caixa_diario`
```sql
-- Campos principais:
id, tenant_id, unidade_id (futuro)
usuario_abertura_id, usuario_fechamento_id
data_abertura, data_fechamento
saldo_inicial, total_entradas, total_saidas, total_sangrias, total_reforcos
saldo_esperado (calculado pela aplicação)
saldo_real, divergencia
status (ABERTO, FECHADO), justificativa_divergencia
created_at, updated_at
```

#### 1.2 Tabela `operacoes_caixa`
```sql
-- Campos principais:
id, caixa_id, tenant_id
tipo (VENDA, SANGRIA, REFORCO, DESPESA)
valor, descricao
destino (DEPOSITO, PAGAMENTO, COFRE) -- para sangrias
origem (TROCO, CAPITAL_GIRO, TRANSFERENCIA) -- para reforços
usuario_id
created_at
```

#### 1.3 Índices e Constraints
- `idx_caixa_tenant_status` → Busca rápida de caixa aberto
- `idx_caixa_aberto_unico` → Garante apenas 1 caixa ABERTO por tenant
- RLS policies para multi-tenant

#### Checklist Fase 1:
- [x] Criar migration 028_caixa_diario.up.sql
- [x] Criar migration 028_caixa_diario.down.sql
- [x] Criar migration 029 (fix saldo_esperado)
- [x] Executar migration (schema_migrations v29)

---

### FASE 2: Domain Layer
**Estimativa:** 3-4 horas  
**Status:** ✅ CONCLUÍDA  
**Arquivos criados:**
- `backend/internal/domain/entity/caixa_diario.go`
- `backend/internal/domain/entity/operacao_caixa.go`
- `backend/internal/domain/port/caixa_diario_repository.go`
- `backend/internal/domain/errors_caixa.go`

#### 2.1 Entidades

**`caixa_diario.go`** (207 linhas)
```go
type CaixaDiario struct {
    ID                       uuid.UUID
    TenantID                 uuid.UUID
    UnidadeID                *uuid.UUID
    UsuarioAberturaID        uuid.UUID
    UsuarioFechamentoID      *uuid.UUID
    DataAbertura             time.Time
    DataFechamento           *time.Time
    SaldoInicial             decimal.Decimal
    TotalEntradas            decimal.Decimal
    TotalSaidas              decimal.Decimal
    TotalSangrias            decimal.Decimal
    TotalReforcos            decimal.Decimal
    SaldoEsperado            decimal.Decimal
    SaldoReal                *decimal.Decimal
    Divergencia              *decimal.Decimal
    Status                   CaixaStatus
    JustificativaDivergencia *string
    CreatedAt                time.Time
    UpdatedAt                time.Time
}
```

**`operacao_caixa.go`** (87 linhas)
```go
type OperacaoCaixa struct {
    ID        uuid.UUID
    CaixaID   uuid.UUID
    TenantID  uuid.UUID
    Tipo      TipoOperacaoCaixa
    Valor     decimal.Decimal
    Descricao string
    Destino   *string
    Origem    *string
    UsuarioID uuid.UUID
    CreatedAt time.Time
}
```

#### 2.2 Repository Interface

```go
type CaixaDiarioRepository interface {
    Create(ctx context.Context, caixa *entity.CaixaDiario) error
    FindByID(ctx context.Context, tenantID, caixaID uuid.UUID) (*entity.CaixaDiario, error)
    FindAberto(ctx context.Context, tenantID uuid.UUID) (*entity.CaixaDiario, error)
    ExistsCaixaAberto(ctx context.Context, tenantID uuid.UUID) (bool, error)
    Update(ctx context.Context, caixa *entity.CaixaDiario) error
    UpdateTotais(ctx context.Context, tenantID, caixaID uuid.UUID, entradas, saidas, sangrias, reforcos, saldoEsperado decimal.Decimal) error
    Fechar(ctx context.Context, caixa *entity.CaixaDiario) error
    ListHistorico(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.CaixaDiario, error)
    CountHistorico(ctx context.Context, tenantID uuid.UUID) (int64, error)
    CreateOperacao(ctx context.Context, operacao *entity.OperacaoCaixa) error
    ListOperacoes(ctx context.Context, tenantID, caixaID uuid.UUID) ([]*entity.OperacaoCaixa, error)
    ListOperacoesByTipo(ctx context.Context, tenantID, caixaID uuid.UUID, tipo entity.TipoOperacaoCaixa) ([]*entity.OperacaoCaixa, error)
    SumOperacoesByTipo(ctx context.Context, tenantID, caixaID uuid.UUID) (map[entity.TipoOperacaoCaixa]decimal.Decimal, error)
}
```

#### 2.3 Domain Errors

```go
var (
    ErrCaixaJaAberto               = errors.New("já existe um caixa aberto")
    ErrCaixaNaoAberto              = errors.New("nenhum caixa aberto")
    ErrCaixaJaFechado              = errors.New("caixa já fechado")
    ErrCaixaJustificativaObrigatoria = errors.New("justificativa obrigatória para divergência maior que R$ 5,00")
    ErrValorInvalido               = errors.New("valor inválido")
    ErrSangriaDestinoObrigatorio   = errors.New("destino é obrigatório para sangria")
    ErrReforcoOrigemObrigatoria    = errors.New("origem é obrigatória para reforço")
)
```

#### Checklist Fase 2:
- [x] Criar entity/caixa_diario.go
- [x] Criar entity/operacao_caixa.go
- [x] Criar port/caixa_diario_repository.go
- [x] Adicionar erros em domain/errors_caixa.go

---

### FASE 3: Infrastructure Layer
**Estimativa:** 4-5 horas  
**Status:** ✅ CONCLUÍDA  
**Arquivos criados:**
- `backend/internal/infra/db/queries/caixa_diario.sql` (17 queries)
- `backend/internal/infra/db/schema/caixa_diario.sql`
- `backend/internal/infra/db/sqlc/caixa_diario.sql.go` (gerado)
- `backend/internal/infra/repository/postgres/caixa_diario_repository.go` (424 linhas)

#### 3.1 Queries SQLC (17 queries implementadas)

```sql
-- CAIXA DIÁRIO
-- name: CreateCaixaDiario :one
-- name: GetCaixaDiarioByID :one
-- name: GetCaixaDiarioAberto :one
-- name: ExistsCaixaAberto :one
-- name: UpdateCaixaDiario :one
-- name: UpdateCaixaDiarioTotais :exec
-- name: FecharCaixaDiario :one
-- name: ListCaixaDiarioHistorico :many
-- name: CountCaixaDiarioHistorico :one

-- OPERAÇÕES
-- name: CreateOperacaoCaixa :one
-- name: ListOperacoesByCaixa :many
-- name: ListOperacoesByCaixaAndTipo :many
-- name: SumOperacoesByTipo :many
-- name: GetLastOperacao :one
```

#### 3.2 Repository PostgreSQL
- ✅ Implementa `port.CaixaDiarioRepository` (verificado em compile-time)
- ✅ Mappers de modelo sqlc → domain entity
- ✅ Tratamento de erros específicos
- ✅ Suporte a filtros e paginação

#### Checklist Fase 3:
- [x] Criar queries em sqlc/queries/caixa_diario.sql
- [x] Criar schema em sqlc/schema/caixa_diario.sql
- [x] Rodar `sqlc generate`
- [x] Criar postgres/caixa_diario_repository.go
- [x] Verificar compile-time interface check

---

### FASE 4: Application Layer (Use Cases)
**Estimativa:** 5-6 horas  
**Status:** ✅ CONCLUÍDA  
**Arquivos criados:**
- `backend/internal/application/dto/caixa_dto.go` (10 DTOs)
- `backend/internal/application/mapper/caixa_mapper.go` (5 mappers)
- `backend/internal/application/usecase/caixa/abrir_caixa.go`
- `backend/internal/application/usecase/caixa/sangria.go`
- `backend/internal/application/usecase/caixa/reforco.go`
- `backend/internal/application/usecase/caixa/fechar_caixa.go`
- `backend/internal/application/usecase/caixa/get_caixa.go`
- `backend/internal/application/usecase/caixa/list_historico.go`
- `backend/internal/application/usecase/caixa/get_totais.go`

#### 4.1 DTOs Implementados

| DTO | Tipo | Descrição |
|-----|------|-----------|
| `AbrirCaixaRequest` | Request | Saldo inicial |
| `SangriaRequest` | Request | Valor, destino, descrição |
| `ReforcoRequest` | Request | Valor, origem, descrição |
| `FecharCaixaRequest` | Request | Saldo real, justificativa |
| `CaixaDiarioResponse` | Response | Dados completos do caixa |
| `OperacaoCaixaResponse` | Response | Dados de operação |
| `CaixaStatusResponse` | Response | Status resumido |
| `CaixaTotaisResponse` | Response | Totais por tipo de operação |
| `HistoricoCaixaResponse` | Response | Lista paginada de histórico |
| `ListHistoricoRequest` | Request | Filtros de paginação |

#### 4.2 Use Cases Implementados

| Use Case | Descrição | Regras |
|----------|-----------|--------|
| `AbrirCaixaUseCase` | Abre caixa com saldo inicial | RN-CAI-001: Apenas 1 aberto por tenant |
| `SangriaUseCase` | Registra retirada de caixa | RN-CAI-002: Destino obrigatório, atualiza totais |
| `ReforcoUseCase` | Registra adição ao caixa | RN-CAI-003: Origem obrigatória, atualiza totais |
| `FecharCaixaUseCase` | Fecha caixa com conferência | RN-CAI-004: Calcula divergência, justificativa se > R$5 |
| `GetCaixaUseCase` | Busca caixa por ID ou aberto | Retorna status atual |
| `ListHistoricoUseCase` | Histórico paginado | Filtros de data e paginação |
| `GetTotaisUseCase` | Totais por tipo de operação | Soma agrupada por tipo |

#### Checklist Fase 4:
- [x] Criar dto/caixa_dto.go (10 DTOs)
- [x] Criar mapper/caixa_mapper.go (5 mappers)
- [x] Criar usecase/caixa/abrir_caixa.go
- [x] Criar usecase/caixa/sangria.go
- [x] Criar usecase/caixa/reforco.go
- [x] Criar usecase/caixa/fechar_caixa.go
- [x] Criar usecase/caixa/get_caixa.go
- [x] Criar usecase/caixa/list_historico.go
- [x] Criar usecase/caixa/get_totais.go
- [ ] Testes unitários dos use cases (futuro)

---

### FASE 5: Interface Layer (HTTP Handlers)
**Estimativa:** 3-4 horas  
**Status:** ✅ CONCLUÍDA  
**Arquivos criados:**
- `backend/internal/infra/http/handler/caixa_handler.go` (9 endpoints com Swagger)

#### 5.1 Endpoints Implementados

| Método | Endpoint | Handler | RBAC |
|--------|----------|---------|------|
| `POST` | `/api/v1/caixa/abrir` | `AbrirCaixa` | owner, manager, employee |
| `GET` | `/api/v1/caixa/status` | `GetStatus` | owner, manager, employee |
| `GET` | `/api/v1/caixa/aberto` | `GetCaixaAberto` | owner, manager, employee |
| `GET` | `/api/v1/caixa/historico` | `ListHistorico` | owner, manager, accountant |
| `GET` | `/api/v1/caixa/totais` | `GetTotais` | owner, manager |
| `POST` | `/api/v1/caixa/sangria` | `RegistrarSangria` | owner, manager, employee* |
| `POST` | `/api/v1/caixa/reforco` | `RegistrarReforco` | owner, manager |
| `POST` | `/api/v1/caixa/fechar` | `FecharCaixa` | owner, manager, employee |
| `GET` | `/api/v1/caixa/:id` | `GetCaixaByID` | owner, manager, accountant |

*employee: limite R$200 para sangria

#### 5.2 Handler Structure

```go
type CaixaHandler struct {
    abrirUC      *caixa.AbrirCaixaUseCase
    fecharUC     *caixa.FecharCaixaUseCase
    sangriaUC    *caixa.SangriaUseCase
    reforcoUC    *caixa.ReforcoUseCase
    getCaixaUC   *caixa.GetCaixaUseCase
    historicoUC  *caixa.ListHistoricoUseCase
    totaisUC     *caixa.GetTotaisUseCase
    logger       *zap.Logger
}
```

#### Checklist Fase 5:
- [x] Criar handler/caixa_handler.go
- [x] Swagger annotations completas
- [x] Registrar rotas em router/routes.go
- [x] Rodar `swag init`
- [x] Testes de integração HTTP

---

### FASE 6: Integração e Wiring
**Estimativa:** 2-3 horas  
**Status:** ✅ CONCLUÍDA

#### 6.1 Tarefas de Integração

| Tarefa | Arquivo | Status |
|--------|---------|--------|
| Registrar rotas | `cmd/api/main.go` | ✅ Concluído |
| Wire repository | DI manual | ✅ Concluído |
| Wire use cases | DI manual | ✅ Concluído |
| Wire handler | DI manual | ✅ Concluído |
| Gerar Swagger | `swag init` | ✅ Concluído |

#### 6.2 Integração com Financeiro (Futuro)

Quando o caixa é fechado, precisa:

1. **Atualizar `fluxo_caixa_diario`**
   - Somar entradas confirmadas do dia
   - Atualizar saldo final

2. **Registrar Divergência**
   - Se `divergencia < 0`: Criar despesa "Quebra de Caixa"
   - Se `divergencia > 0`: Criar receita "Sobra de Caixa"

3. **Audit Log**
   - Registrar ação CAIXA_FECHADO com metadata

#### Checklist Fase 6:
- [x] Registrar rotas em cmd/api/main.go (linha 768)
- [x] Wire dependencies (repository → use cases → handler)
- [x] Rodar `swag init` para atualizar docs
- [x] Testar endpoints via Swagger UI
- [ ] Criar service/financeiro_service.go (interface) — futuro
- [ ] Implementar integração no FecharCaixaUseCase — futuro
- [ ] Atualizar fluxo_caixa_diario — futuro
- [x] Testes de integração HTTP (9/9 endpoints)

---

### FASE 7: Frontend
**Estimativa:** 8-10 horas  
**Status:** ✅ CONCLUÍDA

#### 7.1 Estrutura de Páginas

```
frontend/src/app/(dashboard)/caixa/
├── page.tsx                    # Tela principal (Operação de Caixa)
├── historico/
│   └── page.tsx               # Histórico de caixas fechados
├── [id]/
│   └── page.tsx               # Detalhes de um caixa específico
├── components/
│   ├── caixa-status-card.tsx  # Card de status atual
│   ├── extrato-dia.tsx        # Lista de operações
│   ├── saldo-cards.tsx        # Cards de saldo/sangria/reforço
│   ├── modal-abrir-caixa.tsx  # Modal de abertura
│   ├── modal-sangria.tsx      # Modal de sangria
│   ├── modal-reforco.tsx      # Modal de reforço
│   └── modal-fechar-caixa.tsx # Modal de fechamento
└── hooks/
    └── use-caixa.ts           # Queries e mutations integradas
```

**Arquivos criados:**
- `frontend/src/types/caixa.ts` - Tipos TypeScript
- `frontend/src/services/caixa-service.ts` - Serviço API
- `frontend/src/hooks/use-caixa.ts` - React Query hooks
- `frontend/src/components/caixa/` - 7 componentes
- `frontend/src/app/(dashboard)/caixa/` - 3 páginas

#### 7.2 Componentes Principais

**CaixaStatusCard.tsx**
- Mostra se caixa está aberto ou fechado
- Exibe operador e horário de abertura
- Botões de ação (Sangria, Reforço, Fechar)

**ExtratoDia.tsx**
- Lista cronológica de operações
- Ícones por tipo (abertura, venda, sangria, reforço)
- Valores com cores (+verde, -vermelho)

**FecharCaixaModal.tsx**
- Resumo do dia (saldo inicial, entradas, sangrias, reforços)
- Input para valor contado
- Cálculo automático de divergência
- Campo justificativa (obrigatório se divergência > R$5)

#### 7.3 API Types

```typescript
// types/caixa.ts
interface CaixaDiario {
  id: string;
  status: 'ABERTO' | 'FECHADO';
  saldo_inicial: string;
  total_entradas: string;
  total_sangrias: string;
  total_reforcos: string;
  saldo_esperado: string;
  saldo_real?: string;
  divergencia?: string;
  data_abertura: string;
  data_fechamento?: string;
  usuario_abertura_nome: string;
  usuario_fechamento_nome?: string;
}

interface OperacaoCaixa {
  id: string;
  tipo: 'VENDA' | 'SANGRIA' | 'REFORCO' | 'DESPESA';
  valor: string;
  descricao: string;
  destino?: string;
  origem?: string;
  usuario_nome: string;
  criado_em: string;
}
```

#### Checklist Fase 7:
- [x] Criar page.tsx (tela principal)
- [x] Criar CaixaStatusCard.tsx
- [x] Criar SaldoCards.tsx
- [x] Criar ExtratoDia.tsx
- [x] Criar ModalAbrirCaixa.tsx
- [x] Criar ModalSangria.tsx
- [x] Criar ModalReforco.tsx
- [x] Criar ModalFecharCaixa.tsx
- [x] Criar hooks (React Query) — use-caixa.ts
- [x] Criar tipos TypeScript — types/caixa.ts
- [x] Adicionar na sidebar (Banknote icon)
- [ ] Testes E2E — futuro

---

## 📊 Resumo de Esforço

| Fase | Descrição | Estimativa | Status |
|------|-----------|------------|--------|
| 1 | Database (Migrations) | 2-3h | ✅ Concluída |
| 2 | Domain Layer | 3-4h | ✅ Concluída |
| 3 | Infrastructure Layer | 4-5h | ✅ Concluída |
| 4 | Application Layer | 5-6h | ✅ Concluída |
| 5 | Interface Layer (HTTP) | 3-4h | ✅ Concluída |
| 6 | Integração e Wiring | 2-3h | ✅ Concluída |
| 7 | Frontend | 8-10h | ✅ Concluída |
| - | **TOTAL** | **27-35h** | ✅ 100% |

---

## 🔗 Dependências

### Pré-requisitos
- ✅ Tabela `tenants` (existe)
- ✅ Tabela `users` (existe)
- ✅ Tabela `fluxo_caixa_diario` (existe)
- ✅ Tabela `categorias` (existe)
- ✅ Sistema RBAC (existe)

### Dependências Futuras
- ⏳ Tabela `units` (unidades) — para multi-unidade
- ⏳ Módulo de Vendas (para registrar entradas automáticas)

---

## ✅ Critérios de Aceite

### Backend
- [x] Apenas 1 caixa aberto por tenant
- [x] Sangrias com destino obrigatório
- [x] Reforços com origem obrigatória
- [x] Divergência calculada corretamente
- [x] Justificativa obrigatória se divergência > R$5
- [ ] Audit log para todas operações — futuro
- [x] RBAC validado por endpoint
- [x] Multi-tenant isolado (RLS)

### Frontend
- [x] Tela responsiva
- [x] Validação de formulários (Zod)
- [x] Feedback visual de loading/error
- [x] Atualização automática após mutações
- [x] Formatação de moeda BR

### Integração (Futuro — v1.1.0)
- [ ] Divergência negativa → Despesa no DRE
- [ ] Divergência positiva → Receita no DRE
- [ ] Atualização do fluxo_caixa_diario

---

## 🎉 Conclusão

### Fase 7 — Frontend (CONCLUÍDA)

**Arquivos criados:**

#### Types
- `frontend/src/types/caixa.ts` — Tipos TypeScript espelhando DTOs

#### Service
- `frontend/src/services/caixa-service.ts` — Serviço API

#### Hooks
- `frontend/src/hooks/use-caixa.ts` — React Query hooks

#### Componentes
- `frontend/src/components/caixa/caixa-status-card.tsx` — Card de status
- `frontend/src/components/caixa/saldo-cards.tsx` — Cards de totais
- `frontend/src/components/caixa/extrato-dia.tsx` — Lista de operações
- `frontend/src/components/caixa/modal-abrir-caixa.tsx` — Modal abertura
- `frontend/src/components/caixa/modal-sangria.tsx` — Modal sangria
- `frontend/src/components/caixa/modal-reforco.tsx` — Modal reforço
- `frontend/src/components/caixa/modal-fechar-caixa.tsx` — Modal fechamento
- `frontend/src/components/caixa/index.ts` — Exports

#### Páginas
- `frontend/src/app/(dashboard)/caixa/page.tsx` — Página principal
- `frontend/src/app/(dashboard)/caixa/historico/page.tsx` — Histórico
- `frontend/src/app/(dashboard)/caixa/[id]/page.tsx` — Detalhes do caixa

#### Sidebar
- Item "Caixa Diário" adicionado ao menu Financeiro (ícone Banknote)

---

## 🧪 Testes de Integração HTTP

**Data:** 29/11/2025  
**Status:** ✅ 9/9 endpoints validados

| # | Endpoint | Método | Status | Descrição |
|---|----------|--------|--------|-----------|
| 1 | `/api/v1/caixa/abrir` | POST | ✅ 201 | Abrir caixa com saldo inicial |
| 2 | `/api/v1/caixa/status` | GET | ✅ 200 | Status aberto/fechado |
| 3 | `/api/v1/caixa/aberto` | GET | ✅ 200 | Dados do caixa atual |
| 4 | `/api/v1/caixa/:id` | GET | ✅ 200 | Detalhes por ID |
| 5 | `/api/v1/caixa/sangria` | POST | ✅ 201 | Registrar sangria |
| 6 | `/api/v1/caixa/reforco` | POST | ✅ 201 | Registrar reforço |
| 7 | `/api/v1/caixa/totais` | GET | ✅ 200 | Totais do caixa |
| 8 | `/api/v1/caixa/historico` | GET | ✅ 200 | Histórico paginado |
| 9 | `/api/v1/caixa/fechar` | POST | ✅ 200 | Fechar caixa |

### Correções Aplicadas Durante Testes

1. **`caixa_mapper.go`** — Proteção contra divisão por zero em `ToListCaixaHistoricoResponse`
2. **`caixa_handler.go`** — Defaults de paginação (page=1, pageSize=20)

---

## ✅ Módulo Caixa Diário — IMPLEMENTAÇÃO COMPLETA

**Total de arquivos criados/modificados:** 30+

| Camada | Arquivos | Descrição |
|--------|----------|-----------|
| Migrations | 4 | Tabelas caixa_diario + operacoes_caixa |
| Domain | 4 | Entidades, erros, interface de repositório |
| Infrastructure | 4 | SQLC queries, schema, repositório PostgreSQL |
| Application | 9 | DTOs, mappers, 7 use cases |
| Interface | 1 | Handler HTTP com 9 endpoints |
| Frontend | 13 | Types, service, hooks, 7 componentes, 3 páginas |

### Tarefas Futuras (v1.1.0)
- [ ] Testes unitários dos use cases
- [ ] Testes E2E frontend
- [ ] Audit log para operações
- [ ] Integração divergência → DRE
- [ ] Atualização automática do fluxo_caixa_diario

**Responsável:** Tech Lead  
**Revisão:** Product Owner  
**Última Atualização:** 29/11/2025 - Todas as 7 Fases Concluídas