// Package unit contém os casos de uso relacionados a unidades/filiais
package unit

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
)

// ============================================================================
// DTOs para UserUnit
// ============================================================================

// UserUnitOutput resposta de vínculo usuário-unidade
type UserUnitOutput struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	UnitID        uuid.UUID
	UnitNome      string
	UnitApelido   *string
	UnitMatriz    bool
	UnitAtiva     bool
	IsDefaultUnit bool
	RoleOverride  *string
	TenantID      uuid.UUID
}

// SwitchUnitInput dados para trocar de unidade
type SwitchUnitInput struct {
	UserID uuid.UUID
	UnitID uuid.UUID
}

// ============================================================================
// ListUserUnitsUseCase
// ============================================================================

// ListUserUnitsUseCase caso de uso para listar unidades do usuário
type ListUserUnitsUseCase struct {
	userUnitRepo port.UserUnitRepository
}

// NewListUserUnitsUseCase cria instância do caso de uso
func NewListUserUnitsUseCase(userUnitRepo port.UserUnitRepository) *ListUserUnitsUseCase {
	return &ListUserUnitsUseCase{userUnitRepo: userUnitRepo}
}

// Execute executa o caso de uso
func (uc *ListUserUnitsUseCase) Execute(ctx context.Context, userID uuid.UUID) ([]entity.UserUnitWithDetails, error) {
	units, err := uc.userUnitRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar unidades: %w", err)
	}

	result := make([]entity.UserUnitWithDetails, len(units))
	for i, uu := range units {
		result[i] = *uu
	}

	return result, nil
}

// ============================================================================
// GetUserDefaultUnitUseCase
// ============================================================================

// GetUserDefaultUnitUseCase caso de uso para buscar unidade padrão
type GetUserDefaultUnitUseCase struct {
	userUnitRepo port.UserUnitRepository
}

// NewGetUserDefaultUnitUseCase cria instância do caso de uso
func NewGetUserDefaultUnitUseCase(userUnitRepo port.UserUnitRepository) *GetUserDefaultUnitUseCase {
	return &GetUserDefaultUnitUseCase{userUnitRepo: userUnitRepo}
}

// Execute executa o caso de uso
func (uc *GetUserDefaultUnitUseCase) Execute(ctx context.Context, userID uuid.UUID) (*UserUnitOutput, error) {
	uu, err := uc.userUnitRepo.FindUserDefaultUnit(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar unidade padrão: %w", err)
	}
	if uu == nil {
		return nil, nil // Usuário pode não ter unidade padrão definida
	}

	return toUserUnitOutput(uu), nil
}

// ============================================================================
// SwitchUnitUseCase
// ============================================================================

// SwitchUnitUseCase caso de uso para trocar de unidade
type SwitchUnitUseCase struct {
	userUnitRepo port.UserUnitRepository
	unitRepo     port.UnitRepository
}

// NewSwitchUnitUseCase cria instância do caso de uso
func NewSwitchUnitUseCase(userUnitRepo port.UserUnitRepository, unitRepo port.UnitRepository) *SwitchUnitUseCase {
	return &SwitchUnitUseCase{
		userUnitRepo: userUnitRepo,
		unitRepo:     unitRepo,
	}
}

// Execute executa o caso de uso - retorna dados para novo token
func (uc *SwitchUnitUseCase) Execute(ctx context.Context, input SwitchUnitInput) (*UserUnitOutput, error) {
	// Verificar se usuário tem acesso à unidade
	hasAccess, err := uc.userUnitRepo.CheckAccess(ctx, input.UserID, input.UnitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar acesso: %w", err)
	}
	if !hasAccess {
		return nil, entity.ErrUserSemAcesso
	}

	// Buscar detalhes do vínculo
	uu, err := uc.userUnitRepo.FindByUserAndUnit(ctx, input.UserID, input.UnitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar vínculo: %w", err)
	}
	if uu == nil {
		return nil, entity.ErrUserUnitNaoEncontrado
	}

	// Buscar unidade para detalhes
	unit, err := uc.unitRepo.FindByID(ctx, uuid.Nil, input.UnitID) // tenantID não validado aqui pois já passou no CheckAccess
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar unidade: %w", err)
	}

	return &UserUnitOutput{
		ID:            uu.ID,
		UserID:        uu.UserID,
		UnitID:        uu.UnitID,
		UnitNome:      unit.Nome,
		UnitApelido:   unit.Apelido,
		UnitMatriz:    unit.IsMatriz,
		UnitAtiva:     unit.Ativa,
		IsDefaultUnit: uu.IsDefault,
		RoleOverride:  uu.RoleOverride,
		TenantID:      unit.TenantID,
	}, nil
}

