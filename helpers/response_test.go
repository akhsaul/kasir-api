package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"message": "hello"}

	WriteJSON(rr, http.StatusOK, data)

	if rr.Code != http.StatusOK {
		t.Errorf("WriteJSON should set status code to 200, got: %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("WriteJSON should set Content-Type to application/json, got: %s", contentType)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)
	if response["message"] != "hello" {
		t.Errorf("WriteJSON should encode data correctly, got: %v", response)
	}
}

func TestWriteJSON_DifferentStatusCodes(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
	}{
		{"OK", http.StatusOK},
		{"Created", http.StatusCreated},
		{"BadRequest", http.StatusBadRequest},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			WriteJSON(rr, tc.statusCode, nil)

			if rr.Code != tc.statusCode {
				t.Errorf("WriteJSON should set status code to %d, got: %d", tc.statusCode, rr.Code)
			}
		})
	}
}

func TestWriteSuccess(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]int{"id": 1}

	WriteSuccess(rr, http.StatusOK, "Success", data)

	if rr.Code != http.StatusOK {
		t.Errorf("WriteSuccess should set status code to 200, got: %d", rr.Code)
	}

	var response Response
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Status != "OK" {
		t.Errorf("WriteSuccess should set status to OK, got: %s", response.Status)
	}
	if response.Message != "Success" {
		t.Errorf("WriteSuccess should set message correctly, got: %s", response.Message)
	}
}

func TestWriteSuccess_WithNilData(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteSuccess(rr, http.StatusOK, "Success", nil)

	var response Response
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Status != "OK" {
		t.Errorf("WriteSuccess should set status to OK, got: %s", response.Status)
	}
	if response.Data != nil {
		t.Errorf("WriteSuccess should omit nil data, got: %v", response.Data)
	}
}

func TestWriteSuccess_Created(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]int{"id": 1}

	WriteSuccess(rr, http.StatusCreated, "Created successfully", data)

	if rr.Code != http.StatusCreated {
		t.Errorf("WriteSuccess should set status code to 201, got: %d", rr.Code)
	}
}

func TestWriteError(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/products", nil)

	WriteError(rr, req, http.StatusBadRequest, "Invalid input", errors.New("validation error"))

	if rr.Code != http.StatusBadRequest {
		t.Errorf("WriteError should set status code to 400, got: %d", rr.Code)
	}

	var response Response
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Status != "ERROR" {
		t.Errorf("WriteError should set status to ERROR, got: %s", response.Status)
	}
	if response.Message != "Invalid input" {
		t.Errorf("WriteError should set message correctly, got: %s", response.Message)
	}
}

func TestWriteError_NilError(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)

	WriteError(rr, req, http.StatusNotFound, "Not found", nil)

	if rr.Code != http.StatusNotFound {
		t.Errorf("WriteError should set status code to 404, got: %d", rr.Code)
	}

	var response Response
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Status != "ERROR" {
		t.Errorf("WriteError should set status to ERROR, got: %s", response.Status)
	}
}

func TestWriteError_DifferentStatusCodes(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		message    string
	}{
		{"BadRequest", http.StatusBadRequest, "Bad request"},
		{"NotFound", http.StatusNotFound, "Not found"},
		{"InternalServerError", http.StatusInternalServerError, "Internal error"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			WriteError(rr, req, tc.statusCode, tc.message, nil)

			if rr.Code != tc.statusCode {
				t.Errorf("WriteError should set status code to %d, got: %d", tc.statusCode, rr.Code)
			}

			var response Response
			json.NewDecoder(rr.Body).Decode(&response)

			if response.Message != tc.message {
				t.Errorf("WriteError should set message to %s, got: %s", tc.message, response.Message)
			}
		})
	}
}

func TestValidatePayload_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"name":"Test","price":100}`)
	req := httptest.NewRequest(http.MethodPost, "/api/products", body)

	var input struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}

	result := ValidatePayload(rr, req, &input)

	if !result {
		t.Error("ValidatePayload should return true for valid JSON")
	}
	if input.Name != "Test" {
		t.Errorf("ValidatePayload should decode name correctly, got: %s", input.Name)
	}
	if input.Price != 100 {
		t.Errorf("ValidatePayload should decode price correctly, got: %d", input.Price)
	}
}

func TestValidatePayload_InvalidJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	body := bytes.NewBufferString(`invalid json`)
	req := httptest.NewRequest(http.MethodPost, "/api/products", body)

	var input struct {
		Name string `json:"name"`
	}

	result := ValidatePayload(rr, req, &input)

	if result {
		t.Error("ValidatePayload should return false for invalid JSON")
	}
	if rr.Code != http.StatusBadRequest {
		t.Errorf("ValidatePayload should set status code to 400, got: %d", rr.Code)
	}
}

func TestValidatePayload_EmptyBody(t *testing.T) {
	rr := httptest.NewRecorder()
	body := bytes.NewBufferString(``)
	req := httptest.NewRequest(http.MethodPost, "/api/products", body)

	var input struct {
		Name string `json:"name"`
	}

	result := ValidatePayload(rr, req, &input)

	if result {
		t.Error("ValidatePayload should return false for empty body")
	}
}

func TestValidatePayload_PartialJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"name":"Test"`)
	req := httptest.NewRequest(http.MethodPost, "/api/products", body)

	var input struct {
		Name string `json:"name"`
	}

	result := ValidatePayload(rr, req, &input)

	if result {
		t.Error("ValidatePayload should return false for partial JSON")
	}
}

