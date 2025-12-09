package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

// ProdutoRepositoryPG implementa port.ProdutoRepository usando PostgreSQL
type ProdutoRepositoryPG struct {
	queries *db.Queries
}

// NewProdutoRepository cria nova instância do repositório
func NewProdutoRepository(queries *db.Queries) port.ProdutoRepository {
	return &ProdutoRepositoryPG{queries: queries}
}

// Create cria novo produto
func (r *ProdutoRepositoryPG) Create(ctx context.Context, produto *entity.Produto) error {
	params := db.CreateProdutoParams{
		TenantID:               uuidToPgUUID(produto.TenantID),
		CategoriaProdutoID:     uuidPtrToPgUUID(produto.CategoriaProdutoID),
		Nome:                   produto.Nome,
		Descricao:              produto.Descricao,
		Sku:                    produto.SKU,
		CodigoBarras:           produto.CodigoBarras,
		Preco:                  produto.Preco,
		Custo:                  decimalPtrToNumeric(produto.Custo),
		UnidadeMedida:          string(produto.UnidadeMedida),
		QuantidadeAtual:        produto.QuantidadeAtual,
		QuantidadeMinima:       produto.QuantidadeMinima,
		EstoqueMaximo:          int32Ptr(produto.EstoqueMaximo),
		ValorVendaProfissional: decimalPtrToNumeric(produto.ValorVendaProfissional),
		ValorEntrada:           decimalPtrToNumeric(produto.ValorEntrada),
		FornecedorID:           uuidPtrToPgUUID(produto.FornecedorID),
		Localizacao:            produto.Localizacao,
		Lote:                   produto.Lote,
		DataValidade:           timePtrToDate(produto.DataValidade),
		Ncm:                    produto.NCM,
		PermiteVenda:           produto.PermiteVenda,
		Ativo:                  boolPtr(produto.Ativo),
	}

	result, err := r.queries.CreateProduto(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar produto: %w", err)
	}

	produto.ID = pgUUIDToUUID(result.ID)
	produto.CriadoEm = result.CriadoEm.Time
	produto.AtualizadoEm = result.AtualizadoEm.Time

	return nil
}

