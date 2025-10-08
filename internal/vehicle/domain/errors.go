package domain

import "errors"

var (
	ErrInvalidVehiclePlate = errors.New("invalid vehicle plate")
	ErrInvalidVehicleBrand = errors.New("invalid vehicle brand")
	ErrInvalidVehicleModel = errors.New("invalid vehicle model")
	ErrInvalidVehicleYear  = errors.New("invalid vehicle year")
)
