package servico

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// GetServicoUseCase busca um serviço específico
type GetServicoUseCase struct {
	repo   port.ServicoRepository
	logger *zap.Logger
}

// NewGetServicoUseCase cria uma nova instância do use case
func NewGetServicoUseCase(repo port.ServicoRepository, logger *zap.Logger) *GetServicoUseCase {
	return &GetServicoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute busca um serviço por ID
func (uc *GetServicoUseCase) Execute(
	ctx context.Context,
	tenantID string,
	unitID string,
	servicoID string,
) (*dto.ServicoResponse, error) {
	servico, err := uc.repo.FindByID(ctx, tenantID, unitID, servicoID)
	if err != nil {
		uc.logger.Error("Erro ao buscar serviço",
			zap.String("tenant_id", tenantID),
			zap.String("unit_id", unitID),
			zap.String("servico_id", servicoID),
			zap.Error(err),
		)
		return nil, domain.ErrServicoNotFound
	}

	return mapper.ServicoToResponse(servico), nil
}

// GetServicoStatsUseCase busca estatísticas dos serviços
type GetServicoStatsUseCase struct {
	repo   port.ServicoRepository
	logger *zap.Logger
}

// NewGetServicoStatsUseCase cria uma nova instância do use case
func NewGetServicoStatsUseCase(repo port.ServicoRepository, logger *zap.Logger) *GetServicoStatsUseCase {
	return &GetServicoStatsUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute busca estatísticas dos serviços
func (uc *GetServicoStatsUseCase) Execute(
	ctx context.Context,
	tenantID string,
	unitID string,
) (*dto.ServicoStatsResponse, error) {
	stats, err := uc.repo.GetStats(ctx, tenantID, unitID)
	if err != nil {
		uc.logger.Error("Erro ao buscar estatísticas de serviços",
			zap.String("tenant_id", tenantID),
			zap.String("unit_id", unitID),
			zap.Error(err),
		)
		return nil, err
	}

	return mapper.ServicoStatsToResponse(stats), nil
}
