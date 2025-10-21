package application

import (
	"errors"
)

var (
	ErrVehicleOrCustomerNotFound = errors.New("vehicle or customer not found")
	ErrProductNotFound           = errors.New("product not found")
	ErrServiceNotFound           = errors.New("service not found")
	ErrNoProductsOrServicesFound = errors.New("at least one product and service is required")
	ErrInsufficientStock         = errors.New("insufficient stock for one or more products")
	ErrInvalidQuantity           = errors.New("product or service quantity must be greater than zero")
)
