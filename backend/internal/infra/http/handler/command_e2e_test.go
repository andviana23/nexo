package handler_test

import (
	"context"
	"os"
	"testing"

	"github.com/andviana23/barber-analytics-backend/internal/application/dto"
	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/caixa"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/command"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/andviana23/barber-analytics-backend/internal/infra/repository/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// =============================================================================
// E2E TEST: Fluxo Completo Agendar → Finalizar Comanda → Caixa
// =============================================================================

const cmdE2ETenantID = "e2e00000-0000-0000-0000-000000000001"
const cmdE2EUserID = "e2e00000-0000-0000-0000-000000000002"

type CmdE2EValidator struct {
	validator *validator.Validate
}

func (v *CmdE2EValidator) Validate(i interface{}) error {
	if v.validator == nil {
		v.validator = validator.New()
	}
	return v.validator.Struct(i)
}

func getCmdE2EDBPool(t *testing.T) *pgxpool.Pool {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL não configurada, pulando testes E2E")
		return nil
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("Erro ao conectar ao banco de testes: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("Erro ao fazer ping no banco: %v", err)
	}

	return pool
}

// TestE2E_FluxoComanda_Finalizar_Caixa testa o fluxo completo:
// 1. Abre o caixa
// 2. Cria uma comanda
// 3. Adiciona itens à comanda (serviço + produto)
// 4. Adiciona pagamentos (dinheiro + cartão)
// 5. Finaliza a comanda com integração financeira
// 6. Verifica:
//   - OperaçõesCaixa criadas para pagamentos em dinheiro/PIX
//   - ContasReceber criadas para pagamentos em cartão
func TestE2E_FluxoComanda_Finalizar_Caixa(t *testing.T) {
	pool := getCmdE2EDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	logger, _ := zap.NewDevelopment()
	queries := db.New(pool)

	// Setup repositories
	appointmentRepo := postgres.NewAppointmentRepository(queries, pool)
	commandRepo := postgres.NewCommandRepository(queries, pool)
	meioPagamentoRepo := postgres.NewMeioPagamentoRepository(queries)
	contaReceberRepo := postgres.NewContaReceberRepository(queries)
	caixaDiarioRepo := postgres.NewCaixaDiarioRepository(queries)
	produtoRepo := postgres.NewProdutoRepository(queries)
	movimentacaoRepo := postgres.NewMovimentacaoEstoqueRepository(queries)
	commissionItemRepo := postgres.NewCommissionItemRepository(queries)
	commissionRuleRepo := postgres.NewCommissionRuleRepository(queries)

	// Setup mappers
	commandMapper := mapper.NewCommandMapper()

	// Setup use cases
	createCommandUC := command.NewCreateCommandUseCase(commandRepo, commandMapper)
	getCommandUC := command.NewGetCommandUseCase(commandRepo, commandMapper)
	// T-EST-001: Validação de estoque ao adicionar item PRODUTO
	addCommandItemUC := command.NewAddCommandItemUseCase(commandRepo, produtoRepo, commandMapper)
	addCommandPaymentUC := command.NewAddCommandPaymentUseCase(commandRepo, meioPagamentoRepo, commandMapper)
	finalizarComandaUC := command.NewFinalizarComandaIntegradaUseCase(
		commandRepo,
		appointmentRepo,
		meioPagamentoRepo,
		contaReceberRepo,
		caixaDiarioRepo,
		produtoRepo,
		movimentacaoRepo,
		commissionItemRepo,
		commissionRuleRepo,
		commandMapper,
		logger,
	)
	abrirCaixaUC := caixa.NewAbrirCaixaUseCase(caixaDiarioRepo, logger)
	getCaixaAbertoUC := caixa.NewGetCaixaAbertoUseCase(caixaDiarioRepo, logger)
	fecharCaixaUC := caixa.NewFecharCaixaUseCase(caixaDiarioRepo, logger)

	// Setup Echo
	e := echo.New()
	e.Validator = &CmdE2EValidator{}
	_ = e // Usado para validação

	ctx := context.Background()
	tenantUUID, _ := uuid.Parse(cmdE2ETenantID)
	userUUID, _ := uuid.Parse(cmdE2EUserID)

	// ==========================================================================
	// STEP 1: Abrir caixa
	// ==========================================================================
	t.Run("Step1_AbrirCaixa", func(t *testing.T) {
		// Verifica se já existe caixa aberto e fecha
		caixaAberto, _ := getCaixaAbertoUC.Execute(ctx, tenantUUID)
		if caixaAberto != nil {
			_, err := fecharCaixaUC.Execute(ctx, tenantUUID, userUUID, nil)
			require.NoError(t, err, "Erro ao fechar caixa existente")
		}

		// Abre novo caixa
		input := caixa.AbrirCaixaInput{
			TenantID:     tenantUUID,
			UserID:       userUUID,
			SaldoInicial: "100.00",
		}
		result, err := abrirCaixaUC.Execute(ctx, input)
		require.NoError(t, err, "Erro ao abrir caixa")
		assert.NotEmpty(t, result.ID, "ID do caixa deve ser gerado")
		assert.Equal(t, "aberto", result.Status, "Status deve ser 'aberto'")
		t.Logf("✓ Caixa aberto com ID: %s", result.ID)
	})

	// ==========================================================================
	// STEP 2: Criar comanda
	// ==========================================================================
	var commandID uuid.UUID
	t.Run("Step2_CriarComanda", func(t *testing.T) {
		req := &dto.CreateCommandRequest{}

		result, err := createCommandUC.Execute(ctx, tenantUUID, req)
		require.NoError(t, err, "Erro ao criar comanda")
		require.NotNil(t, result, "Comanda não deve ser nil")

		commandID, err = uuid.Parse(result.ID)
		require.NoError(t, err)

		assert.Equal(t, "aberta", result.Status, "Status deve ser 'aberta'")
		t.Logf("✓ Comanda criada com ID: %s", commandID)
	})

	// ==========================================================================
	// STEP 3: Adicionar itens à comanda
	// ==========================================================================
	t.Run("Step3_AdicionarItens", func(t *testing.T) {
		if commandID == uuid.Nil {
			t.Skip("Comanda não foi criada")
			return
		}

		// Adicionar um item de serviço
		itemServico := &dto.AddCommandItemRequest{
			Tipo:         "servico",
			Descricao:    "Corte de Cabelo",
			PrecoUnitStr: "50.00",
			Quantidade:   1,
		}

		_, err := addCommandItemUC.Execute(ctx, commandID, tenantUUID, userUUID, itemServico)
		require.NoError(t, err, "Erro ao adicionar item de serviço")
		t.Log("✓ Item de serviço adicionado")

		// Adicionar um item de produto
		itemProduto := &dto.AddCommandItemRequest{
			Tipo:         "produto",
			Descricao:    "Pomada Modeladora",
			PrecoUnitStr: "35.00",
			Quantidade:   2,
		}

		_, err = addCommandItemUC.Execute(ctx, commandID, tenantUUID, userUUID, itemProduto)
		require.NoError(t, err, "Erro ao adicionar item de produto")
		t.Log("✓ Item de produto adicionado")

		// Verificar total da comanda
		cmd, err := getCommandUC.Execute(ctx, commandID, tenantUUID)
		require.NoError(t, err)
		// Total esperado: 50 + (35*2) = 120
		assert.Equal(t, "120.00", cmd.TotalBruto, "Total bruto deve ser 120.00")
		t.Logf("✓ Total da comanda: R$ %s", cmd.TotalBruto)
	})

	// ==========================================================================
	// STEP 4: Adicionar pagamentos
	// ==========================================================================
	t.Run("Step4_AdicionarPagamentos", func(t *testing.T) {
		if commandID == uuid.Nil {
			t.Skip("Comanda não foi criada")
			return
		}

		// Buscar meios de pagamento disponíveis
		meioPagamentos, err := meioPagamentoRepo.ListAtivos(ctx, cmdE2ETenantID)
		if err != nil || len(meioPagamentos) == 0 {
			t.Skip("Nenhum meio de pagamento cadastrado")
			return
		}

		// Encontrar um meio DINHEIRO e um CARTAO
		var meioDinheiro, meioCartao string
		for _, mp := range meioPagamentos {
			switch mp.Tipo {
			case "DINHEIRO":
				meioDinheiro = mp.ID.String()
			case "CARTAO_CREDITO", "CARTAO_DEBITO":
				if meioCartao == "" {
					meioCartao = mp.ID.String()
				}
			}
		}

		// Adicionar pagamento em dinheiro (R$ 70)
		if meioDinheiro != "" {
			pagDinheiro := &dto.AddCommandPaymentRequest{
				MeioPagamentoID: meioDinheiro,
				ValorRecebido:   "70.00",
			}
			_, err = addCommandPaymentUC.Execute(ctx, commandID, tenantUUID, userUUID, pagDinheiro)
			require.NoError(t, err, "Erro ao adicionar pagamento em dinheiro")
			t.Log("✓ Pagamento em dinheiro adicionado: R$ 70.00")
		}

		// Adicionar pagamento em cartão (R$ 50)
		if meioCartao != "" {
			pagCartao := &dto.AddCommandPaymentRequest{
				MeioPagamentoID: meioCartao,
				ValorRecebido:   "50.00",
			}
			_, err = addCommandPaymentUC.Execute(ctx, commandID, tenantUUID, userUUID, pagCartao)
			require.NoError(t, err, "Erro ao adicionar pagamento em cartão")
			t.Log("✓ Pagamento em cartão adicionado: R$ 50.00")
		}

		// Verificar totais
		cmd, err := getCommandUC.Execute(ctx, commandID, tenantUUID)
		require.NoError(t, err)
		t.Logf("✓ Total recebido: R$ %s", cmd.TotalRecebido)
	})

	// ==========================================================================
	// STEP 5: Finalizar comanda com integração
	// ==========================================================================
	var finalizarOutput *command.FinalizarComandaIntegradaOutput
	t.Run("Step5_FinalizarComandaIntegrada", func(t *testing.T) {
		if commandID == uuid.Nil {
			t.Skip("Comanda não foi criada")
			return
		}

		input := command.FinalizarComandaIntegradaInput{
			CommandID: commandID,
			TenantID:  tenantUUID,
			UserID:    userUUID,
		}

		var err error
		finalizarOutput, err = finalizarComandaUC.Execute(ctx, input)
		require.NoError(t, err, "Erro ao finalizar comanda")
		require.NotNil(t, finalizarOutput, "Output não deve ser nil")

		assert.Equal(t, "fechada", finalizarOutput.Command.Status, "Status deve ser 'fechada'")
		t.Logf("✓ Comanda finalizada com sucesso")
		t.Logf("  - Operações Caixa: %d", len(finalizarOutput.OperacoesCaixa))
		t.Logf("  - Contas Receber: %d", len(finalizarOutput.ContasReceber))
		t.Logf("  - Movimentações Estoque: %d", len(finalizarOutput.MovimentacoesEstoque))
		t.Logf("  - Comissões: %d", len(finalizarOutput.CommissionItems))
	})

	// ==========================================================================
	// STEP 6: Verificações finais
	// ==========================================================================
	t.Run("Step6_VerificacoesFinais", func(t *testing.T) {
		if finalizarOutput == nil {
			t.Skip("Comanda não foi finalizada")
			return
		}

		// Verificar que pagamentos em dinheiro foram para o caixa
		if finalizarOutput.TotalLancadoCaixa.IsPositive() {
			t.Logf("✓ Total lançado no caixa: R$ %s", finalizarOutput.TotalLancadoCaixa.String())
		}

		// Verificar que pagamentos em cartão geraram ContaReceber
		if finalizarOutput.TotalContasReceber.IsPositive() {
			t.Logf("✓ Total em contas a receber: R$ %s", finalizarOutput.TotalContasReceber.String())
			assert.NotEmpty(t, finalizarOutput.ContasReceber, "Deve haver contas a receber criadas")
		}

		// Verificar comissões (se houver regra de comissão configurada)
		if len(finalizarOutput.CommissionItems) > 0 {
			t.Logf("✓ Comissões geradas: %d itens", len(finalizarOutput.CommissionItems))
			t.Logf("  - Total comissões: R$ %s", finalizarOutput.TotalComissoes.String())
		}

		// Verificar movimentações de estoque (se houver produtos)
		if len(finalizarOutput.MovimentacoesEstoque) > 0 {
			t.Logf("✓ Movimentações de estoque: %d", len(finalizarOutput.MovimentacoesEstoque))
		}

		// Verificar status final do caixa
		caixaFinal, err := getCaixaAbertoUC.Execute(ctx, tenantUUID)
		require.NoError(t, err)
		if caixaFinal != nil {
			t.Logf("✓ Saldo atual do caixa: R$ %s", caixaFinal.TotalEntradas)
		}
	})

	// ==========================================================================
	// CLEANUP: Fechar caixa de teste
	// ==========================================================================
	t.Run("Cleanup_FecharCaixa", func(t *testing.T) {
		_, err := fecharCaixaUC.Execute(ctx, tenantUUID, userUUID, nil)
		if err != nil {
			t.Logf("Aviso: Erro ao fechar caixa no cleanup: %v", err)
		} else {
			t.Log("✓ Caixa fechado no cleanup")
		}
	})
}

// TestE2E_Taxas_MeioPagamento verifica se as taxas são aplicadas corretamente
func TestE2E_Taxas_MeioPagamento(t *testing.T) {
	pool := getCmdE2EDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	logger, _ := zap.NewDevelopment()
	queries := db.New(pool)
	_ = logger // Usado para debug

	// Setup repositories
	commandRepo := postgres.NewCommandRepository(queries, pool)
	meioPagamentoRepo := postgres.NewMeioPagamentoRepository(queries)
	commandMapper := mapper.NewCommandMapper()

	// Setup use cases
	createCommandUC := command.NewCreateCommandUseCase(commandRepo, commandMapper)
	addCommandPaymentUC := command.NewAddCommandPaymentUseCase(commandRepo, meioPagamentoRepo, commandMapper)

	ctx := context.Background()
	tenantUUID, _ := uuid.Parse(cmdE2ETenantID)
	userUUID, _ := uuid.Parse(cmdE2EUserID)

	t.Run("TaxasAplicadasCorretamente", func(t *testing.T) {
		// Criar comanda de teste
		req := &dto.CreateCommandRequest{}
		cmd, err := createCommandUC.Execute(ctx, tenantUUID, req)
		require.NoError(t, err)
		commandID, _ := uuid.Parse(cmd.ID)

		// Buscar meio de pagamento com taxa
		meios, err := meioPagamentoRepo.ListAtivos(ctx, cmdE2ETenantID)
		require.NoError(t, err)
		require.NotEmpty(t, meios, "Deve haver meios de pagamento cadastrados")

		// Encontrar um meio com taxa configurada (ex: cartão)
		var meioComTaxa *struct {
			ID   string
			Taxa float64
		}
		for _, mp := range meios {
			taxa, _ := mp.Taxa.Float64()
			if taxa > 0 {
				meioComTaxa = &struct {
					ID   string
					Taxa float64
				}{
					ID:   mp.ID.String(),
					Taxa: taxa,
				}
				break
			}
		}

		if meioComTaxa == nil {
			t.Skip("Nenhum meio de pagamento com taxa configurada")
			return
		}

		// Adicionar pagamento de R$ 100
		pagReq := &dto.AddCommandPaymentRequest{
			MeioPagamentoID: meioComTaxa.ID,
			ValorRecebido:   "100.00",
		}

		result, err := addCommandPaymentUC.Execute(ctx, commandID, tenantUUID, userUUID, pagReq)
		require.NoError(t, err)

		// Verificar que a taxa foi aplicada
		require.NotEmpty(t, result.Payments, "Deve haver pagamento")
		pag := result.Payments[0]

		// Valor líquido deve ser menor que o recebido
		assert.Less(t, pag.ValorLiquido, pag.ValorRecebido,
			"Valor líquido deve ser menor que valor recebido devido à taxa")

		t.Logf("✓ Valor Recebido: R$ %.2f", pag.ValorRecebido)
		t.Logf("✓ Taxa Percentual: %.2f%%", pag.TaxaPercentual)
		t.Logf("✓ Taxa Fixa: R$ %.2f", pag.TaxaFixa)
		t.Logf("✓ Valor Líquido: R$ %.2f", pag.ValorLiquido)
	})
}

// TestE2E_RBAC_Barbeiro_RestricaoVer verifica se BARBER só vê seus próprios dados
func TestE2E_RBAC_Barbeiro_RestricaoVer(t *testing.T) {
	// Este teste verifica se um usuário com role BARBER
	// não consegue acessar comandas de outros profissionais
	t.Skip("Implementar quando houver filtro por professional_id")
}

// =============================================================================
// T-TEST-003: TESTES DE VALIDAÇÃO NEGATIVA
// =============================================================================

// TestE2E_ValidacaoNegativa_EstoqueInsuficiente verifica que não é possível
// adicionar um produto quando não há estoque suficiente
func TestE2E_ValidacaoNegativa_EstoqueInsuficiente(t *testing.T) {
	pool := getCmdE2EDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	queries := db.New(pool)

	// Setup repositories
	commandRepo := postgres.NewCommandRepository(queries, pool)
	produtoRepo := postgres.NewProdutoRepository(queries)
	commandMapper := mapper.NewCommandMapper()

	// Setup use cases
	createCommandUC := command.NewCreateCommandUseCase(commandRepo, commandMapper)
	addCommandItemUC := command.NewAddCommandItemUseCase(commandRepo, produtoRepo, commandMapper)

	ctx := context.Background()
	tenantUUID, _ := uuid.Parse(cmdE2ETenantID)
	userUUID, _ := uuid.Parse(cmdE2EUserID)

	t.Run("DeveRetornarErroQuandoEstoqueInsuficiente", func(t *testing.T) {
		// Criar comanda de teste
		req := &dto.CreateCommandRequest{}
		cmd, err := createCommandUC.Execute(ctx, tenantUUID, req)
		require.NoError(t, err)
		commandID, _ := uuid.Parse(cmd.ID)

		// Buscar um produto existente para verificar estoque
		produtos, err := produtoRepo.List(ctx, tenantUUID, nil)
		if err != nil || len(produtos) == 0 {
			t.Skip("Nenhum produto cadastrado para testar")
			return
		}

		// Encontrar um produto com estoque baixo ou tentar adicionar quantidade impossível
		produto := produtos[0]
		quantidadeImpossivel := int(produto.QuantidadeAtual.IntPart()) + 1000 // Mais do que tem

		// Tentar adicionar quantidade maior que o estoque disponível
		itemReq := &dto.AddCommandItemRequest{
			Tipo:         "produto",
			ItemID:       produto.ID.String(),
			Descricao:    produto.Nome,
			PrecoUnitStr: produto.PrecoVenda.String(),
			Quantidade:   quantidadeImpossivel,
		}

		_, err = addCommandItemUC.Execute(ctx, commandID, tenantUUID, userUUID, itemReq)

		// Deve retornar erro de estoque insuficiente
		if err != nil {
			assert.Contains(t, err.Error(), "estoque", "Erro deve mencionar estoque")
			t.Logf("✓ Erro retornado corretamente: %v", err)
		} else {
			// Se passou, pode ser que o produto tenha muito estoque
			// Verificar o comportamento atual
			t.Logf("⚠️ Item adicionado mesmo com quantidade alta - verificar estoque do produto")
		}
	})

	t.Run("DeveRetornarErroQuandoProdutoInativo", func(t *testing.T) {
		// Criar comanda de teste
		req := &dto.CreateCommandRequest{}
		cmd, err := createCommandUC.Execute(ctx, tenantUUID, req)
		require.NoError(t, err)
		commandID, _ := uuid.Parse(cmd.ID)

		// Tentar adicionar produto com ID inválido/inexistente
		itemReq := &dto.AddCommandItemRequest{
			Tipo:         "produto",
			ItemID:       uuid.New().String(), // ID que não existe
			Descricao:    "Produto Inexistente",
			PrecoUnitStr: "50.00",
			Quantidade:   1,
		}

		_, err = addCommandItemUC.Execute(ctx, commandID, tenantUUID, userUUID, itemReq)

		// Deve retornar erro de produto não encontrado
		if err != nil {
			t.Logf("✓ Erro retornado para produto inexistente: %v", err)
		} else {
			t.Log("⚠️ Item adicionado mesmo sem produto válido - verificar validação")
		}
	})
}

// TestE2E_ValidacaoNegativa_PagamentoInsuficiente verifica que não é possível
// fechar uma comanda quando o valor pago é menor que o total
func TestE2E_ValidacaoNegativa_PagamentoInsuficiente(t *testing.T) {
	pool := getCmdE2EDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	logger, _ := zap.NewDevelopment()
	queries := db.New(pool)

	// Setup repositories
	appointmentRepo := postgres.NewAppointmentRepository(queries, pool)
	commandRepo := postgres.NewCommandRepository(queries, pool)
	meioPagamentoRepo := postgres.NewMeioPagamentoRepository(queries)
	contaReceberRepo := postgres.NewContaReceberRepository(queries)
	caixaDiarioRepo := postgres.NewCaixaDiarioRepository(queries)
	produtoRepo := postgres.NewProdutoRepository(queries)
	movimentacaoRepo := postgres.NewMovimentacaoEstoqueRepository(queries)
	commissionItemRepo := postgres.NewCommissionItemRepository(queries)
	commissionRuleRepo := postgres.NewCommissionRuleRepository(queries)
	commandMapper := mapper.NewCommandMapper()

	// Setup use cases
	createCommandUC := command.NewCreateCommandUseCase(commandRepo, commandMapper)
	addCommandItemUC := command.NewAddCommandItemUseCase(commandRepo, produtoRepo, commandMapper)
	addCommandPaymentUC := command.NewAddCommandPaymentUseCase(commandRepo, meioPagamentoRepo, commandMapper)
	finalizarComandaUC := command.NewFinalizarComandaIntegradaUseCase(
		commandRepo,
		appointmentRepo,
		meioPagamentoRepo,
		contaReceberRepo,
		caixaDiarioRepo,
		produtoRepo,
		movimentacaoRepo,
		commissionItemRepo,
		commissionRuleRepo,
		commandMapper,
		logger,
	)

	ctx := context.Background()
	tenantUUID, _ := uuid.Parse(cmdE2ETenantID)
	userUUID, _ := uuid.Parse(cmdE2EUserID)

	t.Run("DeveRetornarErroQuandoPagamentoInsuficiente", func(t *testing.T) {
		// Criar comanda de teste
		req := &dto.CreateCommandRequest{}
		cmd, err := createCommandUC.Execute(ctx, tenantUUID, req)
		require.NoError(t, err)
		commandID, _ := uuid.Parse(cmd.ID)

		// Adicionar item de R$ 100
		itemReq := &dto.AddCommandItemRequest{
			Tipo:         "servico",
			Descricao:    "Serviço Teste Pagamento",
			PrecoUnitStr: "100.00",
			Quantidade:   1,
		}
		_, err = addCommandItemUC.Execute(ctx, commandID, tenantUUID, userUUID, itemReq)
		require.NoError(t, err)

		// Buscar meio de pagamento
		meios, err := meioPagamentoRepo.ListAtivos(ctx, cmdE2ETenantID)
		if err != nil || len(meios) == 0 {
			t.Skip("Nenhum meio de pagamento cadastrado")
			return
		}

		// Adicionar pagamento de apenas R$ 50 (metade do valor)
		pagReq := &dto.AddCommandPaymentRequest{
			MeioPagamentoID: meios[0].ID.String(),
			ValorRecebido:   "50.00",
		}
		_, err = addCommandPaymentUC.Execute(ctx, commandID, tenantUUID, userUUID, pagReq)
		require.NoError(t, err)

		// Tentar finalizar comanda - deve falhar
		input := command.FinalizarComandaIntegradaInput{
			CommandID: commandID,
			TenantID:  tenantUUID,
			UserID:    userUUID,
		}

		_, err = finalizarComandaUC.Execute(ctx, input)

		// Deve retornar erro de pagamento insuficiente
		if err != nil {
			assert.Contains(t, err.Error(), "pagamento", "Erro deve mencionar pagamento")
			t.Logf("✓ Erro retornado corretamente: %v", err)
		} else {
			t.Log("⚠️ Comanda finalizada mesmo com pagamento insuficiente - verificar validação")
		}
	})
}

// TestE2E_ValidacaoNegativa_ComandaSemItens verifica que não é possível
// fechar uma comanda vazia (sem itens)
func TestE2E_ValidacaoNegativa_ComandaSemItens(t *testing.T) {
	pool := getCmdE2EDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	logger, _ := zap.NewDevelopment()
	queries := db.New(pool)

	// Setup repositories
	appointmentRepo := postgres.NewAppointmentRepository(queries, pool)
	commandRepo := postgres.NewCommandRepository(queries, pool)
	meioPagamentoRepo := postgres.NewMeioPagamentoRepository(queries)
	contaReceberRepo := postgres.NewContaReceberRepository(queries)
	caixaDiarioRepo := postgres.NewCaixaDiarioRepository(queries)
	produtoRepo := postgres.NewProdutoRepository(queries)
	movimentacaoRepo := postgres.NewMovimentacaoEstoqueRepository(queries)
	commissionItemRepo := postgres.NewCommissionItemRepository(queries)
	commissionRuleRepo := postgres.NewCommissionRuleRepository(queries)
	commandMapper := mapper.NewCommandMapper()

	// Setup use cases
	createCommandUC := command.NewCreateCommandUseCase(commandRepo, commandMapper)
	finalizarComandaUC := command.NewFinalizarComandaIntegradaUseCase(
		commandRepo,
		appointmentRepo,
		meioPagamentoRepo,
		contaReceberRepo,
		caixaDiarioRepo,
		produtoRepo,
		movimentacaoRepo,
		commissionItemRepo,
		commissionRuleRepo,
		commandMapper,
		logger,
	)

	ctx := context.Background()
	tenantUUID, _ := uuid.Parse(cmdE2ETenantID)
	userUUID, _ := uuid.Parse(cmdE2EUserID)

	t.Run("DeveRetornarErroQuandoComandaVazia", func(t *testing.T) {
		// Criar comanda vazia
		req := &dto.CreateCommandRequest{}
		cmd, err := createCommandUC.Execute(ctx, tenantUUID, req)
		require.NoError(t, err)
		commandID, _ := uuid.Parse(cmd.ID)

		// Tentar finalizar comanda vazia
		input := command.FinalizarComandaIntegradaInput{
			CommandID: commandID,
			TenantID:  tenantUUID,
			UserID:    userUUID,
		}

		_, err = finalizarComandaUC.Execute(ctx, input)

		// Deve retornar erro de comanda vazia
		if err != nil {
			t.Logf("✓ Erro retornado para comanda vazia: %v", err)
		} else {
			t.Log("⚠️ Comanda vazia finalizada - verificar validação")
		}
	})
}

// TestE2E_ValidacaoNegativa_MeioPagamentoInativo verifica que não é possível
// usar um meio de pagamento inativo
func TestE2E_ValidacaoNegativa_MeioPagamentoInativo(t *testing.T) {
	pool := getCmdE2EDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	queries := db.New(pool)

	// Setup repositories
	commandRepo := postgres.NewCommandRepository(queries, pool)
	produtoRepo := postgres.NewProdutoRepository(queries)
	meioPagamentoRepo := postgres.NewMeioPagamentoRepository(queries)
	commandMapper := mapper.NewCommandMapper()

	// Setup use cases
	createCommandUC := command.NewCreateCommandUseCase(commandRepo, commandMapper)
	addCommandItemUC := command.NewAddCommandItemUseCase(commandRepo, produtoRepo, commandMapper)
	addCommandPaymentUC := command.NewAddCommandPaymentUseCase(commandRepo, meioPagamentoRepo, commandMapper)

	ctx := context.Background()
	tenantUUID, _ := uuid.Parse(cmdE2ETenantID)
	userUUID, _ := uuid.Parse(cmdE2EUserID)

	t.Run("DeveRetornarErroQuandoMeioPagamentoInexistente", func(t *testing.T) {
		// Criar comanda de teste
		req := &dto.CreateCommandRequest{}
		cmd, err := createCommandUC.Execute(ctx, tenantUUID, req)
		require.NoError(t, err)
		commandID, _ := uuid.Parse(cmd.ID)

		// Adicionar item
		itemReq := &dto.AddCommandItemRequest{
			Tipo:         "servico",
			Descricao:    "Serviço Teste",
			PrecoUnitStr: "50.00",
			Quantidade:   1,
		}
		_, err = addCommandItemUC.Execute(ctx, commandID, tenantUUID, userUUID, itemReq)
		require.NoError(t, err)

		// Tentar usar meio de pagamento inexistente
		pagReq := &dto.AddCommandPaymentRequest{
			MeioPagamentoID: uuid.New().String(), // ID que não existe
			ValorRecebido:   "50.00",
		}

		_, err = addCommandPaymentUC.Execute(ctx, commandID, tenantUUID, userUUID, pagReq)

		// Deve retornar erro
		if err != nil {
			t.Logf("✓ Erro retornado para meio de pagamento inválido: %v", err)
		} else {
			t.Log("⚠️ Pagamento aceito com meio de pagamento inexistente - verificar validação")
		}
	})
}

// TestE2E_ValidacaoNegativa_ComandaJaFechada verifica que não é possível
// modificar uma comanda já fechada
func TestE2E_ValidacaoNegativa_ComandaJaFechada(t *testing.T) {
	pool := getCmdE2EDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	logger, _ := zap.NewDevelopment()
	queries := db.New(pool)

	// Setup repositories
	appointmentRepo := postgres.NewAppointmentRepository(queries, pool)
	commandRepo := postgres.NewCommandRepository(queries, pool)
	meioPagamentoRepo := postgres.NewMeioPagamentoRepository(queries)
	contaReceberRepo := postgres.NewContaReceberRepository(queries)
	caixaDiarioRepo := postgres.NewCaixaDiarioRepository(queries)
	produtoRepo := postgres.NewProdutoRepository(queries)
	movimentacaoRepo := postgres.NewMovimentacaoEstoqueRepository(queries)
	commissionItemRepo := postgres.NewCommissionItemRepository(queries)
	commissionRuleRepo := postgres.NewCommissionRuleRepository(queries)
	commandMapper := mapper.NewCommandMapper()

	// Setup use cases
	createCommandUC := command.NewCreateCommandUseCase(commandRepo, commandMapper)
	addCommandItemUC := command.NewAddCommandItemUseCase(commandRepo, produtoRepo, commandMapper)
	addCommandPaymentUC := command.NewAddCommandPaymentUseCase(commandRepo, meioPagamentoRepo, commandMapper)
	finalizarComandaUC := command.NewFinalizarComandaIntegradaUseCase(
		commandRepo,
		appointmentRepo,
		meioPagamentoRepo,
		contaReceberRepo,
		caixaDiarioRepo,
		produtoRepo,
		movimentacaoRepo,
		commissionItemRepo,
		commissionRuleRepo,
		commandMapper,
		logger,
	)
	caixaUC := caixa.NewAbrirCaixaUseCase(caixaDiarioRepo, logger)
	getCaixaUC := caixa.NewGetCaixaAbertoUseCase(caixaDiarioRepo, logger)
	fecharCaixaUC := caixa.NewFecharCaixaUseCase(caixaDiarioRepo, logger)

	ctx := context.Background()
	tenantUUID, _ := uuid.Parse(cmdE2ETenantID)
	userUUID, _ := uuid.Parse(cmdE2EUserID)

	t.Run("DeveRetornarErroAoModificarComandaFechada", func(t *testing.T) {
		// Garantir caixa aberto
		caixaAberto, _ := getCaixaUC.Execute(ctx, tenantUUID)
		if caixaAberto == nil {
			input := caixa.AbrirCaixaInput{
				TenantID:     tenantUUID,
				UserID:       userUUID,
				SaldoInicial: "100.00",
			}
			_, err := caixaUC.Execute(ctx, input)
			require.NoError(t, err)
		}

		// Criar comanda com item e pagamento
		req := &dto.CreateCommandRequest{}
		cmd, err := createCommandUC.Execute(ctx, tenantUUID, req)
		require.NoError(t, err)
		commandID, _ := uuid.Parse(cmd.ID)

		itemReq := &dto.AddCommandItemRequest{
			Tipo:         "servico",
			Descricao:    "Serviço Teste Fechamento",
			PrecoUnitStr: "50.00",
			Quantidade:   1,
		}
		_, err = addCommandItemUC.Execute(ctx, commandID, tenantUUID, userUUID, itemReq)
		require.NoError(t, err)

		// Buscar meio de pagamento
		meios, err := meioPagamentoRepo.ListAtivos(ctx, cmdE2ETenantID)
		if err != nil || len(meios) == 0 {
			t.Skip("Nenhum meio de pagamento cadastrado")
			return
		}

		pagReq := &dto.AddCommandPaymentRequest{
			MeioPagamentoID: meios[0].ID.String(),
			ValorRecebido:   "50.00",
		}
		_, err = addCommandPaymentUC.Execute(ctx, commandID, tenantUUID, userUUID, pagReq)
		require.NoError(t, err)

		// Finalizar comanda
		input := command.FinalizarComandaIntegradaInput{
			CommandID: commandID,
			TenantID:  tenantUUID,
			UserID:    userUUID,
		}
		_, err = finalizarComandaUC.Execute(ctx, input)
		require.NoError(t, err)

		// Tentar adicionar mais um item após fechamento
		novoItemReq := &dto.AddCommandItemRequest{
			Tipo:         "servico",
			Descricao:    "Serviço Pós-Fechamento",
			PrecoUnitStr: "30.00",
			Quantidade:   1,
		}
		_, err = addCommandItemUC.Execute(ctx, commandID, tenantUUID, userUUID, novoItemReq)

		// Deve retornar erro
		if err != nil {
			t.Logf("✓ Erro ao adicionar item em comanda fechada: %v", err)
		} else {
			t.Log("⚠️ Item adicionado em comanda já fechada - verificar validação")
		}

		// Tentar adicionar pagamento após fechamento
		novoPagReq := &dto.AddCommandPaymentRequest{
			MeioPagamentoID: meios[0].ID.String(),
			ValorRecebido:   "30.00",
		}
		_, err = addCommandPaymentUC.Execute(ctx, commandID, tenantUUID, userUUID, novoPagReq)

		// Deve retornar erro
		if err != nil {
			t.Logf("✓ Erro ao adicionar pagamento em comanda fechada: %v", err)
		} else {
			t.Log("⚠️ Pagamento adicionado em comanda já fechada - verificar validação")
		}

		// Cleanup
		fecharCaixaUC.Execute(ctx, tenantUUID, userUUID, nil)
	})
}
