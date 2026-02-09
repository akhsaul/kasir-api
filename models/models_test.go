package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestProduct_JSONSerialization(t *testing.T) {
	product := &Product{
		ID:    1,
		Name:  "Laptop",
		Price: 1000,
		Stock: 10,
	}

	data, err := json.Marshal(product)
	if err != nil {
		t.Errorf("Product should marshal to JSON, got error: %v", err)
	}

	var parsed Product
	json.Unmarshal(data, &parsed)

	if parsed.ID != 1 {
		t.Errorf("ID should be 1, got: %d", parsed.ID)
	}
	if parsed.Name != "Laptop" {
		t.Errorf("Name should be Laptop, got: %s", parsed.Name)
	}
	if parsed.Price != 1000 {
		t.Errorf("Price should be 1000, got: %d", parsed.Price)
	}
	if parsed.Stock != 10 {
		t.Errorf("Stock should be 10, got: %d", parsed.Stock)
	}
}

func TestProduct_CategoryIDNotInJSON(t *testing.T) {
	categoryID := 1
	product := &Product{
		ID:         1,
		Name:       "Laptop",
		Price:      1000,
		Stock:      10,
		CategoryID: &categoryID,
	}

	data, _ := json.Marshal(product)

	// CategoryID should not appear in JSON (json:"-")
	if string(data) != `{"id":1,"name":"Laptop","price":1000,"stock":10}` {
		t.Errorf("CategoryID should not be in JSON, got: %s", string(data))
	}
}

func TestProduct_WithCategory(t *testing.T) {
	product := &Product{
		ID:    1,
		Name:  "Laptop",
		Price: 1000,
		Stock: 10,
		Category: &ProductCategory{
			Name:        "Electronics",
			Description: "Electronic items",
		},
	}

	data, _ := json.Marshal(product)

	var parsed map[string]interface{}
	json.Unmarshal(data, &parsed)

	category, ok := parsed["category"].(map[string]interface{})
	if !ok {
		t.Error("Category should be in JSON")
	}
	if category["name"] != "Electronics" {
		t.Errorf("Category name should be Electronics, got: %v", category["name"])
	}
}

func TestProduct_WithoutCategory(t *testing.T) {
	product := &Product{
		ID:       1,
		Name:     "Laptop",
		Price:    1000,
		Stock:    10,
		Category: nil,
	}

	data, _ := json.Marshal(product)

	var parsed map[string]interface{}
	json.Unmarshal(data, &parsed)

	if _, ok := parsed["category"]; ok {
		t.Error("Category should be omitted when nil")
	}
}

func TestProductInput_JSONSerialization(t *testing.T) {
	categoryID := 1
	input := &ProductInput{
		Name:       "Laptop",
		Price:      1000,
		Stock:      10,
		CategoryID: &categoryID,
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Errorf("ProductInput should marshal to JSON, got error: %v", err)
	}

	var parsed ProductInput
	json.Unmarshal(data, &parsed)

	if parsed.Name != "Laptop" {
		t.Errorf("Name should be Laptop, got: %s", parsed.Name)
	}
	if parsed.CategoryID == nil || *parsed.CategoryID != 1 {
		t.Error("CategoryID should be 1")
	}
}

func TestProductInput_WithoutCategoryID(t *testing.T) {
	input := &ProductInput{
		Name:  "Laptop",
		Price: 1000,
		Stock: 10,
	}

	data, _ := json.Marshal(input)

	var parsed map[string]interface{}
	json.Unmarshal(data, &parsed)

	if _, ok := parsed["category_id"]; ok {
		t.Error("CategoryID should be omitted when nil")
	}
}

func TestProductCategory_JSONSerialization(t *testing.T) {
	category := &ProductCategory{
		Name:        "Electronics",
		Description: "Electronic items",
	}

	data, err := json.Marshal(category)
	if err != nil {
		t.Errorf("ProductCategory should marshal to JSON, got error: %v", err)
	}

	var parsed ProductCategory
	json.Unmarshal(data, &parsed)

	if parsed.Name != "Electronics" {
		t.Errorf("Name should be Electronics, got: %s", parsed.Name)
	}
	if parsed.Description != "Electronic items" {
		t.Errorf("Description should be Electronic items, got: %s", parsed.Description)
	}
}

