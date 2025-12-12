// Package command contém os use cases relacionados a comandas
package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	"github.com/andviana23/barber-analytics-backend/internal/domain/repository"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// CancelCommandInput define os dados de entrada para cancelamento de comanda
type CancelCommandInput struct {
	CommandID uuid.UUID
	TenantID  uuid.UUID
	UserID    uuid.UUID
	Motivo    string // Motivo do cancelamento (opcional)
}

// CancelCommandOutput define os dados de saída do cancelamento
type CancelCommandOutput struct {
	Command                   *dto.CommandResponse
	EstoqueRevertido          bool
	QuantidadeItensRevertidos int
	MovimentacoesEstoque      []string // IDs das movimentações de estoque criadas (ENTRADA)
}

// CancelCommandUseCase implementa o cancelamento de comanda com reversão de estoque
// T-EST-003: Reverter estoque ao cancelar comanda
type CancelCommandUseCase struct {
	commandRepo        port.CommandRepository
	produtoRepo        port.ProdutoRepository
	movimentacaoRepo   port.MovimentacaoEstoqueRepository
	commissionItemRepo repository.CommissionItemRepository
	contaReceberRepo   port.ContaReceberRepository
	caixaRepo          port.CaixaDiarioRepository
	mapper             *mapper.CommandMapper
	logger             *zap.Logger
}

// NewCancelCommandUseCase cria nova instância do use case
func NewCancelCommandUseCase(
	commandRepo port.CommandRepository,
	produtoRepo port.ProdutoRepository,
	movimentacaoRepo port.MovimentacaoEstoqueRepository,
	commissionItemRepo repository.CommissionItemRepository,
	contaReceberRepo port.ContaReceberRepository,
	caixaRepo port.CaixaDiarioRepository,
	mapper *mapper.CommandMapper,
	logger *zap.Logger,
) *CancelCommandUseCase {
	return &CancelCommandUseCase{
		commandRepo:        commandRepo,
		produtoRepo:        produtoRepo,
		movimentacaoRepo:   movimentacaoRepo,
		commissionItemRepo: commissionItemRepo,
		contaReceberRepo:   contaReceberRepo,
		caixaRepo:          caixaRepo,
		mapper:             mapper,
		logger:             logger,
	}
}

