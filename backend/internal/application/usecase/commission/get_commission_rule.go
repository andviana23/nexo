package commission

import (
	"context"


	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
)

// GetCommissionRuleInput representa a entrada para buscar uma regra de comissão
type GetCommissionRuleInput struct {
	TenantID string
	ID       string
}

// GetCommissionRuleOutput representa a saída da busca
type GetCommissionRuleOutput struct {
	CommissionRule *entity.CommissionRule
}

// GetCommissionRuleUseCase busca uma regra de comissão por ID
type GetCommissionRuleUseCase struct {
	commissionRuleRepo repository.CommissionRuleRepository
}

// NewGetCommissionRuleUseCase cria uma nova instância do use case
func NewGetCommissionRuleUseCase(commissionRuleRepo repository.CommissionRuleRepository) *GetCommissionRuleUseCase {
	return &GetCommissionRuleUseCase{
		commissionRuleRepo: commissionRuleRepo,
	}
}

// Execute executa o use case
func (uc *GetCommissionRuleUseCase) Execute(ctx context.Context, input GetCommissionRuleInput) (*GetCommissionRuleOutput, error) {
	rule, err := uc.commissionRuleRepo.GetByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetCommissionRuleOutput{CommissionRule: rule}, nil
}
