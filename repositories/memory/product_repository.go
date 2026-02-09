package memory

import (
	"strings"
	"sync"

	model "kasir-api/models"
	repository "kasir-api/repositories"
)

// ProductRepository holds in-memory product storage and implements repository.ProductRepository.
type ProductRepository struct {
	mu            sync.RWMutex
	products      map[int]*model.Product
	nextProductID int
	categoryRepo  repository.CategoryRepository
}

// NewProductRepository creates a new in-memory product repository with optional category lookup.
func NewProductRepository(categoryRepo repository.CategoryRepository) *ProductRepository {
	return &ProductRepository{
		products:      make(map[int]*model.Product),
		nextProductID: 1,
		categoryRepo:  categoryRepo,
	}
}

func (r *ProductRepository) enrichWithCategory(p *model.Product) {
	if r.categoryRepo != nil && p.CategoryID != nil {
		if cat, err := r.categoryRepo.GetByID(*p.CategoryID); err == nil {
			p.Category = &model.ProductCategory{
				Name:        cat.Name,
				Description: cat.Description,
			}
		}
	}
}

func (r *ProductRepository) GetAll(name string) ([]*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*model.Product, 0, len(r.products))
	for _, p := range r.products {
		// Filter by name if provided (case-insensitive partial match)
		if name != "" {
			loweredName := strings.ToLower(name)
			loweredProductName := strings.ToLower(p.Name)
			if !strings.Contains(loweredProductName, loweredName) {
				continue
			}
		}
		pCopy := *p
		r.enrichWithCategory(&pCopy)
		products = append(products, &pCopy)
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.products[id]
	if !exists {
		return nil, model.ErrProductNotFound
	}
	pCopy := *p
	r.enrichWithCategory(&pCopy)
	return &pCopy, nil
}

func (r *ProductRepository) Create(product *model.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	product.ID = r.nextProductID
	r.products[product.ID] = product
	r.nextProductID++
	r.enrichWithCategory(product)
	return nil
}

func (r *ProductRepository) Update(product *model.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return model.ErrProductNotFound
	}
	r.products[product.ID] = product
	r.enrichWithCategory(product)
	return nil
}

func (r *ProductRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[id]; !exists {
		return model.ErrProductNotFound
	}
	delete(r.products, id)
	return nil
}
