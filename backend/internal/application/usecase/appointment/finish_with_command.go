package appointment

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// FinishServiceWithCommandInput dados de entrada para finalizar serviço
type FinishServiceWithCommandInput struct {
	TenantID      string
	AppointmentID string
}

// FinishServiceWithCommandOutput resultado da operação
type FinishServiceWithCommandOutput struct {
	Appointment *entity.Appointment
	Command     *entity.Command
	Created     bool // indica se a comanda foi criada (true) ou já existia (false)
}

// FinishServiceWithCommandUseCase implementa a finalização de serviço com criação automática de comanda
type FinishServiceWithCommandUseCase struct {
	appointmentRepo port.AppointmentRepository
	commandRepo     port.CommandRepository
	logger          *zap.Logger
}

// NewFinishServiceWithCommandUseCase cria nova instância do use case
func NewFinishServiceWithCommandUseCase(
	appointmentRepo port.AppointmentRepository,
	commandRepo port.CommandRepository,
	logger *zap.Logger,
) *FinishServiceWithCommandUseCase {
	return &FinishServiceWithCommandUseCase{
		appointmentRepo: appointmentRepo,
		commandRepo:     commandRepo,
		logger:          logger,
	}
}

// Execute finaliza o atendimento e cria comanda automaticamente se não existir
func (uc *FinishServiceWithCommandUseCase) Execute(ctx context.Context, input FinishServiceWithCommandInput) (*FinishServiceWithCommandOutput, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if input.AppointmentID == "" {
		return nil, domain.ErrInvalidID
	}

	// 1. Buscar agendamento
	appointment, err := uc.appointmentRepo.FindByID(ctx, input.TenantID, input.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar agendamento: %w", err)
	}

	// 2. Tentar transição de status
	if err := appointment.FinishService(); err != nil {
		return nil, err
	}

	output := &FinishServiceWithCommandOutput{
		Appointment: appointment,
		Created:     false,
	}

	// 3. Verificar se já tem comanda
	if appointment.CommandID != "" {
		// Já existe comanda, apenas buscar
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
			uc.logger.Warn("Comanda existente não encontrada, criando nova",
				zap.String("command_id", appointment.CommandID),
				zap.Error(err),
			)
			// Continuar para criar nova comanda
		} else {
			output.Command = command
			output.Created = false

			// Persistir appointment atualizado
			if err := uc.appointmentRepo.Update(ctx, appointment); err != nil {
				return nil, fmt.Errorf("erro ao atualizar agendamento: %w", err)
			}

			return output, nil
		}
	}

	// 4. Criar nova comanda
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}
	customerUUID, err := uuid.Parse(appointment.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer_id inválido: %w", err)
	}
	appointmentUUID, err := uuid.Parse(appointment.ID)
	if err != nil {
		return nil, fmt.Errorf("appointment_id inválido: %w", err)
	}

	command, err := entity.NewCommand(tenantUUID, customerUUID, &appointmentUUID)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar comanda: %w", err)
	}

	// 5. Adicionar itens da comanda (serviços do agendamento)
	for _, svc := range appointment.Services {
		serviceUUID, err := uuid.Parse(svc.ServiceID)
		if err != nil {
			uc.logger.Warn("ServiceID inválido, ignorando serviço",
				zap.String("service_id", svc.ServiceID),
				zap.Error(err),
			)
			continue
		}

		// Converter decimal.Decimal para float64
		preco, _ := svc.PriceAtBooking.Value().Float64()

		item, err := entity.NewCommandItem(
			command.ID,
			entity.CommandItemTypeServico,
			serviceUUID,
			svc.ServiceName,
			preco,
			1, // quantidade
		)
		if err != nil {
			uc.logger.Warn("Erro ao criar item da comanda, ignorando serviço",
				zap.String("service_id", svc.ServiceID),
				zap.Error(err),
			)
			continue
		}
		command.Items = append(command.Items, *item)
	}

	// 6. Calcular totais da comanda
	command.RecalculateTotals()

	// 7. Persistir comanda
	if err := uc.commandRepo.Create(ctx, command); err != nil {
		return nil, fmt.Errorf("erro ao criar comanda: %w", err)
	}

	// 8. Atualizar agendamento com command_id
	appointment.CommandID = command.ID.String()

	// 9. Persistir agendamento atualizado
	if err := uc.appointmentRepo.Update(ctx, appointment); err != nil {
		return nil, fmt.Errorf("erro ao atualizar agendamento: %w", err)
	}

	output.Command = command
	output.Created = true

	uc.logger.Info("Serviço finalizado com comanda criada automaticamente",
		zap.String("tenant_id", input.TenantID),
		zap.String("appointment_id", input.AppointmentID),
		zap.String("command_id", command.ID.String()),
		zap.Float64("total", command.Total),
	)

	return output, nil
}
