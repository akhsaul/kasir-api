package helper

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
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
		Success: true,
		Message: message,
		Data:    data,
	})
}

// WriteError writes an error JSON response
func WriteError(w http.ResponseWriter, statusCode int, message string, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	WriteJSON(w, statusCode, Response{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}
