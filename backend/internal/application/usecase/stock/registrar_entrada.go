package stock

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
)

// RegistrarEntradaInput representa a entrada do use case
type RegistrarEntradaInput struct {
	TenantID        uuid.UUID
	UsuarioID       uuid.UUID
	FornecedorID    uuid.UUID
	DataEntrada     time.Time
	Itens           []ItemEntrada
	Observacoes     string
	GerarFinanceiro bool
}

// ItemEntrada representa um item de entrada
type ItemEntrada struct {
	ProdutoID     uuid.UUID
	Quantidade    int
	ValorUnitario decimal.Decimal
}

// RegistrarEntradaOutput representa a saída do use case
type RegistrarEntradaOutput struct {
	MovimentacoesIDs []uuid.UUID
	ValorTotal       decimal.Decimal
	ItensProcessados int
}

// RegistrarEntradaUseCase implementa o caso de uso de registrar entrada de estoque
type RegistrarEntradaUseCase struct {
	produtoRepo      port.ProdutoRepository
	movimentacaoRepo port.MovimentacaoEstoqueRepository
	fornecedorRepo   port.FornecedorRepository
}

// NewRegistrarEntradaUseCase cria uma nova instância do use case
func NewRegistrarEntradaUseCase(
	produtoRepo port.ProdutoRepository,
	movimentacaoRepo port.MovimentacaoEstoqueRepository,
	fornecedorRepo port.FornecedorRepository,
) *RegistrarEntradaUseCase {
	return &RegistrarEntradaUseCase{
		produtoRepo:      produtoRepo,
		movimentacaoRepo: movimentacaoRepo,
		fornecedorRepo:   fornecedorRepo,
	}
}

// Execute executa o caso de uso
func (uc *RegistrarEntradaUseCase) Execute(
	ctx context.Context,
	input RegistrarEntradaInput,
) (*RegistrarEntradaOutput, error) {
	// 1. Validar fornecedor existe e pertence ao tenant
	fornecedor, err := uc.fornecedorRepo.FindByID(ctx, input.TenantID, input.FornecedorID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar fornecedor: %w", err)
	}
	if fornecedor == nil {
		return nil, fmt.Errorf("fornecedor não encontrado")
	}
	if !fornecedor.Ativo {
		return nil, fmt.Errorf("fornecedor está inativo")
	}

	// 2. Processar cada item da entrada
	var movimentacoesIDs []uuid.UUID
	valorTotalGeral := decimal.Zero

	for _, item := range input.Itens {
		// 2.1 Buscar produto
		produto, err := uc.produtoRepo.FindByID(ctx, input.TenantID, item.ProdutoID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar produto %s: %w", item.ProdutoID, err)
		}
		if produto == nil {
			return nil, fmt.Errorf("produto %s não encontrado", item.ProdutoID)
		}

		// 2.2 Adicionar quantidade ao estoque
		quantidadeDecimal := decimal.NewFromInt(int64(item.Quantidade))
		if err := produto.AdicionarEstoque(quantidadeDecimal); err != nil {
			return nil, fmt.Errorf("erro ao adicionar estoque do produto %s: %w", produto.Nome, err)
		}

		// 2.3 Atualizar produto no banco
		if err := uc.produtoRepo.Update(ctx, produto); err != nil {
			return nil, fmt.Errorf("erro ao atualizar produto %s: %w", produto.Nome, err)
		}

		// 2.4 Valor unitário já está em decimal - não precisa converter para centavos
		valorUnitarioDecimal := item.ValorUnitario

		// 2.5 Criar movimentação
		movimentacao, err := entity.NewMovimentacaoEstoque(
			input.TenantID,
			produto.ID,
			input.UsuarioID,
			entity.MovimentacaoEntrada,
			quantidadeDecimal,
			valorUnitarioDecimal,
			input.Observacoes,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar movimentação para produto %s: %w", produto.Nome, err)
		}

		// 2.6 Associar fornecedor à movimentação
		movimentacao.DefinirFornecedor(input.FornecedorID)

		// 2.7 Persistir movimentação
		if err := uc.movimentacaoRepo.Create(ctx, movimentacao); err != nil {
			return nil, fmt.Errorf("erro ao salvar movimentação para produto %s: %w", produto.Nome, err)
		}

		movimentacoesIDs = append(movimentacoesIDs, movimentacao.ID)

		// 2.8 Calcular valor total
		valorItem := item.ValorUnitario.Mul(decimal.NewFromInt(int64(item.Quantidade)))
		valorTotalGeral = valorTotalGeral.Add(valorItem)
	}

	// 3. TODO: Se GerarFinanceiro = true, criar conta a pagar
	// Integração futura com módulo financeiro
	if input.GerarFinanceiro {
		// Será implementado quando integrar com módulo financeiro
		// createPayableInput := ...
		// payableUseCase.Execute(ctx, createPayableInput)
	}

	return &RegistrarEntradaOutput{
		MovimentacoesIDs: movimentacoesIDs,
		ValorTotal:       valorTotalGeral,
		ItensProcessados: len(input.Itens),
	}, nil
}
