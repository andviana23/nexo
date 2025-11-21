# CATEGORIA: FINANCEIRO > FLUXO DE CAIXA COMPENSADO

## Plano de Execução (terceira prioridade)
- **Banco de Dados:** adicionar `d_mais` em `meios_pagamento`; criar `fluxo_caixa_diario` e `compensacoes_bancarias` (status PREVISTO/CONFIRMADO/COMPENSADO, datas de transação/compensação, valor bruto/líquido).
- **Backend:** cálculo D+N por meio de pagamento, geração de compensações e snapshot diário; integrar contas a pagar/receber para previsão; status updates; endpoint diário/intervalo.
- **Frontend:** visão diária (saldo inicial/entradas/saídas/saldo final), linha do tempo de compensações, filtros por período/meio/status.
- **Cálculos aplicados:** “Previsão de Fluxo de Caixa Compensado” (`docs/10-calculos/previsao-fluxo-caixa-compensado.md`). Cruzar com metas mínimas/PE no dashboard (`faturamento-minimo-mensal.md`, `ponto-de-equilibrio.md`) para alertas.

## Análise do Sistema Atual

### Estado Implementado

- ✅ Tabelas: `receitas`, `despesas`, `meios_pagamento` (com `taxa` e `taxa_fixa`)
- ✅ Use Case: `GetFluxoDeCaixaUseCase` (básico: entradas - saídas)
- ✅ Endpoint: `GET /financial/cashflow` (saldo simples)
- ✅ Frontend: hook `useCashflow`
- ⚠️ **Pendente**: Compensação bancária, previsão futura, D+

### Gap Identificado

1. **Tabela `meios_pagamento` não tem campo `dias_compensacao`** → adicionar campo `d_mais`
2. **Sem tabela de `fluxo_caixa_compensado`** → criar snapshot diário com status
3. **Sem lógica de compensação** → implementar cálculo D+N
4. **Sem previsão futura** → incluir recebíveis (contas a receber) e pagáveis (contas a pagar)

---

## Funcionalidade: Fluxo de Caixa Compensado com D+

### Objetivo

Exibir fluxo de caixa considerando:

- **Compensação bancária**: cartão crédito (D+30), débito (D+1), PIX/dinheiro (D+0)
- **Previsão futura**: contas a receber e contas a pagar
- **Visão diária**: saldo inicial, entradas esperadas, saídas, saldo final
- **Status**: `PREVISTO`, `CONFIRMADO`, `COMPENSADO`

---

## BACKEND

### Tarefas Backend

#### 1. Modelagem de Banco de Dados

**Alterar Tabela: `meios_pagamento`**

```sql
-- Migration: 029_alter_meios_pagamento_add_d_mais.up.sql
ALTER TABLE meios_pagamento
ADD COLUMN d_mais INTEGER DEFAULT 0 CHECK (d_mais >= 0);

COMMENT ON COLUMN meios_pagamento.d_mais IS 'Dias para compensação bancária (D+0 para PIX/Dinheiro, D+1 para Débito, D+30 para Crédito)';

-- Atualizar registros existentes (exemplos)
UPDATE meios_pagamento SET d_mais = 0 WHERE tipo IN ('PIX', 'DINHEIRO');
UPDATE meios_pagamento SET d_mais = 1 WHERE tipo = 'DEBITO';
UPDATE meios_pagamento SET d_mais = 30 WHERE tipo = 'CREDITO';
UPDATE meios_pagamento SET d_mais = 1 WHERE tipo = 'TRANSFERENCIA';
```

**Nova Tabela: `fluxo_caixa_diario`**

```sql
-- Migration: 030_create_fluxo_caixa_diario.up.sql
CREATE TABLE fluxo_caixa_diario (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    data DATE NOT NULL,

    -- Saldos
    saldo_inicial DECIMAL(15,2) DEFAULT 0,
    saldo_final DECIMAL(15,2) DEFAULT 0,

    -- Entradas
    entradas_confirmadas DECIMAL(15,2) DEFAULT 0, -- Já compensadas
    entradas_previstas DECIMAL(15,2) DEFAULT 0,   -- Aguardando compensação

    -- Saídas
    saidas_pagas DECIMAL(15,2) DEFAULT 0,         -- Despesas pagas
    saidas_previstas DECIMAL(15,2) DEFAULT 0,     -- Contas a pagar do dia

    -- Metadados
    processado_em TIMESTAMPTZ DEFAULT NOW(),
    criado_em TIMESTAMPTZ DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(tenant_id, data)
);

CREATE INDEX idx_fluxo_caixa_diario_tenant ON fluxo_caixa_diario(tenant_id);
CREATE INDEX idx_fluxo_caixa_diario_data ON fluxo_caixa_diario(tenant_id, data DESC);
```

