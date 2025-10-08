package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/pkg/entity"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type Product struct {
	ID    uuid.UUID
	Name  string
	Price money.Money
	Stock int
	entity.Timestamps
}

func NewProduct(id uuid.UUID, name string, price money.Money, stock int, createdAt, updatedAt time.Time) (*Product, error) {
	product := &Product{
		ID:    idOrNew(id),
		Name:  strings.TrimSpace(name),
		Price: price,
		Stock: stock,
	}

	if err := product.validate(); err != nil {
		return nil, err
	}

	if !createdAt.IsZero() {
		product.CreatedAt = createdAt
	}

	if !updatedAt.IsZero() {
		product.UpdatedAt = updatedAt
	}

	return product, nil
}

func (p *Product) Rename(newName string, now time.Time) error {
	newName = strings.TrimSpace(newName)
	if len(newName) < 3 {
		return ErrInvalidProductName
	}
	p.Name = newName
	p.UpdatedAt = now
	return nil
}

func (p *Product) ChangePrice(price money.Money) error {
	if price.Cents <= 0 {
		return ErrInvalidProductPrice
	}
	p.Price = price
	p.Touch()
	return nil
}

func (p *Product) UpdateStock(stock int, now time.Time) error {
	if stock < 0 {
		return ErrInvalidProductStock
	}
	p.Stock = stock
	p.UpdatedAt = now
	return nil
}

func (p *Product) validate() error {
	if len(p.Name) < 3 {
		return ErrInvalidProductName
	}
	if p.Price.Cents <= 0 {
		return ErrInvalidProductPrice
	}
	if p.Stock < 0 {
		return ErrInvalidProductStock
	}
	return nil
}

func idOrNew(id uuid.UUID) uuid.UUID {
	if id == uuid.Nil {
		return uuid.New()
	}
	return id
}
