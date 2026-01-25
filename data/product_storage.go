package data

import "kasir-api/entity"

// ProductStorage defines the interface for product storage operations
type ProductStorage interface {
	GetAll() ([]*entity.Product, error)
	GetByID(id int) (*entity.Product, error)
	Create(product *entity.Product) error
	Update(product *entity.Product) error
	Delete(id int) error
}
