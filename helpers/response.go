package helper

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"kasir-api/helpers/logger"
)

// Response represents a standard API response.
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log error but can't change status code since WriteHeader was already called
		// In production, you might want to log this error
		return
	}
}

// WriteSuccess writes a successful JSON response.
func WriteSuccess(w http.ResponseWriter, statusCode int, message string, data any) {
	WriteJSON(w, statusCode, Response{
		Status:  "OK",
		Message: message,
		Data:    data,
	})
}

// WriteError writes an error JSON response.
// Logs detailed error information including HTTP method, request path, status code, message, and error details.
// Example log output: [ERROR] POST /api/products - Status: 400 - Message: Invalid JSON - Error: unexpected end of JSON input.
func WriteError(w http.ResponseWriter, r *http.Request, statusCode int, message string, err error) {
	if err != nil {
		logger.Error("%s %s - Status: %d - Message: %s - Error: %v",
			r.Method, r.URL.Path, statusCode, message, err)
	} else {
		logger.Error("%s %s - Status: %d - Message: %s",
			r.Method, r.URL.Path, statusCode, message)
	}
	WriteJSON(w, statusCode, Response{
		Status:  "ERROR",
		Message: message,
	})
}

// ParsePagination extracts page and limit from query parameters with defaults.
func ParsePagination(r *http.Request, defaultLimit int) (page, limit int) {
	page = 1
	limit = defaultLimit

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		if l > 100 {
			l = 100 // max limit cap
		}
		limit = l
	}
	return page, limit
}

// ParseIDFromPath extracts an integer ID from a URL path by trimming the given prefix.
// Returns the parsed ID and true on success, or writes a 404 error and returns 0 and false on failure.
func ParseIDFromPath(w http.ResponseWriter, r *http.Request, prefix string, notFoundErr error) (int, bool) {
	idStr := strings.TrimPrefix(r.URL.Path, prefix)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, r, http.StatusNotFound, notFoundErr.Error(), notFoundErr)
		return 0, false
	}
	return id, true
}

// ValidatePayload decodes JSON payload from request body, validates struct tags, and closes the body.
// Returns false if decoding or validation fails after writing an error response.
func ValidatePayload(w http.ResponseWriter, r *http.Request, v any) bool {
	defer r.Body.Close() //nolint:errcheck // body close error is non-actionable

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		WriteError(w, r, http.StatusBadRequest, "Invalid JSON", err)
		return false
	}

	if err := ValidateStruct(v); err != nil {
		WriteError(w, r, http.StatusBadRequest, FormatValidationErrors(err), err)
		return false
	}

	return true
}
