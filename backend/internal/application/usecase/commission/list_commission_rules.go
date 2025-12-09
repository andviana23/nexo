package commission

import (
	"context"

	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
)

// ListCommissionRulesInput representa a entrada para listar regras de comissão
type ListCommissionRulesInput struct {
	TenantID   string
	ActiveOnly bool
}

// ListCommissionRulesOutput representa a saída da listagem
type ListCommissionRulesOutput struct {
	CommissionRules []*entity.CommissionRule
}

// ListCommissionRulesUseCase lista regras de comissão
type ListCommissionRulesUseCase struct {
	commissionRuleRepo repository.CommissionRuleRepository
}

// NewListCommissionRulesUseCase cria uma nova instância do use case
func NewListCommissionRulesUseCase(commissionRuleRepo repository.CommissionRuleRepository) *ListCommissionRulesUseCase {
	return &ListCommissionRulesUseCase{
		commissionRuleRepo: commissionRuleRepo,
	}
}

// Execute executa o use case
func (uc *ListCommissionRulesUseCase) Execute(ctx context.Context, input ListCommissionRulesInput) (*ListCommissionRulesOutput, error) {
	var rules []*entity.CommissionRule
	var err error

	if input.ActiveOnly {
		rules, err = uc.commissionRuleRepo.ListActive(ctx, input.TenantID)
	} else {
		rules, err = uc.commissionRuleRepo.List(ctx, input.TenantID)
	}

	if err != nil {
		return nil, err
	}

	return &ListCommissionRulesOutput{CommissionRules: rules}, nil
}

// GetEffectiveCommissionRulesInput representa a entrada para buscar regras vigentes
type GetEffectiveCommissionRulesInput struct {
	TenantID string
	Date     time.Time
}

// GetEffectiveCommissionRulesOutput representa a saída da busca de regras vigentes
type GetEffectiveCommissionRulesOutput struct {
	CommissionRules []*entity.CommissionRule
}

// GetEffectiveCommissionRulesUseCase busca regras de comissão vigentes em uma data
type GetEffectiveCommissionRulesUseCase struct {
	commissionRuleRepo repository.CommissionRuleRepository
}

// NewGetEffectiveCommissionRulesUseCase cria uma nova instância do use case
func NewGetEffectiveCommissionRulesUseCase(commissionRuleRepo repository.CommissionRuleRepository) *GetEffectiveCommissionRulesUseCase {
	return &GetEffectiveCommissionRulesUseCase{
		commissionRuleRepo: commissionRuleRepo,
	}
}

// Execute executa o use case
func (uc *GetEffectiveCommissionRulesUseCase) Execute(ctx context.Context, input GetEffectiveCommissionRulesInput) (*GetEffectiveCommissionRulesOutput, error) {
	rules, err := uc.commissionRuleRepo.GetEffective(ctx, input.TenantID, input.Date)
	if err != nil {
		return nil, err
	}

	return &GetEffectiveCommissionRulesOutput{CommissionRules: rules}, nil
}
