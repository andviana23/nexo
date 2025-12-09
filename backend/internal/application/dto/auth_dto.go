package dto

import (
	"time"
)

// =============================================================================
// AUTH DTOs - VALTARIS v1.0
// Baseado em FLUXO_LOGIN.md
// =============================================================================

// LoginRequest - Request para /auth/login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse - Response de /auth/login
type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

// RefreshResponse - Response de /auth/refresh
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

// UserResponse - Dados do usuário autenticado
type UserResponse struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Role     string `json:"role"` // owner, manager, recepcionista, barbeiro, contador
	// CurrentUnitID representa a unidade ativa do usuário (multi-unidade). Pode vir vazio quando não selecionada.
	CurrentUnitID string `json:"current_unit_id,omitempty"`
}

// MeResponse - Response de /auth/me (mesmo que UserResponse)
type MeResponse = UserResponse

// JWTClaims - Claims do JWT
type JWTClaims struct {
	UserID    string `json:"user_id"`
	TenantID  string `json:"tenant_id"`
	UnitID    string `json:"unit_id,omitempty"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// RefreshTokenData - Dados armazenados do refresh token
type RefreshTokenData struct {
	UserID    string
	TenantID  string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
