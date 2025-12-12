package command

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// ListCommandsInput representa os parâmetros de entrada para listar comandas
type ListCommandsInput struct {
	TenantID   uuid.UUID
	Status     *string
	CustomerID *string
	DateFrom   *string // YYYY-MM-DD
	DateTo     *string // YYYY-MM-DD
	Page       int
	PageSize   int
}

// ListCommandsUseCase implementa a listagem de comandas com filtros
type ListCommandsUseCase struct {
	repo   port.CommandRepository
	mapper *mapper.CommandMapper
}

// NewListCommandsUseCase cria uma nova instância do use case
func NewListCommandsUseCase(
	repo port.CommandRepository,
	mapper *mapper.CommandMapper,
) *ListCommandsUseCase {
	return &ListCommandsUseCase{
		repo:   repo,
		mapper: mapper,
	}
}

// Execute lista comandas com filtros e paginação
func (uc *ListCommandsUseCase) Execute(ctx context.Context, input ListCommandsInput) (*dto.CommandListResponse, error) {
	// Validar paginação
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20 // Default
	}

	// Converter filtros para port.CommandFilters
	filters := port.CommandFilters{
		Limit:  input.PageSize,
		Offset: (input.Page - 1) * input.PageSize,
	}

	// Status filter
	if input.Status != nil && *input.Status != "" {
		status := entity.CommandStatus(*input.Status)
		filters.Status = &status
	}

	// Customer ID filter
	if input.CustomerID != nil && *input.CustomerID != "" {
		customerID, err := uuid.Parse(*input.CustomerID)
		if err != nil {
			return nil, fmt.Errorf("invalid customer_id: %w", err)
		}
		filters.CustomerID = &customerID
	}

	// Date filters
	if input.DateFrom != nil && *input.DateFrom != "" {
		filters.DateFrom = input.DateFrom
	}
	if input.DateTo != nil && *input.DateTo != "" {
		filters.DateTo = input.DateTo
	}

	// Buscar comandas
	commands, err := uc.repo.List(ctx, input.TenantID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list commands: %w", err)
	}

	// Converter para DTOs
	var commandDTOs []dto.CommandResponse
	for _, cmd := range commands {
		commandDTOs = append(commandDTOs, *uc.mapper.ToCommandResponse(cmd))
	}

	// Calcular total de páginas (simplificado - idealmente deveria ter um Count no repo)
	totalPages := 1
	if len(commandDTOs) == input.PageSize {
		totalPages = input.Page + 1 // Pode haver mais páginas
	}

	return &dto.CommandListResponse{
		Commands:   commandDTOs,
		Total:      len(commandDTOs), // TODO: Implementar Count no repository
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalPages: totalPages,
	}, nil
}
