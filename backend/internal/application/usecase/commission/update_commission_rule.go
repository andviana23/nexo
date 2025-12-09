package commission

import (
	"context"

	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
	"github.com/shopspring/decimal"
)

// UpdateCommissionRuleInput representa a entrada para atualizar uma regra de comissão
type UpdateCommissionRuleInput struct {
	TenantID        string
	ID              string
	Name            *string
	Description     *string
	Type            *string
	DefaultRate     *string
	MinAmount       *string
	MaxAmount       *string
	CalculationBase *string
	EffectiveFrom   *time.Time
	EffectiveTo     *time.Time
	Priority        *int
	IsActive        *bool
}

// UpdateCommissionRuleOutput representa a saída da atualização
type UpdateCommissionRuleOutput struct {
	CommissionRule *entity.CommissionRule
}

// UpdateCommissionRuleUseCase atualiza uma regra de comissão
type UpdateCommissionRuleUseCase struct {
	commissionRuleRepo repository.CommissionRuleRepository
}

// NewUpdateCommissionRuleUseCase cria uma nova instância do use case
func NewUpdateCommissionRuleUseCase(commissionRuleRepo repository.CommissionRuleRepository) *UpdateCommissionRuleUseCase {
	return &UpdateCommissionRuleUseCase{
		commissionRuleRepo: commissionRuleRepo,
	}
}

// Execute executa o use case
func (uc *UpdateCommissionRuleUseCase) Execute(ctx context.Context, input UpdateCommissionRuleInput) (*UpdateCommissionRuleOutput, error) {
	// Busca a regra existente
	rule, err := uc.commissionRuleRepo.GetByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}

	if rule == nil {
		return nil, domain.ErrCommissionRuleNotFound
	}

	// Atualiza campos se fornecidos
	if input.Name != nil {
		rule.Name = *input.Name
	}

	if input.Description != nil {
		rule.Description = input.Description
	}

	if input.Type != nil {
		rule.Type = *input.Type
	}

	if input.DefaultRate != nil {
		rate, err := decimal.NewFromString(*input.DefaultRate)
		if err != nil {
			return nil, err
		}
		rule.DefaultRate = rate
	}

	if input.MinAmount != nil {
		minAmt, err := decimal.NewFromString(*input.MinAmount)
		if err != nil {
			return nil, err
		}
		rule.MinAmount = &minAmt
	}

	if input.MaxAmount != nil {
		maxAmt, err := decimal.NewFromString(*input.MaxAmount)
		if err != nil {
			return nil, err
		}
		rule.MaxAmount = &maxAmt
	}

	if input.CalculationBase != nil {
		rule.CalculationBase = input.CalculationBase
	}

	if input.EffectiveFrom != nil {
		rule.EffectiveFrom = *input.EffectiveFrom
	}

	if input.EffectiveTo != nil {
		rule.EffectiveTo = input.EffectiveTo
	}

	if input.Priority != nil {
		rule.Priority = input.Priority
	}

	if input.IsActive != nil {
		rule.IsActive = *input.IsActive
	}

	// Valida alterações
	if err := rule.Validate(); err != nil {
		return nil, err
	}

	rule.UpdatedAt = time.Now()

	// Persiste
	updated, err := uc.commissionRuleRepo.Update(ctx, rule)
	if err != nil {
		return nil, err
	}

	return &UpdateCommissionRuleOutput{CommissionRule: updated}, nil
}
