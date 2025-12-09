package command

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// AddCommandPaymentUseCase implementa a adição de pagamento à comanda
type AddCommandPaymentUseCase struct {
	repo              port.CommandRepository
	meioPagamentoRepo port.MeioPagamentoRepository
	mapper            *mapper.CommandMapper
}

// NewAddCommandPaymentUseCase cria uma nova instância do use case
func NewAddCommandPaymentUseCase(
	repo port.CommandRepository,
	meioPagamentoRepo port.MeioPagamentoRepository,
	mapper *mapper.CommandMapper,
) *AddCommandPaymentUseCase {
	return &AddCommandPaymentUseCase{
		repo:              repo,
		meioPagamentoRepo: meioPagamentoRepo,
		mapper:            mapper,
	}
}

// Execute adiciona um pagamento à comanda
// As taxas são buscadas automaticamente do MeioPagamento configurado
func (uc *AddCommandPaymentUseCase) Execute(ctx context.Context, commandID, tenantID, userID uuid.UUID, req *dto.AddCommandPaymentRequest) (*dto.CommandResponse, error) {
	// Buscar comanda existente
	command, err := uc.repo.FindByID(ctx, commandID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get command: %w", err)
	}

	if command == nil {
		return nil, fmt.Errorf("command not found")
	}

	// Buscar meio de pagamento para obter taxas configuradas
	meioPagamentoID, err := uuid.Parse(req.MeioPagamentoID)
	if err != nil {
		return nil, fmt.Errorf("invalid meio_pagamento_id: %w", err)
	}

	meioPagamento, err := uc.meioPagamentoRepo.FindByID(ctx, tenantID.String(), meioPagamentoID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get meio de pagamento: %w", err)
	}
	if meioPagamento == nil {
		return nil, fmt.Errorf("meio de pagamento not found")
	}

	// Verificar se o meio de pagamento está ativo
	if !meioPagamento.Ativo {
		return nil, fmt.Errorf("meio de pagamento '%s' está inativo", meioPagamento.Nome)
	}

	// Extrair taxas do meio de pagamento (decimal → float64)
	taxaPercentual, _ := meioPagamento.Taxa.Float64()
	taxaFixa, _ := meioPagamento.TaxaFixa.Float64()

	// Converter request para entity (com taxas do meio de pagamento)
	payment, err := uc.mapper.FromAddCommandPaymentRequest(req, commandID, tenantID, userID, taxaPercentual, taxaFixa)
	if err != nil {
		return nil, fmt.Errorf("failed to map payment: %w", err)
	}

	// Adicionar pagamento via domain logic
	if err := command.AddPayment(*payment); err != nil {
		return nil, fmt.Errorf("failed to add payment to command: %w", err)
	}

	// Persistir pagamento
	if err := uc.repo.AddPayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to persist payment: %w", err)
	}

	// Atualizar comanda (totais recalculados)
	if err := uc.repo.Update(ctx, command); err != nil {
		return nil, fmt.Errorf("failed to update command: %w", err)
	}

	// Buscar comanda atualizada
	updated, err := uc.repo.FindByID(ctx, commandID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated command: %w", err)
	}

	// Converter para DTO
	response := uc.mapper.ToCommandResponse(updated)
	return response, nil
}
