package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Erros do domínio de Serviço
var (
	ErrServicoNomeVazio         = errors.New("nome do serviço não pode ser vazio")
	ErrServicoNomeMuitoLongo    = errors.New("nome do serviço deve ter no máximo 255 caracteres")
	ErrServicoPrecoInvalido     = errors.New("preço do serviço deve ser maior que zero")
	ErrServicoDuracaoInvalida   = errors.New("duração do serviço deve ser de pelo menos 5 minutos")
	ErrServicoComissaoInvalida  = errors.New("comissão deve estar entre 0 e 100")
	ErrServicoCorInvalida       = errors.New("cor deve estar no formato hexadecimal (#RRGGBB)")
	ErrServicoCategoriaInvalida = errors.New("categoria inválida")
)

// Servico representa um serviço oferecido pela barbearia
type Servico struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	UnitID           uuid.UUID // Identificador da unidade
	CategoriaID      uuid.UUID // Referência à CategoriaServico (pode ser uuid.Nil)
	Nome             string
	Descricao        string
	Preco            decimal.Decimal
	Duracao          int // em minutos
	Comissao         decimal.Decimal
	Cor              string
	Imagem           string
	ProfissionaisIDs []uuid.UUID // Lista de profissionais que executam este serviço
	Observacoes      string
	Tags             []string
	Ativo            bool
	CriadoEm         time.Time
	AtualizadoEm     time.Time

	// Campos auxiliares (preenchidos em JOINs)
	CategoriaNome string // Nome da categoria (do JOIN)
	CategoriaCor  string // Cor da categoria (do JOIN)
}

// NewServico cria um novo serviço com validações
func NewServico(tenantID, unitID uuid.UUID, nome string, preco decimal.Decimal, duracao int) (*Servico, error) {
	// Validar nome
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return nil, ErrServicoNomeVazio
	}
	if len(nome) > 255 {
		return nil, ErrServicoNomeMuitoLongo
	}

	// Validar preço
	if preco.LessThanOrEqual(decimal.Zero) {
		return nil, ErrServicoPrecoInvalido
	}

	// Validar duração (mínimo 5 minutos)
	if duracao < 5 {
		return nil, ErrServicoDuracaoInvalida
	}

	return &Servico{
		ID:               uuid.New(),
		TenantID:         tenantID,
		UnitID:           unitID,
		Nome:             nome,
		Preco:            preco,
		Duracao:          duracao,
		Comissao:         decimal.Zero,
		ProfissionaisIDs: []uuid.UUID{},
		Tags:             []string{},
		Ativo:            true,
		CriadoEm:         time.Now(),
		AtualizadoEm:     time.Now(),
	}, nil
}

// SetCategoria define a categoria do serviço
func (s *Servico) SetCategoria(categoriaID uuid.UUID) {
	s.CategoriaID = categoriaID
	s.AtualizadoEm = time.Now()
}

// SetDescricao define a descrição do serviço
func (s *Servico) SetDescricao(descricao string) {
	s.Descricao = strings.TrimSpace(descricao)
	s.AtualizadoEm = time.Now()
}

// SetComissao define a comissão padrão do serviço
func (s *Servico) SetComissao(comissao decimal.Decimal) error {
	if comissao.LessThan(decimal.Zero) || comissao.GreaterThan(decimal.NewFromInt(100)) {
		return ErrServicoComissaoInvalida
	}
	s.Comissao = comissao
	s.AtualizadoEm = time.Now()
	return nil
}

// SetCor define a cor do serviço no formato hexadecimal
func (s *Servico) SetCor(cor string) error {
	cor = strings.TrimSpace(cor)
	if cor == "" {
		s.Cor = ""
		return nil
	}

	// Validar formato hexadecimal #RRGGBB
	if !isValidHexColor(cor) {
		return ErrServicoCorInvalida
	}

	s.Cor = cor
	s.AtualizadoEm = time.Now()
	return nil
}

// SetImagem define a URL da imagem do serviço
func (s *Servico) SetImagem(imagem string) {
	s.Imagem = strings.TrimSpace(imagem)
	s.AtualizadoEm = time.Now()
}

// SetObservacoes define as observações do serviço
func (s *Servico) SetObservacoes(observacoes string) {
	s.Observacoes = strings.TrimSpace(observacoes)
	s.AtualizadoEm = time.Now()
}

// SetTags define as tags do serviço
func (s *Servico) SetTags(tags []string) {
	// Limpar tags vazias e duplicadas
	cleanTags := make([]string, 0, len(tags))
	seen := make(map[string]bool)
	for _, tag := range tags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		if tag != "" && !seen[tag] {
			cleanTags = append(cleanTags, tag)
			seen[tag] = true
		}
	}
	s.Tags = cleanTags
	s.AtualizadoEm = time.Now()
}

// SetProfissionais define os profissionais que executam este serviço
func (s *Servico) SetProfissionais(profissionaisIDs []uuid.UUID) {
	// Remover duplicatas
	seen := make(map[uuid.UUID]bool)
	unique := make([]uuid.UUID, 0, len(profissionaisIDs))
	for _, id := range profissionaisIDs {
		if !seen[id] {
			unique = append(unique, id)
			seen[id] = true
		}
	}
	s.ProfissionaisIDs = unique
	s.AtualizadoEm = time.Now()
}

// Ativar marca o serviço como ativo
func (s *Servico) Ativar() {
	s.Ativo = true
	s.AtualizadoEm = time.Now()
}

// Desativar marca o serviço como inativo
func (s *Servico) Desativar() {
	s.Ativo = false
	s.AtualizadoEm = time.Now()
}

// Update atualiza os dados principais do serviço
func (s *Servico) Update(nome string, preco decimal.Decimal, duracao int) error {
	// Validar nome
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return ErrServicoNomeVazio
	}
	if len(nome) > 255 {
		return ErrServicoNomeMuitoLongo
	}

	// Validar preço
	if preco.LessThanOrEqual(decimal.Zero) {
		return ErrServicoPrecoInvalido
	}

	// Validar duração
	if duracao < 5 {
		return ErrServicoDuracaoInvalida
	}

	s.Nome = nome
	s.Preco = preco
	s.Duracao = duracao
	s.AtualizadoEm = time.Now()
	return nil
}

// CalcularPrecoComComissao calcula o valor da comissão sobre o preço
func (s *Servico) CalcularPrecoComComissao() decimal.Decimal {
	if s.Comissao.IsZero() {
		return decimal.Zero
	}
	return s.Preco.Mul(s.Comissao).Div(decimal.NewFromInt(100))
}

// FormatarDuracao retorna a duração formatada (ex: "1h 30min")
func (s *Servico) FormatarDuracao() string {
	horas := s.Duracao / 60
	minutos := s.Duracao % 60

	if horas == 0 {
		return fmt.Sprintf("%dmin", minutos)
	}
	if minutos == 0 {
		return fmt.Sprintf("%dh", horas)
	}
	return fmt.Sprintf("%dh %dmin", horas, minutos)
}
