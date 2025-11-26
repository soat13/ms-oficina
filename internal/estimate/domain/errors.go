package domain

import (
	"errors"
)

var (
	ErrRepairIDInvalid           = errors.New("repair ID invalid")
	ErrQuantityInvalid           = errors.New("quantity must be greater than zero")
	ErrItemNotFound              = errors.New("item not found in estimate")
	ErrCannotRemoveAfterApproval = errors.New("cannot remove item after approval")
)
