package postgres

import (
	"database/sql"
	"errors"

	"kasir-api/models"
)

// ProductRepository implements repository.ProductRepository using PostgreSQL.
type ProductRepository struct {
	db *DB
}

// NewProductRepository creates a new ProductRepository.
func NewProductRepository(db *DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll returns all products.
func (r *ProductRepository) GetAll() ([]*model.Product, error) {
	rows, err := r.db.Query(`
		SELECT id, name, price, stock FROM products ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, rows.Err()
}

// GetByID returns a product by ID.
func (r *ProductRepository) GetByID(id int) (*model.Product, error) {
	var p model.Product
	err := r.db.QueryRow(`
		SELECT id, name, price, stock FROM products WHERE id = $1
	`, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

// Create inserts a new product and returns the generated ID.
func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.QueryRow(`
		INSERT INTO products (name, price, stock) VALUES ($1, $2, $3)
		RETURNING id
	`, product.Name, product.Price, product.Stock).Scan(&product.ID)
}

// Update updates an existing product.
func (r *ProductRepository) Update(product *model.Product) error {
	result, err := r.db.Exec(`
		UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4
	`, product.Name, product.Price, product.Stock, product.ID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

// Delete removes a product by ID.
func (r *ProductRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}
