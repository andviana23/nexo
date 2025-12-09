package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

// CompensacaoBancaria representa uma compensação bancária (D+)
type CompensacaoBancaria struct {
	ID        string
	TenantID  uuid.UUID
	ReceitaID string

	DataTransacao   time.Time
	DataCompensacao time.Time
	DataCompensado  *time.Time

	ValorBruto     valueobject.Money
	TaxaPercentual valueobject.Percentage
	TaxaFixa       valueobject.Money
	ValorLiquido   valueobject.Money

	MeioPagamentoID string
	DMais           valueobject.DMais

	Status valueobject.StatusCompensacao

	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewCompensacaoBancaria cria uma nova compensação bancária
func NewCompensacaoBancaria(
	tenantID uuid.UUID, receitaID, meioPagamentoID string,
	dataTransacao time.Time,
	valorBruto valueobject.Money,
	taxaPercentual valueobject.Percentage,
	taxaFixa valueobject.Money,
	dMais valueobject.DMais,
) (*CompensacaoBancaria, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if receitaID == "" {
		return nil, domain.ErrInvalidID
	}
	if valorBruto.IsNegative() {
		return nil, domain.ErrValorNegativo
	}

	now := time.Now()
	dataCompensacao := dMais.CalcularDataCompensacao(dataTransacao)

	comp := &CompensacaoBancaria{
		ID:              uuid.NewString(),
		TenantID:        tenantID,
		ReceitaID:       receitaID,
		DataTransacao:   dataTransacao,
		DataCompensacao: dataCompensacao,
		ValorBruto:      valorBruto,
		TaxaPercentual:  taxaPercentual,
		TaxaFixa:        taxaFixa,
		MeioPagamentoID: meioPagamentoID,
		DMais:           dMais,
		Status:          valueobject.StatusCompensacaoPrevisto,
		CriadoEm:        now,
		AtualizadoEm:    now,
	}

	comp.CalcularValorLiquido()
	return comp, nil
}

// CalcularValorLiquido calcula o valor líquido após taxas
func (c *CompensacaoBancaria) CalcularValorLiquido() {
	taxaPerc := c.ValorBruto.Percentage(c.TaxaPercentual)
	c.ValorLiquido = c.ValorBruto.Sub(taxaPerc).Sub(c.TaxaFixa)
	c.AtualizadoEm = time.Now()
}

// MarcarComoConfirmado marca a compensação como confirmada
func (c *CompensacaoBancaria) MarcarComoConfirmado() error {
	if c.Status == valueobject.StatusCompensacaoCompensado {
		return domain.ErrCompensacaoJaCompensada
	}
	c.Status = valueobject.StatusCompensacaoConfirmado
	c.AtualizadoEm = time.Now()
	return nil
}

// MarcarComoCompensado marca a compensação como compensada
func (c *CompensacaoBancaria) MarcarComoCompensado() error {
	if c.Status == valueobject.StatusCompensacaoCompensado {
		return domain.ErrCompensacaoJaCompensada
	}
	now := time.Now()
	c.DataCompensado = &now
	c.Status = valueobject.StatusCompensacaoCompensado
	c.AtualizadoEm = now
	return nil
}

// Cancelar cancela a compensação
func (c *CompensacaoBancaria) Cancelar() {
	c.Status = valueobject.StatusCompensacaoCancelado
	c.AtualizadoEm = time.Now()
}

// Validate valida as regras de negócio
func (c *CompensacaoBancaria) Validate() error {
	if c.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if c.ReceitaID == "" {
		return domain.ErrInvalidID
	}
	if !c.Status.IsValid() {
		return domain.ErrStatusInvalido
	}
	if c.ValorBruto.IsNegative() {
		return domain.ErrValorNegativo
	}
	if c.DataTransacao.IsZero() || c.DataCompensacao.IsZero() {
		return domain.ErrDataCompensacaoInvalida
	}
	return nil
}

// JaCompensado verifica se já foi compensado
func (c *CompensacaoBancaria) JaCompensado() bool {
	return c.Status == valueobject.StatusCompensacaoCompensado
}

// PodeSerCompensado verifica se pode ser compensado (data chegou)
func (c *CompensacaoBancaria) PodeSerCompensado() bool {
	if c.JaCompensado() {
		return false
	}
	return time.Now().After(c.DataCompensacao) || time.Now().Equal(c.DataCompensacao)
}
