package stock

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
)

// CriarProdutoInput representa a entrada do use case de criação de produto
type CriarProdutoInput struct {
	TenantID               uuid.UUID
	SKU                    string // Deprecated - usar CodigoBarras
	Nome                   string
	Descricao              string
	CodigoBarras           *string                 // Novo campo
	Categoria              entity.CategoriaProduto // Enum legado (deprecated)
	CategoriaProdutoID     *uuid.UUID              // FK para categoria customizada
	CentroCusto            entity.CentroCusto
	UnidadeMedida          entity.UnidadeMedida
	ValorUnitario          decimal.Decimal
	QuantidadeMinima       decimal.Decimal
	EstoqueMaximo          int32            // Novo campo (alinhado com entidade)
	ValorVendaProfissional *decimal.Decimal // Novo campo
	ValorEntrada           *decimal.Decimal // Novo campo
	ControlaValidade       bool
	LeadTimeDias           int
	FornecedorID           *uuid.UUID
}

// CriarProdutoOutput representa a saída do use case
type CriarProdutoOutput struct {
	ID                     uuid.UUID
	TenantID               uuid.UUID
	SKU                    *string
	Nome                   string
	Descricao              *string
	CodigoBarras           *string    // Novo campo
	Categoria              string     // Enum legado
	CategoriaProdutoID     *uuid.UUID // FK para categoria customizada
	CentroCusto            string
	UnidadeMedida          string
	ValorUnitario          decimal.Decimal
	QuantidadeAtual        decimal.Decimal
	QuantidadeMinima       decimal.Decimal
	EstoqueMaximo          int32            // Novo campo
	ValorVendaProfissional *decimal.Decimal // Novo campo
	ValorEntrada           *decimal.Decimal // Novo campo
	ControlaValidade       bool
	LeadTimeDias           int
	Ativo                  bool
}

// CriarProdutoUseCase implementa o caso de uso de criação de produto
type CriarProdutoUseCase struct {
	produtoRepo    port.ProdutoRepository
	fornecedorRepo port.FornecedorRepository
}

// NewCriarProdutoUseCase cria uma nova instância do use case
func NewCriarProdutoUseCase(
	produtoRepo port.ProdutoRepository,
	fornecedorRepo port.FornecedorRepository,
) *CriarProdutoUseCase {
	return &CriarProdutoUseCase{
		produtoRepo:    produtoRepo,
		fornecedorRepo: fornecedorRepo,
	}
}

// Execute executa o caso de uso de criação de produto
func (uc *CriarProdutoUseCase) Execute(ctx context.Context, input CriarProdutoInput) (*CriarProdutoOutput, error) {
	// 1. Validar se SKU já existe para o tenant (se informado)
	if input.SKU != "" {
		existingProduto, err := uc.produtoRepo.FindBySKU(ctx, input.TenantID, input.SKU)
		if err != nil {
			return nil, fmt.Errorf("erro ao verificar SKU existente: %w", err)
		}
		if existingProduto != nil {
			return nil, fmt.Errorf("SKU já existe para este tenant")
		}
	}

	// 1.1. Validar se código de barras já existe (se informado)
	if input.CodigoBarras != nil && *input.CodigoBarras != "" {
		existingProduto, err := uc.produtoRepo.FindByCodigoBarras(ctx, input.TenantID, *input.CodigoBarras)
		if err != nil {
			return nil, fmt.Errorf("erro ao verificar código de barras existente: %w", err)
		}
		if existingProduto != nil {
			return nil, fmt.Errorf("código de barras '%s' já está em uso", *input.CodigoBarras)
		}
	}

	// 2. Validar fornecedor (se informado)
	if input.FornecedorID != nil {
		fornecedor, err := uc.fornecedorRepo.FindByID(ctx, input.TenantID, *input.FornecedorID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar fornecedor: %w", err)
		}
		if fornecedor == nil {
			return nil, fmt.Errorf("fornecedor não encontrado")
		}
		if !fornecedor.Ativo {
			return nil, fmt.Errorf("fornecedor está inativo")
		}
	}

	// 3. Criar entidade de produto (SKU foi removido como campo obrigatório)
	produto, err := entity.NewProduto(
		input.TenantID,
		input.Nome,
		input.Categoria,
		input.CentroCusto,
		input.UnidadeMedida,
		input.ValorUnitario,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar produto: %w", err)
	}

	// 4. Preencher campos opcionais
	if input.SKU != "" {
		produto.SKU = &input.SKU
	}
	if input.Descricao != "" {
		produto.Descricao = &input.Descricao
	}
	produto.CodigoBarras = input.CodigoBarras // Novo campo
	produto.QuantidadeMinima = input.QuantidadeMinima
	produto.EstoqueMaximo = input.EstoqueMaximo                   // Novo campo
	produto.ValorVendaProfissional = input.ValorVendaProfissional // Novo campo
	produto.ValorEntrada = input.ValorEntrada                     // Novo campo
	produto.ControlaValidade = input.ControlaValidade
	produto.CategoriaProdutoID = input.CategoriaProdutoID // FK categoria customizada
	if input.LeadTimeDias > 0 {
		produto.LeadTimeDias = input.LeadTimeDias
	}

	// 5. Salvar no repositório
	if err := uc.produtoRepo.Create(ctx, produto); err != nil {
		return nil, fmt.Errorf("erro ao salvar produto: %w", err)
	}

	// 6. Retornar output
	return &CriarProdutoOutput{
		ID:                     produto.ID,
		TenantID:               produto.TenantID,
		SKU:                    produto.SKU,
		Nome:                   produto.Nome,
		Descricao:              produto.Descricao,
		CodigoBarras:           produto.CodigoBarras,
		Categoria:              string(produto.Categoria),
		CategoriaProdutoID:     produto.CategoriaProdutoID,
		CentroCusto:            string(produto.CentroCusto),
		UnidadeMedida:          string(produto.UnidadeMedida),
		ValorUnitario:          produto.Preco,
		QuantidadeAtual:        produto.QuantidadeAtual,
		QuantidadeMinima:       produto.QuantidadeMinima,
		EstoqueMaximo:          produto.EstoqueMaximo,
		ValorVendaProfissional: produto.ValorVendaProfissional,
		ValorEntrada:           produto.ValorEntrada,
		ControlaValidade:       produto.ControlaValidade,
		LeadTimeDias:           produto.LeadTimeDias,
		Ativo:                  produto.Ativo,
	}, nil
}
