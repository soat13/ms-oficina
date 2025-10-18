package application

import "errors"

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrDuplicateProduct  = errors.New("product with this name already exists")
	ErrInsufficientStock = errors.New("insufficient stock for the product")
)
