package role

import "errors"

var (
	ErrInvalidRole   = errors.New("invalid role")
	ErrRolesRequired = errors.New("at least one role is required")
)
