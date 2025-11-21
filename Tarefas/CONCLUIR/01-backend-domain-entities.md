# 01 - Backend: Domain Entities (BLOQUEADOR)

**Prioridade:** 🔴 CRÍTICA
**Estimativa:** 3-4 dias
**Dependências:** Nenhuma (tabelas DB já criadas)
**Bloqueia:** Todas as outras tarefas de backend

---

## Objetivo

Criar todas as entidades de domínio faltantes para as 19 novas tabelas + completar entidades existentes.

---

## Tarefas Detalhadas

### Grupo 1: Financeiro - DRE

#### 1.1 - Criar `DREMensal` Entity

**Arquivo:** `backend/internal/domain/entity/dre_mensal.go`

```go
package entity

import (
    "time"
    "github.com/shopspring/decimal"
)

type DREMensal struct {
    ID                     string
    TenantID               string
    MesAno                 string // YYYY-MM

    // Receitas
    ReceitaServicos        decimal.Decimal
    ReceitaProdutos        decimal.Decimal
    ReceitaPlanos          decimal.Decimal
    ReceitaTotal           decimal.Decimal

    // Custos Variáveis
    CustoComissoes         decimal.Decimal
    CustoInsumos           decimal.Decimal
    CustoVariavelTotal     decimal.Decimal

    // Despesas
    DespesaFixa            decimal.Decimal
    DespesaVariavel        decimal.Decimal
    DespesaTotal           decimal.Decimal

    // Resultado
    ResultadoBruto         decimal.Decimal
    ResultadoOperacional   decimal.Decimal
    MargemBruta            decimal.Decimal // %
    MargemOperacional      decimal.Decimal // %
    LucroLiquido           decimal.Decimal

    ProcessadoEm           time.Time
    CriadoEm               time.Time
    AtualizadoEm           time.Time
}

func NewDREMensal(tenantID, mesAno string) *DREMensal {
    now := time.Now()
    return &DREMensal{
        ID:       uuid.NewString(),
        TenantID: tenantID,
        MesAno:   mesAno,
        CriadoEm: now,
        AtualizadoEm: now,
    }
}

func (d *DREMensal) Calcular() {
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

    // Lucro Líquido = Resultado Operacional (sem impostos/outros)
    d.LucroLiquido = d.ResultadoOperacional

    // Margens (%)
    if d.ReceitaTotal.GreaterThan(decimal.Zero) {
        d.MargemBruta = d.ResultadoBruto.Div(d.ReceitaTotal).Mul(decimal.NewFromInt(100))
        d.MargemOperacional = d.ResultadoOperacional.Div(d.ReceitaTotal).Mul(decimal.NewFromInt(100))
    }

    d.ProcessadoEm = time.Now()
    d.AtualizadoEm = time.Now()
}

func (d *DREMensal) Validate() error {
    if d.TenantID == "" {
        return ErrTenantIDRequired
    }
    if d.MesAno == "" {
        return ErrMesAnoRequired
    }
    // Validar formato YYYY-MM
    if !isValidMesAno(d.MesAno) {
        return ErrMesAnoInvalid
    }
    return nil
}
```

**Checklist:**

- [ ] Criar arquivo
- [ ] Implementar `NewDREMensal`
- [ ] Implementar `Calcular()`
- [ ] Implementar `Validate()`
- [ ] Adicionar constantes de erro em `errors.go`
- [ ] Criar testes unitários `dre_mensal_test.go`

---

### Grupo 2: Financeiro - Fluxo Compensado

#### 1.2 - Criar `FluxoCaixaDiario` Entity

**Arquivo:** `backend/internal/domain/entity/fluxo_caixa_diario.go`

