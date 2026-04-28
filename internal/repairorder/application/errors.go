package application

import (
	"errors"
)

var (
	ErrVehicleOrCustomerNotFound = errors.New("vehicle or customer not found")
	ErrProductNotFound           = errors.New("product not found")
	ErrServiceNotFound           = errors.New("service not found")
	ErrInsufficientStock         = errors.New("insufficient stock for one or more products")
	ErrInvalidQuantity           = errors.New("product or service quantity must be greater than zero")
	ErrTotalEstimateRequired     = errors.New("total estimate is required to finish execution")
)
