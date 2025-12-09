package valueobject

// CommissionType representa o tipo de comissão
type CommissionType string

const (
	CommissionTypePercentual CommissionType = "PERCENTUAL"
	CommissionTypeFixo       CommissionType = "FIXO"
)

// IsValid verifica se o tipo é válido
func (t CommissionType) IsValid() bool {
	return t == CommissionTypePercentual || t == CommissionTypeFixo
}

// String retorna a string do tipo
func (t CommissionType) String() string {
	return string(t)
}

// CalculationBase representa a base de cálculo da comissão
type CalculationBase string

const (
	CalculationBaseBruto   CalculationBase = "BRUTO"
	CalculationBaseLiquido CalculationBase = "LIQUIDO"
)

// IsValid verifica se a base é válida
func (b CalculationBase) IsValid() bool {
	return b == CalculationBaseBruto || b == CalculationBaseLiquido
}

// String retorna a string da base
func (b CalculationBase) String() string {
	return string(b)
}

// CommissionSource representa a origem da taxa de comissão aplicada
type CommissionSource string

const (
	CommissionSourceServico      CommissionSource = "SERVICO"
	CommissionSourceProfissional CommissionSource = "PROFISSIONAL"
	CommissionSourceRegra        CommissionSource = "REGRA"
	CommissionSourceManual       CommissionSource = "MANUAL"
)

// IsValid verifica se a origem é válida
func (s CommissionSource) IsValid() bool {
	switch s {
	case CommissionSourceServico, CommissionSourceProfissional, CommissionSourceRegra, CommissionSourceManual:
		return true
	}
	return false
}

// String retorna a string da origem
func (s CommissionSource) String() string {
	return string(s)
}

// CommissionItemStatus representa o status de um item de comissão
type CommissionItemStatus string

const (
	CommissionItemStatusPendente   CommissionItemStatus = "PENDENTE"
	CommissionItemStatusProcessado CommissionItemStatus = "PROCESSADO"
	CommissionItemStatusPago       CommissionItemStatus = "PAGO"
	CommissionItemStatusCancelado  CommissionItemStatus = "CANCELADO"
	CommissionItemStatusEstornado  CommissionItemStatus = "ESTORNADO"
)

// IsValid verifica se o status é válido
func (s CommissionItemStatus) IsValid() bool {
	switch s {
	case CommissionItemStatusPendente, CommissionItemStatusProcessado,
		CommissionItemStatusPago, CommissionItemStatusCancelado, CommissionItemStatusEstornado:
		return true
	}
	return false
}

// String retorna a string do status
func (s CommissionItemStatus) String() string {
	return string(s)
}

// CommissionPeriodStatus representa o status de um período de comissão
type CommissionPeriodStatus string

const (
	CommissionPeriodStatusAberto      CommissionPeriodStatus = "ABERTO"
	CommissionPeriodStatusProcessando CommissionPeriodStatus = "PROCESSANDO"
	CommissionPeriodStatusFechado     CommissionPeriodStatus = "FECHADO"
	CommissionPeriodStatusPago        CommissionPeriodStatus = "PAGO"
	CommissionPeriodStatusCancelado   CommissionPeriodStatus = "CANCELADO"
)

// IsValid verifica se o status é válido
func (s CommissionPeriodStatus) IsValid() bool {
	switch s {
	case CommissionPeriodStatusAberto, CommissionPeriodStatusProcessando,
		CommissionPeriodStatusFechado, CommissionPeriodStatusPago, CommissionPeriodStatusCancelado:
		return true
	}
	return false
}

// String retorna a string do status
func (s CommissionPeriodStatus) String() string {
	return string(s)
}

// CanClose verifica se o período pode ser fechado
func (s CommissionPeriodStatus) CanClose() bool {
	return s == CommissionPeriodStatusAberto
}

// CanPay verifica se o período pode ser pago
func (s CommissionPeriodStatus) CanPay() bool {
	return s == CommissionPeriodStatusFechado
}

// AdvanceStatus representa o status de um adiantamento
type AdvanceStatus string

const (
	AdvanceStatusPending   AdvanceStatus = "PENDING"
	AdvanceStatusApproved  AdvanceStatus = "APPROVED"
	AdvanceStatusRejected  AdvanceStatus = "REJECTED"
	AdvanceStatusDeducted  AdvanceStatus = "DEDUCTED"
	AdvanceStatusCancelled AdvanceStatus = "CANCELLED"
)

// IsValid verifica se o status é válido
func (s AdvanceStatus) IsValid() bool {
	switch s {
	case AdvanceStatusPending, AdvanceStatusApproved, AdvanceStatusRejected,
		AdvanceStatusDeducted, AdvanceStatusCancelled:
		return true
	}
	return false
}

// String retorna a string do status
func (s AdvanceStatus) String() string {
	return string(s)
}

// CanApprove verifica se o adiantamento pode ser aprovado
func (s AdvanceStatus) CanApprove() bool {
	return s == AdvanceStatusPending
}

// CanReject verifica se o adiantamento pode ser rejeitado
func (s AdvanceStatus) CanReject() bool {
	return s == AdvanceStatusPending
}

// CanDeduct verifica se o adiantamento pode ser deduzido
func (s AdvanceStatus) CanDeduct() bool {
	return s == AdvanceStatusApproved
}

// CanCancel verifica se o adiantamento pode ser cancelado
func (s AdvanceStatus) CanCancel() bool {
	return s == AdvanceStatusPending || s == AdvanceStatusApproved
}
