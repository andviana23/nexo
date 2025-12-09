package commission

import (
	"context"

	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
	"github.com/google/uuid"
)

// CreateCommissionPeriodInput representa a entrada para criar um período de comissão
type CreateCommissionPeriodInput struct {
	TenantID       string
	UnitID         *string
	ReferenceMonth string // formato: YYYY-MM
	ProfessionalID string
	PeriodStart    time.Time
	PeriodEnd      time.Time
	Notes          *string
}

// CreateCommissionPeriodOutput representa a saída da criação
type CreateCommissionPeriodOutput struct {
	CommissionPeriod *entity.CommissionPeriod
}

// CreateCommissionPeriodUseCase cria um novo período de comissão
type CreateCommissionPeriodUseCase struct {
	commissionPeriodRepo repository.CommissionPeriodRepository
}

// NewCreateCommissionPeriodUseCase cria uma nova instância do use case
func NewCreateCommissionPeriodUseCase(commissionPeriodRepo repository.CommissionPeriodRepository) *CreateCommissionPeriodUseCase {
	return &CreateCommissionPeriodUseCase{
		commissionPeriodRepo: commissionPeriodRepo,
	}
}

// Execute executa o use case
func (uc *CreateCommissionPeriodUseCase) Execute(ctx context.Context, input CreateCommissionPeriodInput) (*CreateCommissionPeriodOutput, error) {
	// Verifica se já existe período aberto para o profissional
	existingPeriod, err := uc.commissionPeriodRepo.GetOpenByProfessional(ctx, input.TenantID, input.ProfessionalID)
	if err != nil {
		return nil, err
	}

	if existingPeriod != nil {
		// Se já existe período aberto, retorna ele
		return &CreateCommissionPeriodOutput{CommissionPeriod: existingPeriod}, nil
	}

	// Converte TenantID para UUID
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, err
	}

	// Cria a entidade validada
	period, err := entity.NewCommissionPeriod(
		tenantUUID,
		input.ReferenceMonth,
		input.ProfessionalID,
		input.PeriodStart,
		input.PeriodEnd,
	)
	if err != nil {
		return nil, err
	}

	// Define campos opcionais
	period.UnitID = input.UnitID
	period.Notes = input.Notes

	// Persiste
	created, err := uc.commissionPeriodRepo.Create(ctx, period)
	if err != nil {
		return nil, err
	}

	return &CreateCommissionPeriodOutput{CommissionPeriod: created}, nil
}
