package domain

import (
	"errors"
)

var (
	ErrRepairIDInvalid                = errors.New("repair ID invalid")
	ErrQuantityInvalid                = errors.New("quantity must be greater than zero")
	ErrItemNotFound                   = errors.New("item not found in estimate")
	ErrCannotChangeItemsAfterApproval = errors.New("cannot change items after approval")
	ErrEstimateMustHaveAtLeastOneItem = errors.New("estimate must have at least one item")
)
