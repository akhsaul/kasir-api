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
		categoryRepo = memory.NewCategoryRepository()
		productRepo = memory.NewProductRepository(categoryRepo)
	}

	// Service layer (logic)
	productService := service.NewProductService(productRepo, categoryRepo)
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

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), rt); err != nil {
		logger.Fatal(err)
	}
}
