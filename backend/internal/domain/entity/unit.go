package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Erros de domínio para Unit
var (
	ErrUnitNomeVazio       = errors.New("nome da unidade não pode ser vazio")
	ErrUnitNomeDuplicado   = errors.New("já existe uma unidade com este nome")
	ErrUnitMatrizNaoDelete = errors.New("não é possível excluir a unidade matriz")
	ErrUnitNaoEncontrada   = errors.New("unidade não encontrada")
)

// Unit representa uma unidade/filial do tenant
type Unit struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	Nome           string
	Apelido        *string
	Descricao      *string
	EnderecoResumo *string
	Cidade         *string
	Estado         *string
	Timezone       string
	Ativa          bool
	IsMatriz       bool
	CriadoEm       time.Time
	AtualizadoEm   time.Time
}

// NewUnit cria uma nova unidade com validações
func NewUnit(tenantID uuid.UUID, nome string) (*Unit, error) {
	if nome == "" {
		return nil, ErrUnitNomeVazio
	}

	return &Unit{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Nome:         nome,
		Timezone:     "America/Sao_Paulo",
		Ativa:        true,
		IsMatriz:     false,
		CriadoEm:     time.Now(),
		AtualizadoEm: time.Now(),
	}, nil
}

// NewMatrizUnit cria uma unidade matriz (principal)
func NewMatrizUnit(tenantID uuid.UUID, nome string) (*Unit, error) {
	unit, err := NewUnit(tenantID, nome)
	if err != nil {
		return nil, err
	}
	unit.IsMatriz = true
	return unit, nil
}

// SetApelido define o apelido da unidade
func (u *Unit) SetApelido(apelido string) {
	if apelido != "" {
		u.Apelido = &apelido
	} else {
		u.Apelido = nil
	}
	u.AtualizadoEm = time.Now()
}

// SetEndereco define o endereço resumido
func (u *Unit) SetEndereco(endereco, cidade, estado string) {
	if endereco != "" {
		u.EnderecoResumo = &endereco
	}
	if cidade != "" {
		u.Cidade = &cidade
	}
	if estado != "" {
		u.Estado = &estado
	}
	u.AtualizadoEm = time.Now()
}

// Ativar ativa a unidade
func (u *Unit) Ativar() {
	u.Ativa = true
	u.AtualizadoEm = time.Now()
}

// Desativar desativa a unidade
func (u *Unit) Desativar() {
	u.Ativa = false
	u.AtualizadoEm = time.Now()
}

// CanDelete verifica se a unidade pode ser excluída
func (u *Unit) CanDelete() error {
	if u.IsMatriz {
		return ErrUnitMatrizNaoDelete
	}
	return nil
}

// DisplayName retorna o nome de exibição (apelido ou nome)
func (u *Unit) DisplayName() string {
	if u.Apelido != nil && *u.Apelido != "" {
		return *u.Apelido
	}
	return u.Nome
}
