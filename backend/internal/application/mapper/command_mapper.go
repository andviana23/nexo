package mapper

import (
	"fmt"
	"strconv"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
)

// CommandMapper lida com conversões entre Entity e DTO
type CommandMapper struct{}

// NewCommandMapper cria uma nova instância do mapper
func NewCommandMapper() *CommandMapper {
	return &CommandMapper{}
}

// ============================================================================
// Entity -> DTO (Response)
// ============================================================================

// ToCommandResponse converte Command entity para CommandResponse DTO
func (m *CommandMapper) ToCommandResponse(command *entity.Command) *dto.CommandResponse {
	response := &dto.CommandResponse{
		ID:                 command.ID.String(),
		CustomerID:         command.CustomerID.String(),
		Status:             string(command.Status),
		Subtotal:           formatMoney(command.Subtotal),
		Desconto:           formatMoney(command.Desconto),
		Total:              formatMoney(command.Total),
		TotalRecebido:      formatMoney(command.TotalRecebido),
		Troco:              formatMoney(command.Troco),
		SaldoDevedor:       formatMoney(command.SaldoDevedor),
		DeixarTrocoGorjeta: command.DeixarTrocoGorjeta,
		DeixarSaldoDivida:  command.DeixarSaldoDivida,
		CriadoEm:           command.CriadoEm,
		AtualizadoEm:       command.AtualizadoEm,
		Items:              []dto.CommandItemResponse{},
		Payments:           []dto.CommandPaymentResponse{},
	}

	if command.AppointmentID != nil {
		aid := command.AppointmentID.String()
		response.AppointmentID = &aid
	}

	if command.Numero != nil {
		response.Numero = command.Numero
	}

	if command.Observacoes != nil {
		response.Observacoes = command.Observacoes
	}

	if command.FechadoEm != nil {
		response.FechadoEm = command.FechadoEm
	}

	if command.FechadoPor != nil {
		fid := command.FechadoPor.String()
		response.FechadoPor = &fid
	}

	// Converter itens
	for _, item := range command.Items {
		response.Items = append(response.Items, m.ToCommandItemResponse(item))
	}

	// Converter pagamentos
	for _, payment := range command.Payments {
		response.Payments = append(response.Payments, m.ToCommandPaymentResponse(payment))
	}

	return response
}

// ToCommandItemResponse converte CommandItem entity para CommandItemResponse DTO
func (m *CommandMapper) ToCommandItemResponse(item entity.CommandItem) dto.CommandItemResponse {
	response := dto.CommandItemResponse{
		ID:                 item.ID.String(),
		CommandID:          item.CommandID.String(),
		Tipo:               string(item.Tipo),
		ItemID:             item.ItemID.String(),
		Descricao:          item.Descricao,
		PrecoUnitario:      formatMoney(item.PrecoUnitario),
		Quantidade:         item.Quantidade,
		DescontoValor:      formatMoney(item.DescontoValor),
		DescontoPercentual: formatPercentage(item.DescontoPercentual),
		PrecoFinal:         formatMoney(item.PrecoFinal),
		CriadoEm:           item.CriadoEm,
	}

	if item.Observacoes != nil {
		response.Observacoes = item.Observacoes
	}

	return response
}

// ToCommandPaymentResponse converte CommandPayment entity para CommandPaymentResponse DTO
func (m *CommandMapper) ToCommandPaymentResponse(payment entity.CommandPayment) dto.CommandPaymentResponse {
	response := dto.CommandPaymentResponse{
		ID:              payment.ID.String(),
		CommandID:       payment.CommandID.String(),
		MeioPagamentoID: payment.MeioPagamentoID.String(),
		ValorRecebido:   formatMoney(payment.ValorRecebido),
		TaxaPercentual:  formatPercentage(payment.TaxaPercentual),
		TaxaFixa:        formatMoney(payment.TaxaFixa),
		ValorLiquido:    formatMoney(payment.ValorLiquido),
		CriadoEm:        payment.CriadoEm,
	}

	if payment.Observacoes != nil {
		response.Observacoes = payment.Observacoes
	}

	cpID := payment.CriadoPor.String()
	response.CriadoPor = &cpID

	return response
}

