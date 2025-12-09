# 🔧 Sprint 2: Backend Core — Módulo Assinaturas

**Sprint:** 2 de 5  
**Status:** ✅ CONCLUÍDO  
**Progresso:** 100%  
**Estimativa:** 20-25 horas  
**Prioridade:** 🔴 CRÍTICA  
**Dependência:** ✅ Sprint 1 (Banco de Dados) deve estar concluída

---

## 📚 Referência Obrigatória

> ⚠️ **ANTES DE INICIAR**, leia completamente:
> 
> - **[FLUXO_ASSINATURA.md](../../docs/11-Fluxos/FLUXO_ASSINATURA.md)** — Fonte da verdade
>   - Seção 2: Página Planos (campos, regras PL-001 a PL-005)
>   - Seção 3: Página Assinantes (status, formas de pagamento)
>   - Seção 6: Fluxos Detalhados (lógica de cada operação)
>   - Seção 8: Regras de Negócio (RN-SUB-*, RN-VENC-*, RN-CANC-*, RN-BEN-*)
> - **[RBAC.md](../../docs/06-seguranca/RBAC.md)** — Permissões por role
> - **[API_PUBLICA.md](../../docs/04-backend/API_PUBLICA.md)** — Padrões de endpoint

---

## 📊 Progresso das Tarefas

| ID | Tarefa | Estimativa | Status | Progresso |
|----|--------|------------|--------|-----------|
| **Entidades** |
| BE-001 | Entidade `Plan` | 30min | ✅ Concluído | 100% |
| BE-002 | Entidade `Subscription` | 1h | ✅ Concluído | 100% |
| BE-003 | Entidade `SubscriptionPayment` | 20min | ✅ Concluído | 100% |
| **Portas/Interfaces** |
| BE-004 | Interface `PlanRepository` | 15min | ✅ Concluído | 100% |
| BE-005 | Interface `SubscriptionRepository` | 30min | ✅ Concluído | 100% |
| **Repositórios** |
| BE-006 | Implementação `PlanRepository` | 1.5h | ✅ Concluído | 100% |
| BE-007 | Implementação `SubscriptionRepository` | 2h | ✅ Concluído | 100% |
| **DTOs** |
| BE-008 | DTOs de Plan | 30min | ✅ Concluído | 100% |
| BE-009 | DTOs de Subscription | 45min | ✅ Concluído | 100% |
| **Use Cases** |
| BE-010 | UseCase: CreatePlan | 30min | ✅ Concluído | 100% |
| BE-011 | UseCase: UpdatePlan | 30min | ✅ Concluído | 100% |
| BE-012 | UseCase: ListPlans | 30min | ✅ Concluído | 100% |
| BE-013 | UseCase: DeactivatePlan | 30min | ✅ Concluído | 100% |
| BE-014 | UseCase: CreateSubscription | 2h | ✅ Concluído | 100% |
| BE-015 | UseCase: CancelSubscription | 1h | ✅ Concluído | 100% |
| BE-016 | UseCase: RenewSubscription | 1h | ✅ Concluído | 100% |
| BE-017 | UseCase: GetSubscriptionMetrics | 30min | ✅ Concluído | 100% |
| **Handlers** |
| BE-018 | Handler: Plans CRUD | 2h | ✅ Concluído | 100% |
| BE-019 | Handler: Subscriptions CRUD | 2h | ✅ Concluído | 100% |
| BE-020 | Handler: Subscription Actions | 1h | ✅ Concluído | 100% |
| **Cron Job** |
| BE-021 | CronJob: Verificar Vencimentos | 1h | ✅ Concluído | 100% |
| **Wire/DI** |
| BE-022 | Configurar Wire (Dependency Injection) | 1h | ✅ Concluído | 100% |

**📈 PROGRESSO SPRINT: 22/22 (100%)**

---

## 📋 Tarefas Detalhadas

### 🏗️ FASE 1: Entidades de Domínio

#### BE-001: Entidade `Plan`

**Objetivo:** Criar entidade de domínio para Planos

