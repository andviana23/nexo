package auth

import (
	"context"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

// =============================================================================
// ME USE CASE - VALTARIS v1.0
// Retorna dados do usuário logado
// =============================================================================

type MeUseCase struct {
	queries *db.Queries
	logger  *zap.Logger
}

func NewMeUseCase(queries *db.Queries, logger *zap.Logger) *MeUseCase {
	return &MeUseCase{
		queries: queries,
		logger:  logger,
	}
}

// Execute retorna dados do usuário autenticado
func (uc *MeUseCase) Execute(ctx context.Context, userID string) (*dto.MeResponse, error) {
	// Converter string para pgtype.UUID
	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		uc.logger.Error("ID de usuário inválido",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, domain.ErrUsuarioNaoEncontrado
	}

	user, err := uc.queries.GetUserByID(ctx, userUUID)
	if err != nil {
		uc.logger.Error("Usuário não encontrado",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, domain.ErrUsuarioNaoEncontrado
	}

	if user.Ativo != nil && !*user.Ativo {
		return nil, domain.ErrContaDesativada
	}

	return &dto.MeResponse{
		ID:       user.ID.String(),
		TenantID: user.TenantID.String(),
		Nome:     user.Nome,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}
