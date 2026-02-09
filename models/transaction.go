package model

import "time"

// Transaction represents a transaction in the kasir system.
type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details,omitempty"`
}

// TransactionDetail represents a detail item in a transaction.
type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      int    `json:"quantity"`
	Price         int    `json:"price"`
	Subtotal      int    `json:"subtotal"`
}

// CheckoutItem represents an item in the checkout request.
type CheckoutItem struct {
	ProductID int `json:"product_id" validate:"gt=0"`
	Quantity  int `json:"quantity" validate:"gt=0"`
}

// CheckoutRequest represents the request body for checkout.
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items" validate:"required,min=1,dive"`
}

// ReportResponse represents the response for daily/range report.
type ReportResponse struct {
	TotalRevenue   int             `json:"total_revenue"`
	TotalTransaksi int             `json:"total_transaksi"`
	ProdukTerlaris *ProdukTerlaris `json:"produk_terlaris"`
}

// ProdukTerlaris represents the best selling product.
type ProdukTerlaris struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}