**Nova Tabela: `compensacoes_bancarias`**

```sql
-- Migration: 031_create_compensacoes_bancarias.up.sql
CREATE TABLE compensacoes_bancarias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    receita_id UUID REFERENCES receitas(id) ON DELETE CASCADE,

    -- Datas
    data_transacao DATE NOT NULL,      -- Data da venda
    data_compensacao DATE NOT NULL,    -- Data prevista de compensação
    data_compensado DATE,               -- Data real de compensação (NULL = pendente)

    -- Valores
    valor_bruto DECIMAL(15,2) NOT NULL,
    taxa_percentual DECIMAL(5,2) DEFAULT 0,
    taxa_fixa DECIMAL(10,2) DEFAULT 0,
    valor_liquido DECIMAL(15,2) NOT NULL, -- Bruto - Taxas

    -- Meio de pagamento
    meio_pagamento_id UUID REFERENCES meios_pagamento(id),
    d_mais INTEGER NOT NULL, -- Dias para compensação

    -- Status
    status VARCHAR(20) DEFAULT 'PREVISTO', -- PREVISTO, CONFIRMADO, COMPENSADO, CANCELADO

    -- Metadados
    criado_em TIMESTAMPTZ DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT chk_status_compensacao CHECK (status IN ('PREVISTO', 'CONFIRMADO', 'COMPENSADO', 'CANCELADO'))
);

CREATE INDEX idx_compensacoes_tenant ON compensacoes_bancarias(tenant_id);
CREATE INDEX idx_compensacoes_data_compensacao ON compensacoes_bancarias(tenant_id, data_compensacao);
CREATE INDEX idx_compensacoes_status ON compensacoes_bancarias(status, data_compensacao);
CREATE INDEX idx_compensacoes_receita ON compensacoes_bancarias(receita_id);
```

#### 2. Domain Layer (Go)

**Entidade CompensacaoBancaria**

```go
// internal/domain/compensacao_bancaria.go
package domain

import (
    "time"
    "github.com/shopspring/decimal"
)

type CompensacaoStatus string

const (
    CompensacaoPrevista   CompensacaoStatus = "PREVISTO"
    CompensacaoConfirmada CompensacaoStatus = "CONFIRMADO"
    CompensacaoCompensada CompensacaoStatus = "COMPENSADO"
    CompensacaoCancelada  CompensacaoStatus = "CANCELADO"
)

type CompensacaoBancaria struct {
    ID                 string
    TenantID           string
    ReceitaID          string

    DataTransacao      time.Time
    DataCompensacao    time.Time
    DataCompensado     *time.Time

    ValorBruto         decimal.Decimal
    TaxaPercentual     decimal.Decimal
    TaxaFixa           decimal.Decimal
    ValorLiquido       decimal.Decimal

    MeioPagamentoID    string
    DMais              int

    Status             CompensacaoStatus

    CriadoEm           time.Time
    AtualizadoEm       time.Time
}

func (c *CompensacaoBancaria) CalcularValorLiquido() {
    taxaPerc := c.ValorBruto.Mul(c.TaxaPercentual).Div(decimal.NewFromInt(100))
    c.ValorLiquido = c.ValorBruto.Sub(taxaPerc).Sub(c.TaxaFixa)
}

func (c *CompensacaoBancaria) MarcarComoCompensado() error {
    if c.Status == CompensacaoCompensada {
        return ErrCompensacaoJaCompensada
    }
    now := time.Now()
    c.DataCompensado = &now
    c.Status = CompensacaoCompensada
    c.AtualizadoEm = now
    return nil
}

func (c *CompensacaoBancaria) Validate() error {
    if c.TenantID == "" {
        return ErrTenantIDRequired
    }
    if c.ReceitaID == "" {
        return ErrReceitaIDRequired
    }
    if c.ValorBruto.LessThanOrEqual(decimal.Zero) {
        return ErrValorInvalido
    }
    if c.DMais < 0 {
        return ErrDMaisInvalido
    }
    return nil
}
```

