package estimate

import (
	"time"

	"github.com/google/uuid"
)

type Created struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	EstimateID    uuid.UUID `json:"estimate_id"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (Created) Topic() string { return "estimate.created" }
