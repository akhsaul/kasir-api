package main

import (
	"fmt"
	"kasir-api/handlers"
	"kasir-api/helpers"
	"kasir-api/repositories"
	"kasir-api/router"
	"kasir-api/services"
	"net/http"
	"os"
)

func main() {
	helper.InitLogger()

	// Repository layer (data)
	repo := repository.NewMemoryRepository()
	categoryRepo := repository.NewCategoryMemoryAdapter(repo)

	// Service layer (logic)
	productService := service.NewProductService(repo)
	categoryService := service.NewCategoryService(categoryRepo)

	// Handler layer (request/response)
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	rt := router.NewRouter(productHandler, categoryHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		helper.Info("Defaulting to port %s", port)
	}

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
