package command

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

// CloseCommandUseCase implementa o fechamento de uma comanda
type CloseCommandUseCase struct {
	repo            port.CommandRepository
	appointmentRepo port.AppointmentRepository
	mapper          *mapper.CommandMapper
}

// NewCloseCommandUseCase cria uma nova instância do use case
func NewCloseCommandUseCase(repo port.CommandRepository, appointmentRepo port.AppointmentRepository, mapper *mapper.CommandMapper) *CloseCommandUseCase {
	return &CloseCommandUseCase{
		repo:            repo,
		appointmentRepo: appointmentRepo,
		mapper:          mapper,
	}
}

// Execute fecha uma comanda (validações + atualização de status)
func (uc *CloseCommandUseCase) Execute(ctx context.Context, commandID, tenantID, userID uuid.UUID, req *dto.CloseCommandRequest) (*dto.CommandResponse, error) {
	// Buscar comanda
	command, err := uc.repo.FindByID(ctx, commandID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get command: %w", err)
	}

	if command == nil {
		return nil, fmt.Errorf("command not found")
	}

	// Configurar opções de fechamento
	if req.DeixarTrocoGorjeta != nil {
		command.DeixarTrocoGorjeta = *req.DeixarTrocoGorjeta
	}

	if req.DeixarSaldoDivida != nil {
		command.DeixarSaldoDivida = *req.DeixarSaldoDivida
	}

	if req.Observacoes != nil {
		command.Observacoes = req.Observacoes
	}

	// Validar se pode fechar
	if err := command.CanClose(); err != nil {
		return nil, fmt.Errorf("cannot close command: %w", err)
	}

	// Fechar comanda (domain logic)
	if err := command.Close(userID); err != nil {
		return nil, fmt.Errorf("failed to close command: %w", err)
	}

	// Persistir
	if err := uc.repo.Update(ctx, command); err != nil {
		return nil, fmt.Errorf("failed to update command: %w", err)
	}

	// Atualizar status do appointment para DONE (se houver appointment_id)
	if command.AppointmentID != nil {
		// Buscar appointment
		appointment, err := uc.appointmentRepo.FindByID(ctx, tenantID.String(), "", command.AppointmentID.String())
		if err == nil && appointment != nil {
			// Atualizar status para DONE
			newStatus := valueobject.AppointmentStatusDone
			appointment.Status = newStatus

			// Persistir atualização do appointment
			if err := uc.appointmentRepo.Update(ctx, appointment); err != nil {
				// Log error mas não falhar o fechamento da comanda
				fmt.Printf("Warning: failed to update appointment status to DONE: %v\n", err)
			}
		}
	} // Buscar comanda fechada
	closed, err := uc.repo.FindByID(ctx, commandID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get closed command: %w", err)
	}

	// Converter para DTO
	response := uc.mapper.ToCommandResponse(closed)
	return response, nil
}
