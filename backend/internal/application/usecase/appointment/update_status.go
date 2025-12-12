package appointment

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UpdateAppointmentStatusInput dados de entrada para atualizar status
type UpdateAppointmentStatusInput struct {
	TenantID      string
	UnitID        string
	AppointmentID string
	NewStatus     valueobject.AppointmentStatus
	Reason        string // Para cancelamento ou no-show
}

// UpdateAppointmentStatusUseCase implementa a atualização de status
type UpdateAppointmentStatusUseCase struct {
	repo        port.AppointmentRepository
	commandRepo port.CommandRepository // Adicionado para validar comanda fechada
	logger      *zap.Logger
}

// NewUpdateAppointmentStatusUseCase cria nova instância do use case
func NewUpdateAppointmentStatusUseCase(
	repo port.AppointmentRepository,
	commandRepo port.CommandRepository,
	logger *zap.Logger,
) *UpdateAppointmentStatusUseCase {
	return &UpdateAppointmentStatusUseCase{
		repo:        repo,
		commandRepo: commandRepo,
		logger:      logger,
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
	appointment, err := uc.repo.FindByID(ctx, input.TenantID, input.UnitID, input.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar agendamento: %w", err)
	}

	// Aplicar transição de status
	switch input.NewStatus {
	case valueobject.AppointmentStatusConfirmed:
		if err := appointment.Confirm(); err != nil {
			return nil, err
		}
	case valueobject.AppointmentStatusCheckedIn:
		// Cliente chegou para o atendimento
		if err := appointment.CheckIn(); err != nil {
			return nil, err
		}
	case valueobject.AppointmentStatusInService:
		if err := appointment.StartService(); err != nil {
			return nil, err
		}
	case valueobject.AppointmentStatusAwaitingPayment:
		// Serviços finalizados, aguardando pagamento
		if err := appointment.FinishService(); err != nil {
			return nil, err
		}
	case valueobject.AppointmentStatusDone:
		// BLOQUEIO: Não permitir marcar como DONE sem comanda fechada
		// Use FinalizarComandaIntegrada para finalizar com integração financeira
		if appointment.CommandID == "" {
			return nil, fmt.Errorf("agendamento não possui comanda vinculada: use a finalização via comanda")
		}

		// Verificar se comanda existe e está fechada
		tenantUUID, err := uuid.Parse(input.TenantID)
		if err != nil {
			return nil, fmt.Errorf("tenant_id inválido: %w", err)
		}
		commandUUID, err := uuid.Parse(appointment.CommandID)
		if err != nil {
			return nil, fmt.Errorf("command_id inválido: %w", err)
		}

		command, err := uc.commandRepo.FindByID(ctx, commandUUID, tenantUUID)
		if err != nil {
			return nil, fmt.Errorf("comanda não encontrada: %w", err)
		}
		if command.Status != entity.CommandStatusClosed {
			return nil, fmt.Errorf("comanda deve estar fechada para finalizar agendamento (status atual: %s)", command.Status)
		}

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
	UnitID        string
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

	appointment, err := uc.repo.FindByID(ctx, input.TenantID, input.UnitID, input.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar agendamento: %w", err)
	}

	return appointment, nil
}
