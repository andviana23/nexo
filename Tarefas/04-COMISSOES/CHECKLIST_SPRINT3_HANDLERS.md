# ✅ CHECKLIST — SPRINT 3: HANDLERS + MOTOR DE CÁLCULO

> **Status:** ❌ Não Iniciado  
> **Dependência:** Sprint 2 (Domain + Repository + UseCases)  
> **Esforço Estimado:** 14 horas  
> **Prioridade:** P0 — Bloqueia Frontend

---

## 📊 RESUMO

```
░░░░░░░░░░░░░░░░░░░░ 0% COMPLETO
```

| Categoria | Completo | Pendente |
|-----------|:--------:|:--------:|
| Handlers | 0/4 | 4 |
| Rotas | 0/20 | 20 |
| Motor de Cálculo | 0/1 | 1 |
| Integração Financeiro | 0/2 | 2 |
| Testes Unitários | 0/5 | 5 |

---

## 1️⃣ HANDLERS

### 1.1 Handler: `CommissionRulesHandler` (Esforço: 2h)

- [ ] Criar `backend/internal/interfaces/http/handler/commission_rules_handler.go`

#### Endpoints

| Método | Rota | Handler | Descrição |
|--------|------|---------|-----------|
| `GET` | `/api/v1/commission-rules` | ListRules | Listar regras |
| `GET` | `/api/v1/commission-rules/:id` | GetRule | Buscar por ID |
| `POST` | `/api/v1/commission-rules` | CreateRule | Criar regra |
| `PUT` | `/api/v1/commission-rules/:id` | UpdateRule | Atualizar |
| `PATCH` | `/api/v1/commission-rules/:id/toggle` | ToggleRule | Ativar/Desativar |
| `DELETE` | `/api/v1/commission-rules/:id` | DeleteRule | Remover |

#### Código Base

```go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "nexo/internal/application/dto"
    "nexo/internal/application/usecase/commission"
)

type CommissionRulesHandler struct {
    createUC *commission.CreateCommissionRuleUseCase
    getUC    *commission.GetCommissionRuleUseCase
    listUC   *commission.ListCommissionRulesUseCase
    updateUC *commission.UpdateCommissionRuleUseCase
    deleteUC *commission.DeleteCommissionRuleUseCase
}

func NewCommissionRulesHandler(/* deps */) *CommissionRulesHandler {
    return &CommissionRulesHandler{/* ... */}
}

// CreateRule cria uma nova regra de comissão
// @Summary Criar regra de comissão
// @Tags Commission Rules
// @Accept json
// @Produce json
// @Param body body dto.CreateCommissionRuleRequest true "Dados da regra"
// @Success 201 {object} dto.CommissionRuleResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/v1/commission-rules [post]
func (h *CommissionRulesHandler) CreateRule(c *gin.Context) {
    tenantID, _ := c.Get("tenant_id")
    userID, _ := c.Get("user_id")
    
    var req dto.CreateCommissionRuleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    result, err := h.createUC.Execute(c.Request.Context(), tenantID.(uuid.UUID), userID.(uuid.UUID), &req)
    if err != nil {
        handleError(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, result)
}

// ListRules lista regras de comissão
func (h *CommissionRulesHandler) ListRules(c *gin.Context) {
    // Implementar com filtros: unit_id, professional_id, service_id, active
}

// GetRule busca regra por ID
func (h *CommissionRulesHandler) GetRule(c *gin.Context) {
    // Implementar
}

// UpdateRule atualiza regra
func (h *CommissionRulesHandler) UpdateRule(c *gin.Context) {
    // Implementar
}

// ToggleRule ativa/desativa regra
func (h *CommissionRulesHandler) ToggleRule(c *gin.Context) {
    // Implementar
}

// DeleteRule remove regra
func (h *CommissionRulesHandler) DeleteRule(c *gin.Context) {
    // Implementar
}
```

#### Checklist

- [ ] NewCommissionRulesHandler
- [ ] CreateRule
- [ ] ListRules (com filtros)
- [ ] GetRule
- [ ] UpdateRule
- [ ] ToggleRule
- [ ] DeleteRule
- [ ] Swagger annotations

---

### 1.2 Handler: `CommissionsHandler` (Esforço: 1.5h)

- [ ] Criar `backend/internal/interfaces/http/handler/commissions_handler.go`

#### Endpoints

| Método | Rota | Handler | Descrição |
|--------|------|---------|-----------|
| `GET` | `/api/v1/commissions` | ListCommissions | Listar comissões |
| `GET` | `/api/v1/commissions/summary` | GetSummary | Resumo do período |
| `GET` | `/api/v1/professionals/:id/commissions` | ListByProfessional | Comissões do barbeiro |

