package financial

import (
	"context"

	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// MarcarPagamentoInput define os dados de entrada para marcar pagamento
type MarcarPagamentoInput struct {
	TenantID       string
	ContaID        string
	DataPagamento  time.Time
	ComprovanteURL string
}

// MarcarPagamentoUseCase implementa a marcação de pagamento de conta a pagar
type MarcarPagamentoUseCase struct {
	repo   port.ContaPagarRepository
	logger *zap.Logger
}

// NewMarcarPagamentoUseCase cria nova instância do use case
func NewMarcarPagamentoUseCase(
	repo port.ContaPagarRepository,
	logger *zap.Logger,
) *MarcarPagamentoUseCase {
	return &MarcarPagamentoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute marca uma conta a pagar como paga
func (uc *MarcarPagamentoUseCase) Execute(ctx context.Context, input MarcarPagamentoInput) (*entity.ContaPagar, error) {
	// Validações de entrada
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.ContaID == "" {
		return nil, fmt.Errorf("ID da conta é obrigatório")
	}

	if input.DataPagamento.IsZero() {
		return nil, fmt.Errorf("data de pagamento é obrigatória")
	}

	// Buscar conta existente
	conta, err := uc.repo.FindByID(ctx, input.TenantID, input.ContaID)
	if err != nil {
		uc.logger.Error("Erro ao buscar conta a pagar",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("conta_id", input.ContaID),
		)
		return nil, fmt.Errorf("erro ao buscar conta: %w", err)
	}

	// Marcar como pago (método do domínio)
	if err := conta.MarcarComoPago(input.DataPagamento, input.ComprovanteURL); err != nil {
		uc.logger.Error("Erro ao marcar conta como paga",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("conta_id", input.ContaID),
		)
		return nil, fmt.Errorf("erro ao marcar conta como paga: %w", err)
	}

	// Atualizar no repositório
	if err := uc.repo.Update(ctx, conta); err != nil {
		uc.logger.Error("Erro ao atualizar conta a pagar",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("conta_id", input.ContaID),
		)
		return nil, fmt.Errorf("erro ao atualizar conta: %w", err)
	}

	uc.logger.Info("Conta a pagar marcada como paga",
		zap.String("tenant_id", input.TenantID),
		zap.String("conta_id", conta.ID),
		zap.Time("data_pagamento", input.DataPagamento),
	)

	return conta, nil
}
