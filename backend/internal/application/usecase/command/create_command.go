package command

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// CreateCommandUseCase implementa a criação de uma comanda
type CreateCommandUseCase struct {
	repo   port.CommandRepository
	mapper *mapper.CommandMapper
}

// NewCreateCommandUseCase cria uma nova instância do use case
func NewCreateCommandUseCase(repo port.CommandRepository, mapper *mapper.CommandMapper) *CreateCommandUseCase {
	return &CreateCommandUseCase{
		repo:   repo,
		mapper: mapper,
	}
}

// Execute cria uma nova comanda
func (uc *CreateCommandUseCase) Execute(ctx context.Context, tenantID uuid.UUID, req *dto.CreateCommandRequest) (*dto.CommandResponse, error) {
	// Converter DTO para Entity
	command, err := uc.mapper.FromCreateCommandRequest(req, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to map request: %w", err)
	}

	// Persistir comanda
	if err := uc.repo.Create(ctx, command); err != nil {
		return nil, fmt.Errorf("failed to create command: %w", err)
	}

	// Buscar comanda criada para retornar completa
	created, err := uc.repo.FindByID(ctx, command.ID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created command: %w", err)
	}

	// Converter para DTO de resposta
	response := uc.mapper.ToCommandResponse(created)
	return response, nil
}
