package main

import (
	"fmt"
	"kasir-api/data"
	"kasir-api/handler"
	"kasir-api/helper"
	"kasir-api/router"
	"kasir-api/service"
	"net/http"
	"os"
)

func main() {
	// Initialize logger
	helper.InitLogger()

	// Initialize layers
	storage := data.NewMemoryStorage()

	// Product components
	productService := service.NewProductService(storage)
	productHandler := handler.NewProductHandler(productService)

	// Category components
	categoryStorage := data.NewCategoryMemoryStorage(storage)
	categoryService := service.NewCategoryService(categoryStorage)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Setup router
	rt := router.NewRouter(productHandler, categoryHandler)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		helper.Info("Defaulting to port %s", port)
	}

	// Start server
	helper.Info("Listening on port %s", port)
	helper.Info("Open http://localhost:%s in the browser", port)
	helper.Info("Available endpoints:")
	helper.Info("  GET     /health")
	helper.Info("  GET     /api/product")
	helper.Info("  POST    /api/product")
	helper.Info("  GET     /api/product/{id}")
	helper.Info("  PUT     /api/product/{id}")
	helper.Info("  DELETE  /api/product/{id}")
	helper.Info("  GET     /api/categories")
	helper.Info("  POST    /api/categories")
	helper.Info("  GET     /api/categories/{id}")
	helper.Info("  PUT     /api/categories/{id}")
	helper.Info("  DELETE  /api/categories/{id}")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), rt); err != nil {
		helper.Fatal(err)
	}
}
