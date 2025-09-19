package domain

import (
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type Service struct {
	ID        uuid.UUID
	Name      string
	Price     money.Money
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewService(id uuid.UUID, name string, price money.Money, currency string, now time.Time) (*Service, error) {
	s := &Service{
		ID:        idOrNew(id),
		Name:      strings.TrimSpace(name),
		Price:     price,
		Currency:  normalizeCurrency(currency),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.validate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Service) Rename(newName string, now time.Time) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return ErrInvalidServiceName
	}
	s.Name = newName
	s.UpdatedAt = now
	return nil
}

func (s *Service) ChangePrice(price money.Money, now time.Time) error {
	s.Price = price
	s.UpdatedAt = now
	return nil
}

func (s *Service) ChangeCurrency(newCurrency string, now time.Time) error {
	c := normalizeCurrency(newCurrency)
	if !isValidCurrency(c) {
		return ErrInvalidServiceCurrency
	}
	s.Currency = c
	s.UpdatedAt = now
	return nil
}

func (s *Service) validate() error {
	if s.Name == "" {
		return ErrInvalidServiceName
	}
	if s.Price.Cents < 1 {
		return ErrInvalidServicePrice
	}
	if !isValidCurrency(s.Currency) {
		return ErrInvalidServiceCurrency
	}
	return nil
}

func idOrNew(id uuid.UUID) uuid.UUID {
	if id == uuid.Nil {
		return uuid.New()
	}
	return id
}

func normalizeCurrency(c string) string {
	if c == "" {
		return "BRL"
	}
	return strings.ToUpper(strings.TrimSpace(c))
}

func isValidCurrency(c string) bool {
	if len(c) != 3 {
		return false
	}
	for _, r := range c {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
