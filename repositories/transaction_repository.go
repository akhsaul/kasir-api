package repository

import (
	"time"

	model "kasir-api/models"
)

// TransactionRepository defines data access for transactions.
// Repository layer: data buat logic. Error database → cek sini.
type TransactionRepository interface {
	Create(transaction *model.Transaction) error
	GetByID(id int) (*model.Transaction, error)
	GetReportByDateRange(startDate, endDate time.Time) (*model.ReportResponse, error)
}
