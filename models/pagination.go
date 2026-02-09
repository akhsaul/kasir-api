package model

// PaginatedResponse wraps a list response with pagination metadata.
type PaginatedResponse struct {
	Items      any `json:"items"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}
