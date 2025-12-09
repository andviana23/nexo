package commission

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CreateAdvanceInput representa a entrada para criar um adiantamento
type CreateAdvanceInput struct {
	TenantID       string
	UnitID         *string
	ProfessionalID string
	Amount         string // Ex: "100.00"
	Reason         *string
	CreatedBy      *string
}

// CreateAdvanceOutput representa a saída da criação
type CreateAdvanceOutput struct {
	Advance *entity.Advance
}

// CreateAdvanceUseCase cria um novo adiantamento
type CreateAdvanceUseCase struct {
	advanceRepo repository.AdvanceRepository
}

// NewCreateAdvanceUseCase cria uma nova instância do use case
func NewCreateAdvanceUseCase(advanceRepo repository.AdvanceRepository) *CreateAdvanceUseCase {
	return &CreateAdvanceUseCase{
		advanceRepo: advanceRepo,
	}
}

// Execute executa o use case
func (uc *CreateAdvanceUseCase) Execute(ctx context.Context, input CreateAdvanceInput) (*CreateAdvanceOutput, error) {
	// Converte TenantID para UUID
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, err
	}

	// Converte Amount para decimal
	amount, err := decimal.NewFromString(input.Amount)
	if err != nil {
		return nil, err
	}

	// Cria a entidade validada
	advance, err := entity.NewAdvance(
		tenantUUID,
		input.ProfessionalID,
		amount,
	)
	if err != nil {
		return nil, err
	}

	// Define campos opcionais
	advance.UnitID = input.UnitID
	advance.Reason = input.Reason
	advance.CreatedBy = input.CreatedBy

	// Persiste
	created, err := uc.advanceRepo.Create(ctx, advance)
	if err != nil {
		return nil, err
	}

	return &CreateAdvanceOutput{Advance: created}, nil
}
