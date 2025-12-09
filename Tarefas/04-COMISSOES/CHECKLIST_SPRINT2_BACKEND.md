# ✅ CHECKLIST — SPRINT 2: DOMAIN + REPOSITORY + USECASES

> **Status:** ❌ Não Iniciado  
> **Dependência:** Sprint 1 (Migrations + Queries)  
> **Esforço Estimado:** 20 horas  
> **Prioridade:** P0 — Bloqueia API

---

## 📊 RESUMO

```
░░░░░░░░░░░░░░░░░░░░ 0% COMPLETO
```

| Categoria | Completo | Pendente |
|-----------|:--------:|:--------:|
| Domain Entities | 0/4 | 4 |
| Value Objects | 0/3 | 3 |
| Repository Interfaces | 0/4 | 4 |
| Repository Impl | 0/4 | 4 |
| UseCases | 0/15 | 15 |
| DTOs | 0/6 | 6 |

---

## 1️⃣ DOMAIN LAYER — ENTITIES

### 1.1 Entity: `CommissionRule` (Esforço: 1h)

- [ ] Criar `backend/internal/domain/entity/commission_rule.go`

```go
package entity

import (
    "encoding/json"
    "errors"
    "time"

    "github.com/google/uuid"
)

// CommissionType define o tipo de comissão
type CommissionType string

const (
    CommissionTypePercentage  CommissionType = "PERCENTAGE"
    CommissionTypeFixed       CommissionType = "FIXED"
    CommissionTypeHybrid      CommissionType = "HYBRID"
    CommissionTypeProgressive CommissionType = "PROGRESSIVE"
)

// ProgressiveTier representa uma faixa de comissão progressiva
type ProgressiveTier struct {
    Min float64 `json:"min"`
    Max float64 `json:"max,omitempty"`
    Pct float64 `json:"pct"`
}

// CommissionRule representa uma regra de comissão
type CommissionRule struct {
    ID             uuid.UUID        `json:"id"`
    TenantID       uuid.UUID        `json:"tenant_id"`
    UnitID         *uuid.UUID       `json:"unit_id,omitempty"`
    ProfessionalID *uuid.UUID       `json:"professional_id,omitempty"`
    ServiceID      *uuid.UUID       `json:"service_id,omitempty"`
    Type           CommissionType   `json:"type"`
    Value          float64          `json:"value"`
    FixedValue     float64          `json:"fixed_value"`
    Tiers          []ProgressiveTier `json:"tiers,omitempty"`
    Priority       int              `json:"priority"`
    Active         bool             `json:"active"`
    CreatedAt      time.Time        `json:"created_at"`
    UpdatedAt      time.Time        `json:"updated_at"`
    CreatedBy      *uuid.UUID       `json:"created_by,omitempty"`
    UpdatedBy      *uuid.UUID       `json:"updated_by,omitempty"`
}

// Validate valida a regra de comissão
func (r *CommissionRule) Validate() error {
    if r.TenantID == uuid.Nil {
        return errors.New("tenant_id é obrigatório")
    }
    
    switch r.Type {
    case CommissionTypePercentage:
        if r.Value < 0 || r.Value > 100 {
            return errors.New("percentual deve ser entre 0 e 100")
        }
    case CommissionTypeFixed:
        if r.Value < 0 {
            return errors.New("valor fixo não pode ser negativo")
        }
    case CommissionTypeHybrid:
        if r.FixedValue < 0 {
            return errors.New("valor fixo não pode ser negativo")
        }
        if r.Value < 0 || r.Value > 100 {
            return errors.New("percentual deve ser entre 0 e 100")
        }
    case CommissionTypeProgressive:
        if len(r.Tiers) == 0 {
            return errors.New("faixas progressivas são obrigatórias")
        }
    default:
        return errors.New("tipo de comissão inválido")
    }
    
    return nil
}

// Calculate calcula a comissão baseado no valor base
func (r *CommissionRule) Calculate(baseValue float64) float64 {
    switch r.Type {
    case CommissionTypePercentage:
        return baseValue * (r.Value / 100)
    case CommissionTypeFixed:
        return r.Value
    case CommissionTypeHybrid:
        return r.FixedValue + (baseValue * (r.Value / 100))
    case CommissionTypeProgressive:
        return r.calculateProgressive(baseValue)
    default:
        return 0
    }
}

func (r *CommissionRule) calculateProgressive(baseValue float64) float64 {
    for _, tier := range r.Tiers {
        if baseValue >= tier.Min && (tier.Max == 0 || baseValue < tier.Max) {
            return baseValue * (tier.Pct / 100)
        }
    }
    // Se não encontrou faixa, usa a última
    if len(r.Tiers) > 0 {
        lastTier := r.Tiers[len(r.Tiers)-1]
        return baseValue * (lastTier.Pct / 100)
    }
    return 0
}

// Specificity retorna o nível de especificidade da regra (maior = mais específica)
func (r *CommissionRule) Specificity() int {
    score := 0
    if r.ServiceID != nil {
        score += 4
    }
    if r.ProfessionalID != nil {
        score += 2
    }
    if r.UnitID != nil {
        score += 1
    }
    return score
}
```

