package dto

import "github.com/shopspring/decimal"

// FinancialDashboardRequest representa a requisição do dashboard
type FinancialDashboardRequest struct {
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
	Month     string `json:"month" validate:"required"` // formato: "2025-11"
}

// DashboardMetricsResponse representa as métricas do dashboard
type DashboardMetricsResponse struct {
	BreakEvenPoint string  `json:"break_even_point"`
	BurnRate       string  `json:"burn_rate"`
	Runway         float64 `json:"runway"`
	AverageTicket  string  `json:"average_ticket"`
}

// DashboardSummaryResponse contém o resumo financeiro
type DashboardSummaryResponse struct {
	// Payables
	PayablesTotal     string `json:"payables_total"`
	PayablesPaid      string `json:"payables_paid"`
	PayablesPending   string `json:"payables_pending"`
	PayablesOverdue   string `json:"payables_overdue"`
	PayablesThisMonth string `json:"payables_this_month"`

	// Receivables
	ReceivablesTotal     string `json:"receivables_total"`
	ReceivablesReceived  string `json:"receivables_received"`
	ReceivablesPending   string `json:"receivables_pending"`
	ReceivablesOverdue   string `json:"receivables_overdue"`
	ReceivablesThisMonth string `json:"receivables_this_month"`

	// Cashflow
	CurrentBalance   string `json:"current_balance"`
	ProjectedBalance string `json:"projected_balance"`
	TotalIncome      string `json:"total_income"`
	TotalExpense     string `json:"total_expense"`
	NetCashflow      string `json:"net_cashflow"`

	// DRE
	ReceitaBruta     string `json:"receita_bruta"`
	ReceitaLiquida   string `json:"receita_liquida"`
	CustoTotal       string `json:"custo_total"`
	LucroOperacional string `json:"lucro_operacional"`
	MargemLucro      string `json:"margem_lucro"`
	EBITDA           string `json:"ebitda"`

	// Metrics
	Metrics DashboardMetricsResponse `json:"metrics"`
}

// FinancialDashboardResponse representa a resposta completa do dashboard
type FinancialDashboardResponse struct {
	Summary   DashboardSummaryResponse `json:"summary"`
	StartDate string                   `json:"start_date"`
	EndDate   string                   `json:"end_date"`
	Month     string                   `json:"month"`
}

// Helper para converter decimal para string
func DecimalToString(d decimal.Decimal) string {
	return d.StringFixed(2)
}
