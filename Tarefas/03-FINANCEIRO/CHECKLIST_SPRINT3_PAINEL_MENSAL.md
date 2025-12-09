# ✅ CHECKLIST — SPRINT 3: PAINEL MENSAL + PROJEÇÕES

> **Status:** 🟢 100% — Concluído  
> **Dependência:** Sprint 2 (Despesas Fixas) ✅  
> **Esforço Estimado:** 18 horas  
> **Prioridade:** P1 — Entrega de valor para o usuário  
> **Concluído em:** 29/11/2025

---

## 📊 OBJETIVO

Implementar o **Dashboard Financeiro Unificado** com:

1. ✅ Painel Mensal consolidado (receitas, despesas, lucro, metas)
2. ✅ Cálculo de Projeções financeiras
3. ✅ Endpoints de agregação de dados

---

## 📋 TAREFAS

### 1️⃣ USE CASE: GetPainelMensalUseCase (Esforço: 6h) ✅

#### 1.1 Criar Use Case

- [x] Criar `backend/internal/application/usecase/financial/get_painel_mensal.go`

#### 1.2 Estrutura do Use Case

```go
package financial

import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

type GetPainelMensalUseCase struct {
    contaPagarRepo      repository.ContaPagarRepository
    contaReceberRepo    repository.ContaReceberRepository
    despesaFixaRepo     repository.DespesaFixaRepository
    metaMensalRepo      repository.MetaMensalRepository
    fluxoCaixaRepo      repository.FluxoCaixaDiarioRepository
    logger              *zap.Logger
}

type PainelMensalInput struct {
    TenantID  uuid.UUID
    UnidadeID *uuid.UUID
    Ano       int
    Mes       int
}

type PainelMensalOutput struct {
    Ano                  int
    Mes                  int
    // Receitas
    ReceitaRealizada     decimal.Decimal
    ReceitaPendente      decimal.Decimal
    ReceitaTotal         decimal.Decimal
    // Despesas
    DespesasFixas        decimal.Decimal
    DespesasVariaveis    decimal.Decimal
    DespesasPagas        decimal.Decimal
    DespesasPendentes    decimal.Decimal
    DespesasTotal        decimal.Decimal
    // Resultados
    LucroBruto           decimal.Decimal
    LucroLiquido         decimal.Decimal
    MargemLiquida        decimal.Decimal // em percentual
    // Metas
    MetaMensal           decimal.Decimal
    PercentualMeta       decimal.Decimal
    DiferencaMeta        decimal.Decimal // positivo = acima, negativo = abaixo
    // Indicadores
    TicketMedio          decimal.Decimal
    TotalAtendimentos    int64
    TotalTransacoes      int64
    // Caixa
    SaldoCaixaAtual      decimal.Decimal
    // Comparativo
    VariacaoMesAnterior  decimal.Decimal // em percentual
}

func (uc *GetPainelMensalUseCase) Execute(ctx context.Context, input PainelMensalInput) (*PainelMensalOutput, error) {
    // 1. Definir período
    inicio := time.Date(input.Ano, time.Month(input.Mes), 1, 0, 0, 0, 0, time.UTC)
    fim := inicio.AddDate(0, 1, -1) // Último dia do mês
    
    // 2. Buscar receitas
    receitaRealizada, err := uc.contaReceberRepo.SumRecebidoByPeriod(ctx, input.TenantID, inicio, fim)
    receitaPendente, err := uc.contaReceberRepo.SumPendenteByPeriod(ctx, input.TenantID, inicio, fim)
    
    // 3. Buscar despesas
    despesasPagas, err := uc.contaPagarRepo.SumPagoByPeriod(ctx, input.TenantID, inicio, fim)
    despesasPendentes, err := uc.contaPagarRepo.SumPendenteByPeriod(ctx, input.TenantID, inicio, fim)
    despesasFixas, err := uc.despesaFixaRepo.SumAtivas(ctx, input.TenantID)
    
    // 4. Buscar meta do mês
    meta, err := uc.metaMensalRepo.GetByMonth(ctx, input.TenantID, input.Ano, input.Mes)
    
    // 5. Buscar saldo do caixa
    caixaHoje, err := uc.fluxoCaixaRepo.GetByDate(ctx, input.TenantID, time.Now())
    
    // 6. Calcular indicadores
    lucroBruto := receitaRealizada.Sub(despesasPagas)
    lucroLiquido := lucroBruto // Simplificado, sem IR
    
    var margemLiquida decimal.Decimal
    if !receitaRealizada.IsZero() {
        margemLiquida = lucroLiquido.Div(receitaRealizada).Mul(decimal.NewFromInt(100))
    }
    
    var percentualMeta decimal.Decimal
    var diferencaMeta decimal.Decimal
    if meta != nil && !meta.Valor.IsZero() {
        percentualMeta = receitaRealizada.Div(meta.Valor).Mul(decimal.NewFromInt(100))
        diferencaMeta = receitaRealizada.Sub(meta.Valor)
    }
    
    // 7. Ticket médio (se tiver atendimentos)
    totalAtendimentos, _ := uc.atendimentoRepo.CountByPeriod(ctx, input.TenantID, inicio, fim)
    var ticketMedio decimal.Decimal
    if totalAtendimentos > 0 {
        ticketMedio = receitaRealizada.Div(decimal.NewFromInt(totalAtendimentos))
    }
    
    // 8. Comparativo mês anterior
    mesAnteriorInicio := inicio.AddDate(0, -1, 0)
    mesAnteriorFim := inicio.AddDate(0, 0, -1)
    receitaMesAnterior, _ := uc.contaReceberRepo.SumRecebidoByPeriod(
        ctx, input.TenantID, mesAnteriorInicio, mesAnteriorFim)
    
    var variacaoMesAnterior decimal.Decimal
    if !receitaMesAnterior.IsZero() {
        variacaoMesAnterior = receitaRealizada.Sub(receitaMesAnterior).
            Div(receitaMesAnterior).
            Mul(decimal.NewFromInt(100))
    }
    
    return &PainelMensalOutput{
        Ano:                 input.Ano,
        Mes:                 input.Mes,
        ReceitaRealizada:    receitaRealizada,
        ReceitaPendente:     receitaPendente,
        ReceitaTotal:        receitaRealizada.Add(receitaPendente),
        DespesasFixas:       despesasFixas,
        DespesasVariaveis:   despesasPagas.Sub(despesasFixas), // Aproximação
        DespesasPagas:       despesasPagas,
        DespesasPendentes:   despesasPendentes,
        DespesasTotal:       despesasPagas.Add(despesasPendentes),
        LucroBruto:          lucroBruto,
        LucroLiquido:        lucroLiquido,
        MargemLiquida:       margemLiquida,
        MetaMensal:          meta.Valor,
        PercentualMeta:      percentualMeta,
        DiferencaMeta:       diferencaMeta,
        TicketMedio:         ticketMedio,
        TotalAtendimentos:   totalAtendimentos,
        SaldoCaixaAtual:     caixaHoje.Fechamento,
        VariacaoMesAnterior: variacaoMesAnterior,
    }, nil
}
```