#### Checklist

- [ ] Struct CommissionRule
- [ ] CommissionType enum
- [ ] ProgressiveTier struct
- [ ] Método Validate()
- [ ] Método Calculate()
- [ ] Método Specificity()

---

### 1.2 Entity: `CommissionPeriod` (Esforço: 1h)

- [ ] Criar `backend/internal/domain/entity/commission_period.go`

```go
package entity

import (
    "errors"
    "time"

    "github.com/google/uuid"
)

// PeriodStatus define o status do período
type PeriodStatus string

const (
    PeriodStatusDraft  PeriodStatus = "DRAFT"
    PeriodStatusClosed PeriodStatus = "CLOSED"
    PeriodStatusPaid   PeriodStatus = "PAID"
)

// CommissionPeriod representa um período de fechamento de comissão
type CommissionPeriod struct {
    ID              uuid.UUID    `json:"id"`
    TenantID        uuid.UUID    `json:"tenant_id"`
    UnitID          *uuid.UUID   `json:"unit_id,omitempty"`
    ProfessionalID  uuid.UUID    `json:"professional_id"`
    StartDate       time.Time    `json:"start_date"`
    EndDate         time.Time    `json:"end_date"`
    TotalServices   float64      `json:"total_services"`
    TotalProducts   float64      `json:"total_products"`
    TotalCommission float64      `json:"total_commission"`
    TotalBonus      float64      `json:"total_bonus"`
    TotalDeductions float64      `json:"total_deductions"`
    NetValue        float64      `json:"net_value"`
    QtyServices     int          `json:"qty_services"`
    QtyProducts     int          `json:"qty_products"`
    Status          PeriodStatus `json:"status"`
    BillID          *uuid.UUID   `json:"bill_id,omitempty"`
    Notes           string       `json:"notes,omitempty"`
    ClosedAt        *time.Time   `json:"closed_at,omitempty"`
    ClosedBy        *uuid.UUID   `json:"closed_by,omitempty"`
    PaidAt          *time.Time   `json:"paid_at,omitempty"`
    CreatedAt       time.Time    `json:"created_at"`
    UpdatedAt       time.Time    `json:"updated_at"`
}

// Validate valida o período
func (p *CommissionPeriod) Validate() error {
    if p.TenantID == uuid.Nil {
        return errors.New("tenant_id é obrigatório")
    }
    if p.ProfessionalID == uuid.Nil {
        return errors.New("professional_id é obrigatório")
    }
    if p.EndDate.Before(p.StartDate) {
        return errors.New("data final deve ser maior ou igual à data inicial")
    }
    return nil
}

// CanClose verifica se o período pode ser fechado
func (p *CommissionPeriod) CanClose() bool {
    return p.Status == PeriodStatusDraft
}

// CanPay verifica se o período pode ser marcado como pago
func (p *CommissionPeriod) CanPay() bool {
    return p.Status == PeriodStatusClosed
}

// Close fecha o período
func (p *CommissionPeriod) Close(closedBy uuid.UUID, billID uuid.UUID) error {
    if !p.CanClose() {
        return errors.New("período não pode ser fechado no status atual")
    }
    now := time.Now()
    p.Status = PeriodStatusClosed
    p.ClosedAt = &now
    p.ClosedBy = &closedBy
    p.BillID = &billID
    p.UpdatedAt = now
    return nil
}

// MarkAsPaid marca o período como pago
func (p *CommissionPeriod) MarkAsPaid() error {
    if !p.CanPay() {
        return errors.New("período não pode ser marcado como pago no status atual")
    }
    now := time.Now()
    p.Status = PeriodStatusPaid
    p.PaidAt = &now
    p.UpdatedAt = now
    return nil
}

// CalculateNetValue calcula o valor líquido
func (p *CommissionPeriod) CalculateNetValue() {
    p.NetValue = p.TotalCommission + p.TotalBonus - p.TotalDeductions
}
```

