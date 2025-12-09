package pricing

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SaveSimulacaoInput define os dados de entrada
type SaveSimulacaoInput struct {
	Simulacao *entity.PrecificacaoSimulacao
}

// SaveSimulacaoUseCase persiste uma simulação no histórico
type SaveSimulacaoUseCase struct {
	repo   port.PrecificacaoSimulacaoRepository
	logger *zap.Logger
}

// NewSaveSimulacaoUseCase cria nova instância
func NewSaveSimulacaoUseCase(
	repo port.PrecificacaoSimulacaoRepository,
	logger *zap.Logger,
) *SaveSimulacaoUseCase {
	return &SaveSimulacaoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute salva a simulação no histórico
func (uc *SaveSimulacaoUseCase) Execute(ctx context.Context, input SaveSimulacaoInput) error {
	if input.Simulacao == nil {
		return fmt.Errorf("simulação é obrigatória")
	}

	if input.Simulacao.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}

	if err := uc.repo.Create(ctx, input.Simulacao); err != nil {
		return fmt.Errorf("erro ao salvar simulação: %w", err)
	}

	uc.logger.Info("Simulação salva no histórico",
		zap.String("tenant_id", input.Simulacao.TenantID.String()),
		zap.String("simulacao_id", input.Simulacao.ID),
	)

	return nil
}