// Execute cancela uma comanda e reverte o estoque se necessário
func (uc *CancelCommandUseCase) Execute(ctx context.Context, input CancelCommandInput) (*CancelCommandOutput, error) {
	uc.logger.Info("Iniciando cancelamento de comanda",
		zap.String("command_id", input.CommandID.String()),
		zap.String("tenant_id", input.TenantID.String()),
		zap.String("motivo", input.Motivo),
	)

	// 1. Buscar comanda
	command, err := uc.commandRepo.FindByID(ctx, input.CommandID, input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar comanda: %w", err)
	}
	if command == nil {
		return nil, errors.New("comanda não encontrada")
	}

	// 2. Verificar se a comanda pode ser cancelada
	// Permitimos cancelar comandas ABERTAS ou FECHADAS
	if command.Status == entity.CommandStatusCanceled {
		return nil, errors.New("comanda já está cancelada")
	}

	// 3. Verificar se precisa reverter estoque (comanda foi fechada e teve produtos)
	precisaReverterEstoque := command.Status == entity.CommandStatusClosed
	var movimentacoesEstoque []string
	var quantidadeItensRevertidos int

	if precisaReverterEstoque {
		// 4. Reverter estoque para cada item PRODUTO
		for _, item := range command.Items {
			if item.Tipo == entity.CommandItemTypeProduto && item.ItemID != uuid.Nil {
				// Buscar produto
				produto, err := uc.produtoRepo.FindByID(ctx, input.TenantID, item.ItemID)
				if err != nil {
					uc.logger.Warn("Erro ao buscar produto para reversão",
						zap.String("item_id", item.ItemID.String()),
						zap.Error(err),
					)
					continue
				}
				if produto == nil {
					continue
				}

				// Calcular quantidade a reverter
				quantidadeReverter := decimal.NewFromInt(int64(item.Quantidade))

				// Criar observação detalhada
				observacao := fmt.Sprintf("Reversão por cancelamento da comanda %s", command.ID.String())
				if input.Motivo != "" {
					observacao = fmt.Sprintf("%s - Motivo: %s", observacao, input.Motivo)
				}

				// Criar movimentação de DEVOLUÇÃO (reversão por cancelamento)
				// Usar Custo do produto como valor unitário
				valorUnitario := decimal.Zero
				if produto.Custo != nil {
					valorUnitario = *produto.Custo
				}

				movimentacao, err := entity.NewMovimentacaoEstoque(
					input.TenantID,
					item.ItemID,
					input.UserID,
					entity.MovimentacaoDevolucao, // Tipo DEVOLUCAO para reversão
					quantidadeReverter,
					valorUnitario,
					observacao,
				)
				if err != nil {
					uc.logger.Error("Erro ao criar movimentação de reversão",
						zap.String("produto_id", item.ItemID.String()),
						zap.Error(err),
					)
					continue
				}

				// Persistir movimentação
				if err := uc.movimentacaoRepo.Create(ctx, movimentacao); err != nil {
					uc.logger.Error("Erro ao salvar movimentação de reversão",
						zap.String("produto_id", item.ItemID.String()),
						zap.Error(err),
					)
					continue
				}

				// Atualizar quantidade do produto
				novaQuantidade := produto.QuantidadeAtual.Add(quantidadeReverter)
				if err := uc.produtoRepo.AtualizarQuantidade(ctx, input.TenantID, item.ItemID, novaQuantidade); err != nil {
					uc.logger.Error("Erro ao atualizar quantidade do produto",
						zap.String("produto_id", item.ItemID.String()),
						zap.Error(err),
					)
					continue
				}

				movimentacoesEstoque = append(movimentacoesEstoque, movimentacao.ID.String())
				quantidadeItensRevertidos++

				uc.logger.Info("Estoque revertido com sucesso",
					zap.String("produto_id", item.ItemID.String()),
					zap.String("quantidade", quantidadeReverter.String()),
					zap.String("nova_quantidade", novaQuantidade.String()),
				)
			}
		}

		// 5. Deletar comissões geradas (soft delete via status)
		if err := uc.deleteCommissionItems(ctx, command.ID, input.TenantID.String()); err != nil {
			uc.logger.Warn("Erro ao deletar itens de comissão",
				zap.String("command_id", command.ID.String()),
				zap.Error(err),
			)
		}

		// 6. Estornar lançamentos financeiros vinculados à comanda fechada
		if err := uc.estornarFinanceiro(ctx, command, input); err != nil {
			uc.logger.Warn("Erro ao estornar financeiro da comanda",
				zap.String("command_id", command.ID.String()),
				zap.Error(err),
			)
		}
	}

	// 7. Cancelar a comanda
	// Se a comanda estava fechada, precisamos alterar diretamente o status
	if command.Status == entity.CommandStatusClosed {
		command.Status = entity.CommandStatusCanceled
	} else {
		// Se estava aberta, usar método do domínio
		if err := command.Cancel(); err != nil {
			return nil, fmt.Errorf("erro ao cancelar comanda: %w", err)
		}
	}

	// 8. Persistir comanda cancelada
	if err := uc.commandRepo.Update(ctx, command); err != nil {
		return nil, fmt.Errorf("erro ao salvar comanda cancelada: %w", err)
	}

	uc.logger.Info("Comanda cancelada com sucesso",
		zap.String("command_id", command.ID.String()),
		zap.Bool("estoque_revertido", precisaReverterEstoque),
		zap.Int("itens_revertidos", quantidadeItensRevertidos),
	)

	return &CancelCommandOutput{
		Command:                   uc.mapper.ToCommandResponse(command),
		EstoqueRevertido:          precisaReverterEstoque && quantidadeItensRevertidos > 0,
		QuantidadeItensRevertidos: quantidadeItensRevertidos,
		MovimentacoesEstoque:      movimentacoesEstoque,
	}, nil
}