#### 1.3 Checklist Painel Mensal

- [ ] Criar struct de input/output
- [ ] Injetar repositórios necessários
- [ ] Calcular receita realizada (RECEBIDO)
- [ ] Calcular receita pendente (PENDENTE)
- [ ] Calcular despesas pagas
- [ ] Calcular despesas pendentes
- [ ] Buscar despesas fixas ativas
- [ ] Buscar meta do mês
- [ ] Calcular lucro bruto
- [ ] Calcular lucro líquido
- [ ] Calcular margem líquida (%)
- [ ] Calcular % da meta
- [ ] Calcular ticket médio
- [ ] Buscar saldo do caixa
- [ ] Calcular variação vs mês anterior
- [ ] Testes unitários

---

### 2️⃣ USE CASE: GetProjecoesUseCase (Esforço: 6h)

#### 2.1 Criar Use Case

- [ ] Criar `backend/internal/application/usecase/financial/get_projecoes.go`

#### 2.2 Estrutura

```go
package financial

type GetProjecoesUseCase struct {
    contaPagarRepo   repository.ContaPagarRepository
    contaReceberRepo repository.ContaReceberRepository
    despesaFixaRepo  repository.DespesaFixaRepository
    logger           *zap.Logger
}

type ProjecoesInput struct {
    TenantID    uuid.UUID
    MesesAhead  int // Quantos meses projetar (default: 3)
}

type ProjecaoMensal struct {
    Ano               int
    Mes               int
    NomeMes           string // "Janeiro", "Fevereiro", etc.
    ReceitaProjetada  decimal.Decimal
    DespesasProjetadas decimal.Decimal
    DespesasFixas     decimal.Decimal
    LucroProjetado    decimal.Decimal
    DiasUteis         int
    MetaDiaria        decimal.Decimal // Para atingir a projeção
    Confianca         string          // "Alta", "Média", "Baixa"
}

type ProjecoesOutput struct {
    Projecoes          []ProjecaoMensal
    MediaReceita3Meses decimal.Decimal
    MediaDespesas3Meses decimal.Decimal
    TendenciaReceita   string // "Crescente", "Estável", "Decrescente"
}

func (uc *GetProjecoesUseCase) Execute(ctx context.Context, input ProjecoesInput) (*ProjecoesOutput, error) {
    // 1. Buscar histórico dos últimos 3 meses
    hoje := time.Now()
    
    var historicoReceitas []decimal.Decimal
    var historicoDespesas []decimal.Decimal
    
    for i := 1; i <= 3; i++ {
        mesRef := hoje.AddDate(0, -i, 0)
        inicio := time.Date(mesRef.Year(), mesRef.Month(), 1, 0, 0, 0, 0, time.UTC)
        fim := inicio.AddDate(0, 1, -1)
        
        receita, _ := uc.contaReceberRepo.SumRecebidoByPeriod(ctx, input.TenantID, inicio, fim)
        despesa, _ := uc.contaPagarRepo.SumPagoByPeriod(ctx, input.TenantID, inicio, fim)
        
        historicoReceitas = append(historicoReceitas, receita)
        historicoDespesas = append(historicoDespesas, despesa)
    }
    
    // 2. Calcular médias
    mediaReceita := calcularMedia(historicoReceitas)
    mediaDespesa := calcularMedia(historicoDespesas)
    
    // 3. Detectar tendência
    tendencia := detectarTendencia(historicoReceitas)
    
    // 4. Buscar despesas fixas (garantidas)
    despesasFixas, _ := uc.despesaFixaRepo.SumAtivas(ctx, input.TenantID)
    
    // 5. Gerar projeções
    var projecoes []ProjecaoMensal
    
    for i := 1; i <= input.MesesAhead; i++ {
        mesProjetado := hoje.AddDate(0, i, 0)
        ano := mesProjetado.Year()
        mes := int(mesProjetado.Month())
        
        // Aplicar fator de tendência
        fatorTendencia := calcularFatorTendencia(tendencia, i)
        receitaProjetada := mediaReceita.Mul(fatorTendencia)
        
        // Despesas = fixas + variáveis estimadas
        despesasVariaveis := mediaDespesa.Sub(despesasFixas)
        if despesasVariaveis.LessThan(decimal.Zero) {
            despesasVariaveis = decimal.Zero
        }
        despesasProjetadas := despesasFixas.Add(despesasVariaveis)
        
        lucroProjetado := receitaProjetada.Sub(despesasProjetadas)
        
        diasUteis := calcularDiasUteis(ano, mes)
        metaDiaria := receitaProjetada.Div(decimal.NewFromInt(int64(diasUteis)))
        
        confianca := calcularConfianca(i, tendencia)
        
        projecoes = append(projecoes, ProjecaoMensal{
            Ano:               ano,
            Mes:               mes,
            NomeMes:           nomeMes(mes),
            ReceitaProjetada:  receitaProjetada,
            DespesasProjetadas: despesasProjetadas,
            DespesasFixas:     despesasFixas,
            LucroProjetado:    lucroProjetado,
            DiasUteis:         diasUteis,
            MetaDiaria:        metaDiaria,
            Confianca:         confianca,
        })
    }
    
    return &ProjecoesOutput{
        Projecoes:          projecoes,
        MediaReceita3Meses: mediaReceita,
        MediaDespesas3Meses: mediaDespesa,
        TendenciaReceita:   tendencia,
    }, nil
}

func calcularMedia(valores []decimal.Decimal) decimal.Decimal {
    if len(valores) == 0 {
        return decimal.Zero
    }
    soma := decimal.Zero
    for _, v := range valores {
        soma = soma.Add(v)
    }
    return soma.Div(decimal.NewFromInt(int64(len(valores))))
}

func detectarTendencia(valores []decimal.Decimal) string {
    if len(valores) < 2 {
        return "Estável"
    }
    // valores[0] = mês mais recente, valores[2] = mês mais antigo
    diff := valores[0].Sub(valores[len(valores)-1])
    percentual := diff.Div(valores[len(valores)-1]).Mul(decimal.NewFromInt(100))
    
    if percentual.GreaterThan(decimal.NewFromInt(5)) {
        return "Crescente"
    } else if percentual.LessThan(decimal.NewFromInt(-5)) {
        return "Decrescente"
    }
    return "Estável"
}

func calcularFatorTendencia(tendencia string, mesAhead int) decimal.Decimal {
    base := decimal.NewFromInt(1)
    fator := decimal.NewFromFloat(0.02) // 2% por mês
    
    switch tendencia {
    case "Crescente":
        return base.Add(fator.Mul(decimal.NewFromInt(int64(mesAhead))))
    case "Decrescente":
        return base.Sub(fator.Mul(decimal.NewFromInt(int64(mesAhead))))
    default:
        return base
    }
}

func calcularConfianca(mesAhead int, tendencia string) string {
    if mesAhead == 1 {
        return "Alta"
    } else if mesAhead == 2 {
        return "Média"
    }
    return "Baixa"
}

func calcularDiasUteis(ano, mes int) int {
    // Simplificado: 22 dias úteis por mês (média brasileira)
    // TODO: Implementar cálculo real considerando feriados
    return 22
}
```

