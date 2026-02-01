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

// GetAll returns all products with category name from JOIN.
func (r *ProductRepository) GetAll() ([]*model.Product, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		var p model.Product
		var categoryID sql.NullInt64
		var categoryName sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID, &categoryName); err != nil {
			return nil, err
		}
		if categoryID.Valid {
			id := int(categoryID.Int64)
			p.CategoryID = &id
		}
		if categoryName.Valid {
			p.CategoryName = categoryName.String
		}
		products = append(products, &p)
	}
	return products, rows.Err()
}

// GetByID returns a product by ID with category name from JOIN.
func (r *ProductRepository) GetByID(id int) (*model.Product, error) {
	var p model.Product
	var categoryID sql.NullInt64
	var categoryName sql.NullString
	err := r.db.QueryRow(`
		SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID, &categoryName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrProductNotFound
		}
		return nil, err
	}
	if categoryID.Valid {
		idVal := int(categoryID.Int64)
		p.CategoryID = &idVal
	}
	if categoryName.Valid {
		p.CategoryName = categoryName.String
	}
	return &p, nil
}

// Create inserts a new product and returns the generated ID.
// If category_id is set, fetches category_name for the response.
func (r *ProductRepository) Create(product *model.Product) error {
	err := r.db.QueryRow(`
		INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4)
		RETURNING id
	`, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	if err != nil {
		return err
	}
	if product.CategoryID != nil {
		var categoryName sql.NullString
		if err := r.db.QueryRow("SELECT name FROM categories WHERE id = $1", *product.CategoryID).Scan(&categoryName); err == nil && categoryName.Valid {
			product.CategoryName = categoryName.String
		}
	}
	return nil
}

// Update updates an existing product.
// If category_id is set, fetches category_name for the response.
func (r *ProductRepository) Update(product *model.Product) error {
	result, err := r.db.Exec(`
		UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5
	`, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return model.ErrProductNotFound
	}
	if product.CategoryID != nil {
		var categoryName sql.NullString
		if err := r.db.QueryRow("SELECT name FROM categories WHERE id = $1", *product.CategoryID).Scan(&categoryName); err == nil && categoryName.Valid {
			product.CategoryName = categoryName.String
		}
	} else {
		product.CategoryName = ""
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
		return model.ErrProductNotFound
	}
	return nil
}
