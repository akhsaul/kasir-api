package service

import (
	"errors"
	"kasir-api/models"
	"kasir-api/repositories"
	"strings"
)

// ProductService handles business logic for products.
// Service layer: logic kode kita. Error logic → cek sini.
type ProductService struct {
	repo         repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

// NewProductService creates a new ProductService.
func NewProductService(repo repository.ProductRepository, categoryRepo repository.CategoryRepository) *ProductService {
	return &ProductService{repo: repo, categoryRepo: categoryRepo}
}

// GetAll retrieves all products.
func (s *ProductService) GetAll() ([]*model.Product, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a product by ID.
func (s *ProductService) GetByID(id int) (*model.Product, error) {
	if id <= 0 {
		return nil, model.ErrIDRequired
	}
	return s.repo.GetByID(id)
}

// Create creates a new product with validation.
func (s *ProductService) Create(product *model.Product) (*model.Product, error) {
	if err := s.validateProduct(product); err != nil {
		return nil, err
	}
	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	return product, nil
}

// Update updates an existing product with validation.
func (s *ProductService) Update(id int, product *model.Product) (*model.Product, error) {
	if id <= 0 {
		return nil, model.ErrIDRequired
	}
	if err := s.validateProduct(product); err != nil {
		return nil, err
	}
	product.ID = id
	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	return product, nil
}

// Delete removes a product by ID.
func (s *ProductService) Delete(id int) error {
	if id <= 0 {
		return model.ErrIDRequired
	}
	return s.repo.Delete(id)
}

func (s *ProductService) validateProduct(product *model.Product) error {
	if strings.TrimSpace(product.Name) == "" {
		return model.ErrNameRequired
	}
	if product.Price <= 0 {
		return model.ErrPriceInvalid
	}
	if product.Stock < 0 {
		return model.ErrStockInvalid
	}
	if product.CategoryID != nil {
		if _, err := s.categoryRepo.GetByID(*product.CategoryID); err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return model.ErrCategoryNotFound
			}
			return err
		}
	}
	return nil
}
