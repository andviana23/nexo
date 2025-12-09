# 🔗 Sprint 3: Integração Asaas — Módulo Assinaturas

**Sprint:** 3 de 5  
**Status:** ✅ CONCLUÍDO  
**Progresso:** 100%  
**Estimativa:** 15-20 horas  
**Prioridade:** 🔴 CRÍTICA  
**Dependência:** ✅ Sprint 2 (Backend Core) deve estar concluída

---

## 📚 Referência Obrigatória

> ⚠️ **ANTES DE INICIAR**, leia completamente:
> 
> - **[FLUXO_ASSINATURA.md](../../docs/11-Fluxos/FLUXO_ASSINATURA.md)** — Fonte da verdade
>   - Seção 4: Integração Asaas (endpoints, webhooks, regras AS-001 a AS-008)
>   - Seção 6.1: Fluxo Nova Assinatura Cartão
>   - Seção 6.6: Fluxo Processar Webhook
>   - Seção 9: Tratamento de Erros
> - **[INTEGRACAO_ASAAS_PLANO.md](../INTEGRACAO_ASAAS_PLANO.md)** — Plano técnico existente

---

## 📊 Progresso das Tarefas

| ID | Tarefa | Estimativa | Status | Progresso |
|----|--------|------------|--------|-----------|
| **Gateway HTTP** |
| AS-001 | Cliente HTTP Base | 1h | ✅ Concluído | 100% |
| AS-002 | Types/DTOs Asaas | 1h | ✅ Concluído | 100% |
| AS-003 | Método: FindOrCreateCustomer | 2h | ✅ Concluído | 100% |
| AS-004 | Método: CreateSubscription | 2h | ✅ Concluído | 100% |
| AS-005 | Método: GeneratePaymentLink | 1h | ✅ Concluído | 100% |
| AS-006 | Método: CancelSubscription | 1h | ✅ Concluído | 100% |
| AS-007 | Retry com Backoff Exponencial | 1h | ✅ Concluído | 100% |
| **Webhooks** |
| AS-008 | Handler: POST /webhooks/asaas | 2h | ✅ Concluído | 100% |
| AS-009 | Validação de Signature | 1h | ✅ Concluído | 100% |
| AS-010 | Processamento de Eventos | 3h | ✅ Concluído | 100% |
| **Integração nos Use Cases** |
| AS-011 | Integrar em CreateSubscription | 1h | ✅ Concluído | 100% |
| AS-012 | Integrar em CancelSubscription | 1h | ✅ Concluído | 100% |
| AS-013 | Fallback para Manual | 1h | ✅ Concluído | 100% |
| **Configuração** |
| AS-014 | Variáveis de Ambiente | 30min | ✅ Concluído | 100% |
| AS-015 | Compilação e Validação | 2h | ✅ Concluído | 100% |

**📈 PROGRESSO SPRINT: 15/15 (100%)**

---

## 📁 Arquivos Criados/Modificados

### Criados:
- `backend/internal/infra/gateway/asaas/types.go` — DTOs Asaas (Customer, Subscription, Webhook)
- `backend/internal/infra/gateway/asaas/customer.go` — Métodos FindOrCreateCustomer, SearchByNameAndPhone
- `backend/internal/infra/gateway/asaas/subscription.go` — Métodos CreateSubscription, CancelSubscription, GetSubscription, CreatePaymentLink
- `backend/internal/infra/gateway/asaas/gateway_adapter.go` — Adapter implementando port.AsaasGateway
- `backend/internal/infra/http/handler/webhook_handler.go` — Handler POST /webhooks/asaas
- `backend/internal/application/usecase/subscription/process_webhook_usecase.go` — Processamento de eventos webhook

### Modificados:
- `backend/internal/infra/gateway/asaas/client.go` — Já existia, validado retry com backoff
- `backend/internal/domain/port/subscription_repository.go` — Removido método duplicado GetByAsaasID
- `backend/internal/domain/port/asaas_gateway.go` — Interface AsaasGateway
- `backend/internal/infra/db/queries/subscriptions.sql` — Adicionada query UpdateSubscription
- `backend/internal/infra/repository/postgres/subscription_repository.go` — Adicionado método Update
- `backend/internal/infra/repository/postgres/helpers.go` — Adicionado timePtrToPgTimestamptz
- `backend/internal/application/usecase/subscription/create_subscription_usecase.go` — Integração Asaas
- `backend/internal/application/usecase/subscription/cancel_subscription_usecase.go` — Integração Asaas
- `backend/cmd/api/main.go` — Injeção AsaasGateway nos use cases
- `backend/.env.example` — Variáveis ASAAS_API_KEY, ASAAS_ENV

---

## 📋 Tarefas Detalhadas

### 🌐 FASE 1: Gateway HTTP

#### AS-001: Cliente HTTP Base

**Objetivo:** Criar cliente HTTP para comunicação com API Asaas

