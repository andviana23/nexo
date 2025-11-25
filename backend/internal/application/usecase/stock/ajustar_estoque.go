package stock

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AjustarEstoqueUseCase implementa lógica de negócio para ajuste manual de estoque
type AjustarEstoqueUseCase struct {
	produtoRepo      port.ProdutoRepository
	movimentacaoRepo port.MovimentacaoEstoqueRepository
}

// NewAjustarEstoqueUseCase cria nova instância do use case
func NewAjustarEstoqueUseCase(
	produtoRepo port.ProdutoRepository,
	movimentacaoRepo port.MovimentacaoEstoqueRepository,
) *AjustarEstoqueUseCase {
	return &AjustarEstoqueUseCase{
		produtoRepo:      produtoRepo,
		movimentacaoRepo: movimentacaoRepo,
	}
}

// Execute executa a lógica de ajuste de estoque
func (uc *AjustarEstoqueUseCase) Execute(
	ctx context.Context,
	tenantID uuid.UUID,
	usuarioID uuid.UUID,
	input dto.AjustarEstoqueRequest,
) (*dto.MovimentacaoResponse, error) {
	// 1. Buscar produto
	produtoID, err := uuid.Parse(input.ProdutoID)
	if err != nil {
		return nil, fmt.Errorf("produto_id inválido: %w", err)
	}

	produto, err := uc.produtoRepo.FindByID(ctx, tenantID, produtoID)
	if err != nil {
		return nil, fmt.Errorf("produto não encontrado: %w", err)
	}

	// 2. Validar ownership
	if produto.TenantID != tenantID {
		return nil, fmt.Errorf("produto não pertence ao tenant")
	}

	// 3. Guardar quantidade anterior para calcular diferença
	quantidadeAnterior := produto.QuantidadeAtual

	// 4. Converter nova quantidade string para decimal
	novaQuantidade, err := decimal.NewFromString(input.NovaQuantidade)
	if err != nil {
		return nil, fmt.Errorf("quantidade inválida: %w", err)
	}

	// 5. Ajustar estoque (validação feita na entidade - requer motivo)
	if err := produto.AjustarEstoque(novaQuantidade, input.Motivo); err != nil {
		return nil, fmt.Errorf("erro ao ajustar estoque: %w", err)
	}

	// 6. Atualizar no banco
	if err := uc.produtoRepo.AtualizarQuantidade(ctx, tenantID, produtoID, produto.QuantidadeAtual); err != nil {
		return nil, fmt.Errorf("erro ao atualizar quantidade: %w", err)
	}

	// 7. Calcular diferença para movimentação
	diferenca := novaQuantidade.Sub(quantidadeAnterior)
	quantidadeAbs := diferenca.Abs()

	// 8. Criar movimentação de ajuste
	movimentacao, err := entity.NewMovimentacaoEstoque(
		tenantID,
		produtoID,
		usuarioID,
		entity.MovimentacaoAjuste,
		quantidadeAbs, // decimal.Decimal (sempre positivo)
		decimal.Zero,  // Ajuste não tem valor unitário
		fmt.Sprintf("Ajuste de estoque: %s (Anterior: %s → Nova: %s)",
			input.Motivo,
			quantidadeAnterior.String(),
			novaQuantidade.String()),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar movimentação: %w", err)
	}

	// 9. Persistir movimentação
	if err := uc.movimentacaoRepo.Create(ctx, movimentacao); err != nil {
		// TODO: Rollback (implementar transação)
		return nil, fmt.Errorf("erro ao registrar movimentação: %w", err)
	}

	// 10. Retornar resposta
	return &dto.MovimentacaoResponse{
		ID:            movimentacao.ID.String(),
		TenantID:      movimentacao.TenantID.String(),
		ProdutoID:     movimentacao.ProdutoID.String(),
		ProdutoNome:   produto.Nome,
		UsuarioID:     movimentacao.UsuarioID.String(),
		Tipo:          string(movimentacao.Tipo),
		Quantidade:    diferenca.String(),
		ValorUnitario: "0.00",
		ValorTotal:    "0.00",
		Observacoes:   movimentacao.Observacoes,
		Data:          movimentacao.DataMovimentacao.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt:     movimentacao.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