#### Checklist

- [ ] Struct CommissionPeriod
- [ ] PeriodStatus enum
- [ ] Método Validate()
- [ ] Método CanClose()
- [ ] Método CanPay()
- [ ] Método Close()
- [ ] Método MarkAsPaid()
- [ ] Método CalculateNetValue()

---

### 1.3 Entity: `Advance` (Esforço: 1h)

- [ ] Criar `backend/internal/domain/entity/advance.go`

```go
package entity

import (
    "errors"
    "time"

    "github.com/google/uuid"
)

// AdvanceStatus define o status do adiantamento
type AdvanceStatus string

const (
    AdvanceStatusPending  AdvanceStatus = "PENDING"
    AdvanceStatusApproved AdvanceStatus = "APPROVED"
    AdvanceStatusRejected AdvanceStatus = "REJECTED"
    AdvanceStatusDeducted AdvanceStatus = "DEDUCTED"
)

// Advance representa um adiantamento/vale
type Advance struct {
    ID              uuid.UUID      `json:"id"`
    TenantID        uuid.UUID      `json:"tenant_id"`
    UnitID          *uuid.UUID     `json:"unit_id,omitempty"`
    ProfessionalID  uuid.UUID      `json:"professional_id"`
    Amount          float64        `json:"amount"`
    RequestDate     time.Time      `json:"request_date"`
    Reason          string         `json:"reason,omitempty"`
    Status          AdvanceStatus  `json:"status"`
    ApprovedAt      *time.Time     `json:"approved_at,omitempty"`
    ApprovedBy      *uuid.UUID     `json:"approved_by,omitempty"`
    RejectedAt      *time.Time     `json:"rejected_at,omitempty"`
    RejectedBy      *uuid.UUID     `json:"rejected_by,omitempty"`
    RejectionReason string         `json:"rejection_reason,omitempty"`
    DeductedIn      *uuid.UUID     `json:"deducted_in,omitempty"`
    DeductedAt      *time.Time     `json:"deducted_at,omitempty"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
    CreatedBy       *uuid.UUID     `json:"created_by,omitempty"`
}

// Validate valida o adiantamento
func (a *Advance) Validate() error {
    if a.TenantID == uuid.Nil {
        return errors.New("tenant_id é obrigatório")
    }
    if a.ProfessionalID == uuid.Nil {
        return errors.New("professional_id é obrigatório")
    }
    if a.Amount <= 0 {
        return errors.New("valor deve ser maior que zero")
    }
    return nil
}

// CanApprove verifica se pode aprovar
func (a *Advance) CanApprove() bool {
    return a.Status == AdvanceStatusPending
}

// CanReject verifica se pode rejeitar
func (a *Advance) CanReject() bool {
    return a.Status == AdvanceStatusPending
}

// CanDeduct verifica se pode deduzir
func (a *Advance) CanDeduct() bool {
    return a.Status == AdvanceStatusApproved
}

// Approve aprova o adiantamento
func (a *Advance) Approve(approvedBy uuid.UUID) error {
    if !a.CanApprove() {
        return errors.New("adiantamento não pode ser aprovado no status atual")
    }
    now := time.Now()
    a.Status = AdvanceStatusApproved
    a.ApprovedAt = &now
    a.ApprovedBy = &approvedBy
    a.UpdatedAt = now
    return nil
}

// Reject rejeita o adiantamento
func (a *Advance) Reject(rejectedBy uuid.UUID, reason string) error {
    if !a.CanReject() {
        return errors.New("adiantamento não pode ser rejeitado no status atual")
    }
    now := time.Now()
    a.Status = AdvanceStatusRejected
    a.RejectedAt = &now
    a.RejectedBy = &rejectedBy
    a.RejectionReason = reason
    a.UpdatedAt = now
    return nil
}

