package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/pkg/entity"
	"github.com/soat13/fase-1-oficina/pkg/money"
	uuidPkg "github.com/soat13/fase-1-oficina/pkg/utils/uuid"
)

type Product struct {
	ID    uuid.UUID
	Name  string
	Price money.Money
	Stock int
	entity.Timestamps
}

func NewProduct(id uuid.UUID, name string, price money.Money, stock int) (*Product, error) {
	name = strings.TrimSpace(name)
	if len(name) < 3 {
		return nil, ErrInvalidProductName
	}
	if price.Cents <= 0 {
		return nil, ErrInvalidProductPrice
	}
	if stock < 0 {
		return nil, ErrInvalidProductStock
	}

	now := time.Now()
	product := &Product{
		ID:         uuidPkg.IDOrNew(id),
		Name:       name,
		Price:      price,
		Stock:      stock,
		Timestamps: entity.NewTimestamps(now, now),
	}

	return product, nil
}

func (p *Product) Rename(newName string) error {
	newName = strings.TrimSpace(newName)
	if len(newName) < 3 {
		return ErrInvalidProductName
	}
	p.Name = newName
	p.Touch()
	return nil
}

func (p *Product) ChangePrice(price money.Money) error {
	p.Price = price
	p.Touch()
	return nil
}

func (p *Product) UpdateStock(stock int) error {
	if stock < 0 {
		return ErrInvalidProductStock
	}
	p.Stock = stock
	p.Touch()
	return nil
}
