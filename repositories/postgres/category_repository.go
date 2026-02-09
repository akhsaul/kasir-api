package postgres

import (
	"database/sql"
	"errors"

	model "kasir-api/models"
)

// CategoryRepository implements repository.CategoryRepository using PostgreSQL.
type CategoryRepository struct {
	db *DB
}

// NewCategoryRepository creates a new CategoryRepository.
func NewCategoryRepository(db *DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll returns all categories.
func (r *CategoryRepository) GetAll() ([]*model.Category, error) {
	rows, err := r.db.Query(`
		SELECT id, name, description FROM categories ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, rows.Err()
}

// GetByID returns a category by ID.
func (r *CategoryRepository) GetByID(id int) (*model.Category, error) {
	var c model.Category
	err := r.db.QueryRow(`
		SELECT id, name, description FROM categories WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrCategoryNotFound
		}
		return nil, err
	}
	return &c, nil
}

// Create inserts a new category and returns the generated ID.
func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.QueryRow(`
		INSERT INTO categories (name, description) VALUES ($1, $2)
		RETURNING id
	`, category.Name, category.Description).Scan(&category.ID)
}

// Update updates an existing category.
func (r *CategoryRepository) Update(category *model.Category) error {
	result, err := r.db.Exec(`
		UPDATE categories SET name = $1, description = $2 WHERE id = $3
	`, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return model.ErrCategoryNotFound
	}
	return nil
}

// Delete removes a category by ID.
func (r *CategoryRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return model.ErrCategoryNotFound
	}
	return nil
}