// Deduct deduz o adiantamento em um período
func (a *Advance) Deduct(periodID uuid.UUID) error {
    if !a.CanDeduct() {
        return errors.New("adiantamento não pode ser deduzido no status atual")
    }
    now := time.Now()
    a.Status = AdvanceStatusDeducted
    a.DeductedIn = &periodID
    a.DeductedAt = &now
    a.UpdatedAt = now
    return nil
}
```

#### Checklist

- [ ] Struct Advance
- [ ] AdvanceStatus enum
- [ ] Método Validate()
- [ ] Método CanApprove()
- [ ] Método CanReject()
- [ ] Método CanDeduct()
- [ ] Método Approve()
- [ ] Método Reject()
- [ ] Método Deduct()

---

### 1.4 Ajuste Entity: `BarberCommission` (Esforço: 0.5h)

- [ ] Atualizar `backend/internal/domain/entity/barber_commission.go`

```go
// Adicionar campos
type BarberCommission struct {
    // ... campos existentes ...
    CommandItemID      *uuid.UUID `json:"command_item_id,omitempty"`
    CommissionPeriodID *uuid.UUID `json:"commission_period_id,omitempty"`
    UnitID             *uuid.UUID `json:"unit_id,omitempty"`
}

// Adicionar método
func (bc *BarberCommission) MarkAsProcessed(periodID uuid.UUID) {
    bc.Status = "PROCESSADO"
    bc.CommissionPeriodID = &periodID
}

func (bc *BarberCommission) MarkAsPaid() {
    bc.Status = "PAGO"
}
```

#### Checklist

- [ ] Adicionar CommandItemID
- [ ] Adicionar CommissionPeriodID
- [ ] Adicionar UnitID
- [ ] Método MarkAsProcessed()
- [ ] Método MarkAsPaid()

---

## 2️⃣ DOMAIN LAYER — VALUE OBJECTS

### 2.1 Value Objects (Esforço: 0.5h)

- [ ] Criar/Atualizar `backend/internal/domain/valueobject/commission_enums.go`

```go
package valueobject

// CommissionBaseMode define a base de cálculo
type CommissionBaseMode string

const (
    CommissionBaseModeGross CommissionBaseMode = "GROSS_TOTAL"
    CommissionBaseModeTable CommissionBaseMode = "TABLE_PRICE"
    CommissionBaseModeNet   CommissionBaseMode = "NET_VALUE"
)

// IsValid valida o modo
func (m CommissionBaseMode) IsValid() bool {
    switch m {
    case CommissionBaseModeGross, CommissionBaseModeTable, CommissionBaseModeNet:
        return true
    default:
        return false
    }
}
```

#### Checklist

- [ ] CommissionBaseMode
- [ ] Método IsValid()

---

## 3️⃣ REPOSITORY LAYER — INTERFACES

### 3.1 Interface: `CommissionRuleRepository` (Esforço: 0.5h)

- [ ] Criar `backend/internal/domain/repository/commission_rule_repository.go`

```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "nexo/internal/domain/entity"
)

type CommissionRuleRepository interface {
    Create(ctx context.Context, rule *entity.CommissionRule) error
    GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.CommissionRule, error)
    List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.CommissionRule, error)
    ListActive(ctx context.Context, tenantID uuid.UUID) ([]*entity.CommissionRule, error)
    ListByUnit(ctx context.Context, tenantID, unitID uuid.UUID) ([]*entity.CommissionRule, error)
    ListByProfessional(ctx context.Context, tenantID, professionalID uuid.UUID) ([]*entity.CommissionRule, error)
    ListByService(ctx context.Context, tenantID, serviceID uuid.UUID) ([]*entity.CommissionRule, error)
    FindApplicable(ctx context.Context, tenantID uuid.UUID, unitID, professionalID, serviceID *uuid.UUID) (*entity.CommissionRule, error)
    Update(ctx context.Context, rule *entity.CommissionRule) error
    Toggle(ctx context.Context, tenantID, id uuid.UUID) (*entity.CommissionRule, error)
    Delete(ctx context.Context, tenantID, id uuid.UUID) error
    Count(ctx context.Context, tenantID uuid.UUID) (int64, error)
}
```

#### Checklist

- [ ] Interface definida
- [ ] Todos os métodos

---

### 3.2 Interface: `CommissionPeriodRepository` (Esforço: 0.5h)

- [ ] Criar `backend/internal/domain/repository/commission_period_repository.go`

```go
package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "nexo/internal/domain/entity"
)

