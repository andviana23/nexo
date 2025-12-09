package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PrecificacaoConfig representa a configuração de precificação do tenant
type PrecificacaoConfig struct {
	ID       string
	TenantID uuid.UUID

	MargemDesejada            valueobject.Percentage // 5-100%
	MarkupAlvo                decimal.Decimal        // >= 1.0
	ImpostoPercentual         valueobject.Percentage
	ComissaoPercentualDefault valueobject.Percentage

	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewPrecificacaoConfig cria uma nova configuração de precificação
func NewPrecificacaoConfig(
	tenantID uuid.UUID,
	margemDesejada, impostoPercentual, comissaoDefault valueobject.Percentage,
	markupAlvo decimal.Decimal,
) (*PrecificacaoConfig, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}

	// Margem deve estar entre 5-100%
	cincoPercent, _ := valueobject.NewPercentage(decimal.NewFromInt(5))
	if margemDesejada.LessThan(cincoPercent) {
		return nil, domain.ErrMargemInvalida
	}

	// Markup deve ser >= 1.0
	if markupAlvo.LessThan(decimal.NewFromInt(1)) {
		return nil, domain.ErrMarkupInvalido
	}

	now := time.Now()
	return &PrecificacaoConfig{
		ID:                        uuid.NewString(),
		TenantID:                  tenantID,
		MargemDesejada:            margemDesejada,
		MarkupAlvo:                markupAlvo,
		ImpostoPercentual:         impostoPercentual,
		ComissaoPercentualDefault: comissaoDefault,
		CriadoEm:                  now,
		AtualizadoEm:              now,
	}, nil
}

// Validate valida as regras de negócio
func (p *PrecificacaoConfig) Validate() error {
	if p.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}

	cincoPercent, _ := valueobject.NewPercentage(decimal.NewFromInt(5))
	if p.MargemDesejada.LessThan(cincoPercent) {
		return domain.ErrMargemInvalida
	}

	if p.MarkupAlvo.LessThan(decimal.NewFromInt(1)) {
		return domain.ErrMarkupInvalido
	}

	return nil
}

// AtualizarMargem atualiza a margem desejada
func (p *PrecificacaoConfig) AtualizarMargem(novaMargem valueobject.Percentage) error {
	cincoPercent, _ := valueobject.NewPercentage(decimal.NewFromInt(5))
	if novaMargem.LessThan(cincoPercent) {
		return domain.ErrMargemInvalida
	}
	p.MargemDesejada = novaMargem
	p.AtualizadoEm = time.Now()
	return nil
}

// AtualizarMarkup atualiza o markup alvo
func (p *PrecificacaoConfig) AtualizarMarkup(novoMarkup decimal.Decimal) error {
	if novoMarkup.LessThan(decimal.NewFromInt(1)) {
		return domain.ErrMarkupInvalido
	}
	p.MarkupAlvo = novoMarkup
	p.AtualizadoEm = time.Now()
	return nil
}

// AtualizarImposto atualiza o percentual de imposto
func (p *PrecificacaoConfig) AtualizarImposto(novoImposto valueobject.Percentage) {
	p.ImpostoPercentual = novoImposto
	p.AtualizadoEm = time.Now()
}

// AtualizarComissaoDefault atualiza a comissão padrão
func (p *PrecificacaoConfig) AtualizarComissaoDefault(novaComissao valueobject.Percentage) {
	p.ComissaoPercentualDefault = novaComissao
	p.AtualizadoEm = time.Now()
}
