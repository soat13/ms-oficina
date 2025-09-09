package domain

import (
	"time"

	"github.com/google/uuid"
	pkgEntity "github.com/soat13/fase-1-oficina/pkg/entity"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type (
	Status string

	ItemType string

	Item struct {
		ID       uuid.UUID
		Name     string
		Price    money.Money
		Quantity int
		Type     ItemType
	}

	Estimate struct {
		ID            uuid.UUID
		RepairOrderID uuid.UUID
		Status        Status
		items         map[uuid.UUID]Item
		pkgEntity.Timestamps
	}
)

const (
	ProductItemType ItemType = "product"
	ServiceItemType ItemType = "service"

	StatusDraft            Status = "draft"
	StatusAwaitingApproval Status = "awaiting_approval"
	StatusApproved         Status = "approved"
	StatusRejected         Status = "rejected"
)

func NewEstimate(repairID uuid.UUID, CratedAt, UpdatedAt time.Time) (*Estimate, error) {
	if repairID == uuid.Nil {
		return nil, ErrRepairIDInvalid
	}

	now := time.Now()

	entity := &Estimate{
		ID:            uuid.New(),
		RepairOrderID: repairID,
		Status:        StatusDraft,
		items:         make(map[uuid.UUID]Item),
		Timestamps:    pkgEntity.NewTimestamps(now, now),
	}

	if !CratedAt.IsZero() {
		entity.CreatedAt = CratedAt
	}

	if !UpdatedAt.IsZero() {
		entity.UpdatedAt = UpdatedAt
	}

	return entity, nil
}

func (e *Estimate) Items() map[uuid.UUID]Item {
	return e.items
}

func (e *Estimate) AddItem(id uuid.UUID, name string, price money.Money, quantity int, itemType ItemType) error {
	if quantity < 1 {
		return ErrQuantityInvalid
	}

	item := Item{
		ID:       id,
		Name:     name,
		Price:    price,
		Quantity: quantity,
		Type:     itemType,
	}

	e.items[item.ID] = item
	return nil
}

func (e *Estimate) Total() money.Money {
	total := money.Money{Cents: 0}
	for _, item := range e.items {
		itemTotal := money.Money{Cents: item.Price.Cents * int64(item.Quantity)}
		total = total.Add(itemTotal)
	}

	return total
}
