package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

func TestNewService(t *testing.T) {
	now := time.Now()
	validName := "Alignment"
	validPrice := money.Money{Cents: 12000}

	tests := []struct {
		name        string
		id          uuid.UUID
		serviceName string
		price       money.Money
		createdAt   time.Time
		updatedAt   time.Time
		wantErr     error
	}{
		{
			name:        "valid service",
			id:          uuid.Nil,
			serviceName: validName,
			price:       validPrice,
			createdAt:   time.Time{},
			updatedAt:   time.Time{},
			wantErr:     nil,
		},
		{
			name:        "valid service with specific ID",
			id:          uuid.New(),
			serviceName: validName,
			price:       validPrice,
			createdAt:   time.Time{},
			updatedAt:   time.Time{},
			wantErr:     nil,
		},
		{
			name:        "valid service with timestamps",
			id:          uuid.Nil,
			serviceName: validName,
			price:       validPrice,
			createdAt:   now,
			updatedAt:   now,
			wantErr:     nil,
		},
		{
			name:        "empty name",
			id:          uuid.Nil,
			serviceName: "",
			price:       validPrice,
			createdAt:   time.Time{},
			updatedAt:   time.Time{},
			wantErr:     ErrInvalidServiceName,
		},
		{
			name:        "name too short",
			id:          uuid.Nil,
			serviceName: "AB",
			price:       validPrice,
			createdAt:   time.Time{},
			updatedAt:   time.Time{},
			wantErr:     ErrInvalidServiceName,
		},
		{
			name:        "whitespace only name",
			id:          uuid.Nil,
			serviceName: "   ",
			price:       validPrice,
			createdAt:   time.Time{},
			updatedAt:   time.Time{},
			wantErr:     ErrInvalidServiceName,
		},
		{
			name:        "invalid price - zero",
			id:          uuid.Nil,
			serviceName: validName,
			price:       money.Money{Cents: 0},
			createdAt:   time.Time{},
			updatedAt:   time.Time{},
			wantErr:     ErrInvalidServicePrice,
		},
		{
			name:        "invalid price - negative",
			id:          uuid.Nil,
			serviceName: validName,
			price:       money.Money{Cents: -100},
			createdAt:   time.Time{},
			updatedAt:   time.Time{},
			wantErr:     ErrInvalidServicePrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewService(tt.id, tt.serviceName, tt.price, tt.createdAt, tt.updatedAt)
			if err != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if got.ID == uuid.Nil {
					t.Error("NewService() ID should not be Nil")
				}
				if got.Name != tt.serviceName {
					t.Errorf("NewService() Name = %v, want %v", got.Name, tt.serviceName)
				}
				if got.Price.Cents != tt.price.Cents {
					t.Errorf("NewService() Price = %v, want %v", got.Price.Cents, tt.price.Cents)
				}
				if !tt.createdAt.IsZero() && got.CreatedAt != tt.createdAt {
					t.Errorf("NewService() CreatedAt = %v, want %v", got.CreatedAt, tt.createdAt)
				}
				if !tt.updatedAt.IsZero() && got.UpdatedAt != tt.updatedAt {
					t.Errorf("NewService() UpdatedAt = %v, want %v", got.UpdatedAt, tt.updatedAt)
				}
			}
		})
	}
}

func TestServiceRename(t *testing.T) {
	service, _ := NewService(uuid.Nil, "Old Name", money.Money{Cents: 12000}, time.Time{}, time.Time{})
	originalUpdatedAt := service.UpdatedAt
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
			wantErr: ErrInvalidServiceName,
		},
		{
			name:    "name too short",
			newName: "AB",
			wantErr: ErrInvalidServiceName,
		},
		{
			name:    "whitespace only name",
			newName: "   ",
			wantErr: ErrInvalidServiceName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Rename(tt.newName)
			if err != tt.wantErr {
				t.Errorf("Rename() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if service.Name != tt.newName {
					t.Errorf("Rename() Name = %v, want %v", service.Name, tt.newName)
				}
				if !service.UpdatedAt.After(originalUpdatedAt) {
					t.Error("Rename() should update UpdatedAt")
				}
			}
		})
	}
}

func TestServiceChangePrice(t *testing.T) {
	service, _ := NewService(uuid.Nil, "Service", money.Money{Cents: 12000}, time.Time{}, time.Time{})
	originalUpdatedAt := service.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	newPrice := money.Money{Cents: 15000}
	err := service.ChangePrice(newPrice)

	if err != nil {
		t.Errorf("ChangePrice() error = %v, want nil", err)
	}
	if service.Price.Cents != newPrice.Cents {
		t.Errorf("ChangePrice() Price = %v, want %v", service.Price.Cents, newPrice.Cents)
	}
	if !service.UpdatedAt.After(originalUpdatedAt) {
		t.Error("ChangePrice() should update UpdatedAt")
	}
}
