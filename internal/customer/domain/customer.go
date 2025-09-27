package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/pkg/entity"
	uuidPkg "github.com/soat13/fase-1-oficina/pkg/utils/uuid"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/document"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/email"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/phone"
)

type Customer struct {
	ID          uuid.UUID
	Name        string
	Document    document.Document
	Email       email.Email
	PhoneNumber phone.PhoneNumber
	entity.Timestamps
}

func NewCustomer(id uuid.UUID, name, documentStr, phoneNumberStr, emailStr string, now time.Time) (*Customer, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidCustomerName
	}
	documentVO, docErr := document.New(documentStr)
	if docErr != nil {
		return nil, docErr
	}
	phoneNumberVO, phoneErr := phone.New(phoneNumberStr)
	if phoneErr != nil {
		return nil, phoneErr
	}
	emailVO, emailErr := email.New(emailStr)
	if emailErr != nil {
		return nil, emailErr
	}

	customer := &Customer{
		ID:          uuidPkg.IDOrNew(id),
		Name:        name,
		Document:    documentVO,
		PhoneNumber: phoneNumberVO,
		Email:       emailVO,
		Timestamps:  entity.NewTimestamps(now, now),
	}

	return customer, nil
}

func (customer *Customer) ChangeName(newName string, now time.Time) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return ErrInvalidCustomerName
	}
	customer.Name = newName
	customer.UpdatedAt = now
	return nil
}

func (customer *Customer) ChangePhoneNumber(newPhoneNumber string, now time.Time) error {
	phoneNumberVO, err := phone.New(newPhoneNumber)
	if err != nil {
		return err
	}
	customer.PhoneNumber = phoneNumberVO
	customer.UpdatedAt = now
	return nil
}

func (customer *Customer) ChangeEmail(newEmail string, now time.Time) error {
	emailVO, err := email.New(newEmail)
	if err != nil {
		return err
	}
	customer.Email = emailVO
	customer.UpdatedAt = now
	return nil
}
