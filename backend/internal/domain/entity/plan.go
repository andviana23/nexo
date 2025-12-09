package entity

import (
	"strings"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Plan representa o modelo de plano de assinatura (template interno)
type Plan struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	Nome            string
	Descricao       *string
	Valor           decimal.Decimal
	Periodicidade   string
	QtdServicos     *int
	LimiteUsoMensal *int
	Ativo           bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewPlan cria um plano validando regras básicas de domínio.
func NewPlan(tenantID uuid.UUID, nome string, valor decimal.Decimal) (*Plan, error) {
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return nil, domain.ErrPlanNameRequired
	}
	if len(nome) < 3 {
		return nil, domain.ErrPlanNameTooShort
	}
	if len(nome) > 100 {
		return nil, domain.ErrPlanNameTooLong
	}

	if valor.LessThanOrEqual(decimal.Zero) {
		return nil, domain.ErrPlanValueInvalid
	}

	now := time.Now()
	return &Plan{
		ID:            uuid.New(),
		TenantID:      tenantID,
		Nome:          nome,
		Valor:         valor,
		Periodicidade: "MENSAL",
		Ativo:         true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// CanBeDeleted aplica a regra PL-003 (não excluir com assinaturas ativas).
func (p *Plan) CanBeDeleted(activeSubscriptionsCount int) bool {
	return activeSubscriptionsCount == 0
}

// IsAvailableForSelection aplica a regra PL-002 (só planos ativos disponíveis).
func (p *Plan) IsAvailableForSelection() bool {
	return p.Ativo
}

// UpdateBasicFields atualiza campos editáveis com validação.
func (p *Plan) UpdateBasicFields(nome string, descricao *string, valor decimal.Decimal, qtdServicos, limiteUsoMensal *int, ativo bool) error {
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return domain.ErrPlanNameRequired
	}
	if len(nome) < 3 {
		return domain.ErrPlanNameTooShort
	}
	if len(nome) > 100 {
		return domain.ErrPlanNameTooLong
	}
	if valor.LessThanOrEqual(decimal.Zero) {
		return domain.ErrPlanValueInvalid
	}

	p.Nome = nome
	p.Descricao = descricao
	p.Valor = valor
	p.QtdServicos = qtdServicos
	p.LimiteUsoMensal = limiteUsoMensal
	p.Ativo = ativo
	p.UpdatedAt = time.Now()
	return nil
}