#### Checklist

- [ ] NewCommissionsHandler
- [ ] ListCommissions (filtros: barbeiro, período, status)
- [ ] GetSummary
- [ ] ListByProfessional
- [ ] RBAC: barbeiro só vê suas comissões

---

### 1.3 Handler: `CommissionPeriodsHandler` (Esforço: 2.5h)

- [ ] Criar `backend/internal/interfaces/http/handler/commission_periods_handler.go`

#### Endpoints

| Método | Rota | Handler | Descrição |
|--------|------|---------|-----------|
| `GET` | `/api/v1/commission-periods` | ListPeriods | Listar períodos |
| `GET` | `/api/v1/commission-periods/:id` | GetPeriod | Buscar por ID |
| `POST` | `/api/v1/commission-periods/preview` | GeneratePreview | Gerar prévia |
| `POST` | `/api/v1/commission-periods` | CreatePeriod | Criar período (DRAFT) |
| `PUT` | `/api/v1/commission-periods/:id` | UpdatePeriod | Atualizar (ajustes) |
| `POST` | `/api/v1/commission-periods/:id/close` | ClosePeriod | Fechar período |
| `DELETE` | `/api/v1/commission-periods/:id` | DeletePeriod | Remover (apenas DRAFT) |

#### Código Base

```go
// ClosePeriod fecha o período e gera conta a pagar
func (h *CommissionPeriodsHandler) ClosePeriod(c *gin.Context) {
    tenantID, _ := c.Get("tenant_id")
    userID, _ := c.Get("user_id")
    periodID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }
    
    var req dto.ClosePeriodRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // UseCase: fecha período + gera conta a pagar
    result, err := h.closeUC.Execute(c.Request.Context(), tenantID.(uuid.UUID), periodID, userID.(uuid.UUID), &req)
    if err != nil {
        handleError(c, err)
        return
    }
    
    c.JSON(http.StatusOK, result)
}
```

#### Checklist

- [ ] NewCommissionPeriodsHandler
- [ ] ListPeriods (filtros: profissional, unidade, status, datas)
- [ ] GetPeriod
- [ ] GeneratePreview
- [ ] CreatePeriod
- [ ] UpdatePeriod
- [ ] ClosePeriod (+ gera contas_a_pagar)
- [ ] DeletePeriod

---

### 1.4 Handler: `AdvancesHandler` (Esforço: 2h)

- [ ] Criar `backend/internal/interfaces/http/handler/advances_handler.go`

#### Endpoints

| Método | Rota | Handler | Descrição |
|--------|------|---------|-----------|
| `GET` | `/api/v1/advances` | ListAdvances | Listar adiantamentos |
| `GET` | `/api/v1/advances/:id` | GetAdvance | Buscar por ID |
| `POST` | `/api/v1/advances` | CreateAdvance | Criar solicitação |
| `POST` | `/api/v1/advances/:id/approve` | ApproveAdvance | Aprovar |
| `POST` | `/api/v1/advances/:id/reject` | RejectAdvance | Rejeitar |
| `DELETE` | `/api/v1/advances/:id` | DeleteAdvance | Remover (apenas PENDING) |

#### Checklist

- [ ] NewAdvancesHandler
- [ ] ListAdvances (filtros: profissional, status)
- [ ] GetAdvance
- [ ] CreateAdvance
- [ ] ApproveAdvance
- [ ] RejectAdvance
- [ ] DeleteAdvance
- [ ] RBAC: barbeiro pode criar, gestor aprova/rejeita

---

## 2️⃣ ROTAS

### 2.1 Registrar Rotas (Esforço: 1h)

- [ ] Atualizar `backend/cmd/api/main.go` ou `routes.go`

