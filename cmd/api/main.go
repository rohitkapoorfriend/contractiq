package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/contractiq/contractiq/internal/application/contract/command"
	"github.com/contractiq/contractiq/internal/application/contract/query"
	appparty "github.com/contractiq/contractiq/internal/application/party"
	apptemplate "github.com/contractiq/contractiq/internal/application/template"
	"github.com/contractiq/contractiq/internal/infrastructure/auth"
	"github.com/contractiq/contractiq/internal/infrastructure/config"
	"github.com/contractiq/contractiq/internal/infrastructure/eventbus"
	"github.com/contractiq/contractiq/internal/infrastructure/persistence/postgres"
	apphttp "github.com/contractiq/contractiq/internal/interfaces/http"
	"github.com/contractiq/contractiq/internal/interfaces/http/handler"
	"github.com/contractiq/contractiq/pkg/clock"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize logger
	var logger *zap.Logger
	if cfg.App.Env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer func() { _ = logger.Sync() }()

	// Connect to database
	db, err := postgres.NewConnection(cfg.DB)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("connected to database")

	// Initialize infrastructure
	clk := clock.New()
	uow := postgres.NewUnitOfWork(db)
	publisher := eventbus.NewInMemoryPublisher(logger)
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiry)

	// Initialize repositories (for query handlers that read outside transactions)
	contractRepo := postgres.NewContractRepository(db)
	templateRepo := postgres.NewTemplateRepository(db)
	partyRepo := postgres.NewPartyRepository(db)

	// Initialize application services
	contractCmdHandler := command.NewHandler(uow, clk, publisher)
	contractQueryHandler := query.NewHandler(contractRepo)
	templateService := apptemplate.NewService(templateRepo, clk)
	partyService := appparty.NewService(partyRepo, clk)
	userService := auth.NewUserService(db, jwtService)

	// Initialize HTTP handlers
	contractHandler := handler.NewContractHandler(contractCmdHandler, contractQueryHandler)
	templateHandler := handler.NewTemplateHandler(templateService)
	partyHandler := handler.NewPartyHandler(partyService)
	authHandler := handler.NewAuthHandler(userService)
	healthHandler := handler.NewHealthHandler(db)

	// Build router
	router := apphttp.NewRouter(apphttp.RouterDeps{
		Logger:          logger,
		JWTService:      jwtService,
		ContractHandler: contractHandler,
		TemplateHandler: templateHandler,
		PartyHandler:    partyHandler,
		AuthHandler:     authHandler,
		HealthHandler:   healthHandler,
		AllowedOrigins:  cfg.CORS.AllowedOrigins,
	})

	// Start server
	server := apphttp.NewServer(
		router,
		cfg.Server.Host,
		cfg.Server.Port,
		cfg.Server.ReadTimeout,
		cfg.Server.WriteTimeout,
		logger,
	)

	// Graceful shutdown
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("received shutdown signal", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server shutdown failed", zap.Error(err))
	}

	logger.Info("server stopped gracefully")
}
