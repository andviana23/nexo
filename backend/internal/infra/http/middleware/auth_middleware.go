package middleware

import (
	"strings"

	"github.com/andviana23/barber-analytics-backend/internal/infra/auth"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// =============================================================================
// AUTH MIDDLEWARE - VALTARIS v1.0
// Baseado em FLUXO_LOGIN.md
// =============================================================================

// JWTMiddleware valida access token e injeta dados no context
func JWTMiddleware(jwtManager *auth.JWTManager, logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extrai token do header Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Token ausente")
				return echo.NewHTTPError(401, "Token não fornecido")
			}

			// Formato esperado: "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Warn("Formato de token inválido")
				return echo.NewHTTPError(401, "Formato de token inválido")
			}

			tokenString := parts[1]

			// Valida token
			claims, err := jwtManager.ValidateAccessToken(tokenString)
			if err != nil {
				logger.Warn("Token inválido",
					zap.Error(err),
				)
				return echo.NewHTTPError(401, "Token inválido ou expirado")
			}

			// Injeta dados no context do Echo
			c.Set("user_id", claims.UserID)
			c.Set("tenant_id", claims.TenantID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			logger.Debug("Token validado com sucesso",
				zap.String("user_id", claims.UserID),
				zap.String("tenant_id", claims.TenantID),
			)

			return next(c)
		}
	}
}

// GetUserID extrai user_id do context
func GetUserID(c echo.Context) string {
	userID, _ := c.Get("user_id").(string)
	return userID
}

// GetTenantID extrai tenant_id do context
func GetTenantID(c echo.Context) string {
	tenantID, _ := c.Get("tenant_id").(string)
	return tenantID
}

// GetUserRole extrai role do context
func GetUserRole(c echo.Context) string {
	role, _ := c.Get("role").(string)
	return role
}

// GetUserEmail extrai email do context
func GetUserEmail(c echo.Context) string {
	email, _ := c.Get("email").(string)
	return email
}

// GetUserIDFromContext é um alias para GetUserID (compatibilidade)
func GetUserIDFromContext(c echo.Context) string {
	return GetUserID(c)
}

// GetTenantIDFromContext é um alias para GetTenantID (compatibilidade)
func GetTenantIDFromContext(c echo.Context) string {
	return GetTenantID(c)
}
