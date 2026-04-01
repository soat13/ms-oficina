package errors

import "errors"

var (
	ErrInvalidStatusTransaction = errors.New("invalid status transaction")
	ErrInvalidID                = errors.New("invalid ID")
	ErrInvalidJSON              = errors.New("invalid JSON body")
	ErrInvalidDocument          = errors.New("invalid document")
	ErrInvalidPhoneNumber       = errors.New("invalid phone number")
)
