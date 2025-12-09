package entity

import (
	"regexp"
	"strings"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain"
	"github.com/google/uuid"
)

// Customer representa um cliente no sistema
type Customer struct {
	ID       string
	TenantID uuid.UUID

	// Dados Básicos (obrigatórios)
	Nome     string
	Telefone string

	// Dados Opcionais
	Email          *string
	CPF            *string
	DataNascimento *time.Time
	Genero         *string // "M", "F", "NB", "PNI"

	// Endereço
	EnderecoLogradouro  *string
	EnderecoNumero      *string
	EnderecoComplemento *string
	EnderecoBairro      *string
	EnderecoCidade      *string
	EnderecoEstado      *string
	EnderecoCEP         *string

	// CRM
	Observacoes *string
	Tags        []string

	// Status
	Ativo bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCustomer cria um novo cliente validado
func NewCustomer(
	tenantID uuid.UUID,
	nome string,
	telefone string,
) (*Customer, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrTenantIDRequired
	}
	if nome == "" {
		return nil, domain.ErrCustomerNameRequired
	}
	if len(nome) < 3 {
		return nil, domain.ErrCustomerNameTooShort
	}
	if telefone == "" {
		return nil, domain.ErrCustomerPhoneRequired
	}

	// Normaliza telefone (remove caracteres não numéricos)
	telefoneNormalizado := normalizePhone(telefone)
	if len(telefoneNormalizado) < 10 || len(telefoneNormalizado) > 11 {
		return nil, domain.ErrCustomerPhoneInvalid
	}

	now := time.Now()
	return &Customer{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		Nome:      nome,
		Telefone:  telefoneNormalizado,
		Ativo:     true,
		Tags:      []string{},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// SetEmail define o email do cliente
func (c *Customer) SetEmail(email string) error {
	if email == "" {
		c.Email = nil
		return nil
	}
	if !isValidEmail(email) {
		return domain.ErrCustomerEmailInvalid
	}
	c.Email = &email
	c.UpdatedAt = time.Now()
	return nil
}

// SetCPF define o CPF do cliente
func (c *Customer) SetCPF(cpf string) error {
	if cpf == "" {
		c.CPF = nil
		return nil
	}
	cpfNormalizado := normalizeCPF(cpf)
	if len(cpfNormalizado) != 11 {
		return domain.ErrCustomerCPFInvalid
	}
	c.CPF = &cpfNormalizado
	c.UpdatedAt = time.Now()
	return nil
}

// SetDataNascimento define a data de nascimento
func (c *Customer) SetDataNascimento(data *time.Time) error {
	if data != nil && data.After(time.Now()) {
		return domain.ErrCustomerBirthDateFuture
	}
	c.DataNascimento = data
	c.UpdatedAt = time.Now()
	return nil
}

// SetGenero define o gênero do cliente
func (c *Customer) SetGenero(genero string) error {
	if genero != "" && genero != "M" && genero != "F" && genero != "NB" && genero != "PNI" {
		return domain.ErrCustomerGenderInvalid
	}
	if genero == "" {
		c.Genero = nil
	} else {
		c.Genero = &genero
	}
	c.UpdatedAt = time.Now()
	return nil
}

// SetEndereco define o endereço completo
func (c *Customer) SetEndereco(
	logradouro, numero, complemento, bairro, cidade, estado, cep *string,
) error {
	// Valida CEP se fornecido
	if cep != nil && *cep != "" {
		cepNormalizado := normalizeCEP(*cep)
		if len(cepNormalizado) != 8 {
			return domain.ErrCustomerCEPInvalid
		}
		c.EnderecoCEP = &cepNormalizado
	} else {
		c.EnderecoCEP = nil
	}

	// Valida UF se fornecido
	if estado != nil && *estado != "" && len(*estado) != 2 {
		return domain.ErrCustomerStateInvalid
	}

	c.EnderecoLogradouro = logradouro
	c.EnderecoNumero = numero
	c.EnderecoComplemento = complemento
	c.EnderecoBairro = bairro
	c.EnderecoCidade = cidade
	c.EnderecoEstado = estado
	c.UpdatedAt = time.Now()
	return nil
}

// SetObservacoes define as observações internas
func (c *Customer) SetObservacoes(obs string) error {
	if len(obs) > 500 {
		return domain.ErrCustomerObservationsTooLong
	}
	if obs == "" {
		c.Observacoes = nil
	} else {
		c.Observacoes = &obs
	}
	c.UpdatedAt = time.Now()
	return nil
}

// AddTag adiciona uma tag ao cliente
func (c *Customer) AddTag(tag string) error {
	if len(c.Tags) >= 10 {
		return domain.ErrCustomerMaxTagsExceeded
	}
	if len(tag) > 50 {
		return domain.ErrCustomerTagTooLong
	}
	// Verifica se já existe
	for _, t := range c.Tags {
		if strings.EqualFold(t, tag) {
			return nil // Já existe, não faz nada
		}
	}
	c.Tags = append(c.Tags, tag)
	c.UpdatedAt = time.Now()
	return nil
}

// RemoveTag remove uma tag do cliente
func (c *Customer) RemoveTag(tag string) {
	for i, t := range c.Tags {
		if strings.EqualFold(t, tag) {
			c.Tags = append(c.Tags[:i], c.Tags[i+1:]...)
			c.UpdatedAt = time.Now()
			return
		}
	}
}

// SetTags substitui todas as tags
func (c *Customer) SetTags(tags []string) error {
	if len(tags) > 10 {
		return domain.ErrCustomerMaxTagsExceeded
	}
	for _, tag := range tags {
		if len(tag) > 50 {
			return domain.ErrCustomerTagTooLong
		}
	}
	c.Tags = tags
	c.UpdatedAt = time.Now()
	return nil
}

// HasTag verifica se cliente tem uma tag
func (c *Customer) HasTag(tag string) bool {
	for _, t := range c.Tags {
		if strings.EqualFold(t, tag) {
			return true
		}
	}
	return false
}

// Inactivate inativa o cliente (soft delete)
func (c *Customer) Inactivate() {
	c.Ativo = false
	c.UpdatedAt = time.Now()
}

// Reactivate reativa o cliente
func (c *Customer) Reactivate() {
	c.Ativo = true
	c.UpdatedAt = time.Now()
}

// UpdateNome atualiza o nome do cliente
func (c *Customer) UpdateNome(nome string) error {
	if nome == "" {
		return domain.ErrCustomerNameRequired
	}
	if len(nome) < 3 {
		return domain.ErrCustomerNameTooShort
	}
	c.Nome = nome
	c.UpdatedAt = time.Now()
	return nil
}

// UpdateTelefone atualiza o telefone do cliente
func (c *Customer) UpdateTelefone(telefone string) error {
	if telefone == "" {
		return domain.ErrCustomerPhoneRequired
	}
	telefoneNormalizado := normalizePhone(telefone)
	if len(telefoneNormalizado) < 10 || len(telefoneNormalizado) > 11 {
		return domain.ErrCustomerPhoneInvalid
	}
	c.Telefone = telefoneNormalizado
	c.UpdatedAt = time.Now()
	return nil
}

// Validate valida as regras de negócio do cliente
func (c *Customer) Validate() error {
	if c.TenantID == uuid.Nil {
		return domain.ErrTenantIDRequired
	}
	if c.Nome == "" {
		return domain.ErrCustomerNameRequired
	}
	if len(c.Nome) < 3 {
		return domain.ErrCustomerNameTooShort
	}
	if c.Telefone == "" {
		return domain.ErrCustomerPhoneRequired
	}

	telefoneNormalizado := normalizePhone(c.Telefone)
	if len(telefoneNormalizado) < 10 || len(telefoneNormalizado) > 11 {
		return domain.ErrCustomerPhoneInvalid
	}

	if c.Email != nil && *c.Email != "" && !isValidEmail(*c.Email) {
		return domain.ErrCustomerEmailInvalid
	}

	if c.CPF != nil && *c.CPF != "" {
		cpfNormalizado := normalizeCPF(*c.CPF)
		if len(cpfNormalizado) != 11 {
			return domain.ErrCustomerCPFInvalid
		}
	}

	if c.EnderecoCEP != nil && *c.EnderecoCEP != "" {
		cepNormalizado := normalizeCEP(*c.EnderecoCEP)
		if len(cepNormalizado) != 8 {
			return domain.ErrCustomerCEPInvalid
		}
	}

	if c.Observacoes != nil && len(*c.Observacoes) > 500 {
		return domain.ErrCustomerObservationsTooLong
	}

	if len(c.Tags) > 10 {
		return domain.ErrCustomerMaxTagsExceeded
	}

	return nil
}

// =============================================================================
// Helpers
// =============================================================================

// normalizePhone remove caracteres não numéricos do telefone
func normalizePhone(phone string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(phone, "")
}

// normalizeCPF remove caracteres não numéricos do CPF
func normalizeCPF(cpf string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(cpf, "")
}

// normalizeCEP remove caracteres não numéricos do CEP
func normalizeCEP(cep string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(cep, "")
}

// isValidEmail valida formato de email
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
