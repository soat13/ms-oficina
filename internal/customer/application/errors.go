package application

import "errors"

var (
	ErrCustomerNotFound  = errors.New("customer not found")
	ErrDuplicateCustomer = errors.New("customer with this document already exists")
)
