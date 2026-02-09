package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	model "kasir-api/models"
	"kasir-api/repositories/memory"
	service "kasir-api/services"
)

func setupTransactionHandler() (*TransactionHandler, *memory.TransactionRepository, *memory.ProductRepository, *memory.CategoryRepository) {
	categoryRepo := memory.NewCategoryRepository()
	productRepo := memory.NewProductRepository(categoryRepo)
	transactionRepo := memory.NewTransactionRepository()
	svc := service.NewTransactionService(transactionRepo, productRepo)
	handler := NewTransactionHandler(svc)
	return handler, transactionRepo, productRepo, categoryRepo
}

func TestNewTransactionHandler(t *testing.T) {
	svc := &service.TransactionService{}
	handler := NewTransactionHandler(svc)

	if handler == nil {
		t.Error("NewTransactionHandler should return a non-nil handler")
	}
	if handler.service != svc {
		t.Error("NewTransactionHandler should set the service")
	}
}

func TestTransactionHandler_HandleCheckout_Success(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()

	productRepo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 10})
	productRepo.Create(&model.Product{Name: "Phone", Price: 500, Stock: 20})

	request := model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
	body, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCheckout(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("HandleCheckout should return 201, got: %d", rr.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	if response["status"] != "OK" {
		t.Errorf("Status should be OK, got: %v", response["status"])
	}
}

func TestTransactionHandler_HandleCheckout_InvalidJSON(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCheckout(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleCheckout with invalid JSON should return 400, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleCheckout_EmptyItems(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	request := model.CheckoutRequest{Items: []model.CheckoutItem{}}
	body, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCheckout(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleCheckout with empty items should return 400, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleCheckout_InvalidQuantity(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()

	productRepo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 10})

	request := model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 1, Quantity: 0},
		},
	}
	body, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCheckout(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleCheckout with invalid quantity should return 400, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleCheckout_ProductNotFound(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	request := model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 999, Quantity: 1},
		},
	}
	body, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCheckout(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleCheckout with non-existent product should return 404, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleCheckout_InsufficientStock(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()

	productRepo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 5})

	request := model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 1, Quantity: 10},
		},
	}
	body, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCheckout(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleCheckout with insufficient stock should return 400, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleGetByID_Success(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()

	productRepo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 10})

	// First create a transaction via checkout
	request := model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 1, Quantity: 2},
		},
	}
	body, _ := json.Marshal(request)

	checkoutReq := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
	checkoutReq.Header.Set("Content-Type", "application/json")
	checkoutRr := httptest.NewRecorder()
	handler.HandleCheckout(checkoutRr, checkoutReq)

	// Then get by ID
	req := httptest.NewRequest(http.MethodGet, "/api/transactions/1", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleGetByID should return 200, got: %d", rr.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	if response["status"] != "OK" {
		t.Errorf("Status should be OK, got: %v", response["status"])
	}
}

