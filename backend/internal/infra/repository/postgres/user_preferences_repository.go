// Package postgres implementa os repositórios usando PostgreSQL e sqlc.
package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
)

// UserPreferencesRepository implementa port.UserPreferencesRepository usando sqlc.
type UserPreferencesRepository struct {
	queries *db.Queries
}

// NewUserPreferencesRepository cria uma nova instância do repositório.
func NewUserPreferencesRepository(queries *db.Queries) *UserPreferencesRepository {
	return &UserPreferencesRepository{
		queries: queries,
	}
}

// Create persiste novas preferências de usuário.
func (r *UserPreferencesRepository) Create(ctx context.Context, prefs *entity.UserPreferences) error {
	userUUID := uuidStringToPgtype(prefs.UserID)

	params := db.CreateUserPreferencesParams{
		UserID:               userUUID,
		AnalyticsEnabled:     prefs.AnalyticsConsent,
		ErrorTrackingEnabled: prefs.ThirdPartyConsent,
		MarketingEnabled:     prefs.MarketingConsent,
	}

	result, err := r.queries.CreateUserPreferences(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar preferências: %w", err)
	}

	prefs.ID = pgUUIDToString(result.ID)
	prefs.CriadoEm = timestamptzToTime(result.CreatedAt)
	prefs.AtualizadoEm = timestamptzToTime(result.UpdatedAt)

	return nil
}

// FindByUserID busca preferências de um usuário.
func (r *UserPreferencesRepository) FindByUserID(ctx context.Context, userID string) (*entity.UserPreferences, error) {
	userUUID := uuidStringToPgtype(userID)

	result, err := r.queries.GetUserPreferencesByUserID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar preferências: %w", err)
	}

	return r.toDomain(&result)
}

// Update atualiza preferências existentes.
func (r *UserPreferencesRepository) Update(ctx context.Context, prefs *entity.UserPreferences) error {
	userUUID := uuidStringToPgtype(prefs.UserID)

	params := db.UpdateUserPreferencesParams{
		UserID:               userUUID,
		AnalyticsEnabled:     prefs.AnalyticsConsent,
		ErrorTrackingEnabled: prefs.ThirdPartyConsent,
		MarketingEnabled:     prefs.MarketingConsent,
	}

	result, err := r.queries.UpdateUserPreferences(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar preferências: %w", err)
	}

	prefs.AtualizadoEm = timestamptzToTime(result.UpdatedAt)
	return nil
}

// Delete remove preferências de um usuário.
func (r *UserPreferencesRepository) Delete(ctx context.Context, userID string) error {
	userUUID := uuidStringToPgtype(userID)

	if err := r.queries.DeleteUserPreferences(ctx, userUUID); err != nil {
		return fmt.Errorf("erro ao deletar preferências: %w", err)
	}

	return nil
}

// toDomain converte modelo sqlc para entidade de domínio.
func (r *UserPreferencesRepository) toDomain(model *db.UserPreference) (*entity.UserPreferences, error) {
	prefs := &entity.UserPreferences{
		ID:                     pgUUIDToString(model.ID),
		UserID:                 pgUUIDToString(model.UserID),
		DataSharingConsent:     false, // Not in DB model
		MarketingConsent:       model.MarketingEnabled,
		AnalyticsConsent:       model.AnalyticsEnabled,
		ThirdPartyConsent:      model.ErrorTrackingEnabled,
		PersonalizedAdsConsent: false, // Not in DB model
		CriadoEm:               timestamptzToTime(model.CreatedAt),
		AtualizadoEm:           timestamptzToTime(model.UpdatedAt),
	}

	return prefs, nil
}
