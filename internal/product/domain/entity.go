package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type Product struct {
	ID        uuid.UUID
	Name      string
	Price     money.Money
	Stock     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProduct(id uuid.UUID, name string, price money.Money, stock int, now time.Time) (*Product, error) {
	p := &Product{
		ID:        idOrNew(id),
		Name:      strings.TrimSpace(name),
		Price:     price,
		Stock:     stock,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := p.validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Product) Rename(newName string, now time.Time) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return ErrInvalidProductName
	}
	p.Name = newName
	p.UpdatedAt = now
	return nil
}

func (p *Product) ChangePrice(price money.Money, now time.Time) error {
	if price.Cents < 1 {
		return ErrInvalidProductPrice
	}
	p.Price = price
	p.UpdatedAt = now
	return nil
}

func (p *Product) ChangeStock(stock int, now time.Time) error {
	if stock < 0 {
		return ErrInvalidProductStock
	}
	p.Stock = stock
	p.UpdatedAt = now
	return nil
}

func (p *Product) validate() error {
	if p.Name == "" {
		return ErrInvalidProductName
	}
	if p.Price.Cents < 1 {
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
