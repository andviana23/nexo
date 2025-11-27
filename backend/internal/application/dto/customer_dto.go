package dto

import "time"

// =============================================================================
// DTOs para Clientes (Customers)
// Conforme FLUXO_CADASTROS_CLIENTE.md
// =============================================================================

// =============================================================================
// Request DTOs
// =============================================================================

// CreateCustomerRequest requisição para criar cliente
type CreateCustomerRequest struct {
	// Campos Obrigatórios
	Nome     string `json:"nome" validate:"required,min=3,max=255"`
	Telefone string `json:"telefone" validate:"required,min=10,max=20"`

	// Campos Opcionais
	Email          *string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	CPF            *string `json:"cpf,omitempty" validate:"omitempty,len=11"`
	DataNascimento *string `json:"data_nascimento,omitempty"` // Format: YYYY-MM-DD
	Genero         *string `json:"genero,omitempty" validate:"omitempty,oneof=M F NB PNI"`

	// Endereço
	EnderecoLogradouro  *string `json:"endereco_logradouro,omitempty" validate:"omitempty,max=255"`
	EnderecoNumero      *string `json:"endereco_numero,omitempty" validate:"omitempty,max=20"`
	EnderecoComplemento *string `json:"endereco_complemento,omitempty" validate:"omitempty,max=100"`
	EnderecoBairro      *string `json:"endereco_bairro,omitempty" validate:"omitempty,max=100"`
	EnderecoCidade      *string `json:"endereco_cidade,omitempty" validate:"omitempty,max=100"`
	EnderecoEstado      *string `json:"endereco_estado,omitempty" validate:"omitempty,len=2"`
	EnderecoCEP         *string `json:"endereco_cep,omitempty" validate:"omitempty,len=8"`

	// CRM
	Observacoes *string  `json:"observacoes,omitempty" validate:"omitempty,max=500"`
	Tags        []string `json:"tags,omitempty" validate:"omitempty,max=10,dive,max=50"`
}

// UpdateCustomerRequest requisição para atualizar cliente
type UpdateCustomerRequest struct {
	Nome           *string `json:"nome,omitempty" validate:"omitempty,min=3,max=255"`
	Telefone       *string `json:"telefone,omitempty" validate:"omitempty,min=10,max=20"`
	Email          *string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	CPF            *string `json:"cpf,omitempty" validate:"omitempty,len=11"`
	DataNascimento *string `json:"data_nascimento,omitempty"`
	Genero         *string `json:"genero,omitempty" validate:"omitempty,oneof=M F NB PNI"`

	// Endereço
	EnderecoLogradouro  *string `json:"endereco_logradouro,omitempty" validate:"omitempty,max=255"`
	EnderecoNumero      *string `json:"endereco_numero,omitempty" validate:"omitempty,max=20"`
	EnderecoComplemento *string `json:"endereco_complemento,omitempty" validate:"omitempty,max=100"`
	EnderecoBairro      *string `json:"endereco_bairro,omitempty" validate:"omitempty,max=100"`
	EnderecoCidade      *string `json:"endereco_cidade,omitempty" validate:"omitempty,max=100"`
	EnderecoEstado      *string `json:"endereco_estado,omitempty" validate:"omitempty,len=2"`
	EnderecoCEP         *string `json:"endereco_cep,omitempty" validate:"omitempty,len=8"`

	// CRM
	Observacoes *string  `json:"observacoes,omitempty" validate:"omitempty,max=500"`
	Tags        []string `json:"tags,omitempty" validate:"omitempty,max=10,dive,max=50"`
}

// UpdateCustomerTagsRequest requisição para atualizar tags
type UpdateCustomerTagsRequest struct {
	Tags []string `json:"tags" validate:"required,max=10,dive,max=50"`
}

// ListCustomersRequest query params para listagem
type ListCustomersRequest struct {
	Search   string   `query:"search" validate:"omitempty,max=100"`
	Ativo    *bool    `query:"ativo"`
	Tags     []string `query:"tags" validate:"omitempty,max=10"`
	OrderBy  string   `query:"order_by" validate:"omitempty,oneof=nome criado_em atualizado_em"`
	Page     int      `query:"page" validate:"omitempty,min=1"`
	PageSize int      `query:"page_size" validate:"omitempty,min=1,max=100"`
}

// SearchCustomersRequest query params para busca rápida
type SearchCustomersRequest struct {
	Query string `query:"q" validate:"required,min=2,max=100"`
}

// =============================================================================
// Response DTOs
// =============================================================================