func TestCategory_JSONSerialization(t *testing.T) {
	category := &Category{
		ID:          1,
		Name:        "Electronics",
		Description: "Electronic items",
	}

	data, err := json.Marshal(category)
	if err != nil {
		t.Errorf("Category should marshal to JSON, got error: %v", err)
	}

	var parsed Category
	json.Unmarshal(data, &parsed)

	if parsed.ID != 1 {
		t.Errorf("ID should be 1, got: %d", parsed.ID)
	}
	if parsed.Name != "Electronics" {
		t.Errorf("Name should be Electronics, got: %s", parsed.Name)
	}
	if parsed.Description != "Electronic items" {
		t.Errorf("Description should be Electronic items, got: %s", parsed.Description)
	}
}

func TestTransaction_JSONSerialization(t *testing.T) {
	now := time.Now()
	transaction := &Transaction{
		ID:          1,
		TotalAmount: 3500,
		CreatedAt:   now,
		Details: []TransactionDetail{
			{ID: 1, TransactionID: 1, ProductID: 1, ProductName: "Laptop", Quantity: 2, Price: 1000, Subtotal: 2000},
			{ID: 2, TransactionID: 1, ProductID: 2, ProductName: "Phone", Quantity: 3, Price: 500, Subtotal: 1500},
		},
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		t.Errorf("Transaction should marshal to JSON, got error: %v", err)
	}

	var parsed Transaction
	json.Unmarshal(data, &parsed)

	if parsed.ID != 1 {
		t.Errorf("ID should be 1, got: %d", parsed.ID)
	}
	if parsed.TotalAmount != 3500 {
		t.Errorf("TotalAmount should be 3500, got: %d", parsed.TotalAmount)
	}
	if len(parsed.Details) != 2 {
		t.Errorf("Details should have 2 items, got: %d", len(parsed.Details))
	}
}

func TestTransaction_WithoutDetails(t *testing.T) {
	transaction := &Transaction{
		ID:          1,
		TotalAmount: 1000,
		CreatedAt:   time.Now(),
	}

	data, _ := json.Marshal(transaction)

	var parsed map[string]interface{}
	json.Unmarshal(data, &parsed)

	if _, ok := parsed["details"]; ok {
		t.Error("Details should be omitted when empty/nil")
	}
}

func TestTransactionDetail_JSONSerialization(t *testing.T) {
	detail := &TransactionDetail{
		ID:            1,
		TransactionID: 1,
		ProductID:     1,
		ProductName:   "Laptop",
		Quantity:      2,
		Price:         1000,
		Subtotal:      2000,
	}

	data, err := json.Marshal(detail)
	if err != nil {
		t.Errorf("TransactionDetail should marshal to JSON, got error: %v", err)
	}

	var parsed TransactionDetail
	json.Unmarshal(data, &parsed)

	if parsed.ID != 1 {
		t.Errorf("ID should be 1, got: %d", parsed.ID)
	}
	if parsed.ProductName != "Laptop" {
		t.Errorf("ProductName should be Laptop, got: %s", parsed.ProductName)
	}
	if parsed.Subtotal != 2000 {
		t.Errorf("Subtotal should be 2000, got: %d", parsed.Subtotal)
	}
}

func TestCheckoutItem_JSONSerialization(t *testing.T) {
	item := &CheckoutItem{
		ProductID: 1,
		Quantity:  2,
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Errorf("CheckoutItem should marshal to JSON, got error: %v", err)
	}

	var parsed CheckoutItem
	json.Unmarshal(data, &parsed)

	if parsed.ProductID != 1 {
		t.Errorf("ProductID should be 1, got: %d", parsed.ProductID)
	}
	if parsed.Quantity != 2 {
		t.Errorf("Quantity should be 2, got: %d", parsed.Quantity)
	}
}