// ToCommandListResponse converte lista de Commands para CommandListResponse
func (m *CommandMapper) ToCommandListResponse(commands []*entity.Command, total, page, pageSize int) *dto.CommandListResponse {
	responses := make([]dto.CommandResponse, 0, len(commands))
	for _, cmd := range commands {
		responses = append(responses, *m.ToCommandResponse(cmd))
	}

	totalPages := (total + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	return &dto.CommandListResponse{
		Commands:   responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ============================================================================
// DTO (Request) -> Entity
// ============================================================================

// FromCreateCommandRequest converte CreateCommandRequest para Command entity
func (m *CommandMapper) FromCreateCommandRequest(req *dto.CreateCommandRequest, tenantID uuid.UUID) (*entity.Command, error) {
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("invalid customer_id: %w", err)
	}

	var appointmentID *uuid.UUID
	if req.AppointmentID != nil {
		aid, err := uuid.Parse(*req.AppointmentID)
		if err != nil {
			return nil, fmt.Errorf("invalid appointment_id: %w", err)
		}
		appointmentID = &aid
	}

	// Criar comanda
	command := &entity.Command{
		ID:            uuid.New(),
		TenantID:      tenantID,
		AppointmentID: appointmentID,
		CustomerID:    customerID,
		Observacoes:   req.Observacoes,
		Items:         []entity.CommandItem{},
		Payments:      []entity.CommandPayment{},
	}

	// Adicionar itens
	for _, itemInput := range req.Items {
		item, err := m.FromCommandItemInput(itemInput, command.ID, tenantID)
		if err != nil {
			return nil, err
		}
		command.Items = append(command.Items, *item)
	}

	// Criar comanda com regras de negócio (NewCommand inicializa com status OPEN)
	domainCommand, err := entity.NewCommand(
		tenantID,
		customerID,
		appointmentID,
	)
	if err != nil {
		return nil, err
	}

	// Configurar observações se houver
	if req.Observacoes != nil {
		domainCommand.Observacoes = req.Observacoes
	}

	// Substituir items gerados pelo NewCommand pelos items do request
	domainCommand.Items = command.Items

	// Recalcular totais
	domainCommand.RecalculateTotals()

	return domainCommand, nil
}

// FromCommandItemInput converte CommandItemInput para CommandItem entity
func (m *CommandMapper) FromCommandItemInput(input dto.CommandItemInput, commandID, tenantID uuid.UUID) (*entity.CommandItem, error) {
	itemID, err := uuid.Parse(input.ItemID)
	if err != nil {
		return nil, fmt.Errorf("invalid item_id: %w", err)
	}

	precoUnitario, err := parseMoney(input.PrecoUnitario)
	if err != nil {
		return nil, fmt.Errorf("invalid preco_unitario: %w", err)
	}

	descontoValor := 0.0
	if input.DescontoValor != nil {
		descontoValor, err = parseMoney(*input.DescontoValor)
		if err != nil {
			return nil, fmt.Errorf("invalid desconto_valor: %w", err)
		}
	}

	descontoPercentual := 0.0
	if input.DescontoPercentual != nil {
		descontoPercentual = *input.DescontoPercentual
	}

	item, err := entity.NewCommandItem(
		commandID,
		entity.CommandItemType(input.Tipo),
		itemID,
		input.Descricao,
		precoUnitario,
		input.Quantidade,
	)
	if err != nil {
		return nil, err
	}

	// Aplicar descontos se houver
	if descontoValor > 0 || descontoPercentual > 0 {
		if err := item.ApplyDiscount(descontoValor, descontoPercentual); err != nil {
			return nil, err
		}
	}

	if input.Observacoes != nil {
		item.Observacoes = input.Observacoes
	}

	return item, nil
}

// FromAddCommandItemRequest converte AddCommandItemRequest para CommandItem entity
func (m *CommandMapper) FromAddCommandItemRequest(req *dto.AddCommandItemRequest, commandID, tenantID uuid.UUID) (*entity.CommandItem, error) {
	itemID, err := uuid.Parse(req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("invalid item_id: %w", err)
	}

	precoUnitario, err := parseMoney(req.PrecoUnitario)
	if err != nil {
		return nil, fmt.Errorf("invalid preco_unitario: %w", err)
	}

	item, err := entity.NewCommandItem(
		commandID,
		entity.CommandItemType(req.Tipo),
		itemID,
		req.Descricao,
		precoUnitario,
		req.Quantidade,
	)
	if err != nil {
		return nil, err
	}

	// Aplicar descontos
	descontoValor := 0.0
	if req.DescontoValor != nil {
		descontoValor, err = parseMoney(*req.DescontoValor)
		if err != nil {
			return nil, fmt.Errorf("invalid desconto_valor: %w", err)
		}
	}

	descontoPercentual := 0.0
	if req.DescontoPercentual != nil {
		descontoPercentual = *req.DescontoPercentual
	}

	if descontoValor > 0 || descontoPercentual > 0 {
		if err := item.ApplyDiscount(descontoValor, descontoPercentual); err != nil {
			return nil, err
		}
	}

	if req.Observacoes != nil {
		item.Observacoes = req.Observacoes
	}

	return item, nil
}

// FromAddCommandPaymentRequest converte AddCommandPaymentRequest para CommandPayment entity
func (m *CommandMapper) FromAddCommandPaymentRequest(req *dto.AddCommandPaymentRequest, commandID, tenantID, userID uuid.UUID, taxaPercentual, taxaFixa float64) (*entity.CommandPayment, error) {
	meioPagamentoID, err := uuid.Parse(req.MeioPagamentoID)
	if err != nil {
		return nil, fmt.Errorf("invalid meio_pagamento_id: %w", err)
	}

	valorRecebido, err := parseMoney(req.ValorRecebido)
	if err != nil {
		return nil, fmt.Errorf("invalid valor_recebido: %w", err)
	}

	payment, err := entity.NewCommandPayment(
		commandID,
		meioPagamentoID,
		valorRecebido,
		taxaPercentual,
		taxaFixa,
		&userID,
	)
	if err != nil {
		return nil, err
	}

	if req.Observacoes != nil {
		payment.Observacoes = req.Observacoes
	}

	return payment, nil
}

// ============================================================================
// Helper Functions
// ============================================================================

// formatMoney formata float64 para string com 2 casas decimais
func formatMoney(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

// formatPercentage formata float64 para string com 2 casas decimais (percentual)
func formatPercentage(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

// parseMoney converte string para float64
func parseMoney(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ptrBoolValue retorna o valor do ponteiro ou false se nil
func ptrBoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
