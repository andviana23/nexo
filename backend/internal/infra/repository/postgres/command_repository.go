package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// CommandRepository implementa port.CommandRepository
type CommandRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

// NewCommandRepository cria uma nova instância do repositório
func NewCommandRepository(queries *db.Queries, pool *pgxpool.Pool) *CommandRepository {
	return &CommandRepository{
		queries: queries,
		pool:    pool,
	}
}

// Create cria uma nova comanda com transação
// G-002: Gera automaticamente o número sequencial da comanda (CMD-YYYY-NNNNN)
func (r *CommandRepository) Create(ctx context.Context, command *entity.Command) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// G-002: Gerar número sequencial se não fornecido
	if command.Numero == nil || *command.Numero == "" {
		nextNum, err := qtx.GetNextCommandNumber(ctx, uuidToUUID(command.TenantID))
		if err != nil {
			return fmt.Errorf("failed to get next command number: %w", err)
		}
		year := time.Now().Year()
		numero := fmt.Sprintf("CMD-%d-%05d", year, nextNum)
		command.Numero = &numero
	}

	// Criar comanda
	params := db.CreateCommandParams{
		ID:            uuidToUUID(command.ID),
		TenantID:      uuidToUUID(command.TenantID),
		CustomerID:    uuidToUUID(command.CustomerID),
		Status:        string(command.Status),
		Subtotal:      decimalFromFloat(command.Subtotal),
		Desconto:      decimalFromFloat(command.Desconto),
		Total:         decimalFromFloat(command.Total),
		TotalRecebido: decimalFromFloat(command.TotalRecebido),
		Troco:         decimalFromFloat(command.Troco),
		SaldoDevedor:  decimalFromFloat(command.SaldoDevedor),
		CriadoEm:      timestampToTimestamptz(command.CriadoEm),
		AtualizadoEm:  timestampToTimestamptz(command.AtualizadoEm),
	}

	if command.AppointmentID != nil {
		params.AppointmentID = uuidToUUID(*command.AppointmentID)
	}

	if command.Numero != nil {
		params.Numero = ptrString(*command.Numero)
	}

	if command.Observacoes != nil {
		params.Observacoes = ptrString(*command.Observacoes)
	}

	params.DeixarTrocoGorjeta = ptrBool(command.DeixarTrocoGorjeta)
	params.DeixarSaldoDivida = ptrBool(command.DeixarSaldoDivida)

	_, err = qtx.CreateCommand(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create command: %w", err)
	}

	// Criar itens
	for _, item := range command.Items {
		itemParams := db.CreateCommandItemParams{
			ID:                 uuidToUUID(item.ID),
			CommandID:          uuidToUUID(command.ID),
			Tipo:               string(item.Tipo),
			ItemID:             uuidToUUID(item.ItemID),
			Descricao:          item.Descricao,
			PrecoUnitario:      decimalFromFloat(item.PrecoUnitario),
			Quantidade:         int32(item.Quantidade),
			DescontoValor:      decimalFromFloat(item.DescontoValor),
			DescontoPercentual: decimalFromFloat(item.DescontoPercentual),
			PrecoFinal:         decimalFromFloat(item.PrecoFinal),
			CriadoEm:           timestampToTimestamptz(item.CriadoEm),
		}

		if item.Observacoes != nil {
			itemParams.Observacoes = ptrString(*item.Observacoes)
		}

		_, err = qtx.CreateCommandItem(ctx, itemParams)
		if err != nil {
			return fmt.Errorf("failed to create command item: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// FindByID busca uma comanda por ID
func (r *CommandRepository) FindByID(ctx context.Context, commandID, tenantID uuid.UUID) (*entity.Command, error) {
	// Buscar comanda
	dbCommand, err := r.queries.GetCommandByID(ctx, db.GetCommandByIDParams{
		ID:       uuidToUUID(commandID),
		TenantID: uuidToUUID(tenantID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get command: %w", err)
	}

	// Converter para entidade
	command := r.dbCommandToEntity(dbCommand)

	// Buscar itens
	dbItems, err := r.queries.GetCommandItems(ctx, db.GetCommandItemsParams{
		CommandID: uuidToUUID(commandID),
		TenantID:  uuidToUUID(tenantID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get command items: %w", err)
	}

	for _, dbItem := range dbItems {
		command.Items = append(command.Items, r.dbItemToEntity(dbItem))
	}

	// Buscar pagamentos
	dbPayments, err := r.queries.GetCommandPayments(ctx, db.GetCommandPaymentsParams{
		CommandID: uuidToUUID(commandID),
		TenantID:  uuidToUUID(tenantID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get command payments: %w", err)
	}

	for _, dbPayment := range dbPayments {
		command.Payments = append(command.Payments, r.dbPaymentToEntity(dbPayment))
	}

	return command, nil
}

// FindByAppointmentID busca uma comanda pelo ID do agendamento
func (r *CommandRepository) FindByAppointmentID(ctx context.Context, appointmentID, tenantID uuid.UUID) (*entity.Command, error) {
	dbCommand, err := r.queries.GetCommandByAppointmentID(ctx, db.GetCommandByAppointmentIDParams{
		AppointmentID: uuidToUUID(appointmentID),
		TenantID:      uuidToUUID(tenantID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get command by appointment: %w", err)
	}

	// Reutilizar FindByID para carregar items e payments
	return r.FindByID(ctx, uuidFromUUID(dbCommand.ID), tenantID)
}

// Update atualiza uma comanda
func (r *CommandRepository) Update(ctx context.Context, command *entity.Command) error {
	params := db.UpdateCommandParams{
		ID:            uuidToUUID(command.ID),
		TenantID:      uuidToUUID(command.TenantID),
		Status:        string(command.Status),
		Subtotal:      decimalFromFloat(command.Subtotal),
		Desconto:      decimalFromFloat(command.Desconto),
		Total:         decimalFromFloat(command.Total),
		TotalRecebido: decimalFromFloat(command.TotalRecebido),
		Troco:         decimalFromFloat(command.Troco),
		SaldoDevedor:  decimalFromFloat(command.SaldoDevedor),
	}

	if command.Observacoes != nil {
		params.Observacoes = ptrString(*command.Observacoes)
	}

	params.DeixarTrocoGorjeta = ptrBool(command.DeixarTrocoGorjeta)
	params.DeixarSaldoDivida = ptrBool(command.DeixarSaldoDivida)

	if command.FechadoEm != nil {
		params.FechadoEm = pgtype.Timestamptz{Time: *command.FechadoEm, Valid: true}
	}

	if command.FechadoPor != nil {
		params.FechadoPor = uuidToUUID(*command.FechadoPor)
	}

	_, err := r.queries.UpdateCommand(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update command: %w", err)
	}

	return nil
}

// Delete remove uma comanda (soft delete via status CANCELED)
func (r *CommandRepository) Delete(ctx context.Context, commandID, tenantID uuid.UUID) error {
	if err := r.queries.DeleteCommand(ctx, db.DeleteCommandParams{
		ID:       uuidToUUID(commandID),
		TenantID: uuidToUUID(tenantID),
	}); err != nil {
		return fmt.Errorf("failed to delete command: %w", err)
	}
	return nil
}

// List lista comandas com filtros
func (r *CommandRepository) List(ctx context.Context, tenantID uuid.UUID, filters port.CommandFilters) ([]*entity.Command, error) {
	params := db.ListCommandsParams{
		TenantID: uuidToUUID(tenantID),
		Limit:    int32(filters.Limit),
		Offset:   int32(filters.Offset),
	}

	// Column2 = status
	if filters.Status != nil {
		params.Column2 = string(*filters.Status)
	}

	// Column3 = customer_id
	if filters.CustomerID != nil {
		params.Column3 = uuidToUUID(*filters.CustomerID)
	}

	// Column4 = date_from
	if filters.DateFrom != nil {
		// Parse string date to pgtype.Date
		params.Column4 = pgtype.Date{Valid: true} // TODO: parse date properly
	}

	// Column5 = date_to
	if filters.DateTo != nil {
		// Parse string date to pgtype.Date
		params.Column5 = pgtype.Date{Valid: true} // TODO: parse date properly
	}

	dbCommands, err := r.queries.ListCommands(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list commands: %w", err)
	}

	var commands []*entity.Command
	for _, dbCmd := range dbCommands {
		commands = append(commands, r.dbCommandToEntity(dbCmd))
	}

	return commands, nil
}

// AddItem adiciona um item à comanda
func (r *CommandRepository) AddItem(ctx context.Context, item *entity.CommandItem) error {
	params := db.CreateCommandItemParams{
		ID:                 uuidToUUID(item.ID),
		CommandID:          uuidToUUID(item.CommandID),
		Tipo:               string(item.Tipo),
		ItemID:             uuidToUUID(item.ItemID),
		Descricao:          item.Descricao,
		PrecoUnitario:      decimalFromFloat(item.PrecoUnitario),
		Quantidade:         int32(item.Quantidade),
		DescontoValor:      decimalFromFloat(item.DescontoValor),
		DescontoPercentual: decimalFromFloat(item.DescontoPercentual),
		PrecoFinal:         decimalFromFloat(item.PrecoFinal),
		CriadoEm:           timestampToTimestamptz(item.CriadoEm),
	}

	if item.Observacoes != nil {
		params.Observacoes = ptrString(*item.Observacoes)
	}

	_, err := r.queries.CreateCommandItem(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}

	return nil
}

// UpdateItem atualiza um item da comanda
func (r *CommandRepository) UpdateItem(ctx context.Context, item *entity.CommandItem) error {
	params := db.UpdateCommandItemParams{
		ID:                 uuidToUUID(item.ID),
		PrecoUnitario:      decimalFromFloat(item.PrecoUnitario),
		Quantidade:         int32(item.Quantidade),
		DescontoValor:      decimalFromFloat(item.DescontoValor),
		DescontoPercentual: decimalFromFloat(item.DescontoPercentual),
		PrecoFinal:         decimalFromFloat(item.PrecoFinal),
	}

	if item.Observacoes != nil {
		params.Observacoes = ptrString(*item.Observacoes)
	}

	_, err := r.queries.UpdateCommandItem(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

// RemoveItem remove um item da comanda
func (r *CommandRepository) RemoveItem(ctx context.Context, itemID, tenantID uuid.UUID) error {
	// Primeiro buscar o item para pegar o command_id e validar tenant
	dbItem, err := r.queries.GetCommandItemByID(ctx, db.GetCommandItemByIDParams{
		ID: uuidToUUID(itemID),
	})
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	// Validar tenant através da comanda
	dbCommand, err := r.queries.GetCommandByID(ctx, db.GetCommandByIDParams{
		ID:       dbItem.CommandID,
		TenantID: uuidToUUID(tenantID),
	})
	if err != nil {
		return fmt.Errorf("failed to validate tenant: %w", err)
	}
	if dbCommand.ID.Bytes != dbItem.CommandID.Bytes {
		return fmt.Errorf("tenant mismatch")
	}

	if err := r.queries.DeleteCommandItem(ctx, db.DeleteCommandItemParams{
		ID:       uuidToUUID(itemID),
		TenantID: uuidToUUID(tenantID),
	}); err != nil {
		return fmt.Errorf("failed to remove item: %w", err)
	}

	return nil
}

// GetItems busca todos os itens de uma comanda
func (r *CommandRepository) GetItems(ctx context.Context, commandID, tenantID uuid.UUID) ([]entity.CommandItem, error) {
	dbItems, err := r.queries.GetCommandItems(ctx, db.GetCommandItemsParams{
		CommandID: uuidToUUID(commandID),
		TenantID:  uuidToUUID(tenantID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	var items []entity.CommandItem
	for _, dbItem := range dbItems {
		items = append(items, r.dbItemToEntity(dbItem))
	}

	return items, nil
}

// AddPayment adiciona um pagamento à comanda
func (r *CommandRepository) AddPayment(ctx context.Context, payment *entity.CommandPayment) error {
	params := db.CreateCommandPaymentParams{
		ID:              uuidToUUID(payment.ID),
		CommandID:       uuidToUUID(payment.CommandID),
		MeioPagamentoID: uuidToUUID(payment.MeioPagamentoID),
		ValorRecebido:   decimalFromFloat(payment.ValorRecebido),
		TaxaPercentual:  decimalFromFloat(payment.TaxaPercentual),
		TaxaFixa:        decimalFromFloat(payment.TaxaFixa),
		ValorLiquido:    decimalFromFloat(payment.ValorLiquido),
		CriadoEm:        pgtype.Timestamptz{Time: payment.CriadoEm, Valid: true},
		CriadoPor:       ptrUUIDToUUID(payment.CriadoPor),
	}

	if payment.Observacoes != nil {
		params.Observacoes = ptrString(*payment.Observacoes)
	}

	if _, err := r.queries.CreateCommandPayment(ctx, params); err != nil {
		return fmt.Errorf("failed to add payment: %w", err)
	}

	return nil
}

// RemovePayment remove um pagamento da comanda
func (r *CommandRepository) RemovePayment(ctx context.Context, paymentID, tenantID uuid.UUID) error {
	// Validar tenant através da comanda
	dbPayment, err := r.queries.GetCommandPaymentByID(ctx, db.GetCommandPaymentByIDParams{
		ID:       uuidToUUID(paymentID),
		TenantID: uuidToUUID(tenantID),
	})
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	dbCommand, err := r.queries.GetCommandByID(ctx, db.GetCommandByIDParams{
		ID:       dbPayment.CommandID,
		TenantID: uuidToUUID(tenantID),
	})
	if err != nil {
		return fmt.Errorf("failed to validate tenant: %w", err)
	}
	if dbCommand.ID.Bytes != dbPayment.CommandID.Bytes {
		return fmt.Errorf("tenant mismatch")
	}

	if err := r.queries.DeleteCommandPayment(ctx, db.DeleteCommandPaymentParams{
		ID:       uuidToUUID(paymentID),
		TenantID: uuidToUUID(tenantID),
	}); err != nil {
		return fmt.Errorf("failed to remove payment: %w", err)
	}

	return nil
}

// GetPayments busca todos os pagamentos de uma comanda
func (r *CommandRepository) GetPayments(ctx context.Context, commandID, tenantID uuid.UUID) ([]entity.CommandPayment, error) {
	dbPayments, err := r.queries.GetCommandPayments(ctx, db.GetCommandPaymentsParams{
		CommandID: uuidToUUID(commandID),
		TenantID:  uuidToUUID(tenantID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}

	var payments []entity.CommandPayment
	for _, dbPayment := range dbPayments {
		payments = append(payments, r.dbPaymentToEntity(dbPayment))
	}

	return payments, nil
}

// ============================================================================
// Conversion Helpers
// ============================================================================

func (r *CommandRepository) dbCommandToEntity(dbCmd db.Command) *entity.Command {
	cmd := &entity.Command{
		ID:            uuidFromUUID(dbCmd.ID),
		TenantID:      uuidFromUUID(dbCmd.TenantID),
		CustomerID:    uuidFromUUID(dbCmd.CustomerID),
		Status:        entity.CommandStatus(dbCmd.Status),
		Subtotal:      floatFromDecimal(dbCmd.Subtotal),
		Desconto:      floatFromDecimal(dbCmd.Desconto),
		Total:         floatFromDecimal(dbCmd.Total),
		TotalRecebido: floatFromDecimal(dbCmd.TotalRecebido),
		Troco:         floatFromDecimal(dbCmd.Troco),
		SaldoDevedor:  floatFromDecimal(dbCmd.SaldoDevedor),
		CriadoEm:      dbCmd.CriadoEm.Time,
		AtualizadoEm:  dbCmd.AtualizadoEm.Time,
		Items:         []entity.CommandItem{},
		Payments:      []entity.CommandPayment{},
	}

	if dbCmd.AppointmentID.Valid {
		aid := uuidFromUUID(dbCmd.AppointmentID)
		cmd.AppointmentID = &aid
	}

	if dbCmd.Numero != nil {
		cmd.Numero = dbCmd.Numero
	}

	if dbCmd.Observacoes != nil {
		cmd.Observacoes = dbCmd.Observacoes
	}

	if dbCmd.DeixarTrocoGorjeta != nil {
		cmd.DeixarTrocoGorjeta = *dbCmd.DeixarTrocoGorjeta
	}

	if dbCmd.DeixarSaldoDivida != nil {
		cmd.DeixarSaldoDivida = *dbCmd.DeixarSaldoDivida
	}

	if dbCmd.FechadoEm.Valid {
		cmd.FechadoEm = &dbCmd.FechadoEm.Time
	}

	if dbCmd.FechadoPor.Valid {
		fid := uuidFromUUID(dbCmd.FechadoPor)
		cmd.FechadoPor = &fid
	}

	return cmd
}

func (r *CommandRepository) dbItemToEntity(dbItem db.CommandItem) entity.CommandItem {
	item := entity.CommandItem{
		ID:                 uuidFromUUID(dbItem.ID),
		CommandID:          uuidFromUUID(dbItem.CommandID),
		Tipo:               entity.CommandItemType(dbItem.Tipo),
		ItemID:             uuidFromUUID(dbItem.ItemID),
		Descricao:          dbItem.Descricao,
		PrecoUnitario:      floatFromDecimal(dbItem.PrecoUnitario),
		Quantidade:         int(dbItem.Quantidade),
		DescontoValor:      floatFromDecimal(dbItem.DescontoValor),
		DescontoPercentual: floatFromDecimal(dbItem.DescontoPercentual),
		PrecoFinal:         floatFromDecimal(dbItem.PrecoFinal),
		CriadoEm:           dbItem.CriadoEm.Time,
	}

	if dbItem.Observacoes != nil {
		item.Observacoes = dbItem.Observacoes
	}

	return item
}

func (r *CommandRepository) dbPaymentToEntity(dbPayment db.CommandPayment) entity.CommandPayment {
	payment := entity.CommandPayment{
		ID:              uuidFromUUID(dbPayment.ID),
		CommandID:       uuidFromUUID(dbPayment.CommandID),
		MeioPagamentoID: uuidFromUUID(dbPayment.MeioPagamentoID),
		ValorRecebido:   floatFromDecimal(dbPayment.ValorRecebido),
		TaxaPercentual:  floatFromDecimal(dbPayment.TaxaPercentual),
		TaxaFixa:        floatFromDecimal(dbPayment.TaxaFixa),
		ValorLiquido:    floatFromDecimal(dbPayment.ValorLiquido),
		CriadoPor:       ptrUUIDFromUUID(dbPayment.CriadoPor),
		CriadoEm:        dbPayment.CriadoEm.Time,
	}

	if dbPayment.Observacoes != nil {
		payment.Observacoes = dbPayment.Observacoes
	}

	return payment
}

// ============================================================================
// Type Conversion Helpers
// ============================================================================

func uuidToUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: u, Valid: true}
}

func uuidFromUUID(u pgtype.UUID) uuid.UUID {
	return u.Bytes
}

func uuidToUUIDPtr(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

func ptrUUIDToUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

func ptrUUIDFromUUID(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	uid := uuid.UUID(u.Bytes)
	return &uid
}

func decimalFromFloat(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

func floatFromDecimal(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

func ptrString(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}
