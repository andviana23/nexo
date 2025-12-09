# Fluxo de Tipos de Recebimento — NEXO v1.0

**Versão:** 1.0
**Data de Criação:** 27/11/2025
**Status:** 🟢 Implementado
**Responsável:** Product + Tech Lead

---

## 📊 Status de Implementação

| Área | Status | Progresso |
|------|--------|-----------|
| Banco de Dados | ✅ Completo | 100% |
| Backend (Go) | ✅ Completo | 100% |
| Frontend (Next.js) | ✅ Completo | 100% |

---

## 📋 Visão Geral

Módulo responsável pelo **cadastro e gestão de meios de pagamento** (formas de recebimento) da barbearia. Permite configurar:

- **Tipo de Pagamento**: Dinheiro, PIX, Crédito, Débito, Transferência
- **Bandeira**: Visa, Master, Elo, Amex, etc.
- **Taxa percentual**: Desconto cobrado pela operadora (%)
- **Taxa fixa**: Valor fixo por transação (R$)
- **D+**: Dias para compensação bancária (ex: D+1, D+2, D+30)

**Localização no Menu:** Cadastros → Tipos de Recebimento

**Prioridade:** 🟡 MÉDIA (Requisito para Comanda/Pagamento)

---

## 🎯 Objetivos do Fluxo

1. ✅ Cadastrar meios de pagamento aceitos pela barbearia
2. ✅ Definir taxas por tipo de pagamento
3. ✅ Configurar prazo de compensação (D+)
4. ✅ Ativar/desativar meios de pagamento
5. ✅ Ordenar exibição na comanda
6. ✅ Respeitar isolamento multi-tenant

---

## 🔐 Regras de Negócio (RN)

### RN-REC-001: Tipos de Pagamento

Tipos permitidos (enum):
- `DINHEIRO` - Pagamento em espécie (D+0)
- `PIX` - Pagamento instantâneo (D+0 ou D+1)
- `CREDITO` - Cartão de crédito (D+30 padrão)
- `DEBITO` - Cartão de débito (D+1 padrão)
- `TRANSFERENCIA` - TED/DOC bancária (D+0 ou D+1)

### RN-REC-002: Bandeiras de Cartão

Para tipos `CREDITO` e `DEBITO`, permitir cadastrar bandeira:
- Visa
- Mastercard
- Elo
- Amex
- Hipercard
- Outros (campo livre)

### RN-REC-003: Cálculo de D+ (Dias para Compensação)

O sistema calcula a **data de compensação** considerando:

1. **D+ configurado**: Ex: D+1, D+2, D+30
2. **Apenas dias úteis**: Pula sábados e domingos
3. **Pula feriados**: (Futuro - tabela de feriados por tenant)

**Exemplo:**
- Venda na sexta-feira com D+1 → Compensação na segunda-feira
- Venda na quinta-feira com D+2 → Compensação na segunda-feira (pula sáb/dom)

### RN-REC-004: Taxa de Pagamento

- **Taxa percentual**: 0% a 100% (ex: 2.49% para crédito)
- **Taxa fixa**: R$ 0,00+ (ex: R$ 0,50 por transação PIX)
- **Cálculo do líquido**:
  ```
  valor_liquido = valor_bruto - (valor_bruto × taxa_percentual / 100) - taxa_fixa
  ```

### RN-REC-005: Validações

- ✅ Nome obrigatório
- ✅ Tipo obrigatório (enum válido)
- ✅ Taxa entre 0% e 100%
- ✅ Taxa fixa >= R$ 0,00
- ✅ D+ >= 0 dias
- ✅ Único por tenant (nome + tipo + bandeira)

### RN-REC-006: Multi-tenant

- Cada tenant tem seus próprios meios de pagamento
- Não compartilha configurações entre tenants
- `tenant_id` obrigatório em todas as queries

---

## 📦 Modelo de Dados

### Tabela: `meios_pagamento`

```sql
CREATE TABLE meios_pagamento (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    nome            VARCHAR(100) NOT NULL,        -- Ex: "Visa Crédito"
    tipo            VARCHAR(30) NOT NULL,         -- DINHEIRO, PIX, CREDITO, DEBITO, TRANSFERENCIA
    bandeira        VARCHAR(50),                  -- Visa, Master, Elo (opcional)
    taxa            NUMERIC(5,2) DEFAULT 0.00,    -- Taxa % (0-100)
    taxa_fixa       NUMERIC(10,2) DEFAULT 0.00,   -- Taxa fixa R$
    d_mais          INTEGER DEFAULT 0,            -- Dias para compensação
    icone           VARCHAR(50),                  -- Ícone Material Icons
    cor             VARCHAR(7),                   -- Cor hexadecimal
    ordem_exibicao  INTEGER DEFAULT 0,            -- Ordem na UI
    observacoes     TEXT,
    ativo           BOOLEAN DEFAULT true,
    criado_em       TIMESTAMPTZ DEFAULT now(),
    atualizado_em   TIMESTAMPTZ DEFAULT now(),
    
    CONSTRAINT chk_taxa_valida CHECK (taxa >= 0 AND taxa <= 100),
    CONSTRAINT chk_taxa_fixa_valida CHECK (taxa_fixa >= 0),
    CONSTRAINT chk_tipo_valido CHECK (tipo IN ('DINHEIRO', 'PIX', 'CREDITO', 'DEBITO', 'TRANSFERENCIA')),
    CONSTRAINT chk_d_mais_valido CHECK (d_mais >= 0)
);
```

