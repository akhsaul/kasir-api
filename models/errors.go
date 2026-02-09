package model

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound is the base "not found" error. Entity-specific errors wrap this,
	// so errors.Is(ErrProductNotFound, ErrNotFound) == true, etc.
	ErrNotFound            = errors.New("not found")
	ErrCategoryNotFound    = fmt.Errorf("category is not found: %w", ErrNotFound)
	ErrProductNotFound     = fmt.Errorf("product is not found: %w", ErrNotFound)
	ErrTransactionNotFound = fmt.Errorf("transaction is not found: %w", ErrNotFound)

	ErrNameRequired = errors.New("name should not be empty")
	ErrPriceInvalid = errors.New("price must be greater than 0")
	ErrStockInvalid = errors.New("stock must be greater than or equal to 0")
	ErrIDRequired   = errors.New("id is required")

	// Transaction errors.
	ErrEmptyCheckout     = errors.New("checkout items cannot be empty")
	ErrInvalidQuantity   = errors.New("quantity must be greater than 0")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidDateRange  = errors.New("invalid date range")
)
