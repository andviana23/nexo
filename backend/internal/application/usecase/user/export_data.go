package user

import (
	"context"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"go.uber.org/zap"
)

// ExportDataResponse contém todos os dados do usuário para portabilidade (LGPD Art. 18, V)
type ExportDataResponse struct {
	User         UserData        `json:"user"`
	Tenant       TenantData      `json:"tenant,omitempty"`
	Preferences  PreferencesData `json:"preferences"`
	AuditLogs    []AuditLogData  `json:"audit_logs,omitempty"`
	ExportedAt   time.Time       `json:"exported_at"`
	ExportFormat string          `json:"export_format"`
}

type UserData struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Nome      string     `json:"nome"`
	Role      string     `json:"role"`
	Ativo     bool       `json:"ativo"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type TenantData struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
	CNPJ string `json:"cnpj,omitempty"`
}

type PreferencesData struct {
	DataSharingConsent     bool      `json:"data_sharing_consent"`
	MarketingConsent       bool      `json:"marketing_consent"`
	AnalyticsConsent       bool      `json:"analytics_consent"`
	ThirdPartyConsent      bool      `json:"third_party_consent"`
	PersonalizedAdsConsent bool      `json:"personalized_ads_consent"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type AuditLogData struct {
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Result    string    `json:"result"`
	IPAddress string    `json:"ip_address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ExportDataUseCase implementa o direito de portabilidade (LGPD Art. 18, V)
type ExportDataUseCase struct {
	prefsRepo port.UserPreferencesRepository
	logger    *zap.Logger
}

func NewExportDataUseCase(
	prefsRepo port.UserPreferencesRepository,
	logger *zap.Logger,
) *ExportDataUseCase {
	return &ExportDataUseCase{
		prefsRepo: prefsRepo,
		logger:    logger,
	}
}

// Execute exporta todos os dados pessoais do usuário
func (uc *ExportDataUseCase) Execute(ctx context.Context, tenantID, userID string) (*ExportDataResponse, error) {
	if tenantID == "" {
		return nil, domain.ErrTenantIDRequired
	}
	if userID == "" {
		return nil, domain.ErrInvalidID
	}

	uc.logger.Info("Exportando dados do usuário (LGPD)",
		zap.String("tenant_id", tenantID),
		zap.String("user_id", userID),
	)

	// Buscar preferências
	prefs, err := uc.prefsRepo.FindByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("Erro ao buscar preferências para export",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		// Continuar sem preferências se não existir
		prefs = nil
	}

	// Montar resposta
	response := &ExportDataResponse{
		User: UserData{
			ID:    userID,
			Email: "[redacted - query users table]", // TODO: integrar com user repository
			Nome:  "[redacted - query users table]",
		},
		Preferences: PreferencesData{
			DataSharingConsent:     false,
			MarketingConsent:       false,
			AnalyticsConsent:       false,
			ThirdPartyConsent:      false,
			PersonalizedAdsConsent: false,
		},
		ExportedAt:   time.Now(),
		ExportFormat: "JSON",
	}

	if prefs != nil {
		response.Preferences = PreferencesData{
			DataSharingConsent:     prefs.DataSharingConsent,
			MarketingConsent:       prefs.MarketingConsent,
			AnalyticsConsent:       prefs.AnalyticsConsent,
			ThirdPartyConsent:      prefs.ThirdPartyConsent,
			PersonalizedAdsConsent: prefs.PersonalizedAdsConsent,
			UpdatedAt:              prefs.AtualizadoEm,
		}
	}

	// TODO: Buscar dados do usuário (users table)
	// TODO: Buscar dados do tenant (tenants table)
	// TODO: Buscar audit logs (últimos 90 dias)

	uc.logger.Info("Exportação de dados concluída",
		zap.String("user_id", userID),
		zap.Time("exported_at", response.ExportedAt),
	)

	return response, nil
}
