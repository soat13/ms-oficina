package application

import "errors"

var (
	ErrVehicleNotFound    = errors.New("vehicle not found")
	ErrDuplicatePlate     = errors.New("vehicle plate already exists")
	ErrCustomerNotFound   = errors.New("customer not found")
)
