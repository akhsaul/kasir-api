package repository

import "kasir-api/models"

// ProductRepository defines data access for products.
// Repository layer: data buat logic. Error database → cek sini.
type ProductRepository interface {
	GetAll() ([]*model.Product, error)
	GetByID(id int) (*model.Product, error)
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(id int) error
}
