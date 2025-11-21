# CATEGORIA: FINANCEIRO > DRE (Demonstrativo de Resultado do Exercício)

## Plano de Execução (alinhado ao item 02-dre, mesma prioridade)
- **Banco de Dados:** snapshot `dre_mensal` (opcional) e flags em categorias/receitas para classificação; integrar com comissões, insumos, payables/receivables.
- **Backend:** serviços de consolidação mensal, comparação m/m, exportação; cálculos conforme DRE completo.
- **Frontend:** mesma tela/fluxo do DRE principal.
- **Cálculos aplicados:** Margem de Lucro (`docs/10-calculos/margem-lucro.md`), Custo de Insumo (`custo-insumo-servico.md`), Markup (`markup.md`), Faturamento Mínimo e Ponto de Equilíbrio (`faturamento-minimo-mensal.md`, `ponto-de-equilibrio.md`), Ticket Médio/LTV/CAC (`ticket-medio.md`, `ltv.md`, `cac.md`).

## Análise do Sistema Atual

### Estado Implementado

- ✅ Tabelas: `receitas`, `despesas`, `categorias`, `meios_pagamento`
- ✅ Repositórios PostgreSQL com 70+ métodos (Save, FindByTenant, SumByPeriod)
- ✅ Domínios: Receita, Despesa, Categoria, MetodoPagamento
- ✅ Use Cases básicos: CreateReceita, CreateDespesa, GetCashflow
- ✅ API endpoints: `/financial/receitas`, `/financial/despesas`, `/financial/cashflow`
- ✅ Frontend: hooks `useReceitas`, `useDespesas`, `useCashflow`
- ⚠️ **Pendente**: Comissões automáticas, Insumos por serviço, DRE consolidado

### Gap Identificado

1. **Sem tabela de DRE consolidado** → criar `dre_mensal` (snapshot)
2. **Sem vínculo de comissões** → precisa integrar com módulo de comissões (Tarefas/FINANCEIRO/05)
3. **Sem custos de insumos** → precisa integrar com estoque (consumo automático)
4. **Sem categorização fixa vs variável** → adicionar flag `tipo_custo` em categorias

---

## Funcionalidade: DRE Completo e Automatizado

### Objetivo

Gerar e exibir DRE mensal com:

- Receitas (Serviços, Produtos, Planos/Mensalidades)
- Custos Variáveis (Comissões, Insumos)
- Despesas Fixas e Variáveis
- Resultado Operacional, Margem de Lucro, Lucro Líquido
- Comparação mês a mês
- Exportação em PDF

---

## BACKEND

### Tarefas Backend

#### 1. Modelagem de Banco de Dados

**Nova Tabela: `dre_mensal`**

```sql
CREATE TABLE dre_mensal (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    mes_ano VARCHAR(7) NOT NULL, -- YYYY-MM

    -- Receitas
    receita_servicos DECIMAL(15,2) DEFAULT 0,
    receita_produtos DECIMAL(15,2) DEFAULT 0,
    receita_planos DECIMAL(15,2) DEFAULT 0,
    receita_total DECIMAL(15,2) DEFAULT 0,

    -- Custos Variáveis
    custo_comissoes DECIMAL(15,2) DEFAULT 0,
    custo_insumos DECIMAL(15,2) DEFAULT 0,
    custo_variavel_total DECIMAL(15,2) DEFAULT 0,

    -- Despesas
    despesa_fixa DECIMAL(15,2) DEFAULT 0,
    despesa_variavel DECIMAL(15,2) DEFAULT 0,
    despesa_total DECIMAL(15,2) DEFAULT 0,

    -- Resultado
    resultado_bruto DECIMAL(15,2) DEFAULT 0, -- Receita - Custo Variável
    resultado_operacional DECIMAL(15,2) DEFAULT 0, -- Bruto - Despesas
    margem_bruta DECIMAL(5,2) DEFAULT 0, -- %
    margem_operacional DECIMAL(5,2) DEFAULT 0, -- %
    lucro_liquido DECIMAL(15,2) DEFAULT 0,

    -- Metadados
    processado_em TIMESTAMPTZ DEFAULT NOW(),
    criado_em TIMESTAMPTZ DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(tenant_id, mes_ano)
);

CREATE INDEX idx_dre_mensal_tenant ON dre_mensal(tenant_id);
CREATE INDEX idx_dre_mensal_mes_ano ON dre_mensal(tenant_id, mes_ano DESC);
```

