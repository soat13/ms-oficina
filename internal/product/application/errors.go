package application

import "errors"

var (
	ErrDuplicateProduct = errors.New("product.duplicate")
	ErrProductNotFound  = errors.New("product.not_found")
)
