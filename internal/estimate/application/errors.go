package application

import "errors"

var (
	ErrInvalidRepairOrderStatus = errors.New("invalid repair order status to perform this operation")
	ErrRepairOrderNotFound      = errors.New("repair order not found")
	ErrSomeCatalogItemsNotFound = errors.New("catalog item not found")
)
