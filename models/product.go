package model

// Product represents a product in the kasir system.
// Model layer: definisi bentuk data.
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}
