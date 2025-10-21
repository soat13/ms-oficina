package events

import (
	"time"

	"github.com/google/uuid"
)

type Canceled struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (Canceled) Topic() string { return "repairOrder.canceled" }

type DiagnosticsFinished struct {
	EventID       uuid.UUID         `json:"event_id"`
	OccurredAt    time.Time         `json:"occurred_at"`
	RepairOrderID uuid.UUID         `json:"repair_order_id"`
	Products      map[uuid.UUID]int `json:"products"`
	Services      map[uuid.UUID]int `json:"services"`
}

func (DiagnosticsFinished) Topic() string { return "repairOrder.diagnostics.finished" }
