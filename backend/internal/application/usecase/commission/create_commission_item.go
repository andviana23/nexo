package commission

import (
	"context"

	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CreateCommissionItemInput representa a entrada para criar um item de comissão
type CreateCommissionItemInput struct {
	TenantID         string
	UnitID           *string
	ProfessionalID   string
	CommandID        *string
	CommandItemID    *string
	AppointmentID    *string
	ServiceID        *string
	ServiceName      *string
	GrossValue       string // Ex: "100.00"
	CommissionRate   string // Ex: "50.00" para 50%
	CommissionType   string // PERCENTUAL, FIXO
	CommissionSource string // SERVICO, PROFISSIONAL, REGRA, MANUAL
	RuleID           *string
	ReferenceDate    time.Time
	Description      *string
}

// CreateCommissionItemOutput representa a saída da criação
type CreateCommissionItemOutput struct {
	CommissionItem *entity.CommissionItem
}

// CreateCommissionItemUseCase cria um novo item de comissão
type CreateCommissionItemUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewCreateCommissionItemUseCase cria uma nova instância do use case
func NewCreateCommissionItemUseCase(commissionItemRepo repository.CommissionItemRepository) *CreateCommissionItemUseCase {
	return &CreateCommissionItemUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *CreateCommissionItemUseCase) Execute(ctx context.Context, input CreateCommissionItemInput) (*CreateCommissionItemOutput, error) {
	// Converte valores para decimal
	grossValue, err := decimal.NewFromString(input.GrossValue)
	if err != nil {
		return nil, err
	}

	commissionRate, err := decimal.NewFromString(input.CommissionRate)
	if err != nil {
		return nil, err
	}

	// Converte TenantID para UUID
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, err
	}

	// Cria a entidade validada
	item, err := entity.NewCommissionItem(
		tenantUUID,
		input.ProfessionalID,
		grossValue,
		commissionRate,
		input.CommissionType,
		input.CommissionSource,
		input.ReferenceDate,
	)
	if err != nil {
		return nil, err
	}

	// Define campos opcionais
	item.UnitID = input.UnitID
	item.CommandID = input.CommandID
	item.CommandItemID = input.CommandItemID
	item.AppointmentID = input.AppointmentID
	item.ServiceID = input.ServiceID
	item.ServiceName = input.ServiceName
	item.RuleID = input.RuleID
	item.Description = input.Description

	// Persiste
	created, err := uc.commissionItemRepo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return &CreateCommissionItemOutput{CommissionItem: created}, nil
}

// CreateCommissionItemBatchInput representa a entrada para criar múltiplos itens
type CreateCommissionItemBatchInput struct {
	Items []CreateCommissionItemInput
}

// CreateCommissionItemBatchOutput representa a saída da criação em lote
type CreateCommissionItemBatchOutput struct {
	CommissionItems []*entity.CommissionItem
}

// CreateCommissionItemBatchUseCase cria múltiplos itens de comissão
type CreateCommissionItemBatchUseCase struct {
	commissionItemRepo repository.CommissionItemRepository
}

// NewCreateCommissionItemBatchUseCase cria uma nova instância do use case
func NewCreateCommissionItemBatchUseCase(commissionItemRepo repository.CommissionItemRepository) *CreateCommissionItemBatchUseCase {
	return &CreateCommissionItemBatchUseCase{
		commissionItemRepo: commissionItemRepo,
	}
}

// Execute executa o use case
func (uc *CreateCommissionItemBatchUseCase) Execute(ctx context.Context, input CreateCommissionItemBatchInput) (*CreateCommissionItemBatchOutput, error) {
	items := make([]*entity.CommissionItem, 0, len(input.Items))

	for _, itemInput := range input.Items {
		grossValue, err := decimal.NewFromString(itemInput.GrossValue)
		if err != nil {
			return nil, err
		}

		commissionRate, err := decimal.NewFromString(itemInput.CommissionRate)
		if err != nil {
			return nil, err
		}

		// Converte TenantID para UUID
		tenantUUID, err := uuid.Parse(itemInput.TenantID)
		if err != nil {
			return nil, err
		}

		item, err := entity.NewCommissionItem(
			tenantUUID,
			itemInput.ProfessionalID,
			grossValue,
			commissionRate,
			itemInput.CommissionType,
			itemInput.CommissionSource,
			itemInput.ReferenceDate,
		)
		if err != nil {
			return nil, err
		}

		item.UnitID = itemInput.UnitID
		item.CommandID = itemInput.CommandID
		item.CommandItemID = itemInput.CommandItemID
		item.AppointmentID = itemInput.AppointmentID
		item.ServiceID = itemInput.ServiceID
		item.ServiceName = itemInput.ServiceName
		item.RuleID = itemInput.RuleID
		item.Description = itemInput.Description

		items = append(items, item)
	}

	created, err := uc.commissionItemRepo.CreateBatch(ctx, items)
	if err != nil {
		return nil, err
	}

	return &CreateCommissionItemBatchOutput{CommissionItems: created}, nil
}
