package data

import "kasir-api/entity"

// CategoryStorage defines the interface for category storage operations
type CategoryStorage interface {
	GetAll() ([]*entity.Category, error)
	GetByID(id int) (*entity.Category, error)
	Create(category *entity.Category) error
	Update(category *entity.Category) error
	Delete(id int) error
}
