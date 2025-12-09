package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Erros de domínio para UserUnit
var (
	ErrUserUnitJaExiste      = errors.New("usuário já está vinculado a esta unidade")
	ErrUserUnitNaoEncontrado = errors.New("vínculo usuário-unidade não encontrado")
	ErrUserSemAcesso         = errors.New("usuário não tem acesso a esta unidade")
)

// UserUnit representa o vínculo entre usuário e unidade
type UserUnit struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	UnitID       uuid.UUID
	IsDefault    bool
	RoleOverride *string // Papel específico nesta unidade (null = usa role do user)
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// UserUnitWithDetails inclui detalhes da unidade
type UserUnitWithDetails struct {
	UserUnit
	UnitNome    string
	UnitApelido *string
	UnitMatriz  bool
	UnitAtiva   bool
	TenantID    uuid.UUID
}

// NewUserUnit cria um novo vínculo usuário-unidade
func NewUserUnit(userID, unitID uuid.UUID, isDefault bool) *UserUnit {
	return &UserUnit{
		ID:           uuid.New(),
		UserID:       userID,
		UnitID:       unitID,
		IsDefault:    isDefault,
		CriadoEm:     time.Now(),
		AtualizadoEm: time.Now(),
	}
}

// SetAsDefault marca este vínculo como padrão
func (uu *UserUnit) SetAsDefault() {
	uu.IsDefault = true
	uu.AtualizadoEm = time.Now()
}

// UnsetDefault remove a marcação de padrão
func (uu *UserUnit) UnsetDefault() {
	uu.IsDefault = false
	uu.AtualizadoEm = time.Now()
}

// SetRoleOverride define um papel específico para esta unidade
func (uu *UserUnit) SetRoleOverride(role string) {
	if role != "" {
		uu.RoleOverride = &role
	} else {
		uu.RoleOverride = nil
	}
	uu.AtualizadoEm = time.Now()
}

// GetEffectiveRole retorna o papel efetivo (override ou fallback)
func (uu *UserUnit) GetEffectiveRole(defaultRole string) string {
	if uu.RoleOverride != nil && *uu.RoleOverride != "" {
		return *uu.RoleOverride
	}
	return defaultRole
}
