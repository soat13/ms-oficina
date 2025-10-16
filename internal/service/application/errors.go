package application

import "errors"

var (
	ErrServiceNotFound  = errors.New("service not found")
	ErrDuplicateService = errors.New("service with this name already exists")
)
