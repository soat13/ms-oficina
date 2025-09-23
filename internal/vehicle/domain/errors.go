package domain

import "errors"

var (
	ErrInvalidPlate      = errors.New("invalid plate")
	ErrInvalidModel      = errors.New("invalid model")
	ErrInvalidBrand      = errors.New("invalid brand")
	ErrInvalidYear       = errors.New("invalid year")
	ErrInvalidCustomerId = errors.New("invalid customer")
)
