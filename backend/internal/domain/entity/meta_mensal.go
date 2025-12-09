package entity

import (
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MetaMensal representa a meta de faturamento mensal da barbearia
type MetaMensal struct {
	ID       string
	TenantID uuid.UUID
	MesAno   valueobject.MesAno

	MetaFaturamento valueobject.Money
	Origem          valueobject.OrigemMeta // MANUAL ou AUTOMATICA
	Status          string                 // PENDENTE, ACEITA, REJEITADA

	// Campos de progresso (calculados)
	Realizado  valueobject.Money
	Percentual valueobject.Percentage

	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewMetaMensal cria uma nova meta mensal
func NewMetaMensal(tenantID uuid.UUID, mesAno valueobject.MesAno, metaFaturamento valueobject.Money, origem valueobject.OrigemMeta) (*MetaMensal, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if metaFaturamento.IsNegative() {
		return nil, domain.ErrMetaNegativa
	}
	if !origem.IsValid() {
		return nil, domain.ErrMetaInvalida
	}

	now := time.Now()
	return &MetaMensal{
		ID:              uuid.NewString(),
		TenantID:        tenantID,
		MesAno:          mesAno,
		MetaFaturamento: metaFaturamento,
		Origem:          origem,
		Status:          "PENDENTE", // Status inicial: PENDENTE, ACEITA ou REJEITADA
		Realizado:       valueobject.Zero(),
		Percentual:      valueobject.ZeroPercent(),
		CriadoEm:        now,
		AtualizadoEm:    now,
	}, nil
}

// CalcularProgresso calcula o percentual realizado
func (m *MetaMensal) CalcularProgresso(realizado valueobject.Money) {
	m.Realizado = realizado
	if m.MetaFaturamento.IsPositive() {
		percentual := realizado.Value().
			Div(m.MetaFaturamento.Value()).
			Mul(decimal.NewFromInt(100))

		// Limita a 200% para evitar valores absurdos
		if percentual.GreaterThan(decimal.NewFromInt(200)) {
			percentual = decimal.NewFromInt(200)
		}
		m.Percentual = valueobject.NewPercentageUnsafe(percentual)
	}
	m.AtualizadoEm = time.Now()
}

// Rejeitar rejeita a meta
func (m *MetaMensal) Rejeitar() {
	m.Status = "REJEITADA"
	m.AtualizadoEm = time.Now()
}

// Aceitar aceita/ativa a meta
func (m *MetaMensal) Aceitar() {
	m.Status = "ACEITA"
	m.AtualizadoEm = time.Now()
}

// AtualizarMeta atualiza o valor da meta
func (m *MetaMensal) AtualizarMeta(novaMeta valueobject.Money) error {
	if novaMeta.IsNegative() {
		return domain.ErrMetaNegativa
	}
	m.MetaFaturamento = novaMeta
	m.Origem = valueobject.OrigemMetaManual
	m.AtualizadoEm = time.Now()
	return nil
}

// Validate valida as regras de negócio
func (m *MetaMensal) Validate() error {
	if m.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if m.MesAno.String() == "" {
		return domain.ErrMesAnoRequired
	}
	if m.MetaFaturamento.IsNegative() {
		return domain.ErrMetaNegativa
	}
	if !m.Origem.IsValid() {
		return domain.ErrMetaInvalida
	}
	return nil
}

// ForaDoAlvo verifica se está fora do alvo (menos de 80%)
func (m *MetaMensal) ForaDoAlvo() bool {
	threshold, _ := valueobject.NewPercentage(decimal.NewFromInt(80))
	return m.Percentual.LessThan(threshold)
}

// Atingiu verifica se atingiu a meta (100% ou mais)
func (m *MetaMensal) Atingiu() bool {
	threshold, _ := valueobject.NewPercentage(decimal.NewFromInt(100))
	return m.Percentual.GreaterThan(threshold) || m.Percentual.Equals(threshold)
}
