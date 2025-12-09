package commission

import (
	"context"

	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
)

// ProcessCommissionItemInput representa a entrada para processar um item de comissão
type ProcessCommissionItemInput struct {
	TenantID string
	ItemID   string
	PeriodID string
}

// ProcessCommissionItemOutput representa a saída do processamento
type ProcessCommissionItemOutput struct {
	CommissionItem *entity.CommissionItem
}

// ProcessCommissionItemUseCase processa um item de comissão
type ProcessCommissionItemUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewProcessCommissionItemUseCase cria uma nova instância do use case
func NewProcessCommissionItemUseCase(commissionItemRepo repository.CommissionItemRepository) *ProcessCommissionItemUseCase {
	return &ProcessCommissionItemUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *ProcessCommissionItemUseCase) Execute(ctx context.Context, input ProcessCommissionItemInput) (*ProcessCommissionItemOutput, error) {
	// Verifica se existe
	item, err := uc.commissionItemRepo.GetByID(ctx, input.TenantID, input.ItemID)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, domain.ErrCommissionItemNotFound
	}

	// Verifica se pode processar
	if !item.CanProcess() {
		return nil, domain.ErrItemNaoPodeProcessar
	}

	// Processa
	processed, err := uc.commissionItemRepo.Process(ctx, input.TenantID, input.ItemID, input.PeriodID)
	if err != nil {
		return nil, err
	}

	return &ProcessCommissionItemOutput{CommissionItem: processed}, nil
}

// AssignItemsToPeriodInput representa a entrada para vincular itens a um período
type AssignItemsToPeriodInput struct {
	TenantID       string
	ProfessionalID string
	PeriodID       string
	StartDate      time.Time
	EndDate        time.Time
}

// AssignItemsToPeriodOutput representa a saída da vinculação
type AssignItemsToPeriodOutput struct {
	ItemsAssigned int64
}

// AssignItemsToPeriodUseCase vincula itens pendentes a um período
type AssignItemsToPeriodUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewAssignItemsToPeriodUseCase cria uma nova instância do use case
func NewAssignItemsToPeriodUseCase(commissionItemRepo repository.CommissionItemRepository) *AssignItemsToPeriodUseCase {
	return &AssignItemsToPeriodUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *AssignItemsToPeriodUseCase) Execute(ctx context.Context, input AssignItemsToPeriodInput) (*AssignItemsToPeriodOutput, error) {
	count, err := uc.commissionItemRepo.AssignToPeriod(ctx, input.TenantID, input.ProfessionalID, input.PeriodID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}

	return &AssignItemsToPeriodOutput{ItemsAssigned: count}, nil
}

// DeleteCommissionItemInput representa a entrada para deletar um item de comissão
type DeleteCommissionItemInput struct {
	TenantID string
	ItemID   string
}

// DeleteCommissionItemOutput representa a saída da exclusão
type DeleteCommissionItemOutput struct {
	Success bool
}

// DeleteCommissionItemUseCase deleta um item de comissão
type DeleteCommissionItemUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewDeleteCommissionItemUseCase cria uma nova instância do use case
func NewDeleteCommissionItemUseCase(commissionItemRepo repository.CommissionItemRepository) *DeleteCommissionItemUseCase {
	return &DeleteCommissionItemUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *DeleteCommissionItemUseCase) Execute(ctx context.Context, input DeleteCommissionItemInput) (*DeleteCommissionItemOutput, error) {
	// Verifica se existe
	item, err := uc.commissionItemRepo.GetByID(ctx, input.TenantID, input.ItemID)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, domain.ErrCommissionItemNotFound
	}

	// Só pode deletar se estiver pendente
	if item.Status != "PENDENTE" {
		return nil, domain.ErrItemNaoPodeCancelar
	}

	// Deleta
	err = uc.commissionItemRepo.Delete(ctx, input.TenantID, input.ItemID)
	if err != nil {
		return nil, err
	}

	return &DeleteCommissionItemOutput{Success: true}, nil
}

// DeleteCommissionItemByCommandItemInput representa a entrada para deletar item por item de comanda
type DeleteCommissionItemByCommandItemInput struct {
	TenantID      string
	CommandItemID string
}

// DeleteCommissionItemByCommandItemOutput representa a saída da exclusão
type DeleteCommissionItemByCommandItemOutput struct {
	Success bool
}

// DeleteCommissionItemByCommandItemUseCase deleta um item de comissão por item de comanda
type DeleteCommissionItemByCommandItemUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewDeleteCommissionItemByCommandItemUseCase cria uma nova instância do use case
func NewDeleteCommissionItemByCommandItemUseCase(commissionItemRepo repository.CommissionItemRepository) *DeleteCommissionItemByCommandItemUseCase {
	return &DeleteCommissionItemByCommandItemUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *DeleteCommissionItemByCommandItemUseCase) Execute(ctx context.Context, input DeleteCommissionItemByCommandItemInput) (*DeleteCommissionItemByCommandItemOutput, error) {
	// Verifica se existe
	item, err := uc.commissionItemRepo.GetByCommandItem(ctx, input.TenantID, input.CommandItemID)
	if err != nil {
		return nil, err
	}

	if item == nil {
		// Não existe, considera sucesso
		return &DeleteCommissionItemByCommandItemOutput{Success: true}, nil
	}

	// Só pode deletar se estiver pendente
	if item.Status != "PENDENTE" {
		return nil, domain.ErrItemNaoPodeCancelar
	}

	// Deleta
	err = uc.commissionItemRepo.DeleteByCommandItem(ctx, input.TenantID, input.CommandItemID)
	if err != nil {
		return nil, err
	}

	return &DeleteCommissionItemByCommandItemOutput{Success: true}, nil
}
