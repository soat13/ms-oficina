package domain

import "errors"

var (
	ErrInvalidProductName  = errors.New("product.invalid_name")
	ErrInvalidProductPrice = errors.New("product.invalid_price")
	ErrInvalidProductStock = errors.New("product.invalid_stock")
)
