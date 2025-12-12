package categoriaproduto

import (
	"context"
	"errors"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// =============================================================================
// CreateCategoriaProdutoUseCase
// =============================================================================

type CreateCategoriaProdutoInput struct {
	TenantID    uuid.UUID
	UnitID      uuid.UUID
	Nome        string
	Descricao   string
	Cor         string
	Icone       string
	CentroCusto string
}

type CreateCategoriaProdutoUseCase struct {
	repo   port.CategoriaProdutoRepository
	logger *zap.Logger
}

func NewCreateCategoriaProdutoUseCase(repo port.CategoriaProdutoRepository, logger *zap.Logger) *CreateCategoriaProdutoUseCase {
	return &CreateCategoriaProdutoUseCase{repo: repo, logger: logger}
}

func (uc *CreateCategoriaProdutoUseCase) Execute(ctx context.Context, input CreateCategoriaProdutoInput) (*entity.CategoriaProdutoEntity, error) {
	// 1. Validar se nome já existe
	exists, err := uc.repo.ExistsWithNome(ctx, input.TenantID, input.UnitID, input.Nome, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("já existe uma categoria com este nome")
	}

	// 2. Criar entidade
	categoria, err := entity.NewCategoriaProduto(input.TenantID, input.UnitID, input.Nome)
	if err != nil {
		return nil, err
	}

	// 3. Preencher campos opcionais
	if input.Descricao != "" {
		categoria.SetDescricao(input.Descricao)
	}
	if input.Cor != "" {
		if err := categoria.SetCor(input.Cor); err != nil {
			return nil, err
		}
	}
	if input.Icone != "" {
		categoria.SetIcone(input.Icone)
	}
	if input.CentroCusto != "" {
		if err := categoria.SetCentroCusto(entity.CentroCusto(input.CentroCusto)); err != nil {
			return nil, err
		}
	}

	// 4. Persistir
	if err := uc.repo.Create(ctx, categoria); err != nil {
		return nil, err
	}

	return categoria, nil
}

// =============================================================================
// ListCategoriasProdutosUseCase
// =============================================================================

type ListCategoriasProdutosUseCase struct {
	repo   port.CategoriaProdutoRepository
	logger *zap.Logger
}

func NewListCategoriasProdutosUseCase(repo port.CategoriaProdutoRepository, logger *zap.Logger) *ListCategoriasProdutosUseCase {
	return &ListCategoriasProdutosUseCase{repo: repo, logger: logger}
}

func (uc *ListCategoriasProdutosUseCase) Execute(ctx context.Context, tenantID, unitID uuid.UUID, apenasAtivas bool) ([]*entity.CategoriaProdutoEntity, error) {
	if apenasAtivas {
		return uc.repo.ListAtivas(ctx, tenantID, unitID)
	}
	return uc.repo.ListAll(ctx, tenantID, unitID)
}

// =============================================================================
// GetCategoriaProdutoUseCase
// =============================================================================

type GetCategoriaProdutoUseCase struct {
	repo   port.CategoriaProdutoRepository
	logger *zap.Logger
}

func NewGetCategoriaProdutoUseCase(repo port.CategoriaProdutoRepository, logger *zap.Logger) *GetCategoriaProdutoUseCase {
	return &GetCategoriaProdutoUseCase{repo: repo, logger: logger}
}

func (uc *GetCategoriaProdutoUseCase) Execute(ctx context.Context, tenantID, unitID, id uuid.UUID) (*entity.CategoriaProdutoEntity, error) {
	categoria, err := uc.repo.FindByID(ctx, tenantID, unitID, id)
	if err != nil {
		return nil, err
	}
	if categoria == nil {
		return nil, errors.New("categoria não encontrada")
	}
	return categoria, nil
}

// =============================================================================
// UpdateCategoriaProdutoUseCase
// =============================================================================

type UpdateCategoriaProdutoInput struct {
	TenantID    uuid.UUID
	UnitID      uuid.UUID
	ID          uuid.UUID
	Nome        string
	Descricao   string
	Cor         string
	Icone       string
	CentroCusto string
	Ativa       bool
}

type UpdateCategoriaProdutoUseCase struct {
	repo   port.CategoriaProdutoRepository
	logger *zap.Logger
}

func NewUpdateCategoriaProdutoUseCase(repo port.CategoriaProdutoRepository, logger *zap.Logger) *UpdateCategoriaProdutoUseCase {
	return &UpdateCategoriaProdutoUseCase{repo: repo, logger: logger}
}

func (uc *UpdateCategoriaProdutoUseCase) Execute(ctx context.Context, input UpdateCategoriaProdutoInput) (*entity.CategoriaProdutoEntity, error) {
	// 1. Buscar categoria existente
	categoria, err := uc.repo.FindByID(ctx, input.TenantID, input.UnitID, input.ID)
	if err != nil {
		return nil, err
	}
	if categoria == nil {
		return nil, errors.New("categoria não encontrada")
	}

	// 2. Validar se novo nome já existe (se mudou)
	if categoria.Nome != input.Nome {
		exists, err := uc.repo.ExistsWithNome(ctx, input.TenantID, input.UnitID, input.Nome, &input.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("já existe uma categoria com este nome")
		}
		categoria.Nome = input.Nome
	}

	// 3. Atualizar campos
	categoria.SetDescricao(input.Descricao)
	if err := categoria.SetCor(input.Cor); err != nil {
		return nil, err
	}
	categoria.SetIcone(input.Icone)
	if err := categoria.SetCentroCusto(entity.CentroCusto(input.CentroCusto)); err != nil {
		return nil, err
	}
	if input.Ativa {
		categoria.Ativar()
	} else {
		categoria.Desativar()
	}

	// 4. Persistir
	if err := uc.repo.Update(ctx, categoria); err != nil {
		return nil, err
	}

	return categoria, nil
}

// =============================================================================
// DeleteCategoriaProdutoUseCase
// =============================================================================

type DeleteCategoriaProdutoUseCase struct {
	repo   port.CategoriaProdutoRepository
	logger *zap.Logger
}

func NewDeleteCategoriaProdutoUseCase(repo port.CategoriaProdutoRepository, logger *zap.Logger) *DeleteCategoriaProdutoUseCase {
	return &DeleteCategoriaProdutoUseCase{repo: repo, logger: logger}
}

func (uc *DeleteCategoriaProdutoUseCase) Execute(ctx context.Context, tenantID, unitID, id uuid.UUID) error {
	// 1. Verificar se categoria existe
	categoria, err := uc.repo.FindByID(ctx, tenantID, unitID, id)
	if err != nil {
		return err
	}
	if categoria == nil {
		return errors.New("categoria não encontrada")
	}

	// 2. Verificar se tem produtos vinculados
	count, err := uc.repo.CountProdutosVinculados(ctx, tenantID, unitID, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("não é possível excluir categoria com produtos vinculados")
	}

	// 3. Deletar
	return uc.repo.Delete(ctx, tenantID, unitID, id)
}

// =============================================================================
// ToggleCategoriaProdutoUseCase
// =============================================================================

type ToggleCategoriaProdutoUseCase struct {
	repo   port.CategoriaProdutoRepository
	logger *zap.Logger
}

func NewToggleCategoriaProdutoUseCase(repo port.CategoriaProdutoRepository, logger *zap.Logger) *ToggleCategoriaProdutoUseCase {
	return &ToggleCategoriaProdutoUseCase{repo: repo, logger: logger}
}

func (uc *ToggleCategoriaProdutoUseCase) Execute(ctx context.Context, tenantID, unitID, id uuid.UUID) (*entity.CategoriaProdutoEntity, error) {
	// 1. Buscar categoria
	categoria, err := uc.repo.FindByID(ctx, tenantID, unitID, id)
	if err != nil {
		return nil, err
	}
	if categoria == nil {
		return nil, errors.New("categoria não encontrada")
	}

	// 2. Toggle status
	if categoria.Ativa {
		categoria.Desativar()
	} else {
		categoria.Ativar()
	}

	// 3. Persistir
	if err := uc.repo.Update(ctx, categoria); err != nil {
		return nil, err
	}

	return categoria, nil
}
