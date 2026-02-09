package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kasir-api/config"
	handler "kasir-api/handlers"
	"kasir-api/helpers/logger"
	"kasir-api/middleware"
	repository "kasir-api/repositories"
	"kasir-api/repositories/memory"
	"kasir-api/repositories/postgres"
	"kasir-api/router"
	service "kasir-api/services"
)

func main() {
	logger.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err)
	}

	var productRepo repository.ProductRepository
	var categoryRepo repository.CategoryRepository
	var transactionRepo repository.TransactionRepository
	var pgDB *postgres.DB

	if cfg.DB.Enabled {
		var err error
		pgDB, err = postgres.NewDB(&cfg.DB)
		if err != nil {
			logger.Fatal(err)
		}
		defer pgDB.Close()
		logger.Info("Using PostgreSQL storage")
		productRepo = postgres.NewProductRepository(pgDB)
		categoryRepo = postgres.NewCategoryRepository(pgDB)
		transactionRepo = postgres.NewTransactionRepository(pgDB)
	} else {
		logger.Info("Using in-memory storage")
		categoryRepo = memory.NewCategoryRepository()
		productRepo = memory.NewProductRepository(categoryRepo)
		transactionRepo = memory.NewTransactionRepository()
	}

	// Service layer (logic)
	productService := service.NewProductService(productRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	transactionService := service.NewTransactionService(transactionRepo, productRepo)

	// Handler layer (request/response)
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	rt := router.NewRouter(productHandler, categoryHandler, transactionHandler)

	if pgDB != nil {
		rt.SetHealthChecker(pgDB)
	}

	// Apply middleware (outermost runs first).
	limiter := middleware.NewRateLimiter(cfg.RateLimit.Rate, cfg.RateLimit.Burst)
	var mux http.Handler = rt
	mux = middleware.RequestLogger(mux)      // log after request completes
	mux = middleware.BodyLimit(1 << 20)(mux) // 1MB max request body
	mux = limiter.Limit(mux)                 // rate limit before processing
	mux = middleware.RequestID(mux)          // assign request ID first
	mux = middleware.RecoverPanic(mux)       // recover from panics

	port := cfg.Server.Port

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Listening on port %s", port)
		logger.Info("Open http://localhost:%s in the browser", port)
		logger.Info("Available endpoints:")
		logger.Info("  GET     /docs                (Swagger UI)")
		logger.Info("  GET     /docs/openapi.yaml   (OpenAPI spec)")
		logger.Info("  GET     /health")
		logger.Info("  GET     /api/products")
		logger.Info("  POST    /api/products")
		logger.Info("  GET     /api/products/{id}")
		logger.Info("  PUT     /api/products/{id}")
		logger.Info("  DELETE  /api/products/{id}")
		logger.Info("  GET     /api/categories")
		logger.Info("  POST    /api/categories")
		logger.Info("  GET     /api/categories/{id}")
		logger.Info("  PUT     /api/categories/{id}")
		logger.Info("  DELETE  /api/categories/{id}")
		logger.Info("  POST    /api/checkout")
		logger.Info("  GET     /api/transactions/{id}")
		logger.Info("  GET     /api/report/hari-ini")
		logger.Info("  GET     /api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	// Graceful shutdown: wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("Received signal %v, shutting down gracefully...", sig)

	// Give active connections 30 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Server stopped")
}