**Alterar Tabela: `categorias`**

```sql
-- Adicionar coluna tipo_custo
ALTER TABLE categorias
ADD COLUMN tipo_custo VARCHAR(20) CHECK (tipo_custo IN ('FIXO', 'VARIAVEL')) DEFAULT 'FIXO';

COMMENT ON COLUMN categorias.tipo_custo IS 'Para despesas: FIXO (aluguel, internet) ou VARIAVEL (marketing, insumos)';
```

**Alterar Tabela: `receitas`**

```sql
-- Adicionar coluna subtipo para classificação DRE
ALTER TABLE receitas
ADD COLUMN subtipo VARCHAR(30) CHECK (subtipo IN ('SERVICO', 'PRODUTO', 'PLANO')) DEFAULT 'SERVICO';

COMMENT ON COLUMN receitas.subtipo IS 'Classificação para DRE: SERVICO, PRODUTO ou PLANO/Mensalidade';
```

#### 2. Domain Layer (Go)

**Entidade DRE**

```go
// internal/domain/dre.go
package domain

import (
    "time"
    "github.com/shopspring/decimal"
)

type DRE struct {
    ID                     string
    TenantID               string
    MesAno                 string // YYYY-MM

    ReceitaServicos        decimal.Decimal
    ReceitaProdutos        decimal.Decimal
    ReceitaPlanos          decimal.Decimal
    ReceitaTotal           decimal.Decimal

    CustoComissoes         decimal.Decimal
    CustoInsumos           decimal.Decimal
    CustoVariavelTotal     decimal.Decimal

    DespesaFixa            decimal.Decimal
    DespesaVariavel        decimal.Decimal
    DespesaTotal           decimal.Decimal

    ResultadoBruto         decimal.Decimal
    ResultadoOperacional   decimal.Decimal
    MargemBruta            decimal.Decimal
    MargemOperacional      decimal.Decimal
    LucroLiquido           decimal.Decimal

    ProcessadoEm           time.Time
    CriadoEm               time.Time
    AtualizadoEm           time.Time
}

func (d *DRE) Calcular() {
    // Receita Total
    d.ReceitaTotal = d.ReceitaServicos.Add(d.ReceitaProdutos).Add(d.ReceitaPlanos)

    // Custo Variável Total
    d.CustoVariavelTotal = d.CustoComissoes.Add(d.CustoInsumos)

    // Despesa Total
    d.DespesaTotal = d.DespesaFixa.Add(d.DespesaVariavel)

    // Resultado Bruto = Receita - Custo Variável
    d.ResultadoBruto = d.ReceitaTotal.Sub(d.CustoVariavelTotal)

    // Resultado Operacional = Bruto - Despesas
    d.ResultadoOperacional = d.ResultadoBruto.Sub(d.DespesaTotal)

    // Lucro Líquido = Resultado Operacional (sem impostos ainda)
    d.LucroLiquido = d.ResultadoOperacional

    // Margens
    if !d.ReceitaTotal.IsZero() {
        d.MargemBruta = d.ResultadoBruto.Div(d.ReceitaTotal).Mul(decimal.NewFromInt(100))
        d.MargemOperacional = d.ResultadoOperacional.Div(d.ReceitaTotal).Mul(decimal.NewFromInt(100))
    }
}

func (d *DRE) Validate() error {
    if d.MesAno == "" {
        return ErrDREMesAnoRequired
    }
    if d.TenantID == "" {
        return ErrTenantIDRequired
    }
    return nil
}
```

**Repository Interface**

```go
// internal/domain/ports/dre_repository.go
type DRERepository interface {
    Save(ctx context.Context, tenantID string, dre *DRE) error
    FindByTenantAndMonth(ctx context.Context, tenantID, mesAno string) (*DRE, error)
    FindByTenantAndPeriod(ctx context.Context, tenantID string, inicio, fim string) ([]*DRE, error)
    Update(ctx context.Context, tenantID string, dre *DRE) error
}
```

#### 3. Application Layer (Use Cases)

**GenerateDREUseCase**