// deleteCommissionItems deleta os itens de comissão vinculados à comanda
func (uc *CancelCommandUseCase) deleteCommissionItems(ctx context.Context, commandID uuid.UUID, tenantID string) error {
	// Buscar itens de comissão vinculados à comanda
	items, err := uc.commissionItemRepo.ListByCommand(ctx, tenantID, commandID)
	if err != nil {
		return fmt.Errorf("erro ao buscar itens de comissão: %w", err)
	}

	// Deletar cada item
	for _, item := range items {
		if err := uc.commissionItemRepo.Delete(ctx, tenantID, item.ID); err != nil {
			uc.logger.Warn("Erro ao deletar item de comissão",
				zap.String("item_id", item.ID),
				zap.Error(err),
			)
		}
	}

	return nil
}

// estornarFinanceiro cancela/estorna contas a receber e cria operação inversa no caixa.
func (uc *CancelCommandUseCase) estornarFinanceiro(ctx context.Context, command *entity.Command, input CancelCommandInput) error {
	// 1) Estornar contas a receber vinculadas
	if uc.contaReceberRepo != nil {
		contas, err := uc.contaReceberRepo.ListByCommandID(ctx, input.TenantID.String(), command.ID.String())
		if err != nil {
			return err
		}

		for _, conta := range contas {
			if conta.Status == valueobject.StatusContaEstornado || conta.Status == valueobject.StatusContaCancelado {
				continue
			}
			conta.Status = valueobject.StatusContaEstornado
			conta.ValorAberto = valueobject.Zero()
			conta.AtualizadoEm = time.Now()

			obs := fmt.Sprintf("Estornada por cancelamento da comanda %s", command.ID.String()[:8])
			if input.Motivo != "" {
				obs = fmt.Sprintf("%s - Motivo: %s", obs, input.Motivo)
			}
			conta.AddObservacao(obs)

			if err := uc.contaReceberRepo.Update(ctx, conta); err != nil {
				uc.logger.Warn("erro ao estornar conta a receber",
					zap.String("conta_receber_id", conta.ID),
					zap.Error(err),
				)
			}
		}
	}

	// 2) Criar operação inversa no caixa (sangria) e ajustar totais
	if uc.caixaRepo != nil {
		caixaAberto, err := uc.caixaRepo.FindAberto(ctx, input.TenantID)
		if err != nil || caixaAberto == nil {
			uc.logger.Warn("não foi possível estornar no caixa (nenhum caixa aberto)",
				zap.String("tenant_id", input.TenantID.String()),
				zap.String("command_id", command.ID.String()),
			)
			return nil
		}

		for _, payment := range command.Payments {
			valor := decimal.NewFromFloat(payment.ValorRecebido)
			if valor.IsZero() || valor.IsNegative() {
				continue
			}

			descricao := fmt.Sprintf("Estorno Comanda #%s", command.ID.String()[:8])
			if input.Motivo != "" {
				descricao = fmt.Sprintf("%s - %s", descricao, input.Motivo)
			}

			operacao, err := entity.NewOperacaoSangria(
				caixaAberto.ID,
				input.TenantID,
				input.UserID,
				valor,
				string(entity.DestinoOutros),
				descricao,
				nil,
			)
			if err != nil {
				uc.logger.Warn("erro ao criar operação de estorno no caixa", zap.Error(err))
				continue
			}

			if err := uc.caixaRepo.CreateOperacao(ctx, operacao); err != nil {
				uc.logger.Warn("erro ao registrar estorno no caixa", zap.Error(err))
				continue
			}

			caixaAberto.TotalSangrias = caixaAberto.TotalSangrias.Add(valor)
			if err := uc.caixaRepo.UpdateTotais(ctx, caixaAberto.ID, input.TenantID, caixaAberto.TotalSangrias, caixaAberto.TotalReforcos, caixaAberto.TotalEntradas); err != nil {
				uc.logger.Warn("erro ao atualizar totais do caixa após estorno", zap.Error(err))
			}
		}
	}

	return nil
}
