package model

// Product represents a product in the kasir system.
// Model layer: definisi bentuk data.
// CategoryID is internal only, tidak diexpose di JSON response.
type Product struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Stock        int    `json:"stock"`
	CategoryID   *int   `json:"-"` // internal only, tidak tampil di response
	CategoryName string `json:"category_name,omitempty"`
}

// ProductInput is the request body for Create/Update product.
// Digunakan untuk parse category_id dari client.
type ProductInput struct {
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID *int   `json:"category_id,omitempty"`
}
