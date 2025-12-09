# ✅ CHECKLIST — SPRINT 6: TESTES, QA & DOCUMENTAÇÃO

> **Status:** ❌ Não Iniciado  
> **Dependência:** Sprint 5 (Frontend Dashboard)  
> **Esforço Estimado:** 12 horas  
> **Prioridade:** P1 — MVP

---

## 📊 RESUMO

```
░░░░░░░░░░░░░░░░░░░░ 0% COMPLETO
```

| Categoria | Completo | Pendente |
|-----------|:--------:|:--------:|
| Testes Backend | 0/6 | 6 |
| Testes E2E | 0/5 | 5 |
| Smoke Tests | 0/1 | 1 |
| Documentação | 0/3 | 3 |
| Revisão QA | 0/4 | 4 |

---

## 1️⃣ TESTES BACKEND (Unit + Integration)

### 1.1 Domain Tests (Esforço: 1.5h)

- [ ] Criar `backend/internal/domain/commission_test.go`

```go
package domain_test

import (
    "testing"
    "nexo/internal/domain"
    "github.com/stretchr/testify/assert"
)

func TestCommissionRule_Validate(t *testing.T) {
    tests := []struct {
        name    string
        rule    domain.CommissionRule
        wantErr bool
    }{
        {
            name: "valid percentage rule",
            rule: domain.CommissionRule{
                TenantID:      "tenant-1",
                Name:          "Padrão",
                Type:          domain.CommissionTypePercentual,
                DefaultRate:   "50.00",
                MinAmount:     nil,
                MaxAmount:     nil,
                EffectiveFrom: time.Now(),
                EffectiveTo:   nil,
            },
            wantErr: false,
        },
        {
            name: "valid fixed rule",
            rule: domain.CommissionRule{
                TenantID:      "tenant-1",
                Name:          "Fixo",
                Type:          domain.CommissionTypeFixo,
                DefaultRate:   "25.00",
                EffectiveFrom: time.Now(),
            },
            wantErr: false,
        },
        {
            name: "invalid - empty name",
            rule: domain.CommissionRule{
                TenantID: "tenant-1",
                Name:     "",
            },
            wantErr: true,
        },
        {
            name: "invalid - rate over 100",
            rule: domain.CommissionRule{
                TenantID:    "tenant-1",
                Name:        "Test",
                Type:        domain.CommissionTypePercentual,
                DefaultRate: "150.00",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.rule.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

#### Checklist

- [ ] TestCommissionRule_Validate
- [ ] TestCommissionItem_Calculate
- [ ] TestCommissionPeriod_CanClose
- [ ] TestAdvance_CanApprove
- [ ] TestAdvance_CanReject
- [ ] TestAdvance_CanDeduct

---

### 1.2 UseCase Tests (Esforço: 2h)

- [ ] Criar `backend/internal/application/usecase/commission_test.go`

```go
package usecase_test

import (
    "context"
    "testing"
    "nexo/internal/application/usecase"
    "nexo/internal/domain"
    "nexo/internal/infrastructure/repository/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestCalculateCommission_Hierarchy(t *testing.T) {
    tests := []struct {
        name               string
        serviceCommission  *string
        barberCommission   *string
        ruleCommission     string
        expectedCommission string
        expectedSource     string
    }{
        {
            name:               "priority 1 - service commission",
            serviceCommission:  ptr("60.00"),
            barberCommission:   ptr("50.00"),
            ruleCommission:     "40.00",
            expectedCommission: "60.00",
            expectedSource:     "servico",
        },
        {
            name:               "priority 2 - barber commission",
            serviceCommission:  nil,
            barberCommission:   ptr("50.00"),
            ruleCommission:     "40.00",
            expectedCommission: "50.00",
            expectedSource:     "profissional",
        },
        {
            name:               "priority 3 - rule commission",
            serviceCommission:  nil,
            barberCommission:   nil,
            ruleCommission:     "40.00",
            expectedCommission: "40.00",
            expectedSource:     "regra",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mocks
            mockRepo := new(mocks.CommissionRepository)
            uc := usecase.NewCalculateCommissionUseCase(mockRepo)

            // Execute
            result, err := uc.Execute(context.Background(), usecase.CalculateCommissionInput{
                ServiceCommission: tt.serviceCommission,
                BarberCommission:  tt.barberCommission,
                RuleCommission:    tt.ruleCommission,
            })

            // Assert
            assert.NoError(t, err)
            assert.Equal(t, tt.expectedCommission, result.Commission)
            assert.Equal(t, tt.expectedSource, result.Source)
        })
    }
}

