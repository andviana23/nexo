package command

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// GetCommandUseCase implementa a busca de uma comanda por ID
type GetCommandUseCase struct {
	repo   port.CommandRepository
	mapper *mapper.CommandMapper
}

// NewGetCommandUseCase cria uma nova instância do use case
func NewGetCommandUseCase(repo port.CommandRepository, mapper *mapper.CommandMapper) *GetCommandUseCase {
	return &GetCommandUseCase{
		repo:   repo,
		mapper: mapper,
	}
}

// Execute busca uma comanda por ID (inclui items e payments)
func (uc *GetCommandUseCase) Execute(ctx context.Context, commandID, tenantID uuid.UUID) (*dto.CommandResponse, error) {
	// Buscar comanda
	command, err := uc.repo.FindByID(ctx, commandID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get command: %w", err)
	}

	if command == nil {
		return nil, fmt.Errorf("command not found")
	}

	// Converter para DTO
	response := uc.mapper.ToCommandResponse(command)
	return response, nil
}