**Entidade FluxoCaixaDiario**

```go
// internal/domain/fluxo_caixa_diario.go
type FluxoCaixaDiario struct {
    ID                   string
    TenantID             string
    Data                 time.Time

    SaldoInicial         decimal.Decimal
    EntradasConfirmadas  decimal.Decimal
    EntradasPrevistas    decimal.Decimal
    SaidasPagas          decimal.Decimal
    SaidasPrevistas      decimal.Decimal
    SaldoFinal           decimal.Decimal

    ProcessadoEm         time.Time
    CriadoEm             time.Time
    AtualizadoEm         time.Time
}

func (f *FluxoCaixaDiario) Calcular() {
    f.SaldoFinal = f.SaldoInicial.
        Add(f.EntradasConfirmadas).
        Add(f.EntradasPrevistas).
        Sub(f.SaidasPagas).
        Sub(f.SaidasPrevistas)
}
```

#### 3. Application Layer (Use Cases)

**CreateCompensacaoUseCase**

```go
// Chamado ao criar receita
type CreateCompensacaoUseCase struct {
    compensacaoRepo CompensacaoBancariaRepository
    meioPagamentoRepo MeioPagamentoRepository
}

func (uc *CreateCompensacaoUseCase) Execute(ctx context.Context, receita *Receita) error {
    // 1. Buscar meio de pagamento para obter D+ e taxas
    meioPagamento, err := uc.meioPagamentoRepo.FindByID(ctx, receita.TenantID, receita.MeioPagamentoID)
    if err != nil {
        return err
    }

    // 2. Calcular data de compensação
    dataCompensacao := receita.Data.AddDate(0, 0, meioPagamento.DMais)

    // 3. Criar compensação
    compensacao := &CompensacaoBancaria{
        ID:              uuid.NewString(),
        TenantID:        receita.TenantID,
        ReceitaID:       receita.ID,
        DataTransacao:   receita.Data,
        DataCompensacao: dataCompensacao,
        ValorBruto:      receita.Valor,
        TaxaPercentual:  meioPagamento.Taxa,
        TaxaFixa:        meioPagamento.TaxaFixa,
        MeioPagamentoID: meioPagamento.ID,
        DMais:           meioPagamento.DMais,
        Status:          CompensacaoPrevista,
        CriadoEm:        time.Now(),
    }

    compensacao.CalcularValorLiquido()

    // 4. Salvar
    return uc.compensacaoRepo.Save(ctx, receita.TenantID, compensacao)
}
```

**GenerateFluxoCaixaDiarioUseCase**

```go
// Gerado via Cron diário
type GenerateFluxoCaixaDiarioUseCase struct {
    fluxoRepo       FluxoCaixaDiarioRepository
    compensacaoRepo CompensacaoBancariaRepository
    despesaRepo     DespesaRepository
    contasPagarRepo ContasAPagarRepository
}

func (uc *GenerateFluxoCaixaDiarioUseCase) Execute(ctx context.Context, tenantID string, data time.Time) error {
    // 1. Buscar saldo inicial (saldo_final do dia anterior)
    diaAnterior := data.AddDate(0, 0, -1)
    fluxoAnterior, _ := uc.fluxoRepo.FindByTenantAndDate(ctx, tenantID, diaAnterior)

    saldoInicial := decimal.Zero
    if fluxoAnterior != nil {
        saldoInicial = fluxoAnterior.SaldoFinal
    }

    // 2. Entradas confirmadas (compensações já compensadas no dia)
    entradasConfirmadas, _ := uc.compensacaoRepo.SumByTenantDateAndStatus(
        ctx, tenantID, data, CompensacaoCompensada,
    )

    // 3. Entradas previstas (compensações previstas para o dia)
    entradasPrevistas, _ := uc.compensacaoRepo.SumByTenantDateAndStatus(
        ctx, tenantID, data, CompensacaoPrevista,
    )

    // 4. Saídas pagas (despesas pagas no dia)
    saidasPagas, _ := uc.despesaRepo.SumByTenantAndDate(ctx, tenantID, data, ExpensePaid)

    // 5. Saídas previstas (contas a pagar com vencimento no dia)
    saidasPrevistas, _ := uc.contasPagarRepo.SumByTenantAndDataVencimento(ctx, tenantID, data)

    // 6. Criar fluxo
    fluxo := &FluxoCaixaDiario{
        ID:                  uuid.NewString(),
        TenantID:            tenantID,
        Data:                data,
        SaldoInicial:        saldoInicial,
        EntradasConfirmadas: entradasConfirmadas,
        EntradasPrevistas:   entradasPrevistas,
        SaidasPagas:         saidasPagas,
        SaidasPrevistas:     saidasPrevistas,
        ProcessadoEm:        time.Now(),
        CriadoEm:            time.Now(),
    }

    fluxo.Calcular()

    // 7. Salvar
    return uc.fluxoRepo.Save(ctx, tenantID, fluxo)
}
```

