package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// CategoriaProdutoEntity representa uma categoria de produto customizável
type CategoriaProdutoEntity struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	Nome         string
	Descricao    string
	Cor          string
	Icone        string
	CentroCusto  CentroCusto
	Ativa        bool
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewCategoriaProduto cria uma nova categoria de produto
func NewCategoriaProduto(tenantID uuid.UUID, nome string) (*CategoriaProdutoEntity, error) {
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant_id é obrigatório")
	}
	if nome == "" {
		return nil, errors.New("nome é obrigatório")
	}
	if len(nome) > 100 {
		return nil, errors.New("nome deve ter no máximo 100 caracteres")
	}

	now := time.Now()
	return &CategoriaProdutoEntity{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Nome:         nome,
		Cor:          "#6B7280", // Cinza padrão
		Icone:        "package",
		CentroCusto:  CentroCustoCMV,
		Ativa:        true,
		CriadoEm:     now,
		AtualizadoEm: now,
	}, nil
}

// Validate valida a entidade
func (c *CategoriaProdutoEntity) Validate() error {
	if c.TenantID == uuid.Nil {
		return errors.New("tenant_id é obrigatório")
	}
	if c.Nome == "" {
		return errors.New("nome é obrigatório")
	}
	if len(c.Nome) > 100 {
		return errors.New("nome deve ter no máximo 100 caracteres")
	}
	if c.Cor != "" && len(c.Cor) != 7 {
		return errors.New("cor deve estar no formato #RRGGBB")
	}
	if !c.CentroCusto.IsValid() {
		return errors.New("centro_custo inválido")
	}
	return nil
}

// SetDescricao define a descrição
func (c *CategoriaProdutoEntity) SetDescricao(descricao string) {
	c.Descricao = descricao
	c.AtualizadoEm = time.Now()
}

// SetCor define a cor
func (c *CategoriaProdutoEntity) SetCor(cor string) error {
	if cor != "" && len(cor) != 7 {
		return errors.New("cor deve estar no formato #RRGGBB")
	}
	c.Cor = cor
	c.AtualizadoEm = time.Now()
	return nil
}

// SetIcone define o ícone
func (c *CategoriaProdutoEntity) SetIcone(icone string) {
	c.Icone = icone
	c.AtualizadoEm = time.Now()
}

// SetCentroCusto define o centro de custo
func (c *CategoriaProdutoEntity) SetCentroCusto(centroCusto CentroCusto) error {
	if !centroCusto.IsValid() {
		return errors.New("centro_custo inválido")
	}
	c.CentroCusto = centroCusto
	c.AtualizadoEm = time.Now()
	return nil
}

// Ativar ativa a categoria
func (c *CategoriaProdutoEntity) Ativar() {
	c.Ativa = true
	c.AtualizadoEm = time.Now()
}

// Desativar desativa a categoria
func (c *CategoriaProdutoEntity) Desativar() {
	c.Ativa = false
	c.AtualizadoEm = time.Now()
}
