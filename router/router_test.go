package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "kasir-api/handlers"
	model "kasir-api/models"
	"kasir-api/repositories/memory"
	service "kasir-api/services"
)

func setupTestRouter() *Router {
	// Create in-memory repositories
	categoryRepo := memory.NewCategoryRepository()
	productRepo := memory.NewProductRepository(categoryRepo)
	transactionRepo := memory.NewTransactionRepository()

	// Create services
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo)
	transactionService := service.NewTransactionService(transactionRepo, productRepo)

	// Create handlers
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Create router
	return NewRouter(productHandler, categoryHandler, transactionHandler)
}

func TestNewRouter(t *testing.T) {
	categoryRepo := memory.NewCategoryRepository()
	productRepo := memory.NewProductRepository(categoryRepo)
	transactionRepo := memory.NewTransactionRepository()

	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo)
	transactionService := service.NewTransactionService(transactionRepo, productRepo)

	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := NewRouter(productHandler, categoryHandler, transactionHandler)

	if router == nil {
		t.Error("NewRouter should return a non-nil router")
	}
	if router.productHandler != productHandler {
		t.Error("NewRouter should set productHandler")
	}
	if router.categoryHandler != categoryHandler {
		t.Error("NewRouter should set categoryHandler")
	}
	if router.transactionHandler != transactionHandler {
		t.Error("NewRouter should set transactionHandler")
	}
}

func TestRouter_HealthEndpoint(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Health endpoint should return 200, got: %d", rr.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	if response["status"] != "OK" {
		t.Errorf("Health endpoint should return status OK, got: %v", response["status"])
	}
}

func TestRouter_HealthEndpoint_WrongMethod(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Health endpoint with wrong method should return 404, got: %d", rr.Code)
	}
}

// Product endpoints tests
func TestRouter_Products_GetAll(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /api/products should return 200, got: %d", rr.Code)
	}
}

func TestRouter_Products_Create(t *testing.T) {
	router := setupTestRouter()

	product := map[string]interface{}{
		"name":  "Laptop",
		"price": 1000,
		"stock": 10,
	}
	body, _ := json.Marshal(product)

	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("POST /api/products should return 201, got: %d", rr.Code)
	}
}

func TestRouter_Products_MethodNotAllowed(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodDelete, "/api/products", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("DELETE /api/products should return 405, got: %d", rr.Code)
	}
}

func TestRouter_Products_GetByID(t *testing.T) {
	router := setupTestRouter()

	// First create a product
	product := map[string]interface{}{
		"name":  "Laptop",
		"price": 1000,
		"stock": 10,
	}
	body, _ := json.Marshal(product)

	createReq := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRr := httptest.NewRecorder()
	router.ServeHTTP(createRr, createReq)

	// Then get it by ID
	req := httptest.NewRequest(http.MethodGet, "/api/products/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /api/products/1 should return 200, got: %d", rr.Code)
	}
}

func TestRouter_Products_Update(t *testing.T) {
	router := setupTestRouter()

	// First create a product
	product := map[string]interface{}{
		"name":  "Laptop",
		"price": 1000,
		"stock": 10,
	}
	body, _ := json.Marshal(product)

	createReq := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRr := httptest.NewRecorder()
	router.ServeHTTP(createRr, createReq)

	// Then update it
	updateProduct := map[string]interface{}{
		"name":  "Updated Laptop",
		"price": 1500,
		"stock": 5,
	}
	updateBody, _ := json.Marshal(updateProduct)

	req := httptest.NewRequest(http.MethodPut, "/api/products/1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PUT /api/products/1 should return 200, got: %d", rr.Code)
	}
}

func TestRouter_Products_Delete(t *testing.T) {
	router := setupTestRouter()

	// First create a product
	product := map[string]interface{}{
		"name":  "Laptop",
		"price": 1000,
		"stock": 10,
	}
	body, _ := json.Marshal(product)

	createReq := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRr := httptest.NewRecorder()
	router.ServeHTTP(createRr, createReq)

	// Then delete it
	req := httptest.NewRequest(http.MethodDelete, "/api/products/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("DELETE /api/products/1 should return 200, got: %d", rr.Code)
	}
}

func TestRouter_ProductsByID_MethodNotAllowed(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/products/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /api/products/1 should return 405, got: %d", rr.Code)
	}
}

// Category endpoints tests
func TestRouter_Categories_GetAll(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /api/categories should return 200, got: %d", rr.Code)
	}
}

func TestRouter_Categories_Create(t *testing.T) {
	router := setupTestRouter()

	category := map[string]interface{}{
		"name":        "Electronics",
		"description": "Electronic items",
	}
	body, _ := json.Marshal(category)

	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("POST /api/categories should return 201, got: %d", rr.Code)
	}
}

func TestRouter_Categories_MethodNotAllowed(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodDelete, "/api/categories", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("DELETE /api/categories should return 405, got: %d", rr.Code)
	}
}

