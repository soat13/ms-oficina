package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/pkg/entity"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type Service struct {
	ID    uuid.UUID
	Name  string
	Price money.Money
	entity.Timestamps
}

func NewService(id uuid.UUID, name string, price money.Money, CratedAt, UpdatedAt time.Time) (*Service, error) {

	service := &Service{
		ID:    idOrNew(id),
		Name:  strings.TrimSpace(name),
		Price: price,
	}

	if err := service.validate(); err != nil {
		return nil, err
	}

	if !CratedAt.IsZero() {
		service.CreatedAt = CratedAt
	}

	if !UpdatedAt.IsZero() {
		service.UpdatedAt = UpdatedAt
	}

	return service, nil
}

func (s *Service) Rename(newName string, now time.Time) error {
	newName = strings.TrimSpace(newName)
	if len(newName) < 3 {
		return ErrInvalidServiceName
	}
	s.Name = newName
	s.UpdatedAt = now
	return nil
}

func (s *Service) ChangePrice(price money.Money) error {
	s.Price = price
	s.Touch()
	return nil
}

func (s *Service) validate() error {
	if len(s.Name) < 3 {
		return ErrInvalidServiceName
	}
	if s.Price.Cents <= 0 {
		return ErrInvalidServicePrice
	}
	return nil
}

func idOrNew(id uuid.UUID) uuid.UUID {
	if id == uuid.Nil {
		return uuid.New()
	}
	return id
}
