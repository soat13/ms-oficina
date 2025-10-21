package errors

import "errors"

var (
	ErrInvalidStatusTransaction  = errors.New("invalid status transaction")
	ErrInvalidID                 = errors.New("invalid ID")
	ErrInvalidJSON               = errors.New("invalid JSON body")
	ErrInvalidDocument           = errors.New("invalid document")
	ErrInvalidPhoneNumber        = errors.New("invalid phone number")
	ErrInvalidEmail              = errors.New("invalid email")
	ErrPasswordTooShort          = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong           = errors.New("password must be at most 72 characters")
	ErrInvalidPasswordHash       = errors.New("invalid password hash")
	ErrNoProductsOrServicesFound = errors.New("at least one product or service is required")
)
