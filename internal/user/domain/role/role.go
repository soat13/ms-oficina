package role

import "strings"

type Roles []Role
type Role string

const (
	attendant = "attendant"
	manager   = "manager"
	mechanic  = "mechanic"
)

var validRoles = []Role{attendant, manager, mechanic}

func New(v string) (Role, error) {
	role := strings.ToLower(strings.TrimSpace(v))
	for _, valid := range validRoles {
		if role == string(valid) {
			return valid, nil
		}
	}

	return "", ErrInvalidRole
}

func NewRoles(values []string) (Roles, error) {
	if len(values) == 0 {
		return nil, ErrRolesRequired
	}

	roles := make(Roles, len(values))
	for i, v := range values {
		role, err := New(v)
		if err != nil {
			return nil, err
		}
		roles[i] = role
	}
	return roles, nil
}

func (rs Roles) Strings() []string {
	result := make([]string, len(rs))
	for i, r := range rs {
		result[i] = string(r)
	}
	return result
}
