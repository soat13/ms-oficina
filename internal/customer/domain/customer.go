package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	documentPkg "github.com/soat13/fase-1-oficina/pkg/document"
	uuidPkg "github.com/soat13/fase-1-oficina/pkg/uuid"
)

type Customer struct {
	ID           uuid.UUID
	Name         string
	Cellphone    string
	Document     string
	DocumentType string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewCustomer(id uuid.UUID, name string, cellphone string, document string, now time.Time) (*Customer, error) {
	customer := &Customer{
		ID:           uuidPkg.IDOrNew(id),
		Name:         strings.TrimSpace(name),
		Cellphone:    strings.TrimSpace(cellphone),
		Document:     document,
		DocumentType: documentPkg.GetDocumentType(document),
		CreatedAt:    now,
		UpdatedAt:    now,
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

func (customer *Customer) validate() error {
	if customer.Name == "" {
		return ErrInvalidCustomerName
	}
	if customer.Cellphone == "" {
		return ErrInvalidCustomerCellphone
	}
	if err := documentPkg.ValidateDocument(customer.Document); err != nil {
		return err
	}

	return nil
}
