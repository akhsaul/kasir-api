package main

import (
	"fmt"
	"kasir-api/config"
	"kasir-api/handlers"
	"kasir-api/helpers/logger"
	"kasir-api/repositories"
	"kasir-api/repositories/memory"
	"kasir-api/repositories/postgres"
	"kasir-api/router"
	"kasir-api/services"
	"net/http"
)

func main() {
	logger.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err)
	}

	var productRepo repository.ProductRepository
	var categoryRepo repository.CategoryRepository

	if cfg.DB.Enabled {
		pgDB, err := postgres.NewDB(&cfg.DB)
		if err != nil {
			logger.Fatal(err)
		}
		defer pgDB.Close()
		logger.Info("Using PostgreSQL storage")
		productRepo = postgres.NewProductRepository(pgDB)
		categoryRepo = postgres.NewCategoryRepository(pgDB)
	} else {
		logger.Info("Using in-memory storage")
		productRepo = memory.NewProductRepository()
		categoryRepo = memory.NewCategoryRepository()
	}

	// Service layer (logic)
	productService := service.NewProductService(productRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// Handler layer (request/response)
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	rt := router.NewRouter(productHandler, categoryHandler)

	port := cfg.Server.Port
	if port == "" {
		port = "8080"
		logger.Info("Defaulting to port %s", port)
	}

	logger.Info("Listening on port %s", port)
	logger.Info("Open http://localhost:%s in the browser", port)
	logger.Info("Available endpoints:")
	logger.Info("  GET     /health")
	logger.Info("  GET     /api/product")
	logger.Info("  POST    /api/product")
	logger.Info("  GET     /api/product/{id}")
	logger.Info("  PUT     /api/product/{id}")
	logger.Info("  DELETE  /api/product/{id}")
	logger.Info("  GET     /api/categories")
	logger.Info("  POST    /api/categories")
	logger.Info("  GET     /api/categories/{id}")
	logger.Info("  PUT     /api/categories/{id}")
	logger.Info("  DELETE  /api/categories/{id}")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), rt); err != nil {
		logger.Fatal(err)
	}
}
