package service

import (
	"time"

	model "kasir-api/models"
	repository "kasir-api/repositories"
)

// TransactionService handles business logic for transactions.
// Service layer: logic kode kita. Error logic → cek sini.
type TransactionService struct {
	repo        repository.TransactionRepository
	productRepo repository.ProductRepository
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(repo repository.TransactionRepository, productRepo repository.ProductRepository) *TransactionService {
	return &TransactionService{repo: repo, productRepo: productRepo}
}

// Checkout processes a checkout request and creates a transaction.
func (s *TransactionService) Checkout(request *model.CheckoutRequest) (*model.Transaction, error) {
	if len(request.Items) == 0 {
		return nil, model.ErrEmptyCheckout
	}

	transaction := &model.Transaction{
		CreatedAt: time.Now(),
	}

	totalAmount := 0
	for _, item := range request.Items {
		if item.Quantity <= 0 {
			return nil, model.ErrInvalidQuantity
		}

		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, err
		}

		if product.Stock < item.Quantity {
			return nil, model.ErrInsufficientStock
		}

		subtotal := product.Price * item.Quantity
		totalAmount += subtotal

		detail := model.TransactionDetail{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Price:       product.Price,
			Subtotal:    subtotal,
		}
		transaction.Details = append(transaction.Details, detail)

		// Update product stock
		product.Stock -= item.Quantity
		if err := s.productRepo.Update(product); err != nil {
			return nil, err
		}
	}

	transaction.TotalAmount = totalAmount

	if err := s.repo.Create(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetByID retrieves a transaction by ID.
func (s *TransactionService) GetByID(id int) (*model.Transaction, error) {
	if id <= 0 {
		return nil, model.ErrTransactionNotFound
	}
	return s.repo.GetByID(id)
}

// GetTodayReport returns the report for today.
func (s *TransactionService) GetTodayReport() (*model.ReportResponse, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)
	return s.repo.GetReportByDateRange(startOfDay, endOfDay)
}

// GetReportByDateRange returns the report for a given date range.
func (s *TransactionService) GetReportByDateRange(startDate, endDate time.Time) (*model.ReportResponse, error) {
	// Add one day to endDate to include the entire end day
	endDate = endDate.AddDate(0, 0, 1)
	return s.repo.GetReportByDateRange(startDate, endDate)
}
