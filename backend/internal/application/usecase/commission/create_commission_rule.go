package commission

import (
	"context"

	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CreateCommissionRuleInput representa a entrada para criar uma regra de comissão
type CreateCommissionRuleInput struct {
	TenantID        string
	UnitID          *string
	Name            string
	Description     *string
	Type            string // PERCENTUAL, FIXO
	DefaultRate     string // Ex: "50.00" para 50%
	MinAmount       *string
	MaxAmount       *string
	CalculationBase *string // BRUTO, LIQUIDO
	EffectiveFrom   *time.Time
	EffectiveTo     *time.Time
	Priority        *int
	CreatedBy       *string
}

// CreateCommissionRuleOutput representa a saída da criação
type CreateCommissionRuleOutput struct {
	CommissionRule *entity.CommissionRule
}

// CreateCommissionRuleUseCase cria uma nova regra de comissão
type CreateCommissionRuleUseCase struct {
	commissionRuleRepo repository.CommissionRuleRepository
}

// NewCreateCommissionRuleUseCase cria uma nova instância do use case
func NewCreateCommissionRuleUseCase(commissionRuleRepo repository.CommissionRuleRepository) *CreateCommissionRuleUseCase {
	return &CreateCommissionRuleUseCase{
		commissionRuleRepo: commissionRuleRepo,
	}
}

// Execute executa o use case
func (uc *CreateCommissionRuleUseCase) Execute(ctx context.Context, input CreateCommissionRuleInput) (*CreateCommissionRuleOutput, error) {
	// Converte DefaultRate para decimal
	defaultRate, err := decimal.NewFromString(input.DefaultRate)
	if err != nil {
		return nil, err
	}

	// Converte TenantID para UUID
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, err
	}

	// Cria a entidade validada
	rule, err := entity.NewCommissionRule(
		tenantUUID,
		input.Name,
		input.Type,
		defaultRate,
	)
	if err != nil {
		return nil, err
	}

	// Define campos opcionais
	rule.UnitID = input.UnitID
	rule.Description = input.Description
	rule.Priority = input.Priority
	rule.CreatedBy = input.CreatedBy

	if input.CalculationBase != nil {
		rule.CalculationBase = input.CalculationBase
	}

	if input.EffectiveFrom != nil {
		rule.EffectiveFrom = *input.EffectiveFrom
	}
	rule.EffectiveTo = input.EffectiveTo

	// Converte MinAmount/MaxAmount
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

	// Valida novamente com todos os campos
	if err := rule.Validate(); err != nil {
		return nil, err
	}

	// Persiste
	created, err := uc.commissionRuleRepo.Create(ctx, rule)
	if err != nil {
		return nil, err
	}

	return &CreateCommissionRuleOutput{CommissionRule: created}, nil
}