```go
// internal/application/usecase/generate_dre_usecase.go
package usecase

type GenerateDREUseCase struct {
    dreRepo         DRERepository
    receitaRepo     ReceitaRepository
    despesaRepo     DespesaRepository
    comissaoService ComissaoService // Integração com módulo de comissões
    estoqueService  EstoqueService  // Integração com consumo de insumos
}

type GenerateDREInput struct {
    TenantID string
    MesAno   string // YYYY-MM
}

func (uc *GenerateDREUseCase) Execute(ctx context.Context, input GenerateDREInput) (*DRE, error) {
    // 1. Buscar receitas do mês por subtipo
    inicio := input.MesAno + "-01"
    fim := CalcularUltimoDiaMes(input.MesAno)

    receitaServicos, _ := uc.receitaRepo.SumByTenantPeriodAndSubtipo(ctx, input.TenantID, inicio, fim, "SERVICO", ReceiptReceived)
    receitaProdutos, _ := uc.receitaRepo.SumByTenantPeriodAndSubtipo(ctx, input.TenantID, inicio, fim, "PRODUTO", ReceiptReceived)
    receitaPlanos, _ := uc.receitaRepo.SumByTenantPeriodAndSubtipo(ctx, input.TenantID, inicio, fim, "PLANO", ReceiptReceived)

    // 2. Buscar comissões do mês
    custoComissoes, _ := uc.comissaoService.SumByTenantAndPeriod(ctx, input.TenantID, inicio, fim)

    // 3. Buscar insumos consumidos (do estoque)
    custoInsumos, _ := uc.estoqueService.SumCustoInsumosByPeriod(ctx, input.TenantID, inicio, fim)

    // 4. Buscar despesas por tipo_custo
    despesaFixa, _ := uc.despesaRepo.SumByTenantPeriodAndTipoCusto(ctx, input.TenantID, inicio, fim, "FIXO", ExpensePaid)
    despesaVariavel, _ := uc.despesaRepo.SumByTenantPeriodAndTipoCusto(ctx, input.TenantID, inicio, fim, "VARIAVEL", ExpensePaid)

    // 5. Criar DRE
    dre := &DRE{
        ID:              uuid.NewString(),
        TenantID:        input.TenantID,
        MesAno:          input.MesAno,
        ReceitaServicos: receitaServicos,
        ReceitaProdutos: receitaProdutos,
        ReceitaPlanos:   receitaPlanos,
        CustoComissoes:  custoComissoes,
        CustoInsumos:    custoInsumos,
        DespesaFixa:     despesaFixa,
        DespesaVariavel: despesaVariavel,
        ProcessadoEm:    time.Now(),
        CriadoEm:        time.Now(),
    }

    dre.Calcular()

    // 6. Salvar
    if err := uc.dreRepo.Save(ctx, input.TenantID, dre); err != nil {
        return nil, err
    }

    return dre, nil
}
```

**GetDREComparisonUseCase**

```go
// Comparar DRE mês a mês
func (uc *GetDREComparisonUseCase) Execute(ctx context.Context, tenantID, mesIni, mesFim string) (*DREComparison, error) {
    dres, _ := uc.dreRepo.FindByTenantAndPeriod(ctx, tenantID, mesIni, mesFim)

    // Calcular variações percentuais entre meses
    comparisons := make([]DREMonthComparison, len(dres)-1)
    for i := 1; i < len(dres); i++ {
        comparisons[i-1] = CompararDRE(dres[i-1], dres[i])
    }

    return &DREComparison{
        Meses:       dres,
        Comparisons: comparisons,
    }, nil
}
```

#### 4. Infrastructure Layer (Repository)

**PostgreSQL Implementation**

