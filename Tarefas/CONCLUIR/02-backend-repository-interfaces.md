# 02 - Backend: Repository Interfaces (BLOQUEADOR)

**Prioridade:** 🔴 CRÍTICA
**Estimativa:** 2 dias
**Dependências:** 01-backend-domain-entities.md
**Bloqueia:** 03-backend-repository-implementations.md

---

## Objetivo

Criar interfaces de repositório (ports) para todas as entidades novas seguindo Clean Architecture.

---

## Estrutura

**Local:** `backend/internal/domain/repository/`

Cada interface deve seguir o padrão:

```go
type XxxRepository interface {
    Save(ctx context.Context, tenantID string, entity *Xxx) error
    FindByID(ctx context.Context, tenantID, id string) (*Xxx, error)
    FindByTenantID(ctx context.Context, tenantID string, filters ...Filter) ([]*Xxx, error)
    Update(ctx context.Context, tenantID string, entity *Xxx) error
    Delete(ctx context.Context, tenantID, id string) error
}
```

---

## Interfaces a Criar

### 1. DREMensalRepository

```go
// backend/internal/domain/repository/dre_repository.go
package repository

type DREMensalRepository interface {
    Save(ctx context.Context, tenantID string, dre *entity.DREMensal) error
    FindByID(ctx context.Context, tenantID, id string) (*entity.DREMensal, error)
    FindByTenantAndMesAno(ctx context.Context, tenantID, mesAno string) (*entity.DREMensal, error)
    FindByTenantAndPeriod(ctx context.Context, tenantID string, inicio, fim string) ([]*entity.DREMensal, error)
    Update(ctx context.Context, tenantID string, dre *entity.DREMensal) error
    Delete(ctx context.Context, tenantID, id string) error
}
```

### 2. FluxoCaixaDiarioRepository

```go
type FluxoCaixaDiarioRepository interface {
    Save(ctx context.Context, tenantID string, fluxo *entity.FluxoCaixaDiario) error
    FindByID(ctx context.Context, tenantID, id string) (*entity.FluxoCaixaDiario, error)
    FindByTenantAndDate(ctx context.Context, tenantID string, data time.Time) (*entity.FluxoCaixaDiario, error)
    FindByTenantAndPeriod(ctx context.Context, tenantID string, inicio, fim time.Time) ([]*entity.FluxoCaixaDiario, error)
    Update(ctx context.Context, tenantID string, fluxo *entity.FluxoCaixaDiario) error
    Delete(ctx context.Context, tenantID, id string) error
}
```

### 3. CompensacaoBancariaRepository

```go
type CompensacaoBancariaRepository interface {
    Save(ctx context.Context, tenantID string, comp *entity.CompensacaoBancaria) error
    FindByID(ctx context.Context, tenantID, id string) (*entity.CompensacaoBancaria, error)
    FindByReceitaID(ctx context.Context, tenantID, receitaID string) (*entity.CompensacaoBancaria, error)
    FindPendentesByTenantAndDataCompensacao(ctx context.Context, tenantID string, data time.Time) ([]*entity.CompensacaoBancaria, error)
    SumByTenantDateAndStatus(ctx context.Context, tenantID string, data time.Time, status entity.CompensacaoStatus) (decimal.Decimal, error)
    Update(ctx context.Context, tenantID string, comp *entity.CompensacaoBancaria) error
    Delete(ctx context.Context, tenantID, id string) error
}
```

### 4. MetaMensalRepository

```go
type MetaMensalRepository interface {
    Save(ctx context.Context, tenantID string, meta *entity.MetaMensal) error
    FindByID(ctx context.Context, tenantID, id string) (*entity.MetaMensal, error)
    FindByTenantAndMesAno(ctx context.Context, tenantID, mesAno string) (*entity.MetaMensal, error)
    FindCurrentByTenant(ctx context.Context, tenantID string) (*entity.MetaMensal, error)
    Update(ctx context.Context, tenantID string, meta *entity.MetaMensal) error
    Delete(ctx context.Context, tenantID, id string) error
}
```

### 5. MetaBarbeiroRepository

```go
type MetaBarbeiroRepository interface {
    Save(ctx context.Context, tenantID string, meta *entity.MetaBarbeiro) error
    FindByID(ctx context.Context, tenantID, id string) (*entity.MetaBarbeiro, error)
    FindByTenantBarbeiroAndMesAno(ctx context.Context, tenantID, barbeiroID, mesAno string) (*entity.MetaBarbeiro, error)
    FindAllByTenantAndMesAno(ctx context.Context, tenantID, mesAno string) ([]*entity.MetaBarbeiro, error)
    Update(ctx context.Context, tenantID string, meta *entity.MetaBarbeiro) error
    Delete(ctx context.Context, tenantID, id string) error
}
```

### 6. MetaTicketMedioRepository

```go
type MetaTicketMedioRepository interface {
    Save(ctx context.Context, tenantID string, meta *entity.MetaTicketMedio) error
    FindByID(ctx context.Context, tenantID, id string) (*entity.MetaTicketMedio, error)
    FindByTenantAndMesAno(ctx context.Context, tenantID, mesAno string, tipo string) (*entity.MetaTicketMedio, error)
    Update(ctx context.Context, tenantID string, meta *entity.MetaTicketMedio) error
    Delete(ctx context.Context, tenantID, id string) error
}
```

