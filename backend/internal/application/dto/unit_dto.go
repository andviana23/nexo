package dto

// ============================================================================
// UNIT DTOs - Multi-Unidade
// ============================================================================

// CreateUnitRequest request para criar unidade
type CreateUnitRequest struct {
	Nome           string `json:"nome" validate:"required,min=2,max=100"`
	Apelido        string `json:"apelido,omitempty" validate:"omitempty,max=50"`
	Descricao      string `json:"descricao,omitempty"`
	EnderecoResumo string `json:"endereco_resumo,omitempty" validate:"omitempty,max=255"`
	Cidade         string `json:"cidade,omitempty" validate:"omitempty,max=100"`
	Estado         string `json:"estado,omitempty" validate:"omitempty,len=2"`
	Timezone       string `json:"timezone,omitempty" validate:"omitempty,max=50"`
	IsMatriz       bool   `json:"is_matriz,omitempty"`
}

// UpdateUnitRequest request para atualizar unidade
type UpdateUnitRequest struct {
	Nome           string `json:"nome,omitempty" validate:"omitempty,min=2,max=100"`
	Apelido        string `json:"apelido,omitempty" validate:"omitempty,max=50"`
	Descricao      string `json:"descricao,omitempty"`
	EnderecoResumo string `json:"endereco_resumo,omitempty" validate:"omitempty,max=255"`
	Cidade         string `json:"cidade,omitempty" validate:"omitempty,max=100"`
	Estado         string `json:"estado,omitempty" validate:"omitempty,len=2"`
	Timezone       string `json:"timezone,omitempty" validate:"omitempty,max=50"`
}

// UnitResponse response de unidade
type UnitResponse struct {
	ID             string  `json:"id"`
	TenantID       string  `json:"tenant_id"`
	Nome           string  `json:"nome"`
	Apelido        *string `json:"apelido,omitempty"`
	Descricao      *string `json:"descricao,omitempty"`
	EnderecoResumo *string `json:"endereco_resumo,omitempty"`
	Cidade         *string `json:"cidade,omitempty"`
	Estado         *string `json:"estado,omitempty"`
	Timezone       string  `json:"timezone"`
	Ativa          bool    `json:"ativa"`
	IsMatriz       bool    `json:"is_matriz"`
	CriadoEm       string  `json:"criado_em"`
	AtualizadoEm   string  `json:"atualizado_em"`
}

// ListUnitsResponse response de lista de unidades
type ListUnitsResponse struct {
	Units []UnitResponse `json:"units"`
	Total int            `json:"total"`
}

// ============================================================================
// USER UNIT DTOs - Vínculo usuário-unidade
// ============================================================================

// UserUnitResponse response de vínculo usuário-unidade
type UserUnitResponse struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	UnitID       string  `json:"unit_id"`
	UnitNome     string  `json:"unit_nome"`
	UnitApelido  *string `json:"unit_apelido,omitempty"`
	UnitMatriz   bool    `json:"unit_matriz"`
	UnitAtiva    bool    `json:"unit_ativa"`
	IsDefault    bool    `json:"is_default"`
	RoleOverride *string `json:"role_override,omitempty"`
	TenantID     string  `json:"tenant_id"`
}

// ListUserUnitsResponse response de lista de unidades do usuário
type ListUserUnitsResponse struct {
	Units []UserUnitResponse `json:"units"`
	Total int                `json:"total"`
}

// SwitchUnitRequest request para trocar de unidade
type SwitchUnitRequest struct {
	UnitID string `json:"unit_id" validate:"required,uuid"`
}

// SwitchUnitResponse response de troca de unidade
type SwitchUnitResponse struct {
	Unit        UserUnitResponse `json:"unit"`
	AccessToken string           `json:"access_token"` // Novo token com unit_id
}

// AddUserToUnitRequest request para adicionar usuário à unidade
type AddUserToUnitRequest struct {
	UserID       string  `json:"user_id" validate:"required,uuid"`
	UnitID       string  `json:"unit_id" validate:"required,uuid"`
	IsDefault    bool    `json:"is_default,omitempty"`
	RoleOverride *string `json:"role_override,omitempty"`
}

// SetDefaultUnitRequest request para definir unidade padrão
type SetDefaultUnitRequest struct {
	UnitID string `json:"unit_id" validate:"required,uuid"`
}
