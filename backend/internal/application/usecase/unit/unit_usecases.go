// Package unit contém os casos de uso relacionados a unidades/filiais
package unit

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// ============================================================================
// DTOs
// ============================================================================

// CreateUnitInput dados para criar unidade
type CreateUnitInput struct {
	TenantID       uuid.UUID
	Nome           string
	Apelido        *string
	Descricao      *string
	EnderecoResumo *string
	Cidade         *string
	Estado         *string
	Timezone       string
	IsMatriz       bool
}

// UpdateUnitInput dados para atualizar unidade
type UpdateUnitInput struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	Nome           *string
	Apelido        *string
	Descricao      *string
	EnderecoResumo *string
	Cidade         *string
	Estado         *string
	Timezone       *string
}

// UnitOutput resposta de unidade
type UnitOutput struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	Nome           string
	Apelido        *string
	Descricao      *string
	EnderecoResumo *string
	Cidade         *string
	Estado         *string
	Timezone       string
	Ativa          bool
	IsMatriz       bool
	CriadoEm       time.Time
	AtualizadoEm   time.Time
}

// ============================================================================
// CreateUnitUseCase
// ============================================================================

// CreateUnitUseCase caso de uso para criar unidade
type CreateUnitUseCase struct {
	unitRepo     port.UnitRepository
	userUnitRepo port.UserUnitRepository
}

// NewCreateUnitUseCase cria instância do caso de uso
func NewCreateUnitUseCase(unitRepo port.UnitRepository, userUnitRepo port.UserUnitRepository) *CreateUnitUseCase {
	return &CreateUnitUseCase{
		unitRepo:     unitRepo,
		userUnitRepo: userUnitRepo,
	}
}

// Execute executa o caso de uso
func (uc *CreateUnitUseCase) Execute(ctx context.Context, input CreateUnitInput) (*UnitOutput, error) {
	// Verificar se já existe unidade com mesmo nome
	existing, err := uc.unitRepo.FindByName(ctx, input.TenantID, input.Nome)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar nome: %w", err)
	}
	if existing != nil {
		return nil, entity.ErrUnitNomeDuplicado
	}

	// Criar entidade
	var unit *entity.Unit
	if input.IsMatriz {
		unit, err = entity.NewMatrizUnit(input.TenantID, input.Nome)
	} else {
		unit, err = entity.NewUnit(input.TenantID, input.Nome)
	}
	if err != nil {
		return nil, err
	}

	// Preencher campos opcionais
	if input.Apelido != nil && *input.Apelido != "" {
		unit.SetApelido(*input.Apelido)
	}
	if input.Descricao != nil && *input.Descricao != "" {
		unit.Descricao = input.Descricao
	}
	if input.EnderecoResumo != nil || input.Cidade != nil || input.Estado != nil {
		endereco := ""
		cidade := ""
		estado := ""
		if input.EnderecoResumo != nil {
			endereco = *input.EnderecoResumo
		}
		if input.Cidade != nil {
			cidade = *input.Cidade
		}
		if input.Estado != nil {
			estado = *input.Estado
		}
		unit.SetEndereco(endereco, cidade, estado)
	}
	if input.Timezone != "" {
		unit.Timezone = input.Timezone
	}

	// Persistir
	if err := uc.unitRepo.Create(ctx, unit); err != nil {
		return nil, fmt.Errorf("erro ao criar unidade: %w", err)
	}

	return toUnitOutput(unit), nil
}

// ============================================================================
// ListUnitsUseCase
// ============================================================================

// ListUnitsUseCase caso de uso para listar unidades
type ListUnitsUseCase struct {
	unitRepo port.UnitRepository
}

// NewListUnitsUseCase cria instância do caso de uso
func NewListUnitsUseCase(unitRepo port.UnitRepository) *ListUnitsUseCase {
	return &ListUnitsUseCase{unitRepo: unitRepo}
}

// Execute executa o caso de uso
func (uc *ListUnitsUseCase) Execute(ctx context.Context, tenantID uuid.UUID, onlyActive bool) ([]UnitOutput, error) {
	var units []*entity.Unit
	var err error

	if onlyActive {
		units, err = uc.unitRepo.ListActive(ctx, tenantID)
	} else {
		units, err = uc.unitRepo.List(ctx, tenantID)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao listar unidades: %w", err)
	}

	result := make([]UnitOutput, len(units))
	for i, u := range units {
		result[i] = *toUnitOutput(u)
	}

	return result, nil
}

// ============================================================================
// GetUnitUseCase
// ============================================================================

// GetUnitUseCase caso de uso para buscar unidade por ID
type GetUnitUseCase struct {
	unitRepo port.UnitRepository
}

// NewGetUnitUseCase cria instância do caso de uso
func NewGetUnitUseCase(unitRepo port.UnitRepository) *GetUnitUseCase {
	return &GetUnitUseCase{unitRepo: unitRepo}
}

// Execute executa o caso de uso
func (uc *GetUnitUseCase) Execute(ctx context.Context, unitID, tenantID uuid.UUID) (*UnitOutput, error) {
	unit, err := uc.unitRepo.FindByID(ctx, tenantID, unitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar unidade: %w", err)
	}
	if unit == nil {
		return nil, entity.ErrUnitNaoEncontrada
	}

	return toUnitOutput(unit), nil
}

// ============================================================================
// UpdateUnitUseCase
// ============================================================================

// UpdateUnitUseCase caso de uso para atualizar unidade
type UpdateUnitUseCase struct {
	unitRepo port.UnitRepository
}

