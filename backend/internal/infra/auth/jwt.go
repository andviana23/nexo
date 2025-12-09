package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/golang-jwt/jwt/v5"
)

// =============================================================================
// JWT MANAGER - VALTARIS v1.0
// Baseado em FLUXO_LOGIN.md
// =============================================================================

var (
	// Access token: 15 minutos (conforme FLUXO_LOGIN.md)
	AccessTokenDuration = 15 * time.Minute
	// Refresh token: 7 dias (conforme FLUXO_LOGIN.md)
	RefreshTokenDuration = 7 * 24 * time.Hour
)

// JWTManager gerencia criação e validação de tokens JWT
type JWTManager struct {
	secretKey string
}

// NewJWTManager cria nova instância do JWT Manager
func NewJWTManager() *JWTManager {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// ATENÇÃO: em produção, DEVE estar no .env
		secret = "valtaris-dev-secret-change-in-production"
	}
	return &JWTManager{secretKey: secret}
}

// GenerateAccessToken gera novo access token JWT (15 minutos)
// unitID é opcional e representa a unidade atualmente selecionada pelo usuário.
func (jm *JWTManager) GenerateAccessToken(userID, tenantID, unitID, email, role string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(AccessTokenDuration)

	claims := jwt.MapClaims{
		"user_id":   userID,
		"tenant_id": tenantID,
		"unit_id":   unitID,
		"email":     email,
		"role":      role,
		"iat":       now.Unix(),
		"exp":       expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.secretKey))
}

// GenerateRefreshToken gera novo refresh token aleatório (7 dias)
// Este token será armazenado em cookie HttpOnly
func (jm *JWTManager) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("erro ao gerar refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ValidateAccessToken valida e extrai claims do access token
func (jm *JWTManager) ValidateAccessToken(tokenString string) (*dto.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifica se o método de assinatura é HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inválido: %v", token.Header["alg"])
		}
		return []byte(jm.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claims inválidos")
	}

	// Extrai claims
	userID, _ := claims["user_id"].(string)
	tenantID, _ := claims["tenant_id"].(string)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)
	unitID, _ := claims["unit_id"].(string)
	iat, _ := claims["iat"].(float64)
	exp, _ := claims["exp"].(float64)

	return &dto.JWTClaims{
		UserID:    userID,
		TenantID:  tenantID,
		UnitID:    unitID,
		Email:     email,
		Role:      role,
		IssuedAt:  int64(iat),
		ExpiresAt: int64(exp),
	}, nil
}
