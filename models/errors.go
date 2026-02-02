package model

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrCategoryNotFound = errors.New("Category is not found")
	ErrProductNotFound  = errors.New("Product is not found")
	ErrNameRequired     = errors.New("name should not be empty")
	ErrPriceInvalid = errors.New("price must be greater than 0")
    ErrStockInvalid = errors.New("stock must be greater than or equal to 0")
	ErrIDRequired   = errors.New("id is required")
)
