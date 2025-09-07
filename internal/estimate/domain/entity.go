package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/pkg/entity"
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
		ID       uuid.UUID
		RepairID uuid.UUID
		Status   Status
		Items    map[uuid.UUID]Item
		entity.Timestamps
	}
)

const (
	StatusDraft            Status = "draft"
	StatusAwaitingApproval Status = "awaiting_approval"
	StatusApproved         Status = "approved"
	StatusRejected         Status = "rejected"
)

func NewEstimate(repairID uuid.UUID) (*Estimate, error) {
	if repairID == uuid.Nil {
		return nil, ErrRepairIDInvalid
	}

	now := time.Now()

	return &Estimate{
		ID:         uuid.New(),
		RepairID:   repairID,
		Status:     StatusDraft,
		Items:      make(map[uuid.UUID]Item),
		Timestamps: entity.NewTimestamps(now, now),
	}, nil
}
