# Fluxo de Estoque — NEXO v2.0

**Versão:** 2.0
**Última Atualização:** 27/11/2025
**Status:** 🟡 Planejado (v1.2.0 - Completo)
**Responsável:** Product + Tech Lead

---

## 📋 Visão Geral

Módulo responsável pelo **gerenciamento completo de estoque 360º**, cobrindo desde a requisição de compra, controle de validade, consumo por barbeiro, até auditorias semanais e prevenção de perdas.

**Diferencial v2.0:**
- **Centro de Custo:** Separação clara para DRE (Insumos vs Limpeza vs Revenda).
- **Validade:** Controle de lotes e alertas de vencimento.
- **Auditoria:** Processo formal de contagem e ajuste.
- **Inteligência:** Previsão de compra e análise de desperdício por barbeiro.
- **Compras:** Fluxo de aprovação e integração financeira.

---

## 🎯 Objetivos do Fluxo

1.  ✅ Cadastro completo (Produtos, Fornecedores, Centros de Custo).
2.  ✅ Controle de Lotes e Validade (FIFO).
3.  ✅ Auditoria Semanal (Checklist + Ajustes).
4.  ✅ Gestão de Compras (Requisição -> Aprovação -> Entrada).
5.  ✅ Rastreabilidade por Barbeiro (Consumo vs Padrão).
6.  ✅ Relatórios de Ruptura (Real vs Técnica).
7.  ✅ Curva ABC e Previsão de Reposição.
8.  ✅ Integração total com Financeiro (DRE).

---

## 🔐 Regras de Negócio (RN)

### RN-EST-001: Cadastro e Categorização
- **Categorias:** `POMADA`, `SHAMPOO`, `CREME`, `LAMINA`, `TOALHA`, `LIMPEZA`, `ESCRITORIO`, `BEBIDA`, `REVENDA`.
- **Centros de Custo:**
    - `CUSTO_SERVICO` (Insumos, Lâminas, Shampoos)
    - `DESPESA_OPERACIONAL` (Limpeza, Escritório)
    - `CUSTO_MERCADORIA_VENDIDA` (Revenda, Bebidas)
- **SKU:** Único por tenant.

### RN-EST-002: Controle de Validade (Lotes)
- Produtos perecíveis (`controla_validade = true`) exigem data de validade na ENTRADA.
- Sistema adota **FEFO** (First Expire, First Out) ou **FIFO** para baixa automática.
- Alertas:
    - 🟡 Vence em 30 dias.
    - 🔴 Vence em 7 dias.
    - ⚫ Vencido (Bloqueio de uso/venda).

### RN-EST-003: Auditoria e Rupturas
- **Ruptura Técnica:** Sistema indica saldo 0.
- **Ruptura Real:** Barbeiro sinaliza falta no app, mesmo com saldo > 0 (indica furto/perda não registrada).
- **Auditoria Semanal:**
    - Gerente recebe checklist dos itens Curva A e B.
    - Contagem cega (sistema não mostra saldo esperado).
    - Divergência > X% exige justificativa.

### RN-EST-004: Compras e Reposição
- **Ponto de Pedido:** `Estoque Mínimo + (Consumo Médio Diário * Lead Time)`.
- **Sugestão de Compra:** Automática baseada no consumo dos últimos 30/90 dias.
- **Fluxo:** Requisição -> Aprovação (Dono) -> Pedido -> Entrada XML/Manual.

### RN-EST-005: Consumo por Barbeiro
- Vínculo de consumo na baixa de serviço (Ficha Técnica).
- Registro de "Retirada de Insumo" pelo barbeiro (ex: pegou um tubo de pomada novo).
- Relatório comparativo: `Consumo Real vs Consumo Padrão (Ficha Técnica)`.

---

## 📊 Diagrama de Fluxo (Mermaid)

