package handler

import (
	"net/http"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	authUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/auth"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	mw "github.com/andviana23/barber-analytics-backend/internal/infra/http/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// =============================================================================
// AUTH HANDLER - VALTARIS v1.0
// Baseado em FLUXO_LOGIN.md
// =============================================================================

type AuthHandler struct {
	loginUC   *authUC.LoginUseCase
	refreshUC *authUC.RefreshUseCase
	meUC      *authUC.MeUseCase
	logoutUC  *authUC.LogoutUseCase
	validator *validator.Validate
	logger    *zap.Logger
}

func NewAuthHandler(
	loginUC *authUC.LoginUseCase,
	refreshUC *authUC.RefreshUseCase,
	meUC *authUC.MeUseCase,
	logoutUC *authUC.LogoutUseCase,
	logger *zap.Logger,
) *AuthHandler {
	return &AuthHandler{
		loginUC:   loginUC,
		refreshUC: refreshUC,
		meUC:      meUC,
		logoutUC:  logoutUC,
		validator: validator.New(),
		logger:    logger,
	}
}

// Login - POST /auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest

	// Bind e validação
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Dados inválidos",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email ou senha inválidos",
		})
	}

	// Executa login
	response, refreshToken, err := h.loginUC.Execute(c.Request().Context(), req)
	if err != nil {
		// Erros específicos do domínio (conforme FLUXO_LOGIN.md)
		switch err {
		case domain.ErrEmailNaoEncontrado:
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Email não encontrado",
			})
		case domain.ErrSenhaIncorreta:
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Senha incorreta",
			})
		case domain.ErrContaDesativada:
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Conta desativada",
			})
		default:
			h.logger.Error("Erro ao fazer login", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Erro interno do servidor",
			})
		}
	}

	// Define cookie HttpOnly com refresh token (7 dias)
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // TODO: em produção = true (HTTPS)
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(7 * 24 * time.Hour.Seconds()), // 7 dias
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, response)
}

// Refresh - POST /auth/refresh
func (h *AuthHandler) Refresh(c echo.Context) error {
	// Extrai refresh token do cookie
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Refresh token não encontrado",
		})
	}

	refreshToken := cookie.Value

	// Executa refresh
	response, err := h.refreshUC.Execute(c.Request().Context(), refreshToken)
	if err != nil {
		switch err {
		case domain.ErrRefreshTokenInvalido:
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Refresh token inválido ou expirado",
			})
		case domain.ErrContaDesativada:
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Conta desativada",
			})
		default:
			h.logger.Error("Erro ao renovar token", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Erro interno do servidor",
			})
		}
	}

	return c.JSON(http.StatusOK, response)
}

// Me - GET /auth/me
func (h *AuthHandler) Me(c echo.Context) error {
	// UserID vem do middleware JWT
	userID := mw.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Não autenticado",
		})
	}

	response, err := h.meUC.Execute(c.Request().Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrUsuarioNaoEncontrado:
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Usuário não encontrado",
			})
		case domain.ErrContaDesativada:
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Conta desativada",
			})
		default:
			h.logger.Error("Erro ao buscar dados do usuário", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Erro interno do servidor",
			})
		}
	}

	return c.JSON(http.StatusOK, response)
}

// Logout - POST /auth/logout
func (h *AuthHandler) Logout(c echo.Context) error {
	// Extrai refresh token do cookie
	cookie, err := c.Cookie("refresh_token")
	if err == nil {
		// Invalida token no banco
		_ = h.logoutUC.Execute(c.Request().Context(), cookie.Value)
	}

	// Remove cookie
	cookie = &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Deleta cookie
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logout realizado com sucesso",
	})
}
