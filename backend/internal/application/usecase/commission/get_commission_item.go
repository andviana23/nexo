package commission

import (
	"context"

	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
)

// GetCommissionItemInput representa a entrada para buscar um item de comissão
type GetCommissionItemInput struct {
	TenantID string
	ID       string
}

// GetCommissionItemOutput representa a saída da busca
type GetCommissionItemOutput struct {
	CommissionItem *entity.CommissionItem
}

// GetCommissionItemUseCase busca um item de comissão por ID
type GetCommissionItemUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewGetCommissionItemUseCase cria uma nova instância do use case
func NewGetCommissionItemUseCase(commissionItemRepo repository.CommissionItemRepository) *GetCommissionItemUseCase {
	return &GetCommissionItemUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *GetCommissionItemUseCase) Execute(ctx context.Context, input GetCommissionItemInput) (*GetCommissionItemOutput, error) {
	item, err := uc.commissionItemRepo.GetByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetCommissionItemOutput{CommissionItem: item}, nil
}

// GetCommissionItemByCommandItemInput representa a entrada para buscar item por item de comanda
type GetCommissionItemByCommandItemInput struct {
	TenantID      string
	CommandItemID string
}

// GetCommissionItemByCommandItemOutput representa a saída da busca
type GetCommissionItemByCommandItemOutput struct {
	CommissionItem *entity.CommissionItem
}

// GetCommissionItemByCommandItemUseCase busca um item de comissão por item de comanda
type GetCommissionItemByCommandItemUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewGetCommissionItemByCommandItemUseCase cria uma nova instância do use case
func NewGetCommissionItemByCommandItemUseCase(commissionItemRepo repository.CommissionItemRepository) *GetCommissionItemByCommandItemUseCase {
	return &GetCommissionItemByCommandItemUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *GetCommissionItemByCommandItemUseCase) Execute(ctx context.Context, input GetCommissionItemByCommandItemInput) (*GetCommissionItemByCommandItemOutput, error) {
	item, err := uc.commissionItemRepo.GetByCommandItem(ctx, input.TenantID, input.CommandItemID)
	if err != nil {
		return nil, err
	}

	return &GetCommissionItemByCommandItemOutput{CommissionItem: item}, nil
}

// ListCommissionItemsInput representa a entrada para listar itens de comissão
type ListCommissionItemsInput struct {
	TenantID       string
	ProfessionalID *string
	PeriodID       *string
	Status         *string
	Limit          int
	Offset         int
}

// ListCommissionItemsOutput representa a saída da listagem
type ListCommissionItemsOutput struct {
	CommissionItems []*entity.CommissionItem
}

// ListCommissionItemsUseCase lista itens de comissão
type ListCommissionItemsUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewListCommissionItemsUseCase cria uma nova instância do use case
func NewListCommissionItemsUseCase(commissionItemRepo repository.CommissionItemRepository) *ListCommissionItemsUseCase {
	return &ListCommissionItemsUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *ListCommissionItemsUseCase) Execute(ctx context.Context, input ListCommissionItemsInput) (*ListCommissionItemsOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 50
	}

	items, err := uc.commissionItemRepo.List(ctx, input.TenantID, input.ProfessionalID, input.PeriodID, input.Status, limit, input.Offset)
	if err != nil {
		return nil, err
	}

	return &ListCommissionItemsOutput{CommissionItems: items}, nil
}

// GetPendingCommissionItemsInput representa a entrada para listar itens pendentes
type GetPendingCommissionItemsInput struct {
	TenantID       string
	ProfessionalID string
}

// GetPendingCommissionItemsOutput representa a saída da listagem
type GetPendingCommissionItemsOutput struct {
	CommissionItems []*entity.CommissionItem
}