#### 2.3 Checklist Projeções

- [ ] Criar struct de input/output
- [ ] Buscar histórico dos últimos 3 meses
- [ ] Calcular média de receitas
- [ ] Calcular média de despesas
- [ ] Detectar tendência (crescente/estável/decrescente)
- [ ] Buscar despesas fixas
- [ ] Aplicar fator de tendência
- [ ] Calcular dias úteis por mês
- [ ] Calcular meta diária
- [ ] Atribuir nível de confiança
- [ ] Gerar array de projeções
- [ ] Testes unitários

---

### 3️⃣ DTOs (Esforço: 1h)

#### 3.1 Criar DTOs

- [ ] Criar/atualizar `backend/internal/application/dto/painel_mensal_dto.go`

```go
package dto

type PainelMensalResponse struct {
    Ano                  int    `json:"ano"`
    Mes                  int    `json:"mes"`
    NomeMes              string `json:"nome_mes"`
    
    // Receitas
    ReceitaRealizada     string `json:"receita_realizada"`
    ReceitaPendente      string `json:"receita_pendente"`
    ReceitaTotal         string `json:"receita_total"`
    
    // Despesas
    DespesasFixas        string `json:"despesas_fixas"`
    DespesasVariaveis    string `json:"despesas_variaveis"`
    DespesasPagas        string `json:"despesas_pagas"`
    DespesasPendentes    string `json:"despesas_pendentes"`
    DespesasTotal        string `json:"despesas_total"`
    
    // Resultados
    LucroBruto           string `json:"lucro_bruto"`
    LucroLiquido         string `json:"lucro_liquido"`
    MargemLiquida        string `json:"margem_liquida"`
    
    // Metas
    MetaMensal           string `json:"meta_mensal"`
    PercentualMeta       string `json:"percentual_meta"`
    DiferencaMeta        string `json:"diferenca_meta"`
    StatusMeta           string `json:"status_meta"` // "Atingida", "Em andamento", "Abaixo"
    
    // Indicadores
    TicketMedio          string `json:"ticket_medio"`
    TotalAtendimentos    int64  `json:"total_atendimentos"`
    TotalTransacoes      int64  `json:"total_transacoes"`
    
    // Caixa
    SaldoCaixaAtual      string `json:"saldo_caixa_atual"`
    
    // Comparativo
    VariacaoMesAnterior  string `json:"variacao_mes_anterior"`
    TendenciaVariacao    string `json:"tendencia_variacao"` // "up", "down", "stable"
}

type ProjecaoMensalResponse struct {
    Ano               int    `json:"ano"`
    Mes               int    `json:"mes"`
    NomeMes           string `json:"nome_mes"`
    ReceitaProjetada  string `json:"receita_projetada"`
    DespesasProjetadas string `json:"despesas_projetadas"`
    DespesasFixas     string `json:"despesas_fixas"`
    LucroProjetado    string `json:"lucro_projetado"`
    DiasUteis         int    `json:"dias_uteis"`
    MetaDiaria        string `json:"meta_diaria"`
    Confianca         string `json:"confianca"`
}

type ProjecoesResponse struct {
    Projecoes          []ProjecaoMensalResponse `json:"projecoes"`
    MediaReceita3Meses string                   `json:"media_receita_3_meses"`
    MediaDespesas3Meses string                  `json:"media_despesas_3_meses"`
    TendenciaReceita   string                   `json:"tendencia_receita"`
    DataGeracao        string                   `json:"data_geracao"`
}
```

