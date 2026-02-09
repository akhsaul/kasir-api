package model

// Product represents a product in the kasir system.
// Model layer: definisi bentuk data.
// CategoryID is internal only, tidak diexpose di JSON response.
type Product struct {
	ID         int              `json:"id"`
	Name       string           `json:"name"`
	Price      int              `json:"price"`
	Stock      int              `json:"stock"`
	CategoryID *int             `json:"-"` // internal only, tidak tampil di response
	Category   *ProductCategory `json:"category,omitempty"`
}

// ProductCategory represents category info embedded in product response.
type ProductCategory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProductInput is the request body for Create/Update product.
// Digunakan untuk parse category_id dari client.
type ProductInput struct {
	Name       string `json:"name" validate:"required"`
	Price      int    `json:"price" validate:"gt=0"`
	Stock      int    `json:"stock" validate:"gte=0"`
	CategoryID *int   `json:"category_id,omitempty" validate:"omitempty,gt=0"`
}
