package application

import "errors"

var (
	ErrCustomerNotFound  = errors.New("customer not found")
	ErrDuplicateCustomer = errors.New("customer with this document or email already exists")
	ErrDuplicateEmail    = errors.New("customer with this email already exists")
)
