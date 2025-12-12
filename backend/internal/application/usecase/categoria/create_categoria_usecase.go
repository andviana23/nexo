package categoria

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateCategoriaServicoUseCase cria uma nova categoria de serviço
type CreateCategoriaServicoUseCase struct {
	repo   port.CategoriaServicoRepository
	logger *zap.Logger
}

// NewCreateCategoriaServicoUseCase cria uma nova instância do use case
func NewCreateCategoriaServicoUseCase(repo port.CategoriaServicoRepository, logger *zap.Logger) *CreateCategoriaServicoUseCase {
	return &CreateCategoriaServicoUseCase{
		repo:   repo,
		logger: logger,
	}
}

// Execute cria uma nova categoria de serviço
func (uc *CreateCategoriaServicoUseCase) Execute(
	ctx context.Context,
	tenantID, unitID string,
	req dto.CreateCategoriaServicoRequest,
) (*dto.CategoriaServicoResponse, error) {
	// Validar tenant_id
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, domain.ErrInvalidTenantID
	}

	// Verificar duplicidade de nome
	exists, err := uc.repo.CheckNomeExists(ctx, tenantID, unitID, req.Nome, "")
	if err != nil {
		uc.logger.Error("Erro ao verificar nome duplicado",
			zap.String("tenant_id", tenantID),
			zap.String("nome", req.Nome),
			zap.Error(err),
		)
		return nil, err
	}

	if exists {
		return nil, domain.ErrCategoriaNomeDuplicate
	}

	// Validar unit_id
	unitUUID, err := uuid.Parse(unitID)
	if err != nil {
		return nil, domain.ErrInvalidUnitID
	}

	// Criar entidade
	categoria, err := entity.NewCategoriaServico(tenantUUID, unitUUID, req.Nome)
	if err != nil {
		return nil, err
	}

	// Definir campos opcionais
	if req.Descricao != nil {
		categoria.SetDescricao(*req.Descricao)
	}
	if req.Cor != nil {
		if err := categoria.SetCor(*req.Cor); err != nil {
			return nil, err
		}
	}
	if req.Icone != nil {
		categoria.SetIcone(*req.Icone)
	}

	// Persistir
	if err := uc.repo.Create(ctx, categoria); err != nil {
		uc.logger.Error("Erro ao criar categoria de serviço",
			zap.String("tenant_id", tenantID),
			zap.String("nome", req.Nome),
			zap.Error(err),
		)
		return nil, err
	}

	uc.logger.Info("Categoria de serviço criada com sucesso",
		zap.String("categoria_id", categoria.ID.String()),
		zap.String("tenant_id", tenantID),
		zap.String("nome", categoria.Nome),
	)

	return mapper.CategoriaServicoToResponse(categoria), nil
}