// NewUpdateUnitUseCase cria instância do caso de uso
func NewUpdateUnitUseCase(unitRepo port.UnitRepository) *UpdateUnitUseCase {
	return &UpdateUnitUseCase{unitRepo: unitRepo}
}

// Execute executa o caso de uso
func (uc *UpdateUnitUseCase) Execute(ctx context.Context, input UpdateUnitInput) (*UnitOutput, error) {
	// Buscar unidade existente
	unit, err := uc.unitRepo.FindByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar unidade: %w", err)
	}
	if unit == nil {
		return nil, entity.ErrUnitNaoEncontrada
	}

	// Verificar duplicidade de nome (se alterado)
	if input.Nome != nil && *input.Nome != "" && *input.Nome != unit.Nome {
		existing, err := uc.unitRepo.FindByName(ctx, input.TenantID, *input.Nome)
		if err != nil {
			return nil, fmt.Errorf("erro ao verificar nome: %w", err)
		}
		if existing != nil && existing.ID != unit.ID {
			return nil, entity.ErrUnitNomeDuplicado
		}
		unit.Nome = *input.Nome
	}

	// Atualizar campos
	if input.Apelido != nil && *input.Apelido != "" {
		unit.SetApelido(*input.Apelido)
	}
	if input.Descricao != nil {
		unit.Descricao = input.Descricao
	}

	endereco := ""
	cidade := ""
	estado := ""
	if input.EnderecoResumo != nil {
		endereco = *input.EnderecoResumo
	} else if unit.EnderecoResumo != nil {
		endereco = *unit.EnderecoResumo
	}
	if input.Cidade != nil {
		cidade = *input.Cidade
	} else if unit.Cidade != nil {
		cidade = *unit.Cidade
	}
	if input.Estado != nil {
		estado = *input.Estado
	} else if unit.Estado != nil {
		estado = *unit.Estado
	}
	if input.EnderecoResumo != nil || input.Cidade != nil || input.Estado != nil {
		unit.SetEndereco(endereco, cidade, estado)
	}

	if input.Timezone != nil && *input.Timezone != "" {
		unit.Timezone = *input.Timezone
	}

	// Persistir
	if err := uc.unitRepo.Update(ctx, unit); err != nil {
		return nil, fmt.Errorf("erro ao atualizar unidade: %w", err)
	}

	return toUnitOutput(unit), nil
}

// ============================================================================
// ToggleUnitUseCase
// ============================================================================

// ToggleUnitUseCase caso de uso para alternar status
type ToggleUnitUseCase struct {
	unitRepo port.UnitRepository
}

// NewToggleUnitUseCase cria instância do caso de uso
func NewToggleUnitUseCase(unitRepo port.UnitRepository) *ToggleUnitUseCase {
	return &ToggleUnitUseCase{unitRepo: unitRepo}
}

// Execute executa o caso de uso
func (uc *ToggleUnitUseCase) Execute(ctx context.Context, unitID, tenantID uuid.UUID) (*UnitOutput, error) {
	unit, err := uc.unitRepo.ToggleStatus(ctx, tenantID, unitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao alternar status: %w", err)
	}

	return toUnitOutput(unit), nil
}

// ============================================================================
// DeleteUnitUseCase
// ============================================================================

// DeleteUnitUseCase caso de uso para excluir unidade
type DeleteUnitUseCase struct {
	unitRepo     port.UnitRepository
	userUnitRepo port.UserUnitRepository
}

// NewDeleteUnitUseCase cria instância do caso de uso
func NewDeleteUnitUseCase(unitRepo port.UnitRepository, userUnitRepo port.UserUnitRepository) *DeleteUnitUseCase {
	return &DeleteUnitUseCase{
		unitRepo:     unitRepo,
		userUnitRepo: userUnitRepo,
	}
}

// Execute executa o caso de uso
func (uc *DeleteUnitUseCase) Execute(ctx context.Context, unitID, tenantID uuid.UUID) error {
	// Buscar unidade
	unit, err := uc.unitRepo.FindByID(ctx, tenantID, unitID)
	if err != nil {
		return fmt.Errorf("erro ao buscar unidade: %w", err)
	}
	if unit == nil {
		return entity.ErrUnitNaoEncontrada
	}

	// Verificar se pode excluir
	if err := unit.CanDelete(); err != nil {
		return err
	}

	// Remover vínculos de usuários
	if err := uc.userUnitRepo.DeleteAllByUnit(ctx, unitID); err != nil {
		return fmt.Errorf("erro ao remover vínculos: %w", err)
	}

	// Excluir unidade
	if err := uc.unitRepo.Delete(ctx, tenantID, unitID); err != nil {
		return fmt.Errorf("erro ao excluir unidade: %w", err)
	}

	return nil
}

// ============================================================================
// Helpers
// ============================================================================

func toUnitOutput(unit *entity.Unit) *UnitOutput {
	return &UnitOutput{
		ID:             unit.ID,
		TenantID:       unit.TenantID,
		Nome:           unit.Nome,
		Apelido:        unit.Apelido,
		Descricao:      unit.Descricao,
		EnderecoResumo: unit.EnderecoResumo,
		Cidade:         unit.Cidade,
		Estado:         unit.Estado,
		Timezone:       unit.Timezone,
		Ativa:          unit.Ativa,
		IsMatriz:       unit.IsMatriz,
		CriadoEm:       unit.CriadoEm,
		AtualizadoEm:   unit.AtualizadoEm,
	}
}
