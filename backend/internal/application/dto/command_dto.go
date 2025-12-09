package dto

import (
	"time"
)

// ============================================================================
// REQUEST DTOs
// ============================================================================

// CreateCommandRequest representa a requisição para criar uma comanda
type CreateCommandRequest struct {
	AppointmentID *string            `json:"appointment_id,omitempty"`
	CustomerID    string             `json:"customer_id" validate:"required,uuid"`
	Items         []CommandItemInput `json:"items" validate:"required,min=1,dive"`
	Observacoes   *string            `json:"observacoes,omitempty"`
}

// CommandItemInput representa um item a ser adicionado à comanda
type CommandItemInput struct {
	Tipo               string   `json:"tipo" validate:"required,oneof=SERVICO PRODUTO PACOTE"`
	ItemID             string   `json:"item_id" validate:"required,uuid"`
	Descricao          string   `json:"descricao" validate:"required"`
	PrecoUnitario      string   `json:"preco_unitario" validate:"required"` // String para evitar problemas com float
	Quantidade         int      `json:"quantidade" validate:"required,min=1"`
	DescontoValor      *string  `json:"desconto_valor,omitempty"`
	DescontoPercentual *float64 `json:"desconto_percentual,omitempty" validate:"omitempty,min=0,max=100"`
	Observacoes        *string  `json:"observacoes,omitempty"`
}

// AddCommandItemRequest representa a requisição para adicionar um item
type AddCommandItemRequest struct {
	Tipo               string   `json:"tipo" validate:"required,oneof=SERVICO PRODUTO PACOTE"`
	ItemID             string   `json:"item_id" validate:"required,uuid"`
	Descricao          string   `json:"descricao" validate:"required"`
	PrecoUnitario      string   `json:"preco_unitario" validate:"required"`
	Quantidade         int      `json:"quantidade" validate:"required,min=1"`
	DescontoValor      *string  `json:"desconto_valor,omitempty"`
	DescontoPercentual *float64 `json:"desconto_percentual,omitempty" validate:"omitempty,min=0,max=100"`
	Observacoes        *string  `json:"observacoes,omitempty"`
}

// UpdateCommandItemRequest representa a requisição para atualizar um item
type UpdateCommandItemRequest struct {
	PrecoUnitario      *string  `json:"preco_unitario,omitempty"`
	Quantidade         *int     `json:"quantidade,omitempty" validate:"omitempty,min=1"`
	DescontoValor      *string  `json:"desconto_valor,omitempty"`
	DescontoPercentual *float64 `json:"desconto_percentual,omitempty" validate:"omitempty,min=0,max=100"`
	Observacoes        *string  `json:"observacoes,omitempty"`
}

// AddCommandPaymentRequest representa a requisição para adicionar pagamento
type AddCommandPaymentRequest struct {
	MeioPagamentoID string  `json:"meio_pagamento_id" validate:"required,uuid"`
	ValorRecebido   string  `json:"valor_recebido" validate:"required"`
	Observacoes     *string `json:"observacoes,omitempty"`
}

// CloseCommandRequest representa a requisição para fechar a comanda
type CloseCommandRequest struct {
	DeixarTrocoGorjeta *bool   `json:"deixar_troco_gorjeta,omitempty"`
	DeixarSaldoDivida  *bool   `json:"deixar_saldo_divida,omitempty"`
	Observacoes        *string `json:"observacoes,omitempty"`
}

// ============================================================================
// RESPONSE DTOs
// ============================================================================

// CommandResponse representa uma comanda completa
type CommandResponse struct {
	ID                 string                   `json:"id"`
	AppointmentID      *string                  `json:"appointment_id,omitempty"`
	CustomerID         string                   `json:"customer_id"`
	Numero             *string                  `json:"numero,omitempty"`
	Status             string                   `json:"status"`
	Subtotal           string                   `json:"subtotal"`
	Desconto           string                   `json:"desconto"`
	Total              string                   `json:"total"`
	TotalRecebido      string                   `json:"total_recebido"`
	Troco              string                   `json:"troco"`
	SaldoDevedor       string                   `json:"saldo_devedor"`
	Observacoes        *string                  `json:"observacoes,omitempty"`
	DeixarTrocoGorjeta bool                     `json:"deixar_troco_gorjeta"`
	DeixarSaldoDivida  bool                     `json:"deixar_saldo_divida"`
	CriadoEm           time.Time                `json:"criado_em"`
	AtualizadoEm       time.Time                `json:"atualizado_em"`
	FechadoEm          *time.Time               `json:"fechado_em,omitempty"`
	FechadoPor         *string                  `json:"fechado_por,omitempty"`
	Items              []CommandItemResponse    `json:"items"`
	Payments           []CommandPaymentResponse `json:"payments"`
}

// CommandItemResponse representa um item da comanda
type CommandItemResponse struct {
	ID                 string    `json:"id"`
	CommandID          string    `json:"command_id"`
	Tipo               string    `json:"tipo"`
	ItemID             string    `json:"item_id"`
	Descricao          string    `json:"descricao"`
	PrecoUnitario      string    `json:"preco_unitario"`
	Quantidade         int       `json:"quantidade"`
	DescontoValor      string    `json:"desconto_valor"`
	DescontoPercentual string    `json:"desconto_percentual"`
	PrecoFinal         string    `json:"preco_final"`
	Observacoes        *string   `json:"observacoes,omitempty"`
	CriadoEm           time.Time `json:"criado_em"`
}

// CommandPaymentResponse representa um pagamento da comanda
type CommandPaymentResponse struct {
	ID              string    `json:"id"`
	CommandID       string    `json:"command_id"`
	MeioPagamentoID string    `json:"meio_pagamento_id"`
	ValorRecebido   string    `json:"valor_recebido"`
	TaxaPercentual  string    `json:"taxa_percentual"`
	TaxaFixa        string    `json:"taxa_fixa"`
	ValorLiquido    string    `json:"valor_liquido"`
	Observacoes     *string   `json:"observacoes,omitempty"`
	CriadoEm        time.Time `json:"criado_em"`
	CriadoPor       *string   `json:"criado_por,omitempty"`
}

// CommandListResponse representa a listagem de comandas
type CommandListResponse struct {
	Commands   []CommandResponse `json:"commands"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
