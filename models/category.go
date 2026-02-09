package model

// Category represents a category in the kasir system.
// Model layer: definisi bentuk data.
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
