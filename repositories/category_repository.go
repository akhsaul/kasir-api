package repository

import model "kasir-api/models"

// CategoryRepository defines data access for categories.
// Repository layer: data buat logic. Error database → cek sini.
type CategoryRepository interface {
	GetAll() ([]*model.Category, error)
	GetByID(id int) (*model.Category, error)
	Create(category *model.Category) error
	Update(category *model.Category) error
	Delete(id int) error
}