func TestTransactionHandler_HandleGetByID_InvalidID(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/invalid", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleGetByID with invalid ID should return 404, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleGetByID_NotFound(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/999", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleGetByID for non-existent transaction should return 404, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleGetByID_ZeroID(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/0", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleGetByID with zero ID should return 404, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleGetTodayReport_Success(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/report/hari-ini", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetTodayReport(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleGetTodayReport should return 200, got: %d", rr.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	if response["status"] != "OK" {
		t.Errorf("Status should be OK, got: %v", response["status"])
	}
}

func TestTransactionHandler_HandleGetReport_Success(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/report?start_date=2024-01-01&end_date=2024-01-31", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetReport(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleGetReport should return 200, got: %d", rr.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	if response["status"] != "OK" {
		t.Errorf("Status should be OK, got: %v", response["status"])
	}
}

func TestTransactionHandler_HandleGetReport_MissingStartDate(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/report?end_date=2024-01-31", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetReport(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleGetReport without start_date should return 400, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleGetReport_MissingEndDate(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/report?start_date=2024-01-01", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetReport(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleGetReport without end_date should return 400, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleGetReport_MissingBothDates(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/report", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetReport(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleGetReport without dates should return 400, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleGetReport_InvalidStartDate(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/report?start_date=invalid&end_date=2024-01-31", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetReport(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleGetReport with invalid start_date should return 400, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleGetReport_InvalidEndDate(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/report?start_date=2024-01-01&end_date=invalid", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetReport(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleGetReport with invalid end_date should return 400, got: %d", rr.Code)
	}
}

func TestTransactionHandler_HandleGetReport_EndBeforeStart(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/report?start_date=2024-01-31&end_date=2024-01-01", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetReport(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleGetReport with end_date before start_date should return 400, got: %d", rr.Code)
	}
}

// ==================== BOUNDARY VALUE TESTS ====================

func TestTransactionHandler_BoundaryValues_ID(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()
	productRepo.Create(&model.Product{Name: "Test", Price: 100, Stock: 100})

	// Create a transaction first
	checkoutReq := model.CheckoutRequest{
		Items: []model.CheckoutItem{{ProductID: 1, Quantity: 1}},
	}
	body, _ := json.Marshal(checkoutReq)
	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.HandleCheckout(rr, req)

	testCases := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{"negative id", "-1", http.StatusNotFound},
		{"zero id", "0", http.StatusNotFound},
		{"float id", "1.5", http.StatusNotFound},
		{"string id", "abc", http.StatusNotFound},
		{"empty id", "", http.StatusNotFound},
		{"max int overflow", "9999999999999999999", http.StatusNotFound},
		{"valid id", "1", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/transactions/"+tc.id, nil)
			rr := httptest.NewRecorder()

			handler.HandleGetByID(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("ID '%s': expected %d, got %d", tc.id, tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestTransactionHandler_BoundaryValues_Quantity(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()
	productRepo.Create(&model.Product{Name: "Test", Price: 100, Stock: 1000})

	testCases := []struct {
		name           string
		quantity       int
		expectedStatus int
	}{
		{"zero quantity", 0, http.StatusBadRequest},
		{"negative quantity", -1, http.StatusBadRequest},
		{"negative large", -999999999, http.StatusBadRequest},
		{"valid quantity", 1, http.StatusCreated},
		{"large valid quantity", 100, http.StatusCreated},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset product stock for each test by updating
			product, _ := productRepo.GetByID(1)
			if product != nil {
				product.Stock = 1000
				productRepo.Update(product)
			}

			input := model.CheckoutRequest{
				Items: []model.CheckoutItem{{ProductID: 1, Quantity: tc.quantity}},
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCheckout(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Quantity %d: expected %d, got %d", tc.quantity, tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestTransactionHandler_BoundaryValues_ProductID(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()
	productRepo.Create(&model.Product{Name: "Test", Price: 100, Stock: 100})

	testCases := []struct {
		name           string
		productID      int
		expectedStatus int
	}{
		{"zero product id", 0, http.StatusBadRequest},
		{"negative product id", -1, http.StatusBadRequest},
		{"non-existent product id", 999, http.StatusNotFound},
		{"valid product id", 1, http.StatusCreated},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset product stock
			product, _ := productRepo.GetByID(1)
			if product != nil {
				product.Stock = 100
				productRepo.Update(product)
			}

			input := model.CheckoutRequest{
				Items: []model.CheckoutItem{{ProductID: tc.productID, Quantity: 1}},
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCheckout(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("ProductID %d: expected %d, got %d", tc.productID, tc.expectedStatus, rr.Code)
			}
		})
	}
}

// ==================== DATE FORMAT TESTS ====================

func TestTransactionHandler_DateFormatValidation(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	testCases := []struct {
		name           string
		startDate      string
		endDate        string
		expectedStatus int
	}{
		// Invalid date formats
		{"invalid start format dd-mm-yyyy", "01-01-2024", "2024-01-31", http.StatusBadRequest},
		{"invalid end format dd-mm-yyyy", "2024-01-01", "31-01-2024", http.StatusBadRequest},
		{"invalid start unix", "1704067200", "2024-01-31", http.StatusBadRequest},
		{"invalid date values", "2024-13-01", "2024-01-31", http.StatusBadRequest},
		{"invalid day", "2024-01-32", "2024-01-31", http.StatusBadRequest},
		// Valid dates
		{"valid dates", "2024-01-01", "2024-01-31", http.StatusOK},
		{"same date", "2024-01-15", "2024-01-15", http.StatusOK},
		{"leap year", "2024-02-29", "2024-03-01", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/api/report?start_date=" + tc.startDate + "&end_date=" + tc.endDate
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rr := httptest.NewRecorder()

			handler.HandleGetReport(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("%s: expected %d, got %d", tc.name, tc.expectedStatus, rr.Code)
			}
		})
	}
}

// ==================== SQL INJECTION TESTS ====================

func TestTransactionHandler_SQLInjection_InDates(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	// URL-encoded SQL injection payloads
	sqlPayloads := []string{
		"2024-01-01%27%3BDROP%20TABLE",
		"2024-01-01%20OR%201%3D1",
		"2024-01-01%27UNION%20SELECT",
		"2024-01-01%3BDELETE%20FROM",
	}

	for _, payload := range sqlPayloads {
		t.Run("sqli_date", func(t *testing.T) {
			url := "/api/report?start_date=" + payload + "&end_date=2024-01-31"
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rr := httptest.NewRecorder()

			handler.HandleGetReport(rr, req)

			// Should be bad request (invalid date), not server error
			if rr.Code == http.StatusInternalServerError {
				t.Errorf("SQL injection in date caused server error")
			}
		})
	}
}

func TestTransactionHandler_SQLInjection_InID(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	// Using URL-safe payloads
	sqlPayloads := []string{
		"1%20OR%201=1",
		"1%3BDROP%20TABLE",
		"1%27OR%271%27=%271",
		"-1%20OR%201=1",
	}

	for i, payload := range sqlPayloads {
		t.Run("sqli_id_"+string(rune('A'+i)), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/transactions/"+payload, nil)
			rr := httptest.NewRecorder()

			handler.HandleGetByID(rr, req)

			if rr.Code == http.StatusInternalServerError {
				t.Errorf("SQL injection in ID may have caused server error")
			}
		})
	}
}

// ==================== INVALID JSON FORMAT TESTS ====================

func TestTransactionHandler_InvalidJSONFormats(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()
	productRepo.Create(&model.Product{Name: "Test", Price: 100, Stock: 100})

	testCases := []struct {
		name string
		body string
	}{
		{"unclosed brace", `{"items":[{"product_id":1,"quantity":1}`},
		{"trailing comma", `{"items":[{"product_id":1,"quantity":1},]}`},
		{"missing colon", `{"items" [{"product_id":1,"quantity":1}]}`},
		{"empty object", `{}`},
		{"null body", "null"},
		{"empty string", ""},
		{"items as string", `{"items":"invalid"}`},
		{"items as number", `{"items":123}`},
		{"product_id as string", `{"items":[{"product_id":"1","quantity":1}]}`},
		{"quantity as string", `{"items":[{"product_id":1,"quantity":"1"}]}`},
		{"negative values in json", `{"items":[{"product_id":-1,"quantity":-1}]}`},
		{"float values", `{"items":[{"product_id":1.5,"quantity":1.5}]}`},
		{"null product_id", `{"items":[{"product_id":null,"quantity":1}]}`},
		{"null quantity", `{"items":[{"product_id":1,"quantity":null}]}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCheckout(rr, req)

			if rr.Code == http.StatusInternalServerError {
				t.Errorf("Invalid JSON '%s' caused server error", tc.name)
			}
		})
	}
}

// ==================== EDGE CASES ====================

func TestTransactionHandler_EdgeCases_EmptyItems(t *testing.T) {
	handler, _, _, _ := setupTransactionHandler()

	testCases := []struct {
		name string
		body string
	}{
		{"empty items array", `{"items":[]}`},
		{"null items", `{"items":null}`},
		{"missing items", `{}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCheckout(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("%s should return 400, got: %d", tc.name, rr.Code)
			}
		})
	}
}

func TestTransactionHandler_EdgeCases_MultipleItems(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()
	productRepo.Create(&model.Product{Name: "Product 1", Price: 100, Stock: 100})
	productRepo.Create(&model.Product{Name: "Product 2", Price: 200, Stock: 100})
	productRepo.Create(&model.Product{Name: "Product 3", Price: 300, Stock: 100})

	// Test with many items
	items := make([]model.CheckoutItem, 0)
	for i := 1; i <= 3; i++ {
		items = append(items, model.CheckoutItem{ProductID: i, Quantity: 1})
	}

	input := model.CheckoutRequest{Items: items}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCheckout(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Multiple items checkout should succeed, got: %d", rr.Code)
	}
}

func TestTransactionHandler_EdgeCases_DuplicateProducts(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()
	productRepo.Create(&model.Product{Name: "Test", Price: 100, Stock: 100})

	// Same product multiple times in items
	input := model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 1, Quantity: 5},
			{ProductID: 1, Quantity: 5},
		},
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCheckout(rr, req)

	// Should either succeed or fail gracefully
	if rr.Code == http.StatusInternalServerError {
		t.Error("Duplicate products in checkout caused server error")
	}
}

func TestTransactionHandler_EdgeCases_ExactStock(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()
	productRepo.Create(&model.Product{Name: "Test", Price: 100, Stock: 10})

	// Order exactly the available stock
	input := model.CheckoutRequest{
		Items: []model.CheckoutItem{{ProductID: 1, Quantity: 10}},
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCheckout(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Exact stock checkout should succeed, got: %d", rr.Code)
	}

	// Verify stock is now 0
	product, _ := productRepo.GetByID(1)
	if product != nil && product.Stock != 0 {
		t.Errorf("Stock should be 0 after exact stock checkout, got: %d", product.Stock)
	}
}

func TestTransactionHandler_EdgeCases_OneMoreThanStock(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()
	productRepo.Create(&model.Product{Name: "Test", Price: 100, Stock: 10})

	// Order one more than available
	input := model.CheckoutRequest{
		Items: []model.CheckoutItem{{ProductID: 1, Quantity: 11}},
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCheckout(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Checkout exceeding stock should return 400, got: %d", rr.Code)
	}
}

// ==================== CONCURRENT REQUEST TESTS ====================

func TestTransactionHandler_ConcurrentCheckouts(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()
	productRepo.Create(&model.Product{Name: "Test", Price: 100, Stock: 1000})

	concurrentRequests := 50
	done := make(chan bool, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			input := model.CheckoutRequest{
				Items: []model.CheckoutItem{{ProductID: 1, Quantity: 1}},
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCheckout(rr, req)
			done <- true
		}()
	}

	for i := 0; i < concurrentRequests; i++ {
		<-done
	}
}

// ==================== LARGE PAYLOAD TESTS ====================

func TestTransactionHandler_LargeNumberOfItems(t *testing.T) {
	handler, _, productRepo, _ := setupTransactionHandler()

	// Create many products
	for i := 1; i <= 100; i++ {
		productRepo.Create(&model.Product{Name: "Product", Price: 100, Stock: 1000})
	}

	// Checkout with many items
	items := make([]model.CheckoutItem, 0)
	for i := 1; i <= 100; i++ {
		items = append(items, model.CheckoutItem{ProductID: i, Quantity: 1})
	}

	input := model.CheckoutRequest{Items: items}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCheckout(rr, req)

	if rr.Code == http.StatusInternalServerError {
		t.Error("Large number of items caused server error")
	}
}
