package appointment

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"go.uber.org/zap"
)

// UpdateAppointmentStatusInput dados de entrada para atualizar status
type UpdateAppointmentStatusInput struct {
	TenantID      string
	AppointmentID string
	NewStatus     valueobject.AppointmentStatus
	Reason        string // Para cancelamento ou no-show
}

// UpdateAppointmentStatusUseCase implementa a atualização de status
type UpdateAppointmentStatusUseCase struct {
	repo   port.AppointmentRepository
	logger *zap.Logger
}

// NewUpdateAppointmentStatusUseCase cria nova instância do use case
func NewUpdateAppointmentStatusUseCase(
	repo port.AppointmentRepository,
	logger *zap.Logger,
) *UpdateAppointmentStatusUseCase {
	return &UpdateAppointmentStatusUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute atualiza o status de um agendamento
func (uc *UpdateAppointmentStatusUseCase) Execute(ctx context.Context, input UpdateAppointmentStatusInput) (*entity.Appointment, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if input.AppointmentID == "" {
		return nil, domain.ErrInvalidID
	}
	if !input.NewStatus.IsValid() {
		return nil, domain.ErrAppointmentInvalidStatus
	}

	// Buscar agendamento
	appointment, err := uc.repo.FindByID(ctx, input.TenantID, input.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar agendamento: %w", err)
	}

	// Aplicar transição de status
	switch input.NewStatus {
	case valueobject.AppointmentStatusConfirmed:
		if err := appointment.Confirm(); err != nil {
			return nil, err
		}
	case valueobject.AppointmentStatusInService:
		if err := appointment.StartService(); err != nil {
			return nil, err
		}
	case valueobject.AppointmentStatusDone:
		if err := appointment.Complete(); err != nil {
			return nil, err
		}
	case valueobject.AppointmentStatusCanceled:
		if err := appointment.Cancel(input.Reason); err != nil {
			return nil, err
		}
	case valueobject.AppointmentStatusNoShow:
		if err := appointment.MarkNoShow(); err != nil {
			return nil, err
		}
	default:
		return nil, domain.ErrAppointmentInvalidStatusTransition
	}

	// Persistir
	if err := uc.repo.Update(ctx, appointment); err != nil {
		return nil, fmt.Errorf("erro ao salvar agendamento: %w", err)
	}

	uc.logger.Info("Status do agendamento atualizado",
		zap.String("tenant_id", input.TenantID),
		zap.String("appointment_id", input.AppointmentID),
		zap.String("new_status", input.NewStatus.String()),
	)

	return appointment, nil
}

// GetAppointmentInput dados de entrada para buscar agendamento
type GetAppointmentInput struct {
	TenantID      string
	AppointmentID string
}

// GetAppointmentUseCase implementa a busca de agendamento
type GetAppointmentUseCase struct {
	repo   port.AppointmentRepository
	logger *zap.Logger
}

// NewGetAppointmentUseCase cria nova instância do use case
func NewGetAppointmentUseCase(
	repo port.AppointmentRepository,
	logger *zap.Logger,
) *GetAppointmentUseCase {
	return &GetAppointmentUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute busca um agendamento por ID
func (uc *GetAppointmentUseCase) Execute(ctx context.Context, input GetAppointmentInput) (*entity.Appointment, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if input.AppointmentID == "" {
		return nil, domain.ErrInvalidID
	}

	appointment, err := uc.repo.FindByID(ctx, input.TenantID, input.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar agendamento: %w", err)
	}

	return appointment, nil
}
