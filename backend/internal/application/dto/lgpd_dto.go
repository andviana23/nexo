package dto

import "time"

// GetUserPreferencesResponse representa as preferências LGPD do usuário
type GetUserPreferencesResponse struct {
	DataSharingConsent     bool      `json:"data_sharing_consent"`
	MarketingConsent       bool      `json:"marketing_consent"`
	AnalyticsConsent       bool      `json:"analytics_consent"`
	ThirdPartyConsent      bool      `json:"third_party_consent"`
	PersonalizedAdsConsent bool      `json:"personalized_ads_consent"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// UpdateUserPreferencesRequest representa a atualização de consentimentos
type UpdateUserPreferencesRequest struct {
	DataSharingConsent     bool `json:"data_sharing_consent"`
	MarketingConsent       bool `json:"marketing_consent"`
	AnalyticsConsent       bool `json:"analytics_consent"`
	ThirdPartyConsent      bool `json:"third_party_consent"`
	PersonalizedAdsConsent bool `json:"personalized_ads_consent"`
}

// ExportUserDataResponse contém todos os dados do usuário (portabilidade LGPD)
type ExportUserDataResponse struct {
	User         UserDataExport        `json:"user"`
	Tenant       TenantDataExport      `json:"tenant,omitempty"`
	Preferences  PreferencesDataExport `json:"preferences"`
	AuditLogs    []AuditLogDataExport  `json:"audit_logs,omitempty"`
	ExportedAt   time.Time             `json:"exported_at"`
	ExportFormat string                `json:"export_format"`
}

type UserDataExport struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Nome      string     `json:"nome"`
	Role      string     `json:"role"`
	Ativo     bool       `json:"ativo"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type TenantDataExport struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
	CNPJ string `json:"cnpj,omitempty"`
}

type PreferencesDataExport struct {
	DataSharingConsent     bool      `json:"data_sharing_consent"`
	MarketingConsent       bool      `json:"marketing_consent"`
	AnalyticsConsent       bool      `json:"analytics_consent"`
	ThirdPartyConsent      bool      `json:"third_party_consent"`
	PersonalizedAdsConsent bool      `json:"personalized_ads_consent"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type AuditLogDataExport struct {
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Result    string    `json:"result"`
	IPAddress string    `json:"ip_address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// DeleteAccountRequest representa a solicitação de exclusão de conta
type DeleteAccountRequest struct {
	Password string `json:"password" validate:"required,min=6"`
	Reason   string `json:"reason,omitempty"`
}

// DeleteAccountResponse confirma a exclusão
type DeleteAccountResponse struct {
	Message   string    `json:"message"`
	DeletedAt time.Time `json:"deleted_at"`
}
