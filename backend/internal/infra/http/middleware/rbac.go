package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// =============================================================================
// RBAC MIDDLEWARE - NEXO v1.0
// Controle de acesso baseado em roles
// =============================================================================

// Role representa uma role do sistema
type Role string

const (
	RoleOwner        Role = "OWNER"        // Dono da barbearia - acesso total
	RoleManager      Role = "MANAGER"      // Gerente - acesso administrativo
	RoleBarber       Role = "BARBER"       // Barbeiro - acesso restrito aos próprios dados
	RoleReceptionist Role = "RECEPTIONIST" // Recepcionista - acesso a agendamentos
)

// RBACConfig configura o middleware RBAC
type RBACConfig struct {
	AllowedRoles []Role      // Roles que podem acessar a rota
	Logger       *zap.Logger // Logger para auditoria
}

// RBAC cria um middleware que valida se o usuário tem uma das roles permitidas
func RBAC(config RBACConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extrair role do contexto (injetada pelo JWTMiddleware)
			userRole, ok := c.Get("role").(string)
			if !ok || userRole == "" {
				if config.Logger != nil {
					config.Logger.Warn("Role não encontrada no token",
						zap.String("path", c.Request().URL.Path),
					)
				}
				return echo.NewHTTPError(http.StatusForbidden, "Acesso negado: role não identificada")
			}

			// Verificar se a role do usuário está na lista de permitidas
			allowed := false
			for _, r := range config.AllowedRoles {
				if string(r) == userRole {
					allowed = true
					break
				}
			}

			if !allowed {
				if config.Logger != nil {
					config.Logger.Warn("Acesso negado por RBAC",
						zap.String("user_role", userRole),
						zap.String("path", c.Request().URL.Path),
						zap.Any("allowed_roles", config.AllowedRoles),
					)
				}
				return echo.NewHTTPError(http.StatusForbidden, "Acesso negado: permissão insuficiente")
			}

			if config.Logger != nil {
				config.Logger.Debug("Acesso permitido por RBAC",
					zap.String("user_role", userRole),
					zap.String("path", c.Request().URL.Path),
				)
			}

			return next(c)
		}
	}
}

// RequireRoles é um helper para criar middleware RBAC com roles específicas
func RequireRoles(logger *zap.Logger, roles ...Role) echo.MiddlewareFunc {
	return RBAC(RBACConfig{
		AllowedRoles: roles,
		Logger:       logger,
	})
}

// RequireOwnerOrManager exige OWNER ou MANAGER
func RequireOwnerOrManager(logger *zap.Logger) echo.MiddlewareFunc {
	return RequireRoles(logger, RoleOwner, RoleManager)
}

// RequireAdminAccess exige OWNER, MANAGER ou RECEPTIONIST (acesso administrativo)
func RequireAdminAccess(logger *zap.Logger) echo.MiddlewareFunc {
	return RequireRoles(logger, RoleOwner, RoleManager, RoleReceptionist)
}

// RequireAnyRole permite qualquer role autenticada (OWNER, MANAGER, BARBER, RECEPTIONIST)
func RequireAnyRole(logger *zap.Logger) echo.MiddlewareFunc {
	return RequireRoles(logger, RoleOwner, RoleManager, RoleBarber, RoleReceptionist)
}

// =============================================================================
// Helpers para verificação de escopo de barbeiro
// =============================================================================

// IsBarber verifica se o usuário atual é um barbeiro
func IsBarber(c echo.Context) bool {
	role, _ := c.Get("role").(string)
	return role == string(RoleBarber)
}

// GetProfessionalIDForBarber retorna o professional_id do barbeiro atual
// Barbeiros só podem ver/editar seus próprios agendamentos
// Para outras roles, retorna string vazia (sem restrição)
func GetProfessionalIDForBarber(c echo.Context) string {
	if !IsBarber(c) {
		return "" // Sem restrição para não-barbeiros
	}
	// O user_id do barbeiro É o professional_id
	userID, _ := c.Get("user_id").(string)
	return userID
}

// EnforceProfessionalScope verifica se o professional_id do request é o mesmo do barbeiro
// Retorna erro se barbeiro tentar acessar dados de outro profissional
func EnforceProfessionalScope(c echo.Context, requestedProfessionalID string) error {
	if !IsBarber(c) {
		return nil // Sem restrição para não-barbeiros
	}

	userID, _ := c.Get("user_id").(string)
	if requestedProfessionalID != "" && requestedProfessionalID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Acesso negado: você só pode acessar seus próprios agendamentos")
	}

	return nil
}
