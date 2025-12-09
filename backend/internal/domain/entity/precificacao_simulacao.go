package entity

import (
	"encoding/json"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

// PrecificacaoSimulacao representa uma simulação de precificação
type PrecificacaoSimulacao struct {
	ID       string
	TenantID uuid.UUID

	ItemID   string
	TipoItem string // SERVICO ou PRODUTO

	CustoMateriais valueobject.Money
	CustoMaoDeObra valueobject.Money
	CustoTotal     valueobject.Money

	MargemDesejada     valueobject.Percentage
	ComissaoPercentual valueobject.Percentage
	ImpostoPercentual  valueobject.Percentage

	PrecoSugerido       valueobject.Money
	PrecoAtual          valueobject.Money
	DiferencaPercentual valueobject.Percentage

	LucroEstimado valueobject.Money
	MargemFinal   valueobject.Percentage

	ParametrosJSON string // JSON com parâmetros adicionais

	CriadoEm time.Time
}

// ParametrosSimulacao representa parâmetros extras da simulação
type ParametrosSimulacao struct {
	TempoMedioMinutos     int     `json:"tempo_medio_minutos,omitempty"`
	QuantidadeMensal      int     `json:"quantidade_mensal,omitempty"`
	CustoPorMinuto        float64 `json:"custo_por_minuto,omitempty"`
	ObservacoesAdicionais string  `json:"observacoes,omitempty"`
}

// NewPrecificacaoSimulacao cria uma nova simulação
func NewPrecificacaoSimulacao(
	tenantID uuid.UUID, itemID, tipoItem string,
	custoMateriais, custoMaoDeObra valueobject.Money,
	margemDesejada, comissao, imposto valueobject.Percentage,
	precoAtual valueobject.Money,
	parametros *ParametrosSimulacao,
) (*PrecificacaoSimulacao, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if itemID == "" {
		return nil, domain.ErrInvalidID
	}
	if tipoItem != "SERVICO" && tipoItem != "PRODUTO" {
		return nil, domain.ErrMetaInvalida
	}

	custoTotal := custoMateriais.Add(custoMaoDeObra)

	// Serializar parâmetros
	var parametrosJSON string
	if parametros != nil {
		jsonBytes, err := json.Marshal(parametros)
		if err != nil {
			parametrosJSON = "{}"
		} else {
			parametrosJSON = string(jsonBytes)
		}
	} else {
		parametrosJSON = "{}"
	}

	sim := &PrecificacaoSimulacao{
		ID:                 uuid.NewString(),
		TenantID:           tenantID,
		ItemID:             itemID,
		TipoItem:           tipoItem,
		CustoMateriais:     custoMateriais,
		CustoMaoDeObra:     custoMaoDeObra,
		CustoTotal:         custoTotal,
		MargemDesejada:     margemDesejada,
		ComissaoPercentual: comissao,
		ImpostoPercentual:  imposto,
		PrecoAtual:         precoAtual,
		ParametrosJSON:     parametrosJSON,
		CriadoEm:           time.Now(),
	}

	sim.CalcularPrecoSugerido()
	return sim, nil
}

// CalcularPrecoSugerido calcula o preço sugerido com base nos custos e margem
func (p *PrecificacaoSimulacao) CalcularPrecoSugerido() {
	// Fórmula: Preço = CustoTotal / (1 - (Margem + Comissão + Imposto) / 100)

	totalDescontos := p.MargemDesejada.Value().
		Add(p.ComissaoPercentual.Value()).
		Add(p.ImpostoPercentual.Value())

	divisor := valueobject.HundredPercent().Value().Sub(totalDescontos)

	// Proteger contra divisão por zero ou valores inválidos
	if divisor.LessThanOrEqual(valueobject.ZeroPercent().Value()) {
		p.PrecoSugerido = valueobject.Zero()
		p.LucroEstimado = valueobject.Zero()
		p.MargemFinal = valueobject.ZeroPercent()
		return
	}

	precoCalculado := p.CustoTotal.Value().
		Mul(valueobject.HundredPercent().Value()).
		Div(divisor)

	p.PrecoSugerido = valueobject.NewMoneyFromDecimal(precoCalculado)

	// Calcular diferença percentual com preço atual
	if p.PrecoAtual.IsPositive() {
		diff := p.PrecoSugerido.Sub(p.PrecoAtual).Value().
			Div(p.PrecoAtual.Value()).
			Mul(valueobject.HundredPercent().Value())

		// Limitar entre -100% e +1000%
		if diff.LessThan(valueobject.HundredPercent().Value().Neg()) {
			diff = valueobject.HundredPercent().Value().Neg()
		}
		if diff.GreaterThan(valueobject.HundredPercent().Value().Mul(valueobject.NewPercentageUnsafe(valueobject.HundredPercent().Value().Mul(valueobject.HundredPercent().Value())).Value())) {
			diff = valueobject.HundredPercent().Value().Mul(valueobject.NewPercentageUnsafe(valueobject.HundredPercent().Value()).Value()).Mul(valueobject.NewPercentageUnsafe(valueobject.HundredPercent().Value()).Value())
		}

		p.DiferencaPercentual = valueobject.NewPercentageUnsafe(diff)
	}

	// Calcular lucro estimado
	comissaoValor := p.PrecoSugerido.Percentage(p.ComissaoPercentual)
	impostoValor := p.PrecoSugerido.Percentage(p.ImpostoPercentual)

	p.LucroEstimado = p.PrecoSugerido.
		Sub(p.CustoTotal).
		Sub(comissaoValor).
		Sub(impostoValor)

	// Calcular margem final
	if p.PrecoSugerido.IsPositive() {
		margemFinalCalc := p.LucroEstimado.Value().
			Div(p.PrecoSugerido.Value()).
			Mul(valueobject.HundredPercent().Value())
		p.MargemFinal = valueobject.NewPercentageUnsafe(margemFinalCalc)
	}
}

// Validate valida as regras de negócio
func (p *PrecificacaoSimulacao) Validate() error {
	if p.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if p.ItemID == "" {
		return domain.ErrInvalidID
	}
	if p.TipoItem != "SERVICO" && p.TipoItem != "PRODUTO" {
		return domain.ErrMetaInvalida
	}
	if p.CustoTotal.IsNegative() {
		return domain.ErrValorNegativo
	}
	return nil
}

// GetParametros deserializa os parâmetros JSON
func (p *PrecificacaoSimulacao) GetParametros() (*ParametrosSimulacao, error) {
	var params ParametrosSimulacao
	if err := json.Unmarshal([]byte(p.ParametrosJSON), &params); err != nil {
		return nil, err
	}
	return &params, nil
}