```mermaid
flowchart TD
    subgraph COMPRAS [Gestão de Compras]
        A[Sugestão de Compra] --> B[Criar Requisição]
        B --> C{Aprovação?}
        C -->|Sim| D[Pedido ao Fornecedor]
        C -->|Não| E[Arquivar]
        D --> F[Recebimento/Entrada]
    end

    subgraph ESTOQUE [Controle Diário]
        F --> G[Entrada (Lote/Validade)]
        G --> H[Estoque Disponível]
        
        H --> I{Tipo Saída?}
        I -->|Venda| J[Baixa Estoque (Revenda)]
        I -->|Serviço| K[Baixa Automática (Ficha Técnica)]
        I -->|Consumo| L[Retirada Barbeiro]
        I -->|Perda/Venc| M[Baixa por Perda]
    end

    subgraph AUDITORIA [Controle e Ajuste]
        N[Auditoria Semanal] --> O[Contagem Cega]
        O --> P{Divergência?}
        P -->|Sim| Q[Ajuste de Estoque (Perda/Sobra)]
        P -->|Não| R[Validado]
        Q --> S[Relatório de Perdas]
    end

    subgraph INTELIGENCIA [Relatórios]
        K --> T[Consumo por Barbeiro]
        L --> T
        H --> U[Previsão de Ruptura]
        G --> V[Contas a Pagar (Financeiro)]
        M --> W[Custo de Desperdício]
    end
```

---

## 🏗️ Arquitetura e Entidades (Atualizado)

### 1. Entity: Produto (Atualizado)

```go
type CentroCusto string

const (
    CentroCustoServico     CentroCusto = "CUSTO_SERVICO"      // Insumos diretos
    CentroCustoOperacional CentroCusto = "DESPESA_OPERACIONAL" // Limpeza, escritório
    CentroCustoCMV         CentroCusto = "CMV"                // Revenda
)

type Produto struct {
    // ... campos existentes ...
    CentroCusto      CentroCusto
    ControlaValidade bool
    LeadTimeDias     int // Tempo médio de reposição
}
```

### 2. Entity: Lote (Novo)

```go
type Lote struct {
    ID             uuid.UUID
    ProdutoID      uuid.UUID
    CodigoLote     string
    DataFabricacao *time.Time
    DataValidade   time.Time
    Quantidade     int
    Ativo          bool // false se vencido ou zerado
}
```

### 3. Entity: Auditoria (Novo)

```go
type Auditoria struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    Responsavel uuid.UUID
    DataInicio  time.Time
    DataFim     *time.Time
    Status      string // ABERTA, FINALIZADA
    Itens       []ItemAuditoria
}

type ItemAuditoria struct {
    ProdutoID        uuid.UUID
    QuantidadeSistema int
    QuantidadeContada int
    Divergencia       int
    Justificativa     string
}
```

### 4. Entity: RequisicaoCompra (Novo)

```go
type RequisicaoCompra struct {
    ID            uuid.UUID
    SolicitanteID uuid.UUID
    Status        string // PENDENTE, APROVADA, REJEITADA, COMPRADA
    Itens         []ItemRequisicao
    ValorEstimado valueobject.Money
}
```

---

## 🚀 Novos Módulos Detalhados

### 1. Gestão de Validade e Lotes
- **Entrada:** Ao registrar entrada, se `produto.controla_validade == true`, exigir Data de Validade. Cria-se um registro na tabela `lotes`.
- **Saída:** O sistema baixa automaticamente do lote com validade mais próxima (FEFO).
- **Cron Job:** Diariamente verifica lotes vencidos -> Marca como `VENCIDO` -> Notifica gerente.

### 2. Auditoria Semanal
- **Checklist:** O sistema gera lista de produtos para contagem (foco em Curva A e B).
- **App Mobile:** Gerente escaneia ou digita a quantidade encontrada.
- **Confronto:** Sistema compara `Qtd Contada` vs `Qtd Sistema`.
- **Ajuste:** Se houver diferença, gera movimentação de `AJUSTE_AUDITORIA` automaticamente ao finalizar.

