package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CategoriaProduto representa as categorias de produtos
type CategoriaProduto string

const (
	CategoriaPomada         CategoriaProduto = "POMADA"
	CategoriaShampoo        CategoriaProduto = "SHAMPOO"
	CategoriaCreme          CategoriaProduto = "CREME"
	CategoriaLamina         CategoriaProduto = "LAMINA"
	CategoriaToalha         CategoriaProduto = "TOALHA"
	CategoriaLimpeza        CategoriaProduto = "LIMPEZA"
	CategoriaEscritorio     CategoriaProduto = "ESCRITORIO"
	CategoriaBebida         CategoriaProduto = "BEBIDA"
	CategoriaRevenda        CategoriaProduto = "REVENDA"
	CategoriaInsumo         CategoriaProduto = "INSUMO" // Mantido por compatibilidade ou genérico
	CategoriaUsoInterno     CategoriaProduto = "USO_INTERNO"
	CategoriaPermanente     CategoriaProduto = "PERMANENTE"
	CategoriaPromocional    CategoriaProduto = "PROMOCIONAL"
	CategoriaKit            CategoriaProduto = "KIT"
	CategoriaProdutoServico CategoriaProduto = "SERVICO"
)

// CentroCusto representa o centro de custo do produto
type CentroCusto string

const (
	CentroCustoServico     CentroCusto = "CUSTO_SERVICO"       // Insumos, Lâminas, Shampoos
	CentroCustoOperacional CentroCusto = "DESPESA_OPERACIONAL" // Limpeza, Escritório
	CentroCustoCMV         CentroCusto = "CMV"                 // Revenda, Bebidas
)

// IsValid verifica se o centro de custo é válido
func (c CentroCusto) IsValid() bool {
	switch c {
	case CentroCustoServico, CentroCustoOperacional, CentroCustoCMV:
		return true
	}
	return false
}

// UnidadeMedida representa as unidades de medida
// Valores conforme constraint do banco: CHECK (unidade_medida IN ('UNIDADE', 'LITRO', 'MILILITRO', 'GRAMA', 'QUILOGRAMA'))
type UnidadeMedida string

const (
	UnidadeUnidade   UnidadeMedida = "UNIDADE"    // Unidade
	UnidadeLitro     UnidadeMedida = "LITRO"      // Litro
	UnidadeMililitro UnidadeMedida = "MILILITRO"  // Mililitro
	UnidadeGrama     UnidadeMedida = "GRAMA"      // Grama
	UnidadeKilograma UnidadeMedida = "QUILOGRAMA" // Quilograma
)

// Erros do domínio
var (
	ErrProdutoNomeVazio           = errors.New("nome do produto não pode ser vazio")
	ErrProdutoValorInvalido       = errors.New("valor unitário deve ser maior que zero")
	ErrProdutoEstoqueNegativo     = errors.New("estoque não pode ser negativo")
	ErrProdutoEstoqueInsuficiente = errors.New("estoque insuficiente para a operação")
	ErrProdutoQuantidadeInvalida  = errors.New("quantidade deve ser maior que zero")
)

// Produto representa um produto/insumo do estoque
type Produto struct {
	ID                     uuid.UUID
	TenantID               uuid.UUID
	CategoriaID            *uuid.UUID // Categoria financeira (receita/despesa)
	CategoriaProdutoID     *uuid.UUID // FK para categorias_produtos customizadas
	FornecedorID           *uuid.UUID // FK para fornecedores
	Nome                   string
	Descricao              *string
	SKU                    *string
	CodigoBarras           *string
	Preco                  decimal.Decimal
	Custo                  *decimal.Decimal
	ValorVendaProfissional *decimal.Decimal // Preço de venda para profissionais
	ValorEntrada           *decimal.Decimal // Custo de aquisição/entrada
	Categoria              CategoriaProduto
	CentroCusto            CentroCusto
	UnidadeMedida          UnidadeMedida
	QuantidadeAtual        decimal.Decimal
	QuantidadeMinima       decimal.Decimal
	EstoqueMaximo          int32 // Quantidade máxima em estoque
	Localizacao            *string
	Lote                   *string
	DataValidade           *time.Time
	NCM                    *string
	PermiteVenda           bool
	ControlaValidade       bool
	LeadTimeDias           int
	Ativo                  bool
	CriadoEm               time.Time
	AtualizadoEm           time.Time
}

