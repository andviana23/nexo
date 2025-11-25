package stock

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ListarAlertasEstoqueBaixoUseCase lista produtos com estoque abaixo do mínimo
type ListarAlertasEstoqueBaixoUseCase struct {
	produtoRepo port.ProdutoRepository
}

// NewListarAlertasEstoqueBaixoUseCase cria nova instância do use case
func NewListarAlertasEstoqueBaixoUseCase(produtoRepo port.ProdutoRepository) *ListarAlertasEstoqueBaixoUseCase {
	return &ListarAlertasEstoqueBaixoUseCase{
		produtoRepo: produtoRepo,
	}
}

// Execute executa a lógica de listar alertas de estoque baixo
func (uc *ListarAlertasEstoqueBaixoUseCase) Execute(ctx context.Context, tenantIDStr string) (*dto.ListAlertasEstoqueBaixoResponse, error) {
	// 1. Validar tenant_id
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, fmt.Errorf("tenant_id inválido: %w", err)
	}

	// 2. Buscar produtos abaixo do mínimo
	produtos, err := uc.produtoRepo.ListAbaixoDoMinimo(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar produtos: %w", err)
	}

	// 3. Converter para DTOs
	alertas := make([]dto.AlertaEstoqueBaixo, len(produtos))
	for i, produto := range produtos {
		// Calcular percentual do estoque
		percentual := 0.0
		if !produto.QuantidadeMinima.IsZero() {
			percentual = produto.QuantidadeAtual.Div(produto.QuantidadeMinima).Mul(decimal.NewFromInt(100)).InexactFloat64()
		}

		// Determinar severidade
		severidade := "baixo"
		if percentual <= 25 {
			severidade = "crítico"
		} else if percentual <= 50 {
			severidade = "alerta"
		}

		alertas[i] = dto.AlertaEstoqueBaixo{
			ProdutoID:        produto.ID.String(),
			SKU:              produto.SKU,
			Nome:             produto.Nome,
			Categoria:        string(produto.Categoria),
			QuantidadeAtual:  produto.QuantidadeAtual.String(),
			QuantidadeMinima: produto.QuantidadeMinima.String(),
			UnidadeMedida:    string(produto.UnidadeMedida),
			Percentual:       fmt.Sprintf("%.1f", percentual),
			Severidade:       severidade,
		}
	}

	// 4. Retornar resposta
	return &dto.ListAlertasEstoqueBaixoResponse{
		Total:   len(alertas),
		Alertas: alertas,
	}, nil
}
