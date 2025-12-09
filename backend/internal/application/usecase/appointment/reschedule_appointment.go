package appointment

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// RescheduleAppointmentInput dados de entrada para reagendar
type RescheduleAppointmentInput struct {
	TenantID       string
	AppointmentID  string
	NewStartTime   time.Time
	ProfessionalID string // Opcional: trocar de profissional
}

// RescheduleAppointmentUseCase implementa o reagendamento
type RescheduleAppointmentUseCase struct {
	repo               port.AppointmentRepository
	professionalReader port.ProfessionalReader
	logger             *zap.Logger
}

// NewRescheduleAppointmentUseCase cria nova instância do use case
func NewRescheduleAppointmentUseCase(
	repo port.AppointmentRepository,
	professionalReader port.ProfessionalReader,
	logger *zap.Logger,
) *RescheduleAppointmentUseCase {
	return &RescheduleAppointmentUseCase{
		repo:               repo,
		professionalReader: professionalReader,
		logger:             logger,
	}
}

// Execute reagenda um agendamento
func (uc *RescheduleAppointmentUseCase) Execute(ctx context.Context, input RescheduleAppointmentInput) (*entity.Appointment, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if input.AppointmentID == "" {
		return nil, domain.ErrInvalidID
	}
	if input.NewStartTime.IsZero() {
		return nil, domain.ErrAppointmentStartTimeRequired
	}

	// Buscar agendamento
	appointment, err := uc.repo.FindByID(ctx, input.TenantID, input.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar agendamento: %w", err)
	}

	// Guardar dados para verificação de conflito
	oldStartTime := appointment.StartTime
	professionalID := appointment.ProfessionalID

	if input.ProfessionalID != "" && input.ProfessionalID != professionalID {
		exists, err := uc.professionalReader.Exists(ctx, input.TenantID, input.ProfessionalID)
		if err != nil {
			return nil, fmt.Errorf("erro ao verificar profissional: %w", err)
		}
		if !exists {
			return nil, domain.ErrAppointmentProfessionalNotFound
		}
		professionalID = input.ProfessionalID
	}

	// Reagendar (atualiza start_time e end_time mantendo duração)
	if err := appointment.Reschedule(input.NewStartTime); err != nil {
		return nil, fmt.Errorf("erro ao reagendar: %w", err)
	}

	// Se mudou de profissional, atualizar
	if professionalID != appointment.ProfessionalID {
		appointment.ProfessionalID = professionalID
	}

	// Verificar conflito de horário
	hasConflict, err := uc.repo.CheckConflict(
		ctx,
		input.TenantID,
		professionalID,
		appointment.StartTime,
		appointment.EndTime,
		appointment.ID, // Excluir o próprio agendamento da verificação
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar conflito: %w", err)
	}
	if hasConflict {
		return nil, domain.ErrAppointmentConflict
	}

	// Verificar conflito com horários bloqueados
	hasBlockedConflict, err := uc.repo.CheckBlockedTimeConflict(
		ctx,
		input.TenantID,
		professionalID,
		appointment.StartTime,
		appointment.EndTime,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar bloqueio: %w", err)
	}
	if hasBlockedConflict {
		return nil, domain.ErrAppointmentBlockedTimeConflict
	}

	// Verificar intervalo mínimo entre agendamentos (RN-AGE-003: 10 minutos)
	const minimumIntervalMinutes = 10
	hasIntervalConflict, err := uc.repo.CheckMinimumIntervalConflict(
		ctx,
		input.TenantID,
		professionalID,
		appointment.StartTime,
		appointment.EndTime,
		appointment.ID,
		minimumIntervalMinutes,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar intervalo mínimo: %w", err)
	}
	if hasIntervalConflict {
		return nil, domain.ErrAppointmentMinimumInterval
	}

	// Persistir
	if err := uc.repo.Update(ctx, appointment); err != nil {
		return nil, fmt.Errorf("erro ao salvar agendamento: %w", err)
	}

	uc.logger.Info("Agendamento reagendado",
		zap.String("tenant_id", input.TenantID),
		zap.String("appointment_id", input.AppointmentID),
		zap.Time("old_start_time", oldStartTime),
		zap.Time("new_start_time", appointment.StartTime),
	)

	return appointment, nil
}
