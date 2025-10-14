package application

import "errors"

var (
	ErrDuplicateService = errors.New("service duplicated")
	ErrServiceNotFound  = errors.New("service not found")
	ErrValidation       = errors.New("service validation")
)
