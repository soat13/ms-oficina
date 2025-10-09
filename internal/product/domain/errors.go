package domain

import "errors"

var (
	ErrInvalidProductName  = errors.New("invalid product name")
	ErrInvalidProductPrice = errors.New("invalid product price")
	ErrInvalidProductStock = errors.New("invalid product stock")
)
