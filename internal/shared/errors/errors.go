package errors

import "errors"

var (
	ErrInvalidStatusTransaction = errors.New("shared.invalid.status.transaction")
	ErrInvalidID                = errors.New("shared.invalid.id")
	ErrInvalidJSON              = errors.New("invalid JSON body")
	ErrInvalidDocument          = errors.New("invalid document")
	ErrInvalidPhoneNumber       = errors.New("invalid phone number")
	ErrInvalidEmail             = errors.New("invalid email")
	ErrPasswordTooShort         = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong          = errors.New("password must be at most 72 characters")
	ErrInvalidPasswordHash      = errors.New("invalid password hash")
)
