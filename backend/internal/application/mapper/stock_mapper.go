// Package mapper contém funções para conversão entre entidades de domínio e DTOs HTTP para Stock.
package mapper

import (
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/stock"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// === MAPEAMENTO DE CATEGORIA ===
// O DTO usa categorias de negócio (POMADA, SHAMPOO, etc.)
// A entidade usa categorias técnicas (INSUMO, REVENDA, etc.)
// Mapeamos DTO -> Entity para criação

// mapCategoriaDTOToEntity converte categoria do DTO para entidade
func mapCategoriaDTOToEntity(categoria string) entity.CategoriaProduto {
	switch categoria {
	case "POMADA", "SHAMPOO", "CREME", "CONSUMIVEL":
		return entity.CategoriaInsumo
	case "REVENDA":
		return entity.CategoriaRevenda
	case "LAMINA", "TOALHA":
		return entity.CategoriaUsoInterno
	default:
		return entity.CategoriaInsumo
	}
}

// mapUnidadeDTOToEntity converte unidade de medida do DTO para entidade
// Aceita tanto valores do banco (UNIDADE, LITRO, etc.) quanto abreviações antigas (UN, L, etc.)
func mapUnidadeDTOToEntity(unidade string) entity.UnidadeMedida {
	switch unidade {
	// Valores do banco de dados (corretos)
	case "UNIDADE":
		return entity.UnidadeUnidade
	case "QUILOGRAMA":
		return entity.UnidadeKilograma
	case "GRAMA":
		return entity.UnidadeGrama
	case "MILILITRO":
		return entity.UnidadeMililitro
	case "LITRO":
		return entity.UnidadeLitro
	// Abreviações antigas (retrocompatibilidade)
	case "UN":
		return entity.UnidadeUnidade
	case "KG":
		return entity.UnidadeKilograma
	case "G":
		return entity.UnidadeGrama
	case "ML":
		return entity.UnidadeMililitro
	case "L":
		return entity.UnidadeLitro
	default:
		return entity.UnidadeUnidade
	}
}

// FromRegistrarEntradaRequest converte DTO Request para Input do use case
func FromRegistrarEntradaRequest(
	tenantID uuid.UUID,
	userID uuid.UUID,
	req dto.RegistrarEntradaRequest,
) (stock.RegistrarEntradaInput, error) {
	// Parse FornecedorID
	fornecedorID, err := uuid.Parse(req.FornecedorID)
	if err != nil {
		return stock.RegistrarEntradaInput{}, fmt.Errorf("fornecedor_id inválido: %w", err)
	}

	// Parse DataEntrada
	dataEntrada, err := time.Parse("2006-01-02", req.DataEntrada)
	if err != nil {
		return stock.RegistrarEntradaInput{}, fmt.Errorf("data_entrada inválida: %w", err)
	}

	// Converter itens
	itens := make([]stock.ItemEntrada, len(req.Itens))
	for i, itemReq := range req.Itens {
		produtoID, err := uuid.Parse(itemReq.ProdutoID)
		if err != nil {
			return stock.RegistrarEntradaInput{}, fmt.Errorf("produto_id inválido no item %d: %w", i, err)
		}

		valorUnitario, err := decimal.NewFromString(itemReq.ValorUnitario)
		if err != nil {
			return stock.RegistrarEntradaInput{}, fmt.Errorf("valor_unitario inválido no item %d: %w", i, err)
		}

		itens[i] = stock.ItemEntrada{
			ProdutoID:     produtoID,
			Quantidade:    itemReq.Quantidade, // Já é int no DTO
			ValorUnitario: valorUnitario,
		}
	}

	return stock.RegistrarEntradaInput{
		TenantID:        tenantID,
		UsuarioID:       userID,
		FornecedorID:    fornecedorID,
		DataEntrada:     dataEntrada,
		Itens:           itens,
		Observacoes:     req.Observacoes,
		GerarFinanceiro: req.GerarFinanceiro,
	}, nil
}

// ToRegistrarEntradaResponse converte Output do use case para DTO Response
func ToRegistrarEntradaResponse(output *stock.RegistrarEntradaOutput) dto.RegistrarEntradaResponse {
	movimentacoesIDs := make([]string, len(output.MovimentacoesIDs))
	for i, id := range output.MovimentacoesIDs {
		movimentacoesIDs[i] = id.String()
	}

	return dto.RegistrarEntradaResponse{
		MovimentacoesIDs: movimentacoesIDs,
		ValorTotal:       output.ValorTotal.StringFixed(2),
		ItensProcessados: output.ItensProcessados,
	}
}

// FromRegistrarSaidaRequest converte DTO Request para Input (neste caso, retorna o próprio DTO)
func FromRegistrarSaidaRequest(
	tenantID uuid.UUID,
	userID uuid.UUID,
	req dto.RegistrarSaidaRequest,
) (dto.RegistrarSaidaRequest, error) {
	// Use cases de Saída, Ajuste e Listar Alertas usam DTOs diretamente
	// Este mapper apenas valida
	return req, nil
}

// FromAjustarEstoqueRequest converte DTO Request para Input
func FromAjustarEstoqueRequest(
	tenantID uuid.UUID,
	userID uuid.UUID,
	req dto.AjustarEstoqueRequest,
) (dto.AjustarEstoqueRequest, error) {
	// Validação básica
	if req.NovaQuantidade == "" {
		return dto.AjustarEstoqueRequest{}, fmt.Errorf("nova_quantidade é obrigatória")
	}
	if req.Motivo == "" {
		return dto.AjustarEstoqueRequest{}, fmt.Errorf("motivo é obrigatório")
	}
	return req, nil
}

// ToMovimentacaoResponse converte MovimentacaoResponse (já é DTO, apenas retorna)
func ToMovimentacaoResponse(resp *dto.MovimentacaoResponse) *dto.MovimentacaoResponse {
	return resp
}

// ToListAlertasEstoqueBaixoResponse converte resposta (já é DTO, apenas retorna)
func ToListAlertasEstoqueBaixoResponse(resp *dto.ListAlertasEstoqueBaixoResponse) *dto.ListAlertasEstoqueBaixoResponse {
	return resp
}

// === PRODUTO MAPPERS ===

// FromCreateProdutoRequest converte DTO Request para Input do use case
func FromCreateProdutoRequest(
	tenantID uuid.UUID,
	req dto.CreateProdutoRequest,
) (stock.CriarProdutoInput, error) {
	// Validar e converter valor unitário
	valorUnitario, err := decimal.NewFromString(req.ValorUnitario)
	if err != nil {
		return stock.CriarProdutoInput{}, fmt.Errorf("valor_unitario inválido: %w", err)
	}

	// Converter quantidade mínima (float para decimal)
	quantidadeMinima := decimal.NewFromFloat(req.QuantidadeMinima)

	// Converter fornecedor ID (se presente)
	var fornecedorID *uuid.UUID
	if req.FornecedorID != nil && *req.FornecedorID != "" {
		id, err := uuid.Parse(*req.FornecedorID)
		if err != nil {
			return stock.CriarProdutoInput{}, fmt.Errorf("fornecedor_id inválido: %w", err)
		}
		fornecedorID = &id
	}

	// Converter categoria_produto_id (obrigatório no novo DTO)
	var categoriaProdutoID *uuid.UUID
	if req.CategoriaProdutoID != "" {
		id, err := uuid.Parse(req.CategoriaProdutoID)
		if err != nil {
			return stock.CriarProdutoInput{}, fmt.Errorf("categoria_produto_id inválido: %w", err)
		}
		categoriaProdutoID = &id
	}

	// Converter quantidade_maxima -> estoque_maximo (se presente)
	var estoqueMaximo int32
	if req.QuantidadeMaxima != nil && *req.QuantidadeMaxima != "" {
		qtdMax, err := decimal.NewFromString(*req.QuantidadeMaxima)
		if err != nil {
			return stock.CriarProdutoInput{}, fmt.Errorf("quantidade_maxima inválida: %w", err)
		}
		estoqueMaximo = int32(qtdMax.IntPart())
	}

	// Converter valor_venda_profissional (se presente)
	var valorVendaProfissional *decimal.Decimal
	if req.ValorVendaProfissional != nil && *req.ValorVendaProfissional != "" {
		val, err := decimal.NewFromString(*req.ValorVendaProfissional)
		if err != nil {
			return stock.CriarProdutoInput{}, fmt.Errorf("valor_venda_profissional inválido: %w", err)
		}
		valorVendaProfissional = &val
	}

	// Converter valor_entrada (se presente)
	var valorEntrada *decimal.Decimal
	if req.ValorEntrada != nil && *req.ValorEntrada != "" {
		val, err := decimal.NewFromString(*req.ValorEntrada)
		if err != nil {
			return stock.CriarProdutoInput{}, fmt.Errorf("valor_entrada inválido: %w", err)
		}
		valorEntrada = &val
	}

	return stock.CriarProdutoInput{
		TenantID:               tenantID,
		Nome:                   req.Nome,
		Descricao:              req.Descricao,
		CodigoBarras:           req.CodigoBarras,
		Categoria:              entity.CategoriaRevenda, // Default - categoria legada
		CategoriaProdutoID:     categoriaProdutoID,
		CentroCusto:            entity.CentroCustoCMV, // Default
		UnidadeMedida:          mapUnidadeDTOToEntity(req.UnidadeMedida),
		ValorUnitario:          valorUnitario,
		QuantidadeMinima:       quantidadeMinima,
		EstoqueMaximo:          estoqueMaximo,
		ValorVendaProfissional: valorVendaProfissional,
		ValorEntrada:           valorEntrada,
		FornecedorID:           fornecedorID,
	}, nil
}

// CriarProdutoOutputToResponse converte Output do use case de criar produto para DTO Response
func CriarProdutoOutputToResponse(output *stock.CriarProdutoOutput) dto.ProdutoResponse {
	var categoriaProdutoIDStr *string
	if output.CategoriaProdutoID != nil {
		s := output.CategoriaProdutoID.String()
		categoriaProdutoIDStr = &s
	}

	var descricao *string
	if output.Descricao != nil {
		descricao = output.Descricao
	}

	var codigoBarras *string
	if output.CodigoBarras != nil {
		codigoBarras = output.CodigoBarras
	}

	var quantidadeMaxima *string
	if output.EstoqueMaximo > 0 {
		s := fmt.Sprintf("%d", output.EstoqueMaximo)
		quantidadeMaxima = &s
	}

	var valorVendaProfissional *string
	if output.ValorVendaProfissional != nil {
		s := output.ValorVendaProfissional.StringFixed(2)
		valorVendaProfissional = &s
	}

	var valorEntrada *string
	if output.ValorEntrada != nil {
		s := output.ValorEntrada.StringFixed(2)
		valorEntrada = &s
	}

	return dto.ProdutoResponse{
		ID:                     output.ID.String(),
		TenantID:               output.TenantID.String(),
		Nome:                   output.Nome,
		Descricao:              descricao,
		CodigoBarras:           codigoBarras,
		CategoriaProdutoID:     categoriaProdutoIDStr,
		UnidadeMedida:          output.UnidadeMedida,
		ValorUnitario:          output.ValorUnitario.StringFixed(2),
		QuantidadeAtual:        output.QuantidadeAtual.StringFixed(2),
		QuantidadeMinima:       output.QuantidadeMinima.StringFixed(2),
		QuantidadeMaxima:       quantidadeMaxima,
		ValorVendaProfissional: valorVendaProfissional,
		ValorEntrada:           valorEntrada,
		EstaBaixo:              output.QuantidadeAtual.LessThanOrEqual(output.QuantidadeMinima),
		Ativo:                  output.Ativo,
		CreatedAt:              time.Now().Format(time.RFC3339),
		UpdatedAt:              time.Now().Format(time.RFC3339),
	}
}

// ToProdutoResponse converte entity.Produto para DTO Response
func ToProdutoResponse(produto *entity.Produto) dto.ProdutoResponse {
	var categoriaProdutoIDStr *string
	if produto.CategoriaProdutoID != nil {
		s := produto.CategoriaProdutoID.String()
		categoriaProdutoIDStr = &s
	}

	var fornecedorIDStr *string
	if produto.FornecedorID != nil {
		s := produto.FornecedorID.String()
		fornecedorIDStr = &s
	}

	var descricao *string
	if produto.Descricao != nil {
		descricao = produto.Descricao
	}

	var codigoBarras *string
	if produto.CodigoBarras != nil {
		codigoBarras = produto.CodigoBarras
	}

	var quantidadeMaxima *string
	if produto.EstoqueMaximo > 0 {
		s := fmt.Sprintf("%d", produto.EstoqueMaximo)
		quantidadeMaxima = &s
	}

	var valorVendaProfissional *string
	if produto.ValorVendaProfissional != nil {
		s := produto.ValorVendaProfissional.StringFixed(2)
		valorVendaProfissional = &s
	}

	var valorEntrada *string
	if produto.ValorEntrada != nil {
		s := produto.ValorEntrada.StringFixed(2)
		valorEntrada = &s
	}

	return dto.ProdutoResponse{
		ID:                     produto.ID.String(),
		TenantID:               produto.TenantID.String(),
		Nome:                   produto.Nome,
		Descricao:              descricao,
		CodigoBarras:           codigoBarras,
		CategoriaProdutoID:     categoriaProdutoIDStr,
		UnidadeMedida:          string(produto.UnidadeMedida),
		ValorUnitario:          produto.Preco.StringFixed(2),
		QuantidadeAtual:        produto.QuantidadeAtual.StringFixed(2),
		QuantidadeMinima:       produto.QuantidadeMinima.StringFixed(2),
		QuantidadeMaxima:       quantidadeMaxima,
		ValorVendaProfissional: valorVendaProfissional,
		ValorEntrada:           valorEntrada,
		FornecedorID:           fornecedorIDStr,
		EstaBaixo:              produto.QuantidadeAtual.LessThanOrEqual(produto.QuantidadeMinima),
		Ativo:                  produto.Ativo,
		CreatedAt:              produto.CriadoEm.Format(time.RFC3339),
		UpdatedAt:              produto.AtualizadoEm.Format(time.RFC3339),
	}
}
