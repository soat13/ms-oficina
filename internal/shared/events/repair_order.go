package events

import (
	"time"

	"github.com/google/uuid"
)

const TopicRepairOrderStockReductionConfirmed = "repairorder-product-stock-reduction-confirmed"

type RepairOrderCanceled struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (RepairOrderCanceled) Topic() string { return "repairorder-canceled" }

type RepairOrderDiagnosticsFinished struct {
	EventID       uuid.UUID         `json:"event_id"`
	OccurredAt    time.Time         `json:"occurred_at"`
	RepairOrderID uuid.UUID         `json:"repair_order_id"`
	Products      map[uuid.UUID]int `json:"products"`
	Services      map[uuid.UUID]int `json:"services"`
}

func (RepairOrderDiagnosticsFinished) Topic() string { return "repairorder-diagnostics-finished" }

type RepairOrderFinished struct {
	EventID       uuid.UUID `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	RepairOrderID uuid.UUID `json:"repair_order_id"`
}

func (RepairOrderFinished) Topic() string { return "repairorder-finished" }
