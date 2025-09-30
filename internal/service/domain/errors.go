package domain

import "errors"

var (
	ErrInvalidServiceName  = errors.New("service.invalid_name")
	ErrInvalidServicePrice = errors.New("service.invalid_price")
)
