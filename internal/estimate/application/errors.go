package application

import "errors"

var (
	ErrInvalidRepairOrderStatus = errors.New("estimate.invalid.repair.order.status")
	ErrEstimateNotFound         = errors.New("estimate.not.found")
)
