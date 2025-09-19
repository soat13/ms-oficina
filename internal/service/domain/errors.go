package domain

import "errors"

var (
	ErrInvalidServiceName     = errors.New("invalid service name")
	ErrInvalidServicePrice    = errors.New("invalid service price (must be > 0)")
	ErrInvalidServiceCurrency = errors.New("invalid service currency")
)
