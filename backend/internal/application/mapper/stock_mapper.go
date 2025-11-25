// Package mapper contém funções para conversão entre entidades de domínio e DTOs HTTP para Stock.
package mapper

import (
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/stock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

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
		ValorTotal:       output.ValorTotal.String(),
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