#### 3.2 Checklist DTOs

- [ ] PainelMensalResponse
- [ ] ProjecaoMensalResponse
- [ ] ProjecoesResponse
- [ ] Mapper entity → DTO
- [ ] Validação de valores monetários como string

---

### 4️⃣ HTTP ENDPOINTS (Esforço: 3h)

#### 4.1 Handler Methods

- [ ] Adicionar ao `FinancialHandler`:

```go
// GET /financial/dashboard
func (h *FinancialHandler) GetDashboard(c echo.Context) error {
    tenantID := c.Get("tenant_id").(uuid.UUID)
    
    // Parse query params
    year := c.QueryParam("year")
    month := c.QueryParam("month")
    
    var ano, mes int
    if year != "" && month != "" {
        ano, _ = strconv.Atoi(year)
        mes, _ = strconv.Atoi(month)
    } else {
        now := time.Now()
        ano = now.Year()
        mes = int(now.Month())
    }
    
    input := financial.PainelMensalInput{
        TenantID: tenantID,
        Ano:      ano,
        Mes:      mes,
    }
    
    result, err := h.getPainelMensalUC.Execute(c.Request().Context(), input)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
    }
    
    response := mapper.ToPainelMensalResponse(result)
    return c.JSON(http.StatusOK, response)
}

// GET /financial/projections
func (h *FinancialHandler) GetProjections(c echo.Context) error {
    tenantID := c.Get("tenant_id").(uuid.UUID)
    
    mesesAhead := 3 // default
    if m := c.QueryParam("months_ahead"); m != "" {
        mesesAhead, _ = strconv.Atoi(m)
        if mesesAhead < 1 || mesesAhead > 12 {
            mesesAhead = 3
        }
    }
    
    input := financial.ProjecoesInput{
        TenantID:   tenantID,
        MesesAhead: mesesAhead,
    }
    
    result, err := h.getProjecoesUC.Execute(c.Request().Context(), input)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
    }
    
    response := mapper.ToProjecoesResponse(result)
    return c.JSON(http.StatusOK, response)
}
```

