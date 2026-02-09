package postgres

import (
	"database/sql"
	"errors"

	model "kasir-api/models"
)

// ProductRepository implements repository.ProductRepository using PostgreSQL.
type ProductRepository struct {
	db *DB
}

// NewProductRepository creates a new ProductRepository.
func NewProductRepository(db *DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll returns all products with category info from JOIN.
// If name is provided, filters products by name (case-insensitive partial match).
func (r *ProductRepository) GetAll(name string) ([]*model.Product, error) {
	var rows *sql.Rows
	var err error

	if name != "" {
		rows, err = r.db.Query(`
			SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name, c.description
			FROM products p
			LEFT JOIN categories c ON p.category_id = c.id
			WHERE p.name ILIKE $1
			ORDER BY p.id
		`, "%"+name+"%")
	} else {
		rows, err = r.db.Query(`
			SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name, c.description
			FROM products p
			LEFT JOIN categories c ON p.category_id = c.id
			ORDER BY p.id
		`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		var p model.Product
		var categoryID sql.NullInt64
		var categoryName, categoryDesc sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID, &categoryName, &categoryDesc); err != nil {
			return nil, err
		}
		if categoryID.Valid {
			id := int(categoryID.Int64)
			p.CategoryID = &id
			p.Category = &model.ProductCategory{
				Name:        categoryName.String,
				Description: categoryDesc.String,
			}
		}
		products = append(products, &p)
	}
	return products, rows.Err()
}

// GetByID returns a product by ID with category info from JOIN.
func (r *ProductRepository) GetByID(id int) (*model.Product, error) {
	var p model.Product
	var categoryID sql.NullInt64
	var categoryName, categoryDesc sql.NullString
	err := r.db.QueryRow(`
		SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID, &categoryName, &categoryDesc)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrProductNotFound
		}
		return nil, err
	}
	if categoryID.Valid {
		idVal := int(categoryID.Int64)
		p.CategoryID = &idVal
		p.Category = &model.ProductCategory{
			Name:        categoryName.String,
			Description: categoryDesc.String,
		}
	}
	return &p, nil
}

// Create inserts a new product and returns the generated ID.
// If category_id is set, fetches category info for the response.
func (r *ProductRepository) Create(product *model.Product) error {
	err := r.db.QueryRow(`
		INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4)
		RETURNING id
	`, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	if err != nil {
		return err
	}
	if product.CategoryID != nil {
		var categoryName, categoryDesc sql.NullString
		if err := r.db.QueryRow("SELECT name, description FROM categories WHERE id = $1", *product.CategoryID).Scan(&categoryName, &categoryDesc); err == nil {
			product.Category = &model.ProductCategory{
				Name:        categoryName.String,
				Description: categoryDesc.String,
			}
		}
	}
	return nil
}

// Update updates an existing product.
// If category_id is set, fetches category info for the response.
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
		var categoryName, categoryDesc sql.NullString
		if err := r.db.QueryRow("SELECT name, description FROM categories WHERE id = $1", *product.CategoryID).Scan(&categoryName, &categoryDesc); err == nil {
			product.Category = &model.ProductCategory{
				Name:        categoryName.String,
				Description: categoryDesc.String,
			}
		}
	} else {
		product.Category = nil
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
