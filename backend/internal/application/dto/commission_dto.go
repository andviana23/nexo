package dto

import "time"

// =============================================================================
// DTOs para Comissões
// Conforme módulo de comissões do NEXO
// =============================================================================

// =============================================================================
// Commission Rule DTOs
// =============================================================================

// CreateCommissionRuleRequest requisição para criar regra de comissão
type CreateCommissionRuleRequest struct {
	UnitID          *string `json:"unit_id,omitempty" validate:"omitempty,uuid"`
	Name            string  `json:"name" validate:"required,min=2,max=255"`
	Description     *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Type            string  `json:"type" validate:"required,oneof=PERCENTUAL FIXO"`
	DefaultRate     string  `json:"default_rate" validate:"required"` // Dinheiro/percentual como string
	MinAmount       *string `json:"min_amount,omitempty"`
	MaxAmount       *string `json:"max_amount,omitempty"`
	CalculationBase *string `json:"calculation_base,omitempty" validate:"omitempty,oneof=BRUTO LIQUIDO"`
	EffectiveFrom   *string `json:"effective_from,omitempty"` // RFC3339
	EffectiveTo     *string `json:"effective_to,omitempty"`   // RFC3339
	Priority        *int    `json:"priority,omitempty"`
}

// UpdateCommissionRuleRequest requisição para atualizar regra de comissão
type UpdateCommissionRuleRequest struct {
	Name            *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Description     *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Type            *string `json:"type,omitempty" validate:"omitempty,oneof=PERCENTUAL FIXO"`
	DefaultRate     *string `json:"default_rate,omitempty"`
	MinAmount       *string `json:"min_amount,omitempty"`
	MaxAmount       *string `json:"max_amount,omitempty"`
	CalculationBase *string `json:"calculation_base,omitempty" validate:"omitempty,oneof=BRUTO LIQUIDO"`
	EffectiveFrom   *string `json:"effective_from,omitempty"`
	EffectiveTo     *string `json:"effective_to,omitempty"`
	Priority        *int    `json:"priority,omitempty"`
	IsActive        *bool   `json:"is_active,omitempty"`
}

// ListCommissionRulesRequest query params para listagem de regras
type ListCommissionRulesRequest struct {
	ActiveOnly bool `query:"active_only"`
}

