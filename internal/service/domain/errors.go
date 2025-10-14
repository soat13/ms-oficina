package domain

import "errors"

var (
	ErrInvalidServiceName  = errors.New("invalid service name")
	ErrInvalidServicePrice = errors.New("invalid service price")
)
