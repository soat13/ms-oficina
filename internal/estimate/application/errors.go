package application

import "errors"

var (
	ErrInvalidRepairOrderStatus = errors.New("invalid repair order status")
	ErrEstimateNotFound         = errors.New("estimate not found")
	ErrProductNotAvailable      = errors.New("product not available")
)
