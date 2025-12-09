package servico

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateServicoUseCase cria um novo serviço
type CreateServicoUseCase struct {
	repo    port.ServicoRepository
	catRepo port.CategoriaServicoRepository
	logger  *zap.Logger
}

// NewCreateServicoUseCase cria uma nova instância do use case
func NewCreateServicoUseCase(
	repo port.ServicoRepository,
	catRepo port.CategoriaServicoRepository,
	logger *zap.Logger,
) *CreateServicoUseCase {
	return &CreateServicoUseCase{
		repo:    repo,
		catRepo: catRepo,
		logger:  logger,
	}
}

// Execute cria um novo serviço
func (uc *CreateServicoUseCase) Execute(
	ctx context.Context,
	tenantID string,
	req dto.CreateServicoRequest,
) (*dto.ServicoResponse, error) {
	// Validar tenant_id
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidTenantID
	}

	// Verificar duplicidade de nome
	exists, err := uc.repo.CheckNomeExists(ctx, tenantID, req.Nome, "")
	if err != nil {
		uc.logger.Error("Erro ao verificar nome duplicado",
			zap.String("tenant_id", tenantID),
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
		categoria, err := uc.catRepo.FindByID(ctx, tenantID, *req.CategoriaID)
		if err != nil || categoria == nil {
			return nil, domain.ErrCategoriaNotFound
		}
	}

	// Converter DTO para entidade
	servico, err := mapper.CreateServicoRequestToEntity(&req, tenantUUID)
	if err != nil {
		return nil, err
	}

	// Persistir
	if err := uc.repo.Create(ctx, servico); err != nil {
		uc.logger.Error("Erro ao criar serviço",
			zap.String("tenant_id", tenantID),
			zap.String("nome", req.Nome),
			zap.Error(err),
		)
		return nil, err
	}

	uc.logger.Info("Serviço criado com sucesso",
		zap.String("servico_id", servico.ID.String()),
		zap.String("tenant_id", tenantID),
		zap.String("nome", servico.Nome),
	)

	return mapper.ServicoToResponse(servico), nil
}
