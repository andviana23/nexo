package main

// @title NEXO - Barber Analytics API
// @version 2.0
// @description API completa para gestão de barbearias com módulos de Agendamento, Comanda, Financeiro, Estoque e Analytics
// @termsOfService https://nexo.barber/terms

// @contact.name API Support
// @contact.url https://nexo.barber/support
// @contact.email support@nexo.barber

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token JWT (formato: "Bearer {token}")

// @x-extension-openapi {"example": "value on a json format"}

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/application/mapper"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/appointment"
	authUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/auth"
	barberturnUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/barberturn"
	blockedtimeUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/blockedtime"
	caixaUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/caixa"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/categoria"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/categoriaproduto"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/command"
	commissionUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/commission"
	customerUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/customer"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/financial"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/meiopagamento"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/metas"
	planUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/plan"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/pricing"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/servico"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/stock"
	subscriptionUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/subscription"
	unitUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/unit"
	"github.com/andviana23/barber-analytics-backend/internal/infra/auth"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/andviana23/barber-analytics-backend/internal/infra/gateway/asaas"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/handler"
	mw "github.com/andviana23/barber-analytics-backend/internal/infra/http/middleware"
	"github.com/andviana23/barber-analytics-backend/internal/infra/repository/postgres"
	"github.com/andviana23/barber-analytics-backend/internal/infra/scheduler"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"

	_ "github.com/andviana23/barber-analytics-backend/docs"
)

// CustomValidator é um wrapper para o validator do go-playground
type CustomValidator struct {
	validate *validator.Validate
}

// Validate implementa a interface echo.Validator
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}

