package service

import (
	"errors"
	"testing"
	"time"

	"kasir-api/mocks"
	model "kasir-api/models"
)

func TestNewTransactionService(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	if service == nil {
		t.Error("NewTransactionService should return a non-nil service")
	}
	if service.repo != transactionRepo {
		t.Error("NewTransactionService should set the transaction repo")
	}
	if service.productRepo != productRepo {
		t.Error("NewTransactionService should set the product repo")
	}
}

func TestTransactionService_Checkout_Success(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	productRepo.Products[1] = &model.Product{ID: 1, Name: "Laptop", Price: 1000, Stock: 10}
	productRepo.Products[2] = &model.Product{ID: 2, Name: "Phone", Price: 500, Stock: 20}

	request := &model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}

	transaction, err := service.Checkout(request)
	if err != nil {
		t.Errorf("Checkout should not return error, got: %v", err)
	}
	if transaction.ID != 1 {
		t.Errorf("Checkout should return transaction with ID 1, got: %d", transaction.ID)
	}
	// Total: 2*1000 + 3*500 = 3500
	if transaction.TotalAmount != 3500 {
		t.Errorf("Checkout should return transaction with total 3500, got: %d", transaction.TotalAmount)
	}
	if len(transaction.Details) != 2 {
		t.Errorf("Checkout should return transaction with 2 details, got: %d", len(transaction.Details))
	}

	// Verify stock is reduced
	if productRepo.Products[1].Stock != 8 {
		t.Errorf("Product 1 stock should be 8, got: %d", productRepo.Products[1].Stock)
	}
	if productRepo.Products[2].Stock != 17 {
		t.Errorf("Product 2 stock should be 17, got: %d", productRepo.Products[2].Stock)
	}
}

func TestTransactionService_Checkout_EmptyItems(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	request := &model.CheckoutRequest{
		Items: []model.CheckoutItem{},
	}

	_, err := service.Checkout(request)
	if !errors.Is(err, model.ErrEmptyCheckout) {
		t.Errorf("Checkout with empty items should return ErrEmptyCheckout, got: %v", err)
	}
}

func TestTransactionService_Checkout_InvalidQuantity(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	testCases := []struct {
		name     string
		quantity int
	}{
		{"zero quantity", 0},
		{"negative quantity", -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := &model.CheckoutRequest{
				Items: []model.CheckoutItem{
					{ProductID: 1, Quantity: tc.quantity},
				},
			}

			_, err := service.Checkout(request)
			if !errors.Is(err, model.ErrInvalidQuantity) {
				t.Errorf("Checkout with %s should return ErrInvalidQuantity, got: %v", tc.name, err)
			}
		})
	}
}

func TestTransactionService_Checkout_ProductNotFound(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	request := &model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 999, Quantity: 1},
		},
	}

	_, err := service.Checkout(request)
	if !errors.Is(err, model.ErrProductNotFound) {
		t.Errorf("Checkout with non-existent product should return ErrProductNotFound, got: %v", err)
	}
}

func TestTransactionService_Checkout_InsufficientStock(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	productRepo.Products[1] = &model.Product{ID: 1, Name: "Laptop", Price: 1000, Stock: 5}

	request := &model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 1, Quantity: 10},
		},
	}

	_, err := service.Checkout(request)
	if !errors.Is(err, model.ErrInsufficientStock) {
		t.Errorf("Checkout with insufficient stock should return ErrInsufficientStock, got: %v", err)
	}
}

func TestTransactionService_Checkout_ProductUpdateError(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	productRepo.Products[1] = &model.Product{ID: 1, Name: "Laptop", Price: 1000, Stock: 10}

	expectedErr := errors.New("database error")
	productRepo.UpdateFunc = func(product *model.Product) error {
		return expectedErr
	}

	request := &model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 1, Quantity: 2},
		},
	}

	_, err := service.Checkout(request)
	if err != expectedErr {
		t.Errorf("Checkout should return the error from product update, got: %v", err)
	}
}

func TestTransactionService_Checkout_TransactionCreateError(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	productRepo.Products[1] = &model.Product{ID: 1, Name: "Laptop", Price: 1000, Stock: 10}

	expectedErr := errors.New("database error")
	transactionRepo.CreateFunc = func(transaction *model.Transaction) error {
		return expectedErr
	}

	request := &model.CheckoutRequest{
		Items: []model.CheckoutItem{
			{ProductID: 1, Quantity: 2},
		},
	}

	_, err := service.Checkout(request)
	if err != expectedErr {
		t.Errorf("Checkout should return the error from transaction create, got: %v", err)
	}
}

