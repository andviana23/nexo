package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Erros do domínio
var (
	ErrFornecedorRazaoSocialVazia = errors.New("razão social não pode ser vazia")
	ErrFornecedorTelefoneInvalido = errors.New("telefone inválido")
)

// Fornecedor representa um fornecedor de produtos
type Fornecedor struct {
	ID       uuid.UUID
	TenantID uuid.UUID
	// Dados principais
	RazaoSocial  string
	NomeFantasia string
	CNPJ         string // 14 dígitos, obrigatório

	// Contato
	Email    string
	Telefone string
	Celular  string

	// Endereço completo
	EnderecoLogradouro  string
	EnderecoNumero      string
	EnderecoComplemento string
	EnderecoBairro      string
	EnderecoCidade      string
	EnderecoEstado      string // UF (2 letras)
	EnderecoCEP         string // 8 dígitos

	// Dados bancários
	Banco   string
	Agencia string
	Conta   string

	// Controle
	Observacoes string
	Ativo       bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewFornecedor cria um novo fornecedor com validações
func NewFornecedor(
	tenantID uuid.UUID,
	razaoSocial string,
	telefone string,
) (*Fornecedor, error) {
	// Validações básicas
	if razaoSocial == "" {
		return nil, ErrFornecedorRazaoSocialVazia
	}
	if telefone == "" {
		return nil, ErrFornecedorTelefoneInvalido
	}

	now := time.Now()
	return &Fornecedor{
		ID:          uuid.New(),
		TenantID:    tenantID,
		RazaoSocial: razaoSocial,
		Telefone:    telefone,
		Ativo:       true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Desativar realiza soft delete do fornecedor
func (f *Fornecedor) Desativar() {
	f.Ativo = false
	f.UpdatedAt = time.Now()
}

// Reativar reativa um fornecedor desativado
func (f *Fornecedor) Reativar() {
	f.Ativo = true
	f.UpdatedAt = time.Now()
}

// AtualizarDados atualiza as informações do fornecedor
func (f *Fornecedor) AtualizarDados(
	razaoSocial, nomeFantasia, cnpj, email, telefone, celular string,
	logradouro, numero, complemento, bairro, cidade, estado, cep string,
	banco, agencia, conta, observacoes string,
) error {
	if razaoSocial == "" {
		return ErrFornecedorRazaoSocialVazia
	}
	if telefone == "" && celular == "" {
		return ErrFornecedorTelefoneInvalido
	}

	f.RazaoSocial = razaoSocial
	f.NomeFantasia = nomeFantasia
	f.CNPJ = cnpj
	f.Email = email
	f.Telefone = telefone
	f.Celular = celular

	f.EnderecoLogradouro = logradouro
	f.EnderecoNumero = numero
	f.EnderecoComplemento = complemento
	f.EnderecoBairro = bairro
	f.EnderecoCidade = cidade
	f.EnderecoEstado = estado
	f.EnderecoCEP = cep

	f.Banco = banco
	f.Agencia = agencia
	f.Conta = conta
	f.Observacoes = observacoes

	f.UpdatedAt = time.Now()

	return nil
}
