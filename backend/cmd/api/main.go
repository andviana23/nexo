package main

import (
	"context"
	"log"
	"os"

	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/appointment"
	authUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/auth"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/financial"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/metas"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/pricing"
	"github.com/andviana23/barber-analytics-backend/internal/application/usecase/stock"
	"github.com/andviana23/barber-analytics-backend/internal/infra/auth"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/handler"
	mw "github.com/andviana23/barber-analytics-backend/internal/infra/http/middleware"
	"github.com/andviana23/barber-analytics-backend/internal/infra/repository/postgres"
	"github.com/andviana23/barber-analytics-backend/internal/infra/scheduler"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	// Load .env file (ignore error in production where env vars are set directly)
	_ = godotenv.Load()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Erro ao inicializar logger: %v", err)
	}
	defer logger.Sync()

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

	// Stock repositories
	produtoRepo := postgres.NewProdutoRepository(queries)
	fornecedorRepo := postgres.NewFornecedorRepositoryPG(queries)
	movimentacaoRepo := postgres.NewMovimentacaoEstoqueRepositoryPG(queries)

	// Appointment repositories and readers
	appointmentRepo := postgres.NewAppointmentRepository(queries, dbPool)
	professionalReader := postgres.NewProfessionalReader(queries)
	customerReader := postgres.NewCustomerReader(queries)
	serviceReader := postgres.NewServiceReader(queries)

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
	marcarCompensacaoUC := financial.NewMarcarCompensacaoUseCase(compensacaoRepo, logger)
	// FluxoCaixa (com dependências de ContaPagar e ContaReceber)
	generateFluxoDiarioUC := financial.NewGenerateFluxoDiarioUseCase(fluxoCaixaRepo, contaPagarRepo, contaReceberRepo, logger)
	getFluxoCaixaUC := financial.NewGetFluxoCaixaUseCase(fluxoCaixaRepo, logger)
	listFluxoCaixaUC := financial.NewListFluxoCaixaUseCase(fluxoCaixaRepo, logger)
	// DRE (com dependências de ContaPagar e ContaReceber)
	generateDREUC := financial.NewGenerateDREUseCase(dreRepo, contaPagarRepo, contaReceberRepo, logger)
	getDREUC := financial.NewGetDREUseCase(dreRepo, logger)
	listDREUC := financial.NewListDREUseCase(dreRepo, logger)

	// Initialize use cases - Stock (4 use cases)
	registrarEntradaUC := stock.NewRegistrarEntradaUseCase(produtoRepo, movimentacaoRepo, fornecedorRepo)
	registrarSaidaUC := stock.NewRegistrarSaidaUseCase(produtoRepo, movimentacaoRepo)
	ajustarEstoqueUC := stock.NewAjustarEstoqueUseCase(produtoRepo, movimentacaoRepo)
	listarAlertasUC := stock.NewListarAlertasEstoqueBaixoUseCase(produtoRepo)

	// Initialize use cases - Appointments (6 use cases)
	createAppointmentUC := appointment.NewCreateAppointmentUseCase(appointmentRepo, serviceReader, professionalReader, customerReader, logger)
	listAppointmentsUC := appointment.NewListAppointmentsUseCase(appointmentRepo, logger)
	getAppointmentUC := appointment.NewGetAppointmentUseCase(appointmentRepo, logger)
	updateAppointmentStatusUC := appointment.NewUpdateAppointmentStatusUseCase(appointmentRepo, logger)
	rescheduleAppointmentUC := appointment.NewRescheduleAppointmentUseCase(appointmentRepo, logger)
	cancelAppointmentUC := appointment.NewCancelAppointmentUseCase(appointmentRepo, logger)

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
		GenerateDRE:         generateDREUC,
		GenerateFluxoDiario: generateFluxoDiarioUC,
		MarcarCompensacoes:  marcarCompensacaoUC,
	}

	// Parse tenant list from ENV (SCHEDULER_TENANTS="tenant1,tenant2,...")
	tenants := scheduler.ParseTenantEnv("SCHEDULER_TENANTS")

	scheduler.RegisterFinancialJobs(sched, logger, financialDeps, tenants)

	// Start scheduler in background
	sched.Start()
	defer sched.Stop(ctx)

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

	// Initialize handlers - Financial (23 use cases)
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
		nil, // getDashboardUC TODO: implementar
		logger,
	)

	// Initialize handlers - Stock (4 use cases)
	stockHandler := handler.NewStockHandler(
		registrarEntradaUC,
		registrarSaidaUC,
		ajustarEstoqueUC,
		listarAlertasUC,
	)

	// Initialize handlers - Auth (4 use cases)
	authHandler := handler.NewAuthHandler(
		loginUC,
		refreshUC,
		meUC,
		logoutUC,
		logger,
	)

	// Initialize handlers - Appointments (6 use cases)
	appointmentHandler := handler.NewAppointmentHandler(
		createAppointmentUC,
		listAppointmentsUC,
		getAppointmentUC,
		updateAppointmentStatusUC,
		rescheduleAppointmentUC,
		cancelAppointmentUC,
		logger,
	)

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:3002", "http://localhost:3006", "http://localhost:8000"},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowCredentials: true, // Permite cookies (refresh token)
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
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

	// Middleware JWT para rotas protegidas
	protected := api.Group("")
	protected.Use(mw.JWTMiddleware(jwtManager, logger))

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

	// Appointment routes - 6 endpoints (PROTEGIDAS com JWT)
	appointmentsGroup := protected.Group("/appointments")
	appointmentsGroup.POST("", appointmentHandler.CreateAppointment)
	appointmentsGroup.GET("", appointmentHandler.ListAppointments)
	appointmentsGroup.GET("/:id", appointmentHandler.GetAppointment)
	appointmentsGroup.PATCH("/:id/status", appointmentHandler.UpdateAppointmentStatus)
	appointmentsGroup.PATCH("/:id/reschedule", appointmentHandler.RescheduleAppointment)
	appointmentsGroup.POST("/:id/cancel", appointmentHandler.CancelAppointment)

	// Financial routes - 19 endpoints (PROTEGIDAS com JWT)
	financialGroup := protected.Group("/financial")

	// ContaPagar (6 endpoints: 5 CRUD + 1 marcarPagamento)
	financialGroup.POST("/payables", financialHandler.CreateContaPagar)
	financialGroup.GET("/payables/:id", financialHandler.GetContaPagar)
	financialGroup.GET("/payables", financialHandler.ListContasPagar)
	financialGroup.PUT("/payables/:id", financialHandler.UpdateContaPagar)
	financialGroup.DELETE("/payables/:id", financialHandler.DeleteContaPagar)
	financialGroup.POST("/payables/:id/payment", financialHandler.MarcarPagamento)

	// ContaReceber (6 endpoints: 5 CRUD + 1 marcarRecebimento)
	financialGroup.POST("/receivables", financialHandler.CreateContaReceber)
	financialGroup.GET("/receivables/:id", financialHandler.GetContaReceber)
	financialGroup.GET("/receivables", financialHandler.ListContasReceber)
	financialGroup.PUT("/receivables/:id", financialHandler.UpdateContaReceber)
	financialGroup.DELETE("/receivables/:id", financialHandler.DeleteContaReceber)
	financialGroup.POST("/receivables/:id/receipt", financialHandler.MarcarRecebimento)

	// Compensação (3 endpoints: Get, List, Delete)
	financialGroup.GET("/compensations/:id", financialHandler.GetCompensacao)
	financialGroup.GET("/compensations", financialHandler.ListCompensacoes)
	financialGroup.DELETE("/compensations/:id", financialHandler.DeleteCompensacao)

	// FluxoCaixa (2 endpoints: Get, List)
	financialGroup.GET("/cashflow/:id", financialHandler.GetFluxoCaixa)
	financialGroup.GET("/cashflow", financialHandler.ListFluxoCaixa)

	// DRE (2 endpoints: Get, List)
	financialGroup.GET("/dre/:month", financialHandler.GetDRE)
	financialGroup.GET("/dre", financialHandler.ListDRE)

	// Stock routes - 4 endpoints
	stockGroup := api.Group("/stock")
	stockGroup.POST("/entries", stockHandler.RegistrarEntrada) // Registrar entrada
	stockGroup.POST("/exit", stockHandler.RegistrarSaida)      // Registrar saída
	stockGroup.POST("/adjust", stockHandler.AjustarEstoque)    // Ajustar estoque
	stockGroup.GET("/alerts", stockHandler.ListarAlertas)      // Listar alertas

	// Placeholder endpoint
	api.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "pong",
		})
	})

	// Start server
	logger.Info("🚀 Servidor iniciado", zap.String("port", port))
	if err := e.Start(":" + port); err != nil {
		logger.Fatal("Erro ao iniciar servidor", zap.Error(err))
	}
}