---

## 🔄 Fluxo de Telas

### Tela 1: Lista de Tipos de Recebimento

**Rota:** `/cadastros/tipos-recebimento`

**Layout:**
```
┌─────────────────────────────────────────────────────────────┐
│  Tipos de Recebimento                    [+ Novo Tipo]      │
├─────────────────────────────────────────────────────────────┤
│  🔍 Buscar...                             [Todos ▼]         │
├─────────────────────────────────────────────────────────────┤
│  ┌────────────────────────────────────────────────────────┐ │
│  │ 💵 Dinheiro           | D+0  | 0%   | Ativo   | ⋮     │ │
│  ├────────────────────────────────────────────────────────┤ │
│  │ 📱 PIX                | D+0  | 0%   | Ativo   | ⋮     │ │
│  ├────────────────────────────────────────────────────────┤ │
│  │ 💳 Visa Crédito       | D+30 | 2.49%| Ativo   | ⋮     │ │
│  ├────────────────────────────────────────────────────────┤ │
│  │ 💳 Master Crédito     | D+30 | 2.49%| Ativo   | ⋮     │ │
│  ├────────────────────────────────────────────────────────┤ │
│  │ 💳 Visa Débito        | D+1  | 1.49%| Ativo   | ⋮     │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

**Ações do Menu (⋮):**
- Editar
- Duplicar
- Desativar/Ativar
- Excluir

### Tela 2: Modal de Cadastro/Edição

**Campos do Formulário:**

| Campo | Tipo | Obrigatório | Validação |
|-------|------|-------------|-----------|
| Nome | Input | Sim | Max 100 chars |
| Tipo | Select | Sim | Enum PaymentType |
| Bandeira | Select/Input | Não | Apenas para CREDITO/DEBITO |
| Taxa (%) | Number | Não | 0-100, 2 decimais |
| Taxa Fixa (R$) | Currency | Não | >= 0 |
| Dias para Recebimento (D+) | Number | Não | >= 0 |
| Ícone | IconPicker | Não | Material Icons |
| Cor | ColorPicker | Não | Hexadecimal |
| Observações | Textarea | Não | Max 500 chars |
| Ativo | Switch | Não | Default: true |

**Layout do Modal:**
```
┌─────────────────────────────────────────────────────────────┐
│  Novo Tipo de Recebimento                              [X]  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Nome *                                                     │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Visa Crédito                                         │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌──────────────────────┐  ┌──────────────────────────┐    │
│  │ Tipo *          ▼    │  │ Bandeira            ▼   │    │
│  │ Crédito              │  │ Visa                    │    │
│  └──────────────────────┘  └──────────────────────────┘    │
│                                                             │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────┐    │
│  │ Taxa (%)       │  │ Taxa Fixa (R$) │  │ D+         │    │
│  │ 2.49           │  │ 0.00           │  │ 30         │    │
│  └────────────────┘  └────────────────┘  └────────────┘    │
│                                                             │
│  💡 Com D+30, um pagamento feito hoje será compensado      │
│     em 30 dias úteis (pulando finais de semana).           │
│                                                             │
│  Observações                                                │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                                                      │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  [  ] Ativo                                                 │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│                              [Cancelar]  [Salvar]           │
└─────────────────────────────────────────────────────────────┘
```

---

## 📡 API Endpoints

### Base URL: `/api/v1/payment-methods`

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/` | Listar meios de pagamento |
| GET | `/:id` | Buscar por ID |
| POST | `/` | Criar novo |
| PATCH | `/:id` | Atualizar existente |
| DELETE | `/:id` | Excluir (soft delete) |
| PATCH | `/:id/toggle` | Ativar/Desativar |

### DTOs

**Request (Create/Update):**
```json
{
  "nome": "Visa Crédito",
  "tipo": "CREDITO",
  "bandeira": "Visa",
  "taxa": "2.49",
  "taxa_fixa": "0.00",
  "d_mais": 30,
  "icone": "credit_card",
  "cor": "#1A73E8",
  "ordem_exibicao": 1,
  "observacoes": "Cartão de crédito Visa",
  "ativo": true
}
```

