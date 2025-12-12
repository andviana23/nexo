package servico

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// DeleteServicoUseCase deleta um serviço
type DeleteServicoUseCase struct {
	repo   port.ServicoRepository
	logger *zap.Logger
}

// NewDeleteServicoUseCase cria uma nova instância do use case
func NewDeleteServicoUseCase(repo port.ServicoRepository, logger *zap.Logger) *DeleteServicoUseCase {
	return &DeleteServicoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute deleta um serviço
func (uc *DeleteServicoUseCase) Execute(
	ctx context.Context,
	tenantID string,
	unitID string,
	servicoID string,
) error {
	// Verificar se serviço existe
	servico, err := uc.repo.FindByID(ctx, tenantID, unitID, servicoID)
	if err != nil {
		return domain.ErrServicoNotFound
	}

	// TODO: Verificar se o serviço está vinculado a agendamentos
	// Se houver agendamentos futuros, retornar erro
	// Por ora, vamos apenas deletar

	// Deletar
	if err := uc.repo.Delete(ctx, tenantID, servico.UnitID.String(), servicoID); err != nil {
		uc.logger.Error("Erro ao deletar serviço",
			zap.String("tenant_id", tenantID),
			zap.String("servico_id", servicoID),
			zap.Error(err),
		)
		return err
	}

	uc.logger.Info("Serviço deletado com sucesso",
		zap.String("servico_id", servicoID),
		zap.String("tenant_id", tenantID),
	)

	return nil
}
