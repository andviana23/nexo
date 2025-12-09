# Fluxo de Assinatura — NEXO / Asaas

> **Versão:** 2.1  
> **Última Atualização:** 04/12/2025  
> **Status:** FONTE DA VERDADE  
> **Criticidade:** 🔴 MÁXIMA — Módulo financeiro core

---

## 📌 Índice

1. [Estrutura do Módulo](#1-estrutura-do-módulo)
2. [Página: Planos](#2-página-planos)
3. [Página: Assinantes](#3-página-assinantes)
4. [Integração Asaas](#4-integração-asaas)
5. [Página: Relatórios](#5-página-relatórios)
6. [Fluxos Detalhados](#6-fluxos-detalhados)
7. [Diagrama Mermaid](#7-diagrama-mermaid)
8. [Regras de Negócio](#8-regras-de-negócio)
9. [Tratamento de Erros](#9-tratamento-de-erros)
10. [Checklist de Implementação](#10-checklist-de-implementação)

---

## 1. Estrutura do Módulo

### 1.1 Navegação

\`\`\`
Sidebar → Assinaturas
├── Planos          (CRUD de modelos de plano)
├── Assinantes      (Gestão de assinaturas ativas)
└── Relatórios      (Análises e métricas)
\`\`\`

### 1.2 Permissões por Página

| Página | Administrador | Gerente | Recepção | Barbeiro |
|--------|---------------|---------|----------|----------|
| Planos | CRUD completo | CRUD completo | Visualizar | ❌ |
| Assinantes | CRUD + Cancelar | CRUD + Cancelar | Criar + Visualizar | ❌ |
| Relatórios | Visualizar | Visualizar | Visualizar | ❌ |

---

## 2. Página: Planos

### 2.1 Objetivo

Criar **modelos de planos** que serão utilizados no momento da venda. Os planos **NÃO são criados no Asaas** neste momento — servem apenas como template interno.

### 2.2 Campos do Cadastro

| Campo | Tipo | Obrigatório | Validação |
|-------|------|-------------|-----------|
| \`nome\` | string | ✅ | 3-100 caracteres |
| \`descricao\` | text | ❌ | max 500 caracteres |
| \`valor\` | decimal | ✅ | > 0, máx 2 casas decimais |
| \`periodicidade\` | enum | ✅ | \`MENSAL\` (fixo v1.0) |
| \`qtd_servicos\` | int | ❌ | ≥ 0 (null = ilimitado) |
| \`limite_uso_mensal\` | int | ❌ | ≥ 0 (null = ilimitado) |
| \`ativo\` | boolean | ✅ | default: true |
| \`tenant_id\` | uuid | ✅ | extraído do contexto |

### 2.3 Regras de Negócio

\`\`\`
REGRA PL-001: Plano NÃO é enviado ao Asaas na criação
REGRA PL-002: Plano inativo NÃO aparece na seleção de nova assinatura
REGRA PL-003: Não pode excluir plano com assinaturas ativas (apenas desativar)
REGRA PL-004: Alteração de valor NÃO afeta assinaturas existentes
REGRA PL-005: Nome do plano deve ser único por tenant
\`\`\`

### 2.4 Fluxo: Criar Plano

\`\`\`
[Admin/Gerente abre Assinaturas → Planos]
   ↓
[Clica "Novo Plano"]
   ↓
[Preenche formulário]
   ↓
[Valida campos obrigatórios]
   ↓
   └── Erro? → [Exibe mensagem de validação] → [Retorna ao form]
   ↓
[Salva plano no banco de dados]
   ↓
[Status: Ativo]
   ↓
[Exibe toast: "Plano criado com sucesso"]
   ↓
[Fim]
\`\`\`

---

## 3. Página: Assinantes

### 3.1 Objetivo

Gestão REAL das assinaturas — tanto locais (PIX/Dinheiro) quanto integradas ao Asaas (Cartão).

### 3.2 Lista de Assinantes

| Coluna | Descrição |
|--------|-----------|
| Cliente | Nome completo + telefone |
| Plano | Nome do plano |
| Status | Badge colorido |
| Vencimento | Data próximo vencimento |
| Forma Pagamento | Cartão / PIX / Dinheiro |
| Ações | Ver, Renovar, Cancelar |

### 3.3 Status Possíveis

| Status | Cor | Descrição | Origem |
|--------|-----|-----------|--------|
| \`ATIVO\` | 🟢 Verde | Pagamento confirmado, dentro da validade | Asaas ou Manual |
| \`AGUARDANDO_PAGAMENTO\` | 🟡 Amarelo | Link gerado, aguardando cliente pagar | Asaas |
| \`INADIMPLENTE\` | 🔴 Vermelho | Venceu sem pagamento | Asaas ou Sistema |
| \`INATIVO\` | ⚫ Cinza | Cancelado pelo cliente ou admin | Asaas ou Manual |
| \`CANCELADO\` | ⚫ Cinza | Cancelado definitivamente | Asaas |

### 3.4 Formas de Pagamento

| Forma | Automático | Responsável | Renovação |
|-------|------------|-------------|-----------|
| **Cartão de Crédito** | ✅ Sim | Asaas | Automática mensal |
| **PIX** | ❌ Não | Recepção | Manual (30 dias) |
| **Dinheiro** | ❌ Não | Recepção | Manual (30 dias) |

### 3.5 Estados do Cliente (flag “assinante”)

| Campo | Valores | Quem altera | Quando |
|-------|---------|------------|--------|
| `cliente_tipo` (ou flag `is_subscriber`) | `CLIENTE_COMUM` / `CLIENTE_ASSINANTE` | Backend | - Sobe para `CLIENTE_ASSINANTE` quando existir assinatura ATIVA (qualquer forma de pagamento).<br>- Retorna para `CLIENTE_COMUM` quando o cliente não possuir mais assinaturas ativas. |

> Esta flag controla acesso aos benefícios e deve ser persistida no cadastro do cliente, não apenas na assinatura.

---

## 4. Integração Asaas

### 4.1 Endpoints Utilizados

| Ação | Método | Endpoint Asaas |
|------|--------|----------------|
| Buscar cliente por Nome e Telefone | GET | \`/customers?name={nome}&mobilePhone={telefone}\` |
| Criar cliente | POST | \`/customers\` |
| Criar assinatura | POST | \`/subscriptions\` |
| Gerar link pagamento | POST | \`/subscriptions/{id}/paymentLink\` |
| Cancelar assinatura | DELETE | \`/subscriptions/{id}\` |
| Consultar assinatura | GET | \`/subscriptions/{id}\` |

### 4.2 Mapeamento de Status

| Status Asaas | Status Sistema | Ação |
|--------------|----------------|------|
| \`ACTIVE\` | \`ATIVO\` | Assinatura válida |
| \`PENDING\` | \`AGUARDANDO_PAGAMENTO\` | Aguardar pagamento |
| \`OVERDUE\` | \`INADIMPLENTE\` | Notificar cliente |
| \`INACTIVE\` | \`INATIVO\` | Assinatura pausada |
| \`CANCELED\` | \`CANCELADO\` | Assinatura encerrada |

### 4.3 Webhooks Obrigatórios

> ⚠️ **IMPORTANTE**: A sequência correta é `CONFIRMED` → `RECEIVED`. O pagamento é **confirmado** primeiro, depois o dinheiro é **recebido** (compensado).

| # | Evento Asaas | Ação no Sistema | Campo BD | Regime Contábil |
|---|--------------|-----------------|----------|------------------|
| 1 | `PAYMENT_CONFIRMED` | ✅ Ativar assinatura, liberar acesso | `confirmed_at` | 📊 DRE (Competência) |
| 2 | `PAYMENT_RECEIVED` | 💰 Registrar recebimento em caixa | `received_at` | 💵 Fluxo de Caixa |
| 3 | `SUBSCRIPTION_ACTIVATED` | Atualizar status para `ATIVO` | - | - |
| 4 | `SUBSCRIPTION_RENEWED` | Atualizar `data_vencimento` +30 dias | - | - |
| 5 | `SUBSCRIPTION_CANCELED` | Atualizar status para `CANCELADO` | - | - |
| 6 | `PAYMENT_OVERDUE` | Atualizar status para `INADIMPLENTE` | - | - |
| 7 | `PAYMENT_REFUNDED` | Atualizar status para `INATIVO`, registrar estorno | - | - |

### 4.3.1 Sequência Correta de Status (Pagamentos)

```
PENDING → CONFIRMED → RECEIVED
   ↓          ↓           ↓
 Criado    Pagou      Compensou
```

| Status | Significado | Webhook | Quando Ocorre |
|--------|-------------|---------|---------------|
| `PENDING` | Cobrança criada, aguardando pagamento | - | Criação da cobrança |
| `CONFIRMED` | Cliente pagou, pagamento confirmado | `payment.confirmed` | Cliente finaliza pagamento |
| `RECEIVED` | Dinheiro compensado e disponível | `payment.received` | D+0 a D+2 após confirmação |
| `OVERDUE` | Vencido sem pagamento | `payment.overdue` | Após data de vencimento |

### 4.3.2 Regimes Contábeis (DRE vs Caixa)

> 🎯 **Separação obrigatória** para relatórios financeiros corretos.

| Regime | Campo Utilizado | Use Case | Relatório |
|--------|-----------------|----------|------------|
| **Competência** | `confirmed_at` | `GenerateDREV2UseCase` | DRE - Receita reconhecida quando cliente pagou |
| **Caixa** | `received_at` | `GenerateFluxoDiarioV2UseCase` | Fluxo de Caixa - Dinheiro disponível para saque |

**Por que liberar acesso no CONFIRMED (não no RECEIVED)?**
- Quando `CONFIRMED`, o cliente já pagou - não faz sentido esperar D+2 para liberar
- Risco de estorno é baixíssimo após confirmação
- Melhor experiência para o cliente

### 4.4 Payload Webhook (Exemplo)

\`\`\`json
{
  "event": "PAYMENT_CONFIRMED",
  "payment": {
    "id": "pay_abc123",
    "subscription": "sub_xyz789",
    "customer": "cus_def456",
    "value": 99.90,
    "status": "CONFIRMED",
    "confirmedDate": "2025-11-27"
  }
}
\`\`\`

### 4.5 Regras de Integração

\`\`\`
REGRA AS-001: Busca de cliente no Asaas é SEMPRE por Nome + Telefone (NUNCA por CPF/CNPJ)
REGRA AS-002: Cliente pode existir no Asaas sem CPF (nome + telefone + email basta)
REGRA AS-003: Assinatura cartão NÃO exige CPF do titular = comprador
REGRA AS-004: Link de pagamento expira em 24h (configurável)
REGRA AS-005: Webhook deve responder 200 em até 5s (senão Asaas reenvia)
REGRA AS-006: Armazenar asaas_customer_id e asaas_subscription_id localmente
REGRA AS-007: NUNCA expor API Key do Asaas no frontend
REGRA AS-008: Todas as chamadas Asaas devem ter retry (3x com backoff)
REGRA AS-009: Se cliente existir no Asaas **e** no sistema, unificar cadastro salvando o mesmo `asaas_customer_id` no cliente local (não criar duplicatas).
REGRA AS-010: `asaas_customer_id` é único por cliente/tenant; migrações e validações devem impedir associação duplicada.
REGRA AS-011: PAYMENT_CONFIRMED libera acesso ao sistema e registra no DRE (competência); usar campo `confirmed_at`.
REGRA AS-012: PAYMENT_RECEIVED registra no Fluxo de Caixa (caixa); usar campo `received_at`. CONFIRMED vem ANTES de RECEIVED.
\`\`\`

---

## 5. Página: Relatórios

### 5.1 Métricas Exibidas

| Métrica | Descrição | Cálculo |
|---------|-----------|---------|
| **Total Ativas** | Assinaturas com status ATIVO | \`COUNT WHERE status = 'ATIVO'\` |
| **Total Inativas** | Assinaturas canceladas ou inativas | \`COUNT WHERE status IN ('INATIVO', 'CANCELADO')\` |
| **Por Forma de Pagamento** | Breakdown por tipo | \`GROUP BY forma_pagamento\` |
| **Por Plano** | Breakdown por plano | \`GROUP BY plano_id\` |
| **Receita Mensal** | Soma dos valores ativos | \`SUM(valor) WHERE status = 'ATIVO'\` |
| **Taxa Cancelamento** | % canceladas no mês | \`(canceladas / total_criadas) * 100\` |
| **Churn** | Canceladas / Ativas início mês | \`(canceladas / ativas_inicio) * 100\` |

### 5.2 Filtros

- Período (data início / data fim)
- Status
- Forma de pagamento
- Plano

---

## 6. Fluxos Detalhados

### 6.1 Fluxo: Nova Assinatura (Cartão de Crédito)

\`\`\`
[Recepção clica "Nova Assinatura"]
   ↓
[Buscar cliente no sistema]
   ↓
   ├── Encontrou? → [Usar cliente existente]
   └── Não encontrou? → [Abrir modal "Novo Cliente"]
                              ↓
                         [Preencher: nome, telefone, email, cpf]
                              ↓
                         [Salvar cliente no sistema]
   ↓
[Verificar cliente no Asaas via API]
   GET /customers?name={nome}&mobilePhone={telefone}
   ↓
   ├── Existe (encontrou por nome + telefone)? → [Recuperar asaas_customer_id]
   │     ↓
   │   [Se cliente local possui outro asaas_customer_id ou está sem ID] → [Unificar: gravar asaas_customer_id no cliente local, evitar duplicação]
   └── Não existe? → [POST /customers → criar no Asaas]
                           ↓
                      [Salvar asaas_customer_id no cliente local]
   ↓
[Marcar cliente como CLIENTE_ASSINANTE se ainda não estiver]
   ↓
[Selecionar Plano]
   ↓
[Selecionar Forma: "Cartão de Crédito"]
   ↓
[POST /subscriptions no Asaas]
   Payload: {
     customer: asaas_customer_id,
     billingType: "CREDIT_CARD",
     value: plano.valor,
     cycle: "MONTHLY",
     description: plano.nome
   }
   ↓
[Receber subscription_id do Asaas]
   ↓
[POST /subscriptions/{id}/paymentLink]
   ↓
[Receber URL do link de pagamento]
   ↓
[Salvar assinatura local]
   {
     cliente_id,
     plano_id,
     asaas_subscription_id,
     forma_pagamento: "CARTAO",
     status: "AGUARDANDO_PAGAMENTO",
     link_pagamento: url
   }
   ↓
[Exibir modal com link + botão "Enviar via WhatsApp"]
   ↓
[Recepção envia link ao cliente]
   ↓
[Cliente paga no checkout Asaas]
   ↓
[Asaas envia webhook: PAYMENT_CONFIRMED]
   ↓
[Backend atualiza assinatura]
   {
     status: "ATIVO",
     data_ativacao: now(),
     data_vencimento: now() + 30 dias
   }
   ↓
[Renovação automática pelo Asaas a cada 30 dias]
   ↓
[Fim]
\`\`\`

### 6.2 Fluxo: Nova Assinatura (PIX)

\`\`\`
[Recepção clica "Nova Assinatura"]
   ↓
[Buscar/Criar cliente no sistema]
   ↓
[Selecionar Plano]
   ↓
[Selecionar Forma: "PIX"]
   ↓
[Exibir formulário de confirmação manual]
   ├── Data da transação (obrigatório)
   ├── Hora da transação (obrigatório)
   └── Código/ID transação (opcional)
   ↓
[Recepção confirma que recebeu o PIX]
   ↓
[Salvar assinatura local]
   {
     cliente_id,
     plano_id,
     asaas_subscription_id: null,  // NÃO cria no Asaas
     forma_pagamento: "PIX",
     status: "ATIVO",
     data_ativacao: now(),
     data_vencimento: now() + 30 dias,
     codigo_transacao: "xxx"
   }
   ↓
[Exibir toast: "Assinatura ativada com sucesso"]
   ↓
[Cron job diário verifica vencimentos]
   ↓
   └── data_vencimento < now()?
         ↓
       [Atualizar status: "INADIMPLENTE"]
   ↓
[Fim]
\`\`\`

### 6.3 Fluxo: Nova Assinatura (Dinheiro)

\`\`\`
[Recepção clica "Nova Assinatura"]
   ↓
[Buscar/Criar cliente no sistema]
   ↓
[Selecionar Plano]
   ↓
[Selecionar Forma: "Dinheiro"]
   ↓
[Confirmar recebimento]
   ↓
[Salvar assinatura local]
   {
     cliente_id,
     plano_id,
     asaas_subscription_id: null,  // NÃO cria no Asaas
     forma_pagamento: "DINHEIRO",
     status: "ATIVO",
     data_ativacao: now(),
     data_vencimento: now() + 30 dias
   }
   ↓
[Registrar no caixa como entrada]
   ↓
[Exibir toast: "Assinatura ativada"]
   ↓
[Cron job diário verifica vencimentos]
   ↓
   └── data_vencimento < now()?
         ↓
       [Atualizar status: "INADIMPLENTE"]
   ↓
[Fim]
\`\`\`

### 6.4 Fluxo: Renovar Assinatura Manual (PIX/Dinheiro)

\`\`\`
[Recepção abre assinatura inadimplente]
   ↓
[Clica "Renovar"]
   ↓
[Confirma forma de pagamento]
   ├── PIX → [Preenche data/hora transação]
   └── Dinheiro → [Confirma recebimento]
   ↓
[Atualizar assinatura]
   {
     status: "ATIVO",
     data_ativacao: now(),
     data_vencimento: now() + 30 dias
   }
   ↓
[Registrar no histórico de pagamentos]
   ↓
[Fim]
\`\`\`

### 6.5 Fluxo: Cancelar Assinatura

\`\`\`
[Admin/Gerente abre assinatura]
   ↓
[Clica "Cancelar Assinatura"]
   ↓
[Modal de confirmação: "Tem certeza?"]
   ↓
[Confirma cancelamento]
   ↓
forma_pagamento == "CARTAO"?
   ├── Sim → [DELETE /subscriptions/{id} no Asaas]
   └── Não → [Apenas atualiza local]
   ↓
[Atualizar assinatura local]
   {
     status: "CANCELADO",
     data_cancelamento: now(),
     cancelado_por: user_id
   }
   ↓
[Exibir toast: "Assinatura cancelada"]
   ↓
[Fim]
\`\`\`

### 6.6 Fluxo: Processar Webhook Asaas

\`\`\`
[Asaas envia POST /webhooks/asaas]
   ↓
[Validar assinatura do webhook (header X-Asaas-Signature)]
   ↓
   └── Inválido? → [Retornar 401] → [Log warning]
   ↓
[Extrair event e payload]
   ↓
[Buscar assinatura por asaas_subscription_id]
   ↓
   └── Não encontrada? → [Retornar 200] → [Log warning: "orphan webhook"]
   ↓
[Switch por evento]
   │
   ├── PAYMENT_CONFIRMED (vem PRIMEIRO):
   │     status = "ATIVO"
   │     confirmed_at = confirmedDate
   │     data_ativacao = confirmedDate
   │     data_vencimento = confirmedDate + 30 dias
   │     → Registrar no DRE (regime competência)
   │     → Liberar acesso ao sistema
   │     
   ├── PAYMENT_RECEIVED (vem DEPOIS):
   │     received_at = creditDate
   │     → Registrar no Fluxo de Caixa (regime caixa)
   │     → Dinheiro disponível para saque
   │     
   ├── SUBSCRIPTION_RENEWED:
   │     data_vencimento = dueDate
   │     
   ├── PAYMENT_OVERDUE:
   │     status = "INADIMPLENTE"
   │     
   ├── SUBSCRIPTION_CANCELED:
   │     status = "CANCELADO"
   │     data_cancelamento = now()
   │     
   └── PAYMENT_REFUNDED:
         status = "INATIVO"
         registrar_estorno(payment_id, value)
   ↓
[Salvar assinatura atualizada]
   ↓
[Retornar 200 OK]
   ↓
[Fim]
\`\`\`

---

## 7. Diagrama Mermaid

### 7.1 Fluxo Completo de Webhooks Asaas (Sequência Correta)

> 📌 Este diagrama mostra a **sequência correta** dos webhooks: `CONFIRMED` vem **ANTES** de `RECEIVED`.

\`\`\`mermaid
graph TD
    Start([👤 Cliente Contrata Assinatura]) --> CreateSub[📝 Criar Assinatura no Asaas via API]
    CreateSub --> CalcProp{💰 Calcular Valor Proporcional?}
    CalcProp -->|Sim| CreateCharge[📄 Criar Primeira Cobrança Status: PENDING]
    CalcProp -->|Não Proporcional| CreateCharge
    CreateCharge --> AwaitPay[⏳ AGUARDANDO PAGAMENTO<br/>Status: PENDING]
    
    AwaitPay --> CheckPay{💳 Cliente Pagou?}
    
    CheckPay -->|❌ Não| CheckDue{📅 Passou do Vencimento?}
    CheckDue -->|Não| AwaitPay
    CheckDue -->|⚠️ Sim| Overdue[🔴 VENCIDAS<br/>Status: OVERDUE]
    Overdue --> Notify[📧 Enviar Notificação de Cobrança]
    Notify --> Retry{💳 Cliente Pagou Após Vencimento?}
    Retry -->|❌ Não| Cancel[🚫 Cancelar Assinatura<br/>🔒 Bloquear Acesso ao Sistema]
    Retry -->|✅ Sim| WebhookConf
    
    CheckPay -->|✅ Sim| WebhookConf[🔔 Webhook payment.confirmed<br/>Cliente pagou confirmado]
    WebhookConf --> UpdateConf[🔄 Atualizar BD<br/>Status: CONFIRMED<br/>confirmed_at = now]
    UpdateConf --> Confirmed[🔵 CONFIRMADAS<br/>Status: CONFIRMED]
    
    Confirmed --> Active[✨ Ativar Funcionalidades<br/>🔓 Liberar Acesso Total<br/>Cliente pode usar sistema]
    Active --> RecordDRE[📊 Registrar no DRE<br/>Regime de Competência]
    RecordDRE --> WaitClearing[⏱️ Aguardar Compensação Bancária<br/>Prazo: D+0 a D+2]
    
    WaitClearing --> WebhookRec[🔔 Webhook payment.received<br/>Dinheiro creditado na conta]
    WebhookRec --> UpdateReceived[🔄 Atualizar BD<br/>Status: RECEIVED<br/>received_at = now]
    UpdateReceived --> Received[🟢 RECEBIDAS<br/>Status: RECEIVED]
    
    Received --> Available[💵 Valor Disponível para Saque<br/>Dinheiro Líquido]
    Available --> RecordFluxo[💰 Registrar no Fluxo de Caixa<br/>Regime de Caixa]
    
    RecordFluxo --> CheckDate{📆 Chegou Data de Vencimento<br/>da Próxima Cobrança?}
    CheckDate -->|✅ Sim| GenNext[🔄 Asaas Gera Próxima Cobrança<br/>Automaticamente - Recorrência Mensal]
    GenNext --> NextPending[📄 Nova Cobrança Criada<br/>Status: PENDING<br/>Vencimento: Próximo Mês]
    NextPending --> AwaitPay
    
    CheckDate -->|⏳ Não| Monitor[👁️ Monitorar Assinatura<br/>Status: ACTIVE<br/>Cliente usando sistema]
    Monitor --> CheckDate
    
    Cancel --> End([⛔ Fim da Assinatura<br/>Cliente Inativo])
    
    style Start fill:#e3f2fd,stroke:#1976d2,stroke-width:3px
    style AwaitPay fill:#fff3e0,stroke:#ff9800,stroke-width:4px
    style Confirmed fill:#e3f2fd,stroke:#2196f3,stroke-width:4px
    style Received fill:#e8f5e9,stroke:#4caf50,stroke-width:4px
    style Overdue fill:#ffebee,stroke:#f44336,stroke-width:4px
    style Active fill:#f3e5f5,stroke:#9c27b0,stroke-width:3px
    style Available fill:#c8e6c9,stroke:#388e3c,stroke-width:3px
    style Cancel fill:#ffcdd2,stroke:#c62828,stroke-width:3px
    style End fill:#fce4ec,stroke:#880e4f,stroke-width:3px
    style WebhookConf fill:#fff9c4,stroke:#f57f17,stroke-width:2px
    style WebhookRec fill:#fff9c4,stroke:#f57f17,stroke-width:2px
    style GenNext fill:#e1bee7,stroke:#8e24aa,stroke-width:2px
    style RecordDRE fill:#ede7f6,stroke:#5e35b1,stroke-width:2px
    style RecordFluxo fill:#e0f2f1,stroke:#00897b,stroke-width:2px
\`\`\`

### 7.2 Fluxo Geral de Assinaturas

\`\`\`mermaid
flowchart TD
    subgraph INICIO["🏠 Início"]
        A[Usuário acessa Assinaturas]
    end

    subgraph CLIENTE["👤 Cliente"]
        B[Nova Assinatura]
        C{Cliente existe<br>no sistema?}
        D[Criar cliente]
        E{Cliente existe<br>no Asaas?}
        F[Criar cliente<br>no Asaas]
        G[Recuperar<br>asaas_customer_id]
    end

    subgraph PLANO["📋 Plano"]
        H[Selecionar Plano]
    end

    subgraph PAGAMENTO["💳 Forma de Pagamento"]
        I{Qual forma?}
        
        subgraph CARTAO["Cartão de Crédito"]
            J[Criar assinatura<br>no Asaas]
            K[Gerar link<br>pagamento]
            L[Enviar link<br>WhatsApp]
            M[Cliente paga]
            N1[Webhook:<br>PAYMENT_CONFIRMED]
            N2[Webhook:<br>PAYMENT_RECEIVED]
            O[Ativar assinatura]
        end
        
        subgraph PIX["PIX"]
            P[Recepção registra<br>data/hora transação]
            Q[Ativar assinatura<br>30 dias]
        end
        
        subgraph DINHEIRO["Dinheiro"]
            R[Confirmar<br>recebimento]
            S[Ativar assinatura<br>30 dias]
        end
    end

    subgraph GESTAO["⚙️ Gestão Contínua"]
        T[Cron job diário<br>verificar vencimentos]
        U{Venceu?}
        V[Status: INADIMPLENTE]
        W[Sincronização<br>webhooks Asaas]
    end

    subgraph RELATORIO["📊 Relatórios"]
        X[Métricas:<br>ativos, inativos,<br>receita, churn]
    end

    A --> B
    B --> C
    C -->|Não| D --> E
    C -->|Sim| E
    E -->|Não| F --> H
    E -->|Sim| G --> H
    H --> I

    I -->|Cartão| J
    J --> K --> L --> M --> N1 --> N2 --> O

    I -->|PIX| P --> Q
    I -->|Dinheiro| R --> S

    O --> W
    Q --> T
    S --> T

    T --> U
    U -->|Sim| V
    U -->|Não| W

    W --> X
    V --> X
\`\`\`

---

## 8. Regras de Negócio

### 8.1 Regras de Criação

| Código | Regra | Validação |
|--------|-------|-----------|
| \`RN-SUB-001\` | Cliente é obrigatório | \`cliente_id NOT NULL\` |
| \`RN-SUB-002\` | Plano é obrigatório | \`plano_id NOT NULL\` |
| \`RN-SUB-003\` | Plano deve estar ativo | \`plano.ativo = true\` |
| \`RN-SUB-004\` | Cliente não pode ter assinatura ativa duplicada do mesmo plano | \`UNIQUE(cliente_id, plano_id) WHERE status = 'ATIVO'\` |
| \`RN-SUB-005\` | Valor mínimo: R$ 1,00 | \`valor >= 1.00\` |

### 8.2 Regras de Vencimento

| Código | Regra |
|--------|-------|
| \`RN-VENC-001\` | Assinatura manual (PIX/Dinheiro) vence em 30 dias corridos |
| \`RN-VENC-002\` | Assinatura cartão segue ciclo do Asaas (30 dias) |
| \`RN-VENC-003\` | Cron job roda diariamente às 00:05 para verificar vencimentos |
| \`RN-VENC-004\` | Assinatura vencida há mais de 3 dias → status INADIMPLENTE |

### 8.3 Regras de Cancelamento

| Código | Regra |
|--------|-------|
| \`RN-CANC-001\` | Apenas Admin/Gerente podem cancelar |
| \`RN-CANC-002\` | Cancelamento no cartão deve refletir no Asaas |
| \`RN-CANC-003\` | Registrar quem cancelou e quando |
| \`RN-CANC-004\` | Assinatura cancelada não pode ser reativada (criar nova) |

### 8.4 Regras de Benefícios

| Código | Regra |
|--------|-------|
| \`RN-BEN-001\` | Cliente com assinatura ativa pode usar serviços do plano |
| \`RN-BEN-002\` | Se \`qtd_servicos\` definido, decrementar a cada uso |
| \`RN-BEN-003\` | Se atingir limite, bloquear uso até renovação |
| \`RN-BEN-004\` | Saldo de serviços NÃO acumula entre meses |
| \`RN-BEN-005\` | Flag do cliente deve estar `CLIENTE_ASSINANTE` enquanto houver assinatura ATIVA; remover flag quando não houver ativa |

### 8.5 Regras de Cliente / Unificação

| Código | Regra | Validação |
|--------|-------|-----------|
| \`RN-CLI-001\` | `asaas_customer_id` único por cliente/tenant | Constraint UNIQUE(tenant_id, asaas_customer_id) |
| \`RN-CLI-002\` | Na criação/renovação, sempre tentar reuse/merge do cliente Asaas antes de criar novo | Busca por nome+telefone e comparação de IDs |
| \`RN-CLI-003\` | Alterar `cliente_tipo` para `CLIENTE_ASSINANTE` quando existir assinatura ATIVA | Atualização no serviço de assinatura / webhook |
| \`RN-CLI-004\` | Rebaixar para `CLIENTE_COMUM` quando o cliente ficar sem assinaturas ATIVAS | Rotina pós-cancelamento/webhook/cron |

---

## 9. Tratamento de Erros

### 9.1 Erros de API Asaas

| Código Asaas | Significado | Ação no Sistema |
|--------------|-------------|-----------------|
| \`400\` | Payload inválido | Log error, exibir mensagem genérica ao usuário |
| \`401\` | API Key inválida | Log critical, alertar DevOps |
| \`404\` | Recurso não encontrado | Sincronizar: remover \`asaas_subscription_id\` local |
| \`422\` | Validação falhou | Exibir campos com erro |
| \`429\` | Rate limit | Retry com backoff exponencial (1s, 2s, 4s) |
| \`500\` | Erro interno Asaas | Retry 3x, depois fallback manual |

### 9.2 Mensagens de Erro (UI)

| Cenário | Mensagem |
|---------|----------|
| Falha criar cliente Asaas | "Não foi possível processar. Tente novamente." |
| Link expirado | "O link de pagamento expirou. Gere um novo." |
| Assinatura não encontrada | "Assinatura não encontrada no sistema." |
| Cliente já tem assinatura ativa | "Este cliente já possui uma assinatura ativa deste plano." |

### 9.3 Fallback para Pagamento Manual

Se a integração Asaas falhar após 3 tentativas:

\`\`\`
[Exibir modal]
   "Ocorreu um erro na integração com o gateway de pagamento.
    Deseja registrar a assinatura manualmente (PIX/Dinheiro)?"
   ↓
   [Sim] → [Fluxo PIX/Dinheiro]
   [Não] → [Cancelar operação]
\`\`\`

---

## 10. Checklist de Implementação

### 10.1 Backend

- [ ] **Entidade:** \`Subscription\` com todos os campos
- [ ] **Entidade:** \`Plan\` com todos os campos
- [ ] **Entidade Cliente:** campo \`asaas_customer_id\` único por tenant + flag \`is_subscriber/cliente_tipo\`
- [ ] **Migração:** adicionar UNIQUE(tenant_id, asaas_customer_id) em clientes e campo boolean/enum para status de assinante
- [ ] **Repository:** \`SubscriptionRepository\` (CRUD + queries)
- [ ] **Repository:** \`PlanRepository\` (CRUD + queries)
- [ ] **Gateway:** \`AsaasGateway\` (client HTTP)
- [ ] **UseCase:** \`CreateSubscriptionUseCase\`
- [ ] **UseCase:** \`CancelSubscriptionUseCase\`
- [ ] **UseCase:** \`RenewSubscriptionUseCase\`
- [ ] **UseCase:** \`ProcessWebhookUseCase\`
- [ ] **Handler:** \`POST /subscriptions\`
- [ ] **Handler:** \`GET /subscriptions\`
- [ ] **Handler:** \`GET /subscriptions/:id\`
- [ ] **Handler:** \`DELETE /subscriptions/:id\`
- [ ] **Handler:** \`POST /subscriptions/:id/renew\`
- [ ] **Handler:** \`POST /webhooks/asaas\`
- [ ] **Handler:** \`GET /plans\`
- [ ] **Handler:** \`POST /plans\`
- [ ] **Handler:** \`PUT /plans/:id\`
- [ ] **Handler:** \`DELETE /plans/:id\`
- [ ] **CronJob:** Verificar vencimentos diariamente
- [ ] **Rotina:** Atualizar \`cliente_tipo\` ao ativar/cancelar assinatura (inclusive via webhook)
- [ ] **Middleware:** Validar webhook signature

### 10.2 Frontend

- [ ] **Página:** \`/assinaturas/planos\` (lista + CRUD)
- [ ] **Página:** \`/assinaturas\` (lista assinantes)
- [ ] **Página:** \`/assinaturas/nova\` (wizard nova assinatura)
- [ ] **Página:** \`/assinaturas/relatorios\` (métricas)
- [ ] **Componente:** \`PlanCard\`
- [ ] **Componente:** \`SubscriptionTable\`
- [ ] **Componente:** \`SubscriptionStatusBadge\`
- [ ] **Componente:** \`PaymentMethodSelector\`
- [ ] **Componente:** \`WhatsAppLinkButton\`
- [ ] **Modal:** \`NewPlanModal\`
- [ ] **Modal:** \`ConfirmCancelModal\`
- [ ] **Modal:** \`ManualPaymentModal\`
- [ ] **Hook:** \`usePlans\`
- [ ] **Hook:** \`useSubscriptions\`
- [ ] **Hook:** \`useCreateSubscription\`
- [ ] **Service:** \`subscription-service.ts\`
- [ ] **Service:** \`plan-service.ts\`

### 10.3 Banco de Dados

\`\`\`sql
-- Tabela: plans
CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    valor DECIMAL(10,2) NOT NULL CHECK (valor > 0),
    periodicidade VARCHAR(20) NOT NULL DEFAULT 'MENSAL',
    qtd_servicos INTEGER,
    limite_uso_mensal INTEGER,
    ativo BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, nome)
);

-- Tabela: subscriptions
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    cliente_id UUID NOT NULL REFERENCES clientes(id),
    plano_id UUID NOT NULL REFERENCES plans(id),
    asaas_customer_id VARCHAR(100),
    asaas_subscription_id VARCHAR(100),
    forma_pagamento VARCHAR(20) NOT NULL CHECK (forma_pagamento IN ('CARTAO', 'PIX', 'DINHEIRO')),
    status VARCHAR(30) NOT NULL DEFAULT 'AGUARDANDO_PAGAMENTO',
    valor DECIMAL(10,2) NOT NULL,
    link_pagamento TEXT,
    codigo_transacao VARCHAR(100),
    data_ativacao TIMESTAMPTZ,
    data_vencimento TIMESTAMPTZ,
    data_cancelamento TIMESTAMPTZ,
    cancelado_por UUID REFERENCES users(id),
    servicos_utilizados INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT valid_status CHECK (status IN (
        'AGUARDANDO_PAGAMENTO', 'ATIVO', 'INADIMPLENTE', 'INATIVO', 'CANCELADO'
    ))
);

-- Índices
CREATE INDEX idx_subscriptions_tenant ON subscriptions(tenant_id);
CREATE INDEX idx_subscriptions_cliente ON subscriptions(cliente_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_vencimento ON subscriptions(data_vencimento);
CREATE INDEX idx_subscriptions_asaas ON subscriptions(asaas_subscription_id);

-- Histórico de pagamentos
CREATE TABLE subscription_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID NOT NULL REFERENCES subscriptions(id),
    asaas_payment_id VARCHAR(100),
    valor DECIMAL(10,2) NOT NULL,
    forma_pagamento VARCHAR(20) NOT NULL,
    status VARCHAR(30) NOT NULL,
    data_pagamento TIMESTAMPTZ,
    codigo_transacao VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
\`\`\`

---

## 📋 Histórico de Alterações

| Versão | Data | Autor | Alterações |
|--------|------|-------|------------|
| 1.0 | 22/11/2025 | Equipe | Versão inicial simplificada |
| 2.0 | 27/11/2025 | Equipe | Documento completo: 3 formas de pagamento, webhooks, regras de negócio, checklist implementação |

---

**FIM DO DOCUMENTO**
