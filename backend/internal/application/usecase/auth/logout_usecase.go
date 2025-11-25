package auth

import (
	"context"

	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"go.uber.org/zap"
)

// =============================================================================
// LOGOUT USE CASE - VALTARIS v1.0
// Invalida refresh token
// =============================================================================

type LogoutUseCase struct {
	queries *db.Queries
	logger  *zap.Logger
}

func NewLogoutUseCase(queries *db.Queries, logger *zap.Logger) *LogoutUseCase {
	return &LogoutUseCase{
		queries: queries,
		logger:  logger,
	}
}

// Execute invalida refresh token no banco
func (uc *LogoutUseCase) Execute(ctx context.Context, refreshToken string) error {
	err := uc.queries.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		// Ignora erro se token já foi deletado ou não existe
		uc.logger.Debug("Erro ao deletar refresh token (pode já estar deletado)",
			zap.Error(err),
		)
	}

	uc.logger.Info("Logout executado com sucesso")
	return nil
}
