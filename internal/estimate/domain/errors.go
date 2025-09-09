package domain

import (
	"errors"
)

var (
	ErrRepairIDInvalid = errors.New("repair ID invalid")
	ErrQuantityInvalid = errors.New("quantity must be greater than zero")
)