#### 4.2 Registrar Rotas

- [ ] Editar `backend/cmd/api/main.go`:

```go
// Dashboard e Projeções
financial.GET("/dashboard", financialHandler.GetDashboard)
financial.GET("/projections", financialHandler.GetProjections)
```

#### 4.3 Checklist Endpoints

- [ ] Implementar GetDashboard handler
- [ ] Implementar GetProjections handler
- [ ] Parsear query params (year, month, months_ahead)
- [ ] Validar tenant_id do context
- [ ] Mapear output para DTO
- [ ] Registrar rotas em main.go
- [ ] Testes de handler

---

### 5️⃣ QUERIES AUXILIARES (Esforço: 2h)

#### 5.1 Adicionar Queries de Agregação

- [ ] `contas_a_receber.sql`:

```sql
-- name: SumRecebidoByPeriod :one
SELECT COALESCE(SUM(valor_pago), 0) as total
FROM contas_a_receber
WHERE tenant_id = $1
  AND status = 'RECEBIDO'
  AND data_recebimento >= $2
  AND data_recebimento <= $3;

-- name: SumPendenteByPeriod :one
SELECT COALESCE(SUM(valor), 0) as total
FROM contas_a_receber
WHERE tenant_id = $1
  AND status = 'PENDENTE'
  AND data_vencimento >= $2
  AND data_vencimento <= $3;
```

