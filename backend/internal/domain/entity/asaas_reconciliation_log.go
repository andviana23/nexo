package entity

import (
	"time"

	"github.com/google/uuid"
)

// AsaasReconciliationLog registra execuções de conciliação Asaas x NEXO
// Alinhado com schema sqlc - tabela asaas_reconciliation_logs
type AsaasReconciliationLog struct {
	ID            string
	TenantID      uuid.UUID
	PeriodStart   time.Time
	PeriodEnd     time.Time
	TotalAsaas    int
	TotalNexo     int
	Divergences   int
	AutoFixed     int
	PendingReview int
	Details       *string // JSON com detalhes das divergências
	CreatedAt     time.Time
}

// NewAsaasReconciliationLog cria um novo log de conciliação
func NewAsaasReconciliationLog(tenantID uuid.UUID, periodStart, periodEnd time.Time) *AsaasReconciliationLog {
	now := time.Now()
	return &AsaasReconciliationLog{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		CreatedAt:   now,
	}
}

// SetResults define os resultados da conciliação
func (r *AsaasReconciliationLog) SetResults(totalAsaas, totalNexo, divergences, autoFixed, pendingReview int, details *string) {
	r.TotalAsaas = totalAsaas
	r.TotalNexo = totalNexo
	r.Divergences = divergences
	r.AutoFixed = autoFixed
	r.PendingReview = pendingReview
	r.Details = details
}

// HasDivergences retorna se há divergências
func (r *AsaasReconciliationLog) HasDivergences() bool {
	return r.Divergences > 0 || r.PendingReview > 0
}
