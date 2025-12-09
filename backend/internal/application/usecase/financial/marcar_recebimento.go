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

// MarcarRecebimentoInput define os dados de entrada para marcar recebimento
type MarcarRecebimentoInput struct {
	TenantID        string
	ContaID         string
	DataRecebimento time.Time
	ComprovanteURL  string
}

// MarcarRecebimentoUseCase implementa a marcação de recebimento de conta a receber
type MarcarRecebimentoUseCase struct {
	repo   port.ContaReceberRepository
	logger *zap.Logger
}

// NewMarcarRecebimentoUseCase cria nova instância do use case
func NewMarcarRecebimentoUseCase(
	repo port.ContaReceberRepository,
	logger *zap.Logger,
) *MarcarRecebimentoUseCase {
	return &MarcarRecebimentoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute marca uma conta a receber como recebida
func (uc *MarcarRecebimentoUseCase) Execute(ctx context.Context, input MarcarRecebimentoInput) (*entity.ContaReceber, error) {
	// Validações de entrada
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	if input.ContaID == "" {
		return nil, fmt.Errorf("ID da conta é obrigatório")
	}

	if input.DataRecebimento.IsZero() {
		return nil, fmt.Errorf("data de recebimento é obrigatória")
	}

	// Buscar conta existente
	conta, err := uc.repo.FindByID(ctx, input.TenantID, input.ContaID)
	if err != nil {
		uc.logger.Error("Erro ao buscar conta a receber",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("conta_id", input.ContaID),
		)
		return nil, fmt.Errorf("erro ao buscar conta: %w", err)
	}

	// Marcar como recebido (método do domínio)
	if err := conta.MarcarComoRecebido(input.DataRecebimento); err != nil {
		uc.logger.Error("Erro ao marcar conta como recebida",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("conta_id", input.ContaID),
		)
		return nil, fmt.Errorf("erro ao marcar conta como recebida: %w", err)
	}

	// Atualizar no repositório
	if err := uc.repo.Update(ctx, conta); err != nil {
		uc.logger.Error("Erro ao atualizar conta a receber",
			zap.Error(err),
			zap.String("tenant_id", input.TenantID),
			zap.String("conta_id", input.ContaID),
		)
		return nil, fmt.Errorf("erro ao atualizar conta: %w", err)
	}

	uc.logger.Info("Conta a receber marcada como recebida",
		zap.String("tenant_id", input.TenantID),
		zap.String("conta_id", conta.ID),
		zap.Time("data_recebimento", input.DataRecebimento),
	)

	return conta, nil
}