**Referência:** [FLUXO_ASSINATURA.md — Seção 4.5](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#45-regras-de-integração)

**Arquivo:** `backend/internal/infrastructure/external/asaas/client.go`

```go
package asaas

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"

    "go.uber.org/zap"
)

type Client struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
    logger     *zap.Logger
}

type Config struct {
    APIKey  string
    BaseURL string // https://sandbox.asaas.com/api/v3 ou https://api.asaas.com/api/v3
    Timeout time.Duration
}

func NewClient(cfg Config, logger *zap.Logger) *Client {
    return &Client{
        apiKey:  cfg.APIKey,
        baseURL: cfg.BaseURL,
        httpClient: &http.Client{
            Timeout: cfg.Timeout,
        },
        logger: logger,
    }
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, int, error) {
    var reqBody io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, 0, fmt.Errorf("marshal request body: %w", err)
        }
        reqBody = bytes.NewBuffer(jsonBody)
    }

    req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
    if err != nil {
        return nil, 0, fmt.Errorf("create request: %w", err)
    }

    req.Header.Set("access_token", c.apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, 0, fmt.Errorf("execute request: %w", err)
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, resp.StatusCode, fmt.Errorf("read response body: %w", err)
    }

    return respBody, resp.StatusCode, nil
}
```

---

#### AS-002: Types/DTOs Asaas

**Objetivo:** Mapear estruturas de request/response da API Asaas

**Arquivo:** `backend/internal/infrastructure/external/asaas/types.go`

```go
package asaas

type CustomerRequest struct {
    Name                 string `json:"name"`
    CpfCnpj              string `json:"cpfCnpj"`
    Email                string `json:"email,omitempty"`
    Phone                string `json:"phone,omitempty"`
    MobilePhone          string `json:"mobilePhone,omitempty"`
    ExternalReference    string `json:"externalReference"`
    NotificationDisabled bool   `json:"notificationDisabled"`
}

type CustomerResponse struct {
    ID string `json:"id"`
    // ... outros campos
}

type SubscriptionRequest struct {
    Customer             string  `json:"customer"`
    BillingType          string  `json:"billingType"` // CREDIT_CARD
    Value                float64 `json:"value"`
    NextDueDate          string  `json:"nextDueDate"`
    Cycle                string  `json:"cycle"` // MONTHLY
    Description          string  `json:"description"`
    ExternalReference    string  `json:"externalReference"`
}

type SubscriptionResponse struct {
    ID string `json:"id"`
    // ... outros campos
}
```

---

#### AS-003: Método `FindOrCreateCustomer`

**Referência:** [FLUXO_ASSINATURA.md — Regra AS-001, AS-002, RN-CLI-002](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#45-regras-de-integração)

**Lógica:**
1. Verificar se cliente local já possui `asaas_customer_id`
   - Se sim, retornar esse ID (evitar duplicação - RN-CLI-002)
2. Se não, buscar cliente no Asaas por `nome + telefone` (AS-001)
3. Se existir no Asaas:
   - Salvar `asaas_customer_id` no cliente local (RN-CLI-009)
   - Retornar ID
4. Se não existir:
   - Criar novo cliente no Asaas (sem exigir CPF - AS-002, AS-003)
   - Salvar `asaas_customer_id` no cliente local
   - Retornar ID

---

#### AS-004: Método `CreateSubscription`

**Referência:** [FLUXO_ASSINATURA.md — Regra AS-002](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#45-regras-de-integração)

**Lógica:**
1. Receber ID do cliente Asaas e dados do plano
2. Criar assinatura com ciclo MENSAL
3. Retornar ID da assinatura Asaas

---

#### AS-005: Método `GeneratePaymentLink`

**Referência:** [FLUXO_ASSINATURA.md — Regra AS-004](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#45-regras-de-integração)

**Lógica:**
1. Obter ID da assinatura
2. Buscar fatura pendente
3. Retornar `invoiceUrl` ou `bankSlipUrl`

---

#### AS-006: Método `CancelSubscription`

**Referência:** [FLUXO_ASSINATURA.md — Regra AS-006](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#45-regras-de-integração)

**Lógica:**
1. Receber ID da assinatura Asaas
2. Enviar DELETE /subscriptions/{id}
3. Validar se foi cancelada (deleted: true)

---

#### AS-007: Retry com Backoff Exponencial

**Referência:** [FLUXO_ASSINATURA.md — Regra AS-008](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#45-regras-de-integração)

**Lógica:**
- Implementar retry para erros 5xx e timeouts
- Máximo 3 tentativas
- Backoff: 1s, 2s, 4s

---

### 🎣 FASE 2: Webhooks

#### AS-008: Handler `POST /webhooks/asaas`

**Arquivo:** `backend/internal/interfaces/http/handlers/webhook_handler.go`

**Lógica:**
1. Receber POST do Asaas
2. Validar signature (AS-009)
3. Identificar tipo de evento
4. Enviar para processamento assíncrono (goroutine ou queue)
5. Retornar 200 OK imediatamente

---

#### AS-009: Validação de Signature

**Referência:** [FLUXO_ASSINATURA.md — Seção 4.4](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#44-webhooks)

**Lógica:**
- Ler header `asaas-access-token`
- Comparar com token configurado no env `ASAAS_WEBHOOK_TOKEN`
- Se diferente, retornar 401 Unauthorized

---

#### AS-010: Processamento de Eventos

**Referência:** [FLUXO_ASSINATURA.md — Seção 6.6](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#66-fluxo-processar-webhook)

**Eventos a tratar:**
- `PAYMENT_CONFIRMED`:
  - Buscar assinatura pelo `externalReference` ou `subscriptionId`
  - Atualizar status = ATIVO
  - Atualizar data_vencimento
  - Registrar pagamento em `subscription_payments`
  - **Marcar cliente como is_subscriber = true (RN-CLI-003)**
  
- `PAYMENT_OVERDUE`:
  - Atualizar status = INADIMPLENTE
  - **Verificar se cliente tem outras assinaturas ativas (RN-CLI-004)**
  - **Se não, atualizar is_subscriber = false**
  
- `SUBSCRIPTION_DELETED`:
  - Atualizar status = CANCELADO
  - **Verificar se cliente tem outras assinaturas ativas (RN-CLI-004)**
  - **Se não, atualizar is_subscriber = false**

---

### 🔄 FASE 3: Integração nos Use Cases
**Alteração:**
- Injetar `AsaasGateway`
- Se forma_pagamento = CARTAO:
  - Chamar `FindOrCreateCustomer` (garante reuso/unificação - RN-CLI-002, RN-CLI-009)
  - Chamar `CreateSubscription`
  - Salvar `asaas_customer_id` e `asaas_subscription_id` no banco
  - Aguardar confirmação via webhook para marcar is_subscriber
- Injetar `AsaasGateway`
- Se forma_pagamento = CARTAO:
  - Chamar `FindOrCreateCustomer`
  - Chamar `CreateSubscription`
  - Salvar `asaas_customer_id` e `asaas_subscription_id` no banco

---

#### AS-012: Integrar em `CancelSubscription`

**Arquivo:** `backend/internal/application/usecase/subscription/cancel_subscription.go`

**Alteração:**
- Se `asaas_subscription_id` existir:
  - Chamar `CancelSubscription` no gateway
  - Se erro, logar mas continuar cancelamento local (soft fail)

---

#### AS-013: Fallback para Manual

**Referência:** [FLUXO_ASSINATURA.md — Regra AS-007](../../docs/11-Fluxos/FLUXO_ASSINATURA.md#45-regras-de-integração)

**Lógica:**
- Se Asaas estiver fora do ar (erro 5xx persistente):
  - Permitir criar assinatura como "AGUARDANDO_INTEGRACAO"
  - Cron job posterior tenta sincronizar

---

### ⚙️ FASE 4: Configuração

#### AS-014: Variáveis de Ambiente

**Arquivo:** `.env.example`

```bash
ASAAS_API_KEY=your_api_key
ASAAS_API_URL=https://sandbox.asaas.com/api/v3
ASAAS_WEBHOOK_TOKEN=your_webhook_token
```

---

#### AS-015: Testes com Sandbox Asaas

**Objetivo:** Validar fluxo completo em ambiente de sandbox

**Checklist:**
- [x] Criar conta sandbox no Asaas
- [x] Gerar API Key
- [x] Configurar URL de webhook (testado localmente com token 123456)
- [x] Simular pagamento confirmado (PAYMENT_CONFIRMED)
- [x] Simular pagamento vencido (PAYMENT_OVERDUE)

**Resultados dos Testes (03/12/2025):**
- ✅ Conexão Sandbox: OK (8 clientes existentes)
- ✅ Webhook PAYMENT_CONFIRMED: Recebido e processado (retornou 200)
- ✅ Webhook PAYMENT_OVERDUE: Recebido e processado (retornou 200)
- ✅ Token inválido: Rejeitado corretamente (retornou 401)
- ✅ Criar cliente no Asaas: `cus_000007273681`
- ✅ Criar assinatura no Asaas: `sub_glqmgn8ixixzg57c`

---

## ✅ Critérios de Conclusão da Sprint

- [x] Gateway HTTP implementado e testado
- [x] Webhook handler recebendo e validando eventos
- [x] Processamento de eventos atualizando status no banco
- [x] Use cases integrados com Asaas
- [x] Teste ponta a ponta na sandbox (criar cliente -> criar assinatura -> webhook)

---

## 🔗 Próxima Sprint

Após conclusão, iniciar **Sprint 4: Frontend**
📂 [04-FRONTEND.md](./04-FRONTEND.md)

---

**FIM DO DOCUMENTO — SPRINT 3: INTEGRAÇÃO ASAAS**
