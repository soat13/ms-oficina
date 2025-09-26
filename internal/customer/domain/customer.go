package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	documentPkg "github.com/soat13/fase-1-oficina/pkg/document"
	"github.com/soat13/fase-1-oficina/pkg/entity"
	uuidPkg "github.com/soat13/fase-1-oficina/pkg/uuid"
)

type Customer struct {
	ID           uuid.UUID
	Name         string
	Document     string
	DocumentType string
	// TODO: ADD CELLPHONE VALUE OBJECT
	Cellphone string
	// TODO: ADD EMAIL VALUE OBJECT
	Email string
	entity.Timestamps
}

func NewCustomer(id uuid.UUID, name, document, cellphone, email string, now time.Time) (*Customer, error) {
	customer := &Customer{
		ID:           uuidPkg.IDOrNew(id),
		Name:         strings.TrimSpace(name),
		Document:     strings.TrimSpace(document),
		DocumentType: documentPkg.GetDocumentType(document),
		Cellphone:    strings.TrimSpace(cellphone),
		Email:        strings.TrimSpace(email),
	}

	err := customer.validate()

	if err != nil {
		return nil, err
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

func (customer *Customer) ChangeCellphone(newCellphone string, now time.Time) error {
	newCellphone = strings.TrimSpace(newCellphone)
	if newCellphone == "" {
		return ErrInvalidCustomerCellphone
	}
	customer.Cellphone = newCellphone
	customer.UpdatedAt = now
	return nil
}

func (customer *Customer) ChangeEmail(newEmail string, now time.Time) error {
	newEmail = strings.TrimSpace(newEmail)
	if newEmail == "" {
		return ErrInvalidCustomerEmail
	}
	customer.Email = newEmail
	customer.UpdatedAt = now
	return nil
}

func (customer *Customer) validate() error {
	if customer.Name == "" {
		return ErrInvalidCustomerName
	}
	if customer.Cellphone == "" {
		return ErrInvalidCustomerCellphone
	}
	if customer.Email == "" {
		return ErrInvalidCustomerEmail
	}
	if err := documentPkg.ValidateDocument(customer.Document); err != nil {
		return err
	}
	// TODO: VALIDATE EMAIL

	return nil
}
