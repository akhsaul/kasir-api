package memory

import (
	"errors"
	"testing"
	"time"

	model "kasir-api/models"
)

func TestNewTransactionRepository(t *testing.T) {
	repo := NewTransactionRepository()

	if repo == nil {
		t.Error("NewTransactionRepository should return a non-nil repository")
	}
	if repo.transactions == nil {
		t.Error("NewTransactionRepository should initialize transactions map")
	}
	if repo.nextID != 1 {
		t.Errorf("NewTransactionRepository should set nextID to 1, got: %d", repo.nextID)
	}
}

func TestTransactionRepository_Create_Success(t *testing.T) {
	repo := NewTransactionRepository()

	transaction := &model.Transaction{
		TotalAmount: 1000,
		CreatedAt:   time.Now(),
		Details: []model.TransactionDetail{
			{ProductID: 1, ProductName: "Laptop", Quantity: 1, Price: 1000, Subtotal: 1000},
		},
	}

	err := repo.Create(transaction)
	if err != nil {
		t.Errorf("Create should not return error, got: %v", err)
	}
	if transaction.ID != 1 {
		t.Errorf("Create should set ID to 1, got: %d", transaction.ID)
	}
	if repo.nextID != 2 {
		t.Errorf("Create should increment nextID to 2, got: %d", repo.nextID)
	}
}

func TestTransactionRepository_Create_AssignsDetailIDs(t *testing.T) {
	repo := NewTransactionRepository()

	transaction := &model.Transaction{
		TotalAmount: 2000,
		CreatedAt:   time.Now(),
		Details: []model.TransactionDetail{
			{ProductID: 1, ProductName: "Laptop", Quantity: 1, Price: 1000, Subtotal: 1000},
			{ProductID: 2, ProductName: "Phone", Quantity: 2, Price: 500, Subtotal: 1000},
		},
	}

	err := repo.Create(transaction)
	if err != nil {
		t.Errorf("Create should not return error, got: %v", err)
	}

	// Check detail IDs
	if transaction.Details[0].ID != 1 {
		t.Errorf("First detail should have ID 1, got: %d", transaction.Details[0].ID)
	}
	if transaction.Details[1].ID != 2 {
		t.Errorf("Second detail should have ID 2, got: %d", transaction.Details[1].ID)
	}

	// Check transaction ID assignment to details
	if transaction.Details[0].TransactionID != 1 {
		t.Errorf("First detail should have TransactionID 1, got: %d", transaction.Details[0].TransactionID)
	}
	if transaction.Details[1].TransactionID != 1 {
		t.Errorf("Second detail should have TransactionID 1, got: %d", transaction.Details[1].TransactionID)
	}
}

func TestTransactionRepository_Create_Multiple(t *testing.T) {
	repo := NewTransactionRepository()

	tx1 := &model.Transaction{TotalAmount: 1000, CreatedAt: time.Now()}
	tx2 := &model.Transaction{TotalAmount: 2000, CreatedAt: time.Now()}

	repo.Create(tx1)
	repo.Create(tx2)

	if tx1.ID != 1 {
		t.Errorf("First transaction should have ID 1, got: %d", tx1.ID)
	}
	if tx2.ID != 2 {
		t.Errorf("Second transaction should have ID 2, got: %d", tx2.ID)
	}
}

func TestTransactionRepository_GetByID_Success(t *testing.T) {
	repo := NewTransactionRepository()

	transaction := &model.Transaction{
		TotalAmount: 1000,
		CreatedAt:   time.Now(),
		Details: []model.TransactionDetail{
			{ProductID: 1, ProductName: "Laptop", Quantity: 1, Price: 1000, Subtotal: 1000},
		},
	}
	repo.Create(transaction)

	retrieved, err := repo.GetByID(1)
	if err != nil {
		t.Errorf("GetByID should not return error, got: %v", err)
	}
	if retrieved.ID != 1 {
		t.Errorf("GetByID should return transaction with ID 1, got: %d", retrieved.ID)
	}
	if retrieved.TotalAmount != 1000 {
		t.Errorf("GetByID should return transaction with total 1000, got: %d", retrieved.TotalAmount)
	}
	if len(retrieved.Details) != 1 {
		t.Errorf("GetByID should return transaction with 1 detail, got: %d", len(retrieved.Details))
	}
}

