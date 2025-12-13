package main

import (
	"context"
	"log"
	"os"
	"strings"

	authUC "github.com/andviana23/barber-analytics-backend/internal/application/usecase/auth"
	"github.com/andviana23/barber-analytics-backend/internal/infra/auth"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/andviana23/barber-analytics-backend/internal/infra/http/handler"
	mw "github.com/andviana23/barber-analytics-backend/internal/infra/http/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
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
	poolCfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		logger.Fatal("Erro ao parsear DATABASE_URL", zap.Error(err))
	}
	if strings.EqualFold(os.Getenv("PGX_SIMPLE_PROTOCOL"), "true") {
		poolCfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	}
	dbPool, err := pgxpool.NewWithConfig(ctx, poolCfg)
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

	// Initialize JWT manager (lê JWT_SECRET do .env)
	jwtManager := auth.NewJWTManager()

	// Initialize use cases - Auth
	loginUC := authUC.NewLoginUseCase(queries, jwtManager, logger)
	refreshUC := authUC.NewRefreshUseCase(queries, jwtManager, logger)
	meUC := authUC.NewMeUseCase(queries, logger)
	logoutUC := authUC.NewLogoutUseCase(queries, logger)

	// Initialize handlers - Auth
	authHandler := handler.NewAuthHandler(
		loginUC,
		refreshUC,
		meUC,
		logoutUC,
		logger,
	)

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8000"},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Unit-ID"},
	}))

	// Health Check Endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status": "ok",
			"app":    "Barber Analytics Pro v2.0 - Auth Only",
		})
	})

	// API Routes
	api := e.Group("/api/v1")

	// Auth routes - PÚBLICAS (sem middleware JWT)
	authGroup := api.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/refresh", authHandler.Refresh)
	authGroup.POST("/logout", authHandler.Logout)
	authGroup.GET("/me", authHandler.Me, mw.JWTMiddleware(jwtManager, logger))

	// Start server
	logger.Info("Iniciando servidor",
		zap.String("port", port),
		zap.String("mode", "auth-only"),
	)

	if err := e.Start(":" + port); err != nil {
		logger.Fatal("Erro ao iniciar servidor", zap.Error(err))
	}
}
