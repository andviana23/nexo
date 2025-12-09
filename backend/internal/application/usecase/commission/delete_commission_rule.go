package commission

import (
	"context"


	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
)

// DeleteCommissionRuleInput representa a entrada para deletar uma regra de comissão
type DeleteCommissionRuleInput struct {
	TenantID string
	ID       string
}

// DeleteCommissionRuleOutput representa a saída da exclusão
type DeleteCommissionRuleOutput struct {
	Success bool
}

// DeleteCommissionRuleUseCase deleta uma regra de comissão
type DeleteCommissionRuleUseCase struct {
	commissionRuleRepo repository.CommissionRuleRepository
}

// NewDeleteCommissionRuleUseCase cria uma nova instância do use case
func NewDeleteCommissionRuleUseCase(commissionRuleRepo repository.CommissionRuleRepository) *DeleteCommissionRuleUseCase {
	return &DeleteCommissionRuleUseCase{
		commissionRuleRepo: commissionRuleRepo,
	}
}

// Execute executa o use case
func (uc *DeleteCommissionRuleUseCase) Execute(ctx context.Context, input DeleteCommissionRuleInput) (*DeleteCommissionRuleOutput, error) {
	// Verifica se existe
	rule, err := uc.commissionRuleRepo.GetByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}

	if rule == nil {
		return nil, domain.ErrCommissionRuleNotFound
	}

	// Deleta
	err = uc.commissionRuleRepo.Delete(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}

	return &DeleteCommissionRuleOutput{Success: true}, nil
}

// DeactivateCommissionRuleInput representa a entrada para desativar uma regra de comissão
type DeactivateCommissionRuleInput struct {
	TenantID string
	ID       string
}

// DeactivateCommissionRuleOutput representa a saída da desativação
type DeactivateCommissionRuleOutput struct {
	Success bool
}

// DeactivateCommissionRuleUseCase desativa uma regra de comissão
type DeactivateCommissionRuleUseCase struct {
	commissionRuleRepo repository.CommissionRuleRepository
}

// NewDeactivateCommissionRuleUseCase cria uma nova instância do use case
func NewDeactivateCommissionRuleUseCase(commissionRuleRepo repository.CommissionRuleRepository) *DeactivateCommissionRuleUseCase {
	return &DeactivateCommissionRuleUseCase{
		commissionRuleRepo: commissionRuleRepo,
	}
}

// Execute executa o use case
func (uc *DeactivateCommissionRuleUseCase) Execute(ctx context.Context, input DeactivateCommissionRuleInput) (*DeactivateCommissionRuleOutput, error) {
	// Verifica se existe
	rule, err := uc.commissionRuleRepo.GetByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}

	if rule == nil {
		return nil, domain.ErrCommissionRuleNotFound
	}

	// Desativa
	err = uc.commissionRuleRepo.Deactivate(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}

	return &DeactivateCommissionRuleOutput{Success: true}, nil
}
