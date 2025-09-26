package application

import "errors"

var (
	ErrVehicleNotFound  = errors.New("vehicle not found")
	ErrDuplicateVehicle = errors.New("vehicle with this plate already exists")
)