```go
// Commission Routes
commissionRulesHandler := handler.NewCommissionRulesHandler(/* deps */)
commissionsHandler := handler.NewCommissionsHandler(/* deps */)
commissionPeriodsHandler := handler.NewCommissionPeriodsHandler(/* deps */)
advancesHandler := handler.NewAdvancesHandler(/* deps */)

// Grupo de rotas protegidas
api := r.Group("/api/v1")
api.Use(authMiddleware.Authenticate())

// Commission Rules (Admin/Manager)
rules := api.Group("/commission-rules")
rules.Use(rbacMiddleware.RequireRole("admin", "manager"))
{
    rules.GET("", commissionRulesHandler.ListRules)
    rules.GET("/:id", commissionRulesHandler.GetRule)
    rules.POST("", commissionRulesHandler.CreateRule)
    rules.PUT("/:id", commissionRulesHandler.UpdateRule)
    rules.PATCH("/:id/toggle", commissionRulesHandler.ToggleRule)
    rules.DELETE("/:id", commissionRulesHandler.DeleteRule)
}

// Commissions (Admin/Manager/Barber)
commissions := api.Group("/commissions")
{
    commissions.GET("", commissionsHandler.ListCommissions)
    commissions.GET("/summary", commissionsHandler.GetSummary)
}

// Professional Commissions (Barber pode ver só as próprias)
api.GET("/professionals/:id/commissions", commissionsHandler.ListByProfessional)

// Commission Periods (Admin/Manager)
periods := api.Group("/commission-periods")
periods.Use(rbacMiddleware.RequireRole("admin", "manager"))
{
    periods.GET("", commissionPeriodsHandler.ListPeriods)
    periods.GET("/:id", commissionPeriodsHandler.GetPeriod)
    periods.POST("/preview", commissionPeriodsHandler.GeneratePreview)
    periods.POST("", commissionPeriodsHandler.CreatePeriod)
    periods.PUT("/:id", commissionPeriodsHandler.UpdatePeriod)
    periods.POST("/:id/close", commissionPeriodsHandler.ClosePeriod)
    periods.DELETE("/:id", commissionPeriodsHandler.DeletePeriod)
}

// Advances
advances := api.Group("/advances")
{
    advances.GET("", advancesHandler.ListAdvances)
    advances.GET("/:id", advancesHandler.GetAdvance)
    advances.POST("", advancesHandler.CreateAdvance) // Barber pode criar
    
    // Apenas Admin/Manager
    advances.POST("/:id/approve", rbacMiddleware.RequireRole("admin", "manager"), advancesHandler.ApproveAdvance)
    advances.POST("/:id/reject", rbacMiddleware.RequireRole("admin", "manager"), advancesHandler.RejectAdvance)
    advances.DELETE("/:id", advancesHandler.DeleteAdvance)
}
```

#### Checklist

- [ ] Rotas commission-rules
- [ ] Rotas commissions
- [ ] Rotas commission-periods
- [ ] Rotas advances
- [ ] RBAC configurado
- [ ] Middleware de autenticação

---

## 3️⃣ MOTOR DE CÁLCULO

### 3.1 Integração com Fechamento de Comanda (Esforço: 4h)

- [ ] Criar/Atualizar `backend/internal/application/usecase/command/close_command.go`

O motor de cálculo é acionado quando uma comanda é fechada.

```go
// No UseCase de fechar comanda, adicionar:
func (uc *CloseCommandUseCase) Execute(ctx context.Context, commandID uuid.UUID) error {
    // 1. Fechar comanda (lógica existente)
    command, err := uc.commandRepo.GetByID(ctx, commandID)
    if err != nil {
        return err
    }
    
    // 2. Marcar como CLOSED
    command.Status = "CLOSED"
    if err := uc.commandRepo.Update(ctx, command); err != nil {
        return err
    }
    
    // 3. NOVO: Calcular comissões
    if err := uc.calculateCommissions(ctx, command); err != nil {
        // Log error mas não falha o fechamento
        log.Printf("Erro ao calcular comissões: %v", err)
    }
    
    return nil
}

func (uc *CloseCommandUseCase) calculateCommissions(ctx context.Context, command *entity.Command) error {
    // Buscar appointment para pegar profissional
    appointment, err := uc.appointmentRepo.GetByID(ctx, command.TenantID, command.AppointmentID)
    if err != nil || appointment == nil {
        return err
    }
    
    professionalID := appointment.ProfessionalID
    
    // Buscar command_items do tipo SERVICO
    items, err := uc.commandItemRepo.ListByCommand(ctx, command.ID)
    if err != nil {
        return err
    }
    
    for _, item := range items {
        if item.Tipo != "SERVICO" {
            continue
        }
        
        // Buscar regra aplicável
        rule, err := uc.ruleRepo.FindApplicable(
            ctx,
            command.TenantID,
            command.UnitID,
            &professionalID,
            &item.ItemID,
        )
        if err != nil {
            continue
        }
        
        // Se não encontrou regra, buscar % do serviço ou profissional
        var commissionValue float64
        if rule != nil {
            commissionValue = rule.Calculate(item.PrecoFinal)
        } else {
            // Fallback: buscar do serviço ou profissional
            commissionValue = uc.calculateFallback(ctx, command.TenantID, professionalID, item)
        }
        
        // Criar registro de comissão
        commission := &entity.BarberCommission{
            TenantID:       command.TenantID,
            BarbeiroID:     professionalID,
            CommandItemID:  &item.ID,
            UnitID:         command.UnitID,
            Valor:          commissionValue,
            Status:         "PENDENTE",
            DataCompetencia: command.CreatedAt,
            Manual:         false,
        }
        
        if err := uc.commissionRepo.CreateFromCommand(ctx, commission); err != nil {
            log.Printf("Erro ao criar comissão para item %s: %v", item.ID, err)
        }
    }
    
    return nil
}

func (uc *CloseCommandUseCase) calculateFallback(ctx context.Context, tenantID, professionalID uuid.UUID, item *entity.CommandItem) float64 {
    // 1. Tentar do serviço
    service, err := uc.serviceRepo.GetByID(ctx, tenantID, item.ItemID)
    if err == nil && service.Comissao > 0 {
        return item.PrecoFinal * (service.Comissao / 100)
    }
    
    // 2. Tentar do profissional
    professional, err := uc.professionalRepo.GetByID(ctx, tenantID, professionalID)
    if err == nil && professional.Comissao > 0 {
        if professional.TipoComissao == "PERCENTUAL" {
            return item.PrecoFinal * (professional.Comissao / 100)
        }
        return professional.Comissao // FIXO
    }
    
    return 0
}
```