func TestTransactionService_GetByID_Success(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	transactionRepo.Transactions[1] = &model.Transaction{
		ID:          1,
		TotalAmount: 1000,
		CreatedAt:   time.Now(),
		Details: []model.TransactionDetail{
			{ID: 1, ProductID: 1, ProductName: "Laptop", Quantity: 1, Price: 1000, Subtotal: 1000},
		},
	}

	transaction, err := service.GetByID(1)
	if err != nil {
		t.Errorf("GetByID should not return error, got: %v", err)
	}
	if transaction.ID != 1 {
		t.Errorf("GetByID should return transaction with ID 1, got: %d", transaction.ID)
	}
	if transaction.TotalAmount != 1000 {
		t.Errorf("GetByID should return transaction with total 1000, got: %d", transaction.TotalAmount)
	}
}

func TestTransactionService_GetByID_InvalidID(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	testCases := []struct {
		name string
		id   int
	}{
		{"zero", 0},
		{"negative", -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.GetByID(tc.id)
			if !errors.Is(err, model.ErrNotFound) {
				t.Errorf("GetByID with %s id should return ErrNotFound, got: %v", tc.name, err)
			}
		})
	}
}

func TestTransactionService_GetByID_NotFound(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	_, err := service.GetByID(999)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("GetByID should return ErrNotFound, got: %v", err)
	}
}

func TestTransactionService_GetTodayReport_Success(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	expectedReport := &model.ReportResponse{
		TotalRevenue:   5000,
		TotalTransaksi: 3,
		ProdukTerlaris: &model.ProdukTerlaris{
			Nama:       "Laptop",
			QtyTerjual: 10,
		},
	}

	transactionRepo.GetReportByDateRangeFunc = func(startDate, endDate time.Time) (*model.ReportResponse, error) {
		// Verify that the date range is for today
		now := time.Now()
		expectedStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		expectedEnd := expectedStart.AddDate(0, 0, 1)

		if !startDate.Equal(expectedStart) {
			t.Errorf("Start date should be start of today, got: %v", startDate)
		}
		if !endDate.Equal(expectedEnd) {
			t.Errorf("End date should be start of tomorrow, got: %v", endDate)
		}
		return expectedReport, nil
	}

	report, err := service.GetTodayReport()
	if err != nil {
		t.Errorf("GetTodayReport should not return error, got: %v", err)
	}
	if report.TotalRevenue != 5000 {
		t.Errorf("GetTodayReport should return report with total revenue 5000, got: %d", report.TotalRevenue)
	}
	if report.TotalTransaksi != 3 {
		t.Errorf("GetTodayReport should return report with 3 transactions, got: %d", report.TotalTransaksi)
	}
}

func TestTransactionService_GetTodayReport_Error(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	expectedErr := errors.New("database error")
	transactionRepo.GetReportByDateRangeFunc = func(startDate, endDate time.Time) (*model.ReportResponse, error) {
		return nil, expectedErr
	}

	_, err := service.GetTodayReport()
	if err != expectedErr {
		t.Errorf("GetTodayReport should return the error from repo, got: %v", err)
	}
}

func TestTransactionService_GetReportByDateRange_Success(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	expectedReport := &model.ReportResponse{
		TotalRevenue:   10000,
		TotalTransaksi: 5,
	}

	transactionRepo.GetReportByDateRangeFunc = func(startDate, endDate time.Time) (*model.ReportResponse, error) {
		return expectedReport, nil
	}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	report, err := service.GetReportByDateRange(startDate, endDate)
	if err != nil {
		t.Errorf("GetReportByDateRange should not return error, got: %v", err)
	}
	if report.TotalRevenue != 10000 {
		t.Errorf("GetReportByDateRange should return report with total revenue 10000, got: %d", report.TotalRevenue)
	}
}

func TestTransactionService_GetReportByDateRange_AddsOneDay(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	transactionRepo.GetReportByDateRangeFunc = func(start, end time.Time) (*model.ReportResponse, error) {
		// Verify that one day is added to endDate
		expectedEnd := endDate.AddDate(0, 0, 1)
		if !end.Equal(expectedEnd) {
			t.Errorf("End date should have one day added, expected: %v, got: %v", expectedEnd, end)
		}
		return &model.ReportResponse{}, nil
	}

	_, err := service.GetReportByDateRange(startDate, endDate)
	if err != nil {
		t.Errorf("GetReportByDateRange should not return error, got: %v", err)
	}
}

func TestTransactionService_GetReportByDateRange_Error(t *testing.T) {
	transactionRepo := mocks.NewMockTransactionRepository()
	productRepo := mocks.NewMockProductRepository()
	service := NewTransactionService(transactionRepo, productRepo)

	expectedErr := errors.New("database error")
	transactionRepo.GetReportByDateRangeFunc = func(startDate, endDate time.Time) (*model.ReportResponse, error) {
		return nil, expectedErr
	}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	_, err := service.GetReportByDateRange(startDate, endDate)
	if err != expectedErr {
		t.Errorf("GetReportByDateRange should return the error from repo, got: %v", err)
	}
}
