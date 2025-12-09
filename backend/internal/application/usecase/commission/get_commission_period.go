package commission

import (
	"context"

	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
)

// GetCommissionPeriodInput representa a entrada para buscar um período de comissão
type GetCommissionPeriodInput struct {
	TenantID string
	ID       string
}

// GetCommissionPeriodOutput representa a saída da busca
type GetCommissionPeriodOutput struct {
	CommissionPeriod *entity.CommissionPeriod
}

// GetCommissionPeriodUseCase busca um período de comissão por ID
type GetCommissionPeriodUseCase struct {
	commissionPeriodRepo repository.CommissionPeriodRepository
}

// NewGetCommissionPeriodUseCase cria uma nova instância do use case
func NewGetCommissionPeriodUseCase(commissionPeriodRepo repository.CommissionPeriodRepository) *GetCommissionPeriodUseCase {
	return &GetCommissionPeriodUseCase{
		commissionPeriodRepo: commissionPeriodRepo,
	}
}

// Execute executa o use case
func (uc *GetCommissionPeriodUseCase) Execute(ctx context.Context, input GetCommissionPeriodInput) (*GetCommissionPeriodOutput, error) {
	period, err := uc.commissionPeriodRepo.GetByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetCommissionPeriodOutput{CommissionPeriod: period}, nil
}

// GetOpenCommissionPeriodInput representa a entrada para buscar o período aberto de um profissional
type GetOpenCommissionPeriodInput struct {
	TenantID       string
	ProfessionalID string
}

// GetOpenCommissionPeriodOutput representa a saída da busca
type GetOpenCommissionPeriodOutput struct {
	CommissionPeriod *entity.CommissionPeriod
}

// GetOpenCommissionPeriodUseCase busca o período aberto de um profissional
type GetOpenCommissionPeriodUseCase struct {
	commissionPeriodRepo repository.CommissionPeriodRepository
}

// NewGetOpenCommissionPeriodUseCase cria uma nova instância do use case
func NewGetOpenCommissionPeriodUseCase(commissionPeriodRepo repository.CommissionPeriodRepository) *GetOpenCommissionPeriodUseCase {
	return &GetOpenCommissionPeriodUseCase{
		commissionPeriodRepo: commissionPeriodRepo,
	}
}

// Execute executa o use case
func (uc *GetOpenCommissionPeriodUseCase) Execute(ctx context.Context, input GetOpenCommissionPeriodInput) (*GetOpenCommissionPeriodOutput, error) {
	period, err := uc.commissionPeriodRepo.GetOpenByProfessional(ctx, input.TenantID, input.ProfessionalID)
	if err != nil {
		return nil, err
	}

	return &GetOpenCommissionPeriodOutput{CommissionPeriod: period}, nil
}

// GetCommissionPeriodSummaryInput representa a entrada para buscar resumo do período
type GetCommissionPeriodSummaryInput struct {
	TenantID string
	PeriodID string
}

// GetCommissionPeriodSummaryOutput representa a saída da busca do resumo
type GetCommissionPeriodSummaryOutput struct {
	Summary *entity.CommissionPeriodSummary
}

// GetCommissionPeriodSummaryUseCase busca o resumo de um período de comissão
type GetCommissionPeriodSummaryUseCase struct {
	commissionPeriodRepo repository.CommissionPeriodRepository
}

// NewGetCommissionPeriodSummaryUseCase cria uma nova instância do use case
func NewGetCommissionPeriodSummaryUseCase(commissionPeriodRepo repository.CommissionPeriodRepository) *GetCommissionPeriodSummaryUseCase {
	return &GetCommissionPeriodSummaryUseCase{
		commissionPeriodRepo: commissionPeriodRepo,
	}
}

// Execute executa o use case
func (uc *GetCommissionPeriodSummaryUseCase) Execute(ctx context.Context, input GetCommissionPeriodSummaryInput) (*GetCommissionPeriodSummaryOutput, error) {
	summary, err := uc.commissionPeriodRepo.GetSummary(ctx, input.TenantID, input.PeriodID)
	if err != nil {
		return nil, err
	}

	return &GetCommissionPeriodSummaryOutput{Summary: summary}, nil
}

// ListCommissionPeriodsInput representa a entrada para listar períodos de comissão
type ListCommissionPeriodsInput struct {
	TenantID       string
	ProfessionalID *string
	Status         *string
	Limit          int
	Offset         int
}

// ListCommissionPeriodsOutput representa a saída da listagem
type ListCommissionPeriodsOutput struct {
	CommissionPeriods []*entity.CommissionPeriod
}

// ListCommissionPeriodsUseCase lista períodos de comissão
type ListCommissionPeriodsUseCase struct {
	commissionPeriodRepo repository.CommissionPeriodRepository
}

// NewListCommissionPeriodsUseCase cria uma nova instância do use case
func NewListCommissionPeriodsUseCase(commissionPeriodRepo repository.CommissionPeriodRepository) *ListCommissionPeriodsUseCase {
	return &ListCommissionPeriodsUseCase{
		commissionPeriodRepo: commissionPeriodRepo,
	}
}

// Execute executa o use case
func (uc *ListCommissionPeriodsUseCase) Execute(ctx context.Context, input ListCommissionPeriodsInput) (*ListCommissionPeriodsOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 50
	}

	periods, err := uc.commissionPeriodRepo.List(ctx, input.TenantID, input.ProfessionalID, input.Status, limit, input.Offset)
	if err != nil {
		return nil, err
	}

	return &ListCommissionPeriodsOutput{CommissionPeriods: periods}, nil
}

// ListCommissionPeriodsByDateRangeInput representa a entrada para listar períodos por intervalo de datas
type ListCommissionPeriodsByDateRangeInput struct {
	TenantID  string
	StartDate time.Time
	EndDate   time.Time
}

// ListCommissionPeriodsByDateRangeOutput representa a saída da listagem
type ListCommissionPeriodsByDateRangeOutput struct {
	CommissionPeriods []*entity.CommissionPeriod
}

// ListCommissionPeriodsByDateRangeUseCase lista períodos de comissão por intervalo de datas
type ListCommissionPeriodsByDateRangeUseCase struct {
	commissionPeriodRepo repository.CommissionPeriodRepository
}

// NewListCommissionPeriodsByDateRangeUseCase cria uma nova instância do use case
func NewListCommissionPeriodsByDateRangeUseCase(commissionPeriodRepo repository.CommissionPeriodRepository) *ListCommissionPeriodsByDateRangeUseCase {
	return &ListCommissionPeriodsByDateRangeUseCase{
		commissionPeriodRepo: commissionPeriodRepo,
	}
}

// Execute executa o use case
func (uc *ListCommissionPeriodsByDateRangeUseCase) Execute(ctx context.Context, input ListCommissionPeriodsByDateRangeInput) (*ListCommissionPeriodsByDateRangeOutput, error) {
	periods, err := uc.commissionPeriodRepo.GetByDateRange(ctx, input.TenantID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}

	return &ListCommissionPeriodsByDateRangeOutput{CommissionPeriods: periods}, nil
}
