> Criado em: 20/11/2025 20:43 (America/Sao_Paulo)

# 🏗️ Arquitetura Barber Analytics Pro v2.0

**Versão:** 2.0  
**Data Criação:** 14/11/2025  
**Última Revisão:** 22/11/2025  
**Status:** Em evolução (documento alinhado ao código atual)  
**Autor:** Arquiteto de Software Sr.

---

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Princípios Arquiteturais](#princípios-arquiteturais)
3. [Stack Tecnológico](#stack-tecnológico)
4. [Arquitetura em Camadas](#arquitetura-em-camadas)
5. [Estrutura de Diretórios](#estrutura-de-diretórios)
6. [Padrões de Design](#padrões-de-design)
7. [Fluxo de Dados](#fluxo-de-dados)
8. [Multi-Tenancy](#multi-tenancy)
9. [Segurança](#segurança)
10. [Escalabilidade](#escalabilidade)
11. [Estado Atual vs Planejado](#estado-atual-vs-planejado)

---

## 🎯 Visão Geral

O Barber Analytics Pro v2.0 é uma plataforma SaaS modular e escalável para gerenciamento completo de barbearias, construída com **Clean Architecture**, **Domain-Driven Design (DDD)** e aderência aos princípios **SOLID**.

### Objetivos Arquiteturais

- ✅ **Independência de Framework**: Lógica de negócio desacoplada de ferramentas
- ✅ **Testabilidade**: Código altamente testável em todos os níveis
- ✅ **Manutenibilidade**: Estrutura clara e padrões consistentes
- ✅ **Escalabilidade**: Suporte a múltiplos tenants e crescimento horizontal
- ✅ **Performance**: Otimizações em queries, cache e processamento assíncrono
- ✅ **Segurança**: Isolamento de dados, auditoria e compliance

---

## 🏛️ Princípios Arquiteturais

### 1. Clean Architecture

```
┌─────────────────────────────────────────┐
│       Presentation Layer (HTTP/UI)      │
├─────────────────────────────────────────┤
│      Application Layer (Use Cases)      │
├─────────────────────────────────────────┤
│    Domain Layer (Business Rules)        │
├─────────────────────────────────────────┤
│  Infrastructure Layer (DB, APIs, etc)   │
└─────────────────────────────────────────┘
```

**Direção de dependências:** Centro (Domain) → Externo (Infrastructure)

### 2. Domain-Driven Design (DDD)

- **Ubiquitous Language**: Linguagem de negócio consistente
- **Bounded Contexts**: Módulos independentes (Financeiro, Assinaturas, Estoque, Lista da Vez)
- **Aggregates**: Entidades relacionadas com raízes claras
- **Value Objects**: Objetos imutáveis sem identidade
- **Repositories**: Abstração de persistência por Aggregate

### 3. SOLID Principles

| Princípio | Aplicação |
|-----------|-----------|
| **S** - SRP | Cada classe tem uma única responsabilidade |
| **O** - OCP | Aberto para extensão, fechado para modificação |
| **L** - LSP | Subtypes são substituíveis por seus tipos base |
| **I** - ISP | Interfaces específicas ao cliente |
| **D** - DIP | Dependências em abstrações, não em implementações |

---

## 🛠️ Stack Tecnológico

### Backend

```yaml
Linguagem: Go 1.24.0
Framework HTTP: Echo v4
Query: SQLC (type-safe SQL)
Autenticação: JWT RS256 (planejado) — hoje mock de tenant em header
Validação: go-playground/validator/v10 (não configurado no server ainda)
Scheduler: robfig/cron v3 (jobs de DRE/Fluxo/Compensações)
Logger: Zap (JSON estruturado)
Trace: OpenTelemetry (planejado)
```

### Banco de Dados

```yaml
Principal: PostgreSQL 14+
Provedor Recomendado: Neon (serverless, backup automático)
Alternativa: Supabase (DB-only mode)
Migrations: golang-migrate/migrate
Backup: Automático (Neon/Supabase) + snapshots periódicos
```

### Frontend (MVP -> V2)

```yaml
Framework: Next.js (App Router)
State Management: TanStack Query (React Query)
UI: MUI + shadcn/ui (mix atual)
Styling: CSS modules + tokens locais (Tailwind não está em uso no repo)
Form Validation: Zod + React Hook Form
```

### DevOps & Infraestrutura

```yaml
Reverse Proxy: NGINX (SSL/TLS via Certbot)
CI/CD: GitHub Actions
Logs & Monitoring: Grafana + Prometheus
APM: Sentry (para exceções e performance)
Hosting: VPS Ubuntu 22.04 LTS
```

---

## 🏗️ Arquitetura em Camadas

### Backend Go (Clean Architecture)

```
backend/
├── cmd/
│   └── api/
│       └── main.go                    # Entrypoint
├── internal/
│   ├── config/                        # Leitura de env
│   ├── domain/                        # Business logic (entities, value objects)
│   ├── application/
│   │   ├── dto/                       # Data Transfer Objects
│   │   ├── mapper/                    # Domain <-> DTO mapping
│   │   └── usecase/                   # Application use cases
│   ├── infrastructure/
│   │   ├── http/                      # HTTP handlers e middlewares
│   │   ├── repository/                # Database repositories
│   │   ├── external/                  # Integrações externas (Asaas, etc)
│   │   └── scheduler/                 # Cron jobs
│   └── ports/                         # Interfaces (abstrações)
├── migrations/                        # SQL migrations
├── tests/                            # Testes integrados
└── go.mod
```

### Camada de Domínio (Domain Layer)

```go
// Entidade - Aggregate Root
type Barbearia struct {
    ID            string
    Nome          string
    CNPJ          string
    Endereco      Endereco           // Value Object
    Barbeiros     []Barbeiro         // Child entities
    Configuracoes Configuracoes      // Value Object
    CriadoEm      time.Time
    AtualizadoEm  time.Time
}

// Entidade - Lista da Vez (Novo Módulo)
type BarbersTurnList struct {
    ID             string
    TenantID       string
    ProfessionalID string
    CurrentPoints  int
    LastTurnAt     time.Time
    IsActive       bool
}

// Value Object - Imutável
type Endereco struct {
    Rua       string
    Numero    int
    Complemento string
    Cidade    string
    UF        string
    CEP       string
}

// Repository Interface (Port)
type BarbeariaRepository interface {
    Save(ctx context.Context, barbearia *Barbearia) error
    FindByID(ctx context.Context, id string) (*Barbearia, error)
    FindByTenantID(ctx context.Context, tenantID string) (*Barbearia, error)
}
```

### Camada de Aplicação (Application Layer)

```go
// Use Case real (financeiro)
type CreateContaPagarUseCase struct {
    repo   port.ContaPagarRepository
    logger *zap.Logger
}

func (uc *CreateContaPagarUseCase) Execute(ctx context.Context, input CreateContaPagarInput) (*entity.ContaPagar, error) {
    if input.TenantID == "" {
        return nil, domain.ErrTenantIDRequired
    }
    conta, err := entity.NewContaPagar(
        input.TenantID,
        input.Descricao,
        input.CategoriaID,
        input.Fornecedor,
        input.Valor,
        input.Tipo,
        input.DataVencimento,
    )
    if err != nil {
        return nil, err
    }
    if err := uc.repo.Create(ctx, conta); err != nil {
        uc.logger.Error("erro ao criar conta pagar", zap.Error(err))
        return nil, err
    }
    return conta, nil
}

// DTO - entrada
type CreateContaPagarInput struct {
    TenantID       string
    Descricao      string
    CategoriaID    string
    Fornecedor     string
    Valor          valueobject.Money
    Tipo           valueobject.TipoCusto
    DataVencimento time.Time
}
```

### Camada de Apresentação (HTTP/Delivery Layer)

```go
// Handler (trecho real de FinancialHandler)
func (h *FinancialHandler) CreateContaPagar(c echo.Context) error {
    ctx := c.Request().Context()
    tenantID, _ := c.Get("tenant_id").(string) // hoje mockado; futuro: JWT

    var req dto.CreateContaPagarRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request"})
    }
    if err := c.Validate(&req); err != nil {
        return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: err.Error()})
    }
    valor, tipo, data, err := mapper.FromCreateContaPagarRequest(req)
    if err != nil {
        return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "conversion_error", Message: err.Error()})
    }
    conta, err := h.createContaPagarUC.Execute(ctx, financial.CreateContaPagarInput{
        TenantID: tenantID, Descricao: req.Descricao, CategoriaID: req.CategoriaID,
        Fornecedor: req.Fornecedor, Valor: valor, Tipo: tipo, DataVencimento: data,
    })
    if err != nil {
        return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
    }
    return c.JSON(http.StatusCreated, mapper.ToContaPagarResponse(conta))
}
```

### Camada de Infraestrutura (Infrastructure Layer)

```go
// Repository Implementation (simplificado)
type PostgresContaPagarRepository struct {
    db *sql.DB
}

func (r *PostgresContaPagarRepository) Create(ctx context.Context, conta *entity.ContaPagar) error {
    query := `
        INSERT INTO contas_a_pagar (id, tenant_id, descricao, valor, data_vencimento, status)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err := r.db.ExecContext(ctx, query,
        conta.ID, conta.TenantID, conta.Descricao,
        conta.Valor.ToDecimal(), conta.DataVencimento, conta.Status,
    )
    return err
}
```

---

## 📂 Estrutura de Diretórios

```
barber-analytics-proV2/
│
├── backend/                        # Backend em Go
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── domain/                # entity, valueobject, port
│   │   ├── application/           # dto, mapper, usecase
│   │   └── infra/                 # http/handler, repository/postgres, scheduler, metrics
│   ├── internal/infra/db/schema   # SQLC schemas/migrations
│   └── go.mod
│
├── frontend/                       # Frontend Next.js (App Router)
│   ├── app/                        # layouts, pages, providers
│   ├── components/                 # layout/ui/cookie banner
│   ├── hooks/                      # React Query hooks (financeiro/metas/precificação)
│   └── lib/services                # clients e schemas zod
│
├── docs/                           # Documentação
│   ├── 02-arquitetura/ARQUITETURA.md (este arquivo)
│   ├── PRD-NEXO.md
│   ├── ROADMAP_MILITAR_NEXO.md
│   └── ...
│
├── .github/workflows/              # CI/CD, backup, testes
└── README.md
```

---

## 🎨 Padrões de Design

### 1. Repository Pattern

Abstração para persistência de dados:

```go
// Port (Interface)
type ContaPagarRepository interface {
    Create(ctx context.Context, conta *entity.ContaPagar) error
    FindByID(ctx context.Context, tenantID, id string) (*entity.ContaPagar, error)
    List(ctx context.Context, tenantID string, filters port.ContaPagarListFilters) ([]*entity.ContaPagar, error)
}

// Adapter (Implementação)
type PostgresContaPagarRepository struct { ... }
```

### 2. Dependency Injection

Injeção de dependências no startup:

```go
func InitializeFinancialHandler(dbPool *pgxpool.Pool, logger *zap.Logger) *handler.FinancialHandler {
    queries := db.New(dbPool)
    contaPagarRepo := postgres.NewContaPagarRepository(queries)
    createUC := financial.NewCreateContaPagarUseCase(contaPagarRepo, logger)
    // ... instanciar demais use cases
    return handler.NewFinancialHandler(createUC, /* outros UCs */, logger)
}
```

### 3. DTO (Data Transfer Object)

Separação entre modelo de domínio e dados transmitidos:

```go
// Domain
type ContaPagar struct {
    ID     string
    Valor  valueobject.Money
    Status valueobject.StatusConta
}

// DTO
type ContaPagarResponse struct {
    ID          string `json:"id"`
    Valor       string `json:"valor"`   // formatado
    Status      string `json:"status"`
    DataVencimento string `json:"data_vencimento"`
}
```

### 4. Middleware Chain

Middleware para cross-cutting concerns:

```go
app.Use(middleware.Logger())
app.Use(middleware.Recovery())
app.Use(middleware.CORSMiddleware())
// TODO: registrar validator, auth JWT + tenant middleware
```

### 5. Service Locator (Opcional)

Para inicialização centralizadas:

```go
type Container struct {
    DB              *sql.DB
    Logger          *zap.Logger
    ContaPagarRepo  port.ContaPagarRepository
    ContaReceberRepo port.ContaReceberRepository
    // ... outros services
}
```

---

## 🔄 Fluxo de Dados

### Fluxo de Requisição HTTP

```
Request HTTP
    ↓
NGINX (Rate Limit, SSL)
    ↓
Echo Router
    ↓
Middleware Chain
  ├── Logger
  ├── Recovery
  ├── Auth (JWT)
  └── Tenant Context
    ↓
Handler (HTTP Layer)
    ├── Bind Request
    ├── Validate Input (Validator)
    └── Call Use Case
    ↓
Use Case (Application Layer)
    ├── Business Logic Validation
    ├── Call Domain Services
    └── Call Repositories
    ↓
Domain Layer
    ├── Business Rules
    ├── Value Object Creation
    └── Entity Validation
    ↓
Repository (Infrastructure)
    └── Database Query (SQLC)
    ↓
Response DTO
    ↓
JSON Response
```

### Fluxo de Processamento Assíncrono (Cron)

```
Scheduler (robfig/cron)
    ↓
Cron Job (ex: GenerateDREMonthly)
    ↓
Use Case (Application Layer)
    ├── Ler contas pagar/receber do período
    ├── Calcular DRE ou Fluxo Diário
    └── Persistir no DB
    ↓
Notificação (opcional)
    └── Log ou Webhook
```

---

## 👥 Multi-Tenancy

### Modelo Selecionado: Column-Based (Tenant per Row)

**Razão**: Simplicidade, segurança, sem complexidade de schema separados.  
**Estado atual**: field `tenant_id` presente, mas middleware ainda é mock (header); falta JWT/RLS.

### Implementação

1. **Coluna tenant_id em todas as tabelas**

```sql
CREATE TABLE contas_a_pagar (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    descricao VARCHAR(255) NOT NULL,
    valor NUMERIC(18, 2) NOT NULL,
    data_vencimento DATE NOT NULL,
    status VARCHAR(20) NOT NULL,
    criado_em TIMESTAMP DEFAULT NOW(),
    UNIQUE(id, tenant_id)
);

CREATE INDEX idx_contas_pagar_tenant_id ON contas_a_pagar(tenant_id);
CREATE INDEX idx_contas_pagar_data ON contas_a_pagar(tenant_id, data_vencimento);
```

2. **Middleware de Tenant**

```go
func TenantMiddleware(c echo.Context) error {
    token := c.Get("user").(*jwt.Token)
    claims := token.Claims.(jwt.MapClaims)
    
    tenantID := claims["tenant_id"].(string)
    c.Set("tenant_id", tenantID)
    
    return c.Next()
}
```

3. **Query Segura**

```go
func (r *PostgresContaPagarRepository) ListByDateRange(
    ctx context.Context, tenantID string, from, to time.Time) ([]*entity.ContaPagar, error) {
    query := `
        SELECT id, tenant_id, descricao, valor, data_vencimento, status
        FROM contas_a_pagar
        WHERE tenant_id = $1 AND data_vencimento BETWEEN $2 AND $3
        ORDER BY data_vencimento DESC
    `
    // ...
    return r.db.QueryContext(ctx, query, tenantID, from, to)
}
```

---

## 🔐 Segurança

### Autenticação

- **Planejado:** JWT RS256 + refresh/rotação
- **Estado atual:** sem auth; tenant vem de header mock para desenvolvimento

### Autorização

- **Planejado:** RBAC por role (Owner, Manager, Employee, Accountant)
- **Estado atual:** inexistente; rotas financeiras expostas sem checagem

### Isolamento de Dados

- **Campo `tenant_id` obrigatório** em entidades e queries
- **Falta:** RLS no banco, audit logs, enforcement em middleware

### Rate Limiting

- **NGINX**: 100 req/s por IP
- **Aplicação**: 50 req/min por endpoint sensível

### HTTPS/TLS

- **Certificados**: Let's Encrypt + Certbot
- **HSTS**: 1 ano
- **CSP**: Restritivo para frontend

---

## 📈 Escalabilidade

### Banco de Dados

- **Índices estratégicos** em `tenant_id`, datas, status
- **Particionamento** de tabelas largas (receitas, despesas) por ano
- **Connection pooling** via pgBouncer (futuro)
- **Read replicas** no Neon (futuro)

### Backend

- **Stateless API** (escalável horizontalmente)
- **Cache de leitura** (Redis, futuro) para dashboards
- **Bulk operations** com batch inserts
- **Async jobs** fora do request cycle

### Frontend

- **Code splitting** automático no Next.js
- **Image optimization** com next/image
- **CDN** para assets estáticos
- **ISR** (Incremental Static Regeneration) para dashboards

### Monitoramento

- **Prometheus** para métricas
- **Grafana** para dashboards
- **Alertas** para SLA violations
- **Logs centralizados** em Loki ou Datadog

---

## 🧭 Estado Atual vs Planejado

| Área                    | Estado em 22/11/2025                                                    | Planejado / Gap                                        |
| ----------------------- | ----------------------------------------------------------------------- | ------------------------------------------------------ |
| Autenticação/RBAC       | Mock de tenant via header; sem JWT/RBAC                                 | JWT RS256 + roles + middleware de tenant               |
| Validator Echo          | Uso de `c.Validate` nos handlers, mas validator não registrado no `main`| Registrar validator global                             |
| Módulos implementados   | Financeiro (payables/receivables/compensação/fluxo/DRE), Metas, Precificação; User prefs parcial | Agendamento, Lista da vez, Comissões, Estoque, CRM, Asaas |
| Repositórios            | Aggregates financeiros com SQLC; `SumByPeriod` e filtros agregados retornam zero (placeholder) | Implementar agregações e filtros completos             |
| Cron/Scheduler          | Jobs DRE/Fluxo/Compensações registrados; tenants via env estática       | Provider real de tenants, jobs de comissões/estoque    |
| LGPD                    | Handlers/UC de export/delete incompletos; rota não exposta              | Integrar rotas `/me/preferences|export|delete` + audit |
| Frontend                | App Router com layout/dashboard básico; hooks/services prontos          | Páginas de Financeiro, Metas, Precificação, Agenda etc. |
| Multi-tenant segurança  | Sem RLS ou enforcement além do campo tenant_id                          | RLS ou validação estrita em todas as queries/handlers  |

> Este quadro deve ser revisado a cada checkpoint (Roadmap Militar). Sempre atualizar o documento quando um gap for fechado.

---

## 🔗 Referências

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design - Eric Evans](https://www.domainlanguage.com/ddd/)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [PostgreSQL Best Practices](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [Echo Framework](https://echo.labstack.com/)
- [Go Best Practices](https://golang.org/doc/effective_go)

---

**Última Atualização:** 22/11/2025  
**Status:** ✅ Alinhado ao estado atual (com gaps mapeados)
