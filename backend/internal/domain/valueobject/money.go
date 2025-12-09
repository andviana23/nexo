package valueobject

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Money representa um valor monetário em centavos (BRL)
// Usa decimal.Decimal para evitar erros de arredondamento
type Money struct {
	value decimal.Decimal
}

// NewMoney cria um novo valor monetário a partir de centavos
func NewMoney(centavos int64) Money {
	return Money{
		value: decimal.NewFromInt(centavos).Div(decimal.NewFromInt(100)),
	}
}

// NewMoneyFromDecimal cria um valor monetário a partir de um decimal
func NewMoneyFromDecimal(valor decimal.Decimal) Money {
	return Money{value: valor}
}

// NewMoneyFromFloat cria um valor monetário a partir de float64
// ATENÇÃO: Use apenas quando absolutamente necessário (ex: parse de JSON)
func NewMoneyFromFloat(valor float64) Money {
	return Money{
		value: decimal.NewFromFloat(valor),
	}
}

// Value retorna o valor como decimal.Decimal
func (m Money) Value() decimal.Decimal {
	return m.value
}

// Centavos retorna o valor em centavos (int64)
func (m Money) Centavos() int64 {
	return m.value.Mul(decimal.NewFromInt(100)).IntPart()
}

// String formata o valor como string (ex: "R$ 10.50")
func (m Money) String() string {
	return fmt.Sprintf("R$ %s", m.value.StringFixed(2))
}

// Raw retorna o valor como string numérica sem formatação (ex: "10.50")
// Usado em DTOs para transferência de dados
func (m Money) Raw() string {
	return m.value.StringFixed(2)
}

// IsPositive verifica se o valor é positivo
func (m Money) IsPositive() bool {
	return m.value.GreaterThan(decimal.Zero)
}

// IsZero verifica se o valor é zero
func (m Money) IsZero() bool {
	return m.value.IsZero()
}

// IsNegative verifica se o valor é negativo
func (m Money) IsNegative() bool {
	return m.value.LessThan(decimal.Zero)
}

// Add adiciona outro Money
func (m Money) Add(other Money) Money {
	return Money{value: m.value.Add(other.value)}
}

// Sub subtrai outro Money
func (m Money) Sub(other Money) Money {
	return Money{value: m.value.Sub(other.value)}
}

// Mul multiplica por um decimal
func (m Money) Mul(multiplicador decimal.Decimal) Money {
	return Money{value: m.value.Mul(multiplicador)}
}

// Div divide por um decimal
func (m Money) Div(divisor decimal.Decimal) Money {
	if divisor.IsZero() {
		return Money{value: decimal.Zero}
	}
	return Money{value: m.value.Div(divisor)}
}

// Percentage calcula a porcentagem do valor
func (m Money) Percentage(percent Percentage) Money {
	return Money{
		value: m.value.Mul(percent.Value()).Div(decimal.NewFromInt(100)),
	}
}

// Equals verifica igualdade
func (m Money) Equals(other Money) bool {
	return m.value.Equal(other.value)
}

// GreaterThan verifica se é maior que outro Money
func (m Money) GreaterThan(other Money) bool {
	return m.value.GreaterThan(other.value)
}

// LessThan verifica se é menor que outro Money
func (m Money) LessThan(other Money) bool {
	return m.value.LessThan(other.value)
}

// Zero retorna um Money com valor zero
func Zero() Money {
	return Money{value: decimal.Zero}
}
