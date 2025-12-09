package command

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// GetCommandByAppointmentUseCase implementa a busca de uma comanda por appointment_id
type GetCommandByAppointmentUseCase struct {
	repo   port.CommandRepository
	mapper *mapper.CommandMapper
}

// NewGetCommandByAppointmentUseCase cria uma nova instância do use case
func NewGetCommandByAppointmentUseCase(repo port.CommandRepository, mapper *mapper.CommandMapper) *GetCommandByAppointmentUseCase {
	return &GetCommandByAppointmentUseCase{
		repo:   repo,
		mapper: mapper,
	}
}

// Execute busca uma comanda por appointment_id (inclui items e payments)
func (uc *GetCommandByAppointmentUseCase) Execute(ctx context.Context, appointmentID, tenantID uuid.UUID) (*dto.CommandResponse, error) {
	// Buscar comanda pelo appointment_id
	command, err := uc.repo.FindByAppointmentID(ctx, appointmentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get command by appointment: %w", err)
	}

	if command == nil {
		return nil, nil // Retorna nil para indicar que não existe comanda (sem erro)
	}

	// Converter para DTO
	response := uc.mapper.ToCommandResponse(command)
	return response, nil
}
