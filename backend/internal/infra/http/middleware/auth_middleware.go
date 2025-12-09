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

// roleMapping mapeia roles do banco (português/minúsculo) para roles RBAC (inglês/maiúsculo)
var roleMapping = map[string]string{
	// Português (banco de dados)
	"owner":         "OWNER",
	"manager":       "MANAGER",
	"barbeiro":      "BARBER",
	"barber":        "BARBER",
	"recepcionista": "RECEPTIONIST",
	"receptionist":  "RECEPTIONIST",
	"contador":      "ACCOUNTANT",
	// Inglês (já normalizado)
	"OWNER":        "OWNER",
	"MANAGER":      "MANAGER",
	"BARBER":       "BARBER",
	"RECEPTIONIST": "RECEPTIONIST",
	"ACCOUNTANT":   "ACCOUNTANT",
}

// normalizeRole converte role do banco para formato RBAC padronizado
func normalizeRole(role string) string {
	if normalized, ok := roleMapping[strings.ToLower(role)]; ok {
		return normalized
	}
	// Se não encontrar mapeamento, retorna uppercase como fallback
	return strings.ToUpper(role)
}

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
			c.Set("role", normalizeRole(claims.Role)) // Normaliza role para RBAC

			// Unidade atual (multi-unidade). Prioridade: claim → header X-Unit-ID → vazio
			unitID := claims.UnitID
			if unitID == "" {
				unitID = c.Request().Header.Get("X-Unit-ID")
			}
			if unitID != "" {
				c.Set("unit_id", unitID)
			}

			logger.Debug("Token validado com sucesso",
				zap.String("user_id", claims.UserID),
				zap.String("tenant_id", claims.TenantID),
				zap.String("original_role", claims.Role),
				zap.String("normalized_role", normalizeRole(claims.Role)),
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

// GetUnitID extrai unit_id do context (pode retornar vazio se não definido)
func GetUnitID(c echo.Context) string {
	unitID, _ := c.Get("unit_id").(string)
	return unitID
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