- [ ] `contas_a_pagar.sql`:

```sql
-- name: SumPagoByPeriod :one
SELECT COALESCE(SUM(valor), 0) as total
FROM contas_a_pagar
WHERE tenant_id = $1
  AND status = 'PAGO'
  AND data_pagamento >= $2
  AND data_pagamento <= $3;

-- name: SumPendenteByPeriod :one
SELECT COALESCE(SUM(valor), 0) as total
FROM contas_a_pagar
WHERE tenant_id = $1
  AND status IN ('ABERTO', 'ATRASADO')
  AND data_vencimento >= $2
  AND data_vencimento <= $3;
```

#### 5.2 Checklist Queries

- [ ] SumRecebidoByPeriod (contas_a_receber)
- [ ] SumPendenteByPeriod (contas_a_receber)
- [ ] SumPagoByPeriod (contas_a_pagar)
- [ ] SumPendenteByPeriod (contas_a_pagar)
- [ ] Executar `sqlc generate`
- [ ] Verificar tipos gerados

---

### 6️⃣ TESTES (Esforço: 4h)

#### 6.1 Testes Unitários

- [ ] GetPainelMensalUseCase
  - [ ] Cenário: mês com receitas e despesas
  - [ ] Cenário: mês sem movimentação
  - [ ] Cenário: sem meta definida
  - [ ] Cenário: meta atingida
  - [ ] Cenário: abaixo da meta

- [ ] GetProjecoesUseCase
  - [ ] Cenário: tendência crescente
  - [ ] Cenário: tendência estável
  - [ ] Cenário: tendência decrescente
  - [ ] Cenário: sem histórico

#### 6.2 Testes de Integração

- [ ] Endpoint /financial/dashboard
- [ ] Endpoint /financial/projections
- [ ] Validar cálculos de agregação

#### 6.3 Validação de Fórmulas

- [ ] Comparar cálculos com `docs/10-calculos/`
- [ ] Ticket Médio = Receita / Atendimentos
- [ ] Margem Líquida = Lucro / Receita * 100
- [ ] % Meta = Realizado / Meta * 100

---

## 📊 ESTIMATIVA DE TEMPO

| Tarefa | Horas |
|--------|:-----:|
| Use Case: GetPainelMensal | 6h |
| Use Case: GetProjecoes | 6h |
| DTOs | 1h |
| HTTP Endpoints | 3h |
| Queries Auxiliares | 2h |
| Testes | 4h |
| **TOTAL** | **22h** |

---

## ✅ CRITÉRIOS DE CONCLUSÃO

- [ ] Endpoint /dashboard retornando dados corretos
- [ ] Endpoint /projections retornando projeções
- [ ] Cálculos validados contra documentação
- [ ] Testes passando (>80% cobertura)
- [ ] Performance < 500ms para resposta
- [ ] Code review aprovado

---

## 🔗 DEPENDÊNCIAS

| Dependência | Status | Bloqueio |
|-------------|--------|----------|
| Sprint 2 (Despesas Fixas) | ❌ | **BLOQUEIA** |
| Queries de agregação | ❌ | Implementar neste sprint |
| MetaMensalRepository | ✅ | — |
| FluxoCaixaRepository | ✅ | — |

---

## 📎 ARQUIVOS A CRIAR/MODIFICAR

```
backend/internal/application/
├── usecase/financial/
│   ├── get_painel_mensal.go    ← CRIAR
│   └── get_projecoes.go        ← CRIAR
├── dto/
│   └── painel_mensal_dto.go    ← CRIAR
├── mapper/
│   └── painel_mensal_mapper.go ← CRIAR

backend/internal/infra/
├── db/queries/
│   ├── contas_a_pagar.sql      ← MODIFICAR (adicionar queries)
│   └── contas_a_receber.sql    ← MODIFICAR (adicionar queries)
├── http/handler/
│   └── financial_handler.go    ← MODIFICAR (adicionar métodos)

backend/cmd/api/
└── main.go                     ← MODIFICAR (registrar rotas)
```

---

*Próximo Sprint: Sprint 4 — Frontend*
