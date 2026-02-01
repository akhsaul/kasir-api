package repository

import "kasir-api/models"

// CategoryMemoryAdapter implements CategoryRepository using MemoryRepository.
// Repository layer: data buat logic. Error database → cek sini.
type CategoryMemoryAdapter struct {
	repo *MemoryRepository
}

// NewCategoryMemoryAdapter creates an adapter for category operations.
func NewCategoryMemoryAdapter(repo *MemoryRepository) *CategoryMemoryAdapter {
	return &CategoryMemoryAdapter{repo: repo}
}

func (c *CategoryMemoryAdapter) GetAll() ([]*model.Category, error) {
	return c.repo.getAllCategories()
}

func (c *CategoryMemoryAdapter) GetByID(id int) (*model.Category, error) {
	return c.repo.getCategoryByID(id)
}

func (c *CategoryMemoryAdapter) Create(category *model.Category) error {
	return c.repo.createCategory(category)
}

func (c *CategoryMemoryAdapter) Update(category *model.Category) error {
	return c.repo.updateCategory(category)
}

func (c *CategoryMemoryAdapter) Delete(id int) error {
	return c.repo.deleteCategory(id)
}
