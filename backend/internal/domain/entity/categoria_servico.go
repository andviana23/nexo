package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Erros do domínio
var (
	ErrCategoriaNomeVazio      = errors.New("nome da categoria não pode ser vazio")
	ErrCategoriaNomeMuitoLongo = errors.New("nome da categoria deve ter no máximo 100 caracteres")
	ErrCategoriaCorInvalida    = errors.New("cor deve estar no formato hexadecimal (#RRGGBB)")
)

// CategoriaServico representa uma categoria de serviço da barbearia
// Exemplos: Cortes de Cabelo, Barba, Tratamentos Capilares, Combos, etc.
// IMPORTANTE: Separada de categorias financeiras (receitas/despesas)
type CategoriaServico struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	UnitID       uuid.UUID // Added for multi-unit isolation
	Nome         string
	Descricao    *string
	Cor          *string // Formato hexadecimal: #RRGGBB
	Icone        *string // Nome do ícone Material Icons (ex: content_cut, face, spa)
	Ativa        bool
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewCategoriaServico cria uma nova categoria de serviço com validações
func NewCategoriaServico(tenantID, unitID uuid.UUID, nome string) (*CategoriaServico, error) {
	// Validações
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return nil, ErrCategoriaNomeVazio
	}
	if len(nome) > 100 {
		return nil, ErrCategoriaNomeMuitoLongo
	}

	return &CategoriaServico{
		ID:           uuid.New(),
		TenantID:     tenantID,
		UnitID:       unitID,
		Nome:         nome,
		Ativa:        true,
		CriadoEm:     time.Now(),
		AtualizadoEm: time.Now(),
	}, nil
}

// SetDescricao define a descrição da categoria
func (c *CategoriaServico) SetDescricao(descricao string) {
	descricao = strings.TrimSpace(descricao)
	if descricao == "" {
		c.Descricao = nil
	} else {
		c.Descricao = &descricao
	}
}

// SetCor define a cor da categoria no formato hexadecimal
func (c *CategoriaServico) SetCor(cor string) error {
	cor = strings.TrimSpace(cor)
	if cor == "" {
		c.Cor = nil
		return nil
	}

	// Validar formato hexadecimal #RRGGBB
	if !isValidHexColor(cor) {
		return ErrCategoriaCorInvalida
	}

	c.Cor = &cor
	return nil
}

// SetIcone define o ícone Material Icons da categoria
func (c *CategoriaServico) SetIcone(icone string) {
	icone = strings.TrimSpace(icone)
	if icone == "" {
		c.Icone = nil
	} else {
		c.Icone = &icone
	}
}

// Ativar marca a categoria como ativa
func (c *CategoriaServico) Ativar() {
	c.Ativa = true
	c.AtualizadoEm = time.Now()
}

// Desativar marca a categoria como inativa
func (c *CategoriaServico) Desativar() {
	c.Ativa = false
	c.AtualizadoEm = time.Now()
}

// Update atualiza os dados da categoria
func (c *CategoriaServico) Update(nome string, descricao, cor, icone *string) error {
	// Validar nome
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return ErrCategoriaNomeVazio
	}
	if len(nome) > 100 {
		return ErrCategoriaNomeMuitoLongo
	}

	// Validar cor se fornecida
	if cor != nil && *cor != "" {
		corValue := strings.TrimSpace(*cor)
		if !isValidHexColor(corValue) {
			return ErrCategoriaCorInvalida
		}
		c.Cor = &corValue
	} else {
		c.Cor = nil
	}

	// Atualizar campos
	c.Nome = nome

	if descricao != nil && *descricao != "" {
		descValue := strings.TrimSpace(*descricao)
		c.Descricao = &descValue
	} else {
		c.Descricao = nil
	}

	if icone != nil && *icone != "" {
		iconeValue := strings.TrimSpace(*icone)
		c.Icone = &iconeValue
	} else {
		c.Icone = nil
	}

	c.AtualizadoEm = time.Now()
	return nil
}

// isValidHexColor valida se a string está no formato hexadecimal #RRGGBB
func isValidHexColor(color string) bool {
	if len(color) != 7 {
		return false
	}
	if color[0] != '#' {
		return false
	}
	for i := 1; i < 7; i++ {
		c := color[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
