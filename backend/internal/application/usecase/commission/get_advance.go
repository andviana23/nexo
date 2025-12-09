package commission

import (
	"context"

	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
)

// GetAdvanceInput representa a entrada para buscar um adiantamento
type GetAdvanceInput struct {
	TenantID string
	ID       string
}

// GetAdvanceOutput representa a saída da busca
type GetAdvanceOutput struct {
	Advance *entity.Advance
}

// GetAdvanceUseCase busca um adiantamento por ID
type GetAdvanceUseCase struct {
	advanceRepo repository.AdvanceRepository
}

// NewGetAdvanceUseCase cria uma nova instância do use case
func NewGetAdvanceUseCase(advanceRepo repository.AdvanceRepository) *GetAdvanceUseCase {
	return &GetAdvanceUseCase{
		advanceRepo: advanceRepo,
	}
}

// Execute executa o use case
func (uc *GetAdvanceUseCase) Execute(ctx context.Context, input GetAdvanceInput) (*GetAdvanceOutput, error) {
	advance, err := uc.advanceRepo.GetByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetAdvanceOutput{Advance: advance}, nil
}

// ListAdvancesInput representa a entrada para listar adiantamentos
type ListAdvancesInput struct {
	TenantID       string
	ProfessionalID *string
	Status         *string
	Limit          int
	Offset         int
}

// ListAdvancesOutput representa a saída da listagem
type ListAdvancesOutput struct {
	Advances []*entity.Advance
}

// ListAdvancesUseCase lista adiantamentos
type ListAdvancesUseCase struct {
	advanceRepo repository.AdvanceRepository
}

// NewListAdvancesUseCase cria uma nova instância do use case
func NewListAdvancesUseCase(advanceRepo repository.AdvanceRepository) *ListAdvancesUseCase {
	return &ListAdvancesUseCase{
		advanceRepo: advanceRepo,
	}
}

// Execute executa o use case
func (uc *ListAdvancesUseCase) Execute(ctx context.Context, input ListAdvancesInput) (*ListAdvancesOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 50
	}

	advances, err := uc.advanceRepo.List(ctx, input.TenantID, input.ProfessionalID, input.Status, limit, input.Offset)
	if err != nil {
		return nil, err
	}

	return &ListAdvancesOutput{Advances: advances}, nil
}

// GetPendingAdvancesInput representa a entrada para listar adiantamentos pendentes
type GetPendingAdvancesInput struct {
	TenantID       string
	ProfessionalID string
}

// GetPendingAdvancesOutput representa a saída da listagem
type GetPendingAdvancesOutput struct {
	Advances     []*entity.Advance
	TotalPending float64
}

// GetPendingAdvancesUseCase lista adiantamentos pendentes de um profissional
type GetPendingAdvancesUseCase struct {
	advanceRepo repository.AdvanceRepository
}

// NewGetPendingAdvancesUseCase cria uma nova instância do use case
func NewGetPendingAdvancesUseCase(advanceRepo repository.AdvanceRepository) *GetPendingAdvancesUseCase {
	return &GetPendingAdvancesUseCase{
		advanceRepo: advanceRepo,
	}
}

// Execute executa o use case
func (uc *GetPendingAdvancesUseCase) Execute(ctx context.Context, input GetPendingAdvancesInput) (*GetPendingAdvancesOutput, error) {
	advances, err := uc.advanceRepo.GetPendingByProfessional(ctx, input.TenantID, input.ProfessionalID)
	if err != nil {
		return nil, err
	}

	total, err := uc.advanceRepo.GetTotalPendingByProfessional(ctx, input.TenantID, input.ProfessionalID)
	if err != nil {
		return nil, err
	}

	return &GetPendingAdvancesOutput{Advances: advances, TotalPending: total}, nil
}

// GetApprovedAdvancesInput representa a entrada para listar adiantamentos aprovados
type GetApprovedAdvancesInput struct {
	TenantID       string
	ProfessionalID string
}

// GetApprovedAdvancesOutput representa a saída da listagem
type GetApprovedAdvancesOutput struct {
	Advances      []*entity.Advance
	TotalApproved float64
}

// GetApprovedAdvancesUseCase lista adiantamentos aprovados (não deduzidos) de um profissional
type GetApprovedAdvancesUseCase struct {
	advanceRepo repository.AdvanceRepository
}

// NewGetApprovedAdvancesUseCase cria uma nova instância do use case
func NewGetApprovedAdvancesUseCase(advanceRepo repository.AdvanceRepository) *GetApprovedAdvancesUseCase {
	return &GetApprovedAdvancesUseCase{
		advanceRepo: advanceRepo,
	}
}

// Execute executa o use case
func (uc *GetApprovedAdvancesUseCase) Execute(ctx context.Context, input GetApprovedAdvancesInput) (*GetApprovedAdvancesOutput, error) {
	advances, err := uc.advanceRepo.GetApprovedByProfessional(ctx, input.TenantID, input.ProfessionalID)
	if err != nil {
		return nil, err
	}

	total, err := uc.advanceRepo.GetTotalApprovedByProfessional(ctx, input.TenantID, input.ProfessionalID)
	if err != nil {
		return nil, err
	}

	return &GetApprovedAdvancesOutput{Advances: advances, TotalApproved: total}, nil
}

// ListAdvancesByDateRangeInput representa a entrada para listar adiantamentos por intervalo de datas
type ListAdvancesByDateRangeInput struct {
	TenantID  string
	StartDate time.Time
	EndDate   time.Time
}

// ListAdvancesByDateRangeOutput representa a saída da listagem
type ListAdvancesByDateRangeOutput struct {
	Advances []*entity.Advance
}

// ListAdvancesByDateRangeUseCase lista adiantamentos por intervalo de datas
type ListAdvancesByDateRangeUseCase struct {
	advanceRepo repository.AdvanceRepository
}

// NewListAdvancesByDateRangeUseCase cria uma nova instância do use case
func NewListAdvancesByDateRangeUseCase(advanceRepo repository.AdvanceRepository) *ListAdvancesByDateRangeUseCase {
	return &ListAdvancesByDateRangeUseCase{
		advanceRepo: advanceRepo,
	}
}

// Execute executa o use case
func (uc *ListAdvancesByDateRangeUseCase) Execute(ctx context.Context, input ListAdvancesByDateRangeInput) (*ListAdvancesByDateRangeOutput, error) {
	advances, err := uc.advanceRepo.GetByDateRange(ctx, input.TenantID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}

	return &ListAdvancesByDateRangeOutput{Advances: advances}, nil
}