type CommissionPeriodRepository interface {
    Create(ctx context.Context, period *entity.CommissionPeriod) error
    GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.CommissionPeriod, error)
    List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.CommissionPeriod, error)
    ListByProfessional(ctx context.Context, tenantID, professionalID uuid.UUID, limit, offset int) ([]*entity.CommissionPeriod, error)
    ListByUnit(ctx context.Context, tenantID, unitID uuid.UUID, limit, offset int) ([]*entity.CommissionPeriod, error)
    ListByStatus(ctx context.Context, tenantID uuid.UUID, status entity.PeriodStatus, limit, offset int) ([]*entity.CommissionPeriod, error)
    ListByDateRange(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) ([]*entity.CommissionPeriod, error)
    GetDraft(ctx context.Context, tenantID, professionalID uuid.UUID) (*entity.CommissionPeriod, error)
    Update(ctx context.Context, period *entity.CommissionPeriod) error
    Close(ctx context.Context, period *entity.CommissionPeriod) error
    MarkAsPaid(ctx context.Context, tenantID, id uuid.UUID) error
    Delete(ctx context.Context, tenantID, id uuid.UUID) error
    SumCommissions(ctx context.Context, tenantID, professionalID uuid.UUID, startDate, endDate time.Time) (float64, int, error)
}
```

#### Checklist

- [ ] Interface definida
- [ ] Todos os métodos

---

### 3.3 Interface: `AdvanceRepository` (Esforço: 0.5h)

- [ ] Criar `backend/internal/domain/repository/advance_repository.go`

```go
package repository

import (
    "context"

    "github.com/google/uuid"
    "nexo/internal/domain/entity"
)

type AdvanceRepository interface {
    Create(ctx context.Context, advance *entity.Advance) error
    GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Advance, error)
    List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.Advance, error)
    ListByProfessional(ctx context.Context, tenantID, professionalID uuid.UUID, limit, offset int) ([]*entity.Advance, error)
    ListByStatus(ctx context.Context, tenantID uuid.UUID, status entity.AdvanceStatus, limit, offset int) ([]*entity.Advance, error)
    ListPending(ctx context.Context, tenantID uuid.UUID) ([]*entity.Advance, error)
    ListApprovedNotDeducted(ctx context.Context, tenantID, professionalID uuid.UUID) ([]*entity.Advance, error)
    Approve(ctx context.Context, advance *entity.Advance) error
    Reject(ctx context.Context, advance *entity.Advance) error
    Deduct(ctx context.Context, advance *entity.Advance) error
    Delete(ctx context.Context, tenantID, id uuid.UUID) error
    SumApprovedNotDeducted(ctx context.Context, tenantID, professionalID uuid.UUID) (float64, error)
}
```

#### Checklist

- [ ] Interface definida
- [ ] Todos os métodos

---

### 3.4 Ajuste Interface: `BarberCommissionRepository` (Esforço: 0.5h)

- [ ] Atualizar `backend/internal/domain/repository/barber_commission_repository.go`

```go
// Adicionar métodos
type BarberCommissionRepository interface {
    // ... métodos existentes ...
    
    CreateFromCommand(ctx context.Context, commission *entity.BarberCommission) error
    ListPendingByProfessional(ctx context.Context, tenantID, professionalID uuid.UUID) ([]*entity.BarberCommission, error)
    ListPendingByPeriod(ctx context.Context, tenantID, professionalID uuid.UUID, startDate, endDate time.Time) ([]*entity.BarberCommission, error)
    MarkAsProcessed(ctx context.Context, tenantID, professionalID uuid.UUID, periodID uuid.UUID, startDate, endDate time.Time) error
    MarkAsPaid(ctx context.Context, periodID uuid.UUID) error
}
```

#### Checklist

- [ ] CreateFromCommand
- [ ] ListPendingByProfessional
- [ ] ListPendingByPeriod
- [ ] MarkAsProcessed
- [ ] MarkAsPaid

---

## 4️⃣ REPOSITORY LAYER — IMPLEMENTATIONS

### 4.1 PostgreSQL: `PGCommissionRuleRepository` (Esforço: 2h)

- [ ] Criar `backend/internal/infra/repository/pg_commission_rule_repository.go`

#### Checklist

- [ ] Implementar todos os métodos da interface
- [ ] Usar queries sqlc
- [ ] Converter entre entity e sqlc types
- [ ] Tratamento de erros adequado

---

### 4.2 PostgreSQL: `PGCommissionPeriodRepository` (Esforço: 2h)

- [ ] Criar `backend/internal/infra/repository/pg_commission_period_repository.go`

#### Checklist

- [ ] Implementar todos os métodos da interface
- [ ] Usar queries sqlc
- [ ] Converter entre entity e sqlc types
- [ ] Tratamento de erros adequado

---

### 4.3 PostgreSQL: `PGAdvanceRepository` (Esforço: 1.5h)

- [ ] Criar `backend/internal/infra/repository/pg_advance_repository.go`

#### Checklist

- [ ] Implementar todos os métodos da interface
- [ ] Usar queries sqlc
- [ ] Converter entre entity e sqlc types
- [ ] Tratamento de erros adequado

---

### 4.4 Ajuste PostgreSQL: `PGBarberCommissionRepository` (Esforço: 1h)

- [ ] Atualizar `backend/internal/infra/repository/pg_barber_commission_repository.go`

#### Checklist

- [ ] Implementar novos métodos
- [ ] CreateFromCommand
- [ ] ListPendingByProfessional
- [ ] ListPendingByPeriod
- [ ] MarkAsProcessed
- [ ] MarkAsPaid

---

## 5️⃣ APPLICATION LAYER — DTOs

### 5.1 DTOs de Comissão (Esforço: 1h)

- [ ] Criar `backend/internal/application/dto/commission_rule_dto.go`
- [ ] Criar `backend/internal/application/dto/commission_period_dto.go`
- [ ] Criar `backend/internal/application/dto/advance_dto.go`

```go
// commission_rule_dto.go
type CreateCommissionRuleRequest struct {
    UnitID         *string `json:"unit_id,omitempty"`
    ProfessionalID *string `json:"professional_id,omitempty"`
    ServiceID      *string `json:"service_id,omitempty"`
    Type           string  `json:"type" validate:"required,oneof=PERCENTAGE FIXED HYBRID PROGRESSIVE"`
    Value          string  `json:"value" validate:"required"`
    FixedValue     string  `json:"fixed_value,omitempty"`
    Tiers          []ProgressiveTierDTO `json:"tiers,omitempty"`
    Priority       int     `json:"priority"`
    Active         bool    `json:"active"`
}