**GetFluxoCaixaCompensadoUseCase**

```go
// Endpoint para consulta
func (uc *GetFluxoCaixaCompensadoUseCase) Execute(
    ctx context.Context,
    tenantID string,
    dataInicio, dataFim time.Time,
) ([]*FluxoCaixaDiario, error) {
    return uc.fluxoRepo.FindByTenantAndPeriod(ctx, tenantID, dataInicio, dataFim)
}
```

**MarcarCompensacaoComoCompensadoUseCase**

```go
// Job automático ou manual
func (uc *MarcarCompensacaoComoCompensadoUseCase) Execute(ctx context.Context, tenantID string) error {
    // Buscar compensações previstas com data_compensacao <= hoje
    hoje := time.Now()
    compensacoes, _ := uc.compensacaoRepo.FindPendentesByTenantAndDataCompensacao(ctx, tenantID, hoje)

    for _, comp := range compensacoes {
        comp.MarcarComoCompensado()
        uc.compensacaoRepo.Update(ctx, tenantID, comp)
    }

    return nil
}
```

#### 4. Infrastructure Layer (Repository)

**PostgresCompensacaoRepository**

```go
func (r *PostgresCompensacaoRepository) Save(ctx context.Context, tenantID string, comp *CompensacaoBancaria) error {
    query := `
        INSERT INTO compensacoes_bancarias (
            id, tenant_id, receita_id,
            data_transacao, data_compensacao,
            valor_bruto, taxa_percentual, taxa_fixa, valor_liquido,
            meio_pagamento_id, d_mais, status, criado_em
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `
    _, err := r.db.ExecContext(ctx, query,
        comp.ID, comp.TenantID, comp.ReceitaID,
        comp.DataTransacao, comp.DataCompensacao,
        comp.ValorBruto, comp.TaxaPercentual, comp.TaxaFixa, comp.ValorLiquido,
        comp.MeioPagamentoID, comp.DMais, comp.Status, comp.CriadoEm,
    )
    return err
}

func (r *PostgresCompensacaoRepository) SumByTenantDateAndStatus(
    ctx context.Context, tenantID string, data time.Time, status CompensacaoStatus,
) (decimal.Decimal, error) {
    query := `
        SELECT COALESCE(SUM(valor_liquido), 0)
        FROM compensacoes_bancarias
        WHERE tenant_id = $1
          AND data_compensacao = $2
          AND status = $3
    `
    var total decimal.Decimal
    err := r.db.QueryRowContext(ctx, query, tenantID, data, status).Scan(&total)
    return total, err
}
```

#### 5. HTTP Layer (Handlers)

**FluxoCaixaCompensadoHandler**

```go
func (h *FluxoCaixaCompensadoHandler) GetFluxoCompensado(c echo.Context) error {
    // GET /financial/cashflow/compensado?inicio=2024-11-01&fim=2024-11-30
    tenantID := c.Get("tenant_id").(string)

    inicio, _ := time.Parse("2006-01-02", c.QueryParam("inicio"))
    fim, _ := time.Parse("2006-01-02", c.QueryParam("fim"))

    fluxos, err := h.getFluxoUseCase.Execute(c.Request().Context(), tenantID, inicio, fim)
    if err != nil {
        return c.JSON(500, ErrorResponse{Message: err.Error()})
    }

    return c.JSON(200, MapFluxosToResponse(fluxos))
}

func (h *FluxoCaixaCompensadoHandler) MarcarCompensacoes(c echo.Context) error {
    // POST /financial/cashflow/marcar-compensacoes
    tenantID := c.Get("tenant_id").(string)

    err := h.marcarCompensacaoUseCase.Execute(c.Request().Context(), tenantID)
    if err != nil {
        return c.JSON(500, ErrorResponse{Message: err.Error()})
    }

    return c.JSON(200, map[string]string{"message": "Compensações atualizadas"})
}
```

#### 6. Cron Jobs

**Job 1: Gerar Fluxo Diário (06:00)**

```go
func (j *GenerateFluxoDiarioJob) Run() {
    ctx := context.Background()
    tenants, _ := j.tenantRepo.FindActive(ctx)

    hoje := time.Now()

    for _, tenant := range tenants {
        err := j.generateFluxoUseCase.Execute(ctx, tenant.ID, hoje)
        if err != nil {
            log.Error("Failed to generate daily cashflow",
                zap.String("tenant_id", tenant.ID),
                zap.Error(err))
        }
    }
}
```

**Job 2: Marcar Compensações (07:00)**

```go
func (j *MarcarCompensacoesJob) Run() {
    ctx := context.Background()
    tenants, _ := j.tenantRepo.FindActive(ctx)

    for _, tenant := range tenants {
        err := j.marcarCompensacaoUseCase.Execute(ctx, tenant.ID)
        if err != nil {
            log.Error("Failed to mark compensations",
                zap.String("tenant_id", tenant.ID),
                zap.Error(err))
        }
    }
}
```

---

## FRONTEND

### Tarefas Frontend

#### 1. Hooks Customizados

**useFluxoCaixaCompensado**

```typescript
// app/hooks/useFluxoCaixaCompensado.ts
export function useFluxoCaixaCompensado(inicio: string, fim: string) {
  return useQuery({
    queryKey: ["fluxo-compensado", inicio, fim],
    queryFn: async () => {
      const res = await api.get("/financial/cashflow/compensado", {
        params: { inicio, fim },
      });
      return res.data;
    },
    staleTime: 1000 * 60 * 5, // 5 min
  });
}

export function useMarcarCompensacoes() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      const res = await api.post("/financial/cashflow/marcar-compensacoes");
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["fluxo-compensado"] });
      toast.success("Compensações atualizadas");
    },
  });
}
```

#### 2. Componentes UI

**FluxoCaixaCompensadoTable**

```tsx
// app/components/financeiro/FluxoCaixaCompensadoTable.tsx
export function FluxoCaixaCompensadoTable({
  fluxos,
}: {
  fluxos: FluxoDiario[];
}) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Data</TableHead>
          <TableHead>Saldo Inicial</TableHead>
          <TableHead>Entradas Confirmadas</TableHead>
          <TableHead>Entradas Previstas</TableHead>
          <TableHead>Saídas Pagas</TableHead>
          <TableHead>Saídas Previstas</TableHead>
          <TableHead>Saldo Final</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {fluxos.map((fluxo) => (
          <TableRow key={fluxo.data}>
            <TableCell>{formatDate(fluxo.data)}</TableCell>
            <TableCell>{formatCurrency(fluxo.saldo_inicial)}</TableCell>
            <TableCell className="text-green-600">
              {formatCurrency(fluxo.entradas_confirmadas)}
            </TableCell>
            <TableCell className="text-green-400">
              {formatCurrency(fluxo.entradas_previstas)}
            </TableCell>
            <TableCell className="text-red-600">
              {formatCurrency(fluxo.saidas_pagas)}
            </TableCell>
            <TableCell className="text-red-400">
              {formatCurrency(fluxo.saidas_previstas)}
            </TableCell>
            <TableCell
              className={
                fluxo.saldo_final > 0
                  ? "text-green-700 font-bold"
                  : "text-red-700 font-bold"
              }
            >
              {formatCurrency(fluxo.saldo_final)}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
```

**FluxoChart (Gráfico de Linhas)**

```tsx
// Gráfico mostrando evolução do saldo ao longo do período
import { LineChart } from "recharts";

export function FluxoChart({ fluxos }: Props) {
  const data = fluxos.map((f) => ({
    data: formatDateShort(f.data),
    saldo: f.saldo_final,
    entradas: f.entradas_confirmadas + f.entradas_previstas,
    saidas: f.saidas_pagas + f.saidas_previstas,
  }));

  return (
    <LineChart data={data}>
      <Line dataKey="saldo" stroke="#10b981" strokeWidth={2} />
      <Line dataKey="entradas" stroke="#3b82f6" strokeWidth={1} />
      <Line dataKey="saidas" stroke="#ef4444" strokeWidth={1} />
    </LineChart>
  );
}
```

#### 3. Páginas

**app/financeiro/fluxo-caixa-compensado/page.tsx**

```tsx
"use client";

export default function FluxoCaixaCompensadoPage() {
  const [periodo, setPeriodo] = useState({
    inicio: getStartOfMonth(),
    fim: getEndOfMonth(),
  });

  const { data: fluxos, isLoading } = useFluxoCaixaCompensado(
    periodo.inicio,
    periodo.fim
  );

  const marcarCompensacoes = useMarcarCompensacoes();

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1>Fluxo de Caixa Compensado</h1>
        <div className="flex gap-2">
          <DateRangePicker value={periodo} onChange={setPeriodo} />
          <Button onClick={() => marcarCompensacoes.mutate()}>
            Marcar Compensações
          </Button>
        </div>
      </div>

      {isLoading && <Skeleton />}

      {fluxos && (
        <>
          <FluxoChart fluxos={fluxos} />
          <FluxoCaixaCompensadoTable fluxos={fluxos} />
        </>
      )}

      {/* Legenda */}
      <div className="text-sm text-muted-foreground">
        <p>
          <strong>Entradas Confirmadas:</strong> Valores já compensados
        </p>
        <p>
          <strong>Entradas Previstas:</strong> Valores aguardando compensação
          (cartão crédito D+30, etc)
        </p>
        <p>
          <strong>Saídas Previstas:</strong> Contas a pagar com vencimento no
          período
        </p>
      </div>
    </div>
  );
}
```

---

## BANCO DE DADOS

### Migrations

**029_alter_meios_pagamento_add_d_mais.up.sql**

```sql
-- Ver seção Backend > Modelagem
```

**030_create_fluxo_caixa_diario.up.sql**

```sql
-- Ver seção Backend > Modelagem
```

**031_create_compensacoes_bancarias.up.sql**

```sql
-- Ver seção Backend > Modelagem
```

### Seeds

**Seed de Meios de Pagamento com D+**

```sql
-- Atualizar meios existentes
UPDATE meios_pagamento SET d_mais = 0 WHERE tipo = 'PIX';
UPDATE meios_pagamento SET d_mais = 0 WHERE tipo = 'DINHEIRO';
UPDATE meios_pagamento SET d_mais = 1 WHERE tipo = 'DEBITO';
UPDATE meios_pagamento SET d_mais = 30 WHERE tipo = 'CREDITO';
UPDATE meios_pagamento SET d_mais = 1 WHERE tipo = 'TRANSFERENCIA';
```

---

## DEPENDÊNCIAS

### Dependências Técnicas

**Backend (Go)**

- ✅ Já instalado: `github.com/shopspring/decimal`, `github.com/labstack/echo/v4`
- ✅ Já instalado: `github.com/robfig/cron/v3` (scheduler)
- 🆕 **Nenhuma nova dependência necessária**

**Frontend (Next.js)**

- ✅ Já instalado: `@tanstack/react-query`, `date-fns`
- 🆕 **Opcional**: `recharts` (para gráficos)
  ```bash
  pnpm add recharts
  ```

**Banco de Dados**

- ✅ PostgreSQL 14+ (Neon)
- ✅ Migrations via `golang-migrate/migrate`

### Dependências de Módulos

1. **Módulo de Meios de Pagamento** (`Tarefas/FINANCEIRO/meios-pagamento.md`)

   - Atualizar para incluir campo `d_mais`
   - Migration e seed de valores padrão

2. **Módulo de Receitas**

   - Atualizar para chamar `CreateCompensacaoUseCase` ao criar receita
   - Vincular `meio_pagamento_id`

3. **Módulo de Contas a Pagar** (`Tarefas/FINANCEIRO/03-contas-a-pagar.md`)

   - Necessário para `saidas_previstas`
   - Fornece método: `ContasAPagarRepo.SumByTenantAndDataVencimento()`

4. **Módulo de Contas a Receber** (`Tarefas/FINANCEIRO/04-contas-a-receber.md`)
   - Opcional para previsões mais completas
   - Pode alimentar `entradas_previstas` futuras

---

## REGRAS DE NEGÓCIO

### Regras Fluxo Compensado

- **RN-FC-001**: PIX e Dinheiro são compensados imediatamente (D+0).
- **RN-FC-002**: Débito é compensado em D+1.
- **RN-FC-003**: Crédito é compensado em D+30 (configurável por bandeira/adquirente).
- **RN-FC-004**: Transferência é compensada em D+1.
- **RN-FC-005**: Compensação é criada automaticamente ao criar receita.
- **RN-FC-006**: Valor líquido = Valor bruto - Taxa percentual - Taxa fixa.
- **RN-FC-007**: Status `PREVISTO` → `COMPENSADO` automaticamente na data de compensação.
- **RN-FC-008**: Fluxo diário é gerado automaticamente via cron às 06:00.
- **RN-FC-009**: Entradas confirmadas são compensações já compensadas.
- **RN-FC-010**: Entradas previstas são compensações ainda pendentes.
- **RN-FC-011**: Saídas previstas vêm de contas a pagar com vencimento no dia.
- **RN-FC-012**: Saldo final = Saldo inicial + Entradas (confirmadas + previstas) - Saídas (pagas + previstas).
- **RN-FC-013**: Job de marcar compensações roda às 07:00 diariamente.
- **RN-FC-014**: Owner/Manager pode marcar compensações manualmente.

---

## CRITÉRIOS DE ACEITE

### Backend

- [ ] Migration `029_alter_meios_pagamento_add_d_mais` aplicada.
- [ ] Migrations `030_create_fluxo_caixa_diario` e `031_create_compensacoes_bancarias` aplicadas.
- [ ] Seeds de D+ aplicados corretamente aos meios de pagamento.
- [ ] Entidade `CompensacaoBancaria` calcula valor líquido corretamente.
- [ ] `CreateCompensacaoUseCase` criado e integrado ao criar receita.
- [ ] `GenerateFluxoCaixaDiarioUseCase` gera fluxo com todas as componentes.
- [ ] `MarcarCompensacaoComoCompensadoUseCase` atualiza status corretamente.
- [ ] Endpoints `/financial/cashflow/compensado` e `/marcar-compensacoes` funcionam.
- [ ] Cron jobs executam nos horários configurados (06:00 e 07:00).
- [ ] Testes unitários cobrem cálculo de compensação e fluxo.
- [ ] Testes de integração garantem isolamento multi-tenant.

### Frontend

- [ ] Hook `useFluxoCaixaCompensado` busca dados do backend.
- [ ] Componente `FluxoCaixaCompensadoTable` exibe todas as colunas.
- [ ] Gráfico `FluxoChart` mostra evolução do saldo.
- [ ] Página `/financeiro/fluxo-caixa-compensado` permite filtrar período.
- [ ] Botão "Marcar Compensações" atualiza status e revalida dados.
- [ ] Cores diferenciadas para valores confirmados (escuro) e previstos (claro).
- [ ] Loading states e error handling implementados.
- [ ] Design System aplicado.

### Banco de Dados

- [ ] Tabela `compensacoes_bancarias` criada com índices otimizados.
- [ ] Tabela `fluxo_caixa_diario` criada com índices.
- [ ] Campo `d_mais` em `meios_pagamento` populado corretamente.
- [ ] Constraint `UNIQUE(tenant_id, data)` em `fluxo_caixa_diario` funciona.
- [ ] Queries de agregação (<500ms p95 com 10k receitas/mês).
- [ ] Rollback migrations funcionam corretamente.

### Integração

- [ ] Criar receita automaticamente gera compensação bancária.
- [ ] Compensações são marcadas como `COMPENSADO` na data correta.
- [ ] Fluxo diário reflete dados de compensações, despesas e contas a pagar.
- [ ] Contas a pagar alimentam `saidas_previstas`.

---

**Status:** Documento técnico completo para implementação do Fluxo de Caixa Compensado.
**Integração:** Conecta-se com DRE, Receitas, Despesas, Contas a Pagar/Receber.
