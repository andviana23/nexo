package dto

// CreatePlanRequest representa o payload de criação de plano
type CreatePlanRequest struct {
	Nome            string  `json:"nome" validate:"required,min=3,max=100"`
	Descricao       *string `json:"descricao,omitempty" validate:"omitempty,max=500"`
	Valor           string  `json:"valor" validate:"required"`
	QtdServicos     *int    `json:"qtd_servicos,omitempty" validate:"omitempty,gte=0"`
	LimiteUsoMensal *int    `json:"limite_uso_mensal,omitempty" validate:"omitempty,gte=0"`
}

// UpdatePlanRequest representa o payload de atualização de plano
type UpdatePlanRequest struct {
	Nome            string  `json:"nome" validate:"required,min=3,max=100"`
	Descricao       *string `json:"descricao,omitempty" validate:"omitempty,max=500"`
	Valor           string  `json:"valor" validate:"required"`
	QtdServicos     *int    `json:"qtd_servicos,omitempty" validate:"omitempty,gte=0"`
	LimiteUsoMensal *int    `json:"limite_uso_mensal,omitempty" validate:"omitempty,gte=0"`
	Ativo           bool    `json:"ativo"`
}

// PlanResponse representa a resposta de um plano
type PlanResponse struct {
	ID              string  `json:"id"`
	Nome            string  `json:"nome"`
	Descricao       *string `json:"descricao,omitempty"`
	Valor           string  `json:"valor"`
	Periodicidade   string  `json:"periodicidade"`
	QtdServicos     *int    `json:"qtd_servicos,omitempty"`
	LimiteUsoMensal *int    `json:"limite_uso_mensal,omitempty"`
	Ativo           bool    `json:"ativo"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}