func TestCheckoutRequest_JSONSerialization(t *testing.T) {
	request := &CheckoutRequest{
		Items: []CheckoutItem{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Errorf("CheckoutRequest should marshal to JSON, got error: %v", err)
	}

	var parsed CheckoutRequest
	json.Unmarshal(data, &parsed)

	if len(parsed.Items) != 2 {
		t.Errorf("Items should have 2 elements, got: %d", len(parsed.Items))
	}
	if parsed.Items[0].ProductID != 1 {
		t.Errorf("First item ProductID should be 1, got: %d", parsed.Items[0].ProductID)
	}
}

func TestCheckoutRequest_JSONDecoding(t *testing.T) {
	jsonStr := `{"items":[{"product_id":1,"quantity":2},{"product_id":2,"quantity":3}]}`

	var request CheckoutRequest
	err := json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		t.Errorf("CheckoutRequest should unmarshal from JSON, got error: %v", err)
	}
	if len(request.Items) != 2 {
		t.Errorf("Items should have 2 elements, got: %d", len(request.Items))
	}
}

func TestReportResponse_JSONSerialization(t *testing.T) {
	report := &ReportResponse{
		TotalRevenue:   10000,
		TotalTransaksi: 5,
		ProdukTerlaris: &ProdukTerlaris{
			Nama:       "Laptop",
			QtyTerjual: 10,
		},
	}

	data, err := json.Marshal(report)
	if err != nil {
		t.Errorf("ReportResponse should marshal to JSON, got error: %v", err)
	}

	var parsed ReportResponse
	json.Unmarshal(data, &parsed)

	if parsed.TotalRevenue != 10000 {
		t.Errorf("TotalRevenue should be 10000, got: %d", parsed.TotalRevenue)
	}
	if parsed.TotalTransaksi != 5 {
		t.Errorf("TotalTransaksi should be 5, got: %d", parsed.TotalTransaksi)
	}
	if parsed.ProdukTerlaris == nil {
		t.Error("ProdukTerlaris should not be nil")
	}
	if parsed.ProdukTerlaris.Nama != "Laptop" {
		t.Errorf("ProdukTerlaris.Nama should be Laptop, got: %s", parsed.ProdukTerlaris.Nama)
	}
}

func TestReportResponse_WithNilProdukTerlaris(t *testing.T) {
	report := &ReportResponse{
		TotalRevenue:   0,
		TotalTransaksi: 0,
		ProdukTerlaris: nil,
	}

	data, _ := json.Marshal(report)

	var parsed map[string]interface{}
	json.Unmarshal(data, &parsed)

	// ProdukTerlaris is a pointer, when nil it becomes JSON null
	if parsed["produk_terlaris"] != nil {
		t.Error("ProdukTerlaris should be null in JSON when nil")
	}
}

func TestProdukTerlaris_JSONSerialization(t *testing.T) {
	produk := &ProdukTerlaris{
		Nama:       "Laptop",
		QtyTerjual: 10,
	}

	data, err := json.Marshal(produk)
	if err != nil {
		t.Errorf("ProdukTerlaris should marshal to JSON, got error: %v", err)
	}

	var parsed ProdukTerlaris
	json.Unmarshal(data, &parsed)

	if parsed.Nama != "Laptop" {
		t.Errorf("Nama should be Laptop, got: %s", parsed.Nama)
	}
	if parsed.QtyTerjual != 10 {
		t.Errorf("QtyTerjual should be 10, got: %d", parsed.QtyTerjual)
	}
}

func TestErrors_AreDefinedCorrectly(t *testing.T) {
	testCases := []struct {
		err      error
		expected string
	}{
		{ErrNotFound, "not found"},
		{ErrCategoryNotFound, "category is not found: not found"},
		{ErrProductNotFound, "product is not found: not found"},
		{ErrTransactionNotFound, "transaction is not found: not found"},
		{ErrNameRequired, "name should not be empty"},
		{ErrPriceInvalid, "price must be greater than 0"},
		{ErrStockInvalid, "stock must be greater than or equal to 0"},
		{ErrIDRequired, "id is required"},
		{ErrEmptyCheckout, "checkout items cannot be empty"},
		{ErrInvalidQuantity, "quantity must be greater than 0"},
		{ErrInsufficientStock, "insufficient stock"},
		{ErrInvalidDateRange, "invalid date range"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			if tc.err.Error() != tc.expected {
				t.Errorf("Error message should be '%s', got: '%s'", tc.expected, tc.err.Error())
			}
		})
	}
}