// CommissionRuleResponse resposta de regra de comissão
type CommissionRuleResponse struct {
	ID              string  `json:"id"`
	TenantID        string  `json:"tenant_id"`
	UnitID          *string `json:"unit_id,omitempty"`
	Name            string  `json:"name"`
	Description     *string `json:"description,omitempty"`
	Type            string  `json:"type"`
	DefaultRate     string  `json:"default_rate"`
	MinAmount       *string `json:"min_amount,omitempty"`
	MaxAmount       *string `json:"max_amount,omitempty"`
	CalculationBase *string `json:"calculation_base,omitempty"`
	EffectiveFrom   string  `json:"effective_from"`
	EffectiveTo     *string `json:"effective_to,omitempty"`
	Priority        *int    `json:"priority,omitempty"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// ListCommissionRulesResponse resposta de listagem de regras
type ListCommissionRulesResponse struct {
	Data  []CommissionRuleResponse `json:"data"`
	Total int                      `json:"total"`
}

// =============================================================================
// Commission Period DTOs
// =============================================================================

// CreateCommissionPeriodRequest requisição para criar período de comissão
type CreateCommissionPeriodRequest struct {
	UnitID         *string `json:"unit_id,omitempty" validate:"omitempty,uuid"`
	ReferenceMonth string  `json:"reference_month" validate:"required"` // formato: YYYY-MM
	ProfessionalID string  `json:"professional_id" validate:"required,uuid"`
	PeriodStart    string  `json:"period_start" validate:"required"` // RFC3339
	PeriodEnd      string  `json:"period_end" validate:"required"`   // RFC3339
	Notes          *string `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

// CloseCommissionPeriodRequest requisição para fechar período
type CloseCommissionPeriodRequest struct {
	Notes *string `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

// ListCommissionPeriodsRequest query params para listagem de períodos
type ListCommissionPeriodsRequest struct {
	ProfessionalID *string `query:"professional_id" validate:"omitempty,uuid"`
	Status         *string `query:"status" validate:"omitempty,oneof=ABERTO PROCESSANDO FECHADO PAGO CANCELADO"`
	Limit          int     `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset         int     `query:"offset" validate:"omitempty,min=0"`
}

// CommissionPeriodResponse resposta de período de comissão
type CommissionPeriodResponse struct {
	ID               string  `json:"id"`
	TenantID         string  `json:"tenant_id"`
	UnitID           *string `json:"unit_id,omitempty"`
	ReferenceMonth   string  `json:"reference_month"`
	ProfessionalID   *string `json:"professional_id,omitempty"`
	ProfessionalName *string `json:"professional_name,omitempty"`
	TotalGross       string  `json:"total_gross"`
	TotalCommission  string  `json:"total_commission"`
	TotalAdvances    string  `json:"total_advances"`
	TotalAdjustments string  `json:"total_adjustments"`
	TotalNet         string  `json:"total_net"`
	ItemsCount       int     `json:"items_count"`
	Status           string  `json:"status"`
	PeriodStart      string  `json:"period_start"`
	PeriodEnd        string  `json:"period_end"`
	ClosedAt         *string `json:"closed_at,omitempty"`
	PaidAt           *string `json:"paid_at,omitempty"`
	ClosedBy         *string `json:"closed_by,omitempty"`
	PaidBy           *string `json:"paid_by,omitempty"`
	Notes            *string `json:"notes,omitempty"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// ListCommissionPeriodsResponse resposta de listagem de períodos
type ListCommissionPeriodsResponse struct {
	Data  []CommissionPeriodResponse `json:"data"`
	Total int                        `json:"total"`
}

// CommissionPeriodSummaryResponse resumo do período
type CommissionPeriodSummaryResponse struct {
	TotalGross      string `json:"total_gross"`
	TotalCommission string `json:"total_commission"`
	TotalAdvances   string `json:"total_advances"`
	TotalNet        string `json:"total_net"`
	ItemsCount      int    `json:"items_count"`
}

// =============================================================================
// Advance DTOs
// =============================================================================

// CreateAdvanceRequest requisição para criar adiantamento
type CreateAdvanceRequest struct {
	UnitID         *string `json:"unit_id,omitempty" validate:"omitempty,uuid"`
	ProfessionalID string  `json:"professional_id" validate:"required,uuid"`
	Amount         string  `json:"amount" validate:"required"` // Dinheiro como string
	Reason         *string `json:"reason,omitempty" validate:"omitempty,max=500"`
}

// RejectAdvanceRequest requisição para rejeitar adiantamento
type RejectAdvanceRequest struct {
	RejectionReason string `json:"rejection_reason" validate:"required,min=5,max=500"`
}

// MarkAdvanceDeductedRequest requisição para marcar como deduzido
type MarkAdvanceDeductedRequest struct {
	PeriodID string `json:"period_id" validate:"required,uuid"`
}

// ListAdvancesRequest query params para listagem de adiantamentos
type ListAdvancesRequest struct {
	ProfessionalID *string `query:"professional_id" validate:"omitempty,uuid"`
	Status         *string `query:"status" validate:"omitempty,oneof=PENDING APPROVED REJECTED DEDUCTED CANCELLED"`
	Limit          int     `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset         int     `query:"offset" validate:"omitempty,min=0"`
}

// AdvanceResponse resposta de adiantamento
type AdvanceResponse struct {
	ID                string  `json:"id"`
	TenantID          string  `json:"tenant_id"`
	UnitID            *string `json:"unit_id,omitempty"`
	ProfessionalID    string  `json:"professional_id"`
	ProfessionalName  *string `json:"professional_name,omitempty"`
	Amount            string  `json:"amount"`
	RequestDate       string  `json:"request_date"`
	Reason            *string `json:"reason,omitempty"`
	Status            string  `json:"status"`
	ApprovedAt        *string `json:"approved_at,omitempty"`
	ApprovedBy        *string `json:"approved_by,omitempty"`
	RejectedAt        *string `json:"rejected_at,omitempty"`
	RejectedBy        *string `json:"rejected_by,omitempty"`
	RejectionReason   *string `json:"rejection_reason,omitempty"`
	DeductedAt        *string `json:"deducted_at,omitempty"`
	DeductionPeriodID *string `json:"deduction_period_id,omitempty"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// ListAdvancesResponse resposta de listagem de adiantamentos
type ListAdvancesResponse struct {
	Data  []AdvanceResponse `json:"data"`
	Total int               `json:"total"`
}

// AdvancesTotalsResponse totais de adiantamentos
type AdvancesTotalsResponse struct {
	Advances      []AdvanceResponse `json:"advances"`
	TotalPending  string            `json:"total_pending"`
	TotalApproved string            `json:"total_approved"`
}

// =============================================================================
// Commission Item DTOs
// =============================================================================

// CreateCommissionItemRequest requisição para criar item de comissão
type CreateCommissionItemRequest struct {
	UnitID           *string `json:"unit_id,omitempty" validate:"omitempty,uuid"`
	ProfessionalID   string  `json:"professional_id" validate:"required,uuid"`
	CommandID        *string `json:"command_id,omitempty" validate:"omitempty,uuid"`
	CommandItemID    *string `json:"command_item_id,omitempty" validate:"omitempty,uuid"`
	AppointmentID    *string `json:"appointment_id,omitempty" validate:"omitempty,uuid"`
	ServiceID        *string `json:"service_id,omitempty" validate:"omitempty,uuid"`
	ServiceName      *string `json:"service_name,omitempty" validate:"omitempty,max=255"`
	GrossValue       string  `json:"gross_value" validate:"required"` // Dinheiro como string
	CommissionRate   string  `json:"commission_rate" validate:"required"`
	CommissionType   string  `json:"commission_type" validate:"required,oneof=PERCENTUAL FIXO"`
	CommissionSource string  `json:"commission_source" validate:"required,oneof=SERVICO PROFISSIONAL REGRA MANUAL"`
	RuleID           *string `json:"rule_id,omitempty" validate:"omitempty,uuid"`
	ReferenceDate    string  `json:"reference_date" validate:"required"` // RFC3339
	Description      *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// CreateCommissionItemBatchRequest requisição para criar múltiplos itens
type CreateCommissionItemBatchRequest struct {
	Items []CreateCommissionItemRequest `json:"items" validate:"required,min=1,dive"`
}

// ProcessCommissionItemRequest requisição para processar item
type ProcessCommissionItemRequest struct {
	PeriodID string `json:"period_id" validate:"required,uuid"`
}

// AssignItemsToPeriodRequest requisição para vincular itens a período
type AssignItemsToPeriodRequest struct {
	ProfessionalID string `json:"professional_id" validate:"required,uuid"`
	PeriodID       string `json:"period_id" validate:"required,uuid"`
	StartDate      string `json:"start_date" validate:"required"` // RFC3339
	EndDate        string `json:"end_date" validate:"required"`   // RFC3339
}

// ListCommissionItemsRequest query params para listagem de itens
type ListCommissionItemsRequest struct {
	ProfessionalID *string `query:"professional_id" validate:"omitempty,uuid"`
	PeriodID       *string `query:"period_id" validate:"omitempty,uuid"`
	Status         *string `query:"status" validate:"omitempty,oneof=PENDENTE PROCESSADO PAGO CANCELADO ESTORNADO"`
	Limit          int     `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset         int     `query:"offset" validate:"omitempty,min=0"`
}

// CommissionItemResponse resposta de item de comissão
type CommissionItemResponse struct {
	ID               string  `json:"id"`
	TenantID         string  `json:"tenant_id"`
	UnitID           *string `json:"unit_id,omitempty"`
	ProfessionalID   string  `json:"professional_id"`
	ProfessionalName *string `json:"professional_name,omitempty"`
	CommandID        *string `json:"command_id,omitempty"`
	CommandItemID    *string `json:"command_item_id,omitempty"`
	AppointmentID    *string `json:"appointment_id,omitempty"`
	ServiceID        *string `json:"service_id,omitempty"`
	ServiceName      *string `json:"service_name,omitempty"`
	GrossValue       string  `json:"gross_value"`
	CommissionRate   string  `json:"commission_rate"`
	CommissionType   string  `json:"commission_type"`
	CommissionValue  string  `json:"commission_value"`
	CommissionSource string  `json:"commission_source"`
	RuleID           *string `json:"rule_id,omitempty"`
	ReferenceDate    string  `json:"reference_date"`
	Description      *string `json:"description,omitempty"`
	Status           string  `json:"status"`
	PeriodID         *string `json:"period_id,omitempty"`
	ProcessedAt      *string `json:"processed_at,omitempty"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// ListCommissionItemsResponse resposta de listagem de itens
type ListCommissionItemsResponse struct {
	Data  []CommissionItemResponse `json:"data"`
	Total int                      `json:"total"`
}

// CommissionSummaryResponse resumo de comissões por profissional
type CommissionSummaryResponse struct {
	ProfessionalID   string `json:"professional_id"`
	ProfessionalName string `json:"professional_name"`
	TotalGross       string `json:"total_gross"`
	TotalCommission  string `json:"total_commission"`
	ItemsCount       int    `json:"items_count"`
}

// CommissionByServiceResponse resumo de comissões por serviço
type CommissionByServiceResponse struct {
	ServiceID       string `json:"service_id"`
	ServiceName     string `json:"service_name"`
	TotalGross      string `json:"total_gross"`
	TotalCommission string `json:"total_commission"`
	ItemsCount      int    `json:"items_count"`
}

// CommissionSummariesResponse resposta com múltiplos resumos
type CommissionSummariesResponse struct {
	ByProfessional []CommissionSummaryResponse   `json:"by_professional,omitempty"`
	ByService      []CommissionByServiceResponse `json:"by_service,omitempty"`
	StartDate      time.Time                     `json:"start_date"`
	EndDate        time.Time                     `json:"end_date"`
}