// CustomerResponse resposta de cliente
type CustomerResponse struct {
	ID             string  `json:"id"`
	TenantID       string  `json:"tenant_id"`
	Nome           string  `json:"nome"`
	Telefone       string  `json:"telefone"`
	Email          *string `json:"email,omitempty"`
	CPF            *string `json:"cpf,omitempty"`
	DataNascimento *string `json:"data_nascimento,omitempty"`
	Genero         *string `json:"genero,omitempty"`

	// Endereço
	EnderecoLogradouro  *string `json:"endereco_logradouro,omitempty"`
	EnderecoNumero      *string `json:"endereco_numero,omitempty"`
	EnderecoComplemento *string `json:"endereco_complemento,omitempty"`
	EnderecoBairro      *string `json:"endereco_bairro,omitempty"`
	EnderecoCidade      *string `json:"endereco_cidade,omitempty"`
	EnderecoEstado      *string `json:"endereco_estado,omitempty"`
	EnderecoCEP         *string `json:"endereco_cep,omitempty"`

	// CRM
	Observacoes *string  `json:"observacoes,omitempty"`
	Tags        []string `json:"tags"`
	Ativo       bool     `json:"ativo"`

	CreatedAt time.Time `json:"criado_em"`
	UpdatedAt time.Time `json:"atualizado_em"`
}

// CustomerSummaryResponse resposta resumida (para selects/listas)
type CustomerSummaryResponse struct {
	ID       string   `json:"id"`
	Nome     string   `json:"nome"`
	Telefone string   `json:"telefone"`
	Email    *string  `json:"email,omitempty"`
	Tags     []string `json:"tags"`
}

// CustomerWithHistoryResponse cliente com histórico
type CustomerWithHistoryResponse struct {
	CustomerResponse
	TotalAtendimentos   int64      `json:"total_atendimentos"`
	TotalGasto          string     `json:"total_gasto"`
	TicketMedio         string     `json:"ticket_medio"`
	UltimoAtendimento   *time.Time `json:"ultimo_atendimento,omitempty"`
	FrequenciaMediaDias *int       `json:"frequencia_media_dias,omitempty"`
}

// ListCustomersResponse resposta de listagem paginada
type ListCustomersResponse struct {
	Data     []CustomerResponse `json:"data"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
	Total    int64              `json:"total"`
}

// CustomerStatsResponse estatísticas de clientes
type CustomerStatsResponse struct {
	TotalAtivos        int64 `json:"total_ativos"`
	TotalInativos      int64 `json:"total_inativos"`
	NovosUltimos30Dias int64 `json:"novos_ultimos_30_dias"`
	TotalGeral         int64 `json:"total_geral"`
}

// =============================================================================
// LGPD Export DTO
// =============================================================================

// CustomerExportResponse dados para exportação LGPD
type CustomerExportResponse struct {
	DadosPessoais struct {
		Nome           string                 `json:"nome"`
		Email          *string                `json:"email,omitempty"`
		Telefone       string                 `json:"telefone"`
		CPF            *string                `json:"cpf,omitempty"`
		DataNascimento *string                `json:"data_nascimento,omitempty"`
		Genero         *string                `json:"genero,omitempty"`
		Endereco       *CustomerAddressExport `json:"endereco,omitempty"`
	} `json:"dados_pessoais"`

	HistoricoAtendimentos []CustomerAppointmentExport `json:"historico_atendimentos"`

	Metricas struct {
		TotalGasto   string `json:"total_gasto"`
		TicketMedio  string `json:"ticket_medio"`
		TotalVisitas int64  `json:"total_visitas"`
	} `json:"metricas"`

	DataExportacao time.Time `json:"data_exportacao"`
}

// CustomerAddressExport endereço para exportação
type CustomerAddressExport struct {
	Logradouro  *string `json:"logradouro,omitempty"`
	Numero      *string `json:"numero,omitempty"`
	Complemento *string `json:"complemento,omitempty"`
	Bairro      *string `json:"bairro,omitempty"`
	Cidade      *string `json:"cidade,omitempty"`
	Estado      *string `json:"estado,omitempty"`
	CEP         *string `json:"cep,omitempty"`
}

// CustomerAppointmentExport atendimento para exportação
type CustomerAppointmentExport struct {
	Data         time.Time `json:"data"`
	Status       string    `json:"status"`
	Profissional string    `json:"profissional"`
	ValorTotal   string    `json:"valor_total"`
}

// =============================================================================
// Validation Check DTOs
// =============================================================================

// CheckPhoneRequest requisição para verificar telefone
type CheckPhoneRequest struct {
	Telefone  string  `query:"telefone" validate:"required"`
	ExcludeID *string `query:"exclude_id" validate:"omitempty,uuid"`
}

// CheckCPFRequest requisição para verificar CPF
type CheckCPFRequest struct {
	CPF       string  `query:"cpf" validate:"required"`
	ExcludeID *string `query:"exclude_id" validate:"omitempty,uuid"`
}

// CheckExistsResponse resposta de verificação de existência
type CheckExistsResponse struct {
	Exists bool `json:"exists"`
}

// =============================================================================
// Gênero Labels
// =============================================================================

var GeneroLabels = map[string]string{
	"M":   "Masculino",
	"F":   "Feminino",
	"NB":  "Não Binário",
	"PNI": "Prefiro não informar",
}

// =============================================================================
// Tags Padrão do Sistema
// =============================================================================

var TagsPadrao = []string{
	"VIP",
	"Recorrente",
	"Inadimplente",
	"Novo",
	"Gastador",
	"Inativo",
}
