package main

import (
	"fmt"
	"kasir-api/data"
	"kasir-api/handler"
	"kasir-api/router"
	"kasir-api/service"
	"log"
	"net/http"
	"os"
)

func main() {
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
		log.Printf("Defaulting to port %s", port)
	}

	// Start server
	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET    /health")
	log.Printf("  GET    /api/product")
	log.Printf("  POST   /api/product")
	log.Printf("  GET    /api/product/{id}")
	log.Printf("  PUT    /api/product/{id}")
	log.Printf("  DELETE /api/product/{id}")
	log.Printf("  GET    /api/categories")
	log.Printf("  POST   /api/categories")
	log.Printf("  GET    /api/categories/{id}")
	log.Printf("  PUT    /api/categories/{id}")
	log.Printf("  DELETE /api/categories/{id}")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), rt))
}
