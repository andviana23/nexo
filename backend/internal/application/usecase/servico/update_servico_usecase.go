package servico

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// UpdateServicoUseCase atualiza um serviço existente
type UpdateServicoUseCase struct {
	repo    port.ServicoRepository
	catRepo port.CategoriaServicoRepository
	logger  *zap.Logger
}

// NewUpdateServicoUseCase cria uma nova instância do use case
func NewUpdateServicoUseCase(
	repo port.ServicoRepository,
	catRepo port.CategoriaServicoRepository,
	logger *zap.Logger,
) *UpdateServicoUseCase {
	return &UpdateServicoUseCase{
		repo:    repo,
		catRepo: catRepo,
		logger:  logger,
	}
}

// Execute atualiza um serviço existente
func (uc *UpdateServicoUseCase) Execute(
	ctx context.Context,
	tenantID string,
	unitID string,
	servicoID string,
	req dto.UpdateServicoRequest,
) (*dto.ServicoResponse, error) {
	// Buscar serviço existente
	servico, err := uc.repo.FindByID(ctx, tenantID, unitID, servicoID)
	if err != nil {
		return nil, domain.ErrServicoNotFound
	}

	// Verificar duplicidade de nome (excluindo o próprio serviço)
	exists, err := uc.repo.CheckNomeExists(ctx, tenantID, unitID, req.Nome, servicoID)
	if err != nil {
		uc.logger.Error("Erro ao verificar nome duplicado",
			zap.String("tenant_id", tenantID),
			zap.String("unit_id", unitID),
			zap.String("servico_id", servicoID),
			zap.String("nome", req.Nome),
			zap.Error(err),
		)
		return nil, err
	}

	if exists {
		return nil, domain.ErrServicoNomeDuplicate
	}

	// Se foi informada uma categoria, verificar se existe
	if req.CategoriaID != nil {
		categoria, err := uc.catRepo.FindByID(ctx, tenantID, unitID, *req.CategoriaID)
		if err != nil || categoria == nil {
			return nil, domain.ErrCategoriaNotFound
		}
	}

	// Aplicar alterações na entidade
	if err := mapper.UpdateServicoRequestToEntity(&req, servico); err != nil {
		return nil, err
	}

	// Persistir
	if err := uc.repo.Update(ctx, servico); err != nil {
		uc.logger.Error("Erro ao atualizar serviço",
			zap.String("tenant_id", tenantID),
			zap.String("servico_id", servicoID),
			zap.Error(err),
		)
		return nil, err
	}

	uc.logger.Info("Serviço atualizado com sucesso",
		zap.String("servico_id", servico.ID.String()),
		zap.String("tenant_id", tenantID),
		zap.String("nome", servico.Nome),
	)

	return mapper.ServicoToResponse(servico), nil
}

// ToggleServicoStatusUseCase ativa/desativa um serviço
type ToggleServicoStatusUseCase struct {
	repo   port.ServicoRepository
	logger *zap.Logger
}

// NewToggleServicoStatusUseCase cria uma nova instância do use case
func NewToggleServicoStatusUseCase(repo port.ServicoRepository, logger *zap.Logger) *ToggleServicoStatusUseCase {
	return &ToggleServicoStatusUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute ativa/desativa um serviço (alternando seu estado atual)
func (uc *ToggleServicoStatusUseCase) Execute(
	ctx context.Context,
	tenantID string,
	unitID string,
	servicoID string,
) (*dto.ServicoResponse, error) {
	// Verificar se serviço existe
	servico, err := uc.repo.FindByID(ctx, tenantID, unitID, servicoID)
	if err != nil {
		return nil, domain.ErrServicoNotFound
	}

	// Alternar status (inverter o valor atual)
	novoStatus := !servico.Ativo

	// Alterar status no repositório
	if err := uc.repo.ToggleStatus(ctx, tenantID, servico.UnitID.String(), servicoID, novoStatus); err != nil {
		uc.logger.Error("Erro ao alterar status do serviço",
			zap.String("tenant_id", tenantID),
			zap.String("servico_id", servicoID),
			zap.Bool("novo_status", novoStatus),
			zap.Error(err),
		)
		return nil, err
	}

	// Atualizar entidade local
	if novoStatus {
		servico.Ativar()
	} else {
		servico.Desativar()
	}

	uc.logger.Info("Status do serviço alterado",
		zap.String("servico_id", servicoID),
		zap.String("tenant_id", tenantID),
		zap.Bool("novo_status", novoStatus),
	)

	return mapper.ServicoToResponse(servico), nil
}
