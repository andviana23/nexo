package appointment

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// CancelAppointmentInput dados de entrada para cancelar agendamento
type CancelAppointmentInput struct {
	TenantID      string
	AppointmentID string
	Reason        string
}

// CancelAppointmentUseCase implementa o cancelamento de agendamentos
type CancelAppointmentUseCase struct {
	repo   port.AppointmentRepository
	logger *zap.Logger
}

// NewCancelAppointmentUseCase cria nova instância do use case
func NewCancelAppointmentUseCase(
	repo port.AppointmentRepository,
	logger *zap.Logger,
) *CancelAppointmentUseCase {
	return &CancelAppointmentUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute cancela um agendamento
func (uc *CancelAppointmentUseCase) Execute(ctx context.Context, input CancelAppointmentInput) (*entity.Appointment, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if input.AppointmentID == "" {
		return nil, domain.ErrInvalidID
	}

	// Buscar agendamento
	appointment, err := uc.repo.FindByID(ctx, input.TenantID, input.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar agendamento: %w", err)
	}

	// Cancelar
	if err := appointment.Cancel(input.Reason); err != nil {
		return nil, fmt.Errorf("erro ao cancelar agendamento: %w", err)
	}

	// Persistir
	if err := uc.repo.Update(ctx, appointment); err != nil {
		return nil, fmt.Errorf("erro ao salvar agendamento: %w", err)
	}

	uc.logger.Info("Agendamento cancelado",
		zap.String("tenant_id", input.TenantID),
		zap.String("appointment_id", input.AppointmentID),
		zap.String("reason", input.Reason),
	)

	return appointment, nil
}