#### Checklist

- [ ] Hook no fechamento de comanda
- [ ] Busca de regra aplicável (hierarquia)
- [ ] Fallback para serviço/profissional
- [ ] Criação de barber_commission
- [ ] Tratamento de erros (não bloqueia fechamento)
- [ ] Log de auditoria

---

## 4️⃣ INTEGRAÇÃO FINANCEIRA

### 4.1 Geração de Conta a Pagar no Fechamento (Esforço: 2h)

- [ ] Atualizar `backend/internal/application/usecase/commission/close_period.go`

```go
func (uc *ClosePeriodUseCase) Execute(ctx context.Context, tenantID, periodID, userID uuid.UUID, req *dto.ClosePeriodRequest) (*dto.CommissionPeriodResponse, error) {
    // 1. Buscar período
    period, err := uc.periodRepo.GetByID(ctx, tenantID, periodID)
    if err != nil {
        return nil, err
    }
    
    if !period.CanClose() {
        return nil, errors.New("período não pode ser fechado")
    }
    
    // 2. Buscar profissional para nome
    professional, err := uc.professionalRepo.GetByID(ctx, tenantID, period.ProfessionalID)
    if err != nil {
        return nil, err
    }
    
    // 3. Buscar categoria "Comissões"
    categoryID, err := uc.categoryRepo.FindByName(ctx, tenantID, "Comissões")
    if err != nil {
        // Criar categoria se não existir
        categoryID, err = uc.categoryRepo.Create(ctx, &entity.Category{
            TenantID: tenantID,
            Nome:     "Comissões",
            Tipo:     "DESPESA",
            TipoCusto: "VARIAVEL",
        })
    }
    
    // 4. Criar conta a pagar
    bill := &entity.ContaPagar{
        TenantID:       tenantID,
        UnitID:         period.UnitID,
        Descricao:      fmt.Sprintf("Comissão %s - %s a %s", professional.Nome, period.StartDate.Format("02/01"), period.EndDate.Format("02/01/2006")),
        CategoriaID:    categoryID,
        Fornecedor:     professional.Nome,
        Valor:          period.NetValue,
        Tipo:           "VARIAVEL",
        DataVencimento: uc.calculateDueDate(period.EndDate),
        Status:         "ABERTO",
    }
    
    if err := uc.billRepo.Create(ctx, bill); err != nil {
        return nil, fmt.Errorf("erro ao criar conta a pagar: %w", err)
    }
    
    // 5. Fechar período
    if err := period.Close(userID, bill.ID); err != nil {
        return nil, err
    }
    
    if err := uc.periodRepo.Close(ctx, period); err != nil {
        return nil, err
    }
    
    // 6. Marcar comissões como PROCESSADO
    if err := uc.commissionRepo.MarkAsProcessed(
        ctx,
        tenantID,
        period.ProfessionalID,
        period.ID,
        period.StartDate,
        period.EndDate,
    ); err != nil {
        log.Printf("Erro ao marcar comissões como processadas: %v", err)
    }
    
    // 7. Deduzir adiantamentos
    if err := uc.deductAdvances(ctx, tenantID, period); err != nil {
        log.Printf("Erro ao deduzir adiantamentos: %v", err)
    }
    
    return dto.ToCommissionPeriodResponse(period), nil
}

func (uc *ClosePeriodUseCase) deductAdvances(ctx context.Context, tenantID uuid.UUID, period *entity.CommissionPeriod) error {
    advances, err := uc.advanceRepo.ListApprovedNotDeducted(ctx, tenantID, period.ProfessionalID)
    if err != nil {
        return err
    }
    
    for _, advance := range advances {
        if err := advance.Deduct(period.ID); err != nil {
            continue
        }
        if err := uc.advanceRepo.Deduct(ctx, advance); err != nil {
            log.Printf("Erro ao deduzir adiantamento %s: %v", advance.ID, err)
        }
    }
    
    return nil
}
```

