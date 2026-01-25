package data

import "kasir-api/entity"

// CategoryMemoryStorage is an adapter that implements CategoryStorage interface using MemoryStorage
type CategoryMemoryStorage struct {
	storage *MemoryStorage
}

// NewCategoryMemoryStorage creates a new CategoryMemoryStorage adapter
func NewCategoryMemoryStorage(storage *MemoryStorage) *CategoryMemoryStorage {
	return &CategoryMemoryStorage{
		storage: storage,
	}
}

// GetAll returns all categories
func (c *CategoryMemoryStorage) GetAll() ([]*entity.Category, error) {
	return c.storage.GetAllCategories()
}

// GetByID returns a category by ID
func (c *CategoryMemoryStorage) GetByID(id int) (*entity.Category, error) {
	return c.storage.GetCategoryByID(id)
}

// Create adds a new category
func (c *CategoryMemoryStorage) Create(category *entity.Category) error {
	return c.storage.CreateCategory(category)
}

// Update modifies an existing category
func (c *CategoryMemoryStorage) Update(category *entity.Category) error {
	return c.storage.UpdateCategory(category)
}

// Delete removes a category by ID
func (c *CategoryMemoryStorage) Delete(id int) error {
	return c.storage.DeleteCategory(id)
}