// ============================================================================
// SetDefaultUnitUseCase
// ============================================================================

// SetDefaultUnitUseCase caso de uso para definir unidade padrão
type SetDefaultUnitUseCase struct {
	userUnitRepo port.UserUnitRepository
}

// NewSetDefaultUnitUseCase cria instância do caso de uso
func NewSetDefaultUnitUseCase(userUnitRepo port.UserUnitRepository) *SetDefaultUnitUseCase {
	return &SetDefaultUnitUseCase{userUnitRepo: userUnitRepo}
}

// Execute executa o caso de uso
func (uc *SetDefaultUnitUseCase) Execute(ctx context.Context, userID, unitID uuid.UUID) error {
	// Verificar se usuário tem acesso à unidade
	hasAccess, err := uc.userUnitRepo.CheckAccess(ctx, userID, unitID)
	if err != nil {
		return fmt.Errorf("erro ao verificar acesso: %w", err)
	}
	if !hasAccess {
		return entity.ErrUserSemAcesso
	}

	// Definir como padrão
	if err := uc.userUnitRepo.SetDefault(ctx, userID, unitID); err != nil {
		return fmt.Errorf("erro ao definir unidade padrão: %w", err)
	}

	return nil
}

// ============================================================================
// AddUserToUnitUseCase
// ============================================================================

// AddUserToUnitInput dados para adicionar usuário à unidade
type AddUserToUnitInput struct {
	TenantID     uuid.UUID
	UserID       uuid.UUID
	UnitID       uuid.UUID
	IsDefault    bool
	RoleOverride *string
}

// AddUserToUnitUseCase caso de uso para adicionar usuário à unidade
type AddUserToUnitUseCase struct {
	userUnitRepo port.UserUnitRepository
	unitRepo     port.UnitRepository
}

// NewAddUserToUnitUseCase cria instância do caso de uso
func NewAddUserToUnitUseCase(userUnitRepo port.UserUnitRepository, unitRepo port.UnitRepository) *AddUserToUnitUseCase {
	return &AddUserToUnitUseCase{
		userUnitRepo: userUnitRepo,
		unitRepo:     unitRepo,
	}
}

// Execute executa o caso de uso
func (uc *AddUserToUnitUseCase) Execute(ctx context.Context, input AddUserToUnitInput) (*entity.UserUnit, error) {
	// Verificar se unidade existe e pertence ao tenant
	unit, err := uc.unitRepo.FindByID(ctx, input.TenantID, input.UnitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar unidade: %w", err)
	}
	if unit == nil {
		return nil, entity.ErrUnitNaoEncontrada
	}

	// Verificar se vínculo já existe
	existing, err := uc.userUnitRepo.FindByUserAndUnit(ctx, input.UserID, input.UnitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar vínculo: %w", err)
	}
	if existing != nil {
		return nil, entity.ErrUserUnitJaExiste
	}

	// Criar vínculo
	userUnit := entity.NewUserUnit(input.UserID, input.UnitID, input.IsDefault)
	if input.RoleOverride != nil {
		userUnit.SetRoleOverride(*input.RoleOverride)
	}

	if err := uc.userUnitRepo.Create(ctx, userUnit); err != nil {
		return nil, fmt.Errorf("erro ao criar vínculo: %w", err)
	}

	return userUnit, nil
}

// ============================================================================
// RemoveUserFromUnitUseCase (alias UnlinkUserFromUnitUseCase)
// ============================================================================

// UnlinkUserFromUnitUseCase caso de uso para remover usuário da unidade
type UnlinkUserFromUnitUseCase struct {
	userUnitRepo port.UserUnitRepository
}

// NewUnlinkUserFromUnitUseCase cria instância do caso de uso
func NewUnlinkUserFromUnitUseCase(userUnitRepo port.UserUnitRepository) *UnlinkUserFromUnitUseCase {
	return &UnlinkUserFromUnitUseCase{userUnitRepo: userUnitRepo}
}

// Execute executa o caso de uso
func (uc *UnlinkUserFromUnitUseCase) Execute(ctx context.Context, userID, unitID uuid.UUID) error {
	if err := uc.userUnitRepo.Delete(ctx, userID, unitID); err != nil {
		return fmt.Errorf("erro ao remover vínculo: %w", err)
	}
	return nil
}

// ============================================================================
// CheckUserAccessToUnitUseCase
// ============================================================================

// CheckUserAccessToUnitUseCase caso de uso para verificar acesso
type CheckUserAccessToUnitUseCase struct {
	userUnitRepo port.UserUnitRepository
}