```go
// internal/infrastructure/persistence/postgres_dre_repository.go
func (r *PostgresDRERepository) Save(ctx context.Context, tenantID string, dre *DRE) error {
    query := `
        INSERT INTO dre_mensal (
            id, tenant_id, mes_ano,
            receita_servicos, receita_produtos, receita_planos, receita_total,
            custo_comissoes, custo_insumos, custo_variavel_total,
            despesa_fixa, despesa_variavel, despesa_total,
            resultado_bruto, resultado_operacional,
            margem_bruta, margem_operacional, lucro_liquido,
            processado_em, criado_em
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
        ON CONFLICT (tenant_id, mes_ano)
        DO UPDATE SET
            receita_servicos = EXCLUDED.receita_servicos,
            -- ... todos os campos
            atualizado_em = NOW()
    `
    // Executar query
}
```

#### 5. HTTP Layer (Handlers)

**DREHandler**

```go
// internal/infrastructure/http/handler/dre_handler.go
func (h *DREHandler) GenerateDRE(c echo.Context) error {
    // POST /financial/dre/generate
    // Body: { "mes_ano": "2024-11" }

    tenantID := c.Get("tenant_id").(string)

    var req GenerateDRERequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(400, ErrorResponse{Message: "Invalid request"})
    }

    dre, err := h.generateDREUseCase.Execute(c.Request().Context(), GenerateDREInput{
        TenantID: tenantID,
        MesAno:   req.MesAno,
    })

    if err != nil {
        return c.JSON(500, ErrorResponse{Message: err.Error()})
    }

    return c.JSON(201, MapDREToResponse(dre))
}

func (h *DREHandler) GetDRE(c echo.Context) error {
    // GET /financial/dre?mes_ano=2024-11
    tenantID := c.Get("tenant_id").(string)
    mesAno := c.QueryParam("mes_ano")

    dre, err := h.dreRepo.FindByTenantAndMonth(c.Request().Context(), tenantID, mesAno)
    if err != nil {
        return c.JSON(404, ErrorResponse{Message: "DRE not found"})
    }

    return c.JSON(200, MapDREToResponse(dre))
}

func (h *DREHandler) GetDREComparison(c echo.Context) error {
    // GET /financial/dre/comparison?inicio=2024-01&fim=2024-11
    tenantID := c.Get("tenant_id").(string)
    inicio := c.QueryParam("inicio")
    fim := c.QueryParam("fim")

    comparison, err := h.getComparisonUseCase.Execute(c.Request().Context(), tenantID, inicio, fim)
    if err != nil {
        return c.JSON(500, ErrorResponse{Message: err.Error()})
    }

    return c.JSON(200, comparison)
}

func (h *DREHandler) ExportDREPDF(c echo.Context) error {
    // POST /financial/dre/export
    // Gerar PDF usando template e retornar link de download
}
```

#### 6. Cron Job (Geração Automática)

**Job Mensal de DRE**

```go
// Executar no dia 1º de cada mês às 05:00
func (j *GenerateDREJob) Run() {
    ctx := context.Background()

    // Para cada tenant ativo
    tenants, _ := j.tenantRepo.FindActive(ctx)

    // Gerar DRE do mês anterior
    mesAnterior := time.Now().AddDate(0, -1, 0).Format("2006-01")

    for _, tenant := range tenants {
        _, err := j.generateDREUseCase.Execute(ctx, GenerateDREInput{
            TenantID: tenant.ID,
            MesAno:   mesAnterior,
        })

        if err != nil {
            log.Error("Failed to generate DRE",
                zap.String("tenant_id", tenant.ID),
                zap.String("mes_ano", mesAnterior),
                zap.Error(err))
        }
    }
}
```

---

## FRONTEND

### Tarefas Frontend

#### 1. Hooks Customizados

**useDRE**

```typescript
// app/hooks/useDRE.ts
export function useDRE(mesAno: string) {
  return useQuery({
    queryKey: ["dre", mesAno],
    queryFn: async () => {
      const res = await api.get(`/financial/dre`, {
        params: { mes_ano: mesAno },
      });
      return res.data;
    },
    staleTime: 1000 * 60 * 5, // 5 min
  });
}

export function useDREComparison(inicio: string, fim: string) {
  return useQuery({
    queryKey: ["dre-comparison", inicio, fim],
    queryFn: async () => {
      const res = await api.get(`/financial/dre/comparison`, {
        params: { inicio, fim },
      });
      return res.data;
    },
  });
}

export function useGenerateDRE() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (mesAno: string) => {
      const res = await api.post(`/financial/dre/generate`, {
        mes_ano: mesAno,
      });
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["dre"] });
      toast.success("DRE gerado com sucesso");
    },
  });
}
```

#### 2. Componentes UI

**DRECard**

```tsx
// app/components/financeiro/DRECard.tsx
export function DRECard({ dre }: { dre: DRE }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>DRE - {dre.mes_ano}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {/* Receitas */}
          <Section title="Receitas">
            <Line label="Serviços" value={dre.receita_servicos} />
            <Line label="Produtos" value={dre.receita_produtos} />
            <Line label="Planos" value={dre.receita_planos} />
            <LineBold label="Total" value={dre.receita_total} />
          </Section>

          {/* Custos Variáveis */}
          <Section title="Custos Variáveis">
            <Line label="Comissões" value={dre.custo_comissoes} />
            <Line label="Insumos" value={dre.custo_insumos} />
            <LineBold label="Total" value={dre.custo_variavel_total} />
          </Section>

          {/* Resultado Bruto */}
          <LineHighlight label="Resultado Bruto" value={dre.resultado_bruto} />

          {/* Despesas */}
          <Section title="Despesas">
            <Line label="Fixas" value={dre.despesa_fixa} />
            <Line label="Variáveis" value={dre.despesa_variavel} />
            <LineBold label="Total" value={dre.despesa_total} />
          </Section>

          {/* Resultado Final */}
          <LineHighlight
            label="Lucro Líquido"
            value={dre.lucro_liquido}
            isPositive={dre.lucro_liquido > 0}
          />

          <div className="flex justify-between text-sm text-muted-foreground">
            <span>Margem Bruta: {dre.margem_bruta}%</span>
            <span>Margem Operacional: {dre.margem_operacional}%</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
```

**DREComparisonTable**

```tsx
// Tabela comparativa mês a mês com variações
export function DREComparisonTable({ comparison }: Props) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Indicador</TableHead>
          {comparison.meses.map((m) => (
            <TableHead key={m.mes_ano}>{m.mes_ano}</TableHead>
          ))}
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow>
          <TableCell>Receita Total</TableCell>
          {comparison.meses.map((m) => (
            <TableCell key={m.mes_ano}>
              {formatCurrency(m.receita_total)}
            </TableCell>
          ))}
        </TableRow>
        {/* ... demais linhas */}
      </TableBody>
    </Table>
  );
}
```

#### 3. Páginas

**app/financeiro/dre/page.tsx**

```tsx
"use client";

export default function DREPage() {
  const [mesAno, setMesAno] = useState(getCurrentMonth());
  const { data: dre, isLoading } = useDRE(mesAno);
  const generateDRE = useGenerateDRE();

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1>Demonstrativo de Resultado (DRE)</h1>
        <div className="flex gap-2">
          <MonthPicker value={mesAno} onChange={setMesAno} />
          <Button onClick={() => generateDRE.mutate(mesAno)}>Gerar DRE</Button>
          <Button variant="outline">Exportar PDF</Button>
        </div>
      </div>

      {isLoading && <Skeleton />}
      {dre && <DRECard dre={dre} />}

      {/* Comparativo últimos 6 meses */}
      <DREComparisonSection />
    </div>
  );
}
```

---

## BANCO DE DADOS

### Migrations

**Arquivo: `026_create_dre_mensal.up.sql`**

```sql
-- Ver seção Backend > Modelagem de Banco de Dados
```

**Arquivo: `027_alter_categorias_add_tipo_custo.up.sql`**

```sql
ALTER TABLE categorias
ADD COLUMN tipo_custo VARCHAR(20) CHECK (tipo_custo IN ('FIXO', 'VARIAVEL')) DEFAULT 'FIXO';
```

**Arquivo: `028_alter_receitas_add_subtipo.up.sql`**

```sql
ALTER TABLE receitas
ADD COLUMN subtipo VARCHAR(30) CHECK (subtipo IN ('SERVICO', 'PRODUTO', 'PLANO')) DEFAULT 'SERVICO';
```

---

## DEPENDÊNCIAS

### Dependências Técnicas

**Backend (Go)**

- ✅ Já instalado: `github.com/shopspring/decimal` (precisão monetária)
- ✅ Já instalado: `github.com/labstack/echo/v4` (HTTP)
- 🆕 **Instalar**: `github.com/jung-kurt/gofpdf` (geração de PDF)
  ```bash
  go get github.com/jung-kurt/gofpdf
  ```

**Frontend (Next.js)**

- ✅ Já instalado: `@tanstack/react-query`, `zod`, `react-hook-form`
- 🆕 **Instalar**: `jspdf` ou `pdfmake` (geração PDF no cliente, opcional)
  ```bash
  pnpm add jspdf
  ```

**Banco de Dados**

- ✅ PostgreSQL 14+ (Neon)
- ✅ Migrations via `golang-migrate/migrate`

### Dependências de Módulos

1. **Módulo de Comissões** (`Tarefas/FINANCEIRO/05-comissoes-automaticas.md`)

   - Precisa estar implementado para calcular `custo_comissoes`
   - Fornece método: `ComissaoService.SumByTenantAndPeriod()`

2. **Módulo de Estoque** (`Tarefas/ESTOQUE/03-consumo-automatico.md`)

   - Precisa rastrear consumo de insumos por serviço
   - Fornece método: `EstoqueService.SumCustoInsumosByPeriod()`

3. **Módulo de Categorias**

   - Atualizar para incluir `tipo_custo` (FIXO/VARIAVEL)
   - Migration e seed de categorias padrão

4. **Módulo de Receitas**
   - Atualizar para incluir `subtipo` (SERVICO/PRODUTO/PLANO)
   - Lógica de classificação automática (ex: receitas de assinatura = PLANO)

---

## REGRAS DE NEGÓCIO

### Regras DRE

- **RN-DRE-001**: DRE é gerado mensalmente de forma automática (dia 1º às 05:00).
- **RN-DRE-002**: Apenas receitas com status `RECEBIDO` entram no cálculo.
- **RN-DRE-003**: Apenas despesas com status `PAGO` entram no cálculo.
- **RN-DRE-004**: Comissões são custos variáveis vinculados a receitas confirmadas.
- **RN-DRE-005**: Insumos consumidos automaticamente (via serviço) entram como custo variável.
- **RN-DRE-006**: Despesas fixas são categorizadas com `tipo_custo = FIXO`.
- **RN-DRE-007**: Despesas variáveis são categorizadas com `tipo_custo = VARIAVEL`.
- **RN-DRE-008**: Margem Bruta = (Resultado Bruto / Receita Total) \* 100.
- **RN-DRE-009**: Margem Operacional = (Resultado Operacional / Receita Total) \* 100.
- **RN-DRE-010**: DRE pode ser regenerado manualmente (sobrescreve anterior).
- **RN-DRE-011**: Exportação PDF requer permissão (Owner/Manager/Accountant).
- **RN-DRE-012**: Comparação mês a mês exibe variação percentual de cada indicador.

---

## CRITÉRIOS DE ACEITE

### Backend

- [ ] Migration `026_create_dre_mensal` aplicada com sucesso.
- [ ] Alterações em `categorias` e `receitas` aplicadas.
- [ ] Entidade `DRE` valida e calcula corretamente todos os campos.
- [ ] `GenerateDREUseCase` agrega dados de receitas, despesas, comissões e insumos.
- [ ] Endpoints `/financial/dre` (GET/POST) respondem corretamente.
- [ ] Cron job gera DRE automaticamente no dia 1º de cada mês.
- [ ] Comparação mês a mês retorna variações percentuais.
- [ ] Exportação PDF funcional (template legível).
- [ ] Testes unitários cobrem cálculo de DRE e validações.
- [ ] Testes de integração garantem isolamento multi-tenant.

### Frontend

- [ ] Hook `useDRE` busca dados do backend.
- [ ] Componente `DRECard` exibe todos os blocos (receitas, custos, despesas, resultado).
- [ ] Página `/financeiro/dre` permite selecionar mês e gerar DRE.
- [ ] Comparativo mês a mês exibe tabela com variações.
- [ ] Botão "Exportar PDF" baixa arquivo formatado.
- [ ] Loading states e error handling implementados.
- [ ] Design System aplicado (tokens, cores, tipografia).

### Banco de Dados

- [ ] Tabela `dre_mensal` criada com índices otimizados.
- [ ] Constraint `UNIQUE(tenant_id, mes_ano)` funciona.
- [ ] Queries de agregação (<500ms p95 com 10k receitas/mês).
- [ ] Rollback migrations funcionam corretamente.

### Integração

- [ ] Módulo de Comissões fornece dados corretos.
- [ ] Módulo de Estoque fornece custo de insumos.
- [ ] Receitas classificadas automaticamente por subtipo.
- [ ] Categorias marcadas corretamente como FIXO/VARIAVEL.

---

**Status:** Documento técnico completo para implementação do DRE.
**Próximo:** Fluxo de Caixa Compensado (documento separado).