func ptr(s string) *string {
    return &s
}
```

#### Checklist

- [ ] TestCalculateCommission_Hierarchy
- [ ] TestCalculateCommission_Percentual
- [ ] TestCalculateCommission_Fixo
- [ ] TestClosePeriod_Success
- [ ] TestClosePeriod_AlreadyClosed
- [ ] TestClosePeriod_WithAdvances

---

### 1.3 Repository Tests (Esforço: 1.5h)

- [ ] Criar `backend/internal/infrastructure/repository/commission_repository_test.go`

```go
package repository_test

import (
    "context"
    "database/sql"
    "testing"
    "nexo/internal/infrastructure/repository"
    "github.com/stretchr/testify/assert"
    _ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    // Apply migrations
    return db
}

func TestCommissionRepository_GetByTenant(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := repository.NewCommissionRepository(db)
    ctx := context.Background()
    
    // Test tenant isolation
    result, err := repo.GetByTenant(ctx, "tenant-1")
    assert.NoError(t, err)
    
    // Ensure no cross-tenant data
    for _, item := range result {
        assert.Equal(t, "tenant-1", item.TenantID)
    }
}

func TestCommissionRepository_TenantIsolation(t *testing.T) {
    // Ensure queries never return data from other tenants
}
```

#### Checklist

- [ ] TestCommissionRepository_Create
- [ ] TestCommissionRepository_Update
- [ ] TestCommissionRepository_GetByTenant
- [ ] TestCommissionRepository_TenantIsolation
- [ ] TestAdvanceRepository_Create
- [ ] TestAdvanceRepository_TenantIsolation

---

### 1.4 Handler Tests (Esforço: 1.5h)

- [ ] Criar `backend/internal/interfaces/http/handlers/commission_handler_test.go`

```go
package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "nexo/internal/interfaces/http/handlers"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestCommissionHandler_Create(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    tests := []struct {
        name         string
        body         map[string]interface{}
        expectedCode int
    }{
        {
            name: "valid creation",
            body: map[string]interface{}{
                "name":         "Regra Teste",
                "type":         "PERCENTUAL",
                "default_rate": "50.00",
            },
            expectedCode: http.StatusCreated,
        },
        {
            name: "invalid - missing name",
            body: map[string]interface{}{
                "type":         "PERCENTUAL",
                "default_rate": "50.00",
            },
            expectedCode: http.StatusBadRequest,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            handler := handlers.NewCommissionHandler(/* mocks */)
            router := gin.New()
            router.POST("/commission-rules", handler.Create)
            
            // Execute
            body, _ := json.Marshal(tt.body)
            req := httptest.NewRequest("POST", "/commission-rules", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            
            // Assert
            assert.Equal(t, tt.expectedCode, w.Code)
        })
    }
}

func TestCommissionHandler_RBAC(t *testing.T) {
    // Test that barbeiros can only access their own data
    // Test that admins can access all data
}
```

#### Checklist

- [ ] TestCommissionHandler_Create
- [ ] TestCommissionHandler_Update
- [ ] TestCommissionHandler_Delete
- [ ] TestCommissionHandler_RBAC_BarbeiroOwnData
- [ ] TestAdvanceHandler_Approve_ManagerOnly
- [ ] TestPeriodHandler_Close_ManagerOnly

---

## 2️⃣ TESTES E2E (Playwright)

### 2.1 Fluxo Completo de Comissões (Esforço: 1.5h)

- [ ] Criar `frontend/tests/comissoes/fluxo-completo.spec.ts`

```typescript
import { test, expect } from '@playwright/test';
import { login } from '../helpers/auth';

