package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

// DREMensal representa o Demonstrativo de Resultado do Exercício mensal
type DREMensal struct {
	ID       string
	TenantID uuid.UUID
	MesAno   valueobject.MesAno

	// Receitas
	ReceitaServicos valueobject.Money
	ReceitaProdutos valueobject.Money
	ReceitaPlanos   valueobject.Money
	ReceitaTotal    valueobject.Money

	// Custos Variáveis
	CustoComissoes     valueobject.Money
	CustoInsumos       valueobject.Money
	CustoVariavelTotal valueobject.Money

	// Despesas
	DespesaFixa     valueobject.Money
	DespesaVariavel valueobject.Money
	DespesaTotal    valueobject.Money

	// Resultado
	ResultadoBruto       valueobject.Money
	ResultadoOperacional valueobject.Money
	MargemBruta          valueobject.Percentage
	MargemOperacional    valueobject.Percentage
	LucroLiquido         valueobject.Money

	ProcessadoEm time.Time
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewDREMensal cria um novo DRE mensal
func NewDREMensal(tenantID uuid.UUID, mesAno valueobject.MesAno) (*DREMensal, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}

	now := time.Now()
	return &DREMensal{
		ID:                   uuid.NewString(),
		TenantID:             tenantID,
		MesAno:               mesAno,
		ReceitaServicos:      valueobject.Zero(),
		ReceitaProdutos:      valueobject.Zero(),
		ReceitaPlanos:        valueobject.Zero(),
		ReceitaTotal:         valueobject.Zero(),
		CustoComissoes:       valueobject.Zero(),
		CustoInsumos:         valueobject.Zero(),
		CustoVariavelTotal:   valueobject.Zero(),
		DespesaFixa:          valueobject.Zero(),
		DespesaVariavel:      valueobject.Zero(),
		DespesaTotal:         valueobject.Zero(),
		ResultadoBruto:       valueobject.Zero(),
		ResultadoOperacional: valueobject.Zero(),
		MargemBruta:          valueobject.ZeroPercent(),
		MargemOperacional:    valueobject.ZeroPercent(),
		LucroLiquido:         valueobject.Zero(),
		CriadoEm:             now,
		AtualizadoEm:         now,
	}, nil
}

// Calcular recalcula todos os valores do DRE
func (d *DREMensal) Calcular() {
	// Receita Total
	d.ReceitaTotal = d.ReceitaServicos.
		Add(d.ReceitaProdutos).
		Add(d.ReceitaPlanos)

	// Custo Variável Total
	d.CustoVariavelTotal = d.CustoComissoes.Add(d.CustoInsumos)

	// Despesa Total
	d.DespesaTotal = d.DespesaFixa.Add(d.DespesaVariavel)

	// Resultado Bruto = Receita - Custo Variável
	d.ResultadoBruto = d.ReceitaTotal.Sub(d.CustoVariavelTotal)

	// Resultado Operacional = Bruto - Despesas
	d.ResultadoOperacional = d.ResultadoBruto.Sub(d.DespesaTotal)

	// Lucro Líquido = Resultado Operacional
	d.LucroLiquido = d.ResultadoOperacional

	// Margens (%)
	if d.ReceitaTotal.IsPositive() {
		margemBrutaDecimal := d.ResultadoBruto.Value().
			Div(d.ReceitaTotal.Value()).
			Mul(valueobject.NewPercentageUnsafe(valueobject.HundredPercent().Value()).Value())
		d.MargemBruta = valueobject.NewPercentageUnsafe(margemBrutaDecimal)

		margemOpDecimal := d.ResultadoOperacional.Value().
			Div(d.ReceitaTotal.Value()).
			Mul(valueobject.HundredPercent().Value())
		d.MargemOperacional = valueobject.NewPercentageUnsafe(margemOpDecimal)
	}

	d.ProcessadoEm = time.Now()
	d.AtualizadoEm = time.Now()
}

// Validate valida as regras de negócio do DRE
func (d *DREMensal) Validate() error {
	if d.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if d.MesAno.String() == "" {
		return domain.ErrMesAnoRequired
	}
	// Valores não podem ser negativos (exceto resultados)
	if d.ReceitaServicos.IsNegative() || d.ReceitaProdutos.IsNegative() || d.ReceitaPlanos.IsNegative() {
		return domain.ErrValorNegativo
	}
	if d.CustoComissoes.IsNegative() || d.CustoInsumos.IsNegative() {
		return domain.ErrValorNegativo
	}
	if d.DespesaFixa.IsNegative() || d.DespesaVariavel.IsNegative() {
		return domain.ErrValorNegativo
	}
	return nil
}

// SetReceitas define os valores de receita
func (d *DREMensal) SetReceitas(servicos, produtos, planos valueobject.Money) {
	d.ReceitaServicos = servicos
	d.ReceitaProdutos = produtos
	d.ReceitaPlanos = planos
	d.AtualizadoEm = time.Now()
}

// SetCustosVariaveis define os custos variáveis
func (d *DREMensal) SetCustosVariaveis(comissoes, insumos valueobject.Money) {
	d.CustoComissoes = comissoes
	d.CustoInsumos = insumos
	d.AtualizadoEm = time.Now()
}

// SetDespesas define as despesas
func (d *DREMensal) SetDespesas(fixa, variavel valueobject.Money) {
	d.DespesaFixa = fixa
	d.DespesaVariavel = variavel
	d.AtualizadoEm = time.Now()
}
