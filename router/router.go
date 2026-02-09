package router

import (
	"net/http"
	"strings"

	"kasir-api/docs"
	handler "kasir-api/handlers"
	helper "kasir-api/helpers"
)

// HealthChecker is an optional interface for checking database connectivity.
type HealthChecker interface {
	Ping() error
}

// Router handles routing logic.
type Router struct {
	productHandler     *handler.ProductHandler
	categoryHandler    *handler.CategoryHandler
	transactionHandler *handler.TransactionHandler
	healthChecker      HealthChecker
}

// NewRouter creates a new Router instance.
func NewRouter(productHandler *handler.ProductHandler, categoryHandler *handler.CategoryHandler, transactionHandler *handler.TransactionHandler) *Router {
	return &Router{
		productHandler:     productHandler,
		categoryHandler:    categoryHandler,
		transactionHandler: transactionHandler,
	}
}

// SetHealthChecker sets an optional database health checker.
func (rt *Router) SetHealthChecker(hc HealthChecker) {
	rt.healthChecker = hc
}

// ServeHTTP implements the http.Handler interface.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	// Health check endpoint
	if path == "/health" && method == http.MethodGet {
		rt.handleHealth(w, r)
		return
	}

	// API docs
	if path == "/docs" && method == http.MethodGet {
		docs.HandleSwaggerUI(w, r)
		return
	}
	if path == "/docs/openapi.yaml" && method == http.MethodGet {
		docs.HandleSpec(w, r)
		return
	}

	// Product endpoints
	if path == "/api/products" {
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
	if strings.HasPrefix(path, "/api/products/") && path != "/api/products/" {
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

	// Checkout endpoint
	if path == "/api/checkout" && method == http.MethodPost {
		rt.transactionHandler.HandleCheckout(w, r)
		return
	}

	// Transaction by ID endpoints
	if strings.HasPrefix(path, "/api/transactions/") && path != "/api/transactions/" {
		switch method {
		case http.MethodGet:
			rt.transactionHandler.HandleGetByID(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Report today endpoint
	if path == "/api/report/hari-ini" && method == http.MethodGet {
		rt.transactionHandler.HandleGetTodayReport(w, r)
		return
	}

	// Report with date range endpoint
	if path == "/api/report" && method == http.MethodGet {
		rt.transactionHandler.HandleGetReport(w, r)
		return
	}

	// Not found
	http.NotFound(w, r)
}

// handleHealth handles the health check endpoint with optional DB connectivity check.
func (rt *Router) handleHealth(w http.ResponseWriter, r *http.Request) {
	status := map[string]string{
		"api": "up",
	}

	if rt.healthChecker != nil {
		if err := rt.healthChecker.Ping(); err != nil {
			status["database"] = "down"
			helper.WriteError(w, r, http.StatusServiceUnavailable, "Database is not reachable", err)
			return
		}
		status["database"] = "up"
	}

	helper.WriteSuccess(w, http.StatusOK, "API Running", status)
}