test.describe('Fluxo Completo de Comissões', () => {
  test.describe('Como Gestor', () => {
    test.beforeEach(async ({ page }) => {
      await login(page, 'admin');
    });

    test('deve configurar regra de comissão', async ({ page }) => {
      await page.goto('/configuracoes/comissoes');
      await page.click('[data-testid="nova-regra"]');
      
      await page.fill('[name=name]', 'Regra Padrão');
      await page.selectOption('[name=type]', 'PERCENTUAL');
      await page.fill('[name=default_rate]', '50');
      
      await page.click('button[type=submit]');
      await expect(page.locator('text=Regra criada')).toBeVisible();
    });

    test('deve visualizar comissões do período', async ({ page }) => {
      await page.goto('/financeiro/comissoes');
      await expect(page.locator('[data-testid="commission-list"]')).toBeVisible();
    });

    test('deve fechar período de comissões', async ({ page }) => {
      await page.goto('/financeiro/comissoes');
      await page.click('[data-testid="fechar-periodo"]');
      await page.click('text=Confirmar');
      
      await expect(page.locator('text=Período fechado')).toBeVisible();
    });

    test('deve gerar contas a pagar', async ({ page }) => {
      // Verificar integração com contas a pagar
      await page.goto('/financeiro/contas-pagar');
      await expect(page.locator('text=Comissão')).toBeVisible();
    });
  });

  test.describe('Como Barbeiro', () => {
    test.beforeEach(async ({ page }) => {
      await login(page, 'barbeiro');
    });

    test('deve ver apenas suas próprias comissões', async ({ page }) => {
      await page.goto('/barbeiro/painel');
      // Verificar que não vê dados de outros barbeiros
    });

    test('não deve acessar configurações', async ({ page }) => {
      const response = await page.goto('/configuracoes/comissoes');
      expect(response?.status()).toBe(403);
    });
  });
});
```

#### Checklist

- [ ] Configurar regra de comissão
- [ ] Listar comissões do período
- [ ] Fechar período
- [ ] Verificar contas a pagar geradas
- [ ] RBAC: barbeiro só vê seus dados
- [ ] RBAC: barbeiro não configura regras

---

### 2.2 Fluxo de Adiantamentos (Esforço: 1h)

- [ ] Criar `frontend/tests/comissoes/adiantamentos-e2e.spec.ts`

```typescript
import { test, expect } from '@playwright/test';
import { login } from '../helpers/auth';

test.describe('Fluxo de Adiantamentos', () => {
  test('barbeiro solicita, gestor aprova, deduz no fechamento', async ({ page }) => {
    // 1. Barbeiro solicita
    await login(page, 'barbeiro');
    await page.goto('/barbeiro/painel');
    await page.click('text=Solicitar Adiantamento');
    await page.fill('[name=amount]', '200');
    await page.fill('[name=reason]', 'Necessidade urgente');
    await page.click('button:has-text("Solicitar")');
    
    await expect(page.locator('text=solicitado com sucesso')).toBeVisible();
    
    // 2. Gestor aprova
    await login(page, 'admin');
    await page.goto('/financeiro/adiantamentos');
    await page.click('[data-testid="aprovar"]');
    
    await expect(page.locator('text=Aprovado')).toBeVisible();
    
    // 3. No fechamento, valor é deduzido
    await page.goto('/financeiro/comissoes');
    await page.click('[data-testid="fechar-periodo"]');
    await page.click('text=Confirmar');
    
    // Verificar que adiantamento foi deduzido
  });

  test('gestor rejeita adiantamento com motivo', async ({ page }) => {
    await login(page, 'admin');
    await page.goto('/financeiro/adiantamentos');
    
    await page.click('[data-testid="rejeitar"]');
    await page.fill('[name=reason]', 'Limite mensal atingido');
    await page.click('text=Confirmar Rejeição');
    
    await expect(page.locator('text=Rejeitado')).toBeVisible();
  });
});
```

#### Checklist

- [ ] Barbeiro solicita adiantamento
- [ ] Gestor aprova
- [ ] Gestor rejeita com motivo
- [ ] Dedução no fechamento
- [ ] Verificar status atualizado

---

### 2.3 Multi-Tenant Isolation (Esforço: 1h)

- [ ] Criar `frontend/tests/comissoes/multi-tenant.spec.ts`

```typescript
import { test, expect } from '@playwright/test';
import { login } from '../helpers/auth';

