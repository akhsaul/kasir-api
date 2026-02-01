package service

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"strings"
)

// CategoryService handles business logic for categories.
// Service layer: logic kode kita. Error logic → cek sini.
type CategoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetAll retrieves all categories.
func (s *CategoryService) GetAll() ([]*model.Category, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a category by ID.
func (s *CategoryService) GetByID(id int) (*model.Category, error) {
	if id <= 0 {
		return nil, model.ErrIDRequired
	}
	return s.repo.GetByID(id)
}

// Create creates a new category with validation.
func (s *CategoryService) Create(category *model.Category) (*model.Category, error) {
	if err := s.validateCategory(category); err != nil {
		return nil, err
	}
	if err := s.repo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

// Update updates an existing category with validation.
func (s *CategoryService) Update(id int, category *model.Category) (*model.Category, error) {
	if id <= 0 {
		return nil, model.ErrIDRequired
	}
	if err := s.validateCategory(category); err != nil {
		return nil, err
	}
	category.ID = id
	if err := s.repo.Update(category); err != nil {
		return nil, err
	}
	return category, nil
}

// Delete removes a category by ID.
func (s *CategoryService) Delete(id int) error {
	if id <= 0 {
		return model.ErrIDRequired
	}
	return s.repo.Delete(id)
}

func (s *CategoryService) validateCategory(category *model.Category) error {
	if strings.TrimSpace(category.Name) == "" {
		return model.ErrNameRequired
	}
	return nil
}