// FindByID busca produto por ID
func (r *ProdutoRepositoryPG) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*entity.Produto, error) {
	result, err := r.queries.GetProdutoByID(ctx, db.GetProdutoByIDParams{
		ID:       uuidToPgUUID(id),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		// Se não encontrou, retorna nil (não é erro)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}

	return r.toDomain(result), nil
}

// FindBySKU busca produto por SKU
func (r *ProdutoRepositoryPG) FindBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*entity.Produto, error) {
	result, err := r.queries.GetProdutoBySKU(ctx, db.GetProdutoBySKUParams{
		Sku:      strPtrToPgText(sku),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		// Se não encontrou, retorna nil (não é erro)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar produto por SKU: %w", err)
	}

	return r.toDomain(result), nil
}

// FindByCodigoBarras busca produto por código de barras
func (r *ProdutoRepositoryPG) FindByCodigoBarras(ctx context.Context, tenantID uuid.UUID, codigoBarras string) (*entity.Produto, error) {
	result, err := r.queries.GetProdutoByCodigoBarras(ctx, db.GetProdutoByCodigoBarrasParams{
		CodigoBarras: strPtrToPgText(codigoBarras),
		TenantID:     uuidToPgUUID(tenantID),
	})
	if err != nil {
		// Se não encontrou, retorna nil (não é erro)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar produto por código de barras: %w", err)
	}

	return r.toDomain(result), nil
}

// ListAll lista todos os produtos do tenant
func (r *ProdutoRepositoryPG) ListAll(ctx context.Context, tenantID uuid.UUID) ([]*entity.Produto, error) {
	results, err := r.queries.ListProdutos(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar produtos: %w", err)
	}

	produtos := make([]*entity.Produto, len(results))
	for i, result := range results {
		produtos[i] = r.toDomain(result)
	}

	return produtos, nil
}

// ListByCategoria lista produtos por categoria
func (r *ProdutoRepositoryPG) ListByCategoria(ctx context.Context, tenantID uuid.UUID, categoria entity.CategoriaProduto) ([]*entity.Produto, error) {
	results, err := r.queries.ListProdutosByCategoria(ctx, db.ListProdutosByCategoriaParams{
		TenantID:         uuidToPgUUID(tenantID),
		CategoriaProduto: string(categoria),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar produtos por categoria: %w", err)
	}

	produtos := make([]*entity.Produto, len(results))
	for i, result := range results {
		produtos[i] = r.toDomain(result)
	}

	return produtos, nil
}

// ListAbaixoDoMinimo lista produtos com estoque abaixo do mínimo
func (r *ProdutoRepositoryPG) ListAbaixoDoMinimo(ctx context.Context, tenantID uuid.UUID) ([]*entity.Produto, error) {
	results, err := r.queries.ListProdutosAbaixoDoMinimo(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar produtos abaixo do mínimo: %w", err)
	}

	produtos := make([]*entity.Produto, len(results))
	for i, result := range results {
		produtos[i] = r.toDomain(result)
	}

	return produtos, nil
}

// Update atualiza produto existente
func (r *ProdutoRepositoryPG) Update(ctx context.Context, produto *entity.Produto) error {
	params := db.UpdateProdutoParams{
		ID:                     uuidToPgUUID(produto.ID),
		TenantID:               uuidToPgUUID(produto.TenantID),
		CategoriaProdutoID:     uuidPtrToPgUUID(produto.CategoriaProdutoID),
		Nome:                   produto.Nome,
		Descricao:              produto.Descricao,
		Sku:                    produto.SKU,
		CodigoBarras:           produto.CodigoBarras,
		Preco:                  produto.Preco,
		Custo:                  decimalPtrToNumeric(produto.Custo),
		UnidadeMedida:          string(produto.UnidadeMedida),
		QuantidadeMinima:       produto.QuantidadeMinima,
		QuantidadeAtual:        produto.QuantidadeAtual,
		EstoqueMaximo:          int32Ptr(produto.EstoqueMaximo),
		ValorVendaProfissional: decimalPtrToNumeric(produto.ValorVendaProfissional),
		ValorEntrada:           decimalPtrToNumeric(produto.ValorEntrada),
		FornecedorID:           uuidPtrToPgUUID(produto.FornecedorID),
		Localizacao:            produto.Localizacao,
		Lote:                   produto.Lote,
		DataValidade:           timePtrToDate(produto.DataValidade),
		Ncm:                    produto.NCM,
		PermiteVenda:           produto.PermiteVenda,
	}

	result, err := r.queries.UpdateProduto(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar produto: %w", err)
	}

	produto.AtualizadoEm = result.AtualizadoEm.Time
	return nil
}

// AtualizarQuantidade atualiza apenas a quantidade do produto
func (r *ProdutoRepositoryPG) AtualizarQuantidade(ctx context.Context, tenantID, id uuid.UUID, novaQuantidade decimal.Decimal) error {
	_, err := r.queries.AtualizarQuantidadeProduto(ctx, db.AtualizarQuantidadeProdutoParams{
		ID:              uuidToPgUUID(id),
		TenantID:        uuidToPgUUID(tenantID),
		QuantidadeAtual: novaQuantidade,
	})
	if err != nil {
		return fmt.Errorf("erro ao atualizar quantidade: %w", err)
	}

	return nil
}

// Delete soft-delete do produto
func (r *ProdutoRepositoryPG) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	err := r.queries.DeleteProduto(ctx, db.DeleteProdutoParams{
		ID:       uuidToPgUUID(id),
		TenantID: uuidToPgUUID(tenantID),
	})
	if err != nil {
		return fmt.Errorf("erro ao deletar produto: %w", err)
	}

	return nil
}

// toDomain converte modelo sqlc para entidade de domínio
func (r *ProdutoRepositoryPG) toDomain(p db.Produto) *entity.Produto {
	var estoqueMaximo int32
	if p.EstoqueMaximo != nil {
		estoqueMaximo = *p.EstoqueMaximo
	}

	return &entity.Produto{
		ID:                     pgUUIDToUUID(p.ID),
		TenantID:               pgUUIDToUUID(p.TenantID),
		CategoriaProdutoID:     pgUUIDToUUIDPtr(p.CategoriaProdutoID),
		FornecedorID:           pgUUIDToUUIDPtr(p.FornecedorID),
		Nome:                   p.Nome,
		Descricao:              p.Descricao,
		SKU:                    p.Sku,
		CodigoBarras:           p.CodigoBarras,
		Preco:                  p.Preco,
		Custo:                  numericToDecimalPtr(p.Custo),
		ValorVendaProfissional: numericToDecimalPtr(p.ValorVendaProfissional),
		ValorEntrada:           numericToDecimalPtr(p.ValorEntrada),
		UnidadeMedida:          entity.UnidadeMedida(p.UnidadeMedida),
		QuantidadeAtual:        p.QuantidadeAtual,
		QuantidadeMinima:       p.QuantidadeMinima,
		EstoqueMaximo:          estoqueMaximo,
		Localizacao:            p.Localizacao,
		Lote:                   p.Lote,
		DataValidade:           dateToTimePtr(p.DataValidade),
		NCM:                    p.Ncm,
		PermiteVenda:           p.PermiteVenda,
		Ativo:                  boolPtrToBool(p.Ativo),
		CriadoEm:               p.CriadoEm.Time,
		AtualizadoEm:           p.AtualizadoEm.Time,
	}
}
