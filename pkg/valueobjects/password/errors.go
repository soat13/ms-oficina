package password

import "errors"

var (
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong     = errors.New("password must be at most 72 characters")
	ErrInvalidPasswordHash = errors.New("invalid password hash")
)
