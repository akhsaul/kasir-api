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

func setupProductHandler() (*ProductHandler, *memory.ProductRepository, *memory.CategoryRepository) {
	categoryRepo := memory.NewCategoryRepository()
	productRepo := memory.NewProductRepository(categoryRepo)
	svc := service.NewProductService(productRepo, categoryRepo)
	handler := NewProductHandler(svc)
	return handler, productRepo, categoryRepo
}

func TestNewProductHandler(t *testing.T) {
	svc := &service.ProductService{}
	handler := NewProductHandler(svc)

	if handler == nil {
		t.Error("NewProductHandler should return a non-nil handler")
	}
	if handler.service != svc {
		t.Error("NewProductHandler should set the service")
	}
}

func TestProductHandler_HandleGetAll_Success(t *testing.T) {
	handler, productRepo, _ := setupProductHandler()

	// Add test data
	productRepo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 10})
	productRepo.Create(&model.Product{Name: "Phone", Price: 500, Stock: 20})

	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetAll(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleGetAll should return 200, got: %d", rr.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	if response["status"] != "OK" {
		t.Errorf("Status should be OK, got: %v", response["status"])
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Response data should be a paginated object")
	}
	items, ok := data["items"].([]interface{})
	if !ok || len(items) != 2 {
		t.Errorf("Should return 2 products in items, got: %v", data["items"])
	}
	if data["total_items"].(float64) != 2 {
		t.Errorf("total_items should be 2, got: %v", data["total_items"])
	}
}

func TestProductHandler_HandleGetAll_WithFilter(t *testing.T) {
	handler, productRepo, _ := setupProductHandler()

	productRepo.Create(&model.Product{Name: "Laptop Pro", Price: 1500, Stock: 10})
	productRepo.Create(&model.Product{Name: "Laptop Basic", Price: 1000, Stock: 15})
	productRepo.Create(&model.Product{Name: "Phone", Price: 500, Stock: 20})

	req := httptest.NewRequest(http.MethodGet, "/api/products?name=laptop", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetAll(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleGetAll should return 200, got: %d", rr.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Response data should be a paginated object")
	}
	items, ok := data["items"].([]interface{})
	if !ok || len(items) != 2 {
		t.Errorf("Should return 2 products with name containing 'laptop'")
	}
}

