package errors

import "errors"

var (
	ErrInvalidStatusTransaction = errors.New("shared.invalid.status.transaction")
	ErrInvalidID                = errors.New("shared.invalid.id")
	ErrInvalidJSON              = errors.New("invalid JSON body")
	ErrInvalidBody              = errors.New("invalid request body")
)