func TestTransactionRepository_GetByID_NotFound(t *testing.T) {
	repo := NewTransactionRepository()

	_, err := repo.GetByID(999)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("GetByID should return ErrNotFound, got: %v", err)
	}
}

func TestTransactionRepository_GetByID_ReturnsCopy(t *testing.T) {
	repo := NewTransactionRepository()

	transaction := &model.Transaction{
		TotalAmount: 1000,
		CreatedAt:   time.Now(),
		Details: []model.TransactionDetail{
			{ProductID: 1, ProductName: "Laptop", Quantity: 1, Price: 1000, Subtotal: 1000},
		},
	}
	repo.Create(transaction)

	retrieved, _ := repo.GetByID(1)
	retrieved.TotalAmount = 9999
	retrieved.Details[0].ProductName = "Modified"

	// Original should not be affected
	original, _ := repo.GetByID(1)
	if original.TotalAmount == 9999 {
		t.Error("GetByID should return a copy, not reference")
	}
	if original.Details[0].ProductName == "Modified" {
		t.Error("GetByID should return a copy of details, not reference")
	}
}

func TestTransactionRepository_GetReportByDateRange_Empty(t *testing.T) {
	repo := NewTransactionRepository()

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

	report, err := repo.GetReportByDateRange(startDate, endDate)
	if err != nil {
		t.Errorf("GetReportByDateRange should not return error, got: %v", err)
	}
	if report.TotalRevenue != 0 {
		t.Errorf("GetReportByDateRange should return 0 total revenue, got: %d", report.TotalRevenue)
	}
	if report.TotalTransaksi != 0 {
		t.Errorf("GetReportByDateRange should return 0 transactions, got: %d", report.TotalTransaksi)
	}
	if report.ProdukTerlaris != nil {
		t.Error("GetReportByDateRange should return nil produk terlaris for empty result")
	}
}

func TestTransactionRepository_GetReportByDateRange_WithData(t *testing.T) {
	repo := NewTransactionRepository()

	// Create transactions within the date range
	tx1 := &model.Transaction{
		TotalAmount: 1000,
		CreatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Details: []model.TransactionDetail{
			{ProductID: 1, ProductName: "Laptop", Quantity: 1, Price: 1000, Subtotal: 1000},
		},
	}
	tx2 := &model.Transaction{
		TotalAmount: 1500,
		CreatedAt:   time.Date(2024, 1, 20, 14, 0, 0, 0, time.UTC),
		Details: []model.TransactionDetail{
			{ProductID: 1, ProductName: "Laptop", Quantity: 1, Price: 1000, Subtotal: 1000},
			{ProductID: 2, ProductName: "Phone", Quantity: 1, Price: 500, Subtotal: 500},
		},
	}

	repo.Create(tx1)
	repo.Create(tx2)

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	report, err := repo.GetReportByDateRange(startDate, endDate)
	if err != nil {
		t.Errorf("GetReportByDateRange should not return error, got: %v", err)
	}
	if report.TotalRevenue != 2500 {
		t.Errorf("GetReportByDateRange should return 2500 total revenue, got: %d", report.TotalRevenue)
	}
	if report.TotalTransaksi != 2 {
		t.Errorf("GetReportByDateRange should return 2 transactions, got: %d", report.TotalTransaksi)
	}
	if report.ProdukTerlaris == nil {
		t.Error("GetReportByDateRange should return produk terlaris")
	}
	if report.ProdukTerlaris != nil && report.ProdukTerlaris.Nama != "Laptop" {
		t.Errorf("Best selling product should be Laptop, got: %s", report.ProdukTerlaris.Nama)
	}
	if report.ProdukTerlaris != nil && report.ProdukTerlaris.QtyTerjual != 2 {
		t.Errorf("Best selling product qty should be 2, got: %d", report.ProdukTerlaris.QtyTerjual)
	}
}

func TestTransactionRepository_GetReportByDateRange_ExcludesOutOfRange(t *testing.T) {
	repo := NewTransactionRepository()

	// Create transactions - one inside, one outside date range
	txInRange := &model.Transaction{
		TotalAmount: 1000,
		CreatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Details: []model.TransactionDetail{
			{ProductID: 1, ProductName: "Laptop", Quantity: 1, Price: 1000, Subtotal: 1000},
		},
	}
	txOutOfRange := &model.Transaction{
		TotalAmount: 2000,
		CreatedAt:   time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC),
		Details: []model.TransactionDetail{
			{ProductID: 2, ProductName: "Phone", Quantity: 4, Price: 500, Subtotal: 2000},
		},
	}

	repo.Create(txInRange)
	repo.Create(txOutOfRange)

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	report, _ := repo.GetReportByDateRange(startDate, endDate)

	if report.TotalRevenue != 1000 {
		t.Errorf("Should only include in-range transaction, got revenue: %d", report.TotalRevenue)
	}
	if report.TotalTransaksi != 1 {
		t.Errorf("Should only include in-range transaction, got count: %d", report.TotalTransaksi)
	}
}