// NewProduto cria um novo produto com validações
// SKU agora é opcional (substituído por código de barras)
func NewProduto(
	tenantID uuid.UUID,
	nome string,
	categoria CategoriaProduto,
	centroCusto CentroCusto,
	unidadeMedida UnidadeMedida,
	preco decimal.Decimal,
) (*Produto, error) {
	// Validações
	if nome == "" {
		return nil, ErrProdutoNomeVazio
	}
	if preco.LessThanOrEqual(decimal.Zero) {
		return nil, ErrProdutoValorInvalido
	}

	now := time.Now()
	return &Produto{
		ID:               uuid.New(),
		TenantID:         tenantID,
		Nome:             nome,
		Categoria:        categoria,
		CentroCusto:      centroCusto,
		UnidadeMedida:    unidadeMedida,
		Preco:            preco,
		QuantidadeAtual:  decimal.Zero,
		QuantidadeMinima: decimal.Zero,
		PermiteVenda:     true,
		ControlaValidade: false,
		LeadTimeDias:     7, // Default
		Ativo:            true,
		CriadoEm:         now,
		AtualizadoEm:     now,
	}, nil
}

// EstaBaixo verifica se o produto está abaixo do estoque mínimo
func (p *Produto) EstaBaixo() bool {
	return p.QuantidadeAtual.LessThanOrEqual(p.QuantidadeMinima)
}

// AdicionarEstoque adiciona quantidade ao estoque (ENTRADA)
func (p *Produto) AdicionarEstoque(quantidade decimal.Decimal) error {
	if quantidade.LessThanOrEqual(decimal.Zero) {
		return ErrProdutoQuantidadeInvalida
	}

	p.QuantidadeAtual = p.QuantidadeAtual.Add(quantidade)
	p.AtualizadoEm = time.Now()
	return nil
}

// RemoverEstoque remove quantidade do estoque (SAIDA) - RN-EST-002
func (p *Produto) RemoverEstoque(quantidade decimal.Decimal) error {
	if quantidade.LessThanOrEqual(decimal.Zero) {
		return ErrProdutoQuantidadeInvalida
	}

	if p.QuantidadeAtual.LessThan(quantidade) {
		return ErrProdutoEstoqueInsuficiente
	}

	p.QuantidadeAtual = p.QuantidadeAtual.Sub(quantidade)
	p.AtualizadoEm = time.Now()
	return nil
}

// AjustarEstoque permite ajuste manual de estoque (inventário)
func (p *Produto) AjustarEstoque(novaQuantidade decimal.Decimal, motivo string) error {
	if novaQuantidade.LessThan(decimal.Zero) {
		return ErrProdutoEstoqueNegativo
	}

	if motivo == "" {
		return errors.New("motivo do ajuste é obrigatório")
	}

	p.QuantidadeAtual = novaQuantidade
	p.AtualizadoEm = time.Now()
	return nil
}

// AtualizarPreco atualiza o preço do produto
func (p *Produto) AtualizarPreco(novoPreco decimal.Decimal) error {
	if novoPreco.LessThanOrEqual(decimal.Zero) {
		return ErrProdutoValorInvalido
	}

	p.Preco = novoPreco
	p.AtualizadoEm = time.Now()
	return nil
}

// Desativar realiza soft delete do produto
func (p *Produto) Desativar() {
	p.Ativo = false
	p.AtualizadoEm = time.Now()
}

// Reativar reativa um produto desativado
func (p *Produto) Reativar() {
	p.Ativo = true
	p.AtualizadoEm = time.Now()
}

// DefinirEstoqueMinimo configura o estoque mínimo para alertas
func (p *Produto) DefinirEstoqueMinimo(quantidade decimal.Decimal) error {
	if quantidade.LessThan(decimal.Zero) {
		return errors.New("estoque mínimo não pode ser negativo")
	}

	p.QuantidadeMinima = quantidade
	p.AtualizadoEm = time.Now()
	return nil
}
