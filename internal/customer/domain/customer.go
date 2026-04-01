package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/oficina-utils/pkg/entity"
	uuidPkg "github.com/soat13/oficina-utils/pkg/utils/uuid"
	"github.com/soat13/oficina-utils/pkg/valueobjects/document"
	"github.com/soat13/oficina-utils/pkg/valueobjects/email"
	"github.com/soat13/oficina-utils/pkg/valueobjects/phone"
)

type Customer struct {
	ID          uuid.UUID
	Name        string
	Document    document.Document
	Email       email.Email
	PhoneNumber phone.PhoneNumber
	entity.Timestamps
}

func NewCustomer(id uuid.UUID, name string, document document.Document, phoneNumber phone.PhoneNumber, email email.Email) (*Customer, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidCustomerName
	}

	now := time.Now()
	customer := &Customer{
		ID:          uuidPkg.IDOrNew(id),
		Name:        name,
		Document:    document,
		PhoneNumber: phoneNumber,
		Email:       email,
		Timestamps:  entity.NewTimestamps(now, now),
	}

	return customer, nil
}

func (customer *Customer) ChangeName(newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return ErrInvalidCustomerName
	}
	customer.Name = newName
	customer.Touch()
	return nil
}

func (customer *Customer) ChangePhoneNumber(phoneNumber phone.PhoneNumber) {
	customer.PhoneNumber = phoneNumber
	customer.Touch()
}

func (customer *Customer) ChangeEmail(newEmail email.Email) {
	customer.Email = newEmail
	customer.Touch()
}
