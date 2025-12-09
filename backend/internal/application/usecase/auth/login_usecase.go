package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/infra/auth"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

// =============================================================================
// LOGIN USE CASE - VALTARIS v1.0
// Baseado em FLUXO_LOGIN.md
// =============================================================================

type LoginUseCase struct {
	queries    *db.Queries
	jwtManager *auth.JWTManager
	logger     *zap.Logger
}

func NewLoginUseCase(queries *db.Queries, jwtManager *auth.JWTManager, logger *zap.Logger) *LoginUseCase {
	return &LoginUseCase{
		queries:    queries,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Execute executa o fluxo de login
// Retorna: (accessToken, refreshToken, user, error)
func (uc *LoginUseCase) Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, string, error) {
	// 1. Buscar usuário por email
	user, err := uc.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		uc.logger.Warn("Login falhou - email não encontrado",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, "", domain.ErrEmailNaoEncontrado
	}

	// 2. Verificar se conta está ativa
	if user.Ativo != nil && !*user.Ativo {
		uc.logger.Warn("Login falhou - conta desativada",
			zap.String("email", req.Email),
			zap.String("user_id", user.ID.String()),
		)
		return nil, "", domain.ErrContaDesativada
	}

	// 3. Verificar senha
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		uc.logger.Warn("Login falhou - senha incorreta",
			zap.String("email", req.Email),
			zap.String("user_id", user.ID.String()),
		)
		return nil, "", domain.ErrSenhaIncorreta
	}

	// 4. Gerar access token (15 minutos)
	// unitID ainda não está implementado no banco; será preenchido após vínculo user↔unit
	accessToken, err := uc.jwtManager.GenerateAccessToken(
		user.ID.String(),
		user.TenantID.String(),
		"", // unit_id opcional até existir tabela/vínculo
		user.Email,
		user.Role,
	)
	if err != nil {
		uc.logger.Error("Erro ao gerar access token",
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
		return nil, "", fmt.Errorf("erro ao gerar token: %w", err)
	}

	// 5. Gerar refresh token (7 dias)
	refreshToken, err := uc.jwtManager.GenerateRefreshToken()
	if err != nil {
		uc.logger.Error("Erro ao gerar refresh token",
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
		return nil, "", fmt.Errorf("erro ao gerar refresh token: %w", err)
	}

	// 6. Armazenar refresh token no banco
	expiresAt := time.Now().Add(auth.RefreshTokenDuration)
	err = uc.queries.SaveRefreshToken(ctx, db.SaveRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		uc.logger.Error("Erro ao salvar refresh token",
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
		return nil, "", fmt.Errorf("erro ao salvar refresh token: %w", err)
	}

	// 7. Atualizar último login
	_ = uc.queries.UpdateLastLogin(ctx, user.ID)

	// 8. Log de sucesso
	uc.logger.Info("Login bem-sucedido",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
		zap.String("tenant_id", user.TenantID.String()),
		zap.String("role", user.Role),
	)

	// 9. Montar response
	response := &dto.LoginResponse{
		AccessToken: accessToken,
		User: dto.UserResponse{
			ID:       user.ID.String(),
			TenantID: user.TenantID.String(),
			Nome:     user.Nome,
			Email:    user.Email,
			Role:     user.Role,
			// current_unit_id vazio até o vínculo ser criado na etapa de banco de dados
			CurrentUnitID: "",
		},
	}

	return response, refreshToken, nil
}
