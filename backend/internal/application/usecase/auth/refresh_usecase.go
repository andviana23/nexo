package auth

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/infra/auth"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"go.uber.org/zap"
)

// =============================================================================
// REFRESH USE CASE - VALTARIS v1.0
// Baseado em FLUXO_LOGIN.md
// =============================================================================

type RefreshUseCase struct {
	queries    *db.Queries
	jwtManager *auth.JWTManager
	logger     *zap.Logger
}

func NewRefreshUseCase(queries *db.Queries, jwtManager *auth.JWTManager, logger *zap.Logger) *RefreshUseCase {
	return &RefreshUseCase{
		queries:    queries,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Execute renova access token usando refresh token
func (uc *RefreshUseCase) Execute(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	// 1. Validar refresh token no banco
	tokenData, err := uc.queries.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		uc.logger.Warn("Refresh token inválido ou não encontrado",
			zap.Error(err),
		)
		return nil, domain.ErrRefreshTokenInvalido
	}

	// 2. Buscar usuário
	user, err := uc.queries.GetUserByID(ctx, tokenData.UserID)
	if err != nil {
		uc.logger.Error("Usuário do refresh token não encontrado",
			zap.String("user_id", tokenData.UserID.String()),
			zap.Error(err),
		)
		return nil, domain.ErrUsuarioNaoEncontrado
	}

	// 3. Verificar se conta está ativa
	if user.Ativo != nil && !*user.Ativo {
		uc.logger.Warn("Tentativa de refresh com conta desativada",
			zap.String("user_id", user.ID.String()),
		)
		return nil, domain.ErrContaDesativada
	}

	// 4. Gerar novo access token
	accessToken, err := uc.jwtManager.GenerateAccessToken(
		user.ID.String(),
		user.TenantID.String(),
		"", // unit_id será populado quando vínculo unidade for implementado
		user.Email,
		user.Role,
	)
	if err != nil {
		uc.logger.Error("Erro ao gerar novo access token",
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	uc.logger.Info("Access token renovado com sucesso",
		zap.String("user_id", user.ID.String()),
	)

	return &dto.RefreshResponse{
		AccessToken: accessToken,
	}, nil
}