// GetPendingCommissionItemsUseCase lista itens de comissão pendentes de um profissional
type GetPendingCommissionItemsUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewGetPendingCommissionItemsUseCase cria uma nova instância do use case
func NewGetPendingCommissionItemsUseCase(commissionItemRepo repository.CommissionItemRepository) *GetPendingCommissionItemsUseCase {
	return &GetPendingCommissionItemsUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *GetPendingCommissionItemsUseCase) Execute(ctx context.Context, input GetPendingCommissionItemsInput) (*GetPendingCommissionItemsOutput, error) {
	items, err := uc.commissionItemRepo.GetPendingByProfessional(ctx, input.TenantID, input.ProfessionalID)
	if err != nil {
		return nil, err
	}

	return &GetPendingCommissionItemsOutput{CommissionItems: items}, nil
}

// ListCommissionItemsByDateRangeInput representa a entrada para listar itens por intervalo de datas
type ListCommissionItemsByDateRangeInput struct {
	TenantID  string
	StartDate time.Time
	EndDate   time.Time
}

// ListCommissionItemsByDateRangeOutput representa a saída da listagem
type ListCommissionItemsByDateRangeOutput struct {
	CommissionItems []*entity.CommissionItem
}

// ListCommissionItemsByDateRangeUseCase lista itens de comissão por intervalo de datas
type ListCommissionItemsByDateRangeUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewListCommissionItemsByDateRangeUseCase cria uma nova instância do use case
func NewListCommissionItemsByDateRangeUseCase(commissionItemRepo repository.CommissionItemRepository) *ListCommissionItemsByDateRangeUseCase {
	return &ListCommissionItemsByDateRangeUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *ListCommissionItemsByDateRangeUseCase) Execute(ctx context.Context, input ListCommissionItemsByDateRangeInput) (*ListCommissionItemsByDateRangeOutput, error) {
	items, err := uc.commissionItemRepo.GetByDateRange(ctx, input.TenantID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}

	return &ListCommissionItemsByDateRangeOutput{CommissionItems: items}, nil
}

// GetCommissionSummaryByProfessionalInput representa a entrada para resumo por profissional
type GetCommissionSummaryByProfessionalInput struct {
	TenantID  string
	StartDate time.Time
	EndDate   time.Time
}

// GetCommissionSummaryByProfessionalOutput representa a saída do resumo
type GetCommissionSummaryByProfessionalOutput struct {
	Summaries []*entity.CommissionSummary
}

// GetCommissionSummaryByProfessionalUseCase retorna resumo de comissões por profissional
type GetCommissionSummaryByProfessionalUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewGetCommissionSummaryByProfessionalUseCase cria uma nova instância do use case
func NewGetCommissionSummaryByProfessionalUseCase(commissionItemRepo repository.CommissionItemRepository) *GetCommissionSummaryByProfessionalUseCase {
	return &GetCommissionSummaryByProfessionalUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *GetCommissionSummaryByProfessionalUseCase) Execute(ctx context.Context, input GetCommissionSummaryByProfessionalInput) (*GetCommissionSummaryByProfessionalOutput, error) {
	summaries, err := uc.commissionItemRepo.GetSummaryByProfessional(ctx, input.TenantID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}

	return &GetCommissionSummaryByProfessionalOutput{Summaries: summaries}, nil
}

// GetCommissionSummaryByServiceInput representa a entrada para resumo por serviço
type GetCommissionSummaryByServiceInput struct {
	TenantID  string
	StartDate time.Time
	EndDate   time.Time
}

// GetCommissionSummaryByServiceOutput representa a saída do resumo
type GetCommissionSummaryByServiceOutput struct {
	Summaries []*entity.CommissionByService
}

// GetCommissionSummaryByServiceUseCase retorna resumo de comissões por serviço
type GetCommissionSummaryByServiceUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewGetCommissionSummaryByServiceUseCase cria uma nova instância do use case
func NewGetCommissionSummaryByServiceUseCase(commissionItemRepo repository.CommissionItemRepository) *GetCommissionSummaryByServiceUseCase {
	return &GetCommissionSummaryByServiceUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *GetCommissionSummaryByServiceUseCase) Execute(ctx context.Context, input GetCommissionSummaryByServiceInput) (*GetCommissionSummaryByServiceOutput, error) {
	summaries, err := uc.commissionItemRepo.GetSummaryByService(ctx, input.TenantID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}

	return &GetCommissionSummaryByServiceOutput{Summaries: summaries}, nil
}