// NewCheckUserAccessToUnitUseCase cria instância do caso de uso
func NewCheckUserAccessToUnitUseCase(userUnitRepo port.UserUnitRepository) *CheckUserAccessToUnitUseCase {
	return &CheckUserAccessToUnitUseCase{userUnitRepo: userUnitRepo}
}

// Execute executa o caso de uso
func (uc *CheckUserAccessToUnitUseCase) Execute(ctx context.Context, userID, unitID uuid.UUID) (bool, error) {
	return uc.userUnitRepo.CheckAccess(ctx, userID, unitID)
}

// ============================================================================
// GetDefaultUnitUseCase
// ============================================================================

// GetDefaultUnitUseCase caso de uso para buscar unidade padrão do usuário
type GetDefaultUnitUseCase struct {
	userUnitRepo port.UserUnitRepository
}

// NewGetDefaultUnitUseCase cria instância do caso de uso
func NewGetDefaultUnitUseCase(userUnitRepo port.UserUnitRepository) *GetDefaultUnitUseCase {
	return &GetDefaultUnitUseCase{userUnitRepo: userUnitRepo}
}

// Execute executa o caso de uso
func (uc *GetDefaultUnitUseCase) Execute(ctx context.Context, userID uuid.UUID) (*entity.UserUnitWithDetails, error) {
	return uc.userUnitRepo.FindUserDefaultUnit(ctx, userID)
}

// ============================================================================
// ListUnitUsersUseCase
// ============================================================================

// ListUnitUsersUseCase caso de uso para listar usuários de uma unidade
type ListUnitUsersUseCase struct {
	userUnitRepo port.UserUnitRepository
}

// NewListUnitUsersUseCase cria instância do caso de uso
func NewListUnitUsersUseCase(userUnitRepo port.UserUnitRepository) *ListUnitUsersUseCase {
	return &ListUnitUsersUseCase{userUnitRepo: userUnitRepo}
}

// Execute executa o caso de uso
func (uc *ListUnitUsersUseCase) Execute(ctx context.Context, unitID uuid.UUID) ([]entity.UserUnitWithDetails, error) {
	users, err := uc.userUnitRepo.ListByUnit(ctx, unitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar usuários: %w", err)
	}

	result := make([]entity.UserUnitWithDetails, len(users))
	for i, uu := range users {
		result[i] = entity.UserUnitWithDetails{UserUnit: *uu}
	}

	return result, nil
}

// ============================================================================
// LinkUserToUnitUseCase (alias AddUserToUnitUseCase)
// ============================================================================

// LinkUserToUnitInput dados para vincular usuário à unidade
type LinkUserToUnitInput struct {
	UserID       uuid.UUID
	UnitID       uuid.UUID
	IsDefault    bool
	RoleOverride *string
}

// LinkUserToUnitUseCase caso de uso para vincular usuário à unidade
type LinkUserToUnitUseCase struct {
	userUnitRepo port.UserUnitRepository
}

// NewLinkUserToUnitUseCase cria instância do caso de uso
func NewLinkUserToUnitUseCase(userUnitRepo port.UserUnitRepository) *LinkUserToUnitUseCase {
	return &LinkUserToUnitUseCase{userUnitRepo: userUnitRepo}
}

// Execute executa o caso de uso
func (uc *LinkUserToUnitUseCase) Execute(ctx context.Context, input LinkUserToUnitInput) (*entity.UserUnit, error) {
	// Verificar se vínculo já existe
	existing, err := uc.userUnitRepo.FindByUserAndUnit(ctx, input.UserID, input.UnitID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar vínculo: %w", err)
	}
	if existing != nil {
		return nil, entity.ErrUserUnitJaExiste
	}

	// Criar vínculo
	userUnit := entity.NewUserUnit(input.UserID, input.UnitID, input.IsDefault)
	if input.RoleOverride != nil {
		userUnit.SetRoleOverride(*input.RoleOverride)
	}

	if err := uc.userUnitRepo.Create(ctx, userUnit); err != nil {
		return nil, fmt.Errorf("erro ao criar vínculo: %w", err)
	}

	return userUnit, nil
}

// ============================================================================
// Helpers
// ============================================================================

func toUserUnitOutput(uu *entity.UserUnitWithDetails) *UserUnitOutput {
	return &UserUnitOutput{
		ID:            uu.ID,
		UserID:        uu.UserID,
		UnitID:        uu.UnitID,
		UnitNome:      uu.UnitNome,
		UnitApelido:   uu.UnitApelido,
		UnitMatriz:    uu.UnitMatriz,
		UnitAtiva:     uu.UnitAtiva,
		IsDefaultUnit: uu.IsDefault,
		RoleOverride:  uu.RoleOverride,
		TenantID:      uu.TenantID,
	}
}
