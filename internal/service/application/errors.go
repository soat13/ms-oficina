package application

import "errors"

var (
	ErrDuplicateService = errors.New("service.duplicate")
	ErrServiceNotFound  = errors.New("service.not_found")
	ErrValidation       = errors.New("service.validation")
)