### 3. Ruptura Real vs Técnica
- **Botão de Pânico (App Barbeiro):** "Informar falta de insumo".
- Se Barbeiro informa falta, mas Sistema diz `Qtd > 0` -> **Ruptura Real** (Erro de estoque/Furto).
- Se Sistema diz `Qtd = 0` -> **Ruptura Técnica** (Falha de reposição).

### 4. Previsão de Reposição
- **Cálculo:**
  `Consumo Médio Diário (CMD) = Consumo últimos 30 dias / 30`
  `Estoque de Segurança (ES) = CMD * Dias de Segurança (ex: 5)`
  `Ponto de Pedido = (CMD * Lead Time) + ES`
- **Sugestão:** Se `Estoque Atual <= Ponto de Pedido`, sugerir compra.

### 5. Histórico por Barbeiro
- Cada baixa de insumo (automática por serviço ou manual por retirada) é vinculada ao `barbeiro_id`.
- **KPIs:**
  - Custo de Insumo por Serviço (R$)
  - Desvio Padrão (Quem gasta muito mais que a média?)
  - Índice de Desperdício (Retiradas manuais sem serviço vinculado).

---

## 📊 Modelo de Dados (SQL Atualizado)

```sql
-- Atualização Tabela Produtos
ALTER TABLE produtos ADD COLUMN centro_custo VARCHAR(50);
ALTER TABLE produtos ADD COLUMN controla_validade BOOLEAN DEFAULT false;
ALTER TABLE produtos ADD COLUMN lead_time_dias INT DEFAULT 7;

-- Tabela Lotes
CREATE TABLE lotes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    produto_id UUID NOT NULL REFERENCES produtos(id),
    codigo_lote VARCHAR(50),
    data_validade DATE NOT NULL,
    quantidade_inicial INT NOT NULL,
    quantidade_atual INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_lotes_validade ON lotes(data_validade);

-- Tabela Auditorias
CREATE TABLE auditorias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    responsavel_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'ABERTA',
    data_inicio TIMESTAMP DEFAULT NOW(),
    data_fim TIMESTAMP
);

-- Itens Auditoria
CREATE TABLE itens_auditoria (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auditoria_id UUID NOT NULL REFERENCES auditorias(id),
    produto_id UUID NOT NULL REFERENCES produtos(id),
    qtd_sistema INT NOT NULL,
    qtd_contada INT NOT NULL,
    divergencia INT GENERATED ALWAYS AS (qtd_contada - qtd_sistema) STORED,
    justificativa TEXT
);

-- Requisições de Compra
CREATE TABLE requisicoes_compra (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    solicitante_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'PENDENTE', -- PENDENTE, APROVADA, COMPRADA, CANCELADA
    observacoes TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE itens_requisicao (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requisicao_id UUID NOT NULL REFERENCES requisicoes_compra(id),
    produto_id UUID NOT NULL REFERENCES produtos(id),
    qtd_sugerida INT,
    qtd_aprovada INT
);
```

---

## 🌐 Novos Endpoints API

### Auditoria
- `POST /api/v1/estoque/auditorias/iniciar`
- `POST /api/v1/estoque/auditorias/{id}/contagem` (Lança item contado)
- `POST /api/v1/estoque/auditorias/{id}/finalizar` (Processa ajustes)

### Compras
- `GET /api/v1/estoque/sugestao-compra` (Algoritmo de previsão)
- `POST /api/v1/estoque/requisicoes`
- `PATCH /api/v1/estoque/requisicoes/{id}/aprovar`

### Relatórios
- `GET /api/v1/estoque/relatorios/validade` (Itens vencendo)
- `GET /api/v1/estoque/relatorios/consumo-barbeiro`
- `GET /api/v1/estoque/relatorios/rupturas`

---

**Status:** 🟡 Planejado (v1.2.0)
**Prioridade:** Média/Alta (Gestão Eficiente)
