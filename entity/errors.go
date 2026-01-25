package entity

import "errors"

var (
	ErrNotFound     = errors.New("product not found")
	ErrNameRequired = errors.New("name is required")
	ErrPriceInvalid = errors.New("price must be greater than 0")
	ErrStockInvalid = errors.New("stock must be greater than or equal to 0")
	ErrIDRequired   = errors.New("id is required")
)
