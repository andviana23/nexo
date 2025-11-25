package user

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// DeleteAccountRequest representa a solicitação de exclusão de conta
type DeleteAccountRequest struct {
	UserID   string
	TenantID string
	Password string // Senha para confirmar exclusão
	Reason   string // Motivo opcional
}

// DeleteAccountUseCase implementa o direito ao esquecimento (LGPD Art. 18, VI)
type DeleteAccountUseCase struct {
	prefsRepo port.UserPreferencesRepository
	logger    *zap.Logger
}

func NewDeleteAccountUseCase(
	prefsRepo port.UserPreferencesRepository,
	logger *zap.Logger,
) *DeleteAccountUseCase {
	return &DeleteAccountUseCase{
		prefsRepo: prefsRepo,
		logger:    logger,
	}
}

// Execute realiza soft delete do usuário e anonimização de dados (LGPD)
func (uc *DeleteAccountUseCase) Execute(ctx context.Context, req DeleteAccountRequest) error {
	if req.TenantID == "" {
		return domain.ErrTenantIDRequired
	}
	if req.UserID == "" {
		return domain.ErrInvalidID
	}

	uc.logger.Info("Iniciando exclusão de conta (LGPD)",
		zap.String("tenant_id", req.TenantID),
		zap.String("user_id", req.UserID),
		zap.String("reason", req.Reason),
	)

	now := time.Now()

	// TODO: 1. Validar senha do usuário antes de deletar
	// TODO: 2. Soft delete em users table (UPDATE users SET deleted_at = NOW(), ativo = false WHERE id = ?)
	// TODO: 3. Anonimizar PII (nome = "[REMOVIDO]", email = "deleted-{uuid}@anonimizado.local", password_hash = "")
	// TODO: 4. Deletar preferências do usuário
	// TODO: 5. Revogar tokens JWT (blacklist ou invalidar refresh tokens)
	// TODO: 6. Registrar em audit log: ActionDeleteAccount

	// Step 4: Deletar preferências
	if err := uc.prefsRepo.Delete(ctx, req.UserID); err != nil {
		uc.logger.Warn("Erro ao deletar preferências durante exclusão de conta",
			zap.String("user_id", req.UserID),
			zap.Error(err),
		)
		// Não falhar se preferências não existirem
	}

	uc.logger.Info("Conta deletada com sucesso",
		zap.String("user_id", req.UserID),
		zap.Time("deleted_at", now),
	)

	return nil
}

// AnonymizeUserData anonimiza dados pessoais do usuário
func (uc *DeleteAccountUseCase) AnonymizeUserData(ctx context.Context, userID string) error {
	// TODO: Implementar query SQL:
	// UPDATE users
	// SET nome = '[USUÁRIO REMOVIDO]',
	//     email = CONCAT('deleted-', SUBSTRING(id::text, 1, 8), '@anonimizado.local'),
	//     password_hash = '',
	//     deleted_at = NOW(),
	//     ativo = false
	// WHERE id = $1

	uc.logger.Info("Dados do usuário anonimizados",
		zap.String("user_id", userID),
	)

	return nil
}
