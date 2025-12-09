package port

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
)

// CommandRepository define operações de persistência para comandas
type CommandRepository interface {
	// Create cria uma nova comanda com seus itens
	Create(ctx context.Context, command *entity.Command) error

	// FindByID busca uma comanda por ID (inclui items e payments)
	FindByID(ctx context.Context, commandID, tenantID uuid.UUID) (*entity.Command, error)

	// FindByAppointmentID busca uma comanda pelo ID do agendamento
	FindByAppointmentID(ctx context.Context, appointmentID, tenantID uuid.UUID) (*entity.Command, error)

	// Update atualiza uma comanda existente
	Update(ctx context.Context, command *entity.Command) error

	// Delete remove uma comanda (soft delete via status)
	Delete(ctx context.Context, commandID, tenantID uuid.UUID) error

	// List lista comandas com filtros
	List(ctx context.Context, tenantID uuid.UUID, filters CommandFilters) ([]*entity.Command, error)

	// AddItem adiciona um item à comanda
	AddItem(ctx context.Context, item *entity.CommandItem) error

	// UpdateItem atualiza um item da comanda
	UpdateItem(ctx context.Context, item *entity.CommandItem) error

	// RemoveItem remove um item da comanda
	RemoveItem(ctx context.Context, itemID, tenantID uuid.UUID) error

	// GetItems busca todos os itens de uma comanda
	GetItems(ctx context.Context, commandID, tenantID uuid.UUID) ([]entity.CommandItem, error)

	// AddPayment adiciona um pagamento à comanda
	AddPayment(ctx context.Context, payment *entity.CommandPayment) error

	// RemovePayment remove um pagamento da comanda
	RemovePayment(ctx context.Context, paymentID, tenantID uuid.UUID) error

	// GetPayments busca todos os pagamentos de uma comanda
	GetPayments(ctx context.Context, commandID, tenantID uuid.UUID) ([]entity.CommandPayment, error)
}

// CommandFilters representa filtros para busca de comandas
type CommandFilters struct {
	Status     *entity.CommandStatus
	CustomerID *uuid.UUID
	DateFrom   *string // YYYY-MM-DD
	DateTo     *string // YYYY-MM-DD
	Limit      int
	Offset     int
}