**Referência:** [FLUXO_ASSINATURA.md — Seção 2.2 e 2.3](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#22-campos-do-cadastro)

**Arquivo:** `backend/internal/domain/plan/plan.go`

```go
package plan

import (
    "time"

    "github.com/google/uuid"
)

type Periodicidade string

const (
    PeriodicidadeMensal Periodicidade = "MENSAL"
)

type Plan struct {
    ID              uuid.UUID
    TenantID        uuid.UUID
    Nome            string
    Descricao       *string
    Valor           string // Decimal como string (ex: "99.90")
    Periodicidade   Periodicidade
    QtdServicos     *int
    LimiteUsoMensal *int
    Ativo           bool
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// Regras de negócio do FLUXO_ASSINATURA.md
// REGRA PL-003: Não pode excluir plano com assinaturas ativas
func (p *Plan) CanBeDeleted(activeSubscriptionsCount int) bool {
    return activeSubscriptionsCount == 0
}

// REGRA PL-002: Plano inativo NÃO aparece na seleção
func (p *Plan) IsAvailableForSelection() bool {
    return p.Ativo
}
```

**Critérios de Aceite:**
- [x] Struct com todos os campos do FLUXO_ASSINATURA.md
- [x] Regras PL-002 e PL-003 implementadas como métodos
- [x] Valor como string (nunca float para dinheiro)

---

#### BE-002: Entidade `Subscription`

**Objetivo:** Criar entidade de domínio para Assinaturas

**Referência:** [FLUXO_ASSINATURA.md — Seção 3 e 8](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#3-página-assinantes)

**Arquivo:** `backend/internal/domain/subscription/subscription.go`

```go
package subscription

import (
    "time"

    "github.com/google/uuid"
)

type Status string

const (
    StatusAguardandoPagamento Status = "AGUARDANDO_PAGAMENTO"
    StatusAtivo               Status = "ATIVO"
    StatusInadimplente        Status = "INADIMPLENTE"
    StatusInativo             Status = "INATIVO"
    StatusCancelado           Status = "CANCELADO"
)

type FormaPagamento string

const (
    FormaPagamentoCartao   FormaPagamento = "CARTAO"
    FormaPagamentoPix      FormaPagamento = "PIX"
    FormaPagamentoDinheiro FormaPagamento = "DINHEIRO"
)

type Subscription struct {
    ID                   uuid.UUID
    TenantID             uuid.UUID
    ClienteID            uuid.UUID
    PlanoID              uuid.UUID
    AsaasCustomerID      *string
    AsaasSubscriptionID  *string
    FormaPagamento       FormaPagamento
    Status               Status
    Valor                string // Decimal como string
    LinkPagamento        *string
    CodigoTransacao      *string
    DataAtivacao         *time.Time
    DataVencimento       *time.Time
    DataCancelamento     *time.Time
    CanceladoPor         *uuid.UUID
    ServicosUtilizados   int
    CreatedAt            time.Time
    UpdatedAt            time.Time

    // Campos de JOIN (não persistidos)
    PlanoNome      string
    ClienteNome    string
    ClienteTelefone string
}

// RN-BEN-001: Cliente com assinatura ativa pode usar serviços
func (s *Subscription) CanUseServices() bool {
    return s.Status == StatusAtivo
}

// RN-BEN-003: Se atingir limite, bloquear uso
func (s *Subscription) HasReachedServiceLimit(planoLimite *int) bool {
    if planoLimite == nil {
        return false // Ilimitado
    }
    return s.ServicosUtilizados >= *planoLimite
}

// RN-VENC-004: Vencida há mais de 3 dias
func (s *Subscription) ShouldBecomeInadimplente() bool {
    if s.DataVencimento == nil || s.Status != StatusAtivo {
        return false
    }
    return time.Now().After(s.DataVencimento.AddDate(0, 0, 3))
}

// RN-CANC-004: Cancelada não pode ser reativada
func (s *Subscription) CanBeReactivated() bool {
    return s.Status != StatusCancelado
}

// Verifica se precisa de integração Asaas
func (s *Subscription) RequiresAsaasIntegration() bool {
    return s.FormaPagamento == FormaPagamentoCartao
}

// RN-BEN-005: Verifica se cliente deve ser marcado como assinante
func (s *Subscription) ShouldMarkClientAsSubscriber() bool {
    return s.Status == StatusAtivo
}
```

**Critérios de Aceite:**
- [x] Todos os status do FLUXO_ASSINATURA.md implementados
- [x] Todas as formas de pagamento implementadas
- [x] Regras de negócio RN-BEN-*, RN-VENC-*, RN-CANC-* como métodos
- [x] Campos de JOIN para evitar N+1 queries
- [x] Método para verificar se deve atualizar flag is_subscriber do cliente

---

#### BE-003: Entidade `SubscriptionPayment`

**Arquivo:** `backend/internal/domain/subscription/subscription_payment.go`

```go
package subscription

import (
    "time"

    "github.com/google/uuid"
)

type PaymentStatus string

const (
    PaymentStatusPendente   PaymentStatus = "PENDENTE"
    PaymentStatusConfirmado PaymentStatus = "CONFIRMADO"
    PaymentStatusEstornado  PaymentStatus = "ESTORNADO"
    PaymentStatusCancelado  PaymentStatus = "CANCELADO"
)

type SubscriptionPayment struct {
    ID              uuid.UUID
    TenantID        uuid.UUID
    SubscriptionID  uuid.UUID
    AsaasPaymentID  *string
    Valor           string
    FormaPagamento  FormaPagamento
    Status          PaymentStatus
    DataPagamento   *time.Time
    CodigoTransacao *string
    Observacao      *string
    CreatedAt       time.Time
}
```

---

### 🔌 FASE 2: Portas/Interfaces

#### BE-004: Interface `PlanRepository`

**Arquivo:** `backend/internal/domain/plan/repository.go`

```go
package plan

import (
    "context"

    "github.com/google/uuid"
)

type Repository interface {
    Create(ctx context.Context, plan *Plan) error
    GetByID(ctx context.Context, id, tenantID uuid.UUID) (*Plan, error)
    ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*Plan, error)
    ListActiveByTenant(ctx context.Context, tenantID uuid.UUID) ([]*Plan, error)
    Update(ctx context.Context, plan *Plan) error
    Deactivate(ctx context.Context, id, tenantID uuid.UUID) error
    CountActiveSubscriptions(ctx context.Context, planID, tenantID uuid.UUID) (int, error)
}
```

---

#### BE-005: Interface `SubscriptionRepository`

**Arquivo:** `backend/internal/domain/subscription/repository.go`

```go
package subscription

import (
    "context"
    "time"

    "github.com/google/uuid"
)

type Repository interface {
    // CRUD
    Create(ctx context.Context, sub *Subscription) error
    GetByID(ctx context.Context, id, tenantID uuid.UUID) (*Subscription, error)
    ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*Subscription, error)
    ListByStatus(ctx context.Context, tenantID uuid.UUID, status Status) ([]*Subscription, error)
    
    // Status management
    UpdateStatus(ctx context.Context, id, tenantID uuid.UUID, status Status) error
    Activate(ctx context.Context, id, tenantID uuid.UUID, dataAtivacao, dataVencimento time.Time) error
    Cancel(ctx context.Context, id, tenantID uuid.UUID, canceladoPor uuid.UUID) error
    
    // Asaas integration
    UpdateAsaasIDs(ctx context.Context, id, tenantID uuid.UUID, customerID, subscriptionID, linkPagamento string) error
    GetByAsaasSubscriptionID(ctx context.Context, asaasSubscriptionID string) (*Subscription, error)
    
    // Service usage (RN-BEN-002)
    IncrementServicosUtilizados(ctx context.Context, id, tenantID uuid.UUID) error
    ResetServicosUtilizados(ctx context.Context, id, tenantID uuid.UUID) error
    
    // Cron job queries
    ListOverdue(ctx context.Context, tenantID uuid.UUID) ([]*Subscription, error)
    ListExpiringSoon(ctx context.Context, tenantID uuid.UUID, days int) ([]*Subscription, error)
    
    // Metrics (Seção 5.1 do FLUXO_ASSINATURA.md)
    GetMetrics(ctx context.Context, tenantID uuid.UUID) (*SubscriptionMetrics, error)
    
    // Payments
    CreatePayment(ctx context.Context, payment *SubscriptionPayment) error
    ListPaymentsBySubscription(ctx context.Context, subscriptionID, tenantID uuid.UUID) ([]*SubscriptionPayment, error)
}

type SubscriptionMetrics struct {
    TotalAtivas       int
    TotalInativas     int
    TotalInadimplentes int
    ReceitaMensal     string
}
```

---

### 🗃️ FASE 3: Repositórios

#### BE-006: Implementação `PlanRepository`

**Arquivo:** `backend/internal/infrastructure/database/plan_repository.go`

**Critérios de Aceite:**
- [x] Implementar todos os métodos da interface
- [x] Usar queries sqlc geradas
- [x] Nunca fazer query sem tenant_id
- [x] Mappers entre sqlc types e domain types

---

#### BE-007: Implementação `SubscriptionRepository`

**Arquivo:** `backend/internal/infrastructure/database/subscription_repository.go`

**Critérios de Aceite:**
- [x] Implementar todos os métodos da interface
- [x] Usar queries sqlc geradas
- [x] Nunca fazer query sem tenant_id
- [x] Tratamento de campos opcionais (nullables)

---

### 📦 FASE 4: DTOs

#### BE-008: DTOs de Plan

**Referência:** [FLUXO_ASSINATURA.md — Seção 2.2](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#22-campos-do-cadastro)

**Arquivo:** `backend/internal/application/dto/plan_dto.go`

```go
package dto

type CreatePlanRequest struct {
    Nome            string  `json:"nome" validate:"required,min=3,max=100"`
    Descricao       *string `json:"descricao,omitempty" validate:"omitempty,max=500"`
    Valor           string  `json:"valor" validate:"required"`
    QtdServicos     *int    `json:"qtd_servicos,omitempty" validate:"omitempty,gte=0"`
    LimiteUsoMensal *int    `json:"limite_uso_mensal,omitempty" validate:"omitempty,gte=0"`
}

type UpdatePlanRequest struct {
    Nome            string  `json:"nome" validate:"required,min=3,max=100"`
    Descricao       *string `json:"descricao,omitempty" validate:"omitempty,max=500"`
    Valor           string  `json:"valor" validate:"required"`
    QtdServicos     *int    `json:"qtd_servicos,omitempty" validate:"omitempty,gte=0"`
    LimiteUsoMensal *int    `json:"limite_uso_mensal,omitempty" validate:"omitempty,gte=0"`
    Ativo           bool    `json:"ativo"`
}

type PlanResponse struct {
    ID              string  `json:"id"`
    Nome            string  `json:"nome"`
    Descricao       *string `json:"descricao,omitempty"`
    Valor           string  `json:"valor"`
    Periodicidade   string  `json:"periodicidade"`
    QtdServicos     *int    `json:"qtd_servicos,omitempty"`
    LimiteUsoMensal *int    `json:"limite_uso_mensal,omitempty"`
    Ativo           bool    `json:"ativo"`
    CreatedAt       string  `json:"created_at"`
    UpdatedAt       string  `json:"updated_at"`
}
```

**⚠️ Regras dos DTOs (copilot-instructions.md):**
- `valor` como string (nunca float para dinheiro)
- Sem `tenant_id` no payload (extraído do contexto)
- `omitempty` em campos opcionais
- snake_case no JSON

---

#### BE-009: DTOs de Subscription

**Referência:** [FLUXO_ASSINATURA.md — Seção 6.1, 6.2, 6.3](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#6-fluxos-detalhados)

**Arquivo:** `backend/internal/application/dto/subscription_dto.go`

```go
package dto

type CreateSubscriptionRequest struct {
    ClienteID       string  `json:"cliente_id" validate:"required,uuid"`
    PlanoID         string  `json:"plano_id" validate:"required,uuid"`
    FormaPagamento  string  `json:"forma_pagamento" validate:"required,oneof=CARTAO PIX DINHEIRO"`
    CodigoTransacao *string `json:"codigo_transacao,omitempty"` // Para PIX
}

type RenewSubscriptionRequest struct {
    FormaPagamento  string  `json:"forma_pagamento" validate:"required,oneof=PIX DINHEIRO"`
    CodigoTransacao *string `json:"codigo_transacao,omitempty"`
}

type SubscriptionResponse struct {
    ID                  string  `json:"id"`
    ClienteID           string  `json:"cliente_id"`
    ClienteNome         string  `json:"cliente_nome"`
    ClienteTelefone     string  `json:"cliente_telefone"`
    PlanoID             string  `json:"plano_id"`
    PlanoNome           string  `json:"plano_nome"`
    FormaPagamento      string  `json:"forma_pagamento"`
    Status              string  `json:"status"`
    Valor               string  `json:"valor"`
    LinkPagamento       *string `json:"link_pagamento,omitempty"`
    DataAtivacao        *string `json:"data_ativacao,omitempty"`
    DataVencimento      *string `json:"data_vencimento,omitempty"`
    ServicosUtilizados  int     `json:"servicos_utilizados"`
    CreatedAt           string  `json:"created_at"`
}

type SubscriptionMetricsResponse struct {
    TotalAtivas        int    `json:"total_ativas"`
    TotalInativas      int    `json:"total_inativas"`
    TotalInadimplentes int    `json:"total_inadimplentes"`
    ReceitaMensal      string `json:"receita_mensal"`
}
```

---

### 🎯 FASE 5: Use Cases

#### BE-010 a BE-013: Use Cases de Plan

**Arquivos:**
- `backend/internal/application/usecase/plan/create_plan.go`
- `backend/internal/application/usecase/plan/update_plan.go`
- `backend/internal/application/usecase/plan/list_plans.go`
- `backend/internal/application/usecase/plan/deactivate_plan.go`

**Regras do FLUXO_ASSINATURA.md a implementar:**
- **PL-003:** Não pode excluir plano com assinaturas ativas
- **PL-005:** Nome do plano deve ser único por tenant

---

#### BE-014: UseCase `CreateSubscription`

**Referência:** [FLUXO_ASSINATURA.md — Seção 6.1, 6.2, 6.3](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#61-fluxo-nova-assinatura-cartão-de-crédito)

**Arquivo:** `backend/internal/application/usecase/subscription/create_subscription.go`

**Lógica:**
```
1. Validar cliente existe
2. Validar plano existe e está ativo (PL-002)
3. Verificar se cliente já tem assinatura ativa do mesmo plano (RN-SUB-004)
4. Se forma = CARTAO:
   - Verificar se cliente já tem asaas_customer_id
   - Se não, buscar/criar no Asaas (RN-CLI-002)
   - Chamar AsaasGateway para criar assinatura
   - Gerar link de pagamento
   - Status = AGUARDANDO_PAGAMENTO
5. Se forma = PIX ou DINHEIRO:
   - Status = ATIVO
   - Data vencimento = now() + 30 dias
   - Marcar cliente como is_subscriber = true (RN-CLI-003)
6. Salvar assinatura
7. Se forma = DINHEIRO: registrar no caixa como entrada
```

---

#### BE-015: UseCase `CancelSubscription`

**Referência:** [FLUXO_ASSINATURA.md — Seção 6.5](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#65-fluxo-cancelar-assinatura)

**Regras:**
- **RN-CANC-001:** Apenas Admin/Gerente podem cancelar
- **RN-CANC-002:** Se forma = CARTAO, cancelar também no Asaas
- **RN-CANC-003:** Registrar quem cancelou e quando
- **RN-CLI-004:** Verificar se cliente possui outras assinaturas ativas; se não, remover flag is_subscriber

---

#### BE-016: UseCase `RenewSubscription`

**Referência:** [FLUXO_ASSINATURA.md — Seção 6.4](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#64-fluxo-renovar-assinatura-manual-pixdinheiro)

**Lógica:**
```
1. Validar assinatura existe
2. Validar status permite renovação (não pode ser CANCELADO - RN-CANC-004)
3. Atualizar status = ATIVO
4. Data ativação = now()
5. Data vencimento = now() + 30 dias
6. Resetar servicos_utilizados = 0 (RN-BEN-004)
7. Registrar pagamento no histórico
```

---

#### BE-017: UseCase `GetSubscriptionMetrics`

**Referência:** [FLUXO_ASSINATURA.md — Seção 5.1](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#51-métricas-exibidas)

**Métricas:**
- Total ativas
- Total inativas
- Total inadimplentes
- Receita mensal (soma dos valores ativos)

---

### 🌐 FASE 6: Handlers HTTP

#### BE-018: Handler Plans CRUD

**Arquivo:** `backend/internal/interfaces/http/handlers/plan_handler.go`

**Endpoints:**
```
GET    /api/v1/plans           → ListPlans
GET    /api/v1/plans/:id       → GetPlanByID
POST   /api/v1/plans           → CreatePlan
PUT    /api/v1/plans/:id       → UpdatePlan
DELETE /api/v1/plans/:id       → DeactivatePlan
```

**RBAC:**
- GET: Admin, Gerente, Recepção (visualizar)
- POST/PUT/DELETE: Admin, Gerente (modificar)

---

#### BE-019: Handler Subscriptions CRUD

**Arquivo:** `backend/internal/interfaces/http/handlers/subscription_handler.go`

**Endpoints:**
```
GET    /api/v1/subscriptions                → ListSubscriptions
GET    /api/v1/subscriptions/:id            → GetSubscriptionByID
POST   /api/v1/subscriptions                → CreateSubscription
GET    /api/v1/subscriptions/metrics        → GetMetrics
```

**RBAC:**
- GET: Admin, Gerente, Recepção
- POST: Admin, Gerente, Recepção

---

#### BE-020: Handler Subscription Actions

**Arquivo:** (mesmo arquivo do BE-019)

**Endpoints:**
```
POST   /api/v1/subscriptions/:id/renew      → RenewSubscription
DELETE /api/v1/subscriptions/:id            → CancelSubscription
```

**RBAC:**
- POST renew: Admin, Gerente, Recepção
- DELETE: Admin, Gerente (RN-CANC-001)

---

### ⏰ FASE 7: Cron Job

#### BE-021: CronJob Verificar Vencimentos

**Referência:** [FLUXO_ASSINATURA.md — Seção 8.2](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#82-regras-de-vencimento)

**Arquivo:** `backend/internal/interfaces/cron/subscription_cron.go`

**Lógica:**
```
Executar diariamente às 00:05 (RN-VENC-003)

1. Buscar assinaturas com status = ATIVO
2. Para cada assinatura:
   - Se data_vencimento + 3 dias < now():
     - Atualizar status = INADIMPLENTE (RN-VENC-004)
     - Verificar se cliente tem outras assinaturas ativas (RN-CLI-004)
     - Se não, atualizar is_subscriber = false
3. Log quantidade de assinaturas atualizadas
```

---

### 🔗 FASE 8: Wire/DI

#### BE-022: Configurar Wire

**Arquivo:** `backend/cmd/api/wire.go`

**Adicionar:**
- PlanRepository
- SubscriptionRepository
- Plan Use Cases
- Subscription Use Cases
- Plan Handler
- Subscription Handler
- Subscription Cron Job

---

## ✅ Critérios de Conclusão da Sprint

- [x] 3 entidades de domínio criadas
- [x] 2 interfaces de repositório definidas
- [x] 2 repositórios implementados com sqlc
- [x] DTOs para Plan e Subscription
- [x] 8 use cases implementados
- [x] 3 handlers com todos os endpoints
- [x] Cron job de vencimentos
- [x] Wire configurado
- [x] Testes unitários para use cases críticos

---

## 🔗 Próxima Sprint

Após conclusão, iniciar **Sprint 3: Integração Asaas**
📂 [03-INTEGRACAO-ASAAS.md](./03-INTEGRACAO-ASAAS.md)

---

**FIM DO DOCUMENTO — SPRINT 2: BACKEND CORE**