type CommissionRuleResponse struct {
    ID             string `json:"id"`
    UnitID         string `json:"unit_id,omitempty"`
    ProfessionalID string `json:"professional_id,omitempty"`
    ServiceID      string `json:"service_id,omitempty"`
    Type           string `json:"type"`
    Value          string `json:"value"`
    FixedValue     string `json:"fixed_value,omitempty"`
    Priority       int    `json:"priority"`
    Active         bool   `json:"active"`
    CreatedAt      string `json:"created_at"`
}

// commission_period_dto.go
type CreatePeriodPreviewRequest struct {
    UnitID         *string `json:"unit_id,omitempty"`
    ProfessionalID string  `json:"professional_id" validate:"required,uuid"`
    StartDate      string  `json:"start_date" validate:"required,datetime=2006-01-02"`
    EndDate        string  `json:"end_date" validate:"required,datetime=2006-01-02"`
}

type CommissionPeriodResponse struct {
    ID              string `json:"id"`
    ProfessionalID  string `json:"professional_id"`
    StartDate       string `json:"start_date"`
    EndDate         string `json:"end_date"`
    TotalCommission string `json:"total_commission"`
    TotalBonus      string `json:"total_bonus"`
    TotalDeductions string `json:"total_deductions"`
    NetValue        string `json:"net_value"`
    Status          string `json:"status"`
}

type ClosePeriodRequest struct {
    Notes string `json:"notes,omitempty"`
}

// advance_dto.go
type CreateAdvanceRequest struct {
    UnitID         *string `json:"unit_id,omitempty"`
    ProfessionalID string  `json:"professional_id" validate:"required,uuid"`
    Amount         string  `json:"amount" validate:"required"`
    RequestDate    string  `json:"request_date,omitempty"`
    Reason         string  `json:"reason,omitempty"`
}

type AdvanceResponse struct {
    ID             string `json:"id"`
    ProfessionalID string `json:"professional_id"`
    Amount         string `json:"amount"`
    RequestDate    string `json:"request_date"`
    Reason         string `json:"reason,omitempty"`
    Status         string `json:"status"`
    ApprovedAt     string `json:"approved_at,omitempty"`
}

