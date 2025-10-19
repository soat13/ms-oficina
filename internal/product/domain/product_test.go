package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

func TestNewProduct(t *testing.T) {
	validName := "Product A"
	validPrice := money.Money{Cents: 10000}
	validStock := 50

	tests := []struct {
		name        string
		id          uuid.UUID
		productName string
		price       money.Money
		stock       int
		wantErr     error
	}{
		{
			name:        "valid product",
			id:          uuid.Nil,
			productName: validName,
			price:       validPrice,
			stock:       validStock,
			wantErr:     nil,
		},
		{
			name:        "valid product with specific ID",
			id:          uuid.New(),
			productName: validName,
			price:       validPrice,
			stock:       validStock,
			wantErr:     nil,
		},
		{
			name:        "valid product with zero stock",
			id:          uuid.Nil,
			productName: validName,
			price:       validPrice,
			stock:       0,
			wantErr:     nil,
		},
		{
			name:        "empty name",
			id:          uuid.Nil,
			productName: "",
			price:       validPrice,
			stock:       validStock,
			wantErr:     ErrInvalidProductName,
		},
		{
			name:        "name too short",
			id:          uuid.Nil,
			productName: "AB",
			price:       validPrice,
			stock:       validStock,
			wantErr:     ErrInvalidProductName,
		},
		{
			name:        "invalid price - zero",
			id:          uuid.Nil,
			productName: validName,
			price:       money.Money{Cents: 0},
			stock:       validStock,
			wantErr:     ErrInvalidProductPrice,
		},
		{
			name:        "invalid price - negative",
			id:          uuid.Nil,
			productName: validName,
			price:       money.Money{Cents: -100},
			stock:       validStock,
			wantErr:     ErrInvalidProductPrice,
		},
		{
			name:        "invalid stock - negative",
			id:          uuid.Nil,
			productName: validName,
			price:       validPrice,
			stock:       -1,
			wantErr:     ErrInvalidProductStock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProduct(tt.id, tt.productName, tt.price, tt.stock)
			if err != tt.wantErr {
				t.Errorf("NewProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if got.ID == uuid.Nil {
					t.Error("NewProduct() ID should not be Nil")
				}
				if got.Name != tt.productName {
					t.Errorf("NewProduct() Name = %v, want %v", got.Name, tt.productName)
				}
				if got.Price.Cents != tt.price.Cents {
					t.Errorf("NewProduct() Price = %v, want %v", got.Price.Cents, tt.price.Cents)
				}
				if got.Stock != tt.stock {
					t.Errorf("NewProduct() Stock = %v, want %v", got.Stock, tt.stock)
				}
				if got.CreatedAt.IsZero() {
					t.Error("NewProduct() CreatedAt should not be zero")
				}
				if got.UpdatedAt.IsZero() {
					t.Error("NewProduct() UpdatedAt should not be zero")
				}
			}
		})
	}
}

func TestProductRename(t *testing.T) {
	product, _ := NewProduct(uuid.Nil, "Old Name", money.Money{Cents: 10000}, 50)
	originalUpdatedAt := product.UpdatedAt
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
			wantErr: ErrInvalidProductName,
		},
		{
			name:    "name too short",
			newName: "AB",
			wantErr: ErrInvalidProductName,
		},
		{
			name:    "whitespace only name",
			newName: "   ",
			wantErr: ErrInvalidProductName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := product.Rename(tt.newName)
			if err != tt.wantErr {
				t.Errorf("Rename() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if product.Name != tt.newName {
					t.Errorf("Rename() Name = %v, want %v", product.Name, tt.newName)
				}
				if !product.UpdatedAt.After(originalUpdatedAt) {
					t.Error("Rename() should update UpdatedAt")
				}
			}
		})
	}
}

func TestProductChangePrice(t *testing.T) {
	product, _ := NewProduct(uuid.Nil, "Product", money.Money{Cents: 10000}, 50)
	originalUpdatedAt := product.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	newPrice := money.Money{Cents: 15000}
	err := product.ChangePrice(newPrice)

	if err != nil {
		t.Errorf("ChangePrice() error = %v, want nil", err)
	}
	if product.Price.Cents != newPrice.Cents {
		t.Errorf("ChangePrice() Price = %v, want %v", product.Price.Cents, newPrice.Cents)
	}
	if !product.UpdatedAt.After(originalUpdatedAt) {
		t.Error("ChangePrice() should update UpdatedAt")
	}
}

func TestProductUpdateStock(t *testing.T) {
	product, _ := NewProduct(uuid.Nil, "Product", money.Money{Cents: 10000}, 50)
	originalUpdatedAt := product.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	tests := []struct {
		name     string
		newStock int
		wantErr  error
	}{
		{
			name:     "valid stock update",
			newStock: 100,
			wantErr:  nil,
		},
		{
			name:     "zero stock",
			newStock: 0,
			wantErr:  nil,
		},
		{
			name:     "negative stock",
			newStock: -1,
			wantErr:  ErrInvalidProductStock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := product.UpdateStock(tt.newStock)
			if err != tt.wantErr {
				t.Errorf("UpdateStock() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if product.Stock != tt.newStock {
					t.Errorf("UpdateStock() Stock = %v, want %v", product.Stock, tt.newStock)
				}
				if !product.UpdatedAt.After(originalUpdatedAt) {
					t.Error("UpdateStock() should update UpdatedAt")
				}
			}
		})
	}
}
