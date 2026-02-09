package memory

import (
	"sync"
	"time"

	model "kasir-api/models"
)

// TransactionRepository holds in-memory transaction storage and implements repository.TransactionRepository.
type TransactionRepository struct {
	mu           sync.RWMutex
	transactions map[int]*model.Transaction
	nextID       int
}

// NewTransactionRepository creates a new in-memory transaction repository.
func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{
		transactions: make(map[int]*model.Transaction),
		nextID:       1,
	}
}

// Create inserts a new transaction with its details.
func (r *TransactionRepository) Create(transaction *model.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	transaction.ID = r.nextID
	r.nextID++

	// Assign IDs to details
	for i := range transaction.Details {
		transaction.Details[i].ID = i + 1
		transaction.Details[i].TransactionID = transaction.ID
	}

	// Store a copy
	stored := *transaction
	stored.Details = make([]model.TransactionDetail, len(transaction.Details))
	copy(stored.Details, transaction.Details)
	r.transactions[transaction.ID] = &stored

	return nil
}

// GetByID returns a transaction by ID with its details.
func (r *TransactionRepository) GetByID(id int) (*model.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, exists := r.transactions[id]
	if !exists {
		return nil, model.ErrTransactionNotFound
	}

	// Return a copy
	result := *t
	result.Details = make([]model.TransactionDetail, len(t.Details))
	copy(result.Details, t.Details)
	return &result, nil
}

// GetReportByDateRange returns report data for a given date range.
func (r *TransactionRepository) GetReportByDateRange(startDate, endDate time.Time) (*model.ReportResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	report := &model.ReportResponse{}
	productQty := make(map[string]int)

	for _, t := range r.transactions {
		// Check if transaction is within date range [startDate, endDate)
		// Consistent with PostgreSQL: created_at >= startDate AND created_at < endDate
		if t.CreatedAt.Before(startDate) || !t.CreatedAt.Before(endDate) {
			continue
		}

		report.TotalRevenue += t.TotalAmount
		report.TotalTransaksi++

		for _, d := range t.Details {
			productQty[d.ProductName] += d.Quantity
		}
	}

	// Find best selling product
	var maxQty int
	var bestProduct string
	for name, qty := range productQty {
		if qty > maxQty {
			maxQty = qty
			bestProduct = name
		}
	}

	if bestProduct != "" {
		report.ProdukTerlaris = &model.ProdukTerlaris{
			Nama:       bestProduct,
			QtyTerjual: maxQty,
		}
	}

	return report, nil
}
