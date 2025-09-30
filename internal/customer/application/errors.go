package application

import "errors"

var (
	ErrCustomerNotFound  = errors.New("customer not found")
	ErrDuplicateDocument = errors.New("customer with this document already exists")
	ErrDuplicateEmail    = errors.New("customer with this email already exists")
)
