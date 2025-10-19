package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/document"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/email"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/phone"
)

func TestNewCustomer(t *testing.T) {
	validName := "John Doe"
	validDocument, _ := document.New("12345678900")
	validPhone, _ := phone.New("11987654321")
	validEmail, _ := email.New("test@example.com")

	tests := []struct {
		name         string
		id           uuid.UUID
		customerName string
		document     document.Document
		phoneNumber  phone.PhoneNumber
		email        email.Email
		wantErr      error
	}{
		{
			name:         "valid customer",
			id:           uuid.Nil,
			customerName: validName,
			document:     validDocument,
			phoneNumber:  validPhone,
			email:        validEmail,
			wantErr:      nil,
		},
		{
			name:         "valid customer with specific ID",
			id:           uuid.New(),
			customerName: validName,
			document:     validDocument,
			phoneNumber:  validPhone,
			email:        validEmail,
			wantErr:      nil,
		},
		{
			name:         "empty name",
			id:           uuid.Nil,
			customerName: "",
			document:     validDocument,
			phoneNumber:  validPhone,
			email:        validEmail,
			wantErr:      ErrInvalidCustomerName,
		},
		{
			name:         "whitespace only name",
			id:           uuid.Nil,
			customerName: "   ",
			document:     validDocument,
			phoneNumber:  validPhone,
			email:        validEmail,
			wantErr:      ErrInvalidCustomerName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCustomer(tt.id, tt.customerName, tt.document, tt.phoneNumber, tt.email)
			if err != tt.wantErr {
				t.Errorf("NewCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if got.ID == uuid.Nil {
					t.Error("NewCustomer() ID should not be Nil")
				}
				if got.Name != tt.customerName {
					t.Errorf("NewCustomer() Name = %v, want %v", got.Name, tt.customerName)
				}
				if got.Email != tt.email {
					t.Errorf("NewCustomer() Email = %v, want %v", got.Email, tt.email)
				}
				if got.Document != tt.document {
					t.Errorf("NewCustomer() Document = %v, want %v", got.Document, tt.document)
				}
				if got.PhoneNumber != tt.phoneNumber {
					t.Errorf("NewCustomer() PhoneNumber = %v, want %v", got.PhoneNumber, tt.phoneNumber)
				}
				if got.CreatedAt.IsZero() {
					t.Error("NewCustomer() CreatedAt should not be zero")
				}
				if got.UpdatedAt.IsZero() {
					t.Error("NewCustomer() UpdatedAt should not be zero")
				}
			}
		})
	}
}

func TestCustomerChangeName(t *testing.T) {
	doc, _ := document.New("12345678900")
	phone, _ := phone.New("11987654321")
	em, _ := email.New("test@example.com")

	customer, _ := NewCustomer(uuid.Nil, "Old Name", doc, phone, em)
	originalUpdatedAt := customer.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	tests := []struct {
		name    string
		newName string
		wantErr error
	}{
		{
			name:    "valid name change",
			newName: "New Name",
			wantErr: nil,
		},
		{
			name:    "empty name",
			newName: "",
			wantErr: ErrInvalidCustomerName,
		},
		{
			name:    "whitespace only name",
			newName: "   ",
			wantErr: ErrInvalidCustomerName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := customer.ChangeName(tt.newName)
			if err != tt.wantErr {
				t.Errorf("ChangeName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if customer.Name != tt.newName {
					t.Errorf("ChangeName() Name = %v, want %v", customer.Name, tt.newName)
				}
				if !customer.UpdatedAt.After(originalUpdatedAt) {
					t.Error("ChangeName() should update UpdatedAt")
				}
			}
		})
	}
}

func TestCustomerChangePhoneNumber(t *testing.T) {
	doc, _ := document.New("12345678900")
	oldPhone, _ := phone.New("11987654321")
	newPhone, _ := phone.New("11987654322")
	em, _ := email.New("test@example.com")

	customer, _ := NewCustomer(uuid.Nil, "Test Customer", doc, oldPhone, em)
	originalUpdatedAt := customer.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	customer.ChangePhoneNumber(newPhone)

	if customer.PhoneNumber != newPhone {
		t.Errorf("ChangePhoneNumber() PhoneNumber = %v, want %v", customer.PhoneNumber, newPhone)
	}
	if !customer.UpdatedAt.After(originalUpdatedAt) {
		t.Error("ChangePhoneNumber() should update UpdatedAt")
	}
}

func TestCustomerChangeEmail(t *testing.T) {
	doc, _ := document.New("12345678900")
	phone, _ := phone.New("11987654321")
	oldEmail, _ := email.New("old@example.com")
	newEmail, _ := email.New("new@example.com")

	customer, _ := NewCustomer(uuid.Nil, "Test Customer", doc, phone, oldEmail)
	originalUpdatedAt := customer.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	customer.ChangeEmail(newEmail)

	if customer.Email != newEmail {
		t.Errorf("ChangeEmail() Email = %v, want %v", customer.Email, newEmail)
	}
	if !customer.UpdatedAt.After(originalUpdatedAt) {
		t.Error("ChangeEmail() should update UpdatedAt")
	}
}
