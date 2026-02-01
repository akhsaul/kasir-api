package repository

import (
	"kasir-api/models"
	"sync"
)

// MemoryRepository holds in-memory data and implements ProductRepository.
// Category access is done via CategoryMemoryAdapter.
type MemoryRepository struct {
	mu             sync.RWMutex
	products       map[int]*model.Product
	categories     map[int]*model.Category
	nextProductID  int
	nextCategoryID int
}

// NewMemoryRepository creates a new in-memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		products:       make(map[int]*model.Product),
		categories:     make(map[int]*model.Category),
		nextProductID:  1,
		nextCategoryID: 1,
	}
}

// ProductRepository implementation

func (m *MemoryRepository) GetAll() ([]*model.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	products := make([]*model.Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, p)
	}
	return products, nil
}

func (m *MemoryRepository) GetByID(id int) (*model.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, exists := m.products[id]
	if !exists {
		return nil, model.ErrNotFound
	}
	return p, nil
}

func (m *MemoryRepository) Create(product *model.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	product.ID = m.nextProductID
	m.products[product.ID] = product
	m.nextProductID++
	return nil
}

func (m *MemoryRepository) Update(product *model.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.products[product.ID]; !exists {
		return model.ErrNotFound
	}
	m.products[product.ID] = product
	return nil
}

func (m *MemoryRepository) Delete(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.products[id]; !exists {
		return model.ErrNotFound
	}
	delete(m.products, id)
	return nil
}

// Internal category methods (used by CategoryMemoryAdapter)

func (m *MemoryRepository) getAllCategories() ([]*model.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	categories := make([]*model.Category, 0, len(m.categories))
	for _, c := range m.categories {
		categories = append(categories, c)
	}
	return categories, nil
}

func (m *MemoryRepository) getCategoryByID(id int) (*model.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, exists := m.categories[id]
	if !exists {
		return nil, model.ErrNotFound
	}
	return c, nil
}

func (m *MemoryRepository) createCategory(category *model.Category) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	category.ID = m.nextCategoryID
	m.categories[category.ID] = category
	m.nextCategoryID++
	return nil
}

func (m *MemoryRepository) updateCategory(category *model.Category) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.categories[category.ID]; !exists {
		return model.ErrNotFound
	}
	m.categories[category.ID] = category
	return nil
}

func (m *MemoryRepository) deleteCategory(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.categories[id]; !exists {
		return model.ErrNotFound
	}
	delete(m.categories, id)
	return nil
}
