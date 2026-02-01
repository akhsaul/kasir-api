package router

import (
	"kasir-api/handlers"
	"kasir-api/helpers"
	"net/http"
	"strings"
)

// Router handles routing logic
type Router struct {
	productHandler  *handler.ProductHandler
	categoryHandler *handler.CategoryHandler
}

// NewRouter creates a new Router instance
func NewRouter(productHandler *handler.ProductHandler, categoryHandler *handler.CategoryHandler) *Router {
	return &Router{
		productHandler:  productHandler,
		categoryHandler: categoryHandler,
	}
}

// ServeHTTP implements the http.Handler interface
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	// Health check endpoint
	if path == "/health" && method == http.MethodGet {
		rt.handleHealth(w)
		return
	}

	// Product endpoints
	if path == "/api/product" {
		switch method {
		case http.MethodGet:
			rt.productHandler.HandleGetAll(w, r)
		case http.MethodPost:
			rt.productHandler.HandleCreate(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Product by ID endpoints
	if strings.HasPrefix(path, "/api/product/") && path != "/api/product/" {
		switch method {
		case http.MethodGet:
			rt.productHandler.HandleGetByID(w, r)
		case http.MethodPut:
			rt.productHandler.HandleUpdate(w, r)
		case http.MethodDelete:
			rt.productHandler.HandleDelete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Category endpoints
	if path == "/api/categories" {
		switch method {
		case http.MethodGet:
			rt.categoryHandler.HandleGetAll(w, r)
		case http.MethodPost:
			rt.categoryHandler.HandleCreate(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Category by ID endpoints
	if strings.HasPrefix(path, "/api/categories/") && path != "/api/categories/" {
		switch method {
		case http.MethodGet:
			rt.categoryHandler.HandleGetByID(w, r)
		case http.MethodPut:
			rt.categoryHandler.HandleUpdate(w, r)
		case http.MethodDelete:
			rt.categoryHandler.HandleDelete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Not found
	http.NotFound(w, r)
}

// handleHealth handles the health check endpoint
func (rt *Router) handleHealth(w http.ResponseWriter) {
	helper.WriteSuccess(w, http.StatusOK, "API Running", nil)
}
