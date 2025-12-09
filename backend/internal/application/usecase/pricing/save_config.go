package pricing

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// SaveConfigPrecificacaoInput define os dados de entrada
type SaveConfigPrecificacaoInput struct {
	TenantID          string
	MargemDesejada    valueobject.Percentage
	MarkupAlvo        decimal.Decimal
	ImpostoPercentual valueobject.Percentage
	ComissaoDefault   valueobject.Percentage
}

// SaveConfigPrecificacaoUseCase salva configuração de precificação
type SaveConfigPrecificacaoUseCase struct {
	repo   port.PrecificacaoConfigRepository
	logger *zap.Logger
}

// NewSaveConfigPrecificacaoUseCase cria nova instância
func NewSaveConfigPrecificacaoUseCase(
	repo port.PrecificacaoConfigRepository,
	logger *zap.Logger,
) *SaveConfigPrecificacaoUseCase {
	return &SaveConfigPrecificacaoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute salva ou atualiza a configuração de precificação
func (uc *SaveConfigPrecificacaoUseCase) Execute(ctx context.Context, input SaveConfigPrecificacaoInput) (*entity.PrecificacaoConfig, error) {
	if input.TenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}

	// Buscar configuração existente
	existing, err := uc.repo.FindByTenantID(ctx, input.TenantID)
	if err == nil && existing != nil {
		// Atualizar configuração existente
		if errUpdate := existing.AtualizarMargem(input.MargemDesejada); errUpdate != nil {
			return nil, fmt.Errorf("erro ao atualizar margem: %w", errUpdate)
		}
		if errUpdate := existing.AtualizarMarkup(input.MarkupAlvo); errUpdate != nil {
			return nil, fmt.Errorf("erro ao atualizar markup: %w", errUpdate)
		}
		existing.AtualizarImposto(input.ImpostoPercentual)
		existing.AtualizarComissaoDefault(input.ComissaoDefault)

		if errUpdate := uc.repo.Update(ctx, existing); errUpdate != nil {
			return nil, fmt.Errorf("erro ao salvar configuração: %w", errUpdate)
		}
		return existing, nil
	}

	// Converter tenant_id de string para uuid.UUID
	tenantUUID, err := uuid.Parse(input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}

	// Criar nova configuração
	config, err := entity.NewPrecificacaoConfig(
		tenantUUID,
		input.MargemDesejada,
		input.ImpostoPercentual,
		input.ComissaoDefault,
		input.MarkupAlvo,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar configuração: %w", err)
	}

	if err := uc.repo.Create(ctx, config); err != nil {
		return nil, fmt.Errorf("erro ao salvar configuração: %w", err)
	}

	uc.logger.Info("Configuração de precificação salva",
		zap.String("tenant_id", input.TenantID),
	)

	return config, nil
}
