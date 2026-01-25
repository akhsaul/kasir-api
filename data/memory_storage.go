package data

import (
	"kasir-api/entity"
	"sync"
)

// MemoryStorage implements ProductStorage and CategoryStorage interfaces with in-memory storage
type MemoryStorage struct {
	mu             sync.RWMutex
	products       map[int]*entity.Product
	categories     map[int]*entity.Category
	nextProductID  int
	nextCategoryID int
}

// NewMemoryStorage creates a new instance of MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		products:       make(map[int]*entity.Product),
		categories:     make(map[int]*entity.Category),
		nextProductID:  1,
		nextCategoryID: 1,
	}
}

// GetAll returns all products
func (m *MemoryStorage) GetAll() ([]*entity.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	products := make([]*entity.Product, 0, len(m.products))
	for _, product := range m.products {
		products = append(products, product)
	}
	return products, nil
}

// GetByID returns a product by ID
func (m *MemoryStorage) GetByID(id int) (*entity.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	product, exists := m.products[id]
	if !exists {
		return nil, entity.ErrNotFound
	}
	return product, nil
}

// Create adds a new product
func (m *MemoryStorage) Create(product *entity.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	product.ID = m.nextProductID
	m.products[product.ID] = product
	m.nextProductID++
	return nil
}

// Update modifies an existing product
func (m *MemoryStorage) Update(product *entity.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.products[product.ID]; !exists {
		return entity.ErrNotFound
	}
	m.products[product.ID] = product
	return nil
}

// Delete removes a product by ID
func (m *MemoryStorage) Delete(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.products[id]; !exists {
		return entity.ErrNotFound
	}
	delete(m.products, id)
	return nil
}

// Category Storage Methods (implementing CategoryStorage interface)

// GetAllCategories GetAll returns all categories (for CategoryStorage interface)
func (m *MemoryStorage) GetAllCategories() ([]*entity.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	categories := make([]*entity.Category, 0, len(m.categories))
	for _, category := range m.categories {
		categories = append(categories, category)
	}
	return categories, nil
}

// GetCategoryByID GetByID returns a category by ID (for CategoryStorage interface)
func (m *MemoryStorage) GetCategoryByID(id int) (*entity.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	category, exists := m.categories[id]
	if !exists {
		return nil, entity.ErrNotFound
	}
	return category, nil
}

// CreateCategory Create adds a new category (for CategoryStorage interface)
func (m *MemoryStorage) CreateCategory(category *entity.Category) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	category.ID = m.nextCategoryID
	m.categories[category.ID] = category
	m.nextCategoryID++
	return nil
}

// UpdateCategory Update modifies an existing category (for CategoryStorage interface)
func (m *MemoryStorage) UpdateCategory(category *entity.Category) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.categories[category.ID]; !exists {
		return entity.ErrNotFound
	}
	m.categories[category.ID] = category
	return nil
}

// DeleteCategory Delete removes a category by ID (for CategoryStorage interface)
func (m *MemoryStorage) DeleteCategory(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.categories[id]; !exists {
		return entity.ErrNotFound
	}
	delete(m.categories, id)
	return nil
}