test.describe('Isolamento Multi-Tenant', () => {
  test('tenant A não vê dados do tenant B', async ({ page }) => {
    // Login no tenant A
    await login(page, 'admin-tenant-a');
    await page.goto('/financeiro/comissoes');
    
    const commissions = await page.locator('[data-testid="commission-row"]').all();
    for (const row of commissions) {
      const tenantId = await row.getAttribute('data-tenant');
      expect(tenantId).toBe('tenant-a');
    }
    
    // Login no tenant B
    await login(page, 'admin-tenant-b');
    await page.goto('/financeiro/comissoes');
    
    const commissionsB = await page.locator('[data-testid="commission-row"]').all();
    for (const row of commissionsB) {
      const tenantId = await row.getAttribute('data-tenant');
      expect(tenantId).toBe('tenant-b');
    }
  });
});
```

#### Checklist

- [ ] Comissões isoladas por tenant
- [ ] Regras isoladas por tenant
- [ ] Adiantamentos isolados por tenant

---

### 2.4 Relatórios e DRE (Esforço: 1h)

- [ ] Criar `frontend/tests/comissoes/relatorios.spec.ts`

```typescript
import { test, expect } from '@playwright/test';

test.describe('Integração com DRE', () => {
  test('comissões aparecem no DRE mensal', async ({ page }) => {
    await page.goto('/financeiro/dre');
    
    // Selecionar mês/ano
    await page.selectOption('[name=month]', '01');
    await page.selectOption('[name=year]', '2025');
    
    // Verificar que custo_comissoes está presente
    await expect(page.locator('text=Comissões')).toBeVisible();
    
    // Verificar valor
    const valor = await page.locator('[data-testid="custo-comissoes"]').textContent();
    expect(parseFloat(valor?.replace(/[^0-9,]/g, '').replace(',', '.'))).toBeGreaterThan(0);
  });

  test('relatório de comissões por profissional', async ({ page }) => {
    await page.goto('/financeiro/relatorios/comissoes');
    
    // Filtrar por profissional
    await page.selectOption('[name=professional]', 'barbeiro-1');
    await page.click('text=Gerar Relatório');
    
    // Verificar dados
    await expect(page.locator('[data-testid="report-table"]')).toBeVisible();
  });
});
```

#### Checklist

- [ ] Comissões no DRE
- [ ] Relatório por profissional
- [ ] Relatório por período
- [ ] Export para CSV/PDF

---

## 3️⃣ SMOKE TESTS

### 3.1 Script de Smoke Test (Esforço: 1h)

- [ ] Criar `scripts/smoke_tests_comissoes.sh`

```bash
#!/bin/bash
set -e

BASE_URL="${BASE_URL:-http://localhost:3000}"
API_URL="${API_URL:-http://localhost:8080/api/v1}"

echo "🧪 SMOKE TESTS - MÓDULO COMISSÕES"
echo "=================================="

# Auth
TOKEN=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"123456"}' \
  | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Auth failed"
  exit 1
fi
echo "✅ Auth OK"

