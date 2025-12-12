// Package command contém os use cases relacionados a comandas
package command

import (
	"context"
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

// FinalizarComandaIntegradaInput define os dados de entrada para finalização integrada
type FinalizarComandaIntegradaInput struct {
	CommandID          uuid.UUID
	TenantID           uuid.UUID
	UserID             uuid.UUID
	DeixarTrocoGorjeta *bool
	DeixarSaldoDivida  *bool
	Observacoes        *string
}

// FinalizarComandaIntegradaOutput define os dados de saída
type FinalizarComandaIntegradaOutput struct {
	Command              *dto.CommandResponse
	ContasReceber        []string // IDs das contas a receber criadas
	OperacoesCaixa       []string // IDs das operações de caixa criadas
	CommissionItems      []string // IDs dos itens de comissão criados
	MovimentacoesEstoque []string // IDs das movimentações de estoque criadas
	TotalLancadoCaixa    decimal.Decimal
	TotalContasReceber   decimal.Decimal
	TotalComissoes       decimal.Decimal
}

// FinalizarComandaIntegradaUseCase implementa a finalização integrada de comanda
// que gera lançamentos financeiros, registra no caixa, abate estoque e gera comissões
type FinalizarComandaIntegradaUseCase struct {
	commandRepo        port.CommandRepository
	appointmentRepo    port.AppointmentRepository
	meioPagamentoRepo  port.MeioPagamentoRepository
	contaReceberRepo   port.ContaReceberRepository
	compensacaoRepo    port.CompensacaoBancariaRepository
	caixaRepo          port.CaixaDiarioRepository
	produtoRepo        port.ProdutoRepository
	movimentacaoRepo   port.MovimentacaoEstoqueRepository
	commissionItemRepo repository.CommissionItemRepository
	commissionRuleRepo repository.CommissionRuleRepository
	// COM-001: Dependências para hierarquia de regras de comissão
	serviceReader      port.ServiceReader
	professionalReader port.ProfessionalReader
	mapper             *mapper.CommandMapper
	logger             *zap.Logger
}

// NewFinalizarComandaIntegradaUseCase cria nova instância do use case
func NewFinalizarComandaIntegradaUseCase(
	commandRepo port.CommandRepository,
	appointmentRepo port.AppointmentRepository,
	meioPagamentoRepo port.MeioPagamentoRepository,
	contaReceberRepo port.ContaReceberRepository,
	compensacaoRepo port.CompensacaoBancariaRepository,
	caixaRepo port.CaixaDiarioRepository,
	produtoRepo port.ProdutoRepository,
	movimentacaoRepo port.MovimentacaoEstoqueRepository,
	commissionItemRepo repository.CommissionItemRepository,
	commissionRuleRepo repository.CommissionRuleRepository,
	// COM-001: Novos readers para hierarquia de comissões
	serviceReader port.ServiceReader,
	professionalReader port.ProfessionalReader,
	mapper *mapper.CommandMapper,
	logger *zap.Logger,
) *FinalizarComandaIntegradaUseCase {
	return &FinalizarComandaIntegradaUseCase{
		commandRepo:        commandRepo,
		appointmentRepo:    appointmentRepo,
		meioPagamentoRepo:  meioPagamentoRepo,
		contaReceberRepo:   contaReceberRepo,
		compensacaoRepo:    compensacaoRepo,
		caixaRepo:          caixaRepo,
		produtoRepo:        produtoRepo,
		movimentacaoRepo:   movimentacaoRepo,
		commissionItemRepo: commissionItemRepo,
		commissionRuleRepo: commissionRuleRepo,
		serviceReader:      serviceReader,
		professionalReader: professionalReader,
		mapper:             mapper,
		logger:             logger,
	}
}

// Execute executa a finalização integrada da comanda
// 1. Valida se há caixa aberto (obrigatório para registro financeiro)
// 2. Valida pagamentos e itens
// 3. Para cada item PRODUTO: abate estoque (MovimentacaoEstoque tipo SAIDA)
// 4. Para cada item SERVICO: cria CommissionItem para o profissional
// 5. Para cada pagamento:
//   - DINHEIRO/PIX: cria OperacaoCaixa
//   - Outros: cria ContaReceber com D+ do meio de pagamento
//
// 6. Fecha a comanda
// 7. Atualiza o agendamento para DONE (se vinculado)
func (uc *FinalizarComandaIntegradaUseCase) Execute(ctx context.Context, input FinalizarComandaIntegradaInput) (*FinalizarComandaIntegradaOutput, error) {
	// G-003: Verificar se há caixa aberto ANTES de qualquer operação
	// Comanda só pode ser fechada com caixa aberto para garantir integridade financeira
	caixaAberto, err := uc.caixaRepo.FindAberto(ctx, input.TenantID)
	if err != nil {
		uc.logger.Error("erro ao verificar caixa aberto", zap.Error(err))
		return nil, fmt.Errorf("erro ao verificar caixa: %w", err)
	}
	if caixaAberto == nil {
		uc.logger.Warn("tentativa de fechar comanda sem caixa aberto",
			zap.String("tenant_id", input.TenantID.String()),
			zap.String("command_id", input.CommandID.String()))
		return nil, fmt.Errorf("não é possível fechar a comanda: caixa não está aberto. Abra o caixa antes de finalizar vendas")
	}

	// Buscar comanda com itens e pagamentos
	command, err := uc.commandRepo.FindByID(ctx, input.CommandID, input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar comanda: %w", err)
	}
	if command == nil {
		return nil, fmt.Errorf("comanda não encontrada")
	}

	// Aplicar opções de fechamento
	if input.DeixarTrocoGorjeta != nil {
		command.DeixarTrocoGorjeta = *input.DeixarTrocoGorjeta
	}
	if input.DeixarSaldoDivida != nil {
		command.DeixarSaldoDivida = *input.DeixarSaldoDivida
	}
	if input.Observacoes != nil {
		command.Observacoes = input.Observacoes
	}

	// Validar se pode fechar
	if err := command.CanClose(); err != nil {
		return nil, fmt.Errorf("não é possível fechar a comanda: %w", err)
	}

	output := &FinalizarComandaIntegradaOutput{
		ContasReceber:        make([]string, 0),
		OperacoesCaixa:       make([]string, 0),
		CommissionItems:      make([]string, 0),
		MovimentacoesEstoque: make([]string, 0),
		TotalLancadoCaixa:    decimal.Zero,
		TotalContasReceber:   decimal.Zero,
		TotalComissoes:       decimal.Zero,
	}

	// Totais por tipo de item para rateio das receitas por origem
	var totalServicos decimal.Decimal
	var totalProdutos decimal.Decimal

	// Buscar profissional do appointment (para comissões)
	var professionalID string
	var professionalInfo *port.ProfessionalInfo
	var appointmentDate *time.Time // Data do agendamento para reference_date
	if command.AppointmentID != nil {
		appointment, err := uc.appointmentRepo.FindByID(ctx, input.TenantID.String(), "", command.AppointmentID.String())
		if err == nil && appointment != nil {
			professionalID = appointment.ProfessionalID
			// COM-001: Buscar info completa do profissional (inclui comissão)
			professionalInfo, _ = uc.professionalReader.FindByID(ctx, input.TenantID.String(), professionalID)
			// COM-004: Usar data do agendamento como reference_date
			if !appointment.StartTime.IsZero() {
				t := appointment.StartTime
				appointmentDate = &t
			}
		}
	}

	// Obter unit_id (TODO: quando Command tiver unit_id, buscar de lá)
	// Por enquanto usa nil (sem unidade) - vai direto para regra global
	var unitID *string

	// T-EST-002 & T-COM-001: Processar cada item da comanda
	for _, item := range command.Items {
		switch item.Tipo {
		case entity.CommandItemTypeProduto:
			totalProdutos = totalProdutos.Add(decimal.NewFromFloat(item.PrecoFinal))
			// T-EST-002: Abater estoque para produtos
			if err := uc.processarEstoqueProduto(ctx, input.TenantID, input.UserID, &item, output); err != nil {
				uc.logger.Warn("erro ao abater estoque do produto",
					zap.String("item_id", item.ID.String()),
					zap.Error(err))
				// Continua mesmo com erro - log mas não bloqueia fechamento
			}

			// COM-005: Gerar comissão para produtos (usa comissão do profissional ou regra global)
			if professionalID != "" {
				ruleResult := uc.buscarRegraComissaoProduto(
					ctx,
					input.TenantID,
					unitID,
					professionalID,
					professionalInfo,
				)

				if ruleResult != nil {
					if err := uc.processarComissaoProduto(
						ctx, input.TenantID, professionalID, command, &item,
						ruleResult, appointmentDate, output,
					); err != nil {
						uc.logger.Warn("erro ao gerar comissão do produto",
							zap.String("item_id", item.ID.String()),
							zap.String("source", ruleResult.Source),
							zap.Error(err))
					}
				} else {
					uc.logger.Debug("nenhuma regra de comissão encontrada para o produto",
						zap.String("item_id", item.ID.String()),
						zap.String("produto_id", item.ItemID.String()))
				}
			}

		case entity.CommandItemTypeServico:
			totalServicos = totalServicos.Add(decimal.NewFromFloat(item.PrecoFinal))
			// COM-001: Buscar regra usando hierarquia de 5 níveis
			if professionalID != "" {
				ruleResult := uc.buscarRegraComissaoHierarquica(
					ctx,
					input.TenantID,
					unitID,
					professionalID,
					professionalInfo,
					&item,
				)

				if ruleResult != nil {
					if err := uc.processarComissaoServicoHierarquicaComData(
						ctx, input.TenantID, professionalID, command, &item,
						ruleResult, appointmentDate, output,
					); err != nil {
						uc.logger.Warn("erro ao gerar comissão do serviço",
							zap.String("item_id", item.ID.String()),
							zap.String("source", ruleResult.Source),
							zap.String("calculation_base", ruleResult.CalculationBase),
							zap.Error(err))
						// Continua mesmo com erro - log mas não bloqueia fechamento
					}
				} else {
					uc.logger.Warn("nenhuma regra de comissão encontrada para o serviço",
						zap.String("item_id", item.ID.String()),
						zap.String("servico_id", item.ItemID.String()),
						zap.String("profissional_id", professionalID))
				}
			}
		}
	}

	// Calcular proporção de serviços vs produtos para rateio das receitas por pagamento
	totalItens := totalServicos.Add(totalProdutos)
	ratioServicos := decimal.NewFromInt(1)
	if !totalItens.IsZero() {
		if totalServicos.IsZero() {
			ratioServicos = decimal.Zero
		} else if totalProdutos.IsZero() {
			ratioServicos = decimal.NewFromInt(1)
		} else {
			ratioServicos = totalServicos.Div(totalItens)
		}
	}

	// Mapear contas já criadas por pagamento+origem para idempotência
	existingByPaymentOrigem := make(map[string]struct{})
	if uc.contaReceberRepo != nil {
		existing, err := uc.contaReceberRepo.ListByCommandID(ctx, input.TenantID.String(), command.ID.String())
		if err == nil {
			for _, c := range existing {
				if c.CommandPaymentID != nil {
					key := *c.CommandPaymentID + "|" + c.Origem
					existingByPaymentOrigem[key] = struct{}{}
				}
			}
		}
	}

	// Processar cada pagamento (caixaAberto já foi validado no início)
	// TODOS os pagamentos são registrados no caixa do dia, independente do tipo
	for _, payment := range command.Payments {
		// Buscar meio de pagamento para obter tipo e D+
		meioPagamento, err := uc.meioPagamentoRepo.FindByID(ctx, input.TenantID.String(), payment.MeioPagamentoID.String())
		if err != nil {
			uc.logger.Warn("meio de pagamento não encontrado, usando padrões",
				zap.String("meio_pagamento_id", payment.MeioPagamentoID.String()),
				zap.Error(err))
			// Continua com valores padrão
			meioPagamento = &entity.MeioPagamento{
				Nome:  "Pagamento",
				Tipo:  entity.TipoPagamentoOutro,
				DMais: 0,
			}
		}

		valorRecebido := decimal.NewFromFloat(payment.ValorRecebido)

		// Usar o nome cadastrado do meio de pagamento na descrição
		nomeDescricao := meioPagamento.Nome
		if nomeDescricao == "" {
			nomeDescricao = string(meioPagamento.Tipo)
		}

		// TODOS os pagamentos são registrados no caixa diário
		operacao, err := entity.NewOperacaoVenda(
			caixaAberto.ID,
			input.TenantID,
			input.UserID,
			valorRecebido,
			fmt.Sprintf("Comanda #%s - %s", command.ID.String()[:8], nomeDescricao),
		)
		if err != nil {
			uc.logger.Error("erro ao criar operação de venda", zap.Error(err))
			continue
		}

		if err := uc.caixaRepo.CreateOperacao(ctx, operacao); err != nil {
			uc.logger.Error("erro ao registrar operação no caixa", zap.Error(err))
			continue
		}

		// Atualizar totais do caixa
		novoTotalEntradas := caixaAberto.TotalEntradas.Add(valorRecebido)
		caixaAberto.TotalEntradas = novoTotalEntradas // Atualizar local para próximos pagamentos
		if err := uc.caixaRepo.UpdateTotais(ctx, caixaAberto.ID, input.TenantID, caixaAberto.TotalSangrias, caixaAberto.TotalReforcos, novoTotalEntradas); err != nil {
			uc.logger.Error("erro ao atualizar totais do caixa", zap.Error(err))
		}

		output.OperacoesCaixa = append(output.OperacoesCaixa, operacao.ID.String())
		output.TotalLancadoCaixa = output.TotalLancadoCaixa.Add(valorRecebido)

		uc.logger.Info("operação de venda registrada no caixa",
			zap.String("operacao_id", operacao.ID.String()),
			zap.String("meio_pagamento", nomeDescricao),
			zap.String("tipo", string(meioPagamento.Tipo)),
			zap.String("valor", valorRecebido.String()))

		// Gerar contas a receber por origem (SERVICO/PRODUTO) para este pagamento.
		// - Pagamentos D+0: status RECEBIDO no ato
		// - Pagamentos D+N: status CONFIRMADO aguardando compensação bancária
		if uc.contaReceberRepo != nil {
			now := time.Now()
			competencia := now.Format("2006-01")

			// Rateio proporcional do valor deste pagamento entre serviços e produtos
			valorServicos := valorRecebido.Mul(ratioServicos).Round(2)
			valorProdutos := valorRecebido.Sub(valorServicos)

			splits := []struct {
				origem string
				valor  decimal.Decimal
			}{
				{origem: "SERVICO", valor: valorServicos},
				{origem: "PRODUTO", valor: valorProdutos},
			}

			for _, sp := range splits {
				if sp.valor.IsZero() || sp.valor.IsNegative() {
					continue
				}

				paymentIDStr := payment.ID.String()
				key := paymentIDStr + "|" + sp.origem
				if _, exists := existingByPaymentOrigem[key]; exists {
					continue
				}

				valorMoney := valueobject.NewMoneyFromDecimal(sp.valor)
				dataVencimento := now
				status := valueobject.StatusContaRecebido
				var dataRecebimento *time.Time

				if meioPagamento.DMais > 0 {
					dataVencimento = now.AddDate(0, 0, meioPagamento.DMais)
					status = valueobject.StatusContaConfirmado
				} else {
					dataRecebimento = &now
				}

				descricaoConta := fmt.Sprintf("Comanda #%s - %s (%s)", command.ID.String()[:8], nomeDescricao, sp.origem)
				contaReceber, err := entity.NewContaReceber(
					input.TenantID,
					sp.origem,
					nil,
					descricaoConta,
					valorMoney,
					dataVencimento,
				)
				if err != nil {
					uc.logger.Warn("erro ao criar conta a receber da comanda", zap.Error(err))
					continue
				}

				commandIDStr := command.ID.String()
				contaReceber.CommandID = &commandIDStr
				contaReceber.CommandPaymentID = &paymentIDStr
				contaReceber.CompetenciaMes = &competencia
				contaReceber.Status = status

				if status == valueobject.StatusContaRecebido {
					contaReceber.ValorPago = valorMoney
					contaReceber.ValorAberto = valueobject.Zero()
					contaReceber.DataRecebimento = dataRecebimento
					contaReceber.ReceivedAt = dataRecebimento
				} else {
					contaReceber.ValorPago = valueobject.Zero()
					contaReceber.ValorAberto = valorMoney
					contaReceber.ConfirmedAt = &now
				}

				if err := uc.contaReceberRepo.Create(ctx, contaReceber); err != nil {
					uc.logger.Warn("erro ao persistir conta a receber da comanda", zap.Error(err))
					continue
				}

				existingByPaymentOrigem[key] = struct{}{}
				output.ContasReceber = append(output.ContasReceber, contaReceber.ID)
				output.TotalContasReceber = output.TotalContasReceber.Add(sp.valor)

				uc.logger.Info("conta a receber criada para comanda",
					zap.String("conta_receber_id", contaReceber.ID),
					zap.String("origem", sp.origem),
					zap.String("valor", sp.valor.String()),
					zap.Int("d_mais", meioPagamento.DMais))

				// Gerar compensação bancária automática para pagamentos D+
				if uc.compensacaoRepo != nil && meioPagamento.DMais > 0 {
					// Idempotência: evitar criar compensação duplicada
					if existingComp, err := uc.compensacaoRepo.FindByReceitaID(ctx, input.TenantID.String(), contaReceber.ID); err == nil && existingComp != nil {
						continue
					}

					taxaPercentual := valueobject.NewPercentageUnsafe(meioPagamento.Taxa)

					taxaFixaDec := meioPagamento.TaxaFixa
					if !taxaFixaDec.IsZero() && !valorRecebido.IsZero() {
						ratio := sp.valor.Div(valorRecebido)
						taxaFixaDec = meioPagamento.TaxaFixa.Mul(ratio).Round(2)
					}

					taxaFixaMoney := valueobject.NewMoneyFromDecimal(taxaFixaDec)
					dMaisVO := valueobject.NewDMaisUnsafe(meioPagamento.DMais)

					comp, err := entity.NewCompensacaoBancaria(
						input.TenantID,
						contaReceber.ID,
						payment.MeioPagamentoID.String(),
						now,
						valorMoney,
						taxaPercentual,
						taxaFixaMoney,
						dMaisVO,
					)
					if err != nil {
						uc.logger.Warn("erro ao criar compensação bancária automática", zap.Error(err))
						continue
					}

					_ = comp.MarcarComoConfirmado()
					if err := uc.compensacaoRepo.Create(ctx, comp); err != nil {
						uc.logger.Warn("erro ao persistir compensação bancária automática", zap.Error(err))
						continue
					}

					uc.logger.Info("compensação bancária criada para pagamento D+",
						zap.String("compensacao_id", comp.ID),
						zap.String("receita_id", comp.ReceitaID),
						zap.Int("d_mais", meioPagamento.DMais),
						zap.String("data_compensacao", comp.DataCompensacao.Format("2006-01-02")))
				}
			}
		}
	}

	// Fechar comanda (domain logic)
	if err := command.Close(input.UserID); err != nil {
		return nil, fmt.Errorf("falha ao fechar comanda: %w", err)
	}

	// Persistir atualização da comanda
	if err := uc.commandRepo.Update(ctx, command); err != nil {
		return nil, fmt.Errorf("falha ao atualizar comanda: %w", err)
	}

	// Atualizar status do appointment para DONE (se houver)
	if command.AppointmentID != nil {
		appointment, err := uc.appointmentRepo.FindByID(ctx, input.TenantID.String(), "", command.AppointmentID.String())
		if err == nil && appointment != nil {
			appointment.Status = valueobject.AppointmentStatusDone
			if err := uc.appointmentRepo.Update(ctx, appointment); err != nil {
				uc.logger.Warn("falha ao atualizar status do agendamento",
					zap.String("appointment_id", command.AppointmentID.String()),
					zap.Error(err))
			}
		}
	}

	// Buscar comanda atualizada para retorno
	closedCommand, err := uc.commandRepo.FindByID(ctx, input.CommandID, input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar comanda fechada: %w", err)
	}

	output.Command = uc.mapper.ToCommandResponse(closedCommand)

	uc.logger.Info("comanda finalizada com integração financeira completa",
		zap.String("command_id", command.ID.String()),
		zap.Int("contas_receber_criadas", len(output.ContasReceber)),
		zap.Int("operacoes_caixa_criadas", len(output.OperacoesCaixa)),
		zap.Int("comissoes_criadas", len(output.CommissionItems)),
		zap.Int("movimentacoes_estoque", len(output.MovimentacoesEstoque)),
		zap.String("total_caixa", output.TotalLancadoCaixa.String()),
		zap.String("total_contas_receber", output.TotalContasReceber.String()),
		zap.String("total_comissoes", output.TotalComissoes.String()))

	return output, nil
}

// processarEstoqueProduto abate estoque para um item do tipo PRODUTO
// T-EST-002: Abater estoque ao fechar comanda
func (uc *FinalizarComandaIntegradaUseCase) processarEstoqueProduto(
	ctx context.Context,
	tenantID, userID uuid.UUID,
	item *entity.CommandItem,
	output *FinalizarComandaIntegradaOutput,
) error {
	// Buscar produto para obter dados atuais
	produto, err := uc.produtoRepo.FindByID(ctx, tenantID, item.ItemID)
	if err != nil {
		return fmt.Errorf("produto não encontrado: %w", err)
	}
	if produto == nil {
		return fmt.Errorf("produto %s não encontrado", item.ItemID.String())
	}

	quantidade := decimal.NewFromInt(int64(item.Quantidade))
	valorUnitario := decimal.NewFromFloat(item.PrecoUnitario)

	// Criar movimentação de saída
	movimentacao, err := entity.NewMovimentacaoEstoque(
		tenantID,
		item.ItemID, // ProdutoID
		userID,
		entity.MovimentacaoSaida,
		quantidade,
		valorUnitario,
		fmt.Sprintf("Venda - Comanda item %s", item.ID.String()[:8]),
	)
	if err != nil {
		return fmt.Errorf("erro ao criar movimentação: %w", err)
	}

	// Persistir movimentação
	if err := uc.movimentacaoRepo.Create(ctx, movimentacao); err != nil {
		return fmt.Errorf("erro ao persistir movimentação: %w", err)
	}

	// Atualizar quantidade do produto (abater estoque)
	novaQuantidade := produto.QuantidadeAtual.Sub(quantidade)
	if novaQuantidade.IsNegative() {
		uc.logger.Warn("estoque ficará negativo após venda",
			zap.String("produto_id", produto.ID.String()),
			zap.String("quantidade_atual", produto.QuantidadeAtual.String()),
			zap.String("quantidade_vendida", quantidade.String()))
		// Permite venda mesmo com estoque negativo (warn apenas)
	}

	if err := uc.produtoRepo.AtualizarQuantidade(ctx, tenantID, item.ItemID, novaQuantidade); err != nil {
		return fmt.Errorf("erro ao atualizar quantidade do produto: %w", err)
	}

	output.MovimentacoesEstoque = append(output.MovimentacoesEstoque, movimentacao.ID.String())

	uc.logger.Info("estoque abatido",
		zap.String("produto_id", item.ItemID.String()),
		zap.String("quantidade", quantidade.String()),
		zap.String("movimentacao_id", movimentacao.ID.String()))

	return nil
}

// processarComissaoServico cria CommissionItem para um item do tipo SERVICO
// T-COM-001: Gerar commission_items ao fechar comanda
func (uc *FinalizarComandaIntegradaUseCase) processarComissaoServico(
	ctx context.Context,
	tenantID uuid.UUID,
	professionalID string,
	command *entity.Command,
	item *entity.CommandItem,
	rule *entity.CommissionRule,
	output *FinalizarComandaIntegradaOutput,
) error {
	// Calcular valor bruto do serviço
	grossValue := decimal.NewFromFloat(item.PrecoFinal)

	// Criar item de comissão
	commissionItem, err := entity.NewCommissionItem(
		tenantID,
		professionalID,
		grossValue,
		rule.DefaultRate,
		rule.Type, // PERCENTUAL ou FIXO
		"SERVICO", // CommissionSource
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("erro ao criar item de comissão: %w", err)
	}

	// Associar à comanda e item
	commandIDStr := command.ID.String()
	itemIDStr := item.ID.String()
	commissionItem.CommandID = &commandIDStr
	commissionItem.CommandItemID = &itemIDStr
	commissionItem.ServiceID = func() *string { s := item.ItemID.String(); return &s }()
	commissionItem.ServiceName = &item.Descricao
	ruleIDStr := rule.ID
	commissionItem.RuleID = &ruleIDStr

	// Persistir
	created, err := uc.commissionItemRepo.Create(ctx, commissionItem)
	if err != nil {
		return fmt.Errorf("erro ao persistir item de comissão: %w", err)
	}

	output.CommissionItems = append(output.CommissionItems, created.ID)
	output.TotalComissoes = output.TotalComissoes.Add(created.CommissionValue)

	uc.logger.Info("comissão gerada",
		zap.String("commission_item_id", created.ID),
		zap.String("professional_id", professionalID),
		zap.String("valor_bruto", grossValue.String()),
		zap.String("comissao", created.CommissionValue.String()),
		zap.String("tipo", rule.Type))

	return nil
}

// =============================================================================
// COM-001: Hierarquia de Regras de Comissão (5 níveis)
// =============================================================================

// CommissionRuleResult encapsula o resultado da busca de regra de comissão
type CommissionRuleResult struct {
	Rule            *entity.CommissionRule
	Source          string // SERVICO, CATEGORIA, PROFISSIONAL, REGRA
	CalculationBase string // BRUTO, LIQUIDO
}

// buscarRegraComissaoHierarquica implementa a hierarquia de 5 níveis para regras de comissão
// Prioridade: 1) Serviço → 2) Categoria → 3) Profissional → 4) Unidade → 5) Tenant Global
// Retorna a regra encontrada, source e base de cálculo
func (uc *FinalizarComandaIntegradaUseCase) buscarRegraComissaoHierarquica(
	ctx context.Context,
	tenantID uuid.UUID,
	unitID *string,
	professionalID string,
	professionalInfo *port.ProfessionalInfo,
	item *entity.CommandItem,
) *CommissionRuleResult {
	tenantIDStr := tenantID.String()

	// Buscar info do serviço (necessário para comissão específica e categoria)
	var servico *port.ServiceInfo
	if item.ItemID != uuid.Nil {
		servico, _ = uc.serviceReader.FindByID(ctx, tenantIDStr, item.ItemID.String())
	}

	// 1º NÍVEL: Serviço tem comissão específica?
	if servico != nil && servico.Comissao != nil && *servico.Comissao != "" {
		// Converter comissão do serviço para decimal
		comissaoServico, err := decimal.NewFromString(*servico.Comissao)
		if err == nil && comissaoServico.GreaterThan(decimal.Zero) {
			uc.logger.Debug("usando comissão do serviço (nível 1)",
				zap.String("servico_id", servico.ID),
				zap.String("comissao", comissaoServico.String()))
			return &CommissionRuleResult{
				Rule: &entity.CommissionRule{
					ID:          "servico-" + servico.ID, // ID virtual
					TenantID:    tenantID,
					Name:        "Comissão do Serviço: " + servico.Name,
					Type:        "PERCENTUAL", // Comissão de serviço é sempre percentual
					DefaultRate: comissaoServico,
					IsActive:    true,
				},
				Source:          "SERVICO",
				CalculationBase: "BRUTO", // Comissão de serviço usa valor bruto por padrão
			}
		}
	}

	// 2º NÍVEL: Profissional tem comissão específica para a CATEGORIA do serviço?
	if servico != nil && servico.CategoriaID != nil && *servico.CategoriaID != "" && professionalID != "" {
		comissaoCategoria, err := uc.professionalReader.GetCategoryCommission(ctx, tenantIDStr, professionalID, *servico.CategoriaID)
		if err == nil && comissaoCategoria != nil && *comissaoCategoria != "" {
			comissaoCat, err := decimal.NewFromString(*comissaoCategoria)
			if err == nil && comissaoCat.GreaterThan(decimal.Zero) {
				uc.logger.Debug("usando comissão da categoria (nível 2)",
					zap.String("categoria_id", *servico.CategoriaID),
					zap.String("profissional_id", professionalID),
					zap.String("comissao", comissaoCat.String()))
				return &CommissionRuleResult{
					Rule: &entity.CommissionRule{
						ID:          "categoria-" + *servico.CategoriaID + "-" + professionalID, // ID virtual
						TenantID:    tenantID,
						Name:        "Comissão por Categoria",
						Type:        "PERCENTUAL", // Comissão por categoria é sempre percentual
						DefaultRate: comissaoCat,
						IsActive:    true,
					},
					Source:          "CATEGORIA",
					CalculationBase: "BRUTO", // Comissão por categoria usa valor bruto por padrão
				}
			}
		}
	}

	// 3º NÍVEL: Profissional tem comissão padrão configurada?
	if professionalInfo != nil && professionalInfo.Comissao != nil && *professionalInfo.Comissao != "" {
		comissaoProf, err := decimal.NewFromString(*professionalInfo.Comissao)
		if err == nil && comissaoProf.GreaterThan(decimal.Zero) {
			tipoComissao := "PERCENTUAL"
			if professionalInfo.TipoComissao != nil && *professionalInfo.TipoComissao != "" {
				tipoComissao = *professionalInfo.TipoComissao
			}
			uc.logger.Debug("usando comissão do profissional (nível 3)",
				zap.String("profissional_id", professionalID),
				zap.String("comissao", comissaoProf.String()),
				zap.String("tipo", tipoComissao))
			return &CommissionRuleResult{
				Rule: &entity.CommissionRule{
					ID:          "profissional-" + professionalID, // ID virtual
					TenantID:    tenantID,
					Name:        "Comissão do Profissional: " + professionalInfo.Name,
					Type:        tipoComissao,
					DefaultRate: comissaoProf,
					IsActive:    true,
				},
				Source:          "PROFISSIONAL",
				CalculationBase: "BRUTO", // Comissão de profissional usa valor bruto por padrão
			}
		}
	}

	// 4º NÍVEL: Regra específica da unidade?
	if unitID != nil && *unitID != "" {
		rule, err := uc.commissionRuleRepo.GetEffectiveByUnit(ctx, tenantIDStr, *unitID, time.Now())
		if err == nil && rule != nil {
			calcBase := "BRUTO"
			if rule.CalculationBase != nil {
				calcBase = *rule.CalculationBase
			}
			uc.logger.Debug("usando regra da unidade (nível 4)",
				zap.String("unit_id", *unitID),
				zap.String("rule_id", rule.ID),
				zap.String("calculation_base", calcBase))
			return &CommissionRuleResult{
				Rule:            rule,
				Source:          "REGRA",
				CalculationBase: calcBase,
			}
		}
	}

	// 5º NÍVEL: Regra global do tenant
	rule, err := uc.commissionRuleRepo.GetEffectiveGlobal(ctx, tenantIDStr, time.Now())
	if err == nil && rule != nil {
		calcBase := "BRUTO"
		if rule.CalculationBase != nil {
			calcBase = *rule.CalculationBase
		}
		uc.logger.Debug("usando regra global do tenant (nível 5)",
			zap.String("tenant_id", tenantIDStr),
			zap.String("rule_id", rule.ID),
			zap.String("calculation_base", calcBase))
		return &CommissionRuleResult{
			Rule:            rule,
			Source:          "REGRA",
			CalculationBase: calcBase,
		}
	}

	// Nenhuma regra encontrada
	return nil
}

// processarComissaoServicoHierarquica cria CommissionItem usando a hierarquia de regras
// COM-001, COM-002 & COM-003: Gera comissão com source correta e base de cálculo
func (uc *FinalizarComandaIntegradaUseCase) processarComissaoServicoHierarquica(
	ctx context.Context,
	tenantID uuid.UUID,
	professionalID string,
	command *entity.Command,
	item *entity.CommandItem,
	ruleResult *CommissionRuleResult,
	output *FinalizarComandaIntegradaOutput,
) error {
	rule := ruleResult.Rule
	commissionSource := ruleResult.Source
	calculationBase := ruleResult.CalculationBase

	// COM-003: Determinar valor base para cálculo
	var baseValue decimal.Decimal
	if calculationBase == "LIQUIDO" {
		// Calcular valor líquido proporcional do item
		baseValue = uc.calcularValorLiquidoProporcional(command, item)
		uc.logger.Debug("usando base LIQUIDO para comissão",
			zap.String("item_id", item.ID.String()),
			zap.String("preco_final", decimal.NewFromFloat(item.PrecoFinal).String()),
			zap.String("valor_liquido_proporcional", baseValue.String()))
	} else {
		// BRUTO (padrão): usar preço final do item
		baseValue = decimal.NewFromFloat(item.PrecoFinal)
	}

	// GrossValue sempre é o preço final do item (para registro)
	grossValue := decimal.NewFromFloat(item.PrecoFinal)

	// Criar item de comissão com base correta
	commissionItem, err := entity.NewCommissionItem(
		tenantID,
		professionalID,
		baseValue, // Valor usado para cálculo (pode ser bruto ou líquido)
		rule.DefaultRate,
		rule.Type,        // PERCENTUAL ou FIXO
		commissionSource, // SERVICO, PROFISSIONAL, ou REGRA
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("erro ao criar item de comissão: %w", err)
	}

	// Sobrescrever GrossValue com valor bruto real (para auditoria)
	commissionItem.GrossValue = grossValue

	// Associar à comanda e item
	commandIDStr := command.ID.String()
	itemIDStr := item.ID.String()
	commissionItem.CommandID = &commandIDStr
	commissionItem.CommandItemID = &itemIDStr
	commissionItem.ServiceID = func() *string { s := item.ItemID.String(); return &s }()
	commissionItem.ServiceName = &item.Descricao

	// RuleID apenas se for regra do banco (não serviço/profissional)
	if commissionSource == "REGRA" {
		commissionItem.RuleID = &rule.ID
	}

	// Persistir
	created, err := uc.commissionItemRepo.Create(ctx, commissionItem)
	if err != nil {
		return fmt.Errorf("erro ao persistir item de comissão: %w", err)
	}

	output.CommissionItems = append(output.CommissionItems, created.ID)
	output.TotalComissoes = output.TotalComissoes.Add(created.CommissionValue)

	uc.logger.Info("comissão gerada com hierarquia",
		zap.String("commission_item_id", created.ID),
		zap.String("professional_id", professionalID),
		zap.String("source", commissionSource),
		zap.String("calculation_base", calculationBase),
		zap.String("valor_bruto", grossValue.String()),
		zap.String("valor_base_calculo", baseValue.String()),
		zap.String("comissao", created.CommissionValue.String()),
		zap.String("tipo", rule.Type),
		zap.String("taxa", rule.DefaultRate.String()))

	return nil
}

// calcularValorLiquidoProporcional calcula o valor líquido proporcional de um item
// COM-003: Usado quando CalculationBase = LIQUIDO
// Fórmula: ValorLiquidoItem = (PrecoFinalItem / TotalComanda) * TotalValorLiquidoPagamentos
func (uc *FinalizarComandaIntegradaUseCase) calcularValorLiquidoProporcional(
	command *entity.Command,
	item *entity.CommandItem,
) decimal.Decimal {
	// Se comanda não tem total ou pagamentos, retorna valor bruto
	if command.Total <= 0 || len(command.Payments) == 0 {
		return decimal.NewFromFloat(item.PrecoFinal)
	}

	// Somar valor líquido de todos os pagamentos
	totalLiquido := decimal.Zero
	for _, payment := range command.Payments {
		totalLiquido = totalLiquido.Add(decimal.NewFromFloat(payment.ValorLiquido))
	}

	// Se não há valor líquido, retorna valor bruto
	if totalLiquido.IsZero() {
		return decimal.NewFromFloat(item.PrecoFinal)
	}

	// Calcular proporção do item no total da comanda
	precoItem := decimal.NewFromFloat(item.PrecoFinal)
	totalComanda := decimal.NewFromFloat(command.Total)

	// Proporção = PrecoItem / TotalComanda
	// ValorLiquidoItem = Proporção * TotalLiquido
	proporcao := precoItem.Div(totalComanda)
	valorLiquidoItem := proporcao.Mul(totalLiquido).Round(2)

	return valorLiquidoItem
}

// =============================================================================
// COM-004: Funções de comissão com data de referência correta
// =============================================================================

// processarComissaoServicoHierarquicaComData cria CommissionItem com data do agendamento
// COM-004: Usa appointmentDate como reference_date quando disponível
func (uc *FinalizarComandaIntegradaUseCase) processarComissaoServicoHierarquicaComData(
	ctx context.Context,
	tenantID uuid.UUID,
	professionalID string,
	command *entity.Command,
	item *entity.CommandItem,
	ruleResult *CommissionRuleResult,
	appointmentDate *time.Time,
	output *FinalizarComandaIntegradaOutput,
) error {
	rule := ruleResult.Rule
	commissionSource := ruleResult.Source
	calculationBase := ruleResult.CalculationBase

	// COM-003: Determinar valor base para cálculo
	var baseValue decimal.Decimal
	if calculationBase == "LIQUIDO" {
		baseValue = uc.calcularValorLiquidoProporcional(command, item)
		uc.logger.Debug("usando base LIQUIDO para comissão",
			zap.String("item_id", item.ID.String()),
			zap.String("preco_final", decimal.NewFromFloat(item.PrecoFinal).String()),
			zap.String("valor_liquido_proporcional", baseValue.String()))
	} else {
		baseValue = decimal.NewFromFloat(item.PrecoFinal)
	}

	grossValue := decimal.NewFromFloat(item.PrecoFinal)

	// COM-004: Usar data do agendamento ou fallback para data atual
	referenceDate := time.Now()
	if appointmentDate != nil {
		referenceDate = *appointmentDate
	}

	// Criar item de comissão com base correta e data correta
	commissionItem, err := entity.NewCommissionItem(
		tenantID,
		professionalID,
		baseValue,
		rule.DefaultRate,
		rule.Type,
		commissionSource,
		referenceDate, // Usa data do agendamento
	)
	if err != nil {
		return fmt.Errorf("erro ao criar item de comissão: %w", err)
	}

	commissionItem.GrossValue = grossValue

	commandIDStr := command.ID.String()
	itemIDStr := item.ID.String()
	commissionItem.CommandID = &commandIDStr
	commissionItem.CommandItemID = &itemIDStr
	commissionItem.ServiceID = func() *string { s := item.ItemID.String(); return &s }()
	commissionItem.ServiceName = &item.Descricao

	if commissionSource == "REGRA" {
		commissionItem.RuleID = &rule.ID
	}

	created, err := uc.commissionItemRepo.Create(ctx, commissionItem)
	if err != nil {
		return fmt.Errorf("erro ao persistir item de comissão: %w", err)
	}

	output.CommissionItems = append(output.CommissionItems, created.ID)
	output.TotalComissoes = output.TotalComissoes.Add(created.CommissionValue)

	uc.logger.Info("comissão de serviço gerada",
		zap.String("commission_item_id", created.ID),
		zap.String("professional_id", professionalID),
		zap.String("source", commissionSource),
		zap.String("reference_date", referenceDate.Format("2006-01-02")),
		zap.String("comissao", created.CommissionValue.String()))

	return nil
}

// =============================================================================
// COM-005: Comissão de Produtos
// =============================================================================

// buscarRegraComissaoProduto busca regra de comissão para produtos
// Hierarquia simplificada: 1) Profissional → 2) Unidade → 3) Global
// Produtos não têm comissão específica por item nem por categoria
func (uc *FinalizarComandaIntegradaUseCase) buscarRegraComissaoProduto(
	ctx context.Context,
	tenantID uuid.UUID,
	unitID *string,
	professionalID string,
	professionalInfo *port.ProfessionalInfo,
) *CommissionRuleResult {
	tenantIDStr := tenantID.String()

	// 1º NÍVEL: Profissional tem comissão padrão configurada?
	if professionalInfo != nil && professionalInfo.Comissao != nil && *professionalInfo.Comissao != "" {
		comissaoProf, err := decimal.NewFromString(*professionalInfo.Comissao)
		if err == nil && comissaoProf.GreaterThan(decimal.Zero) {
			tipoComissao := "PERCENTUAL"
			if professionalInfo.TipoComissao != nil && *professionalInfo.TipoComissao != "" {
				tipoComissao = *professionalInfo.TipoComissao
			}
			uc.logger.Debug("usando comissão do profissional para produto",
				zap.String("profissional_id", professionalID),
				zap.String("comissao", comissaoProf.String()))
			return &CommissionRuleResult{
				Rule: &entity.CommissionRule{
					ID:          "profissional-produto-" + professionalID,
					TenantID:    tenantID,
					Name:        "Comissão do Profissional (Produto): " + professionalInfo.Name,
					Type:        tipoComissao,
					DefaultRate: comissaoProf,
					IsActive:    true,
				},
				Source:          "PROFISSIONAL",
				CalculationBase: "BRUTO",
			}
		}
	}

	// 2º NÍVEL: Regra específica da unidade?
	if unitID != nil && *unitID != "" {
		rule, err := uc.commissionRuleRepo.GetEffectiveByUnit(ctx, tenantIDStr, *unitID, time.Now())
		if err == nil && rule != nil {
			calcBase := "BRUTO"
			if rule.CalculationBase != nil {
				calcBase = *rule.CalculationBase
			}
			return &CommissionRuleResult{
				Rule:            rule,
				Source:          "REGRA",
				CalculationBase: calcBase,
			}
		}
	}

	// 3º NÍVEL: Regra global do tenant
	rule, err := uc.commissionRuleRepo.GetEffectiveGlobal(ctx, tenantIDStr, time.Now())
	if err == nil && rule != nil {
		calcBase := "BRUTO"
		if rule.CalculationBase != nil {
			calcBase = *rule.CalculationBase
		}
		return &CommissionRuleResult{
			Rule:            rule,
			Source:          "REGRA",
			CalculationBase: calcBase,
		}
	}

	return nil
}

// processarComissaoProduto cria CommissionItem para um item do tipo PRODUTO
// COM-005: Gera comissão de produtos usando a mesma estrutura de serviços
func (uc *FinalizarComandaIntegradaUseCase) processarComissaoProduto(
	ctx context.Context,
	tenantID uuid.UUID,
	professionalID string,
	command *entity.Command,
	item *entity.CommandItem,
	ruleResult *CommissionRuleResult,
	appointmentDate *time.Time,
	output *FinalizarComandaIntegradaOutput,
) error {
	rule := ruleResult.Rule
	commissionSource := ruleResult.Source
	calculationBase := ruleResult.CalculationBase

	// Determinar valor base para cálculo (produtos usam valor bruto por padrão)
	var baseValue decimal.Decimal
	if calculationBase == "LIQUIDO" {
		baseValue = uc.calcularValorLiquidoProporcional(command, item)
	} else {
		baseValue = decimal.NewFromFloat(item.PrecoFinal)
	}

	grossValue := decimal.NewFromFloat(item.PrecoFinal)

	// Usar data do agendamento ou fallback para data atual
	referenceDate := time.Now()
	if appointmentDate != nil {
		referenceDate = *appointmentDate
	}

	// Criar item de comissão para produto
	commissionItem, err := entity.NewCommissionItem(
		tenantID,
		professionalID,
		baseValue,
		rule.DefaultRate,
		rule.Type,
		commissionSource,
		referenceDate,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar item de comissão de produto: %w", err)
	}

	commissionItem.GrossValue = grossValue

	commandIDStr := command.ID.String()
	itemIDStr := item.ID.String()
	commissionItem.CommandID = &commandIDStr
	commissionItem.CommandItemID = &itemIDStr
	// Para produtos, ServiceID é o ProdutoID e ServiceName é a descrição
	commissionItem.ServiceID = func() *string { s := item.ItemID.String(); return &s }()
	commissionItem.ServiceName = &item.Descricao

	if commissionSource == "REGRA" {
		commissionItem.RuleID = &rule.ID
	}

	created, err := uc.commissionItemRepo.Create(ctx, commissionItem)
	if err != nil {
		return fmt.Errorf("erro ao persistir item de comissão de produto: %w", err)
	}

	output.CommissionItems = append(output.CommissionItems, created.ID)
	output.TotalComissoes = output.TotalComissoes.Add(created.CommissionValue)

	uc.logger.Info("comissão de produto gerada",
		zap.String("commission_item_id", created.ID),
		zap.String("professional_id", professionalID),
		zap.String("produto_id", item.ItemID.String()),
		zap.String("source", commissionSource),
		zap.String("reference_date", referenceDate.Format("2006-01-02")),
		zap.String("comissao", created.CommissionValue.String()))

	return nil
}
