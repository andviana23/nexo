package entity

import "github.com/shopspring/decimal"

// CommissionPeriodSummary representa o resumo de totais de um período de comissão
type CommissionPeriodSummary struct {
	TotalGross       decimal.Decimal
	TotalCommission  decimal.Decimal
	TotalAdvances    decimal.Decimal
	TotalAdjustments decimal.Decimal
	TotalNet         decimal.Decimal
	ItemsCount       int
}

// CommissionSummary representa o resumo de comissões por profissional
type CommissionSummary struct {
	ProfessionalID   string
	ProfessionalName string
	TotalGross       decimal.Decimal
	TotalCommission  decimal.Decimal
	ItemsCount       int
}

// CommissionByService representa o resumo de comissões por serviço
type CommissionByService struct {
	ServiceID       string
	ServiceName     string
	TotalGross      decimal.Decimal
	TotalCommission decimal.Decimal
	ItemsCount      int
}
