package postgres

import (
	"database/sql"
	"errors"
	"time"

	model "kasir-api/models"
)

// TransactionRepository implements repository.TransactionRepository using PostgreSQL.
type TransactionRepository struct {
	db *DB
}

// NewTransactionRepository creates a new TransactionRepository.
func NewTransactionRepository(db *DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create inserts a new transaction with its details.
func (r *TransactionRepository) Create(transaction *model.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck // rollback after commit is a no-op

	// Insert transaction
	err = tx.QueryRow(`
		INSERT INTO transactions (total_amount, created_at) VALUES ($1, $2)
		RETURNING id
	`, transaction.TotalAmount, transaction.CreatedAt).Scan(&transaction.ID)
	if err != nil {
		return err
	}

	// Insert transaction details
	for i := range transaction.Details {
		detail := &transaction.Details[i]
		detail.TransactionID = transaction.ID
		err = tx.QueryRow(`
			INSERT INTO transaction_details (transaction_id, product_id, product_name, quantity, price, subtotal)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, detail.TransactionID, detail.ProductID, detail.ProductName, detail.Quantity, detail.Price, detail.Subtotal).Scan(&detail.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetByID returns a transaction by ID with its details.
func (r *TransactionRepository) GetByID(id int) (*model.Transaction, error) {
	var t model.Transaction
	err := r.db.QueryRow(`
		SELECT id, total_amount, created_at FROM transactions WHERE id = $1
	`, id).Scan(&t.ID, &t.TotalAmount, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrTransactionNotFound
		}
		return nil, err
	}

	// Get transaction details
	rows, err := r.db.Query(`
		SELECT id, transaction_id, product_id, product_name, quantity, price, subtotal
		FROM transaction_details WHERE transaction_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d model.TransactionDetail
		if err := rows.Scan(&d.ID, &d.TransactionID, &d.ProductID, &d.ProductName, &d.Quantity, &d.Price, &d.Subtotal); err != nil {
			return nil, err
		}
		t.Details = append(t.Details, d)
	}

	return &t, rows.Err()
}

// GetReportByDateRange returns report data for a given date range.
func (r *TransactionRepository) GetReportByDateRange(startDate, endDate time.Time) (*model.ReportResponse, error) {
	report := &model.ReportResponse{}

	// Get total revenue and total transactions
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// Get best selling product
	var produkTerlaris model.ProdukTerlaris
	err = r.db.QueryRow(`
		SELECT product_name, COALESCE(SUM(quantity), 0) as total_qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY product_name
		ORDER BY total_qty DESC
		LIMIT 1
	`, startDate, endDate).Scan(&produkTerlaris.Nama, &produkTerlaris.QtyTerjual)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		// No transactions in range, produk_terlaris will be nil
	} else {
		report.ProdukTerlaris = &produkTerlaris
	}

	return report, nil
}