func TestValidatePayload_MissingFields(t *testing.T) {
	rr := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"name":"Test"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/products", body)

	var input struct {
		Name  string `json:"name"`
		Price int    `json:"price"` // missing in JSON, should be zero value
	}

	result := ValidatePayload(rr, req, &input)

	if !result {
		t.Error("ValidatePayload should return true for valid JSON with missing optional fields")
	}
	if input.Price != 0 {
		t.Errorf("Missing field should be zero value, got: %d", input.Price)
	}
}

func TestValidatePayload_ExtraFields(t *testing.T) {
	rr := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"name":"Test","extra":"field"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/products", body)

	var input struct {
		Name string `json:"name"`
	}

	result := ValidatePayload(rr, req, &input)

	if !result {
		t.Error("ValidatePayload should return true and ignore extra fields")
	}
	if input.Name != "Test" {
		t.Errorf("ValidatePayload should decode name correctly, got: %s", input.Name)
	}
}

func TestValidatePayload_NestedJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"name":"Test","category":{"id":1,"name":"Electronics"}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/products", body)

	var input struct {
		Name     string `json:"name"`
		Category struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"category"`
	}

	result := ValidatePayload(rr, req, &input)

	if !result {
		t.Error("ValidatePayload should return true for nested JSON")
	}
	if input.Category.ID != 1 {
		t.Errorf("ValidatePayload should decode nested category id correctly, got: %d", input.Category.ID)
	}
}

func TestValidatePayload_ArrayJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"items":[{"product_id":1,"quantity":2}]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/checkout", body)

	var input struct {
		Items []struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		} `json:"items"`
	}

	result := ValidatePayload(rr, req, &input)

	if !result {
		t.Error("ValidatePayload should return true for array JSON")
	}
	if len(input.Items) != 1 {
		t.Errorf("ValidatePayload should decode items correctly, got: %d items", len(input.Items))
	}
}

func TestResponse_JSONSerialization(t *testing.T) {
	response := Response{
		Status:  "OK",
		Message: "Success",
		Data:    map[string]int{"id": 1},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Response should marshal to JSON, got error: %v", err)
	}

	var parsed Response
	json.Unmarshal(data, &parsed)

	if parsed.Status != "OK" {
		t.Errorf("Status should be OK, got: %s", parsed.Status)
	}
	if parsed.Message != "Success" {
		t.Errorf("Message should be Success, got: %s", parsed.Message)
	}
}

func TestResponse_JSONSerialization_OmitsEmptyData(t *testing.T) {
	response := Response{
		Status:  "OK",
		Message: "Success",
		Data:    nil,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Response should marshal to JSON, got error: %v", err)
	}

	// Check that data field is omitted
	if bytes.Contains(data, []byte(`"data"`)) {
		t.Error("Response should omit nil data field")
	}
}

// Custom body that fails on close for testing error handling
type failingCloseBody struct {
	*bytes.Reader
}

func (f *failingCloseBody) Close() error {
	return errors.New("close error")
}

func TestValidatePayload_BodyCloseError(t *testing.T) {
	rr := httptest.NewRecorder()

	// Create a request with a body that fails on close
	body := &failingCloseBody{Reader: bytes.NewReader([]byte(`{"name":"Test"}`))}
	req := httptest.NewRequest(http.MethodPost, "/api/products", body)

	var input struct {
		Name string `json:"name"`
	}

	// This should still work, body close errors are silently ignored
	result := ValidatePayload(rr, req, &input)

	if !result {
		t.Error("ValidatePayload should return true even if body close fails")
	}
}

// Custom writer that fails on write for testing error handling
type failingWriter struct {
	header http.Header
}

func (f *failingWriter) Header() http.Header {
	if f.header == nil {
		f.header = make(http.Header)
	}
	return f.header
}

func (f *failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write error")
}

func (f *failingWriter) WriteHeader(statusCode int) {}

func TestWriteJSON_WriteError(t *testing.T) {
	// This tests the error case in WriteJSON when encoding fails
	// The function returns silently on error, so we just verify no panic
	fw := &failingWriter{}
	data := map[string]string{"message": "hello"}

	// Should not panic
	WriteJSON(fw, http.StatusOK, data)
}

func TestWriteJSON_EncodesNilCorrectly(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteJSON(rr, http.StatusOK, nil)

	body, _ := io.ReadAll(rr.Body)
	if string(body) != "null\n" {
		t.Errorf("WriteJSON should encode nil as null, got: %s", string(body))
	}
}
