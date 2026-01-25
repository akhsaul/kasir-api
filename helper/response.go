package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteJSON writes a JSON response with the given status code
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log error but can't change status code since WriteHeader was already called
		// In production, you might want to log this error
		return
	}
}

// WriteSuccess writes a successful JSON response
func WriteSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	WriteJSON(w, statusCode, Response{
		Status:  "OK",
		Message: message,
		Data:    data,
	})
}

// WriteError writes an error JSON response
// Logs detailed error information including HTTP method, request path, status code, message, and error details
// Example log output: [ERROR] POST /api/product - Status: 400 - Message: Invalid JSON - Error: unexpected end of JSON input
func WriteError(w http.ResponseWriter, r *http.Request, statusCode int, message string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s %s - Status: %d - Message: %s - Error: %v",
			r.Method, r.URL.Path, statusCode, message, err)
	} else {
		log.Printf("[ERROR] %s %s - Status: %d - Message: %s",
			r.Method, r.URL.Path, statusCode, message)
	}
	WriteJSON(w, statusCode, Response{
		Status:  "ERROR",
		Message: message,
	})
}

// ValidatePayload decodes JSON payload from request body and closes the body
// Returns false if decoding fails after writing error response
func ValidatePayload(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	isSuccess := true
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		WriteError(w, r, http.StatusBadRequest, "Invalid JSON", err)
		isSuccess = false
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			// Body close error, already processed the request
		}
	}()
	return isSuccess
}