# 1. Commission Rules
echo -e "\n📋 Commission Rules"
RULE_RESPONSE=$(curl -s -X POST "$API_URL/commission-rules" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Rule","type":"PERCENTUAL","default_rate":"50.00"}')
RULE_ID=$(echo $RULE_RESPONSE | jq -r '.id')

if [ "$RULE_ID" == "null" ]; then
  echo "❌ Create rule failed"
  exit 1
fi
echo "✅ Create rule: $RULE_ID"

# List rules
LIST_RULES=$(curl -s "$API_URL/commission-rules" \
  -H "Authorization: Bearer $TOKEN")
COUNT=$(echo $LIST_RULES | jq 'length')
echo "✅ List rules: $COUNT rules found"

# 2. Commission Items
echo -e "\n💰 Commission Items"
ITEMS=$(curl -s "$API_URL/commission-items?status=PENDENTE" \
  -H "Authorization: Bearer $TOKEN")
ITEMS_COUNT=$(echo $ITEMS | jq 'length')
echo "✅ Pending items: $ITEMS_COUNT"

# 3. Advances
echo -e "\n💵 Advances"
ADVANCE=$(curl -s -X POST "$API_URL/advances" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"professional_id":"prof-1","amount":"100.00","reason":"Test"}')
ADVANCE_ID=$(echo $ADVANCE | jq -r '.id')

if [ "$ADVANCE_ID" != "null" ]; then
  echo "✅ Create advance: $ADVANCE_ID"
  
  # Approve
  curl -s -X POST "$API_URL/advances/$ADVANCE_ID/approve" \
    -H "Authorization: Bearer $TOKEN" > /dev/null
  echo "✅ Approve advance"
fi

# 4. Close Period
echo -e "\n📅 Period Management"
PERIOD=$(curl -s -X POST "$API_URL/commission-periods/close" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"month":1,"year":2025}')
PERIOD_ID=$(echo $PERIOD | jq -r '.id')
echo "✅ Close period: $PERIOD_ID"

# 5. DRE Integration
echo -e "\n📊 DRE Integration"
DRE=$(curl -s "$API_URL/dre?month=1&year=2025" \
  -H "Authorization: Bearer $TOKEN")
COMISSOES=$(echo $DRE | jq -r '.custo_comissoes')
echo "✅ DRE custo_comissoes: R$ $COMISSOES"

# Cleanup
curl -s -X DELETE "$API_URL/commission-rules/$RULE_ID" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

echo -e "\n=================================="
echo "✅ ALL SMOKE TESTS PASSED!"
```

#### Checklist

- [ ] Auth
- [ ] CRUD Commission Rules
- [ ] List Commission Items
- [ ] CRUD Advances
- [ ] Approve/Reject Advances
- [ ] Close Period
- [ ] DRE Integration

---

## 4️⃣ DOCUMENTAÇÃO

### 4.1 API Documentation (Esforço: 1h)

- [ ] Atualizar `backend/docs/swagger.yaml`

```yaml
# Adicionar ao swagger existente

paths:
  /commission-rules:
    get:
      tags: [Commissions]
      summary: Lista regras de comissão
      security:
        - BearerAuth: []
      responses:
        200:
          description: Lista de regras
    post:
      tags: [Commissions]
      summary: Cria nova regra de comissão
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateCommissionRuleRequest'

  /advances:
    get:
      tags: [Advances]
      summary: Lista adiantamentos
    post:
      tags: [Advances]
      summary: Solicita novo adiantamento

  /advances/{id}/approve:
    post:
      tags: [Advances]
      summary: Aprova adiantamento
      
  /advances/{id}/reject:
    post:
      tags: [Advances]
      summary: Rejeita adiantamento

  /commission-periods/close:
    post:
      tags: [Commissions]
      summary: Fecha período de comissões

components:
  schemas:
    CreateCommissionRuleRequest:
      type: object
      required: [name, type, default_rate]
      properties:
        name:
          type: string
        type:
          type: string
          enum: [PERCENTUAL, FIXO]
        default_rate:
          type: string
```

#### Checklist

- [ ] Documentar /commission-rules
- [ ] Documentar /commission-items
- [ ] Documentar /advances
- [ ] Documentar /commission-periods
- [ ] Schemas de request/response

---

### 4.2 User Guide (Esforço: 0.5h)

- [ ] Criar `docs/07-produto-e-funcionalidades/GUIA_COMISSOES.md`

#### Checklist

- [ ] Como configurar regras
- [ ] Como visualizar comissões
- [ ] Como solicitar adiantamentos
- [ ] Como fechar período
- [ ] FAQ

---

### 4.3 Technical Documentation (Esforço: 0.5h)

- [ ] Atualizar `docs/02-arquitetura/MODELO_DE_DADOS.md`

#### Checklist

- [ ] Diagrama das novas tabelas
- [ ] Relacionamentos
- [ ] Índices criados

---

## 5️⃣ REVISÃO QA

### 5.1 Checklist de Qualidade

- [ ] **Performance**: Queries otimizadas com índices
- [ ] **Segurança**: Tenant isolation em todas as queries
- [ ] **RBAC**: Permissões verificadas em todos os handlers
- [ ] **UX**: Loading states, empty states, error handling
- [ ] **Acessibilidade**: Labels, focus, keyboard navigation
- [ ] **Mobile**: Responsividade verificada
- [ ] **Logs**: Operações críticas logadas

### 5.2 Code Review Checklist

- [ ] Sem SQL manual (sqlc only)
- [ ] Sem `tenant_id` em payloads
- [ ] Dinheiro como string, não float
- [ ] Handlers finos, lógica em use cases
- [ ] DTOs com validação
- [ ] Testes cobrindo casos edge

### 5.3 Security Review

- [ ] Tenant isolation testada
- [ ] RBAC testado
- [ ] Input sanitization
- [ ] Rate limiting em endpoints críticos

### 5.4 Integration Verification

- [ ] DRE recebe custo_comissoes
- [ ] Contas a Pagar criadas corretamente
- [ ] Fluxo de Caixa impactado
- [ ] Notifications enviadas (se aplicável)

---

## 📝 NOTAS

### Cobertura de Testes Esperada

| Camada | Mínimo | Ideal |
|--------|--------|-------|
| Domain | 90% | 100% |
| UseCases | 80% | 95% |
| Repository | 70% | 85% |
| Handlers | 70% | 80% |
| E2E | Core flows | All flows |

### Comandos de Teste

```bash
# Backend unit tests
cd backend && go test ./... -v -cover

# Backend with coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# E2E tests
cd frontend && pnpm test:e2e

# Smoke tests
./scripts/smoke_tests_comissoes.sh
```

### Arquivos Criados

| Arquivo | Status |
|---------|--------|
| `domain/commission_test.go` | ❌ |
| `usecase/commission_test.go` | ❌ |
| `repository/commission_repository_test.go` | ❌ |
| `handlers/commission_handler_test.go` | ❌ |
| `tests/comissoes/fluxo-completo.spec.ts` | ❌ |
| `tests/comissoes/adiantamentos-e2e.spec.ts` | ❌ |
| `tests/comissoes/multi-tenant.spec.ts` | ❌ |
| `tests/comissoes/relatorios.spec.ts` | ❌ |
| `scripts/smoke_tests_comissoes.sh` | ❌ |
| `docs/.../GUIA_COMISSOES.md` | ❌ |

---

## ✅ CRITÉRIOS DE ACEITE FINAL

Antes de marcar o módulo como completo:

1. [ ] Todos os testes de unidade passando
2. [ ] Todos os testes E2E passando
3. [ ] Smoke tests passando em staging
4. [ ] Code review aprovado
5. [ ] Security review aprovado
6. [ ] Documentação atualizada
7. [ ] Performance validada (queries < 100ms)
8. [ ] Zero erros de lint/type

---

*Checklist criado em: 05/12/2025*