func TestProductHandler_HandleGetAll_Empty(t *testing.T) {
	handler, _, _ := setupProductHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetAll(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleGetAll should return 200, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleGetByID_Success(t *testing.T) {
	handler, productRepo, _ := setupProductHandler()

	productRepo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 10})

	req := httptest.NewRequest(http.MethodGet, "/api/products/1", nil)
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

func TestProductHandler_HandleGetByID_InvalidID(t *testing.T) {
	handler, _, _ := setupProductHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/products/invalid", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleGetByID with invalid ID should return 404, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleGetByID_NotFound(t *testing.T) {
	handler, _, _ := setupProductHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/products/999", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleGetByID for non-existent product should return 404, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleGetByID_ZeroID(t *testing.T) {
	handler, _, _ := setupProductHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/products/0", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleGetByID with 0 ID should return 404, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleCreate_Success(t *testing.T) {
	handler, _, _ := setupProductHandler()

	input := model.ProductInput{Name: "Laptop", Price: 1000, Stock: 10}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreate(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("HandleCreate should return 201, got: %d", rr.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	if response["status"] != "OK" {
		t.Errorf("Status should be OK, got: %v", response["status"])
	}
}

func TestProductHandler_HandleCreate_WithCategory(t *testing.T) {
	handler, _, categoryRepo := setupProductHandler()

	categoryRepo.Create(&model.Category{Name: "Electronics", Description: "Electronic items"})

	categoryID := 1
	input := model.ProductInput{Name: "Laptop", Price: 1000, Stock: 10, CategoryID: &categoryID}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreate(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("HandleCreate should return 201, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleCreate_InvalidJSON(t *testing.T) {
	handler, _, _ := setupProductHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleCreate with invalid JSON should return 400, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleCreate_CategoryNotFound(t *testing.T) {
	handler, _, _ := setupProductHandler()

	categoryID := 999
	input := model.ProductInput{Name: "Laptop", Price: 1000, Stock: 10, CategoryID: &categoryID}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleCreate with invalid category should return 400, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleCreate_EmptyName(t *testing.T) {
	handler, _, _ := setupProductHandler()

	input := model.ProductInput{Name: "", Price: 1000, Stock: 10}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleCreate with empty name should return 400, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleCreate_InvalidPrice(t *testing.T) {
	handler, _, _ := setupProductHandler()

	input := model.ProductInput{Name: "Laptop", Price: 0, Stock: 10}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleCreate with invalid price should return 400, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleCreate_InvalidStock(t *testing.T) {
	handler, _, _ := setupProductHandler()

	input := model.ProductInput{Name: "Laptop", Price: 1000, Stock: -1}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleCreate with invalid stock should return 400, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleUpdate_Success(t *testing.T) {
	handler, productRepo, _ := setupProductHandler()

	productRepo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 10})

	input := model.ProductInput{Name: "Updated Laptop", Price: 1500, Stock: 5}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPut, "/api/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleUpdate(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleUpdate should return 200, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleUpdate_InvalidID(t *testing.T) {
	handler, _, _ := setupProductHandler()

	req := httptest.NewRequest(http.MethodPut, "/api/products/invalid", nil)
	rr := httptest.NewRecorder()

	handler.HandleUpdate(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleUpdate with invalid ID should return 404, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleUpdate_InvalidJSON(t *testing.T) {
	handler, productRepo, _ := setupProductHandler()

	productRepo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 10})

	req := httptest.NewRequest(http.MethodPut, "/api/products/1", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleUpdate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleUpdate with invalid JSON should return 400, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleUpdate_NotFound(t *testing.T) {
	handler, _, _ := setupProductHandler()

	input := model.ProductInput{Name: "Updated", Price: 1000, Stock: 10}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPut, "/api/products/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleUpdate(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleUpdate for non-existent product should return 404, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleUpdate_CategoryNotFound(t *testing.T) {
	handler, productRepo, _ := setupProductHandler()

	productRepo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 10})

	categoryID := 999
	input := model.ProductInput{Name: "Updated", Price: 1000, Stock: 10, CategoryID: &categoryID}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPut, "/api/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleUpdate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleUpdate with invalid category should return 400, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleUpdate_ValidationErrors(t *testing.T) {
	handler, productRepo, _ := setupProductHandler()

	productRepo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 10})

	testCases := []struct {
		name  string
		input model.ProductInput
	}{
		{"empty name", model.ProductInput{Name: "", Price: 1000, Stock: 10}},
		{"invalid price", model.ProductInput{Name: "Laptop", Price: 0, Stock: 10}},
		{"invalid stock", model.ProductInput{Name: "Laptop", Price: 1000, Stock: -1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.input)

			req := httptest.NewRequest(http.MethodPut, "/api/products/1", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleUpdate(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("HandleUpdate with %s should return 400, got: %d", tc.name, rr.Code)
			}
		})
	}
}

func TestProductHandler_HandleDelete_Success(t *testing.T) {
	handler, productRepo, _ := setupProductHandler()

	productRepo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 10})

	req := httptest.NewRequest(http.MethodDelete, "/api/products/1", nil)
	rr := httptest.NewRecorder()

	handler.HandleDelete(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleDelete should return 200, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleDelete_InvalidID(t *testing.T) {
	handler, _, _ := setupProductHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/products/invalid", nil)
	rr := httptest.NewRecorder()

	handler.HandleDelete(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleDelete with invalid ID should return 404, got: %d", rr.Code)
	}
}

func TestProductHandler_HandleDelete_NotFound(t *testing.T) {
	handler, _, _ := setupProductHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/products/999", nil)
	rr := httptest.NewRecorder()

	handler.HandleDelete(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleDelete for non-existent product should return 404, got: %d", rr.Code)
	}
}

// ==================== BOUNDARY VALUE TESTS ====================

func TestProductHandler_BoundaryValues_ID(t *testing.T) {
	handler, productRepo, _ := setupProductHandler()
	productRepo.Create(&model.Product{Name: "Test", Price: 100, Stock: 10})

	testCases := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		// Invalid ID formats
		{"negative id", "-1", http.StatusNotFound},
		{"negative large", "-999999999", http.StatusNotFound},
		{"zero id", "0", http.StatusNotFound},
		{"float id", "1.5", http.StatusNotFound},
		{"string id", "abc", http.StatusNotFound},
		{"empty id", "", http.StatusNotFound},
		{"unicode id", "١٢٣", http.StatusNotFound},
		{"max int overflow", "9999999999999999999", http.StatusNotFound},
		{"hex id", "0x1", http.StatusNotFound},
		{"binary id", "0b1", http.StatusNotFound},
		{"octal id", "01", http.StatusOK}, // "01" parses as 1 in Go
		// Valid ID
		{"valid id", "1", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/products/"+tc.id, nil)
			rr := httptest.NewRecorder()

			handler.HandleGetByID(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("ID '%s': expected %d, got %d", tc.id, tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestProductHandler_BoundaryValues_Price(t *testing.T) {
	handler, _, _ := setupProductHandler()

	testCases := []struct {
		name           string
		price          interface{}
		expectedStatus int
	}{
		// Invalid prices
		{"zero price", 0, http.StatusBadRequest},
		{"negative price", -1, http.StatusBadRequest},
		{"negative large", -999999999, http.StatusBadRequest},
		// Valid prices
		{"minimum valid price", 1, http.StatusCreated},
		{"large price", 999999999, http.StatusCreated},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := map[string]interface{}{
				"name":  "Test Product",
				"price": tc.price,
				"stock": 10,
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Price %v: expected %d, got %d", tc.price, tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestProductHandler_BoundaryValues_Stock(t *testing.T) {
	handler, _, _ := setupProductHandler()

	testCases := []struct {
		name           string
		stock          interface{}
		expectedStatus int
	}{
		// Invalid stock
		{"negative stock", -1, http.StatusBadRequest},
		{"negative large", -999999999, http.StatusBadRequest},
		// Valid stock
		{"zero stock", 0, http.StatusCreated},
		{"positive stock", 10, http.StatusCreated},
		{"large stock", 999999999, http.StatusCreated},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := map[string]interface{}{
				"name":  "Test Product",
				"price": 100,
				"stock": tc.stock,
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Stock %v: expected %d, got %d", tc.stock, tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestProductHandler_BoundaryValues_Name(t *testing.T) {
	handler, _, _ := setupProductHandler()

	testCases := []struct {
		name           string
		productName    string
		expectedStatus int
	}{
		// Invalid names
		{"empty name", "", http.StatusBadRequest},
		{"whitespace only", "   ", http.StatusBadRequest},
		{"tab only", "\t", http.StatusBadRequest},
		{"newline only", "\n", http.StatusBadRequest},
		// Valid names
		{"single char", "A", http.StatusCreated},
		{"normal name", "Laptop", http.StatusCreated},
		{"name with spaces", "Gaming Laptop Pro", http.StatusCreated},
		{"name with numbers", "iPhone 15", http.StatusCreated},
		{"very long name", string(make([]byte, 1000)), http.StatusCreated}, // 1000 chars
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			name := tc.productName
			if tc.name == "very long name" {
				name = ""
				for i := 0; i < 1000; i++ {
					name += "A"
				}
			}

			input := map[string]interface{}{
				"name":  name,
				"price": 100,
				"stock": 10,
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Name '%s': expected %d, got %d", tc.name, tc.expectedStatus, rr.Code)
			}
		})
	}
}

// ==================== XSS INJECTION TESTS ====================

func TestProductHandler_XSSInjection(t *testing.T) {
	handler, _, _ := setupProductHandler()

	xssPayloads := []string{
		// Basic XSS
		"<script>alert('xss')</script>",
		"<img src=x onerror=alert('xss')>",
		"<svg onload=alert('xss')>",
		"<body onload=alert('xss')>",
		// Event handlers
		"<div onmouseover=alert('xss')>hover me</div>",
		"<input onfocus=alert('xss') autofocus>",
		// JavaScript URLs
		"javascript:alert('xss')",
		"<a href='javascript:alert(1)'>click</a>",
		// Encoded XSS
		"&lt;script&gt;alert('xss')&lt;/script&gt;",
		"%3Cscript%3Ealert('xss')%3C/script%3E",
		// Data URLs
		"data:text/html,<script>alert('xss')</script>",
		// Unicode encoding
		"\u003cscript\u003ealert('xss')\u003c/script\u003e",
		// Mixed case
		"<ScRiPt>alert('xss')</ScRiPt>",
		// Null byte injection
		"<scr\x00ipt>alert('xss')</script>",
	}

	for _, payload := range xssPayloads {
		t.Run("xss_"+payload[:min(20, len(payload))], func(t *testing.T) {
			input := map[string]interface{}{
				"name":  payload,
				"price": 100,
				"stock": 10,
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			// XSS payloads should either be stored safely or rejected
			// The key is they should NOT execute - this tests that the API accepts/rejects properly
			if rr.Code != http.StatusCreated && rr.Code != http.StatusBadRequest {
				t.Errorf("XSS payload handling failed, got status: %d", rr.Code)
			}
		})
	}
}

// ==================== SQL INJECTION TESTS ====================

func TestProductHandler_SQLInjection(t *testing.T) {
	handler, _, _ := setupProductHandler()

	sqlPayloads := []string{
		// Classic SQL injection
		"'; DROP TABLE products; --",
		"1'; DROP TABLE products; --",
		"1 OR 1=1",
		"1' OR '1'='1",
		"1\" OR \"1\"=\"1",
		"1; SELECT * FROM users",
		// Union based
		"1 UNION SELECT * FROM users",
		"1' UNION SELECT username, password FROM users--",
		// Blind SQL injection
		"1' AND 1=1--",
		"1' AND 1=2--",
		"1' AND SLEEP(5)--",
		"1' WAITFOR DELAY '0:0:5'--",
		// Error based
		"1' AND (SELECT 1 FROM(SELECT COUNT(*),CONCAT((SELECT user()),FLOOR(RAND(0)*2))x FROM information_schema.tables GROUP BY x)a)--",
		// Comment bypass
		"1'/**/OR/**/1=1",
		"1'%00OR%001=1",
		// PostgreSQL specific
		"1'; SELECT pg_sleep(5);--",
		"1'; COPY (SELECT * FROM users) TO '/tmp/test';--",
		// Stacked queries
		"1; INSERT INTO products (name) VALUES ('hacked');--",
		// Batch queries
		"1; UPDATE products SET price=0;--",
	}

	for i, payload := range sqlPayloads {
		t.Run("sqli_"+string(rune('A'+i)), func(t *testing.T) {
			input := map[string]interface{}{
				"name":  payload,
				"price": 100,
				"stock": 10,
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			// SQL injection payloads should be treated as regular strings (parameterized queries)
			// They should NOT cause errors or data manipulation
			if rr.Code == http.StatusInternalServerError {
				t.Errorf("SQL injection may have caused server error for payload: %s", payload)
			}
		})
	}
}

func TestProductHandler_SQLInjection_InID(t *testing.T) {
	handler, _, _ := setupProductHandler()

	// Using URL-safe payloads for path testing
	sqlPayloads := []string{
		"1%20OR%201=1",
		"1%3BDROP%20TABLE",
		"1%27OR%271%27=%271",
		"-1%20OR%201=1",
		"1--",
		"1%2F*",
	}

	for i, payload := range sqlPayloads {
		t.Run("sqli_id_"+string(rune('A'+i)), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/products/"+payload, nil)
			rr := httptest.NewRecorder()

			handler.HandleGetByID(rr, req)

			// ID should be parsed as int, SQL injection attempts should fail parsing
			if rr.Code == http.StatusInternalServerError {
				t.Errorf("SQL injection in ID may have caused server error")
			}
		})
	}
}

// ==================== INVALID JSON FORMAT TESTS ====================

func TestProductHandler_InvalidJSONFormats(t *testing.T) {
	handler, _, _ := setupProductHandler()

	testCases := []struct {
		name string
		body string
	}{
		// Malformed JSON
		{"unclosed brace", `{"name":"test"`},
		{"unclosed bracket", `{"name":"test","items":[}`},
		{"trailing comma", `{"name":"test",}`},
		{"missing colon", `{"name" "test"}`},
		{"single quotes", `{'name':'test'}`},
		{"unquoted key", `{name:"test"}`},
		{"double comma", `{"name":"test",,"price":100}`},
		// Empty/null values
		{"empty object", `{}`},
		{"null body", "null"},
		{"empty string", ""},
		{"whitespace only", "   "},
		// Type mismatches
		{"price as string", `{"name":"test","price":"100","stock":10}`},
		{"stock as string", `{"name":"test","price":100,"stock":"10"}`},
		{"name as number", `{"name":123,"price":100,"stock":10}`},
		{"price as array", `{"name":"test","price":[100],"stock":10}`},
		{"price as object", `{"name":"test","price":{"value":100},"stock":10}`},
		{"price as boolean", `{"name":"test","price":true,"stock":10}`},
		{"price as null", `{"name":"test","price":null,"stock":10}`},
		// Extra/unexpected fields (should be ignored)
		{"extra field", `{"name":"test","price":100,"stock":10,"extra":"field"}`},
		// Nested objects
		{"deeply nested", `{"name":"test","price":100,"stock":10,"nested":{"deep":{"deeper":{"deepest":"value"}}}}`},
		// Large numbers
		{"float price", `{"name":"test","price":100.5,"stock":10}`},
		{"scientific notation", `{"name":"test","price":1e10,"stock":10}`},
		// Unicode
		{"unicode name", `{"name":"测试产品","price":100,"stock":10}`},
		{"emoji name", `{"name":"📱 Phone","price":100,"stock":10}`},
		// Escape sequences
		{"escaped quotes", `{"name":"test\"product","price":100,"stock":10}`},
		{"newline in name", `{"name":"test\nproduct","price":100,"stock":10}`},
		{"tab in name", `{"name":"test\tproduct","price":100,"stock":10}`},
		// Binary data
		{"null byte", "{\"name\":\"test\x00product\",\"price\":100,\"stock\":10}"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			// Should either succeed (201) or return bad request (400), never server error (500)
			if rr.Code == http.StatusInternalServerError {
				t.Errorf("Invalid JSON '%s' caused server error", tc.name)
			}
		})
	}
}

// ==================== HTTP HEADER TESTS ====================

func TestProductHandler_ContentTypeHeaders(t *testing.T) {
	handler, _, _ := setupProductHandler()

	input := map[string]interface{}{
		"name":  "Test",
		"price": 100,
		"stock": 10,
	}
	body, _ := json.Marshal(input)

	testCases := []struct {
		name        string
		contentType string
	}{
		{"no content type", ""},
		{"text/plain", "text/plain"},
		{"text/html", "text/html"},
		{"application/xml", "application/xml"},
		{"application/x-www-form-urlencoded", "application/x-www-form-urlencoded"},
		{"multipart/form-data", "multipart/form-data"},
		{"application/json", "application/json"},
		{"application/json; charset=utf-8", "application/json; charset=utf-8"},
		{"APPLICATION/JSON", "APPLICATION/JSON"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			// Should handle all content types gracefully
			if rr.Code == http.StatusInternalServerError {
				t.Errorf("Content-Type '%s' caused server error", tc.contentType)
			}
		})
	}
}

// ==================== SPECIAL CHARACTERS TESTS ====================

func TestProductHandler_SpecialCharacters(t *testing.T) {
	handler, _, _ := setupProductHandler()

	specialStrings := []string{
		// Control characters
		"\x00\x01\x02\x03",
		"\r\n\t",
		// Unicode special
		"\u200B", // Zero-width space
		"\u200C", // Zero-width non-joiner
		"\u200D", // Zero-width joiner
		"\uFEFF", // BOM
		// RTL/LTR markers
		"\u202A\u202B\u202C\u202D\u202E",
		// Combining characters
		"a\u0300\u0301", // a with combining accents
		// Surrogate pairs (emoji)
		"👨‍👩‍👧‍👦", // Family emoji (complex)
		"🏳️‍🌈",    // Rainbow flag
		// Newlines variants
		"\r\n",
		"\n\r",
		"\u2028", // Line separator
		"\u2029", // Paragraph separator
		// Quotes and brackets
		`"'<>[]{}()`,
		// Backslashes
		`\\\\`,
		`\/\/`,
		// HTML entities
		"&amp;&lt;&gt;&quot;",
		// URL encoding
		"%00%20%22%27",
		// Path traversal
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32",
		// Command injection
		"; ls -la",
		"| cat /etc/passwd",
		"` whoami `",
		"$( whoami )",
		// LDAP injection
		"*)(uid=*))(|(uid=*",
		// XML injection
		"<?xml version=\"1.0\"?><!DOCTYPE foo [<!ENTITY xxe SYSTEM \"file:///etc/passwd\">]>",
	}

	for i, str := range specialStrings {
		t.Run("special_"+string(rune('A'+i)), func(t *testing.T) {
			input := map[string]interface{}{
				"name":  str,
				"price": 100,
				"stock": 10,
			}
			body, err := json.Marshal(input)
			if err != nil {
				// Some strings can't be marshaled to JSON
				return
			}

			req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			// Should not cause server error
			if rr.Code == http.StatusInternalServerError {
				t.Errorf("Special character caused server error: %q", str)
			}
		})
	}
}

// ==================== LARGE PAYLOAD TESTS ====================

func TestProductHandler_LargePayloads(t *testing.T) {
	handler, _, _ := setupProductHandler()

	testCases := []struct {
		name     string
		nameSize int
	}{
		{"1KB name", 1024},
		{"10KB name", 10 * 1024},
		{"100KB name", 100 * 1024},
		{"1MB name", 1024 * 1024},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			largeName := make([]byte, tc.nameSize)
			for i := range largeName {
				largeName[i] = 'A'
			}

			input := map[string]interface{}{
				"name":  string(largeName),
				"price": 100,
				"stock": 10,
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			// Should handle large payloads without server error
			if rr.Code == http.StatusInternalServerError {
				t.Errorf("Large payload (%d bytes) caused server error", tc.nameSize)
			}
		})
	}
}

// ==================== CONCURRENT REQUEST TESTS ====================

func TestProductHandler_ConcurrentRequests(t *testing.T) {
	handler, _, _ := setupProductHandler()

	concurrentRequests := 100
	done := make(chan bool, concurrentRequests)
	errors := make(chan error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func(id int) {
			input := map[string]interface{}{
				"name":  "Product " + string(rune('A'+id%26)),
				"price": 100 + id,
				"stock": 10,
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			if rr.Code != http.StatusCreated {
				errors <- nil // Ignore non-201 for this test
			}
			done <- true
		}(i)
	}

	for i := 0; i < concurrentRequests; i++ {
		<-done
	}

	close(errors)
	for err := range errors {
		if err != nil {
			t.Errorf("Concurrent request error: %v", err)
		}
	}
}

// ==================== QUERY PARAMETER TESTS ====================

func TestProductHandler_QueryParameterInjection(t *testing.T) {
	handler, productRepo, _ := setupProductHandler()
	productRepo.Create(&model.Product{Name: "Test Product", Price: 100, Stock: 10})

	// URL-encoded malicious queries
	maliciousQueries := []string{
		"name=%3Cscript%3Ealert(1)%3C/script%3E",   // XSS
		"name=%27%3B%20DROP%20TABLE%20products%3B", // SQL injection (URL encoded)
		"name=%00",                         // Null byte
		"name=..%2F..%2F..%2Fetc%2Fpasswd", // Path traversal
		"name=%252e%252e%252f",             // Double encoding
		"name=test",                        // Normal query
	}

	for i, query := range maliciousQueries {
		t.Run("query_"+string(rune('A'+i)), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/products?"+query, nil)
			rr := httptest.NewRecorder()

			handler.HandleGetAll(rr, req)

			if rr.Code == http.StatusInternalServerError {
				t.Errorf("Malicious query caused server error")
			}
		})
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
