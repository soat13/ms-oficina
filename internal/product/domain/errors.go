package domain

import "errors"

var (
	ErrInvalidProductName  = errors.New("invalid product name")
	ErrInvalidProductPrice = errors.New("invalid product price (must be > 0)")
	ErrInvalidProductStock = errors.New("invalid product stock (must be >= 0)")
)
