package application

import "errors"

var (
	ErrEstimateNotFound         = errors.New("estimate not found")
	ErrProductNotAvailable      = errors.New("product not available")
	ErrProductOrServiceNotFound = errors.New("product or service not found")
)