```go
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

**Checklist:**

- [ ] Criar arquivo
- [ ] Implementar `NewFluxoCaixaDiario`
- [ ] Implementar `Calcular()`
- [ ] Implementar `Validate()`
- [ ] Testes unitários

#### 1.3 - Criar `CompensacaoBancaria` Entity

**Arquivo:** `backend/internal/domain/entity/compensacao_bancaria.go`

```go
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
```

**Checklist:**

- [ ] Criar arquivo
- [ ] Implementar entity completa
- [ ] Métodos de cálculo e transição de status
- [ ] Validações
- [ ] Testes

---

### Grupo 3: Metas

#### 1.4 - Criar `MetaMensal` Entity

**Arquivo:** `backend/internal/domain/entity/meta_mensal.go`

#### 1.5 - Criar `MetaBarbeiro` Entity

**Arquivo:** `backend/internal/domain/entity/meta_barbeiro.go`

#### 1.6 - Criar `MetaTicketMedio` Entity

**Arquivo:** `backend/internal/domain/entity/meta_ticket_medio.go`

**Checklist (para cada):**

- [ ] Criar arquivo
- [ ] Implementar entity
- [ ] Métodos de validação
- [ ] Métodos de cálculo de progresso
- [ ] Testes

---

### Grupo 4: Precificação

#### 1.7 - Criar `PrecificacaoConfig` Entity

**Arquivo:** `backend/internal/domain/entity/precificacao_config.go`

#### 1.8 - Criar `PrecificacaoSimulacao` Entity

**Arquivo:** `backend/internal/domain/entity/precificacao_simulacao.go`

**Checklist:**

- [ ] Criar arquivos
- [ ] Implementar entities
- [ ] Métodos de cálculo de preço
- [ ] Validações (margem, markup, etc)
- [ ] Testes

---

### Grupo 5: Contas

#### 1.9 - Criar `ContaAPagar` Entity

**Arquivo:** `backend/internal/domain/entity/conta_a_pagar.go`

```go
type ContaAPagar struct {
    ID              string
    TenantID        string
    Descricao       string
    CategoriaID     string
    Fornecedor      string
    Valor           decimal.Decimal

    Tipo            string // FIXA, VARIAVEL
    Recorrente      bool
    Periodicidade   string // MENSAL, TRIMESTRAL, ANUAL

    DataVencimento  time.Time
    DataPagamento   *time.Time
    Status          string // ABERTO, PAGO, ATRASADO, CANCELADO

    ComprovanteURL  string
    PixCode         string
    Observacoes     string

    CriadoEm        time.Time
    AtualizadoEm    time.Time
}

func (c *ContaAPagar) MarcarComoPago(dataPagamento time.Time, comprovante string) error {
    if c.Status == "PAGO" {
        return ErrContaJaPaga
    }
    c.DataPagamento = &dataPagamento
    c.ComprovanteURL = comprovante
    c.Status = "PAGO"
    c.AtualizadoEm = time.Now()
    return nil
}

func (c *ContaAPagar) VerificarAtraso() {
    if c.Status == "ABERTO" && time.Now().After(c.DataVencimento) {
        c.Status = "ATRASADO"
    }
}
```

#### 1.10 - Criar `ContaAReceber` Entity

**Arquivo:** `backend/internal/domain/entity/conta_a_receber.go`

**Checklist:**

- [ ] Criar arquivos
- [ ] Implementar entities
- [ ] Métodos de transição de status
- [ ] Métodos de verificação de atraso
- [ ] Lógica de recorrência
- [ ] Testes

---

### Grupo 6: Completar Existentes

#### 1.11 - Completar `BarberCommission` Entity

Atualmente existe apenas uma entity básica. Completar com:

- Métodos de cálculo (fixo/percentual/degrau)
- Validações completas
- Status de pagamento

#### 1.12 - Completar `FinancialSnapshot` Entity

Adicionar métodos de agregação e comparação.

#### 1.13 - Completar `CronRunLog` Entity

Adicionar métodos de registro de execução.

---

## Critérios de Aceite

- [ ] Todas as 19 entidades criadas
- [ ] Cada entity tem método `Validate()`
- [ ] Cada entity tem construtor `New*()`
- [ ] Métodos de cálculo implementados onde aplicável
- [ ] Erros de domínio definidos em `errors.go`
- [ ] Testes unitários para cada entity (cobertura > 80%)
- [ ] Nenhum import de camadas externas (infra/http)
- [ ] Uso correto de `decimal.Decimal` para dinheiro
- [ ] Uso correto de ponteiros para campos opcionais

---

## Arquivos a Criar

```
backend/internal/domain/entity/
├── dre_mensal.go
├── dre_mensal_test.go
├── fluxo_caixa_diario.go
├── fluxo_caixa_diario_test.go
├── compensacao_bancaria.go
├── compensacao_bancaria_test.go
├── meta_mensal.go
├── meta_mensal_test.go
├── meta_barbeiro.go
├── meta_barbeiro_test.go
├── meta_ticket_medio.go
├── meta_ticket_medio_test.go
├── precificacao_config.go
├── precificacao_config_test.go
├── precificacao_simulacao.go
├── precificacao_simulacao_test.go
├── conta_a_pagar.go
├── conta_a_pagar_test.go
├── conta_a_receber.go
└── conta_a_receber_test.go
```

---

## Observações

- Usar sempre `decimal.Decimal` para valores monetários
- Validar tenant_id em TODAS as entities
- Seguir padrão de naming das entities existentes
- Adicionar comentários GoDoc em todas as structs e métodos públicos
- Não acessar banco de dados nas entities (pure domain logic)
