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

// RegistrarSaidaUseCase implementa lógica de negócio para saída de estoque
type RegistrarSaidaUseCase struct {
	produtoRepo      port.ProdutoRepository
	movimentacaoRepo port.MovimentacaoEstoqueRepository
}

// NewRegistrarSaidaUseCase cria nova instância do use case
func NewRegistrarSaidaUseCase(
	produtoRepo port.ProdutoRepository,
	movimentacaoRepo port.MovimentacaoEstoqueRepository,
) *RegistrarSaidaUseCase {
	return &RegistrarSaidaUseCase{
		produtoRepo:      produtoRepo,
		movimentacaoRepo: movimentacaoRepo,
	}
}

// Execute executa a lógica de registrar saída de estoque
func (uc *RegistrarSaidaUseCase) Execute(
	ctx context.Context,
	tenantID uuid.UUID,
	usuarioID uuid.UUID,
	input dto.RegistrarSaidaRequest,
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

	// 2. Validar que produto pertence ao tenant
	if produto.TenantID != tenantID {
		return nil, fmt.Errorf("produto não pertence ao tenant")
	}

	// 3. Converter quantidade string para decimal
	quantidade, err := decimal.NewFromString(input.Quantidade)
	if err != nil {
		return nil, fmt.Errorf("quantidade inválida: %w", err)
	}

	// 4. Remover estoque (validação de quantidade é feita na entidade)
	if err := produto.RemoverEstoque(quantidade); err != nil {
		return nil, fmt.Errorf("erro ao remover estoque: %w", err)
	}

	// 5. Atualizar quantidade no banco
	if err := uc.produtoRepo.AtualizarQuantidade(ctx, tenantID, produtoID, produto.QuantidadeAtual); err != nil {
		return nil, fmt.Errorf("erro ao atualizar quantidade: %w", err)
	}

	// 6. Determinar tipo de movimentação
	tipoMov := entity.MovimentacaoSaida
	switch input.Motivo {
	case "CONSUMO_INTERNO":
		tipoMov = entity.MovimentacaoConsumoInterno
	case "PERDA":
		tipoMov = entity.MovimentacaoPerda
	case "DEVOLUCAO":
		tipoMov = entity.MovimentacaoDevolucao
	}

	// 7. Criar movimentação de saída (usando decimal para valor unitário)
	valorUnitario := produto.Preco // Já é decimal.Decimal

	movimentacao, err := entity.NewMovimentacaoEstoque(
		tenantID,
		produtoID,
		usuarioID,
		tipoMov,
		quantidade,    // decimal.Decimal
		valorUnitario, // decimal.Decimal
		input.Observacoes,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar movimentação: %w", err)
	}

	// 8. Persistir movimentação
	if err := uc.movimentacaoRepo.Create(ctx, movimentacao); err != nil {
		// TODO: Rollback da quantidade (implementar transação)
		return nil, fmt.Errorf("erro ao registrar movimentação: %w", err)
	}

	// 9. Retornar resposta
	return &dto.MovimentacaoResponse{
		ID:            movimentacao.ID.String(),
		TenantID:      movimentacao.TenantID.String(),
		ProdutoID:     movimentacao.ProdutoID.String(),
		ProdutoNome:   produto.Nome,
		UsuarioID:     movimentacao.UsuarioID.String(),
		Tipo:          string(movimentacao.Tipo),
		Quantidade:    quantidade.String(),
		ValorUnitario: movimentacao.ValorUnitario.StringFixed(2),
		ValorTotal:    movimentacao.ValorTotal.StringFixed(2),
		Observacoes:   movimentacao.Observacoes,
		Data:          movimentacao.DataMovimentacao.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt:     movimentacao.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