func TestTransactionRepository_GetReportByDateRange_BeforeStartDate(t *testing.T) {
	repo := NewTransactionRepository()

	tx := &model.Transaction{
		TotalAmount: 1000,
		CreatedAt:   time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
		Details: []model.TransactionDetail{
			{ProductID: 1, ProductName: "Laptop", Quantity: 1, Price: 1000, Subtotal: 1000},
		},
	}
	repo.Create(tx)

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	report, _ := repo.GetReportByDateRange(startDate, endDate)

	if report.TotalRevenue != 0 {
		t.Errorf("Should exclude transactions before start date, got revenue: %d", report.TotalRevenue)
	}
}

func TestTransactionRepository_GetReportByDateRange_AtEndDate(t *testing.T) {
	repo := NewTransactionRepository()

	// Transaction exactly at end date boundary should be excluded (end date is exclusive)
	tx := &model.Transaction{
		TotalAmount: 1000,
		CreatedAt:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		Details: []model.TransactionDetail{
			{ProductID: 1, ProductName: "Laptop", Quantity: 1, Price: 1000, Subtotal: 1000},
		},
	}
	repo.Create(tx)

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	report, _ := repo.GetReportByDateRange(startDate, endDate)

	if report.TotalRevenue != 0 {
		t.Errorf("Should exclude transactions at end date boundary, got revenue: %d", report.TotalRevenue)
	}
}

func TestTransactionRepository_GetReportByDateRange_MultipleBestSelling(t *testing.T) {
	repo := NewTransactionRepository()

	// Create transactions with same qty for different products
	tx := &model.Transaction{
		TotalAmount: 2000,
		CreatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Details: []model.TransactionDetail{
			{ProductID: 1, ProductName: "Laptop", Quantity: 2, Price: 500, Subtotal: 1000},
			{ProductID: 2, ProductName: "Phone", Quantity: 2, Price: 500, Subtotal: 1000},
		},
	}
	repo.Create(tx)

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	report, _ := repo.GetReportByDateRange(startDate, endDate)

	if report.ProdukTerlaris == nil {
		t.Error("Should return one of the best selling products")
	}
	if report.ProdukTerlaris != nil && report.ProdukTerlaris.QtyTerjual != 2 {
		t.Errorf("Best selling qty should be 2, got: %d", report.ProdukTerlaris.QtyTerjual)
	}
}

func TestTransactionRepository_Concurrency(t *testing.T) {
	repo := NewTransactionRepository()

	// Test concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			transaction := &model.Transaction{
				TotalAmount: 1000,
				CreatedAt:   time.Now(),
				Details: []model.TransactionDetail{
					{ProductID: 1, ProductName: "Product", Quantity: 1, Price: 1000, Subtotal: 1000},
				},
			}
			repo.Create(transaction)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all transactions were created
	if repo.nextID != 11 {
		t.Errorf("Should have created 10 transactions, nextID should be 11, got: %d", repo.nextID)
	}
}

func TestTransactionRepository_Create_StoresCopy(t *testing.T) {
	repo := NewTransactionRepository()

	transaction := &model.Transaction{
		TotalAmount: 1000,
		CreatedAt:   time.Now(),
		Details: []model.TransactionDetail{
			{ProductID: 1, ProductName: "Laptop", Quantity: 1, Price: 1000, Subtotal: 1000},
		},
	}
	repo.Create(transaction)

	// Modify original after creation
	transaction.TotalAmount = 9999
	transaction.Details[0].ProductName = "Modified"

	// Stored copy should not be affected
	stored, _ := repo.GetByID(1)
	if stored.TotalAmount == 9999 {
		t.Error("Create should store a copy, not reference")
	}
	if stored.Details[0].ProductName == "Modified" {
		t.Error("Create should store a copy of details, not reference")
	}
}
