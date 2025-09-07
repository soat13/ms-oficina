package domain

import (
	"errors"
)

var (
	ErrOperationNotAllowed        = errors.New("operation not allowed in current status")
	ErrQuantityMustBeNonNegative  = errors.New("quantity must be greater than zero")
	ErrInvalidItemType            = errors.New("invalid item type")
	ErrCustomerOrVehicleIDInvalid = errors.New("customer or vehicle invalid")
)