type RejectAdvanceRequest struct {
    Reason string `json:"reason" validate:"required"`
}
```

#### Checklist DTOs

- [ ] CreateCommissionRuleRequest
- [ ] UpdateCommissionRuleRequest
- [ ] CommissionRuleResponse
- [ ] CreatePeriodPreviewRequest
- [ ] CommissionPeriodResponse
- [ ] ClosePeriodRequest
- [ ] CreateAdvanceRequest
- [ ] AdvanceResponse
- [ ] RejectAdvanceRequest

---

## 6️⃣ APPLICATION LAYER — USE CASES

### 6.1 UseCases de CommissionRule (Esforço: 2h)

- [ ] Criar `backend/internal/application/usecase/commission/create_rule.go`
- [ ] Criar `backend/internal/application/usecase/commission/get_rule.go`
- [ ] Criar `backend/internal/application/usecase/commission/list_rules.go`
- [ ] Criar `backend/internal/application/usecase/commission/update_rule.go`
- [ ] Criar `backend/internal/application/usecase/commission/delete_rule.go`

#### Checklist

- [ ] CreateCommissionRuleUseCase
- [ ] GetCommissionRuleUseCase
- [ ] ListCommissionRulesUseCase
- [ ] UpdateCommissionRuleUseCase
- [ ] DeleteCommissionRuleUseCase

---

### 6.2 UseCases de CommissionPeriod (Esforço: 3h)

- [ ] Criar `backend/internal/application/usecase/commission/generate_preview.go`
- [ ] Criar `backend/internal/application/usecase/commission/create_period.go`
- [ ] Criar `backend/internal/application/usecase/commission/close_period.go`
- [ ] Criar `backend/internal/application/usecase/commission/list_periods.go`

#### Checklist

- [ ] GeneratePreviewUseCase
- [ ] CreatePeriodUseCase
- [ ] ClosePeriodUseCase (integra com contas_a_pagar)
- [ ] ListPeriodsUseCase

---

### 6.3 UseCases de Advance (Esforço: 2h)

- [ ] Criar `backend/internal/application/usecase/advance/create_advance.go`
- [ ] Criar `backend/internal/application/usecase/advance/approve_advance.go`
- [ ] Criar `backend/internal/application/usecase/advance/reject_advance.go`
- [ ] Criar `backend/internal/application/usecase/advance/list_advances.go`

#### Checklist

- [ ] CreateAdvanceUseCase
- [ ] ApproveAdvanceUseCase
- [ ] RejectAdvanceUseCase
- [ ] ListAdvancesUseCase

---

### 6.4 UseCase: Motor de Cálculo (Esforço: 3h)

- [ ] Criar `backend/internal/application/usecase/commission/calculate_commission.go`

```go
// CalculateCommissionUseCase é chamado quando uma comanda é fechada
type CalculateCommissionUseCase struct {
    ruleRepo       repository.CommissionRuleRepository
    commissionRepo repository.BarberCommissionRepository
    // ... outras deps
}

// Execute calcula e grava as comissões de uma comanda
func (uc *CalculateCommissionUseCase) Execute(ctx context.Context, command *entity.Command) error {
    // 1. Buscar profissional via appointment
    // 2. Para cada command_item de serviço:
    //    - Buscar regra aplicável (hierarquia)
    //    - Calcular valor base (GROSS/TABLE/NET)
    //    - Calcular comissão
    //    - Gravar barber_commission
    return nil
}
```

#### Checklist

- [ ] CalculateCommissionUseCase
- [ ] Hierarquia de regras
- [ ] Base de cálculo configurável
- [ ] Integração com fechamento de comanda

---

## 📝 NOTAS

### Próximos Passos

Após completar esta sprint:
1. Iniciar Sprint 3 (Handlers + Motor de Cálculo)
2. Checklist: `CHECKLIST_SPRINT3_HANDLERS.md`

### Arquivos Criados

| Arquivo | Status |
|---------|--------|
| `entity/commission_rule.go` | ❌ |
| `entity/commission_period.go` | ❌ |
| `entity/advance.go` | ❌ |
| `repository/commission_rule_repository.go` | ❌ |
| `repository/commission_period_repository.go` | ❌ |
| `repository/advance_repository.go` | ❌ |
| `pg_commission_rule_repository.go` | ❌ |
| `pg_commission_period_repository.go` | ❌ |
| `pg_advance_repository.go` | ❌ |
| `dto/commission_rule_dto.go` | ❌ |
| `dto/commission_period_dto.go` | ❌ |
| `dto/advance_dto.go` | ❌ |
| `usecase/commission/*.go` | ❌ |
| `usecase/advance/*.go` | ❌ |

---

*Checklist criado em: 05/12/2025*
