package command

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// RemoveCommandItemUseCase implementa a remoção de item da comanda
type RemoveCommandItemUseCase struct {
	repo   port.CommandRepository
	mapper *mapper.CommandMapper
}

// NewRemoveCommandItemUseCase cria uma nova instância do use case
func NewRemoveCommandItemUseCase(repo port.CommandRepository, mapper *mapper.CommandMapper) *RemoveCommandItemUseCase {
	return &RemoveCommandItemUseCase{
		repo:   repo,
		mapper: mapper,
	}
}

// Execute remove um item da comanda e recalcula totais
func (uc *RemoveCommandItemUseCase) Execute(ctx context.Context, commandID, itemID, tenantID uuid.UUID) (*dto.CommandResponse, error) {
	// Buscar comanda
	command, err := uc.repo.FindByID(ctx, commandID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get command: %w", err)
	}

	if command == nil {
		return nil, fmt.Errorf("command not found")
	}

	// Remover item via domain logic
	if err := command.RemoveItem(itemID); err != nil {
		return nil, fmt.Errorf("failed to remove item from command: %w", err)
	}

	// Persistir remoção
	if err := uc.repo.RemoveItem(ctx, itemID, tenantID); err != nil {
		return nil, fmt.Errorf("failed to persist item removal: %w", err)
	}

	// Atualizar comanda (totais recalculados)
	if err := uc.repo.Update(ctx, command); err != nil {
		return nil, fmt.Errorf("failed to update command: %w", err)
	}

	// Buscar comanda atualizada
	updated, err := uc.repo.FindByID(ctx, commandID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated command: %w", err)
	}

	// Converter para DTO
	response := uc.mapper.ToCommandResponse(updated)
	return response, nil
}

// RemoveCommandPaymentUseCase implementa a remoção de pagamento da comanda
type RemoveCommandPaymentUseCase struct {
	repo   port.CommandRepository
	mapper *mapper.CommandMapper
}

// NewRemoveCommandPaymentUseCase cria uma nova instância do use case
func NewRemoveCommandPaymentUseCase(repo port.CommandRepository, mapper *mapper.CommandMapper) *RemoveCommandPaymentUseCase {
	return &RemoveCommandPaymentUseCase{
		repo:   repo,
		mapper: mapper,
	}
}

// Execute remove um pagamento da comanda e recalcula totais
func (uc *RemoveCommandPaymentUseCase) Execute(ctx context.Context, commandID, paymentID, tenantID uuid.UUID) (*dto.CommandResponse, error) {
	// Buscar comanda
	command, err := uc.repo.FindByID(ctx, commandID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get command: %w", err)
	}

	if command == nil {
		return nil, fmt.Errorf("command not found")
	}

	// Remover pagamento via domain logic
	if err := command.RemovePayment(paymentID); err != nil {
		return nil, fmt.Errorf("failed to remove payment from command: %w", err)
	}

	// Persistir remoção
	if err := uc.repo.RemovePayment(ctx, paymentID, tenantID); err != nil {
		return nil, fmt.Errorf("failed to persist payment removal: %w", err)
	}

	// Atualizar comanda (totais recalculados)
	if err := uc.repo.Update(ctx, command); err != nil {
		return nil, fmt.Errorf("failed to update command: %w", err)
	}

	// Buscar comanda atualizada
	updated, err := uc.repo.FindByID(ctx, commandID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated command: %w", err)
	}

	// Converter para DTO
	response := uc.mapper.ToCommandResponse(updated)
	return response, nil
}