func main() {
	// Load .env file (ignore error in production where env vars are set directly)
	_ = godotenv.Load()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Erro ao inicializar logger: %v", err)
	}
	defer logger.Sync()

	// Initialize Sentry
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN_BACKEND"),
		Environment:      os.Getenv("SENTRY_ENV"),
		TracesSampleRate: 1.0, // Enable performance and APM
	})
	if err != nil {
		logger.Error("Erro ao inicializar Sentry", zap.Error(err))
	}
	defer sentry.Flush(2 * time.Second)

	// Load environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logger.Fatal("DATABASE_URL não configurada")
	}

	// Initialize database connection
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		logger.Fatal("Erro ao conectar ao banco de dados", zap.Error(err))
	}
	defer dbPool.Close()

	// Test connection
	if err := dbPool.Ping(ctx); err != nil {
		logger.Fatal("Erro ao fazer ping no banco", zap.Error(err))
	}
	logger.Info("Conectado ao banco de dados PostgreSQL")

	// Initialize sqlc queries
	queries := db.New(dbPool)

	// Initialize repositories
	metaMensalRepo := postgres.NewMetaMensalRepository(queries)
	metaBarbeiroRepo := postgres.NewMetaBarbeiroRepository(queries)
	metaTicketMedioRepo := postgres.NewMetasTicketMedioRepository(queries)

	// Pricing repositories
	precificacaoConfigRepo := postgres.NewPrecificacaoConfigRepository(queries)
	precificacaoSimulacaoRepo := postgres.NewPrecificacaoSimulacaoRepository(queries)

	// Financial repositories
	contaPagarRepo := postgres.NewContaPagarRepository(queries)
	contaReceberRepo := postgres.NewContaReceberRepository(queries)
	compensacaoRepo := postgres.NewCompensacaoBancariaRepository(queries)
	fluxoCaixaRepo := postgres.NewFluxoCaixaDiarioRepository(queries)
	dreRepo := postgres.NewDREMensalRepository(queries)
	despesaFixaRepo := postgres.NewDespesaFixaRepository(queries)

	// Asaas Webhook Log repository (Migration 041)
	webhookLogRepo := postgres.NewAsaasWebhookLogRepository(queries)

	// Stock repositories
	produtoRepo := postgres.NewProdutoRepository(queries)
	fornecedorRepo := postgres.NewFornecedorRepositoryPG(queries)
	movimentacaoRepo := postgres.NewMovimentacaoEstoqueRepositoryPG(queries)

	// Appointment repositories and readers
	appointmentRepo := postgres.NewAppointmentRepository(queries, dbPool)
	professionalReader := postgres.NewProfessionalReader(queries)
	customerReader := postgres.NewCustomerReader(queries)
	serviceReader := postgres.NewServiceReader(queries)

	// Blocked Time repository
	blockedTimeRepo := postgres.NewBlockedTimeRepository(queries)

	// Command repository
	commandRepo := postgres.NewCommandRepository(queries, dbPool)

	// Customer repository
	customerRepo := postgres.NewCustomerRepository(queries)

	// Professional repository
	professionalRepo := postgres.NewProfessionalRepository(queries)

	// Barber Turn (Lista da Vez) repository
	barberTurnRepo := postgres.NewBarberTurnRepository(queries)

	// Categoria Servico repository
	categoriaServicoRepo := postgres.NewCategoriaServicoRepository(queries)

	// Categoria Produto repository
	categoriaProdutoRepo := postgres.NewCategoriaProdutoRepository(queries)

	// Servico repository
	servicoRepo := postgres.NewServicoRepository(queries)

	// Assinaturas (Planos e Subscriptions)
	planRepo := postgres.NewPlanRepository(queries)
	subscriptionRepo := postgres.NewSubscriptionRepository(queries, dbPool)
	subscriptionPaymentRepo := postgres.NewSubscriptionPaymentRepository(queries, dbPool)

	// MeioPagamento repository
	meioPagamentoRepo := postgres.NewMeioPagamentoRepository(queries)

	// Caixa Diário repository
	caixaDiarioRepo := postgres.NewCaixaDiarioRepository(queries)

	// Unit repositories
	unitRepo := postgres.NewUnitRepository(queries)
	userUnitRepo := postgres.NewUserUnitRepository(queries)

	// Commission repositories
	commissionRuleRepo := postgres.NewCommissionRuleRepository(queries)
	commissionPeriodRepo := postgres.NewCommissionPeriodRepository(queries)
	commissionItemRepo := postgres.NewCommissionItemRepository(queries)
	advanceRepo := postgres.NewAdvanceRepository(queries)

	// Initialize Asaas Gateway (payment gateway integration)
	asaasClient := asaas.NewClient(asaas.Config{
		APIKey:      os.Getenv("ASAAS_API_KEY"),
		Environment: os.Getenv("ASAAS_ENV"), // "sandbox" or "production"
	}, logger)
	asaasGateway := asaas.NewGatewayAdapter(asaasClient, logger)

	// Initialize use cases - Meta Mensal
	setMetaMensalUC := metas.NewSetMetaMensalUseCase(metaMensalRepo, logger)
	getMetaMensalUC := metas.NewGetMetaMensalUseCase(metaMensalRepo, logger)
	listMetasMensaisUC := metas.NewListMetasMensaisUseCase(metaMensalRepo, logger)
	updateMetaMensalUC := metas.NewUpdateMetaMensalUseCase(metaMensalRepo, logger)
	deleteMetaMensalUC := metas.NewDeleteMetaMensalUseCase(metaMensalRepo, logger)

	// Initialize use cases - Meta Barbeiro
	setMetaBarbeiroUC := metas.NewSetMetaBarbeiroUseCase(metaBarbeiroRepo, logger)
	getMetaBarbeiroUC := metas.NewGetMetaBarbeiroUseCase(metaBarbeiroRepo, logger)
	listMetasBarbeiroUC := metas.NewListMetasBarbeiroUseCase(metaBarbeiroRepo, logger)
	updateMetaBarbeiroUC := metas.NewUpdateMetaBarbeiroUseCase(metaBarbeiroRepo, logger)
	deleteMetaBarbeiroUC := metas.NewDeleteMetaBarbeiroUseCase(metaBarbeiroRepo, logger)

	// Initialize use cases - Meta Ticket Médio
	setMetaTicketMedioUC := metas.NewSetMetaTicketUseCase(metaTicketMedioRepo, logger)
	getMetaTicketMedioUC := metas.NewGetMetaTicketMedioUseCase(metaTicketMedioRepo, logger)
	listMetasTicketMedioUC := metas.NewListMetasTicketMedioUseCase(metaTicketMedioRepo, logger)
	updateMetaTicketMedioUC := metas.NewUpdateMetaTicketMedioUseCase(metaTicketMedioRepo, logger)
	deleteMetaTicketMedioUC := metas.NewDeleteMetaTicketMedioUseCase(metaTicketMedioRepo, logger)

	// Initialize use cases - Pricing (9 use cases)
	saveConfigUC := pricing.NewSaveConfigPrecificacaoUseCase(precificacaoConfigRepo, logger)
	getConfigUC := pricing.NewGetPrecificacaoConfigUseCase(precificacaoConfigRepo, logger)
	updateConfigUC := pricing.NewUpdatePrecificacaoConfigUseCase(precificacaoConfigRepo, logger)
	deleteConfigUC := pricing.NewDeletePrecificacaoConfigUseCase(precificacaoConfigRepo, logger)
	simularPrecoUC := pricing.NewSimularPrecoUseCase(precificacaoConfigRepo, precificacaoSimulacaoRepo, logger)
	saveSimulacaoUC := pricing.NewSaveSimulacaoUseCase(precificacaoSimulacaoRepo, logger)
	getSimulacaoUC := pricing.NewGetSimulacaoUseCase(precificacaoSimulacaoRepo, logger)
	listSimulacoesUC := pricing.NewListSimulacoesUseCase(precificacaoSimulacaoRepo, logger)
	deleteSimulacaoUC := pricing.NewDeleteSimulacaoUseCase(precificacaoSimulacaoRepo, logger)

	// Initialize use cases - Financial (23 use cases)
	// ContaPagar
	createContaPagarUC := financial.NewCreateContaPagarUseCase(contaPagarRepo, logger)
	getContaPagarUC := financial.NewGetContaPagarUseCase(contaPagarRepo, logger)
	listContasPagarUC := financial.NewListContasPagarUseCase(contaPagarRepo, logger)
	updateContaPagarUC := financial.NewUpdateContaPagarUseCase(contaPagarRepo, logger)
	deleteContaPagarUC := financial.NewDeleteContaPagarUseCase(contaPagarRepo, logger)
	marcarPagamentoUC := financial.NewMarcarPagamentoUseCase(contaPagarRepo, logger)
	// ContaReceber
	createContaReceberUC := financial.NewCreateContaReceberUseCase(contaReceberRepo, logger)
	getContaReceberUC := financial.NewGetContaReceberUseCase(contaReceberRepo, logger)
	listContasReceberUC := financial.NewListContasReceberUseCase(contaReceberRepo, logger)
	updateContaReceberUC := financial.NewUpdateContaReceberUseCase(contaReceberRepo, logger)
	deleteContaReceberUC := financial.NewDeleteContaReceberUseCase(contaReceberRepo, logger)
	marcarRecebimentoUC := financial.NewMarcarRecebimentoUseCase(contaReceberRepo, logger)
	// Compensação
	createCompensacaoUC := financial.NewCreateCompensacaoUseCase(compensacaoRepo, logger)
	getCompensacaoUC := financial.NewGetCompensacaoUseCase(compensacaoRepo, logger)
	listCompensacoesUC := financial.NewListCompensacoesUseCase(compensacaoRepo, logger)
	deleteCompensacaoUC := financial.NewDeleteCompensacaoUseCase(compensacaoRepo, logger)
	marcarCompensacaoUC := financial.NewMarcarCompensacaoUseCase(compensacaoRepo, contaReceberRepo, logger)
	// FluxoCaixa (com dependências de ContaPagar, ContaReceber e Compensacao)
	generateFluxoDiarioUC := financial.NewGenerateFluxoDiarioUseCase(fluxoCaixaRepo, contaPagarRepo, contaReceberRepo, compensacaoRepo, logger)
	generateFluxoDiarioV2UC := financial.NewGenerateFluxoDiarioV2UseCase(fluxoCaixaRepo, contaPagarRepo, contaReceberRepo, compensacaoRepo, logger)
	getFluxoCaixaUC := financial.NewGetFluxoCaixaUseCase(fluxoCaixaRepo, logger)
	listFluxoCaixaUC := financial.NewListFluxoCaixaUseCase(fluxoCaixaRepo, logger)
	// DRE (com dependências de ContaPagar, ContaReceber e CommissionItem)
	generateDREUC := financial.NewGenerateDREUseCase(dreRepo, contaPagarRepo, contaReceberRepo, commissionItemRepo, logger)
	generateDREV2UC := financial.NewGenerateDREV2UseCase(dreRepo, contaPagarRepo, contaReceberRepo, subscriptionPaymentRepo, logger)
	getDREUC := financial.NewGetDREUseCase(dreRepo, logger)
	listDREUC := financial.NewListDREUseCase(dreRepo, logger)
	// DespesaFixa (7 use cases)
	createDespesaFixaUC := financial.NewCreateDespesaFixaUseCase(despesaFixaRepo, logger)
	getDespesaFixaUC := financial.NewGetDespesaFixaUseCase(despesaFixaRepo, logger)
	listDespesasFixasUC := financial.NewListDespesasFixasUseCase(despesaFixaRepo, logger)
	updateDespesaFixaUC := financial.NewUpdateDespesaFixaUseCase(despesaFixaRepo, logger)
	toggleDespesaFixaUC := financial.NewToggleDespesaFixaUseCase(despesaFixaRepo, logger)
	deleteDespesaFixaUC := financial.NewDeleteDespesaFixaUseCase(despesaFixaRepo, logger)
	gerarContasFromDespesasUC := financial.NewGerarContasFromDespesasFixasUseCase(despesaFixaRepo, contaPagarRepo, logger)
	// Dashboard e Projeções (2 use cases)
	getPainelMensalUC := financial.NewGetPainelMensalUseCase(contaPagarRepo, contaReceberRepo, despesaFixaRepo, metaMensalRepo, fluxoCaixaRepo, logger)
	getProjecoesUC := financial.NewGetProjecoesUseCase(contaPagarRepo, contaReceberRepo, despesaFixaRepo, logger)

	// Initialize use cases - Unit (13 use cases)
	createUnitUC := unitUC.NewCreateUnitUseCase(unitRepo, userUnitRepo)
	listUnitsUC := unitUC.NewListUnitsUseCase(unitRepo)
	getUnitUC := unitUC.NewGetUnitUseCase(unitRepo)
	updateUnitUC := unitUC.NewUpdateUnitUseCase(unitRepo)
	deleteUnitUC := unitUC.NewDeleteUnitUseCase(unitRepo, userUnitRepo)
	toggleUnitUC := unitUC.NewToggleUnitUseCase(unitRepo)
	linkUserToUnitUC := unitUC.NewLinkUserToUnitUseCase(userUnitRepo)
	unlinkUserFromUnitUC := unitUC.NewUnlinkUserFromUnitUseCase(userUnitRepo)
	listUserUnitsUC := unitUC.NewListUserUnitsUseCase(userUnitRepo)
	setDefaultUnitUC := unitUC.NewSetDefaultUnitUseCase(userUnitRepo)
	getDefaultUnitUC := unitUC.NewGetDefaultUnitUseCase(userUnitRepo)
	checkUserAccessToUnitUC := unitUC.NewCheckUserAccessToUnitUseCase(userUnitRepo)
	listUnitUsersUC := unitUC.NewListUnitUsersUseCase(userUnitRepo)

	// Initialize use cases - Stock (5 use cases)
	criarProdutoUC := stock.NewCriarProdutoUseCase(produtoRepo, fornecedorRepo)
	registrarEntradaUC := stock.NewRegistrarEntradaUseCase(produtoRepo, movimentacaoRepo, fornecedorRepo)
	registrarSaidaUC := stock.NewRegistrarSaidaUseCase(produtoRepo, movimentacaoRepo)
	ajustarEstoqueUC := stock.NewAjustarEstoqueUseCase(produtoRepo, movimentacaoRepo)
	listarAlertasUC := stock.NewListarAlertasEstoqueBaixoUseCase(produtoRepo)

	// Initialize use cases - Appointments (7 use cases)
	// G-001: createAppointmentUC agora recebe commandRepo para criar comanda automaticamente
	createAppointmentUC := appointment.NewCreateAppointmentUseCase(appointmentRepo, commandRepo, serviceReader, professionalReader, customerReader, logger)
	listAppointmentsUC := appointment.NewListAppointmentsUseCase(appointmentRepo, logger)
	getAppointmentUC := appointment.NewGetAppointmentUseCase(appointmentRepo, logger)
	updateAppointmentStatusUC := appointment.NewUpdateAppointmentStatusUseCase(appointmentRepo, commandRepo, logger)
	rescheduleAppointmentUC := appointment.NewRescheduleAppointmentUseCase(appointmentRepo, professionalReader, logger)
	cancelAppointmentUC := appointment.NewCancelAppointmentUseCase(appointmentRepo, logger)
	finishWithCommandUC := appointment.NewFinishServiceWithCommandUseCase(appointmentRepo, commandRepo, logger)

	// Initialize use cases - Blocked Times (3 use cases)
	createBlockedTimeUC := blockedtimeUC.NewCreateBlockedTimeUseCase(blockedTimeRepo)
	listBlockedTimesUC := blockedtimeUC.NewListBlockedTimesUseCase(blockedTimeRepo)
	deleteBlockedTimeUC := blockedtimeUC.NewDeleteBlockedTimeUseCase(blockedTimeRepo)

	// Initialize mapper - Commands
	commandMapper := mapper.NewCommandMapper()

	// Initialize use cases - Commands (10 use cases)
	createCommandUC := command.NewCreateCommandUseCase(commandRepo, commandMapper)
	getCommandUC := command.NewGetCommandUseCase(commandRepo, commandMapper)
	listCommandsUC := command.NewListCommandsUseCase(commandRepo, commandMapper)
	getCommandByAppointmentUC := command.NewGetCommandByAppointmentUseCase(commandRepo, commandMapper)
	// T-EST-001: Validação de estoque ao adicionar item PRODUTO
	addCommandItemUC := command.NewAddCommandItemUseCase(commandRepo, produtoRepo, commandMapper)
	removeCommandItemUC := command.NewRemoveCommandItemUseCase(commandRepo, commandMapper)
	addCommandPaymentUC := command.NewAddCommandPaymentUseCase(commandRepo, meioPagamentoRepo, commandMapper)
	removeCommandPaymentUC := command.NewRemoveCommandPaymentUseCase(commandRepo, commandMapper)
	closeCommandUC := command.NewCloseCommandUseCase(commandRepo, appointmentRepo, commandMapper)
	// T-EST-002, T-COM-001: Finalização integrada com estoque e comissões
	// COM-001: Agora com hierarquia de 4 níveis para regras de comissão
	finalizarComandaIntegradaUC := command.NewFinalizarComandaIntegradaUseCase(
		commandRepo,
		appointmentRepo,
		meioPagamentoRepo,
		contaReceberRepo,
		compensacaoRepo,
		caixaDiarioRepo,
		produtoRepo,
		movimentacaoRepo,
		commissionItemRepo,
		commissionRuleRepo,
		serviceReader,      // COM-001: Para buscar comissão do serviço
		professionalReader, // COM-001: Para buscar comissão do profissional
		commandMapper,
		logger,
	)
	// T-EST-003: Cancelamento de comanda com reversão de estoque
	cancelCommandUC := command.NewCancelCommandUseCase(
		commandRepo,
		produtoRepo,
		movimentacaoRepo,
		commissionItemRepo,
		contaReceberRepo,
		caixaDiarioRepo,
		commandMapper,
		logger,
	)

	// Initialize use cases - Customer (12 use cases)
	createCustomerUC := customerUC.NewCreateCustomerUseCase(customerRepo, logger)
	updateCustomerUC := customerUC.NewUpdateCustomerUseCase(customerRepo, logger)
	listCustomersUC := customerUC.NewListCustomersUseCase(customerRepo, logger)
	getCustomerUC := customerUC.NewGetCustomerUseCase(customerRepo, logger)
	getCustomerWithHistoryUC := customerUC.NewGetCustomerWithHistoryUseCase(customerRepo, logger)
	inactivateCustomerUC := customerUC.NewInactivateCustomerUseCase(customerRepo, logger)
	searchCustomersUC := customerUC.NewSearchCustomersUseCase(customerRepo, logger)
	exportCustomerDataUC := customerUC.NewExportCustomerDataUseCase(customerRepo, logger)
	getCustomerStatsUC := customerUC.NewGetCustomerStatsUseCase(customerRepo, logger)
	checkPhoneDuplicateUC := customerUC.NewCheckPhoneDuplicateUseCase(customerRepo, logger)
	checkCPFDuplicateUC := customerUC.NewCheckCPFDuplicateUseCase(customerRepo, logger)

	// Initialize use cases - Barber Turn / Lista da Vez (9 use cases)
	listBarberTurnUC := barberturnUC.NewListBarberTurnUseCase(barberTurnRepo, logger)
	addBarberTurnUC := barberturnUC.NewAddBarberToTurnListUseCase(barberTurnRepo, logger)
	recordTurnUC := barberturnUC.NewRecordTurnUseCase(barberTurnRepo, logger)
	toggleBarberStatusUC := barberturnUC.NewToggleBarberStatusUseCase(barberTurnRepo, logger)
	removeBarberTurnUC := barberturnUC.NewRemoveBarberFromTurnListUseCase(barberTurnRepo, logger)
	resetTurnListUC := barberturnUC.NewResetTurnListUseCase(barberTurnRepo, logger)
	getTurnHistoryUC := barberturnUC.NewGetTurnHistoryUseCase(barberTurnRepo, logger)
	getHistorySummaryUC := barberturnUC.NewGetHistorySummaryUseCase(barberTurnRepo, logger)
	getAvailableBarbersUC := barberturnUC.NewGetAvailableBarbersUseCase(barberTurnRepo, logger)

	// Initialize use cases - Categoria Servico (5 use cases)
	createCategoriaUC := categoria.NewCreateCategoriaServicoUseCase(categoriaServicoRepo, logger)
	getCategoriaUC := categoria.NewGetCategoriaServicoUseCase(categoriaServicoRepo, logger)
	listCategoriasUC := categoria.NewListCategoriasServicosUseCase(categoriaServicoRepo, logger)
	updateCategoriaUC := categoria.NewUpdateCategoriaServicoUseCase(categoriaServicoRepo, logger)
	deleteCategoriaUC := categoria.NewDeleteCategoriaServicoUseCase(categoriaServicoRepo, logger)

	// Initialize use cases - Categoria Produto (6 use cases)
	createCategoriaProdutoUC := categoriaproduto.NewCreateCategoriaProdutoUseCase(categoriaProdutoRepo, logger)
	listCategoriasProdutosUC := categoriaproduto.NewListCategoriasProdutosUseCase(categoriaProdutoRepo, logger)
	getCategoriaProdutoUC := categoriaproduto.NewGetCategoriaProdutoUseCase(categoriaProdutoRepo, logger)
	updateCategoriaProdutoUC := categoriaproduto.NewUpdateCategoriaProdutoUseCase(categoriaProdutoRepo, logger)
	deleteCategoriaProdutoUC := categoriaproduto.NewDeleteCategoriaProdutoUseCase(categoriaProdutoRepo, logger)
	toggleCategoriaProdutoUC := categoriaproduto.NewToggleCategoriaProdutoUseCase(categoriaProdutoRepo, logger)

	// Initialize use cases - Servico
	createServicoUC := servico.NewCreateServicoUseCase(servicoRepo, categoriaServicoRepo, logger)
	getServicoUC := servico.NewGetServicoUseCase(servicoRepo, logger)
	getServicoStatsUC := servico.NewGetServicoStatsUseCase(servicoRepo, logger)
	listServicosUC := servico.NewListServicosUseCase(servicoRepo, logger)
	updateServicoUC := servico.NewUpdateServicoUseCase(servicoRepo, categoriaServicoRepo, logger)
	deleteServicoUC := servico.NewDeleteServicoUseCase(servicoRepo, logger)
	toggleServicoStatusUC := servico.NewToggleServicoStatusUseCase(servicoRepo, logger)

	// Use cases - Planos de Assinatura
	createPlanUC := planUC.NewCreatePlanUseCase(planRepo, logger)
	getPlanUC := planUC.NewGetPlanUseCase(planRepo)
	listPlansUC := planUC.NewListPlansUseCase(planRepo)
	updatePlanUC := planUC.NewUpdatePlanUseCase(planRepo, logger)
	deactivatePlanUC := planUC.NewDeactivatePlanUseCase(planRepo)

	// Use cases - Assinaturas
	createSubscriptionUC := subscriptionUC.NewCreateSubscriptionUseCase(planRepo, subscriptionRepo, customerRepo, asaasGateway, logger)
	listSubscriptionsUC := subscriptionUC.NewListSubscriptionsUseCase(subscriptionRepo)
	getSubscriptionUC := subscriptionUC.NewGetSubscriptionUseCase(subscriptionRepo)
	cancelSubscriptionUC := subscriptionUC.NewCancelSubscriptionUseCase(subscriptionRepo, asaasGateway, logger)
	renewSubscriptionUC := subscriptionUC.NewRenewSubscriptionUseCase(subscriptionRepo)
	subscriptionMetricsUC := subscriptionUC.NewGetSubscriptionMetricsUseCase(subscriptionRepo)
	overdueSubscriptionsUC := subscriptionUC.NewProcessOverdueSubscriptionsUseCase(subscriptionRepo)
	// processWebhookUC (V1) mantido para retrocompatibilidade se necessário
	_ = subscriptionUC.NewProcessWebhookUseCase(subscriptionRepo, subscriptionPaymentRepo, logger)
	// T-ASAAS-001: ProcessWebhookV2 agora lança no caixa quando PAYMENT_RECEIVED
	processWebhookUCV2 := subscriptionUC.NewProcessWebhookUseCaseV2(
		subscriptionRepo,
		subscriptionPaymentRepo,
		contaReceberRepo,
		webhookLogRepo,
		caixaDiarioRepo, // T-ASAAS-001: Adicionar caixa para lançar pagamentos
		logger,
	)
	// T-ASAAS-002: Reconciliação automática Asaas <-> NEXO
	reconcileAsaasUC := subscriptionUC.NewReconcileAsaasUseCase(
		subscriptionPaymentRepo,
		contaReceberRepo,
		webhookLogRepo,
		nil, // ReconciliationLogRepository - opcional, pode ser implementado depois
		logger,
	)

	// Initialize use cases - MeioPagamento (6 use cases)
	createMeioPagamentoUC := meiopagamento.NewCreateMeioPagamentoUseCase(meioPagamentoRepo)
	getMeioPagamentoUC := meiopagamento.NewGetMeioPagamentoUseCase(meioPagamentoRepo)
	listMeiosPagamentoUC := meiopagamento.NewListMeiosPagamentoUseCase(meioPagamentoRepo)
	updateMeioPagamentoUC := meiopagamento.NewUpdateMeioPagamentoUseCase(meioPagamentoRepo)
	deleteMeioPagamentoUC := meiopagamento.NewDeleteMeioPagamentoUseCase(meioPagamentoRepo)
	toggleMeioPagamentoUC := meiopagamento.NewToggleMeioPagamentoUseCase(meioPagamentoRepo)

	// Initialize use cases - Caixa Diário (8 use cases)
	abrirCaixaUC := caixaUC.NewAbrirCaixaUseCase(caixaDiarioRepo, logger)
	sangriaUC := caixaUC.NewSangriaUseCase(caixaDiarioRepo, contaPagarRepo, logger)
	reforcoUC := caixaUC.NewReforcoUseCase(caixaDiarioRepo, logger)
	fecharCaixaUC := caixaUC.NewFecharCaixaUseCase(caixaDiarioRepo, logger)
	getCaixaAbertoUC := caixaUC.NewGetCaixaAbertoUseCase(caixaDiarioRepo, logger)
	getCaixaByIDUC := caixaUC.NewGetCaixaByIDUseCase(caixaDiarioRepo, logger)
	listHistoricoCaixaUC := caixaUC.NewListHistoricoUseCase(caixaDiarioRepo, logger)
	getTotaisCaixaUC := caixaUC.NewGetTotaisCaixaUseCase(caixaDiarioRepo, logger)

	// Initialize use cases - Commission (31 use cases)
	// Commission Rules (7)
	createCommissionRuleUC := commissionUC.NewCreateCommissionRuleUseCase(commissionRuleRepo)
	getCommissionRuleUC := commissionUC.NewGetCommissionRuleUseCase(commissionRuleRepo)
	listCommissionRulesUC := commissionUC.NewListCommissionRulesUseCase(commissionRuleRepo)
	getEffectiveRulesUC := commissionUC.NewGetEffectiveCommissionRulesUseCase(commissionRuleRepo)
	updateCommissionRuleUC := commissionUC.NewUpdateCommissionRuleUseCase(commissionRuleRepo)
	deleteCommissionRuleUC := commissionUC.NewDeleteCommissionRuleUseCase(commissionRuleRepo)
	deactivateCommissionRuleUC := commissionUC.NewDeactivateCommissionRuleUseCase(commissionRuleRepo)
	// Commission Periods (8)
	createCommissionPeriodUC := commissionUC.NewCreateCommissionPeriodUseCase(commissionPeriodRepo)
	getCommissionPeriodUC := commissionUC.NewGetCommissionPeriodUseCase(commissionPeriodRepo)
	getOpenCommissionPeriodUC := commissionUC.NewGetOpenCommissionPeriodUseCase(commissionPeriodRepo)
	getCommissionPeriodSummaryUC := commissionUC.NewGetCommissionPeriodSummaryUseCase(commissionPeriodRepo)
	listCommissionPeriodsUC := commissionUC.NewListCommissionPeriodsUseCase(commissionPeriodRepo)
	// T-COM-002: Fechar período com geração de ContaPagar
	// COM-004: Incluir advanceRepo para dedução automática de adiantamentos
	closeCommissionPeriodUC := commissionUC.NewCloseCommissionPeriodUseCase(
		commissionPeriodRepo,
		commissionItemRepo,
		advanceRepo,
		contaPagarRepo,
		professionalReader,
		logger,
	)
	markPeriodPaidUC := commissionUC.NewMarkPeriodAsPaidUseCase(commissionPeriodRepo)
	deleteCommissionPeriodUC := commissionUC.NewDeleteCommissionPeriodUseCase(commissionPeriodRepo)
	// Advances (10)
	createAdvanceUC := commissionUC.NewCreateAdvanceUseCase(advanceRepo)
	getAdvanceUC := commissionUC.NewGetAdvanceUseCase(advanceRepo)
	listAdvancesUC := commissionUC.NewListAdvancesUseCase(advanceRepo)
	getPendingAdvancesUC := commissionUC.NewGetPendingAdvancesUseCase(advanceRepo)
	getApprovedAdvancesUC := commissionUC.NewGetApprovedAdvancesUseCase(advanceRepo)
	approveAdvanceUC := commissionUC.NewApproveAdvanceUseCase(advanceRepo)
	rejectAdvanceUC := commissionUC.NewRejectAdvanceUseCase(advanceRepo)
	markAdvanceDeductedUC := commissionUC.NewMarkAdvanceDeductedUseCase(advanceRepo)
	cancelAdvanceUC := commissionUC.NewCancelAdvanceUseCase(advanceRepo)
	deleteAdvanceUC := commissionUC.NewDeleteAdvanceUseCase(advanceRepo)
	// Commission Items (10)
	createCommissionItemUC := commissionUC.NewCreateCommissionItemUseCase(commissionItemRepo)
	createCommissionItemBatchUC := commissionUC.NewCreateCommissionItemBatchUseCase(commissionItemRepo)
	getCommissionItemUC := commissionUC.NewGetCommissionItemUseCase(commissionItemRepo)
	listCommissionItemsUC := commissionUC.NewListCommissionItemsUseCase(commissionItemRepo)
	getPendingCommissionItemsUC := commissionUC.NewGetPendingCommissionItemsUseCase(commissionItemRepo)
	getSummaryByProfessionalUC := commissionUC.NewGetCommissionSummaryByProfessionalUseCase(commissionItemRepo)
	getSummaryByServiceUC := commissionUC.NewGetCommissionSummaryByServiceUseCase(commissionItemRepo)
	processCommissionItemUC := commissionUC.NewProcessCommissionItemUseCase(commissionItemRepo)
	assignItemsToPeriodUC := commissionUC.NewAssignItemsToPeriodUseCase(commissionItemRepo)
	deleteCommissionItemUC := commissionUC.NewDeleteCommissionItemUseCase(commissionItemRepo)

	// Initialize JWT Manager
	jwtManager := auth.NewJWTManager()

	// Initialize use cases - Auth (4 use cases)
	loginUC := authUC.NewLoginUseCase(queries, jwtManager, logger)
	refreshUC := authUC.NewRefreshUseCase(queries, jwtManager, logger)
	meUC := authUC.NewMeUseCase(queries, logger)
	logoutUC := authUC.NewLogoutUseCase(queries, logger)

	// Initialize scheduler for cron jobs
	sched := scheduler.New(logger)

	// Register financial cron jobs
	financialDeps := scheduler.FinancialJobDeps{
		GenerateDRE:              generateDREUC,
		GenerateDREV2:            generateDREV2UC,
		GenerateFluxoDiario:      generateFluxoDiarioUC,
		GenerateFluxoDiarioV2:    generateFluxoDiarioV2UC,
		MarcarCompensacoes:       marcarCompensacaoUC,
		GerarContasDespesasFixas: gerarContasFromDespesasUC,
	}

	// Parse tenant list from ENV (SCHEDULER_TENANTS="tenant1,tenant2,...")
	tenants := scheduler.ParseTenantEnv("SCHEDULER_TENANTS")

	scheduler.RegisterFinancialJobs(sched, logger, financialDeps, tenants)

	subscriptionDeps := scheduler.SubscriptionJobDeps{
		ProcessOverdue: overdueSubscriptionsUC,
	}
	scheduler.RegisterSubscriptionJobs(sched, logger, subscriptionDeps, tenants)

	// Start scheduler in background
	sched.Start()
	defer sched.Stop(ctx)

	// Initialize handlers - Unit
	unitHandler := handler.NewUnitHandler(
		createUnitUC,
		listUnitsUC,
		getUnitUC,
		updateUnitUC,
		deleteUnitUC,
		toggleUnitUC,
		linkUserToUnitUC,
		unlinkUserFromUnitUC,
		listUserUnitsUC,
		setDefaultUnitUC,
		getDefaultUnitUC,
		checkUserAccessToUnitUC,
		listUnitUsersUC,
		logger,
	)

	// Initialize handlers - Metas completo (15 use cases)
	metasHandler := handler.NewMetasHandler(
		setMetaMensalUC,
		getMetaMensalUC,
		listMetasMensaisUC,
		updateMetaMensalUC,
		deleteMetaMensalUC,
		setMetaBarbeiroUC,
		getMetaBarbeiroUC,
		listMetasBarbeiroUC,
		updateMetaBarbeiroUC,
		deleteMetaBarbeiroUC,
		setMetaTicketMedioUC,
		getMetaTicketMedioUC,
		listMetasTicketMedioUC,
		updateMetaTicketMedioUC,
		deleteMetaTicketMedioUC,
		logger,
	)

	// Initialize handlers - Pricing (9 use cases)
	pricingHandler := handler.NewPricingHandler(
		// Config
		saveConfigUC,
		getConfigUC,
		updateConfigUC,
		deleteConfigUC,
		// Simulação
		simularPrecoUC,
		saveSimulacaoUC,
		getSimulacaoUC,
		listSimulacoesUC,
		deleteSimulacaoUC,
		logger,
	)

	// Initialize handlers - Financial (25 use cases)
	financialHandler := handler.NewFinancialHandler(
		// ContaPagar
		createContaPagarUC,
		getContaPagarUC,
		listContasPagarUC,
		updateContaPagarUC,
		deleteContaPagarUC,
		marcarPagamentoUC,
		// ContaReceber
		createContaReceberUC,
		getContaReceberUC,
		listContasReceberUC,
		updateContaReceberUC,
		deleteContaReceberUC,
		marcarRecebimentoUC,
		// Compensação
		createCompensacaoUC,
		getCompensacaoUC,
		listCompensacoesUC,
		deleteCompensacaoUC,
		marcarCompensacaoUC,
		// FluxoCaixa
		generateFluxoDiarioUC,
		getFluxoCaixaUC,
		listFluxoCaixaUC,
		// DRE
		generateDREUC,
		getDREUC,
		listDREUC,
		// Dashboard e Projeções
		getPainelMensalUC,
		getProjecoesUC,
		logger,
	)

	// Initialize handlers - Stock (5 use cases)
	stockHandler := handler.NewStockHandler(
		produtoRepo,
		criarProdutoUC,
		registrarEntradaUC,
		registrarSaidaUC,
		ajustarEstoqueUC,
		listarAlertasUC,
	)

	// Initialize handlers - Fornecedores
	fornecedorHandler := handler.NewFornecedorHandler(fornecedorRepo, logger)

	// Initialize handlers - Despesa Fixa (7 use cases)
	despesaFixaHandler := handler.NewDespesaFixaHandler(
		createDespesaFixaUC,
		getDespesaFixaUC,
		listDespesasFixasUC,
		updateDespesaFixaUC,
		toggleDespesaFixaUC,
		deleteDespesaFixaUC,
		gerarContasFromDespesasUC,
		despesaFixaRepo,
		logger,
	)

	// Initialize handlers - Auth (4 use cases)
	authHandler := handler.NewAuthHandler(
		loginUC,
		refreshUC,
		meUC,
		logoutUC,
		logger,
	)

	// Initialize handlers - Appointments (7 use cases)
	appointmentHandler := handler.NewAppointmentHandler(
		createAppointmentUC,
		listAppointmentsUC,
		getAppointmentUC,
		updateAppointmentStatusUC,
		rescheduleAppointmentUC,
		cancelAppointmentUC,
		finishWithCommandUC,
		logger,
	)

	// Initialize handlers - Blocked Times (3 use cases)
	blockedTimeHandler := handler.NewBlockedTimeHandler(
		createBlockedTimeUC,
		listBlockedTimesUC,
		deleteBlockedTimeUC,
		logger,
	)

	// Initialize handlers - Commands (11 use cases)
	commandHandler := handler.NewCommandHandler(
		createCommandUC,
		getCommandUC,
		listCommandsUC,
		getCommandByAppointmentUC,
		addCommandItemUC,
		removeCommandItemUC,
		addCommandPaymentUC,
		removeCommandPaymentUC,
		closeCommandUC,
		finalizarComandaIntegradaUC,
		cancelCommandUC, // T-EST-003: Cancelamento com reversão
		logger,
	)

	// Initialize handlers - Customers (11 use cases)
	customerHandler := handler.NewCustomerHandler(
		createCustomerUC,
		updateCustomerUC,
		listCustomersUC,
		getCustomerUC,
		getCustomerWithHistoryUC,
		inactivateCustomerUC,
		searchCustomersUC,
		exportCustomerDataUC,
		getCustomerStatsUC,
		checkPhoneDuplicateUC,
		checkCPFDuplicateUC,
		logger,
	)

	// Initialize handlers - Barber Turn / Lista da Vez (9 use cases)
	barberTurnHandler := handler.NewBarberTurnHandler(
		listBarberTurnUC,
		addBarberTurnUC,
		recordTurnUC,
		toggleBarberStatusUC,
		removeBarberTurnUC,
		resetTurnListUC,
		getTurnHistoryUC,
		getHistorySummaryUC,
		getAvailableBarbersUC,
		logger,
	)

	// Initialize handlers - Professionals (8 endpoints)
	professionalHandler := handler.NewProfessionalHandler(professionalRepo, logger)

	// Initialize handlers - Categoria Servico (5 endpoints)
	categoriaServicoHandler := handler.NewCategoriaServicoHandler(
		createCategoriaUC,
		getCategoriaUC,
		listCategoriasUC,
		updateCategoriaUC,
		deleteCategoriaUC,
		logger,
	)

	// Initialize handlers - Categoria Produto (6 endpoints)
	categoriaProdutoHandler := handler.NewCategoriaProdutoHandler(
		createCategoriaProdutoUC,
		listCategoriasProdutosUC,
		getCategoriaProdutoUC,
		updateCategoriaProdutoUC,
		deleteCategoriaProdutoUC,
		toggleCategoriaProdutoUC,
		logger,
	)

	servicoHandler := handler.NewServicoHandler(
		createServicoUC,
		getServicoUC,
		listServicosUC,
		updateServicoUC,
		deleteServicoUC,
		toggleServicoStatusUC,
		getServicoStatsUC,
		logger,
	)

	planHandler := handler.NewPlanHandler(
		createPlanUC,
		getPlanUC,
		listPlansUC,
		updatePlanUC,
		deactivatePlanUC,
		logger,
	)

	subscriptionHandler := handler.NewSubscriptionHandler(
		createSubscriptionUC,
		listSubscriptionsUC,
		getSubscriptionUC,
		cancelSubscriptionUC,
		renewSubscriptionUC,
		subscriptionMetricsUC,
		reconcileAsaasUC, // T-ASAAS-002
		logger,
	)

	// Initialize webhook handler for Asaas (using V2 with full audit logging)
	webhookHandler := handler.NewWebhookHandlerV2(
		processWebhookUCV2,
		os.Getenv("ASAAS_WEBHOOK_TOKEN"),
		logger,
	)

	// Initialize handlers - MeioPagamento (6 use cases)
	meioPagamentoHandler := handler.NewMeioPagamentoHandler(
		createMeioPagamentoUC,
		getMeioPagamentoUC,
		listMeiosPagamentoUC,
		updateMeioPagamentoUC,
		toggleMeioPagamentoUC,
		deleteMeioPagamentoUC,
		logger,
	)

	// Initialize handlers - Caixa Diário (8 use cases)
	caixaHandler := handler.NewCaixaHandler(
		abrirCaixaUC,
		sangriaUC,
		reforcoUC,
		fecharCaixaUC,
		getCaixaAbertoUC,
		getCaixaByIDUC,
		listHistoricoCaixaUC,
		getTotaisCaixaUC,
		logger,
	)

	// Initialize handlers - Commission (31 use cases)
	commissionHandler := handler.NewCommissionHandler(
		// Commission Rule UseCases
		createCommissionRuleUC,
		getCommissionRuleUC,
		listCommissionRulesUC,
		getEffectiveRulesUC,
		updateCommissionRuleUC,
		deleteCommissionRuleUC,
		deactivateCommissionRuleUC,
		// Commission Period UseCases
		createCommissionPeriodUC,
		getCommissionPeriodUC,
		getOpenCommissionPeriodUC,
		getCommissionPeriodSummaryUC,
		listCommissionPeriodsUC,
		closeCommissionPeriodUC,
		markPeriodPaidUC,
		deleteCommissionPeriodUC,
		// Advance UseCases
		createAdvanceUC,
		getAdvanceUC,
		listAdvancesUC,
		getPendingAdvancesUC,
		getApprovedAdvancesUC,
		approveAdvanceUC,
		rejectAdvanceUC,
		markAdvanceDeductedUC,
		cancelAdvanceUC,
		deleteAdvanceUC,
		// Commission Item UseCases
		createCommissionItemUC,
		createCommissionItemBatchUC,
		getCommissionItemUC,
		listCommissionItemsUC,
		getPendingCommissionItemsUC,
		getSummaryByProfessionalUC,
		getSummaryByServiceUC,
		processCommissionItemUC,
		assignItemsToPeriodUC,
		deleteCommissionItemUC,
		logger,
	)

	// Create Echo instance
	e := echo.New()

	// Add Sentry middleware
	e.Use(sentryecho.New(sentryecho.Options{}))

	// Configure Validator
	e.Validator = &CustomValidator{validate: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:3002", "http://localhost:3006", "http://localhost:8000"},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowCredentials: true, // Permite cookies (refresh token)
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Unit-ID"},
	}))

	// Health Check Endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status": "ok",
			"app":    "Barber Analytics Pro v2.0",
		})
	})

	// API Routes
	api := e.Group("/api/v1")

	// Auth routes - PÚBLICAS (sem middleware JWT)
	authGroup := api.Group("/auth")
	authGroup.POST("/login", authHandler.Login)                                // POST /api/v1/auth/login
	authGroup.POST("/refresh", authHandler.Refresh)                            // POST /api/v1/auth/refresh
	authGroup.POST("/logout", authHandler.Logout)                              // POST /api/v1/auth/logout
	authGroup.GET("/me", authHandler.Me, mw.JWTMiddleware(jwtManager, logger)) // GET /api/v1/auth/me (protegido)

	// Webhook routes - PÚBLICAS (validação por token no header)
	webhooksGroup := api.Group("/webhooks")
	webhooksGroup.POST("/asaas", webhookHandler.HandleAsaasWebhook) // POST /api/v1/webhooks/asaas

	// Middleware JWT para rotas protegidas
	protected := api.Group("")
	protected.Use(mw.JWTMiddleware(jwtManager, logger))

	// =============================================================================
	// T-ASAAS-003: Middleware de verificação de assinatura
	// Bloqueia tenants inadimplentes (assinatura vencida > 5 dias)
	// =============================================================================
	subscriptionChecker := mw.NewDefaultSubscriptionChecker(subscriptionRepo, logger)
	requireActiveSubscription := mw.RequireActiveSubscription(subscriptionChecker, logger)

	// Grupo guarded: JWT + verificação de assinatura ativa
	// Usado em rotas críticas de negócio (agendamentos, comandas, caixa, financeiro)
	guarded := api.Group("")
	guarded.Use(mw.JWTMiddleware(jwtManager, logger))
	guarded.Use(requireActiveSubscription)

	// Unit routes - PROTEGIDAS (JWT)
	unitGroup := protected.Group("/units")
	// User routes
	unitGroup.GET("/me", unitHandler.ListMyUnits)
	unitGroup.POST("/switch", unitHandler.SwitchUnit)
	unitGroup.POST("/default", unitHandler.SetDefaultUnit)
	// Admin routes
	unitGroup.POST("", unitHandler.Create, mw.RequireAdminAccess(logger))
	unitGroup.GET("", unitHandler.List, mw.RequireAdminAccess(logger))
	unitGroup.GET("/:id", unitHandler.GetByID, mw.RequireAdminAccess(logger))
	unitGroup.PUT("/:id", unitHandler.Update, mw.RequireAdminAccess(logger))
	unitGroup.DELETE("/:id", unitHandler.Delete, mw.RequireAdminAccess(logger))
	unitGroup.PATCH("/:id/toggle", unitHandler.Toggle, mw.RequireAdminAccess(logger))
	unitGroup.POST("/users", unitHandler.AddUserToUnit, mw.RequireAdminAccess(logger))
	unitGroup.DELETE("/:id/users/:userId", unitHandler.RemoveUserFromUnit, mw.RequireAdminAccess(logger))
	unitGroup.GET("/:id/users", unitHandler.ListUnitUsers, mw.RequireAdminAccess(logger))

	// Metas routes - 15 endpoints completos (PROTEGIDAS)
	metasGroup := protected.Group("/metas")

	// Meta Mensal (5 endpoints)
	metasGroup.POST("/monthly", metasHandler.SetMetaMensal)
	metasGroup.GET("/monthly/:id", metasHandler.GetMetaMensal)
	metasGroup.GET("/monthly", metasHandler.ListMetasMensais)
	metasGroup.PUT("/monthly/:id", metasHandler.UpdateMetaMensal)
	metasGroup.DELETE("/monthly/:id", metasHandler.DeleteMetaMensal)

	// Meta Barbeiro (5 endpoints)
	metasGroup.POST("/barbers", metasHandler.SetMetaBarbeiro)
	metasGroup.GET("/barbers/:id", metasHandler.GetMetaBarbeiro)
	metasGroup.GET("/barbers", metasHandler.ListMetasBarbeiro)
	metasGroup.PUT("/barbers/:id", metasHandler.UpdateMetaBarbeiro)
	metasGroup.DELETE("/barbers/:id", metasHandler.DeleteMetaBarbeiro)

	// Meta Ticket Médio (5 endpoints)
	metasGroup.POST("/ticket", metasHandler.SetMetaTicket)
	metasGroup.GET("/ticket/:id", metasHandler.GetMetaTicket)
	metasGroup.GET("/ticket", metasHandler.ListMetasTicket)
	metasGroup.PUT("/ticket/:id", metasHandler.UpdateMetaTicket)
	metasGroup.DELETE("/ticket/:id", metasHandler.DeleteMetaTicket)

	// Pricing routes - 9 endpoints (PROTEGIDAS com JWT)
	pricingGroup := protected.Group("/pricing")
	pricingHandler.RegisterRoutes(pricingGroup)

	// Appointment routes - 12 endpoints (PROTEGIDAS com JWT + RBAC + ASSINATURA ATIVA)
	// Regras:
	// - BARBER só vê/edita seus próprios agendamentos (filtro aplicado no handler)
	// - OWNER, MANAGER, RECEPTIONIST podem criar/listar/editar qualquer agendamento
	// - Apenas OWNER e MANAGER podem alterar status para DONE e NO_SHOW
	// - T-ASAAS-003: Requer assinatura ativa (grupo guarded)
	appointmentsGroup := guarded.Group("/appointments")
	appointmentsGroup.Use(mw.UnitMiddleware())
	// Rotas com acesso geral (todos os roles autenticados podem acessar, mas BARBER tem escopo limitado)
	appointmentsGroup.POST("", appointmentHandler.CreateAppointment, mw.RequireAnyRole(logger))
	appointmentsGroup.GET("", appointmentHandler.ListAppointments, mw.RequireAnyRole(logger))
	appointmentsGroup.GET("/:id", appointmentHandler.GetAppointment, mw.RequireAnyRole(logger))
	appointmentsGroup.PATCH("/:id/status", appointmentHandler.UpdateAppointmentStatus, mw.RequireAdminAccess(logger))
	appointmentsGroup.PATCH("/:id/reschedule", appointmentHandler.RescheduleAppointment, mw.RequireAdminAccess(logger))
	// Transições de status específicas
	appointmentsGroup.POST("/:id/confirm", appointmentHandler.ConfirmAppointment, mw.RequireAnyRole(logger))
	appointmentsGroup.POST("/:id/cancel", appointmentHandler.CancelAppointment, mw.RequireAdminAccess(logger))
	appointmentsGroup.POST("/:id/check-in", appointmentHandler.CheckInAppointment, mw.RequireAnyRole(logger))
	appointmentsGroup.POST("/:id/start", appointmentHandler.StartServiceAppointment, mw.RequireAnyRole(logger))
	appointmentsGroup.POST("/:id/finish", appointmentHandler.FinishServiceAppointment, mw.RequireAnyRole(logger))
	appointmentsGroup.POST("/:id/complete", appointmentHandler.CompleteAppointment, mw.RequireAdminAccess(logger))
	appointmentsGroup.POST("/:id/no-show", appointmentHandler.NoShowAppointment, mw.RequireOwnerOrManager(logger))

	// Blocked Times routes - 3 endpoints (PROTEGIDAS com JWT)
	blockedTimesGroup := protected.Group("/blocked-times")
	blockedTimesGroup.POST("", blockedTimeHandler.CreateBlockedTime)
	blockedTimesGroup.GET("", blockedTimeHandler.ListBlockedTimes)
	blockedTimesGroup.DELETE("/:id", blockedTimeHandler.DeleteBlockedTime)

	// Command routes - 11 endpoints (PROTEGIDAS com JWT + RBAC + ASSINATURA ATIVA)
	// T-ASAAS-003: Requer assinatura ativa (grupo guarded)
	commandsGroup := guarded.Group("/commands")
	commandsGroup.POST("", commandHandler.CreateCommand, mw.RequireAnyRole(logger))
	commandsGroup.GET("", commandHandler.ListCommands, mw.RequireAnyRole(logger)) // LIST - filtros e paginação
	commandsGroup.GET("/by-appointment/:appointmentId", commandHandler.GetCommandByAppointment, mw.RequireAnyRole(logger))
	commandsGroup.GET("/:id", commandHandler.GetCommand, mw.RequireAnyRole(logger))
	commandsGroup.POST("/:id/items", commandHandler.AddCommandItem, mw.RequireAnyRole(logger))
	commandsGroup.DELETE("/:id/items/:itemId", commandHandler.RemoveCommandItem, mw.RequireAdminAccess(logger))
	commandsGroup.POST("/:id/payments", commandHandler.AddCommandPayment, mw.RequireAnyRole(logger))
	commandsGroup.DELETE("/:id/payments/:paymentId", commandHandler.RemoveCommandPayment, mw.RequireAdminAccess(logger))
	commandsGroup.POST("/:id/close", commandHandler.CloseCommand, mw.RequireAdminAccess(logger))
	commandsGroup.POST("/:id/close-integrated", commandHandler.CloseCommandIntegrated, mw.RequireAdminAccess(logger)) // T-EST-002, T-COM-001
	commandsGroup.POST("/:id/cancel", commandHandler.CancelCommand, mw.RequireAdminAccess(logger))                    // T-EST-003: Cancelar comanda com reversão de estoque

	// Customer routes - 11 endpoints (PROTEGIDAS com JWT)
	customersGroup := protected.Group("/customers")
	customersGroup.POST("", customerHandler.CreateCustomer)
	customersGroup.GET("", customerHandler.ListCustomers)
	customersGroup.GET("/search", customerHandler.SearchCustomers)
	customersGroup.GET("/stats", customerHandler.GetCustomerStats)
	customersGroup.GET("/check-phone", customerHandler.CheckPhoneExists)
	customersGroup.GET("/check-cpf", customerHandler.CheckCPFExists)
	customersGroup.GET("/:id", customerHandler.GetCustomer)
	customersGroup.GET("/:id/history", customerHandler.GetCustomerWithHistory)
	customersGroup.GET("/:id/export", customerHandler.ExportCustomerData)
	customersGroup.PUT("/:id", customerHandler.UpdateCustomer)
	customersGroup.DELETE("/:id", customerHandler.InactivateCustomer)

	// Professional routes - 8 endpoints (PROTEGIDAS com JWT)
	professionalsGroup := protected.Group("/professionals")
	professionalsGroup.Use(mw.UnitMiddleware())
	professionalsGroup.GET("", professionalHandler.ListProfessionals)
	professionalsGroup.POST("", professionalHandler.CreateProfessional)
	professionalsGroup.GET("/check-email", professionalHandler.CheckEmailExists)
	professionalsGroup.GET("/check-cpf", professionalHandler.CheckCpfExists)
	professionalsGroup.GET("/:id", professionalHandler.GetProfessional)
	professionalsGroup.PUT("/:id", professionalHandler.UpdateProfessional)
	professionalsGroup.PUT("/:id/status", professionalHandler.UpdateProfessionalStatus)
	professionalsGroup.DELETE("/:id", professionalHandler.DeleteProfessional)

	// Financial routes - 19 endpoints (PROTEGIDAS com JWT + RBAC + ASSINATURA ATIVA)
	// T-ASAAS-003: Requer assinatura ativa (grupo guarded)
	financialGroup := guarded.Group("/financial")

	// ContaPagar (6 endpoints: 5 CRUD + 1 marcarPagamento) - Apenas Admin
	financialGroup.POST("/payables", financialHandler.CreateContaPagar, mw.RequireAdminAccess(logger))
	financialGroup.GET("/payables/:id", financialHandler.GetContaPagar, mw.RequireAdminAccess(logger))
	financialGroup.GET("/payables", financialHandler.ListContasPagar, mw.RequireAdminAccess(logger))
	financialGroup.PUT("/payables/:id", financialHandler.UpdateContaPagar, mw.RequireOwnerOrManager(logger))
	financialGroup.DELETE("/payables/:id", financialHandler.DeleteContaPagar, mw.RequireOwnerOrManager(logger))
	financialGroup.POST("/payables/:id/payment", financialHandler.MarcarPagamento, mw.RequireOwnerOrManager(logger))

	// ContaReceber (6 endpoints: 5 CRUD + 1 marcarRecebimento) - Apenas Admin
	financialGroup.POST("/receivables", financialHandler.CreateContaReceber, mw.RequireAdminAccess(logger))
	financialGroup.GET("/receivables/:id", financialHandler.GetContaReceber, mw.RequireAdminAccess(logger))
	financialGroup.GET("/receivables", financialHandler.ListContasReceber, mw.RequireAdminAccess(logger))
	financialGroup.PUT("/receivables/:id", financialHandler.UpdateContaReceber, mw.RequireOwnerOrManager(logger))
	financialGroup.DELETE("/receivables/:id", financialHandler.DeleteContaReceber, mw.RequireOwnerOrManager(logger))
	financialGroup.POST("/receivables/:id/receipt", financialHandler.MarcarRecebimento, mw.RequireOwnerOrManager(logger))

	// Compensação (3 endpoints: Get, List, Delete) - Apenas Admin
	financialGroup.GET("/compensations/:id", financialHandler.GetCompensacao, mw.RequireAdminAccess(logger))
	financialGroup.GET("/compensations", financialHandler.ListCompensacoes, mw.RequireAdminAccess(logger))
	financialGroup.DELETE("/compensations/:id", financialHandler.DeleteCompensacao, mw.RequireOwnerOrManager(logger))

	// FluxoCaixa (2 endpoints: Get, List) - Apenas Admin
	financialGroup.GET("/cashflow/:id", financialHandler.GetFluxoCaixa, mw.RequireAdminAccess(logger))
	financialGroup.GET("/cashflow", financialHandler.ListFluxoCaixa, mw.RequireAdminAccess(logger))

	// DRE (2 endpoints: Get, List) - Apenas Owner/Manager
	financialGroup.GET("/dre/:month", financialHandler.GetDRE, mw.RequireOwnerOrManager(logger))
	financialGroup.GET("/dre", financialHandler.ListDRE, mw.RequireOwnerOrManager(logger))

	// Dashboard e Projeções (2 endpoints: dashboard, projections) - Apenas Admin
	financialGroup.GET("/dashboard", financialHandler.GetDashboard, mw.RequireAdminAccess(logger))
	financialGroup.GET("/projections", financialHandler.GetProjections, mw.RequireAdminAccess(logger))

	// Despesas Fixas (8 endpoints: CRUD + toggle + summary + generate)
	despesaFixaHandler.RegisterRoutes(financialGroup)

	// Barber Turn (Lista da Vez) routes - 9 endpoints (PROTEGIDAS com JWT)
	turnGroup := protected.Group("/barber-turn")
	turnGroup.GET("/list", barberTurnHandler.ListBarbersTurn)                              // GET /api/v1/barber-turn/list
	turnGroup.POST("/add", barberTurnHandler.AddBarberToTurnList)                          // POST /api/v1/barber-turn/add
	turnGroup.POST("/record", barberTurnHandler.RecordTurn)                                // POST /api/v1/barber-turn/record
	turnGroup.PUT("/:professional_id/toggle-status", barberTurnHandler.ToggleBarberStatus) // PUT /api/v1/barber-turn/:professional_id/toggle-status
	turnGroup.DELETE("/:professional_id", barberTurnHandler.RemoveBarberFromTurnList)      // DELETE /api/v1/barber-turn/:professional_id
	turnGroup.POST("/reset", barberTurnHandler.ResetTurnList)                              // POST /api/v1/barber-turn/reset
	turnGroup.GET("/history", barberTurnHandler.GetTurnHistory)                            // GET /api/v1/barber-turn/history
	turnGroup.GET("/history/summary", barberTurnHandler.GetHistorySummary)                 // GET /api/v1/barber-turn/history/summary
	turnGroup.GET("/available", barberTurnHandler.GetAvailableBarbers)                     // GET /api/v1/barber-turn/available

	// Stock routes - 7 endpoints (PROTEGIDAS com JWT + RBAC + ASSINATURA ATIVA)
	// T-ASAAS-003: Requer assinatura ativa (grupo guarded)
	stockGroup := guarded.Group("/stock")
	stockGroup.GET("/items", stockHandler.ListProdutos, mw.RequireAnyRole(logger))               // GET /api/v1/stock/items - Listar produtos
	stockGroup.GET("/items/:id", stockHandler.GetProduto, mw.RequireAnyRole(logger))             // GET /api/v1/stock/items/:id - Buscar produto
	stockGroup.POST("/products", stockHandler.CreateProduto, mw.RequireOwnerOrManager(logger))   // POST /api/v1/stock/products - Criar produto
	stockGroup.POST("/entries", stockHandler.RegistrarEntrada, mw.RequireOwnerOrManager(logger)) // POST /api/v1/stock/entries - Registrar entrada
	stockGroup.POST("/exit", stockHandler.RegistrarSaida, mw.RequireAdminAccess(logger))         // POST /api/v1/stock/exit - Registrar saída
	stockGroup.POST("/adjust", stockHandler.AjustarEstoque, mw.RequireOwnerOrManager(logger))    // POST /api/v1/stock/adjust - Ajustar estoque
	stockGroup.GET("/alerts", stockHandler.ListarAlertas, mw.RequireAdminAccess(logger))         // GET /api/v1/stock/alerts - Listar alertas

	// Fornecedores routes - 7 endpoints (PROTEGIDAS com JWT)
	fornecedoresGroup := protected.Group("/fornecedores")
	fornecedorHandler.RegisterRoutes(fornecedoresGroup)

	// Categoria Servico routes - 5 endpoints (PROTEGIDAS com JWT)
	categoriasGroup := protected.Group("/categorias-servicos")
	categoriasGroup.Use(mw.UnitMiddleware())
	categoriasGroup.POST("", categoriaServicoHandler.Create)       // POST /api/v1/categorias-servicos
	categoriasGroup.GET("", categoriaServicoHandler.List)          // GET /api/v1/categorias-servicos
	categoriasGroup.GET("/:id", categoriaServicoHandler.GetByID)   // GET /api/v1/categorias-servicos/:id
	categoriasGroup.PUT("/:id", categoriaServicoHandler.Update)    // PUT /api/v1/categorias-servicos/:id
	categoriasGroup.DELETE("/:id", categoriaServicoHandler.Delete) // DELETE /api/v1/categorias-servicos/:id

	// Categoria Produto routes - 6 endpoints (PROTEGIDAS com JWT)
	categoriasProdutosGroup := protected.Group("/categorias-produtos")
	categoriasProdutosGroup.Use(mw.UnitMiddleware())
	categoriasProdutosGroup.POST("", categoriaProdutoHandler.Create)             // POST /api/v1/categorias-produtos
	categoriasProdutosGroup.GET("", categoriaProdutoHandler.List)                // GET /api/v1/categorias-produtos
	categoriasProdutosGroup.GET("/:id", categoriaProdutoHandler.GetByID)         // GET /api/v1/categorias-produtos/:id
	categoriasProdutosGroup.PUT("/:id", categoriaProdutoHandler.Update)          // PUT /api/v1/categorias-produtos/:id
	categoriasProdutosGroup.DELETE("/:id", categoriaProdutoHandler.Delete)       // DELETE /api/v1/categorias-produtos/:id
	categoriasProdutosGroup.PATCH("/:id/toggle", categoriaProdutoHandler.Toggle) // PATCH /api/v1/categorias-produtos/:id/toggle

	// Servicos routes - 9 endpoints (PROTEGIDAS com JWT)
	servicosGroup := protected.Group("/servicos")
	servicosGroup.Use(mw.UnitMiddleware())
	servicosGroup.POST("", servicoHandler.Create)                          // POST /api/v1/servicos
	servicosGroup.GET("", servicoHandler.List)                             // GET /api/v1/servicos
	servicosGroup.GET("/stats", servicoHandler.GetStats)                   // GET /api/v1/servicos/stats
	servicosGroup.GET("/:id", servicoHandler.GetByID)                      // GET /api/v1/servicos/:id
	servicosGroup.PUT("/:id", servicoHandler.Update)                       // PUT /api/v1/servicos/:id
	servicosGroup.DELETE("/:id", servicoHandler.Delete)                    // DELETE /api/v1/servicos/:id
	servicosGroup.PATCH("/:id/toggle-status", servicoHandler.ToggleStatus) // PATCH /api/v1/servicos/:id/toggle-status

	// Plans (assinaturas) routes
	plansGroup := protected.Group("/plans")
	plansGroup.GET("", planHandler.List, mw.RequireAdminAccess(logger))          // GET /api/v1/plans
	plansGroup.GET("/:id", planHandler.Get, mw.RequireAdminAccess(logger))       // GET /api/v1/plans/:id
	plansGroup.POST("", planHandler.Create, mw.RequireOwnerOrManager(logger))    // POST /api/v1/plans
	plansGroup.PUT("/:id", planHandler.Update, mw.RequireOwnerOrManager(logger)) // PUT /api/v1/plans/:id
	plansGroup.DELETE("/:id", planHandler.Deactivate, mw.RequireOwnerOrManager(logger))

	// Subscriptions routes
	subscriptionsGroup := protected.Group("/subscriptions")
	subscriptionsGroup.GET("/metrics", subscriptionHandler.Metrics, mw.RequireAdminAccess(logger))
	subscriptionsGroup.GET("", subscriptionHandler.List, mw.RequireAdminAccess(logger))
	subscriptionsGroup.GET("/:id", subscriptionHandler.Get, mw.RequireAdminAccess(logger))
	subscriptionsGroup.POST("", subscriptionHandler.Create, mw.RequireAdminAccess(logger))
	subscriptionsGroup.POST("/:id/renew", subscriptionHandler.Renew, mw.RequireAdminAccess(logger))
	subscriptionsGroup.POST("/reconcile", subscriptionHandler.Reconcile, mw.RequireOwnerOrManager(logger)) // T-ASAAS-002
	subscriptionsGroup.DELETE("/:id", subscriptionHandler.Cancel, mw.RequireOwnerOrManager(logger))

	// MeioPagamento routes - 6 endpoints (PROTEGIDAS com JWT)
	meiosPagamentoGroup := protected.Group("/meios-pagamento")
	meiosPagamentoGroup.POST("", meioPagamentoHandler.Create)             // POST /api/v1/meios-pagamento
	meiosPagamentoGroup.GET("", meioPagamentoHandler.List)                // GET /api/v1/meios-pagamento
	meiosPagamentoGroup.GET("/:id", meioPagamentoHandler.Get)             // GET /api/v1/meios-pagamento/:id
	meiosPagamentoGroup.PUT("/:id", meioPagamentoHandler.Update)          // PUT /api/v1/meios-pagamento/:id
	meiosPagamentoGroup.DELETE("/:id", meioPagamentoHandler.Delete)       // DELETE /api/v1/meios-pagamento/:id
	meiosPagamentoGroup.PATCH("/:id/toggle", meioPagamentoHandler.Toggle) // PATCH /api/v1/meios-pagamento/:id/toggle

	// Caixa Diário routes - 9 endpoints (PROTEGIDAS com JWT + ASSINATURA ATIVA)
	// T-ASAAS-003: Requer assinatura ativa (grupo guarded)
	caixaHandler.RegisterRoutes(guarded)

	// Commission routes - 35+ endpoints (PROTEGIDAS com JWT + RBAC)
	// Regras:
	// - BARBER pode ver suas próprias comissões e solicitar adiantamentos
	// - OWNER, MANAGER podem gerenciar todas as comissões, aprovar/rejeitar adiantamentos
	commissionsGroup := protected.Group("/commissions")
	commissionHandler.RegisterRoutes(commissionsGroup)

	// Placeholder endpoint
	api.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "pong",
		})
	})

	// Swagger documentation UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	logger.Info("🚀 Servidor iniciado", zap.String("port", port))
	if err := e.Start(":" + port); err != nil {
		logger.Fatal("Erro ao iniciar servidor", zap.Error(err))
	}
}
