package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/user/domain/role"
	"github.com/soat13/fase-1-oficina/pkg/entity"
	uuidPkg "github.com/soat13/fase-1-oficina/pkg/utils/uuid"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/document"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/email"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/password"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/phone"
)

type User struct {
	ID          uuid.UUID
	Name        string
	Document    document.Document
	Email       email.Email
	PhoneNumber phone.PhoneNumber
	Password    password.Password
	Roles       role.Roles
	entity.Timestamps
}

func NewUser(id uuid.UUID, name string, document document.Document, phoneNumber phone.PhoneNumber, email email.Email, password password.Password, roles role.Roles) (*User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidUserName
	}

	now := time.Now()
	user := &User{
		ID:          uuidPkg.IDOrNew(id),
		Name:        name,
		Document:    document,
		PhoneNumber: phoneNumber,
		Email:       email,
		Password:    password,
		Roles:       roles,
		Timestamps:  entity.NewTimestamps(now, now),
	}

	return user, nil
}

func (user *User) ChangeName(newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return ErrInvalidUserName
	}
	user.Name = newName
	user.Touch()
	return nil
}

func (user *User) ChangePhoneNumber(phoneNumber phone.PhoneNumber) {
	user.PhoneNumber = phoneNumber
	user.Touch()
}

func (user *User) ChangeEmail(email email.Email) {
	user.Email = email
	user.Touch()
}

func (user *User) ChangePassword(password password.Password) {
	user.Password = password
	user.Touch()
}

func (user *User) ChangeRoles(roles role.Roles) error {
	user.Roles = roles
	user.Touch()
	return nil
}