**Response:**
```json
{
  "id": "uuid",
  "tenant_id": "uuid",
  "nome": "Visa Crédito",
  "tipo": "CREDITO",
  "bandeira": "Visa",
  "taxa": "2.49",
  "taxa_fixa": "0.00",
  "d_mais": 30,
  "icone": "credit_card",
  "cor": "#1A73E8",
  "ordem_exibicao": 1,
  "observacoes": "Cartão de crédito Visa",
  "ativo": true,
  "criado_em": "2025-11-27T10:00:00Z",
  "atualizado_em": "2025-11-27T10:00:00Z"
}
```

---

## 📊 Função: Calcular Data de Compensação

### Lógica de D+ (Dias Úteis)

```go
// CalculateSettlementDate calcula a data de compensação baseada em D+
// Considera apenas dias úteis (pula sábados e domingos)
func CalculateSettlementDate(transactionDate time.Time, dPlus int) time.Time {
    if dPlus == 0 {
        return transactionDate
    }
    
    result := transactionDate
    daysAdded := 0
    
    for daysAdded < dPlus {
        result = result.AddDate(0, 0, 1)
        weekday := result.Weekday()
        
        // Pula sábado e domingo
        if weekday != time.Saturday && weekday != time.Sunday {
            daysAdded++
        }
    }
    
    return result
}
```

### Exemplos:

| Data Transação | D+ | Data Compensação |
|----------------|-----|------------------|
| Seg 25/11/2025 | D+1 | Ter 26/11/2025 |
| Sex 28/11/2025 | D+1 | Seg 01/12/2025 |
| Qui 27/11/2025 | D+2 | Seg 01/12/2025 |
| Seg 25/11/2025 | D+30| Ter 07/01/2026 |

---

## 🧪 Critérios de Aceite

### CA-001: Listagem
- [x] Exibe todos os meios de pagamento do tenant
- [x] Permite filtrar por tipo
- [x] Permite buscar por nome
- [x] Ordena por `ordem_exibicao`
- [x] Mostra badge de ativo/inativo

### CA-002: Criação
- [x] Valida campos obrigatórios
- [x] Salva com taxa padrão 0%
- [x] Salva com D+0 padrão
- [x] Mostra bandeira apenas para CREDITO/DEBITO
- [x] Toast de sucesso

### CA-003: Edição
- [x] Carrega dados existentes
- [x] Atualiza campos alterados
- [x] Atualiza `updated_at`
- [x] Toast de sucesso

### CA-004: Exclusão
- [x] Confirmação antes de excluir
- [x] Soft delete (desativa)
- [x] Toast de sucesso

### CA-005: Toggle Ativo
- [x] Alterna status com um clique
- [x] Atualiza UI imediatamente
- [x] Toast de confirmação

---

## 🔗 Dependências

### Upstream (Este módulo depende de):
- `tenants` - Isolamento multi-tenant

### Downstream (Módulos que dependem deste):
- **Comanda** - Seleção de forma de pagamento
- **Compensações Bancárias** - Cálculo de D+
- **Fluxo de Caixa** - Previsão de recebimentos
- **Relatórios** - Análise por forma de pagamento

---

## 📱 Responsividade

### Desktop (>1024px)
- Lista em tabela com todas as colunas
- Modal centralizado 480px

### Tablet (768-1024px)
- Lista em cards compactos
- Modal full-width

### Mobile (<768px)
- Lista em cards empilhados
- Modal full-screen

---

## 🚀 Seeds de Teste

```sql
-- Meios de pagamento padrão para tenant E2E
INSERT INTO meios_pagamento (tenant_id, nome, tipo, bandeira, taxa, taxa_fixa, d_mais, icone, ordem_exibicao) VALUES
('TENANT_E2E', 'Dinheiro', 'DINHEIRO', NULL, 0, 0, 0, 'payments', 1),
('TENANT_E2E', 'PIX', 'PIX', NULL, 0, 0, 0, 'qr_code', 2),
('TENANT_E2E', 'Visa Crédito', 'CREDITO', 'Visa', 2.49, 0, 30, 'credit_card', 3),
('TENANT_E2E', 'Master Crédito', 'CREDITO', 'Mastercard', 2.49, 0, 30, 'credit_card', 4),
('TENANT_E2E', 'Elo Crédito', 'CREDITO', 'Elo', 2.99, 0, 30, 'credit_card', 5),
('TENANT_E2E', 'Visa Débito', 'DEBITO', 'Visa', 1.49, 0, 1, 'credit_card', 6),
('TENANT_E2E', 'Master Débito', 'DEBITO', 'Mastercard', 1.49, 0, 1, 'credit_card', 7);
```

---

## 📜 Histórico de Alterações

| Data | Autor | Alteração |
|------|-------|-----------|
| 27/11/2025 | Copilot | Criação do documento |

---

**Gerente de Projeto:** Andrey  
**Tech Lead:** Copilot  
**Data de Início:** 27/11/2025  
**Última Atualização:** 27/11/2025