func TestRouter_Categories_GetByID(t *testing.T) {
	router := setupTestRouter()

	// First create a category
	category := map[string]interface{}{
		"name":        "Electronics",
		"description": "Electronic items",
	}
	body, _ := json.Marshal(category)

	createReq := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRr := httptest.NewRecorder()
	router.ServeHTTP(createRr, createReq)

	// Then get it by ID
	req := httptest.NewRequest(http.MethodGet, "/api/categories/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /api/categories/1 should return 200, got: %d", rr.Code)
	}
}

func TestRouter_Categories_Update(t *testing.T) {
	router := setupTestRouter()

	// First create a category
	category := map[string]interface{}{
		"name":        "Electronics",
		"description": "Electronic items",
	}
	body, _ := json.Marshal(category)

	createReq := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRr := httptest.NewRecorder()
	router.ServeHTTP(createRr, createReq)

	// Then update it
	updateCategory := map[string]interface{}{
		"name":        "Updated Electronics",
		"description": "Updated description",
	}
	updateBody, _ := json.Marshal(updateCategory)

	req := httptest.NewRequest(http.MethodPut, "/api/categories/1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PUT /api/categories/1 should return 200, got: %d", rr.Code)
	}
}

func TestRouter_Categories_Delete(t *testing.T) {
	router := setupTestRouter()

	// First create a category
	category := map[string]interface{}{
		"name":        "Electronics",
		"description": "Electronic items",
	}
	body, _ := json.Marshal(category)

	createReq := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRr := httptest.NewRecorder()
	router.ServeHTTP(createRr, createReq)

	// Then delete it
	req := httptest.NewRequest(http.MethodDelete, "/api/categories/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("DELETE /api/categories/1 should return 200, got: %d", rr.Code)
	}
}

func TestRouter_CategoriesByID_MethodNotAllowed(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/categories/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /api/categories/1 should return 405, got: %d", rr.Code)
	}
}

// Transaction endpoints tests
func TestRouter_Checkout(t *testing.T) {
	router := setupTestRouter()

	// First create a product
	product := map[string]interface{}{
		"name":  "Laptop",
		"price": 1000,
		"stock": 10,
	}
	productBody, _ := json.Marshal(product)

	createReq := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(productBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRr := httptest.NewRecorder()
	router.ServeHTTP(createRr, createReq)

	// Then checkout
	checkout := model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 1, Quantity: 2},
		},
	}
	checkoutBody, _ := json.Marshal(checkout)

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(checkoutBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("POST /api/checkout should return 201, got: %d", rr.Code)
	}
}

func TestRouter_Checkout_WrongMethod(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/checkout", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("GET /api/checkout should return 404, got: %d", rr.Code)
	}
}

func TestRouter_Transactions_GetByID(t *testing.T) {
	router := setupTestRouter()

	// First create a product and checkout
	product := map[string]interface{}{
		"name":  "Laptop",
		"price": 1000,
		"stock": 10,
	}
	productBody, _ := json.Marshal(product)

	createReq := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(productBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRr := httptest.NewRecorder()
	router.ServeHTTP(createRr, createReq)

	checkout := model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 1, Quantity: 2},
		},
	}
	checkoutBody, _ := json.Marshal(checkout)

	checkoutReq := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(checkoutBody))
	checkoutReq.Header.Set("Content-Type", "application/json")
	checkoutRr := httptest.NewRecorder()
	router.ServeHTTP(checkoutRr, checkoutReq)

	// Then get transaction by ID
	req := httptest.NewRequest(http.MethodGet, "/api/transactions/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /api/transactions/1 should return 200, got: %d", rr.Code)
	}
}

func TestRouter_TransactionsByID_MethodNotAllowed(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/transactions/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /api/transactions/1 should return 405, got: %d", rr.Code)
	}
}

func TestRouter_Report_Today(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/report/hari-ini", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /api/report/hari-ini should return 200, got: %d", rr.Code)
	}
}

func TestRouter_Report_Today_WrongMethod(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/report/hari-ini", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("POST /api/report/hari-ini should return 404, got: %d", rr.Code)
	}
}

func TestRouter_Report_DateRange(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/report?start_date=2024-01-01&end_date=2024-01-31", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /api/report should return 200, got: %d", rr.Code)
	}
}

func TestRouter_Report_DateRange_WrongMethod(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/report", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("POST /api/report should return 404, got: %d", rr.Code)
	}
}

func TestRouter_NotFound(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Unknown endpoint should return 404, got: %d", rr.Code)
	}
}

func TestRouter_ProductsPath_TrailingSlash(t *testing.T) {
	router := setupTestRouter()

	// /api/products/ should match the products collection, not products by ID
	req := httptest.NewRequest(http.MethodGet, "/api/products/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// The path /api/products/ does not match /api/products and not /api/products/{id}
	// so it should be 404
	if rr.Code != http.StatusNotFound {
		t.Errorf("GET /api/products/ should return 404, got: %d", rr.Code)
	}
}

func TestRouter_CategoriesPath_TrailingSlash(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/categories/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("GET /api/categories/ should return 404, got: %d", rr.Code)
	}
}

func TestRouter_TransactionsPath_TrailingSlash(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("GET /api/transactions/ should return 404, got: %d", rr.Code)
	}
}
