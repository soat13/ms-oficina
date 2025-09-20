package domain

import (
	"errors"
)

var (
	ErrCustomerOrVehicleIDInvalid = errors.New("repair.order.customer.or.vehicle.id.invalid")
	ErrInvalidStatusTransition    = errors.New("repair.order.invalid.status.transition")
)