### 7. PrecificacaoConfigRepository

```go
type PrecificacaoConfigRepository interface {
    Save(ctx context.Context, tenantID string, config *entity.PrecificacaoConfig) error
    FindByTenantID(ctx context.Context, tenantID string) (*entity.PrecificacaoConfig, error)
    Update(ctx context.Context, tenantID string, config *entity.PrecificacaoConfig) error
}
```

### 8. PrecificacaoSimulacaoRepository

```go
type PrecificacaoSimulacaoRepository interface {
    Save(ctx context.Context, tenantID string, sim *entity.PrecificacaoSimulacao) error
    FindByID(ctx context.Context, tenantID, id string) (*entity.PrecificacaoSimulacao, error)
    FindByTenantAndItemID(ctx context.Context, tenantID, itemID string, limit int) ([]*entity.PrecificacaoSimulacao, error)
    FindRecent(ctx context.Context, tenantID string, limit int) ([]*entity.PrecificacaoSimulacao, error)
}
```

### 9. ContaAPagarRepository

```go
type ContaAPagarRepository interface {
    Save(ctx context.Context, tenantID string, conta *entity.ContaAPagar) error
    FindByID(ctx context.Context, tenantID, id string) (*entity.ContaAPagar, error)
    FindByTenantID(ctx context.Context, tenantID string, filters ContaPagarFilters) ([]*entity.ContaAPagar, error)
    FindVencendoEm(ctx context.Context, tenantID string, dias int) ([]*entity.ContaAPagar, error)
    SumByTenantAndDataVencimento(ctx context.Context, tenantID string, data time.Time) (decimal.Decimal, error)
    Update(ctx context.Context, tenantID string, conta *entity.ContaAPagar) error
    Delete(ctx context.Context, tenantID, id string) error
}
```

### 10. ContaAReceberRepository

```go
type ContaAReceberRepository interface {
    Save(ctx context.Context, tenantID string, conta *entity.ContaAReceber) error
    FindByID(ctx context.Context, tenantID, id string) (*entity.ContaAReceber, error)
    FindByAssinaturaID(ctx context.Context, tenantID, assinaturaID string) ([]*entity.ContaAReceber, error)
    FindByTenantID(ctx context.Context, tenantID string, filters ContaReceberFilters) ([]*entity.ContaAReceber, error)
    FindInadimplentes(ctx context.Context, tenantID string) ([]*entity.ContaAReceber, error)
    Update(ctx context.Context, tenantID string, conta *entity.ContaAReceber) error
    Delete(ctx context.Context, tenantID, id string) error
}
```

---

## Estender Repositórios Existentes

### Estender ReceitaRepository

Adicionar métodos de agregação:

```go
SumByTenantAndPeriod(ctx context.Context, tenantID string, inicio, fim time.Time) (decimal.Decimal, error)
SumByTenantPeriodAndSubtipo(ctx context.Context, tenantID string, inicio, fim time.Time, subtipo string) (decimal.Decimal, error)
CountByTenantAndPeriod(ctx context.Context, tenantID string, inicio, fim time.Time) (int64, error)
```

### Estender DespesaRepository

Adicionar métodos:

```go
SumByTenantAndPeriod(ctx context.Context, tenantID string, inicio, fim time.Time) (decimal.Decimal, error)
SumByTenantPeriodAndTipoCusto(ctx context.Context, tenantID string, inicio, fim time.Time, tipoCusto string) (decimal.Decimal, error)
```

### Estender MeioPagamentoRepository

Adicionar:

```go
FindByIDWithDMais(ctx context.Context, tenantID, id string) (*entity.MeioPagamento, error)
```

---

## Critérios de Aceite

- [ ] Todas as 10 interfaces criadas
- [ ] Métodos existentes estendidos (Receita, Despesa, MeioPagamento)
- [ ] Cada interface documenta parâmetros e retornos (GoDoc)
- [ ] Nenhuma lógica de negócio nas interfaces (apenas contratos)
- [ ] Context como primeiro parâmetro sempre
- [ ] TenantID sempre presente para garantir multi-tenancy
- [ ] Métodos de agregação retornam `decimal.Decimal` ou tipos primitivos
- [ ] Erros customizados definidos (ErrNotFound, etc)

---

## Arquivos a Criar

```
backend/internal/domain/repository/
├── dre_repository.go
├── fluxo_caixa_repository.go
├── compensacao_repository.go
├── meta_mensal_repository.go
├── meta_barbeiro_repository.go
├── meta_ticket_repository.go
├── precificacao_config_repository.go
├── precificacao_simulacao_repository.go
├── conta_a_pagar_repository.go
└── conta_a_receber_repository.go
```

---

## Observações

- Interfaces devem ser agnósticas de implementação (PostgreSQL/SQLite/etc)
- Não usar tipos específicos de DB (sql.Row, pgx.Conn, etc)
- Retornar sempre entidades de domínio, nunca DTOs
- Filters podem ser structs específicas (ContaPagarFilters, etc)