#### Checklist

- [ ] Criar conta a pagar automática
- [ ] Categoria "Comissões" (criar se não existir)
- [ ] Fornecedor = Nome do profissional
- [ ] Marcar comissões como PROCESSADO
- [ ] Deduzir adiantamentos
- [ ] Transaction para consistência

---

### 4.2 Atualização da DRE (Esforço: 1h)

- [ ] Criar/Atualizar hook para atualizar `dre_mensal.custo_comissoes`

```go
// Pode ser um job que roda ao final do dia ou ao fechar período
func (uc *UpdateDREUseCase) UpdateCommissionCosts(ctx context.Context, tenantID uuid.UUID, month string) error {
    // Buscar total de comissões PROCESSADO do mês
    startDate, endDate := parseMonthRange(month)
    
    total, err := uc.commissionRepo.SumByPeriod(ctx, tenantID, startDate, endDate, "PROCESSADO")
    if err != nil {
        return err
    }
    
    // Atualizar DRE
    dre, err := uc.dreRepo.GetByMonth(ctx, tenantID, month)
    if err != nil {
        // Criar se não existir
        dre = &entity.DREMensal{
            TenantID: tenantID,
            MesAno:   month,
        }
    }
    
    dre.CustoComissoes = total
    dre.RecalculateResults()
    
    return uc.dreRepo.Upsert(ctx, dre)
}
```

#### Checklist

- [ ] Somar comissões PROCESSADO do mês
- [ ] Atualizar campo custo_comissoes
- [ ] Recalcular resultado operacional

---

## 5️⃣ TESTES UNITÁRIOS

### 5.1 Testes do Motor de Cálculo (Esforço: 2h)

- [ ] Criar `backend/internal/application/usecase/commission/calculate_commission_test.go`

```go
func TestCalculateCommission_PercentageRule(t *testing.T) {
    // Regra: 50%
    // Base: R$ 100
    // Esperado: R$ 50
}

func TestCalculateCommission_FixedRule(t *testing.T) {
    // Regra: R$ 15 fixo
    // Esperado: R$ 15
}

func TestCalculateCommission_HybridRule(t *testing.T) {
    // Regra: R$ 500 + 30%
    // Base: R$ 1000
    // Esperado: R$ 500 + R$ 300 = R$ 800
}

func TestCalculateCommission_ProgressiveRule(t *testing.T) {
    // Regra: < 5k = 40%, >= 5k = 50%
    // Base: R$ 6000
    // Esperado: R$ 3000
}

func TestCalculateCommission_Hierarchy(t *testing.T) {
    // Testar hierarquia: Serviço > Profissional > Unidade > Tenant
}

func TestCalculateCommission_Fallback(t *testing.T) {
    // Sem regra, usa profissionais.comissao
}
```

#### Checklist Testes

- [ ] TestCalculateCommission_PercentageRule
- [ ] TestCalculateCommission_FixedRule
- [ ] TestCalculateCommission_HybridRule
- [ ] TestCalculateCommission_ProgressiveRule
- [ ] TestCalculateCommission_Hierarchy
- [ ] TestCalculateCommission_Fallback
- [ ] TestClosePeriod_GeneratesBill
- [ ] TestClosePeriod_DeductsAdvances

---

## 📝 NOTAS

### Próximos Passos

Após completar esta sprint:
1. Iniciar Sprint 4 (Frontend Config + Fechamento)
2. Checklist: `CHECKLIST_SPRINT4_FRONTEND_CONFIG.md`

### Arquivos Criados

| Arquivo | Status |
|---------|--------|
| `handler/commission_rules_handler.go` | ❌ |
| `handler/commissions_handler.go` | ❌ |
| `handler/commission_periods_handler.go` | ❌ |
| `handler/advances_handler.go` | ❌ |
| `usecase/command/close_command.go` (update) | ❌ |
| `usecase/commission/close_period.go` (update) | ❌ |
| `*_test.go` | ❌ |

---

*Checklist criado em: 05/12/2025*
