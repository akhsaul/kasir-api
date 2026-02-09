package memory

import (
	"sync"

	model "kasir-api/models"
)

// CategoryRepository holds in-memory category storage and implements repository.CategoryRepository.
type CategoryRepository struct {
	mu             sync.RWMutex
	categories     map[int]*model.Category
	nextCategoryID int
}

// NewCategoryRepository creates a new in-memory category repository.
func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		categories:     make(map[int]*model.Category),
		nextCategoryID: 1,
	}
}

func (r *CategoryRepository) GetAll() ([]*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]*model.Category, 0, len(r.categories))
	for _, c := range r.categories {
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, exists := r.categories[id]
	if !exists {
		return nil, model.ErrCategoryNotFound
	}
	return c, nil
}

func (r *CategoryRepository) Create(category *model.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	category.ID = r.nextCategoryID
	r.categories[category.ID] = category
	r.nextCategoryID++
	return nil
}

func (r *CategoryRepository) Update(category *model.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.categories[category.ID]; !exists {
		return model.ErrCategoryNotFound
	}
	r.categories[category.ID] = category
	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.categories[id]; !exists {
		return model.ErrCategoryNotFound
	}
	delete(r.categories, id)
	return nil
}
