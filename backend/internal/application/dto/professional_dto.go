// Package dto contém DTOs para profissionais
package dto

import "time"

// =============================================================================
// Request DTOs
// =============================================================================

// CreateProfessionalRequest representa os dados para criar um profissional
type CreateProfessionalRequest struct {
	Nome            string   `json:"nome" validate:"required,min=2,max=100"`
	Email           string   `json:"email" validate:"required,email"`
	Telefone        string   `json:"telefone" validate:"required,min=10,max=15"`
	CPF             string   `json:"cpf" validate:"required,min=11,max=14"` // CPF (11) ou CNPJ (14)
	Especialidades  []string `json:"especialidades" validate:"omitempty"`
	Comissao        string   `json:"comissao" validate:"omitempty"` // String para precisão decimal
	TipoComissao    string   `json:"tipo_comissao" validate:"omitempty,oneof=PERCENTUAL FIXO"`
	Foto            *string  `json:"foto,omitempty"`
	DataAdmissao    string   `json:"data_admissao" validate:"required"`
	Status          string   `json:"status" validate:"omitempty,oneof=ATIVO INATIVO FERIAS LICENCA DEMITIDO"`
	HorarioTrabalho *string  `json:"horario_trabalho,omitempty"` // JSON string
	Observacoes     *string  `json:"observacoes,omitempty"`
	Tipo            string   `json:"tipo" validate:"required,oneof=BARBEIRO GERENTE RECEPCIONISTA OUTRO"`
}

// UpdateProfessionalRequest representa os dados para atualizar um profissional
type UpdateProfessionalRequest struct {
	Nome            string   `json:"nome" validate:"required,min=2,max=100"`
	Email           string   `json:"email" validate:"required,email"`
	Telefone        string   `json:"telefone" validate:"required,min=10,max=15"`
	CPF             string   `json:"cpf" validate:"required,min=11,max=14"` // CPF (11) ou CNPJ (14)
	Especialidades  []string `json:"especialidades" validate:"omitempty"`
	Comissao        string   `json:"comissao" validate:"omitempty"`
	TipoComissao    string   `json:"tipo_comissao" validate:"omitempty,oneof=PERCENTUAL FIXO"`
	Foto            *string  `json:"foto,omitempty"`
	DataAdmissao    string   `json:"data_admissao" validate:"required"`
	DataDemissao    *string  `json:"data_demissao,omitempty"`
	Status          string   `json:"status" validate:"omitempty,oneof=ATIVO INATIVO FERIAS LICENCA DEMITIDO"`
	HorarioTrabalho *string  `json:"horario_trabalho,omitempty"`
	Observacoes     *string  `json:"observacoes,omitempty"`
	Tipo            string   `json:"tipo" validate:"required,oneof=BARBEIRO GERENTE RECEPCIONISTA OUTRO"`
}

// UpdateProfessionalStatusRequest representa os dados para atualizar status
type UpdateProfessionalStatusRequest struct {
	Status       string  `json:"status" validate:"required,oneof=ATIVO INATIVO FERIAS LICENCA DEMITIDO"`
	DataDemissao *string `json:"data_demissao,omitempty"`
}

// ListProfessionalsRequest query params para listagem
type ListProfessionalsRequest struct {
	Search         string `query:"search" validate:"omitempty,max=100"`
	Status         string `query:"status" validate:"omitempty,oneof=ATIVO INATIVO FERIAS LICENCA DEMITIDO"`
	Tipo           string `query:"tipo" validate:"omitempty,oneof=BARBEIRO GERENTE RECEPCIONISTA OUTRO"`
	OrderBy        string `query:"order_by" validate:"omitempty,oneof=nome criado_em data_admissao"`
	OrderDirection string `query:"order_direction" validate:"omitempty,oneof=asc desc"`
	Page           int    `query:"page" validate:"omitempty,min=1"`
	PageSize       int    `query:"page_size" validate:"omitempty,min=1,max=100"`
}

// CheckEmailRequest query params para verificar email
type CheckEmailProfessionalRequest struct {
	Email     string  `query:"email" validate:"required,email"`
	ExcludeID *string `query:"exclude_id" validate:"omitempty,uuid"`
}

// CheckCpfRequest query params para verificar cpf
type CheckCpfProfessionalRequest struct {
	CPF       string  `query:"cpf" validate:"required,len=11"`
	ExcludeID *string `query:"exclude_id" validate:"omitempty,uuid"`
}

// =============================================================================
// Response DTOs
// =============================================================================

// ProfessionalResponse resposta de profissional
type ProfessionalResponse struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	UserID          *string   `json:"user_id,omitempty"`
	Nome            string    `json:"nome"`
	Email           string    `json:"email"`
	Telefone        string    `json:"telefone"`
	CPF             string    `json:"cpf"`
	Especialidades  []string  `json:"especialidades"`
	Comissao        string    `json:"comissao"`
	TipoComissao    string    `json:"tipo_comissao"`
	Foto            *string   `json:"foto,omitempty"`
	DataAdmissao    string    `json:"data_admissao"`
	DataDemissao    *string   `json:"data_demissao,omitempty"`
	Status          string    `json:"status"`
	HorarioTrabalho *string   `json:"horario_trabalho,omitempty"`
	Observacoes     *string   `json:"observacoes,omitempty"`
	Tipo            string    `json:"tipo"`
	CriadoEm        time.Time `json:"criado_em"`
	AtualizadoEm    time.Time `json:"atualizado_em"`
}

// ListProfessionalsResponse resposta de listagem de profissionais
type ListProfessionalsResponse struct {
	Data     []ProfessionalResponse `json:"data"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Total    int64                  `json:"total"`
}

// CheckExistsResponse resposta de verificação de existência
type CheckExistsProfessionalResponse struct {
	Exists bool `json:"exists"`
}
