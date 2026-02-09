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

func setupCategoryHandler() (*CategoryHandler, *memory.CategoryRepository) {
	repo := memory.NewCategoryRepository()
	svc := service.NewCategoryService(repo)
	handler := NewCategoryHandler(svc)
	return handler, repo
}

func TestNewCategoryHandler(t *testing.T) {
	svc := &service.CategoryService{}
	handler := NewCategoryHandler(svc)

	if handler == nil {
		t.Error("NewCategoryHandler should return a non-nil handler")
	}
	if handler.service != svc {
		t.Error("NewCategoryHandler should set the service")
	}
}

func TestCategoryHandler_HandleGetAll_Success(t *testing.T) {
	handler, repo := setupCategoryHandler()

	// Add test data
	repo.Create(&model.Category{Name: "Electronics", Description: "Electronic items"})
	repo.Create(&model.Category{Name: "Food", Description: "Food items"})

	req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
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
		t.Errorf("Should return 2 categories in items, got: %v", data["items"])
	}
	if data["total_items"].(float64) != 2 {
		t.Errorf("total_items should be 2, got: %v", data["total_items"])
	}
}

func TestCategoryHandler_HandleGetAll_Empty(t *testing.T) {
	handler, _ := setupCategoryHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetAll(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleGetAll should return 200, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleGetByID_Success(t *testing.T) {
	handler, repo := setupCategoryHandler()

	repo.Create(&model.Category{Name: "Electronics", Description: "Electronic items"})

	req := httptest.NewRequest(http.MethodGet, "/api/categories/1", nil)
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

func TestCategoryHandler_HandleGetByID_InvalidID(t *testing.T) {
	handler, _ := setupCategoryHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/categories/invalid", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleGetByID with invalid ID should return 404, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleGetByID_NotFound(t *testing.T) {
	handler, _ := setupCategoryHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/categories/999", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleGetByID for non-existent category should return 404, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleGetByID_IDRequired(t *testing.T) {
	handler, _ := setupCategoryHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/categories/0", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetByID(rr, req)

	// ID 0 or negative should return error
	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleGetByID with 0 ID should return 400, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleCreate_Success(t *testing.T) {
	handler, _ := setupCategoryHandler()

	category := model.Category{Name: "Electronics", Description: "Electronic items"}
	body, _ := json.Marshal(category)

	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
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

func TestCategoryHandler_HandleCreate_InvalidJSON(t *testing.T) {
	handler, _ := setupCategoryHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleCreate with invalid JSON should return 400, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleCreate_EmptyName(t *testing.T) {
	handler, _ := setupCategoryHandler()

	category := model.Category{Name: "", Description: "desc"}
	body, _ := json.Marshal(category)

	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleCreate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleCreate with empty name should return 400, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleUpdate_Success(t *testing.T) {
	handler, repo := setupCategoryHandler()

	repo.Create(&model.Category{Name: "Electronics", Description: "Electronic items"})

	updated := model.Category{Name: "Updated Electronics", Description: "Updated desc"}
	body, _ := json.Marshal(updated)

	req := httptest.NewRequest(http.MethodPut, "/api/categories/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleUpdate(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleUpdate should return 200, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleUpdate_InvalidID(t *testing.T) {
	handler, _ := setupCategoryHandler()

	req := httptest.NewRequest(http.MethodPut, "/api/categories/invalid", nil)
	rr := httptest.NewRecorder()

	handler.HandleUpdate(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleUpdate with invalid ID should return 404, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleUpdate_InvalidJSON(t *testing.T) {
	handler, repo := setupCategoryHandler()

	repo.Create(&model.Category{Name: "Electronics", Description: "Electronic items"})

	req := httptest.NewRequest(http.MethodPut, "/api/categories/1", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleUpdate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleUpdate with invalid JSON should return 400, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleUpdate_NotFound(t *testing.T) {
	handler, _ := setupCategoryHandler()

	updated := model.Category{Name: "Updated", Description: "desc"}
	body, _ := json.Marshal(updated)

	req := httptest.NewRequest(http.MethodPut, "/api/categories/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleUpdate(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleUpdate for non-existent category should return 404, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleUpdate_EmptyName(t *testing.T) {
	handler, repo := setupCategoryHandler()

	repo.Create(&model.Category{Name: "Electronics", Description: "Electronic items"})

	updated := model.Category{Name: "", Description: "desc"}
	body, _ := json.Marshal(updated)

	req := httptest.NewRequest(http.MethodPut, "/api/categories/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleUpdate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleUpdate with empty name should return 400, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleDelete_Success(t *testing.T) {
	handler, repo := setupCategoryHandler()

	repo.Create(&model.Category{Name: "Electronics", Description: "Electronic items"})

	req := httptest.NewRequest(http.MethodDelete, "/api/categories/1", nil)
	rr := httptest.NewRecorder()

	handler.HandleDelete(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleDelete should return 200, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleDelete_InvalidID(t *testing.T) {
	handler, _ := setupCategoryHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/categories/invalid", nil)
	rr := httptest.NewRecorder()

	handler.HandleDelete(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleDelete with invalid ID should return 404, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleDelete_NotFound(t *testing.T) {
	handler, _ := setupCategoryHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/categories/999", nil)
	rr := httptest.NewRecorder()

	handler.HandleDelete(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleDelete for non-existent category should return 404, got: %d", rr.Code)
	}
}

func TestCategoryHandler_HandleDelete_IDRequired(t *testing.T) {
	handler, _ := setupCategoryHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/categories/0", nil)
	rr := httptest.NewRecorder()

	handler.HandleDelete(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleDelete with 0 ID should return 400, got: %d", rr.Code)
	}
}

// ==================== BOUNDARY VALUE TESTS ====================

func TestCategoryHandler_BoundaryValues_ID(t *testing.T) {
	handler, repo := setupCategoryHandler()
	repo.Create(&model.Category{Name: "Test", Description: "Test"})

	testCases := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		// Parseable but invalid IDs (service layer returns error)
		{"negative id", "-1", http.StatusBadRequest},
		{"negative large", "-999999999", http.StatusBadRequest},
		{"zero id", "0", http.StatusBadRequest},
		// Non-parseable IDs (ParseIDFromPath returns 404)
		{"float id", "1.5", http.StatusNotFound},
		{"string id", "abc", http.StatusNotFound},
		{"empty id", "", http.StatusNotFound},
		{"max int overflow", "9999999999999999999", http.StatusNotFound},
		{"hex id", "0x1", http.StatusNotFound},
		// Valid ID
		{"valid id", "1", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/categories/"+tc.id, nil)
			rr := httptest.NewRecorder()

			handler.HandleGetByID(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("ID '%s': expected %d, got %d", tc.id, tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestCategoryHandler_BoundaryValues_Name(t *testing.T) {
	handler, _ := setupCategoryHandler()

	testCases := []struct {
		name           string
		categoryName   string
		expectedStatus int
	}{
		// Invalid names
		{"empty name", "", http.StatusBadRequest},
		{"whitespace only", "   ", http.StatusBadRequest},
		{"tab only", "\t", http.StatusBadRequest},
		{"newline only", "\n", http.StatusBadRequest},
		// Valid names
		{"single char", "A", http.StatusCreated},
		{"normal name", "Electronics", http.StatusCreated},
		{"name with spaces", "Home & Garden", http.StatusCreated},
		{"name with numbers", "Category 123", http.StatusCreated},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := map[string]interface{}{
				"name":        tc.categoryName,
				"description": "Test description",
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
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

func TestCategoryHandler_XSSInjection(t *testing.T) {
	handler, _ := setupCategoryHandler()

	xssPayloads := []string{
		"<script>alert('xss')</script>",
		"<img src=x onerror=alert('xss')>",
		"<svg onload=alert('xss')>",
		"javascript:alert('xss')",
		"<a href='javascript:alert(1)'>click</a>",
		"&lt;script&gt;alert('xss')&lt;/script&gt;",
		"\u003cscript\u003ealert('xss')\u003c/script\u003e",
		"<ScRiPt>alert('xss')</ScRiPt>",
	}

	for _, payload := range xssPayloads {
		t.Run("xss_name", func(t *testing.T) {
			input := map[string]interface{}{
				"name":        payload,
				"description": "Test",
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			if rr.Code != http.StatusCreated && rr.Code != http.StatusBadRequest {
				t.Errorf("XSS payload handling failed, got status: %d", rr.Code)
			}
		})

		t.Run("xss_description", func(t *testing.T) {
			input := map[string]interface{}{
				"name":        "Test Category",
				"description": payload,
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			if rr.Code != http.StatusCreated && rr.Code != http.StatusBadRequest {
				t.Errorf("XSS payload in description failed, got status: %d", rr.Code)
			}
		})
	}
}

// ==================== SQL INJECTION TESTS ====================

func TestCategoryHandler_SQLInjection(t *testing.T) {
	handler, _ := setupCategoryHandler()

	sqlPayloads := []string{
		"'; DROP TABLE categories; --",
		"1 OR 1=1",
		"1' OR '1'='1",
		"1 UNION SELECT * FROM users",
		"1'; SELECT pg_sleep(5);--",
		"1; DELETE FROM categories;--",
	}

	for i, payload := range sqlPayloads {
		t.Run("sqli_"+string(rune('A'+i)), func(t *testing.T) {
			input := map[string]interface{}{
				"name":        payload,
				"description": "Test",
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			if rr.Code == http.StatusInternalServerError {
				t.Errorf("SQL injection may have caused server error for payload: %s", payload)
			}
		})
	}
}

func TestCategoryHandler_SQLInjection_InID(t *testing.T) {
	handler, _ := setupCategoryHandler()

	// Using URL-safe payloads
	sqlPayloads := []string{
		"1%20OR%201=1",
		"1%3BDROP%20TABLE",
		"1%27OR%271%27=%271",
		"-1%20OR%201=1",
	}

	for i, payload := range sqlPayloads {
		t.Run("sqli_id_"+string(rune('A'+i)), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/categories/"+payload, nil)
			rr := httptest.NewRecorder()

			handler.HandleGetByID(rr, req)

			if rr.Code == http.StatusInternalServerError {
				t.Errorf("SQL injection in ID may have caused server error")
			}
		})
	}
}

// ==================== INVALID JSON FORMAT TESTS ====================

func TestCategoryHandler_InvalidJSONFormats(t *testing.T) {
	handler, _ := setupCategoryHandler()

	testCases := []struct {
		name string
		body string
	}{
		{"unclosed brace", `{"name":"test"`},
		{"trailing comma", `{"name":"test",}`},
		{"missing colon", `{"name" "test"}`},
		{"single quotes", `{'name':'test'}`},
		{"empty object", `{}`},
		{"null body", "null"},
		{"empty string", ""},
		{"name as number", `{"name":123,"description":"test"}`},
		{"name as array", `{"name":["test"],"description":"test"}`},
		{"name as null", `{"name":null,"description":"test"}`},
		{"unicode name", `{"name":"测试分类","description":"test"}`},
		{"emoji name", `{"name":"📦 Category","description":"test"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			if rr.Code == http.StatusInternalServerError {
				t.Errorf("Invalid JSON '%s' caused server error", tc.name)
			}
		})
	}
}

// ==================== SPECIAL CHARACTERS TESTS ====================

func TestCategoryHandler_SpecialCharacters(t *testing.T) {
	handler, _ := setupCategoryHandler()

	specialStrings := []string{
		"\x00\x01\x02\x03",
		"\r\n\t",
		"\u200B",
		"👨‍👩‍👧‍👦",
		`"'<>[]{}()`,
		`\\\\`,
		"&amp;&lt;&gt;",
		"../../../etc/passwd",
		"; ls -la",
		"| cat /etc/passwd",
	}

	for i, str := range specialStrings {
		t.Run("special_"+string(rune('A'+i)), func(t *testing.T) {
			input := map[string]interface{}{
				"name":        str,
				"description": "Test",
			}
			body, err := json.Marshal(input)
			if err != nil {
				return
			}

			req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			if rr.Code == http.StatusInternalServerError {
				t.Errorf("Special character caused server error: %q", str)
			}
		})
	}
}

// ==================== LARGE PAYLOAD TESTS ====================

func TestCategoryHandler_LargePayloads(t *testing.T) {
	handler, _ := setupCategoryHandler()

	testCases := []struct {
		name     string
		nameSize int
	}{
		{"1KB name", 1024},
		{"10KB name", 10 * 1024},
		{"100KB name", 100 * 1024},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			largeName := make([]byte, tc.nameSize)
			for i := range largeName {
				largeName[i] = 'A'
			}

			input := map[string]interface{}{
				"name":        string(largeName),
				"description": "Test",
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)

			if rr.Code == http.StatusInternalServerError {
				t.Errorf("Large payload (%d bytes) caused server error", tc.nameSize)
			}
		})
	}
}

// ==================== CONCURRENT REQUEST TESTS ====================

func TestCategoryHandler_ConcurrentRequests(t *testing.T) {
	handler, _ := setupCategoryHandler()

	concurrentRequests := 50
	done := make(chan bool, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func(id int) {
			input := map[string]interface{}{
				"name":        "Category " + string(rune('A'+id%26)),
				"description": "Description",
			}
			body, _ := json.Marshal(input)

			req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreate(rr, req)
			done <- true
		}(i)
	}

	for i := 0; i < concurrentRequests; i++ {
		<-done
	}
}
